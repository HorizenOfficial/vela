package communication

import (
	"context"
	"time"
)

// Per-request execution budget: the caller tells the executor how much longer it
// intends to wait, and the executor shortens its own guest execution bound to fit.
//
// Why relative milliseconds and not a deadline timestamp: the executor runs inside
// a Nitro enclave whose clock is independent of the parent instance's, so an
// absolute time from the other side is not interpretable. Transit over vsock is
// sub-millisecond, so treating "remaining" as still-remaining on arrival is exact
// enough for a multi-second bound.
//
// Why it can only ever SHORTEN execution: the value crosses the TEE trust boundary,
// so the enclave keeps the last word. The executor arms its stores from its OWN
// configured timeout (see pkg/wasm) and this budget only ever cuts that short via
// the context. A caller claiming an hour of patience gets the configured bound; a
// caller claiming a millisecond gets ignored (see minUsableExecutionBudget).
//
// Why an interrupt caused by this budget is classified transient rather than as a
// signed on-chain failure: a signed failure charges the user MinFeePerRequest. If a
// host-supplied value could produce one, a compromised manager could burn every
// user's minimum fee by sending a tiny budget — turning griefing into a financial
// attack without ever forging a signature. Routing budget-driven interrupts to the
// retry path removes that capability structurally rather than by policy.
const (
	// executionBudgetMargin is reserved out of the caller's budget for everything the
	// executor does around the guest call itself: decrypting state, encrypting
	// events, signing the payload, and writing the response back. Without it the
	// guest could consume the entire budget and the caller would time out while a
	// perfectly good answer was being signed.
	executionBudgetMargin = 2 * time.Second

	// minUsableExecutionBudget is the least post-margin time worth starting work
	// with. Below it the request would be interrupted almost immediately and retried
	// on the next poll, forever — so the budget is ignored and the executor's own
	// bound applies instead. This is the guard against an absurd or hostile value.
	minUsableExecutionBudget = 1 * time.Second
)

// contextForExecutionBudget derives a context carrying the caller's execution
// budget. Returns ctx unchanged (with a no-op cancel) when no usable budget was
// supplied, so callers can always defer the returned cancel.
//
// Note this never extends: context.WithTimeout keeps the earlier of its own
// deadline and the parent's, so a generous budget cannot outlive the connection
// context it derives from.
func contextForExecutionBudget(ctx context.Context, budgetMs int64) (context.Context, context.CancelFunc) {
	if budgetMs <= 0 {
		return ctx, func() {}
	}

	usable := time.Duration(budgetMs)*time.Millisecond - executionBudgetMargin
	if usable < minUsableExecutionBudget {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, usable)
}
