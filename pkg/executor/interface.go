// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"context"
	"log"
	"os"

	"github.com/horizen-pes/pkg/communication"

	"github.com/horizen-pes/pkg/common"
	"github.com/magiconair/properties"
)

const confFileName = "executor.conf"

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
}

// DefaultConfig returns the default configuration (possibly overridden by env variables)
func DefaultConfig() *Config {
	serverAddress := os.Getenv("EXECUTOR_IP_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost"
	}
	serverPort := os.Getenv("EXECUTOR_IP_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	return &Config{
		ServerType: "tcp",
		ServerAddr: serverAddress + ":" + serverPort,
	}
}

func ReadConfig() *Config {

	if !fileExists(confFileName) {
		log.Printf("File %s not found. Using default configuration for Executor.\n", confFileName)
		return DefaultConfig()
	}
	// Load properties from file
	config, err := properties.LoadFile(confFileName, properties.UTF8)
	if err != nil {
		panic(err)
	}

	return &Config{
		ServerType: config.MustGetString("ServerType"),
		ServerAddr: config.MustGetString("ServerAddr"),
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
