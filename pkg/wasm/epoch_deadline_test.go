package wasm

import (
	"context"
	goruntime "runtime"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	wasmtime "github.com/bytecodealliance/wasmtime-go/v47"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Tests for the wall-clock execution bound (epoch interruption). See
// docs/design/WASM_HOST_ABI.md for the resulting guest-visible contract.
//
// The fixtures here spin FOREVER, unlike spinningProcessRequestWat in
// wasmtime_runtime_test.go which counts down 20M iterations and therefore
// terminates on its own. A bounded spinner cannot test an execution bound: it
// would pass even with the bound removed.
//
// testGuestTimeout is short so the suite stays fast, but not shorter than the
// epoch granularity can honour: epochTicksFor rounds up and adds a tick, so a
// budget of T is interrupted somewhere in [T, T+epochTickInterval]. At 200ms that
// is ~300ms, leaving requireBoundedElapsed's lower bound (T/2 = 100ms) a 200ms
// margin against scheduling jitter on a loaded machine.
//
// Do not lower this further without lowering epochTickInterval too: 150ms arms the
// same 3 ticks as 200ms and saves nothing, and 100ms leaves only 50ms of margin.
const testGuestTimeout = 200 * time.Millisecond

// guestEntryDelay is how long to wait for a guest to actually be executing before
// interrupting it from another goroutine. It is not a latency tuning knob: if it
// were too short the interrupt could land before the guest entered wasm, and tests
// like TestCloseInterruptsRunningGuest would pass without exercising the
// interruption they exist to prove.
const guestEntryDelay = 150 * time.Millisecond

// newBoundedRuntime returns a runtime with the short test timeout already applied.
func newBoundedRuntime(t *testing.T, maxCachedModules int) *WasmtimeRuntime {
	t.Helper()
	r := NewWasmtimeRuntime(testLogger, maxCachedModules)
	r.SetGuestExecutionTimeout(testGuestTimeout)
	return r
}

// foreverSpinningProcessRequestWat never returns from process_request. The
// i32.const after the loop is unreachable and exists only to satisfy validation.
const foreverSpinningProcessRequestWat = `(module
  (memory (export "memory") 1)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    (loop $l br $l)
    i32.const 100)
)`

// foreverSpinningAllocateWat hangs in the allocate export, i.e. inside
// writeToMemory, before the request export is ever entered. Those Func.Call
// sites need the same bound as the main exports, or a guest could hang the
// executor from its allocator.
const foreverSpinningAllocateWat = `(module
  (memory (export "memory") 1)
  (func (export "allocate") (param i32) (result i32)
    (loop $l br $l)
    i32.const 500)
  (func (export "deallocate") (param i32 i32))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100)
)`

// spinningStartSectionWat hangs during Instantiate rather than in an export, so
// it is only bounded if compileAndInstantiate arms the store before
// instantiating (measured fact M3 in the plan). Reachable from Deploy and from
// any cache-miss reload.
const spinningStartSectionWat = `(module
  (memory (export "memory") 1)
  (func $spin (loop $l br $l))
  (start $spin)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100)
)`

// spinningDeallocateWat succeeds at everything the host reads a result from, and
// only hangs in `deallocate` — the fire-and-forget cleanup whose errors the host
// drops. It models a guest that returns a valid result and then refuses to let go
// of its memory, which is the one way an interrupt can land without any call the
// host checks having failed.
const spinningDeallocateWat = `(module
  (memory (export "memory") 1)

  ;; same result layout as dispatchTestWat: state [1]
  (data (i32.const 100)
    "\46\00\00\00"
    "\7b\22\73\74\61\74\65\22\3a\5b\31\5d\2c\22\65\76\65\6e\74\73\22\3a\5b\5d\2c"
    "\22\61\70\70\45\76\65\6e\74\73\22\3a\5b\5d\2c\22\77\69\74\68\64\72\61\77\61"
    "\6c\73\22\3a\5b\5d\2c\22\66\75\65\6c\22\3a\22\30\78\31\22\7d"
  )

  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32) (loop $l br $l))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100)
)`

// deployThenSpinningDeallocateWat deploys successfully and then hangs in
// `deallocate`. The deploy path caches the instance instead of evicting it, so it
// needs its own fixture to prove the interrupted-cleanup case is handled there too.
const deployThenSpinningDeallocateWat = `(module
  (memory (export "memory") 1)

  ;; DeployResult at offset 100: 4-byte LE length (26 = 0x1a) + {"state":[1],"fuel":"0x1"}
  (data (i32.const 100)
    "\1a\00\00\00"
    "\7b\22\73\74\61\74\65\22\3a\5b\31\5d\2c\22\66\75\65\6c\22\3a\22\30\78\31\22\7d"
  )

  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32) (loop $l br $l))
  (func (export "deploy") (param i64 i32 i32) (result i32) i32.const 100)
)`

// foreverSpinningDeployWat hangs in the deploy export. Deploy returns a plain
// error rather than a *RequestFailure, so it carries the classification as a
// sentinel instead of a failure code.
const foreverSpinningDeployWat = `(module
  (memory (export "memory") 1)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))
  (func (export "deploy") (param i64 i32 i32) (result i32)
    (loop $l br $l)
    i32.const 100)
)`

// requireBoundedElapsed asserts the call was interrupted within a sane window:
// not absurdly early (which is what an unarmed store does — it traps instantly,
// measured fact M1), and comfortably before the shutdown backstop.
//
// The lower bound is testGuestTimeout/2 rather than testGuestTimeout because the
// ticker's phase is independent of when a store is armed: the first tick after
// arming can land almost immediately, so the observable floor is one tick below
// the nominal timeout.
func requireBoundedElapsed(t *testing.T, elapsed time.Duration) {
	t.Helper()
	require.Greater(t, elapsed, testGuestTimeout/2,
		"interrupted far too early — the deadline is not being armed from the configured timeout")
	require.Less(t, elapsed, shutdownExecLockTimeout,
		"guest was not interrupted within the configured bound")
}

// TestRunawayGuestIsInterrupted is the core acceptance criterion: a guest export
// that never returns fails with a classified error within the bound, the module
// is dropped from the cache, and — the part that matters for the platform — an
// unrelated app is served normally afterwards instead of being blocked behind it.
func TestRunawayGuestIsInterrupted(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	spinning, err := wasmtime.Wat2Wasm(foreverSpinningProcessRequestWat)
	require.NoError(t, err)
	healthy, err := wasmtime.Wat2Wasm(dispatchTestWat)
	require.NoError(t, err)

	runawayApp := common.NewApplicationId(1)
	otherApp := common.NewApplicationId(2)
	ctx := context.Background()

	start := time.Now()
	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, runawayApp, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), spinning)
	elapsed := time.Since(start)

	require.NotNil(t, failure, "a guest that never returns must fail the request")
	require.Equal(t, apperrors.CodeGuestExecutionTimeout, failure.RequestError,
		"a wall-clock deadline expiry is a timeout, not a generic request failure")
	requireBoundedElapsed(t, elapsed)

	// An interrupted guest's heap is in an unknown state, exactly like a panic.
	require.False(t, isCached(runtime, runawayApp),
		"an interrupted module must not stay cached")

	// The whole point of the ticket: the runaway app must not wedge the executor.
	state, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, otherApp, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), healthy)
	require.Nil(t, failure, "an unrelated app must still be served after a runaway guest")
	require.Equal(t, []byte{1}, state)
}

