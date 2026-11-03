package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"testing"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela-common-go/wasm/types"
	"github.com/HorizenOfficial/vela/app/simple/app"
	"github.com/HorizenOfficial/vela/app/subtype"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/testutil"
	vela_wasm "github.com/HorizenOfficial/vela/pkg/wasm"
	ethCommon "github.com/ethereum/go-ethereum/common"
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
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module
	initialStateBytes, fuel, err := runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(5)))

	var initialState app.ApplicationInternalState
	err = json.Unmarshal(initialStateBytes, &initialState)
	require.NoError(t, err)
	require.Equal(t, appId, common.NewApplicationId(initialState.AppID))
	require.Empty(t, initialState.Accounts)

	// 2. User1 Deposits funds
	deposit1Amount := big.NewInt(1000)
	depositState1Bytes, depositEvents, depositAppEvents, fuel, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), velacommon.ETH_TOKEN, deposit1Amount, initialStateBytes, wasmBytes)
	require.Nil(t, failure)
	require.NotNil(t, depositState1Bytes)
	require.Len(t, depositEvents, 1)
	require.Len(t, depositAppEvents, 1)
	require.Equal(t, subtype.FromString("deposit_received"), depositAppEvents[0].EventSubType)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(35)))

	var depositState app.ApplicationInternalState
	err = json.Unmarshal(depositState1Bytes, &depositState)
	require.NoError(t, err)
	require.Len(t, depositState.Accounts, 1)
	require.Equal(t, deposit1Amount.String(), depositState.Accounts[user1Address.Hex()].Balances["0x0000000000000000000000000000000000000000"].String())

	// 2. User2 Deposits funds (more than previous user)
	deposit2Amount := big.NewInt(2000)
	depositState2Bytes, depositEvents, _, fuel, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user2Address), velacommon.ETH_TOKEN, deposit2Amount, depositState1Bytes, wasmBytes)
	require.Nil(t, failure)
	require.NotNil(t, depositState2Bytes)
	require.Len(t, depositEvents, 1)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(35)))

	err = json.Unmarshal(depositState2Bytes, &depositState)
	require.NoError(t, err)
	require.Len(t, depositState.Accounts, 2)
	require.Equal(t, deposit2Amount.String(), depositState.Accounts[user2Address.Hex()].Balances["0x0000000000000000000000000000000000000000"].String())

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

	withdrawStateBytes, withdrawEvents, _, withdrawals, _, fuel, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, depositState2Bytes, wasmBytes)
	require.Nil(t, err)
	require.NotNil(t, withdrawStateBytes)
	require.Len(t, withdrawEvents, 1)
	require.Len(t, withdrawals, 1)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(50)))

	var withdrawState app.ApplicationInternalState
	err = json.Unmarshal(withdrawStateBytes, &withdrawState)
	require.NoError(t, err)

	diffBalance := new(big.Int).Sub(deposit1Amount, withdrawAmount)
	require.Equal(t, diffBalance.String(), withdrawState.Accounts[user1Address.Hex()].Balances["0x0000000000000000000000000000000000000000"].String())

	require.Equal(t, ethCommon.Address(recipient1Address), withdrawals[0].DestinationAddress)
	require.Equal(t, withdrawAmount, withdrawals[0].Amount.ToInt())

	// 4. Generate deanonymization report via ProcessRequest with type "deanonymize"
	deanonPayload := app.PayloadInstructions{
		Type:        "deanonymize",
		Deanonymize: &app.DeanonymizeInstruction{IncludeTag: "my_custom_tag"},
	}
	payloadBytes, err = json.Marshal(deanonPayload)
	require.NoError(t, err)

	deanonStateBytes, _, _, _, reportBytes, fuel, failure := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Deanonymize, payloadBytes, withdrawStateBytes, wasmBytes)
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

	compareStateBytes, events, _, withdrawals, _, fuel, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, withdrawStateBytes, wasmBytes)
	require.Nil(t, err)
	require.NotNil(t, compareStateBytes)
	require.Len(t, events, 1)
	require.Len(t, withdrawals, 0)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(50)))

	var eventData map[string]interface{}
	err = json.Unmarshal(events[0].Data, &eventData)
	require.NoError(t, err)
	t.Log("Event:\n", testutil.PrettyPrintJSON(eventData))
}

func TestSimpleAppIntegration_NullPayload(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module to get an initial state
	initialStateBytes, fuel, err := runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.NoError(t, err)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(5)))

	t.Run("null payload json", func(t *testing.T) {
		nullPayload := []byte{}
		_, _, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, nullPayload, initialStateBytes, wasmBytes)
		require.Nil(t, err)
	})
}

