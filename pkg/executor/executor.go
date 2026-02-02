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
	"github.com/horizen-pes/pkg/executor/kms"
	"github.com/horizen-pes/pkg/logger"
)

// NsmSession is an interface abstracting nsm.Session for testability.
type NsmSession interface {
	Send(req request.Request) (response.Response, error)
	Close() error
}

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
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - recoveryType: 0 for unsafe/dev mode, 1 for AWS KMS
//   - kmsClient: KMS client (required for Type 1, can be nil for Type 0)
//   - enclaveHandle: Enclave handle for attestation (required for Type 1, can be nil for Type 0)
//   - kmsKeyARN: AWS KMS key ARN (required for Type 1, can be empty for Type 0)
func GenerateEnclaveKeySet(
	ctx context.Context,
	recoveryType int,
	kmsClient kms.KMSClient,
	enclaveHandle kms.EnclaveHandle,
	kmsKeyARN string,
) (*EnclaveKeySet, *common.EnclaveKeySetRecovery, error) {
	switch recoveryType {
	case 0:
		// Type 0: Unsafe/development mode - masterKey stored in plaintext
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

	case 1:
		// Type 1: AWS KMS - masterKey encrypted with KMS using Nitro attestation
		// Validate required dependencies
		if kmsClient == nil {
			return nil, nil, fmt.Errorf("KMS client is required for Type 1 recovery")
		}
		if enclaveHandle == nil {
			return nil, nil, fmt.Errorf("enclave handle is required for Type 1 recovery")
		}
		if kmsKeyARN == "" {
			return nil, nil, fmt.Errorf("KMS key ARN is required for Type 1 recovery")
		}

		// 1. Generate attestation document from NSM
		attestationDoc, err := enclaveHandle.Attest(nil)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate attestation: %w", err)
		}

		// 2. Request KMS to generate a data key with attestation
		// KMS will validate the attestation (PCR values) before responding
		dataKeyOutput, err := kmsClient.GenerateDataKeyWithAttestation(ctx, kmsKeyARN, attestationDoc)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate data key from KMS: %w", err)
		}

		// Validate KMS response contains required data for recovery
		if len(dataKeyOutput.CiphertextBlob) == 0 {
			return nil, nil, fmt.Errorf("KMS returned empty CiphertextBlob - cannot create recoverable keyset")
		}
		if len(dataKeyOutput.CiphertextForRecipient) == 0 {
			return nil, nil, fmt.Errorf("KMS returned empty CiphertextForRecipient - attestation may have failed")
		}

		// 3. Decrypt the data key using enclave's private RSA key
		// The CiphertextForRecipient is encrypted for the enclave's public key from attestation
		masterKeyBytes, err := enclaveHandle.DecryptKMSEnvelopedKey(dataKeyOutput.CiphertextForRecipient)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decrypt KMS enveloped key: %w", err)
		}

		if len(masterKeyBytes) != 32 {
			return nil, nil, fmt.Errorf("invalid master key length from KMS: got %d, expected 32", len(masterKeyBytes))
		}

		var masterKey cryptotypes.AES256Key
		copy(masterKey[:], masterKeyBytes)

		// 4. Create a new EnclaveKeySet
		keySet, err := CreateNewKeySet()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create key set: %w", err)
		}

		// 5. Serialize the EnclaveKeySet
		serializedKeySet, err := keySet.Serialize()
		if err != nil {
			return nil, nil, err
		}

		// 6. Encrypt the serialized EnclaveKeySet with the master key
		encryptedKeySet, err := crypto.EncryptWithAES(masterKey, serializedKeySet)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt key set: %w", err)
		}

		// 7. Create the recovery structure
		// Store the KMS CiphertextBlob
		// This can only be decrypted by KMS with valid attestation
		recovery := &common.EnclaveKeySetRecovery{
			RecoveryType:       1,
			KeySetCiphertext:   encryptedKeySet,
			RecoveryCiphertext: dataKeyOutput.CiphertextBlob,
		}

		return keySet, recovery, nil

	default:
		return nil, nil, fmt.Errorf("unsupported recovery type: %d", recoveryType)
	}
}

