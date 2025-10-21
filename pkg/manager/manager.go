package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
)

// As of now we support only one app having this ID
const (
	admittedAppID = "1"
)

// SecureProcessorManager is an implementation of the Manager interface
type SecureProcessorManager struct {
	config           *Config
	blockchainClient blockchain.Client
	executorClient   communication.ExecutorClient
	dataLayer        storage.DataLayer
	mu               sync.RWMutex
	isRunning        bool
	stopChan         chan struct{}
	wg               sync.WaitGroup
	endReorgTime     time.Time
}

// NewSecureProcessorManager creates a new SecureProcessorManager
func NewSecureProcessorManager(config *Config, blockchainClient blockchain.Client, dataLayer storage.DataLayer, executorClient communication.ExecutorClient) *SecureProcessorManager {
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

	// Connect to the blockchain
	if err := m.blockchainClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to blockchain node: %w", err)
	}

	// Start the blockchain polling loop in a goroutine
	m.wg.Add(1)
	go m.pollBlockchain(ctx)

	log.Printf("Manager starting - Ethereum address: " + m.config.PrivateKey.PublicKey().Address())

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
			err := m.processRequestFromChain(ctx)
			if err != nil {
				log.Fatalf("Manager: Error processing requests from chain: %v, exiting", err)
			}
		}
	}
}

// processRequestFromChain retrieves the next pending request from the blockchain and processes it
func (m *SecureProcessorManager) processRequestFromChain(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.isRunning {
		log.Printf("Manager is not started yet, skipping")
		return nil
	}

	// Get next pending request from the blockchain
	request, stateRoot, err := m.blockchainClient.GetNextPendingRequest(ctx)
	if err != nil {
		log.Printf("Manager: Failed to get pending request: %v", err)
		return nil
	}

	localStateRoot, err := m.dataLayer.LastVersionID()
	if err != nil {
		if dbErr, ok := err.(*storageErrors.Error); ok && dbErr.Code == storageErrors.NoVersionInDb {
			localStateRoot = make([]byte, 32) // Initialize to zero state root if no version exists
		} else {
			log.Printf("Manager: Failed to get local state root: %v", err)
			return nil
		}
	}

	if !bytes.Equal(localStateRoot, stateRoot[:]) {
		log.Printf("Manager: State root mismatch, expected %x, got %x. Checking if it is a REORG.", localStateRoot, stateRoot)

		isReorg, err := m.checkIfReorg(stateRoot)
		if err != nil {
			log.Printf("Manager: Failed to check for reorg: %v", err)
			return nil
		}

		if isReorg {
			if m.endReorgTime.IsZero() {
				log.Printf("Manager: Starting REORG timeout %d", m.config.ReorgTimeout)
				m.endReorgTime = time.Now().Add(time.Duration(m.config.ReorgTimeout) * time.Second)
				return nil
			}
			if time.Now().Before(m.endReorgTime) {
				log.Printf("Manager: REORG timeout not expired yet. Keep waiting...")
				return nil
			}
			log.Printf("Manager: REORG not solved within timeout => Rollback the DB")
			if err := m.dataLayer.Rollback(stateRoot[:]); err != nil {
				log.Printf("Manager: Error while rolling back the DB: %v. ", err)
				return nil
			}

		} else {
			return fmt.Errorf("unrecoverable disalignment between DB and chain, no matching state root found in db")
		}

	}

	if !m.endReorgTime.IsZero() {
		log.Printf("Manager: State roots match, REORG resolved")
		m.endReorgTime = time.Time{}
	}

	if request == nil {
		return nil
	}

	log.Printf("Manager: processing request %x", request.RequestID)

	if err := m.processRequest(ctx, request); err != nil {
		// Log the error and mark the request as failed
		log.Printf("Manager: Failed to process request %s: %v", request.RequestID, err)
		if err = m.blockchainClient.MarkRequestFailed(ctx, request.RequestID); err != nil {
			log.Printf("Manager: Failed to mark request %s as failed: %v", request.RequestID, err)
		}
	} else {
		log.Printf("Manager: Processed request %s", request.RequestID)
	}
	return nil

}

