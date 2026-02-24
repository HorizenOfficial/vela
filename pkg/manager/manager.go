package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/admin"
	"github.com/horizen-pes/pkg/blockchain"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
)

// Version is the current version of the manager
const Version = "0.1.0"

// As of now we support only one app having this ID
var (
	admittedAppID = common.NewApplicationId(1)
)

type ExecutorHandShake struct {
	isComplete chan struct{}
	once       sync.Once
	err        error
}

// SecureProcessorManager is an implementation of the Manager interface
type SecureProcessorManager struct {
	config            *Config
	blockchainClient  blockchain.Client
	executorClient    communication.ExecutorClient
	dataLayer         storage.DataLayer
	adminServer       admin.AdminCommandServer
	mu                sync.RWMutex
	isRunning         bool
	executorHandShake *ExecutorHandShake
	stopChan          chan struct{} // Channel to signal the polling loop to stop
	wg                sync.WaitGroup
	endReorgTime      time.Time
	log               logger.Logger
	fatalErrChan      chan error // Channel to signal fatal errors to main
}

// NewSecureProcessorManager creates a new SecureProcessorManager
func NewSecureProcessorManager(config *Config, blockchainClient blockchain.Client, dataLayer storage.DataLayer, executorClient communication.ExecutorClient, adminServer admin.AdminCommandServer, log logger.Logger) *SecureProcessorManager {
	manager := &SecureProcessorManager{
		config:           config,
		blockchainClient: blockchainClient,
		executorClient:   executorClient,
		dataLayer:        dataLayer,
		adminServer:      adminServer,
		stopChan:         make(chan struct{}),
		executorHandShake: &ExecutorHandShake{
			isComplete: make(chan struct{}),
		},
		log:          log,
		fatalErrChan: make(chan error, 1), // buffered to avoid blocking
	}
	// Set up the executor client
	manager.executorClient.SetClientRequestHandler(manager)

	// Set up the admin server command handler
	if adminServer != nil {
		adminServer.SetCmdHandler(manager)
	}

	return manager
}

// FatalErrChan returns a channel that receives fatal errors requiring shutdown.
func (m *SecureProcessorManager) FatalErrChan() <-chan error {
	return m.fatalErrChan
}

func (m *SecureProcessorManager) waitForExecutorHandshake() error {
	m.log.Info("Manager: Waiting for the executor to complete the handshake...")

	select {
	case <-m.executorHandShake.isComplete:
		// Handshake completed (successfully or with error)
	case <-time.After(time.Duration(m.config.HandshakeTimeout) * time.Second):
		// Handshake timed out
		m.executorHandShake.once.Do(func() {
			m.log.Warn("Manager: Handshake timed out after %d seconds", m.config.HandshakeTimeout)
			m.executorHandShake.err = fmt.Errorf("handshake timed out after %d seconds", m.config.HandshakeTimeout)
			close(m.executorHandShake.isComplete)
		})
	}

	if m.executorHandShake.err != nil {
		m.log.Error("Manager: ... executor completed handshake with error: %v", m.executorHandShake.err)
		return m.executorHandShake.err
	}
	m.log.Info("Manager: ... executor completed handshake! Continuing execution.")
	return nil
}

func (m *SecureProcessorManager) completeExecutorHandshake(handshakeErr error) {
	// Channels can only be closed once. If you try to close it a second time it panics.
	// We have just one go routine that completes the handshake, but to be on the safe side
	// sync.Once guarantees that no matter how many goroutines attempt to mark the handshake complete, the
	// channel is only closed once.
	m.executorHandShake.once.Do(func() {
		if handshakeErr != nil {
			m.log.Error("Manager: Handshake completed with error: %v", handshakeErr)
			m.executorHandShake.err = handshakeErr
		} else {
			m.log.Info("Manager: setting handshake as completed")
		}
		// When a channel is closed all current receivers unblock immediately, all future receives <-ch also
		// return instantly (zero value) and the channel transitions into a permanent done state
		close(m.executorHandShake.isComplete)
	})
}

