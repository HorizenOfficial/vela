package executor

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// These pin the two-channel rule for the guest execution bound (see
// PLAN_epoch_interruption.md, decision D1):
//
//   - A guest that exceeded its wall-clock budget is a REQUEST FAILURE: signed and
//     submitted on-chain, so the on-chain queue head advances. Treating it as
//     transient would re-deliver the same request forever and let one
//     non-terminating app stall every other application.
//   - A guest interrupted because the host stopped waiting (shutdown, cancelled
//     context, caller budget) is TRANSIENT: a plain Go error, retried on the next
//     poll, nothing signed. Charging a user MinFeePerRequest for our own shutdown
//     would be indefensible.
//
// The distinction is invisible in pkg/wasm — both are traps — so it has to be
// asserted here, at the boundary where the choice is actually made.

// boundStubRuntime returns whatever failure the test asks for, so the executor's
// classification can be exercised without a real WASM guest.
type boundStubRuntime struct {
	processFailure *apperrors.RequestFailure
	depositFailure *apperrors.RequestFailure
	deployErr      error
}

func (r *boundStubRuntime) Deploy(ctx context.Context, appId common.ApplicationIdType, constructorParams []byte, wasm []byte) ([]byte, *big.Int, error) {
	if r.deployErr != nil {
		return nil, big.NewInt(0), r.deployErr
	}
	return []byte(`{"appId":1,"accounts":{},"nonce":0}`), big.NewInt(1), nil
}

func (r *boundStubRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, tokenAddress ethCommon.Address, depositAmount *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.AppEvent, *big.Int, *apperrors.RequestFailure) {
	if r.depositFailure != nil {
		return nil, nil, nil, big.NewInt(0), r.depositFailure
	}
	return state, nil, nil, big.NewInt(1), nil
}

func (r *boundStubRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.AppEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	if r.processFailure != nil {
		return nil, nil, nil, nil, nil, big.NewInt(0), r.processFailure
	}
	return state, nil, nil, nil, nil, big.NewInt(1), nil
}

func (r *boundStubRuntime) Close() error { return nil }

// newBoundTestRequest builds a TrustProcess request: its payload is not encrypted
// and it skips the fee checks, so a test can reach the runtime call without also
// having to set up user keys.
func newBoundTestRequest() *common.Request {
	req := newProcessRequest()
	req.RequestType = common.TrustProcess
	req.Payload = []byte("payload")
	req.MaxFeeValue = common.NewBig(0)
	return req
}

func TestCancelledGuestExecutionIsTransient(t *testing.T) {
	runtime := &boundStubRuntime{
		processFailure: apperrors.New(apperrors.CodeGuestExecutionCancelled,
			apperrors.ErrGuestExecutionCancelled.Error()),
	}
	exec := newTestExecutor(t, runtime)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, nil, nil, wasmModule)

	payload, state, report, err := exec.HandleProcessRequest(
		context.Background(), newBoundTestRequest(), appState, wasmModule)

	require.Error(t, err, "an abandoned request must surface as a plain error so it is retried")
	require.Nil(t, payload, "nothing may be signed: the request was abandoned, not judged to have failed")
	require.Nil(t, state)
	require.Nil(t, report)
}

func TestTimedOutGuestExecutionIsSignedOnChain(t *testing.T) {
	runtime := &boundStubRuntime{
		processFailure: apperrors.New(apperrors.CodeGuestExecutionTimeout,
			apperrors.ErrGuestExecutionTimeout.Error()),
	}
	exec := newTestExecutor(t, runtime)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, nil, nil, wasmModule)

	payload, _, _, err := exec.HandleProcessRequest(
		context.Background(), newBoundTestRequest(), appState, wasmModule)

	require.NoError(t, err, "a guest that burned its whole budget is a request failure, not an executor error")
	require.NotNil(t, payload, "the failure must be signed so the on-chain queue head advances")
	require.NotZero(t, payload.ErrorCode, "a signed failure must carry a non-zero error code")
	require.Equal(t, apperrors.ErrGuestExecutionTimeout.Error(), payload.ErrorMsg,
		"the signed message must be the stable classified string, not wasmtime's text")
	require.LessOrEqual(t, len(payload.ErrorMsg), 100,
		"the signed message must survive the on-chain truncation intact")
	require.NotEmpty(t, payload.Signature)
}

func TestCancelledDepositIsTransient(t *testing.T) {
	runtime := &boundStubRuntime{
		depositFailure: apperrors.New(apperrors.CodeGuestExecutionCancelled,
			apperrors.ErrGuestExecutionCancelled.Error()),
	}
	exec := newTestExecutor(t, runtime)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, nil, nil, wasmModule)

	req := newBoundTestRequest()
	req.AssetAmount = common.NewBig(100) // routes through the deposit path first

	payload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, wasmModule)

	require.Error(t, err, "a deposit interrupted by shutdown must be retried, not signed")
	require.Nil(t, payload)
}

func TestCancelledDeployIsTransient(t *testing.T) {
	runtime := &boundStubRuntime{
		deployErr: fmt.Errorf("failed to call deploy: %w", apperrors.ErrGuestExecutionCancelled),
	}
	exec := newTestExecutor(t, runtime)

	req, wasmModule := newDeployRequest(t)

	payload, state, err := exec.HandleDeployApp(context.Background(), req, nil, wasmModule)

	require.Error(t, err, "a deploy interrupted by shutdown must be retried, not signed")
	require.Nil(t, payload)
	require.Nil(t, state)
}

func TestTimedOutDeployIsSignedOnChain(t *testing.T) {
	runtime := &boundStubRuntime{
		deployErr: fmt.Errorf("failed to call deploy: %w", apperrors.ErrGuestExecutionTimeout),
	}
	exec := newTestExecutor(t, runtime)

	req, wasmModule := newDeployRequest(t)

	payload, _, err := exec.HandleDeployApp(context.Background(), req, nil, wasmModule)

	require.NoError(t, err)
	require.NotNil(t, payload, "a deploy whose constructor ran away must be signed, like any other request failure")
	require.NotZero(t, payload.ErrorCode)
}

// TestGuestExecutionTimeoutConfigValidation covers the startup guard. Unlike the
// runtime setter, which clamps so a late change cannot take a running enclave
// down, this must fail loudly before the executor serves anything.
func TestGuestExecutionTimeoutConfigValidation(t *testing.T) {
	baseConfig := func() *Config {
		return &Config{
			ChannelType:         "tcp",
			ChannelParams:       common.TcpChannelConnectionParams{Ip: "localhost", Port: 4000},
			FuelPricePerUnit:    big.NewInt(1),
			MinFeePerRequest:    big.NewInt(10),
			CommunicationParams: common.CommunicationParams{RequestTimeoutSec: 30},
			KeySetRecoveryType:  common.RecoveryTypeUnsafe,
		}
	}

	t.Run("zero selects the default", func(t *testing.T) {
		c := baseConfig()
		c.GuestExecutionTimeoutMs = 0
		require.NoError(t, c.Validate())
	})

	t.Run("in-range value accepted", func(t *testing.T) {
		c := baseConfig()
		c.GuestExecutionTimeoutMs = 5_000
		require.NoError(t, c.Validate())
	})

	t.Run("negative rejected", func(t *testing.T) {
		c := baseConfig()
		c.GuestExecutionTimeoutMs = -1
		err := c.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS")
	})

	t.Run("above ceiling rejected", func(t *testing.T) {
		c := baseConfig()
		c.GuestExecutionTimeoutMs = common.MaxGuestExecutionTimeoutMs + 1
		err := c.Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS")
	})
}
