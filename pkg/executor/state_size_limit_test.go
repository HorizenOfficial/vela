package executor

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const testMaxAppStateSize = 1024

// stateOfSize returns a JSON app state of exactly n bytes.
func stateOfSize(t *testing.T, n int) string {
	t.Helper()
	const prefix, suffix = `{"appId":1,"accounts":{},"nonce":0,"pad":"`, `"}`
	require.GreaterOrEqual(t, n, len(prefix)+len(suffix))
	return prefix + strings.Repeat("x", n-len(prefix)-len(suffix)) + suffix
}

// oversizedOutputRuntime returns a new state of outSize bytes from the calls whose grow flag is set.
type oversizedOutputRuntime struct {
	MockRuntime
	t            *testing.T
	outSize      int
	growDeposit  bool
	growProcess  bool
	growDeploy   bool
	depositCalls int
	processCalls int
}

func (r *oversizedOutputRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, tokenAddress ethCommon.Address, depositAmount *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.AppEvent, *big.Int, *apperrors.RequestFailure) {
	r.depositCalls++
	if r.growDeposit {
		return []byte(stateOfSize(r.t, r.outSize)), nil, nil, big.NewInt(1), nil
	}
	return state, nil, nil, big.NewInt(1), nil
}

func (r *oversizedOutputRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.AppEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	r.processCalls++
	if r.growProcess {
		return []byte(stateOfSize(r.t, r.outSize)), nil, nil, nil, nil, big.NewInt(1), nil
	}
	return state, nil, nil, nil, nil, big.NewInt(1), nil
}

func (r *oversizedOutputRuntime) Deploy(ctx context.Context, appID common.ApplicationIdType, constructorParams []byte, wasm []byte) ([]byte, *big.Int, error) {
	if r.growDeploy {
		return []byte(stateOfSize(r.t, r.outSize)), big.NewInt(1), nil
	}
	return r.MockRuntime.Deploy(ctx, appID, constructorParams, wasm)
}

func newStateLimitExecutor(t *testing.T, runtime Runtime) *StatelessExecutor {
	t.Helper()
	exec := newTestExecutor(t, runtime)
	exec.config.MaxAppStateSize = testMaxAppStateSize
	return exec
}

// newTrustProcessRequest returns a request that reaches the guest without payload encryption.
func newTrustProcessRequest() *common.Request {
	req := newProcessRequest()
	req.RequestType = common.TrustProcess
	req.Payload = []byte(`{}`)
	return req
}

func requireStateTooLarge(t *testing.T, payload *common.UpdatePayload, newAppState *common.ApplicationState, err error, prevRoot [32]byte) {
	t.Helper()
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Nil(t, newAppState, "the oversized state must not be stored")
	require.Equal(t, uint8(apperrors.CodeAppStateTooLarge.Category.Category), payload.ErrorCode)
	require.Contains(t, payload.ErrorMsg, "application state too large")
	require.LessOrEqual(t, len(payload.ErrorMsg), 100)
	require.Equal(t, prevRoot, payload.PrevStateRoot)
	require.Equal(t, prevRoot, payload.NewStateRoot, "state unchanged on error")
	require.Len(t, payload.Signature, 65, "the failure is signed so the manager can post it")
}

func TestAppStateTooLargeCode(t *testing.T) {
	require.Equal(t, "APP_STATE_TOO_LARGE", apperrors.CodeAppStateTooLarge.Code)
	require.Equal(t, apperrors.CategoryRequestFuncFailedMeta, apperrors.CodeAppStateTooLarge.Category)
}

func TestHandleProcessRequest_InputStateTooLarge_FailsWithoutCallingTheGuest(t *testing.T) {
	for _, withDeposit := range []bool{false, true} {
		runtime := &countingRuntime{MockRuntime: *NewMockRuntime(testLogger)}
		exec := newStateLimitExecutor(t, runtime)
		appState := buildEncryptedAppStateFromJSON(t, exec, stateOfSize(t, testMaxAppStateSize+1), nil, nil)
		req := newTrustProcessRequest()
		if withDeposit {
			req.AssetAmount = common.NewBig(1000)
		}

		payload, newAppState, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
		requireStateTooLarge(t, payload, newAppState, err, appState.StateRoot)
		require.Zero(t, runtime.depositCalls, "deposit=%v", withDeposit)
		require.Zero(t, runtime.processCalls, "deposit=%v", withDeposit)
	}
}

