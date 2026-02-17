package executor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/hf/nsm"
	"github.com/hf/nsm/request"
	"github.com/hf/nsm/response"

	"github.com/horizen-pes/pkg/admin"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/appdata"
	"github.com/horizen-pes/pkg/common/apperrors"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/crypto"
	"github.com/horizen-pes/pkg/logger"
)

// NsmSession is an interface abstracting nsm.Session for testability.
type NsmSession interface {
	Send(req request.Request) (response.Response, error)
	Close() error
}

// As of now we support only one app having this ID
var (
	admittedAppID = common.NewApplicationId(1)
	admittedProtocolVersion = uint8(0)
	emptyStateRoot = [32]byte{}
)

func CreateNewKeySet() (*EnclaveKeySet, error) {
	communicationKey, err := crypto.GeneratePrivateKeyP521()
	if err != nil {
		return nil, fmt.Errorf("failed to generate communication key: %w", err)
	}
	signingKey, err := crypto.GeneratePrivateKeySecp256k1()
	if err != nil {
		return nil, fmt.Errorf("failed to generate signing key: %w", err)
	}
	stateKey, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate shared key: %w", err)
	}
	keySet := &EnclaveKeySet{
		CommunicationKey: *communicationKey,
		SigningKey:       *signingKey,
		StateKey:         stateKey,
	}

	return keySet, nil
}

func (e *StatelessExecutor) DumpPublicKeys() {
	if e.keySet == nil {
		e.log.Info("Executor: nothing to print, keyset is null")
		return
	}

	keyP521StrPub := crypto.ExportPublicKeyP521ToHex(e.keySet.CommunicationKey.PublicKey())
	keySecp256k1StrPub := crypto.ExportPublicKeySecp256k1ToHex(e.keySet.SigningKey.PublicKey())
	keySecp256k1StrAddress := e.keySet.SigningKey.PublicKey().Address()
	e.log.Info("###: Communication key P521 (public): 0x" + keyP521StrPub)
	e.log.Info("###: Signing key Secp256k1 (public):  0x" + keySecp256k1StrPub)
	e.log.Info("###:             Secp256k1 (address): 0x" + keySecp256k1StrAddress)
}

func (e *StatelessExecutor) DumpPrivateKeys() {
	if e.keySet == nil {
		e.log.Info("Executor: nothing to print, keyset is null")
		return
	}

	keyP521StrPriv := crypto.ExportPrivateKeyP521ToHex(&e.keySet.CommunicationKey)
	keySecp256k1StrPriv := crypto.ExportPrivateKeySecp256k1ToHex(&e.keySet.SigningKey)
	e.log.Info("###: Communication key P521 (private): 0x" + keyP521StrPriv)
	e.log.Info("###: Signing key Secp256k1 (private):  0x" + keySecp256k1StrPriv)
	e.log.Info("###: State key AES256 (raw): %x", e.keySet.StateKey)
}

// GenerateEnclaveKeySet creates a new keyset from scratch.
// It will be used at the first initialization.
// It returns the new keyset and the recovery structure.
func GenerateEnclaveKeySet(recoveryType int) (*EnclaveKeySet, *common.EnclaveKeySetRecovery, error) {
	switch recoveryType {
	case 0:
		// 1. Create a new AES master key.
		masterKey, err := crypto.GenerateAESKey()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate master key: %w", err)
		}

		// 2. Create a new EnclaveKeySet.
		keySet, err := CreateNewKeySet()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create key set: %w", err)
		}

		// 3. Serialize the EnclaveKeySet.
		serializedKeySet, err := keySet.Serialize()
		if err != nil {
			return nil, nil, err
		}

		// 4. Encrypt the serialized EnclaveKeySet with the master key.
		encryptedKeySet, err := crypto.EncryptWithAES(masterKey, serializedKeySet)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt key set: %w", err)
		}

		// 5. Create the recovery structure.
		recovery := &common.EnclaveKeySetRecovery{
			RecoveryType:       0,
			KeySetCiphertext:   encryptedKeySet,
			RecoveryCiphertext: masterKey[:],
		}

		return keySet, recovery, nil
	default:
		return nil, nil, fmt.Errorf("unsupported recovery type: %d", recoveryType)
	}
}

