package system

import (
	"bytes"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/authorityservice/deployartifact"
	"github.com/HorizenOfficial/vela/pkg/common"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/testutil"
	"github.com/HorizenOfficial/vela/pkg/testutil/fullstack"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

// These tests exercise the external smart-contract trigger feature end-to-end
// against a simulated chain with the real ProcessorEndpoint deployed, a real
// manager and executor, and the real WASM runtime (app/trigger). They assert
// that the components work together — not the individual units.

// uploadTriggerDescriptor saves the WASM into the suite's shared artifacts dir
// (where the manager reads deploy artifacts from) and returns the deploy
// descriptor payload referencing it.
func uploadTriggerDescriptor(t *testing.T, artifactsPath string, wasm []byte, params json.RawMessage) []byte {
	t.Helper()
	store, err := deployartifact.NewStore(artifactsPath)
	require.NoError(t, err)

	resp, err := store.SaveWASM(bytes.NewReader(wasm))
	require.NoError(t, err)

	descriptor := common.DeployDescriptor{
		Mode:              common.DeployModeArtifactRef,
		ArtifactID:        resp.ArtifactID,
		WasmSHA256:        resp.WasmSHA256,
		ConstructorParams: params,
	}
	payload, err := json.Marshal(descriptor)
	require.NoError(t, err)
	return payload
}

// TestTrigger_DeployAppWithTrigger is the happy-path deploy test: it deploys a
// TestTrigger, then deploys the trigger app registering that trigger, and
// verifies the deploy committed both private (manager DB) and on-chain (state
// root, trigger registration mappings) state.
func TestTrigger_DeployAppWithTrigger(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("skipping trigger fullstack test under CI_FLAG")
	}

	suite := fullstack.NewWasmRuntimeSuite(t)
	defer func() { _ = suite.Cleanup() }()

	wasm := buildWasmApp(t, "trigger")

	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	timeout := 100 * time.Second

	// Deploy the trigger contract; it must not be registered to any app yet.
	trigger := suite.DeployTestTrigger(false, false)
	require.NotEqual(t, ethCommon.Address{}, trigger)
	require.Equal(t, uint64(0), suite.GetTriggerAppId(trigger), "trigger must be unregistered before deploy")

	// Submit a deploy request that registers the trigger for the new app. The
	// contract assigns RequestID and ApplicationID, written back into deployReq.
	deployReq := &common.Request{
		RequestType:  common.Deploy,
		Payload:      uploadTriggerDescriptor(t, suite.GetArtifactsPath(), wasm, nil),
		Sender:       suite.GetDeployerAddress(),
		Timestamp:    common.ToBig(big.NewInt(time.Now().Unix())),
		TokenAddress: velacommon.ETH_TOKEN,
		AssetAmount:  common.NewBig(0),
		MaxFeeValue:  common.NewBig(100),
	}
	require.NoError(t, suite.SubmitDeployRequestWithTrigger(deployReq, trigger))

	appID := deployReq.ApplicationID

	// Private side: app state persisted in the manager's data layer.
	_, err := suite.WaitForAppStateInDB(appID, timeout)
	require.NoError(t, err)

	// On-chain side: application state root committed (non-zero). This succeeding
	// means the contract accepted the deploy's stateUpdate.
	_, err = suite.WaitForAppStateInBlockchain(appID, timeout)
	require.NoError(t, err)

	// Deploy request marked completed (observed via SubmitStateUpdate interception).
	require.NoError(t, suite.AssertRequestCompleted(deployReq.RequestID, timeout))

	// Trigger registration committed on-chain, both directions of the mapping.
	require.Equal(t, trigger, suite.GetTriggerContract(appID), "triggerContracts[appId] must point at the trigger")
	require.Equal(t, uint64(appID), suite.GetTriggerAppId(trigger), "triggersToAppIds[trigger] must point at the app")

	// No silent stateUpdate reverts during the happy path.
	suite.AssertNoStateUpdateErrors(t)
}

