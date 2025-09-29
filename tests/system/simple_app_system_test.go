package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/crypto"
	"github.com/horizen-pes/pkg/testutil"
)

// buildAndLoadSimpleApp is a helper function to build the wasm module and read its bytecode.
func buildAndLoadSimpleApp(t *testing.T) []byte {
	// Get the project root directory to construct absolute paths
	_, b, _, ok := runtime.Caller(0)
	require.True(t, ok)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	appDir := filepath.Join(projectRoot, "app", "simple")

	// Build the wasm module
	cmd := exec.Command("make", "build")
	cmd.Dir = appDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "failed to build wasm module: %s", string(output))

	// Load wasm bytecode for the wasm app
	wasmPath := filepath.Join(appDir, "build", "simple_app.wasm")
	wasmBytecode, err := os.ReadFile(wasmPath)
	require.NoError(t, err)
	require.NotEmpty(t, wasmBytecode)

	return wasmBytecode
}

func TestDeploySimpleApp(t *testing.T) {
	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	// 1. Build and load wasm bytecode
	wasmBytecode := buildAndLoadSimpleApp(t)

	// 2. Start executor
	err := suite.StartExecutor()
	require.NoError(t, err)

	// 3. Start manager
	err = suite.StartManager()
	require.NoError(t, err)

	// 4. Add user keys to registry
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	err = suite.AddUserKeys("test-user", userKey.PublicKey().Bytes())
	require.NoError(t, err)

	RequestID := "233"
	ApplicationId := "1"

	// 5. Submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: ApplicationId,
		RequestID:     RequestID,
		Payload:       wasmBytecode,
		Sender:        "test-user",
		Timestamp:     time.Now().Unix(),
	}
	err = suite.SubmitRequest(deployReq)
	require.NoError(t, err)

	// 6. Assert app state created in DB
	appState, err := suite.WaitForAppStateInDB(ApplicationId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 7. Assert app state created in blockchain
	appState, err = suite.WaitForAppStateInBlockchain(ApplicationId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 8. Assert request marked as done
	err = suite.AssertRequestCompleted(RequestID, 10*time.Second)
	require.NoError(t, err)
}

func TestWasmtimeRuntimeSimpleAppFullSystemFlow(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	suite := testutil.NewSystemTestSuite(t, "wasm-runtime")
	defer suite.Cleanup()

	wasmBytecode := buildAndLoadSimpleApp(t)

	testutil.ExecTestAppFullSystemFlow(t, suite, wasmBytecode)
}
