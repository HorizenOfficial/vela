package app

import (
	"encoding/json"
	"fmt"

	"github.com/horizen-pes/app/simple/utils"
	"github.com/horizen-pes/pkg/common"
)

// --- High-Level Application Logic ---

func LoadModule(appId int64) LoadModuleResult {
	initialState := &ApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[string]*AccountState),
	}
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		utils.LogError("LoadModule: failed to marshal initial state: %v", err)
		return LoadModuleResult{
			Error: fmt.Sprintf("failed to marshal initial state: %v", err),
		}
	}
	fuel := NewUint256(5)
	utils.LogDebug("LoadModule: appId=%d, stateSize=%d, fuel=%v", appId, len(stateJSON), fuel)
	return LoadModuleResult{
		State: stateJSON,
		Fuel:  fuel,
	}
}

func DepositFunds(senderPtr *Address, value *Uint256, stateJSON string) DepositResult {
	if senderPtr == nil {
		utils.LogError("DepositFunds: sender address is nil")
		return DepositResult{Error: "Sender address is nil"}
	}

	//This should never happens but just in case
	if value == nil {
		utils.LogError("DepositFunds: value is nil")
		return DepositResult{Error: "value is nil"}
	}

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		// we could add the stateJSON to the error, but it is not safe if it is very large
		utils.LogError("DepositFunds: failed to parse application state: %v", err)
		return DepositResult{Error: fmt.Sprintf("Failed to parse application state: %v", err)}
	}

	senderHex := senderPtr.Hex()

	// safe Map Access & Initialization
	acc, exists := currentState.Accounts[senderHex]
	if !exists {
		acc = &AccountState{
			Address: *senderPtr,
			Balance: NewUint256(0),
		}
		currentState.Accounts[senderHex] = acc
	}

	// Update balance in-place
	if acc.Balance.AddOverflow(*acc.Balance, *value) {
		return DepositResult{Error: fmt.Sprintf("Overflow while adding amount %s to balance: %s", value, acc.Balance)}
	}

	// Create deposit event
	type Event struct {
		Type   string   `json:"type"`
		Amount *Uint256 `json:"amount"`
	}

	eventData := Event{
		Type:   "deposit",
		Amount: value,
	}

	eventDataBytes, err := json.Marshal(eventData)
	if err != nil {
		utils.LogError("DepositFunds: failed to serialize event data: %v", err)
		return DepositResult{Error: fmt.Sprintf("Failed to serialize event data: %+v, err: %v", eventData, err)}
	}

	events := []PlainEvent{{
		UserID:       *senderPtr,
		EventSubType: "deposit",
		Data:         eventDataBytes,
	}}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		utils.LogError("DepositFunds: failed to serialize new state: %v", err)
		return DepositResult{Error: fmt.Sprintf("Failed to serialize new state: %v", err)}
	}
	fuel := NewUint256(35)
	utils.LogDebug("DepositFunds: sender=%s, value=%v, newBalance=%v, eventsCount=%d, stateSize=%d, fuel=%v",
		senderHex, value, acc.Balance, len(events), len(newStateBytes), fuel)
	return DepositResult{State: newStateBytes, Events: events, Fuel: fuel}
}

