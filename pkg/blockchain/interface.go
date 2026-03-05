// Package blockchain defines the interface for interacting with a blockchain system.
package blockchain

import (
	"context"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/HorizenOfficial/vela/pkg/common"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
)

// Client defines the interface for interacting with the blockchain
type Client interface {
	// SubmitRequest submits a request to the blockchain, returning submitted request id as a string and block number in which it is submitted
	SubmitRequest(ctx context.Context, protocolVersion uint8, applicationId common.ApplicationIdType, requestType common.RequestType, payload []byte, depositAmount *big.Int, maxFeeValue *big.Int) (common.RequestIdType, uint64, error)
	// GetPendingRequests gets pending requests from the blockchain
	GetPendingRequests(ctx context.Context) ([]*common.Request, error)
	// GetNextPendingRequest gets next pending request and current state root from the blockchain
	GetNextPendingRequest(ctx context.Context) (*common.Request, [32]byte, error)
	// SubmitStateUpdate submits a state update to the blockchain
	SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error
	//GetTeePublicKey gets the public key from the blockchain needed to encrypt payloads
	GetTeePublicKey(ctx context.Context) (*cryptotypes.PublicKeyP521, error)
	// ChainID returns the connected chain ID.
	ChainID(ctx context.Context) (*big.Int, error)
	// LatestBlockNumber returns the latest block number from the chain.
	LatestBlockNumber(ctx context.Context) (uint64, error)

	// GetPendingPayments returns the pending payment balance for the given address.
	GetPendingPayments(ctx context.Context, addr ethCommon.Address) (*big.Int, error)
	// WithdrawPayments calls withdrawPayments on the ProcessorEndpoint contract for the given payee.
	WithdrawPayments(ctx context.Context, payee ethCommon.Address) error

	// Close closes the blockchain client
	Close() error
	// Connect connects to the blockchain. It is idempotent: calling it on an
	// already-connected client returns nil.
	Connect(ctx context.Context) error
	// IsConnected returns true if the client has successfully connected.
	IsConnected() bool
}
