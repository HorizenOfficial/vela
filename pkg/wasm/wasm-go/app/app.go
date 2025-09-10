package app

import (
	"encoding/json"
	"fmt"
	"payment-app/utils"

	"github.com/horizen-pes/pkg/common"
	wasmCommon "github.com/horizen-pes/pkg/wasm/common"
)

// --- High-Level Application Logic ---

func LoadModule(appId string) []byte {
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

func DepositFunds(sender string, value uint64, stateJSON string) wasmCommon.DepositResult {
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return wasmCommon.DepositResult{Error: "Failed to parse application state"}
	}
	if !utils.IsValidAddress(sender) {
		return wasmCommon.DepositResult{Error: fmt.Sprintf("sender address is not valid: %s", sender)}
	}

	var events []common.PlainEvent

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
		eventData := wasmCommon.DepositEvent{
			Type:    "deposit",
			Amount:  value,
			Balance: currentState.Accounts[sender].Balance,
			Nonce:   currentState.Nonce,
		}
		eventDataBytes, err := json.Marshal(eventData)
		if err != nil {
			return wasmCommon.DepositResult{Error: "Failed to serialize event data"}
		}

		events = append(events, common.PlainEvent{
			UserID: sender,
			Data:   eventDataBytes,
		})
	}

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
	if !utils.IsValidAddress(sender) {
		return wasmCommon.ProcessResult{Error: fmt.Sprintf("sender address is not valid: %s", sender)}
	}

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal

	// Process payload instructions if payload is not empty
	if payloadJSON != "" {
		var instructions PayloadInstructions
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			return wasmCommon.ProcessResult{Error: "Failed to parse payload instructions"}
		}

		switch instructions.Type {
		case "transfer":
			if instructions.Transfer == nil {
				return wasmCommon.ProcessResult{Error: "Transfer instruction is missing"}
			}

			if !utils.IsValidAddress(instructions.Transfer.To) {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Transfer destination addrss is not valid: %s", instructions.Transfer.To)}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return wasmCommon.ProcessResult{Error: fmt.Sprintf("Account does not exist: %s", sender)}
			}
			if currentState.Accounts[sender].Balance < instructions.Transfer.Amount {
				return wasmCommon.ProcessResult{Error: "Insufficient balance for transfer"}
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
			senderEventData := wasmCommon.SenderEvent{
				Type:    "transfer_sent",
				To:      instructions.Transfer.To,
				Amount:  instructions.Transfer.Amount,
				Balance: currentState.Accounts[sender].Balance,
				Nonce:   currentState.Nonce,
			}
			senderEventDataBytes, err := json.Marshal(senderEventData)
			if err != nil {
				return wasmCommon.ProcessResult{Error: "Failed to serialize sender event data"}
			}

			recipientEventData := wasmCommon.RecipientEvent{
				Type:    "transfer_received",
				From:    sender,
				Amount:  instructions.Transfer.Amount,
				Balance: currentState.Accounts[instructions.Transfer.To].Balance,
				Nonce:   currentState.Nonce,
			}
			recipientEventDataBytes, err := json.Marshal(recipientEventData)
			if err != nil {
				return wasmCommon.ProcessResult{Error: "Failed to serialize recipient event data"}
			}

			events = append(events, common.PlainEvent{
				UserID: sender,
				Data:   senderEventDataBytes,
			})

			events = append(events, common.PlainEvent{
				UserID: instructions.Transfer.To,
				Data:   recipientEventDataBytes,
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
			currentState.Nonce++

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
				Nonce:   currentState.Nonce,
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
			return wasmCommon.ProcessResult{Error: "Unsupported instruction type"}
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

func GenerateDeanonymizationReport(appId, requestId, stateJSON string) wasmCommon.DeanonymizationResult {
	// Deserialize current state
	var currentState ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return wasmCommon.DeanonymizationResult{Error: "Failed to parse application state"}
	}

	// Create deanonymization report
	report := wasmCommon.UnencryptedDeanonimizationReportData{
		ApplicationID: appId,
		RequestID:     requestId,
		Accounts:      convertAccountsToWasmCommon(currentState.Accounts),
		Nonce:         currentState.Nonce,
	}

	// Serialize the report
	reportBytes, err := json.Marshal(report)
	if err != nil {
		return wasmCommon.DeanonymizationResult{Error: "Failed to serialize deanonymization report"}
	}
	return wasmCommon.DeanonymizationResult{Report: reportBytes}
}

// convertAccountsToWasmCommon converts app-level AccountState map to wasm/common.AccountState map
func convertAccountsToWasmCommon(src map[string]*AccountState) map[string]*wasmCommon.AccountState {
	if src == nil {
		return nil
	}
	out := make(map[string]*wasmCommon.AccountState, len(src))
	for addr, acc := range src {
		if acc == nil {
			out[addr] = nil
			continue
		}
		out[addr] = &wasmCommon.AccountState{
			Address: acc.Address,
			Balance: acc.Balance,
		}
	}
	return out
}

// AccountState represents the state of a user account
type AccountState struct {
	Address string `json:"address"`
	Balance uint64 `json:"balance"`
}

// ApplicationInternalState represents the internal state of the application
type ApplicationInternalState struct {
	AppID    string                   `json:"appId"`
	Accounts map[string]*AccountState `json:"accounts"`
	Nonce    uint64                   `json:"nonce"`
}

// TransferInstruction represents instructions for transferring funds
type TransferInstruction struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

// WithdrawInstruction represents instructions for withdrawing funds
type WithdrawInstruction struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

// PayloadInstructions represents the deserialized payload instructions
type PayloadInstructions struct {
	Type     string               `json:"type"`
	Transfer *TransferInstruction `json:"transfer,omitempty"`
	Withdraw *WithdrawInstruction `json:"withdraw,omitempty"`
}
