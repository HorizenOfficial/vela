package app

import (
	"encoding/json"
	"fmt"

	"github.com/horizen-pes/app/simple/utils"

	"github.com/horizen-pes/pkg/common"
	wasmCommon "github.com/horizen-pes/pkg/wasm/common"
)

// AccountState represents the state of a user account
type AccountState struct {
	Address string `json:"address"`
	Balance uint64 `json:"balance"`
}

// ApplicationInternalState represents the internal state of the application
type ApplicationInternalState struct {
	AppID    string                   `json:"appId"`
	Accounts map[string]*AccountState `json:"accounts"`
}

// WithdrawInstruction represents instructions for withdrawing funds
type WithdrawInstruction struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

type CompareInstructions struct {
	TargetAddress string `json:"targetAddress"`
}

// PayloadInstructions represents the deserialized payload instructions
type PayloadInstructions struct {
	Type            string               `json:"type"`
	CompareAccounts *CompareInstructions `json:"compare,omitempty"`
	Withdraw        *WithdrawInstruction `json:"withdraw,omitempty"`
}

// --- High-Level Application Logic ---

func LoadModule(appId string) []byte {
	initialState := &ApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[string]*AccountState),
	}
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return utils.WasmSerializationError
	}
	return stateJSON
}

func DepositFunds(sender string, value uint64, stateJSON string) wasmCommon.DepositResult {
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return wasmCommon.DepositResult{Error: "Failed to parse application state"}
	}

	// Ensure sender account exists
	if currentState.Accounts[sender] == nil {
		currentState.Accounts[sender] = &AccountState{
			Address: sender,
			Balance: 0,
		}
	}

	// Add deposit to sender's balance
	currentState.Accounts[sender].Balance += value

	// Create deposit event
	eventData := map[string]interface{}{
		"type":   "deposit",
		"amount": value,
	}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		return wasmCommon.DepositResult{Error: "Failed to serialize event data"}
	}

	events := []common.PlainEvent{{
		UserID: sender,
		Data:   eventDataBytes,
	}}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		return wasmCommon.DepositResult{Error: "Failed to serialize new state"}
	}
	return wasmCommon.DepositResult{State: newStateBytes, Events: events}
}

func ProcessRequest(sender, payloadJSON, stateJSON string) wasmCommon.ProcessResult {
	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return wasmCommon.ProcessResult{Error: "Failed to parse application state"}
	}

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal

	if payloadJSON != "" {
		var instructions PayloadInstructions
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			return wasmCommon.ProcessResult{Error: "Failed to parse payload instructions"}
		}

		switch instructions.Type {
		case "compare_addresses":
			if instructions.CompareAccounts == nil {
				return wasmCommon.ProcessResult{Error: "Compare instruction is missing"}
			}
			targetAddress := instructions.CompareAccounts.TargetAddress

			// Validate accounts to be compared
			if currentState.Accounts[sender] == nil {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Account %s does not exist!", sender)}
			}
			if currentState.Accounts[targetAddress] == nil {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Account %s does not exist!", targetAddress)}
			}

			targetBalance := currentState.Accounts[targetAddress].Balance
			senderBalance := currentState.Accounts[sender].Balance

			var cmp = ""
			if targetBalance < senderBalance {
				cmp = "richer than"
			} else if targetBalance > senderBalance {
				cmp = "poorer than"
			} else {
				cmp = "as wealthy as"
			}
			sentence := sender + " is " + cmp + " " + targetAddress

			// Create action event
			eventData := map[string]interface{}{
				"type":     "compare_accounts",
				"sentence": sentence,
			}
			eventDataBytes, err := json.Marshal(eventData)
			if err != nil {
				return wasmCommon.ProcessResult{Error: "Failed to serialize event data"}
			}

			events = append(events, common.PlainEvent{
				UserID: sender,
				Data:   eventDataBytes,
			})

		case "withdraw":
			if instructions.Withdraw == nil {
				return wasmCommon.ProcessResult{Error: "Withdraw instruction is missing"}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return wasmCommon.ProcessResult{Error: "Account does not exist"}
			}

			if currentState.Accounts[sender].Balance < instructions.Withdraw.Amount {
				return wasmCommon.ProcessResult{Error: "Insufficient balance for withdrawal"}
			}

			// Execute withdrawal
			currentState.Accounts[sender].Balance -= instructions.Withdraw.Amount

			// Create withdrawal
			withdrawals = append(withdrawals, common.Withdrawal{
				DestinationAddress: instructions.Withdraw.To,
				Amount:             instructions.Withdraw.Amount,
			})

			// Create event for sender
			withdrawEventData := wasmCommon.WithdrawalEvent{
				Type:    "withdrawal",
				To:      instructions.Withdraw.To,
				Amount:  instructions.Withdraw.Amount,
				Balance: currentState.Accounts[sender].Balance,
			}
			withdrawEventDataBytes, err := json.Marshal(withdrawEventData)
			if err != nil {
				return wasmCommon.ProcessResult{Error: "Failed to serialize withdraw event data"}
			}

			events = append(events, common.PlainEvent{
				UserID: sender,
				Data:   withdrawEventDataBytes,
			})

		default:
			return wasmCommon.ProcessResult{Error: fmt.Sprintf("Unsupported instruction type: [%s]", instructions.Type)}
		}
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		return wasmCommon.ProcessResult{Error: "Failed to serialize new state"}
	}
	return wasmCommon.ProcessResult{
		State:       newStateBytes,
		Events:      events,
		Withdrawals: withdrawals,
	}
}

func GenerateDeanonymizationReport(stateJSON string) wasmCommon.DeanonymizationResult {
	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return wasmCommon.DeanonymizationResult{Error: "Failed to parse application state"}
	}

	// Create deanonymization report
	report := map[string]interface{}{
		"accounts": currentState.Accounts,
	}

	// Serialize the report
	reportBytes, err := json.Marshal(report)
	if err != nil {
		return wasmCommon.DeanonymizationResult{Error: "Failed to serialize deanonymization report"}
	}
	return wasmCommon.DeanonymizationResult{Report: reportBytes}
}
