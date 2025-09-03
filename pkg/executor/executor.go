package executor

import (
	"bytes"
	"context"
	"crypto/ecdh"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/crypto"
)

// StatelessExecutor implements the Executor interface
type StatelessExecutor struct {
	config  *Config
	runtime Runtime
	server  communication.ExecutorServer
}

// NewStatelessExecutor creates a new stateless executor
func NewStatelessExecutor(config *Config, runtime Runtime, server communication.ExecutorServer) *StatelessExecutor {
	executor := &StatelessExecutor{
		config:  config,
		runtime: runtime,
		server:  server,
	}
	// Set this executor as the request handler
	executor.server.SetRequestHandler(executor)

	return executor
}

// Start starts the executor server
func (e *StatelessExecutor) Start(ctx context.Context) error {
	switch e.config.ServerType {
	case "tcp":
		log.Printf("Executor: Starting TCP executor server on %s", e.config.ServerAddr)
	case "v-sock":
		log.Printf("Executor: Starting v-socket executor server on CID %d, Port %d", e.config.ServerCid, e.config.ServerPort)
	}
	return e.server.Start(ctx)
}

// Stop stops the executor server
func (e *StatelessExecutor) Stop() error {
	log.Printf("Executor: Stopping stateless executor")
	return e.server.Stop()
}

// Close closes the executor and its runtime
func (e *StatelessExecutor) Close() error {
	if err := e.Stop(); err != nil {
		log.Printf("Executor: Error stopping server: %v", err)
	}
	if e.runtime != nil {
		return e.runtime.Close()
	}
	return nil
}

// HandleProcessRequest is the main workflow: decrypt state -> invoke WASM -> encrypt new state -> sign -> respond
func (e *StatelessExecutor) HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	log.Printf("Executor: Processing request %s for application %s", req.RequestID, req.ApplicationID)

	// Decrypt the encrypted state
	decryptedState, err := DecryptState(appState.EncryptedState, e.config.StateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt state: %w", err)
	}

	// Verify state consistency
	if sha256.Sum256(decryptedState) != [32]byte(appState.StateRoot) {
		return nil, nil, fmt.Errorf("state root mismatch: got %x, want %x", decryptedState, appState.StateRoot)
	}

	// Decrypt the request payload
	decryptedPayload, err := DecryptPayload(e.config.CommunicationKey, req.Payload, senderKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt request payload: %w", err)
	}

	// If the request contains a deposit, handle it first
	var tempState = decryptedState
	var depositEvents []PlainEvent
	if req.Value > 0 {
		tempState, depositEvents, err = e.runtime.Deposit(ctx, req.ApplicationID, req.Sender, req.Value, decryptedState, wasmModule)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to process deposit in WASM runtime: %w", err)
		}
		log.Printf("Executor: Successfully processed deposit for request %s", req.RequestID)
	}

	// Invoke WASM method to process the request
	newState, events, withdrawals, err := e.runtime.ProcessRequest(ctx, req.ApplicationID, req.Sender, decryptedPayload, tempState, wasmModule)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process request in WASM runtime: %w", err)
	}
	log.Printf("Executor: Successfully processed request %s", req.RequestID)

	// Encrypt the new state and events
	encryptedNewState, err := crypto.EncryptWithAES(e.config.StateKey, newState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt new application state: %w", err)
	}
	// Encrypt events if they are not empty
	events = append(depositEvents, events...)
	encryptedEvents, err := EncryptEvents(ctx, events, req.ApplicationID, e.config.CommunicationKey, e.server)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt events: %w", err)
	}
	log.Printf("Executor: Successfully encrypted new application state")

	// Create state root hash
	newStateRoot := sha256.Sum256(newState)

	// Create the update payload
	updatePayload := &common.UpdatePayload{
		ApplicationID: req.ApplicationID,
		RequestID:     req.RequestID,
		PrevStateRoot: appState.StateRoot,
		NewStateRoot:  newStateRoot[:],
		Events:        encryptedEvents,
		Withdrawals:   withdrawals,
	}

	// Sign the update payload (produce attestation)
	signature, err := e.signUpdatePayload(updatePayload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign update payload: %w", err)
	}
	updatePayload.Signature = signature

	// Create the new application state
	newAppState := &common.ApplicationState{
		ApplicationID:  req.ApplicationID,
		StateRoot:      newStateRoot[:],
		EncryptedState: encryptedNewState,
	}

	log.Printf("Executor: Successfully processed request %s", req.RequestID)
	return updatePayload, newAppState, nil
}

