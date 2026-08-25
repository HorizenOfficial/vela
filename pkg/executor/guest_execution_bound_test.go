package executor

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	"github.com/HorizenOfficial/vela/pkg/communication"
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

	// failOnCall, when non-zero, restricts processFailure to that ProcessRequest
	// call only (1-based); every other call succeeds. Used to place a failure at a
	// chosen position inside a batch. Zero means "fail every call".
	failOnCall int
	calls      int

	// delay makes each ProcessRequest consume wall-clock time, so a test can watch a
	// shared batch budget actually drain. An instant stub never shrinks the deadline.
	delay time.Duration
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
	r.calls++
	if r.delay > 0 {
		time.Sleep(r.delay)
	}
	if r.processFailure != nil && (r.failOnCall == 0 || r.calls == r.failOnCall) {
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

// TestTimedOutDeployIsSignedOnChain also pins that the REASON survives the 100-char
// truncation of the signed ErrorMsg.
//
// This is not hypothetical. The on-chain ErrorCode for a timeout is WASM_INTERNAL,
// shared with every other host-side WASM failure, so the message is the only thing
// that identifies a timeout — and docs/design/WASM_HOST_ABI.md tells integrators to
// match on it. The deploy path builds its message by nesting prefixes, and the
// start-section case ("failed to load or get module: failed to load module: failed to
// instantiate WASM module: guest execution timed out") is 113 characters, so a naive
// pass-through loses the reason exactly the way
// JIRA_TASK_trap_error_classification.md describes.
func TestTimedOutDeployIsSignedOnChain(t *testing.T) {
	for name, deployErr := range map[string]error{
		"deploy export": fmt.Errorf("failed to call deploy: %w", apperrors.ErrGuestExecutionTimeout),
		// The longest nesting the deploy path can produce: an interrupt during
		// instantiation, wrapped by compileAndInstantiate, then getOrLoadModule.
		"start section": fmt.Errorf("failed to load module: failed to instantiate WASM module: %w",
			apperrors.ErrGuestExecutionTimeout),
	} {
		t.Run(name, func(t *testing.T) {
			exec := newTestExecutor(t, &boundStubRuntime{deployErr: deployErr})
			req, wasmModule := newDeployRequest(t)

			payload, _, err := exec.HandleDeployApp(context.Background(), req, nil, wasmModule)

			require.NoError(t, err)
			require.NotNil(t, payload, "a deploy whose guest ran away must be signed, like any other request failure")
			require.NotZero(t, payload.ErrorCode)
			require.LessOrEqual(t, len(payload.ErrorMsg), 100,
				"the signed message is truncated to 100 characters on-chain")
			require.Contains(t, payload.ErrorMsg, apperrors.ErrGuestExecutionTimeout.Error(),
				"the reason must survive truncation intact, or a timeout is indistinguishable on-chain")
		})
	}
}

// The batch path shares the classification with the single-request path through the
// one check in executeRequest's softFailure. The two tests below pin the *different*
// consequences that follow from it inside a batch, which is where the distinction
// stops being academic:
//
//   - a timeout is a soft failure, so it becomes one error payload and the batch
//     CONTINUES — the queue advances past every request in the batch;
//   - a cancellation is a hard failure, so the batch STOPS and submits only what it
//     had already collected — the rest stay pending and are retried.

func TestBatchContinuesPastAGuestTimeout(t *testing.T) {
	runtime := &boundStubRuntime{
		processFailure: apperrors.New(apperrors.CodeGuestExecutionTimeout,
			apperrors.ErrGuestExecutionTimeout.Error()),
		failOnCall: 2, // the middle request of three
	}
	exec := newTestExecutor(t, runtime)

	user, _, userPub := newBatchTestUser(t)
	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	requests := []*common.Request{newBoundTestRequest(), newBoundTestRequest(), newBoundTestRequest()}

	payloads, sig, finalState, _, err := exec.HandleBatchProcessRequest(
		context.Background(), requests, appState, wasmModule)

	require.NoError(t, err)
	require.Len(t, payloads, 3,
		"a timeout is a soft failure: it must not stop the batch, so all three requests settle")
	require.Zero(t, payloads[0].ErrorCode)
	require.NotZero(t, payloads[1].ErrorCode, "the timed-out request must carry an error code")
	require.Zero(t, payloads[2].ErrorCode, "the request after the timeout must still succeed")
	require.NotEmpty(t, sig)
	require.NotNil(t, finalState)
}

func TestBatchStopsAtAGuestCancellation(t *testing.T) {
	runtime := &boundStubRuntime{
		processFailure: apperrors.New(apperrors.CodeGuestExecutionCancelled,
			apperrors.ErrGuestExecutionCancelled.Error()),
		failOnCall: 2, // the middle request of three
	}
	exec := newTestExecutor(t, runtime)

	user, _, userPub := newBatchTestUser(t)
	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	requests := []*common.Request{newBoundTestRequest(), newBoundTestRequest(), newBoundTestRequest()}

	payloads, _, _, _, err := exec.HandleBatchProcessRequest(
		context.Background(), requests, appState, wasmModule)

	require.NoError(t, err, "the batch itself did not fail — it stopped early and returns what it has")
	require.Len(t, payloads, 1,
		"a cancellation is a hard failure: the batch must stop, submitting only the request that completed")
	require.Zero(t, payloads[0].ErrorCode,
		"and nothing may be signed for the abandoned request — it stays pending and is retried")
}

// TestBatchStopsWhenBudgetCannotCoverAnotherRequest covers the budget-aware loop.
//
// A batch shares one caller budget. Without this check the loop starts a request it
// cannot possibly finish: the guest runs until the budget expires, is interrupted,
// classified as abandonment, and discarded — burning a full guest bound of execution
// and evicting the module (an interrupt is a trap) only for the work to be redone on
// the next poll. Stopping before starting it costs nothing and settles exactly the
// same requests.
func TestBatchStopsWhenBudgetCannotCoverAnotherRequest(t *testing.T) {
	const guestBound = 1 * time.Second

	runtime := &boundStubRuntime{delay: 700 * time.Millisecond}
	exec := newTestExecutor(t, runtime)
	exec.config.GuestExecutionTimeoutMs = guestBound.Milliseconds()

	user, _, userPub := newBatchTestUser(t)
	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	requests := []*common.Request{newBoundTestRequest(), newBoundTestRequest(), newBoundTestRequest()}

	// Room for the first request's bound, but not a second one after it has run.
	ctx, cancel := context.WithTimeout(context.Background(), guestBound+600*time.Millisecond)
	defer cancel()

	payloads, _, _, _, err := exec.HandleBatchProcessRequest(ctx, requests, appState, wasmModule)

	require.NoError(t, err)
	require.Len(t, payloads, 1,
		"the loop must stop once the remaining budget cannot cover another guest bound")
	require.Zero(t, payloads[0].ErrorCode)
	require.Equal(t, 1, runtime.calls,
		"the second request must never be started, not started and then interrupted")
}

// TestBatchWithoutADeadlineProcessesEverything is the other half: the budget-aware
// check must not fire when no caller deadline was supplied at all.
func TestBatchWithoutADeadlineProcessesEverything(t *testing.T) {
	runtime := &boundStubRuntime{}
	exec := newTestExecutor(t, runtime)
	exec.config.GuestExecutionTimeoutMs = 10_000

	user, _, userPub := newBatchTestUser(t)
	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	requests := []*common.Request{newBoundTestRequest(), newBoundTestRequest(), newBoundTestRequest()}

	payloads, _, _, _, err := exec.HandleBatchProcessRequest(
		context.Background(), requests, appState, wasmModule)

	require.NoError(t, err)
	require.Len(t, payloads, 3, "with no deadline the whole batch must be processed")
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

// TestGuestTimeoutMustBeatCallerBudget guards the forward-progress invariant.
//
// The caller's budget becomes a context deadline at
// requestTimeout - communication.ExecutionBudgetMargin. If the guest's own bound is
// not below that, the BUDGET wins the race on every slow request — and a budget
// expiry is host-side abandonment, i.e. transient. So a non-terminating guest would
// be retried on every poll forever, nothing would ever be signed, and the on-chain
// FIFO head would never advance. That is the stall signing a timeout exists to
// prevent, reintroduced purely by configuration.
//
// A plain `timeout >= requestTimeout` check is NOT enough: it leaves a silent window
// one margin wide (28000..29999 ms against the 30 s default) where validation passes
// and the queue still wedges.
func TestGuestTimeoutMustBeatCallerBudget(t *testing.T) {
	const patienceSec = 30
	budgetMs := int64(patienceSec)*1000 - communication.ExecutionBudgetMargin.Milliseconds()

	configWith := func(timeoutMs int64, patience time.Duration) *Config {
		return &Config{
			ChannelType:             "tcp",
			ChannelParams:           common.TcpChannelConnectionParams{Ip: "localhost", Port: 4000},
			FuelPricePerUnit:        big.NewInt(1),
			MinFeePerRequest:        big.NewInt(10),
			CommunicationParams:     common.CommunicationParams{RequestTimeoutSec: patience},
			KeySetRecoveryType:      common.RecoveryTypeUnsafe,
			GuestExecutionTimeoutMs: timeoutMs,
		}
	}

	t.Run("comfortably below the budget is accepted", func(t *testing.T) {
		require.NoError(t, configWith(10_000, patienceSec).Validate())
	})

	t.Run("the default is resolved and accepted", func(t *testing.T) {
		// 0 means the default, so the check must resolve it rather than skip.
		require.NoError(t, configWith(0, patienceSec).Validate())
	})

	t.Run("equal to the budget is rejected", func(t *testing.T) {
		err := configWith(budgetMs, patienceSec).Validate()
		require.Error(t, err, "a bound equal to the budget lets the budget win the race")
		require.Contains(t, err.Error(), "EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS")
	})

	t.Run("inside the silent window is rejected", func(t *testing.T) {
		// The case a bare `>= requestTimeout` check misses entirely.
		err := configWith(29_000, patienceSec).Validate()
		require.Error(t, err, "29s beats the 28s budget deadline and must not pass validation")
	})

	t.Run("caller too impatient for the default is rejected", func(t *testing.T) {
		// A request timeout at or below the margin drops the budget entirely, so the
		// executor would run its full bound long after the caller stopped waiting.
		require.Error(t, configWith(0, 3).Validate())
	})

	t.Run("a non-positive request timeout is left to its own check", func(t *testing.T) {
		// Guarded so this does not mask the existing "must be > 0" error with a
		// confusing message about the guest bound.
		err := configWith(10_000, 0).Validate()
		require.Error(t, err)
		require.Contains(t, err.Error(), "REQUEST_TIMEOUT_SEC")
	})
}
