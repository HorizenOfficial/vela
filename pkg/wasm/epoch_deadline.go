package wasm

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	wasmtime "github.com/bytecodealliance/wasmtime-go/v47"
)

// Wall-clock bound on guest execution, built on wasmtime's epoch interruption.
//
// Why epochs and not something else: wasmtime's alternatives are fuel (bounds
// *work*, not time, and a guest blocked in a host call burns none) and a watchdog
// thread killing the process (unacceptable inside the enclave). Epoch
// interruption costs a compare-and-branch at loop back-edges and function
// entries, so a healthy guest pays essentially nothing.
//
// How it fits together:
//
//   - The engine is built with SetEpochInterruption(true) (see newPinnedEngine),
//     which makes the compiler emit those checks.
//   - One goroutine per runtime bumps the engine epoch every epochTickInterval.
//   - Before any guest code runs on a store, that store is armed with a deadline
//     of N ticks (beginGuestExecution). Exceeding it traps the guest with
//     TrapCode Interrupt.
//   - A watcher goroutine per operation forces that same interrupt early when the
//     host stops caring about the result (context cancelled, runtime closing).
//
// Two measured properties of wasmtime-go v47 drive the shape of this and must not
// be "simplified" away without re-measuring:
//
//  1. With epoch interruption enabled, a store that was NEVER armed traps
//     immediately — its default deadline is already reached. Enabling the engine
//     flag and arming the stores are therefore one indivisible change: every path
//     that runs guest code has to arm, or it fails instantly and permanently.
//  2. The deadline is ABSOLUTE (epoch-at-arming + N), not a per-call allowance, so
//     idle wall-clock time between calls consumes it. Arming has to happen per
//     operation, immediately before entering the guest — arming a store once when
//     it is created would interrupt the next request that arrives late enough.
const (
	// epochTickInterval is the epoch-bump period, and therefore the granularity of
	// every deadline: a guest is interrupted at the first tick after its budget is
	// spent, so the effective timeout overshoots by up to one interval. 100 ms keeps
	// that overshoot negligible against the multi-second timeouts this bounds, at
	// the cost of one cheap atomic increment per interval per runtime.
	// Defined in pkg/common so pkg/executor can size its handshake safety margin from
	// the real value without linking libwasmtime, exactly as for the timeout bounds
	// below. Changing it there changes it here.
	epochTickInterval = common.GuestExecutionEpochTick

	// defaultGuestExecutionTimeout and maxGuestExecutionTimeout mirror the bounds
	// pkg/executor validates EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS against; see the
	// pkg/common constants for the reasoning behind the values.
	defaultGuestExecutionTimeout = common.DefaultGuestExecutionTimeoutMs * time.Millisecond
	maxGuestExecutionTimeout     = common.MaxGuestExecutionTimeoutMs * time.Millisecond
)

// epochTicksFor converts a timeout into the epoch-deadline tick count to arm a
// store with.
//
// The +1 is not slack, it is correctness: the ticker's phase is independent of
// when a store is armed, so the first bump after arming can land immediately.
// Without it a guest given 400 ms could be interrupted after ~300 ms. With it,
// a guest is never interrupted before its configured allowance, and the overshoot
// stays bounded by one tick interval.
//
// timeout must be positive; callers get that from SetGuestExecutionTimeout, which
// clamps. A zero deadline would mean "already expired" to wasmtime, not
// "unbounded", and would fail every call.
func epochTicksFor(timeout time.Duration) uint64 {
	ticks := (timeout + epochTickInterval - 1) / epochTickInterval // ceil
	if ticks < 1 {
		ticks = 1
	}
	return uint64(ticks) + 1
}

// startEpochTicker launches the goroutine that advances the engine epoch, which
// is what makes armed deadlines expire. It is stopped by closing r.closing and
// joined via r.epochTickerDone.
//
// The engine is passed in rather than read from the runtime inside the loop:
// Close sets r.engine to nil, and reading the field concurrently would be a data
// race. Close joins this goroutine BEFORE releasing the engine, so the captured
// pointer stays valid for the goroutine's whole life — IncrementEpoch on a closed
// engine panics.
func (r *WasmtimeRuntime) startEpochTicker(engine *wasmtime.Engine) {
	go func() {
		defer close(r.epochTickerDone)

		ticker := time.NewTicker(epochTickInterval)
		defer ticker.Stop()

		for {
			select {
			case <-r.closing:
				return
			case <-ticker.C:
				engine.IncrementEpoch()
			}
		}
	}()
}

