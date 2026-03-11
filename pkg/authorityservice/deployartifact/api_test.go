package deployartifact

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestHandleUpload_Success(t *testing.T) {
	store, err := NewStore(t.TempDir())
	require.NoError(t, err)

	api := NewAPI(store, 0, testLogger())
	body, contentType := buildMultipartBody(t, "wasm", "app.wasm", []byte("wasm-content"))

	req := httptest.NewRequest(http.MethodPost, "/deploy/upload", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	api.HandleUpload(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp UploadResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

	sum := sha256.Sum256([]byte("wasm-content"))
	expectedSHA := hex.EncodeToString(sum[:])
	require.Equal(t, "sha256:"+expectedSHA, resp.ArtifactID)
	require.Equal(t, expectedSHA, resp.WasmSHA256)
}

func TestHandleUpload_MissingWASM(t *testing.T) {
	store, err := NewStore(t.TempDir())
	require.NoError(t, err)

	api := NewAPI(store, 0, testLogger())

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/deploy/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	api.HandleUpload(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleUpload_SizeExceeded(t *testing.T) {
	store, err := NewStore(t.TempDir())
	require.NoError(t, err)

	// 1MB max upload body.
	api := NewAPI(store, 1, testLogger())
	content := bytes.Repeat([]byte{0xAB}, 2*1024*1024)
	body, contentType := buildMultipartBody(t, "wasm", "app.wasm", content)

	req := httptest.NewRequest(http.MethodPost, "/deploy/upload", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	api.HandleUpload(rr, req)
	require.Equal(t, http.StatusRequestEntityTooLarge, rr.Code)
}

func TestHandleUpload_DuplicateUploadSameArtifactID(t *testing.T) {
	store, err := NewStore(t.TempDir())
	require.NoError(t, err)

	api := NewAPI(store, 0, testLogger())
	content := []byte("same-content")

	firstResp := uploadAndDecode(t, api, content)
	secondResp := uploadAndDecode(t, api, content)

	require.Equal(t, firstResp.ArtifactID, secondResp.ArtifactID)
	require.Equal(t, firstResp.WasmSHA256, secondResp.WasmSHA256)
}

func TestHandleUpload_BadMethod(t *testing.T) {
	store, err := NewStore(t.TempDir())
	require.NoError(t, err)

	api := NewAPI(store, 0, testLogger())
	req := httptest.NewRequest(http.MethodGet, "/deploy/upload", nil)
	rr := httptest.NewRecorder()

	api.HandleUpload(rr, req)
	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func uploadAndDecode(t *testing.T, api *API, content []byte) UploadResponse {
	t.Helper()
	body, contentType := buildMultipartBody(t, "wasm", "app.wasm", content)
	req := httptest.NewRequest(http.MethodPost, "/deploy/upload", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()
	api.HandleUpload(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var resp UploadResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	return resp
}

func buildMultipartBody(t *testing.T, fieldName, fileName string, content []byte) (io.Reader, string) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileWriter, err := writer.CreateFormFile(fieldName, fileName)
	require.NoError(t, err)
	_, err = fileWriter.Write(content)
	require.NoError(t, err)

	require.NoError(t, writer.Close())
	return &body, writer.FormDataContentType()
}

func testLogger() logger.Logger {
	return logger.NewLogger(&logger.Config{
		Kind:         "zerolog",
		Console:      true,
		ConsoleLevel: "error",
		ConsoleColor: false,
	})
}
