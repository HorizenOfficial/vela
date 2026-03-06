package executor

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/appdata"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/crypto"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// newProcessRequest creates a valid Process request with sensible defaults.
func newProcessRequest() *common.Request {
	return &common.Request{
		ProtocolVersion: 0,
		ApplicationID:   common.NewApplicationId(1),
		RequestID:       commontestutil.GenerateRandomRequestID(),
		RequestType:     common.Process,
		Payload:         nil, // set per test
		MaxFeeValue:     common.NewBig(1000),
		DepositAmount:   common.NewBig(0),
	}
}

// buildEncryptedAppState creates a valid encrypted ApplicationState using the executor's keys.
// Optionally adds a user key to the keystore.
func buildEncryptedAppState(t *testing.T, exec *StatelessExecutor, userAddr *ethCommon.Address, userKey *cryptotypes.PublicKeyP521, wasmModule ...[]byte) *common.ApplicationState {
	t.Helper()

	moduleForFingerprint := []byte("wasm")
	if len(wasmModule) > 0 {
		moduleForFingerprint = wasmModule[0]
	}

	appState := []byte(`{"appId":1,"accounts":{},"nonce":0}`)
	ad := appdata.NewAppData(appState)
	ad.SetWasmFingerprint(sha256.Sum256(moduleForFingerprint))
	if userAddr != nil && userKey != nil {
		ad.AddKey(*userAddr, *userKey)
	}

	serialized, err := ad.Serialize()
	require.NoError(t, err)

	encrypted, err := crypto.EncryptWithAES(exec.keySet.StateKey, serialized)
	require.NoError(t, err)

	stateRoot := sha256.Sum256(serialized)
	return &common.ApplicationState{
		ApplicationID:  common.NewApplicationId(1),
		StateRoot:      stateRoot,
		EncryptedState: encrypted,
	}
}

// encryptPayload encrypts a plaintext payload using the user's private key and the executor's communication public key.
func encryptPayload(t *testing.T, userPrivKey *cryptotypes.PrivateKeyP521, execPubKey *cryptotypes.PublicKeyP521, plaintext []byte) []byte {
	t.Helper()
	encrypted, err := crypto.Encrypt(userPrivKey, execPubKey, plaintext)
	require.NoError(t, err)
	return encrypted
}

func TestHandleProcessRequest_WrongProtocolVersion(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()
	req.ProtocolVersion = 42

	_, _, _, err := exec.HandleProcessRequest(context.Background(), req, nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "protocol version 42 is not admitted")
}

func TestHandleProcessRequest_WrongApplicationID(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()
	req.ApplicationID = common.NewApplicationId(999)

	_, _, _, err := exec.HandleProcessRequest(context.Background(), req, nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "application id 999 is not admitted")
}

func TestHandleProcessRequest_FeeBelowMinimum(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()
	req.MaxFeeValue = common.NewBig(1) // below MinFeePerRequest=10

	_, _, _, err := exec.HandleProcessRequest(context.Background(), req, nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "request fee is below minimum fee")
}

func TestHandleProcessRequest_NilAppState(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryAppNotDeployedMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "state not found for application")
	require.Equal(t, [32]byte{}, payload.PrevStateRoot)
	require.Equal(t, [32]byte{}, payload.NewStateRoot)
}

func TestHandleProcessRequest_UnsupportedRequestType(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()
	req.RequestType = common.Deploy // Deploy is not handled by HandleProcessRequest

	appState := buildEncryptedAppState(t, exec, nil, nil)

	_, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported request type")
}

func TestHandleProcessRequest_DecryptStateFails(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()

	// Provide garbage encrypted state that can't be decrypted
	appState := &common.ApplicationState{
		ApplicationID:  common.NewApplicationId(1),
		StateRoot:      [32]byte{0x01},
		EncryptedState: []byte("not-valid-encrypted-data"),
	}

	_, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decrypt state")
}

func TestHandleProcessRequest_StateRootMismatch(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()

	appState := buildEncryptedAppState(t, exec, nil, nil)
	// Corrupt the state root so hash check fails
	appState.StateRoot = [32]byte{0xFF}

	_, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "state root mismatch")
}

type countingRuntime struct {
	MockRuntime
	depositCalls int
	processCalls int
}

