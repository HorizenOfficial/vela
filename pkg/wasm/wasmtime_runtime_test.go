package wasm

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/horizen-pes/pkg/common"
	appState "github.com/horizen-pes/pkg/common/appstate"
	appCommon "github.com/horizen-pes/pkg/wasm/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWasmtimeRuntime_LoadModule(t *testing.T) {
	// Load the compiled WASM module
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err, "Failed to read WASM file")

	// Create runtime
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	// Test LoadModule
	ctx := context.Background()
	appId := "test-app"

	state, stateRoot, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err, "LoadModule should succeed")
	require.NotNil(t, state, "State should not be nil")
	require.NotNil(t, stateRoot, "State root should not be nil")
	require.Len(t, stateRoot, 32, "State root should be 32 bytes")

	// Verify the state is valid JSON
	var stateData appState.ApplicationInternalState
	err = json.Unmarshal(state, &stateData)
	require.NoError(t, err, "State should be valid JSON")

	// Check that the state contains expected fields
	assert.Equal(t, appId, stateData.AppID)
	assert.Equal(t, 0, len(stateData.Accounts))
	assert.True(t, stateData.Nonce >= 0)
}

func TestWasmtimeRuntime_Deposit(t *testing.T) {

	type TestAccountState struct {
		Balance uint64 `json:"balance"`
	}

	type TestStateData struct {
		AppId    string                      `json:"appId"`
		Accounts map[string]TestAccountState `json:"accounts"`
		Nonce    uint64                      `json:"nonce"`
	}

	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err, "Failed to read WASM file")

	// Create runtime
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	// Load module first
	ctx := context.Background()
	appId := "test-app"
	sender := fmt.Sprintf("0xadd%037x", 1)
	value := uint64(1000000000000000000) // 1 ETH

	initialState, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err, "LoadModule should succeed")

	// Test Deposit
	newState, events, err := runtime.Deposit(ctx, appId, sender, value, initialState, wasmBytes)
	require.NoError(t, err, "Deposit should succeed")
	require.NotNil(t, newState, "New state should not be nil")
	require.Len(t, events, 1, "Should generate one event")

	// Verify the event
	event := events[0]
	assert.Equal(t, sender, event.UserID)

	var eventData appCommon.DepositEvent
	err = json.Unmarshal(event.Data, &eventData)
	require.NoError(t, err, "Event data should be valid JSON")
	assert.Equal(t, "deposit", eventData.Type)
	assert.Equal(t, value, eventData.Amount)

	// Verify the state was updated
	var stateData TestStateData
	err = json.Unmarshal(newState, &stateData)
	require.NoError(t, err, "New state should be valid JSON")

	require.Contains(t, stateData.Accounts, sender)
	assert.Equal(t, value, stateData.Accounts[sender].Balance)
}

func TestWasmtimeRuntime_ProcessRequest_Transfer(t *testing.T) {

	type TestAccountState struct {
		Balance uint64 `json:"balance"`
	}

	type TestStateData struct {
		AppId    string                      `json:"appId"`
		Accounts map[string]TestAccountState `json:"accounts"`
		Nonce    uint64                      `json:"nonce"`
	}

	type TestTransferEventData struct {
		Type   string `json:"type"`
		From   string `json:"from,omitempty"`
		To     string `json:"to,omitempty"`
		Amount uint64 `json:"amount"`
	}

	// Load the compiled WASM module
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err, "Failed to read WASM file")

	// Create runtime
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	sender := fmt.Sprintf("0xadd%037x", 1)
	recipient := fmt.Sprintf("0xadd%037x", 2)
	depositValue := uint64(2000000000000000000) // 2 ETH
	transferValue := uint64(500000000000000000) // 0.5 ETH

	// Load module and make a deposit first
	initialState, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err, "LoadModule should succeed")

	stateAfterDeposit, _, err := runtime.Deposit(ctx, appId, sender, depositValue, initialState, wasmBytes)
	require.NoError(t, err, "Deposit should succeed")

	// Create transfer payload
	transferPayload := appState.PayloadInstructions{
		Type: "transfer",
		Transfer: &appState.TransferInstruction{
			To:     recipient,
			Amount: transferValue,
		},
	}
	payloadBytes, err := json.Marshal(transferPayload)
	require.NoError(t, err, "Should marshal transfer payload")

	// Test ProcessRequest for transfer
	newState, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, sender, payloadBytes, stateAfterDeposit, wasmBytes)
	require.NoError(t, err, "ProcessRequest should succeed")
	require.NotNil(t, newState, "New state should not be nil")
	require.Len(t, events, 2, "Should generate two events (sender and recipient)")
	require.Len(t, withdrawals, 0, "Should not generate withdrawals")

	// Verify sender event
	var senderEventData TestTransferEventData
	err = json.Unmarshal(events[0].Data, &senderEventData)
	require.NoError(t, err)
	assert.Equal(t, "transfer_sent", senderEventData.Type)
	assert.Equal(t, recipient, senderEventData.To)
	assert.Equal(t, transferValue, senderEventData.Amount)

	// Verify recipient event
	var recipientEventData TestTransferEventData
	err = json.Unmarshal(events[1].Data, &recipientEventData)
	require.NoError(t, err)
	assert.Equal(t, "transfer_received", recipientEventData.Type)
	assert.Equal(t, sender, recipientEventData.From)
	assert.Equal(t, transferValue, recipientEventData.Amount)

	// Verify the state was updated
	var stateData TestStateData
	err = json.Unmarshal(newState, &stateData)
	require.NoError(t, err, "New state should be valid JSON")

	assert.Equal(t, depositValue-transferValue, stateData.Accounts[sender].Balance)
	assert.Equal(t, transferValue, stateData.Accounts[recipient].Balance)
}

