package system

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/crypto"
)

func TestDeployApp(t *testing.T) {
	suite := NewSystemTestSuite(t, "mock-runtime")
	defer suite.Cleanup()

	// 1. Start executor
	err := suite.StartExecutor()
	require.NoError(t, err)

	// 2. Start manager
	err = suite.StartManager()
	require.NoError(t, err)

	// 3. Add user keys to registry
	userKey, _ := crypto.GeneratePrivateKeyP521()
	err = suite.AddUserKeys("test-user", userKey.PublicKey().Bytes())
	require.NoError(t, err)

	// 4. Submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: "test-app",
		RequestID:     "test-deploy-1",
		Payload:       []byte("deploy-payload"),
		Sender:        "test-user",
		Timestamp:     time.Now().Unix(),
	}
	err = suite.SubmitRequest(deployReq)
	require.NoError(t, err)

	// 5. Assert app state created in DB
	appState, err := suite.WaitForAppStateInDB("test-app", 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 6. Assert app state created in blockchain
	appState, err = suite.WaitForAppStateInBlockchain("test-app", 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 7. Assert request marked as done
	err = suite.AssertRequestCompleted("test-deploy-1", 10*time.Second)
	require.NoError(t, err)
}
