package benchmark

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func rec(block, submitBlock uint64, tx string, gas uint64, failed bool) CompletionRecord {
	return CompletionRecord{
		Block:       block,
		SubmitBlock: submitBlock,
		TxHash:      tx,
		GasUsed:     gas,
		Failed:      failed,
	}
}

func TestComputeSaturationMetrics_Empty(t *testing.T) {
	_, err := ComputeSaturationMetrics(nil)
	require.Error(t, err)
}

func TestComputeSaturationMetrics_SingleRequestPerTx(t *testing.T) {
	// Today's implementation: one state-update tx per request.
	records := []CompletionRecord{
		rec(10, 5, "0xa", 100_000, false),
		rec(12, 5, "0xb", 110_000, false),
		rec(14, 5, "0xc", 120_000, false),
	}
	m, err := ComputeSaturationMetrics(records)
	require.NoError(t, err)

	require.Equal(t, 3, m.Completed)
	require.Equal(t, 0, m.Failed)
	require.EqualValues(t, 10, m.FirstBlock)
	require.EqualValues(t, 14, m.LastBlock)
	require.EqualValues(t, 5, m.BlocksSpanned)
	require.InDelta(t, 0.6, m.RequestsPerBlock, 1e-9)
	require.Equal(t, 3, m.StateUpdateTxs)
	require.InDelta(t, 1.0, m.MeanBatchSize, 1e-9)
	require.Equal(t, map[int]int{1: 3}, m.BatchSizeDist)
	require.EqualValues(t, 330_000, m.TotalGas)
	require.InDelta(t, 110_000, m.GasPerRequest, 1e-9)
}

func TestComputeSaturationMetrics_BatchedTxs(t *testing.T) {
	// Future batch implementation: several completions share one tx; gas must
	// be counted once per tx, batch sizes from the per-tx grouping.
	records := []CompletionRecord{
		rec(10, 5, "0xa", 100_000, false),
		rec(11, 5, "0xb", 150_000, false),
		rec(11, 5, "0xb", 150_000, false),
		rec(12, 5, "0xc", 200_000, false),
		rec(12, 5, "0xc", 200_000, false),
		rec(12, 5, "0xc", 200_000, false),
	}
	m, err := ComputeSaturationMetrics(records)
	require.NoError(t, err)

	require.Equal(t, 6, m.Completed)
	require.EqualValues(t, 3, m.BlocksSpanned)
	require.InDelta(t, 2.0, m.RequestsPerBlock, 1e-9)
	require.Equal(t, 3, m.StateUpdateTxs)
	require.InDelta(t, 2.0, m.MeanBatchSize, 1e-9)
	require.Equal(t, map[int]int{1: 1, 2: 1, 3: 1}, m.BatchSizeDist)
	require.EqualValues(t, 450_000, m.TotalGas)
	require.InDelta(t, 75_000, m.GasPerRequest, 1e-9)
}

func TestComputeSaturationMetrics_FailedCountedSeparately(t *testing.T) {
	records := []CompletionRecord{
		rec(10, 5, "0xa", 100_000, false),
		rec(11, 5, "0xb", 100_000, true),
		rec(12, 5, "0xc", 100_000, false),
	}
	m, err := ComputeSaturationMetrics(records)
	require.NoError(t, err)

	require.Equal(t, 2, m.Completed)
	require.Equal(t, 1, m.Failed)
	// Requests/block counts successful completions only.
	require.InDelta(t, 2.0/3.0, m.RequestsPerBlock, 1e-9)
}

func TestComputeSaturationMetrics_LatencyBlocks(t *testing.T) {
	// All submitted at block 5, completed at 10..15 → latencies 5,6,7,8,9,10.
	records := []CompletionRecord{
		rec(10, 5, "0xa", 1, false),
		rec(11, 5, "0xb", 1, false),
		rec(12, 5, "0xc", 1, false),
		rec(13, 5, "0xd", 1, false),
		rec(14, 5, "0xe", 1, false),
		rec(15, 5, "0xf", 1, false),
	}
	m, err := ComputeSaturationMetrics(records)
	require.NoError(t, err)

	// Nearest-rank percentiles: p50 = 3rd of 6 = 7, p95 = 6th of 6 = 10.
	require.InDelta(t, 7, m.LatencyBlocksP50, 1e-9)
	require.InDelta(t, 10, m.LatencyBlocksP95, 1e-9)
}

func TestParseStageTimings(t *testing.T) {
	log := `{"level":"info","time":"...","message":"Manager: executor round-trip for request 0xabc: executor_roundtrip_ms=42"}
{"level":"info","time":"...","message":"BlockChainClient: tx 0xdef mined: mine_wait_ms=900"}
{"level":"info","time":"...","message":"Manager: unrelated line"}
{"level":"info","time":"...","message":"Manager: executor round-trip for request 0x123: executor_roundtrip_ms=58"}
{"level":"info","time":"...","message":"BlockChainClient: tx 0x456 mined: mine_wait_ms=1100"}`

	tm := ParseStageTimings(log)
	require.Equal(t, []int64{42, 58}, tm.ExecutorRoundTripMs)
	require.Equal(t, []int64{900, 1100}, tm.MineWaitMs)
	require.EqualValues(t, 100, tm.SumExecutorRoundTrip().Milliseconds())
	require.EqualValues(t, 2000, tm.SumMineWait().Milliseconds())
}
