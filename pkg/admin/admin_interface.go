package admin

import (
	"context"
	"encoding/json"
	"fmt"
)


type AdminCommandServer interface {
	Start(ctx context.Context, identityLogTag string) error
	Stop() error
	// SetCmdHandler sets the handler for incoming requests
	SetCmdHandler(handler AdminCmdHandler)
}

// AdminCmdHandler is the interface for handling admin commands
// Both executor and manager implement this with their own command handling logic
type AdminCmdHandler interface {
	ExecuteCommand(ctx context.Context, msg AdminMessage) (interface{}, error)
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
	Data json.RawMessage  `json:"data"`
}

// NewAdminMessage creates an AdminMessage by marshaling the provided data into json.RawMessage.
func NewAdminMessage(msgType AdminMessageType, data interface{}) (AdminMessage, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return AdminMessage{}, fmt.Errorf("failed to marshal admin message data: %w", err)
	}
	return AdminMessage{
		Type: msgType,
		Data: raw,
	}, nil
}

// IsSupportedCommand checks if a message type is in the list of supported commands.
// This is a generic helper function used by both executor and manager.
func IsSupportedCommand(msgType AdminMessageType, supportedCommands []AdminMessageType) bool {
	for _, supportedType := range supportedCommands {
		if msgType == supportedType {
			return true
		}
	}
	return false
}