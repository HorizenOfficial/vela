package manager

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/HorizenOfficial/vela/pkg/admin"
	"github.com/HorizenOfficial/vela/pkg/blockchain"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/communication"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/storage"
	storageErrors "github.com/HorizenOfficial/vela/pkg/storage/errors"
	"github.com/HorizenOfficial/vela/pkg/version"
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
	// handShakeMu guards the executorHandShake pointer (NOT its internal fields,
	// which are protected by their own sync.Once). The comm layer dispatches
	// message handlers on detached goroutines (shared_impl.go), so a late
	// handshake handler from a drained connection can call
	// completeExecutorHandshake concurrently with reconnectExecutor swapping in a
	// fresh tracker — an unsynchronized pointer read/write without this lock.
	handShakeMu       sync.RWMutex
	stopChan          chan struct{} // Channel to signal the polling loop to stop
	wg                sync.WaitGroup
	endReorgTime      time.Time
	log               logger.Logger
	fatalErrChan      chan error // Channel to signal fatal errors to main

	// runningPcr0 is the hex-encoded PCR0 reported by the executor at handshake
	// (the "running image"); empty in non-Nitro (TCP/dev) mode. Guarded by
	// identityMu, a dedicated lock: the handshake handlers that write these run
	// on the comm client's reader goroutine while Start holds m.mu, so reusing
	// m.mu here would deadlock.
	identityMu  sync.RWMutex
	runningPcr0 string
	// runningImageHash caches keccak256(pcr0) — the value compared against the
	// on-chain activeImage — computed once per handshake in setExecutorIdentity.
	// nil when no PCR0 was reported (dev/TCP mode); a bad-hex PCR0 is rejected at
	// handshake and never reaches here.
	runningImageHash *ethCommon.Hash
	// reportedSigner is the executor's signing-key address reported at handshake
	// (hex, "0x…"), used for the signer-continuity check against on-chain teeSigner.
	reportedSigner string

	// draining is set once a PCR0/activeImage mismatch is observed: the manager
	// stops dispatching new requests, signals the executor to shut down, and
	// waits to reconnect to the swapped-in enclave.
	//
	// signerVerified records that the post-handshake signer-continuity check has
	// passed for the current connection; reset on each (re)handshake.
	//
	// Both are guarded by identityMu (NOT m.mu): the handshake handlers that
	// reset them run on the comm reader goroutine while Start holds m.mu, so
	// using m.mu here would deadlock (see Task 3). Lock order is always
	// m.mu → identityMu; nothing takes identityMu then m.mu.
	draining       bool
	signerVerified bool
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
		log:           log,
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

// currentHandShake returns the active handshake tracker under handShakeMu.
func (m *SecureProcessorManager) currentHandShake() *ExecutorHandShake {
	m.handShakeMu.RLock()
	defer m.handShakeMu.RUnlock()
	return m.executorHandShake
}

// resetHandShake installs a fresh single-shot handshake tracker for a new
// connection and returns it. The swap is guarded by handShakeMu so it cannot
// race a concurrent pointer read from a (possibly stale) handshake handler
// goroutine.
func (m *SecureProcessorManager) resetHandShake() *ExecutorHandShake {
	m.handShakeMu.Lock()
	defer m.handShakeMu.Unlock()
	m.executorHandShake = &ExecutorHandShake{isComplete: make(chan struct{})}
	return m.executorHandShake
}

func (m *SecureProcessorManager) waitForExecutorHandshake() error {
	m.log.Info("Manager: Waiting for the executor to complete the handshake...")
	hs := m.currentHandShake()

	select {
	case <-hs.isComplete:
		// Handshake completed (successfully or with error)
	case <-time.After(time.Duration(m.config.HandshakeTimeout) * time.Second):
		// Handshake timed out
		hs.once.Do(func() {
			m.log.Warn("Manager: Handshake timed out after %d seconds", m.config.HandshakeTimeout)
			hs.err = fmt.Errorf("handshake timed out after %d seconds", m.config.HandshakeTimeout)
			close(hs.isComplete)
		})
	}

	if hs.err != nil {
		m.log.Error("Manager: ... executor completed handshake with error: %v", hs.err)
		return hs.err
	}
	m.log.Info("Manager: ... executor completed handshake! Continuing execution.")
	return nil
}

