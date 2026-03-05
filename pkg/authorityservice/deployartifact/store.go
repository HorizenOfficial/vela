package deployartifact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	blobsFolder = "blobs"
	metaFolder  = "meta_global"
)

type Store struct {
	rootPath  string
	blobsPath string
	metaPath  string
}

func NewStore(rootPath string) (*Store, error) {
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return nil, fmt.Errorf("deploy artifacts path is required")
	}

	s := &Store{
		rootPath:  rootPath,
		blobsPath: filepath.Join(rootPath, blobsFolder),
		metaPath:  filepath.Join(rootPath, metaFolder),
	}

	if err := s.ensureDirectories(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) SaveWASM(reader io.Reader) (*UploadResponse, error) {
	if reader == nil {
		return nil, fmt.Errorf("wasm reader is nil")
	}

	if err := s.ensureDirectories(); err != nil {
		return nil, err
	}

	tempFile, err := os.CreateTemp(s.blobsPath, "upload-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp blob file: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	hasher := sha256.New()
	size, err := copyAndHash(tempFile, hasher, reader)
	if err != nil {
		_ = tempFile.Close()
		return nil, fmt.Errorf("failed to read upload body: %w", err)
	}
	if size == 0 {
		_ = tempFile.Close()
		return nil, ErrEmptyWASM
	}

	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return nil, fmt.Errorf("failed to sync temp blob file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp blob file: %w", err)
	}

	shaHex := hex.EncodeToString(hasher.Sum(nil))
	artifactID := artifactIDPrefix + shaHex
	blobPath := filepath.Join(s.blobsPath, shaHex+".wasm")

	if err := renameIfMissing(tempPath, blobPath); err != nil {
		return nil, fmt.Errorf("failed to persist wasm blob: %w", err)
	}

	if err := syncDir(s.blobsPath); err != nil {
		return nil, fmt.Errorf("failed to sync blobs directory: %w", err)
	}

	metadata := Metadata{
		ArtifactID: artifactID,
		WasmSHA256: shaHex,
		WasmSize:   uint64(size),
		CreatedAt:  time.Now().UTC().Truncate(time.Second),
	}
	metaPath := filepath.Join(s.metaPath, shaHex+".json")
	if err := writeJSONAtomicallyIfMissing(metaPath, metadata); err != nil {
		return nil, fmt.Errorf("failed to persist artifact metadata: %w", err)
	}
	if err := syncDir(s.metaPath); err != nil {
		return nil, fmt.Errorf("failed to sync metadata directory: %w", err)
	}

	return &UploadResponse{
		ArtifactID: artifactID,
		WasmSHA256: shaHex,
		WasmSize:   uint64(size),
	}, nil
}

func (s *Store) ensureDirectories() error {
	if err := os.MkdirAll(s.blobsPath, 0o755); err != nil {
		return fmt.Errorf("failed to create blobs directory: %w", err)
	}
	if err := os.MkdirAll(s.metaPath, 0o755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}
	return nil
}

func copyAndHash(dst *os.File, hasher hash.Hash, src io.Reader) (int64, error) {
	return io.Copy(io.MultiWriter(dst, hasher), src)
}

func renameIfMissing(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		return nil
	}

	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	if _, err := os.Stat(dst); err == nil {
		// Another concurrent writer won the race; this is still a success.
		return nil
	}
	return fmt.Errorf("rename %s -> %s failed", src, dst)
}

func writeJSONAtomicallyIfMissing(dstPath string, value any) error {
	if _, err := os.Stat(dstPath); err == nil {
		return nil
	}

	dir := filepath.Dir(dstPath)
	tempFile, err := os.CreateTemp(dir, "meta-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp metadata file: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	encoded, err := json.Marshal(value)
	if err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to encode metadata json: %w", err)
	}

	if _, err := tempFile.Write(encoded); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to write metadata file: %w", err)
	}
	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to sync metadata file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close metadata file: %w", err)
	}

	return renameIfMissing(tempPath, dstPath)
}

func syncDir(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}
