package executor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/hf/nsm"

	"github.com/HorizenOfficial/vela/pkg/admin"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/appdata"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	"github.com/HorizenOfficial/vela/pkg/communication"
	"github.com/HorizenOfficial/vela/pkg/crypto"
	"github.com/HorizenOfficial/vela/pkg/executor/kms"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/nsmutil"
	"github.com/HorizenOfficial/vela/pkg/version"
)

// NsmSession is an interface abstracting nsm.Session for testability.
type NsmSession = nsmutil.Session

var (
	admittedProtocolVersion = uint8(0)
	emptyStateRoot          = [32]byte{}
)

const (
	deployDescriptorFailureMsg = "failed to deploy application"
	deployLoadFailureMsg       = "failed to load or get module"
)

func CreateNewKeySet() (*EnclaveKeySet, error) {
	// If all three fixed key env vars are set, import from hex instead of generating.
	// This allows using deterministic keys in dev/test environments.
	fixedSigningKey := os.Getenv("EXECUTOR_FIXED_SIGNING_KEY")
	fixedCommKey := os.Getenv("EXECUTOR_FIXED_COMMUNICATION_KEY")
	fixedStateKey := os.Getenv("EXECUTOR_FIXED_STATE_KEY")

	if fixedSigningKey != "" && fixedCommKey != "" && fixedStateKey != "" {
		return importFixedKeySet(fixedSigningKey, fixedCommKey, fixedStateKey)
	}

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

// importFixedKeySet creates an EnclaveKeySet from fixed hex-encoded private keys.
// !!! UNSAFE !!! use this method only for DEV environment !!!
func importFixedKeySet(signingHex, commHex, stateHex string) (*EnclaveKeySet, error) {
	signingKey, err := crypto.ImportPrivateKeySecp256k1FromHex(signingHex)
	if err != nil {
		return nil, fmt.Errorf("failed to import fixed signing key: %w", err)
	}
	communicationKey, err := crypto.ImportPrivateKeyP521FromHex(commHex)
	if err != nil {
		return nil, fmt.Errorf("failed to import fixed communication key: %w", err)
	}
	stateKey, err := crypto.ImportKeyAESFromHex(stateHex)
	if err != nil {
		return nil, fmt.Errorf("failed to import fixed state key: %w", err)
	}
	return &EnclaveKeySet{
		CommunicationKey: *communicationKey,
		SigningKey:       *signingKey,
		StateKey:         *stateKey,
	}, nil
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
	recoveryType common.RecoveryType,
	kmsClient kms.KMSClient,
	enclaveHandle kms.EnclaveHandle,
	kmsKeyARN string,
) (*EnclaveKeySet, *common.EnclaveKeySetRecovery, error) {
	if recoveryType != common.RecoveryTypeUnsafe && recoveryType != common.RecoveryTypeKMS {
		return nil, nil, fmt.Errorf("unsupported recovery type: %d", recoveryType)
	}

	// Validate required dependencies for Type 1.
	if recoveryType == common.RecoveryTypeKMS {
		if kmsClient == nil {
			return nil, nil, fmt.Errorf("KMS client is required for Type 1 recovery")
		}
		if enclaveHandle == nil {
			return nil, nil, fmt.Errorf("enclave handle is required for Type 1 recovery")
		}
		if kmsKeyARN == "" {
			return nil, nil, fmt.Errorf("KMS key ARN is required for Type 1 recovery")
		}
	}

	// 1. Create a new EnclaveKeySet.
	keySet, err := CreateNewKeySet()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create key set: %w", err)
	}

	// 2. Serialize the EnclaveKeySet.
	serializedKeySet, err := keySet.Serialize()
	if err != nil {
		return nil, nil, err
	}

	var masterKey cryptotypes.AES256Key
	var recoveryCiphertext []byte
	var cleanup func()
	defer func() {
		if cleanup != nil {
			cleanup()
		}
	}()

	switch recoveryType {
	case common.RecoveryTypeUnsafe:
		// Type 0: Unsafe/development mode - masterKey stored in plaintext
		// 3. Create a new AES master key.
		generatedKey, err := crypto.GenerateAESKey()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate master key: %w", err)
		}
		masterKey = generatedKey
		cleanup = func() { zeroAESKey(&masterKey) }

		// Create recovery ciphertext (plaintext master key for Type 0).
		recoveryKey := make([]byte, len(masterKey))
		copy(recoveryKey, masterKey[:])
		recoveryCiphertext = recoveryKey

	case common.RecoveryTypeKMS:
		// Type 1: AWS KMS - masterKey encrypted with KMS using Nitro attestation
		// 3. Generate attestation document from NSM.
		attestationDoc, err := enclaveHandle.Attest(nil)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate attestation: %w", err)
		}

		// 4. Request KMS to generate a data key with attestation.
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

		// 5. Decrypt the data key using enclave's private RSA key.
		// The CiphertextForRecipient is encrypted for the enclave's public key from attestation
		masterKeyBytes, err := enclaveHandle.DecryptKMSEnvelopedKey(dataKeyOutput.CiphertextForRecipient)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decrypt KMS enveloped key: %w", err)
		}
		cleanup = func() {
			zeroBytes(masterKeyBytes)
			zeroAESKey(&masterKey)
		}

		if len(masterKeyBytes) != 32 {
			return nil, nil, fmt.Errorf("invalid master key length from KMS: got %d, expected 32", len(masterKeyBytes))
		}

		copy(masterKey[:], masterKeyBytes)
		// Store the KMS CiphertextBlob (decryptable only by KMS with valid attestation).
		recoveryCiphertext = dataKeyOutput.CiphertextBlob
	}

	// Encrypt the serialized EnclaveKeySet with the master key.
	encryptedKeySet, err := crypto.EncryptWithAES(masterKey, serializedKeySet)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt key set: %w", err)
	}

	recovery := &common.EnclaveKeySetRecovery{
		RecoveryType:       recoveryType,
		KeySetCiphertext:   encryptedKeySet,
		RecoveryCiphertext: recoveryCiphertext,
	}

	return keySet, recovery, nil
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
	if recovery == nil {
		return nil, fmt.Errorf("recovery is required")
	}

	switch recovery.RecoveryType {
	case common.RecoveryTypeUnsafe:
		// Type 0: Unsafe/development mode - masterKey stored in plaintext
		// 1. Use the RecoveryCiphertext as the master key to decrypt the KeySetCiphertext.
		var masterKey cryptotypes.AES256Key
		copy(masterKey[:], recovery.RecoveryCiphertext)
		defer zeroAESKey(&masterKey)

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

	case common.RecoveryTypeKMS:
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
		defer zeroBytes(masterKeyBytes)

		if len(masterKeyBytes) != 32 {
			return nil, fmt.Errorf("invalid master key length from KMS: got %d, expected 32", len(masterKeyBytes))
		}

		var masterKey cryptotypes.AES256Key
		copy(masterKey[:], masterKeyBytes)
		defer zeroAESKey(&masterKey)

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

func zeroBytes(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

func zeroAESKey(key *cryptotypes.AES256Key) {
	if key == nil {
		return
	}
	for i := range key {
		key[i] = 0
	}
}

// StatelessExecutor implements the Executor interface
type StatelessExecutor struct {
	config  *Config
	runtime Runtime
	server  communication.ExecutorServer
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
//   - log: Logger instance
//   - kmsClient: KMS client for Type 1 recovery (can be nil for Type 0)
//   - enclaveHandle: Enclave handle for attestation (can be nil for Type 0)
func NewStatelessExecutor(
	config *Config,
	runtime Runtime,
	server communication.ExecutorServer,
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
		MsgToSignBuilder: msgBuilder,
		log:              log,
		kmsClient:        kmsClient,
		enclaveHandle:    enclaveHandle,
	}
	// Set this executor as the request handler
	executor.server.SetRequestHandler(executor)
	// Set the connection handler to perform handshake
	executor.server.SetConnectionHandler(executor.handleNewConnection)

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

	return e.server.Start(ctx, "Executor")
}

// Stop stops the executor server
func (e *StatelessExecutor) Stop() error {
	e.log.Info("Executor: Stopping stateless executor")

	err := e.server.Stop()
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

// validateRequest performs common sanity checks on request fields.
// These checks should have been done by the manager or by the contract, so if they fail here
// it means that there is a tampering or a bug in the manager. For all these checks we do not
// set the request as failed because it is not a fault of the request, but rather a fault of
// the manager or a tampering attempt.
//
// TODO ST As a safer alternative, in order to check that the request was not tampered, we could use the requestId,
// that is a hash of the request parameters + an index from the ProcessorEndpoint.
// The index could be added to the request and then the requestId could be reconstructed here and compared with the
// one in the request. This would ensure that the request parameters were not tampered, otherwise the hash would not match.
func (e *StatelessExecutor) validateRequest(req *common.Request) error {
	if req.ProtocolVersion != admittedProtocolVersion {
		return fmt.Errorf("protocol version %d is not admitted", req.ProtocolVersion)
	}
	// TRUSTPROCESS requests are enqueued by the on-chain trigger contract with
	// maxFeeValue = 0 (set by _enqueueTrustedRequest in ProcessorEndpoint).
	// They are authenticated on-chain and do not go through the normal user
	// fee path, so the minimum fee check must be skipped for them.
	if req.RequestType != common.TrustProcess && req.MaxFeeValue.ToInt().Cmp(e.config.MinFeePerRequest) < 0 {
		return fmt.Errorf("request fee is below minimum fee")
	}
	return nil
}

// HandleProcessRequest is the main workflow: decrypt state -> invoke WASM -> encrypt new state -> sign -> respond
// Returns UpdatePayload, ApplicationState, optional DeanonymizationReport, and error
func (e *StatelessExecutor) HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
	e.log.Info("Executor: Processing request %s for application %d", req.RequestID, req.ApplicationID)

	// Also validated inside executeRequest; checked here first so tampered
	// fields fail before any state access.
	if err := e.validateRequest(req); err != nil {
		return nil, nil, nil, err
	}

	if req.RequestType != common.Process && req.RequestType != common.AssociateKey && req.RequestType != common.Deanonymize && req.RequestType != common.TrustProcess {
		return nil, nil, nil, fmt.Errorf("unsupported request type: %s", req.RequestType)
	}

	// State presence, wasm presence, state root and wasm fingerprint checks
	appData, _, err := e.loadVerifiedAppData(appState, wasmModule, req.ApplicationID, fmt.Sprintf("request %s", req.RequestID))
	if err != nil {
		return nil, nil, nil, err
	}

	// The manager pairs request and state by applicationId; a mismatch is
	// evidence of a bug or tampering.
	if req.ApplicationID != appState.ApplicationID {
		return nil, nil, nil, fmt.Errorf("request applicationId %d does not match state applicationId %d", req.ApplicationID, appState.ApplicationID)
	}

	// Execute the request against the decrypted app data
	outcome, err := e.executeRequest(ctx, req, appData, appState.StateRoot, wasmModule)
	if err != nil {
		return nil, nil, nil, err
	}

	// Sign the payload (success or soft-failure error payload) to produce the attestation
	signature, err := e.signUpdatePayload(outcome.payload)
	if err != nil {
		e.log.Error("Executor: error signing update payload: %v", err)
		return nil, nil, nil, err
	}
	outcome.payload.Signature = signature

	// Soft failure: signed error payload, state unchanged, nothing to store
	if outcome.payload.ErrorCode != 0 {
		return outcome.payload, nil, nil, nil
	}

	// Success outcomes always carry the serialized app data (Serialize never
	// returns empty bytes); a nil value here is a bug in executeRequest.
	if outcome.newSerialized == nil {
		return nil, nil, nil, fmt.Errorf("bug: success outcome without serialized app data for request %s", req.RequestID)
	}

	// Encrypt the new app data
	newApplicationState, err := e.buildEncryptedApplicationState(req.ApplicationID, outcome.payload.NewStateRoot, outcome.newSerialized)
	if err != nil {
		return nil, nil, nil, err
	}

	return outcome.payload, newApplicationState, outcome.report, nil
}

// HandleDeployApp implements the RequestHandler interface
func (e *StatelessExecutor) HandleDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	e.log.Info("Executor: Deploying application for request %s", req.RequestID)

	if err := e.validateRequest(req); err != nil {
		return nil, nil, err
	}

	if req.RequestType != common.Deploy {
		return nil, nil, fmt.Errorf("request type %s is not deploy", req.RequestType)
	}

	if appState != nil {
		// This means that the application was already deployed. This is the request fault, so we set the request as failed with the appropriate error code and message.
		errorPayload, err := e.processErrorResponse(req,
			appState.StateRoot,
			apperrors.New(apperrors.CodeApplicationAlreadyDeployed,
				fmt.Sprintf("application %s was already deployed", req.ApplicationID)))
		return errorPayload, nil, err
	}

	descriptor, err := common.DecodeDeployDescriptorStrict(req.Payload)
	if err != nil {
		errorPayload, respErr := e.errorResponse(req, emptyStateRoot,
			apperrors.CodeInternalFallback,
			deployDescriptorFailureMsg, err)
		return errorPayload, nil, respErr
	}

	if len(wasmModule) == 0 {
		errorPayload, err := e.processErrorResponse(
			req,
			emptyStateRoot,
			apperrors.New(apperrors.CodeWasmModuleEmpty, "wasm module is empty"),
		)
		return errorPayload, nil, err
	}

	moduleSHA := sha256.Sum256(wasmModule)
	if got := hex.EncodeToString(moduleSHA[:]); got != descriptor.WasmSHA256 {
		return nil, nil, fmt.Errorf("wasm fingerprint mismatch")
	}

	// Deploy the module with constructor params and get initial state
	initialAppState, fuel, err := e.runtime.Deploy(ctx, req.ApplicationID, descriptor.ConstructorParams, wasmModule)
	if err != nil {
		// The runtime's error message carries the specific failure mode
		// (compile failure, invalid guest result, JSON parse error, etc.).
		// errorResponse logs it locally and folds it into the signed-payload
		// message so downstream consumers see more than just the error code.
		errorPayload, respErr := e.errorResponse(req, emptyStateRoot,
			apperrors.CodeFailedLoadingOrGettingModule,
			deployLoadFailureMsg, err)
		return errorPayload, nil, respErr
	}

	// Check if there is enough ETH to cover the fuel costs
	// TODO make a helper function?
	applicationFee := new(big.Int).Mul(fuel, e.config.FuelPricePerUnit)

	// Application fee must be minimum fee at least
	if applicationFee.Cmp(e.config.MinFeePerRequest) < 0 {
		applicationFee = new(big.Int).Set(e.config.MinFeePerRequest)
	}

	if req.MaxFeeValue.ToInt().Cmp(applicationFee) < 0 {
		errorPayload, err := e.processErrorResponse(req,
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
	initialAppData.SetWasmFingerprint(moduleSHA)

	//serialize the new app data
	initialAppDataBytes, err := initialAppData.Serialize()
	if err != nil {
		errorPayload, respErr := e.errorResponse(req, emptyStateRoot,
			apperrors.CodeAppDataSerializationFailure,
			"failed to serialize new app data", err)
		return errorPayload, nil, respErr
	}
	// Create app data root hash
	initialAppDataRoot := sha256.Sum256(initialAppDataBytes)

	// Encrypt the initial state and build the application state
	newAppState, err := e.buildEncryptedApplicationState(req.ApplicationID, initialAppDataRoot, initialAppDataBytes)
	if err != nil {
		return nil, nil, err
	}
	e.log.Info("Executor: Successfully encrypted initial app data")

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
	appData, _, err := e.decryptAndDeserializeState(encState)
	return appData, err
}

// buildEncryptedApplicationState encrypts serialized app data with the TEE
// state key and wraps it into the ApplicationState handed back to the manager.
func (e *StatelessExecutor) buildEncryptedApplicationState(appID common.ApplicationIdType, stateRoot [32]byte, serialized []byte) (*common.ApplicationState, error) {
	encryptedState, err := crypto.EncryptWithAES(e.keySet.StateKey, serialized)
	if err != nil {
		e.log.Error("Executor: error encrypting app data for application %d: %v", appID, err)
		return nil, err
	}

	return &common.ApplicationState{
		ApplicationID:  appID,
		StateRoot:      stateRoot,
		EncryptedState: encryptedState,
	}, nil
}

// decryptAndDeserializeState decrypts the state, verifies it against the state
// root and parses it. It returns the parsed app data together with the
// decrypted bytes it was parsed from.
func (e *StatelessExecutor) decryptAndDeserializeState(encState *common.ApplicationState) (*appdata.AppData, []byte, error) {
	// Decrypt the encrypted state
	decryptedState, err := e.DecryptState(encState.EncryptedState, e.keySet.StateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt state: %w", err)
	}

	// Verify state consistency
	hash := sha256.Sum256(decryptedState)
	if hash != encState.StateRoot {
		return nil, nil, fmt.Errorf("state root mismatch: got %x, want %x", hash, encState.StateRoot)
	}

	// Parse the app data
	appData, err := appdata.DeserializeAppData(decryptedState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize state: %w", err)
	}

	return appData, decryptedState, nil
}

// loadVerifiedAppData runs the checks every state-consuming flow needs before
// touching the WASM runtime: state presence, wasm presence, decryption, state
// root verification, deserialization and wasm fingerprint verification.
// It returns the parsed app data and the decrypted bytes it was parsed from, so
// callers that need to re-deserialize the same state (batch execution) can
// reuse them.
//
// applicationID and logCtx are used for diagnostics only; logCtx identifies the
// caller (e.g. "request <id>" or "batch").
func (e *StatelessExecutor) loadVerifiedAppData(appState *common.ApplicationState, wasmModule []byte, applicationID common.ApplicationIdType, logCtx string) (*appdata.AppData, []byte, error) {
	// App existence is validated on-chain (validApplicationId modifier in
	// ProcessorEndpoint), so a missing state here means tampering or
	// manager-side state loss.
	if appState == nil {
		e.log.Error("Executor: state not found for application %d (%s): app existence is enforced on-chain, check the manager DB", applicationID, logCtx)
		return nil, nil, fmt.Errorf("state not found for application %d", applicationID)
	}

	if len(wasmModule) == 0 {
		return nil, nil, fmt.Errorf("empty wasm module for application %d", applicationID)
	}

	appData, serializedState, err := e.decryptAndDeserializeState(appState)
	if err != nil {
		return nil, nil, err
	}

	expectedWasmFingerprint := appData.GetWasmFingerprint()
	currentWasmFingerprint := sha256.Sum256(wasmModule)
	if currentWasmFingerprint != expectedWasmFingerprint {
		e.log.Warn(
			"Executor: Wasm fingerprint mismatch for application %d (%s) (gotPrefix=%x expectedPrefix=%x)",
			applicationID,
			logCtx,
			currentWasmFingerprint[:4],
			expectedWasmFingerprint[:4],
		)
		return nil, nil, fmt.Errorf("wasm fingerprint mismatch for application %d", applicationID)
	}

	return appData, serializedState, nil
}

// buildUnsignedErrorPayload creates an unsigned error payload for failed requests.
// The payload has: prevStateRoot == newStateRoot (unchanged), empty events/withdrawals,
// refund = MaxFeeValue - MinFee, and the error code/message from the failure.
func (e *StatelessExecutor) buildUnsignedErrorPayload(req *common.Request, stateRoot [32]byte, failure *apperrors.RequestFailure) *common.UpdatePayload {
	// Truncate error message if longer than 100 characters
	errorMsg := failure.ExternalMessage()
	if len(errorMsg) > 100 {
		errorMsg = errorMsg[:100]
	}

	// Create the error payload with state unchanged
	return &common.UpdatePayload{
		ApplicationID:  req.ApplicationID,
		RequestID:      req.RequestID,
		PrevStateRoot:  stateRoot,
		NewStateRoot:   stateRoot,        // State unchanged on error
		Events:         nil,              // Empty events on error
		AppEvents:      nil,              // Empty appevents on error
		Withdrawals:    nil,              // Empty withdrawals on error
		RefundAmount:   req.MaxFeeValue,  // Refund and fees are calculated on chain
		ApplicationFee: common.NewBig(0), // No application fee on error, but the actual fee handling is done on chain
		ErrorCode:      failure.Category(),
		ErrorMsg:       errorMsg,
	}
}

// buildErrorPayload creates a signed error payload for failed requests.
func (e *StatelessExecutor) buildErrorPayload(req *common.Request, stateRoot [32]byte, failure *apperrors.RequestFailure) (*common.UpdatePayload, error) {
	updatePayload := e.buildUnsignedErrorPayload(req, stateRoot, failure)

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

// errorResponse builds a signed error response from a failure code + base
// message, optionally folding an underlying cause into both the Executor log
// and the payload message. It exists so call sites don't have to remember to
// log AND wrap the cause every time they construct an apperrors.RequestFailure.
//
// When cause is nil, this is equivalent to processErrorResponse with
// apperrors.New(code, baseMsg). When cause is non-nil, it:
//  1. logs the cause at Error level with request/app context, and
//  2. appends ": <cause>" to baseMsg in the signed payload so downstream
//     consumers (wallet, subgraph) see the specific failure, not just the
//     generic error code category.
func (e *StatelessExecutor) errorResponse(req *common.Request, stateRoot [32]byte, code apperrors.FailureCode, baseMsg string, cause error) (*common.UpdatePayload, error) {
	msg := baseMsg
	if cause != nil {
		e.log.Error("Executor: request %s (app %d) failed (%s): %v",
			req.RequestID, req.ApplicationID, code.Code, cause)
		msg = fmt.Sprintf("%s: %v", baseMsg, cause)
	}
	return e.processErrorResponse(req, stateRoot, apperrors.New(code, msg))
}

// signUpdatePayload signs the update payload to produce an attestation
func (e *StatelessExecutor) signUpdatePayload(payload *common.UpdatePayload) ([]byte, error) {
	// Serialize the payload for signing
	hash, err := e.BuildMsgHash(payload)

	if err != nil {
		return nil, fmt.Errorf("failed to create message to sign: %w", err)
	}

	return e.signHash(hash)
}

// signHash signs a prepared message hash with the TEE signing key and adjusts
// the recovery byte to the 27/28 format expected by Solidity's ecrecover.
func (e *StatelessExecutor) signHash(hash []byte) ([]byte, error) {
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

func (e *StatelessExecutor) encryptEvents(ctx context.Context, events []common.PlainEvent, appId common.ApplicationIdType, key *cryptotypes.PrivateKeyP521, server communication.ExecutorServer, keyStore appdata.KeyStore, seedStore appdata.SeedStore) ([]common.Event, *apperrors.RequestFailure, error) {
	if len(events) == 0 {
		return nil, nil, nil // No events to encrypt
	}
	encryptedEvents := make([]common.Event, len(events))

	for i, event := range events {
		// retrieve user Secp521r1_PubKey
		userKey, exists := keyStore[event.UserID]

		if !exists {
			return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, "no Secp521r1_PubKey found"), nil
		}
		// Encrypt the event data
		encryptedData, err := crypto.Encrypt(key, userKey, event.Data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt deanonymization report: %w", err)
		}

		// Use a privacy-preserving subtype if a seed is registered for this user;
		// otherwise fall back to the WASM-provided subtype.
		eventSubType := event.EventSubType
		if seed, hasSeed := seedStore[event.UserID]; hasSeed {
			privSubtype, err := GenerateRandomSubtype(seed, DefaultSubtypeN)
			if err != nil {
				e.log.Warn("Executor: failed to generate random subtype for user %s, falling back to WASM subtype: %v", event.UserID, err)
			} else {
				eventSubType = privSubtype
			}
		}

		// Create the encrypted event
		encryptedEvents[i] = common.Event{
			ApplicationID: appId,
			UserID:        event.UserID,
			EventSubType:  eventSubType,
			EncryptedData: encryptedData,
		}
	}

	e.log.Info("Executor: Successfully encrypted %d events", len(events))
	return encryptedEvents, nil, nil
}

func (e *StatelessExecutor) encryptDeanonymizationReport(applicationId common.ApplicationIdType, requestId common.RequestIdType, key *cryptotypes.PrivateKeyP521, requester ethCommon.Address, reportData []byte, keyStore appdata.KeyStore) ([]byte, *apperrors.RequestFailure, error) {
	if len(reportData) == 0 {
		return nil, apperrors.New(apperrors.CodeNoReportDataFound, "no report data found"), nil
	}

	// retrieve user Secp521r1_PubKey
	requesterPublicKey, exists := keyStore[requester]
	if !exists {
		return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, "no Secp521r1_PubKey found"), nil
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
		return nil, apperrors.New(apperrors.CodeJsonMarshalError, "failed to marshal updated deanonymization report"), nil
	}

	// Encrypt the report data
	encryptedReport, err := crypto.Encrypt(key, requesterPublicKey, updatedReportData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt deanonymization report: %w", err)
	}

	return encryptedReport, nil, nil
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
		return nil, apperrors.New(apperrors.CodePubKeyNotRegistered, "no Secp521r1_PubKey found")
	}

	decryptedPayload, err := crypto.Decrypt(userKey, decryptionKey, payload)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeWrongKey, "payload decryption failed")
	}

	e.log.Info("Executor: Successfully decrypted request payload")
	return decryptedPayload, nil
}