func (m *SecureProcessorManager) completeExecutorHandshake(handshakeErr error) {
	// Channels can only be closed once. If you try to close it a second time it panics.
	// We have just one go routine that completes the handshake, but to be on the safe side
	// sync.Once guarantees that no matter how many goroutines attempt to mark the handshake complete, the
	// channel is only closed once.
	hs := m.currentHandShake()
	hs.once.Do(func() {
		if handshakeErr != nil {
			m.log.Error("Manager: Handshake completed with error: %v", handshakeErr)
			hs.err = handshakeErr
		} else {
			m.log.Info("Manager: setting handshake as completed")
		}
		// When a channel is closed all current receivers unblock immediately, all future receives <-ch also
		// return instantly (zero value) and the channel transitions into a permanent done state
		close(hs.isComplete)
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
	if !m.isRunning {
		m.mu.Unlock()
		return nil
	}

	// Flip isRunning and signal the polling loop to stop under m.mu, then RELEASE
	// m.mu before waiting. processRequestFromChain does its blocking work
	// (blockchain reconnect, maintainExecutorForActiveImage) lock-free and only
	// takes m.mu at dispatch; if Stop held m.mu across wg.Wait(), a tick already in
	// that pre-lock section would block forever at m.mu.Lock() while Stop blocked
	// forever in wg.Wait() — a lock-ordering deadlock. With isRunning already
	// false, a mid-flight tick instead acquires m.mu, sees !isRunning, and exits
	// without dispatching.
	m.isRunning = false
	select {
	case <-m.stopChan:
	// already closed
	default:
		close(m.stopChan)
	}
	m.mu.Unlock()

	// Wait for the polling loop to stop (without holding m.mu).
	m.wg.Wait()

	// Teardown runs without m.mu: the polling goroutine has exited, the clients
	// have their own internal locks, and isRunning is already false so no request
	// path will dispatch. If a teardown step fails, restore isRunning=true so the
	// stop is not considered complete and the operator can retry Stop().
	stopFailed := func(err error) error {
		m.mu.Lock()
		m.isRunning = true
		m.mu.Unlock()
		return err
	}

	// Stop the admin server
	if m.adminServer != nil {
		if err := m.adminServer.Stop(); err != nil {
			m.log.Warn("Manager: Error stopping admin server: %v", err)
		}
	}

	// Close the executor client
	if err := m.executorClient.Close(); err != nil {
		return stopFailed(fmt.Errorf("failed to close executor client: %w", err))
	}

	// Close the blockchain client
	if err := m.blockchainClient.Close(); err != nil {
		return stopFailed(fmt.Errorf("failed to close blockchain client: %w", err))
	}

	// Close the data layer
	if err := m.dataLayer.Close(); err != nil {
		return stopFailed(fmt.Errorf("failed to close data layer: %w", err))
	}

	return nil
}

// HandleGetKeysetRecoveryRequest implements the ClientRequestHandler interface.
func (m *SecureProcessorManager) HandleGetKeysetRecoveryRequest(ctx context.Context, peerProtocolVersion uint32) (*common.EnclaveKeySetRecovery, error) {
	m.log.Info("Manager: Received GetKeysetRecovery message from executor (protocol version %d)", peerProtocolVersion)

	// Reject an incompatible executor before any keyset-recovery data is read or
	// exchanged. Fail the handshake with a typed error so startup aborts cleanly
	// instead of timing out.
	if !communication.IsCompatible(peerProtocolVersion) {
		err := &communication.IncompatibleProtocolError{Local: communication.WireProtocolVersion, Peer: peerProtocolVersion}
		m.log.Error("Manager: %v — failing handshake, no recovery data exchanged", err)
		m.completeExecutorHandshake(err)
		return nil, err
	}

	recv, err := m.dataLayer.GetEnclaveKeySetRecovery(ctx)
	if err != nil {
		m.log.Info("Manager: could not get keyset recovery data: %v", err)
		return nil, err
	}

	return recv, nil
}

// HandleSetKeysetRecoveryRequest implements communication.ClientRequestHandler.
func (m *SecureProcessorManager) HandleSetKeysetRecoveryRequest(ctx context.Context, recv *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr, pcr0 string) error {
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
	m.log.Info("Manager: Executor's pcr0: %s", pcr0)
	if err := m.setExecutorIdentity(pcr0, signingKeyAddr); err != nil {
		m.log.Error("Manager: %v", err)
		m.completeExecutorHandshake(err)
		return err
	}

	// set the handshake as completed
	m.completeExecutorHandshake(nil)
	return nil
}

// HandleKeysetRecoveryResult implements communication.ClientRequestHandler.
func (m *SecureProcessorManager) HandleKeysetRecoveryResult(ctx context.Context, result error, commPubKey, signingKeyAddr, pcr0 string) error {
	if result != nil {
		m.log.Error("Manager: Received KeysetRecoveryResult message from executor with error: %v", result)
		m.completeExecutorHandshake(result)
	} else {
		m.log.Info("Manager: Received KeysetRecoveryResult message from executor with success")
		m.log.Info("Manager: Executor's public communication key (P521): %s", commPubKey)
		m.log.Info("Manager: Executor's signing key address (Secp256k1): %s", signingKeyAddr)
		if err := m.setExecutorIdentity(pcr0, signingKeyAddr); err != nil {
			m.log.Error("Manager: %v", err)
			m.completeExecutorHandshake(err)
			return err
		}
		m.completeExecutorHandshake(nil)
	}
	return nil
}

// setExecutorIdentity records the running image PCR0 and signing-key address the
// executor reported at handshake. The PCR0 is the "running image" compared (as
// keccak256(pcr0)) against the on-chain activeImage; the signer feeds the
// signer-continuity check. A fresh handshake re-arms the signer check (a
// reconnect may bring a different enclave).
//
// The image hash is computed once here rather than on every activeImage check.
// An empty PCR0 (dev/TCP mode) leaves the hash nil; a non-empty PCR0 that is not
// valid hex is rejected — the identity is left unchanged and the handshake fails.
func (m *SecureProcessorManager) setExecutorIdentity(pcr0, signingKeyAddr string) error {
	var imageHash *ethCommon.Hash
	if pcr0 != "" {
		rawPcr0, err := hex.DecodeString(pcr0)
		if err != nil {
			return fmt.Errorf("executor reported PCR0 %q is not valid hex: %w", pcr0, err)
		}
		h := ethCrypto.Keccak256Hash(rawPcr0)
		imageHash = &h
	}

	m.identityMu.Lock()
	m.runningPcr0 = pcr0
	m.runningImageHash = imageHash
	m.reportedSigner = signingKeyAddr
	m.signerVerified = false // re-arm the signer check for this (re)handshake
	m.identityMu.Unlock()

	m.log.Info("Manager: Executor's running image PCR0: %q, signer: %q", pcr0, signingKeyAddr)
	return nil
}

// runningImageHashCopy returns the cached keccak256(pcr0) image hash under the
// identity lock, or nil in dev/TCP mode (no PCR0 reported).
func (m *SecureProcessorManager) runningImageHashCopy() *ethCommon.Hash {
	m.identityMu.RLock()
	defer m.identityMu.RUnlock()
	return m.runningImageHash
}

// RunningPcr0 returns the hex-encoded PCR0 the executor reported at handshake
// (the running image), or "" if not reported (non-Nitro/dev mode).
func (m *SecureProcessorManager) RunningPcr0() string {
	m.identityMu.RLock()
	defer m.identityMu.RUnlock()
	return m.runningPcr0
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
	if unmarshalErr := json.Unmarshal(execData, &execResp); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to parse executor response: %v", unmarshalErr)
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

// forwardShutdown asks the executor to exit cleanly (drain step 3). The executor
// acks and then closes the connection, so a disconnect immediately after sending
// is an expected outcome, not a failure — callers should not treat it as fatal.
func (m *SecureProcessorManager) forwardShutdown(ctx context.Context) error {
	_, err := m.forwardToExecutor(ctx, admin.AdminCmdShutdown, nil)
	if err != nil {
		// The executor may drop the connection before/around acking; log and
		// swallow — the shutdown intent was delivered.
		m.log.Warn("Manager: shutdown ack not received (executor likely already exiting): %v", err)
	}
	return nil
}

// forwardGetVersion forwards a GetVersion command to the executor only.
func (m *SecureProcessorManager) forwardGetVersion(ctx context.Context) (string, error) {
	execData, err := m.forwardToExecutor(ctx, admin.AdminCmdGetVersion, nil)
	if err != nil {
		return "", err
	}
	var v string
	if err := json.Unmarshal(execData, &v); err != nil {
		return "", fmt.Errorf("failed to parse executor version: %v", err)
	}
	return v, nil
}

// ExecuteCommand processes an admin command and returns the result.
// Supported commands: KeyAttestation, GetVersion, SetLogLevel, GetLogLevel.
// KeyAttestation is always forwarded to the Executor.
// For GetVersion/SetLogLevel/GetLogLevel, msg.Target controls routing:
//   - "manager": applied locally only
//   - "executor": forwarded to Executor only
//   - "all" or "": applied locally AND forwarded to Executor (aggregated response)
func (m *SecureProcessorManager) ExecuteCommand(ctx context.Context, msg admin.AdminMessage) (interface{}, error) {
	// Validate and default the target once, before dispatching.
	if err := admin.ValidateTarget(msg.Target); err != nil {
		return nil, err
	}
	target := msg.Target
	if target == "" {
		m.log.Warn("Manager: %s received without target, defaulting to 'all' (applies to both manager and executor)", msg.Type)
		target = admin.TargetAll
	}

	switch msg.Type {
	case admin.KeyAttestationRequestMessage:
		if target != admin.TargetExecutor && target != admin.TargetAll {
			return nil, fmt.Errorf("key_attestation is only supported on the executor; valid targets: '%s', '%s'", admin.TargetExecutor, admin.TargetAll)
		}
		m.log.Info("Manager: KeyAttestation command received, forwarding to executor")
		respData, err := m.forwardToExecutor(ctx, admin.AdminCmdKeyAttestation, msg.Data)
		if err != nil {
			return nil, err
		}
		return json.RawMessage(respData), nil

	case admin.GetVersionRequestMessage:
		if target == admin.TargetExecutor {
			return m.forwardGetVersion(ctx)
		}
		if target == admin.TargetAll {
			resp := admin.AggregatedGetVersionResponse{}
			mgrVersion, mgrErr := m.GetVersion(ctx)
			if mgrErr != nil {
				resp.ManagerError = mgrErr.Error()
			} else {
				resp.Manager = mgrVersion
			}
			execVersion, execErr := m.forwardGetVersion(ctx)
			if execErr != nil {
				resp.ExecutorError = execErr.Error()
			} else {
				resp.Executor = execVersion
			}
			resp.ExecutorPcr0 = m.RunningPcr0()
			return resp, nil
		}
		return m.GetVersion(ctx)

	case admin.SetLogLevelRequestMessage:
		if msg.Data == nil || string(msg.Data) == "null" {
			return nil, fmt.Errorf("missing request data for set_log_level")
		}
		var req admin.SetLogLevelRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			return nil, fmt.Errorf("invalid request data: %w", err)
		}

		if target == admin.TargetExecutor {
			return m.forwardSetLogLevel(ctx, req.Level)
		}

		if target == admin.TargetAll {
			resp := admin.AggregatedSetLogLevelResponse{}
			mgrResult, mgrErr := m.SetLogLevel(ctx, req.Level)
			if mgrErr != nil {
				resp.ManagerError = mgrErr.Error()
			} else {
				resp.Manager = mgrResult.Level
			}

			execResp, execErr := m.forwardSetLogLevel(ctx, req.Level)
			if execErr != nil {
				resp.ExecutorError = execErr.Error()
			} else {
				resp.Executor = execResp.Level
			}

			return resp, nil
		}

		result, err := m.SetLogLevel(ctx, req.Level)
		if err != nil {
			return nil, err
		}
		return result, nil

	case admin.GetLogLevelRequestMessage:
		if target == admin.TargetExecutor {
			execLevel, err := m.forwardGetLogLevel(ctx)
			if err != nil {
				return nil, err
			}
			return execLevel, nil
		}

		if target == admin.TargetAll {
			resp := admin.AggregatedGetLogLevelResponse{}
			mgrLevel, mgrErr := m.GetLogLevel(ctx)
			if mgrErr != nil {
				resp.ManagerError = mgrErr.Error()
			} else {
				resp.Manager = mgrLevel
			}

			execLevel, execErr := m.forwardGetLogLevel(ctx)
			if execErr != nil {
				resp.ExecutorError = execErr.Error()
			} else {
				resp.Executor = execLevel
			}

			return resp, nil
		}

		return m.GetLogLevel(ctx)

	case admin.SetWasmCacheSizeRequestMessage, admin.GetWasmCacheSizeRequestMessage:
		if target != admin.TargetExecutor && target != admin.TargetAll {
			return nil, fmt.Errorf("%s is only supported on the executor; valid targets: '%s', '%s'", msg.Type, admin.TargetExecutor, admin.TargetAll)
		}
		cmdType := admin.AdminCmdSetWasmCacheSize
		if msg.Type == admin.GetWasmCacheSizeRequestMessage {
			cmdType = admin.AdminCmdGetWasmCacheSize
		}
		respData, err := m.forwardToExecutor(ctx, cmdType, msg.Data)
		if err != nil {
			return nil, err
		}
		return json.RawMessage(respData), nil

	default:
		return nil, fmt.Errorf("unsupported command type: %v", msg.Type)
	}
}

// GetVersion returns the current version of the manager.
func (m *SecureProcessorManager) GetVersion(ctx context.Context) (string, error) {
	m.log.Info("Manager: GetVersion command received, returning version %s", version.Version)
	return version.Version, nil
}

// SetLogLevel changes the manager's log level at runtime.
func (m *SecureProcessorManager) SetLogLevel(ctx context.Context, level string) (admin.SetLogLevelResponse, error) {
	return admin.HandleSetLogLevel(m.log, "Manager", level)
}

// GetLogLevel returns the current log level of the manager.
func (m *SecureProcessorManager) GetLogLevel(ctx context.Context) (string, error) {
	return admin.HandleGetLogLevel(m.log, "Manager")
}

// SignerContinuityError signals that the executor's signing-key address no
// longer matches the on-chain teeSigner — i.e. the enclave regenerated its
// keyset instead of recovering it (silent orphan). It is fatal.
type SignerContinuityError struct {
	Reported string
	OnChain  string
}

func (e *SignerContinuityError) Error() string {
	return fmt.Sprintf("signer-continuity check failed: executor signer %q != on-chain teeSigner %q (keyset likely regenerated)", e.Reported, e.OnChain)
}

func (m *SecureProcessorManager) isDraining() bool {
	m.identityMu.RLock()
	defer m.identityMu.RUnlock()
	return m.draining
}

// maintainExecutorForActiveImage reconciles the connected enclave against the
// on-chain activeImage before any request is dispatched (level-based on
// activeImage, so a Manager restart converges to the same outcome). It returns
// skip=true when the caller must not dispatch this tick (transient read error,
// draining, or reconnecting), and a non-nil error only for the fatal
// signer-continuity mismatch.
func (m *SecureProcessorManager) maintainExecutorForActiveImage(ctx context.Context) (skip bool, fatal error) {
	// Bring the executor back before dispatching if either the channel is down.
	// Two triggers: (1) we drained it for a TEE swap; (2) the channel dropped
	// underneath us (executor crash/restart, same image — draining is never set
	// in this case). Both need a re-dial + re-handshake, which also refreshes the
	// reported PCR0/signer we verify below; without this the Manager would retry
	// into a dead connection forever and the relaunched executor, waiting for a
	// handshake that never comes, could not recover its keyset.
	if m.isDraining() || !m.executorClient.IsConnected() {
		if err := m.reconnectExecutor(ctx); err != nil {
			m.log.Warn("Manager: executor not connected yet, will retry next tick: %v", err)
			return true, nil
		}
		// Reconnected + re-handshaked; fall through to re-evaluate the image.
	}

	if runningHash := m.runningImageHashCopy(); runningHash != nil {
		// Nitro executor: reconcile the running image against the on-chain
		// activeImage. Dev/TCP executors report no PCR0 (nil hash) and skip this
		// block, but still run the signer-continuity check below (it does not need
		// PCR0). The hash was computed once at handshake (setExecutorIdentity).
		activeImage, err := m.blockchainClient.GetActiveImage(ctx)
		if err != nil {
			// Transient chain read error during the swap window must not be fatal.
			m.log.Warn("Manager: could not read activeImage (retry next tick): %v", err)
			return true, nil
		}

		if !bytes.Equal(runningHash[:], activeImage[:]) {
			m.log.Warn("Manager: running image keccak(pcr0)=%s != activeImage %x — TEE swap detected, draining",
				runningHash.Hex(), activeImage)
			m.drainExecutor(ctx)
			return true, nil
		}
	}

	// Verify signer continuity once per handshake before dispatching anything.
	// Runs in every channel mode (it compares the reported signing key against
	// on-chain teeSigner and does not depend on PCR0), so TCP deployments and the
	// system suite exercise it too.
	if skip, err := m.verifySignerContinuity(ctx); err != nil {
		return false, err
	} else if skip {
		return true, nil
	}
	return false, nil
}

// drainExecutor enters drain mode: it signals the executor to exit cleanly and
// closes the channel. Idempotent within a drain episode (the shutdown is sent
// once; reconnectExecutor clears the flag on a successful re-handshake).
func (m *SecureProcessorManager) drainExecutor(ctx context.Context) {
	m.identityMu.Lock()
	already := m.draining
	m.draining = true
	m.identityMu.Unlock()
	if already {
		return
	}

	m.log.Warn("Manager: entering drain — signaling executor shutdown and closing channel")
	// forwardShutdown tolerates the post-ack disconnect; never fatal.
	_ = m.forwardShutdown(ctx)
	if err := m.executorClient.Close(); err != nil {
		m.log.Warn("Manager: error closing executor channel during drain: %v", err)
	}
}

// reconnectExecutor re-dials the executor and performs a fresh handshake. It
// re-arms the single-shot handshake state and, on success, exits drain mode.
// Called only from the polling goroutine (before m.mu), so it may block on the
// handshake without risking Stop() or the handshake-handler deadlock.
func (m *SecureProcessorManager) reconnectExecutor(ctx context.Context) error {
	m.log.Info("Manager: attempting to reconnect to executor...")
	// Install a fresh single-shot handshake tracker for the new connection.
	// resetHandShake swaps the pointer under handShakeMu so it cannot race a late
	// handshake handler from the drained connection (the comm layer dispatches
	// handlers on detached goroutines, so one may outlive the old reader).
	m.resetHandShake()

	if err := m.executorClient.Connect(ctx, "Manager"); err != nil {
		return fmt.Errorf("executor connect failed: %w", err)
	}
	if err := m.waitForExecutorHandshake(); err != nil {
		// Tear down the half-open connection so the next tick can re-dial. Without
		// this, the client stays connected and every future Connect fails with
		// "already connected", wedging the drain permanently (e.g. a frozen or
		// half-open executor that never closes its end).
		if cerr := m.executorClient.Close(); cerr != nil {
			m.log.Warn("Manager: error closing executor after failed reconnect handshake: %v", cerr)
		}
		return fmt.Errorf("executor handshake failed: %w", err)
	}

	m.identityMu.Lock()
	m.draining = false
	m.identityMu.Unlock()
	m.log.Info("Manager: reconnected to executor and completed handshake")
	return nil
}

// verifySignerContinuity compares the executor's handshake-reported signing
// address against the on-chain teeSigner, once per handshake. It returns
// skip=true when the caller must not dispatch this tick (a transient teeSigner
// read error — consistent with the activeImage read), and a non-nil error only
// for a confirmed mismatch on a set teeSigner, which is fatal.
//
// It is deferred (skip=false, no error — dispatch proceeds) while teeSigner is
// unset (the legitimate bootstrap window before the first updateTee) or the
// executor reported no signing address: neither proves a regeneration, and the
// first state update reverts on-chain if the signer is actually wrong.
func (m *SecureProcessorManager) verifySignerContinuity(ctx context.Context) (skip bool, fatal error) {
	m.identityMu.RLock()
	already := m.signerVerified
	reported := m.reportedSigner
	m.identityMu.RUnlock()
	if already {
		return false, nil
	}

	onchain, err := m.blockchainClient.GetTeeSigner(ctx)
	if err != nil {
		// Transient read error: retry next tick rather than dispatch unverified
		// (matches the activeImage read above).
		m.log.Warn("Manager: could not read on-chain teeSigner (retry next tick): %v", err)
		return true, nil
	}

	if onchain == (ethCommon.Address{}) {
		// teeSigner not yet set (pre-updateTee bootstrap) — nothing to compare.
		m.log.Info("Manager: on-chain teeSigner not set yet, deferring signer-continuity check")
		return false, nil
	}

	if reported == "" {
		m.log.Warn("Manager: executor reported no signing address; deferring signer-continuity check")
		return false, nil
	}

	if !ethCommon.IsHexAddress(reported) || ethCommon.HexToAddress(reported) != onchain {
		return false, &SignerContinuityError{Reported: reported, OnChain: onchain.Hex()}
	}

	m.identityMu.Lock()
	m.signerVerified = true
	m.identityMu.Unlock()
	m.log.Info("Manager: signer-continuity check passed (signer %s matches on-chain teeSigner)", onchain.Hex())
	return false, nil
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

	// Reconcile the connected enclave against the on-chain activeImage (TEE swap
	// observation, drain, reconnect, and signer-continuity check) before taking
	// m.mu. reconnectExecutor may block on the handshake, which must not be done
	// under m.mu.
	if skip, err := m.maintainExecutorForActiveImage(ctx); err != nil {
		return err // fatal (signer-continuity mismatch)
	} else if skip {
		return nil // draining / reconnecting / transient read — no dispatch this tick
	}

	m.mu.Lock()
	defer m.mu.Unlock()

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

	if request == nil { // No request in queue, nothing to process
		return nil
	}

	appID := request.ApplicationID

	localStateRoot, err := m.dataLayer.LastVersionID(appID)
	if err != nil {
		if dbErr, ok := err.(*storageErrors.Error); ok && dbErr.Code == storageErrors.NoVersionInDb {
			localStateRoot = make([]byte, 32) // Initialize to zero state root if no version exists
		} else {
			m.log.Error("Manager: Failed to get local state root for app %d: %v", appID, err)
			return nil
		}
	}

	if !bytes.Equal(localStateRoot, stateRoot[:]) {
		m.log.Warn("Manager: State root mismatch for app %d, expected %x, got %x. Checking if it is a REORG.", appID, localStateRoot, stateRoot)

		isReorg, err := m.checkIfReorg(appID, stateRoot)
		if err != nil {
			m.log.Error("Manager: Failed to check for reorg: %v", err)
			return nil
		}

		if isReorg {
			if m.endReorgTime.IsZero() {
				m.log.Info("Manager: Starting REORG timeout %d for app %d", m.config.ReorgTimeout, appID)
				m.endReorgTime = time.Now().Add(time.Duration(m.config.ReorgTimeout) * time.Second)
				return nil
			}
			if time.Now().Before(m.endReorgTime) {
				m.log.Info("Manager: REORG timeout not expired yet for app %d. Keep waiting...", appID)
				return nil
			}
			m.log.Info("Manager: REORG not solved within timeout for app %d => Rollback the DB", appID)
			if err := m.dataLayer.Rollback(appID, stateRoot[:]); err != nil {
				m.log.Error("Manager: Error while rolling back the DB: %v", err)
				return fmt.Errorf("fatal: rollback failed: %w", err)
			}

		} else {
			m.log.Error("Manager: unrecoverable disalignment between DB and chain for app %d, no matching state root found in db", appID)
			emptyStateRoot := [32]byte{}
			if bytes.Equal(localStateRoot, emptyStateRoot[:]) {
				m.log.Error("Manager: the DB is empty but the chain state root is non-zero, check if the database file is correct and restart the manager")
			}
			return fmt.Errorf("unrecoverable disalignment between DB and chain, no matching state root found in db")
		}

	}

	if !m.endReorgTime.IsZero() {
		m.log.Info("Manager: State roots match for app %d, REORG resolved", appID)
		m.endReorgTime = time.Time{}
	}

	m.log.Info("Manager: processing request %s", request.RequestID)

	if err := m.processRequest(ctx, request); err != nil {
		m.log.Error("Manager: Failed to process request %s: %v", request.RequestID, err)
	} else {
		m.log.Info("Manager: Processed request %s", request.RequestID)
	}
	return nil

}

func (m *SecureProcessorManager) checkIfReorg(appID common.ApplicationIdType, stateRoot [32]byte) (bool, error) {
	if stateRoot == [32]byte{} {
		m.log.Info("Manager: State root is zero for app %d, REORG", appID)
		// Don't look for older db versions, just mark as reorged and wait for next poll
		return true, nil
	}

	oldVersions, err := m.dataLayer.ListVersions(appID)
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
func (m *SecureProcessorManager) processRequest(ctx context.Context, req *common.Request) error {
	if !m.isRunning {
		m.log.Warn("Manager is not started yet, skipping")
		return errors.New("manager is not running")
	}

	switch req.RequestType {
	case common.Deploy:
		return m.processDeployApp(ctx, req)
	default: //Now it is the executor that has to check if a requestType is correct
		return m.processProcessRequest(ctx, req)
	}
}

func (m *SecureProcessorManager) submitStateOnChain(ctx context.Context, updatePayload *common.UpdatePayload) error {
	// Submit the state update to the blockchain
	if err := m.blockchainClient.SubmitStateUpdate(ctx, updatePayload); err != nil {
		m.log.Error("Failed to submit state update for error: %v", err)
		if updatePayload.ErrorCode == 0 {
			m.log.Info("Rollback the application state to previous version")
			if rollbackErr := m.dataLayer.Rollback(updatePayload.ApplicationID, updatePayload.PrevStateRoot[:]); rollbackErr != nil {
				// If this happens, the local db and the chain are out of sync and cannot be recovered automatically.
				// Log and return err to let REORG detection handle it on the next poll.
				m.log.Error("Failed to rollback application state: %v. Will retry via REORG detection.", rollbackErr)
				return rollbackErr
			}

			if _, ok := err.(blockchain.ReorgError); ok {
				m.log.Warn("REORG, wait for next poll")
				return nil
			}
		}
		return err
	}

	m.log.Info("Manager: Processed request %s for application %d", updatePayload.RequestID, updatePayload.ApplicationID)
	return nil
}

// processDeployApp processes a deploy app request
func (m *SecureProcessorManager) processDeployApp(ctx context.Context, req *common.Request) error {
	m.log.Info("Processing deploy app request: %s", req.RequestID)
	if !m.isRunning {
		m.log.Warn("Manager is not started yet, skipping")
		return errors.New("manager is not running")
	}

	// Check if app was already deployed. With per-app state roots, the executor handles
	// appState == nil gracefully (uses emptyStateRoot for signed error payloads).
	state, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err != nil {
		if !storageErrors.IsNotFound(err) {
			m.log.Warn("Manager: Got error while getting application state: %v", err)
			// Even if it is not the expected error in case of a deploy, continue for backwards compatibility
		}
	}

	wasmModule, resolveErr := m.resolveDeployWASM(ctx, req)
	if resolveErr != nil {
		if dErr, ok := resolveErr.(*deployResolutionError); ok {
			m.log.Warn(
				"Manager: continuing deploy with unresolved artifact requestId=%s applicationId=%d kind=%s err=%v",
				req.RequestID,
				req.ApplicationID,
				dErr.kind,
				resolveErr,
			)
		} else {
			m.log.Warn(
				"Manager: continuing deploy with unresolved artifact requestId=%s applicationId=%d err=%v",
				req.RequestID,
				req.ApplicationID,
				resolveErr,
			)
		}
		wasmModule = nil
	}

	// Deploy the application
	updatePayload, appState, err := m.executorClient.SendDeployApp(ctx, req, state, wasmModule)
	if err != nil {
		return err
	}

	// Check if this is a signed error payload
	if updatePayload.ErrorCode != 0 {
		m.log.Info("Manager: Received signed error payload for deploy request %s (error code: %d)", req.RequestID, updatePayload.ErrorCode)
		// Submit the signed error payload to the blockchain (no state to store)
		return m.submitStateOnChain(ctx, updatePayload)
	}

	// Store the application state and WASM bytecode
	err = m.initAppStorage(ctx, appState, wasmModule)
	if err != nil {
		m.log.Error("failed to submit state update: %v", err)
		return err
	}

	m.log.Info("Deployed application, submit the state update to the blockchain")
	return m.submitStateOnChain(ctx, updatePayload)

}

// processProcessRequest processes a process request (including deanonymization requests)
func (m *SecureProcessorManager) processProcessRequest(ctx context.Context, req *common.Request) error {
	m.log.Info("Processing Process app request: %s (type: %s)", req.RequestID, req.RequestType)
	if !m.isRunning {
		m.log.Warn("Manager is not started yet, skipping")
		return errors.New("manager is not running")
	}

	// Get the application state
	//TODO ST manager shouldn't be able return a tampered state. Same for wasm
	appState, err := m.dataLayer.GetApplicationState(ctx, req.ApplicationID)
	if err != nil {
		m.log.Error("GetApplicationState returns an error: %v", err)
		if !storageErrors.IsNotFound(err) {
			// Other errors are likely db errors, retry on next poll
			return err
		}
		// if application state not found, pass a nil state to the executor to let it handle the error
		appState = nil
	}

	// Get the WASM module for the application
	var wasmBytes []byte
	if appState != nil {
		wasmBytes, err = m.dataLayer.GetWASMBytecode(ctx, req.ApplicationID)
		if err != nil {
			m.log.Error("GetWASMBytecode returns an error: %v", err)
			return err
		}
	}
	// Process the request
	updatePayload, updatedState, deanonymizationReport, err := m.executorClient.SendProcessRequest(ctx, req, appState, wasmBytes)
	if err != nil {
		// Fallback case: executor couldn't create signed payload. Retry on next poll, no state to store
		return err
	}

	// Check if this is a signed error payload
	if updatePayload.ErrorCode != 0 {
		m.log.Info("Manager: Received signed error payload for request %s (error code: %d)", req.RequestID, updatePayload.ErrorCode)
		// Submit the signed error payload to the blockchain (no state to store)
		return m.submitStateOnChain(ctx, updatePayload)
	}

	// If a deanonymization report was generated, save it to filesystem
	// Note: The executor already validates that reports are only and always generated for Deanonymize requests
	if deanonymizationReport != nil {
		if err := m.saveDeanonymizationReport(deanonymizationReport, req); err != nil {
			// Treat as transient: log and retry on next poll instead of failing the request.
			m.log.Error("Failed to save deanonymization report: %v", err)
			return err
		}
		m.log.Info("Manager: Saved deanonymization report %s for application %d", deanonymizationReport.ReportID, req.ApplicationID)
	}

	// Store the updated application state
	err = m.storeStateToStorage(ctx, updatedState)
	if err != nil {
		m.log.Error("failed to submit state update: %v", err)
		return err
	}

	return m.submitStateOnChain(ctx, updatePayload)
}

// initAppStorage stores the application state and WASM bytecode for a newly deployed application.
func (m *SecureProcessorManager) initAppStorage(ctx context.Context, state *common.ApplicationState, wasmModule []byte) error {
	versionID := state.StateRoot[:]
	m.log.Info("Storing app state and WASM, versionID %x", versionID)
	return m.dataLayer.StoreWithWasm(
		ctx,
		versionID,
		state,
		&common.WASMData{ApplicationID: state.ApplicationID, Bytecode: wasmModule},
	)
}

// storeStateToStorage stores the updated application state after processing a request.
func (m *SecureProcessorManager) storeStateToStorage(ctx context.Context, state *common.ApplicationState) error {
	versionID := state.StateRoot[:]
	m.log.Info("Storing app state, versionID %x", versionID)
	return m.dataLayer.Store(ctx, versionID, state)
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
