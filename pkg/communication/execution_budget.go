package communication

import (
	"context"
	"fmt"
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
	// ExecutionBudgetMargin is reserved out of the caller's budget for everything the
	// executor does around the guest call itself: decrypting state, encrypting
	// events, signing the payload, and writing the response back. Without it the
	// guest could consume the entire budget and the caller would time out while a
	// perfectly good answer was being signed.
	//
	// Exported because pkg/executor validates its configured guest bound against
	// it during the handshake, using the request timeout the manager reports there
	// (GetKeysetRecoveryResponseData.RequestTimeoutMs): the guest bound, plus its own
	// safety margin for epoch overshoot and pre-arming work, must be below
	// timeout - margin, or the budget wins every race and no runaway guest is ever
	// settled on-chain. Erring large is deliberate — too small discards a request
	// that actually succeeded, whereas too large only shortens an allowance nothing
	// depends on.
	ExecutionBudgetMargin = 2 * time.Second

	// minUsableExecutionBudget is the least post-margin time worth starting work
	// with. Below it the request would be interrupted almost immediately and retried
	// on the next poll, forever — so the budget is ignored and the executor's own
	// bound applies instead. This is the guard against an absurd or hostile value.
	minUsableExecutionBudget = 1 * time.Second

	// MaxExecutionBudgetMs is the largest budget worth honouring, in milliseconds. An
	// hour is already orders of magnitude beyond any plausible request timeout, so
	// anything above it is a bug or a hostile value either way.
	//
	// It exists for a second, sharper reason: budgetMs feeds
	// `time.Duration(budgetMs) * time.Millisecond`, which overflows int64 above about
	// 9.2e12 ms. Most wrapped values land negative and are harmlessly ignored, but a
	// crafted one wraps to a small POSITIVE duration and would silently cut a
	// legitimate request short. Capping the input keeps the multiplication away from
	// that entirely, rather than depending on where the wrap happens to land. The
	// field crosses the TEE boundary, so it must not be trusted to be sane.
	MaxExecutionBudgetMs = 60 * 60 * 1000
)

// validateExecutionBudget checks the ExecutionBudgetMs field shared by the process,
// deploy and batch request types.
//
// Zero is legitimate and must stay accepted: the field is `omitempty`, so a peer
// that does not set it sends nothing and it decodes to zero — rejecting it would
// fail every request from an older manager. Zero simply means no caller deadline,
// leaving the executor's own guest bound as the only limit.
func validateExecutionBudget(budgetMs int64) error {
	if budgetMs < 0 {
		return fmt.Errorf("executionBudgetMs must not be negative, got %d", budgetMs)
	}
	if budgetMs > MaxExecutionBudgetMs {
		return fmt.Errorf("executionBudgetMs must not exceed %d ms, got %d", int64(MaxExecutionBudgetMs), budgetMs)
	}
	return nil
}

// contextForExecutionBudget derives a context carrying the caller's execution
// budget. Returns ctx unchanged (with a no-op cancel) when no usable budget was
// supplied, so callers can always defer the returned cancel.
//
// Note this never extends: context.WithTimeout keeps the earlier of its own
// deadline and the parent's, so a generous budget cannot outlive the connection
// context it derives from.
func contextForExecutionBudget(ctx context.Context, budgetMs int64) (context.Context, context.CancelFunc) {
	// Anything outside the sane range is treated as "not supplied": the executor's own
	// bound then applies, which is always a safe fallback. The upper check also keeps
	// the multiplication below from overflowing — see MaxExecutionBudgetMs.
	if budgetMs <= 0 || budgetMs > MaxExecutionBudgetMs {
		return ctx, func() {}
	}

	usable := time.Duration(budgetMs)*time.Millisecond - ExecutionBudgetMargin
	if usable < minUsableExecutionBudget {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, usable)
}