// HandleAdminCommand handles admin commands forwarded from the manager through
// the communication channel. It dispatches on cmdType to the appropriate handler.
func (e *StatelessExecutor) HandleAdminCommand(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
	switch cmdType {
	case admin.AdminCmdKeyAttestation:
		attestation, err := e.CreateKeyAttestation(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := json.Marshal(attestation)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal attestation response: %w", err)
		}
		return resp, nil

	case admin.AdminCmdSetLogLevel:
		var req admin.SetLogLevelRequest
		if data != nil && string(data) != "null" {
			if err := json.Unmarshal(data, &req); err != nil {
				return nil, fmt.Errorf("invalid request data: %w", err)
			}
		}
		result, err := admin.HandleSetLogLevel(e.log, "Executor", req.Level)
		if err != nil {
			return nil, err
		}
		resp, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		return resp, nil

	case admin.AdminCmdGetLogLevel:
		level, err := admin.HandleGetLogLevel(e.log, "Executor")
		if err != nil {
			return nil, err
		}
		resp, err := json.Marshal(level)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		return resp, nil

	case admin.AdminCmdGetVersion:
		resp, err := json.Marshal(version.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal version: %w", err)
		}
		return resp, nil

	case admin.AdminCmdSetWasmCacheSize:
		var req admin.SetWasmCacheSizeRequest
		if data != nil && string(data) != "null" {
			if err := json.Unmarshal(data, &req); err != nil {
				return nil, fmt.Errorf("invalid request data: %w", err)
			}
		}
		if req.MaxCachedModules < 0 {
			return nil, fmt.Errorf("maxCachedModules must be >= 0 (0 = unlimited)")
		}
		cc, ok := e.runtime.(moduleCacheController)
		if !ok {
			return nil, fmt.Errorf("runtime does not support module caching")
		}
		cc.SetMaxCachedModules(req.MaxCachedModules)
		resp, err := json.Marshal(admin.SetWasmCacheSizeResponse{MaxCachedModules: req.MaxCachedModules})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		return resp, nil

	case admin.AdminCmdGetWasmCacheSize:
		cc, ok := e.runtime.(moduleCacheController)
		if !ok {
			return nil, fmt.Errorf("runtime does not support module caching")
		}
		maxCached := cc.GetMaxCachedModules()
		resp, err := json.Marshal(admin.GetWasmCacheSizeResponse{MaxCachedModules: maxCached})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unsupported admin command type: %s", cmdType)
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

	signerAddress, err := hex.DecodeString(keySet.SigningKey.PublicKey().Address()[2:]) // remove 0x prefix
	if err != nil {
		e.log.Error("Executor: error while decoding signer address: %v", err)
		return nil, fmt.Errorf("failed to generate attestation")
	}

	doc, err := nsmutil.AttestWithSession(nsmSessionOpener, keySet.CommunicationKey.PublicKey().Bytes(), signerAddress)
	if err != nil {
		e.log.Warn("Executor: failed to generate attestation: %v", err)
		return nil, fmt.Errorf("failed to generate attestation: %w", err)
	}

	return doc, nil

}