func TestWasmtimeRuntime_ProcessRequest_Withdrawal(t *testing.T) {

	type TestAccountState struct {
		Balance uint64 `json:"balance"`
	}

	type TestStateData struct {
		AppId    string                      `json:"appId"`
		Accounts map[string]TestAccountState `json:"accounts"`
		Nonce    uint64                      `json:"nonce"`
	}

	type TestWithdrawalEventData struct {
		Type   string `json:"type"`
		To     string `json:"to"`
		Amount uint64 `json:"amount"`
	}

	// Load the compiled WASM module
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err, "Failed to read WASM file")

	// Create runtime
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	sender := fmt.Sprintf("0xadd%037x", 1)
	depositValue := uint64(1000000000000000000) // 1 ETH
	withdrawValue := uint64(500000000000000000) // 0.5 ETH
	withdrawAddress := "0x1234567890123456789012345678901234567890"

	// Load module and make a deposit first
	initialState, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err, "LoadModule should succeed")

	stateAfterDeposit, _, err := runtime.Deposit(ctx, appId, sender, depositValue, initialState, wasmBytes)
	require.NoError(t, err, "Deposit should succeed")

	// Create withdrawal payload
	withdrawPayload := appState.PayloadInstructions{
		Type: "withdraw",
		Withdraw: &appState.WithdrawInstruction{
			To:     withdrawAddress,
			Amount: withdrawValue,
		},
	}
	payloadBytes, err := json.Marshal(withdrawPayload)
	require.NoError(t, err, "Should marshal withdrawal payload")

	// Test ProcessRequest for withdrawal
	newState, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, sender, payloadBytes, stateAfterDeposit, wasmBytes)
	require.NoError(t, err, "ProcessRequest should succeed")
	require.NotNil(t, newState, "New state should not be nil")
	require.Len(t, events, 1, "Should generate one event")
	require.Len(t, withdrawals, 1, "Should generate one withdrawal")

	// Verify withdrawal event
	event := events[0]
	assert.Equal(t, sender, event.UserID)

	var eventData TestWithdrawalEventData
	err = json.Unmarshal(event.Data, &eventData)
	require.NoError(t, err, "Event data should be valid JSON")
	assert.Equal(t, "withdrawal", eventData.Type)
	assert.Equal(t, withdrawAddress, eventData.To)
	assert.Equal(t, withdrawValue, eventData.Amount)

	// Verify withdrawal
	withdrawal := withdrawals[0]
	assert.Equal(t, withdrawAddress, withdrawal.DestinationAddress)
	assert.Equal(t, withdrawValue, withdrawal.Amount)

	// Verify the state was updated
	var stateData TestStateData
	err = json.Unmarshal(newState, &stateData)
	require.NoError(t, err, "New state should be valid JSON")

	assert.Equal(t, depositValue-withdrawValue, stateData.Accounts[sender].Balance)
}

