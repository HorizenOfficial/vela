package executor

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/appdata"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	"github.com/HorizenOfficial/vela/pkg/crypto"
	"github.com/stretchr/testify/require"
)

// newBatchTestUser generates a user key pair for batch tests.
func newBatchTestUser(t *testing.T) (ethCommon.Address, *cryptotypes.PrivateKeyP521, *cryptotypes.PublicKeyP521) {
	t.Helper()
	priv, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	addr := ethCommon.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	return addr, priv, priv.PublicKey()
}

// newDepositRequest creates a Process request carrying only a deposit.
func newDepositRequest(sender ethCommon.Address, amount uint64) *common.Request {
	req := newProcessRequest()
	req.Sender = sender
	req.AssetAmount = common.NewBig(amount)
	return req
}

// newTransferRequest creates a Process request whose payload instructs the mock
// runtime to transfer `amount` from sender to a fixed recipient. The payload is
// encrypted toward the enclave with the sender's key.
func newTransferRequest(t *testing.T, exec *StatelessExecutor, sender ethCommon.Address, senderPriv *cryptotypes.PrivateKeyP521, amount uint64) *common.Request {
	t.Helper()
	instructions := testPayloadInstructions{
		Type: "transfer",
		Transfer: &testTransferInstruction{
			To:     ethCommon.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
			Amount: common.NewBig(amount),
		},
	}
	plaintext, err := json.Marshal(instructions)
	require.NoError(t, err)

	req := newProcessRequest()
	req.Sender = sender
	req.Payload = encryptPayload(t, senderPriv, exec.keySet.CommunicationKey.PublicKey(), plaintext)
	return req
}

// decryptFinalState decrypts the final application state and returns the mock
// runtime's internal state.
func decryptFinalState(t *testing.T, exec *StatelessExecutor, appState *common.ApplicationState) testApplicationInternalState {
	t.Helper()
	serialized, err := crypto.DecryptWithAES(exec.keySet.StateKey, appState.EncryptedState)
	require.NoError(t, err)

	ad, err := appdata.DeserializeAppData(serialized)
	require.NoError(t, err)

	var state testApplicationInternalState
	require.NoError(t, json.Unmarshal(ad.GetAppState(), &state))
	return state
}

func TestHandleBatchProcessRequest_EmptyBatch(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))

	payloads, sig, finalState, reports, processed, err := exec.HandleBatchProcessRequest(context.Background(), nil, nil, []byte("wasm"))
	require.Error(t, err)
	require.Nil(t, payloads)
	require.Nil(t, sig)
	require.Nil(t, finalState)
	require.Nil(t, reports)
	require.Zero(t, processed)
}

func TestHandleBatchProcessRequest_NilAppState_HardFailure(t *testing.T) {
	// App existence is validated on-chain (validApplicationId modifier), so a nil
	// state here means tampering or manager-side state loss: hard failure, no
	// signed error payloads, requests stay pending.
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, _ := newBatchTestUser(t)

	requests := []*common.Request{
		newDepositRequest(user, 10),
		newDepositRequest(user, 20),
	}

	payloads, sig, finalState, reports, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, nil, []byte("wasm"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "state not found for application")
	require.Nil(t, payloads)
	require.Nil(t, sig)
	require.Nil(t, finalState)
	require.Nil(t, reports)
	require.Zero(t, processed)
}

func TestHandleBatchProcessRequest_EmptyWasmModule_HardFailure(t *testing.T) {
	// A missing wasm module means the executor cannot load the app: hard failure,
	// no payloads, the whole batch stays pending on-chain.
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, userPub := newBatchTestUser(t)

	appState := buildEncryptedAppState(t, exec, &user, userPub, []byte("wasm"))

	requests := []*common.Request{
		newDepositRequest(user, 10),
		newDepositRequest(user, 20),
	}

	payloads, sig, finalState, reports, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty wasm module for application")
	require.Nil(t, payloads)
	require.Nil(t, sig)
	require.Nil(t, finalState)
	require.Nil(t, reports)
	require.Zero(t, processed)
}

func TestHandleBatchProcessRequest_WrongWasmFingerprint_HardFailure(t *testing.T) {
	// The state commits to the fingerprint of "wasm", but a different module is
	// supplied at execution time: tampering, hard failure, batch stays pending.
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, userPub := newBatchTestUser(t)

	appState := buildEncryptedAppState(t, exec, &user, userPub, []byte("wasm"))

	requests := []*common.Request{
		newDepositRequest(user, 10),
		newDepositRequest(user, 20),
	}

	payloads, sig, finalState, reports, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, []byte("different-wasm"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "wasm fingerprint mismatch for application")
	require.Nil(t, payloads)
	require.Nil(t, sig)
	require.Nil(t, finalState)
	require.Nil(t, reports)
	require.Zero(t, processed)
}

