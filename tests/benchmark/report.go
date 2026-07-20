package benchmark

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// RunParams pins the parameters of a benchmark run so the report is
// self-describing (see THROUGHPUT_BENCHMARKS.md §1 — results are only
// comparable at equal parameters).
type RunParams struct {
	Scenario        string
	Implementation  string // e.g. "baseline"
	Requests        int
	BlockTime       time.Duration
	PollingInterval time.Duration
}

// FormatReport renders a saturation-mode run as markdown.
func FormatReport(p RunParams, m SaturationMetrics, tm StageTimings, wallClock time.Duration) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# Benchmark: %s (%s)\n\n", p.Scenario, p.Implementation)
	fmt.Fprintf(&b, "Mode: saturation. Requests: %d. Block time: %s. Polling interval: %s.\n\n",
		p.Requests, p.BlockTime, p.PollingInterval)
	b.WriteString("> App parameters are pinned (minimal WASM work, small state/payloads); these numbers do not measure the WASM dimension. Block-denominated metrics are the headline; wall-clock values depend on the emulated block time.\n\n")

	fmt.Fprintf(&b, "| Metric | Value |\n|---|---|\n")
	fmt.Fprintf(&b, "| Requests/block | **%.2f** |\n", m.RequestsPerBlock)
	fmt.Fprintf(&b, "| Completed / failed | %d / %d |\n", m.Completed, m.Failed)
	fmt.Fprintf(&b, "| Drain window (blocks) | %d (block %d → %d) |\n", m.BlocksSpanned, m.FirstBlock, m.LastBlock)
	fmt.Fprintf(&b, "| State-update txs | %d |\n", m.StateUpdateTxs)
	fmt.Fprintf(&b, "| Mean batch size | %.2f |\n", m.MeanBatchSize)
	fmt.Fprintf(&b, "| Batch size distribution | %s |\n", formatDist(m.BatchSizeDist))
	fmt.Fprintf(&b, "| Gas/request | %.0f |\n", m.GasPerRequest)
	fmt.Fprintf(&b, "| Latency p50 / p95 (blocks, saturation) | %.0f / %.0f |\n", m.LatencyBlocksP50, m.LatencyBlocksP95)
	fmt.Fprintf(&b, "| Wall clock (submission → last completion) | %s |\n", wallClock.Round(time.Millisecond))

	if len(tm.ExecutorRoundTripMs) > 0 || len(tm.MineWaitMs) > 0 {
		execShare := float64(tm.SumExecutorRoundTrip()) / float64(wallClock) * 100
		mineShare := float64(tm.SumMineWait()) / float64(wallClock) * 100
		fmt.Fprintf(&b, "| Executor round-trip total | %s (%.0f%% of wall clock, %d calls) |\n",
			tm.SumExecutorRoundTrip().Round(time.Millisecond), execShare, len(tm.ExecutorRoundTripMs))
		fmt.Fprintf(&b, "| Mine-wait total | %s (%.0f%% of wall clock, %d txs) |\n",
			tm.SumMineWait().Round(time.Millisecond), mineShare, len(tm.MineWaitMs))
	}

	if m.Failed > 0 {
		fmt.Fprintf(&b, "\n**WARNING: %d requests failed — this run is not a valid comparison point.**\n", m.Failed)
	}

	return b.String()
}

func formatDist(dist map[int]int) string {
	sizes := make([]int, 0, len(dist))
	for s := range dist {
		sizes = append(sizes, s)
	}
	sort.Ints(sizes)
	parts := make([]string, 0, len(sizes))
	for _, s := range sizes {
		parts = append(parts, fmt.Sprintf("%d×[%d req/tx]", dist[s], s))
	}
	return strings.Join(parts, ", ")
}