// TestRunawayAllocateIsInterrupted covers the writeToMemory call sites: a guest
// can hang in allocate before its request export is ever entered.
func TestRunawayAllocateIsInterrupted(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(foreverSpinningAllocateWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)

	start := time.Now()
	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		context.Background(), appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	elapsed := time.Since(start)

	require.NotNil(t, failure, "a guest hanging in allocate must fail the request")
	require.Equal(t, apperrors.CodeGuestExecutionTimeout, failure.RequestError,
		"an interrupt in allocate must classify as a timeout, not as a memory-write error")
	requireBoundedElapsed(t, elapsed)
	require.False(t, isCached(runtime, appId))
}

// TestSpinningStartSectionIsInterrupted pins measured fact M3: guest code can run
// during Instantiate, so the store must be armed before instantiating and not
// only before calling an export.
func TestSpinningStartSectionIsInterrupted(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(spinningStartSectionWat)
	require.NoError(t, err)

	start := time.Now()
	_, err = runtime.getOrLoadModule(context.Background(), common.NewApplicationId(1), wasmBytes)
	elapsed := time.Since(start)

	require.Error(t, err, "a module whose start section never returns must fail to load")
	require.ErrorIs(t, err, apperrors.ErrGuestExecutionTimeout,
		"a start-section interrupt must be classified, not surface as a bare instantiate failure")
	requireBoundedElapsed(t, elapsed)
}

// TestRunawayDeployIsInterrupted covers the deploy path, which returns a plain
// error and therefore carries the classification as a sentinel.
func TestRunawayDeployIsInterrupted(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(foreverSpinningDeployWat)
	require.NoError(t, err)

	start := time.Now()
	_, _, err = runtime.Deploy(context.Background(), common.NewApplicationId(1), []byte("{}"), wasmBytes)
	elapsed := time.Since(start)

	require.Error(t, err, "a deploy whose constructor never returns must fail")
	require.ErrorIs(t, err, apperrors.ErrGuestExecutionTimeout)
	requireBoundedElapsed(t, elapsed)
}