// RestoreEnclaveKeySet recovers a keyset from a recovery previously stored.
// It will be used when starting the executors after the first init.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - recovery: The recovery structure containing encrypted keyset data
//   - kmsClient: KMS client (required for Type 1, can be nil for Type 0)
//   - enclaveHandle: Enclave handle for attestation (required for Type 1, can be nil for Type 0)
func RestoreEnclaveKeySet(
	ctx context.Context,
	recovery *common.EnclaveKeySetRecovery,
	kmsClient kms.KMSClient,
	enclaveHandle kms.EnclaveHandle,
) (*EnclaveKeySet, error) {
	switch recovery.RecoveryType {
	case 0:
		// Type 0: Unsafe/development mode - masterKey stored in plaintext
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

	case 1:
		// Type 1: AWS KMS - masterKey encrypted with KMS using Nitro attestation
		// Validate required dependencies
		if kmsClient == nil {
			return nil, fmt.Errorf("KMS client is required for Type 1 recovery")
		}
		if enclaveHandle == nil {
			return nil, fmt.Errorf("enclave handle is required for Type 1 recovery")
		}

		// 1. Generate fresh attestation document
		// A new enclave instance needs to prove its identity to KMS
		attestationDoc, err := enclaveHandle.Attest(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to generate attestation: %w", err)
		}

		// 2. Decrypt the master key using KMS with attestation
		// RecoveryCiphertext contains the KMS CiphertextBlob from generation
		// KMS will validate the attestation (PCR values) before decrypting
		ciphertextForRecipient, err := kmsClient.DecryptWithAttestation(ctx, recovery.RecoveryCiphertext, attestationDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt master key from KMS: %w", err)
		}

		// 3. Decrypt the enveloped key using enclave's private RSA key
		masterKeyBytes, err := enclaveHandle.DecryptKMSEnvelopedKey(ciphertextForRecipient)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt KMS enveloped key: %w", err)
		}

		if len(masterKeyBytes) != 32 {
			return nil, fmt.Errorf("invalid master key length from KMS: got %d, expected 32", len(masterKeyBytes))
		}

		var masterKey cryptotypes.AES256Key
		copy(masterKey[:], masterKeyBytes)

		// 4. Decrypt the keyset using the recovered master key
		decryptedKeySet, err := crypto.DecryptWithAES(masterKey, recovery.KeySetCiphertext)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt key set: %w", err)
		}

		// 5. Deserialize the decrypted data into an EnclaveKeySet
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

	// KMS dependencies (nil for Type 0, required for Type 1)
	kmsClient     kms.KMSClient
	enclaveHandle kms.EnclaveHandle
}

