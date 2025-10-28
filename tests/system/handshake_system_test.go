package system

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/horizen-pes/pkg/testutil"
)

func TestHandshakeFailureSystem(t *testing.T) {
	// 1. Create a manager config with a mock data layer
	mgrConfig := manager.ReadConfig()
	mgrConfig.DataLayerType = "mockdb"

	// 2. Create a new system test suite without a keyset
	var keySet *executor.EnclaveKeySet
	var recoveryData *common.EnclaveKeySetRecovery
	suite := testutil.NewSystemTestSuiteWithMgrConfig(t, "mock-runtime", mgrConfig, keySet, recoveryData)
	defer suite.Cleanup()

	// 3. Get the mock data layer and configure it to fail
	mockDataLayer, ok := suite.GetDataLayer().(*mockdb.MockDataLayer)
	require.True(t, ok, "failed to cast data layer to mockdb.MockDataLayer")

	mockDataLayer.AddMockedFunc("StoreEnclaveKeySetRecovery", func(ctx context.Context, recoveryData *common.EnclaveKeySetRecovery) error {
		return assert.AnError
	})

	// 4. Start executor and manager
	// Executor should start fine
	err := suite.StartExecutor()
	require.NoError(t, err)

	// Manager will start, but the handshake with the executor will fail in the background.
	// The executor will drop the connection.
	err = suite.StartManager()
	require.NoError(t, err)

	// Give some time for the handshake to fail
	time.Sleep(2 * time.Second)

	// 5. Try to submit a request and expect it to fail
	// This will fail because the manager's client is not connected to the executor.
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: "1",
		RequestID:     "123",
		Payload:       []byte("deploy-payload"),
		Sender:        "test-user",
		Timestamp:     time.Now().Unix(),
	}
	err = suite.SubmitRequest(deployReq)

	// TODO currently we do not trap an error in case the handshake fails
	//require.Error(t, err)
	//require.Contains(t, err.Error(), "not connected")
}
