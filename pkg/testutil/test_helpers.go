package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela/pkg/blockchain"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/executor"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/manager"
	"github.com/stretchr/testify/require"
)

// SystemTestSuite is the mock-based system test suite. It uses an in-memory
// MockClient for blockchain interactions, suitable for fast WASM-stack tests.
type SystemTestSuite struct {
	*TestSuiteCore                          // shared infrastructure
	blockchainClient *blockchain.MockClient // mock-specific
}

// NewSystemTestSuite creates a self-contained system test suite.
// It accepts logger configs (not logger instances) so the suite can inject the
// ephemeral log-server port into RemoteLogParams before creating the loggers.
// This guarantees zeronetwork loggers connect to the correct address, and
// console/file loggers simply ignore the injected params.
func NewSystemTestSuite(t *testing.T, appType string, mgrLogCfg *logger.Config, excLogCfg *logger.Config) *SystemTestSuite {
	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)
	ctx := context.Background()
	keySet, newRecoveryData, err := executor.GenerateEnclaveKeySet(ctx, execConfig.KeySetRecoveryType, nil, nil, "")
	require.NoError(t, err)
	return NewSystemTestSuiteWithConfigs(t, appType, mgrConfig, execConfig, keySet, newRecoveryData, mgrLogCfg, excLogCfg)
}

func NewSystemTestSuiteWithConfigs(
	t *testing.T,
	appType string,
	mgrConfig *manager.Config,
	execConfig *executor.Config,
	keySet *executor.EnclaveKeySet,
	recoveryData *common.EnclaveKeySetRecovery,
	mgrLogCfg *logger.Config,
	excLogCfg *logger.Config,
) *SystemTestSuite {
	// Create mock blockchain client
	blockchainClient := blockchain.NewMockClient()

	// Build shared core infrastructure, passing the MockClient as the blockchain.Client
	core := NewTestSuiteCore(t, appType, mgrConfig, execConfig, blockchainClient, keySet, recoveryData, mgrLogCfg, excLogCfg)

	// Subscribe the event channel to MockClient events
	blockchainClient.SubscribeToEvents(core.ctx, core.eventChannel)

	return &SystemTestSuite{
		TestSuiteCore:    core,
		blockchainClient: blockchainClient,
	}
}

// --- Mock-specific methods (not shared with fullstack suite) ---

func (s *SystemTestSuite) SubmitRequest(req *common.Request) error {
	return s.blockchainClient.SendRequestToChain(s.ctx, req)
}

func (s *SystemTestSuite) WaitForAppStateInBlockchain(appID common.ApplicationIdType, timeout time.Duration) (*common.ApplicationState, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			state, err := s.blockchainClient.GetApplicationState(s.ctx, appID)
			if err == nil {
				return state, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %d in blockchain", appID)

		}
	}
}

func (s *SystemTestSuite) AssertRequestCompleted(requestID common.RequestIdType, timeout time.Duration) error {
	return s.blockchainClient.WaitForRequestCompletion(requestID, timeout)
}

// WaitForWithdrawal waits for a withdrawal to be processed
func (s *SystemTestSuite) WaitForWithdrawal(appID common.ApplicationIdType, timeout time.Duration) (*common.Withdrawal, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			// Check if withdrawal exists in blockchain
			withdrawals, err := s.blockchainClient.GetWithdrawals(s.ctx, appID)
			if err == nil && withdrawals != nil && len(*withdrawals) > 0 {
				return &(*withdrawals)[0], nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for withdrawal for app %d", appID)
		}
	}
}

func (s *SystemTestSuite) GetFailedRequest() []*common.Request {
	failedReq := s.blockchainClient.GetFailedRequests()
	return failedReq
}

func (s *SystemTestSuite) GetRequestUpdatePayload(reqId common.RequestIdType) (*common.UpdatePayload, error) {
	return s.blockchainClient.GetRequestUpdatePayload(s.ctx, reqId)
}

func (s *SystemTestSuite) Cleanup() error {
	s.CleanupCore()
	s.blockchainClient.ClearAllData()
	return nil
}
