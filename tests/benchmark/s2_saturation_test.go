package benchmark

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/processorendpoint"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
)

// triggerRoundTrip pairs a user request with the TRUSTPROCESS it spawned.
// Per THROUGHPUT_BENCHMARKS.md §4.3, a trigger-app request is complete only
// when its TRUSTPROCESS is finalized.
type triggerRoundTrip struct {
	fireID           [32]byte
	trustedID        [32]byte
	appID            uint64 // application the round-trip belongs to
	submitBlock      uint64 // block of the user's submitRequest tx
	fireCompTx       string // state-update tx finalizing the user request
	trustedCompTx    string // state-update tx finalizing the TRUSTPROCESS
	trustedCompBlock uint64
	failed           bool // either leg finalized with an error
}

// scanTriggerRoundTrips reads the window's logs and correlates each
// user-submitted request with its spawned TRUSTPROCESS. The correlation hook
// is on-chain: _enqueueTrustedRequest emits RequestSubmitted *inside the
// spawning request's own state-update transaction* (ProcessorEndpoint.sol),
// so a tx containing both RequestCompleted(F) and RequestSubmitted(T) means F
// spawned T. This holds for any number of in-flight round-trips. The on-chain
// hook also survives batch execution (T1), but with several completions per tx
// the same-tx grouping alone cannot pair spawner and spawned — a guard below
// fails loudly so the pairing gets extended when batch lands.
func (e *benchEnv) scanTriggerRoundTrips(t *testing.T, fromBlock uint64) []triggerRoundTrip {
	t.Helper()
	ctx := context.Background()
	head, err := e.client().BlockNumber(ctx)
	require.NoError(t, err)
	if head < fromBlock {
		return nil
	}
	logs, err := e.client().FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(head),
		Addresses: []ethCommon.Address{e.processor},
	})
	require.NoError(t, err)

	pe := processorendpoint.NewProcessorEndpoint()
	inst := pe.Instance(e.suite.GetSimTestHelper().Client(), e.processor)

	type completion struct {
		block  uint64
		tx     string
		appID  uint64
		failed bool
	}
	userSubmitted := make(map[[32]byte]uint64)     // requestId -> submission block (user txs)
	spawned := make(map[[32]byte][32]byte)         // spawning requestId -> trusted requestId
	completions := make(map[[32]byte]completion)   // requestId -> finalization
	submittedInTx := make(map[string][][32]byte)   // txHash -> RequestSubmitted ids
	completedInTx := make(map[string][][32]byte)   // txHash -> RequestCompleted ids

	for _, lg := range logs {
		txHash := lg.TxHash.Hex()
		var sub processorendpoint.ProcessorEndpointRequestSubmitted
		if err := inst.UnpackLog(&sub, processorendpoint.ProcessorEndpointRequestSubmittedEventName, lg); err == nil {
			submittedInTx[txHash] = append(submittedInTx[txHash], sub.RequestId)
			userSubmitted[sub.RequestId] = lg.BlockNumber // reclassified below if the tx is a state update
			continue
		}
		var comp processorendpoint.ProcessorEndpointRequestCompleted
		if err := inst.UnpackLog(&comp, processorendpoint.ProcessorEndpointRequestCompletedEventName, lg); err == nil {
			completedInTx[txHash] = append(completedInTx[txHash], comp.RequestId)
			completions[comp.RequestId] = completion{block: lg.BlockNumber, tx: txHash, appID: comp.ApplicationId, failed: comp.ErrorCode != 0}
		}
	}

	// A RequestSubmitted inside a tx that also finalizes a request is a
	// trusted enqueue, not a user submission: bind it to the spawning request.
	for txHash, completedIDs := range completedInTx {
		for _, trustedID := range submittedInTx[txHash] {
			delete(userSubmitted, trustedID)
			require.Len(t, completedIDs, 1,
				"correlation assumes one finalization per state-update tx today; extend for batch txs")
			spawned[completedIDs[0]] = trustedID
		}
	}

	// Assemble finished round-trips: user request completed AND its trusted
	// follow-up completed.
	var trips []triggerRoundTrip
	for fireID, submitBlock := range userSubmitted {
		fireComp, fireDone := completions[fireID]
		if !fireDone {
			continue
		}
		trustedID, hasTrusted := spawned[fireID]
		if !hasTrusted {
			continue // trusted not enqueued (or fire failed) — not a finished round-trip
		}
		trustedComp, trustedDone := completions[trustedID]
		if !trustedDone {
			continue
		}
		trips = append(trips, triggerRoundTrip{
			fireID:           fireID,
			trustedID:        trustedID,
			appID:            fireComp.appID,
			submitBlock:      submitBlock,
			fireCompTx:       fireComp.tx,
			trustedCompTx:    trustedComp.tx,
			trustedCompBlock: trustedComp.block,
			failed:           fireComp.failed || trustedComp.failed,
		})
	}
	return trips
}

// roundTripRecords converts finished round-trips into CompletionRecords: one
// logical record per user request, completing at the TRUSTPROCESS block, with
// GasUsed covering BOTH state-update transactions.
func (e *benchEnv) roundTripRecords(t *testing.T, trips []triggerRoundTrip) []CompletionRecord {
	t.Helper()
	ctx := context.Background()
	gasByTx := make(map[string]uint64)
	gasOf := func(tx string) uint64 {
		if g, ok := gasByTx[tx]; ok {
			return g
		}
		receipt, err := e.client().TransactionReceipt(ctx, ethCommon.HexToHash(tx))
		require.NoError(t, err)
		gasByTx[tx] = receipt.GasUsed
		return receipt.GasUsed
	}

	records := make([]CompletionRecord, 0, len(trips))
	for _, trip := range trips {
		records = append(records, CompletionRecord{
			RequestID:   trip.fireID,
			AppID:       trip.appID,
			Block:       trip.trustedCompBlock,
			TxHash:      trip.trustedCompTx, // unique per round-trip; carries the combined gas
			GasUsed:     gasOf(trip.fireCompTx) + gasOf(trip.trustedCompTx),
			Failed:      trip.failed,
			SubmitBlock: trip.submitBlock,
		})
	}
	return records
}

