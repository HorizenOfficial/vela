package benchmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/authorityservice/deployartifact"
	blockchainTestutil "github.com/HorizenOfficial/vela/pkg/blockchain/testutil"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/processorendpoint"
	"github.com/HorizenOfficial/vela/pkg/common"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/testutil"
	"github.com/HorizenOfficial/vela/pkg/testutil/fullstack"
)

// The benchmarks are gated behind VELA_BENCH: they need Wasmtime + TinyGo and
// run for minutes, so they are excluded from the normal test pass (both quick
// and full). Run with:
//
//	VELA_BENCH=1 go test -v -timeout 30m -run TestBench ./tests/benchmark/
//
// Tunables (env): VELA_BENCH_N (requests, default 30),
// VELA_BENCH_BLOCK_MS (emulated block time, default 1000 — must be >= the
// manager polling interval, which the fullstack suite pins to 1s).
func skipUnlessBench(t *testing.T) {
	t.Helper()
	if os.Getenv("VELA_BENCH") == "" {
		t.Skip("benchmark: set VELA_BENCH=1 to run")
	}
}

// envInt reads an integer tunable; a set-but-malformed value fails the test
// rather than silently running with the default (runs take minutes).
func envInt(t *testing.T, name string, def int) int {
	t.Helper()
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	require.NoError(t, err, "invalid %s: %q", name, v)
	return n
}

// benchEnv bundles the running system plus everything the scenarios need to
// submit load and collect metrics.
type benchEnv struct {
	suite      *fullstack.FullStackSystemTestSuite
	ch         *testutil.CryptoHelper
	mgrLogFile string
	blockTime  time.Duration
	processor  ethCommon.Address
}

// newBenchEnv boots the fullstack stack (manager + executor + simulated chain)
// with the emulated block time and the manager log captured to a file for
// stage-timing extraction.
func newBenchEnv(t *testing.T, blockTime time.Duration) *benchEnv {
	t.Helper()
	t.Setenv("MANAGER_ARTIFACTS_PATH", t.TempDir())
	t.Setenv("EXECUTOR_KEYSET_RECOVERY_TYPE", "0")

	mgrLogFile := filepath.Join(t.TempDir(), "manager.log")
	mgrLogCfg := &logger.Config{
		Kind:         "zerolog",
		Console:      false,
		ConsoleLevel: "info",
		FileName:     mgrLogFile,
		FileLevel:    "info",
	}
	// The temp dir is deleted when the test ends; on failure, dump the tail of
	// the manager log first — it carries the poll-loop errors that never reach
	// the console.
	t.Cleanup(func() {
		if !t.Failed() {
			return
		}
		content, err := os.ReadFile(mgrLogFile)
		if err != nil {
			t.Logf("could not read manager log: %v", err)
			return
		}
		const tail = 8000
		if len(content) > tail {
			content = content[len(content)-tail:]
		}
		t.Logf("manager log tail:\n%s", content)
	})
	excLogCfg := &logger.Config{Kind: "zerolog", Console: true, ConsoleLevel: "warn", ConsoleColor: false}

	suite := fullstack.NewFullStackSystemTestSuite(t, "wasm-runtime", mgrLogCfg, excLogCfg,
		blockchainTestutil.WithAutoMineInterval(blockTime))

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	// The contract caps the pending queue at 10 by default; saturation
	// pre-fill needs room for the full workload (plus in-flight TRUSTPROCESS).
	suite.GetSimTestHelper().SetQueueThreshold(big.NewInt(100_000))

	return &benchEnv{
		suite:      suite,
		ch:         testutil.NewCryptoHelper(),
		mgrLogFile: mgrLogFile,
		blockTime:  blockTime,
		processor:  suite.GetSimTestHelper().ProcessorContractAddress,
	}
}

func (e *benchEnv) client() interface {
	ethereum.LogFilterer
	ethereum.BlockNumberReader
	ethereum.TransactionReader
} {
	return e.suite.GetSimTestHelper().Client()
}

func (e *benchEnv) currentBlock(t *testing.T) uint64 {
	t.Helper()
	n, err := e.client().BlockNumber(context.Background())
	require.NoError(t, err)
	return n
}

// buildWasmApp compiles a TinyGo guest app (same helper as tests/system).
func buildWasmApp(t *testing.T, appName string) []byte {
	t.Helper()
	_, b, _, ok := runtime.Caller(0)
	require.True(t, ok)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	appDir := filepath.Join(projectRoot, "app", appName)

	cmd := exec.Command("make", "build")
	cmd.Dir = appDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "failed to build wasm module: %s", string(output))

	wasmPath := filepath.Join(appDir, "build", appName+"_app.wasm")
	wasmBytecode, err := os.ReadFile(wasmPath)
	require.NoError(t, err)
	require.NotEmpty(t, wasmBytecode)
	return wasmBytecode
}

