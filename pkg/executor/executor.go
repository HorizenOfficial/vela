package executor

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	dbtypes "github.com/horizen-pes/pkg/common/dbstate"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/crypto"
)

// StatelessExecutor implements the Executor interface
type StatelessExecutor struct {
	config  *Config
	runtime Runtime
	server  communication.ExecutorServer
	*MsgToSignBuilder
}

// NewStatelessExecutor creates a new stateless executor
func NewStatelessExecutor(config *Config, runtime Runtime, server communication.ExecutorServer) (*StatelessExecutor, error) {
	msgBuilder, err := NewMsgToSignBuilder()
	if err != nil {
		return nil, fmt.Errorf("failed to setup msg to sign builder: %w", err)
	}

	executor := &StatelessExecutor{
		config:           config,
		runtime:          runtime,
		server:           server,
		MsgToSignBuilder: msgBuilder,
	}
	// Set this executor as the request handler
	executor.server.SetRequestHandler(executor)

	return executor, nil
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
func (e *StatelessExecutor) HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	log.Printf("Executor: Processing request %s for application %s", req.RequestID, req.ApplicationID)

	// Decrypt the encrypted state
	decryptedState, err := DecryptState(appState.EncryptedState, e.config.StateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt state: %w", err)
	}

	// Verify state consistency
	stateRoot := sha256.Sum256(decryptedState)
	if stateRoot != appState.StateRoot {
		return nil, nil, fmt.Errorf("state root mismatch: got %x, want %x", stateRoot, appState.StateRoot)
	}

	//parse the state
	dbState, err := dbtypes.DeserializeDBState(decryptedState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize state: %w", err)
	}

	// Decrypt the request payload
	decryptedPayload, err := e.decryptPayload(e.config.CommunicationKey, req.Payload, req.Sender, dbState.GetKeyStore())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt request payload: %w", err)
	}

	// If the request contains a deposit, handle it first
	var tempState = dbState.GetAppState()
	var depositEvents []common.PlainEvent
	if req.Value > 0 {
		tempState, depositEvents, err = e.runtime.Deposit(ctx, req.ApplicationID, req.Sender, req.Value, tempState, wasmModule)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to process deposit in WASM runtime: %w", err)
		}
		log.Printf("Executor: Successfully processed deposit for request %s", req.RequestID)
	}

	// Invoke WASM method to process the request
	newWasmState, events, withdrawals, err := e.runtime.ProcessRequest(ctx, req.ApplicationID, req.Sender, decryptedPayload, tempState, wasmModule)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process request in WASM runtime: %w", err)
	}
	log.Printf("Executor: Successfully processed request %s", req.RequestID)

	//recreate new dbstate
	dbState.SetAppState(newWasmState)
	newDbState, err := dbState.Serialize()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize new state: %w", err)
	}

	// Encrypt the new state and events
	encryptedNewState, err := crypto.EncryptWithAES(e.config.StateKey, newDbState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt new application state: %w", err)
	}
	// Encrypt events if they are not empty
	events = append(depositEvents, events...)
	encryptedEvents, err := e.encryptEvents(ctx, events, req.ApplicationID, e.config.CommunicationKey, e.server, dbState.GetKeyStore())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt events: %w", err)
	}
	log.Printf("Executor: Successfully encrypted new application state")

	// Create state root hash
	newStateRoot := sha256.Sum256(newDbState)

	// Create the update payload
	updatePayload := &common.UpdatePayload{
		ApplicationID: req.ApplicationID,
		RequestID:     req.RequestID,
		PrevStateRoot: appState.StateRoot,
		NewStateRoot:  newStateRoot,
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
	newApplicationState := &common.ApplicationState{
		ApplicationID:  req.ApplicationID,
		StateRoot:      newStateRoot,
		EncryptedState: encryptedNewState,
	}

	log.Printf("Executor: Successfully processed request %s", req.RequestID)
	return updatePayload, newApplicationState, nil
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
		PrevStateRoot: [32]byte{}, // No previous state root for new applications
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
func (e *StatelessExecutor) HandleGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error) {
	log.Printf("Executor: Generating deanonymization report for request %s", req.RequestID)

	// Decrypt the encrypted state
	decryptedState, err := DecryptState(appState.EncryptedState, e.config.StateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state: %w", err)
	}

	// Verify state consistency
	hash := sha256.Sum256(decryptedState)
	if hash != appState.StateRoot {
		return nil, fmt.Errorf("state root mismatch: got %x, want %x", hash, appState.StateRoot)
	}

	// Parse the state
	dbState, err := dbtypes.DeserializeDBState(decryptedState)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize state: %w", err)
	}

	// Decrypt the request payload
	decryptedPayload, err := e.decryptPayload(e.config.CommunicationKey, req.Payload, req.Sender, dbState.GetKeyStore())
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt Payload: %w", err)
	}

	// Generate the report using the runtime
	reportData, err := e.runtime.GenerateDeanonymizationReport(ctx, req.ApplicationID, req.RequestID, decryptedPayload, dbState.GetAppState(), wasmModule)
	if err != nil {
		return nil, fmt.Errorf("failed to generate deanonymization report: %w", err)
	}

	// Encrypt the report
	encryptedReport, err := e.encryptDeanonymizationReport(e.config.CommunicationKey, req.Sender, reportData, dbState.GetKeyStore())
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
	hash, err := e.BuildMsgHash(payload)

	if err != nil {
		return nil, fmt.Errorf("failed to create message to sign: %w", err)
	}

	// Sign the hash
	signature, err := e.config.SignatureKey.Sign(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign update payload: %w", err)
	}

	// The signature should be 65 bytes long (R + S + V)
	if len(signature) != 65 {
		return nil, fmt.Errorf("signature has wrong length: got %d, want 65", len(signature))
	}

	// Adjust V value to be 27 or 28, because ecrecover in Solidity expects this format
	if signature[64] < 27 {
		signature[64] += 27
	}
	return signature, nil
}