// SetGuestExecutionTimeout updates the wall-clock bound applied to guest
// operations started from now on. A non-positive value selects the default and a
// value above the ceiling is clamped to it, both reported.
//
// Clamps rather than rejects, like SetMaxGuestMemoryBytes: startup validation in
// pkg/executor.Config.Validate is where a bad value fails loudly, whereas a late
// admin-driven change must not be able to take a running enclave down. In
// particular a non-positive value must never reach the store: wasmtime reads a
// zero deadline as "already expired", so it would fail every subsequent call
// rather than removing the bound.
func (r *WasmtimeRuntime) SetGuestExecutionTimeout(timeout time.Duration) {
	switch {
	case timeout <= 0:
		r.log.Warn("Runtime: requested guest execution timeout %v is not positive, using default %v", timeout, defaultGuestExecutionTimeout)
		timeout = defaultGuestExecutionTimeout
	case timeout > maxGuestExecutionTimeout:
		r.log.Warn("Runtime: requested guest execution timeout %v above ceiling, using %v", timeout, maxGuestExecutionTimeout)
		timeout = maxGuestExecutionTimeout
	}

	r.log.Info("Runtime: Setting guest execution timeout from %v to %v",
		time.Duration(r.guestExecutionTimeoutNs.Load()), timeout)
	r.guestExecutionTimeoutNs.Store(int64(timeout))
}

// GetGuestExecutionTimeout returns the wall-clock bound applied to new guest
// operations.
func (r *WasmtimeRuntime) GetGuestExecutionTimeout() time.Duration {
	return time.Duration(r.guestExecutionTimeoutNs.Load())
}

// guestExecution is the armed-and-watched window around one guest operation: the
// store carries the deadline, and a watcher goroutine forces the interrupt early
// if the host stops waiting for the result. One operation may make several guest
// calls (allocate, the export itself, deallocate) and they share the window,
// which is the right granularity — what is being protected is the executor's
// ability to answer the manager at all, not any individual call.
//
// Always paired with a deferred end(), which joins the watcher.
type guestExecution struct {
	runtime *WasmtimeRuntime
	log     interface{ Warn(string, ...interface{}) }

	// done is closed by end() to release the watcher when the operation finishes
	// on its own; stopped is closed by the watcher as it exits, so end() can join.
	done    chan struct{}
	stopped chan struct{}

	// cancelled records that the interrupt was forced by the host rather than by
	// the deadline expiring. It decides the classification, so it is written before
	// the epoch is bumped and read after the guest call returns.
	cancelled atomic.Bool

	// ticks is what the store was armed with, and therefore how many bumps are
	// enough to guarantee the deadline is passed when forcing an interrupt.
	ticks uint64

	// cleanupInterrupted records an interrupt during a fire-and-forget cleanup call
	// — a deallocate, whose error the host otherwise drops. Without it, a guest that
	// returns a valid result and then hangs in deallocate would stay cached with an
	// allocator that trapped mid-free, which is the damaged heap evictFaultedModule
	// exists to drop. See noteCleanupCall.
	cleanupInterrupted atomic.Bool
}

// noteCleanupCall records the outcome of a cleanup call whose error is not
// propagated to the caller. Only an interrupt is recorded: any other trap from a
// deallocate is a guest bug that the checked call sites already classify.
func (g *guestExecution) noteCleanupCall(err error) {
	if err != nil && isEpochInterrupt(err) {
		g.cleanupInterrupted.Store(true)
	}
}

// cleanupWasInterrupted reports whether any cleanup call in this window was
// interrupted, meaning the instance must not serve another request even though the
// operation itself may have succeeded.
func (g *guestExecution) cleanupWasInterrupted() bool {
	return g.cleanupInterrupted.Load()
}

// guestExecutionAllowance returns how long the operation about to start may run:
// the configured bound, or whatever is left of the request's shared budget when
// that is less (see common.WithGuestExecutionBudget).
//
// The distinction that matters is between "no budget was set" and "the budget is
// spent". The first means this operation is not part of a request that bounds it as
// a whole, so the configured bound applies on its own, exactly as before. The
// second means the request has already had its full allowance across earlier
// operations, and the answer is a non-positive duration — epochTicksFor floors that
// at one tick, so the guest is interrupted at the first epoch check rather than
// being handed another full bound. Either way the interrupt is the guest exhausting
// its own allowance, which is signed; the caller's budget is a separate mechanism
// that works through cancellation and stays transient.
func (r *WasmtimeRuntime) guestExecutionAllowance(ctx context.Context) time.Duration {
	bound := r.GetGuestExecutionTimeout()

	remaining, bounded := common.GuestExecutionBudgetRemaining(ctx)
	if bounded && remaining < bound {
		return remaining
	}
	return bound
}

