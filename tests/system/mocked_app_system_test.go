package system

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
	testutil "github.com/horizen-pes/pkg/testutil"
)

func TestDeployApp(t *testing.T) {
	suite := testutil.NewSystemTestSuite(t, "mock-runtime", logger.NewLogger("zerolog"))
	defer suite.Cleanup()

	// 1. Start executor
	err := suite.StartExecutor()
	require.NoError(t, err)

	// 2. Start manager
	err = suite.StartManager()
	require.NoError(t, err)

	RequestID := "2333"
	ApplicationId := "1"

	// 4. Submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: ApplicationId,
		RequestID:     RequestID,
		Payload:       []byte("deploy-payload"),
		Sender:        "test-user",
		Timestamp:     time.Now().Unix(),
	}
	err = suite.SubmitRequest(deployReq)
	require.NoError(t, err)

	// 5. Assert app state created in DB
	appState, err := suite.WaitForAppStateInDB(ApplicationId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 6. Assert app state created in blockchain
	appState, err = suite.WaitForAppStateInBlockchain(ApplicationId, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 7. Assert request marked as done
	err = suite.AssertRequestCompleted(RequestID, 10*time.Second)
	require.NoError(t, err)
}

func TestMockRuntimeAppFullSystemFlow(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running test in CI environment")
	}

	suite := testutil.NewSystemTestSuite(t, "mock-runtime", logger.NewLogger("zeroLog"))
	defer suite.Cleanup()
	// Load wasm bytecode for the wasm app
	wasmBytecode := []byte("mock-runtime-app-bytecode")

	testutil.ExecTestAppFullSystemFlow(t, suite, wasmBytecode)
}
