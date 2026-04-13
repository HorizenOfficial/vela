package system

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/authorityservice/deployartifact"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/testutil"
	"github.com/stretchr/testify/require"
)

// uploadArtifactAndBuildDescriptorPayloadWithParams uploads a WASM artifact and builds
// a deploy descriptor with optional constructor params.
func uploadArtifactAndBuildDescriptorPayloadWithParams(t *testing.T, suite *testutil.SystemTestSuite, wasmBytecode []byte, constructorParams json.RawMessage) []byte {
	t.Helper()

	store, err := deployartifact.NewStore(suite.GetArtifactsPath())
	require.NoError(t, err)
	uploadAPI := deployartifact.NewAPI(store, 50, logger.NewLogger(&logger.Config{Kind: "zerolog", Console: false}))

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("wasm", "app.wasm")
	require.NoError(t, err)
	_, err = fileWriter.Write(wasmBytecode)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/deploy/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	uploadAPI.HandleUpload(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())

	var uploadResp deployartifact.UploadResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &uploadResp))

	localSum := sha256.Sum256(wasmBytecode)
	localSHA := hex.EncodeToString(localSum[:])
	localArtifactID, err := common.BuildArtifactID(localSHA)
	require.NoError(t, err)
	require.Equal(t, localSHA, uploadResp.WasmSHA256)
	require.Equal(t, localArtifactID, uploadResp.ArtifactID)

	descriptor := common.DeployDescriptor{
		Mode:              common.DeployModeArtifactRef,
		ArtifactID:        uploadResp.ArtifactID,
		WasmSHA256:        uploadResp.WasmSHA256,
		ConstructorParams: constructorParams,
	}
	payload, err := json.Marshal(descriptor)
	require.NoError(t, err)
	return payload
}

// uploadArtifactAndBuildDescriptorPayload uploads a WASM artifact and builds
// a deploy descriptor without constructor params.
func uploadArtifactAndBuildDescriptorPayload(t *testing.T, suite *testutil.SystemTestSuite, wasmBytecode []byte) []byte {
	return uploadArtifactAndBuildDescriptorPayloadWithParams(t, suite, wasmBytecode, nil)
}
