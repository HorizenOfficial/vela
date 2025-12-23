package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"sync"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/app/simple/app"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
	"github.com/horizen-pes/pkg/testutil"
	pes_wasm "github.com/horizen-pes/pkg/wasm"
	"github.com/stretchr/testify/require"
)

var testLogger logger.Logger

func TestMain(m *testing.M) {
	// Initialize once
	//	testLogger = logger.NewLogger(&logger.Config{Kind: "printf"})
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

var _ = []common.Withdrawal{}

const (
	wasmModulePath = "build/simple_app.wasm"
)

var (
	appId             = common.NewApplicationId(1)
	user1Address      = app.HexToAddress("0xadd0000000000000000000000000000000000001")
	user2Address      = app.HexToAddress("0xadd0000000000000000000000000000000000002")
	recipient1Address = app.HexToAddress("0xadd0000000000000000000000000000000000003")
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
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module
	initialStateBytes, fuel, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(5)))

	var initialState app.ApplicationInternalState
	err = json.Unmarshal(initialStateBytes, &initialState)
	require.NoError(t, err)
	require.Equal(t, appId, common.NewApplicationId(initialState.AppID))
	require.Empty(t, initialState.Accounts)

	// 2. User1 Deposits funds
	deposit1Amount := big.NewInt(1000)
	depositState1Bytes, depositEvents, fuel, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), deposit1Amount, initialStateBytes, wasmBytes)
	require.Nil(t, failure)
	require.NotNil(t, depositState1Bytes)
	require.Len(t, depositEvents, 1)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(35)))

	var depositState app.ApplicationInternalState
	err = json.Unmarshal(depositState1Bytes, &depositState)
	require.NoError(t, err)
	require.Len(t, depositState.Accounts, 1)
	require.Equal(t, deposit1Amount, depositState.Accounts[user1Address.Hex()].Balance)

	// 2. User2 Deposits funds (more than previous user)
	deposit2Amount := big.NewInt(2000)
	depositState2Bytes, depositEvents, fuel, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user2Address), deposit2Amount, depositState1Bytes, wasmBytes)
	require.Nil(t, failure)
	require.NotNil(t, depositState2Bytes)
	require.Len(t, depositEvents, 1)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(35)))

	err = json.Unmarshal(depositState2Bytes, &depositState)
	require.NoError(t, err)
	require.Len(t, depositState.Accounts, 2)
	require.Equal(t, deposit2Amount, depositState.Accounts[user2Address.Hex()].Balance)

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

	withdrawStateBytes, withdrawEvents, withdrawals, fuel, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), payloadBytes, depositState2Bytes, wasmBytes)
	require.Nil(t, err)
	require.NotNil(t, withdrawStateBytes)
	require.Len(t, withdrawEvents, 1)
	require.Len(t, withdrawals, 1)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(50)))

	var withdrawState app.ApplicationInternalState
	err = json.Unmarshal(withdrawStateBytes, &withdrawState)
	require.NoError(t, err)

	diffBalance := new(big.Int).Sub(deposit1Amount, withdrawAmount)
	require.Equal(t, diffBalance, withdrawState.Accounts[user1Address.Hex()].Balance)

	require.Equal(t, ethCommon.Address(recipient1Address), withdrawals[0].DestinationAddress)
	require.Equal(t, withdrawAmount, withdrawals[0].Amount)

	// 4. Generate deanonymization report
	payloadJSON := `{"tag":"my_custom_tag"}`
	payloadBytes = []byte(payloadJSON)
	reportBytes, fuel, failure := runtime.GenerateDeanonymizationReport(ctx, appId, payloadBytes, withdrawStateBytes, wasmBytes)
	require.Nil(t, failure)
	require.NotNil(t, reportBytes)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(20)))

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

	compareStateBytes, events, withdrawals, fuel, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), payloadBytes, withdrawStateBytes, wasmBytes)
	require.Nil(t, err)
	require.NotNil(t, compareStateBytes)
	require.Len(t, events, 1)
	require.Len(t, withdrawals, 0)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(50)))

	var eventData map[string]interface{}
	err = json.Unmarshal(events[0].Data, &eventData)
	require.NoError(t, err)
	t.Log("Event:\n", testutil.PrettyPrintJSON(eventData))

	// check we have no memory leaks
	mem_map_entries, total_allocated_bytes, err := runtime.GetAllocatedMemoryStats(ctx, appId, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, int64(0), mem_map_entries)
	require.Equal(t, int64(0), total_allocated_bytes)
	t.Logf("stats - memory map entries: %d, total bytes allocated: %d", mem_map_entries, total_allocated_bytes)

	// use an alternative implementation of the function
	mem_map_entries, total_allocated_bytes, err = runtime.GetAllocatedMemoryStats2(ctx, appId, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, int64(0), mem_map_entries)
	require.Equal(t, int64(0), total_allocated_bytes)
	t.Logf("stats2 - memory map entries: %d, total bytes allocated: %d", mem_map_entries, total_allocated_bytes)
}

