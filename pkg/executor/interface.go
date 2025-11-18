// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"context"

	"github.com/horizen-pes/pkg/communication"

	"github.com/horizen-pes/pkg/common"
)

// Executor defines the interface for the WASM Executor
type Executor interface {
	// RequestHandler interface to handle requests from the communication layer
	communication.RequestHandler
	// Start starts the executor server
	Start(ctx context.Context) error
	// Close closes the executor
	Close() error
}

// Runtime defines the interface for a WASM runtime
type Runtime interface {
	// LoadModule loads a module from bytecode. Must return the initial application state (or an empty byte array if any)
	LoadModule(ctx context.Context, appId string, wasm []byte) ([]byte, error)
	// Deposit processes a deposit
	Deposit(ctx context.Context, appId string, sender string, value uint64, state []byte, wasm []byte) ([]byte, []common.PlainEvent, error)
	// ProcessRequest processes a request and returns the new state
	ProcessRequest(ctx context.Context, appId string, sender string, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, appId string, payload []byte, state []byte, wasm []byte) ([]byte, error)
	// Close closes the WASM runtime
	Close() error
}
