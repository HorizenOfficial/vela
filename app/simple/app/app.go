package app

import (
	"encoding/json"
	"fmt"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	wasmCommon "github.com/horizen-pes/pkg/wasm/common"
)

// AccountState represents the state of a user account
type AccountState struct {
	Address ethCommon.Address `json:"address"`
	Balance *big.Int          `json:"balance"`
}

// ApplicationInternalState represents the internal state of the application
type ApplicationInternalState struct {
	AppID    int64                               `json:"appId"`
	Accounts map[ethCommon.Address]*AccountState `json:"accounts"`
}

// WithdrawInstruction represents instructions for withdrawing funds
type WithdrawInstruction struct {
	To     ethCommon.Address `json:"to"`
	Amount *big.Int          `json:"amount"`
}

type CompareInstructions struct {
	TargetAddress ethCommon.Address `json:"targetAddress"`
}

// PayloadInstructions represents the deserialized payload instructions
type PayloadInstructions struct {
	Type            string               `json:"type"`
	CompareAccounts *CompareInstructions `json:"compare,omitempty"`
	Withdraw        *WithdrawInstruction `json:"withdraw,omitempty"`
}

// ReportPayloadInstructions represent a specific information on how to generate a report
// In this simple app its a custom tag to add to the report
type ReportPayloadInstructions struct {
	IncludeTag string `json:"tag,omitempty"`
}

// DeanonymizationReport represents the structure of the deanonymization report.
type DeanonymizationReport struct {
	Tag      string                              `json:"tag,omitempty"`
	Accounts map[ethCommon.Address]*AccountState `json:"accounts"`
}

// --- High-Level Application Logic ---

func LoadModule(appId int64) []byte {
	initialState := &ApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[ethCommon.Address]*AccountState),
	}
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return []byte(wasmCommon.WasmSerializationError)
	}
	return stateJSON
}

func DepositFunds(senderPtr *ethCommon.Address, value *big.Int, stateJSON string) wasmCommon.DepositResult {
	if senderPtr == nil {
		return wasmCommon.DepositResult{Error: "Sender address is nil"}
	}

	sender := *senderPtr
	//This should never happens but just in case
	if value == nil {
		return wasmCommon.DepositResult{Error: "value is nil"}
	}

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return wasmCommon.DepositResult{Error: fmt.Sprintf("Failed to parse application state: %s", stateJSON)}
	}

	// Ensure sender account exists
	if currentState.Accounts[sender] == nil {
		currentState.Accounts[sender] = &AccountState{
			Address: sender,
			Balance: big.NewInt(0),
		}
	}

	// Add deposit to sender's balance
	currentState.Accounts[sender].Balance.Add(currentState.Accounts[sender].Balance, value)

	// Create deposit event
	eventData := map[string]interface{}{
		"type":   "deposit",
		"amount": value,
	}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		return wasmCommon.DepositResult{Error: fmt.Sprintf("Failed to serialize event data: %+v", eventData)}
	}

	events := []common.PlainEvent{{
		UserID: sender,
		Data:   eventDataBytes,
	}}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		return wasmCommon.DepositResult{Error: fmt.Sprintf("Failed to serialize new state: %+v", &currentState)}
	}
	return wasmCommon.DepositResult{State: newStateBytes, Events: events}
}

func ProcessRequest(senderPtr *ethCommon.Address, payloadJSON, stateJSON string) wasmCommon.ProcessResult {
	// Deserialize current state
	if senderPtr == nil {
		return wasmCommon.ProcessResult{Error: "Sender address is nil"}
	}

	sender := *senderPtr

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return wasmCommon.ProcessResult{Error: fmt.Sprintf("Failed to parse application state: %s", stateJSON)}
	}

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal

	if payloadJSON != "" {
		var instructions PayloadInstructions
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			return wasmCommon.ProcessResult{Error: fmt.Sprintf("Failed to parse payload instructions: %s", payloadJSON)}
		}

		switch instructions.Type {
		case "compare_addresses":
			if instructions.CompareAccounts == nil {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Compare instruction is missing in payload: %s", payloadJSON)}
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
			switch targetBalance.Cmp(senderBalance) {
			case -1:
				cmp = "richer than"
			case 1:
				cmp = "poorer than"
			case 0:
				cmp = "as wealthy as"

			}

			sentence := sender.Hex() + " is " + cmp + " " + targetAddress.Hex()

			// Create action event
			eventData := map[string]interface{}{
				"type":     "compare_accounts",
				"sentence": sentence,
			}
			eventDataBytes, err := json.Marshal(eventData)
			if err != nil {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Failed to serialize event data: %+v", eventData)}
			}

			events = append(events, common.PlainEvent{
				UserID: sender,
				Data:   eventDataBytes,
			})

		case "withdraw":
			if instructions.Withdraw == nil {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Withdraw instruction is missing in payload: %s", payloadJSON)}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Account %s does not exist", sender)}
			}

			if currentState.Accounts[sender].Balance.Cmp(instructions.Withdraw.Amount) < 0 {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Insufficient balance for withdrawal for account %s", sender)}
			}

			// Execute withdrawal
			currentState.Accounts[sender].Balance.Sub(currentState.Accounts[sender].Balance, instructions.Withdraw.Amount)

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
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Failed to serialize withdraw event data: %+v", withdrawEventData)}
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
		return wasmCommon.ProcessResult{Error: fmt.Sprintf("Failed to serialize new state: %+v", currentState)}
	}
	return wasmCommon.ProcessResult{
		State:       newStateBytes,
		Events:      events,
		Withdrawals: withdrawals,
	}
}

func GenerateDeanonymizationReport(payloadJSON, stateJSON string) wasmCommon.DeanonymizationResult {
	// Deserialize payload
	var payload ReportPayloadInstructions
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return wasmCommon.DeanonymizationResult{Error: fmt.Sprintf("Failed to parse payload: %s", payloadJSON)}
	}

	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return wasmCommon.DeanonymizationResult{Error: fmt.Sprintf("Failed to parse application state: %s", stateJSON)}
	}

	// Create deanonymization report
	report := DeanonymizationReport{
		Accounts: currentState.Accounts,
	}

	// read contents of the payload and decide how to build the report. In this simple case just add a tag if any
	if payload.IncludeTag != "" {
		report.Tag = payload.IncludeTag
	}

	// Serialize the report
	reportBytes, err := json.Marshal(report)
	if err != nil {
		return wasmCommon.DeanonymizationResult{Error: fmt.Sprintf("Failed to serialize deanonymization report: %+v", report)}
	}
	return wasmCommon.DeanonymizationResult{Report: reportBytes}
}
