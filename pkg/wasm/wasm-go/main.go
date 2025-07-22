package main

import (
	"encoding/json"
	"strconv"
	"unsafe"
)

const WasmSerializationError = "{}"

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

// Memory management functions
//
//export allocate
func allocate(size int32) *byte {
	// Allocate memory and return pointer
	data := make([]byte, size)
	return &data[0]
}

//export deallocate
func deallocate(ptr *byte, size int32) {
	// no-op for deallocation in Go, as Go's garbage collector handles memory management
}

// Helper function to convert pointer and length to string
func ptrToString(ptr *byte, length int32) string {
	if ptr == nil || length == 0 {
		return ""
	}
	return string(unsafe.Slice(ptr, length))
}

// Helper function to convert string to allocated memory pointer
func stringToPtr(s string) *byte {
	if s == "" {
		return nil
	}
	data := []byte(s)
	data = append(data, 0x00) // null-terminate the string
	ptr := allocate(int32(len(data)))
	copy(unsafe.Slice(ptr, len(data)), data)
	return ptr
}

//export load_module
func load_module(appIdPtr *byte, appIdLen int32) *byte {
	appId := ptrToString(appIdPtr, appIdLen)

	initialState := &ApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[string]*AccountState),
		Nonce:    0,
	}

	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return stringToPtr(WasmSerializationError)
	}

	return stringToPtr(string(stateJSON))
}

//export deposit
func deposit(appIdPtr *byte, appIdLen int32, senderPtr *byte, senderLen int32, value int64, statePtr *byte, stateLen int32) *byte {
	_ = ptrToString(appIdPtr, appIdLen) // appId not used in this function
	sender := ptrToString(senderPtr, senderLen)
	stateJSON := ptrToString(statePtr, stateLen)

	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		result := DepositResult{
			Error: "Failed to parse application state",
		}
		return serializeAndWriteResult(result)
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
		result := DepositResult{
			Error: "Failed to serialize new state",
		}
		return serializeAndWriteResult(result)
	}

	result := DepositResult{
		State:  newStateBytes,
		Events: events,
	}

	return serializeAndWriteResult(result)
}

//export process_request
func process_request(appIdPtr *byte, appIdLen int32, senderPtr *byte, senderLen int32, payloadPtr *byte, payloadLen int32, statePtr *byte, stateLen int32) *byte {
	_ = ptrToString(appIdPtr, appIdLen) // appId not used in this function
	sender := ptrToString(senderPtr, senderLen)
	payloadJSON := ptrToString(payloadPtr, payloadLen)
	stateJSON := ptrToString(statePtr, stateLen)

	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		result := ProcessResult{
			Error: "Failed to parse application state",
		}
		return serializeAndWriteResult(result)
	}

	var events []PlainEvent
	var withdrawals []Withdrawal

	// Process payload instructions if payload is not empty
	if payloadJSON != "" {
		var instructions PayloadInstructions
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			result := ProcessResult{
				Error: "Failed to parse payload instructions",
			}
			return serializeAndWriteResult(result)
		}

		switch instructions.Type {
		case "transfer":
			if instructions.Transfer == nil {
				result := ProcessResult{
					Error: "Transfer instruction is missing",
				}
				return serializeAndWriteResult(result)
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				result := ProcessResult{
					Error: "Account does not exist",
				}
				return serializeAndWriteResult(result)
			}
			if currentState.Accounts[sender].Balance < instructions.Transfer.Amount {
				result := ProcessResult{
					Error: "Insufficient balance for transfer",
				}
				return serializeAndWriteResult(result)
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
				result := ProcessResult{
					Error: "Withdraw instruction is missing",
				}
				return serializeAndWriteResult(result)
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				result := ProcessResult{
					Error: "Account does not exist",
				}
				return serializeAndWriteResult(result)
			}
			if currentState.Accounts[sender].Balance < instructions.Withdraw.Amount {
				result := ProcessResult{
					Error: "Insufficient balance for withdrawal",
				}
				return serializeAndWriteResult(result)
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
			result := ProcessResult{
				Error: "Unsupported instruction type",
			}
			return serializeAndWriteResult(result)
		}
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		result := ProcessResult{
			Error: "Failed to serialize new state",
		}
		return serializeAndWriteResult(result)
	}

	result := ProcessResult{
		State:       newStateBytes,
		Events:      events,
		Withdrawals: withdrawals,
	}

	return serializeAndWriteResult(result)
}

//export generate_deanonymization_report
func generate_deanonymization_report(appIdPtr *byte, appIdLen int32, requestIdPtr *byte, requestIdLen int32, statePtr *byte, stateLen int32) *byte {
	appId := ptrToString(appIdPtr, appIdLen)
	requestId := ptrToString(requestIdPtr, requestIdLen)
	stateJSON := ptrToString(statePtr, stateLen)

	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		result := DeanonymizationResult{
			Error: "Failed to parse application state",
		}
		return serializeAndWriteResult(result)
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
		result := DeanonymizationResult{
			Error: "Failed to serialize deanonymization report",
		}
		return serializeAndWriteResult(result)
	}

	result := DeanonymizationResult{
		Report: reportBytes,
	}
	return serializeAndWriteResult(result)
}

func serializeAndWriteResult(result any) *byte {
	reportJSON, err := json.Marshal(result)
	if err != nil {
		return stringToPtr(WasmSerializationError)
	}
	return stringToPtr(string(reportJSON))
}

// Main function is required but not used in WASM
func main() {}