// deployTriggerApp deploys app/trigger registering `trigger` with the given
// constructor params, waits for the deploy to commit (DB + on-chain), and returns
// the contract-assigned appID.
func deployTriggerApp(t *testing.T, suite *fullstack.FullStackSystemTestSuite, wasm []byte, trigger ethCommon.Address, params json.RawMessage, timeout time.Duration) common.ApplicationIdType {
	t.Helper()
	deployReq := &common.Request{
		RequestType:  common.Deploy,
		Payload:      uploadTriggerDescriptor(t, suite.GetArtifactsPath(), wasm, params),
		Sender:       suite.GetDeployerAddress(),
		Timestamp:    common.ToBig(big.NewInt(time.Now().Unix())),
		TokenAddress: velacommon.ETH_TOKEN,
		AssetAmount:  common.NewBig(0),
		MaxFeeValue:  common.NewBig(100),
	}
	require.NoError(t, suite.SubmitDeployRequestWithTrigger(deployReq, trigger))

	appID := deployReq.ApplicationID
	_, err := suite.WaitForAppStateInDB(appID, timeout)
	require.NoError(t, err)
	_, err = suite.WaitForAppStateInBlockchain(appID, timeout)
	require.NoError(t, err)
	require.NoError(t, suite.AssertRequestCompleted(deployReq.RequestID, timeout))
	return appID
}

// createTriggerAccount creates a funded on-chain account and wires its keys into
// the crypto helper, so its address is known before deploy. Call
// associateTriggerKey after the app is deployed to register the P521 key on the
// executor (required to submit requests or to receive encrypted events).
func createTriggerAccount(t *testing.T, suite *fullstack.FullStackSystemTestSuite, ch *testutil.CryptoHelper) ethCommon.Address {
	t.Helper()
	addr, secp, err := suite.CreateFundedAccount()
	require.NoError(t, err)
	// Same secp256k1 key on-chain (for tx) and in the crypto helper (for seed),
	// so the seed signature recovers to addr.
	ch.RegisterUserSigningKey(addr, secp)
	_, err = ch.GenerateUserKey(addr)
	require.NoError(t, err)
	return addr
}

// associateTriggerKey registers an account's P521 key on the executor for the
// given app via AssociateKey, so it can submit encrypted Process requests and
// receive encrypted events.
func associateTriggerKey(t *testing.T, suite *fullstack.FullStackSystemTestSuite, ch *testutil.CryptoHelper, addr ethCommon.Address, appID common.ApplicationIdType, timeout time.Duration) {
	t.Helper()
	userP521, err := ch.GetUserKey(addr)
	require.NoError(t, err)
	execPub, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	assoc, err := ch.CreateAssociateKeyRequest(appID, commontestutil.GenerateRandomRequestID(), addr, userP521.PublicKey(), execPub)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(assoc)) // fills assoc.RequestID from chain
	require.NoError(t, suite.AssertRequestCompleted(assoc.RequestID, timeout))
}

