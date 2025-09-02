package main

import (
	"encoding/json"
	"strconv"

	"payment-app/utils"
)


// AccountState represents the state of a user account
type AccountState struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
}

// ApplicationInternalState represents the internal state of the application
type ApplicationInternalState struct {
	AppID    string                   `json:"appId"`
	Accounts map[string]*AccountState `json:"accounts"`
	Nonce    int64                    `json:"nonce"`
}

// TransferInstruction represents instructions for transferring funds
type TransferInstruction struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

// WithdrawInstruction represents instructions for withdrawing funds
type WithdrawInstruction struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

// PayloadInstructions represents the deserialized payload instructions
type PayloadInstructions struct {
	Type     string               `json:"type"`
	Transfer *TransferInstruction `json:"transfer,omitempty"`
	Withdraw *WithdrawInstruction `json:"withdraw,omitempty"`
}

// PlainEvent represents an emitted event
type PlainEvent struct {
	UserID string `json:"userId"`
	Data   []byte `json:"data"`
}

// Withdrawal represents a withdrawal instruction
type Withdrawal struct {
	DestinationAddress string `json:"destinationAddress"`
	Amount             string `json:"amount"`
}

// DepositResult represents the result of a deposit operation
type DepositResult struct {
	State  []byte       `json:"state"`
	Events []PlainEvent `json:"events"`
	Error  string       `json:"error,omitempty"`
}

// ProcessResult represents the result of a process request operation
type ProcessResult struct {
	State       []byte       `json:"state"`
	Events      []PlainEvent `json:"events"`
	Withdrawals []Withdrawal `json:"withdrawals"`
	Error       string       `json:"error,omitempty"`
}

// DeanonymizationResult represents the result of generating deanonymization report
type DeanonymizationResult struct {
	Report []byte `json:"report"`
	Error  string `json:"error,omitempty"`
}

// --- WASM-Exposed Functions (Bridge to Application Logic) ---

// These functions handle the WASM I/O and call the high-level logic functions.

//export load_module
func load_module(appIdPtr *byte, appIdLen int32) *byte {
	appId := utils.PtrToString(appIdPtr, appIdLen)
	stateBytes := loadModule(appId)
	return utils.StringToPtr(stateBytes)
}

//export deposit
func deposit(appIdPtr *byte, appIdLen int32, senderPtr *byte, senderLen int32, value int64, statePtr *byte, stateLen int32) *byte {
	_ = utils.PtrToString(appIdPtr, appIdLen)
	sender := utils.PtrToString(senderPtr, senderLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := depositFunds(sender, value, stateJSON)
	return utils.SerializeAndWriteResult(result)
}

//export process_request
func process_request(appIdPtr *byte, appIdLen int32, senderPtr *byte, senderLen int32, payloadPtr *byte, payloadLen int32, statePtr *byte, stateLen int32) *byte {
	_ = utils.PtrToString(appIdPtr, appIdLen)
	sender := utils.PtrToString(senderPtr, senderLen)
	payloadJSON := utils.PtrToString(payloadPtr, payloadLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := processRequest(sender, payloadJSON, stateJSON)
	return utils.SerializeAndWriteResult(result)
}

//export generate_deanonymization_report
func generate_deanonymization_report(appIdPtr *byte, appIdLen int32, requestIdPtr *byte, requestIdLen int32, statePtr *byte, stateLen int32) *byte {
	appId := utils.PtrToString(appIdPtr, appIdLen)
	requestId := utils.PtrToString(requestIdPtr, requestIdLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := generateDeanonymizationReport(appId, requestId, stateJSON)
	return utils.SerializeAndWriteResult(result)
}

// --- High-Level Application Logic ---

func loadModule(appId string) []byte {
	initialState := &ApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[string]*AccountState),
		Nonce:    0,
	}
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return utils.WasmSerializationError
	}
	return stateJSON
}

func depositFunds(sender string, value int64, stateJSON string) DepositResult {
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return DepositResult{Error: "Failed to parse application state"}
	}

	var events []PlainEvent

	// Handle deposit
	if value > 0 {
		// Ensure sender account exists
		if currentState.Accounts[sender] == nil {
			currentState.Accounts[sender] = &AccountState{
				Address: sender,
				Balance: 0,
			}
		}

		// Add deposit to sender's balance
		currentState.Accounts[sender].Balance += value
		currentState.Nonce++

		// Create deposit event
		eventData := map[string]interface{}{
			"type":    "deposit",
			"amount":  value,
			"balance": currentState.Accounts[sender].Balance,
			"nonce":   currentState.Nonce,
		}
		eventDataBytes, _ := json.Marshal(eventData)

		events = append(events, PlainEvent{
			UserID: sender,
			Data:   eventDataBytes,
		})
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		return DepositResult{Error: "Failed to serialize new state"}
	}
	return DepositResult{State: newStateBytes, Events: events}
}