// TestInterruptedCleanupEvictsModule covers the one interrupt that no checked call
// reports. The host drops the error from every deferred `deallocate` (`_, _ =`), so
// a guest that returns a valid result and then hangs in cleanup would otherwise
// stay cached with an allocator that trapped mid-free — precisely the damaged heap
// evictFaultedModule exists to drop, and the state in which the comment on
// errGuestFault warns a module could "commit a plausible-but-incorrect state root
// on-chain".
//
// The request itself must still SUCCEED: the result bytes are copied out of guest
// memory before any deallocate runs, so the state transition is validly produced.
// Only the instance is suspect, and only the instance is thrown away.
func TestInterruptedCleanupEvictsModule(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(spinningDeallocateWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)

	state, _, _, _, _, _, failure := runtime.ProcessRequest(
		context.Background(), appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)

	require.Nil(t, failure,
		"the result was extracted before cleanup ran, so the request itself succeeded")
	require.Equal(t, []byte{1}, state)

	require.False(t, isCached(runtime, appId),
		"a module whose cleanup was interrupted must be evicted: its allocator trapped mid-free")
}

// TestInterruptedCleanupIsNotCachedAfterDeploy is the deploy-path counterpart of
// TestInterruptedCleanupEvictsModule. Deploy caches the instance rather than
// evicting it, and it does so before its deferred deallocate has run — so the
// "was cleanup interrupted" answer is only final in a later defer. The deploy
// itself must still succeed and return its state; only the instance is discarded.
func TestInterruptedCleanupIsNotCachedAfterDeploy(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(deployThenSpinningDeallocateWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)

	state, _, err := runtime.Deploy(context.Background(), appId, []byte("{}"), wasmBytes)

	require.NoError(t, err, "the deploy result was extracted before cleanup ran, so the deploy succeeded")
	require.Equal(t, []byte{1}, state)

	require.False(t, isCached(runtime, appId),
		"an instance whose cleanup was interrupted must not be cached: its allocator trapped mid-free")
}

// TestContextCancellationInterruptsGuest is acceptance criterion 3. The
// classification differs from a deadline expiry on purpose: we abandoned the
// request, so it must be retried rather than signed as a permanent on-chain
// failure that charges the user a fee.
func TestContextCancellationInterruptsGuest(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	// Far longer than the cancellation below, so a passing test cannot be the
	// deadline firing by coincidence.
	runtime.SetGuestExecutionTimeout(30 * time.Second)

	wasmBytes, err := wasmtime.Wat2Wasm(foreverSpinningProcessRequestWat)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(guestEntryDelay)
		cancel()
	}()
	defer cancel()

	start := time.Now()
	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, common.NewApplicationId(1), ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	elapsed := time.Since(start)

	require.NotNil(t, failure, "a cancelled guest call must fail")
	require.Equal(t, apperrors.CodeGuestExecutionCancelled, failure.RequestError,
		"cancellation is host-side abandonment and must not be signed as a guest timeout")
	require.Less(t, elapsed, 10*time.Second,
		"cancellation must interrupt promptly, not wait out the configured timeout")
}