func TestWasmtimeRuntime_GenerateDeanonymizationReport(t *testing.T) {

	type TestAccountState struct {
		Balance uint64 `json:"balance"`
	}

	type TestDeanonymizationReport struct {
		ApplicationId string                      `json:"applicationId"`
		RequestId     string                      `json:"requestId"`
		Accounts      map[string]TestAccountState `json:"accounts"`
		Nonce         uint64                      `json:"nonce"`
	}

	// Load the compiled WASM module
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err, "Failed to read WASM file")

	// Create runtime
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	requestId := "deanon-1"
	sender := fmt.Sprintf("0xadd%037x", 1)
	value := uint64(1000000000000000000) // 1 ETH

	// Load module and make a deposit first to have some state
	initialState, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err, "LoadModule should succeed")

	stateWithData, _, err := runtime.Deposit(ctx, appId, sender, value, initialState, wasmBytes)
	require.NoError(t, err, "Deposit should succeed")

	// Test GenerateDeanonymizationReport
	report, err := runtime.GenerateDeanonymizationReport(ctx, appId, requestId, nil, stateWithData, wasmBytes)
	require.NoError(t, err, "GenerateDeanonymizationReport should succeed")
	require.NotNil(t, report, "Report should not be nil")

	var reportData TestDeanonymizationReport
	err = json.Unmarshal(report, &reportData)
	require.NoError(t, err, "Report should be valid JSON")

	assert.Equal(t, appId, reportData.ApplicationId)
	assert.Equal(t, requestId, reportData.RequestId)
	require.Contains(t, reportData.Accounts, sender)
	assert.Equal(t, value, reportData.Accounts[sender].Balance)
}

func TestWasmtimeRuntime_FullWorkflow(t *testing.T) {

	type TestAccountState struct {
		Balance uint64 `json:"balance"`
	}

	type TestStateData struct {
		AppId    string                      `json:"appId"`
		Accounts map[string]TestAccountState `json:"accounts"`
		Nonce    uint64                      `json:"nonce"`
	}

	// Load the compiled WASM module
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err, "Failed to read WASM file")

	// Create runtime
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "payment-app"
	user1 := fmt.Sprintf("0xadd%037x", 1)
	user2 := fmt.Sprintf("0xadd%037x", 2)

	t.Log("Step 1: Load module")
	state, stateRoot, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err, "LoadModule should succeed")
	require.NotNil(t, state)
	require.NotNil(t, stateRoot)

	t.Log("Step 2: Make deposit for user1")
	depositValue := uint64(2000000000000000000) // 2 ETH
	state, events, err := runtime.Deposit(ctx, appId, user1, depositValue, state, wasmBytes)
	require.NoError(t, err, "Deposit should succeed")
	require.Len(t, events, 1)

	t.Log("Step 3: Transfer from user1 to user2")
	transferValue := uint64(500000000000000000) // 0.5 ETH
	transferPayload := appState.PayloadInstructions{
		Type: "transfer",
		Transfer: &appState.TransferInstruction{
			To:     user2,
			Amount: transferValue,
		},
	}
	payloadBytes, err := json.Marshal(transferPayload)
	require.NoError(t, err)

	state, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, user1, payloadBytes, state, wasmBytes)
	require.NoError(t, err, "Transfer should succeed")
	require.Len(t, events, 2)
	require.Len(t, withdrawals, 0)

	t.Log("Step 4: Withdraw from user2")
	// Create withdrawal payload
	withdrawValue := uint64(250000000000000000) // 0.25 ETH
	withdrawAddress := "0x1234567890123456789012345678901234567890"
	withdrawPayload := appState.PayloadInstructions{
		Type: "withdraw",
		Withdraw: &appState.WithdrawInstruction{
			To:     withdrawAddress,
			Amount: withdrawValue,
		},
	}
	payloadBytes, err = json.Marshal(withdrawPayload)
	require.NoError(t, err)

	state, events, withdrawals, err = runtime.ProcessRequest(ctx, appId, user2, payloadBytes, state, wasmBytes)
	require.NoError(t, err, "Withdrawal should succeed")
	require.Len(t, events, 1)
	require.Len(t, withdrawals, 1)

	t.Log("Step 5: Generate deanonymization report")
	report, err := runtime.GenerateDeanonymizationReport(ctx, appId, "deanon-1", nil, state, wasmBytes)
	require.NoError(t, err, "Deanonymization report should succeed")
	require.NotNil(t, report)

	var stateData TestStateData
	err = json.Unmarshal(state, &stateData)
	require.NoError(t, err)

	assert.Equal(t, depositValue-transferValue, stateData.Accounts[user1].Balance)
	assert.Equal(t, transferValue-withdrawValue, stateData.Accounts[user2].Balance)

	t.Log("Full workflow completed successfully!")
}

