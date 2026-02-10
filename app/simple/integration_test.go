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
	"github.com/horizen-cce-common-go/wasm/types"
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
	appId                = common.NewApplicationId(1)
	user1Address, _      = types.HexToAddress("0xadd0000000000000000000000000000000000001")
	user2Address, _      = types.HexToAddress("0xadd0000000000000000000000000000000000002")
	recipient1Address, _ = types.HexToAddress("0xadd0000000000000000000000000000000000003")
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
	require.Equal(t, deposit1Amount.String(), depositState.Accounts[user1Address.Hex()].Balance.String())

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
	require.Equal(t, deposit2Amount.String(), depositState.Accounts[user2Address.Hex()].Balance.String())

	// 3. Process a withdraw request for user1
	withdrawAmount := big.NewInt(200)
	withdrawInstruction := app.WithdrawInstruction{
		To:     recipient1Address,
		Amount: new(types.Uint256).SetBytes(withdrawAmount.Bytes()),
	}
	payload := app.PayloadInstructions{
		Type:     "withdraw",
		Withdraw: &withdrawInstruction,
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	withdrawStateBytes, withdrawEvents, withdrawals, _, fuel, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, depositState2Bytes, wasmBytes)
	require.Nil(t, err)
	require.NotNil(t, withdrawStateBytes)
	require.Len(t, withdrawEvents, 1)
	require.Len(t, withdrawals, 1)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(50)))

	var withdrawState app.ApplicationInternalState
	err = json.Unmarshal(withdrawStateBytes, &withdrawState)
	require.NoError(t, err)

	diffBalance := new(big.Int).Sub(deposit1Amount, withdrawAmount)
	require.Equal(t, diffBalance.String(), withdrawState.Accounts[user1Address.Hex()].Balance.String())

	require.Equal(t, ethCommon.Address(recipient1Address), withdrawals[0].DestinationAddress)
	require.Equal(t, withdrawAmount, withdrawals[0].Amount.ToInt())

	// 4. Generate deanonymization report via ProcessRequest with type "deanonymize"
	deanonPayload := app.PayloadInstructions{
		Type:        "deanonymize",
		Deanonymize: &app.DeanonymizeInstruction{IncludeTag: "my_custom_tag"},
	}
	payloadBytes, err = json.Marshal(deanonPayload)
	require.NoError(t, err)

	deanonStateBytes, _, _, reportBytes, fuel, failure := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Deanonymize, payloadBytes, withdrawStateBytes, wasmBytes)
	require.Nil(t, failure)
	require.NotNil(t, deanonStateBytes)
	require.NotNil(t, reportBytes)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(20)))

	var report map[string]interface{}
	err = json.Unmarshal(reportBytes, &report)
	require.NoError(t, err)
	require.Equal(t, "my_custom_tag", report["tag"])
	t.Log(" Report:\n", testutil.PrettyPrintJSON(report))
	withdrawStateBytes = deanonStateBytes // Update state for next test

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

	compareStateBytes, events, withdrawals, _, fuel, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, withdrawStateBytes, wasmBytes)
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
		_, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, nullPayload, initialStateBytes, wasmBytes)
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
				Amount: types.NewUint256(2000),
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Insufficient balance")
	})

	t.Run("withdraw from non-existent account", func(t *testing.T) {
		// A user that never deposited tries to withdraw
		nonExistentUser, _ := types.HexToAddress("0xadd0000000000000000000000000000000000099")
		payload := app.PayloadInstructions{
			Type: "withdraw",
			Withdraw: &app.WithdrawInstruction{
				To:     recipient1Address,
				Amount: types.NewUint256(100),
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(nonExistentUser), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "does not exist")
	})

	t.Run("compare from non-existent account", func(t *testing.T) {
		nonExistentUser, _ := types.HexToAddress("0xadd0000000000000000000000000000000000099")
		payload := app.PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &app.CompareInstructions{
				TargetAddress: user1Address, // user1 exists
			},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(nonExistentUser), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Account "+nonExistentUser.Hex()+" does not exist!")
	})

	t.Run("withdraw with missing instruction", func(t *testing.T) {
		payload := app.PayloadInstructions{Type: "withdraw"} // Withdraw field is nil
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Withdraw instruction is missing")
	})

	t.Run("compare with missing instruction", func(t *testing.T) {
		payload := app.PayloadInstructions{Type: "compare_addresses"} // CompareAccounts field is nil
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Compare instruction is missing")
	})

	t.Run("unsupported instruction type", func(t *testing.T) {
		payload := `{"type":"invalid_instruction"}`
		_, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, []byte(payload), populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Unsupported instruction type")
	})

	t.Run("invalid payload json", func(t *testing.T) {
		payload := `{"type":"withdraw","withdraw":{"to":`
		_, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, []byte(payload), populatedStateBytes, wasmBytes)
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

	// Load the module to verify it works (state used in later tests is created separately)
	_, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	t.Run("deposit with nil state", func(t *testing.T) {
		// This should fail inside the wasm module because a nil state is not valid JSON.
		// This test verifies that the runtime correctly handles passing a nil slice to wasm.
		_, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), big.NewInt(1000), nil, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})

	t.Run("deanonymize with nil state", func(t *testing.T) {
		// The deanonymize instruction via ProcessRequest should fail with nil state
		deanonPayload := app.PayloadInstructions{
			Type: "deanonymize",
		}
		payloadBytes, err := json.Marshal(deanonPayload)
		require.NoError(t, err)

		_, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Deanonymize, payloadBytes, nil, wasmBytes)
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
		_, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payload, initialStateBytes, nil)
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
		_, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payload, invalidState, wasmBytes)
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

	for i := range numGoroutines {
		go func(goroutineIndex int) {
			defer wg.Done()
			stateBytes := initialStateBytes
			for j := 0; j < iterationsPerGoroutine; j++ {
				iterationIndex := goroutineIndex*iterationsPerGoroutine + j
				depositAmount := big.NewInt(1)
				userAddress, err := types.HexToAddress(fmt.Sprintf("0xadd%037d", iterationIndex))
				require.NoError(t, err)

				runtimeMutex.Lock()
				newStateBytes, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(userAddress), depositAmount, stateBytes, wasmBytes)
				require.Nil(t, err, "deposit failed at iteration %d", iterationIndex)
				stateBytes = newStateBytes
				runtimeMutex.Unlock()

				// Process a withdraw request for the current user
				withdrawAmount := big.NewInt(1)
				withdrawInstruction := app.WithdrawInstruction{
					To:     recipient1Address,
					Amount: new(types.Uint256).SetBytes(withdrawAmount.Bytes()),
				}
				withdrawPayload := app.PayloadInstructions{
					Type:     "withdraw",
					Withdraw: &withdrawInstruction,
				}
				withdrawPayloadBytes, ret := json.Marshal(withdrawPayload)
				require.NoError(t, ret, "failed to marshal withdraw payload at iteration %d", iterationIndex)

				runtimeMutex.Lock()
				processStateBytes, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(userAddress), common.Process, withdrawPayloadBytes, stateBytes, wasmBytes)
				require.Nil(t, err, "ProcessRequest failed at iteration %d", iterationIndex)
				stateBytes = processStateBytes
				runtimeMutex.Unlock()

				// Generate deanonymization report via ProcessRequest with type "deanonymize"
				deanonPayload := app.PayloadInstructions{
					Type:        "deanonymize",
					Deanonymize: &app.DeanonymizeInstruction{IncludeTag: fmt.Sprintf("memory_stress_report_%d", iterationIndex)},
				}
				deanonPayloadBytes, ret := json.Marshal(deanonPayload)
				require.NoError(t, ret, "failed to marshal deanonymize payload at iteration %d", iterationIndex)
				runtimeMutex.Lock()
				deanonStateBytes, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(userAddress), common.Deanonymize, deanonPayloadBytes, stateBytes, wasmBytes)
				require.Nil(t, err, "Deanonymize ProcessRequest failed at iteration %d", iterationIndex)
				stateBytes = deanonStateBytes
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

