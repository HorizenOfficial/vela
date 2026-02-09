package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	"github.com/horizen-pes/pkg/logger"
)

// Local mirror types used in tests to avoid importing wasm-go/app

type testAccountState struct {
	Address ethCommon.Address `json:"address"`
	Balance *common.Big       `json:"balance"`
}

type testApplicationInternalState struct {
	AppID    common.ApplicationIdType                `json:"appId"`
	Accounts map[ethCommon.Address]*testAccountState `json:"accounts"`
	Nonce    uint64                                  `json:"nonce"`
}

type testTransferInstruction struct {
	To     ethCommon.Address `json:"to"`
	Amount *common.Big       `json:"amount"`
}

type testWithdrawInstruction struct {
	To     ethCommon.Address `json:"to"`
	Amount *common.Big       `json:"amount"`
}

type testDeanonymizeInstruction struct {
	// A dummy tag for the deanonymization request
	Tag string `json:"tag"`
}

type testPayloadInstructions struct {
	Type        string                      `json:"type"`
	Transfer    *testTransferInstruction    `json:"transfer,omitempty"`
	Withdraw    *testWithdrawInstruction    `json:"withdraw,omitempty"`
	Deanonymize *testDeanonymizeInstruction `json:"deanonymize,omitempty"`
}

// MockRuntime implements a simple mock runtime that mimics a wasm application
// It supports deposits, fund transfers, withdrawals, and events with serialized state persistence

type MockRuntime struct {
	fuel *big.Int
	log  logger.Logger
}

// NewMockRuntime creates a new mock runtime instance
func NewMockRuntime(log logger.Logger) *MockRuntime {
	log.Info("Initializing mock runtime")
	return &MockRuntime{fuel: big.NewInt(10), log: log}
}

// LoadModule loads a WASM module and returns initial state
func (r *MockRuntime) LoadModule(ctx context.Context, appId common.ApplicationIdType, wasm []byte) ([]byte, *big.Int, error) {
	r.log.Info("Mock Runtime: Loading mock runtime module for application %d (wasm size: %d bytes)", appId, len(wasm))

	initialState := &testApplicationInternalState{
		AppID:    appId,
		Accounts: make(map[ethCommon.Address]*testAccountState),
		Nonce:    0,
	}
	stateBytes, err := json.Marshal(initialState)
	if err != nil {
		return nil, r.fuel, fmt.Errorf("failed to marshal initial state: %w", err)
	}

	r.log.Info("Mock Runtime: Successfully loaded mock runtime module for application %d", appId)
	return stateBytes, r.fuel, nil
}

func (r *MockRuntime) Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, depositAmount *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, *big.Int, *apperrors.RequestFailure) {
	r.log.Info("Mock Runtime: Processing deposit for application %d ( value: %s wei for sender: %s )", appId, depositAmount.String(), sender)

	var currentState testApplicationInternalState
	if err := json.Unmarshal(state, &currentState); err != nil {
		return nil, nil, r.fuel, apperrors.New(apperrors.CodeJsonUnmarshalError, "failed to deserialize state", err)
	}

	accounts := currentState.Accounts
	nonce := currentState.Nonce

	var events []common.PlainEvent
	if depositAmount.Sign() == 1 {
		// Ensure sender account exists
		acct := ensureAccount(accounts, sender)
		// Update balance
		balance := new(big.Int).Add(acct.Balance.ToInt(), depositAmount)
		acct.Balance = common.ToBig(balance)
		// Increment nonce
		nonce++
		currentState.Nonce = nonce

		depositEvent := common.PlainEvent{
			UserID:       sender,
			EventSubType: "deposit",
			Data:         []byte(fmt.Sprintf(`{"type":"deposit","amount":"0x%s","balance":"0x%s","nonce":%d}`, depositAmount.Text(16), balance.Text(16), nonce)),
		}
		events = append(events, depositEvent)
	}

	newSerializedState, err := json.Marshal(currentState)
	if err != nil {
		return nil, nil, r.fuel, apperrors.New(apperrors.CodeJsonMarshalError, "failed to serialize new state", err)
	}

	r.log.Info("Mock Runtime: Successfully processed deposit for sender %s, generated %d events", sender, len(events))
	return newSerializedState, events, r.fuel, nil
}

