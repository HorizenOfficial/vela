package main_test

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"os/exec"
	"testing"

	"github.com/horizen-pes/app/simple/app"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/testutil"
	pes_wasm "github.com/horizen-pes/pkg/wasm"
	"github.com/stretchr/testify/require"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

var _ = []common.Withdrawal{}

const (
	wasmModulePath = "build/simple_app.wasm"

)


var (
	appId = common.NewApplicationId(1)
	user1Address      = ethCommon.HexToAddress("0xadd0000000000000000000000000000000000001")
	user2Address      = ethCommon.HexToAddress("0xadd0000000000000000000000000000000000002")
	recipient1Address = ethCommon.HexToAddress("0xadd0000000000000000000000000000000000003")

)

// buildAndLoadWasmModule runs `make build` to compile and load the wasm module.
func buildAndLoadWasmModule(t *testing.T) []byte {
	t.Helper()
	cmd := exec.Command("make", "build")
	cmd.Dir = "." // Run in the current directory
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "failed to build wasm module: %s", string(output))

	// Read the wasm module
	wasmBytes, err := os.ReadFile(wasmModulePath)
	require.NoError(t, err)
	require.NotEmpty(t, wasmBytes)
	return wasmBytes
}

func TestSimpleAppIntegration(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module
	initialStateBytes, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	var initialState app.ApplicationInternalState
	err = json.Unmarshal(initialStateBytes, &initialState)
	require.NoError(t, err)
	require.Equal(t, appId, common.NewApplicationId(initialState.AppID))
	require.Empty(t, initialState.Accounts)

	// 2. User1 Deposits funds
	deposit1Amount := big.NewInt(1000)
	depositState1Bytes, depositEvents, err := runtime.Deposit(ctx, appId, user1Address, deposit1Amount, initialStateBytes, wasmBytes)
	require.NoError(t, err)
	require.NotNil(t, depositState1Bytes)
	require.Len(t, depositEvents, 1)

	var depositState app.ApplicationInternalState
	err = json.Unmarshal(depositState1Bytes, &depositState)
	require.NoError(t, err)
	require.Len(t, depositState.Accounts, 1)
	require.Equal(t, deposit1Amount, depositState.Accounts[user1Address].Balance)

	// 2. User2 Deposits funds (more than previous user)
	deposit2Amount := big.NewInt(2000)
	depositState2Bytes, depositEvents, err := runtime.Deposit(ctx, appId, user2Address, deposit2Amount, depositState1Bytes, wasmBytes)
	require.NoError(t, err)
	require.NotNil(t, depositState2Bytes)
	require.Len(t, depositEvents, 1)

	err = json.Unmarshal(depositState2Bytes, &depositState)
	require.NoError(t, err)
	require.Len(t, depositState.Accounts, 2)
	require.Equal(t, deposit2Amount, depositState.Accounts[user2Address].Balance)

	// 3. Process a withdraw request for user1
	withdrawAmount := big.NewInt(200)
	withdrawInstruction := app.WithdrawInstruction{
		To:     recipient1Address,
		Amount: withdrawAmount,
	}
	payload := app.PayloadInstructions{
		Type:     "withdraw",
		Withdraw: &withdrawInstruction,
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	withdrawStateBytes, withdrawEvents, withdrawals, err := runtime.ProcessRequest(ctx, appId, user1Address, payloadBytes, depositState2Bytes, wasmBytes)
	require.NoError(t, err)
	require.NotNil(t, withdrawStateBytes)
	require.Len(t, withdrawEvents, 1)
	require.Len(t, withdrawals, 1)

	var withdrawState app.ApplicationInternalState
	err = json.Unmarshal(withdrawStateBytes, &withdrawState)
	require.NoError(t, err)

	diffBalance := new(big.Int).Sub(deposit1Amount, withdrawAmount)
	require.Equal(t, diffBalance, withdrawState.Accounts[user1Address].Balance)

	require.Equal(t, recipient1Address, withdrawals[0].DestinationAddress)
	require.Equal(t, withdrawAmount, withdrawals[0].Amount)

	// 4. Generate deanonymization report
	payloadJSON := `{"tag":"my_custom_tag"}`
	payloadBytes = []byte(payloadJSON)
	reportBytes, err := runtime.GenerateDeanonymizationReport(ctx, appId, payloadBytes, withdrawStateBytes, wasmBytes)
	require.NoError(t, err)
	require.NotNil(t, reportBytes)

	var report map[string]interface{}
	err = json.Unmarshal(reportBytes, &report)
	require.NoError(t, err)
	require.Equal(t, "my_custom_tag", report["tag"])
	t.Log(" Report:\n", testutil.PrettyPrintJSON(report))

	// 5. Compare addresses
	compareInstructions := app.CompareInstructions{
		TargetAddress: user2Address,
	}
	payload = app.PayloadInstructions{
		Type:            "compare_addresses",
		CompareAccounts: &compareInstructions,
	}

	payloadBytes, err = json.Marshal(payload)
	require.NoError(t, err)

	compareStateBytes, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, user1Address, payloadBytes, withdrawStateBytes, wasmBytes)
	require.NoError(t, err)
	require.NotNil(t, compareStateBytes)
	require.Len(t, events, 1)
	require.Len(t, withdrawals, 0)

	var eventData map[string]interface{}
	err = json.Unmarshal(events[0].Data, &eventData)
	require.NoError(t, err)
	t.Log("Event:\n", testutil.PrettyPrintJSON(eventData))

}

