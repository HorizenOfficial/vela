package main

import (
	"context"
	"fmt"

	"testing"
	"time"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/executor"
	"github.com/horizen-pes/pkg/manager"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/horizen-pes/pkg/wasm"
	"github.com/stretchr/testify/require"
)

type SystemTestSuite struct {
	t                *testing.T
	manager          manager.Manager
	executor         executor.Executor
	blockchainClient *blockchain.MockClient
	dataLayer        *mockdb.MockDataLayer
	ctx              context.Context
	cancel           context.CancelFunc
	eventChannel     chan interface{}
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
	exec, err := executor.NewStatelessExecutor(execConfig, runtime, server)
	require.NoError(t, err)

	// Create event channel
	eventChannel := make(chan interface{}, 100)
	blockchainClient.SubscribeToEvents(ctx, eventChannel)

	return &SystemTestSuite{
		t:                t,
		manager:          mgr,
		executor:         exec,
		blockchainClient: blockchainClient,
		dataLayer:        dataLayer,
		ctx:              ctx,
		cancel:           cancel,
		eventChannel:     eventChannel,
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
