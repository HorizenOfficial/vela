// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"context"

	"github.com/horizen-pes/pkg/common"
)

// Executor defines the interface for the WASM Executor
type Executor interface {
	// ProcessRequest processes a request and returns the response
	ProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error)
	// DeployApp deploys a new application
	DeployApp(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.DeanonymizationReport, error)
	// Close closes the executor
	Close() error
}

// WASMRuntime defines the interface for a WASM runtime
type WASMRuntime interface {
	// LoadModule loads a WASM module from bytecode
	LoadModule(ctx context.Context, appId string, wasm []byte) ([]byte, []byte, error) //todo: think about adding properties to manage module initialization
	// ProcessRequest processes a request and returns the new state
	ProcessRequest(ctx context.Context, appId string, payload []byte, state []byte, wasm []byte) ([]byte, []common.Event, []common.Withdrawal, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, req *common.Request, state []byte, wasm []byte) ([]byte, error)
	// Close closes the WASM runtime
	Close() error
}

// Event represents an event to be emitted
type Event struct {
	// ApplicationID is the ID of the application
	ApplicationID string `json:"applicationId"`
	// UserID is the ID of the user associated with the event
	UserID string `json:"userId"`
	// EncryptedData is the encrypted event data
	EncryptedData []byte `json:"encryptedData"`
}
