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

	"github.com/horizen-pes/pkg/authorityservice/deployartifact"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func uploadArtifactAndBuildDescriptorPayload(t *testing.T, suite *testutil.SystemTestSuite, appID common.ApplicationIdType, wasmBytecode []byte) []byte {
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
		Mode:          common.DeployModeArtifactRef,
		ApplicationID: appID,
		ArtifactID:    uploadResp.ArtifactID,
		WasmSHA256:    uploadResp.WasmSHA256,
	}
	payload, err := json.Marshal(descriptor)
	require.NoError(t, err)
	return payload
}