// HandleDeployApp implements the RequestHandler interface
func (e *StatelessExecutor) HandleDeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, error) {
	log.Printf("Executor: Deploying application for request %s", req.RequestID)

	// For deployment, we need to initialize the application with the WASM module
	// The payload should contain the WASM bytecode
	wasmModule := req.Payload

	// Load the module and get initial state
	initialState, stateRoot, err := e.runtime.LoadModule(ctx, req.ApplicationID, wasmModule)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load WASM module: %w", err)
	}

	// Encrypt the initial state
	encryptedState, err := crypto.EncryptWithAES(e.config.StateKey, initialState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt initial application state: %w", err)
	}
	log.Printf("Executor: Successfully encrypted initial application state")

	// Create the application state
	appState := &common.ApplicationState{
		ApplicationID:  req.ApplicationID,
		StateRoot:      stateRoot,
		EncryptedState: encryptedState,
	}

	// Create the update payload
	updatePayload := &common.UpdatePayload{
		ApplicationID: req.ApplicationID,
		RequestID:     req.RequestID,
		PrevStateRoot: nil, // No previous state root for new applications
		NewStateRoot:  stateRoot,
	}

	// Sign the update payload (produce attestation)
	signature, err := e.signUpdatePayload(updatePayload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign update payload: %w", err)
	}
	updatePayload.Signature = signature

	log.Printf("Executor: Successfully deployed application %s", req.ApplicationID)
	return updatePayload, appState, nil
}

// HandleGenerateDeanonymizationReport implements the RequestHandler interface
func (e *StatelessExecutor) HandleGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.DeanonymizationReport, error) {
	log.Printf("Executor: Generating deanonymization report for request %s", req.RequestID)

	// Decrypt the encrypted state
	decryptedState, err := DecryptState(appState.EncryptedState, e.config.StateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state: %w", err)
	}

	// Verify state consistency
	hash := sha256.Sum256(decryptedState)
	if !bytes.Equal(hash[:], appState.StateRoot) {
		return nil, fmt.Errorf("state root mismatch: got %x, want %x", decryptedState, appState.StateRoot)
	}

	// Decrypt the request payload
	decryptedPayload, err := DecryptPayload(e.config.CommunicationKey, req.Payload, senderKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt Payload: %w", err)
	}

	// Generate the report using the runtime
	reportData, err := e.runtime.GenerateDeanonymizationReport(ctx, req.ApplicationID, req.RequestID, decryptedPayload, decryptedState, wasmModule)
	if err != nil {
		return nil, fmt.Errorf("failed to generate deanonymization report: %w", err)
	}

	// Encrypt the report
	encryptedReport, err := EncryptDeanonymizationReport(e.config.CommunicationKey, senderKey, reportData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt deanonymization report: %w", err)
	}
	log.Printf("Executor: Successfully encrypted deanonymization report")

	report := &common.DeanonymizationReport{
		ApplicationID:   req.ApplicationID,
		ReportID:        req.RequestID,
		EncryptedReport: encryptedReport,
	}

	log.Printf("Executor: Successfully generated deanonymization report %s", req.RequestID)
	return report, nil
}