func (e *StatelessExecutor) encryptEvents(ctx context.Context, events []common.PlainEvent, appId string, key *cryptotypes.PrivateKeyP521, server communication.ExecutorServer, keyStore dbtypes.KeyStore) ([]common.Event, error) {
	if len(events) == 0 {
		return nil, nil // No events to encrypt
	}
	encryptedEvents := make([]common.Event, len(events))

	// Request user public keys from the server
	var users []string
	for _, event := range events {
		if event.UserID != "" {
			users = append(users, event.UserID)
		}
	}

	for i, event := range events {
		// retrieve user Secp521r1_PubKey
		userKey, exists := keyStore[event.UserID]
		if !exists {
			return nil, fmt.Errorf("no Secp521r1_PubKey found for user %s", event.UserID)
		}
		// Encrypt the event data
		encryptedData, err := crypto.Encrypt(key, userKey, event.Data)
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

func (e *StatelessExecutor) encryptDeanonymizationReport(key *cryptotypes.PrivateKeyP521, requester string, reportData []byte, keyStore dbtypes.KeyStore) ([]byte, error) {
	if len(reportData) == 0 {
		return reportData, nil // No report to encrypt
	}

	// retrieve user Secp521r1_PubKey
	requesterPublicKey, exists := keyStore[requester]
	if !exists {
		return nil, fmt.Errorf("no Secp521r1_PubKey found for user %s", requester)
	}

	// Encrypt the report data
	encryptedReport, err := crypto.Encrypt(key, requesterPublicKey, reportData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt deanonymization report: %w", err)
	}

	return encryptedReport, nil
}

func DecryptState(encryptedState []byte, decryptionKey cryptotypes.AES256Key) ([]byte, error) {
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

func (e *StatelessExecutor) decryptPayload(decryptionKey *cryptotypes.PrivateKeyP521, payload []byte, sender string, keyStore dbtypes.KeyStore) ([]byte, error) {
	if len(payload) == 0 {
		return payload, nil // No payload to decrypt
	}

	// retrieve sender Secp521r1_PubKey
	userKey, exists := keyStore[sender]
	if !exists {
		return nil, fmt.Errorf("no Secp521r1_PubKey found for sender %s", sender)
	}

	decryptedPayload, err := crypto.Decrypt(userKey, decryptionKey, payload)
	if err != nil {
		return nil, fmt.Errorf("payload decryption failed: %w", err)
	}

	log.Printf("Executor: Successfully decrypted request payload")
	return decryptedPayload, nil
}