func (r *countingRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, depositAmount *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, *big.Int, *apperrors.RequestFailure) {
	r.depositCalls++
	return r.MockRuntime.Deposit(ctx, appId, sender, depositAmount, state, wasm)
}

func (r *countingRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	r.processCalls++
	return r.MockRuntime.ProcessRequest(ctx, appId, sender, requestType, payload, state, wasm)
}

func TestHandleProcessRequest_EmptyWasmModuleReturnsSignedError(t *testing.T) {
	runtime := &countingRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)

	req := newProcessRequest()
	req.Payload = []byte("encrypted-stuff")
	req.Sender = ethCommon.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	appState := buildEncryptedAppState(t, exec, nil, nil)

	payload, newAppState, report, err := exec.HandleProcessRequest(context.Background(), req, appState, nil)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryWasmInternalMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "wasm module is empty")
	require.Equal(t, appState.StateRoot, payload.PrevStateRoot)
	require.Equal(t, appState.StateRoot, payload.NewStateRoot)
	require.Len(t, payload.Signature, 65)
	require.Nil(t, newAppState)
	require.Nil(t, report)
	require.Equal(t, 0, runtime.depositCalls)
	require.Equal(t, 0, runtime.processCalls)
}

func TestHandleProcessRequest_WasmFingerprintMismatchReturnsSignedError(t *testing.T) {
	runtime := &countingRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)

	req := newProcessRequest()
	req.Payload = []byte("encrypted-stuff")
	req.Sender = ethCommon.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")

	appState := buildEncryptedAppState(t, exec, nil, nil, []byte("expected-wasm"))

	payload, newAppState, report, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("tampered-wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryWasmInternalMeta.Category), payload.ErrorCode)
	require.Equal(t, "wasm fingerprint mismatch", payload.ErrorMsg)
	require.NotContains(t, payload.ErrorMsg, "got ")
	require.NotContains(t, payload.ErrorMsg, "want ")
	require.Equal(t, appState.StateRoot, payload.PrevStateRoot)
	require.Equal(t, appState.StateRoot, payload.NewStateRoot)
	require.Len(t, payload.Signature, 65)
	require.Nil(t, newAppState)
	require.Nil(t, report)
	require.Equal(t, 0, runtime.depositCalls)
	require.Equal(t, 0, runtime.processCalls)
}

func TestHandleProcessRequest_EmptyWasmModuleWithDepositReturnsSignedErrorBeforeDeposit(t *testing.T) {
	runtime := &countingRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)

	req := newProcessRequest()
	req.Sender = ethCommon.HexToAddress("0xABABABABABABABABABABABABABABABABABABABAB")
	req.DepositAmount = common.ToBig(big.NewInt(1000))

	appState := buildEncryptedAppState(t, exec, nil, nil)

	payload, newAppState, report, err := exec.HandleProcessRequest(context.Background(), req, appState, nil)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryWasmInternalMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "wasm module is empty")
	require.Equal(t, appState.StateRoot, payload.PrevStateRoot)
	require.Equal(t, appState.StateRoot, payload.NewStateRoot)
	require.Len(t, payload.Signature, 65)
	require.Nil(t, newAppState)
	require.Nil(t, report)
	require.Equal(t, 0, runtime.depositCalls)
	require.Equal(t, 0, runtime.processCalls)
}

func TestHandleProcessRequest_WasmFingerprintMismatchWithDepositReturnsSignedErrorBeforeDeposit(t *testing.T) {
	runtime := &countingRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)

	req := newProcessRequest()
	req.Sender = ethCommon.HexToAddress("0xCDCDCDCDCDCDCDCDCDCDCDCDCDCDCDCDCDCDCDCD")
	req.DepositAmount = common.ToBig(big.NewInt(1000))

	appState := buildEncryptedAppState(t, exec, nil, nil, []byte("expected-wasm"))

	payload, newAppState, report, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("tampered-wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryWasmInternalMeta.Category), payload.ErrorCode)
	require.Equal(t, "wasm fingerprint mismatch", payload.ErrorMsg)
	require.NotContains(t, payload.ErrorMsg, "got ")
	require.NotContains(t, payload.ErrorMsg, "want ")
	require.Equal(t, appState.StateRoot, payload.PrevStateRoot)
	require.Equal(t, appState.StateRoot, payload.NewStateRoot)
	require.Len(t, payload.Signature, 65)
	require.Nil(t, newAppState)
	require.Nil(t, report)
	require.Equal(t, 0, runtime.depositCalls)
	require.Equal(t, 0, runtime.processCalls)
}

