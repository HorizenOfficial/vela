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

// blockingWasiSleepWat does not spin at all: it asks the HOST to sleep, via the
// WASI poll_oneoff clock subscription that TinyGo's time.Sleep lowers to.
//
// This is the one escape epoch interruption cannot close by itself. Epoch checks
// are emitted into compiled guest code, so they are only reached while the guest
// is executing; a guest parked inside a host call is off in native code, where
// bumping the epoch does nothing until it comes back. The watcher, the caller's
// execution budget and Close all pull that same dead lever, so without a host-side
// restriction the guest runs for as long as it asked for — holding execLock, which
// is global, so every other application waits behind it.
//
// The subscription layout (48 bytes at offset 0) is written with explicit stores
// rather than a hand-encoded data segment so it stays auditable:
//
//	[0]  userdata u64     = 0
//	[8]  eventtype u8     = 0 (clock); the union is 8-aligned, so it starts at 16
//	[16] clock id u32     = 1 (monotonic)
//	[24] timeout u64 (ns) = wasiSleepRequest
//	[32] precision u64    = 0
//	[40] flags u16        = 0 (relative)
//
// The out-events buffer is at 104 (8-aligned, as wasmtime requires) and the
// nevents out-param at 200. The errno is dropped: a guest whose sleep is refused
// carries on and returns its result, which is what makes the exclusion visible as
// a fast, successful call rather than as a failure.
const blockingWasiSleepWat = `(module
  (import "wasi_snapshot_preview1" "poll_oneoff"
    (func $poll (param i32 i32 i32 i32) (result i32)))
  (memory (export "memory") 1)

  ;; same result layout as dispatchTestWat: state [1]
  (data (i32.const 300)
    "\54\00\00\00"
    "{\22state\22:[1],\22events\22:[],\22appEvents\22:[],\22withdrawals\22:[],\22report\22:null,\22fuel\22:\220x0\22}")

  (func $sleep
    (i64.store   (i32.const 0)  (i64.const 0))
    (i32.store8  (i32.const 8)  (i32.const 0))
    (i32.store   (i32.const 16) (i32.const 1))
    (i64.store   (i32.const 24) (i64.const 5000000000))
    (i64.store   (i32.const 32) (i64.const 0))
    (i32.store16 (i32.const 40) (i32.const 0))
    (drop (call $poll (i32.const 0) (i32.const 104) (i32.const 1) (i32.const 200))))

  (func (export "allocate") (param i32) (result i32) i32.const 900)
  (func (export "deallocate") (param i32 i32))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    (call $sleep)
    i32.const 300)
)`

// wasiSleepRequest is the sleep blockingWasiSleepWat asks the host for. It is far
// above both testGuestTimeout and shutdownExecLockTimeout, so a guest that is
// allowed to serve it is unmistakable in the elapsed time.
const wasiSleepRequest = 5 * time.Second

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
// TestBlockingWasiCallCannotOutlastTheGuestBound is the exclusion test for the
// restricted WASI surface: it proves a VIOLATING guest is refused, not merely that
// legitimate ones still run. Every other fixture in this file spins in pure wasm,
// where epoch interruption was never in doubt — this one blocks in a host call,
// which is where the bound had no reach at all.
//
// The load-time import check makes this a refusal rather than a neutered call, so
// the assertion is on elapsed time: whatever else happens, the host must never
// serve the sleep the guest asked for.
func TestBlockingWasiCallCannotOutlastTheGuestBound(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(blockingWasiSleepWat)
	require.NoError(t, err)

	start := time.Now()
	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		context.Background(), common.NewApplicationId(1), ethCommon.Address{},
		common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	elapsed := time.Since(start)

	require.Less(t, elapsed, wasiSleepRequest/2,
		"a guest blocked in a WASI host call outran its execution bound: the host served "+
			"the sleep it asked for, holding execLock and stalling every other application")
	require.NotNil(t, failure, "a guest that can block the host must not be loaded at all")
	require.Equal(t, apperrors.CodeFailedLoadingOrGettingModule, failure.RequestError,
		"the refusal must be a signed load failure so the request settles and the queue advances")
}

// TestDisallowedImportIsRejectedAtLoad pins the shape of the restriction: a guest is
// refused for DECLARING a host function outside the allowed set, whether or not it
// ever calls it, and whether it comes from WASI or from some other namespace.
//
// Declaring is the right trigger. An import that is merely present is still an
// import the guest can reach at any time, and refusing at load makes it a signed
// on-chain failure at deploy — the earliest, loudest point, where the queue still
// advances. The alternative, letting it load and neutering the call, would change
// guest semantics silently: a guest tested against stock wasmtime would behave
// differently inside the enclave, which is a worse failure than a refused deploy.
func TestDisallowedImportIsRejectedAtLoad(t *testing.T) {
	cases := []struct {
		name    string
		wat     string
		offence string
	}{
		{
			name:    "a blocking WASI import",
			wat:     blockingWasiSleepWat,
			offence: "wasi_snapshot_preview1.poll_oneoff",
		},
		{
			name:    "a host function outside WASI",
			offence: "vela_host.do_something",
			wat: `(module
  (import "vela_host" "do_something" (func (param i32) (result i32)))
  (memory (export "memory") 1)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100)
)`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wasmBytes, err := wasmtime.Wat2Wasm(tc.wat)
			require.NoError(t, err)

			runtime := newBoundedRuntime(t, 0)
			defer runtime.Close()

			_, err = runtime.getOrLoadModule(context.Background(), common.NewApplicationId(1), wasmBytes)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.offence,
				"the error must name the offending import so the app author can act on it")
		})
	}
}