// RestoreEnclaveKeySet recovers a keyset from a recovery previously stored.
// It will be used when starting the executors after the first init.
func RestoreEnclaveKeySet(recovery *common.EnclaveKeySetRecovery) (*EnclaveKeySet, error) {
	switch recovery.RecoveryType {
	case 0:
		// 1. Use the RecoveryCiphertext as the master key to decrypt the KeySetCiphertext.
		var masterKey cryptotypes.AES256Key
		copy(masterKey[:], recovery.RecoveryCiphertext)

		decryptedKeySet, err := crypto.DecryptWithAES(masterKey, recovery.KeySetCiphertext)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt key set: %w", err)
		}

		// 2. Deserialize the decrypted data into an EnclaveKeySet.
		keySet, err := DeserializeEnclaveKeySet(decryptedKeySet)
		if err != nil {
			return nil, err
		}

		return keySet, nil
	default:
		return nil, fmt.Errorf("unsupported recovery type: %d", recovery.RecoveryType)
	}
}

// StatelessExecutor implements the Executor interface
type StatelessExecutor struct {
	config       *Config
	runtime      Runtime
	server       communication.ExecutorServer
	admCmdServer admin.AdminCommandServer
	*MsgToSignBuilder
	keySet *EnclaveKeySet
	log    logger.Logger
}

// NewStatelessExecutor creates a new stateless executor
func NewStatelessExecutor(config *Config, runtime Runtime, server communication.ExecutorServer, admCmdServer admin.AdminCommandServer, log logger.Logger) (*StatelessExecutor, error) {
	msgBuilder, err := NewMsgToSignBuilder()
	if err != nil {
		return nil, fmt.Errorf("failed to setup msg to sign builder: %w", err)
	}

	executor := &StatelessExecutor{
		config:           config,
		runtime:          runtime,
		server:           server,
		admCmdServer:     admCmdServer,
		MsgToSignBuilder: msgBuilder,
		log:              log,
	}
	// Set this executor as the request handler
	executor.server.SetRequestHandler(executor)
	// Set the connection handler to perform handshake
	executor.server.SetConnectionHandler(executor.handleNewConnection)

	executor.admCmdServer.SetCmdHandler(executor)

	return executor, nil
}

func (e *StatelessExecutor) handleNewConnection(ctx context.Context, conn communication.ServerConnection) {
	e.log.Info("Executor: New connection established, starting handshake")
	keySet, err := e.performHandshake(ctx, conn)
	if err != nil {
		e.log.Error("Executor: Handshake failed: %v", err)
		conn.Close()
		return
	}
	e.keySet = keySet
	e.log.Info("Executor: Handshake successful")
}

func (e *StatelessExecutor) performHandshake(ctx context.Context, conn communication.ServerConnection) (*EnclaveKeySet, error) {
	e.log.Info("Executor: Performing key recovery handshake")

	found, recoveryData, err := conn.GetKeysetRecovery(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get keyset recovery data: %w", err)
	}

	var keySet *EnclaveKeySet
	if found {
		e.log.Info("Executor: Keyset recovery data found, restoring keyset...")
		keySet, err = RestoreEnclaveKeySet(recoveryData)
		if err != nil {
			// Notify manager of failure
			if notifyErr := conn.KeysetRecoveryResult(ctx, err, "", ""); notifyErr != nil {
				e.log.Error("Executor: failed to send keyset recovery failure to manager: %v", notifyErr)
			}
			return nil, fmt.Errorf("failed to restore enclave keyset: %w", err)
		}

		commPubKey := crypto.ExportPublicKeyP521ToHex(keySet.CommunicationKey.PublicKey())
		signingKeyAddr := keySet.SigningKey.PublicKey().Address()

		e.log.Info("Executor: Keyset restored successfully, confirming ")
		err = conn.KeysetRecoveryResult(ctx, nil, commPubKey, signingKeyAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm keyset recovery: %w", err)
		}
	} else {
		e.log.Info("Executor: Keyset recovery data not found, generating new keyset...")
		var newRecoveryData *common.EnclaveKeySetRecovery
		keySet, newRecoveryData, err = GenerateEnclaveKeySet(e.config.KeySetRecoveryType)
		if err != nil {
			return nil, fmt.Errorf("failed to generate enclave keyset: %w", err)
		}

		commPubKey := crypto.ExportPublicKeyP521ToHex(keySet.CommunicationKey.PublicKey())
		signingKeyAddr := keySet.SigningKey.PublicKey().Address()

		err = conn.SetKeysetRecovery(ctx, newRecoveryData, commPubKey, signingKeyAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to set keyset recovery data: %w", err)
		}
		e.log.Info("Executor: New keyset generated and sent to manager for storage")
	}

	e.keySet = keySet
	// this recomputes the pub keys and address, but it is cheap anyway and one time only operation
	e.DumpPublicKeys()

	return keySet, nil
}