func TestSimpleAppIntegration_NullPayload(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module to get an initial state
	initialStateBytes, fuel, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(5)))

	t.Run("null payload json", func(t *testing.T) {
		nullPayload := []byte{}
		_, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), nullPayload, initialStateBytes, wasmBytes)
		require.Nil(t, err)
	})
}

func TestSimpleAppIntegration_NegativeScenarios(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module to get an initial state
	initialStateBytes, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	// 2. Create a populated state for testing
	// User1 deposits 1000
	populatedStateBytes, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), big.NewInt(1000), initialStateBytes, wasmBytes)
	require.Nil(t, err)
	// User2 deposits 500
	populatedStateBytes, _, _, err = runtime.Deposit(ctx, appId, ethCommon.Address(user2Address), big.NewInt(500), populatedStateBytes, wasmBytes)
	require.Nil(t, err)

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

		_, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Insufficient balance for withdrawal")
	})

	t.Run("withdraw from non-existent account", func(t *testing.T) {
		// A user that never deposited tries to withdraw
		nonExistentUser := app.HexToAddress("0xadd0000000000000000000000000000000000099")
		payload := app.PayloadInstructions{
			Type: "withdraw",
			Withdraw: &app.WithdrawInstruction{
				To:     recipient1Address,
				Amount: big.NewInt(100),
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(nonExistentUser), payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "does not exist")
	})

	t.Run("compare from non-existent account", func(t *testing.T) {
		nonExistentUser := app.HexToAddress("0xadd0000000000000000000000000000000000099")
		payload := app.PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &app.CompareInstructions{
				TargetAddress: user1Address, // user1 exists
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(nonExistentUser), payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Account "+nonExistentUser.Hex()+" does not exist!")
	})

	t.Run("withdraw with missing instruction", func(t *testing.T) {
		payload := app.PayloadInstructions{Type: "withdraw"} // Withdraw field is nil
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Withdraw instruction is missing")
	})

	t.Run("compare with missing instruction", func(t *testing.T) {
		payload := app.PayloadInstructions{Type: "compare_addresses"} // CompareAccounts field is nil
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Compare instruction is missing")
	})

	t.Run("unsupported instruction type", func(t *testing.T) {
		payload := `{"type":"invalid_instruction"}`
		_, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), []byte(payload), populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Unsupported instruction type")
	})

	t.Run("invalid payload json", func(t *testing.T) {
		payload := `{"type":"withdraw","withdraw":{"to":`
		_, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), []byte(payload), populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse payload instructions")
	})
}

