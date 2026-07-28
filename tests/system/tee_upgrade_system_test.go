package system

import (
	"context"
	"errors"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/common"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/executor"
	"github.com/HorizenOfficial/vela/pkg/manager"
	"github.com/HorizenOfficial/vela/pkg/testutil"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

// Task 9 — TEE-upgrade integration tests (CI-runnable subset).
//
// The full propose→timelock→apply→drain→reconnect swap and KMS recovery under a
// new PCR0 are Nitro-only: a TCP/dev executor reports an empty PCR0 (the dev
// marker), so swap observation is skipped and the drain path cannot be driven
// end-to-end without real NSM hardware. The drain/reconnect state machine itself
// is covered by unit tests in pkg/manager (Task 5). The tests here cover the
// upgrade guarantees that DO run over TCP:
//
//   - R2 key continuity: the keyset is recovered — never regenerated — across a
//     coupled manager+executor restart (TestTeeUpgrade_KeysetRecoveredAcrossRestart).
//   - Task 6 guard: with EXECUTOR_EXPECT_EXISTING_KEYSET set, a "no recovery
//     data" handshake fails fast instead of generating a keyset
//     (TestTeeUpgrade_ExpectExistingKeyset_HandshakeFailsFast).
//   - D2 signer continuity: a regenerated keyset (different signer) is caught by
//     the manager as a fatal error rather than silently orphaning state
//     (TestTeeUpgrade_SignerMismatch_Fatal).

// TestTeeUpgrade_KeysetRecoveredAcrossRestart proves R2: after a coupled
// manager+executor restart, the executor recovers its original keyset from the
// persisted recovery blob rather than generating a new one. The proof is that a
// post-restart request's update payload still verifies against the ORIGINAL
// signing key (depositToSimpleApp validates every payload signature against
// suite.GetExecutorSigningKey()); a regenerated keyset would produce a signature
// that no longer matches. It also checks the pre-restart app state is still
// readable and that dispatch resumes (the state root advances).
func TestTeeUpgrade_KeysetRecoveredAcrossRestart(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("Skipping long running WASM test in CI environment")
	}

	mgrLogCfg := consoleLogConfig()
	excLogCfg := consoleLogConfig()

	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	mgrConfig.BlockchainPollingInterval = 1 // faster dispatch for the test
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)

	// Pre-generate the keyset so the suite exposes its keys and the recovery blob
	// is stored (LevelDB) at startup — the executor restores it on both the
	// initial handshake and the post-restart handshake.
	ctx := context.Background()
	keySet, recoveryData, err := executor.GenerateEnclaveKeySet(ctx, execConfig.KeySetRecoveryType, nil, nil, "")
	require.NoError(t, err)

	suite := testutil.NewSystemTestSuiteWithConfigs(t, "wasm-runtime", mgrConfig, execConfig, keySet, recoveryData, mgrLogCfg, excLogCfg)
	defer suite.Cleanup()

	wasmBytecode := buildWasmApp(t, "simple")
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	timeout := 60 * time.Second
	appID := common.NewApplicationId(1)

	cryptoHelper := testutil.NewCryptoHelper()
	userAddress, err := cryptoHelper.GenerateUserIdentity()
	require.NoError(t, err)

	deploySimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), wasmBytecode)

	executorPubKey, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	registerUserKey(t, suite, cryptoHelper, appID, userAddress, executorPubKey, timeout)

	// Pre-restart deposit — establishes state and validates the signer.
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, big.NewInt(1000000000000000000))

	preState, err := suite.WaitForAppStateInDB(appID, timeout)
	require.NoError(t, err)

	// Coupled restart: both processes die and come back; the executor must
	// recover the keyset from the persisted recovery blob.
	require.NoError(t, suite.RestartAll())

	// State survives the restart and is immediately readable.
	postState, err := suite.WaitForAppStateInDB(appID, timeout)
	require.NoError(t, err)
	require.Equal(t, preState.StateRoot, postState.StateRoot, "app state root must survive the restart unchanged")

	// Dispatch resumes AND the recovered keyset still signs validly: this deposit
	// completes and its payload signature validates against the ORIGINAL signer
	// (asserted inside depositToSimpleApp). A regenerated keyset would fail here.
	depositToSimpleApp(t, suite, cryptoHelper, appID, commontestutil.GenerateRandomRequestID(), userAddress, big.NewInt(1000000000000000000))

	finalState, err := suite.WaitForAppStateInDB(appID, timeout)
	require.NoError(t, err)
	require.NotEqual(t, postState.StateRoot, finalState.StateRoot, "state root must advance after the post-restart request (dispatch resumed)")
}

