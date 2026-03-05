package app

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/HorizenOfficial/vela-common-go/wasm/types"
	"github.com/HorizenOfficial/vela/pkg/common"
	wasmCommon "github.com/HorizenOfficial/vela/pkg/wasm/common"
	"github.com/stretchr/testify/require"
)

// host-side event types for JSON compatibility tests (app-specific, not framework types)
type hostDepositEvent struct {
	Type    string      `json:"type"`
	Amount  *common.Big `json:"amount"`
	Balance *common.Big `json:"balance"`
	Nonce   uint64      `json:"nonce"`
}

type hostSenderEvent struct {
	Type    string            `json:"type"`
	To      ethCommon.Address `json:"to"`
	Amount  *common.Big       `json:"amount"`
	Balance *common.Big       `json:"balance"`
	Nonce   uint64            `json:"nonce"`
}

type hostRecipientEvent struct {
	Type    string            `json:"type"`
	From    ethCommon.Address `json:"from"`
	Amount  *common.Big       `json:"amount"`
	Balance *common.Big       `json:"balance"`
	Nonce   uint64            `json:"nonce"`
}

const (
	testAppId = uint64(1)
)

var (
	user1Address, _ = types.HexToAddress("0xadd0000000000000000000000000000000000001")
	user2Address, _ = types.HexToAddress("0xadd0000000000000000000000000000000000002")
	user3Address, _ = types.HexToAddress("0xadd0000000000000000000000000000000000003")
)

func getInitialState(t *testing.T) (string, ApplicationInternalState) {
	t.Helper()
	initialState := ApplicationInternalState{
		AppID:    testAppId,
		Accounts: make(map[string]*AccountState),
	}
	stateBytes, err := json.Marshal(initialState)
	require.NoError(t, err)
	return string(stateBytes), initialState
}

func getPopulatedState(t *testing.T) (string, ApplicationInternalState) {
	t.Helper()
	state := ApplicationInternalState{
		AppID: testAppId,
		Accounts: map[string]*AccountState{
			user1Address.Hex(): {Address: user1Address, Balance: types.NewUint256(1000)},
			user2Address.Hex(): {Address: user2Address, Balance: types.NewUint256(500)},
		},
	}
	stateBytes, err := json.Marshal(state)
	require.NoError(t, err)
	return string(stateBytes), state
}

func TestLoadModule(t *testing.T) {
	result := LoadModule(int64(testAppId))
	require.NotNil(t, result.State)
	require.NotNil(t, result.Fuel)

	var state ApplicationInternalState
	err := json.Unmarshal(result.State, &state)
	require.NoError(t, err)

	require.Equal(t, testAppId, state.AppID)
	require.Empty(t, state.Accounts)
}

func TestDepositFunds(t *testing.T) {
	t.Run("deposit to new account", func(t *testing.T) {
		stateJSON, _ := getInitialState(t)
		depositAmount := types.NewUint256(100)

		result := DepositFunds(&user1Address, depositAmount, stateJSON)
		require.Empty(t, result.Error)
		require.NotNil(t, result.State)
		require.Len(t, result.Events, 1)

		var newState ApplicationInternalState
		err := json.Unmarshal(result.State, &newState)
		require.NoError(t, err)

		require.Len(t, newState.Accounts, 1)
		require.Equal(t, depositAmount, newState.Accounts[user1Address.Hex()].Balance)
		require.Equal(t, user1Address, newState.Accounts[user1Address.Hex()].Address)

		event := result.Events[0]
		require.Equal(t, user1Address, event.UserID)
		var eventData struct {
			Type   string         `json:"type"`
			Amount *types.Uint256 `json:"amount"`
		}
		err = json.Unmarshal(event.Data, &eventData)
		require.NoError(t, err)
		require.Equal(t, "deposit", eventData.Type)

		require.Equal(t, depositAmount, eventData.Amount)
	})

	t.Run("deposit to existing account", func(t *testing.T) {
		_, state := getInitialState(t)
		initialBalance := types.NewUint256(50)
		state.Accounts[user1Address.Hex()] = &AccountState{Address: user1Address, Balance: initialBalance}

		stateBytes, err := json.Marshal(state)
		require.NoError(t, err)
		stateJSON := string(stateBytes)

		depositAmount := types.NewUint256(100)
		result := DepositFunds(&user1Address, depositAmount, stateJSON)
		require.Empty(t, result.Error)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)

		sum := types.NewUint256(0).Add(*initialBalance, *depositAmount)
		require.Equal(t, sum, newState.Accounts[user1Address.Hex()].Balance)
	})

	t.Run("deposit with invalid state", func(t *testing.T) {
		result := DepositFunds(&user1Address, types.NewUint256(100), "{invalid json}")
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse application state")
	})
}

