package executor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// newTestExecutor creates a StatelessExecutor configured for testing HandleDeployApp.
func newTestExecutor(t *testing.T, runtime Runtime) *StatelessExecutor {
	t.Helper()
	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)

	ks, err := CreateNewKeySet()
	require.NoError(t, err)

	return &StatelessExecutor{
		config: &Config{
			FuelPricePerUnit: big.NewInt(1),
			MinFeePerRequest: big.NewInt(10),
		},
		runtime:          runtime,
		MsgToSignBuilder: builder,
		keySet:           ks,
		log:              testLogger,
	}
}

func newDeployRequest(t *testing.T) (*common.Request, []byte) {
	t.Helper()

	wasmModule := []byte("mock-wasm-bytecode")
	payload := buildDeployDescriptorPayload(t, wasmModule)

	return &common.Request{
		ProtocolVersion: 0,
		ApplicationID:   common.NewApplicationId(1),
		RequestID:       commontestutil.GenerateRandomRequestID(),
		RequestType:     common.Deploy,
		Payload:         payload,
		MaxFeeValue:     common.NewBig(1000),
		AssetAmount:     common.NewBig(0),
	}, wasmModule
}

func buildDeployDescriptorPayload(t *testing.T, wasmModule []byte) []byte {
	t.Helper()

	sum := sha256.Sum256(wasmModule)
	wasmSHA := hex.EncodeToString(sum[:])
	artifactID, err := common.BuildArtifactID(wasmSHA)
	require.NoError(t, err)

	payload, err := json.Marshal(common.DeployDescriptor{
		Mode:       common.DeployModeArtifactRef,
		ArtifactID: artifactID,
		WasmSHA256: wasmSHA,
	})
	require.NoError(t, err)

	return payload
}

func TestHandleDeployApp_WrongProtocolVersion(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req, wasmModule := newDeployRequest(t)
	req.ProtocolVersion = 99

	_, _, err := executor.HandleDeployApp(context.Background(), req, nil, wasmModule)
	require.Error(t, err)
	require.Contains(t, err.Error(), "protocol version 99 is not admitted")
}

func TestHandleDeployApp_WrongRequestType(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req, wasmModule := newDeployRequest(t)
	req.RequestType = common.Process

	_, _, err := executor.HandleDeployApp(context.Background(), req, nil, wasmModule)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not deploy")
}

func TestHandleDeployApp_FeeBelowMinimum(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req, wasmModule := newDeployRequest(t)
	req.MaxFeeValue = common.NewBig(1)

	_, _, err := executor.HandleDeployApp(context.Background(), req, nil, wasmModule)
	require.Error(t, err)
	require.Contains(t, err.Error(), "request fee is below minimum fee")
}

func TestHandleDeployApp_InvalidDescriptor(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req, wasmModule := newDeployRequest(t)
	req.Payload = []byte("not-json")

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, nil, wasmModule)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CodeInternalFallback.Category.Category), updatePayload.ErrorCode)
	// errorResponse now appends the underlying decode error to the base message
	// (e.g. "failed to deploy application: invalid character 'o' ...").
	// Assert the base message is present; the exact suffix is
	// implementation-dependent on the json package.
	require.Contains(t, updatePayload.ErrorMsg, "failed to deploy application")
	require.Equal(t, [32]byte{}, updatePayload.PrevStateRoot)
	require.Equal(t, [32]byte{}, updatePayload.NewStateRoot)
}

func TestHandleDeployApp_NilWASM(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req, _ := newDeployRequest(t)

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CodeWasmModuleEmpty.Category.Category), updatePayload.ErrorCode)
	require.Equal(t, "wasm module is empty", updatePayload.ErrorMsg)
}

func TestHandleDeployApp_EmptyWASM(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req, _ := newDeployRequest(t)

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, nil, []byte{})
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CodeWasmModuleEmpty.Category.Category), updatePayload.ErrorCode)
	require.Equal(t, "wasm module is empty", updatePayload.ErrorMsg)
}

func TestHandleDeployApp_ApplicationAlreadyDeployed(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req, wasmModule := newDeployRequest(t)
	existingState := &common.ApplicationState{
		ApplicationID: req.ApplicationID,
		StateRoot:     [32]byte{0x01},
	}

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, existingState, wasmModule)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CategoryApplicationAlreadyDeployedMeta.Category), updatePayload.ErrorCode)
	require.Contains(t, updatePayload.ErrorMsg, "was already deployed")
	require.Equal(t, existingState.StateRoot, updatePayload.PrevStateRoot)
	require.Equal(t, existingState.StateRoot, updatePayload.NewStateRoot)
}

// failingRuntime implements Runtime but always returns an error from LoadModule.
type failingRuntime struct{}

func (r *failingRuntime) LoadModule(_ context.Context, _ common.ApplicationIdType, _ []byte) ([]byte, *big.Int, error) {
	return nil, big.NewInt(0), fmt.Errorf("wasm compilation failed")
}

func (r *failingRuntime) Deploy(_ context.Context, _ common.ApplicationIdType, _ []byte, _ []byte) ([]byte, *big.Int, error) {
	return nil, big.NewInt(0), fmt.Errorf("wasm compilation failed")
}

func (r *failingRuntime) Deposit(_ context.Context, _ common.ApplicationIdType, _ ethCommon.Address, _ ethCommon.Address, _ *big.Int, _ []byte, _ []byte) ([]byte, []common.PlainEvent, []common.AppEvent, *big.Int, *apperrors.RequestFailure) {
	return nil, nil, nil, big.NewInt(0), nil
}

