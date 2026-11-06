package benchmark

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/HorizenOfficial/vela/pkg/common"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
)

// TestBenchS3Mixed — scenario S3 (THROUGHPUT_BENCHMARKS.md §3): 10 apps,
// 5 normal + 5 trigger, requests uniformly round-robined across them.
// With the current global FIFO the apps interleave in the queue — this is the
// workload that discriminates the multi-app fetch strategies (2a/2b/2c).
func TestBenchS3Mixed(t *testing.T) {
	skipUnlessBench(t)

	const appsPerKind = 5
	n := envInt(t, "VELA_BENCH_N", 30)
	require.Zero(t, n%(2*appsPerKind), "VELA_BENCH_N must be a multiple of %d for uniform distribution", 2*appsPerKind)
	blockTime := time.Duration(envInt(t, "VELA_BENCH_BLOCK_MS", 1000)) * time.Millisecond
	setupTimeout := 240 * time.Second

	env := newBenchEnv(t, blockTime)
	defer func() { _ = env.suite.Cleanup() }()

	execPub, err := env.suite.GetExecutorCommunicationKey()
	require.NoError(t, err)

	// One user for the normal apps, one firing user + one fee collector shared
	// across the trigger apps (keys associated per app below).
	newAccount := func() ethCommon.Address {
		addr, secp, err := env.suite.CreateFundedAccount()
		require.NoError(t, err)
		env.ch.RegisterUserSigningKey(addr, secp)
		_, err = env.ch.GenerateUserKey(addr)
		require.NoError(t, err)
		return addr
	}
	associate := func(addr ethCommon.Address, appID common.ApplicationIdType) {
		userP521, err := env.ch.GetUserKey(addr)
		require.NoError(t, err)
		assoc, err := env.ch.CreateAssociateKeyRequest(appID, commontestutil.GenerateRandomRequestID(), addr, userP521.PublicKey(), execPub)
		require.NoError(t, err)
		require.NoError(t, env.suite.SubmitRequest(assoc))
		require.NoError(t, env.suite.AssertRequestCompleted(assoc.RequestID, setupTimeout))
	}
	normalUser := newAccount()
	fireUser := newAccount()
	feeCollector := newAccount()

	// Deploy 5 normal apps.
	simpleWasm := buildWasmApp(t, "simple")
	var normalApps []common.ApplicationIdType
	for i := 0; i < appsPerKind; i++ {
		appID := env.deployApp(t, simpleWasm, nil, nil, setupTimeout)
		associate(normalUser, appID)
		normalApps = append(normalApps, appID)
	}

	// Deploy 5 trigger apps, each with its own GuardedTrigger, and fund the
	// firing user on each.
	triggerWasm := buildWasmApp(t, "trigger")
	sink := ethCommon.HexToAddress("0x00000000000000000000000000000000000000F1")
	var triggerApps []common.ApplicationIdType
	for i := 0; i < appsPerKind; i++ {
		trigger := env.suite.DeployGuardedTrigger(sink)
		params, err := json.Marshal(map[string]any{
			"triggerAddress": trigger.Hex(),
			"feeCollector":   feeCollector.Hex(),
		})
		require.NoError(t, err)
		appID := env.deployApp(t, triggerWasm, params, &trigger, setupTimeout)
		associate(fireUser, appID)
		associate(feeCollector, appID)

		deposit, err := env.ch.CreateDepositRequest(appID, commontestutil.GenerateRandomRequestID(), fireUser, big.NewInt(1_000_000_000_000), execPub)
		require.NoError(t, err)
		require.NoError(t, env.suite.SubmitRequest(deposit))
		require.NoError(t, env.suite.AssertRequestCompleted(deposit.RequestID, setupTimeout))
		triggerApps = append(triggerApps, appID)
	}

	// Measurement window starts after setup traffic has fully completed.
	fromBlock := env.currentBlock(t) + 1
	logOffset := env.logOffset(t)

	// Pre-fill: round-robin across all 10 apps, alternating normal/trigger so
	// the queue interleaves app kinds maximally.
	start := time.Now()
	perApp := n / (2 * appsPerKind)
	txs := make([]*ethTypes.Transaction, 0, n)
	for round := 0; round < perApp; round++ {
		for i := 0; i < appsPerKind; i++ {
			dep, err := env.ch.CreateDepositRequest(normalApps[i], commontestutil.GenerateRandomRequestID(), normalUser, big.NewInt(1000), execPub)
			require.NoError(t, err)
			tx, err := env.suite.SubmitRequestNoWait(dep)
			require.NoError(t, err)
			txs = append(txs, tx)

			fire, err := env.ch.CreateProcessRequest(triggerApps[i], commontestutil.GenerateRandomRequestID(), fireUser, []byte(`{"type":"fire"}`), execPub)
			require.NoError(t, err)
			tx, err = env.suite.SubmitRequestNoWait(fire)
			require.NoError(t, err)
			txs = append(txs, tx)
		}
	}
	env.confirmSubmissions(t, txs)

	nNormal := n / 2
	nTrigger := n / 2
	normalAppSet := make(map[uint64]bool, len(normalApps))
	for _, id := range normalApps {
		normalAppSet[uint64(id)] = true
	}

	// Drain: trigger round-trips need two state updates each.
	drainTimeout := time.Duration(n)*20*blockTime + 3*time.Minute
	deadline := time.Now().Add(drainTimeout)
	var normalRecords []CompletionRecord
	var trips []triggerRoundTrip
	for {
		normalRecords = nil
		for _, r := range env.scanCompletions(t, fromBlock) {
			if normalAppSet[r.AppID] {
				normalRecords = append(normalRecords, r)
			}
		}
		trips = env.scanTriggerRoundTrips(t, fromBlock)
		if len(normalRecords) >= nNormal && len(trips) >= nTrigger {
			break
		}
		require.False(t, time.Now().After(deadline),
			"timeout: %d/%d normal, %d/%d trigger round-trips after %s",
			len(normalRecords), nNormal, len(trips), nTrigger, drainTimeout)
		time.Sleep(250 * time.Millisecond)
	}
	wallClock := time.Since(start)

	env.fillGasUsed(t, normalRecords)
	records := append(normalRecords, env.roundTripRecords(t, trips)...)

	metrics, err := ComputeSaturationMetrics(records)
	require.NoError(t, err)

	// Report before validity assertions so a failed run still leaves its
	// diagnostics (the report flags a non-zero failed count loudly).
	writeReport(t, RunParams{
		Scenario:        "S3-mixed-10-apps",
		Implementation:  "baseline",
		Requests:        n,
		BlockTime:       blockTime,
		PollingInterval: time.Second,
	}, metrics, env.stageTimings(t, logOffset), wallClock,
		"> Trigger-app requests comprise two on-chain state updates (fire + TRUSTPROCESS); gas/request covers both. The state-update-tx and batch-size rows count the TRUSTPROCESS leg only — the mine-wait call count reflects the full transaction count.\n",
		perAppSection(records))

	require.Equal(t, 0, metrics.Failed, "failed requests invalidate the comparison (see report)")
	env.suite.AssertNoStateUpdateErrors(t)
}

// perAppSection renders the per-app breakdown (completions and latency) as a
// markdown section appended to the aggregate report.
func perAppSection(records []CompletionRecord) string {
	byApp := make(map[uint64][]CompletionRecord)
	for _, r := range records {
		byApp[r.AppID] = append(byApp[r.AppID], r)
	}
	appIDs := make([]uint64, 0, len(byApp))
	for id := range byApp {
		appIDs = append(appIDs, id)
	}
	sort.Slice(appIDs, func(i, j int) bool { return appIDs[i] < appIDs[j] })

	var b strings.Builder
	b.WriteString("## Per-app breakdown\n\n| App | Completed | Failed | Latency p50 (blocks) | Latency p95 (blocks) |\n|---|---|---|---|---|\n")
	for _, id := range appIDs {
		m, err := ComputeSaturationMetrics(byApp[id])
		if err != nil {
			continue
		}
		fmt.Fprintf(&b, "| %d | %d | %d | %.0f | %.0f |\n", id, m.Completed, m.Failed, m.LatencyBlocksP50, m.LatencyBlocksP95)
	}
	return b.String()
}
