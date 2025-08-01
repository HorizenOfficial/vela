package manager

import (
	"context"
	"fmt"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/storage"
	"log"
	"sync"
	"time"
)

// SecureProcessorManager is an implementation of the Manager interface
type SecureProcessorManager struct {
	config           *Config
	blockchainClient blockchain.Client
	executorClient   communication.ExecutorClient
	dataLayer        storage.ApplicationStateStore
	mu               sync.Mutex
	isRunning        bool
	stopChan         chan struct{}
	wg               sync.WaitGroup
}

// NewSecureProcessorManager creates a new SecureProcessorManager
func NewSecureProcessorManager(config *Config, blockchainClient blockchain.Client, dataLayer storage.ApplicationStateStore, executorClient communication.ExecutorClient) *SecureProcessorManager {
	manager := &SecureProcessorManager{
		config:           config,
		blockchainClient: blockchainClient,
		executorClient:   executorClient,
		dataLayer:        dataLayer,
		stopChan:         make(chan struct{}),
	}
	// Set up the executor client
	manager.executorClient.SetClientRequestHandler(manager)

	return manager
}

// Start starts the manager
func (m *SecureProcessorManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("manager is already running")
	}

	// Connect to the executor
	if err := m.executorClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to executor: %w", err)
	}

	// Start the blockchain polling loop in a goroutine
	m.wg.Add(1)
	go m.pollBlockchain(ctx)

	m.isRunning = true
	return nil
}

// Stop stops the manager
func (m *SecureProcessorManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	// Wait for the polling loop to stop
	m.wg.Wait()

	// Close the executor client
	if err := m.executorClient.Close(); err != nil {
		return fmt.Errorf("failed to close executor client: %w", err)
	}

	// Close the blockchain client
	if err := m.blockchainClient.Close(); err != nil {
		return fmt.Errorf("failed to close blockchain client: %w", err)
	}

	// Close the data layer
	if err := m.dataLayer.Close(); err != nil {
		return fmt.Errorf("failed to close data layer: %w", err)
	}

	m.isRunning = false
	return nil
}

// pollBlockchain polls the blockchain for new requests
func (m *SecureProcessorManager) pollBlockchain(ctx context.Context) {
	defer m.wg.Done()

	ticker := time.NewTicker(time.Duration(m.config.BlockchainPollingInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Get pending requests from the blockchain
			requests, err := m.blockchainClient.GetPendingRequests(ctx)
			if err != nil {
				log.Printf("Manager: Failed to get pending requests: %v", err)
				continue
			}

			// Process each request
			for _, req := range requests {
				if err := m.processRequest(ctx, req); err != nil {
					// Log the error and mark the request as failed
					log.Printf("Manager: Failed to process request %s: %v", req.RequestID, err)
					if err = m.blockchainClient.MarkRequestFailed(ctx, req.RequestID); err != nil {
						log.Printf("Manager: Failed to mark request %s as failed: %v", req.RequestID, err)
					}
				} else {
					log.Printf("Manager: Processed and marked as completed request %s", req.RequestID)
				}
			}
		}
	}
}

// processRequest processes a request
func (m *SecureProcessorManager) processRequest(ctx context.Context, req *common.Request) error {
	switch req.RequestType {
	case common.Deploy:
		return m.processDeployApp(ctx, req)
	case common.Process:
		return m.processProcessRequest(ctx, req)
	case common.Deanonymize:
		return m.processDeanonymization(ctx, req)
	default:
		return fmt.Errorf("unsupported request type: %s", req.RequestType)
	}
}

// processDeployApp processes a deploy app request
func (m *SecureProcessorManager) processDeployApp(ctx context.Context, req *common.Request) error {
	// Deploy the application
	updatePayload, appState, err := m.executorClient.SendDeployApp(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to deploy application: %w", err)
	}

	// Store the application state
	err = m.dataLayer.StoreApplicationState(ctx, appState)
	if err != nil {
		return fmt.Errorf("failed to persist application state: %w", err)
	}

	// Store the WASM bytecode
	err = m.dataLayer.StoreWASMBytecode(ctx, appState.ApplicationID, req.Payload)
	if err != nil {
		return fmt.Errorf("failed to persist WASM bytecode: %w", err)
	}

	// Submit the state update to the blockchain
	if err = m.blockchainClient.SubmitStateUpdate(ctx, updatePayload); err != nil {
		return fmt.Errorf("failed to submit state update: %w", err)
	}

	log.Printf("Manager: Deployed application %s", appState.ApplicationID)
	return nil
}

// processProcessRequest processes a process request
func (m *SecureProcessorManager) processProcessRequest(ctx context.Context, req *common.Request) error {
	// Get the application state
	appState, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed to get application state: %w", err)
	}

	// Get the WASM module for the application
	wasmBytes, err := m.dataLayer.GetWASMBytecode(ctx, req.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed to get wasm bytes: %w", err)
	}

	// Get the user key for the request sender
	senderKey, err := m.dataLayer.GetUserKey(ctx, req.Sender)
	if err != nil {
		return fmt.Errorf("failed to get user key for sender %s: %w", req.Sender, err)
	}

	// Process the request
	updatePayload, updatedState, err := m.executorClient.SendProcessRequest(ctx, req, appState, senderKey, wasmBytes)
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}

	// Store the updated application state
	err = m.dataLayer.StoreApplicationState(ctx, updatedState)
	if err != nil {
		return fmt.Errorf("failed to persist application state: %w", err)
	}

	// Submit the state update to the blockchain
	if err = m.blockchainClient.SubmitStateUpdate(ctx, updatePayload); err != nil {
		return fmt.Errorf("failed to submit state update: %w", err)
	}

	log.Printf("Manager: Processed request %s for application %s", req.RequestID, req.ApplicationID)
	return nil
}

// processDeanonymization processes a deanonymization request
func (m *SecureProcessorManager) processDeanonymization(ctx context.Context, req *common.Request) error {
	// Get the application state
	appState, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed to get application state: %w", err)
	}

	// Get the WASM module for the application
	wasmBytes, err := m.dataLayer.GetWASMBytecode(ctx, req.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed to get wasm bytes: %w", err)
	}

	// Get the user key for the request sender
	senderKey, err := m.dataLayer.GetUserKey(ctx, req.Sender)
	if err != nil {
		return fmt.Errorf("failed to get user key for sender %s: %w", req.Sender, err)
	}

	// Generate the deanonymization report
	report, err := m.executorClient.SendGenerateDeanonymizationReport(ctx, req, appState, senderKey, wasmBytes)
	if err != nil {
		return fmt.Errorf("failed to generate deanonymization report: %w", err)
	}

	// Store the deanonymization report
	err = m.dataLayer.StoreDeanonymizationReport(ctx, report)
	if err != nil {
		return fmt.Errorf("failed to persist deanonymization report: %w", err)
	}

	// Submit the deanonymization report to the blockchain
	err = m.blockchainClient.SubmitDeanonymizationReport(ctx, report)
	if err != nil {
		return fmt.Errorf("failed to submit deanonymization report: %w", err)
	}

	log.Printf("Manager: Generated deanonymization report %s for application %s", report.ReportID, req.ApplicationID)
	return nil
}

func (m *SecureProcessorManager) HandleGetUserKeys(ctx context.Context, users []string) (map[string][]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil, fmt.Errorf("manager is not running")
	}
	userKeys := make(map[string][]byte)
	// Get user keys from the blockchain or data layer
	for _, user := range users {
		userKey, err := m.dataLayer.GetUserKey(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to get user {%s} key: %w", user, err)
		}
		userKeys[user] = userKey
	}

	return userKeys, nil
}
