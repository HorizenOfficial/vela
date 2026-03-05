package deployartifact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStoreSaveWASM_Success(t *testing.T) {
	root := t.TempDir()
	store, err := NewStore(root)
	require.NoError(t, err)

	payload := []byte("test-wasm-binary")
	resp, err := store.SaveWASM(strings.NewReader(string(payload)))
	require.NoError(t, err)

	sum := sha256.Sum256(payload)
	expectedSHA := hex.EncodeToString(sum[:])
	require.Equal(t, "sha256:"+expectedSHA, resp.ArtifactID)
	require.Equal(t, expectedSHA, resp.WasmSHA256)
	require.EqualValues(t, len(payload), resp.WasmSize)

	blobPath := filepath.Join(root, "blobs", expectedSHA+".wasm")
	blob, err := os.ReadFile(blobPath)
	require.NoError(t, err)
	require.Equal(t, payload, blob)

	metaPath := filepath.Join(root, "meta_global", expectedSHA+".json")
	metaRaw, err := os.ReadFile(metaPath)
	require.NoError(t, err)

	var meta Metadata
	require.NoError(t, json.Unmarshal(metaRaw, &meta))
	require.Equal(t, resp.ArtifactID, meta.ArtifactID)
	require.Equal(t, resp.WasmSHA256, meta.WasmSHA256)
	require.Equal(t, resp.WasmSize, meta.WasmSize)
	require.False(t, meta.CreatedAt.IsZero())
	require.Equal(t, time.UTC, meta.CreatedAt.Location())
}

func TestStoreSaveWASM_IdempotentReupload(t *testing.T) {
	store, err := NewStore(t.TempDir())
	require.NoError(t, err)

	payload := "same-bytes"
	first, err := store.SaveWASM(strings.NewReader(payload))
	require.NoError(t, err)
	second, err := store.SaveWASM(strings.NewReader(payload))
	require.NoError(t, err)

	require.Equal(t, first, second)
}

func TestStoreSaveWASM_EmptyPayload(t *testing.T) {
	store, err := NewStore(t.TempDir())
	require.NoError(t, err)

	_, err = store.SaveWASM(strings.NewReader(""))
	require.ErrorIs(t, err, ErrEmptyWASM)
}

func TestNewStore_RejectsEmptyPath(t *testing.T) {
	_, err := NewStore(" ")
	require.Error(t, err)
}
