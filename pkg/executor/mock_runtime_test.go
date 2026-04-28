package executor

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"math/big"

	"github.com/HorizenOfficial/vela/pkg/logger"

	ethCommon "github.com/ethereum/go-ethereum/common"
	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/common"
)

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

	serializedState, fuel, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	if len(serializedState) == 0 {
		t.Error("Expected non-empty serialized state")
	}

	if fuel.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Expected 10 fuel, got %s", fuel.String())
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
	depositAmount := big.NewInt(1000000000000000000)
	wasmBytes := []byte("mock-wasm-bytecode")

	// Load module
	serializedState, fuel, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}
	if fuel.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Expected 10 fuel, got %s", fuel.String())
	}

	// Make a deposit
	ctx := context.Background()
	newState, events, _, fuel, failure := runtime.Deposit(ctx, appId, sender, velacommon.NativeTokenAddress(), depositAmount, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("ProcessRequest failed: %v", failure)
	}
	if fuel.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Expected 10 fuel, got %s", fuel.String())
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

	if state.Accounts[sender].Balance.ToInt().Cmp(depositAmount) != 0 {
		t.Errorf("Expected balance %s, got %s", depositAmount, state.Accounts[sender].Balance)
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
	serializedState, fuel, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	recipient := ethCommon.HexToAddress("0x0987654321098765432109876543210987654321")
	depositAmount := big.NewInt(2000000000000000000)    // 2 ETH
	transferAmount := common.NewBig(500000000000000000) // 0.5 ETH

	// make a deposit
	ctx := context.Background()
	serializedState, _, _, _, failure := runtime.Deposit(ctx, appId, sender, velacommon.NativeTokenAddress(), depositAmount, serializedState, wasmBytes)
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

	newState, events, _, withdrawals, _, fuel, failure := runtime.ProcessRequest(ctx, appId, sender, common.Process, payload, serializedState, wasmBytes)
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

	if fuel.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Expected 10 fuel, got %s", fuel.String())
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
	expectedSenderBalance := common.ToBig(new(big.Int).Sub(depositAmount, transferAmount.ToInt()))
	if state.Accounts[sender].Balance.ToInt().Cmp(expectedSenderBalance.ToInt()) != 0 {
		t.Errorf("Expected sender balance %s, got %s", expectedSenderBalance, state.Accounts[sender].Balance)
	}

	// Check recipient balance
	if state.Accounts[recipient] == nil {
		t.Fatal("Expected recipient account to exist")
	}
	if state.Accounts[recipient].Balance.ToInt().Cmp(transferAmount.ToInt()) != 0 {
		t.Errorf("Expected recipient balance %s, got %s", transferAmount, state.Accounts[recipient].Balance)
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
	serializedState, _, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	withdrawTo := ethCommon.HexToAddress("0x0987654321098765432109876543210987654321")
	depositAmount := big.NewInt(2000000000000000000)    // 2 ETH
	withdrawAmount := common.NewBig(500000000000000000) // 0.5 ETH

	// make a deposit
	ctx := context.Background()
	serializedState, _, _, _, failure := runtime.Deposit(ctx, appId, sender, velacommon.NativeTokenAddress(), depositAmount, serializedState, wasmBytes)
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

	newState, events, _, withdrawals, _, fuel, failure := runtime.ProcessRequest(ctx, appId, sender, common.Process, payload, serializedState, wasmBytes)
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

	if fuel.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Expected 10 fuel, got %s", fuel.String())
	}

	// Verify withdrawals
	if len(withdrawals) != 1 {
		t.Errorf("Expected 1 withdrawal, got %d", len(withdrawals))
	}

	if withdrawals[0].DestinationAddress != withdrawTo {
		t.Errorf("Expected withdrawal destination %s, got %s", withdrawTo, withdrawals[0].DestinationAddress)
	}

	if withdrawals[0].Amount.ToInt().Cmp(withdrawAmount.ToInt()) != 0 {
		t.Errorf("Expected withdrawal amount %s, got %s", withdrawAmount, withdrawals[0].Amount)
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
	expectedBalance := common.ToBig(new(big.Int).Sub(depositAmount, withdrawAmount.ToInt()))
	if state.Accounts[sender].Balance.ToInt().Cmp(expectedBalance.ToInt()) != 0 {
		t.Errorf("Expected sender balance %s, got %s", expectedBalance, state.Accounts[sender].Balance)
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
	serializedState, _, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	sender := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	recipient := ethCommon.HexToAddress("0x0987654321098765432109876543210987654321")
	transferAmount := common.NewBig(1000000000000000000) // 1 ETH

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
	_, _, _, _, _, _, failure := runtime.ProcessRequest(ctx, appId, sender, common.Process, payload, serializedState, wasmBytes)
	if failure == nil {
		t.Error("Expected error for insufficient balance, got nil")
	}

	if failure.Error() != "sender account 0x1234567890123456789012345678901234567890 does not exist" {
		t.Errorf("Expected specific error message, got: %v", failure)
	}
}

func TestMockRuntime_DeanonymizationViaProcessRequest(t *testing.T) {
	runtime := NewMockRuntime(testLogger)
	defer runtime.Close()

	appId := common.NewApplicationId(123)
	wasmBytes := []byte("mock-wasm-bytecode")
	sender1 := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	sender2 := ethCommon.HexToAddress("0x0987654321098765432109876543210987654321")
	depositAmount := big.NewInt(1000000000000000000) // 1 ETH

	// Load module first
	serializedState, _, err := runtime.LoadModule(context.Background(), appId, wasmBytes)
	if err != nil {
		t.Fatalf("LoadModule failed: %v", err)
	}

	// Deposit for sender1
	ctx := context.Background()
	serializedState, _, _, _, failure := runtime.Deposit(ctx, appId, sender1, velacommon.NativeTokenAddress(), depositAmount, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("First deposit failed: %v", failure)
	}

	// Deposit for sender2
	serializedState, _, _, _, failure = runtime.Deposit(ctx, appId, sender2, velacommon.NativeTokenAddress(), depositAmount, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("Second deposit failed: %v", failure)
	}

	// Generate deanonymization report via ProcessRequest with type "deanonymize"
	deanonPayload := testPayloadInstructions{
		Type: "deanonymize",
		Deanonymize: &testDeanonymizeInstruction{
			Tag: "dummytag",
		},
	}
	payload, err := json.Marshal(deanonPayload)
	if err != nil {
		t.Fatalf("Failed to marshal deanonymize payload: %v", err)
	}

	_, _, _, _, reportBytes, _, failure := runtime.ProcessRequest(context.Background(), appId, sender1, common.Deanonymize, payload, serializedState, wasmBytes)
	if failure != nil {
		t.Fatalf("Deanonymize ProcessRequest failed: %v", failure)
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
	tag, ok := report["tag"].(string)
	if !ok || tag != "dummytag" {
		t.Errorf("Expected tag in the report")
	}
}

func TestMockRuntime_Close(t *testing.T) {
	runtime := NewMockRuntime(testLogger)

	err := runtime.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