// Start starts the executor server
func (e *StatelessExecutor) Start(ctx context.Context) error {
	switch e.config.ChannelType {
	case "tcp":
		e.log.Info("Executor: Starting TCP executor server")
	case "vsock":
		e.log.Info("Executor: Starting v-socket executor server")
	}

	if err := e.server.Start(ctx, "Executor"); err != nil {
		return err
	}

	switch e.config.ChannelType {
	case "tcp":
		e.log.Info("Executor: Starting TCP admin executor server on %s", e.config.AdminChannelParams.(common.TcpChannelConnectionParams).Url())
	case "vsock":
		e.log.Info("Executor: Starting v-socket admin executor server on CID %d, Port %d",
			e.config.AdminChannelParams.(common.VSockChannelConnectionParams).CID,
			e.config.AdminChannelParams.(common.VSockChannelConnectionParams).Port)
	}
	return e.admCmdServer.Start(ctx, "Executor")
}

// Stop stops the executor servers
func (e *StatelessExecutor) Stop() error {
	e.log.Info("Executor: Stopping stateless executor")

	err := e.admCmdServer.Stop()
	if err != nil {
		e.log.Warn("Executor: Error stopping admin server: %v", err)
	}

	err = e.server.Stop()
	if err != nil {
		e.log.Warn("Executor: Error stopping server: %v", err)
	}

	return err
}

// Close closes the executor and its runtime
func (e *StatelessExecutor) Close() error {
	if err := e.Stop(); err != nil {
		e.log.Error("Executor: Error stopping server: %v", err)
	}
	if e.runtime != nil {
		return e.runtime.Close()
	}
	return nil
}