// TestTrigger_RequestTriggersTrustProcessRoundTrip is the core round-trip with an
// ETH unshield/re-shield: a PROCESS "fire" withdraws 1% of the user's ETH balance
// to the registered (guarded) trigger and emits an AppEvent. On-chain the trigger
// receives the ETH (unshield), sweeps all but 1 wei to its sink, and the leftover
// 1 wei is re-shielded into app custody; the trigger then enqueues a TRUSTPROCESS
// carrying the re-shielded amount, which the WASM trusted_request credits to a
// fee-collector account (emitting an encrypted balance event). The test verifies
// the full money flow on-chain plus the decrypted fee-collector event, and that
// the loop terminates after exactly one TRUSTPROCESS.
func TestTrigger_RequestTriggersTrustProcessRoundTrip(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("skipping trigger fullstack test under CI_FLAG")
	}

	suite := fullstack.NewWasmRuntimeSuite(t)
	defer func() { _ = suite.Cleanup() }()

	wasm := buildWasmApp(t, "trigger")
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	timeout := 100 * time.Second
	ch := testutil.NewCryptoHelper()

	// Fixed sink EOA: the trigger sweeps all received ETH but 1 wei here.
	sink := ethCommon.HexToAddress("0x00000000000000000000000000000000000000F1")

	// Create accounts up front so their addresses can go into the deploy params:
	// the firing user, and the internal fee-collector that receives the re-shield.
	user := createTriggerAccount(t, suite, ch)
	feeCollector := createTriggerAccount(t, suite, ch)

	// Deploy the loop-safe guarded trigger (sweeps to sink) and the app, passing
	// the trigger + fee-collector addresses as constructor params.
	trigger := suite.DeployGuardedTrigger(sink)
	require.NotEqual(t, ethCommon.Address{}, trigger)
	params, err := json.Marshal(map[string]any{
		"triggerAddress": trigger.Hex(),
		"feeCollector":   feeCollector.Hex(),
	})
	require.NoError(t, err)
	appID := deployTriggerApp(t, suite, wasm, trigger, params, timeout)

	// Register both accounts' keys now that the app exists: the user to submit
	// requests, the fee collector to receive its encrypted balance event.
	associateTriggerKey(t, suite, ch, user, appID, timeout)
	associateTriggerKey(t, suite, ch, feeCollector, appID, timeout)

	// Deposit 10000 wei so 1% = 100 wei is a clean withdrawal amount.
	const depositAmount = 10000
	const withdrawAmount = depositAmount / 100 // counter=1 => 1%
	const reshielded = 1                       // trigger keeps 1 wei, swept back
	const toSink = withdrawAmount - reshielded
	execPub, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	depositReq, err := ch.CreateDepositRequest(appID, commontestutil.GenerateRandomRequestID(), user, big.NewInt(depositAmount), execPub)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(depositReq))
	require.NoError(t, suite.AssertRequestCompleted(depositReq.RequestID, timeout))

	require.EqualValues(t, depositAmount, suite.GetAppCustody(appID, velacommon.ETH_TOKEN).Int64(), "custody must equal the deposit")
	require.EqualValues(t, 0, suite.GetEthBalance(sink).Int64(), "sink must start empty")

	// Submit the PROCESS "fire": withdraws 1% to the trigger and emits an AppEvent.
	fireReq, err := ch.CreateProcessRequest(appID, commontestutil.GenerateRandomRequestID(), user, []byte(`{"type":"fire"}`), execPub)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(fireReq)) // fills fireReq.RequestID from chain
	require.NoError(t, suite.AssertRequestCompleted(fireReq.RequestID, timeout))

	firePayload, err := suite.GetRequestUpdatePayload(fireReq.RequestID)
	require.NoError(t, err)
	require.NotEmpty(t, firePayload.AppEvents, "PROCESS fire must emit an AppEvent to fire the trigger")
	require.NotEqual(t, firePayload.PrevStateRoot, firePayload.NewStateRoot, "PROCESS must change the state root")
	require.Len(t, firePayload.Withdrawals, 1, "fire must produce one withdrawal (the unshield to the trigger)")
	require.Equal(t, trigger, firePayload.Withdrawals[0].DestinationAddress, "withdrawal must target the trigger")
	require.EqualValues(t, withdrawAmount, firePayload.Withdrawals[0].Amount.ToInt().Int64(), "withdrawal must be 1% of balance")

	// The trigger enqueued a TRUSTPROCESS; the manager processed it via trusted_request.
	tpID, tpPayload, err := suite.WaitForTrustProcessRequest(timeout)
	require.NoError(t, err)
	require.NotNil(t, tpPayload)

	// State-root chain continuity: the TRUSTPROCESS builds on the PROCESS result,
	// and itself changes state (trusted_request credited the fee collector).
	require.Equal(t, firePayload.NewStateRoot, tpPayload.PrevStateRoot, "TRUSTPROCESS must build on the PROCESS state root")
	require.NotEqual(t, tpPayload.PrevStateRoot, tpPayload.NewStateRoot, "TRUSTPROCESS must change the state root")

	// trusted_request emits no AppEvents — this is what terminates the loop.
	require.Empty(t, tpPayload.AppEvents, "trusted_request must emit no AppEvents (loop terminator)")

	// On-chain state root reflects the final (TRUSTPROCESS) update.
	require.Equal(t, tpPayload.NewStateRoot, suite.GetStateRoot(appID), "on-chain root must match the TRUSTPROCESS result")

	// The guard worked: exactly one TRUSTPROCESS, no infinite re-enqueue.
	require.Equal(t, 1, suite.TrustProcessCount(), "exactly one TRUSTPROCESS expected (guarded trigger stops the loop)")
	require.NotEqual(t, common.RequestIdType{}, tpID)

	// Money flow on-chain: custody dropped by the unshield then regained the
	// re-shielded wei; the sink received the swept remainder.
	require.EqualValues(t, depositAmount-withdrawAmount+reshielded, suite.GetAppCustody(appID, velacommon.ETH_TOKEN).Int64(), "custody = deposit - unshield + reshield")
	require.EqualValues(t, toSink, suite.GetEthBalance(sink).Int64(), "sink received the swept remainder")

	// The re-shielded amount was credited to the fee collector: decrypt the event
	// it emitted to that account and check the reported balance.
	feeEvent := findUserEvent(tpPayload.Events, feeCollector)
	require.NotNil(t, feeEvent, "trusted_request must emit a fee-collector event")
	decrypted, err := ch.DecryptEvent(feeCollector, feeEvent, execPub)
	require.NoError(t, err)
	var body struct {
		Token   string `json:"token"`
		Amount  string `json:"amount"`
		Balance string `json:"balance"`
	}
	require.NoError(t, json.Unmarshal(decrypted, &body))
	require.Equal(t, ethCommon.Address{}, ethCommon.HexToAddress(body.Token), "fee credited in ETH")
	require.EqualValues(t, reshielded, hexutil.MustDecodeBig(body.Amount).Int64(), "fee event amount = re-shielded wei")
	require.EqualValues(t, reshielded, hexutil.MustDecodeBig(body.Balance).Int64(), "fee collector balance = re-shielded wei")

	suite.AssertNoStateUpdateErrors(t)
}

