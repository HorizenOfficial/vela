package common

import (
	"github.com/horizen-pes/pkg/common"
)

// WasmSerializationError is a generic error for failed WASM serialization.
const WasmSerializationError = "{}"

// TODO add applicationId to the definitions where appropriate, in future we will have many differnt apps

// DepositResult represents the result of a deposit operation
type DepositResult struct {
	State  []byte              `json:"state"`
	Events []common.PlainEvent `json:"events"`
	Error  string              `json:"error,omitempty"`
}

// ProcessResult represents the result of a process request operation
type ProcessResult struct {
	State       []byte              `json:"state"`
	Events      []common.PlainEvent `json:"events"`
	Withdrawals []common.Withdrawal `json:"withdrawals"`
	Error       string              `json:"error,omitempty"`
}

// DeanonymizationResult represents the result of generating deanonymization report
type DeanonymizationResult struct {
	Report []byte `json:"report"`
	Error  string `json:"error,omitempty"`
}

type DepositEvent struct {
	Type    string `json:"type"`
	Amount  uint64 `json:"amount"`
	Balance uint64 `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

type SenderEvent struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Amount  uint64 `json:"amount"`
	Balance uint64 `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

type RecipientEvent struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	Amount  uint64 `json:"amount"`
	Balance uint64 `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

type WithdrawalEvent = SenderEvent

// statistics for memory allocation in WASM modules
type MemoryStats struct {
	MapSize              int64 `json:"mapSize"`
	CumulativeMemorySize int64 `json:"cumulativeMemorySize"`
}
