package common

import (
	"github.com/horizen-pes/pkg/common"
)

// WasmSerializationError is a generic error for failed WASM serialization.
const WasmSerializationError = `{"error":"wasm serialization error"}`

// TODO add applicationId to the definitions where appropriate, in future we will have many different apps

// LoadModuleResult represents the result of a load module operation
type LoadModuleResult struct {
	State []byte      `json:"state"`
	Fuel  *common.Big `json:"fuel"`
	Error string      `json:"error,omitempty"`
}

// DepositResult represents the result of a deposit operation
type DepositResult struct {
	State  []byte              `json:"state"`
	Events []common.PlainEvent `json:"events"`
	Fuel   *common.Big         `json:"fuel"`
	Error  string              `json:"error,omitempty"`
}

// ProcessResult represents the result of a process request operation
type ProcessResult struct {
	State       []byte              `json:"state"`
	Events      []common.PlainEvent `json:"events"`
	Withdrawals []common.Withdrawal `json:"withdrawals"`
	Report      []byte              `json:"report,omitempty"` // Optional deanonymization report
	Fuel        *common.Big         `json:"fuel"`
	Error       string              `json:"error,omitempty"`
}

// DeanonymizationResult represents the result of generating deanonymization report
type DeanonymizationResult struct {
	Report []byte      `json:"report"`
	Fuel   *common.Big `json:"fuel"`
	Error  string      `json:"error,omitempty"`
}

// optional mem statistics
type MemoryStats struct {
	MapSize              int64 `json:"mapSize"`
	CumulativeMemorySize int64 `json:"cumulativeMemorySize"`
}
