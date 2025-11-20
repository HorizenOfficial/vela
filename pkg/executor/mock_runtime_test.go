package executor

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"math/big"

	"github.com/horizen-pes/pkg/logger"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
)

/*
// Local mirror types used in tests to avoid importing wasm-go/app
type testAccountState struct {
	Address string `json:"address"`
	Balance uint64 `json:"balance"`
}

type testApplicationInternalState struct {
	AppID    string                       `json:"appId"`
	Accounts map[string]*testAccountState `json:"accounts"`
	Nonce    uint64                       `json:"nonce"`
}

type testTransferInstruction struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

type testWithdrawInstruction struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

type testPayloadInstructions struct {
	Type     string                   `json:"type"`
	Transfer *testTransferInstruction `json:"transfer,omitempty"`
	Withdraw *testWithdrawInstruction `json:"withdraw,omitempty"`
}
*/

var testLogger logger.Logger

func TestMain(m *testing.M) {
	// Initialize once, by default it writes on stderr
	//testLogger = logger.NewLogger(&logger.Config{Kind: "printf"})

	testLogger = logger.NewLogger(
		&logger.Config{
			Kind:         "zerolog",
			ConsoleColor: false, // colors can print escape chars on tty
			Console:      true,
			ConsoleLevel: "trace",
			//FileName:     "qqq.log",
			//FileLevel:    "info",
		},
	)
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestMockRuntime_LoadModule(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	defer runtime.Close()

	appId := common.NewApplicationId(1)
	wasmBytes := []byte("mock-wasm-bytecode")

	serializedState, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	if len(serializedState) == 0 {
		t.Error("Expected non-empty serialized state")
	}

	// Verify we can deserialize the initial state
	var state testApplicationInternalState
	err = json.Unmarshal(serializedState, &state)
	if err != nil {
		t.Fatalf("Failed to deserialize initial state: %v", err)
	}

	if state.AppID != appId {
		t.Errorf("Expected AppID %d, got %d", appId, state.AppID)
	}

	if len(state.Accounts) != 0 {
		t.Errorf("Expected empty accounts map, got %d accounts", len(state.Accounts))
	}

	if state.Nonce != 0 {
		t.Errorf("Expected nonce 0, got %d", state.Nonce)
	}
}

func TestMockRuntime_ProcessRequest_Deposit(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	defer runtime.Close()

	appId := common.NewApplicationId(123)
	sender := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	value := big.NewInt(1000000000000000000)
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module
	serializedState, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	// Make a deposit
	ctx := context.Background()
	newState, events, failure := runtime.Deposit(ctx, appId, sender, value, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("ProcessRequest failed: %v", failure)
	}

	// Verify events
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].UserID != sender {
		t.Errorf("Expected event UserID %s, got %s", sender, events[0].UserID)
	}

	// Verify state update
	var state testApplicationInternalState
	err = json.Unmarshal(newState, &state)
	if err != nil {
		t.Fatalf("Failed to deserialize new state: %v", err)
	}

	if state.Accounts[sender] == nil {
		t.Fatal("Expected sender account to exist")
	}

	if state.Accounts[sender].Balance.Cmp(value) != 0 {
		t.Errorf("Expected balance %d, got %d", value, state.Accounts[sender].Balance)
	}

	if state.Nonce != 1 {
		t.Errorf("Expected nonce 1, got %d", state.Nonce)
	}
}