func TestSimpleAppIntegration_NilData(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()

	// Load the module to get an initial state for some tests
	initialStateBytes, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	t.Run("deposit with nil state", func(t *testing.T) {
		// This should fail inside the wasm module because a nil state is not valid JSON.
		// This test verifies that the runtime correctly handles passing a nil slice to wasm.
		_, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), big.NewInt(1000), nil, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})

	t.Run("generate report with nil payload", func(t *testing.T) {
		// The simple_app expects a JSON object for the payload. Passing nil results in an
		// empty string, which is invalid JSON, causing an error inside wasm.
		_, _, err := runtime.GenerateDeanonymizationReport(ctx, appId, nil, initialStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse payload")
	})

	t.Run("generate report with nil state", func(t *testing.T) {
		// This should fail inside the wasm module because a nil state is not valid JSON.
		_, _, err := runtime.GenerateDeanonymizationReport(ctx, appId, []byte("{}"), nil, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})
}

func TestSimpleAppIntegration_InvalidWasm(t *testing.T) {
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()
	// A dummy state, as we expect failures before the state is even used.
	initialStateBytes := []byte(`{"appId":"simple_app_test","accounts":{}}`)

	t.Run("load module with nil wasm", func(t *testing.T) {
		_, _, err := runtime.LoadModule(ctx, appId, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to compile WASM module")
	})

	t.Run("load module with invalid wasm", func(t *testing.T) {
		_, _, err := runtime.LoadModule(ctx, appId, []byte("invalid wasm"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to compile WASM module")
	})

	// For other functions, the error will come from getOrLoadModule -> LoadModule
	t.Run("deposit with nil wasm", func(t *testing.T) {
		_, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), big.NewInt(1000), initialStateBytes, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load module")
	})

	t.Run("process request with nil wasm", func(t *testing.T) {
		payload := []byte(`{"type":"withdraw","withdraw":{"to":"some_address","amount":100}}`)
		_, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), payload, initialStateBytes, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load module")
	})

	t.Run("generate report with nil wasm", func(t *testing.T) {
		_, _, err := runtime.GenerateDeanonymizationReport(ctx, appId, []byte("{}"), initialStateBytes, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load module")
	})
}

func TestSimpleAppIntegration_InvalidState(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()
	invalidState := []byte("{invalid-state}")

	t.Run("deposit with invalid state", func(t *testing.T) {
		_, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), big.NewInt(1000), invalidState, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})

	t.Run("process request with invalid state", func(t *testing.T) {
		payload := []byte(`{"type":"withdraw","withdraw":{"to":"some_address","amount":100}}`)
		_, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), payload, invalidState, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})

	t.Run("generate report with invalid state", func(t *testing.T) {
		_, _, err := runtime.GenerateDeanonymizationReport(ctx, appId, []byte("{}"), invalidState, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})
}

func TestSimpleAppIntegration_MemoryStress(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime with limited memory to make leaks surface faster.
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()

	initialStateBytes, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	// Mutex to protect access to the shared WASM runtime instance
	// Note this might not be enough, due to the stateful nature of the underlying Wasm instance and store which are being shared and reused across all goroutines.
	// See TODO at the end
	var runtimeMutex sync.Mutex

	// Repeatedly make calls and check mem is ok
	// TODO - we use 1 only routine until we will have a correct handling of wasm runtime concurrency (see TODO below), otherwise we might have sporadic failures
	const numGoroutines = 1
	const iterationsPerGoroutine = 40
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineIndex int) {
			defer wg.Done()
			stateBytes := initialStateBytes
			for j := 0; j < iterationsPerGoroutine; j++ {
				iterationIndex := goroutineIndex*iterationsPerGoroutine + j
				depositAmount := big.NewInt(1)
				userAddress := app.HexToAddress(fmt.Sprintf("0xadd%039d", iterationIndex))

				runtimeMutex.Lock()
				newStateBytes, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(userAddress), depositAmount, stateBytes, wasmBytes)
				require.Nil(t, err, "deposit failed at iteration %d", iterationIndex)
				stateBytes = newStateBytes
				runtimeMutex.Unlock()

				// Process a withdraw request for the current user
				withdrawAmount := big.NewInt(1)
				withdrawInstruction := app.WithdrawInstruction{
					To:     recipient1Address,
					Amount: withdrawAmount,
				}
				withdrawPayload := app.PayloadInstructions{
					Type:     "withdraw",
					Withdraw: &withdrawInstruction,
				}
				withdrawPayloadBytes, ret := json.Marshal(withdrawPayload)
				require.NoError(t, ret, "failed to marshal withdraw payload at iteration %d", iterationIndex)

				runtimeMutex.Lock()
				processStateBytes, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(userAddress), withdrawPayloadBytes, stateBytes, wasmBytes)
				require.Nil(t, err, "ProcessRequest failed at iteration %d", iterationIndex)
				stateBytes = processStateBytes
				runtimeMutex.Unlock()

				// Generate deanonymization report
				reportPayloadJSON := fmt.Sprintf(`{"tag":"memory_stress_report_%d"}`, iterationIndex)
				reportPayloadBytes := []byte(reportPayloadJSON)
				runtimeMutex.Lock()
				_, _, err = runtime.GenerateDeanonymizationReport(ctx, appId, reportPayloadBytes, stateBytes, wasmBytes)
				require.Nil(t, err, "GenerateDeanonymizationReport failed at iteration %d", iterationIndex)
				runtimeMutex.Unlock()
			}
		}(i)
	}

	wg.Wait()

	mem_map_entries, total_allocated_bytes, err := runtime.GetAllocatedMemoryStats2(ctx, appId, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, int64(0), mem_map_entries)
	require.Equal(t, int64(0), total_allocated_bytes)
	t.Logf("stats2 - memory map entries: %d, total bytes allocated: %d", mem_map_entries, total_allocated_bytes)

	// TODO -  The correct approach would be creating a new wasmtime.Instance for each concurrent operation.
	// While the wasmtime.Store and compiled wasmtime.Module can be shared, the wasmtime.Instance must be unique per goroutine.
	// We should:
	// 1. Separate the compiled module from the instance. The ApplicationModule should only contain the compiled wasmtime.Module,
	// not the wasmtime.Instance.
	// 2. Instantiate on demand. The Deposit, ProcessRequest, and other execution functions should get the shared wasmtime.Module
	// from the cache, and then create a brand new wasmtime.Instance for the duration of that specific function call
}