func TestProcessRequest(t *testing.T) {
	stateJSON, _ := getPopulatedState(t)

	t.Run("withdraw success", func(t *testing.T) {
		withdrawAmount := types.NewUint256(200)
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user3Address,
				Amount: withdrawAmount,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, types.NewUint256(800), newState.Accounts[user1Address.Hex()].Balance)

		withdrawal := result.Withdrawals[0]
		require.Equal(t, user3Address, withdrawal.DestinationAddress)
		require.Equal(t, withdrawAmount, withdrawal.Amount)
	})

	t.Run("withdraw insufficient balance", func(t *testing.T) {
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user3Address,
				Amount: types.NewUint256(2000),
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Insufficient balance")
	})

	t.Run("withdraw from non-existent account", func(t *testing.T) {
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user1Address,
				Amount: types.NewUint256(100),
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		nonexistent, _ := types.HexToAddress("0xadd0000000000000000000000000000000009999")
		result := ProcessRequest(&nonexistent, int32(common.Process), string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "does not exist")
	})

	t.Run("withdraw with missing instruction", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "withdraw"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Withdraw instruction is missing")
	})

	t.Run("withdraw zero amount", func(t *testing.T) {
		withdrawAmount := types.NewUint256(0)
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user3Address,
				Amount: withdrawAmount,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, types.NewUint256(1000), newState.Accounts[user1Address.Hex()].Balance) // balance should not change

		withdrawal := result.Withdrawals[0]
		require.Equal(t, user3Address, withdrawal.DestinationAddress)
		require.Equal(t, withdrawAmount, withdrawal.Amount)
	})

	t.Run("withdraw exact balance", func(t *testing.T) {
		withdrawAmount := types.NewUint256(500)
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user3Address,
				Amount: withdrawAmount,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user2Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, types.NewUint256(0), newState.Accounts[user2Address.Hex()].Balance)

		withdrawal := result.Withdrawals[0]
		require.Equal(t, user3Address, withdrawal.DestinationAddress)
		require.Equal(t, withdrawAmount, withdrawal.Amount)
	})

	t.Run("compare addresses richer", func(t *testing.T) {
		instruction := PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &CompareInstructions{
				TargetAddress: user2Address,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Empty(t, result.Withdrawals)

		var eventData map[string]interface{}
		err = json.Unmarshal(result.Events[0].Data, &eventData)
		require.NoError(t, err)
		require.Contains(t, eventData["sentence"], " richer than ")
	})

	t.Run("compare addresses poorer", func(t *testing.T) {
		instruction := PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &CompareInstructions{
				TargetAddress: user1Address,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user2Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)

		var eventData map[string]interface{}
		err = json.Unmarshal(result.Events[0].Data, &eventData)
		require.NoError(t, err)
		require.Contains(t, eventData["sentence"], " poorer than ")
	})

	t.Run("compare addresses equal wealth", func(t *testing.T) {
		// Create a new state for this test
		state := ApplicationInternalState{
			AppID: testAppId,
			Accounts: map[string]*AccountState{
				user1Address.Hex(): {Address: user1Address, Balance: types.NewUint256(1000)},
				user2Address.Hex(): {Address: user2Address, Balance: types.NewUint256(1000)},
			},
		}
		stateBytes, err := json.Marshal(state)
		require.NoError(t, err)
		localStateJSON := string(stateBytes)

		instruction := PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &CompareInstructions{
				TargetAddress: user2Address,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Process), string(payloadBytes), localStateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)

		var eventData map[string]interface{}
		err = json.Unmarshal(result.Events[0].Data, &eventData)
		require.NoError(t, err)
		require.Contains(t, eventData["sentence"], " as wealthy as ")
	})

	t.Run("compare with non-existent account", func(t *testing.T) {
		nonexistent, _ := types.HexToAddress("0xadd0000000000000000000000000000000009999")
		instruction := PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &CompareInstructions{
				TargetAddress: nonexistent,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&nonexistent, int32(common.Process), string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Equal(t, "Account "+nonexistent.Hex()+" does not exist!", result.Error)
	})

	t.Run("compare with missing instruction", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "compare_addresses"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Compare instruction is missing")
	})

	t.Run("unsupported instruction type", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "invalid_type"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Process), string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Equal(t, "Unsupported instruction type: [invalid_type]", result.Error)
	})

	t.Run("invalid payload json", func(t *testing.T) {
		result := ProcessRequest(&user1Address, int32(common.Process), "{invalid json}", stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse payload instructions")
	})

	t.Run("invalid state json", func(t *testing.T) {
		result := ProcessRequest(&user1Address, int32(common.Process), "{}", "{invalid json}")
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse application state")
	})

	t.Run("empty payload", func(t *testing.T) {
		result := ProcessRequest(&user1Address, int32(common.Process), "", stateJSON)
		require.Empty(t, result.Error)
		require.Empty(t, result.Events)
		require.Empty(t, result.Withdrawals)

		var newState ApplicationInternalState
		err := json.Unmarshal(result.State, &newState)
		require.NoError(t, err)

		var originalState ApplicationInternalState
		err = json.Unmarshal([]byte(stateJSON), &originalState)
		require.NoError(t, err)

		require.Equal(t, originalState, newState)
	})
}