// beginGuestExecution arms store with the operation's allowance and starts the
// watcher. The caller MUST defer end().
//
// Must be called with execLock held — see the guestExecution watcher for why that
// is what makes the captured engine pointer safe.
func (r *WasmtimeRuntime) beginGuestExecution(ctx context.Context, store *wasmtime.Store) *guestExecution {
	timeout := r.guestExecutionAllowance(ctx)
	ticks := epochTicksFor(timeout)

	// Arming is per operation, never once per store: the deadline is absolute, so a
	// module that sat in the cache while several ticks passed would otherwise be
	// interrupted on its next, perfectly healthy, request.
	store.SetEpochDeadline(ticks)

	g := &guestExecution{
		runtime: r,
		log:     r.log,
		done:    make(chan struct{}),
		stopped: make(chan struct{}),
		ticks:   ticks,
	}

	// Captured rather than read inside the goroutine: Close sets r.engine to nil,
	// and this goroutine outlives neither execLock nor the operation. end() joins it
	// before the operation returns and releases execLock, and Close can only reach
	// the engine after acquiring execLock, so the pointer cannot be freed underneath.
	engine := r.engine

	go func() {
		defer close(g.stopped)

		select {
		case <-g.done:
			// Finished on its own, nothing to interrupt. The common case.
			return
		case <-ctx.Done():
		case <-r.closing:
		}

		// Host-side abandonment: classified as transient, so the request is retried
		// rather than signed as a permanent on-chain failure. Set before bumping, so
		// the flag is visible to whoever reads the resulting trap.
		g.cancelled.Store(true)

		// Bump by the full armed budget: the guest may have consumed any part of it
		// already, and this is guaranteed to carry the epoch past the deadline
		// regardless. Extra bumps are harmless — they only expire deadlines that are
		// already armed, and this store's is about to be abandoned anyway.
		for i := uint64(0); i < g.ticks; i++ {
			engine.IncrementEpoch()
		}
	}()

	return g
}

// end releases the watcher and waits for it to exit. Deferred immediately after
// beginGuestExecution, and — where the operation also holds execLock — deferred
// AFTER the execLock release so that LIFO joins the watcher first. That ordering
// is what lets the watcher touch the engine without a lock of its own.
func (g *guestExecution) end() {
	close(g.done)
	<-g.stopped
}

// deallocateGuest frees a guest buffer, recording an interrupt so the module can be
// evicted afterwards. Replaces the fire-and-forget
// `_, _ = module.deallocate.Call(...)` closures: the outcome is still not
// propagated to the caller (a failed free must not turn a good result into a failed
// request), but it is no longer discarded either.
//
// Safe to defer directly; a nil deallocate export or module is a no-op, matching
// the nil checks the call sites used to do inline.
func (r *WasmtimeRuntime) deallocateGuest(g *guestExecution, module *ApplicationModule, ptr int32, size int32) {
	if module == nil || module.deallocate == nil || module.store == nil || ptr == 0 {
		return
	}
	_, err := module.deallocate.Call(module.store, ptr, size)
	if err != nil {
		g.noteCleanupCall(err)
		r.log.Warn("Wasmtime Runtime: failed to deallocate guest buffer at %d (%d bytes): %v", ptr, size, err)
	}
}

// isEpochInterrupt reports whether err is (or wraps) a guest interrupted by an
// epoch deadline, as opposed to any other trap.
//
// Unwraps, so it still recognises the interrupt through the layers writeToMemory
// adds (fmt.Errorf with %w, then guestFaultError).
func isEpochInterrupt(err error) bool {
	var trap *wasmtime.Trap
	if !errors.As(err, &trap) {
		return false
	}
	code := trap.Code()
	return code != nil && *code == wasmtime.Interrupt
}

