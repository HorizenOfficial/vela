package system

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/horizen-pes/pkg/testutil"
)

func TestHandshakeFailureSystem(t *testing.T) {
	log := logger.NewLogger("zerolog")
	// 1. Create a manager config with a mock data layer
	mgrConfig := manager.ReadConfig(log)
	mgrConfig.DataLayerType = "mockdb"

	// 2. Create a new system test suite without a keyset, it will try to get recovery data from datalayer, will not
	// find anything stored there, will create a new keyset but it will fail storing it in datalayer
	var keySet *executor.EnclaveKeySet = nil
	var recoveryData *common.EnclaveKeySetRecovery = nil
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "mock-runtime", mgrConfig, executor.DefaultConfig(), keySet, recoveryData, log)
	defer suite.Cleanup()

	// 3. Get the mock data layer and configure it to fail
	mockDataLayer, ok := suite.GetDataLayer().(*mockdb.MockDataLayer)
	require.True(t, ok, "failed to cast data layer to mockdb.MockDataLayer")

	mockDataLayer.AddMockedFunc("StoreEnclaveKeySetRecovery", func(ctx context.Context, recoveryData *common.EnclaveKeySetRecovery) error {
		return fmt.Errorf("Test error for StoreEnclaveKeySetRecovery")
	})

	// 4. Start executor and manager
	// Executor should start fine
	err := suite.StartExecutor()
	require.NoError(t, err)

	// Manager will fail to start because the handshake with the executor will fail in the background.
	// The executor will drop the connection.
	err = suite.StartManager()
	require.Error(t, err)
	require.Contains(t, err.Error(), "executor handshake failed")
}
