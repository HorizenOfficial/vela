// Package blockchain defines the interface for interacting with a blockchain system.
package blockchain

import (
	"context"
	"math/big"

	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
)

// Client defines the interface for interacting with the blockchain
type Client interface {
	// SubmitRequest submits a request to the blockchain
	SubmitRequest(ctx context.Context, protocolVersion uint8, applicationId *big.Int, requestType *common.RequestType, payload []byte, value *big.Int) (string, error)
	// GetPendingRequests gets pending requests from the blockchain
	GetPendingRequests(ctx context.Context) ([]*common.Request, error)
	// GetNextPendingRequest gets next pending request and current state root from the blockchain
	GetNextPendingRequest(ctx context.Context) (*common.Request, [32]byte, error)
	// MarkRequestFailed marks a request as failed
	MarkRequestFailed(ctx context.Context, requestID string) error
	// SubmitStateUpdate submits a state update to the blockchain
	SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error
	// SubmitDeanonymizationReport submits a deanonymization report to the blockchain
	SubmitDeanonymizationReport(ctx context.Context, update *common.DeanonymizationReport) error
	// GetUserEvents gets decryptable user events in the given block range
	GetUserEvents(ctx context.Context, privKey cryptotypes.PrivateKeyP521, applicationId big.Int, fromBlock uint64, toBlock uint64, filter func([]byte) bool, stopAtFirst bool) ([][]byte, error)

	// Close closes the blockchain client
	Close() error
	// Connect connects to the blockchain
	Connect(ctx context.Context) error
}