// NewStatelessExecutor creates a new stateless executor.
//
// Parameters:
//   - config: Executor configuration
//   - runtime: WASM runtime for executing modules
//   - server: Communication server for manager connections
//   - admCmdServer: Admin command server
//   - log: Logger instance
//   - kmsClient: KMS client for Type 1 recovery (can be nil for Type 0)
//   - enclaveHandle: Enclave handle for attestation (can be nil for Type 0)
func NewStatelessExecutor(
	config *Config,
	runtime Runtime,
	server communication.ExecutorServer,
	admCmdServer admin.AdminCommandServer,
	log logger.Logger,
	kmsClient kms.KMSClient,
	enclaveHandle kms.EnclaveHandle,
) (*StatelessExecutor, error) {
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
		kmsClient:        kmsClient,
		enclaveHandle:    enclaveHandle,
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
	e.log.Info("Executor: Performing key recovery handshake (Type %d)", e.config.KeySetRecoveryType)

	found, recoveryData, err := conn.GetKeysetRecovery(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get keyset recovery data: %w", err)
	}

	var keySet *EnclaveKeySet
	if found {
		e.log.Info("Executor: Keyset recovery data found (Type %d), restoring keyset...", recoveryData.RecoveryType)
		keySet, err = RestoreEnclaveKeySet(ctx, recoveryData, e.kmsClient, e.enclaveHandle)
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
		e.log.Info("Executor: Keyset recovery data not found, generating new keyset (Type %d)...", e.config.KeySetRecoveryType)
		var newRecoveryData *common.EnclaveKeySetRecovery
		keySet, newRecoveryData, err = GenerateEnclaveKeySet(
			ctx,
			e.config.KeySetRecoveryType,
			e.kmsClient,
			e.enclaveHandle,
			e.config.KMSKeyARN,
		)
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

	err :=  e.admCmdServer.Stop();
	if err != nil {
		e.log.Warn("Executor: Error stopping admin server: %v", err)
	}

	err = e.server.Stop();
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
func (e *StatelessExecutor) HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure) {
	e.log.Info("Executor: Processing request %s for application %d", req.RequestID, req.ApplicationID)

	// Decrypte and parse the app data
	appData, err := e.fromEncryptedStateToAppData(appState)
	if err != nil {
		return nil, nil, apperrors.New(apperrors.CodeEncryptedToAppDataFailure, "failed to decrypt data", err)
	}

	// If the request contains a deposit, handle it first
	var tempState = appData.GetAppState()
	var depositEvents []common.PlainEvent
	var totalFuel *big.Int = big.NewInt(0)
	if req.DepositAmount.Sign() > 0 {
		newState, depEvents, reqFuel, failure := e.runtime.Deposit(ctx, req.ApplicationID, req.Sender, req.DepositAmount, tempState, wasmModule)
		if failure != nil {
			return nil, nil, failure
		}
		tempState = newState
		depositEvents = depEvents
		totalFuel = totalFuel.Add(totalFuel, reqFuel)
		e.log.Info("Executor: Successfully processed deposit for request %s", req.RequestID)
	}

	applicationFee := new(big.Int).Mul(totalFuel, e.config.FuelPricePerUnit)
	if req.MaxFeeValue.Cmp(applicationFee) < 0 {
		return nil, nil, apperrors.New(
			apperrors.CodeInsufficientFuel,
			fmt.Sprintf(
				"insufficient fuel: required %s wei, provided %s wei",
				applicationFee.String(),
				req.MaxFeeValue.String(),
			),
			nil,
		)
	}

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal

	if req.RequestType == common.AssociateKey {
		//request  of type associate key: the payload is not encrypted and contains the new key
		e.log.Info("Associating new key - RequestID %s", req.RequestID)

		keyToAssociate, err := cryptotypes.NewPublicKeyP521(req.Payload)
		if err != nil {
			return nil, nil, apperrors.New(apperrors.CodeParsingKeyError, "failed to parse keyP521 in request payload", err)
		}

		totalFuel = totalFuel.Add(totalFuel, big.NewInt(10))

		appData.AddKey(req.Sender, *keyToAssociate)
	} else {
		//any other case: decrypt the payload and forward to the WASM to obtain the new state

		// Decrypt the request payload
		decryptedPayload, failure := e.decryptPayload(&e.keySet.CommunicationKey, req.Payload, req.Sender, appData.GetKeyStore())
		if failure != nil {
			return nil, nil, failure
		}

		// Invoke WASM method to process the request
		newState, reqEvents, reqWithdrawals, reqFuel, failure := e.runtime.ProcessRequest(ctx, req.ApplicationID, req.Sender, decryptedPayload, tempState, wasmModule)
		if failure != nil {
			return nil, nil, failure
		}
		tempState = newState
		events = reqEvents
		withdrawals = reqWithdrawals
		totalFuel = totalFuel.Add(totalFuel, reqFuel)
		e.log.Info("Executor: Successfully processed request %s", req.RequestID)
	}

	// Check if there is enough ETH to cover the fuel costs
	applicationFee = new(big.Int).Mul(totalFuel, e.config.FuelPricePerUnit)
	if req.MaxFeeValue.Cmp(applicationFee) < 0 {
		return nil, nil, apperrors.New(
			apperrors.CodeInsufficientFuel,
			fmt.Sprintf(
				"insufficient fuel: required %s wei, provided %s wei",
				applicationFee.String(),
				req.MaxFeeValue.String(),
			),
			nil,
		)
	}

	// Application fee must be minumum fee at least
	if applicationFee.Cmp(e.config.MinFeePerRequest) < 0 {
		applicationFee = new(big.Int).Set(e.config.MinFeePerRequest)
	}

	// Compute refundAmount = req.MaxFeeValue - applicationFee
	refundAmount := new(big.Int).Sub(req.MaxFeeValue, applicationFee)

	//set the updated state
	appData.SetAppState(tempState)

	//increment appNonce
	appData.IncrementNonce()

	//serialize the new app data
	newAppData, err := appData.Serialize()
	if err != nil {
		return nil, nil, apperrors.New(apperrors.CodeAppDataSerializationFailure, "failed to serialize new appData", err)
	}

	// Encrypt the new app data and events
	encryptedNewAppData, err := crypto.EncryptWithAES(e.keySet.StateKey, newAppData)
	if err != nil {
		return nil, nil, apperrors.New(apperrors.CodeAppDataEncryptionFailure, "failed to encrypt new appData", err)
	}
	// Encrypt events if they are not empty
	events = append(depositEvents, events...)
	encryptedEvents, failure := e.encryptEvents(ctx, events, req.ApplicationID, &e.keySet.CommunicationKey, e.server, appData.GetKeyStore())
	if failure != nil {
		return nil, nil, failure
	}
	e.log.Info("Executor: Successfully encrypted new application data")

	// Create appdata root hash
	newStateRoot := sha256.Sum256(newAppData)

	// Create the update payload
	updatePayload := &common.UpdatePayload{
		ApplicationID:  req.ApplicationID,
		RequestID:      req.RequestID,
		PrevStateRoot:  appState.StateRoot,
		NewStateRoot:   newStateRoot,
		Events:         encryptedEvents,
		Withdrawals:    withdrawals,
		RefundAmount:   refundAmount,
		ApplicationFee: applicationFee,
	}

	// Sign the update payload (produce attestation)
	signature, err := e.signUpdatePayload(updatePayload)
	if err != nil {
		return nil, nil, apperrors.New(apperrors.CodePayloadUpdateSigningFailure, "failed to sign update payload", err)
	}
	updatePayload.Signature = signature

	// Create the new application state
	newApplicationState := &common.ApplicationState{
		ApplicationID:  req.ApplicationID,
		StateRoot:      newStateRoot,
		EncryptedState: encryptedNewAppData,
	}

	e.log.Info("Executor: Successfully processed request %s", req.RequestID)
	return updatePayload, newApplicationState, nil
}

// HandleDeployApp implements the RequestHandler interface
func (e *StatelessExecutor) HandleDeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, *apperrors.RequestFailure) {
	e.log.Info("Executor: Deploying application for request %s", req.RequestID)

	// For deployment, we need to initialize the application with the WASM module
	// The payload should contain the WASM bytecode
	wasmModule := req.Payload

	// Load the module and get initial state
	initialAppState, fuel, err := e.runtime.LoadModule(ctx, req.ApplicationID, wasmModule)
	if err != nil {
		return nil, nil, apperrors.New(apperrors.CodeFailedLoadingOrGettingModule, "failed to load or get module", err)
	}

	// Check if there is enough ETH to cover the fuel costs // TODO make a helper function?
	applicationFee := new(big.Int).Mul(fuel, e.config.FuelPricePerUnit)
	if req.MaxFeeValue.Cmp(applicationFee) < 0 {
		return nil, nil, apperrors.New(
			apperrors.CodeInsufficientFuel,
			fmt.Sprintf(
				"insufficient fuel: required %s wei, provided %s wei",
				applicationFee.String(),
				req.MaxFeeValue.String(),
			),
			nil,
		)
	}

	// Application fee must be minumum fee at least
	if applicationFee.Cmp(e.config.MinFeePerRequest) < 0 {
		applicationFee = new(big.Int).Set(e.config.MinFeePerRequest)
	}

	// Compute refundAmount = req.MaxFeeValue - applicationFee
	refundAmount := new(big.Int).Sub(req.MaxFeeValue, applicationFee)

	initialAppData := appdata.NewAppData(initialAppState)

	//serialize the new app data
	initialAppDataBytes, err := initialAppData.Serialize()
	if err != nil {
		return nil, nil, apperrors.New(apperrors.CodeAppDataSerializationFailure, "failed to serialize new app data", err)
	}
	// Create app data root hash
	initialAppDataRoot := sha256.Sum256(initialAppDataBytes)

	// Encrypt the initial state
	encryptedState, err := crypto.EncryptWithAES(e.keySet.StateKey, initialAppDataBytes)
	if err != nil {
		return nil, nil, apperrors.New(apperrors.CodeAppDataEncryptionFailure, "failed to encrypt initial app data", err)
	}
	e.log.Info("Executor: Successfully encrypted initial app data")

	// Create the application state
	appState := &common.ApplicationState{
		ApplicationID:  req.ApplicationID,
		StateRoot:      initialAppDataRoot,
		EncryptedState: encryptedState,
	}

	// Create the update payload
	updatePayload := &common.UpdatePayload{
		ApplicationID:  req.ApplicationID,
		RequestID:      req.RequestID,
		PrevStateRoot:  [32]byte{}, // No previous state root for new applications
		NewStateRoot:   initialAppDataRoot,
		RefundAmount:   refundAmount,
		ApplicationFee: applicationFee,
	}

	// Sign the update payload (produce attestation)
	signature, err := e.signUpdatePayload(updatePayload)
	if err != nil {
		return nil, nil, apperrors.New(apperrors.CodePayloadUpdateSigningFailure, "failed to sign update payload", err)
	}
	updatePayload.Signature = signature

	e.log.Info("Executor: Successfully deployed application %d", req.ApplicationID)
	return updatePayload, appState, nil
}