// Start starts the manager
func (m *SecureProcessorManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("manager is already running")
	}

	m.log.Info("Manager: startup sequence begin")

	// Start the admin command server early so operators can connect even while
	// the executor handshake or blockchain connection is still in progress.
	if m.adminServer != nil {
		adminParams := m.config.AdminChannelParams.(common.TcpChannelConnectionParams)
		m.log.Info("Manager: starting admin command server on %s", adminParams.Url())
		if err := m.adminServer.Start(ctx, "Manager"); err != nil {
			m.log.Warn("Manager: failed to start admin server: %v", err)
			// Don't fail startup if admin server fails - it's not critical
		} else {
			m.log.Info("Manager: admin command server started")
		}
	}

	// Connect to the executor
	m.log.Info("Manager: connecting to executor...")
	if err := m.executorClient.Connect(ctx, "Manager"); err != nil {
		return fmt.Errorf("failed to connect to executor: %w", err)
	}
	m.log.Info("Manager: connected to executor")

	// The handshake ensures the executor has a valid keyset before the manager starts
	// processing any requests from the blockchain, which rely on this keyset for cryptographic operations.
	// The process is as follows:
	// 1. The executor initiates the handshake upon a new connection from the manager.
	// 2. The executor requests the key set recovery data from the manager.
	// 3. If the manager has the data (i.e., this is not the first run), it sends the data to the executor.
	//    The executor restores its keyset and notifies the manager of the successful recovery.
	// 4. If the manager does not have the data (i.e., this is the first run), the executor generates a new
	//    keyset and sends the recovery data to the manager, which stores it for future use.
	if err := m.waitForExecutorHandshake(); err != nil {
		return fmt.Errorf("Manager: executor handshake failed: %w", err)
	}

	// Attempt initial blockchain connection. If it fails, the manager still
	// starts — the polling loop will retry on every tick until connected.
	m.log.Info("Manager: connecting to blockchain node at %s...", m.config.RpcURL)
	if err := m.blockchainClient.Connect(ctx); err != nil {
		m.log.Warn("Manager: initial blockchain connect failed (will retry during polling): %v", err)
	} else {
		m.log.Info("Manager: connected to blockchain node")
	}

	// Start the blockchain polling loop in a goroutine
	m.wg.Add(1)
	go m.pollBlockchain(ctx)

	m.log.Info("Manager: startup sequence complete - Ethereum address: %s", m.config.PrivateKey.PublicKey().Address())

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

	// Signal the polling loop to stop
	select {
	case <-m.stopChan:
	// already closed
	default:
		close(m.stopChan)
	}

	// Wait for the polling loop to stop
	m.wg.Wait()

	// Stop the admin server
	if m.adminServer != nil {
		if err := m.adminServer.Stop(); err != nil {
			m.log.Warn("Manager: Error stopping admin server: %v", err)
		}
	}

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

// HandleGetKeysetRecoveryRequest implements the ClientRequestHandler interface.
func (m *SecureProcessorManager) HandleGetKeysetRecoveryRequest(ctx context.Context) (*common.EnclaveKeySetRecovery, error) {
	m.log.Info("Manager: Received GetKeysetRecovery message from executor")

	recv, err := m.dataLayer.GetEnclaveKeySetRecovery(ctx)
	if err != nil {
		m.log.Info("Manager: could not get keyset recovery data: %v", err)
		return nil, err
	}

	return recv, nil
}