func (m *SecureProcessorManager) checkIfReorg(stateRoot [32]byte) (bool, error) {
	if stateRoot == [32]byte{} {
		log.Printf("Manager: State root is zero, REORG")
		// Don't look for older db versions, just mark as reorged and wait for next poll
		return true, nil
	}

	oldVersions, err := m.dataLayer.ListVersions()
	if err != nil {
		log.Printf("Manager: Failed to get db old versions: %v", err)
		return false, err
	}

	for _, oldVersion := range oldVersions[1:] {
		if bytes.Equal(oldVersion, stateRoot[:]) {
			log.Printf("Manager: Found matching state root %x in db, REORG. Checking if the timeout is expired", stateRoot)
			return true, nil
		}
	}
	return false, nil

}

// processRequest processes a request
func (m *SecureProcessorManager) processRequest(ctx context.Context, req *common.Request) error {
	if !m.isRunning {
		log.Printf("Manager is not started yet, skipping")
		return fmt.Errorf("Manager is not started yet")
	}

	// check admitted appID. TODO: This check will be removed in future
	if req.ApplicationID != admittedAppID {
		return fmt.Errorf("application id %s is not admitted", req.ApplicationID)
	}

	switch req.RequestType {
	case common.Deploy:
		return m.processDeployApp(ctx, req)
	case common.Process, common.AssociateKey:
		return m.processProcessRequest(ctx, req)
	case common.Deanonymize:
		return m.processDeanonymization(ctx, req)
	default:
		return fmt.Errorf("unsupported request type: %s", req.RequestType)
	}
}

func (m *SecureProcessorManager) submitStateOnChain(ctx context.Context, updatePayload *common.UpdatePayload) error {
	// Submit the state update to the blockchain
	if err := m.blockchainClient.SubmitStateUpdate(ctx, updatePayload); err != nil {
		log.Printf("Failed to submit state update for error: %v", err)
		log.Printf("Rollback the application state to previous version")
		if err := m.dataLayer.Rollback(updatePayload.PrevStateRoot[:]); err != nil {
			// If this happens, the local db and the chain are out of sync and cannot be recovered automatically
			log.Fatalf("Failed to rollback application state: %v", err)
		}

		if _, ok := err.(blockchain.ReorgError); ok {
			log.Printf("REORG, do not call MarkFailed, wait for next poll")
			return nil
		}
		return fmt.Errorf("failed to submit state update: %w", err)
	}

	log.Printf("Manager: Processed request %s for application %s", updatePayload.RequestID, updatePayload.ApplicationID)
	return nil
}

// processDeployApp processes a deploy app request
func (m *SecureProcessorManager) processDeployApp(ctx context.Context, req *common.Request) error {
	log.Printf("Processing deploy app request: %s", req.RequestID)
	if !m.isRunning {
		log.Printf("Manager is not started yet, skipping")
		return fmt.Errorf("Manager is not started yet")
	}

	// check if app was already deployed
	_, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err == nil {
		return fmt.Errorf("application %s was already deployed", req.ApplicationID)
	}

	// Deploy the application
	updatePayload, appState, err := m.executorClient.SendDeployApp(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to deploy application: %w", err)
	}

	// Store the application state and WASM bytecode
	versionID := appState.StateRoot[:]
	err = m.dataLayer.Store(
		ctx,
		versionID[:],
		[]*common.ApplicationState{appState},
		[]*common.WASMData{{ApplicationID: appState.ApplicationID, Bytecode: req.Payload}},
	)
	if err != nil {
		log.Printf("failed to submit state update: %v", err)
		return nil
	}

	log.Printf("Deployed application, submit the state update to the blockchain")
	return m.submitStateOnChain(ctx, updatePayload)

}