func ProcessRequest(senderPtr *Address, requestType int32, payloadJSON, stateJSON string) ProcessResult {
	// Deserialize current state
	if senderPtr == nil {
		utils.LogError("ProcessRequest: sender address is nil")
		return ProcessResult{Error: "Sender address is nil"}
	}

	sender := *senderPtr
	senderHex := sender.Hex()

	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		utils.LogError("ProcessRequest: failed to parse application state: %v", err)
		return ProcessResult{Error: fmt.Sprintf("Failed to parse application state: %v", err)}
	}

	var events []PlainEvent
	var withdrawals []Withdrawal

	// Determine instruction type: requestType has priority over payload
	var instructionType string
	var instructions PayloadInstructions

	// Parse payload if present (for additional options like deanonymize tag)
	if payloadJSON != "" && payloadJSON != "{}" {
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			utils.LogError("ProcessRequest: failed to parse payload instructions: %v", err)
			return ProcessResult{Error: fmt.Sprintf("Failed to parse payload instructions: %v", err)}
		}
	}

	// requestType takes precedence over payload type
	if requestType == int32(common.Deanonymize) {
		instructionType = "deanonymize"
	} else {
		instructionType = instructions.Type
	}

	if instructionType != "" {
		switch instructionType {
		case "compare_addresses":
			if instructions.CompareAccounts == nil {
				utils.LogError("ProcessRequest: compare instruction is missing in payload")
				return ProcessResult{Error: fmt.Sprintf("Compare instruction is missing in payload: %s", payloadJSON)}
			}
			targetAddress := instructions.CompareAccounts.TargetAddress
			targetAddressHex := targetAddress.Hex()

			// Validate accounts to be compared
			if currentState.Accounts[senderHex] == nil {
				utils.LogError("ProcessRequest: account %s does not exist", senderHex)
				return ProcessResult{Error: fmt.Sprintf("Account %s does not exist!", senderHex)}
			}
			if currentState.Accounts[targetAddressHex] == nil {
				utils.LogError("ProcessRequest: account %s does not exist", targetAddressHex)
				return ProcessResult{Error: fmt.Sprintf("Account %s does not exist!", targetAddressHex)}
			}

			targetBalance := currentState.Accounts[targetAddressHex].Balance
			senderBalance := currentState.Accounts[senderHex].Balance

			var cmp = ""
			switch targetBalance.Cmp(*senderBalance) {
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
				utils.LogError("ProcessRequest: failed to serialize event data: %v", err)
				return ProcessResult{Error: fmt.Sprintf("Failed to serialize event data: %+v, err: %v", eventData, err)}
			}

			events = append(events, PlainEvent{
				UserID:       sender,
				EventSubType: "compare_accounts",
				Data:         eventDataBytes,
			})

		case "withdraw":
			if instructions.Withdraw == nil {
				utils.LogError("ProcessRequest: withdraw instruction is missing in payload")
				return ProcessResult{Error: fmt.Sprintf("Withdraw instruction is missing in payload: %s", payloadJSON)}
			}
			if instructions.Withdraw.Amount == nil {
				return ProcessResult{Error: "Withdraw amount is nil"}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[senderHex] == nil {
				utils.LogError("ProcessRequest: account %s does not exist", senderHex)
				return ProcessResult{Error: fmt.Sprintf("Account %s does not exist", senderHex)}
			}

			if currentState.Accounts[senderHex].Balance.Cmp(*instructions.Withdraw.Amount) < 0 {
				utils.LogError("ProcessRequest: insufficient balance for account %s", senderHex)
				return ProcessResult{Error: fmt.Sprintf("Insufficient balance %s for withdrawal %s for account %s",
					currentState.Accounts[senderHex].Balance, *instructions.Withdraw.Amount, senderHex)}
			}

			// Execute withdrawal
			currentState.Accounts[senderHex].Balance.Sub(*currentState.Accounts[senderHex].Balance, *instructions.Withdraw.Amount)

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
				utils.LogError("ProcessRequest: failed to serialize withdraw event data: %v", err)
				return ProcessResult{Error: fmt.Sprintf("Failed to serialize withdraw event data: %+v, err: %v", withdrawEventData, err)}
			}

			events = append(events, PlainEvent{
				UserID:       sender,
				EventSubType: "withdrawal",
				Data:         withdrawEventDataBytes,
			})

		case "deanonymize":
			// Generate deanonymization report
			report := DeanonymizationReport{
				Accounts: currentState.Accounts,
			}

			// Add optional tag from payload
			if instructions.Deanonymize != nil && instructions.Deanonymize.IncludeTag != "" {
				report.Tag = instructions.Deanonymize.IncludeTag
			}

			// Serialize the report
			reportBytes, err := json.Marshal(report)
			if err != nil {
				return ProcessResult{Error: fmt.Sprintf("Failed to serialize deanonymization report: %v", err)}
			}

			// Serialize the updated state (state changes due to nonce increment handled by executor)
			newStateBytes, err := json.Marshal(currentState)
			if err != nil {
				return ProcessResult{Error: fmt.Sprintf("Failed to serialize new state: %v", err)}
			}
			return ProcessResult{
				State:  newStateBytes,
				Report: reportBytes,
				Fuel:   NewUint256(20),
			}

		default:
			utils.LogError("ProcessRequest: unsupported instruction type: %s", instructions.Type)
			return ProcessResult{Error: fmt.Sprintf("Unsupported instruction type: [%s]", instructionType)}
		}
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		utils.LogError("ProcessRequest: failed to serialize new state: %v", err)
		return ProcessResult{Error: fmt.Sprintf("Failed to serialize new state: %v", err)}
	}
	fuel := NewUint256(50)
	utils.LogDebug("ProcessRequest: sender=%s, eventsCount=%d, withdrawalsCount=%d, stateSize=%d, fuel=%v",
		senderHex, len(events), len(withdrawals), len(newStateBytes), fuel)
	return ProcessResult{
		State:       newStateBytes,
		Events:      events,
		Withdrawals: withdrawals,
		Fuel:        fuel,
	}
}

func GetAllocatedMemoryStats() MemoryStats {
	map_size, total_bytes := utils.GetAllocatedMemoryStats()
	return MemoryStats{
		MapSize:              map_size,
		CumulativeMemorySize: total_bytes,
	}
}
