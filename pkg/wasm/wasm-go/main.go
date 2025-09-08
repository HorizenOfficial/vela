package main

import (
	"encoding/json"
	"fmt"

	"payment-app/utils"

	appCommon "github.com/horizen-pes/pkg/wasm/common"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/appstate"
)

// --- WASM-Exposed Functions (Bridge to Application Logic) ---

// These functions handle the WASM I/O and call the high-level logic functions.

//export load_module
func load_module(appIdPtr *byte, appIdLen int32) *byte {
	appId := utils.PtrToString(appIdPtr, appIdLen)
	stateBytes := loadModule(appId)
	return utils.StringToPtr(stateBytes)
}

//export deposit
func deposit(appIdPtr *byte, appIdLen int32, senderPtr *byte, senderLen int32, value uint64, statePtr *byte, stateLen int32) *byte {
	// TODO: in future we must use the appId for adding it to the generated event
	_ = utils.PtrToString(appIdPtr, appIdLen)
	sender := utils.PtrToString(senderPtr, senderLen)
	stateJSON := utils.PtrToString(statePtr, stateLen)
	result := depositFunds(sender, value, stateJSON)
	return utils.SerializeAndWriteResult(result)
}

//export process_request
func process_request(appIdPtr *byte, appIdLen int32, senderPtr *byte, senderLen int32, payloadPtr *byte, payloadLen int32, statePtr *byte, stateLen int32) *byte {
	// TODO: in future we must use the appId for setting it in the generated event
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
	initialState := &appstate.ApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[string]*appstate.AccountState),
		Nonce:    0,
	}
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return utils.WasmSerializationError
	}
	return stateJSON
}

func depositFunds(sender string, value uint64, stateJSON string) appCommon.DepositResult {
	var currentState appstate.ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return appCommon.DepositResult{Error: "Failed to parse application state"}
	}
	if !utils.IsValidAddress(sender) {
		return appCommon.DepositResult{Error: fmt.Sprintf("sender address is not valid: %s", sender)}
	}

	var events []common.PlainEvent

	// Handle deposit
	if value > 0 {
		// Ensure sender account exists
		if currentState.Accounts[sender] == nil {
			currentState.Accounts[sender] = &appstate.AccountState{
				Address: sender,
				Balance: 0,
			}
		}

		// Add deposit to sender's balance
		currentState.Accounts[sender].Balance += value
		currentState.Nonce++

		// Create deposit event
		eventData := appCommon.DepositEvent{
			Type:    "deposit",
			Amount:  value,
			Balance: currentState.Accounts[sender].Balance,
			Nonce:   currentState.Nonce,
		}
		eventDataBytes, err := json.Marshal(eventData)
		if err != nil {
			return appCommon.DepositResult{Error: "Failed to serialize event data"}
		}

		events = append(events, common.PlainEvent{
			UserID: sender,
			Data:   eventDataBytes,
		})
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(&currentState)
	if err != nil {
		return appCommon.DepositResult{Error: "Failed to serialize new state"}
	}
	return appCommon.DepositResult{State: newStateBytes, Events: events}
}