// requireMemoryClean is a helper that checks guest memory is fully deallocated.
func requireMemoryClean(t *testing.T, runtime *pes_wasm.WasmtimeRuntime, wasmBytes []byte, errMsg string) {
	t.Helper()
	ctx := context.Background()
	mapEntries, totalBytes, err := runtime.GetAllocatedMemoryStats2(ctx, appId, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, int64(0), mapEntries, errMsg)
	require.Equal(t, int64(0), totalBytes, errMsg)
}

// TestSimpleAppIntegration_MemoryCleanBetweenOps verifies that BytesToPtr allocations
// (created by SerializeAndWriteResult inside the WASM guest) are fully deallocated after
// each host call. This exercises the allocate -> BytesToPtr -> extractResultBytes -> deallocate
// round-trip that cannot be unit-tested in native 64-bit Go.
func TestSimpleAppIntegration_MemoryCleanBetweenOps(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()

	// LoadModule: exercises SerializeAndWriteResult for LoadModuleResult
	stateBytes, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after LoadModule")

	// Deposit: exercises SerializeAndWriteResult for DepositResult (with events)
	stateBytes, _, _, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), big.NewInt(5000), stateBytes, wasmBytes)
	require.Nil(t, failure)
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after Deposit")

	// ProcessRequest (withdraw): exercises SerializeAndWriteResult for ProcessResult (with events + withdrawals)
	payload := app.PayloadInstructions{
		Type: "withdraw",
		Withdraw: &app.WithdrawInstruction{
			To:     recipient1Address,
			Amount: types.NewUint256(100),
		},
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	stateBytes, _, _, _, _, failure2 := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, stateBytes, wasmBytes)
	require.Nil(t, failure2)
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after ProcessRequest (withdraw)")

	// Deposit a second user so we can compare
	stateBytes, _, _, failure = runtime.Deposit(ctx, appId, ethCommon.Address(user2Address), big.NewInt(3000), stateBytes, wasmBytes)
	require.Nil(t, failure)
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after second Deposit")

	// ProcessRequest (compare): different result shape
	payload2 := app.PayloadInstructions{
		Type: "compare_addresses",
		CompareAccounts: &app.CompareInstructions{
			TargetAddress: user2Address,
		},
	}
	payloadBytes2, err := json.Marshal(payload2)
	require.NoError(t, err)

	stateBytes, _, _, _, _, failure2 = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes2, stateBytes, wasmBytes)
	require.Nil(t, failure2)
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after ProcessRequest (compare)")

	// 4. Generate deanonymization report via ProcessRequest with type "deanonymize"
	deanonPayload := app.PayloadInstructions{
		Type:        "deanonymize",
		Deanonymize: &app.DeanonymizeInstruction{IncludeTag: "my_custom_tag"},
	}
	payloadBytes, err = json.Marshal(deanonPayload)
	require.NoError(t, err)

	_, _, _, _, _, failure = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Deanonymize, payloadBytes, stateBytes, wasmBytes)

	require.Nil(t, failure)
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after GenerateDeanonymizationReport")
}