// descriptorPayload stores the WASM in the shared artifacts dir and returns
// the deploy descriptor payload referencing it.
func (e *benchEnv) descriptorPayload(t *testing.T, wasm []byte, params json.RawMessage) []byte {
	t.Helper()
	store, err := deployartifact.NewStore(e.suite.GetArtifactsPath())
	require.NoError(t, err)
	resp, err := store.SaveWASM(bytes.NewReader(wasm))
	require.NoError(t, err)
	payload, err := json.Marshal(common.DeployDescriptor{
		Mode:       common.DeployModeArtifactRef,
		ArtifactID: resp.ArtifactID,
		WasmSHA256: resp.WasmSHA256,
		ConstructorParams: params,
	})
	require.NoError(t, err)
	return payload
}

// deployApp submits a deploy request (optionally registering a trigger) and
// waits for it to commit. Returns the contract-assigned appID.
func (e *benchEnv) deployApp(t *testing.T, wasm []byte, params json.RawMessage, trigger *ethCommon.Address, timeout time.Duration) common.ApplicationIdType {
	t.Helper()
	req := &common.Request{
		RequestType:  common.Deploy,
		Payload:      e.descriptorPayload(t, wasm, params),
		Sender:       e.suite.GetDeployerAddress(),
		Timestamp:    common.ToBig(big.NewInt(time.Now().Unix())),
		TokenAddress: velacommon.ETH_TOKEN,
		AssetAmount:  common.NewBig(0),
		MaxFeeValue:  common.NewBig(100),
	}
	if trigger != nil {
		require.NoError(t, e.suite.SubmitDeployRequestWithTrigger(req, *trigger))
	} else {
		require.NoError(t, e.suite.SubmitRequest(req))
	}
	_, err := e.suite.WaitForAppStateInDB(req.ApplicationID, timeout)
	require.NoError(t, err)
	_, err = e.suite.WaitForAppStateInBlockchain(req.ApplicationID, timeout)
	require.NoError(t, err)
	require.NoError(t, e.suite.AssertRequestCompleted(req.RequestID, timeout))
	return req.ApplicationID
}

// newAppUser creates a funded account wired into the crypto helper and
// registers its P521 key on the executor for the app.
func (e *benchEnv) newAppUser(t *testing.T, appID common.ApplicationIdType, timeout time.Duration) ethCommon.Address {
	t.Helper()
	addr, secp, err := e.suite.CreateFundedAccount()
	require.NoError(t, err)
	e.ch.RegisterUserSigningKey(addr, secp)
	_, err = e.ch.GenerateUserKey(addr)
	require.NoError(t, err)

	userP521, err := e.ch.GetUserKey(addr)
	require.NoError(t, err)
	execPub, err := e.suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	assoc, err := e.ch.CreateAssociateKeyRequest(appID, commontestutil.GenerateRandomRequestID(), addr, userP521.PublicKey(), execPub)
	require.NoError(t, err)
	require.NoError(t, e.suite.SubmitRequest(assoc))
	require.NoError(t, e.suite.AssertRequestCompleted(assoc.RequestID, timeout))
	return addr
}

// scanCompletions reads RequestSubmitted and RequestCompleted events from
// fromBlock to the chain head and joins them into completion records
// (without gas — see fillGasUsed). Requests submitted before fromBlock (setup
// traffic) produce no RequestSubmitted entry and are excluded.
func (e *benchEnv) scanCompletions(t *testing.T, fromBlock uint64) []CompletionRecord {
	t.Helper()
	ctx := context.Background()
	head, err := e.client().BlockNumber(ctx)
	require.NoError(t, err)
	if head < fromBlock {
		return nil // no blocks in the window yet
	}
	logs, err := e.client().FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(head),
		Addresses: []ethCommon.Address{e.processor},
	})
	require.NoError(t, err)

	pe := processorendpoint.NewProcessorEndpoint()
	inst := pe.Instance(e.suite.GetSimTestHelper().Client(), e.processor)

	submitted := make(map[[32]byte]uint64) // requestId -> submission block
	var completions []CompletionRecord

	for _, lg := range logs {
		var sub processorendpoint.ProcessorEndpointRequestSubmitted
		if err := inst.UnpackLog(&sub, processorendpoint.ProcessorEndpointRequestSubmittedEventName, lg); err == nil {
			submitted[sub.RequestId] = lg.BlockNumber
			continue
		}
		var comp processorendpoint.ProcessorEndpointRequestCompleted
		if err := inst.UnpackLog(&comp, processorendpoint.ProcessorEndpointRequestCompletedEventName, lg); err == nil {
			subBlock, known := submitted[comp.RequestId]
			if !known {
				// Setup-phase request completing late, or an internally
				// enqueued request (TRUSTPROCESS) — scenario code decides how
				// to handle those; skip requests we did not see submitted.
				continue
			}
			completions = append(completions, CompletionRecord{
				RequestID:   comp.RequestId,
				AppID:       comp.ApplicationId,
				Block:       lg.BlockNumber,
				TxHash:      lg.TxHash.Hex(),
				Failed:      comp.ErrorCode != 0,
				SubmitBlock: subBlock,
			})
		}
	}
	return completions
}