func TestSimpleAppIntegration_NullPayload(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module to get an initial state
	initialStateBytes, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	t.Run("null payload json", func(t *testing.T) {
		nullPayload := []byte{}
		_, _, _, err := runtime.ProcessRequest(ctx, appId, user1Address, nullPayload, initialStateBytes, wasmBytes)
		require.NoError(t, err)
	})
}

func TestSimpleAppIntegration_NegativeScenarios(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module to get an initial state
	initialStateBytes, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	// 2. Create a populated state for testing
	// User1 deposits 1000
	populatedStateBytes, _, err := runtime.Deposit(ctx, appId, user1Address, big.NewInt(1000), initialStateBytes, wasmBytes)
	require.NoError(t, err)
	// User2 deposits 500
	populatedStateBytes, _, err = runtime.Deposit(ctx, appId, user2Address, big.NewInt(500), populatedStateBytes, wasmBytes)
	require.NoError(t, err)

	// --- Test Cases ---

	t.Run("withdraw with insufficient balance", func(t *testing.T) {
		// User1 has 1000, tries to withdraw 2000
		payload := app.PayloadInstructions{
			Type: "withdraw",
			Withdraw: &app.WithdrawInstruction{
				To:     recipient1Address,
				Amount: big.NewInt(2000),
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, err = runtime.ProcessRequest(ctx, appId, user1Address, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Insufficient balance for withdrawal")
	})

	t.Run("withdraw from non-existent account", func(t *testing.T) {
		// A user that never deposited tries to withdraw
		nonExistentUser := ethCommon.HexToAddress("0xadd0000000000000000000000000000000000099")
		payload := app.PayloadInstructions{
			Type: "withdraw",
			Withdraw: &app.WithdrawInstruction{
				To:     recipient1Address,
				Amount: big.NewInt(100),
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, err = runtime.ProcessRequest(ctx, appId, nonExistentUser, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "does not exist")
	})

	t.Run("compare from non-existent account", func(t *testing.T) {
		nonExistentUser := ethCommon.HexToAddress("0xadd0000000000000000000000000000000000099")
		payload := app.PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &app.CompareInstructions{
				TargetAddress: user1Address, // user1 exists
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, err = runtime.ProcessRequest(ctx, appId, nonExistentUser, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Account " + nonExistentUser.Hex() + " does not exist!")
	})

	t.Run("withdraw with missing instruction", func(t *testing.T) {
		payload := app.PayloadInstructions{Type: "withdraw"} // Withdraw field is nil
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, err = runtime.ProcessRequest(ctx, appId, user1Address, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Withdraw instruction is missing")
	})

	t.Run("compare with missing instruction", func(t *testing.T) {
		payload := app.PayloadInstructions{Type: "compare_addresses"} // CompareAccounts field is nil
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, err = runtime.ProcessRequest(ctx, appId, user1Address, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Compare instruction is missing")
	})

	t.Run("unsupported instruction type", func(t *testing.T) {
		payload := `{"type":"invalid_instruction"}`
		_, _, _, err := runtime.ProcessRequest(ctx, appId, user1Address, []byte(payload), populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Unsupported instruction type")
	})

	t.Run("invalid payload json", func(t *testing.T) {
		payload := `{"type":"withdraw","withdraw":{"to":`
		_, _, _, err := runtime.ProcessRequest(ctx, appId, user1Address, []byte(payload), populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse payload instructions")
	})
}

func TestSimpleAppIntegration_NilData(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()

	// Load the module to get an initial state for some tests
	initialStateBytes, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	t.Run("deposit with nil state", func(t *testing.T) {
		// This should fail inside the wasm module because a nil state is not valid JSON.
		// This test verifies that the runtime correctly handles passing a nil slice to wasm.
		_, _, err := runtime.Deposit(ctx, appId, user1Address, big.NewInt(1000), nil, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})

	t.Run("generate report with nil payload", func(t *testing.T) {
		// The simple_app expects a JSON object for the payload. Passing nil results in an
		// empty string, which is invalid JSON, causing an error inside wasm.
		_, err := runtime.GenerateDeanonymizationReport(ctx, appId, nil, initialStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse payload")
	})

	t.Run("generate report with nil state", func(t *testing.T) {
		// This should fail inside the wasm module because a nil state is not valid JSON.
		_, err := runtime.GenerateDeanonymizationReport(ctx, appId, []byte("{}"), nil, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})
}

func TestSimpleAppIntegration_InvalidWasm(t *testing.T) {
	runtime := pes_wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	// A dummy state, as we expect failures before the state is even used.
	initialStateBytes := []byte(`{"appId":"simple_app_test","accounts":{}}`)

	t.Run("load module with nil wasm", func(t *testing.T) {
		_, err := runtime.LoadModule(ctx, appId, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to compile WASM module")
	})

	t.Run("load module with invalid wasm", func(t *testing.T) {
		_, err := runtime.LoadModule(ctx, appId, []byte("invalid wasm"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to compile WASM module")
	})

	// For other functions, the error will come from getOrLoadModule -> LoadModule
	t.Run("deposit with nil wasm", func(t *testing.T) {
		_, _, err := runtime.Deposit(ctx, appId, user1Address, big.NewInt(1000), initialStateBytes, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load module")
	})

	t.Run("process request with nil wasm", func(t *testing.T) {
		payload := []byte(`{"type":"withdraw","withdraw":{"to":"some_address","amount":100}}`)
		_, _, _, err := runtime.ProcessRequest(ctx, appId, user1Address, payload, initialStateBytes, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load module")
	})

	t.Run("generate report with nil wasm", func(t *testing.T) {
		_, err := runtime.GenerateDeanonymizationReport(ctx, appId, []byte("{}"), initialStateBytes, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load module")
	})
}

func TestSimpleAppIntegration_InvalidState(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	invalidState := []byte("{invalid-state}")

	t.Run("deposit with invalid state", func(t *testing.T) {
		_, _, err := runtime.Deposit(ctx, appId, user1Address, big.NewInt(1000), invalidState, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})

	t.Run("process request with invalid state", func(t *testing.T) {
		payload := []byte(`{"type":"withdraw","withdraw":{"to":"some_address","amount":100}}`)
		_, _, _, err := runtime.ProcessRequest(ctx, appId, user1Address, payload, invalidState, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})

	t.Run("generate report with invalid state", func(t *testing.T) {
		_, err := runtime.GenerateDeanonymizationReport(ctx, appId, []byte("{}"), invalidState, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})
}