func TestWasmtimeRuntime_ConcurrentModuleLoading(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	var wg sync.WaitGroup
	errors := make(chan error, 3)

	// Load the same module concurrently
	for i := range 3 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			appId := fmt.Sprintf("concurrent-app-%d", id)
			_, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		require.NoError(t, err, "Concurrent module loading should succeed")
	}
}

func TestWasmtimeRuntime_ModuleCaching(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "cached-app"

	// First load
	module1, err := runtime.getOrLoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	// Second load should use cached module, pass empty bytes to verify cache hit
	module2, err := runtime.getOrLoadModule(ctx, appId, []byte{})
	require.NoError(t, err)

	assert.Equal(t, module1, module2, "Cached module should be the same as loaded module")
}

func TestWasmtimeRuntime_LargeStateHandling(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "large-state-app"

	// Load module and make many deposits to create large state
	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	// Make 100 deposits to create a large state
	// Tested up to 6k, after 6k app is very slow and failing randomly: TODO check this!
	for i := range 100 {
		user := fmt.Sprintf("0xadd%037x", i)
		value := uint64(1000000000000000000) // 1 ETH
		state, _, err = runtime.Deposit(ctx, appId, user, value, state, wasmBytes)
		require.NoError(t, err)
	}

	// Verify the large state can still be processed
	transferPayload := appState.PayloadInstructions{
		Type: "transfer",
		Transfer: &appState.TransferInstruction{
			To:     fmt.Sprintf("0xadd%037x", 1),
			Amount: uint64(500000000000000000),
		},
	}
	payloadBytes, err := json.Marshal(transferPayload)
	require.NoError(t, err)

	sender := fmt.Sprintf("0xadd%037x", 0)
	_, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, sender, payloadBytes, state, wasmBytes)
	require.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Len(t, withdrawals, 0)
}

func TestWasmtimeRuntime_InvalidWasmModule(t *testing.T) {
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "invalid-app"
	invalidWasm := []byte("invalid wasm bytes")

	_, _, err := runtime.LoadModule(ctx, appId, invalidWasm)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to compile WASM module")
}

func TestWasmtimeRuntime_EmptyWasmModule(t *testing.T) {
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "empty-app"
	emptyWasm := []byte{}

	_, _, err := runtime.LoadModule(ctx, appId, emptyWasm)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to compile WASM module")
}

func TestWasmtimeRuntime_NilInputs(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	user1 := fmt.Sprintf("0xadd%037x", 1)

	t.Run("NilWasmBytes", func(t *testing.T) {
		_, _, err := runtime.LoadModule(ctx, "test-app", nil)
		assert.Error(t, err)
	})

	t.Run("EmptyAppId", func(t *testing.T) {
		_, _, err := runtime.LoadModule(ctx, "", wasmBytes)
		// This might succeed depending on implementation, but state should be testable
		if err == nil {
			// Verify we can't use empty app ID for operations
			_, _, err = runtime.Deposit(ctx, "", user1, 1000, []byte("{}"), wasmBytes)
			assert.Error(t, err)
		}
	})

	t.Run("NilState", func(t *testing.T) {
		_, _, err := runtime.Deposit(ctx, "test-app", user1, 1000, nil, wasmBytes)
		assert.Error(t, err)
	})
}

