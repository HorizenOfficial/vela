package executor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/horizen-pes/pkg/common/appstate"
)

func TestMockRuntime_LoadModule(t *testing.T) {
	runtime := NewMockRuntime()
	defer runtime.Close()

	appId := "test-app-123"
	wasmBytes := []byte("mock-wasm-bytecode")

	serializedState, stateRoot, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	if len(serializedState) == 0 {
		t.Error("Expected non-empty serialized state")
	}

	if stateRoot == ([32]byte{}) {
		t.Error("Expected non-empty state root")
	}

	// Verify we can deserialize the initial state
	var state appstate.ApplicationInternalState
	err = json.Unmarshal(serializedState, &state)
	if err != nil {
		t.Fatalf("Failed to deserialize initial state: %v", err)
	}

	if state.AppID != appId {
		t.Errorf("Expected AppID %s, got %s", appId, state.AppID)
	}

	if len(state.Accounts) != 0 {
		t.Errorf("Expected empty accounts map, got %d accounts", len(state.Accounts))
	}

	if state.Nonce != 0 {
		t.Errorf("Expected nonce 0, got %d", state.Nonce)
	}
}

func TestMockRuntime_ProcessRequest_Deposit(t *testing.T) {
	runtime := NewMockRuntime()
	defer runtime.Close()

	appId := "test-app-123"
	sender := "0x1234567890123456789012345678901234567890"
	value := uint64(1000000000000000000)
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module
	serializedState, _, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	// Make a deposit
	ctx := context.Background()
	newState, events, err := runtime.Deposit(ctx, appId, sender, value, serializedState, wasmBytes)
	if err != nil {
		t.Fatalf("ProcessRequest failed: %v", err)
	}

	// Verify events
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].UserID != sender {
		t.Errorf("Expected event UserID %s, got %s", sender, events[0].UserID)
	}

	// Verify state update
	var state appstate.ApplicationInternalState
	err = json.Unmarshal(newState, &state)
	if err != nil {
		t.Fatalf("Failed to deserialize new state: %v", err)
	}

	if state.Accounts[sender] == nil {
		t.Fatal("Expected sender account to exist")
	}

	if state.Accounts[sender].Balance != value {
		t.Errorf("Expected balance %d, got %d", value, state.Accounts[sender].Balance)
	}

	if state.Nonce != 1 {
		t.Errorf("Expected nonce 1, got %d", state.Nonce)
	}
}

func TestMockRuntime_ProcessRequest_Transfer(t *testing.T) {
	runtime := NewMockRuntime()
	defer runtime.Close()

	appId := "test-app-123"
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module
	serializedState, _, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := "0x1234567890123456789012345678901234567890"
	recipient := "0x0987654321098765432109876543210987654321"
	depositAmount := uint64(2000000000000000000) // 2 ETH
	transferAmount := uint64(500000000000000000) // 0.5 ETH

	// make a deposit
	ctx := context.Background()
	serializedState, _, err = runtime.Deposit(ctx, appId, sender, depositAmount, serializedState, wasmBytes)
	if err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}

	// make a transfer
	transferInstructions := appstate.PayloadInstructions{
		Type: "transfer",
		Transfer: &appstate.TransferInstruction{
			To:     recipient,
			Amount: transferAmount,
		},
	}

	payload, err := json.Marshal(transferInstructions)
	if err != nil {
		t.Fatalf("Failed to marshal transfer instructions: %v", err)
	}

	newState, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, sender, payload, serializedState, wasmBytes)
	if err != nil {
		t.Fatalf("Transfer failed: %v", err)
	}

	// Verify events (should have 2: one for sender, one for recipient)
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}

	// Verify no withdrawals
	if len(withdrawals) != 0 {
		t.Errorf("Expected 0 withdrawals, got %d", len(withdrawals))
	}

	// Verify state update
	var state appstate.ApplicationInternalState
	err = json.Unmarshal(newState, &state)
	if err != nil {
		t.Fatalf("Failed to deserialize new state: %v", err)
	}

	// Check sender balance
	if state.Accounts[sender] == nil {
		t.Fatal("Expected sender account to exist")
	}
	expectedSenderBalance := depositAmount - transferAmount
	if state.Accounts[sender].Balance != expectedSenderBalance {
		t.Errorf("Expected sender balance %d, got %d", expectedSenderBalance, state.Accounts[sender].Balance)
	}

	// Check recipient balance
	if state.Accounts[recipient] == nil {
		t.Fatal("Expected recipient account to exist")
	}
	if state.Accounts[recipient].Balance != transferAmount {
		t.Errorf("Expected recipient balance %d, got %d", transferAmount, state.Accounts[recipient].Balance)
	}

	if state.Nonce != 2 {
		t.Errorf("Expected nonce 2, got %d", state.Nonce)
	}
}

