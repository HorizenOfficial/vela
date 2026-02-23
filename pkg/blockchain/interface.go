// Package blockchain defines the interface for interacting with a blockchain system.
package blockchain

import (
	"context"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
)

// Client defines the interface for interacting with the blockchain
type Client interface {
	// GetNextPendingRequest gets next pending request and current state root from the blockchain
	GetNextPendingRequest(ctx context.Context) (*common.Request, [32]byte, error)
	// MarkRequestFailed marks a request as failed
	MarkRequestFailed(ctx context.Context, requestID common.RequestIdType, requestFailure *apperrors.RequestFailure) error
	// SubmitStateUpdate submits a state update to the blockchain
	SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error

	// Close closes the blockchain client
	Close() error
	// Connect connects to the blockchain
	Connect(ctx context.Context) error
}