func TestHandleProcessRequest_AssociateKey_FingerprintCheckHappensBeforePayloadValidation(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))

	req := newProcessRequest()
	req.RequestType = common.AssociateKey
	req.Payload = []byte("too-short")

	appState := buildEncryptedAppState(t, exec, nil, nil, []byte("expected-wasm"))

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("tampered-wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryWasmInternalMeta.Category), payload.ErrorCode)
	require.Equal(t, "wasm fingerprint mismatch", payload.ErrorMsg)
	require.NotContains(t, payload.ErrorMsg, "got ")
	require.NotContains(t, payload.ErrorMsg, "want ")
}

// failingDepositRuntime is a runtime where Deposit always fails.
type failingDepositRuntime struct {
	MockRuntime
}

func (r *failingDepositRuntime) Deposit(_ context.Context, _ common.ApplicationIdType, _ ethCommon.Address, _ *big.Int, _ []byte, _ []byte) ([]byte, []common.PlainEvent, *big.Int, *apperrors.RequestFailure) {
	return nil, nil, big.NewInt(10), apperrors.New(apperrors.CodeDepositFailed, "deposit rejected by app")
}

func TestHandleProcessRequest_DepositFails(t *testing.T) {
	runtime := &failingDepositRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)
	req := newProcessRequest()
	req.DepositAmount = common.ToBig(big.NewInt(1000))

	appState := buildEncryptedAppState(t, exec, nil, nil)

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryDepositFailedMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "deposit rejected by app")
}

// expensiveDepositRuntime returns high fuel from Deposit to trigger insufficient fuel after deposit.
type expensiveDepositRuntime struct {
	MockRuntime
}

func (r *expensiveDepositRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, depositAmount *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, *big.Int, *apperrors.RequestFailure) {
	newState, events, _, failure := r.MockRuntime.Deposit(ctx, appId, sender, depositAmount, state, wasm)
	return newState, events, big.NewInt(999999), failure
}

func TestHandleProcessRequest_InsufficientFuelAfterDeposit(t *testing.T) {
	runtime := &expensiveDepositRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)
	req := newProcessRequest()
	req.DepositAmount = common.ToBig(big.NewInt(1000))
	req.MaxFeeValue = common.NewBig(100) // less than 999999 * 1

	sender := ethCommon.HexToAddress("0x1111111111111111111111111111111111111111")
	req.Sender = sender

	appState := buildEncryptedAppState(t, exec, nil, nil)

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryInsufficientFuelMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "insufficient fuel")
}

func TestHandleProcessRequest_AssociateKey_InvalidPayloadLength(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()
	req.RequestType = common.AssociateKey
	req.Payload = []byte("too-short") // not 133 bytes

	appState := buildEncryptedAppState(t, exec, nil, nil)

	_, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid payload length")
}

func TestHandleProcessRequest_AssociateKey_InvalidKeyBytes(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()
	req.RequestType = common.AssociateKey
	req.Payload = make([]byte, 133) // 133 bytes but not a valid P521 key

	appState := buildEncryptedAppState(t, exec, nil, nil)

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryWrongKeySentMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "failed to parse keyP521")
}

func TestHandleProcessRequest_DecryptPayloadFails_NoPubKey(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	req := newProcessRequest()
	req.Payload = []byte("encrypted-stuff")
	req.Sender = ethCommon.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

	// Build state without the sender's key in the keystore
	appState := buildEncryptedAppState(t, exec, nil, nil)

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryPubKeyNotRegisteredMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "no Secp521r1_PubKey found")
}

func TestHandleProcessRequest_DecryptPayloadFails_WrongKey(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))

	sender := ethCommon.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	// Register the user's key in the keystore
	appState := buildEncryptedAppState(t, exec, &sender, userKey.PublicKey())

	// Encrypt payload with a DIFFERENT key (not the user's), causing decryption to fail
	otherKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	badEncrypted := encryptPayload(t, otherKey, exec.keySet.CommunicationKey.PublicKey(), []byte("hello"))

	req := newProcessRequest()
	req.Payload = badEncrypted
	req.Sender = sender

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryWrongKeySentMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "payload decryption failed")
}

// failingProcessRuntime is a runtime where ProcessRequest always fails.
type failingProcessRuntime struct {
	MockRuntime
}