// HandleSetKeysetRecoveryRequest implements communication.ClientRequestHandler.
func (m *SecureProcessorManager) HandleSetKeysetRecoveryRequest(ctx context.Context, recv *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr string) error {
	m.log.Info("Manager: Received SetKeysetRecovery message from executor")

	err := m.dataLayer.StoreEnclaveKeySetRecovery(ctx, recv)

	if err != nil {
		m.log.Error("Manager: Failed to set keyset recovery: %v", err)
		m.completeExecutorHandshake(err)
		return err
	}

	m.log.Info("Manager: KeysetRecovery data stored in data layer")
	m.log.Info("Manager: Executor's public communication key (P521): %s", commPubKey)
	m.log.Info("Manager: Executor's signing key address (Secp256k1): %s", signingKeyAddr)

	// set the handshake as completed
	m.completeExecutorHandshake(nil)
	return nil
}

// HandleKeysetRecoveryResult implements communication.ClientRequestHandler.
func (m *SecureProcessorManager) HandleKeysetRecoveryResult(ctx context.Context, result error, commPubKey, signingKeyAddr string) error {
	if result != nil {
		m.log.Error("Manager: Received KeysetRecoveryResult message from executor with error: %v", result)
		m.completeExecutorHandshake(result)
	} else {
		m.log.Info("Manager: Received KeysetRecoveryResult message from executor with success")
		m.log.Info("Manager: Executor's public communication key (P521): %s", commPubKey)
		m.log.Info("Manager: Executor's signing key address (Secp256k1): %s", signingKeyAddr)
		m.completeExecutorHandshake(nil)
	}
	return nil
}

// forwardToExecutor sends an admin command to the executor via the existing
// communication channel (ForwardAdminCommand) and returns the response data.
func (m *SecureProcessorManager) forwardToExecutor(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
	return m.executorClient.ForwardAdminCommand(ctx, cmdType, data)
}

// forwardSetLogLevel forwards a SetLogLevel command to the executor only.
func (m *SecureProcessorManager) forwardSetLogLevel(ctx context.Context, level string) (*admin.SetLogLevelResponse, error) {
	fwdData, err := json.Marshal(admin.SetLogLevelRequest{Level: level})
	if err != nil {
		return nil, fmt.Errorf("failed to build executor request: %v", err)
	}
	execData, execErr := m.forwardToExecutor(ctx, admin.AdminCmdSetLogLevel, fwdData)
	if execErr != nil {
		return nil, execErr
	}
	var execResp admin.SetLogLevelResponse
	if err := json.Unmarshal(execData, &execResp); err != nil {
		return nil, fmt.Errorf("failed to parse executor response: %v", err)
	}
	return &execResp, nil
}

// forwardGetLogLevel forwards a GetLogLevel command to the executor only.
func (m *SecureProcessorManager) forwardGetLogLevel(ctx context.Context) (string, error) {
	execData, execErr := m.forwardToExecutor(ctx, admin.AdminCmdGetLogLevel, nil)
	if execErr != nil {
		return "", execErr
	}
	var execLevel string
	if err := json.Unmarshal(execData, &execLevel); err != nil {
		return "", fmt.Errorf("failed to parse executor log level: %v", err)
	}
	return execLevel, nil
}

