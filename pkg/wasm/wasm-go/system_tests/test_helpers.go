package main_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/horizen-pes/pkg/wasm"
	"github.com/stretchr/testify/require"
)

type SystemTestSuite struct {
	t                  *testing.T
	manager            manager.Manager
	executor           executor.Executor
	blockchainClient   *blockchain.MockClient
	dataLayer          *mockdb.MockDataLayer
	eventChannel       chan interface{}
	ctx                context.Context
	cancel             context.CancelFunc
	executorCommKey    *cryptotypes.PrivateKeyP521
	executorSigningKey *cryptotypes.PrivateKeySecp256k1
}

func NewSystemTestSuite(t *testing.T, appType string) *SystemTestSuite {
	ctx, cancel := context.WithCancel(context.Background())

	blockchainClient := blockchain.NewMockClient()
	dataLayer := mockdb.NewMockDataLayer()

	factory := communication.NewTCPConnectionFactory("localhost:8080")
	executorClient := communication.NewClient(factory)

	config := manager.DefaultConfig()
	config.ExecutorConnectionType = "tcp"
	config.ExecutorConnectionParams = map[string]string{"url": "http://localhost:8080"}
	mgr := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer, executorClient)

	execConfig := executor.DefaultConfig()
	execConfig.ServerType = "tcp"
	execConfig.ServerAddr = "localhost:8080"

	server := communication.NewServer(factory)
	var rt executor.Runtime
	switch appType {
	case "wasmtime-payment":
		rt = wasm.NewWasmtimeRuntime()
	case "mock-runtime":
		rt = executor.NewMockRuntime()
	default:
		t.Fatalf("Unknown app type: %s", appType)
	}
	exec := executor.NewStatelessExecutor(execConfig, rt, server)

	eventChannel := make(chan interface{}, 100)
	blockchainClient.SubscribeToEvents(ctx, eventChannel)

	return &SystemTestSuite{
		t:                  t,
		manager:            mgr,
		executor:           exec,
		blockchainClient:   blockchainClient,
		dataLayer:          dataLayer,
		eventChannel:       eventChannel,
		ctx:                ctx,
		cancel:             cancel,
		executorCommKey:    execConfig.CommunicationKey,
		executorSigningKey: execConfig.SignatureKey,
	}
}

func (s *SystemTestSuite) StartManager() error {
	go func() { _ = s.manager.Start(s.ctx) }()
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SystemTestSuite) StartExecutor() error {
	go func() { _ = s.executor.Start(s.ctx) }()
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SystemTestSuite) AddUserKeys(userID string, publicKey []byte) error {
	if err := s.blockchainClient.RegisterPublicKey(s.ctx, userID, publicKey); err == nil {
		return s.dataLayer.StoreUserKey(s.ctx, userID, publicKey)
	}
	return fmt.Errorf("failed to add user key")
}

func (s *SystemTestSuite) SubmitRequest(req *common.Request) error {
	return s.blockchainClient.SubmitRequest(s.ctx, req)
}

func (s *SystemTestSuite) WaitForAppStateInDB(appID string, timeout time.Duration) (*common.ApplicationState, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-ticker.C:
			if state, err := s.dataLayer.GetApplicationState(s.ctx, appID); err == nil {
				return state, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %s", appID)
		}
	}
}

func (s *SystemTestSuite) WaitForAppStateInBlockchain(appID string, timeout time.Duration) (*common.ApplicationState, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-ticker.C:
			if state, err := s.blockchainClient.GetApplicationState(s.ctx, appID); err == nil {
				return state, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %s in blockchain", appID)
		}
	}
}

func (s *SystemTestSuite) AssertRequestCompleted(requestID string, timeout time.Duration) error {
	return s.blockchainClient.WaitForRequestCompletion(requestID, timeout)
}

func (s *SystemTestSuite) WaitForEvent(userID string, timeout time.Duration) (*common.Event, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeoutCh := time.After(timeout)
	for {
		select {
		case event := <-s.eventChannel:
			if evt, ok := event.(common.Event); ok && evt.UserID == userID {
				return &evt, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for event for user %s", userID)
		}
	}
}

func (s *SystemTestSuite) WaitForDeanonymizationReport(reportID string, timeout time.Duration) (*common.DeanonymizationReport, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-ticker.C:
			if report, err := s.blockchainClient.GetDeanonymizationReport(s.ctx, reportID); err == nil && report != nil {
				return report, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for deanonymization report %s", reportID)
		}
	}
}

func (s *SystemTestSuite) WaitForWithdrawal(appID string, timeout time.Duration) (*common.Withdrawal, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeoutCh := time.After(timeout)
	for {
		select {
		case <-ticker.C:
			if withdrawals, err := s.blockchainClient.GetWithdrawals(s.ctx, appID); err == nil && withdrawals != nil && len(*withdrawals) > 0 {
				return &(*withdrawals)[0], nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for withdrawal for app %s", appID)
		}
	}
}

func (s *SystemTestSuite) GetRequestUpdatePayload(reqId string) (*common.UpdatePayload, error) {
	return s.blockchainClient.GetRequestUpdatePayload(s.ctx, reqId)
}

func (s *SystemTestSuite) GetExecutorCommunicationKey() (*cryptotypes.PublicKeyP521, error) {
	if s.executorCommKey == nil {
		return nil, fmt.Errorf("executor communication key not initialized")
	}
	return s.executorCommKey.PublicKey(), nil
}

func (s *SystemTestSuite) GetExecutorSigningKey() (*cryptotypes.PublicKeySecp256k1, error) {
	if s.executorSigningKey == nil {
		return nil, fmt.Errorf("executor signing key not initialized")
	}
	return s.executorSigningKey.PublicKey(), nil
}

func (s *SystemTestSuite) LoadWasmModule(t *testing.T, moduleFilename string) []byte {
	wasmBytes, err := os.ReadFile("payment_app.wasm")
	require.NoError(t, err, "Failed to read WASM file")
	return wasmBytes
}

func (s *SystemTestSuite) Cleanup() error {
	s.cancel()
	if s.manager != nil {
		s.manager.Stop()
	}
	if s.executor != nil {
		s.executor.Close()
	}
	s.blockchainClient.ClearAllData()
	return nil
}