// HandleProcessRequest is the main workflow: decrypt state -> invoke WASM -> encrypt new state -> sign -> respond
// Returns UpdatePayload, ApplicationState, optional DeanonymizationReport, and error
func (e *StatelessExecutor) HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
	e.log.Info("Executor: Processing request %s for application %d", req.RequestID, req.ApplicationID)
	//TODO ST In order to check that the request was not tampered, we could use the requestId, that is a hash of the request parameters + an index from the ProcessorEndpoint.
	// The index could be added to the request and then the requestId could be reconstructed here and compared with the one in the request. 
	// This would ensure that the request parameters were not tampered with, otherwise the hash would not match.

	// Sanity checks that the request has not be tampered. The following are checks that should have been done by the manager or by the contract, so if they fail here
	// it means that there is a tampering or a bug in the manager. For all these checks we do not set the request as failed because it is not a fault of the request, but rather a fault of the manager or a tampering attempt.
	if req.ProtocolVersion != admittedProtocolVersion {
		return nil, nil, nil, fmt.Errorf("protocol version %d is not admitted", req.ProtocolVersion)
	}
	if req.ApplicationID != admittedAppID {
		// This won't set the request as failed because the contract doesn't accept requests with other application id, so the only way there is a wrong appId is if the request was tampered or 
		// if there is a bug in the manager. In both cases, it is not a fault of the request and so it is not set as failed.
		// When the multi app will be implemented, this will be changed accordingly
		return nil, nil, nil, fmt.Errorf("application id %s is not admitted", req.ApplicationID)
	}

	if len(req.Payload) == 0 {	
		return nil, nil, nil, fmt.Errorf("request payload is empty")
	}

	if req.RequestType != common.Process && req.RequestType != common.AssociateKey && req.RequestType != common.Deanonymize {
		return nil, nil, nil, fmt.Errorf("unsupported request type: %d", req.RequestType)
	}
	
	// This check should be updated the moment the minimum fee can be updated
	if req.MaxFeeValue.ToInt().Cmp(e.config.MinFeePerRequest) < 0 {
		return nil, nil, nil, fmt.Errorf("request fee is below minimum fee")
	}

	if appState == nil {
		errorPayload, err := e.processErrorResponse(req, 
				emptyStateRoot, 
				apperrors.New(apperrors.CodeAppStateNotFound, fmt.Sprintf("state not found for application %s", req.ApplicationID)))
		return errorPayload, nil, nil, err
	}


	// Decrypt and parse the app data
	appData, err := e.fromEncryptedStateToAppData(appState)
	if err != nil {
		return nil, nil, nil, err
	}

	// If the request contains a deposit, handle it first
	var tempState = appData.GetAppState()
	var depositEvents []common.PlainEvent
	var totalFuel *big.Int = big.NewInt(0)
	if req.DepositAmount.ToInt().Sign() > 0 {
		newState, depEvents, reqFuel, failure := e.runtime.Deposit(ctx, req.ApplicationID, req.Sender, req.DepositAmount.ToInt(), tempState, wasmModule)
		if failure != nil {
			errorPayload, err := e.processErrorResponse(req, appState.StateRoot, failure)
			return errorPayload, nil, nil, err
		}
		tempState = newState
		depositEvents = depEvents
		totalFuel = totalFuel.Add(totalFuel, reqFuel)
		e.log.Info("Executor: Successfully processed deposit for request %s", req.RequestID)
	}

	applicationFee := new(big.Int).Mul(totalFuel, e.config.FuelPricePerUnit)
	if req.MaxFeeValue.ToInt().Cmp(applicationFee) < 0 {
		errorPayload, err :=  e.processErrorResponse(req, appState.StateRoot,
			apperrors.New(apperrors.CodeInsufficientFuel,
				fmt.Sprintf("insufficient fuel: required %s wei, provided %s wei",
				applicationFee.String(),
				req.MaxFeeValue.ToInt().String(),
			),
		))
		return errorPayload, nil, nil, err
	}

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal
	var reportData []byte

	if req.RequestType == common.AssociateKey {
		//request  of type associate key: the payload is not encrypted and contains the new key
		e.log.Info("Associating new key - RequestID %s", req.RequestID)

		if len(req.Payload) != 133 {	
			return nil, nil, nil, fmt.Errorf("invalid payload length")
		}

		keyToAssociate, err := cryptotypes.NewPublicKeyP521(req.Payload)
		if err != nil {
			errorPayload, err := e.processErrorResponse(req, 
				appState.StateRoot, 
				apperrors.New(apperrors.CodeParsingKeyError, "failed to parse keyP521 in request payload"))
			return errorPayload, nil, nil, err	
		}

		totalFuel = totalFuel.Add(totalFuel, big.NewInt(10))

		appData.AddKey(req.Sender, *keyToAssociate)
	} else {
		//any other case: decrypt the payload and forward to the WASM to obtain the new state

		// Decrypt the request payload
		decryptedPayload, failure := e.decryptPayload(&e.keySet.CommunicationKey, req.Payload, req.Sender, appData.GetKeyStore())
		if failure != nil {
			errorPayload, err := e.processErrorResponse(req, 
				appState.StateRoot, 
				failure)
			return errorPayload, nil, nil, err	
		}

		// Invoke WASM method to process the request
		newState, reqEvents, reqWithdrawals, reqReportData, reqFuel, failure := e.runtime.ProcessRequest(ctx, req.ApplicationID, req.Sender, req.RequestType, decryptedPayload, tempState, wasmModule)
		if failure != nil {
			errorPayload, err := e.processErrorResponse(req, 
				appState.StateRoot, 
				failure)
			return errorPayload, nil, nil, err	
		}
		tempState = newState
		events = reqEvents
		withdrawals = reqWithdrawals
		totalFuel = totalFuel.Add(totalFuel, reqFuel)

		// Validate report generation rules:
		// 1. Reports must only be generated for Deanonymize requests
		// 2. Deanonymize requests must always generate a report
		if len(reqReportData) > 0 {
			if req.RequestType != common.Deanonymize {
				errorPayload, err := e.processErrorResponse(req, 
					appState.StateRoot, 
					apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("WASM module generated unexpected report for request type %s", req.RequestType)))
				return errorPayload, nil, nil, err	
			}
			reportData = reqReportData
		} else if req.RequestType == common.Deanonymize {
			errorPayload, err := e.processErrorResponse(req, 
				appState.StateRoot, 
				apperrors.New(apperrors.CodeRequestFuncFailed,
						"WASM module failed to generate report for Deanonymize request"))
			return errorPayload, nil, nil, err	
		}
		e.log.Info("Executor: Successfully processed request %s", req.RequestID)
	}

	// Check if there is enough ETH to cover the fuel costs
	applicationFee = new(big.Int).Mul(totalFuel, e.config.FuelPricePerUnit)


	// Application fee must be minumum fee at least
	if applicationFee.Cmp(e.config.MinFeePerRequest) < 0 {
		applicationFee = new(big.Int).Set(e.config.MinFeePerRequest)
	}

	if req.MaxFeeValue.ToInt().Cmp(applicationFee) < 0 {
		errorPayload, err :=  e.processErrorResponse(req,
			appState.StateRoot,
			apperrors.New(apperrors.CodeInsufficientFuel,
				fmt.Sprintf("insufficient fuel: required %s wei, provided %s wei",
				applicationFee.String(),
				req.MaxFeeValue.ToInt().String(),
			),
		))
		return errorPayload, nil, nil, err
	}

	// Compute refundAmount = req.MaxFeeValue - applicationFee
	refundAmount := new(big.Int).Sub(req.MaxFeeValue.ToInt(), applicationFee)

	//set the updated state
	appData.SetAppState(tempState)

	//increment appNonce
	appData.IncrementNonce()

	//serialize the new app data
	newAppData, err := appData.Serialize()
	if err != nil {
		errorPayload, err := e.processErrorResponse(req, appState.StateRoot, apperrors.New(apperrors.CodeAppDataSerializationFailure, "failed to serialize new app data"))
		return errorPayload, nil, nil, err
	}

	// Encrypt the new app data and events
	encryptedNewAppData, err := crypto.EncryptWithAES(e.keySet.StateKey, newAppData)
	if err != nil {
		e.log.Error("Executor: error encrypting initial app data: %v", err)
		return nil, nil, nil, err
	}
	// Encrypt events if they are not empty
	events = append(depositEvents, events...)
	encryptedEvents, failure := e.encryptEvents(ctx, events, req.ApplicationID, &e.keySet.CommunicationKey, e.server, appData.GetKeyStore())
	if failure != nil {
		errorPayload, err := e.processErrorResponse(req, 
			appState.StateRoot, 
			failure)
		return errorPayload, nil, nil, err	
	}
	e.log.Info("Executor: Successfully encrypted new application data")

	// Create appdata root hash
	newStateRoot := sha256.Sum256(newAppData)

	// Determine if a report was generated
	reportGenerated := len(reportData) > 0

	// Create the update payload
	updatePayload := &common.UpdatePayload{
		ApplicationID:  req.ApplicationID,
		RequestID:      req.RequestID,
		PrevStateRoot:  appState.StateRoot,
		NewStateRoot:   newStateRoot,
		Events:         encryptedEvents,
		Withdrawals:    withdrawals,
		RefundAmount:   common.ToBig(refundAmount),
		ApplicationFee: common.ToBig(applicationFee),
	}

	// Sign the update payload (produce attestation)
	signature, err := e.signUpdatePayload(updatePayload)
	if err != nil {
		e.log.Error("Executor: error signing update payload: %v", err)
		return nil, nil, nil, err
	}
	updatePayload.Signature = signature

	// Create the new application state
	newApplicationState := &common.ApplicationState{
		ApplicationID:  req.ApplicationID,
		StateRoot:      newStateRoot,
		EncryptedState: encryptedNewAppData,
	}

	// If a report was generated, encrypt it and create the DeanonymizationReport
	var deanonymizationReport *common.DeanonymizationReport
	if reportGenerated {
		encryptedReport, failure := e.encryptDeanonymizationReport(
			req.ApplicationID,
			req.RequestID,
			&e.keySet.CommunicationKey,
			req.Sender,
			reportData,
			appData.GetKeyStore(),
		)
		if failure != nil {
			errorPayload, err := e.processErrorResponse(req, 
				appState.StateRoot, 
				failure)
			return errorPayload, nil, nil, err	
		}
		e.log.Info("Executor: Successfully encrypted deanonymization report")

		deanonymizationReport = &common.DeanonymizationReport{
			ApplicationID:   req.ApplicationID,
			ReportID:        req.RequestID,
			EncryptedReport: encryptedReport,
			Authority:       req.Sender,
		}
		e.log.Info("Executor: Successfully generated deanonymization report %s", req.RequestID)
	}

	e.log.Info("Executor: Successfully processed request %s", req.RequestID)
	return updatePayload, newApplicationState, deanonymizationReport, nil
}