// findUserEvent returns the first event addressed to addr, or nil.
func findUserEvent(events []common.Event, addr ethCommon.Address) *common.Event {
	for i := range events {
		if events[i].UserID == addr {
			return &events[i]
		}
	}
	return nil
}

// TestTrigger_FailedTrustProcessDoesNotLoop is the negative round-trip: the trigger
// enqueues a TRUSTPROCESS whose payload the real WASM trusted_request rejects, so
// the executor marks it FAILED. It verifies the failure is inert end-to-end — the
// manager does not loop (exactly one TRUSTPROCESS ever pulled, on-chain queue
// drained) and the failed TRUSTPROCESS leaves the app state root unchanged. This
// exercises executor+manager+contract together; the contract tests (trigger.spec.ts)
// only hand-craft the failing stateUpdate.
func TestTrigger_FailedTrustProcessDoesNotLoop(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("skipping trigger fullstack test under CI_FLAG")
	}

	suite := fullstack.NewWasmRuntimeSuite(t)
	defer func() { _ = suite.Cleanup() }()

	wasm := buildWasmApp(t, "trigger")
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	timeout := 100 * time.Second
	ch := testutil.NewCryptoHelper()

	// TestTrigger returns its configured payload UNCONDITIONALLY. That would loop
	// with a succeeding trusted_request, but here the TRUSTPROCESS FAILS, and the
	// contract's error path skips _invokeTrigger — so there is no re-enqueue.
	trigger := suite.DeployTestTrigger(false, false)
	require.NotEqual(t, ethCommon.Address{}, trigger)

	user := createTriggerAccount(t, suite, ch)
	// feeCollector is unused on the failure path (trusted_request rejects before
	// crediting) but a valid address is required in the deploy params.
	params, err := json.Marshal(map[string]any{
		"triggerAddress": trigger.Hex(),
		"feeCollector":   user.Hex(),
	})
	require.NoError(t, err)
	appID := deployTriggerApp(t, suite, wasm, trigger, params, timeout)
	associateTriggerKey(t, suite, ch, user, appID, timeout)

	// Arm the trigger with a malformed (non-32-byte) trusted payload: the WASM
	// trusted_request rejects it, so the executor marks the TRUSTPROCESS FAILED.
	suite.SetTrustedPayload(trigger, []byte{0xab, 0xcd})

	// "ping" changes state and fires the trigger, enqueuing one TRUSTPROCESS.
	execPub, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	pingReq, err := ch.CreateProcessRequest(appID, commontestutil.GenerateRandomRequestID(), user, []byte(`{"type":"ping"}`), execPub)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(pingReq))
	require.NoError(t, suite.AssertRequestCompleted(pingReq.RequestID, timeout))

	pingPayload, err := suite.GetRequestUpdatePayload(pingReq.RequestID)
	require.NoError(t, err)
	pingRoot := pingPayload.NewStateRoot

	// The trigger enqueued a TRUSTPROCESS; the manager pulled it and the executor
	// marked it FAILED (trusted_request rejected the malformed payload).
	tpID, err := suite.WaitForFailedTrustProcessRequest(timeout)
	require.NoError(t, err)
	require.NotEqual(t, common.RequestIdType{}, tpID)

	// No loop: give the manager time to (not) re-enqueue, then assert exactly one
	// TRUSTPROCESS was ever pulled and the on-chain trigger queue is drained.
	time.Sleep(2 * time.Second)
	require.Equal(t, 1, suite.TrustProcessCount(), "a failed TRUSTPROCESS must not be re-enqueued (no loop)")
	require.EqualValues(t, 0, suite.GetTriggerQueueSize().Int64(), "trigger queue must drain after the failed TRUSTPROCESS")

	// The failed TRUSTPROCESS committed no state change: the on-chain root stays at
	// the ping's root (the error path leaves state unchanged).
	require.Equal(t, pingRoot, suite.GetStateRoot(appID), "failed TRUSTPROCESS must not change the state root")

	// The failure was a clean signed error stateUpdate, not an on-chain revert.
	suite.AssertNoStateUpdateErrors(t)
}

