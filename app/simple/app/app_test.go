package app

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testAppId = int64(1)
)

var (
	user1Address, _ = HexToAddress("0xadd0000000000000000000000000000000000001")
	user2Address, _ = HexToAddress("0xadd0000000000000000000000000000000000002")
	user3Address, _ = HexToAddress("0xadd0000000000000000000000000000000000003")
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
			user1Address.Hex(): {Address: user1Address, Balance: big.NewInt(1000)},
			user2Address.Hex(): {Address: user2Address, Balance: big.NewInt(500)},
		},
	}
	stateBytes, err := json.Marshal(state)
	require.NoError(t, err)
	return string(stateBytes), state
}

func TestLoadModule(t *testing.T) {
	result := LoadModule(testAppId)
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
		depositAmount := big.NewInt(100)

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
			Type   string   `json:"type"`
			Amount *big.Int `json:"amount"`
		}
		err = json.Unmarshal(event.Data, &eventData)
		require.NoError(t, err)
		require.Equal(t, "deposit", eventData.Type)

		require.Equal(t, depositAmount, eventData.Amount)
	})

	t.Run("deposit to existing account", func(t *testing.T) {
		_, state := getInitialState(t)
		initialBalance := big.NewInt(50)
		state.Accounts[user1Address.Hex()] = &AccountState{Address: user1Address, Balance: initialBalance}

		stateBytes, err := json.Marshal(state)
		require.NoError(t, err)
		stateJSON := string(stateBytes)

		depositAmount := big.NewInt(100)
		result := DepositFunds(&user1Address, depositAmount, stateJSON)
		require.Empty(t, result.Error)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)

		sum := big.NewInt(0).Add(initialBalance, depositAmount)
		require.Equal(t, sum, newState.Accounts[user1Address.Hex()].Balance)
	})

	t.Run("deposit with invalid state", func(t *testing.T) {
		result := DepositFunds(&user1Address, big.NewInt(100), "{invalid json}")
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse application state")
	})
}

