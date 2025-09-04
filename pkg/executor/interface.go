// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"context"

	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/crypto"

	"github.com/horizen-pes/pkg/common"
)

// Config defines the configuration for the executor application
type Config struct {
	// ServerType is the type of server to use (tcp / vsock)
	ServerType string
	// ServerAddr is the address for the TCP server
	ServerAddr string
	// ServerCid is the cid for the v-socket server
	ServerCid uint32
	// ServerPort is the port for the v-socket server
	ServerPort uint32
	// StateKey is the key to use for signing update payloads
	StateKey cryptotypes.AES256Key
	// CommunicationKey is the key to use for encrypting payloads
	CommunicationKey *cryptotypes.PrivateKeyP521
	// SignatureKey is used to sign UpdatePayloads and DeanonymizationReports
	SignatureKey *cryptotypes.PrivateKeySecp256k1 // Key used for signing updatePayload
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	stateKey, _ := crypto.GenerateAESKey()
	communicationKey, _ := crypto.GeneratePrivateKeyP521()
	signatureKey, _ := crypto.GeneratePrivateKeySecp256k1()

	return &Config{
		ServerType:       "tcp",
		ServerAddr:       "localhost:8080",
		StateKey:         stateKey,
		CommunicationKey: communicationKey,
		SignatureKey:     signatureKey,
	}
}

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
	// LoadModule loads a module from bytecode
	LoadModule(ctx context.Context, appId string, wasm []byte) ([]byte, []byte, error)
	// Deposit processes a deposit
	Deposit(ctx context.Context, appId string, sender string, value uint64, state []byte, wasm []byte) ([]byte, []common.PlainEvent, error)
	// ProcessRequest processes a request and returns the new state
	ProcessRequest(ctx context.Context, appId string, sender string, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, appId string, requestId string, payload []byte, state []byte, wasm []byte) ([]byte, error)
	// Close closes the WASM runtime
	Close() error
}
