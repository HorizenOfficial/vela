// Package benchmark contains the throughput benchmark harness for Vela
// (see docs/design/THROUGHPUT_BENCHMARKS.md). The pure computation lives in
// non-test files so it is built and vetted with the normal toolchain; the
// harness and the scenarios live in _test files and are gated behind the
// VELA_BENCH environment variable.
package benchmark

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"
)

// CompletionRecord is one request finalized on-chain, as observed from the
// RequestCompleted event and its transaction receipt.
type CompletionRecord struct {
	RequestID   [32]byte
	AppID       uint64
	Block       uint64 // block of the state-update tx that finalized it
	TxHash      string // hash of the state-update tx
	GasUsed     uint64 // gas used by that tx (same value on records sharing a tx)
	Failed      bool   // finalized with a non-zero error code
	SubmitBlock uint64 // block the submitRequest tx landed in
}

// SaturationMetrics is the saturation-mode result set. Block-denominated
// numbers are the headline: they are invariant to the emulated block time.
type SaturationMetrics struct {
	Completed        int         // successfully completed requests
	Failed           int         // requests finalized with an error
	FirstBlock       uint64      // earliest completion block
	LastBlock        uint64      // latest completion block
	BlocksSpanned    uint64      // LastBlock - FirstBlock + 1
	RequestsPerBlock float64     // Completed / BlocksSpanned
	StateUpdateTxs   int         // distinct state-update transactions
	MeanBatchSize    float64     // finalized requests per state-update tx
	BatchSizeDist    map[int]int // batch size -> number of txs with that size
	TotalGas         uint64      // gas summed once per distinct tx
	GasPerRequest    float64     // TotalGas / (Completed + Failed)
	LatencyBlocksP50 float64     // submit block -> completion block, nearest-rank
	LatencyBlocksP95 float64
}

// ComputeSaturationMetrics aggregates completion records into the
// saturation-mode metrics.
func ComputeSaturationMetrics(records []CompletionRecord) (SaturationMetrics, error) {
	if len(records) == 0 {
		return SaturationMetrics{}, fmt.Errorf("no completion records")
	}

	m := SaturationMetrics{
		FirstBlock:    records[0].Block,
		BatchSizeDist: make(map[int]int),
	}

	txGas := make(map[string]uint64)
	txCount := make(map[string]int)
	latencies := make([]float64, 0, len(records))

	for _, r := range records {
		if r.Failed {
			m.Failed++
		} else {
			m.Completed++
		}
		if r.Block < m.FirstBlock {
			m.FirstBlock = r.Block
		}
		if r.Block > m.LastBlock {
			m.LastBlock = r.Block
		}
		txGas[r.TxHash] = r.GasUsed
		txCount[r.TxHash]++
		latencies = append(latencies, float64(r.Block-r.SubmitBlock))
	}

	m.BlocksSpanned = m.LastBlock - m.FirstBlock + 1
	m.RequestsPerBlock = float64(m.Completed) / float64(m.BlocksSpanned)
	m.StateUpdateTxs = len(txGas)
	m.MeanBatchSize = float64(len(records)) / float64(m.StateUpdateTxs)
	for _, n := range txCount {
		m.BatchSizeDist[n]++
	}
	for _, g := range txGas {
		m.TotalGas += g
	}
	m.GasPerRequest = float64(m.TotalGas) / float64(len(records))
	m.LatencyBlocksP50 = percentileNearestRank(latencies, 50)
	m.LatencyBlocksP95 = percentileNearestRank(latencies, 95)

	return m, nil
}

// percentileNearestRank returns the nearest-rank percentile (p in [0,100]) of
// the values. Returns 0 for an empty slice.
func percentileNearestRank(values []float64, p int) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	rank := (p*len(sorted) + 99) / 100 // ceil(p/100 * n)
	if rank < 1 {
		rank = 1
	}
	if rank > len(sorted) {
		rank = len(sorted)
	}
	return sorted[rank-1]
}

// StageTimings holds the per-stage durations parsed from the manager's log
// file (the two instrumentation statements: executor round-trip in the
// manager, mine wait in the blockchain client).
type StageTimings struct {
	ExecutorRoundTripMs []int64
	MineWaitMs          []int64
}

var (
	executorRoundTripRe = regexp.MustCompile(`executor_roundtrip_ms=(\d+)`)
	mineWaitRe          = regexp.MustCompile(`mine_wait_ms=(\d+)`)
)

// ParseStageTimings extracts the instrumentation durations from raw log
// content. Lines without the known keys are ignored.
func ParseStageTimings(logContent string) StageTimings {
	var tm StageTimings
	for _, match := range executorRoundTripRe.FindAllStringSubmatch(logContent, -1) {
		v, err := strconv.ParseInt(match[1], 10, 64)
		if err == nil {
			tm.ExecutorRoundTripMs = append(tm.ExecutorRoundTripMs, v)
		}
	}
	for _, match := range mineWaitRe.FindAllStringSubmatch(logContent, -1) {
		v, err := strconv.ParseInt(match[1], 10, 64)
		if err == nil {
			tm.MineWaitMs = append(tm.MineWaitMs, v)
		}
	}
	return tm
}

// SumExecutorRoundTrip returns the total time spent in executor round-trips.
func (tm StageTimings) SumExecutorRoundTrip() time.Duration {
	return sumMs(tm.ExecutorRoundTripMs)
}

// SumMineWait returns the total time spent waiting for transaction inclusion.
func (tm StageTimings) SumMineWait() time.Duration {
	return sumMs(tm.MineWaitMs)
}

func sumMs(values []int64) time.Duration {
	var total int64
	for _, v := range values {
		total += v
	}
	return time.Duration(total) * time.Millisecond
}