// HandleDeployApp implements the RequestHandler interface
func (e *StatelessExecutor) HandleDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error) {
	e.log.Info("Executor: Deploying application for request %s", req.RequestID)
	//TODO ST In order to check that the request was not tampered, we could use the requestId, that is a hash of the request parameters + an index from the ProcessorEndpoint.
	// The index could be added to the request and then the requestId could be reconstructed here and compared with the one in the request. 
	// This would ensure that the request parameters were not tampered with, otherwise the hash would not match.

	// Sanity checks that the request has not be tampered. The following are checks that should have been done by the manager or by the contract, so if they fail here
	// it means that there is a tampering or a bug in the manager. For all these checks we do not set the request as failed because it is not a fault of the request, but rather a fault of the manager or a tampering attempt.
	if req.ProtocolVersion != admittedProtocolVersion {
		return nil, nil, fmt.Errorf("protocol version %d is not admitted", req.ProtocolVersion)
	}

	if req.ApplicationID != admittedAppID {
		// This won't set the request as failed because the contract doesn't accept requests with other application id, so the only way there is a wrong appId is if the request was tampered or 
		// if there is a bug in the manager. In both cases, it is not a fault of the request and so it is not set as failed.
		// When the multi app will be implemented, this will be changed accordingly
		return nil, nil, fmt.Errorf("application id %s is not admitted", req.ApplicationID)
	}

	if req.RequestType != common.Deploy {	
		return nil, nil, fmt.Errorf("request type %s is not deploy", req.RequestType)
	}

	if len(req.Payload) == 0 {	
		return nil, nil, fmt.Errorf("request payload is empty")
	}

	// This check should be updated the moment the minimum fee can be updated
	if req.MaxFeeValue.ToInt().Cmp(e.config.MinFeePerRequest) < 0 {
		return nil, nil, fmt.Errorf("request fee is below minimum fee")
	}


	if appState != nil {
		// This means that the application was already deployed. This is the request fault, so we set the request as failed with the appropriate error code and message.
		errorPayload, err := e.processErrorResponse(req, 
				appState.StateRoot, 
				apperrors.New(apperrors.CodeApplicationAlreadyDeployed,
				fmt.Sprintf("application %s was already deployed", req.ApplicationID)))
		return errorPayload, nil, err
	}


	// For deployment, we need to initialize the application with the WASM module
	// The payload should contain the WASM bytecode
	wasmModule := req.Payload

	// Load the module and get initial state
	initialAppState, fuel, err := e.runtime.LoadModule(ctx, req.ApplicationID, wasmModule)
	if err != nil {

		errorPayload, err := e.processErrorResponse(req, 
				emptyStateRoot, 
				apperrors.New(apperrors.CodeFailedLoadingOrGettingModule,
				"failed to load or get module"))
		return errorPayload, nil, err
	}

	// Check if there is enough ETH to cover the fuel costs // TODO make a helper function?
	applicationFee := new(big.Int).Mul(fuel, e.config.FuelPricePerUnit)

	// Application fee must be minumum fee at least
	if applicationFee.Cmp(e.config.MinFeePerRequest) < 0 {
		applicationFee = new(big.Int).Set(e.config.MinFeePerRequest)
	}

	if req.MaxFeeValue.ToInt().Cmp(applicationFee) < 0 {
		errorPayload, err :=  e.processErrorResponse(req, 
			emptyStateRoot,
			apperrors.New(apperrors.CodeInsufficientFuel,
				fmt.Sprintf("insufficient fuel: required %s wei, provided %s wei",
				applicationFee.String(),
				req.MaxFeeValue.ToInt().String(),
			),
		))
		return errorPayload, nil, err
	}

	// Compute refundAmount = req.MaxFeeValue - applicationFee
	refundAmount := new(big.Int).Sub(req.MaxFeeValue.ToInt(), applicationFee)

	initialAppData := appdata.NewAppData(initialAppState)

	//serialize the new app data
	initialAppDataBytes, err := initialAppData.Serialize()
	if err != nil {
		errorPayload, err := e.processErrorResponse(req, 
			emptyStateRoot, 
			apperrors.New(apperrors.CodeAppDataSerializationFailure, "failed to serialize new app data"))
		return errorPayload, nil, err
	}
	// Create app data root hash
	initialAppDataRoot := sha256.Sum256(initialAppDataBytes)

	// Encrypt the initial state
	encryptedState, err := crypto.EncryptWithAES(e.keySet.StateKey, initialAppDataBytes)
	if err != nil {
		e.log.Error("Executor: error encrypting initial app data: %v", err)
		return nil, nil, err
	}
	e.log.Info("Executor: Successfully encrypted initial app data")

	// Create the application state
	newAppState := &common.ApplicationState{
		ApplicationID:  req.ApplicationID,
		StateRoot:      initialAppDataRoot,
		EncryptedState: encryptedState,
	}

	// Create the update payload
	updatePayload := &common.UpdatePayload{
		ApplicationID:  req.ApplicationID,
		RequestID:      req.RequestID,
		PrevStateRoot:  emptyStateRoot, // No previous state root for new applications
		NewStateRoot:   initialAppDataRoot,
		RefundAmount:   common.ToBig(refundAmount),
		ApplicationFee: common.ToBig(applicationFee),
	}

	// Sign the update payload (produce attestation)
	signature, err := e.signUpdatePayload(updatePayload)
	if err != nil {
		e.log.Error("Executor: error signing update payload: %v", err)
		return nil, nil, err
	}
	updatePayload.Signature = signature

	e.log.Info("Executor: Successfully deployed application %d", req.ApplicationID)
	return updatePayload, newAppState, nil
}

