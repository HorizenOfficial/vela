package executor

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/horizen-pes/pkg/common"
)

// Local mirror types to avoid importing the wasm-go app package
type AccountState struct {
	Address string `json:"address"`
	Balance uint64 `json:"balance"`
}

type ApplicationInternalState struct {
	AppID    string                   `json:"appId"`
	Accounts map[string]*AccountState `json:"accounts"`
	Nonce    uint64                   `json:"nonce"`
}

type TransferInstruction struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

type WithdrawInstruction struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

type PayloadInstructions struct {
	Type     string               `json:"type"`
	Transfer *TransferInstruction `json:"transfer,omitempty"`
	Withdraw *WithdrawInstruction `json:"withdraw,omitempty"`
}

// MockRuntime implements a simple mock runtime that mimics a wasm application
// It supports deposits, fund transfers, withdrawals, and events with serializeed state persistence
type MockRuntime struct {
}

// NewMockRuntime creates a new mock runtime instance
func NewMockRuntime() *MockRuntime {
	log.Println("Initializing mock runtime")
	return &MockRuntime{}
}

// LoadModule loads a WASM module and returns initial state and state root
func (r *MockRuntime) LoadModule(ctx context.Context, appId string, wasm []byte) ([]byte, [32]byte, error) {
	log.Printf("Mock Runtime: Loading mock runtime module for application %s (wasm size: %d bytes)", appId, len(wasm))

	// Create initial application state
	initialState := &ApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[string]*AccountState),
		Nonce:    0,
	}

	// Serialize the state
	stateBytes, err := json.Marshal(initialState)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("failed to marshal initial state: %w", err)
	}

	// Create state root hash
	stateRoot := sha256.Sum256(stateBytes)

	log.Printf("Mock Runtime: Successfully loaded mock runtime module for application %s", appId)
	return stateBytes, stateRoot, nil
}

func (r *MockRuntime) Deposit(ctx context.Context, appId string, sender string, value uint64, state []byte, wasm []byte) ([]byte, []common.PlainEvent, error) {
	log.Printf("Mock Runtime: Processing deposit for application %s ( value: %d wei for sender: %s )", appId, value, sender)

	// Deserialize the current state
	var currentState ApplicationInternalState
	err := json.Unmarshal(state, &currentState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize state: %w", err)
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

		// Create a deposit event for sender
		depositEvent := common.PlainEvent{
			UserID: sender,
			Data:   []byte(fmt.Sprintf(`{"type":"deposit","amount":%d,"balance":%d,"nonce":%d}`, value, currentState.Accounts[sender].Balance, currentState.Nonce)),
		}
		events = append(events, depositEvent)
	}

	// Serialize the updated state
	newSerializedState, err := json.Marshal(&currentState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize new state: %w", err)
	}

	log.Printf("Mock Runtime: Successfully processed deposit for sender %s, generated %d events", sender, len(events))
	return newSerializedState, events, nil
}

