package app

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testAppId    = "test_app"
	user1Address = "0xadd0000000000000000000000000000000000001"
	user2Address = "0xadd0000000000000000000000000000000000002"
	user3Address = "0xadd0000000000000000000000000000000000003"
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
			user1Address: {Address: user1Address, Balance: big.NewInt(1000)},
			user2Address: {Address: user2Address, Balance: big.NewInt(500)},
		},
	}
	stateBytes, err := json.Marshal(state)
	require.NoError(t, err)
	return string(stateBytes), state
}

func TestLoadModule(t *testing.T) {
	stateBytes := LoadModule(testAppId)
	require.NotNil(t, stateBytes)

	var state ApplicationInternalState
	err := json.Unmarshal(stateBytes, &state)
	require.NoError(t, err)

	require.Equal(t, testAppId, state.AppID)
	require.Empty(t, state.Accounts)
}

func TestDepositFunds(t *testing.T) {
	t.Run("deposit to new account", func(t *testing.T) {
		stateJSON, _ := getInitialState(t)
		depositAmount := big.NewInt(100)

		result := DepositFunds(user1Address, depositAmount, stateJSON)
		require.Empty(t, result.Error)
		require.NotNil(t, result.State)
		require.Len(t, result.Events, 1)

		var newState ApplicationInternalState
		err := json.Unmarshal(result.State, &newState)
		require.NoError(t, err)

		require.Len(t, newState.Accounts, 1)
		require.Equal(t, depositAmount, newState.Accounts[user1Address].Balance)
		require.Equal(t, user1Address, newState.Accounts[user1Address].Address)

		event := result.Events[0]
		require.Equal(t, user1Address, event.UserID)
		var eventData map[string]interface{}
		err = json.Unmarshal(event.Data, &eventData)
		require.NoError(t, err)
		require.Equal(t, "deposit", eventData["type"])
		floatValue, _ := depositAmount.Float64()
		
		require.Equal(t, floatValue, eventData["amount"])
	})

	t.Run("deposit to existing account", func(t *testing.T) {
		_, state := getInitialState(t)
		initialBalance := big.NewInt(50)
		state.Accounts[user1Address] = &AccountState{Address: user1Address, Balance: initialBalance}
		stateBytes, err := json.Marshal(state)
		require.NoError(t, err)
		stateJSON := string(stateBytes)

		depositAmount := big.NewInt(100)
		result := DepositFunds(user1Address, depositAmount, stateJSON)
		require.Empty(t, result.Error)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)

		sum := big.NewInt(0).Add(initialBalance, depositAmount)
		require.Equal(t, sum, newState.Accounts[user1Address].Balance)
	})

	t.Run("deposit with invalid state", func(t *testing.T) {
		result := DepositFunds(user1Address, big.NewInt(100), "{invalid json}")
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

		result := ProcessRequest(user1Address, string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(800), newState.Accounts[user1Address].Balance)

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

		result := ProcessRequest(user1Address, string(payloadBytes), stateJSON)
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

		result := ProcessRequest("nonexistent", string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "does not exist")
	})

	t.Run("withdraw with missing instruction", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "withdraw"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(user1Address, string(payloadBytes), stateJSON)
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

		result := ProcessRequest(user1Address, string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(1000), newState.Accounts[user1Address].Balance) // balance should not change

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

		result := ProcessRequest(user2Address, string(payloadBytes), stateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)
		require.Len(t, result.Withdrawals, 1)

		var newState ApplicationInternalState
		err = json.Unmarshal(result.State, &newState)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(0), newState.Accounts[user2Address].Balance)

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

		result := ProcessRequest(user1Address, string(payloadBytes), stateJSON)
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

		result := ProcessRequest(user2Address, string(payloadBytes), stateJSON)
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
				user1Address: {Address: user1Address, Balance: big.NewInt(1000)},
				user2Address: {Address: user2Address, Balance: big.NewInt(1000)},
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

		result := ProcessRequest(user1Address, string(payloadBytes), localStateJSON)
		require.Empty(t, result.Error)
		require.Len(t, result.Events, 1)

		var eventData map[string]interface{}
		err = json.Unmarshal(result.Events[0].Data, &eventData)
		require.NoError(t, err)
		require.Contains(t, eventData["sentence"], " as wealthy as ")
	})

	t.Run("compare with non-existent account", func(t *testing.T) {
		instruction := PayloadInstructions{
			Type: "compare_addresses",
			CompareAccounts: &CompareInstructions{
				TargetAddress: "nonexistent",
			},
		}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(user1Address, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Equal(t, "Account nonexistent does not exist!", result.Error)
	})

	t.Run("compare with missing instruction", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "compare_addresses"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(user1Address, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Compare instruction is missing")
	})

	t.Run("unsupported instruction type", func(t *testing.T) {
		instruction := PayloadInstructions{Type: "invalid_type"}
		payloadBytes, err := json.Marshal(instruction)
		require.NoError(t, err)

		result := ProcessRequest(user1Address, string(payloadBytes), stateJSON)
		require.NotEmpty(t, result.Error)
		require.Equal(t, "Unsupported instruction type: [invalid_type]", result.Error)
	})

	t.Run("invalid payload json", func(t *testing.T) {
		result := ProcessRequest(user1Address, "{invalid json}", stateJSON)
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse payload instructions")
	})

	t.Run("invalid state json", func(t *testing.T) {
		result := ProcessRequest(user1Address, "{}", "{invalid json}")
		require.NotEmpty(t, result.Error)
		require.Contains(t, result.Error, "Failed to parse application state")
	})

	t.Run("empty payload", func(t *testing.T) {
		result := ProcessRequest(user1Address, "", stateJSON)
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

		var report map[string]interface{}
		err := json.Unmarshal(result.Report, &report)
		require.NoError(t, err)

		require.Equal(t, "SIMPLE_REPORT", report["tag"])

		accountsData, ok := report["accounts"].(map[string]interface{})
		require.True(t, ok)
		require.Len(t, accountsData, len(state.Accounts))
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