func TestDeanonymizationViaProcessRequest(t *testing.T) {
	stateJSON, state := getPopulatedState(t)

	t.Run("successful report generation via process request with payload", func(t *testing.T) {
		payload := PayloadInstructions{
			Type:        "deanonymize",
			Deanonymize: &DeanonymizeInstruction{IncludeTag: "SIMPLE_REPORT"},
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Deanonymize), string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.NotNil(t, result.Report)

		var report DeanonymizationReport
		err = json.Unmarshal(result.Report, &report)
		require.NoError(t, err)

		require.Equal(t, "SIMPLE_REPORT", report.Tag)

		// Check if the accounts in the report match the expected ones (from state)
		require.Equal(t, len(state.Accounts), len(report.Accounts))
		for _, expectedAcc := range state.Accounts {
			found := false
			for _, reportAcc := range report.Accounts {
				if reportAcc.Address == expectedAcc.Address {
					found = reportAcc.Balance.Cmp(*expectedAcc.Balance) == 0
					break
				}
			}
			require.True(t, found, "Account %s not found in report", expectedAcc.Address.Hex())
		}
	})

	t.Run("successful report generation with empty payload and int32(common.Deanonymize)", func(t *testing.T) {
		// When RequestType is Deanonymize, empty payload should still generate a report
		result := ProcessRequest(&user1Address, int32(common.Deanonymize), "{}", stateJSON)
		require.Empty(t, result.Error)
		require.NotNil(t, result.Report)

		var report DeanonymizationReport
		err := json.Unmarshal(result.Report, &report)
		require.NoError(t, err)

		// Tag should be empty since no payload specified it
		require.Empty(t, report.Tag)

		// Check if the accounts in the report match the expected ones (from state)
		require.Equal(t, len(state.Accounts), len(report.Accounts))
	})

	t.Run("report with invalid state", func(t *testing.T) {
		payload := PayloadInstructions{
			Type: "deanonymize",
		}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, int32(common.Deanonymize), string(payloadBytes), "{invalid json}")
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse application state")
	})
}

