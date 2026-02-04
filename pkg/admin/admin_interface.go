package admin

import (
	"context"
)


type AdminCommandServer interface {
	Start(ctx context.Context, identityLogTag string) error
	Stop() error
	// SetRequestHandler sets the handler for incoming requests
	SetCmdHandler(handler AdminCmdHandler)

}

// AdminCmdHandler is the interface for handling executor admin commands
type AdminCmdHandler interface {
	CreateKeyAttestation(ctx context.Context) ([]byte, error)
}

// ManagerCmdHandler is the interface for handling manager admin commands
type ManagerCmdHandler interface {
	GetVersion(ctx context.Context) (string, error)
}

// AdminMessageType represents the command being sent
type AdminMessageType int

const (
	// KeyAttestationRequestMessage represents a request to generate a key attestation (executor)
	KeyAttestationRequestMessage AdminMessageType = iota

	// AdminResponseMessage represents a successful response
	AdminResponseMessage
	// AdminErrorMessage represents an error message
	AdminErrorMessage

	// GetVersionRequestMessage represents a request to get the version (manager)
	GetVersionRequestMessage
)


// AdminMessage represents an admin command message
type AdminMessage struct {
	Type AdminMessageType `json:"type"`
	Data interface{} `json:"data"`
}