// TestBenchS2HotTriggerApp — scenario S2 (THROUGHPUT_BENCHMARKS.md §3): one
// trigger app, 100% fire rate, saturation mode. Every request unshields to
// the GuardedTrigger and completes only when the spawned TRUSTPROCESS is
// finalized, so each logical request costs two state-update transactions.
func TestBenchS2HotTriggerApp(t *testing.T) {
	skipUnlessBench(t)

	n := envInt(t, "VELA_BENCH_N", 30)
	blockTime := time.Duration(envInt(t, "VELA_BENCH_BLOCK_MS", 1000)) * time.Millisecond
	setupTimeout := 120 * time.Second

	env := newBenchEnv(t, blockTime)
	defer func() { _ = env.suite.Cleanup() }()

	// Setup: guarded trigger + trigger app + firing user + fee collector.
	wasm := buildWasmApp(t, "trigger")
	sink := ethCommon.HexToAddress("0x00000000000000000000000000000000000000F1")
	trigger := env.suite.DeployGuardedTrigger(sink)

	// Accounts must exist before deploy: their addresses go into the params.
	userAddr, userSecp, err := env.suite.CreateFundedAccount()
	require.NoError(t, err)
	env.ch.RegisterUserSigningKey(userAddr, userSecp)
	_, err = env.ch.GenerateUserKey(userAddr)
	require.NoError(t, err)
	feeAddr, feeSecp, err := env.suite.CreateFundedAccount()
	require.NoError(t, err)
	env.ch.RegisterUserSigningKey(feeAddr, feeSecp)
	_, err = env.ch.GenerateUserKey(feeAddr)
	require.NoError(t, err)

	params, err := json.Marshal(map[string]any{
		"triggerAddress": trigger.Hex(),
		"feeCollector":   feeAddr.Hex(),
	})
	require.NoError(t, err)
	appID := env.deployApp(t, wasm, params, &trigger, setupTimeout)

	execPub, err := env.suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	for _, addr := range []ethCommon.Address{userAddr, feeAddr} {
		userP521, err := env.ch.GetUserKey(addr)
		require.NoError(t, err)
		assoc, err := env.ch.CreateAssociateKeyRequest(appID, commontestutil.GenerateRandomRequestID(), addr, userP521.PublicKey(), execPub)
		require.NoError(t, err)
		require.NoError(t, env.suite.SubmitRequest(assoc))
		require.NoError(t, env.suite.AssertRequestCompleted(assoc.RequestID, setupTimeout))
	}

	// Deposit enough that N "fire"s (each withdrawing Counter% of the running
	// balance) always withdraw >= 1 wei.
	depositReq, err := env.ch.CreateDepositRequest(appID, commontestutil.GenerateRandomRequestID(), userAddr, big.NewInt(1_000_000_000_000), execPub)
	require.NoError(t, err)
	require.NoError(t, env.suite.SubmitRequest(depositReq))
	require.NoError(t, env.suite.AssertRequestCompleted(depositReq.RequestID, setupTimeout))

	// Measurement window starts after setup traffic has fully completed.
	fromBlock := env.currentBlock(t) + 1
	logOffset := env.logOffset(t)

	// Pre-fill: N fires submitted without waiting for mining.
	start := time.Now()
	txs := make([]*ethTypes.Transaction, 0, n)
	for i := 0; i < n; i++ {
		req, err := env.ch.CreateProcessRequest(appID, commontestutil.GenerateRandomRequestID(), userAddr, []byte(`{"type":"fire"}`), execPub)
		require.NoError(t, err)
		tx, err := env.suite.SubmitRequestNoWait(req)
		require.NoError(t, err)
		txs = append(txs, tx)
	}
	env.confirmSubmissions(t, txs)

	// Drain: each round-trip needs two state updates; generous bound.
	drainTimeout := time.Duration(n)*20*blockTime + 2*time.Minute
	deadline := time.Now().Add(drainTimeout)
	var trips []triggerRoundTrip
	for {
		trips = env.scanTriggerRoundTrips(t, fromBlock)
		if len(trips) >= n {
			break
		}
		require.False(t, time.Now().After(deadline),
			"timeout waiting for round-trips: %d/%d after %s", len(trips), n, drainTimeout)
		time.Sleep(250 * time.Millisecond)
	}
	wallClock := time.Since(start)

	records := env.roundTripRecords(t, trips)
	metrics, err := ComputeSaturationMetrics(records)
	require.NoError(t, err)

	// Report before validity assertions so a failed run still leaves its
	// diagnostics (the report flags a non-zero failed count loudly).
	writeReport(t, RunParams{
		Scenario:        "S2-hot-trigger-app",
		Implementation:  "baseline",
		Requests:        n,
		BlockTime:       blockTime,
		PollingInterval: time.Second,
	}, metrics, env.stageTimings(t, logOffset), wallClock,
		"> Each logical request comprises two on-chain state updates (user request + TRUSTPROCESS); gas/request covers both. The state-update-tx and batch-size rows count the TRUSTPROCESS leg only — the mine-wait call count reflects the full transaction count.\n")

	require.Equal(t, 0, metrics.Failed, "failed round-trips invalidate the comparison (see report)")
	require.Equal(t, n, env.suite.TrustProcessCount(), "every fire must spawn exactly one TRUSTPROCESS")
	env.suite.AssertNoStateUpdateErrors(t)
}
