package executor

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	commontestutil "github.com/horizen-pes/pkg/common/testutil"
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

// newDeployRequest creates a valid deploy request with sensible defaults.
func newDeployRequest() *common.Request {
	return &common.Request{
		ProtocolVersion: 0,
		ApplicationID:   common.NewApplicationId(1),
		RequestID:       commontestutil.GenerateRandomRequestID(),
		RequestType:     common.Deploy,
		Payload:         []byte("mock-wasm-bytecode"),
		MaxFeeValue:     common.NewBig(1000),
		DepositAmount:   common.NewBig(0),
	}
}

func TestHandleDeployApp_WrongProtocolVersion(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req := newDeployRequest()
	req.ProtocolVersion = 99

	_, _, err := executor.HandleDeployApp(context.Background(), req, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "protocol version 99 is not admitted")
}

func TestHandleDeployApp_WrongApplicationID(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req := newDeployRequest()
	req.ApplicationID = common.NewApplicationId(999)

	_, _, err := executor.HandleDeployApp(context.Background(), req, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "application id 999 is not admitted")
}

func TestHandleDeployApp_WrongRequestType(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req := newDeployRequest()
	req.RequestType = common.Process

	_, _, err := executor.HandleDeployApp(context.Background(), req, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not deploy")
}

func TestHandleDeployApp_FeeBelowMinimum(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req := newDeployRequest()
	req.MaxFeeValue = common.NewBig(1) // below MinFeePerRequest=10

	_, _, err := executor.HandleDeployApp(context.Background(), req, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "request fee is below minimum fee")
}

func TestHandleDeployApp_ApplicationAlreadyDeployed(t *testing.T) {
	executor := newTestExecutor(t, NewMockRuntime(testLogger))

	req := newDeployRequest()
	existingState := &common.ApplicationState{
		ApplicationID: req.ApplicationID,
		StateRoot:     [32]byte{0x01},
	}

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, existingState)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CategoryApplicationAlreadyDeployedMeta.Category), updatePayload.ErrorCode)
	require.Contains(t, updatePayload.ErrorMsg, "was already deployed")
	// State should be unchanged on error
	require.Equal(t, existingState.StateRoot, updatePayload.PrevStateRoot)
	require.Equal(t, existingState.StateRoot, updatePayload.NewStateRoot)
}

// failingRuntime implements Runtime but always returns an error from LoadModule.
type failingRuntime struct{}

func (r *failingRuntime) LoadModule(_ context.Context, _ common.ApplicationIdType, _ []byte) ([]byte, *big.Int, error) {
	return nil, big.NewInt(0), fmt.Errorf("wasm compilation failed")
}

func (r *failingRuntime) Deposit(_ context.Context, _ common.ApplicationIdType, _ ethCommon.Address, _ *big.Int, _ []byte, _ []byte) ([]byte, []common.PlainEvent, *big.Int, *apperrors.RequestFailure) {
	return nil, nil, big.NewInt(0), nil
}

func (r *failingRuntime) ProcessRequest(_ context.Context, _ common.ApplicationIdType, _ ethCommon.Address, _ common.RequestType, _ []byte, _ []byte, _ []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	return nil, nil, nil, nil, big.NewInt(0), nil
}

func (r *failingRuntime) Close() error { return nil }

func TestHandleDeployApp_LoadModuleFails(t *testing.T) {
	executor := newTestExecutor(t, &failingRuntime{})

	req := newDeployRequest()

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, nil)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CategoryWasmInternalMeta.Category), updatePayload.ErrorCode)
	require.Contains(t, updatePayload.ErrorMsg, "failed to load or get module")
	// State root should be empty (no previous state)
	require.Equal(t, [32]byte{}, updatePayload.PrevStateRoot)
	require.Equal(t, [32]byte{}, updatePayload.NewStateRoot)
}

// expensiveRuntime returns high fuel from LoadModule to trigger insufficient fuel errors.
type expensiveRuntime struct {
	MockRuntime
}

func (r *expensiveRuntime) LoadModule(ctx context.Context, appId common.ApplicationIdType, wasm []byte) ([]byte, *big.Int, error) {
	state, _, err := r.MockRuntime.LoadModule(ctx, appId, wasm)
	// Return very high fuel cost
	return state, big.NewInt(999999), err
}

func TestHandleDeployApp_InsufficientFuel(t *testing.T) {
	runtime := &expensiveRuntime{MockRuntime: *NewMockRuntime(testLogger)}
	executor := newTestExecutor(t, runtime)

	req := newDeployRequest()
	req.MaxFeeValue = common.NewBig(100) // way less than 999999 * 1 fuel price

	updatePayload, _, err := executor.HandleDeployApp(context.Background(), req, nil)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.Equal(t, uint8(apperrors.CategoryInsufficientFuelMeta.Category), updatePayload.ErrorCode)
	require.Contains(t, updatePayload.ErrorMsg, "insufficient fuel")
	require.Equal(t, [32]byte{}, updatePayload.PrevStateRoot)
	require.Equal(t, [32]byte{}, updatePayload.NewStateRoot)
}

func TestHandleDeployApp_Success(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	executor := newTestExecutor(t, runtime)

	req := newDeployRequest()

	updatePayload, newAppState, err := executor.HandleDeployApp(context.Background(), req, nil)
	require.NoError(t, err)
	require.NotNil(t, updatePayload)
	require.NotNil(t, newAppState)

	// No error on success
	require.Equal(t, uint8(0), updatePayload.ErrorCode)
	require.Empty(t, updatePayload.ErrorMsg)

	// PrevStateRoot should be empty for new app
	require.Equal(t, [32]byte{}, updatePayload.PrevStateRoot)
	// NewStateRoot should be non-empty
	require.NotEqual(t, [32]byte{}, updatePayload.NewStateRoot)
	// NewStateRoot should match the state returned
	require.Equal(t, updatePayload.NewStateRoot, newAppState.StateRoot)

	// Signature should be 65 bytes
	require.Len(t, updatePayload.Signature, 65)

	// Encrypted state should be non-empty
	require.NotEmpty(t, newAppState.EncryptedState)

	// RefundAmount + ApplicationFee == MaxFeeValue
	refund := updatePayload.RefundAmount.ToInt()
	fee := updatePayload.ApplicationFee.ToInt()
	total := new(big.Int).Add(refund, fee)
	require.Equal(t, req.MaxFeeValue.ToInt().Cmp(total), 0, "refund + fee should equal MaxFeeValue")
}