func TestProcessRequest(t *testing.T) {
	stateJSON, _ := getPopulatedState(t)

	t.Run("withdraw success", func(t *testing.T) {
		withdrawAmount := big.NewInt(200)
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user3Address,
				Amount: withdrawAmount,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(800), newState.Accounts[user1Address.Hex()].Balance)

		withdrawal := result.Withdrawals[0]
		require.Equal(t, user3Address, withdrawal.DestinationAddress)
		require.Equal(t, withdrawAmount, withdrawal.Amount)
	})

	t.Run("withdraw insufficient balance", func(t *testing.T) {
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user3Address,
				Amount: big.NewInt(2000),
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Insufficient balance for withdrawal")
	})

	t.Run("withdraw from non-existent account", func(t *testing.T) {
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user1Address,
				Amount: big.NewInt(100),
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		nonexistent, _ := HexToAddress("0xadd0000000000000000000000000000000009999")
		result := ProcessRequest(&nonexistent, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "does not exist")
	})

	t.Run("withdraw with missing instruction", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "withdraw"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Withdraw instruction is missing")
	})

	t.Run("withdraw zero amount", func(t *testing.T) {
		withdrawAmount := big.NewInt(0)
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user3Address,
				Amount: withdrawAmount,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(1000), newState.Accounts[user1Address.Hex()].Balance) // balance should not change

		withdrawal := result.Withdrawals[0]
		require.Equal(t, user3Address, withdrawal.DestinationAddress)
		require.Equal(t, withdrawAmount, withdrawal.Amount)
	})

	t.Run("withdraw exact balance", func(t *testing.T) {
		withdrawAmount := big.NewInt(500)
		instruction := PayloadInstructions{
			Type: "withdraw",
			Withdraw: &WithdrawInstruction{
				To:     user3Address,
				Amount: withdrawAmount,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user2Address, string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(0), newState.Accounts[user2Address.Hex()].Balance)

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

		result := ProcessRequest(&user1Address, string(payloadBytes), stateJSON)
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

		result := ProcessRequest(&user2Address, string(payloadBytes), stateJSON)
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
				user1Address.Hex(): {Address: user1Address, Balance: big.NewInt(1000)},
				user2Address.Hex(): {Address: user2Address, Balance: big.NewInt(1000)},
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

		result := ProcessRequest(&user1Address, string(payloadBytes), localStateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)

		var eventData map[string]interface{}
		err = json.Unmarshal(result.Events[0].Data, &eventData)
		require.NoError(t, err)
		require.Contains(t, eventData["sentence"], " as wealthy as ")
	})

	t.Run("compare with non-existent account", func(t *testing.T) {
		nonexistent, _ := HexToAddress("0xadd0000000000000000000000000000000009999")
		instruction := PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &CompareInstructions{
				TargetAddress: nonexistent,
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&nonexistent, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Equal(t, "Account "+nonexistent.Hex()+" does not exist!", result.Error)
	})

	t.Run("compare with missing instruction", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "compare_addresses"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Compare instruction is missing")
	})

	t.Run("unsupported instruction type", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "invalid_type"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(&user1Address, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Equal(t, "Unsupported instruction type: [invalid_type]", result.Error)
	})

	t.Run("invalid payload json", func(t *testing.T) {
		result := ProcessRequest(&user1Address, "{invalid json}", stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse payload instructions")
	})

	t.Run("invalid state json", func(t *testing.T) {
		result := ProcessRequest(&user1Address, "{}", "{invalid json}")
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse application state")
	})

	t.Run("empty payload", func(t *testing.T) {
		result := ProcessRequest(&user1Address, "", stateJSON)
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

func TestGenerateDeanonymizationReport(t *testing.T) {
	stateJSON, state := getPopulatedState(t)
	payloadJSON := `{"tag":"SIMPLE_REPORT"}`

	t.Run("successful report generation", func(t *testing.T) {
		result := GenerateDeanonymizationReport(payloadJSON, stateJSON)
		require.Empty(t, result.Error)
		require.NotNil(t, result.Report)

		var report DeanonymizationReport
		err := json.Unmarshal(result.Report, &report)
		require.NoError(t, err)

		require.Equal(t, "SIMPLE_REPORT", report.Tag)

		// Check if the accounts in the report match the expected ones (from state)
		require.Equal(t, len(state.Accounts), len(report.Accounts))
		for _, expectedAcc := range state.Accounts {
			found := false
			for _, reportAcc := range report.Accounts {
				if reportAcc.Address == expectedAcc.Address {
					found = reportAcc.Balance.Cmp(expectedAcc.Balance) == 0
					break
				}
			}
			require.True(t, found, "Account %s not found in report", expectedAcc.Address.Hex())
		}
	})

	t.Run("report with invalid state", func(t *testing.T) {
		result := GenerateDeanonymizationReport("{}", "{invalid json}")
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse application state")
	})

	t.Run("report with invalid payload", func(t *testing.T) {
		result := GenerateDeanonymizationReport("{invalid json}", "{}")
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse payload")
	})
}

func TestHexToAddress_ValidFullLength(t *testing.T) {
	s := "0x00112233445566778899aabbccddeeff00112233"
	addr, err := HexToAddress(s)
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
	addr1, err := HexToAddress(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	addr2, err := HexToAddress("0x" + s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if addr1 != addr2 {
		t.Fatalf("addresses differ with/without prefix")
	}
}

func TestHexToAddress_InvalidLength(t *testing.T) {
	_, err := HexToAddress("0x1234")
	if err == nil {
		t.Fatalf("expected error for short address")
	}

	_, err = HexToAddress("0x" + strings.Repeat("11", 21))
	if err == nil {
		t.Fatalf("expected error for long address")
	}
}

func TestHexToAddress_OddLengthRejected(t *testing.T) {
	_, err := HexToAddress("abc")
	if err == nil {
		t.Fatalf("expected error for odd-length input")
	}
}

func TestHexToAddress_InvalidHex(t *testing.T) {
	_, err := HexToAddress("0xzz1122")
	if err == nil {
		t.Fatalf("expected error for invalid hex, got nil")
	}
}

func TestHexToAddress_Empty(t *testing.T) {
	_, err := HexToAddress("")
	if err == nil {
		t.Fatalf("expected error for empty string, got nil")
	}
}