// TestJsonCompatibility ensures that the local WASM guest types have the same JSON representation
// as the corresponding host types. This is crucial because the guest and host communicate
// via JSON serialization/deserialization, and any discrepancy would lead to communication failures.
func TestJsonCompatibility(t *testing.T) {
	// Address compatibility
	addrStr := "0x1234567890123456789012345678901234567890"
	guestAddr, err := types.HexToAddress(addrStr)
	require.NoError(t, err)
	hostAddr := ethCommon.HexToAddress(addrStr)

	// PlainEvent
	guestPlainEvent := types.PlainEvent{
		UserID: guestAddr,
		Data:   []byte("test data"),
	}
	jsonBytes, err := json.Marshal(guestPlainEvent)
	require.NoError(t, err)

	var hostPlainEvent common.PlainEvent
	err = json.Unmarshal(jsonBytes, &hostPlainEvent)
	require.NoError(t, err)
	require.Equal(t, hostAddr, hostPlainEvent.UserID)
	require.Equal(t, guestPlainEvent.Data, hostPlainEvent.Data)

	// Withdrawal
	guestWithdrawal := types.Withdrawal{
		DestinationAddress: guestAddr,
		Amount:             types.NewUint256(100),
	}
	jsonBytes, err = json.Marshal(guestWithdrawal)
	require.NoError(t, err)

	var hostWithdrawal common.Withdrawal
	err = json.Unmarshal(jsonBytes, &hostWithdrawal)
	require.NoError(t, err)
	require.Equal(t, hostAddr, hostWithdrawal.DestinationAddress)
	require.Equal(t, guestWithdrawal.Amount.String(), hostWithdrawal.Amount.String())

	// LoadModuleResult
	guestLoadModuleResult := types.LoadModuleResult{
		State: []byte("state"),
		Fuel:  types.NewUint256(10),
		Error: "load error",
	}
	jsonBytes, err = json.Marshal(guestLoadModuleResult)
	require.NoError(t, err)

	var hostLoadModuleResult wasmCommon.LoadModuleResult
	err = json.Unmarshal(jsonBytes, &hostLoadModuleResult)
	require.NoError(t, err)
	require.Equal(t, guestLoadModuleResult.State, hostLoadModuleResult.State)
	require.Equal(t, guestLoadModuleResult.Fuel.String(), hostLoadModuleResult.Fuel.String())
	require.Equal(t, guestLoadModuleResult.Error, hostLoadModuleResult.Error)

	// DepositResult
	guestDepositResult := types.DepositResult{
		State:  []byte("new state"),
		Events: []types.PlainEvent{guestPlainEvent},
		Fuel:   types.NewUint256(20),
		Error:  "deposit error",
	}
	jsonBytes, err = json.Marshal(guestDepositResult)
	require.NoError(t, err)

	var hostDepositResult wasmCommon.DepositResult
	err = json.Unmarshal(jsonBytes, &hostDepositResult)
	require.NoError(t, err)
	require.Equal(t, guestDepositResult.State, hostDepositResult.State)
	require.Len(t, hostDepositResult.Events, 1)
	require.Equal(t, hostAddr, hostDepositResult.Events[0].UserID)
	require.Equal(t, guestDepositResult.Fuel.String(), hostDepositResult.Fuel.String())
	require.Equal(t, guestDepositResult.Error, hostDepositResult.Error)

	// ProcessResult
	guestProcessResult := types.ProcessResult{
		State:       []byte("state"),
		Events:      []types.PlainEvent{guestPlainEvent},
		Withdrawals: []types.Withdrawal{guestWithdrawal},
		Fuel:        types.NewUint256(50),
		Error:       "error",
	}
	jsonBytes, err = json.Marshal(guestProcessResult)
	require.NoError(t, err)

	var hostProcessResult wasmCommon.ProcessResult
	err = json.Unmarshal(jsonBytes, &hostProcessResult)
	require.NoError(t, err)
	require.Equal(t, guestProcessResult.State, hostProcessResult.State)
	require.Len(t, hostProcessResult.Events, 1)
	require.Equal(t, hostAddr, hostProcessResult.Events[0].UserID)
	require.Len(t, hostProcessResult.Withdrawals, 1)
	require.Equal(t, hostAddr, hostProcessResult.Withdrawals[0].DestinationAddress)
	require.Equal(t, guestProcessResult.Fuel.String(), hostProcessResult.Fuel.String())
	require.Equal(t, guestProcessResult.Error, hostProcessResult.Error)

	// DepositEvent
	guestDepositEvent := DepositEvent{
		Type:    "deposit",
		Amount:  types.NewUint256(100),
		Balance: types.NewUint256(1100),
		Nonce:   5,
	}
	jsonBytes, err = json.Marshal(guestDepositEvent)
	require.NoError(t, err)

	var hDepositEvent hostDepositEvent
	err = json.Unmarshal(jsonBytes, &hDepositEvent)
	require.NoError(t, err)
	require.Equal(t, guestDepositEvent.Type, hDepositEvent.Type)
	require.Equal(t, guestDepositEvent.Amount.String(), hDepositEvent.Amount.String())
	require.Equal(t, guestDepositEvent.Balance.String(), hDepositEvent.Balance.String())
	require.Equal(t, guestDepositEvent.Nonce, hDepositEvent.Nonce)

	// SenderEvent / WithdrawalEvent
	guestSenderEvent := SenderEvent{
		Type:    "sender",
		To:      guestAddr,
		Amount:  types.NewUint256(200),
		Balance: types.NewUint256(800),
		Nonce:   6,
	}
	jsonBytes, err = json.Marshal(guestSenderEvent)
	require.NoError(t, err)

	var hSenderEvent hostSenderEvent
	err = json.Unmarshal(jsonBytes, &hSenderEvent)
	require.NoError(t, err)
	require.Equal(t, guestSenderEvent.Type, hSenderEvent.Type)
	require.Equal(t, hostAddr, hSenderEvent.To)
	require.Equal(t, guestSenderEvent.Amount.String(), hSenderEvent.Amount.String())
	require.Equal(t, guestSenderEvent.Balance.String(), hSenderEvent.Balance.String())
	require.Equal(t, guestSenderEvent.Nonce, hSenderEvent.Nonce)

	// RecipientEvent
	guestRecipientEvent := RecipientEvent{
		Type:    "recipient",
		From:    guestAddr,
		Amount:  types.NewUint256(200),
		Balance: types.NewUint256(1200),
		Nonce:   7,
	}
	jsonBytes, err = json.Marshal(guestRecipientEvent)
	require.NoError(t, err)

	var hRecipientEvent hostRecipientEvent
	err = json.Unmarshal(jsonBytes, &hRecipientEvent)
	require.NoError(t, err)
	require.Equal(t, guestRecipientEvent.Type, hRecipientEvent.Type)
	require.Equal(t, hostAddr, hRecipientEvent.From)
	require.Equal(t, guestRecipientEvent.Amount.String(), hRecipientEvent.Amount.String())
	require.Equal(t, guestRecipientEvent.Balance.String(), hRecipientEvent.Balance.String())
	require.Equal(t, guestRecipientEvent.Nonce, hRecipientEvent.Nonce)

	// MemoryStats
	guestMemoryStats := types.MemoryStats{
		MapSize:              10,
		CumulativeMemorySize: 100,
	}
	jsonBytes, err = json.Marshal(guestMemoryStats)
	require.NoError(t, err)

	var hostMemoryStats wasmCommon.MemoryStats
	err = json.Unmarshal(jsonBytes, &hostMemoryStats)
	require.NoError(t, err)
	require.Equal(t, guestMemoryStats.MapSize, hostMemoryStats.MapSize)
	require.Equal(t, guestMemoryStats.CumulativeMemorySize, hostMemoryStats.CumulativeMemorySize)
}

