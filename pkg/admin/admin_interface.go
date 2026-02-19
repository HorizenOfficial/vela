package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/horizen-pes/pkg/logger"
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
	// ExecuteCommand processes an admin command and returns the result.
	// The returned data must be ready for json.Marshal.
	ExecuteCommand(ctx context.Context, msg AdminMessage) (interface{}, error)
}

// AdminMessageType represents the command being sent
type AdminMessageType int

const (
	// AdminResponseMessage represents a successful response
	AdminResponseMessage AdminMessageType = iota
	// AdminErrorMessage represents an error message
	AdminErrorMessage

	// KeyAttestationRequestMessage represents a request to generate a key attestation (executor)
	KeyAttestationRequestMessage

	// GetVersionRequestMessage represents a request to get the version (manager)
	GetVersionRequestMessage

	// SetLogLevelRequestMessage represents a request to change the log level
	SetLogLevelRequestMessage

	// GetLogLevelRequestMessage represents a request to get the current log level
	GetLogLevelRequestMessage

	// ExecutorKeyAttestationRequestMessage represents a proxy request forwarded by the manager to the executor for key attestation
	ExecutorKeyAttestationRequestMessage

	// ExecutorSetLogLevelRequestMessage represents a proxy request forwarded by the manager to the executor to change the log level
	ExecutorSetLogLevelRequestMessage

	// ExecutorGetLogLevelRequestMessage represents a proxy request forwarded by the manager to the executor to get the current log level
	ExecutorGetLogLevelRequestMessage
)

// SupportedLogLevels lists all valid log level strings accepted by SetLogLevel.
const SupportedLogLevels = "trace, debug, info, warn, error, fatal, panic, disabled"

// SetLogLevelRequest is the payload for SetLogLevelRequestMessage.
type SetLogLevelRequest struct {
	Level  string `json:"level"`
	Target string `json:"target,omitempty"`
}

// GetLogLevelRequest is the payload for GetLogLevelRequestMessage.
type GetLogLevelRequest struct {
	Target string `json:"target,omitempty"`
}

// TargetAll is the target value that means "apply to all components".
// Currently handled as "apply to self" since proxy is not yet implemented.
// When the Manager gains proxy capability, target "all" on the Manager will
// apply the command locally AND forward it to the Executor.
const TargetAll = "all"

// ValidateTarget checks the target field against the component's own identity.
// It accepts the component's own name (e.g. "manager" or "executor"), empty, or "all".
// Returns an error for mismatches or unknown values.
func ValidateTarget(target, self string) error {
	switch target {
	case "", self, TargetAll:
		return nil
	case "manager", "executor":
		other := "manager"
		if self == "manager" {
			other = "executor"
		}
		return fmt.Errorf("this is the %s admin server, target '%s' is not supported; connect to the %s admin server directly", self, target, other)
	default:
		return fmt.Errorf("unknown target '%s'; valid targets: 'manager', 'executor', 'all'", target)
	}
}

// SetLogLevelResponse is the successful response payload for SetLogLevel.
type SetLogLevelResponse struct {
	Success bool   `json:"success"`
	Level   string `json:"level"`
}

// HandleSetLogLevel validates and applies a log level change.
// It is the shared implementation used by both Manager and Executor.
func HandleSetLogLevel(log logger.Logger, componentName, level string) (interface{}, error) {
	log.Info("%s: SetLogLevel command received, level=%s", componentName, level)
	if _, ok := log.(*logger.ZeroNetworkLogger); !ok {
		return nil, fmt.Errorf("SetLogLevel is only supported with the ZeroNetworkLogger")
	}
	if level == "" {
		return nil, fmt.Errorf("invalid log level: level must not be empty; supported levels: %s", SupportedLogLevels)
	}
	if err := log.SetLevel(level); err != nil {
		return nil, fmt.Errorf("invalid log level '%s'; supported levels: %s", level, SupportedLogLevels)
	}
	return SetLogLevelResponse{Success: true, Level: level}, nil
}

// HandleGetLogLevel returns the current log level.
// It is the shared implementation used by both Manager and Executor.
func HandleGetLogLevel(log logger.Logger, componentName string) (string, error) {
	log.Info("%s: GetLogLevel command received", componentName)
	if _, ok := log.(*logger.ZeroNetworkLogger); !ok {
		return "", fmt.Errorf("GetLogLevel is only supported with the ZeroNetworkLogger")
	}
	return log.GetLevel(), nil
}

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

