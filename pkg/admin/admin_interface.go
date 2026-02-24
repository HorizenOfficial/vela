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
// The manager implements this to handle admin commands, forwarding executor-specific ones over the communication channel
type AdminCmdHandler interface {
	// ExecuteCommand processes an admin command and returns the result.
	// The returned data must be ready for json.Marshal.
	ExecuteCommand(ctx context.Context, msg AdminMessage) (interface{}, error)
}

// AdminMessageType represents the command being sent
type AdminMessageType string

const (
	// AdminResponseMessage represents a successful response
	AdminResponseMessage AdminMessageType = "response"
	// AdminErrorMessage represents an error message
	AdminErrorMessage AdminMessageType = "error"

	// KeyAttestationRequestMessage represents a request to generate a key attestation (forwarded to executor)
	KeyAttestationRequestMessage AdminMessageType = "key_attestation"

	// GetVersionRequestMessage represents a request to get the version (manager)
	GetVersionRequestMessage AdminMessageType = "get_version"

	// SetLogLevelRequestMessage represents a request to change the log level (manager)
	SetLogLevelRequestMessage AdminMessageType = "set_log_level"

	// GetLogLevelRequestMessage represents a request to get the current log level (manager)
	GetLogLevelRequestMessage AdminMessageType = "get_log_level"
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
