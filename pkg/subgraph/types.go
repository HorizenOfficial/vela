package subgraph

import (
	"context"
	"math/big"

	"github.com/horizen-pes/pkg/common"
)

// Client defines the subgraph operations used by the services.
type Client interface {
	GetRequestCompletedByID(ctx context.Context, requestID common.RequestIdType) (*RequestCompleted, error)
	GetUserEvents(ctx context.Context, applicationID common.ApplicationIdType, eventSubType string, limit int, before *big.Int) ([]UserEvent, error)
}

// RequestCompleted is the projection returned by the subgraph.
type RequestCompleted struct {
	RequestID       common.RequestIdType
	Status          common.RequestResultStatus
	ErrorCode       uint8
	ErrorMessage    string
	ApplicationFees *big.Int
	BlockNumber     uint64
}

// UserEvent is the projection returned by the subgraph.
type UserEvent struct {
	ApplicationID common.ApplicationIdType
	RequestID     common.RequestIdType
	EventSubType  string
	EncryptedData []byte
	BlockNumber   uint64
	LogIndex      uint64
	SortKey       *big.Int
}