// TestNormalGuestCallsAreUnaffected is the other half of every restriction: that
// legitimate guests still work. The idle period in the middle pins measured fact
// M2 — the epoch deadline is absolute, so wall-clock time passing between calls
// consumes it. Arming once per store instead of once per operation would make
// this fail while leaving the runaway tests green.
func TestNormalGuestCallsAreUnaffected(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(dispatchTestWat)
	require.NoError(t, err)

	appId := common.NewApplicationId(1)
	ctx := context.Background()

	state, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	require.Nil(t, failure, "a healthy guest must not be interrupted")
	require.Equal(t, []byte{1}, state)

	// Idle for longer than the whole timeout, keeping the module cached.
	time.Sleep(testGuestTimeout * 2)

	state, _, _, _, _, _, failure = runtime.ProcessRequest(
		ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	require.Nil(t, failure,
		"a cached module must be re-armed per call: idle wall-clock time must not consume the deadline")
	require.Equal(t, []byte{1}, state)

	// Several consecutive calls must all succeed, not just the first after arming.
	for i := 0; i < 5; i++ {
		_, _, _, _, _, _, failure = runtime.ProcessRequest(
			ctx, appId, ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
		require.Nil(t, failure, "consecutive healthy call %d must succeed", i)
	}
}

// TestCloseInterruptsRunningGuest is acceptance criterion 2: Close must complete
// because the guest was interrupted, NOT because it gave up after
// shutdownExecLockTimeout. TestCloseDoesNotHangOnStuckGuest covers the fallback;
// this asserts the fallback is not what fires.
func TestCloseInterruptsRunningGuest(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)

	// Long enough that only Close-driven interruption can end this call.
	runtime.SetGuestExecutionTimeout(30 * time.Second)

	wasmBytes, err := wasmtime.Wat2Wasm(foreverSpinningProcessRequestWat)
	require.NoError(t, err)

	callReturned := make(chan struct{})
	go func() {
		defer close(callReturned)
		runtime.ProcessRequest(context.Background(), common.NewApplicationId(1),
			ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	}()

	// Let the guest actually enter the spin before shutting down, or Close could win
	// the lock before any guest code ran and the test would prove nothing.
	time.Sleep(guestEntryDelay)

	start := time.Now()
	closeErr := runtime.Close()
	elapsed := time.Since(start)

	require.NoError(t, closeErr,
		"Close must acquire execLock by interrupting the guest, not report that it gave up")
	require.Less(t, elapsed, shutdownExecLockTimeout,
		"Close fell back to the bounded-acquire timeout instead of interrupting the guest")

	select {
	case <-callReturned:
	case <-time.After(5 * time.Second):
		t.Fatal("the interrupted guest call never returned")
	}
}

// TestNoGoroutineLeak guards the ticker and the per-call watchers: both are
// started by the runtime and must be joined by Close, or every executor restart
// in a test binary — and every runtime in production — leaks them.
func TestNoGoroutineLeak(t *testing.T) {
	wasmBytes, err := wasmtime.Wat2Wasm(dispatchTestWat)
	require.NoError(t, err)
	spinning, err := wasmtime.Wat2Wasm(foreverSpinningProcessRequestWat)
	require.NoError(t, err)

	settle(t)
	before := goruntime.NumGoroutine()

	runtime := newBoundedRuntime(t, 0)
	ctx := context.Background()

	// A healthy call (watcher exits via the done channel) and an interrupted one
	// (watcher exits after forcing the interrupt) — both watcher exit paths.
	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, common.NewApplicationId(1), ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	require.Nil(t, failure)

	_, _, _, _, _, _, failure = runtime.ProcessRequest(
		ctx, common.NewApplicationId(2), ethCommon.Address{}, common.Process, []byte("{}"), []byte("{}"), spinning)
	require.NotNil(t, failure)

	require.NoError(t, runtime.Close())

	settle(t)
	after := goruntime.NumGoroutine()
	require.LessOrEqual(t, after, before+1,
		"goroutines leaked across runtime lifecycle (before=%d after=%d): check the epoch ticker and the per-call watchers", before, after)
}

// settle waits for goroutine counts to stop moving, so the leak check does not
// race with goroutines that are already on their way out (the log-pipe readers
// have a short drain of their own).
func settle(t *testing.T) {
	t.Helper()
	prev := -1
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		n := goruntime.NumGoroutine()
		if n == prev {
			return
		}
		prev = n
		time.Sleep(100 * time.Millisecond)
	}
}

// TestSetGuestExecutionTimeoutClamps mirrors TestSetMaxGuestMemoryBytes: a late
// admin-driven change must never take a running enclave down, so out-of-range
// values are clamped and reported rather than rejected. Startup validation in
// pkg/executor is where a bad value fails loudly.
func TestSetGuestExecutionTimeoutClamps(t *testing.T) {
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	require.Equal(t, defaultGuestExecutionTimeout, runtime.GetGuestExecutionTimeout(),
		"a fresh runtime must be bounded by the default, never unbounded")

	runtime.SetGuestExecutionTimeout(2 * time.Second)
	require.Equal(t, 2*time.Second, runtime.GetGuestExecutionTimeout())

	// Non-positive means "use the default" — never "no bound". A zero deadline
	// would arm an already-expired store (measured fact M5) and fail every call.
	for _, invalid := range []time.Duration{0, -1, -time.Hour} {
		runtime.SetGuestExecutionTimeout(invalid)
		require.Equal(t, defaultGuestExecutionTimeout, runtime.GetGuestExecutionTimeout(),
			"non-positive timeout %v must fall back to the default", invalid)
	}

	runtime.SetGuestExecutionTimeout(maxGuestExecutionTimeout + time.Second)
	require.Equal(t, maxGuestExecutionTimeout, runtime.GetGuestExecutionTimeout(),
		"an absurd timeout must be clamped to the ceiling")
}
