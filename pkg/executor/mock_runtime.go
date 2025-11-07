package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
)

// MockRuntime implements a simple mock runtime that mimics a wasm application
// It supports deposits, fund transfers, withdrawals, and events with serialized state persistence
type MockRuntime struct{}

// NewMockRuntime creates a new mock runtime instance
func NewMockRuntime() *MockRuntime {
	log.Println("Initializing mock runtime")
	return &MockRuntime{}
}

// LoadModule loads a WASM module and returns initial state
func (r *MockRuntime) LoadModule(ctx context.Context, appId common.ApplicationIdType, wasm []byte) ([]byte, error) {
	log.Printf("Mock Runtime: Loading mock runtime module for application %d (wasm size: %d bytes)", appId, len(wasm))

	// Create initial application state (generic map-based representation)
	initialState := map[string]interface{}{
		"appId":    appId,
		"accounts": map[string]map[string]interface{}{},
		"nonce":    uint64(0),
	}

	stateBytes, err := json.Marshal(initialState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal initial state: %w", err)
	}

	log.Printf("Mock Runtime: Successfully loaded mock runtime module for application %d", appId)
	return stateBytes, nil
}

func (r *MockRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, value *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, error) {
	log.Printf("Mock Runtime: Processing deposit for application %d ( value: %d wei for sender: %s )", appId, value, sender)

	var currentState map[string]interface{}
	if err := json.Unmarshal(state, &currentState); err != nil {
		return nil, nil, fmt.Errorf("failed to deserialize state: %w", err)
	}

	accounts := ensureAccounts(currentState)
	nonce := ensureNonce(currentState)

	var events []common.PlainEvent
	if value.Sign() == 1 {
		// Ensure sender account exists
		acct := ensureAccount(accounts, sender)
		// Update balance
		balance := toUint64(acct["balance"]) + value.Uint64()
		acct["balance"] = balance
		// Increment nonce
		nonce++
		currentState["nonce"] = nonce

		depositEvent := common.PlainEvent{
			UserID: sender,
			Data:   []byte(fmt.Sprintf(`{"type":"deposit","amount":%d,"balance":%d,"nonce":%d}`, value, balance, nonce)),
		}
		events = append(events, depositEvent)
	}

	newSerializedState, err := json.Marshal(currentState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize new state: %w", err)
	}

	log.Printf("Mock Runtime: Successfully processed deposit for sender %s, generated %d events", sender, len(events))
	return newSerializedState, events, nil
}

// ProcessRequest processes a request and returns the new state, events, and withdrawals
func (r *MockRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, error) {
	log.Printf("Mock Runtime: Processing request for application %d (payload size: %d, state size: %d)", appId, len(payload), len(state))

	var currentState map[string]interface{}
	if err := json.Unmarshal(state, &currentState); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to deserialize state: %w", err)
	}

	accounts := ensureAccounts(currentState)
	nonce := ensureNonce(currentState)

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal

	if len(payload) > 0 {
		var instructions map[string]interface{}
		if err := json.Unmarshal(payload, &instructions); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal payload instructions: %w", err)
		}

		typ := instructions["type"].(string)
		switch typ {
		case "transfer":
			transfer := instructions["transfer"].(map[string]interface{})
			if transfer == nil {
				return nil, nil, nil, fmt.Errorf("transfer instruction is nil")
			}
			to := ethCommon.HexToAddress(transfer["to"].(string))
			amount := toUint64(transfer["amount"])

			// Ensure sender exists and has balance
			senderAcct := accounts[sender]
			if senderAcct == nil {
				return nil, nil, nil, fmt.Errorf("sender account %s does not exist", sender)
			}
			if toUint64(senderAcct["balance"]) < amount {
				return nil, nil, nil, fmt.Errorf("insufficient balance for transfer")
			}

			// Ensure recipient account
			recipientAcct := ensureAccount(accounts, to)

			// Execute transfer
			senderAcct["balance"] = toUint64(senderAcct["balance"]) - amount
			recipientAcct["balance"] = toUint64(recipientAcct["balance"]) + amount
			nonce++
			currentState["nonce"] = nonce

			// Events
			senderEvent := common.PlainEvent{
				UserID: sender,
				Data: []byte(fmt.Sprintf(`{"type":"transfer_sent","to":"%s","amount":%d,"balance":%d,"nonce":%d}`,
					to, amount, toUint64(senderAcct["balance"]), nonce)),
			}
			recipientEvent := common.PlainEvent{
				UserID: to,
				Data: []byte(fmt.Sprintf(`{"type":"transfer_received","from":"%s","amount":%d,"balance":%d,"nonce":%d}`,
					sender, amount, toUint64(recipientAcct["balance"]), nonce)),
			}
			events = append(events, senderEvent, recipientEvent)

		case "withdraw":
			withdraw := instructions["withdraw"].(map[string]interface{})
			if withdraw == nil {
				return nil, nil, nil, fmt.Errorf("withdraw instruction is nil")
			}
			to := ethCommon.HexToAddress(withdraw["to"].(string))
			amount := toUint64(withdraw["amount"])

			senderAcct := accounts[sender]
			if senderAcct == nil {
				return nil, nil, nil, fmt.Errorf("sender account %s does not exist", sender)
			}
			if toUint64(senderAcct["balance"]) < amount {
				return nil, nil, nil, fmt.Errorf("insufficient balance for withdrawal")
			}

			// Execute withdrawal
			senderAcct["balance"] = toUint64(senderAcct["balance"]) - amount
			nonce++
			currentState["nonce"] = nonce

			withdrawals = append(withdrawals, common.Withdrawal{DestinationAddress: to, Amount: new(big.Int).SetUint64(amount)})

			withdrawEvent := common.PlainEvent{
				UserID: sender,
				Data: []byte(fmt.Sprintf(`{"type":"withdrawal","to":"%s","amount":%d,"balance":%d,"nonce":%d}`,
					to, amount, toUint64(senderAcct["balance"]), nonce)),
			}
			events = append(events, withdrawEvent)

		default:
			return nil, nil, nil, fmt.Errorf("unknown instruction type: %s", typ)
		}
	}

	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to serialize new state: %w", err)
	}

	log.Printf("Mock Runtime: Successfully processed request for application %d, generated %d events and %d withdrawals", appId, len(events), len(withdrawals))
	return newStateBytes, events, withdrawals, nil
}