// ExecuteCommand processes an admin command and returns the result.
// Supported commands: KeyAttestation, GetVersion, SetLogLevel, GetLogLevel.
// KeyAttestation is always forwarded to the Executor.
// For SetLogLevel/GetLogLevel, target controls routing:
//   - "manager": applied locally only
//   - "executor": forwarded to Executor only
//   - "all" or "": applied locally AND forwarded to Executor (aggregated response)
func (m *SecureProcessorManager) ExecuteCommand(ctx context.Context, msg admin.AdminMessage) (interface{}, error) {
	switch msg.Type {
	case admin.KeyAttestationRequestMessage:
		m.log.Info("Manager: KeyAttestation command received, forwarding to executor")
		respData, err := m.forwardToExecutor(ctx, admin.AdminCmdKeyAttestation, msg.Data)
		if err != nil {
			return nil, err
		}
		return json.RawMessage(respData), nil

	case admin.GetVersionRequestMessage:
		return m.GetVersion(ctx)

	case admin.SetLogLevelRequestMessage:
		var req admin.SetLogLevelRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			return nil, fmt.Errorf("invalid request data: %w", err)
		}
		if err := admin.ValidateTarget(req.Target); err != nil {
			return nil, err
		}

		if req.Target == "" {
			m.log.Warn("Manager: SetLogLevel received without target, defaulting to 'all' (applies to both manager and executor)")
			req.Target = admin.TargetAll
		}

		if req.Target == admin.TargetExecutor {
			return m.forwardSetLogLevel(ctx, req.Level)
		}

		if req.Target == admin.TargetAll {
			// Apply locally first.
			mgrResult, mgrErr := m.SetLogLevel(ctx, req.Level)
			if mgrErr != nil {
				return nil, mgrErr
			}
			mgrResp, ok := mgrResult.(admin.SetLogLevelResponse)
			if !ok {
				return nil, fmt.Errorf("unexpected response type from local SetLogLevel")
			}

			// Forward to executor via communication channel.
			execResp, execErr := m.forwardSetLogLevel(ctx, req.Level)
			if execErr != nil {
				return nil, fmt.Errorf("manager level set to '%s' successfully, but executor failed: %v", req.Level, execErr)
			}

			return admin.AggregatedSetLogLevelResponse{
				Manager:  mgrResp,
				Executor: *execResp,
			}, nil
		}

		return m.SetLogLevel(ctx, req.Level)

	case admin.GetLogLevelRequestMessage:
		var req admin.GetLogLevelRequest
		if msg.Data != nil && string(msg.Data) != "null" {
			if err := json.Unmarshal(msg.Data, &req); err != nil {
				return nil, fmt.Errorf("invalid request data: %w", err)
			}
		}
		if err := admin.ValidateTarget(req.Target); err != nil {
			return nil, err
		}

		if req.Target == "" {
			m.log.Warn("Manager: GetLogLevel received without target, defaulting to 'all' (applies to both manager and executor)")
			req.Target = admin.TargetAll
		}

		if req.Target == admin.TargetExecutor {
			execLevel, err := m.forwardGetLogLevel(ctx)
			if err != nil {
				return nil, err
			}
			return execLevel, nil
		}

		if req.Target == admin.TargetAll {
			mgrLevel, mgrErr := m.GetLogLevel(ctx)
			if mgrErr != nil {
				return nil, mgrErr
			}

			execLevel, execErr := m.forwardGetLogLevel(ctx)
			if execErr != nil {
				return nil, fmt.Errorf("failed to get executor log level: %v", execErr)
			}

			return admin.AggregatedGetLogLevelResponse{
				Manager:  mgrLevel,
				Executor: execLevel,
			}, nil
		}

		return m.GetLogLevel(ctx)

	default:
		return nil, fmt.Errorf("unsupported command type: %v", msg.Type)
	}
}

// GetVersion returns the current version of the manager.
func (m *SecureProcessorManager) GetVersion(ctx context.Context) (string, error) {
	m.log.Info("Manager: GetVersion command received, returning version %s", Version)
	return Version, nil
}

// SetLogLevel changes the manager's log level at runtime.
func (m *SecureProcessorManager) SetLogLevel(ctx context.Context, level string) (interface{}, error) {
	return admin.HandleSetLogLevel(m.log, "Manager", level)
}

// GetLogLevel returns the current log level of the manager.
func (m *SecureProcessorManager) GetLogLevel(ctx context.Context) (string, error) {
	return admin.HandleGetLogLevel(m.log, "Manager")
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
		case <-m.stopChan:
			return
		case <-ticker.C:
			err := m.processRequestFromChain(ctx)
			if err != nil {
				m.log.Error("Manager: Fatal error processing requests from chain: %v, initiating shutdown", err)
				select {
				case m.fatalErrChan <- err:
				default:
					// Channel already has an error, don't block
				}
				return
			}
		}
	}
}

