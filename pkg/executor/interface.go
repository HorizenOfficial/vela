// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	ethCommon "github.com/ethereum/go-ethereum/common"
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
	// KeySetRecoveryType is the type of recovery mechanism to use for the keyset
	KeySetRecoveryType int
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
	recTypeVar := os.Getenv("EXECUTOR_KEYSET_RECOVERY_TYPE")
	recType, err := strconv.Atoi(recTypeVar)
	if err != nil {
		fmt.Printf("Failed to convert EXECUTOR_KEYSET_RECOVERY_TYPE for error %v, using default value\n", err)
		recType = 0
	}

	return &Config{
		ServerType:         "tcp",
		ServerAddr:         serverAddress + ":" + serverPort,
		KeySetRecoveryType: recType,
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
		ServerType:         config.MustGetString("ServerType"),
		ServerAddr:         config.MustGetString("ServerAddr"),
		KeySetRecoveryType: config.GetInt("KeySetRecoveryType", 0),
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
	LoadModule(ctx context.Context, appId common.ApplicationIdType, wasm []byte) ([]byte, error)
	// Deposit processes a deposit
	Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, value *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, error)
	// ProcessRequest processes a request and returns the new state
	ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.Withdrawal, error)
	// GenerateDeanonymizationReport generates a deanonymization report
	GenerateDeanonymizationReport(ctx context.Context, appId common.ApplicationIdType, payload []byte, state []byte, wasm []byte) ([]byte, error)
	// Close closes the WASM runtime
	Close() error
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
