package app

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// --- High-Level Application Logic ---

func LoadModule(appId int64) []byte {
	initialState := &ApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[string]*AccountState),
	}
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return []byte(WasmSerializationError)
	}
	return stateJSON
}

func DepositFunds(senderPtr *Address, value *big.Int, stateJSON string) DepositResult {
	if senderPtr == nil {
		return DepositResult{Error: "Sender address is nil"}
	}

	sender := *senderPtr
	senderHex := sender.Hex()
	//This should never happens but just in case
	if value == nil {
		return DepositResult{Error: "value is nil"}
	}

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return DepositResult{Error: fmt.Sprintf("Failed to parse application state: %s", stateJSON)}
	}

	// Ensure sender account exists
	if currentState.Accounts[senderHex] == nil {
		currentState.Accounts[senderHex] = &AccountState{
			Address: sender,
			Balance: big.NewInt(0),
		}
	}

	// Add deposit to sender's balance
	currentState.Accounts[senderHex].Balance.Add(currentState.Accounts[senderHex].Balance, value)

	// Create deposit event
	eventData := map[string]interface{}{
		"type":   "deposit",
		"amount": value,
	}
	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		return DepositResult{Error: fmt.Sprintf("Failed to serialize event data: %+v", eventData)}
	}

	events := []PlainEvent{{
		UserID: sender,
		Data:   eventDataBytes,
	}}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		return DepositResult{Error: fmt.Sprintf("Failed to serialize new state: %+v", &currentState)}
	}
	return DepositResult{State: newStateBytes, Events: events}
}

func ProcessRequest(senderPtr *Address, payloadJSON, stateJSON string) ProcessResult {
	// Deserialize current state
	if senderPtr == nil {
		return ProcessResult{Error: "Sender address is nil"}
	}

	sender := *senderPtr
	senderHex := sender.Hex()

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return ProcessResult{Error: fmt.Sprintf("Failed to parse application state: %s", stateJSON)}
	}

	var events []PlainEvent
	var withdrawals []Withdrawal

	if payloadJSON != "" {
		var instructions PayloadInstructions
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			return ProcessResult{Error: fmt.Sprintf("Failed to parse payload instructions: %s", payloadJSON)}
		}

		switch instructions.Type {
		case "compare_addresses":
			if instructions.CompareAccounts == nil {
				return ProcessResult{Error: fmt.Sprintf("Compare instruction is missing in payload: %s", payloadJSON)}
			}
			targetAddress := instructions.CompareAccounts.TargetAddress
			targetAddressHex := targetAddress.Hex()

			// Validate accounts to be compared
			if currentState.Accounts[senderHex] == nil {
				return ProcessResult{Error: fmt.Sprintf("Account %s does not exist!", sender.Hex())}
			}
			if currentState.Accounts[targetAddressHex] == nil {
				return ProcessResult{Error: fmt.Sprintf("Account %s does not exist!", targetAddress.Hex())}
			}

			targetBalance := currentState.Accounts[targetAddressHex].Balance
			senderBalance := currentState.Accounts[senderHex].Balance

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
				return ProcessResult{Error: fmt.Sprintf("Failed to serialize event data: %+v", eventData)}
			}

			events = append(events, PlainEvent{
				UserID: sender,
				Data:   eventDataBytes,
			})

		case "withdraw":
			if instructions.Withdraw == nil {
				return ProcessResult{Error: fmt.Sprintf("Withdraw instruction is missing in payload: %s", payloadJSON)}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[senderHex] == nil {
				return ProcessResult{Error: fmt.Sprintf("Account %s does not exist", sender.Hex())}
			}

			if currentState.Accounts[senderHex].Balance.Cmp(instructions.Withdraw.Amount) < 0 {
				return ProcessResult{Error: fmt.Sprintf("Insufficient balance for withdrawal for account %s", sender.Hex())}
			}

			// Execute withdrawal
			currentState.Accounts[senderHex].Balance.Sub(currentState.Accounts[senderHex].Balance, instructions.Withdraw.Amount)

			// Create withdrawal
			withdrawals = append(withdrawals, Withdrawal{
				DestinationAddress: instructions.Withdraw.To,
				Amount:             instructions.Withdraw.Amount,
			})

			// Create event for sender
			withdrawEventData := WithdrawalEvent{
				Type:    "withdrawal",
				To:      instructions.Withdraw.To,
				Amount:  instructions.Withdraw.Amount,
				Balance: currentState.Accounts[senderHex].Balance,
			}
			withdrawEventDataBytes, err := json.Marshal(withdrawEventData)
			if err != nil {
				return ProcessResult{Error: fmt.Sprintf("Failed to serialize withdraw event data: %+v", withdrawEventData)}
			}

			events = append(events, PlainEvent{
				UserID: sender,
				Data:   withdrawEventDataBytes,
			})

		default:
			return ProcessResult{Error: fmt.Sprintf("Unsupported instruction type: [%s]", instructions.Type)}
		}
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		return ProcessResult{Error: fmt.Sprintf("Failed to serialize new state: %+v", currentState)}
	}
	return ProcessResult{
		State:       newStateBytes,
		Events:      events,
		Withdrawals: withdrawals,
	}
}

func GenerateDeanonymizationReport(payloadJSON, stateJSON string) DeanonymizationResult {
	// Deserialize payload
	var payload ReportPayloadInstructions
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return DeanonymizationResult{Error: fmt.Sprintf("Failed to parse payload: %s", payloadJSON)}
	}

	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return DeanonymizationResult{Error: fmt.Sprintf("Failed to parse application state: %s", stateJSON)}
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
		return DeanonymizationResult{Error: fmt.Sprintf("Failed to serialize deanonymization report: %+v", report)}
	}
	return DeanonymizationResult{Report: reportBytes}
}
