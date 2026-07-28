package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/HorizenOfficial/vela/pkg/logger"
)

// AggregatedGetLogLevelResponse is returned when target="all" for GetLogLevel,
// containing the log level of both the Manager and the Executor.
// If one component fails, its error field is populated and the other result is still returned.
type AggregatedGetLogLevelResponse struct {
	Manager       string `json:"manager,omitempty"`
	ManagerError  string `json:"managerError,omitempty"`
	Executor      string `json:"executor,omitempty"`
	ExecutorError string `json:"executorError,omitempty"`
}

// AggregatedSetLogLevelResponse is returned when target="all" for SetLogLevel,
// containing the result from both the Manager and the Executor.
// If one component fails, its error field is populated and the other result is still returned.
type AggregatedSetLogLevelResponse struct {
	Manager       string `json:"manager,omitempty"`
	ManagerError  string `json:"managerError,omitempty"`
	Executor      string `json:"executor,omitempty"`
	ExecutorError string `json:"executorError,omitempty"`
}

// AggregatedGetVersionResponse is returned when target="all" for GetVersion,
// containing the version of both the Manager and the Executor.
// If one component fails, its error field is populated and the other result is still returned.
type AggregatedGetVersionResponse struct {
	Manager       string `json:"manager,omitempty"`
	ManagerError  string `json:"managerError,omitempty"`
	Executor      string `json:"executor,omitempty"`
	// ExecutorPcr0 is the hex-encoded PCR0 of the running enclave image the
	// executor reported at handshake, or "" in non-Nitro (TCP/dev) mode.
	ExecutorPcr0 string `json:"executorPcr0,omitempty"`
	ExecutorError string `json:"executorError,omitempty"`
}

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

	// GetVersionRequestMessage represents a request to get the version
	GetVersionRequestMessage AdminMessageType = "get_version"

	// SetLogLevelRequestMessage represents a request to change the log level
	SetLogLevelRequestMessage AdminMessageType = "set_log_level"

	// GetLogLevelRequestMessage represents a request to get the current log level
	GetLogLevelRequestMessage AdminMessageType = "get_log_level"

	// SetWasmCacheSizeRequestMessage represents a request to change the WASM module cache size (executor only)
	SetWasmCacheSizeRequestMessage AdminMessageType = "set_wasm_cache_size"

	// GetWasmCacheSizeRequestMessage represents a request to get the WASM module cache size (executor only)
	GetWasmCacheSizeRequestMessage AdminMessageType = "get_wasm_cache_size"
)

// Admin command type constants for ForwardAdminCommand dispatching through the
// communication channel. These are used by both the Manager (sender) and
// Executor (handler).
const (
	AdminCmdKeyAttestation = "key_attestation"
	AdminCmdGetVersion     = "get_version"
	AdminCmdSetLogLevel    = "set_log_level"
	AdminCmdGetLogLevel       = "get_log_level"
	AdminCmdSetWasmCacheSize  = "set_wasm_cache_size"
	AdminCmdGetWasmCacheSize  = "get_wasm_cache_size"
	// AdminCmdShutdown asks the Executor to exit cleanly (drain step 3). It is
	// forwarded by the Manager over the communication channel, not exposed as an
	// operator-facing verb on the Manager admin server.
	AdminCmdShutdown = "shutdown"
)

// ShutdownResponse is the Executor's acknowledgement of an AdminCmdShutdown.
type ShutdownResponse struct {
	// Stopping is always true; the executor exits after sending this ack.
	Stopping bool `json:"stopping"`
}

// SupportedLogLevels lists all valid log level strings accepted by SetLogLevel.
const SupportedLogLevels = "trace, debug, info, warn, error, fatal, panic, disabled"

// SetLogLevelRequest is the payload for SetLogLevelRequestMessage.
type SetLogLevelRequest struct {
	Level string `json:"level"`
}

// GetLogLevelRequest is the payload for GetLogLevelRequestMessage.
// Currently empty; retained for forward compatibility.
type GetLogLevelRequest struct{}

// GetVersionRequest is the payload for GetVersionRequestMessage.
// Currently empty; retained for forward compatibility.
type GetVersionRequest struct{}

// Target constants for admin commands. The Manager admin server is the single
// entry point; these values control which component(s) a command applies to.
const (
	TargetManager  = "manager"  // Apply to Manager only
	TargetExecutor = "executor" // Forward to Executor only
	TargetAll      = "all"      // Apply to Manager locally AND forward to Executor
)

// ValidateTarget checks the target field for admin commands received by the
// Manager admin server. Accepted values: "manager", "executor", "all", or ""
// (empty defaults to "all" at the call-site). Returns an error for unknown values.
func ValidateTarget(target string) error {
	switch target {
	case "", TargetManager, TargetExecutor, TargetAll:
		return nil
	default:
		return fmt.Errorf("unknown target '%s'; valid targets: '%s', '%s', '%s'", target, TargetManager, TargetExecutor, TargetAll)
	}
}

// SetLogLevelResponse is the successful response payload for SetLogLevel.
type SetLogLevelResponse struct {
	Level string `json:"level"`
}

// HandleSetLogLevel validates and applies a log level change.
// It is the shared implementation used by both Manager and Executor.
func HandleSetLogLevel(log logger.Logger, componentName, level string) (SetLogLevelResponse, error) {
	log.Info("%s: SetLogLevel command received, level=%s", componentName, level)
	if _, ok := log.(*logger.ZeroNetworkLogger); !ok {
		return SetLogLevelResponse{}, fmt.Errorf("SetLogLevel is only supported with the ZeroNetworkLogger")
	}
	if level == "" {
		return SetLogLevelResponse{}, fmt.Errorf("invalid log level: level must not be empty; supported levels: %s", SupportedLogLevels)
	}
	if err := log.SetLevel(level); err != nil {
		return SetLogLevelResponse{}, fmt.Errorf("invalid log level '%s'; supported levels: %s", level, SupportedLogLevels)
	}
	return SetLogLevelResponse{Level: level}, nil
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

// SetWasmCacheSizeRequest is the payload for SetWasmCacheSizeRequestMessage.
type SetWasmCacheSizeRequest struct {
	MaxCachedModules int `json:"maxCachedModules"`
}

// SetWasmCacheSizeResponse is the response for SetWasmCacheSizeRequestMessage.
type SetWasmCacheSizeResponse struct {
	MaxCachedModules int `json:"maxCachedModules"`
}

// GetWasmCacheSizeResponse is the response for GetWasmCacheSizeRequestMessage.
type GetWasmCacheSizeResponse struct {
	MaxCachedModules int `json:"maxCachedModules"`
}

// AdminMessage represents an admin command message.
// Target controls routing: "manager", "executor", "all", or "" (defaults to "all").
type AdminMessage struct {
	Type   AdminMessageType `json:"type"`
	Target string           `json:"target,omitempty"`
	Data   json.RawMessage  `json:"data"`
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