// processRequestFromChain retrieves the next pending request from the blockchain and processes it
func (m *SecureProcessorManager) processRequestFromChain(ctx context.Context) error {
	// Attempt blockchain (re)connect before acquiring the manager lock.
	// Connect has its own internal mutex and may block for up to the
	// configured connect timeout; holding m.mu.RLock during that period
	// would delay Stop().
	if !m.blockchainClient.IsConnected() {
		if err := m.blockchainClient.Connect(ctx); err != nil {
			m.log.Warn("Manager: blockchain not connected, retrying next poll: %v", err)
			return nil
		}
		m.log.Info("Manager: connected to blockchain node at %s", m.config.RpcURL)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.isRunning {
		m.log.Warn("Manager is not started yet, skipping")
		return nil
	}

	// Get next pending request from the blockchain
	request, stateRoot, err := m.blockchainClient.GetNextPendingRequest(ctx)
	if err != nil {
		m.log.Error("Manager: Failed to get pending request: %v", err)
		return nil
	}

	localStateRoot, err := m.dataLayer.LastVersionID()
	if err != nil {
		if dbErr, ok := err.(*storageErrors.Error); ok && dbErr.Code == storageErrors.NoVersionInDb {
			localStateRoot = make([]byte, 32) // Initialize to zero state root if no version exists
		} else {
			m.log.Error("Manager: Failed to get local state root: %v", err)
			return nil
		}
	}

	if !bytes.Equal(localStateRoot, stateRoot[:]) {
		m.log.Warn("Manager: State root mismatch, expected %x, got %x. Checking if it is a REORG.", localStateRoot, stateRoot)

		isReorg, err := m.checkIfReorg(stateRoot)
		if err != nil {
			m.log.Error("Manager: Failed to check for reorg: %v", err)
			return nil
		}

		if isReorg {
			if m.endReorgTime.IsZero() {
				m.log.Info("Manager: Starting REORG timeout %d", m.config.ReorgTimeout)
				m.endReorgTime = time.Now().Add(time.Duration(m.config.ReorgTimeout) * time.Second)
				return nil
			}
			if time.Now().Before(m.endReorgTime) {
				m.log.Info("Manager: REORG timeout not expired yet. Keep waiting...")
				return nil
			}
			m.log.Info("Manager: REORG not solved within timeout => Rollback the DB")
			if err := m.dataLayer.Rollback(stateRoot[:]); err != nil {
				m.log.Error("Manager: Error while rolling back the DB: %v", err)
				return fmt.Errorf("fatal: rollback failed: %w", err)
			}

		} else {
			m.log.Error("Manager: unrecoverable disalignment between DB and chain, no matching state root found in db")
			emptyStateRoot := [32]byte{}
			if bytes.Equal(localStateRoot, emptyStateRoot[:]) {
				m.log.Error("Manager: the DB is empty but the chain state root is non-zero, check if the database file is correct and restart the manager")
			}
			return fmt.Errorf("unrecoverable disalignment between DB and chain, no matching state root found in db")
		}

	}

	if !m.endReorgTime.IsZero() {
		m.log.Info("Manager: State roots match, REORG resolved")
		m.endReorgTime = time.Time{}
	}

	if request == nil {
		return nil
	}

	m.log.Info("Manager: processing request %s", request.RequestID)

	if rf := m.processRequest(ctx, request); rf != nil {
		// Log the error and mark the request as failed
		m.log.Error("Manager: Failed to process request %s: %v", request.RequestID, rf)
		if err = m.blockchainClient.MarkRequestFailed(ctx, request.RequestID, rf); err != nil {
			m.log.Error("Manager: Failed to mark request %s as failed: %v", request.RequestID, err)
		}
	} else {
		m.log.Info("Manager: Processed request %s", request.RequestID)
	}
	return nil

}

func (m *SecureProcessorManager) checkIfReorg(stateRoot [32]byte) (bool, error) {
	if stateRoot == [32]byte{} {
		m.log.Info("Manager: State root is zero, REORG")
		// Don't look for older db versions, just mark as reorged and wait for next poll
		return true, nil
	}

	oldVersions, err := m.dataLayer.ListVersions()
	if err != nil {
		m.log.Error("Manager: Failed to get db old versions: %v", err)
		return false, err
	}

	for _, oldVersion := range oldVersions[1:] {
		if bytes.Equal(oldVersion, stateRoot[:]) {
			m.log.Info("Manager: Found matching state root %x in db, REORG. Checking if the timeout is expired", stateRoot)
			return true, nil
		}
	}
	return false, nil

}

// processRequest processes a request
func (m *SecureProcessorManager) processRequest(ctx context.Context, req *common.Request) *apperrors.RequestFailure {
	if !m.isRunning {
		m.log.Warn("Manager is not started yet, skipping")
		return apperrors.New(apperrors.CodeManagerNotRunning, "manager is not running", nil)
	}

	// check admitted appID. TODO: This check will be removed in future
	if req.ApplicationID != admittedAppID {
		return apperrors.New(apperrors.CodeAppNotAdmitted, fmt.Sprintf("application id %s is not admitted", req.ApplicationID), nil)
	}

	switch req.RequestType {
	case common.Deploy:
		return m.processDeployApp(ctx, req)
	case common.Process, common.AssociateKey, common.Deanonymize:
		return m.processProcessRequest(ctx, req)
	default:
		return apperrors.New(apperrors.CodeRequestTypeNotPermitted, fmt.Sprintf("unsupported request type: %s", req.RequestType), nil)
	}
}

func (m *SecureProcessorManager) submitStateOnChain(ctx context.Context, updatePayload *common.UpdatePayload) *apperrors.RequestFailure {
	// Submit the state update to the blockchain
	if err := m.blockchainClient.SubmitStateUpdate(ctx, updatePayload); err != nil {
		m.log.Error("Failed to submit state update for error: %v", err)
		m.log.Info("Rollback the application state to previous version")
		if err := m.dataLayer.Rollback(updatePayload.PrevStateRoot[:]); err != nil {
			// If this happens, the local db and the chain are out of sync and cannot be recovered automatically.
			// Log and return nil to let REORG detection handle it on the next poll.
			m.log.Error("Failed to rollback application state: %v. Will retry via REORG detection.", err)
			return nil
		}

		if _, ok := err.(blockchain.ReorgError); ok {
			m.log.Warn("REORG, do not call MarkFailed, wait for next poll")
			return nil
		}
		return apperrors.New(apperrors.CodeSubmittingStateUpdateFailed, fmt.Sprintf("failed to submit state update: %v", err), err)
	}

	m.log.Info("Manager: Processed request %s for application %d", updatePayload.RequestID, updatePayload.ApplicationID)
	return nil
}

// processDeployApp processes a deploy app request
func (m *SecureProcessorManager) processDeployApp(ctx context.Context, req *common.Request) *apperrors.RequestFailure {
	m.log.Info("Processing deploy app request: %s", req.RequestID)
	if !m.isRunning {
		m.log.Warn("Manager is not started yet, skipping")
		return apperrors.New(apperrors.CodeManagerNotRunning, "manager is not started yet", nil)
	}

	// check if app was already deployed
	_, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err == nil {
		return apperrors.New(apperrors.CodeApplicationAlreadyDeployed, fmt.Sprintf("application %s was already deployed", req.ApplicationID), err)
	}

	//if the payload is empty, try to retrieve the wasm locally
	if len(req.Payload) == 0 {
		m.log.Info("Empty payload received - trying to retrieve wasm locally")
		wasmFilePath := filepath.Join(m.config.InputWasmPath, req.ApplicationID.String()+".wasm")
		if !common.FileExists(wasmFilePath) {
			return apperrors.New(apperrors.CodeWasmNotFound, fmt.Sprintf("failed to deploy application - wasm not found in both payload or local path: %v", wasmFilePath), err)
		}
		wasmBytesFromFile, err := os.ReadFile(wasmFilePath)
		if err != nil {
			return apperrors.New(apperrors.CodeErrorReadingWasm, fmt.Sprintf("failed to deploy application - Error reading wasm file: %v", err), err)

		}
		req.Payload = wasmBytesFromFile
	}

	// Deploy the application
	updatePayload, appState, failure := m.executorClient.SendDeployApp(ctx, req)
	if failure != nil {
		return failure
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
		m.log.Error("failed to submit state update: %v", err)
		return nil
	}

	m.log.Info("Deployed application, submit the state update to the blockchain")
	return m.submitStateOnChain(ctx, updatePayload)

}

// processProcessRequest processes a process request (including deanonymization requests)
func (m *SecureProcessorManager) processProcessRequest(ctx context.Context, req *common.Request) *apperrors.RequestFailure {
	m.log.Info("Processing Process app request: %s (type: %s)", req.RequestID, req.RequestType)
	if !m.isRunning {
		m.log.Warn("Manager is not started yet, skipping")
		return apperrors.New(apperrors.CodeManagerNotRunning, "manager is not started yet", nil)
	}

	// Get the application state
	appState, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err != nil {
		m.log.Error("GetApplicationState returns an error: %v", err)
		if strings.Contains(err.Error(), "application state not found") {
			// This can happen if the deploy transaction was not mined yet, mark request as failed
			return apperrors.New(apperrors.CodeAppStateNotFound, fmt.Sprintf("state not found for application %s", req.ApplicationID), err)
		}
		// Other errors are likely db errors, retry on next poll
		return nil
	}

	// Get the WASM module for the application
	wasmBytes, err := m.dataLayer.GetWASMBytecode(ctx, req.ApplicationID)
	if err != nil {
		m.log.Error("GetWASMBytecode returns an error: %v", err)
		return nil
	}

	// Process the request
	updatePayload, updatedState, deanonymizationReport, failure := m.executorClient.SendProcessRequest(ctx, req, appState, wasmBytes)
	if failure != nil {
		return failure
	}

	// If a deanonymization report was generated, save it to filesystem
	// Note: The executor already validates that reports are only and always generated for Deanonymize requests
	if deanonymizationReport != nil {
		if err := m.saveDeanonymizationReport(deanonymizationReport, req); err != nil {
			// Treat as transient: log and retry on next poll instead of failing the request.
			m.log.Error("Failed to save deanonymization report: %v", err)
			return nil
		}
		m.log.Info("Manager: Saved deanonymization report %s for application %d", deanonymizationReport.ReportID, req.ApplicationID)
	}

	// Store the updated application state
	versionID := updatedState.StateRoot[:]
	m.log.Info("VersionID %x", versionID[:])

	err = m.dataLayer.Store(
		ctx,
		versionID[:],
		[]*common.ApplicationState{updatedState},
		nil,
	)
	if err != nil {
		m.log.Error("failed to submit state update: %v", err)
		return nil
	}

	return m.submitStateOnChain(ctx, updatePayload)
}

// saveDeanonymizationReport saves a deanonymization report to the filesystem
func (m *SecureProcessorManager) saveDeanonymizationReport(report *common.DeanonymizationReport, req *common.Request) error {
	if err := os.MkdirAll(m.config.DeanonymizationReportPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory for deanonymization reports: %w", err)
	}

	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal deanonymization report to JSON: %w", err)
	}

	filePath := filepath.Join(m.config.DeanonymizationReportPath, common.ReportFilename(req.ApplicationID, req.RequestID))
	if err := os.WriteFile(filePath, reportJSON, 0644); err != nil {
		return fmt.Errorf("failed to write deanonymization report to file: %w", err)
	}

	return nil
}