func (r *failingProcessRuntime) ProcessRequest(_ context.Context, _ common.ApplicationIdType, _ ethCommon.Address, _ common.RequestType, _ []byte, _ []byte, _ []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	return nil, nil, nil, nil, big.NewInt(10), apperrors.New(apperrors.CodeRequestFuncFailed, "wasm execution failed")
}

func TestHandleProcessRequest_RuntimeProcessRequestFails(t *testing.T) {
	runtime := &failingProcessRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)

	sender := ethCommon.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appState := buildEncryptedAppState(t, exec, &sender, userKey.PublicKey())

	plainPayload, err := json.Marshal(testPayloadInstructions{Type: "transfer"})
	require.NoError(t, err)
	encrypted := encryptPayload(t, userKey, exec.keySet.CommunicationKey.PublicKey(), plainPayload)

	req := newProcessRequest()
	req.Payload = encrypted
	req.Sender = sender

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryRequestFuncFailedMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "wasm execution failed")
}

// reportOnProcessRuntime generates a report even on Process requests (not Deanonymize).
type reportOnProcessRuntime struct {
	MockRuntime
}

func (r *reportOnProcessRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	newState, events, withdrawals, _, fuel, failure := r.MockRuntime.ProcessRequest(ctx, appId, sender, requestType, payload, state, wasm)
	// Always return a report, even for non-Deanonymize requests
	return newState, events, withdrawals, []byte(`{"unexpected":"report"}`), fuel, failure
}

func TestHandleProcessRequest_UnexpectedReportForNonDeanonymize(t *testing.T) {
	runtime := &reportOnProcessRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)

	sender := ethCommon.HexToAddress("0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appState := buildEncryptedAppState(t, exec, &sender, userKey.PublicKey())

	// Empty payload (no instructions) to get a simple pass-through
	encrypted := encryptPayload(t, userKey, exec.keySet.CommunicationKey.PublicKey(), []byte("{}"))

	req := newProcessRequest()
	req.RequestType = common.Process
	req.Payload = encrypted
	req.Sender = sender

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryRequestFuncFailedMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "unexpected report")
}

// noReportDeanonymizeRuntime returns no report on Deanonymize requests.
type noReportDeanonymizeRuntime struct {
	MockRuntime
}

func (r *noReportDeanonymizeRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	newState, events, withdrawals, _, fuel, failure := r.MockRuntime.ProcessRequest(ctx, appId, sender, requestType, payload, state, wasm)
	// Force empty report even for Deanonymize
	return newState, events, withdrawals, nil, fuel, failure
}

func TestHandleProcessRequest_MissingReportForDeanonymize(t *testing.T) {
	runtime := &noReportDeanonymizeRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)

	sender := ethCommon.HexToAddress("0xEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appState := buildEncryptedAppState(t, exec, &sender, userKey.PublicKey())

	encrypted := encryptPayload(t, userKey, exec.keySet.CommunicationKey.PublicKey(), []byte("{}"))

	req := newProcessRequest()
	req.RequestType = common.Deanonymize
	req.Payload = encrypted
	req.Sender = sender

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryRequestFuncFailedMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "failed to generate report for Deanonymize")
}

// expensiveProcessRuntime returns high fuel from ProcessRequest.
type expensiveProcessRuntime struct {
	MockRuntime
}

func (r *expensiveProcessRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	newState, events, withdrawals, report, _, failure := r.MockRuntime.ProcessRequest(ctx, appId, sender, requestType, payload, state, wasm)
	return newState, events, withdrawals, report, big.NewInt(999999), failure
}

func TestHandleProcessRequest_InsufficientFuelAfterProcessing(t *testing.T) {
	runtime := &expensiveProcessRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	exec := newTestExecutor(t, runtime)

	sender := ethCommon.HexToAddress("0x1111111111111111111111111111111111111111")
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appState := buildEncryptedAppState(t, exec, &sender, userKey.PublicKey())

	encrypted := encryptPayload(t, userKey, exec.keySet.CommunicationKey.PublicKey(), []byte("{}"))

	req := newProcessRequest()
	req.Payload = encrypted
	req.Sender = sender
	req.MaxFeeValue = common.NewBig(100) // less than 999999 * 1

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(apperrors.CategoryInsufficientFuelMeta.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "insufficient fuel")
}