func TestSimpleAppIntegration_NegativeScenarios(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	ctx := context.Background()

	// 1. Load the module to get an initial state
	initialStateBytes, _, err := runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.NoError(t, err)

	// 2. Create a populated state for testing
	// User1 deposits 1000
	populatedStateBytes, _, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), velacommon.ETH_TOKEN, big.NewInt(1000), initialStateBytes, wasmBytes)
	require.Nil(t, err)
	// User2 deposits 500
	populatedStateBytes, _, _, _, err = runtime.Deposit(ctx, appId, ethCommon.Address(user2Address), velacommon.ETH_TOKEN, big.NewInt(500), populatedStateBytes, wasmBytes)
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

		_, _, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
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

		_, _, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(nonExistentUser), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
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

		_, _, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(nonExistentUser), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Account "+nonExistentUser.Hex()+" does not exist!")
	})

	t.Run("withdraw with missing instruction", func(t *testing.T) {
		payload := app.PayloadInstructions{Type: "withdraw"} // Withdraw field is nil
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Withdraw instruction is missing")
	})

	t.Run("compare with missing instruction", func(t *testing.T) {
		payload := app.PayloadInstructions{Type: "compare_addresses"} // CompareAccounts field is nil
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		_, _, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Compare instruction is missing")
	})

	t.Run("unsupported instruction type", func(t *testing.T) {
		payload := `{"type":"invalid_instruction"}`
		_, _, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, []byte(payload), populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Unsupported instruction type")
	})

	t.Run("invalid payload json", func(t *testing.T) {
		payload := `{"type":"withdraw","withdraw":{"to":`
		_, _, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, []byte(payload), populatedStateBytes, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse payload instructions")
	})
}

func TestSimpleAppIntegration_NilData(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	ctx := context.Background()

	// Load the module to verify it works (state used in later tests is created separately)
	_, _, err := runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.NoError(t, err)

	t.Run("deposit with nil state", func(t *testing.T) {
		// This should fail inside the wasm module because a nil state is not valid JSON.
		// This test verifies that the runtime correctly handles passing a nil slice to wasm.
		_, _, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), velacommon.ETH_TOKEN, big.NewInt(1000), nil, wasmBytes)
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

		_, _, _, _, _, _, err = runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Deanonymize, payloadBytes, nil, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})
}

func TestSimpleAppIntegration_InvalidWasm(t *testing.T) {
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	ctx := context.Background()
	// A dummy state, as we expect failures before the state is even used.
	initialStateBytes := []byte(`{"appId":"simple_app_test","accounts":{}}`)

	t.Run("deploy with nil wasm", func(t *testing.T) {
		_, _, err := runtime.Deploy(ctx, appId, nil, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to compile WASM module")
	})

	t.Run("deploy with invalid wasm", func(t *testing.T) {
		_, _, err := runtime.Deploy(ctx, appId, nil, []byte("invalid wasm"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to compile WASM module")
	})

	// For other functions, the error will come from getOrLoadModule -> compileAndInstantiate
	t.Run("deposit with nil wasm", func(t *testing.T) {
		_, _, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), velacommon.ETH_TOKEN, big.NewInt(1000), initialStateBytes, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load module")
	})

	t.Run("process request with nil wasm", func(t *testing.T) {
		payload := []byte(`{"type":"withdraw","withdraw":{"to":"some_address","amount":100}}`)
		_, _, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payload, initialStateBytes, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load module")
	})
}

func TestSimpleAppIntegration_InvalidState(t *testing.T) {
	// Build and load the wasm module
	wasmBytes := buildAndLoadWasmModule(t)

	// Create a new wasmtime runtime
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	ctx := context.Background()
	invalidState := []byte("{invalid-state}")

	t.Run("deposit with invalid state", func(t *testing.T) {
		_, _, _, _, err := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), velacommon.ETH_TOKEN, big.NewInt(1000), invalidState, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})

	t.Run("process request with invalid state", func(t *testing.T) {
		payload := []byte(`{"type":"withdraw","withdraw":{"to":"some_address","amount":100}}`)
		_, _, _, _, _, _, err := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payload, invalidState, wasmBytes)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Failed to parse application state")
	})
}