// HandleGenerateDeanonymizationReport implements the RequestHandler interface
func (e *StatelessExecutor) HandleGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, *apperrors.RequestFailure) {
	e.log.Info("Executor: Generating deanonymization report for request %s", req.RequestID)

	// Decrypte and parse the app data
	appData, err := e.fromEncryptedStateToAppData(appState)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeEncryptedToAppDataFailure, "failed to decrypt data", err)
	}

	// Decrypt the request payload
	decryptedPayload, failure := e.decryptPayload(&e.keySet.CommunicationKey, req.Payload, req.Sender, appData.GetKeyStore())
	if failure != nil {
		return nil, failure
	}

	// Generate the report using the runtime
	reportData, fuel, failure := e.runtime.GenerateDeanonymizationReport(ctx, req.ApplicationID, decryptedPayload, appData.GetAppState(), wasmModule)
	if failure != nil {
		return nil, failure
	}

	// Check if there is enough ETH to cover the fuel costs
	applicationFee := new(big.Int).Mul(fuel, e.config.FuelPricePerUnit)
	if req.MaxFeeValue.Cmp(applicationFee) < 0 {
		return nil, apperrors.New(
			apperrors.CodeInsufficientFuel,
			fmt.Sprintf(
				"insufficient fuel: required %s wei, provided %s wei",
				applicationFee.String(),
				req.MaxFeeValue.String(),
			),
			nil,
		)
	}

	// Application fee must be minumum fee at least
	if applicationFee.Cmp(e.config.MinFeePerRequest) < 0 {
		applicationFee = new(big.Int).Set(e.config.MinFeePerRequest)
	}

	// Compute refundAmount = req.MaxFeeValue - applicationFee
	refundAmount := new(big.Int).Sub(req.MaxFeeValue, applicationFee)

	// Encrypt the report
	encryptedReport, failure := e.encryptDeanonymizationReport(
		req.ApplicationID,
		req.RequestID,
		&e.keySet.CommunicationKey,
		req.Sender,
		reportData,
		appData.GetKeyStore(),
	)
	if failure != nil {
		return nil, failure
	}
	e.log.Info("Executor: Successfully encrypted deanonymization report")

	report := &common.DeanonymizationReport{
		ApplicationID:   req.ApplicationID,
		ReportID:        req.RequestID,
		EncryptedReport: encryptedReport,
		Authority:       req.Sender,
		RefundAmount:    refundAmount,
		ApplicationFee:  applicationFee,
	}

	e.log.Info("Executor: Successfully generated deanonymization report %s", req.RequestID)
	return report, nil
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
			return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, fmt.Sprintf("no Secp521r1_PubKey found for user %s", event.UserID), nil)
		}
		// Encrypt the event data
		encryptedData, err := crypto.Encrypt(key, userKey, event.Data)
		if err != nil {
			return nil, apperrors.New(apperrors.CodeWrongKey, "failed to encrypt event data", err)
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
		return nil, apperrors.New(apperrors.CodeNoReportDataFound, "no report data found", nil)
	}

	// retrieve user Secp521r1_PubKey
	requesterPublicKey, exists := keyStore[requester]
	if !exists {
		return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, fmt.Sprintf("no Secp521r1_PubKey found for user %s", requester), nil)
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
		return nil, apperrors.New(apperrors.CodeJsonMarshalError, "failed to marshal updated deanonymization report", err)
	}

	// Encrypt the report data
	encryptedReport, err := crypto.Encrypt(key, requesterPublicKey, updatedReportData)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeWrongKey, "failed to encrypt deanonymization report", err)
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
		return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, fmt.Sprintf("no Secp521r1_PubKey found for sender %s", sender), nil)
	}

	decryptedPayload, err := crypto.Decrypt(userKey, decryptionKey, payload)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeWrongKey, "payload decryption failed", err)
	}

	e.log.Info("Executor: Successfully decrypted request payload")
	return decryptedPayload, nil
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