func (e *StatelessExecutor) fromEncryptedStateToAppData(encState *common.ApplicationState) (*appdata.AppData, error) {
	// Decrypt the encrypted state
	decryptedState, err := e.DecryptState(encState.EncryptedState, e.keySet.StateKey)
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

// buildErrorPayload creates a signed error payload for failed requests.
// The payload has: prevStateRoot == newStateRoot (unchanged), empty events/withdrawals,
// refund = MaxFeeValue - MinFee, and the error code/message from the failure.
func (e *StatelessExecutor) buildErrorPayload(req *common.Request, stateRoot [32]byte, failure *apperrors.RequestFailure) (*common.UpdatePayload, error) {
	// Truncate error message if longer than 100 characters
	errorMsg := failure.ExternalMessage()
	if len(errorMsg) > 100 {
		errorMsg = errorMsg[:100]
	}

	// Create the error payload with state unchanged
	updatePayload := &common.UpdatePayload{
		ApplicationID:   req.ApplicationID,
		RequestID:       req.RequestID,
		PrevStateRoot:   stateRoot,
		NewStateRoot:    stateRoot, // State unchanged on error
		Events:          nil,       // Empty events on error
		Withdrawals:     nil,       // Empty withdrawals on error
		RefundAmount:    req.MaxFeeValue,		// Refund and fees are calculated on chain
		ApplicationFee:  common.NewBig(0), // No application fee on error, but the actual fee handling is done on chain 
		ErrorCode:       failure.Category(),
		ErrorMsg:        errorMsg,
	}

	// Sign the error payload
	signature, err := e.signUpdatePayload(updatePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to sign error payload: %w", err)
	}
	updatePayload.Signature = signature

	return updatePayload, nil
}

// processErrorResponse creates a signed error response for HandleProcessRequest.
// If building the error payload fails, it falls back to returning the original error.
func (e *StatelessExecutor) processErrorResponse(req *common.Request, stateRoot [32]byte, failure *apperrors.RequestFailure) (*common.UpdatePayload, error) {
	payload, err := e.buildErrorPayload(req, stateRoot, failure)
	if err != nil {
		e.log.Error("Executor: Failed to build error payload: %v, returning unsigned error", err)
		return nil, err
	}
	e.log.Info("Executor: Returning signed error payload for request %s (error code: %d)", req.RequestID, payload.ErrorCode)
	return payload, nil
}


// signUpdatePayload signs the update payload to produce an attestation
func (e *StatelessExecutor) signUpdatePayload(payload *common.UpdatePayload) ([]byte, error) {
	// Serialize the payload for signing
	hash, err := e.BuildMsgHash(payload)

	if err != nil {
		return nil, fmt.Errorf("failed to create message to sign: %w", err)
	}

	// Sign the hash
	signature, err := e.keySet.SigningKey.Sign(hash)
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

func (e *StatelessExecutor) encryptEvents(ctx context.Context, events []common.PlainEvent, appId common.ApplicationIdType, key *cryptotypes.PrivateKeyP521, server communication.ExecutorServer, keyStore appdata.KeyStore) ([]common.Event, *apperrors.RequestFailure) {
	if len(events) == 0 {
		return nil, nil // No events to encrypt
	}
	encryptedEvents := make([]common.Event, len(events))

	for i, event := range events {
		// retrieve user Secp521r1_PubKey
		userKey, exists := keyStore[event.UserID]

		if !exists {
			return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, fmt.Sprintf("no Secp521r1_PubKey found for user %s", event.UserID))
		}
		// Encrypt the event data
		encryptedData, err := crypto.Encrypt(key, userKey, event.Data)
		if err != nil {
			return nil, apperrors.New(apperrors.CodeWrongKey, "failed to encrypt event data")
		}

		// Create the encrypted event
		encryptedEvents[i] = common.Event{
			ApplicationID: appId,
			UserID:        event.UserID,
			EventSubType:  event.EventSubType,
			EncryptedData: encryptedData,
		}
	}

	e.log.Info("Executor: Successfully encrypted %d events", len(events))
	return encryptedEvents, nil
}

func (e *StatelessExecutor) encryptDeanonymizationReport(applicationId common.ApplicationIdType, requestId common.RequestIdType, key *cryptotypes.PrivateKeyP521, requester ethCommon.Address, reportData []byte, keyStore appdata.KeyStore) ([]byte, *apperrors.RequestFailure) {
	if len(reportData) == 0 {
		return nil, apperrors.New(apperrors.CodeNoReportDataFound, "no report data found")
	}

	// retrieve user Secp521r1_PubKey
	requesterPublicKey, exists := keyStore[requester]
	if !exists {
		return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, fmt.Sprintf("no Secp521r1_PubKey found for user %s", requester))
	}

	// Unencrypted deanonymization reports are specific to the application, we can not assume a defined struct of the reportData.
	// therefore we decide to add the raw bytes into a separate field
	data := common.DecryptedReport{
		ApplicationID:   applicationId,
		RequestID:       requestId,
		ReportDataBytes: reportData,
	}

	// Marshal the updated report data
	updatedReportData, err := json.Marshal(data)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeJsonMarshalError, "failed to marshal updated deanonymization report")
	}

	// Encrypt the report data
	encryptedReport, err := crypto.Encrypt(key, requesterPublicKey, updatedReportData)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeWrongKey, "failed to encrypt deanonymization report")
	}

	return encryptedReport, nil
}