// TestSimpleAppIntegration_LargeResultRoundTrip exercises BytesToPtr with a large JSON payload.
// By depositing to many accounts and then generating a report, the guest serializes a large
// result that stresses the allocate -> BytesToPtr -> extractResultBytes -> deallocate pipeline.
func TestSimpleAppIntegration_LargeResultRoundTrip(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()

	ctx := context.Background()

	stateBytes, _, err := runtime.Deploy(ctx, appId, nil, wasmBytes)
	require.NoError(t, err)

	// Create 100 accounts with deposits to build a large state
	const numAccounts = 100
	for i := range numAccounts {
		addr, err := types.HexToAddress(fmt.Sprintf("0xadd%037d", i))
		require.NoError(t, err)

		newState, _, _, _, failure := runtime.Deposit(ctx, appId, ethCommon.Address(addr), velacommon.ETH_TOKEN, big.NewInt(int64(1000+i)), stateBytes, wasmBytes)
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

	_, _, _, _, reportBytes, _, failure := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Deanonymize, payloadBytes, stateBytes, wasmBytes)
	// Generate report: SerializeAndWriteResult must handle the large report via BytesToPtr
	require.Nil(t, failure)
	require.NotNil(t, reportBytes)

	// Verify report contains all accounts
	var report app.DeanonymizationReport
	err = json.Unmarshal(reportBytes, &report)
	require.NoError(t, err)
	require.Equal(t, "my_custom_tag", report.Tag)
	require.Len(t, report.Accounts, numAccounts)
}

// --- ERC-20 integration tests ---
// These tests verify that token-aware operations cross the host→guest WASM memory boundary correctly.

var (
	erc20TokenAddress, _ = types.HexToAddress("0xdead000000000000000000000000000000000001")
	ethZeroAddress       = types.Address{}
	ethZeroAddressHex    = ethZeroAddress.Hex()
)

// TestDeployWithConstructorParams verifies the deploy export receives constructor params
// through WASM memory and initializes the app state with a token allowlist.
func TestDeployWithConstructorParams(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	ctx := context.Background()

	// Deploy with an ERC-20 token in the allowlist
	params := fmt.Sprintf(`{"allowedTokens":["%s"]}`, erc20TokenAddress.Hex())
	stateBytes, fuel, err := runtime.Deploy(ctx, appId, []byte(params), wasmBytes)
	require.NoError(t, err)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(5)))

	var state app.ApplicationInternalState
	err = json.Unmarshal(stateBytes, &state)
	require.NoError(t, err)
	require.Equal(t, appId, common.NewApplicationId(state.AppID))
	require.Empty(t, state.Accounts)
	// Both ETH and the ERC-20 token should be in the allowlist
	require.True(t, state.AllowedTokens[ethZeroAddressHex], "ETH should always be allowed")
	require.True(t, state.AllowedTokens[erc20TokenAddress.Hex()], "ERC-20 token should be allowed")
}

// TestDeployWithEmptyParams verifies deploy works with no constructor params (ETH-only app).
func TestDeployWithEmptyParams(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	ctx := context.Background()

	stateBytes, fuel, err := runtime.Deploy(ctx, appId, []byte(""), wasmBytes)
	require.NoError(t, err)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(5)))

	var state app.ApplicationInternalState
	err = json.Unmarshal(stateBytes, &state)
	require.NoError(t, err)
	require.True(t, state.AllowedTokens[ethZeroAddressHex], "ETH should always be allowed")
	require.Len(t, state.AllowedTokens, 1, "only ETH should be allowed")
}

// TestERC20DepositThroughWasm verifies that the tokenAddress bytes cross the
// host→guest WASM memory boundary correctly and the per-token balance is credited.
func TestERC20DepositThroughWasm(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	ctx := context.Background()

	// Deploy with ERC-20 token allowed
	params := fmt.Sprintf(`{"allowedTokens":["%s"]}`, erc20TokenAddress.Hex())
	stateBytes, _, err := runtime.Deploy(ctx, appId, []byte(params), wasmBytes)
	require.NoError(t, err)

	// Deposit ERC-20 token
	depositAmount := big.NewInt(500)
	newStateBytes, events, _, fuel, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), ethCommon.Address(erc20TokenAddress), depositAmount, stateBytes, wasmBytes)
	require.Nil(t, failure)
	require.NotNil(t, newStateBytes)
	require.Len(t, events, 1)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(35)))

	var state app.ApplicationInternalState
	err = json.Unmarshal(newStateBytes, &state)
	require.NoError(t, err)

	// ERC-20 balance should be credited
	acct := state.Accounts[user1Address.Hex()]
	require.NotNil(t, acct)
	require.Equal(t, depositAmount.String(), acct.Balances[erc20TokenAddress.Hex()].String())
	// ETH balance should not exist
	require.Nil(t, acct.Balances[ethZeroAddressHex])
}