// processProcessRequest processes a process request
func (m *SecureProcessorManager) processProcessRequest(ctx context.Context, req *common.Request) error {
	log.Printf("Processing Process app request: %s", req.RequestID)
	if !m.isRunning {
		log.Printf("Manager is not started yet, skipping")
		return fmt.Errorf("Manager is not started yet")
	}

	// Get the application state
	appState, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err != nil {
		log.Printf("GetApplicationState returns an error: %v", err)
		if strings.Contains(err.Error(), "application state not found") {
			// This can happen if the deploy transaction was not mined yet, mark request as failed
			return err
		}
		// Other errors are likely db errors, retry on next poll
		return nil
	}

	// Get the WASM module for the application
	wasmBytes, err := m.dataLayer.GetWASMBytecode(ctx, req.ApplicationID)
	if err != nil {
		log.Printf("GetWASMBytecode returns an error: %v", err)
		return nil
	}

	// Process the request
	updatePayload, updatedState, err := m.executorClient.SendProcessRequest(ctx, req, appState, wasmBytes)
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}

	// Store the updated application state
	versionID := updatedState.StateRoot[:]
	log.Printf("VersionID %x", versionID[:])

	err = m.dataLayer.Store(
		ctx,
		versionID[:],
		[]*common.ApplicationState{updatedState},
		nil,
	)
	if err != nil {
		log.Printf("failed to submit state update: %v", err)
		return nil
	}

	return m.submitStateOnChain(ctx, updatePayload)
}

// processDeanonymization processes a deanonymization request
func (m *SecureProcessorManager) processDeanonymization(ctx context.Context, req *common.Request) error {
	if !m.isRunning {
		log.Printf("Manager is not started yet, skipping")
		return fmt.Errorf("Manager is not started yet")
	}

	// Get the application state
	appState, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err != nil {
		log.Printf("GetApplicationState returns an error: %v", err)
		if strings.Contains(err.Error(), "application state not found") {
			// This can happen if the deploy transaction was not mined yet, mark request as failed
			return err
		}
		// Other errors are likely db errors, retry on next poll
		return nil
	}

	// Get the WASM module for the application
	wasmBytes, err := m.dataLayer.GetWASMBytecode(ctx, req.ApplicationID)
	if err != nil {
		log.Printf("GetWASMBytecode returns an error: %v", err)
		return nil
	}

	// Generate the deanonymization report
	report, err := m.executorClient.SendGenerateDeanonymizationReport(ctx, req, appState, wasmBytes)
	if err != nil {
		return fmt.Errorf("failed to generate deanonymization report: %w", err)
	}

	// If a path is configured, save the deanonymization report to the filesystem
	// Do not bail out from the method even if an error occurs, we must continue saving the report in the data layer
	if m.config.DeanonymizationReportPath != "" {
		// Ensure the directory exists
		if err := os.MkdirAll(m.config.DeanonymizationReportPath, 0755); err != nil {
			log.Printf("Failed to create directory for deanonymization reports: %v", err)
		} else {
			// Marshal the report to JSON
			reportJSON, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				log.Printf("Failed to marshal deanonymization report to JSON: %v", err)
			} else {
				// The request ID is unique across applications, but for cleaner organization we use a folder name that includes both the app ID and the request ID
				filePath := filepath.Join(m.config.DeanonymizationReportPath, req.ApplicationID+"_"+req.RequestID)
				// Write the report to the file
				if err := os.WriteFile(filePath, reportJSON, 0644); err != nil {
					log.Printf("Failed to write deanonymization report to file: %v", err)
				} else {
					log.Printf("Saved deanonymization report %s to %s", report.ReportID, filePath)
				}
			}
		}
	}

	// Store the deanonymization report
	err = m.dataLayer.StoreDeanonymizationReport(ctx, report)
	if err != nil {
		// TODO must we return an error?
		log.Printf("StoreDeanonymizationReport returns an error: %v", err)
		return nil
	}

	// Submit the deanonymization report to the blockchain
	err = m.blockchainClient.SubmitDeanonymizationReport(ctx, report)
	if err != nil {
		// TODO must we return an error?
		log.Printf("SubmitDeanonymizationReport returns an error: %v", err)
		return nil
	}

	log.Printf("Manager: Generated deanonymization report %s for application %s", report.ReportID, req.ApplicationID)
	return nil
}