func TestMockRuntime_ProcessRequest_Withdrawal(t *testing.T) {
	runtime := NewMockRuntime()
	defer runtime.Close()

	appId := "test-app-123"
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module
	serializedState, _, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := "0x1234567890123456789012345678901234567890"
	withdrawTo := "0x0987654321098765432109876543210987654321"
	depositAmount := uint64(2000000000000000000) // 2 ETH
	withdrawAmount := uint64(500000000000000000) // 0.5 ETH

	// make a deposit
	ctx := context.Background()
	serializedState, _, err = runtime.Deposit(ctx, appId, sender, depositAmount, serializedState, wasmBytes)
	if err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}

	// make a withdrawal
	withdrawInstructions := appstate.PayloadInstructions{
		Type: "withdraw",
		Withdraw: &appstate.WithdrawInstruction{
			To:     withdrawTo,
			Amount: withdrawAmount,
		},
	}

	payload, err := json.Marshal(withdrawInstructions)
	if err != nil {
		t.Fatalf("Failed to marshal withdraw instructions: %v", err)
	}

	newState, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, sender, payload, serializedState, wasmBytes)
	if err != nil {
		t.Fatalf("Withdrawal failed: %v", err)
	}

	// Verify events (should have 1 for sender)
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].UserID != sender {
		t.Errorf("Expected event UserID %s, got %s", sender, events[0].UserID)
	}

	// Verify withdrawals
	if len(withdrawals) != 1 {
		t.Errorf("Expected 1 withdrawal, got %d", len(withdrawals))
	}

	if withdrawals[0].DestinationAddress != withdrawTo {
		t.Errorf("Expected withdrawal destination %s, got %s", withdrawTo, withdrawals[0].DestinationAddress)
	}

	if withdrawals[0].Amount != 500000000000000000 {
		t.Errorf("Expected withdrawal amount 500000000000000000, got %d", withdrawals[0].Amount)
	}

	// Verify state update
	var state appstate.ApplicationInternalState
	err = json.Unmarshal(newState, &state)
	if err != nil {
		t.Fatalf("Failed to deserialize new state: %v", err)
	}

	// Check sender balance
	if state.Accounts[sender] == nil {
		t.Fatal("Expected sender account to exist")
	}
	expectedBalance := depositAmount - withdrawAmount
	if state.Accounts[sender].Balance != expectedBalance {
		t.Errorf("Expected sender balance %d, got %d", expectedBalance, state.Accounts[sender].Balance)
	}

	if state.Nonce != 2 {
		t.Errorf("Expected nonce 2, got %d", state.Nonce)
	}
}

func TestMockRuntime_ProcessRequest_InsufficientBalance(t *testing.T) {
	runtime := NewMockRuntime()
	defer runtime.Close()

	appId := "test-app-123"
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module first
	serializedState, _, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := "0x1234567890123456789012345678901234567890"
	recipient := "0x0987654321098765432109876543210987654321"
	transferAmount := uint64(1000000000000000000) // 1 ETH

	// Try to transfer without any balance
	transferInstructions := appstate.PayloadInstructions{
		Type: "transfer",
		Transfer: &appstate.TransferInstruction{
			To:     recipient,
			Amount: transferAmount,
		},
	}

	payload, err := json.Marshal(transferInstructions)
	if err != nil {
		t.Fatalf("Failed to marshal transfer instructions: %v", err)
	}

	ctx := context.Background()
	_, _, _, err = runtime.ProcessRequest(ctx, appId, sender, payload, serializedState, wasmBytes)
	if err == nil {
		t.Error("Expected error for insufficient balance, got nil")
	}

	if err.Error() != "sender account 0x1234567890123456789012345678901234567890 does not exist" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestMockRuntime_GenerateDeanonymizationReport(t *testing.T) {
	runtime := NewMockRuntime()
	defer runtime.Close()

	appId := "test-app-123"
	wasmBytes := []byte("mock-wasm-bytecode")
	sender1 := "0x1234567890123456789012345678901234567890"
	sender2 := "0x0987654321098765432109876543210987654321"
	value := uint64(1000000000000000000) // 1 ETH

	// Load module first
	serializedState, _, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	// Deposit for sender1
	ctx := context.Background()
	serializedState, _, err = runtime.Deposit(ctx, appId, sender1, value, serializedState, wasmBytes)
	if err != nil {
		t.Fatalf("First deposit failed: %v", err)
	}

	// Deposit for sender2
	serializedState, _, err = runtime.Deposit(ctx, appId, sender2, value, serializedState, wasmBytes)
	if err != nil {
		t.Fatalf("Second deposit failed: %v", err)
	}

	// Generate deanonymization report
	var requestID string = "report-req"

	reportBytes, err := runtime.GenerateDeanonymizationReport(context.Background(), appId, "report-req", []byte(""), serializedState, wasmBytes)
	if err != nil {
		t.Fatalf("GenerateDeanonymizationReport failed: %v", err)
	}

	// Parse the report
	var report map[string]interface{}
	if err := json.Unmarshal(reportBytes, &report); err != nil {
		t.Fatalf("Failed to unmarshal report: %v", err)
	}

	// Verify report contents
	if report["applicationId"] != appId {
		t.Errorf("Expected applicationId %s, got %v", appId, report["applicationId"])
	}

	if report["requestId"] != requestID {
		t.Errorf("Expected requestId %s, got %v", requestID, report["requestId"])
	}

	totalAccounts, ok := report["accounts"].(map[string]interface{})
	if !ok || len(totalAccounts) != 2 {
		t.Errorf("Expected totalAccounts 2, got %v", report["totalAccounts"])
	}

	nonce, ok := report["nonce"].(float64)
	if !ok || int(nonce) != 2 {
		t.Errorf("Expected nonce 2, got %v", report["nonce"])
	}
}

func TestMockRuntime_Close(t *testing.T) {
	runtime := NewMockRuntime()

	err := runtime.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