func processRequest(sender, payloadJSON, stateJSON string) ProcessResult {
	// Deserialize current state
        var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return ProcessResult{Error: "Failed to parse application state"}
	}

	var events []PlainEvent
	var withdrawals []Withdrawal

	// Process payload instructions if payload is not empty
	if payloadJSON != "" {
		var instructions PayloadInstructions
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			return ProcessResult{Error: "Failed to parse payload instructions"}
		}

		switch instructions.Type {
		case "transfer":
			if instructions.Transfer == nil {
				return ProcessResult{Error: "Transfer instruction is missing"}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return ProcessResult{Error: "Account does not exist"}
			}
			if currentState.Accounts[sender].Balance < instructions.Transfer.Amount {
				return ProcessResult{Error: "Insufficient balance for transfer"}
			}

			// Ensure recipient account exists
			if currentState.Accounts[instructions.Transfer.To] == nil {
				currentState.Accounts[instructions.Transfer.To] = &AccountState{
					Address: instructions.Transfer.To,
					Balance: 0,
				}
			}

			// Execute transfer
			currentState.Accounts[sender].Balance -= instructions.Transfer.Amount
			currentState.Accounts[instructions.Transfer.To].Balance += instructions.Transfer.Amount
			currentState.Nonce++

			// Create events for both parties
			senderEventData := map[string]interface{}{
				"type":    "transfer_sent",
				"to":      instructions.Transfer.To,
				"amount":  instructions.Transfer.Amount,
				"balance": currentState.Accounts[sender].Balance,
				"nonce":   currentState.Nonce,
			}
			senderEventDataBytes, _ := json.Marshal(senderEventData)

			recipientEventData := map[string]interface{}{
				"type":    "transfer_received",
				"from":    sender,
				"amount":  instructions.Transfer.Amount,
				"balance": currentState.Accounts[instructions.Transfer.To].Balance,
				"nonce":   currentState.Nonce,
			}
			recipientEventDataBytes, _ := json.Marshal(recipientEventData)

			events = append(events, PlainEvent{
				UserID: sender,
				Data:   senderEventDataBytes,
			})

			events = append(events, PlainEvent{
				UserID: instructions.Transfer.To,
				Data:   recipientEventDataBytes,
			})

		case "withdraw":
			if instructions.Withdraw == nil {
				return ProcessResult{Error: "Withdraw instruction is missing"}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return ProcessResult{Error: "Account does not exist"}
			}
			if currentState.Accounts[sender].Balance < instructions.Withdraw.Amount {
				return ProcessResult{Error: "Insufficient balance for withdrawal"}
			}

			// Execute withdrawal
			currentState.Accounts[sender].Balance -= instructions.Withdraw.Amount
			currentState.Nonce++

			// Create withdrawal
			withdrawals = append(withdrawals, Withdrawal{
				DestinationAddress: instructions.Withdraw.To,
				Amount:             strconv.FormatInt(instructions.Withdraw.Amount, 10),
			})

			// Create event for sender
			withdrawEventData := map[string]interface{}{
				"type":    "withdrawal",
				"to":      instructions.Withdraw.To,
				"amount":  instructions.Withdraw.Amount,
				"balance": currentState.Accounts[sender].Balance,
				"nonce":   currentState.Nonce,
			}
			withdrawEventDataBytes, _ := json.Marshal(withdrawEventData)

			events = append(events, PlainEvent{
				UserID: sender,
				Data:   withdrawEventDataBytes,
			})

		default:
			return ProcessResult{Error: "Unsupported instruction type"}
		}
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		return ProcessResult{Error: "Failed to serialize new state"}
	}
	return ProcessResult{
		State:       newStateBytes,
		Events:      events,
		Withdrawals: withdrawals,
	}
}

func generateDeanonymizationReport(appId, requestId, stateJSON string) DeanonymizationResult {
	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return DeanonymizationResult{Error: "Failed to parse application state"}
	}

	// Create deanonymization report
	report := map[string]interface{}{
		"applicationId": appId,
		"requestId":     requestId,
		"accounts":      currentState.Accounts,
		"nonce":         currentState.Nonce,
	}

	// Serialize the report
	reportBytes, err := json.Marshal(report)
	if err != nil {
		return DeanonymizationResult{Error: "Failed to serialize deanonymization report"}
	}
	return DeanonymizationResult{Report: reportBytes}
}

// Main function is required but not used in WASM
func main() {}
