package common

import (
	"context"
	"time"
)

// The per-request guest execution budget: how much guest wall-clock time ONE
// request may consume in total, across every guest operation it triggers.
//
// Why a budget per request and not per operation. The execution bound is armed
// per guest operation, and one request opens several: loading the module runs its
// start section, a request carrying a deposit calls the guest twice, and each of
// those is armed with the full bound. Their sum is what the manager waits for, so
// with a 10s bound a single request could occupy the executor for 30s while the
// manager gave up at 28s. Whichever deadline wins decides how the interrupt is
// classified: the enclave's own bound expiring is a signed on-chain failure, so
// the request settles and the FIFO queue advances, whereas the caller giving up is
// host-side abandonment — nothing is signed, and the identical request is redelivered
// on the next poll, forever. Bounding the request as a whole is what keeps the
// enclave the one that wins.
//
// Why the deadline is carried as a context VALUE rather than by
// context.WithDeadline. A real deadline would cancel the context, and cancellation
// is precisely the classification that must not happen here: it is the transient,
// nothing-is-signed path. This value only tells the runtime how much of the bound
// is left to arm the next operation with; expiry surfaces as the guest exhausting
// its own allowance, which is signed.
//
// Why the value derives from the executor's OWN configured bound and never from
// anything the manager sends. A signed failure charges the request's author
// MinFeePerRequest. If a host-supplied number could shorten this budget, a
// compromised manager could make every request expire and burn users' fees without
// forging a signature. The caller's execution budget still bounds the request, but
// it is applied through context cancellation and stays transient. See
// pkg/communication/execution_budget.go.
//
// It lives in pkg/common for the same reason DefaultGuestExecutionTimeoutMs does:
// pkg/executor sets it and pkg/wasm reads it, and pkg/executor must not gain a
// dependency on pkg/wasm (and libwasmtime with it) to do so.

type guestBudgetKeyType struct{}

// guestBudgetKey is an unexported zero-size type, so no other package can collide
// with it or read the value except through the accessors below.
var guestBudgetKey guestBudgetKeyType

// WithGuestExecutionBudget marks the start of one request's guest execution budget,
// returning a context carrying the instant at which it is exhausted. Call it once
// per request, before any guest work: every operation started under the returned
// context shares the one budget.
//
// A non-positive budget is ignored, returning ctx unchanged, so that a caller which
// has not configured one gets the previous behaviour (each operation bounded on its
// own) rather than a budget that is already spent.
func WithGuestExecutionBudget(ctx context.Context, budget time.Duration) context.Context {
	if budget <= 0 {
		return ctx
	}
	return context.WithValue(ctx, guestBudgetKey, time.Now().Add(budget))
}

// GuestExecutionBudgetRemaining reports how much of the request's guest budget is
// left. ok is false when no budget was set, which means "not bounded as a request"
// — the caller then applies its own per-operation bound unchanged.
//
// The remaining value may be zero or negative, meaning the request has already
// spent its whole allowance; that is a distinct case from ok == false and callers
// must treat it as exhausted rather than as unbounded.
func GuestExecutionBudgetRemaining(ctx context.Context) (time.Duration, bool) {
	deadline, ok := ctx.Value(guestBudgetKey).(time.Time)
	if !ok {
		return 0, false
	}
	return time.Until(deadline), true
}
