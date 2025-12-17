package communication

import (
	"context"
)


type AdminCommandServer interface {
	Start(ctx context.Context, identityLogTag string) error
	Stop() error
	// SetRequestHandler sets the handler for incoming requests
	SetCmdHandler(handler AdminCmdHandler)

}

type AdminCmdHandler interface {
	CreateKeyAttestation(ctx context.Context) ([]byte, error)
}

// AdminMessageType represents the command being sent
type AdminMessageType int

const (
	// KeyAttestationRequestMessage represents a request to generate a key attestation
	KeyAttestationRequestMessage AdminMessageType = iota

	// AdminResponseMessage represents a successful response 
	AdminResponseMessage 
	// AdminErrorMessage represents an error message
	AdminErrorMessage
)


// AdminMessage represents an admin command message
type AdminMessage struct {
	Type AdminMessageType `json:"type"`
	Data interface{} `json:"data"`
}