func (e *StatelessExecutor) DecryptState(encryptedState []byte, decryptionKey cryptotypes.AES256Key) ([]byte, error) {
	if len(encryptedState) == 0 {
		return encryptedState, nil // nothing to decrypt
	}

	decryptedState, err := crypto.DecryptWithAES(decryptionKey, encryptedState)
	if err != nil {
		return nil, fmt.Errorf("state decryption failed: %w", err)
	}

	e.log.Info("Executor: Successfully decrypted application state")
	return decryptedState, nil
}

func (e *StatelessExecutor) decryptPayload(decryptionKey *cryptotypes.PrivateKeyP521, payload []byte, sender ethCommon.Address, keyStore appdata.KeyStore) ([]byte, *apperrors.RequestFailure) {
	if len(payload) == 0 {
		return payload, nil // No payload to decrypt
	}

	// retrieve sender Secp521r1_PubKey
	userKey, exists := keyStore[sender]
	if !exists {
		return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, fmt.Sprintf("no Secp521r1_PubKey found for sender %s", sender))
	}

	decryptedPayload, err := crypto.Decrypt(userKey, decryptionKey, payload)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeWrongKey, "payload decryption failed")
	}

	e.log.Info("Executor: Successfully decrypted request payload")
	return decryptedPayload, nil
}

