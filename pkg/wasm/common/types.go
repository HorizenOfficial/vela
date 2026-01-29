package common

import (
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
)

// WasmSerializationError is a generic error for failed WASM serialization.
const WasmSerializationError = `{"error":"wasm serialization error"}`

// TODO add applicationId to the definitions where appropriate, in future we will have many different apps

// LoadModuleResult represents the result of a load module operation
type LoadModuleResult struct {
	State []byte   `json:"state"`
	Fuel  *big.Int `json:"fuel"`
	Error string   `json:"error,omitempty"`
}

// DepositResult represents the result of a deposit operation
type DepositResult struct {
	State  []byte              `json:"state"`
	Events []common.PlainEvent `json:"events"`
	Fuel   *big.Int            `json:"fuel"`
	Error  string              `json:"error,omitempty"`
}

// ProcessResult represents the result of a process request operation
type ProcessResult struct {
	State       []byte              `json:"state"`
	Events      []common.PlainEvent `json:"events"`
	Withdrawals []common.Withdrawal `json:"withdrawals"`
	Report      []byte              `json:"report,omitempty"` // Optional deanonymization report
	Fuel        *big.Int            `json:"fuel"`
	Error       string              `json:"error,omitempty"`
}

// DeanonymizationResult represents the result of generating deanonymization report
type DeanonymizationResult struct {
	Report []byte   `json:"report"`
	Fuel   *big.Int `json:"fuel"`
	Error  string   `json:"error,omitempty"`
}

type DepositEvent struct {
	Type    string   `json:"type"`
	Amount  *big.Int `json:"amount"`
	Balance *big.Int `json:"balance"`
	Nonce   uint64   `json:"nonce"`
}

type SenderEvent struct {
	Type    string            `json:"type"`
	To      ethCommon.Address `json:"to"`
	Amount  *big.Int          `json:"amount"`
	Balance *big.Int          `json:"balance"`
	Nonce   uint64            `json:"nonce"`
}

type RecipientEvent struct {
	Type    string            `json:"type"`
	From    ethCommon.Address `json:"from"`
	Amount  *big.Int          `json:"amount"`
	Balance *big.Int          `json:"balance"`
	Nonce   uint64            `json:"nonce"`
}

type WithdrawalEvent = SenderEvent

// optional mem statistics
type MemoryStats struct {
	MapSize              int64 `json:"mapSize"`
	CumulativeMemorySize int64 `json:"cumulativeMemorySize"`
}