// interruptSentinel maps an interrupted guest call to the sentinel describing why,
// or returns nil if err is not an interrupt.
//
// The full wasmtime error, backtrace included, is logged here rather than
// propagated: what the caller gets back is a short, stable string, because it ends
// up in the signed on-chain ErrorMsg (truncated to 100 characters). A wasm
// backtrace would both overflow that budget and couple a signed value to
// wasmtime's error formatting.
//
// KNOWN RACE, deliberately not synchronised. g.cancelled is read after the guest
// call has already returned, so a deadline expiry that coincides — within
// microseconds — with a cancellation is reported as cancelled, i.e. transient,
// when it was really a timeout. It is left alone because it self-corrects and
// costs at most one wasted retry: the retry is a fresh call, and a cancellation
// arriving again in that same window is uncorrelated with the deadline, so the
// second attempt classifies as a timeout and is signed. The failure mode is a
// slightly delayed on-chain settlement, never a permanent stall.
//
// Making it deterministic would mean latching the reason at the moment the
// interrupt is forced (a generation counter the call site compares), which is more
// machinery than a one-retry delay justifies.
//
// The "self-corrects" argument depends on the two deadlines not being correlated,
// which in turn depends on the enclave's own bound expiring first. That is enforced
// in two places, and if either goes the argument above stops holding and this
// becomes a permanent stall rather than a delayed settlement: one request consumes
// at most ONE bound however many guest operations it makes
// (common.WithGuestExecutionBudget — before that, their sum could outlast the
// caller), and the bound plus its safety margin is validated at handshake to fit
// inside the caller's budget (pkg/executor.Config.checkGuestBoundFitsCallerBudget,
// which runs there rather than at start-up because only the manager knows its own
// timeout).
func (g *guestExecution) interruptSentinel(err error) error {
	if !isEpochInterrupt(err) {
		return nil
	}

	if g.cancelled.Load() {
		g.log.Warn("Wasmtime Runtime: guest execution interrupted, host stopped waiting for the result: %v", err)
		return apperrors.ErrGuestExecutionCancelled
	}

	g.log.Warn("Wasmtime Runtime: guest execution exceeded the %v limit and was interrupted: %v",
		g.runtime.GetGuestExecutionTimeout(), err)
	return apperrors.ErrGuestExecutionTimeout
}

// failure classifies an error from a guest call into a *RequestFailure.
//
// An interrupt takes its own code and a fixed message; anything else keeps the
// call site's own code and message byte for byte, so this can be dropped in at
// existing sites without changing what they report. (Mapping the OTHER trap codes to
// stable messages is a separate, tracked task.)
func (g *guestExecution) failure(err error, fallback apperrors.FailureCode, format string, args ...interface{}) *apperrors.RequestFailure {
	if sentinel := g.interruptSentinel(err); sentinel != nil {
		return apperrors.New(interruptCode(sentinel), sentinel.Error())
	}
	return apperrors.New(fallback, fmt.Sprintf(format, args...))
}

// wrapErr is failure's counterpart for the entry points that return a plain error
// (Deploy and the module-load path). On an interrupt it wraps the sentinel, so
// callers can test it with errors.Is; otherwise it wraps the original error,
// producing exactly the message the call site produced before.
//
// Takes a fixed context string rather than a format: a printf-style wrapper with a
// non-wrapping branch cannot support %w (go vet rejects it), and silently
// degrading %w to %v would drop the error chains that isGuestTrap and errGuestFault
// depend on.
func (g *guestExecution) wrapErr(err error, context string) error {
	if sentinel := g.interruptSentinel(err); sentinel != nil {
		return fmt.Errorf("%s: %w", context, sentinel)
	}
	return fmt.Errorf("%s: %w", context, err)
}

// interruptCode maps an execution-bound sentinel to its on-chain failure code.
func interruptCode(sentinel error) apperrors.FailureCode {
	if errors.Is(sentinel, apperrors.ErrGuestExecutionCancelled) {
		return apperrors.CodeGuestExecutionCancelled
	}
	return apperrors.CodeGuestExecutionTimeout
}

// executionBoundFailure classifies an error that may wrap an execution-bound
// sentinel — used where the interrupt happened deeper down and has already been
// converted to a sentinel (a module load that hit a spinning start section, say),
// so the trap itself is no longer available to inspect.
//
// Returns ok=false when err is unrelated, leaving the caller's own
// classification in place.
func executionBoundFailure(err error) (*apperrors.RequestFailure, bool) {
	switch {
	case errors.Is(err, apperrors.ErrGuestExecutionCancelled):
		return apperrors.New(apperrors.CodeGuestExecutionCancelled, apperrors.ErrGuestExecutionCancelled.Error()), true
	case errors.Is(err, apperrors.ErrGuestExecutionTimeout):
		return apperrors.New(apperrors.CodeGuestExecutionTimeout, apperrors.ErrGuestExecutionTimeout.Error()), true
	}
	return nil, false
}