func processRequest(sender, payloadJSON, stateJSON string) appCommon.ProcessResult {
	// Deserialize current state
	var currentState appstate.ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return appCommon.ProcessResult{Error: "Failed to parse application state"}
	}
	if !utils.IsValidAddress(sender) {
		return appCommon.ProcessResult{Error: fmt.Sprintf("sender address is not valid: %s", sender)}
	}

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal

	// Process payload instructions if payload is not empty
	if payloadJSON != "" {
		var instructions appstate.PayloadInstructions
		if err := json.Unmarshal([]byte(payloadJSON), &instructions); err != nil {
			return appCommon.ProcessResult{Error: "Failed to parse payload instructions"}
		}

		switch instructions.Type {
		case "transfer":
			if instructions.Transfer == nil {
				return appCommon.ProcessResult{Error: "Transfer instruction is missing"}
			}

			if !utils.IsValidAddress(instructions.Transfer.To) {
				return appCommon.ProcessResult{Error: fmt.Sprintf("Transfer destination addrss is not valid: %s", instructions.Transfer.To)}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return appCommon.ProcessResult{Error: fmt.Sprintf("Account does not exist: %s", sender)}
			}
			if currentState.Accounts[sender].Balance < instructions.Transfer.Amount {
				return appCommon.ProcessResult{Error: "Insufficient balance for transfer"}
			}

			// Ensure recipient account exists
			if currentState.Accounts[instructions.Transfer.To] == nil {
				currentState.Accounts[instructions.Transfer.To] = &appstate.AccountState{
					Address: instructions.Transfer.To,
					Balance: 0,
				}
			}

			// Execute transfer
			currentState.Accounts[sender].Balance -= instructions.Transfer.Amount
			currentState.Accounts[instructions.Transfer.To].Balance += instructions.Transfer.Amount
			currentState.Nonce++

			// Create events for both parties
			senderEventData := appCommon.SenderEvent{
				Type:    "transfer_sent",
				To:      instructions.Transfer.To,
				Amount:  instructions.Transfer.Amount,
				Balance: currentState.Accounts[sender].Balance,
				Nonce:   currentState.Nonce,
			}
			senderEventDataBytes, err := json.Marshal(senderEventData)
			if err != nil {
				return appCommon.ProcessResult{Error: "Failed to serialize sender event data"}
			}

			recipientEventData := appCommon.RecipientEvent{
				Type:    "transfer_received",
				From:    sender,
				Amount:  instructions.Transfer.Amount,
				Balance: currentState.Accounts[instructions.Transfer.To].Balance,
				Nonce:   currentState.Nonce,
			}
			recipientEventDataBytes, err := json.Marshal(recipientEventData)
			if err != nil {
				return appCommon.ProcessResult{Error: "Failed to serialize recipient event data"}
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
				return appCommon.ProcessResult{Error: "Withdraw instruction is missing"}
			}

			// Validate sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return appCommon.ProcessResult{Error: "Account does not exist"}
			}

			if currentState.Accounts[sender].Balance < instructions.Withdraw.Amount {
				return appCommon.ProcessResult{Error: "Insufficient balance for withdrawal"}
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
			withdrawEventData := appCommon.WithdrawalEvent{
				Type:    "withdrawal",
				To:      instructions.Withdraw.To,
				Amount:  instructions.Withdraw.Amount,
				Balance: currentState.Accounts[sender].Balance,
				Nonce:   currentState.Nonce,
			}
			withdrawEventDataBytes, err := json.Marshal(withdrawEventData)
			if err != nil {
				return appCommon.ProcessResult{Error: "Failed to serialize withdraw event data"}
			}

			events = append(events, common.PlainEvent{
				UserID: sender,
				Data:   withdrawEventDataBytes,
			})

		default:
			return appCommon.ProcessResult{Error: "Unsupported instruction type"}
		}
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		return appCommon.ProcessResult{Error: "Failed to serialize new state"}
	}
	return appCommon.ProcessResult{
		State:       newStateBytes,
		Events:      events,
		Withdrawals: withdrawals,
	}
}

func generateDeanonymizationReport(appId, requestId, stateJSON string) appCommon.DeanonymizationResult {
	// Deserialize current state
	var currentState appstate.ApplicationInternalState
	if err := json.Unmarshal([]byte(stateJSON), &currentState); err != nil {
		return appCommon.DeanonymizationResult{Error: "Failed to parse application state"}
	}

	// Create deanonymization report
	report := appCommon.UnencryptedDeanonimizationReportData{
		ApplicationID: appId,
		RequestID:     requestId,
		Accounts:      currentState.Accounts,
		Nonce:         currentState.Nonce,
	}

	// Serialize the report
	reportBytes, err := json.Marshal(report)
	if err != nil {
		return appCommon.DeanonymizationResult{Error: "Failed to serialize deanonymization report"}
	}
	return appCommon.DeanonymizationResult{Report: reportBytes}
}

// Main function is required but not used in WASM
func main() {}
