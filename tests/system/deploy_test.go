package system

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/blockchain"
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
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	err = suite.AddUserKeys(suite.simNode.Submitter, userKey.PublicKey().Bytes())
	require.NoError(t, err)

	// 4. Submit deploy request
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: "1",
		Payload:       []byte("deploy-payload"),
		Sender:        suite.simNode.Submitter.From.Hex(),
		Timestamp:     time.Now().Unix(),
	}
	requestID, err := suite.SubmitRequest(deployReq)
	require.NoError(t, err)

	// 5. Assert request marked as done
	result, err := suite.AssertRequestCompleted(requestID, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, blockchain.RequestCompleted, result, "Request should be marked as completed")

	// 6. Assert app state created in DB
	appState, err := suite.WaitForAppStateInDB(deployReq.ApplicationID, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appState)

	// 7. Assert app state created in blockchain
	appStateinBC, err := suite.WaitForAppStateInBlockchain(deployReq.ApplicationID, 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, appStateinBC)
	require.Equal(t, appState.StateRoot, appStateinBC.StateRoot, "App state should be the same in data layer and in blockchain")

}
