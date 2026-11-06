// Package executor provides functionality for executing WASM modules within a secure environment.
// It uses Wasmtime-go as the runtime for WASM Modules.
package executor

import (
	"context"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	"github.com/HorizenOfficial/vela/pkg/communication"

	"github.com/HorizenOfficial/vela/pkg/common"
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

// moduleCacheController is an optional interface for runtimes that support
// module cache management (LRU eviction). Used by admin commands to get/set
// cache limits at runtime without importing the concrete wasm package.
type moduleCacheController interface {
	SetMaxCachedModules(max int)
	GetMaxCachedModules() int
}

// Runtime defines the interface for a WASM runtime
type Runtime interface {
	// Deploy loads a module and initializes it with constructor parameters
	Deploy(ctx context.Context, appId common.ApplicationIdType, constructorParams []byte, wasm []byte) ([]byte, *big.Int, error)
	// Deposit processes a deposit with token awareness (tokenAddress = 0x0 for ETH)
	Deposit(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, tokenAddress ethCommon.Address, depositAmount *big.Int, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.AppEvent, *big.Int, *apperrors.RequestFailure)
	// ProcessRequest processes a request and returns the new state, events, withdrawals, and optionally a deanonymization report
	ProcessRequest(ctx context.Context, appId common.ApplicationIdType, sender ethCommon.Address, requestType common.RequestType, payload []byte, state []byte, wasm []byte) ([]byte, []common.PlainEvent, []common.AppEvent, []common.Withdrawal, []byte, *big.Int, *apperrors.RequestFailure)
	// Close closes the WASM runtime
	Close() error
}