// fillGasUsed resolves GasUsed for each record from its transaction receipt
// (fetched once per distinct tx).
func (e *benchEnv) fillGasUsed(t *testing.T, records []CompletionRecord) {
	t.Helper()
	ctx := context.Background()
	gasByTx := make(map[string]uint64)
	for i := range records {
		gas, seen := gasByTx[records[i].TxHash]
		if !seen {
			receipt, err := e.client().TransactionReceipt(ctx, ethCommon.HexToHash(records[i].TxHash))
			require.NoError(t, err)
			gas = receipt.GasUsed
			gasByTx[records[i].TxHash] = gas
		}
		records[i].GasUsed = gas
	}
}

// confirmSubmissions waits for the pre-fill transactions to be mined and
// fails loudly if any reverted (e.g. QueueThresholdExceeded): a silently
// reverted submission would otherwise surface as a drain timeout.
func (e *benchEnv) confirmSubmissions(t *testing.T, txs []*ethTypes.Transaction) {
	t.Helper()
	ctx := context.Background()
	for i, tx := range txs {
		receipt, err := bind.WaitMined(ctx, e.suite.GetSimTestHelper().Client(), tx.Hash())
		require.NoError(t, err, "submission %d not mined", i)
		require.EqualValues(t, 1, receipt.Status, "submission %d reverted on-chain (tx %s)", i, tx.Hash().Hex())
	}
}

// waitForCompletions polls the chain until `expected` benchmark requests are
// finalized (or the timeout expires) and returns the records with gas filled.
func (e *benchEnv) waitForCompletions(t *testing.T, fromBlock uint64, expected int, timeout time.Duration) []CompletionRecord {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		records := e.scanCompletions(t, fromBlock)
		if len(records) >= expected {
			e.fillGasUsed(t, records)
			return records
		}
		require.False(t, time.Now().After(deadline),
			"timeout waiting for completions: %d/%d after %s", len(records), expected, timeout)
		time.Sleep(250 * time.Millisecond)
	}
}

// logOffset returns the current size of the manager log file. Capture it just
// before the measured phase and pass it to stageTimings so setup traffic is
// excluded from the stage timings.
func (e *benchEnv) logOffset(t *testing.T) int64 {
	t.Helper()
	info, err := os.Stat(e.mgrLogFile)
	if os.IsNotExist(err) {
		return 0
	}
	require.NoError(t, err)
	return info.Size()
}

// stageTimings parses the manager log file (from the given byte offset) for
// the two instrumentation lines. The chain scan can observe the last
// completion before the manager's own WaitMined poll returns and writes its
// log line, so settle briefly before reading.
func (e *benchEnv) stageTimings(t *testing.T, fromOffset int64) StageTimings {
	t.Helper()
	time.Sleep(2 * time.Second)
	content, err := os.ReadFile(e.mgrLogFile)
	require.NoError(t, err)
	if fromOffset > int64(len(content)) {
		fromOffset = int64(len(content))
	}
	return ParseStageTimings(string(content[fromOffset:]))
}

// writeReport renders the report, logs it, and saves it under
// tests/benchmark/results/. Optional extra markdown sections are appended
// (e.g. the per-app breakdown of the mixed scenario).
func writeReport(t *testing.T, p RunParams, m SaturationMetrics, tm StageTimings, wallClock time.Duration, extra ...string) {
	t.Helper()
	report := FormatReport(p, m, tm, wallClock)
	for _, section := range extra {
		report += "\n" + section
	}
	t.Logf("\n%s", report)

	_, b, _, ok := runtime.Caller(0)
	require.True(t, ok)
	resultsDir := filepath.Join(filepath.Dir(b), "results")
	require.NoError(t, os.MkdirAll(resultsDir, 0o755))
	name := fmt.Sprintf("%s_%s_%s.md", time.Now().Format("2006-01-02"), p.Scenario, p.Implementation)
	require.NoError(t, os.WriteFile(filepath.Join(resultsDir, name), []byte(report), 0o644))
	t.Logf("report written to %s", filepath.Join(resultsDir, name))
}