func TestHexToAddress_ValidFullLength(t *testing.T) {
	s := "0x00112233445566778899aabbccddeeff00112233"
	addr, err := types.HexToAddress(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte{
		0x00, 0x11, 0x22, 0x33, 0x44,
		0x55, 0x66, 0x77, 0x88, 0x99,
		0xaa, 0xbb, 0xcc, 0xdd, 0xee,
		0xff, 0x00, 0x11, 0x22, 0x33,
	}

	if !bytes.Equal(addr[:], expected) {
		t.Fatalf("address mismatch\nexpected: %x\ngot:      %x", expected, addr[:])
	}
}

func TestHexToAddress_NoPrefix(t *testing.T) {
	s := "00112233445566778899aabbccddeeff00112233"
	_, err := types.HexToAddress(s)
	if err == nil {
		t.Fatalf("expected error for address without 0x prefix")
	}

	// With prefix should still work
	_, err = types.HexToAddress("0x" + s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHexToAddress_InvalidLength(t *testing.T) {
	_, err := types.HexToAddress("0x1234")
	if err == nil {
		t.Fatalf("expected error for short address")
	}

	_, err = types.HexToAddress("0x" + strings.Repeat("11", 21))
	if err == nil {
		t.Fatalf("expected error for long address")
	}
}

func TestHexToAddress_OddLengthRejected(t *testing.T) {
	_, err := types.HexToAddress("abc")
	if err == nil {
		t.Fatalf("expected error for odd-length input")
	}
}

func TestHexToAddress_InvalidHex(t *testing.T) {
	_, err := types.HexToAddress("0xzz1122")
	if err == nil {
		t.Fatalf("expected error for invalid hex, got nil")
	}
}

func TestHexToAddress_Empty(t *testing.T) {
	_, err := types.HexToAddress("")
	if err == nil {
		t.Fatalf("expected error for empty string, got nil")
	}
}

func TestBigIntUint256JSONRoundTrip(t *testing.T) {
	// Step 1: start with a big.Int value
	orig := new(big.Int)
	orig.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10) // 2^256-1

	// Step 2: marshal *common.Big into JSON
	type HostStruct struct {
		Amount *common.Big `json:"amount"`
	}
	hostObj := HostStruct{Amount: common.ToBig(orig)}

	jsonData, err := json.Marshal(hostObj)
	if err != nil {
		t.Fatalf("failed to marshal host JSON: %v", err)
	}
	t.Logf("JSON output from host (big.Int): %s", string(jsonData))

	// Step 3: unmarshal JSON into Uint256 (like WASM side)
	type WASMStruct struct {
		Amount types.Uint256 `json:"amount"`
	}
	var wasmObj WASMStruct
	if err := json.Unmarshal(jsonData, &wasmObj); err != nil {
		t.Fatalf("failed to unmarshal into Uint256: %v", err)
	}

	// Step 4: marshal Uint256 back to JSON (WASM → host)
	jsonData2, err := json.Marshal(wasmObj)
	if err != nil {
		t.Fatalf("failed to marshal Uint256 back to JSON: %v", err)
	}
	t.Logf("JSON output from Uint256: %s", string(jsonData2))

	// Step 5: unmarshal back into big.Int (host)
	var hostObj2 HostStruct
	if err := json.Unmarshal(jsonData2, &hostObj2); err != nil {
		t.Fatalf("failed to unmarshal JSON back into big.Int: %v", err)
	}

	// Step 6: compare
	if orig.Cmp(hostObj2.Amount.ToInt()) != 0 {
		t.Errorf("round-trip mismatch:\noriginal: %s\nfinal:    %s", orig.String(), hostObj2.Amount.String())
	} else {
		t.Logf("Round-trip successful: value preserved exactly")
	}
}

// TestUint256BigJSONCompatibility verifies that Uint256 and common.Big produce
// identical JSON representations and can be unmarshaled interchangeably.
func TestUint256BigJSONCompatibility(t *testing.T) {
	// Test values covering edge cases
	testValues := []struct {
		name    string
		decimal string
		hex     string
	}{
		{"zero", "0", "0x0"},
		{"one", "1", "0x1"},
		{"small", "255", "0xff"},
		{"medium", "12345", "0x3039"},
		{"large", "12345678901234567890", "0xab54a98ceb1f0ad2"},
		{"max_uint64", "18446744073709551615", "0xffffffffffffffff"},
		{"max_uint128", "340282366920938463463374607431768211455", "0xffffffffffffffffffffffffffffffff"},
		{"max_uint256", "115792089237316195423570985008687907853269984665640564039457584007913129639935", "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
	}

	for _, tt := range testValues {
		t.Run(tt.name, func(t *testing.T) {
			// Create big.Int from decimal
			bi, ok := new(big.Int).SetString(tt.decimal, 10)
			require.True(t, ok, "Failed to parse decimal: %s", tt.decimal)

			// Create common.Big and Uint256 from the same value
			bigVal := common.ToBig(bi)
			uint256Val := new(types.Uint256).SetBytes(bi.Bytes())

			// Marshal both to JSON
			bigJSON, err := json.Marshal(bigVal)
			require.NoError(t, err, "Failed to marshal common.Big")

			uint256JSON, err := json.Marshal(uint256Val)
			require.NoError(t, err, "Failed to marshal Uint256")

			// Verify JSON is byte-for-byte identical
			require.Equal(t, string(bigJSON), string(uint256JSON),
				"JSON mismatch for %s:\n  common.Big: %s\n  Uint256:    %s",
				tt.name, string(bigJSON), string(uint256JSON))

			// Verify JSON matches expected hex format
			expectedJSON := `"` + tt.hex + `"`
			require.Equal(t, expectedJSON, string(bigJSON),
				"Unexpected JSON format for %s", tt.name)

			// Test cross-type unmarshal: Big JSON → Uint256
			var u256FromBig types.Uint256
			err = json.Unmarshal(bigJSON, &u256FromBig)
			require.NoError(t, err, "Failed to unmarshal Big JSON into Uint256")
			require.Equal(t, tt.decimal, u256FromBig.String(),
				"Value mismatch after Big→Uint256 unmarshal")

			// Test cross-type unmarshal: Uint256 JSON → Big
			var bigFromU256 common.Big
			err = json.Unmarshal(uint256JSON, &bigFromU256)
			require.NoError(t, err, "Failed to unmarshal Uint256 JSON into Big")
			require.Equal(t, tt.decimal, bigFromU256.String(),
				"Value mismatch after Uint256→Big unmarshal")
		})
	}
}

// TestUint256BigStructCompatibility verifies that structs containing Uint256
// and common.Big fields produce compatible JSON.
func TestUint256BigStructCompatibility(t *testing.T) {
	// Simulate host-side struct (uses common.Big)
	type HostStruct struct {
		Amount  *common.Big `json:"amount"`
		Fee     *common.Big `json:"fee"`
		Balance *common.Big `json:"balance"`
	}

	// Simulate WASM-side struct (uses Uint256)
	type WASMStruct struct {
		Amount  types.Uint256 `json:"amount"`
		Fee     types.Uint256 `json:"fee"`
		Balance types.Uint256 `json:"balance"`
	}

	// Create host struct with test values
	hostStruct := HostStruct{
		Amount:  common.NewBig(1000000),
		Fee:     common.NewBig(100),
		Balance: common.NewBig(999900),
	}

	// Marshal host struct
	hostJSON, err := json.Marshal(hostStruct)
	require.NoError(t, err)
	t.Logf("Host JSON: %s", string(hostJSON))

	// Unmarshal into WASM struct
	var wasmStruct WASMStruct
	err = json.Unmarshal(hostJSON, &wasmStruct)
	require.NoError(t, err)

	// Verify values match
	require.Equal(t, "1000000", wasmStruct.Amount.String())
	require.Equal(t, "100", wasmStruct.Fee.String())
	require.Equal(t, "999900", wasmStruct.Balance.String())

	// Marshal WASM struct back
	wasmJSON, err := json.Marshal(wasmStruct)
	require.NoError(t, err)
	t.Logf("WASM JSON: %s", string(wasmJSON))

	// JSON should be identical
	require.JSONEq(t, string(hostJSON), string(wasmJSON))

	// Unmarshal back into host struct
	var hostStruct2 HostStruct
	err = json.Unmarshal(wasmJSON, &hostStruct2)
	require.NoError(t, err)

	// Final values should match original
	require.Equal(t, 0, hostStruct.Amount.ToInt().Cmp(hostStruct2.Amount.ToInt()))
	require.Equal(t, 0, hostStruct.Fee.ToInt().Cmp(hostStruct2.Fee.ToInt()))
	require.Equal(t, 0, hostStruct.Balance.ToInt().Cmp(hostStruct2.Balance.ToInt()))
}
