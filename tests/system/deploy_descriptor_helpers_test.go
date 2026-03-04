package system

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func stageArtifactAndBuildDescriptorPayload(t *testing.T, suite *testutil.SystemTestSuite, appID common.ApplicationIdType, wasmBytecode []byte) []byte {
	t.Helper()

	sum := sha256.Sum256(wasmBytecode)
	shaHex := hex.EncodeToString(sum[:])

	artifactID, err := common.BuildArtifactID(shaHex)
	require.NoError(t, err)

	blobPath := filepath.Join(suite.GetArtifactsPath(), "blobs", shaHex+".wasm")
	require.NoError(t, os.MkdirAll(filepath.Dir(blobPath), 0o755))
	require.NoError(t, os.WriteFile(blobPath, wasmBytecode, 0o600))

	descriptor := common.DeployDescriptor{
		Mode:          common.DeployModeArtifactRef,
		ApplicationID: appID,
		ArtifactID:    artifactID,
		WasmSHA256:    shaHex,
		WasmSize:      uint64(len(wasmBytecode)),
	}
	payload, err := json.Marshal(descriptor)
	require.NoError(t, err)
	return payload
}