// ProcessRequest processes a request and returns the new state, events, withdrawals, and optionally a deanonymization report
func (r *MockRuntime) ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure) {
	r.log.Info("Mock Runtime: Processing request for application %d (type: %s, payload size: %d, state size: %d)", appId, requestType, len(payload), len(state))

	var currentState testApplicationInternalState
	if err := json.Unmarshal(state, &currentState); err != nil {
		return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeJsonUnmarshalError, "failed to deserialize state", err)
	}

	accounts := currentState.Accounts
	nonce := currentState.Nonce

	var events []common.PlainEvent
	var withdrawals []common.Withdrawal
	var report []byte

	// Determine the instruction type from payload or request type
	var typ string
	var instructions testPayloadInstructions
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &instructions); err != nil {
			return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeJsonUnmarshalError, "failed to unmarshal payload instructions", err)
		}
		typ = instructions.Type
	}

	// If payload type is empty but request type is Deanonymize, use that
	if typ == "" && requestType == common.Deanonymize {
		typ = "deanonymize"
	}

	if typ != "" {
		switch typ {
		case "transfer":
			transfer := instructions.Transfer
			if transfer == nil {
				return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeRequestFuncFailed, "transfer instruction is nil", nil)
			}
			to := transfer.To
			amount := transfer.Amount

			// Ensure sender exists and has balance
			senderAcct := accounts[sender]
			if senderAcct == nil {
				return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("sender account %s does not exist", sender), nil)
			}
			if senderAcct.Balance.ToInt().Cmp(amount.ToInt()) < 0 {
				return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("sender account %s has insufficient balance", sender), nil)
			}

			// Ensure recipient account
			recipientAcct := ensureAccount(accounts, to)

			// Execute transfer
			senderAcct.Balance = common.ToBig(new(big.Int).Sub(senderAcct.Balance.ToInt(), amount.ToInt()))
			recipientAcct.Balance = common.ToBig(new(big.Int).Add(recipientAcct.Balance.ToInt(), amount.ToInt()))
			nonce++
			currentState.Nonce = nonce

			// Events
			senderEvent := common.PlainEvent{
				UserID:       sender,
				EventSubType: "transfer_sent",
				Data: []byte(fmt.Sprintf(`{"type":"transfer_sent","to":"%s","amount":"0x%s","balance":"0x%s","nonce":%d}`,
					to, amount.ToInt().Text(16), senderAcct.Balance.ToInt().Text(16), nonce)),
			}
			recipientEvent := common.PlainEvent{
				UserID:       to,
				EventSubType: "transfer_received",
				Data: []byte(fmt.Sprintf(`{"type":"transfer_received","from":"%s","amount":"0x%s","balance":"0x%s","nonce":%d}`,
					sender, amount.ToInt().Text(16), recipientAcct.Balance.ToInt().Text(16), nonce)),
			}
			events = append(events, senderEvent, recipientEvent)

		case "withdraw":
			withdraw := instructions.Withdraw
			if withdraw == nil {
				return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeRequestFuncFailed, "withdraw instruction is nil", nil)
			}
			to := withdraw.To
			amount := withdraw.Amount

			senderAcct := accounts[sender]
			if senderAcct == nil {
				return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("sender account %s does not exist", sender), nil)
			}
			if senderAcct.Balance.ToInt().Cmp(amount.ToInt()) < 0 {
				return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeRequestFuncFailed, "request function execution failed", nil)
			}

			// Execute withdrawal
			senderAcct.Balance = common.ToBig(new(big.Int).Sub(senderAcct.Balance.ToInt(), amount.ToInt()))
			nonce++
			currentState.Nonce = nonce

			withdrawals = append(withdrawals, common.Withdrawal{DestinationAddress: to, Amount: amount})

			withdrawEvent := common.PlainEvent{
				UserID:       sender,
				EventSubType: "withdrawal",
				Data: []byte(fmt.Sprintf(`{"type":"withdrawal","to":"%s","amount":"0x%s","balance":"0x%s","nonce":%d}`,
					to, amount.ToInt().Text(16), senderAcct.Balance.ToInt().Text(16), nonce)),
			}
			events = append(events, withdrawEvent)

		case "deanonymize":
			// Generate deanonymization report
			var tag string
			if instructions.Deanonymize != nil {
				tag = instructions.Deanonymize.Tag
			}

			reportData := map[string]interface{}{
				"accounts": currentState.Accounts,
				"tag":      tag,
			}
			var err error
			report, err = json.Marshal(reportData)
			if err != nil {
				return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeJsonMarshalError, "failed to marshal deanonymization report", err)
			}
			r.log.Info("Mock Runtime: Generated deanonymization report for application %d", appId)

		default:
			return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeRequestFuncFailed, fmt.Sprintf("unknown instruction type: %s", typ), nil)
		}
	}

	newStateBytes, err := json.Marshal(currentState)
	if err != nil {
		return nil, nil, nil, nil, r.fuel, apperrors.New(apperrors.CodeJsonMarshalError, "failed to serialize new state", err)
	}

	r.log.Info("Mock Runtime: Successfully processed request for application %d, generated %d events and %d withdrawals", appId, len(events), len(withdrawals))
	return newStateBytes, events, withdrawals, report, r.fuel, nil
}

// Close closes the mock runtime and cleans up resources
func (r *MockRuntime) Close() error {
	r.log.Info("Mock Runtime: Closing mock runtime")
	r.log.Info("Mock Runtime: Mock runtime closed successfully")
	return nil
}

// --- helpers ---

func ensureAccount(accounts map[ethCommon.Address]*testAccountState, addr ethCommon.Address) *testAccountState {
	acct := accounts[addr]
	if acct == nil {
		acct = &testAccountState{Address: addr, Balance: common.NewBig(0)}
		accounts[addr] = acct
	}
	return acct
}