func TestMockRuntime_ProcessRequest_Transfer(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	defer runtime.Close()

	appId := common.NewApplicationId(123)
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module
	serializedState, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	recipient := ethCommon.HexToAddress("0x0987654321098765432109876543210987654321")
	depositAmount := big.NewInt(2000000000000000000) // 2 ETH
	transferAmount := big.NewInt(500000000000000000) // 0.5 ETH

	// make a deposit
	ctx := context.Background()
	serializedState, _, failure := runtime.Deposit(ctx, appId, sender, depositAmount, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("Deposit failed: %v", failure)
	}

	// make a transfer
	transferInstructions := testPayloadInstructions{
		Type: "transfer",
		Transfer: &testTransferInstruction{
			To:     recipient,
			Amount: transferAmount,
		},
	}

	payload, err := json.Marshal(transferInstructions)
	if err != nil {
		t.Fatalf("Failed to marshal transfer instructions: %v", err)
	}

	newState, events, withdrawals, failure := runtime.ProcessRequest(ctx, appId, sender, payload, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("Transfer failed: %v", failure)
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
	var state testApplicationInternalState
	err = json.Unmarshal(newState, &state)
	if err != nil {
		t.Fatalf("Failed to deserialize new state: %v", err)
	}

	// Check sender balance
	if state.Accounts[sender] == nil {
		t.Fatal("Expected sender account to exist")
	}
	expectedSenderBalance := new(big.Int).Sub(depositAmount, transferAmount)
	if state.Accounts[sender].Balance.Cmp(expectedSenderBalance) != 0 {
		t.Errorf("Expected sender balance %d, got %d", expectedSenderBalance, state.Accounts[sender].Balance)
	}

	// Check recipient balance
	if state.Accounts[recipient] == nil {
		t.Fatal("Expected recipient account to exist")
	}
	if state.Accounts[recipient].Balance.Cmp(transferAmount) != 0 {
		t.Errorf("Expected recipient balance %d, got %d", transferAmount, state.Accounts[recipient].Balance)
	}

	if state.Nonce != 2 {
		t.Errorf("Expected nonce 2, got %d", state.Nonce)
	}
}

func TestMockRuntime_ProcessRequest_Withdrawal(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	defer runtime.Close()

	appId := common.NewApplicationId(123)
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module
	serializedState, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	withdrawTo := ethCommon.HexToAddress("0x0987654321098765432109876543210987654321")
	depositAmount := big.NewInt(2000000000000000000) // 2 ETH
	withdrawAmount := big.NewInt(500000000000000000) // 0.5 ETH

	// make a deposit
	ctx := context.Background()
	serializedState, _, failure := runtime.Deposit(ctx, appId, sender, depositAmount, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("Deposit failed: %v", failure)
	}

	// make a withdrawal
	withdrawInstructions := testPayloadInstructions{
		Type: "withdraw",
		Withdraw: &testWithdrawInstruction{
			To:     withdrawTo,
			Amount: withdrawAmount,
		},
	}

	payload, err := json.Marshal(withdrawInstructions)
	if err != nil {
		t.Fatalf("Failed to marshal withdraw instructions: %v", err)
	}

	newState, events, withdrawals, failure := runtime.ProcessRequest(ctx, appId, sender, payload, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("Withdrawal failed: %v", failure)
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

	if withdrawals[0].Amount.Cmp(withdrawAmount) != 0 {
		t.Errorf("Expected withdrawal amount %d, got %d", withdrawAmount, withdrawals[0].Amount)
	}

	// Verify state update
	var state testApplicationInternalState
	err = json.Unmarshal(newState, &state)
	if err != nil {
		t.Fatalf("Failed to deserialize new state: %v", err)
	}

	// Check sender balance
	if state.Accounts[sender] == nil {
		t.Fatal("Expected sender account to exist")
	}
	expectedBalance := new(big.Int).Sub(depositAmount, withdrawAmount)
	if state.Accounts[sender].Balance.Cmp(expectedBalance) != 0 {
		t.Errorf("Expected sender balance %d, got %d", expectedBalance, state.Accounts[sender].Balance)
	}

	if state.Nonce != 2 {
		t.Errorf("Expected nonce 2, got %d", state.Nonce)
	}
}

func TestMockRuntime_ProcessRequest_InsufficientBalance(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	defer runtime.Close()

	appId := common.NewApplicationId(123)
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module first
	serializedState, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	recipient := ethCommon.HexToAddress("0x0987654321098765432109876543210987654321")
	transferAmount := big.NewInt(1000000000000000000) // 1 ETH

	// Try to transfer without any balance
	transferInstructions := testPayloadInstructions{
		Type: "transfer",
		Transfer: &testTransferInstruction{
			To:     recipient,
			Amount: transferAmount,
		},
	}

	payload, err := json.Marshal(transferInstructions)
	if err != nil {
		t.Fatalf("Failed to marshal transfer instructions: %v", err)
	}

	ctx := context.Background()
	_, _, _, failure := runtime.ProcessRequest(ctx, appId, sender, payload, serializedState, wasmBytes)
	if failure == nil {
		t.Error("Expected error for insufficient balance, got nil")
	}

	if failure.Error() != "sender account 0x1234567890123456789012345678901234567890 does not exist" {
		t.Errorf("Expected specific error message, got: %v", failure)
	}
}

func TestMockRuntime_GenerateDeanonymizationReport(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	defer runtime.Close()

	appId := common.NewApplicationId(123)
	wasmBytes := []byte("mock-wasm-bytecode")
	sender1 := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	sender2 := ethCommon.HexToAddress("0x0987654321098765432109876543210987654321")
	value := big.NewInt(1000000000000000000) // 1 ETH

	// Load module first
	serializedState, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	// Deposit for sender1
	ctx := context.Background()
	serializedState, _, failure := runtime.Deposit(ctx, appId, sender1, value, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("First deposit failed: %v", failure)
	}

	// Deposit for sender2
	serializedState, _, failure = runtime.Deposit(ctx, appId, sender2, value, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("Second deposit failed: %v", failure)
	}

	// Generate deanonymization report
	reportBytes, failure := runtime.GenerateDeanonymizationReport(context.Background(), appId, []byte(""), serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("GenerateDeanonymizationReport failed: %v", failure)
	}

	// Parse the report
	var report map[string]interface{}
	if err := json.Unmarshal(reportBytes, &report); err != nil {
		t.Fatalf("Failed to unmarshal report: %v", err)
	}

	// Verify report contents
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
	runtime := NewMockRuntime(testLogger)

	err := runtime.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