func TestHandleBatchProcessRequest_MultipleSuccess_StateRootsChain(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, userPub := newBatchTestUser(t)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	requests := []*common.Request{
		newDepositRequest(user, 10),
		newDepositRequest(user, 20),
		newDepositRequest(user, 30),
	}

	payloads, sig, finalState, _, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, wasmModule)
	require.NoError(t, err)
	require.Len(t, payloads, len(requests))
	require.Equal(t, len(requests), processed)

	// State roots chain: first from the input state, then payload to payload
	require.Equal(t, appState.StateRoot, payloads[0].PrevStateRoot)
	for i := 1; i < len(payloads); i++ {
		require.Equal(t, payloads[i-1].NewStateRoot, payloads[i].PrevStateRoot)
	}
	for i, p := range payloads {
		require.Equal(t, uint8(0), p.ErrorCode, "payload %d should be a success payload", i)
		require.Empty(t, p.Signature, "payload %d must not be individually signed", i)
		require.NotEqual(t, p.PrevStateRoot, p.NewStateRoot, "payload %d must change the state", i)
	}

	// One batch signature covering all entries
	require.Len(t, sig, 65)

	// Final state is the last state in the chain and decrypts to the accumulated balance
	require.NotNil(t, finalState)
	require.Equal(t, payloads[2].NewStateRoot, finalState.StateRoot)
	state := decryptFinalState(t, exec, finalState)
	require.Equal(t, int64(60), state.Accounts[user].Balance.ToInt().Int64())
}

func TestHandleBatchProcessRequest_SoftFailureMidBatch(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, userPriv, userPub := newBatchTestUser(t)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	requests := []*common.Request{
		newDepositRequest(user, 10),
		newTransferRequest(t, exec, user, userPriv, 1000), // insufficient balance -> soft failure
		newDepositRequest(user, 20),
	}

	payloads, sig, finalState, _, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, wasmModule)
	require.NoError(t, err)
	require.Len(t, payloads, len(requests))
	require.Equal(t, len(requests), processed)

	// Request 2 got an error payload with state unchanged
	require.NotEqual(t, uint8(0), payloads[1].ErrorCode)
	require.Equal(t, payloads[0].NewStateRoot, payloads[1].PrevStateRoot)
	require.Equal(t, payloads[1].PrevStateRoot, payloads[1].NewStateRoot)

	// Request 3 continues from request 1's state
	require.Equal(t, uint8(0), payloads[2].ErrorCode)
	require.Equal(t, payloads[0].NewStateRoot, payloads[2].PrevStateRoot)

	require.Len(t, sig, 65)
	require.Equal(t, payloads[2].NewStateRoot, finalState.StateRoot)
	state := decryptFinalState(t, exec, finalState)
	require.Equal(t, int64(30), state.Accounts[user].Balance.ToInt().Int64())
}

func TestHandleBatchProcessRequest_HardFailureMidBatch_WrongApplicationId(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, userPub := newBatchTestUser(t)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	tampered := newDepositRequest(user, 30)
	tampered.ApplicationID = common.NewApplicationId(999)

	requests := []*common.Request{
		newDepositRequest(user, 10),
		newDepositRequest(user, 20),
		tampered,
		newDepositRequest(user, 40), // never reached
	}

	payloads, sig, finalState, _, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, wasmModule)
	require.NoError(t, err)

	// Batch stopped at the tampered request: only the first two results are returned
	require.Len(t, payloads, 2)
	require.Equal(t, 2, processed)
	require.Len(t, sig, 65)
	require.Equal(t, payloads[0].RequestID, requests[0].RequestID)
	require.Equal(t, payloads[1].RequestID, requests[1].RequestID)

	require.Equal(t, payloads[1].NewStateRoot, finalState.StateRoot)
	state := decryptFinalState(t, exec, finalState)
	require.Equal(t, int64(30), state.Accounts[user].Balance.ToInt().Int64())
}

func TestHandleBatchProcessRequest_HardFailureOnFirstRequest(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, userPub := newBatchTestUser(t)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	bad := newDepositRequest(user, 10)
	bad.ProtocolVersion = 42

	requests := []*common.Request{bad, newDepositRequest(user, 20)}

	payloads, sig, finalState, reports, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, wasmModule)
	require.NoError(t, err)
	require.Empty(t, payloads)
	require.Nil(t, sig)
	require.Nil(t, finalState)
	require.Nil(t, reports)
	require.Zero(t, processed)
}