func TestHandleProcessRequest_AssociateKey_Success(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))

	sender := ethCommon.HexToAddress("0x2222222222222222222222222222222222222222")

	appState := buildEncryptedAppState(t, exec, nil, nil)

	// Generate a valid P521 key to associate
	keyToAssociate, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	pubKeyBytes := keyToAssociate.PublicKey().Bytes()
	require.Len(t, pubKeyBytes, 133, "P521 uncompressed public key should be 133 bytes")

	req := newProcessRequest()
	req.RequestType = common.AssociateKey
	req.Payload = pubKeyBytes
	req.Sender = sender

	payload, newAppState, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.NotNil(t, newAppState)
	require.Equal(t, uint8(0), payload.ErrorCode)
	require.Empty(t, payload.ErrorMsg)
	require.Len(t, payload.Signature, 65)
}

func TestHandleProcessRequest_Deanonymize_Success(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	exec := newTestExecutor(t, runtime)

	sender := ethCommon.HexToAddress("0x4444444444444444444444444444444444444444")
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appState := buildEncryptedAppState(t, exec, &sender, userKey.PublicKey())

	// Build a deanonymize payload and encrypt it
	deanonPayload, err := json.Marshal(testPayloadInstructions{
		Type:        "deanonymize",
		Deanonymize: &testDeanonymizeInstruction{Tag: "test-tag"},
	})
	require.NoError(t, err)
	encrypted := encryptPayload(t, userKey, exec.keySet.CommunicationKey.PublicKey(), deanonPayload)

	req := newProcessRequest()
	req.RequestType = common.Deanonymize
	req.Payload = encrypted
	req.Sender = sender

	payload, newAppState, report, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.NotNil(t, newAppState)
	require.NotNil(t, report, "deanonymization report should be returned")

	// Verify update payload
	require.Equal(t, uint8(0), payload.ErrorCode)
	require.Empty(t, payload.ErrorMsg)
	require.Equal(t, appState.StateRoot, payload.PrevStateRoot)
	require.NotEqual(t, [32]byte{}, payload.NewStateRoot)
	require.Equal(t, payload.NewStateRoot, newAppState.StateRoot)
	require.Len(t, payload.Signature, 65)
	require.NotEmpty(t, newAppState.EncryptedState)

	// RefundAmount + ApplicationFee == MaxFeeValue
	refund := payload.RefundAmount.ToInt()
	fee := payload.ApplicationFee.ToInt()
	total := new(big.Int).Add(refund, fee)
	require.Equal(t, 0, req.MaxFeeValue.ToInt().Cmp(total), "refund + fee should equal MaxFeeValue")

	// Verify deanonymization report fields
	require.Equal(t, req.ApplicationID, report.ApplicationID)
	require.Equal(t, req.RequestID, report.ReportID)
	require.Equal(t, sender, report.Authority)
	require.NotEmpty(t, report.EncryptedReport)
}

func TestHandleProcessRequest_Process_Success(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	exec := newTestExecutor(t, runtime)

	sender := ethCommon.HexToAddress("0x3333333333333333333333333333333333333333")
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appState := buildEncryptedAppState(t, exec, &sender, userKey.PublicKey())

	// Empty instructions - MockRuntime will just pass through with no state change
	encrypted := encryptPayload(t, userKey, exec.keySet.CommunicationKey.PublicKey(), []byte("{}"))

	req := newProcessRequest()
	req.Payload = encrypted
	req.Sender = sender

	payload, newAppState, report, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.NotNil(t, newAppState)
	require.Nil(t, report)

	require.Equal(t, uint8(0), payload.ErrorCode)
	require.Empty(t, payload.ErrorMsg)
	require.Equal(t, appState.StateRoot, payload.PrevStateRoot)
	require.NotEqual(t, [32]byte{}, payload.NewStateRoot)
	require.Equal(t, payload.NewStateRoot, newAppState.StateRoot)
	require.Len(t, payload.Signature, 65)
	require.NotEmpty(t, newAppState.EncryptedState)

	// RefundAmount + ApplicationFee == MaxFeeValue
	refund := payload.RefundAmount.ToInt()
	fee := payload.ApplicationFee.ToInt()
	total := new(big.Int).Add(refund, fee)
	require.Equal(t, 0, req.MaxFeeValue.ToInt().Cmp(total), "refund + fee should equal MaxFeeValue")
}
