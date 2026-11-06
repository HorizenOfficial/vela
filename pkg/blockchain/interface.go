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
	SubmitRequest(ctx context.Context, protocolVersion uint8, applicationId common.ApplicationIdType, requestType common.RequestType, payload []byte, tokenAddress ethCommon.Address, assetAmount *big.Int, maxFeeValue *big.Int) (common.RequestIdType, uint64, error)
	// SubmitDeployRequest submits a deploy request to the blockchain, returning the derived application id, submitted request id and block number
	SubmitDeployRequest(ctx context.Context, protocolVersion uint8, payload []byte, maxFeeValue *big.Int) (common.ApplicationIdType, common.RequestIdType, uint64, error)
	// GetPendingRequests gets pending requests from the blockchain
	GetPendingRequests(ctx context.Context) ([]*common.Request, error)
	// GetPendingRequestsWithStateRoot fetches up to maxCount pending requests for the
	// application selected on-chain (round-robin, see docs/design/BATCH_EXECUTION.md
	// section 4.3), together with the selected applicationId and its on-chain state root.
	// The manager does not choose the application — the contract does.
	GetPendingRequestsWithStateRoot(ctx context.Context, maxCount uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error)
	// SubmitStateUpdate submits a state update to the blockchain
	SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error
	// SubmitBatchStateUpdate submits a batch of per-request update payloads together
	// with a single batch signature covering all entry hashes in one transaction.
	SubmitBatchStateUpdate(ctx context.Context, updates []*common.UpdatePayload, batchSignature []byte) error
	//GetTeePublicKey gets the public key from the blockchain needed to encrypt payloads
	GetTeePublicKey(ctx context.Context) (*cryptotypes.PublicKeyP521, error)
	// ChainID returns the connected chain ID.
	ChainID(ctx context.Context) (*big.Int, error)
	// LatestBlockNumber returns the latest block number from the chain.
	LatestBlockNumber(ctx context.Context) (uint64, error)

	// GetPendingClaims returns the pending claim balance for the given token and address.
	GetPendingClaims(ctx context.Context, tokenAddress ethCommon.Address, addr ethCommon.Address) (*big.Int, error)
	// Claim calls claim on the ProcessorEndpoint contract for the given token and payee.
	Claim(ctx context.Context, tokenAddress ethCommon.Address, payee ethCommon.Address) error

	// Close closes the blockchain client
	Close() error
	// Connect connects to the blockchain. It is idempotent: calling it on an
	// already-connected client returns nil.
	Connect(ctx context.Context) error
	// IsConnected returns true if the client has successfully connected.
	IsConnected() bool
}