// TestGuestKeepsTheWasiImportsItActuallyUses is the other half of the exclusion
// test: the restriction must not cost a legitimate guest anything it uses. An
// import a module declares but the linker does not define fails at INSTANTIATION,
// not at call time, so removing one too many breaks every guest at load.
//
// The list is the exact import set of the shipped TinyGo guests, read off
// app/simple/build/simple_app.wasm and app/trigger/build/trigger_app.wasm with
// wasm2wat. It is spelled out as a WAT fixture rather than loaded from those
// artifacts so this runs in the quick suite without a TinyGo build; the real
// guests are covered end to end by app/simple/integration_test.go.
const realGuestWasiImportsWat = `(module
  (import "wasi_snapshot_preview1" "args_get" (func (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "args_sizes_get" (func (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "clock_time_get" (func (param i32 i64 i32) (result i32)))
  (import "wasi_snapshot_preview1" "environ_get" (func (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "environ_sizes_get" (func (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "fd_write" (func (param i32 i32 i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "proc_exit" (func (param i32)))
  (import "wasi_snapshot_preview1" "random_get" (func (param i32 i32) (result i32)))
  (memory (export "memory") 1)
  (func (export "allocate") (param i32) (result i32) i32.const 500)
  (func (export "deallocate") (param i32 i32))
  (func (export "process_request") (param i64 i32 i32 i32 i32 i32 i32 i32) (result i32)
    i32.const 100)
)`

func TestGuestKeepsTheWasiImportsItActuallyUses(t *testing.T) {
	wasmBytes, err := wasmtime.Wat2Wasm(realGuestWasiImportsWat)
	require.NoError(t, err)

	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	_, err = runtime.getOrLoadModule(context.Background(), common.NewApplicationId(1), wasmBytes)
	require.NoError(t, err,
		"a guest importing only what the real TinyGo guests import must still instantiate")
}

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

// fdReadProbeWat reads into a 32-byte buffer from whichever fd it is given. It
// imports fd_read directly, which checkGuestImportsAllowed refuses — so it is
// instantiated through a bare linker rather than through the runtime's own load
// path. That is the point: it probes the WASI configuration itself, not the import
// gate layered above it.
const fdReadProbeWat = `(module
  (import "wasi_snapshot_preview1" "fd_read"
    (func $fd_read (param i32 i32 i32 i32) (result i32)))
  (memory (export "memory") 1)
  ;; iovec at 0: {buf = 64, len = 32}; bytes-read out-param at 16
  (func (export "read_fd") (param $fd i32) (result i32)
    (i32.store (i32.const 0) (i32.const 64))
    (i32.store (i32.const 4) (i32.const 32))
    (call $fd_read (local.get $fd) (i32.const 0) (i32.const 1) (i32.const 16)))
)`

// TestConfiguredWasiCannotBlockOnRead pins the configuration invariant the whole
// WASI surface rests on: no descriptor a guest can read is one that can make it
// wait. Today that holds because configureWasiLogPipes gives a guest nothing but
// the write-only log pipes, leaving stdin at wasmtime's default empty stream.
//
// It is guarded here rather than left implicit because it is one line away from
// being lost: an InheritStdin, a SetStdinFile or a PreopenDir added to
// configureWasiLogPipes would hand a guest something that can block, and a blocking
// read is beyond the execution bound entirely — epoch interruption only reaches
// guest code, so a guest parked in a host call holds execLock past its own
// deadline, past the caller's budget and past shutdown, stalling every application.
//
// The assertion is a wall-clock timeout, not an errno: what must never happen is
// waiting. Which errno comes back is wasmtime's business and may change.
func TestConfiguredWasiCannotBlockOnRead(t *testing.T) {
	// The default bound, not the short test one: a store carries an epoch deadline
	// from creation, and this test must not race it.
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	store := runtime.newModuleStore()
	cleanup, err := runtime.configureWasiLogPipes(context.Background(), common.NewApplicationId(1), store)
	require.NoError(t, err)
	defer cleanup()

	linker := wasmtime.NewLinker(runtime.engine)
	require.NoError(t, linker.DefineWasi())

	wasmBytes, err := wasmtime.Wat2Wasm(fdReadProbeWat)
	require.NoError(t, err)
	module, err := wasmtime.NewModule(runtime.engine, wasmBytes)
	require.NoError(t, err)
	instance, err := linker.Instantiate(store, module)
	require.NoError(t, err)
	readFd := instance.GetFunc(store, "read_fd")
	require.NotNil(t, readFd)

	// stdin, stdout, stderr, and the first fd beyond them: every descriptor a guest
	// could name without one having been preopened for it.
	for _, fd := range []int32{0, 1, 2, 3} {
		done := make(chan struct{})
		go func() {
			defer close(done)
			_, _ = readFd.Call(store, fd)
		}()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			// Deliberately fatal: the goroutine is wedged in a host call and cannot be
			// interrupted, which is precisely the condition being reported.
			t.Fatalf("fd_read on fd %d did not return: the WASI configuration exposes a "+
				"readable descriptor, so a guest can park the host thread past the "+
				"execution bound (check configureWasiLogPipes for InheritStdin, "+
				"SetStdinFile or PreopenDir)", fd)
		}
	}
}