// signUpdatePayload signs the update payload to produce an attestation
func (e *StatelessExecutor) signUpdatePayload(payload *common.UpdatePayload) ([]byte, error) {
	// Serialize the payload for signing
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload for signing: %w", err)
	}

	// Create a hash of the payload
	hash := ethCrypto.Keccak256(payloadBytes)

	// Sign the hash
	signature, err := e.config.SignatureKey.Sign(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign update payload: %w", err)
	}

	// The signature should be 65 bytes long (R + S + V)
	if len(signature) != 65 {
		return nil, fmt.Errorf("signature has wrong length: got %d, want 65", len(signature))
	}

	return signature, nil
}

func EncryptEvents(ctx context.Context, events []PlainEvent, appId string, key *common.PrivateKeyP521, server communication.ExecutorServer) ([]common.Event, error) {
	if len(events) == 0 {
		return nil, nil // No events to encrypt
	}
	curve := ecdh.P521()
	encryptedEvents := make([]common.Event, len(events))

	// Request user public keys from the server
	var users []string
	for _, event := range events {
		if event.UserID != "" {
			users = append(users, event.UserID)
		}
	}
	keys, err := server.SendGetUserKeys(ctx, users)
	if err != nil {
		return nil, fmt.Errorf("failed to get user keys: %w", err)
	}

	for i, event := range events {
		// create a user public key
		userKeyBytes, exists := keys[event.UserID]
		if !exists {
			return nil, fmt.Errorf("no public key found for user %s", event.UserID)
		}
		pubkey, err := curve.NewPublicKey(userKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to create public key from user key: %w", err)
		}
		userKey := common.PublicKeyP521{PublicKey: pubkey}

		// Encrypt the event data
		encryptedData, err := crypto.Encrypt(key, &userKey, event.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt event data: %w", err)
		}

		// Create the encrypted event
		encryptedEvents[i] = common.Event{
			ApplicationID: appId,
			UserID:        event.UserID,
			EncryptedData: encryptedData,
		}
	}

	log.Printf("Executor: Successfully encrypted %d events", len(events))
	return encryptedEvents, nil
}

func EncryptDeanonymizationReport(key *common.PrivateKeyP521, senderKey []byte, reportData []byte) ([]byte, error) {
	if len(reportData) == 0 {
		return reportData, nil // No report to encrypt
	}

	curve := ecdh.P521()
	pubkey, err := curve.NewPublicKey(senderKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key from sender key: %w", err)
	}
	senderPublicKey := common.PublicKeyP521{PublicKey: pubkey}

	// Encrypt the report data
	encryptedReport, err := crypto.Encrypt(key, &senderPublicKey, reportData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt deanonymization report: %w", err)
	}

	return encryptedReport, nil
}

func DecryptState(encryptedState []byte, decryptionKey common.AES256Key) ([]byte, error) {
	if len(encryptedState) == 0 {
		return encryptedState, nil // nothing to decrypt
	}

	decryptedState, err := crypto.DecryptWithAES(decryptionKey, encryptedState)
	if err != nil {
		return nil, fmt.Errorf("state decryption failed: %w", err)
	}

	log.Printf("Executor: Successfully decrypted application state")
	return decryptedState, nil
}

func DecryptPayload(decryptionKey *common.PrivateKeyP521, payload []byte, senderKey []byte) ([]byte, error) {
	if len(payload) == 0 {
		return payload, nil // No payload to decrypt
	}

	curve := ecdh.P521()
	pubkey, err := curve.NewPublicKey(senderKey)
	key := common.PublicKeyP521{PublicKey: pubkey}
	if err != nil {
		return nil, fmt.Errorf("failed to create public key from sender key: %w", err)
	}

	decryptedPayload, err := crypto.Decrypt(&key, decryptionKey, payload)
	if err != nil {
		return nil, fmt.Errorf("payload decryption failed: %w", err)
	}

	log.Printf("Executor: Successfully decrypted request payload")
	return decryptedPayload, nil
}
