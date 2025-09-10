package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"payment-app/app"

	"github.com/horizen-pes/pkg/wasm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readWasm(t *testing.T) []byte {
	t.Helper()
	wasmPath := filepath.Join("payment_app.wasm")
	wasmBytes, err := os.ReadFile(wasmPath)
	require.NoError(t, err, "Failed to read WASM file")
	return wasmBytes
}

func TestIntegration_LoadModule(t *testing.T) {
	wasmBytes := readWasm(t)
	runtime := wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"

	state, stateRoot, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)
	require.NotNil(t, state)
	require.NotNil(t, stateRoot)

	var stateData app.ApplicationInternalState
	require.NoError(t, json.Unmarshal(state, &stateData))
	assert.Equal(t, appId, stateData.AppID)
}

func TestIntegration_Deposit(t *testing.T) {
	wasmBytes := readWasm(t)
	runtime := wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	sender := fmt.Sprintf("0xadd%037x", 1)
	value := uint64(1_000_000_000_000_000_000)

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)

	newState, events, err := runtime.Deposit(ctx, appId, sender, value, state, wasmBytes)
	require.NoError(t, err)
	require.Len(t, events, 1)

	var stateData app.ApplicationInternalState
	require.NoError(t, json.Unmarshal(newState, &stateData))
	require.Contains(t, stateData.Accounts, sender)
	assert.Equal(t, value, stateData.Accounts[sender].Balance)
}

func TestIntegration_ProcessRequest_Transfer(t *testing.T) {
	wasmBytes := readWasm(t)
	runtime := wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	sender := fmt.Sprintf("0xadd%037x", 1)
	recipient := fmt.Sprintf("0xadd%037x", 2)
	depositValue := uint64(2_000_000_000_000_000_000)
	transferValue := uint64(500_000_000_000_000_000)

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)
	state, _, err = runtime.Deposit(ctx, appId, sender, depositValue, state, wasmBytes)
	require.NoError(t, err)

	payload := app.PayloadInstructions{
		Type:     "transfer",
		Transfer: &app.TransferInstruction{To: recipient, Amount: transferValue},
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	newState, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, sender, payloadBytes, state, wasmBytes)
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.Len(t, withdrawals, 0)

	var stateData app.ApplicationInternalState
	require.NoError(t, json.Unmarshal(newState, &stateData))
	assert.Equal(t, depositValue-transferValue, stateData.Accounts[sender].Balance)
	assert.Equal(t, transferValue, stateData.Accounts[recipient].Balance)
}

func TestIntegration_ProcessRequest_Withdrawal(t *testing.T) {
	wasmBytes := readWasm(t)
	runtime := wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	ctx := context.Background()
	appId := "test-app"
	sender := fmt.Sprintf("0xadd%037x", 1)
	depositValue := uint64(1_000_000_000_000_000_000)
	withdrawValue := uint64(500_000_000_000_000_000)
	withdrawAddress := "0x1234567890123456789012345678901234567890"

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)
	state, _, err = runtime.Deposit(ctx, appId, sender, depositValue, state, wasmBytes)
	require.NoError(t, err)

	payload := app.PayloadInstructions{
		Type:     "withdraw",
		Withdraw: &app.WithdrawInstruction{To: withdrawAddress, Amount: withdrawValue},
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	newState, events, withdrawals, err := runtime.ProcessRequest(ctx, appId, sender, payloadBytes, state, wasmBytes)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Len(t, withdrawals, 1)
	assert.Equal(t, withdrawAddress, withdrawals[0].DestinationAddress)
	assert.Equal(t, withdrawValue, withdrawals[0].Amount)

	var stateData app.ApplicationInternalState
	require.NoError(t, json.Unmarshal(newState, &stateData))
	assert.Equal(t, depositValue-withdrawValue, stateData.Accounts[sender].Balance)
}

func TestIntegration_GenerateDeanonymizationReport(t *testing.T) {
	wasmBytes := readWasm(t)
	runtime := wasm.NewWasmtimeRuntime()
	defer runtime.Close()

	type reportStruct struct {
		ApplicationID string                       `json:"applicationId"`
		RequestID     string                       `json:"requestId"`
		Accounts      map[string]*app.AccountState `json:"accounts"`
		Nonce         uint64                       `json:"nonce"`
	}

	ctx := context.Background()
	appId := "test-app"
	requestId := "deanon-1"
	sender := fmt.Sprintf("0xadd%037x", 1)
	value := uint64(1_000_000_000_000_000_000)

	state, _, err := runtime.LoadModule(ctx, appId, wasmBytes)
	require.NoError(t, err)
	state, _, err = runtime.Deposit(ctx, appId, sender, value, state, wasmBytes)
	require.NoError(t, err)

	reportBytes, err := runtime.GenerateDeanonymizationReport(ctx, appId, requestId, []byte{}, state, wasmBytes)
	require.NoError(t, err)
	require.NotNil(t, reportBytes)

	var report reportStruct
	require.NoError(t, json.Unmarshal(reportBytes, &report))
	assert.Equal(t, appId, report.ApplicationID)
	assert.Equal(t, requestId, report.RequestID)
	require.Contains(t, report.Accounts, sender)
	assert.Equal(t, value, report.Accounts[sender].Balance)
}
