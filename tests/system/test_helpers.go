package system

import (
	"context"
	"fmt"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage"
	"log"
	"testing"
	"time"
)

type SystemTestSuite struct {
	t                *testing.T
	manager          manager.Manager
	executor         executor.Executor
	blockchainClient *blockchain.MockClient
	dataLayer        *storage.MockDataLayer
	eventChannel     chan interface{}
	ctx              context.Context
	cancel           context.CancelFunc
	executorCommKey  *common.PrivateKeyP521 // Executor's communication key for testing
}

func NewSystemTestSuite(t *testing.T) *SystemTestSuite {
	ctx, cancel := context.WithCancel(context.Background())

	// Create mock components
	blockchainClient := blockchain.NewMockClient()
	dataLayer := storage.NewMockDataLayer()

	// Create an executor client (TCP for testing)
	factory := communication.NewTCPConnectionFactory("localhost:8080")
	executorClient := communication.NewClient(factory)

	// Create manager
	config := manager.DefaultConfig()
	config.ExecutorConnectionType = "tcp"
	config.ExecutorConnectionParams = map[string]string{"url": "http://localhost:8080"}
	mgr := manager.NewSecureProcessorManager(config, blockchainClient, dataLayer, executorClient)

	// Create executor
	execConfig := executor.DefaultConfig()
	execConfig.ServerType = "tcp"
	execConfig.ServerAddr = "localhost:8080"

	server := communication.NewServer(factory)
	runtime := executor.NewMockRuntime()
	exec := executor.NewStatelessExecutor(execConfig, runtime, server)

	// Create event channel
	eventChannel := make(chan interface{}, 100)
	blockchainClient.SubscribeToEvents(ctx, eventChannel)

	return &SystemTestSuite{
		t:                t,
		manager:          mgr,
		executor:         exec,
		blockchainClient: blockchainClient,
		dataLayer:        dataLayer,
		eventChannel:     eventChannel,
		ctx:              ctx,
		cancel:           cancel,
		executorCommKey:  execConfig.CommunicationKey, // Store the executor's communication key
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

func (s *SystemTestSuite) AssertEventPublished(expectedEvent interface{}, timeout time.Duration) error {
	timeoutCh := time.After(timeout)

	for {
		select {
		case event := <-s.eventChannel:
			if event != nil {
				return nil
			}
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for event")
		}
	}
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

// GetExecutorCommunicationKey returns the executor's communication public key for encryption
func (s *SystemTestSuite) GetExecutorCommunicationKey() (*common.PublicKeyP521, error) {
	if s.executorCommKey == nil {
		return nil, fmt.Errorf("executor communication key not initialized")
	}
	return s.executorCommKey.PublicKey(), nil
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