// GenerateDeanonymizationReport generates a deanonymization report
func (r *MockRuntime) GenerateDeanonymizationReport(ctx context.Context, appId common.ApplicationIdType, payload []byte, state []byte, wasm []byte) ([]byte, error) {
	log.Printf("Mock Runtime: Generating deanonymization report for application %d", appId)

	var currentState map[string]interface{}
	if err := json.Unmarshal(state, &currentState); err != nil {
		return nil, fmt.Errorf("failed to deserialize state for deanonymization: %w", err)
	}

	report := map[string]interface{}{
		"accounts": currentState["accounts"],
		"nonce":    currentState["nonce"],
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deanonymization report: %w", err)
	}

	log.Printf("Mock Runtime: Successfully generated deanonymization report for application %d", appId)
	return reportBytes, nil
}

// Close closes the mock runtime and cleans up resources
func (r *MockRuntime) Close() error {
	log.Printf("Mock Runtime: Closing mock runtime")
	log.Printf("Mock Runtime: Mock runtime closed successfully")
	return nil
}

// --- helpers ---

func ensureAccounts(state map[string]interface{}) map[ethCommon.Address]map[string]interface{} {
	accAny, ok := state["accounts"]
	if !ok || accAny == nil {
		m := map[ethCommon.Address]map[string]interface{}{}
		state["accounts"] = m
		return m
	}
	// try typed map
	if m, ok := accAny.(map[ethCommon.Address]map[string]interface{}); ok {
		return m
	}
	// convert from generic map[string]interface{}
	res := map[ethCommon.Address]map[string]interface{}{}
	if gm, ok := accAny.(map[string]interface{}); ok {
		for k, v := range gm {
			if sub, ok := v.(map[string]interface{}); ok {
				res[ethCommon.HexToAddress(k)] = sub
			}
		}
	}
	state["accounts"] = res
	return res
}

func ensureAccount(accounts map[ethCommon.Address]map[string]interface{}, addr ethCommon.Address) map[string]interface{} {
	acct := accounts[addr]
	if acct == nil {
		acct = map[string]interface{}{"address": addr, "balance": uint64(0)}
		accounts[addr] = acct
	}
	return acct
}

func ensureNonce(state map[string]interface{}) uint64 {
	if n, ok := state["nonce"]; ok {
		return toUint64(n)
	}
	state["nonce"] = uint64(0)
	return 0
}

func toUint64(v interface{}) uint64 {
	switch x := v.(type) {
	case uint64:
		return x
	case uint32:
		return uint64(x)
	case int:
		return uint64(x)
	case int64:
		return uint64(x)
	case float64:
		return uint64(x)
	default:
		return 0
	}
}
