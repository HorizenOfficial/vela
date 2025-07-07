// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"context"
	"github.com/horizen-pes/pkg/communication"

	"github.com/horizen-pes/pkg/common"
)

// Config defines the configuration for the executor application
type Config struct {
	// ServerType is the type of server to use (tcp / vsock)
	ServerType string
	// ServerAddr is the address for the TCP server
	ServerAddr string
	// ServerPort is the port for the v-socket server
	ServerPort uint32
	// SigningKey is the key to use for signing update payloads
	SigningKey []byte
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		ServerType: "tcp",
		ServerAddr: "localhost:8080",
		ServerPort: 5000,
		SigningKey: []byte("dummy-signing-key"),
	}
}

// Executor defines the interface for the WASM Executor
type Executor interface {
	// RequestHandler interface to handle requests from the communication layer
	communication.RequestHandler
	// Close closes the executor
	Close() error
}

// Runtime defines the interface for a WASM runtime
type Runtime interface {
	// LoadModule loads a module from bytecode
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