// ProcessRequest processes a request and returns the new state, events, and withdrawals
func (r *MockRuntime) ProcessRequest(ctx context.Context, appId string, sender string, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, error) {
	log.Printf("Mock Runtime: Processing request for application %s (payload size: %d, state size: %d)", appId, len(payload), len(state))

	// deserialize the current state
	var currentState ApplicationInternalState
	err := json.Unmarshal(state, &currentState)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to deserialize state: %w", err)
	}

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal

	// Process payload instructions if payload is not empty
	if len(payload) > 0 {
		var instructions PayloadInstructions
		if err := json.Unmarshal(payload, &instructions); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal payload instructions: %w", err)
		}

		switch instructions.Type {
		case "transfer":
			if instructions.Transfer == nil {
				return nil, nil, nil, fmt.Errorf("transfer instruction is nil")
			}

			log.Printf("Mock Runtime: Processing transfer from %s to %s of %d wei", sender, instructions.Transfer.To, instructions.Transfer.Amount)

			// Ensure sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return nil, nil, nil, fmt.Errorf("sender account %s does not exist", sender)
			}
			if currentState.Accounts[sender].Balance < instructions.Transfer.Amount {
				return nil, nil, nil, fmt.Errorf("insufficient balance for transfer")
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
			senderEvent := common.PlainEvent{
				UserID: sender,
				Data: []byte(fmt.Sprintf(`{"type":"transfer_sent","to":"%s","amount":%d,"balance":%d,"nonce":%d}`,
					instructions.Transfer.To, instructions.Transfer.Amount, currentState.Accounts[sender].Balance, currentState.Nonce)),
			}
			recipientEvent := common.PlainEvent{
				UserID: instructions.Transfer.To,
				Data: []byte(fmt.Sprintf(`{"type":"transfer_received","from":"%s","amount":%d,"balance":%d,"nonce":%d}`,
					sender, instructions.Transfer.Amount, currentState.Accounts[instructions.Transfer.To].Balance, currentState.Nonce)),
			}
			events = append(events, senderEvent, recipientEvent)

		case "withdraw":
			if instructions.Withdraw == nil {
				return nil, nil, nil, fmt.Errorf("withdraw instruction is nil")
			}

			log.Printf("Mock Runtime: Processing withdrawal from %s to %s of %d wei", sender, instructions.Withdraw.To, instructions.Withdraw.Amount)

			// Ensure sender account exists and has sufficient balance
			if currentState.Accounts[sender] == nil {
				return nil, nil, nil, fmt.Errorf("sender account %s does not exist", sender)
			}
			if currentState.Accounts[sender].Balance < instructions.Withdraw.Amount {
				return nil, nil, nil, fmt.Errorf("insufficient balance for withdrawal")
			}

			// Execute withdrawal
			currentState.Accounts[sender].Balance -= instructions.Withdraw.Amount
			currentState.Nonce++

			// Create withdrawal
			withdrawal := common.Withdrawal{
				DestinationAddress: instructions.Withdraw.To,
				Amount:             instructions.Withdraw.Amount,
			}
			withdrawals = append(withdrawals, withdrawal)

			// Create event for sender
			withdrawEvent := common.PlainEvent{
				UserID: sender,
				Data: []byte(fmt.Sprintf(`{"type":"withdrawal","to":"%s","amount":%d,"balance":%d,"nonce":%d}`,
					instructions.Withdraw.To, instructions.Withdraw.Amount, currentState.Accounts[sender].Balance, currentState.Nonce)),
			}
			events = append(events, withdrawEvent)

		default:
			return nil, nil, nil, fmt.Errorf("unknown instruction type: %s", instructions.Type)
		}
	}

	// Serialize the updated state
	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to serialize new state: %w", err)
	}

	log.Printf("Mock Runtime: Successfully processed request for application %s, generated %d events and %d withdrawals", appId, len(events), len(withdrawals))
	return newStateBytes, events, withdrawals, nil
}

// GenerateDeanonymizationReport generates a deanonymization report
func (r *MockRuntime) GenerateDeanonymizationReport(ctx context.Context, appId string, requestId string, payload []byte, state []byte, wasm []byte) ([]byte, error) {
	log.Printf("Mock Runtime: Generating deanonymization report, id: %s,for application %s", requestId, appId)

	// deserialize the current state to access account information
	var currentState ApplicationInternalState
	err := json.Unmarshal(state, &currentState)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize state for deanonymization: %w", err)
	}

	// Create a deanonymization report with full account information
	report := map[string]interface{}{
		"applicationId": appId,
		"requestId":     requestId,
		"accounts":      currentState.Accounts,
		"nonce":         currentState.Nonce,
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deanonymization report: %w", err)
	}

	log.Printf("Mock Runtime: Successfully generated deanonymization report for application %s", appId)
	return reportBytes, nil
}

// Close closes the mock runtime and cleans up resources
func (r *MockRuntime) Close() error {
	log.Printf("Mock Runtime: Closing mock runtime")
	// In a real implementation, this would clean up resources, close connections, etc.
	log.Printf("Mock Runtime: Mock runtime closed successfully")
	return nil
}