// TestRequestBudgetBoundsTheWholeRequest is the core of the per-request budget: a
// request that declares a budget must be interrupted when THAT budget is spent,
// not when each individual operation's bound is.
//
// The runtime bound here is the 10s default and the request budget is a fraction of
// it, so the two are impossible to confuse: if operations were still armed from the
// runtime bound alone the guest would spin for ten seconds.
//
// The classification matters as much as the timing. Exhausting the request's own
// allowance is the guest's doing, so it must surface as CodeGuestExecutionTimeout,
// which is SIGNED — the request settles on-chain and the queue advances. It must
// NOT come back as CodeGuestExecutionCancelled, which is the transient,
// nothing-is-signed path reserved for the host giving up.
func TestRequestBudgetBoundsTheWholeRequest(t *testing.T) {
	const requestBudget = 300 * time.Millisecond

	// The default 10s bound, deliberately far above the request budget.
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	wasmBytes, err := wasmtime.Wat2Wasm(foreverSpinningProcessRequestWat)
	require.NoError(t, err)

	ctx := common.WithGuestExecutionBudget(context.Background(), requestBudget)

	start := time.Now()
	_, _, _, _, _, _, failure := runtime.ProcessRequest(
		ctx, common.NewApplicationId(1), ethCommon.Address{},
		common.Process, []byte("{}"), []byte("{}"), wasmBytes)
	elapsed := time.Since(start)

	require.NotNil(t, failure)
	require.Equal(t, apperrors.CodeGuestExecutionTimeout, failure.RequestError,
		"spending the request's own budget is the guest's doing and must be signed, "+
			"not reported as host-side cancellation")
	require.Less(t, elapsed, requestBudget*4,
		"the operation was armed from the runtime bound instead of the remaining request budget")
	require.Greater(t, elapsed, requestBudget/2,
		"interrupted far too early — the budget is not being honoured as the arming source")
}

// TestRequestBudgetIsSharedAcrossOperations pins the property the request budget
// exists for: consecutive operations DRAW DOWN one allowance instead of each
// receiving a fresh one. Without it, a request that loads a module, takes a deposit
// and then processes gets three full bounds, and their sum is what the manager
// waits through.
//
// Asserted on the armed tick count rather than on elapsed time, so it states the
// arithmetic directly and cannot be blurred by scheduling jitter.
func TestRequestBudgetIsSharedAcrossOperations(t *testing.T) {
	const (
		requestBudget = 2 * time.Second
		spent         = 700 * time.Millisecond
	)

	// Default 10s bound, so the budget — not the bound — is what binds.
	runtime := NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	ctx := common.WithGuestExecutionBudget(context.Background(), requestBudget)
	store := runtime.newModuleStore()

	first := runtime.beginGuestExecution(ctx, store)
	first.end()

	time.Sleep(spent)

	second := runtime.beginGuestExecution(ctx, store)
	second.end()

	require.Less(t, second.ticks, first.ticks,
		"the second operation was armed with a fresh allowance: the budget is per operation, "+
			"not per request, so their sum can outlast the caller")
	require.LessOrEqual(t, second.ticks, epochTicksFor(requestBudget-spent),
		"the second operation was armed with more than the budget had left")
}

// TestWithoutARequestBudgetOperationsUseTheFullBound is the other half: the budget
// is opt-in, and a caller that sets none must keep the previous behaviour rather
// than inherit an allowance that is already spent.
func TestWithoutARequestBudgetOperationsUseTheFullBound(t *testing.T) {
	runtime := newBoundedRuntime(t, 0)
	defer runtime.Close()

	store := runtime.newModuleStore()
	g := runtime.beginGuestExecution(context.Background(), store)
	defer g.end()

	require.Equal(t, epochTicksFor(testGuestTimeout), g.ticks,
		"with no request budget an operation must be armed from the configured bound")
}