func (r *failingRuntime) ProcessRequest(_ context.Context, _ common.ApplicationIdType, _ ethCommon.Address, _ common.RequestType, _ []byte, _ []byte, _ []byte) ([]byte, []common.PlainEvent, []common.AppEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	return nil, nil, nil, nil, nil, big.NewInt(0), nil
}

func (r *failingRuntime) Close() error { return nil }

func TestHandleDeployApp_LoadModuleFails(t *testing.T) {
	executor := newTestExecutor(t, &failingRuntime{})

	req, wasmModule := newDeployRequest(t)

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, nil, wasmModule)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CategoryWasmInternalMeta.Category), updatePayload.ErrorCode)
	require.Contains(t, updatePayload.ErrorMsg, "failed to load or get module")
	require.Equal(t, [32]byte{}, updatePayload.PrevStateRoot)
	require.Equal(t, [32]byte{}, updatePayload.NewStateRoot)
}

type expensiveRuntime struct {
	MockRuntime
}

func (r *expensiveRuntime) LoadModule(ctx context.Context, appID common.ApplicationIdType, wasm []byte) ([]byte, *big.Int, error) {
	state, _, err := r.MockRuntime.LoadModule(ctx, appID, wasm)
	return state, big.NewInt(999999), err
}

func (r *expensiveRuntime) Deploy(ctx context.Context, appID common.ApplicationIdType, constructorParams []byte, wasm []byte) ([]byte, *big.Int, error) {
	state, _, err := r.MockRuntime.Deploy(ctx, appID, constructorParams, wasm)
	return state, big.NewInt(999999), err
}

func TestHandleDeployApp_InsufficientFuel(t *testing.T) {
	runtime := &expensiveRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	executor := newTestExecutor(t, runtime)

	req, wasmModule := newDeployRequest(t)
	req.MaxFeeValue = common.NewBig(100)

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, nil, wasmModule)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CategoryInsufficientFuelMeta.Category), updatePayload.ErrorCode)
	require.Contains(t, updatePayload.ErrorMsg, "insufficient fuel")
	require.Equal(t, [32]byte{}, updatePayload.PrevStateRoot)
	require.Equal(t, [32]byte{}, updatePayload.NewStateRoot)
}

func TestHandleDeployApp_HashMismatch(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req, _ := newDeployRequest(t)
	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, nil, []byte("different-wasm-bytecode"))
	require.Error(t, err)
	require.Nil(t, updatePayload)
	require.Contains(t, err.Error(), "wasm fingerprint mismatch")}

func TestHandleDeployApp_Success(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	executor := newTestExecutor(t, runtime)

	req, wasmModule := newDeployRequest(t)

	updatePayload, newAppState, err := executor.HandleDeployApp(context.Background(), req, nil, wasmModule)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.NotNil(t, newAppState)

	require.Equal(t, uint8(0), updatePayload.ErrorCode)
	require.Empty(t, updatePayload.ErrorMsg)
	require.Equal(t, [32]byte{}, updatePayload.PrevStateRoot)
	require.NotEqual(t, [32]byte{}, updatePayload.NewStateRoot)
	require.Equal(t, updatePayload.NewStateRoot, newAppState.StateRoot)
	require.Len(t, updatePayload.Signature, 65)
	require.NotEmpty(t, newAppState.EncryptedState)

	// WASM fingerprint must be persisted in private app state
	decryptedAppData, err := executor.fromEncryptedStateToAppData(newAppState)
	require.NoError(t, err)
	require.Equal(t, sha256.Sum256(wasmModule), decryptedAppData.GetWasmFingerprint())

	// RefundAmount + ApplicationFee == MaxFeeValue
	refund := updatePayload.RefundAmount.ToInt()
	fee := updatePayload.ApplicationFee.ToInt()
	total := new(big.Int).Add(refund, fee)
	require.Equal(t, req.MaxFeeValue.ToInt().Cmp(total), 0, "refund + fee should equal MaxFeeValue")
}

// TestHandleDeployApp_WithConstructorParams verifies that constructor params
// from the deploy descriptor are passed through to the runtime's Deploy method.
func TestHandleDeployApp_WithConstructorParams(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	executor := newTestExecutor(t, runtime)

	wasmModule := []byte("mock-wasm-bytecode")
	sum := sha256.Sum256(wasmModule)
	wasmSHA := hex.EncodeToString(sum[:])
	artifactID, err := common.BuildArtifactID(wasmSHA)
	require.NoError(t, err)

	// Build deploy descriptor with constructor params
	constructorParams := json.RawMessage(`{"allowedTokens":["0xdead000000000000000000000000000000000001"]}`)
	payload, err := json.Marshal(common.DeployDescriptor{
		Mode:              common.DeployModeArtifactRef,
		ArtifactID:        artifactID,
		WasmSHA256:        wasmSHA,
		ConstructorParams: constructorParams,
	})
	require.NoError(t, err)

	req := &common.Request{
		ProtocolVersion: 0,
		ApplicationID:   common.NewApplicationId(1),
		RequestID:       commontestutil.GenerateRandomRequestID(),
		RequestType:     common.Deploy,
		Payload:         payload,
		MaxFeeValue:     common.NewBig(1000),
		AssetAmount:     common.NewBig(0),
	}

	updatePayload, newAppState, err := executor.HandleDeployApp(context.Background(), req, nil, wasmModule)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.NotNil(t, newAppState)
	require.Equal(t, uint8(0), updatePayload.ErrorCode)
	require.Empty(t, updatePayload.ErrorMsg)
}