// ExecuteCommand implements admin.AdminCmdHandler interface.
// Handles admin commands for the executor, currently only KeyAttestationRequestMessage.
func (e *StatelessExecutor) ExecuteCommand(ctx context.Context, msg admin.AdminMessage) (interface{}, error) {
	switch msg.Type {
	case admin.KeyAttestationRequestMessage:
		return e.CreateKeyAttestation(ctx)
	default:
		return nil, fmt.Errorf("unsupported command type: %v", msg.Type)
	}
}

func (e *StatelessExecutor) CreateKeyAttestation(ctx context.Context) ([]byte, error) {
	return e.createKeyAttestationInternal(ctx, func() (NsmSession, error) {
		s, err := nsm.OpenDefaultSession()
		if err != nil {
			return nil, err
		}
		return s, nil
	})
}

func (e *StatelessExecutor) createKeyAttestationInternal(ctx context.Context, nsmSessionOpener func() (NsmSession, error)) ([]byte, error) {
	keySet := e.keySet
	if keySet == nil {
		return nil, fmt.Errorf("keyset is empty")
	}

	session, err := nsmSessionOpener()
	if err != nil {
		e.log.Info("Executor: error opening nms session: %v", err)
		return nil, fmt.Errorf("failed to generate attestation")
	}
	defer func() {
		if err := session.Close(); err != nil {
			e.log.Warn("Executor: error closing nms session: %v", err)
		}
	}()

	signerAddress, err := hex.DecodeString(keySet.SigningKey.PublicKey().Address()[2:]) // remove 0x prefix
	if err != nil {
		e.log.Error("Executor: error while decoding signer address: %v", err)
		return nil, fmt.Errorf("failed to generate attestation")
	}

	res, err := session.Send(&request.Attestation{PublicKey: keySet.CommunicationKey.PublicKey().Bytes(), UserData: signerAddress})
	if err != nil {
		e.log.Warn("failed to get attestation: %v", err)
		return nil, fmt.Errorf("failed to generate attestation")
	}
	if res.Error != "" {
		e.log.Warn("NSM device returned an error: %s", res.Error)
		return nil, fmt.Errorf("failed to generate attestation")
	}
	if res.Attestation == nil || res.Attestation.Document == nil {
		e.log.Warn("NSM device did not return an attestation")
		return nil, fmt.Errorf("failed to generate attestation")
	}

	return res.Attestation.Document, nil

}