// TestERC20DepositRejectedForNonAllowlistedToken verifies that depositing a token
// not in the app's allowlist returns an error through the WASM boundary.
func TestERC20DepositRejectedForNonAllowlistedToken(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	ctx := context.Background()

	// Deploy ETH-only (no ERC-20 in allowlist)
	stateBytes, _, err := runtime.Deploy(ctx, appId, []byte(""), wasmBytes)
	require.NoError(t, err)

	// Attempt to deposit an ERC-20 token
	depositAmount := big.NewInt(100)
	_, _, _, _, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), ethCommon.Address(erc20TokenAddress), depositAmount, stateBytes, wasmBytes)
	require.NotNil(t, failure, "deposit of non-allowlisted token should fail")
}

// TestERC20WithdrawalThroughWasm verifies token-aware withdrawals survive
// the guest→host JSON serialization and carry the correct TokenAddress.
func TestERC20WithdrawalThroughWasm(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	ctx := context.Background()

	// Deploy with ERC-20 token allowed
	params := fmt.Sprintf(`{"allowedTokens":["%s"]}`, erc20TokenAddress.Hex())
	stateBytes, _, err := runtime.Deploy(ctx, appId, []byte(params), wasmBytes)
	require.NoError(t, err)

	// Deposit ERC-20 token
	depositAmount := big.NewInt(1000)
	stateBytes, _, _, _, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), ethCommon.Address(erc20TokenAddress), depositAmount, stateBytes, wasmBytes)
	require.Nil(t, failure)

	// Withdraw 300 of ERC-20 token to recipient
	withdrawAmount := big.NewInt(300)
	withdrawInstruction := app.WithdrawInstruction{
		To:           recipient1Address,
		TokenAddress: erc20TokenAddress,
		Amount:       new(types.Uint256).SetBytes(withdrawAmount.Bytes()),
	}
	payload := app.PayloadInstructions{
		Type:     "withdraw",
		Withdraw: &withdrawInstruction,
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	newStateBytes, events, _, withdrawals, _, fuel, processFailure := runtime.ProcessRequest(ctx, appId, ethCommon.Address(user1Address), common.Process, payloadBytes, stateBytes, wasmBytes)
	require.Nil(t, processFailure)
	require.NotNil(t, newStateBytes)
	require.Len(t, events, 1)
	require.Len(t, withdrawals, 1)
	require.Equal(t, 0, fuel.Cmp(big.NewInt(50)))

	// Verify withdrawal carries the correct token address
	require.Equal(t, ethCommon.Address(erc20TokenAddress), withdrawals[0].TokenAddress)
	require.Equal(t, ethCommon.Address(recipient1Address), withdrawals[0].DestinationAddress)
	require.Equal(t, withdrawAmount, withdrawals[0].Amount.ToInt())

	// Verify ERC-20 balance was debited
	var state app.ApplicationInternalState
	err = json.Unmarshal(newStateBytes, &state)
	require.NoError(t, err)
	expectedBalance := new(big.Int).Sub(depositAmount, withdrawAmount)
	require.Equal(t, expectedBalance.String(), state.Accounts[user1Address.Hex()].Balances[erc20TokenAddress.Hex()].String())
}

// TestMixedETHAndERC20Deposits verifies that ETH and ERC-20 deposits are tracked independently
// when both cross the WASM boundary.
func TestMixedETHAndERC20Deposits(t *testing.T) {
	wasmBytes := buildAndLoadWasmModule(t)
	runtime := vela_wasm.NewWasmtimeRuntime(testLogger, 0)
	defer runtime.Close()
	ctx := context.Background()

	// Deploy with ERC-20 token allowed
	params := fmt.Sprintf(`{"allowedTokens":["%s"]}`, erc20TokenAddress.Hex())
	stateBytes, _, err := runtime.Deploy(ctx, appId, []byte(params), wasmBytes)
	require.NoError(t, err)

	// Deposit ETH
	ethAmount := big.NewInt(1000)
	stateBytes, _, _, _, failure := runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), velacommon.ETH_TOKEN, ethAmount, stateBytes, wasmBytes)
	require.Nil(t, failure)

	// Deposit ERC-20
	erc20Amount := big.NewInt(500)
	stateBytes, _, _, _, failure = runtime.Deposit(ctx, appId, ethCommon.Address(user1Address), ethCommon.Address(erc20TokenAddress), erc20Amount, stateBytes, wasmBytes)
	require.Nil(t, failure)

	// Verify both balances are independent
	var state app.ApplicationInternalState
	err = json.Unmarshal(stateBytes, &state)
	require.NoError(t, err)

	acct := state.Accounts[user1Address.Hex()]
	require.NotNil(t, acct)
	require.Equal(t, ethAmount.String(), acct.Balances[ethZeroAddressHex].String())
	require.Equal(t, erc20Amount.String(), acct.Balances[erc20TokenAddress.Hex()].String())
}
