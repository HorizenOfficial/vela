package executor

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/appdata"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
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

	// Decrypte and parse the app data
	appData, err := e.fromEncryptedStateToAppData(appState)
	if err != nil {
		return nil, nil, err
	}

	// If the request contains a deposit, handle it first
	var tempState = appData.GetAppState()
	var depositEvents []common.PlainEvent
	if req.Value > 0 {
		tempState, depositEvents, err = e.runtime.Deposit(ctx, req.ApplicationID, req.Sender, req.Value, tempState, wasmModule)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to process deposit in WASM runtime: %w", err)
		}
		log.Printf("Executor: Successfully processed deposit for request %s", req.RequestID)
	}

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal
	if req.RequestType == common.AssociateKey {
		//request  of type associate key: the payload is not encrypted and contains the new key
		log.Printf("Associating new key - RequestID %s", req.RequestID)

		keyToAssociate, err := cryptotypes.NewPublicKeyP521(req.Payload)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse keyP521 in request payload: %w", err)
		}

		appData.AddKey(ethCommon.HexToAddress(req.Sender), *keyToAssociate)

	} else {
		//any other case: decrypt the payload and forward to the WASM to obtain the new state

		// Decrypt the request payload
		senderAddress := ethCommon.HexToAddress(req.Sender)
		decryptedPayload, err := e.decryptPayload(e.config.CommunicationKey, req.Payload, senderAddress, appData.GetKeyStore())
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decrypt request payload: %w", err)
		}

		// Invoke WASM method to process the request
		var newWasmState []byte
		newWasmState, events, withdrawals, err = e.runtime.ProcessRequest(ctx, req.ApplicationID, req.Sender, decryptedPayload, tempState, wasmModule)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to process request in WASM runtime: %w", err)
		}
		log.Printf("Executor: Successfully processed request %s", req.RequestID)

		appData.SetAppState(newWasmState)
	}

	//increment appNonce
	appData.IncrementNonce()

	//serialize the new app data
	newAppData, err := appData.Serialize()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize new appData: %w", err)
	}

	// Encrypt the new app data and events
	encryptedNewAppData, err := crypto.EncryptWithAES(e.config.StateKey, newAppData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt new appData: %w", err)
	}
	// Encrypt events if they are not empty
	events = append(depositEvents, events...)
	encryptedEvents, err := e.encryptEvents(ctx, events, req.ApplicationID, e.config.CommunicationKey, e.server, appData.GetKeyStore())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt events: %w", err)
	}
	log.Printf("Executor: Successfully encrypted new application data")

	// Create appdata root hash
	newStateRoot := sha256.Sum256(newAppData)

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
		EncryptedState: encryptedNewAppData,
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
	initialAppState, _, err := e.runtime.LoadModule(ctx, req.ApplicationID, wasmModule)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load WASM module: %w", err)
	}

	initialAppData := appdata.NewAppData(initialAppState)

	//serialize the new app data
	initialAppDataBytes, err := initialAppData.Serialize()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize new app data: %w", err)
	}
	// Create app data root hash
	initialAppDataRoot := sha256.Sum256(initialAppDataBytes)

	// Encrypt the initial state
	encryptedState, err := crypto.EncryptWithAES(e.config.StateKey, initialAppDataBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt initial app data: %w", err)
	}
	log.Printf("Executor: Successfully encrypted initial app data")

	// Create the application state
	appState := &common.ApplicationState{
		ApplicationID:  req.ApplicationID,
		StateRoot:      initialAppDataRoot,
		EncryptedState: encryptedState,
	}

	// Create the update payload
	updatePayload := &common.UpdatePayload{
		ApplicationID: req.ApplicationID,
		RequestID:     req.RequestID,
		PrevStateRoot: [32]byte{}, // No previous state root for new applications
		NewStateRoot:  initialAppDataRoot,
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

	// Decrypte and parse the app data
	appData, err := e.fromEncryptedStateToAppData(appState)
	if err != nil {
		return nil, err
	}

	// Decrypt the request payload
	senderAddress := ethCommon.HexToAddress(req.Sender)
	decryptedPayload, err := e.decryptPayload(e.config.CommunicationKey, req.Payload, senderAddress, appData.GetKeyStore())
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt Payload: %w", err)
	}

	// Generate the report using the runtime
	reportData, err := e.runtime.GenerateDeanonymizationReport(ctx, req.ApplicationID, req.RequestID, decryptedPayload, appData.GetAppState(), wasmModule)
	if err != nil {
		return nil, fmt.Errorf("failed to generate deanonymization report: %w", err)
	}

	// Encrypt the report
	encryptedReport, err := e.encryptDeanonymizationReport(e.config.CommunicationKey, senderAddress, reportData, appData.GetKeyStore())
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

func (e *StatelessExecutor) fromEncryptedStateToAppData(encState *common.ApplicationState) (*appdata.AppData, error) {
	// Decrypt the encrypted state
	decryptedState, err := DecryptState(encState.EncryptedState, e.config.StateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt state: %w", err)
	}

	// Verify state consistency
	hash := sha256.Sum256(decryptedState)
	if hash != encState.StateRoot {
		return nil, fmt.Errorf("state root mismatch: got %x, want %x", hash, encState.StateRoot)
	}

	// Parse the app data
	appData, err := appdata.DeserializeAppData(decryptedState)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize state: %w", err)
	}

	return appData, nil
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

func (e *StatelessExecutor) encryptEvents(ctx context.Context, events []common.PlainEvent, appId string, key *cryptotypes.PrivateKeyP521, server communication.ExecutorServer, keyStore appdata.KeyStore) ([]common.Event, error) {
	if len(events) == 0 {
		return nil, nil // No events to encrypt
	}
	encryptedEvents := make([]common.Event, len(events))

	for i, event := range events {
		// retrieve user Secp521r1_PubKey
		userKey, exists := keyStore[ethCommon.HexToAddress(event.UserID)]
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

func (e *StatelessExecutor) encryptDeanonymizationReport(key *cryptotypes.PrivateKeyP521, requester ethCommon.Address, reportData []byte, keyStore appdata.KeyStore) ([]byte, error) {
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

func (e *StatelessExecutor) decryptPayload(decryptionKey *cryptotypes.PrivateKeyP521, payload []byte, sender ethCommon.Address, keyStore appdata.KeyStore) ([]byte, error) {
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