// submitMockDeploy submits a mock-runtime deploy request and waits for it to be
// marked completed on-chain. Used as a lightweight "dispatch works" probe.
func submitMockDeploy(t *testing.T, suite *testutil.SystemTestSuite, appID common.ApplicationIdType, timeout time.Duration) {
	t.Helper()
	reqID := commontestutil.GenerateRandomRequestID()
	deployReq := &common.Request{
		RequestType:   common.Deploy,
		ApplicationID: appID,
		RequestID:     reqID,
		Payload:       uploadArtifactAndBuildDescriptorPayload(t, suite, []byte("mock-runtime-app-bytecode")),
		Sender:        deployRequestSender,
		Timestamp:     common.ToBig(new(big.Int).SetInt64(time.Now().Unix())),
		TokenAddress:  velacommon.ETH_TOKEN,
		AssetAmount:   common.NewBig(0),
		MaxFeeValue:   common.NewBig(100),
	}
	require.NoError(t, suite.SubmitRequest(deployReq))
	require.NoError(t, suite.AssertRequestCompleted(reqID, timeout))
}

// TestTeeUpgrade_ExecutorRestartReconnects proves the Task 5 production
// crash/relaunch reconnect path: when ONLY the executor restarts (the manager
// stays up), the manager's polling loop detects the dropped channel, re-dials
// the relaunched executor, and re-runs the handshake — during which the
// executor recovers its ORIGINAL keyset from the still-running manager's
// persisted recovery blob. This is the real drain/relaunch path rather than the
// coupled restart that TestTeeUpgrade_KeysetRecoveredAcrossRestart approximates.
//
// The on-chain teeSigner is pinned to the executor's original signer, so a
// successful post-restart dispatch AND the absence of a fatal
// SignerContinuityError together prove the keyset was recovered, not
// regenerated (a regenerated keyset would produce a different signer and trip
// the continuity guard).
func TestTeeUpgrade_ExecutorRestartReconnects(t *testing.T) {
	mgrLogCfg := consoleLogConfig()
	excLogCfg := consoleLogConfig()

	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	mgrConfig.BlockchainPollingInterval = 1 // detect the drop + reconnect quickly
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)

	ctx := context.Background()
	keySet, recoveryData, err := executor.GenerateEnclaveKeySet(ctx, execConfig.KeySetRecoveryType, nil, nil, "")
	require.NoError(t, err)

	suite := testutil.NewSystemTestSuiteWithConfigs(t, "mock-runtime", mgrConfig, execConfig, keySet, recoveryData, mgrLogCfg, excLogCfg)
	defer suite.Cleanup()

	// Pin on-chain teeSigner to the executor's real signer so the
	// signer-continuity check actively verifies recovery on the post-restart
	// reconnect (a regenerated keyset would trip a fatal error instead).
	signerKey, err := suite.GetExecutorSigningKey()
	require.NoError(t, err)
	suite.SetTeeSigner(ethCommon.HexToAddress(signerKey.Address()))

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	timeout := 30 * time.Second

	// Baseline: a deploy dispatches before the restart.
	submitMockDeploy(t, suite, common.NewApplicationId(1), timeout)

	// Relaunch ONLY the executor; the manager keeps running and must reconnect
	// via its polling-loop re-dial.
	require.NoError(t, suite.RestartExecutorOnly())

	// Dispatch resumes: this deploy only completes if the manager reconnected,
	// re-handshook, and the executor recovered its keyset.
	submitMockDeploy(t, suite, common.NewApplicationId(2), timeout)

	// The recovered signer still matches on-chain teeSigner — no fatal fired.
	select {
	case err := <-suite.ManagerFatalErrChan():
		t.Fatalf("unexpected fatal error after executor reconnect: %v", err)
	default:
	}
}

