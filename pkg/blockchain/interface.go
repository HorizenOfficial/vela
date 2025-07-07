// Package blockchain defines the interface for interacting with a blockchain system.
package blockchain

import (
	"context"

	"github.com/horizen-pes/pkg/common"
)

// Client defines the interface for interacting with the blockchain
type Client interface {
	// GetPendingRequests gets pending requests from the blockchain
	GetPendingRequests(ctx context.Context) ([]*common.Request, error)
	// MarkRequestCompleted marks a request as completed
	MarkRequestCompleted(ctx context.Context, requestID string) error
	// MarkRequestFailed marks a request as failed
	MarkRequestFailed(ctx context.Context, requestID string) error
	// SubmitStateUpdate submits a state update to the blockchain
	SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error
	// SubmitDeanonymizationReport submits a deanonymization report to the blockchain
	SubmitDeanonymizationReport(ctx context.Context, update *common.DeanonymizationReport) error
	// GetPublicKey gets the public key for an address
	GetPublicKey(ctx context.Context, address string) ([]byte, error)
	// Close closes the blockchain client
	Close() error
}

// TestClient defines the interface for testing the blockchain client
type TestClient interface {
	Client
	// SubmitRequest submits a request to the blockchain
	SubmitRequest(ctx context.Context, req *common.Request) error
	// RegisterPublicKey registers a public key for an address
	RegisterPublicKey(ctx context.Context, address string, publicKey []byte) error
	// SubscribeToEvents subscribes to events from the blockchain
	SubscribeToEvents(ctx context.Context, eventCh chan<- interface{}) error
	// GetApplicationState gets the state of an application
	GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error)
	// GetDeanonymizationReport gets a deanonymization report
	GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error)
	// GetWithdrawals gets withdrawal requests
	GetWithdrawals(ctx context.Context, applicationID string) (*[]common.Withdrawal, error)
}
