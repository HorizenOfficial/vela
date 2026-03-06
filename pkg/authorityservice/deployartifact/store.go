package deployartifact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	blobsFolder = "blobs"
)

type Store struct {
	rootPath  string
	blobsPath string
}

func NewStore(rootPath string) (*Store, error) {
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return nil, fmt.Errorf("deploy artifacts path is required")
	}

	s := &Store{
		rootPath:  rootPath,
		blobsPath: filepath.Join(rootPath, blobsFolder),
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

func syncDir(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}
