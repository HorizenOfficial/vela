package communication

import (
	"context"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/stretchr/testify/require"
)

// The execution budget is the link between the caller's patience and the guest
// execution bound inside the enclave (see PLAN_epoch_interruption.md, decision D5).
// It can only ever SHORTEN execution: it arrives from outside the TEE, so the
// enclave's own configured timeout stays the ceiling, and an interrupt caused by a
// budget is classified as host-side abandonment — transient, never a signed
// on-chain failure that would charge the user a fee.

func TestExecutionBudgetNotSupplied(t *testing.T) {
	ctx, cancel := contextForExecutionBudget(context.Background(), 0)
	defer cancel()

	_, hasDeadline := ctx.Deadline()
	require.False(t, hasDeadline,
		"an absent budget must leave the context alone so the enclave's own timeout applies")
}

func TestExecutionBudgetAppliesDeadline(t *testing.T) {
	const budgetMs = 30_000

	start := time.Now()
	ctx, cancel := contextForExecutionBudget(context.Background(), budgetMs)
	defer cancel()

	deadline, hasDeadline := ctx.Deadline()
	require.True(t, hasDeadline, "a usable budget must become a context deadline")

	remaining := deadline.Sub(start)
	require.Less(t, remaining, time.Duration(budgetMs)*time.Millisecond,
		"the deadline must leave the executor margin to encrypt and sign the response")
	require.Greater(t, remaining, time.Duration(budgetMs)*time.Millisecond-executionBudgetMargin-time.Second,
		"the margin must not eat more than it needs")
}

// TestExecutionBudgetIgnoresUnusableValues covers the hostile case: the budget
// crosses the trust boundary, so a compromised manager could send a value that is
// impossible to serve. Honouring it would make every request fail and be retried
// forever. The enclave falls back to its own timeout instead.
func TestExecutionBudgetIgnoresUnusableValues(t *testing.T) {
	for _, budgetMs := range []int64{1, 10, 500, -1, -30_000} {
		ctx, cancel := contextForExecutionBudget(context.Background(), budgetMs)
		_, hasDeadline := ctx.Deadline()
		require.False(t, hasDeadline,
			"budget %d ms leaves nothing usable and must be ignored, not applied", budgetMs)
		cancel()
	}
}

// TestExecutionBudgetNeverExtends is the security property: the budget is a
// shortening-only input. A caller asking for more time than the enclave allows
// cannot widen the bound, because the enclave arms its store from its own
// configured timeout and only the context can cut that short.
func TestExecutionBudgetNeverExtends(t *testing.T) {
	parent, parentCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer parentCancel()

	// A caller claiming an hour of patience.
	ctx, cancel := contextForExecutionBudget(parent, 3_600_000)
	defer cancel()

	deadline, hasDeadline := ctx.Deadline()
	require.True(t, hasDeadline)
	require.LessOrEqual(t, time.Until(deadline), 5*time.Second,
		"a budget must never push the deadline beyond what the parent context already allows")
}

// newValidatableRequest returns a request that passes common.Request.Validate, so
// these tests fail only on the budget field they are actually about.
func newValidatableRequest() *common.Request {
	return &common.Request{
		ApplicationID: common.NewApplicationId(1),
		RequestType:   common.Process,
		Timestamp:     common.NewBig(1),
		AssetAmount:   common.NewBig(0),
		MaxFeeValue:   common.NewBig(10),
	}
}

func TestExecutionBudgetValidation(t *testing.T) {
	t.Run("negative process budget rejected", func(t *testing.T) {
		d := &ProcessRequestData{Request: newValidatableRequest(), ExecutionBudgetMs: -1}
		require.Error(t, d.Validate())
	})

	t.Run("zero process budget accepted", func(t *testing.T) {
		d := &ProcessRequestData{Request: newValidatableRequest(), ExecutionBudgetMs: 0}
		require.NoError(t, d.Validate())
	})

	t.Run("negative deploy budget rejected", func(t *testing.T) {
		d := &DeployAppRequestData{Request: newValidatableRequest(), ExecutionBudgetMs: -1}
		require.Error(t, d.Validate())
	})
}