func TestHandleBatchProcessRequest_DepositDiscardedWhenProcessFails(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, userPriv, userPub := newBatchTestUser(t)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	// Deposit 10 succeeds within the request, but the process step fails:
	// the whole request must fail and the deposit must be discarded.
	depositAndFail := newTransferRequest(t, exec, user, userPriv, 1000)
	depositAndFail.AssetAmount = common.NewBig(10)

	requests := []*common.Request{
		depositAndFail,
		newDepositRequest(user, 5),
	}

	payloads, _, finalState, _, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, wasmModule)
	require.NoError(t, err)
	require.Len(t, payloads, 2)
	require.Equal(t, 2, processed)

	// Request 1 failed with state unchanged from the initial root
	require.NotEqual(t, uint8(0), payloads[0].ErrorCode)
	require.Equal(t, appState.StateRoot, payloads[0].PrevStateRoot)
	require.Equal(t, appState.StateRoot, payloads[0].NewStateRoot)

	// Request 2 chains from the initial root; only its deposit is in the final state
	require.Equal(t, appState.StateRoot, payloads[1].PrevStateRoot)
	state := decryptFinalState(t, exec, finalState)
	require.Equal(t, int64(5), state.Accounts[user].Balance.ToInt().Int64())
}

func TestHandleBatchProcessRequest_SingleRequestParity(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, userPub := newBatchTestUser(t)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	req := newDepositRequest(user, 10)

	singlePayload, singleState, _, err := exec.HandleProcessRequest(context.Background(), req, appState, wasmModule)
	require.NoError(t, err)

	batchPayloads, sig, batchState, _, processed, err := exec.HandleBatchProcessRequest(context.Background(), []*common.Request{req}, appState, wasmModule)
	require.NoError(t, err)
	require.Len(t, batchPayloads, 1)
	require.Equal(t, 1, processed)
	require.Len(t, sig, 65)

	// Same deterministic outcome as the single-request path, minus the per-payload signature
	require.Equal(t, singlePayload.PrevStateRoot, batchPayloads[0].PrevStateRoot)
	require.Equal(t, singlePayload.NewStateRoot, batchPayloads[0].NewStateRoot)
	require.Equal(t, singlePayload.RefundAmount.ToInt(), batchPayloads[0].RefundAmount.ToInt())
	require.Equal(t, singlePayload.ApplicationFee.ToInt(), batchPayloads[0].ApplicationFee.ToInt())
	require.Equal(t, singlePayload.ErrorCode, batchPayloads[0].ErrorCode)
	require.Empty(t, batchPayloads[0].Signature)
	require.Equal(t, singleState.StateRoot, batchState.StateRoot)
}

func TestHandleBatchProcessRequest_BatchSignatureRecovery(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, userPub := newBatchTestUser(t)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	requests := []*common.Request{
		newDepositRequest(user, 10),
		newDepositRequest(user, 20),
	}

	payloads, sig, _, _, _, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, wasmModule)
	require.NoError(t, err)
	require.Len(t, sig, 65)

	// The batch signature must recover the executor's TEE signing address from
	// the hash covering all entry hashes.
	batchHash, err := exec.BuildBatchMsgHash(payloads)
	require.NoError(t, err)

	recoverySig := append([]byte{}, sig...)
	require.GreaterOrEqual(t, recoverySig[64], byte(27))
	recoverySig[64] -= 27
	pubKey, err := ethCrypto.SigToPub(batchHash, recoverySig)
	require.NoError(t, err)

	recovered := ethCrypto.PubkeyToAddress(*pubKey).Hex()
	require.True(t, strings.EqualFold(exec.keySet.SigningKey.PublicKey().Address(), recovered))
}

func TestHandleBatchProcessRequest_SigningFailureDiscardsBatch(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))
	user, _, userPub := newBatchTestUser(t)

	wasmModule := []byte("wasm")
	appState := buildEncryptedAppState(t, exec, &user, userPub, wasmModule)

	requests := []*common.Request{
		newDepositRequest(user, 10),
		newDepositRequest(user, 20),
	}

	// A zero-value builder cannot hash any payload, so batch signing fails
	// after the batch was fully processed.
	exec.MsgToSignBuilder = &MsgToSignBuilder{}

	payloads, sig, finalState, reports, processed, err := exec.HandleBatchProcessRequest(context.Background(), requests, appState, wasmModule)
	require.Error(t, err)
	require.Nil(t, payloads)
	require.Nil(t, sig)
	require.Nil(t, finalState)
	require.Nil(t, reports)
	require.Zero(t, processed)
}