func TestWasmtimeRuntime_InvalidPayloads(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	user1 := fmt.Sprintf("0xadd%037x", 1)
	user2 := fmt.Sprintf("0xadd%037x", 2)

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	t.Run("InvalidJSON", func(t *testing.T) {
		invalidPayload := []byte("invalid json")
		_, _, _, err := runtime.ProcessRequest(ctx, appId, user1, invalidPayload, state, wasmBytes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to parse payload instructions")
	})

	t.Run("MissingTransferFields", func(t *testing.T) {
		incompletePayload := map[string]interface{}{
			"type": "transfer",
			// Missing transfer field
		}
		payloadBytes, _ := json.Marshal(incompletePayload)
		_, _, _, err := runtime.ProcessRequest(ctx, appId, user1, payloadBytes, state, wasmBytes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Transfer instruction is missing")
	})

	t.Run("AccountDoesNotExist", func(t *testing.T) {
		negativePayload := map[string]interface{}{
			"type": "transfer",
			"transfer": map[string]interface{}{
				"to":     user2,
				"amount": uint64(500),
			},
		}
		payloadBytes, _ := json.Marshal(negativePayload)
		_, _, _, err := runtime.ProcessRequest(ctx, appId, user1, payloadBytes, state, wasmBytes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Account does not exist")
	})
}

func TestWasmtimeRuntime_InsufficientFunds(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	user1 := fmt.Sprintf("0xadd%037x", 1)
	user2 := fmt.Sprintf("0xadd%037x", 2)
	value := uint64(12345678901234567890) // # fits in uint64, > max int64

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	state, _, err = runtime.Deposit(ctx, appId, user1, value/2, state, wasmBytes)
	require.NoError(t, err)

	// Try to transfer without enough funds
	transferPayload := appState.PayloadInstructions{
		Type: "transfer",
		Transfer: &appState.TransferInstruction{
			To:     user2,
			Amount: value,
		},
	}
	payloadBytes, _ := json.Marshal(transferPayload)

	_, _, _, err = runtime.ProcessRequest(ctx, appId, user1, payloadBytes, state, wasmBytes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Insufficient balance for transfer")
}

func TestWasmtimeRuntime_LargePayload(t *testing.T) {
	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "memory-test-app"

	// Create an extremely large payload
	largePayload := make([]byte, 10*1024*1024) // 10MB
	for i := range largePayload {
		largePayload[i] = byte(i % 256)
	}

	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	user1 := fmt.Sprintf("0xadd%037x", 1)
	_, _, _, err = runtime.ProcessRequest(ctx, appId, user1, largePayload, state, wasmBytes)
	// should not panic but return an error that the payload does not conform to the expected format
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to parse payload instructions")
}

func TestWasmtimeRuntime_InvalidStateFormat(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	user1 := fmt.Sprintf("0xadd%037x", 1)

	// Use corrupted state
	corruptedState := []byte("corrupted state data")

	state, events, err := runtime.Deposit(ctx, appId, user1, 1000, corruptedState, wasmBytes)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Failed to parse application state")
	require.Equal(t, []byte(nil), state)
	require.Equal(t, []common.PlainEvent(nil), events)
}

func TestWasmtimeRuntime_StateRootConsistency(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "consistency-test-app"

	// Load same module multiple times and verify state root consistency
	for i := 0; i < 5; i++ {
		state, stateRoot, err := runtime.LoadModule(ctx, fmt.Sprintf("%s-%d", appId, i), wasmBytes)
		require.NoError(t, err)

		// Verify state root is exactly 32 bytes (SHA256)
		assert.Len(t, stateRoot, 32)

		// Verify state root is deterministic
		expectedHash := sha256.Sum256(state)
		assert.Equal(t, expectedHash, stateRoot)
	}
}

func TestWasmtimeRuntime_ZeroValueOperations(t *testing.T) {
	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "zero-value-app"
	user1 := fmt.Sprintf("0xadd%037x", 1)

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	// Test zero value deposit
	newState, events, err := runtime.Deposit(ctx, appId, user1, 0, state, wasmBytes)

	require.NoError(t, err, "Deposit with zero value should succeed")
	require.Len(t, events, 0, "Zero value deposit should not generate any events")
	require.NotNil(t, state, "State should not be nil after zero value deposit")
	require.Equal(t, state, newState)
}

func TestWasmtimeRuntime_InvalidInstruction(t *testing.T) {

	wasmPath := filepath.Join("wasm-go", "payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err)

	runtime := NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "invalid-instruction-app"

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	invalidPayload := map[string]interface{}{
		"type": "unsupported_instruction",
	}
	payloadBytes, err := json.Marshal(invalidPayload)
	require.NoError(t, err)

	user1 := fmt.Sprintf("0xadd%037x", 1)
	_, _, _, err = runtime.ProcessRequest(ctx, appId, user1, payloadBytes, state, wasmBytes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported instruction type")
}