// TestTrigger_FailedProcessEnqueuesNoTrustProcess is the other failure path: the
// FIRST request ("fire") itself fails in the executor (no funds to withdraw), so
// the contract's error path skips _invokeTrigger entirely — no TRUSTPROCESS is
// ever enqueued, even though the trigger is registered and armed. End-to-end
// counterpart of the contract test "a failed PROCESS does not invoke the trigger".
func TestTrigger_FailedProcessEnqueuesNoTrustProcess(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("skipping trigger fullstack test under CI_FLAG")
	}

	suite := fullstack.NewWasmRuntimeSuite(t)
	defer func() { _ = suite.Cleanup() }()

	wasm := buildWasmApp(t, "trigger")
	require.NoError(t, suite.StartExecutor())
	require.NoError(t, suite.StartManager())

	timeout := 100 * time.Second
	ch := testutil.NewCryptoHelper()

	trigger := suite.DeployTestTrigger(false, false)
	require.NotEqual(t, ethCommon.Address{}, trigger)

	user := createTriggerAccount(t, suite, ch)
	params, err := json.Marshal(map[string]any{
		"triggerAddress": trigger.Hex(),
		"feeCollector":   user.Hex(),
	})
	require.NoError(t, err)
	appID := deployTriggerApp(t, suite, wasm, trigger, params, timeout)
	associateTriggerKey(t, suite, ch, user, appID, timeout)

	// Arm the trigger: were a failed request to (wrongly) invoke it, this
	// unconditional payload would enqueue a TRUSTPROCESS.
	suite.SetTrustedPayload(trigger, []byte{0xab, 0xcd})

	// Baseline root after setup; the failed request must not change it.
	rootBefore := suite.GetStateRoot(appID)

	// "fire" with no prior deposit: the WASM has no balance to withdraw and returns
	// an error → the executor marks the PROCESS FAILED before any trigger invocation.
	execPub, err := suite.GetExecutorCommunicationKey()
	require.NoError(t, err)
	fireReq, err := ch.CreateProcessRequest(appID, commontestutil.GenerateRandomRequestID(), user, []byte(`{"type":"fire"}`), execPub)
	require.NoError(t, err)
	require.NoError(t, suite.SubmitRequest(fireReq))
	require.NoError(t, suite.AssertRequestFailed(fireReq.RequestID, timeout))

	// The failed PROCESS invoked no trigger: no TRUSTPROCESS enqueued, ever.
	time.Sleep(2 * time.Second)
	require.Equal(t, 0, suite.TrustProcessCount(), "a failed PROCESS must not enqueue a TRUSTPROCESS")
	require.EqualValues(t, 0, suite.GetTriggerQueueSize().Int64(), "trigger queue must stay empty after a failed PROCESS")

	// State unchanged by the failed request.
	require.Equal(t, rootBefore, suite.GetStateRoot(appID), "failed PROCESS must not change the state root")

	// Clean signed error, not an on-chain revert.
	suite.AssertNoStateUpdateErrors(t)
}