// TestTeeUpgrade_ExpectExistingKeyset_HandshakeFailsFast proves the Task 6
// key-continuity guard end-to-end: with EXECUTOR_EXPECT_EXISTING_KEYSET set and
// no recovery data at the manager (wiped/wrong data folder), the executor aborts
// the handshake instead of generating a fresh keyset. The manager's Start fails,
// and nothing is generated or stored.
func TestTeeUpgrade_ExpectExistingKeyset_HandshakeFailsFast(t *testing.T) {
	mgrLogCfg := consoleLogConfig()
	excLogCfg := consoleLogConfig()

	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	mgrConfig.DataLayerType = "mockdb" // empty store → manager reports found=false
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)
	execConfig.ExpectExistingKeyset = true // the upgrade guard

	// No keyset / recovery data seeded: the manager has nothing to serve.
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "mock-runtime", mgrConfig, execConfig, nil, nil, mgrLogCfg, excLogCfg)
	defer suite.Cleanup()

	require.NoError(t, suite.StartExecutor())

	// The executor sees found=false with the guard set, aborts the handshake, and
	// drops the connection — the manager's handshake wait fails.
	err = suite.StartManager()
	require.Error(t, err)
	require.Contains(t, err.Error(), "executor handshake failed")

	// Nothing was generated or stored: the store is still empty.
	_, getErr := suite.GetDataLayer().GetEnclaveKeySetRecovery(suite.Context())
	require.Error(t, getErr, "no recovery data must have been generated or stored")
}

// TestTeeUpgrade_SignerMismatch_Fatal proves the D2 signer-continuity guard: if
// the executor comes up with a keyset whose signer does not match the on-chain
// teeSigner (i.e. it regenerated instead of recovering), the manager surfaces a
// fatal SignerContinuityError on the polling loop rather than silently accepting
// the orphaning keyset.
func TestTeeUpgrade_SignerMismatch_Fatal(t *testing.T) {
	mgrLogCfg := consoleLogConfig()
	excLogCfg := consoleLogConfig()

	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	mgrConfig.DataLayerType = "mockdb"
	mgrConfig.BlockchainPollingInterval = 1
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)

	// No recovery data + guard off → the executor generates a fresh keyset with a
	// random signer at handshake.
	suite := testutil.NewSystemTestSuiteWithConfigs(t, "mock-runtime", mgrConfig, execConfig, nil, nil, mgrLogCfg, excLogCfg)
	defer suite.Cleanup()

	// On-chain teeSigner is set to a fixed "old" address that cannot match the
	// freshly generated signer, simulating a post-swap enclave that regenerated
	// its keyset instead of recovering the attested one.
	oldSigner := ethCommon.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	suite.SetTeeSigner(oldSigner)

	require.NoError(t, suite.StartExecutor())
	// Handshake succeeds (a keyset is generated and stored); the mismatch is only
	// detectable against the chain, so Start returns and the polling loop catches it.
	require.NoError(t, suite.StartManager())

	select {
	case err := <-suite.ManagerFatalErrChan():
		require.Error(t, err)
		var sce *manager.SignerContinuityError
		require.True(t, errors.As(err, &sce), "expected a SignerContinuityError, got %T: %v", err, err)
	case <-time.After(30 * time.Second):
		t.Fatal("timed out waiting for a fatal signer-continuity error")
	}
}