// TestSimpleAppIntegration_ErrorPathMemory verifies that error results returned by the guest
// (which still use SerializeAndWriteResult -> BytesToPtr) do not leak memory.
func TestSimpleAppIntegration_ErrorPathMemory(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()

	stateBytes, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	// Deposit so user1 has a balance
	stateBytes, _, _, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), big.NewInt(100), stateBytes, wasmBytes)
	require.Nil(t, failure)
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after initial deposit")

	// Trigger error: withdraw more than balance.
	// The guest returns an error result via SerializeAndWriteResult -> BytesToPtr.
	payload := app.PayloadInstructions{
		Type: "withdraw",
		Withdraw: &app.WithdrawInstruction{
			To:     recipient1Address,
			Amount: types.NewUint256(9999),
		},
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	_, _, _, _, _, failure2 := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, stateBytes, wasmBytes)
	require.NotNil(t, failure2, "expected error for insufficient balance")
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after error result from ProcessRequest")

	// Trigger error: non-existent account
	_, _, _, _, _, failure2 = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user2Address), common.Process, payloadBytes, stateBytes, wasmBytes)
	require.NotNil(t, failure2, "expected error for non-existent account")
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after error result for non-existent account")

	// Trigger error: invalid state
	_, _, _, failure = runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), big.NewInt(100), []byte("{bad-json}"), wasmBytes)
	require.NotNil(t, failure, "expected error for invalid state")
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after error result for invalid state")
}

// TestSimpleAppIntegration_LargeResultRoundTrip exercises BytesToPtr with a large JSON payload.
// By depositing to many accounts and then generating a report, the guest serializes a large
// result that stresses the allocate -> BytesToPtr -> extractResultBytes -> deallocate pipeline.
func TestSimpleAppIntegration_LargeResultRoundTrip(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := pes_wasm.NewWasmtimeRuntime(testLogger)
	defer runtime.Close()

	ctx := context.Background()

	stateBytes, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	// Create 100 accounts with deposits to build a large state
	const numAccounts = 100
	for i := range numAccounts {
		addr, err := types.HexToAddress(fmt.Sprintf("0xadd%037d", i))
		require.NoError(t, err)

		newState, _, _, failure := runtime.Deposit(ctx, appId, ethCommon.Address(addr), big.NewInt(int64(1000+i)), stateBytes, wasmBytes)
		require.Nil(t, failure, "deposit failed for account %d", i)
		stateBytes = newState
	}

	// Verify large state can be deserialized
	var state app.ApplicationInternalState
	err = json.Unmarshal(stateBytes, &state)
	require.NoError(t, err)
	require.Len(t, state.Accounts, numAccounts)

	// 4. Generate deanonymization report via ProcessRequest with type "deanonymize"
	deanonPayload := app.PayloadInstructions{
		Type:        "deanonymize",
		Deanonymize: &app.DeanonymizeInstruction{IncludeTag: "my_custom_tag"},
	}
	payloadBytes, err := json.Marshal(deanonPayload)
	require.NoError(t, err)

	_, _, _, reportBytes, _, failure := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Deanonymize, payloadBytes, stateBytes, wasmBytes)
	// Generate report: SerializeAndWriteResult must handle the large report via BytesToPtr
	require.Nil(t, failure)
	require.NotNil(t, reportBytes)

	// Verify report contains all accounts
	var report app.DeanonymizationReport
	err = json.Unmarshal(reportBytes, &report)
	require.NoError(t, err)
	require.Equal(t, "my_custom_tag", report.Tag)
	require.Len(t, report.Accounts, numAccounts)

	// Verify no memory leaked despite the large allocation
	requireMemoryClean(t, runtime, wasmBytes, "memory leak after large result round-trip")
}
