package system

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	executorCommKey    *cryptotypes.PrivateKeyP521      // Executor's communication key for testing
	executorSigningKey *cryptotypes.PrivateKeySecp256k1 // Executor's signing key for testing
}

func NewSystemTestSuite(t *testing.T, appType string) *SystemTestSuite {
	ctx, cancel := context.WithCancel(context.Background())

	// Create mock components
	blockchainClient := blockchain.NewMockClient()
	dataLayer := mockdb.NewMockDataLayer()

	// Create an executor client (TCP for testing)
	factory := communication.NewTCPConnectionFactory("localhost:8080")
	executorClient := communication.NewClient(factory)

	// Create manager
	config := manager.ReadConfig()
	mgr := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer, executorClient)

	// Create executor
	execConfig := executor.DefaultConfig() // just to generate keys

	server := communication.NewServer(factory)
	var runtime executor.Runtime
	switch appType {
	case "wasmtime-payment":
		runtime = wasm.NewWasmtimeRuntime()
	case "mock-runtime":
		runtime = executor.NewMockRuntime()
	default:
		t.Fatalf("Unknown app type: %s", appType)
	}
	exec := executor.NewStatelessExecutor(execConfig, runtime, server)

	// Create event channel
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
		executorCommKey:    execConfig.CommunicationKey, // Store the executor's communication key
		executorSigningKey: execConfig.SignatureKey,     // Store the executor's signing key
	}
}

func (s *SystemTestSuite) StartManager() error {
	go func() {
		if err := s.manager.Start(s.ctx); err != nil {
			s.t.Errorf("Manager failed: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SystemTestSuite) StartExecutor() error {
	go func() {
		if err := s.executor.Start(s.ctx); err != nil {
			s.t.Errorf("Executor failed: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SystemTestSuite) AddUserKeys(userID string, publicKey []byte) error {
	// Register in a blockchain client
	err := s.blockchainClient.RegisterPublicKey(s.ctx, userID, publicKey)
	if err == nil {
		// Register in data layer
		return s.dataLayer.StoreUserKey(s.ctx, userID, publicKey)
	}
	return err
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
			state, err := s.dataLayer.GetApplicationState(s.ctx, appID)
			if err == nil {
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
			state, err := s.blockchainClient.GetApplicationState(s.ctx, appID)
			if err == nil {
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

// WaitForEvent waits for a specific event to be published for a user
func (s *SystemTestSuite) WaitForEvent(userID string, timeout time.Duration) (*common.Event, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case event := <-s.eventChannel:
			log.Printf("TESTING: Received event: %+v", event)
			if evt, ok := event.(common.Event); ok && evt.UserID == userID {
				return &evt, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for event for user %s", userID)
		}
	}
}

// WaitForDeanonymizationReport waits for a deanonymization report to be generated
func (s *SystemTestSuite) WaitForDeanonymizationReport(reportID string, timeout time.Duration) (*common.DeanonymizationReport, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			// Check if deanonymization report exists in blockchain
			report, err := s.blockchainClient.GetDeanonymizationReport(s.ctx, reportID)
			if err == nil && report != nil {
				return report, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for deanonymization report %s", reportID)
		}
	}
}

// WaitForWithdrawal waits for a withdrawal to be processed
func (s *SystemTestSuite) WaitForWithdrawal(appID string, timeout time.Duration) (*common.Withdrawal, error) {
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
			return nil, fmt.Errorf("timeout waiting for withdrawal for app %s", appID)
		}
	}
}

func (s *SystemTestSuite) GetRequestUpdatePayload(reqId string) (*common.UpdatePayload, error) {
	// Get the update payload for the request
	return s.blockchainClient.GetRequestUpdatePayload(s.ctx, reqId)
}

// GetExecutorCommunicationKey returns the executor's communication public key for encryption
func (s *SystemTestSuite) GetExecutorCommunicationKey() (*cryptotypes.PublicKeyP521, error) {
	if s.executorCommKey == nil {
		return nil, fmt.Errorf("executor communication key not initialized")
	}
	return s.executorCommKey.PublicKey(), nil
}

// GetExecutorSigningKey returns the executor's signing public key for encryption
func (s *SystemTestSuite) GetExecutorSigningKey() (*cryptotypes.PublicKeySecp256k1, error) {
	if s.executorSigningKey == nil {
		return nil, fmt.Errorf("executor signing key not initialized")
	}
	return s.executorSigningKey.PublicKey(), nil
}

func (s *SystemTestSuite) LoadWasmModule(t *testing.T, moduleFilename string) []byte {
	// Load the compiled WASM module
	wasmPath := filepath.Join("wasm", moduleFilename)
	wasmBytes, err := os.ReadFile(wasmPath)
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
