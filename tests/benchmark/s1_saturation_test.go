package benchmark

import (
	"math/big"
	"testing"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
)

// TestBenchS1HotNormalApp — scenario S1 (THROUGHPUT_BENCHMARKS.md §3):
// one normal app, saturation mode. The queue is pre-filled with N small ETH
// deposits from a single user, then drained; the headline metric is
// requests/block over the drain window.
func TestBenchS1HotNormalApp(t *testing.T) {
	skipUnlessBench(t)

	n := envInt(t, "VELA_BENCH_N", 30)
	blockTime := time.Duration(envInt(t, "VELA_BENCH_BLOCK_MS", 1000)) * time.Millisecond
	setupTimeout := 120 * time.Second

	env := newBenchEnv(t, blockTime)
	defer func() { _ = env.suite.Cleanup() }()

	// Setup (excluded from measurement): deploy app, register one user.
	wasm := buildWasmApp(t, "simple")
	appID := env.deployApp(t, wasm, nil, nil, setupTimeout)
	user := env.newAppUser(t, appID, setupTimeout)
	execPub, err := env.suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// Measurement window starts after setup traffic has fully completed.
	fromBlock := env.currentBlock(t) + 1
	logOffset := env.logOffset(t)

	// Pre-fill: N deposits submitted without waiting for mining.
	start := time.Now()
	txs := make([]*ethTypes.Transaction, 0, n)
	for i := 0; i < n; i++ {
		req, err := env.ch.CreateDepositRequest(appID, commontestutil.GenerateRandomRequestID(), user, big.NewInt(1000), execPub)
		require.NoError(t, err)
		tx, err := env.suite.SubmitRequestNoWait(req)
		require.NoError(t, err)
		txs = append(txs, tx)
	}
	env.confirmSubmissions(t, txs)

	// Drain: generous bound — baseline needs a few blocks per request.
	drainTimeout := time.Duration(n)*10*blockTime + 2*time.Minute
	records := env.waitForCompletions(t, fromBlock, n, drainTimeout)
	wallClock := time.Since(start)

	metrics, err := ComputeSaturationMetrics(records)
	require.NoError(t, err)

	// Report before validity assertions so a failed run still leaves its
	// diagnostics (the report flags a non-zero failed count loudly).
	writeReport(t, RunParams{
		Scenario:        "S1-hot-normal-app",
		Implementation:  "baseline",
		Requests:        n,
		BlockTime:       blockTime,
		PollingInterval: time.Second,
	}, metrics, env.stageTimings(t, logOffset), wallClock)

	require.Equal(t, 0, metrics.Failed, "failed requests invalidate the comparison (see report)")
	env.suite.AssertNoStateUpdateErrors(t)
}