func TestHandleProcessRequest_StateAtTheLimitIsAccepted(t *testing.T) {
	runtime := &oversizedOutputRuntime{MockRuntime: *NewMockRuntime(testLogger), t: t, outSize: testMaxAppStateSize, growDeposit: true, growProcess: true}
	exec := newStateLimitExecutor(t, runtime)
	appState := buildEncryptedAppStateFromJSON(t, exec, stateOfSize(t, testMaxAppStateSize), nil, nil)
	req := newTrustProcessRequest()
	req.AssetAmount = common.NewBig(1000)

	payload, newAppState, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, uint8(0), payload.ErrorCode, payload.ErrorMsg)
	require.NotNil(t, newAppState)
	require.Equal(t, 1, runtime.depositCalls)
	require.Equal(t, 1, runtime.processCalls)
}

func TestHandleProcessRequest_DepositOutputStateTooLarge(t *testing.T) {
	runtime := &oversizedOutputRuntime{MockRuntime: *NewMockRuntime(testLogger), t: t, outSize: testMaxAppStateSize + 1, growDeposit: true}
	exec := newStateLimitExecutor(t, runtime)
	appState := buildEncryptedAppStateFromJSON(t, exec, stateOfSize(t, 100), nil, nil)
	req := newTrustProcessRequest()
	req.AssetAmount = common.NewBig(1000)

	payload, newAppState, _, err := exec.HandleProcessRequest(context.Background(), req, appState, []byte("wasm"))
	requireStateTooLarge(t, payload, newAppState, err, appState.StateRoot)
	require.Equal(t, 1, runtime.depositCalls)
	require.Zero(t, runtime.processCalls, "the request is not processed after a rejected deposit")
}

func TestHandleProcessRequest_ProcessOutputStateTooLarge(t *testing.T) {
	runtime := &oversizedOutputRuntime{MockRuntime: *NewMockRuntime(testLogger), t: t, outSize: testMaxAppStateSize + 1, growProcess: true}
	exec := newStateLimitExecutor(t, runtime)
	appState := buildEncryptedAppStateFromJSON(t, exec, stateOfSize(t, 100), nil, nil)

	payload, newAppState, _, err := exec.HandleProcessRequest(context.Background(), newTrustProcessRequest(), appState, []byte("wasm"))
	requireStateTooLarge(t, payload, newAppState, err, appState.StateRoot)
	require.Equal(t, 1, runtime.processCalls)
}

func TestHandleProcessRequest_ZeroLimitMeansNoLimit(t *testing.T) {
	runtime := &oversizedOutputRuntime{MockRuntime: *NewMockRuntime(testLogger), t: t, outSize: 4 * testMaxAppStateSize, growProcess: true}
	exec := newTestExecutor(t, runtime)
	appState := buildEncryptedAppStateFromJSON(t, exec, stateOfSize(t, 2*testMaxAppStateSize), nil, nil)

	payload, newAppState, _, err := exec.HandleProcessRequest(context.Background(), newTrustProcessRequest(), appState, []byte("wasm"))
	require.NoError(t, err)
	require.Equal(t, uint8(0), payload.ErrorCode, payload.ErrorMsg)
	require.NotNil(t, newAppState)
}

func TestHandleDeployApp_InitialStateTooLarge(t *testing.T) {
	runtime := &oversizedOutputRuntime{MockRuntime: *NewMockRuntime(testLogger), t: t, outSize: testMaxAppStateSize + 1, growDeploy: true}
	exec := newStateLimitExecutor(t, runtime)
	req, wasmModule := newDeployRequest(t)

	payload, newAppState, err := exec.HandleDeployApp(context.Background(), req, nil, wasmModule)
	requireStateTooLarge(t, payload, newAppState, err, [32]byte{})
}
