package communication

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
)

const MsgDelimiter = byte('\n')

// MessageType represents the type of message being sent
type MessageType int

const (
	// ProcessRequestMessage represents a request to process an action
	ProcessRequestMessage MessageType = iota
	// ProcessResponseMessage represents a response to a process request
	ProcessResponseMessage
	// DeployAppRequestMessage represents a request to deploy an application
	DeployAppRequestMessage
	// DeployAppResponseMessage represents a response to a deploy app request
	DeployAppResponseMessage
	// GetKeysetRecoveryRequestMessage represents a handshake message from executor to manager
	GetKeysetRecoveryRequestMessage
	// GetKeysetRecoveryResponseMessage represents a handshake message from manager to executor
	GetKeysetRecoveryResponseMessage
	// SetKeysetRecoveryRequestMessage represents a request to set the keyset recovery data from executor to manager
	SetKeysetRecoveryRequestMessage
	// SetKeysetRecoveryResponseMessage represents a response to a set keyset recovery request from manager to executor
	SetKeysetRecoveryResponseMessage
	// KeysetRecoveryResultMessage represent a confirmation from executor to manager confirming the recovery of keyset
	KeysetRecoveryResultMessage

	// AdminCommandRequestMessage represents a request to execute an admin command on the executor.
	// Sent from manager (client) to executor (server) through the existing communication channel.
	AdminCommandRequestMessage

	// AdminCommandResponseMessage represents a response to an admin command request.
	// Sent from executor (server) to manager (client) through the existing communication channel.
	AdminCommandResponseMessage

	// BatchProcessRequestMessage represents a request to process a batch of requests
	// for a single application in one vsock round-trip.
	BatchProcessRequestMessage
	// BatchProcessResponseMessage represents a response to a batch process request.
	BatchProcessResponseMessage

	// ErrorMessage represents an error message
	ErrorMessage
)

// Message represents a message exchanged between components
type Message struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
	Data interface{} `json:"data"`
}

// PendingRequest represents a request waiting for response
type PendingRequest struct {
	ResponseChan chan Message
	Timeout      time.Time
}

// ProcessRequestData represents data for a process request message
type ProcessRequestData struct {
	// Request is the request to process
	Request *common.Request `json:"request"`
	// ApplicationState is the current state of the application
	ApplicationState *common.ApplicationState `json:"applicationState"`
	// WasmModule is the WASM module to execute
	WasmModule []byte `json:"wasmModule"`
	// ExecutionBudgetMs is how much longer the caller will wait for this request,
	// in milliseconds relative to send time. 0 (or absent) means "not supplied",
	// so an older peer that does not set it simply falls back to the executor's own
	// configured bound. See contextForExecutionBudget for how it is applied and why
	// it can only shorten execution.
	ExecutionBudgetMs int64 `json:"executionBudgetMs,omitempty"`
}

func (prd *ProcessRequestData) Validate() error {
	if prd.Request == nil {
		return fmt.Errorf("Request is required")
	}
	if err := prd.Request.Validate(); err != nil {
		return fmt.Errorf("invalid Request: %w", err)
	}
	if err := validateExecutionBudget(prd.ExecutionBudgetMs); err != nil {
		return err
	}
	return nil
}

// ProcessResponseData represents data for a process response message
type ProcessResponseData struct {
	// UpdatePayload is the update payload
	UpdatePayload *common.UpdatePayload `json:"updatePayload"`
	// UpdatedApplicationState is the updated application state
	UpdatedApplicationState *common.ApplicationState `json:"updatedApplicationState"`
	// DeanonymizationReport is the optional deanonymization report (present if request type was Deanonymize)
	DeanonymizationReport *common.DeanonymizationReport `json:"deanonymizationReport,omitempty"`
}

// BatchProcessRequestData represents data for a batch process request message.
// A batch is always scoped to a single application: one WASM module, one
// encrypted state, and one ordered list of requests to process sequentially.
type BatchProcessRequestData struct {
	// Requests is the ordered list of requests to process
	Requests []*common.Request `json:"requests"`
	// ApplicationState is the current state of the application
	ApplicationState *common.ApplicationState `json:"applicationState"`
	// WasmModule is the WASM module to execute
	WasmModule []byte `json:"wasmModule"`
	// ExecutionBudgetMs is how much longer the caller will wait for the WHOLE batch,
	// in milliseconds relative to send time. A batch is one message and therefore one
	// budget: its requests share it, so a batch too large to finish within it settles
	// the requests that fit and leaves the rest pending (see
	// contextForExecutionBudget). 0 (or absent) means "not supplied", which leaves the
	// executor's own guest bound as the only limit.
	ExecutionBudgetMs int64 `json:"executionBudgetMs,omitempty"`
}

func (bpr *BatchProcessRequestData) Validate() error {
	if len(bpr.Requests) == 0 {
		return fmt.Errorf("Requests is required")
	}
	if err := validateExecutionBudget(bpr.ExecutionBudgetMs); err != nil {
		return err
	}
	// A batch is scoped to one application, so the state is what defines that
	// scope: without it there is nothing to validate the requests against.
	// The executor rejects a nil state as a hard failure anyway, so failing
	// here keeps the same outcome (no payloads, requests stay pending).
	if bpr.ApplicationState == nil {
		return fmt.Errorf("ApplicationState is required")
	}

	applicationID := bpr.ApplicationState.ApplicationID
	for i, req := range bpr.Requests {
		if req == nil {
			return fmt.Errorf("Requests[%d] is required", i)
		}
		if req.ApplicationID != applicationID {
			return fmt.Errorf("invalid Requests[%d]: ApplicationID mismatch", i)
		}
		if err := req.Validate(); err != nil {
			return fmt.Errorf("invalid Requests[%d]: %w", i, err)
		}
	}
	return nil
}

// BatchProcessResponseData represents data for a batch process response message.
// Individual update payloads are not signed on their own; a single batch
// signature covers all entry hashes. Only the final encrypted state is returned.
type BatchProcessResponseData struct {
	// UpdatePayloads is one update payload per handled request (unsigned individually),
	// in input order. Fewer payloads than input requests means a hard failure stopped
	// the batch at request len(UpdatePayloads) and the rest stay pending.
	UpdatePayloads []*common.UpdatePayload `json:"updatePayloads"`
	// BatchSignature is the single TEE signature covering all entry hashes
	BatchSignature []byte `json:"batchSignature"`
	// UpdatedApplicationState is the final application state (only the last state, not per-request)
	UpdatedApplicationState *common.ApplicationState `json:"updatedApplicationState"`
	// DeanonymizationReports are the optional deanonymization reports produced within the batch
	DeanonymizationReports []*common.DeanonymizationReport `json:"deanonymizationReports,omitempty"`
}

// DeployAppRequestData represents data for a deploy app request message
type DeployAppRequestData struct {
	// Request is the request to deploy an application
	Request *common.Request `json:"request"`
	// ApplicationState is the current state of the application
	ApplicationState *common.ApplicationState `json:"applicationState"`
	// WasmModule is the resolved WASM bytecode if available; nil means artifact unavailable/unresolved.
	WasmModule []byte `json:"wasmModule"`
	// ExecutionBudgetMs is how much longer the caller will wait for this deploy, in
	// milliseconds relative to send time. See ProcessRequestData.ExecutionBudgetMs.
	ExecutionBudgetMs int64 `json:"executionBudgetMs,omitempty"`
}

func (dad *DeployAppRequestData) Validate() error {
	if dad.Request == nil {
		return fmt.Errorf("Request is required")
	}
	if err := dad.Request.Validate(); err != nil {
		return fmt.Errorf("invalid Request: %w", err)
	}
	if err := validateExecutionBudget(dad.ExecutionBudgetMs); err != nil {
		return err
	}
	return nil
}

// DeployAppResponseData represents data for a deploy app response message
type DeployAppResponseData struct {
	// UpdatePayload is the update payload
	UpdatePayload *common.UpdatePayload `json:"updatePayload"`
	// ApplicationState initialized application state
	ApplicationState *common.ApplicationState `json:"applicationState"`
}

// ErrorData represents data for an error message
type ErrorData struct {
	// Code is the error code
	Code string `json:"code"`
	// Message is the error message
	Message string `json:"message"`
}

// GetKeysetRecoveryRequestData represents data for a message from executor to manager
type GetKeysetRecoveryRequestData struct {
}

// GetKeysetRecoveryResponseData represents data for a message from manager to executor
type GetKeysetRecoveryResponseData struct {
	DataFound      bool                          `json:"dataFound"`
	KeySetRecovery *common.EnclaveKeySetRecovery `json:"keySetRecovery"`

	// RequestTimeoutMs is how long the manager waits for any request it sends the
	// executor, in milliseconds — the value every later ExecutionBudgetMs will be
	// derived from. The executor validates its guest execution bound against it
	// during the handshake: the enclave cannot read the manager's configuration, and
	// checking against its own unrelated timeout gave false assurance. Optional: an
	// older manager omits it and it decodes to 0, meaning "not reported".
	RequestTimeoutMs int64 `json:"requestTimeoutMs,omitempty"`
}

// Validate rejects a reply whose DataFound flag and KeySetRecovery payload disagree.
//
// Nothing else pairs the two: the flag is what the executor branches on, and the payload
// is what it then reads, so a reply claiming data while carrying none is dereferenced as
// a nil pointer during the handshake — before any keyset exists, in a goroutine with no
// recover, which takes the whole enclave down. This is the first message the executor
// accepts from outside the TEE, so it must not be trusted to be self-consistent.
//
// A false flag with no payload is the ordinary first-connection reply and stays valid.
func (gkr *GetKeysetRecoveryResponseData) Validate() error {
	if gkr.DataFound && gkr.KeySetRecovery == nil {
		return fmt.Errorf("dataFound is set but keySetRecovery is missing")
	}
	return nil
}

// SetKeysetRecoveryRequestData represents data for a set keyset recovery request message
type SetKeysetRecoveryRequestData struct {
	KeySetRecovery *common.EnclaveKeySetRecovery `json:"keySetRecovery"`
	CommPubKey     string                        `json:"commPubKey"`
	SigningKeyAddr string                        `json:"signingKeyAddr"`
}

// SetKeysetRecoveryResponseData represents data for a set keyset recovery response message
type SetKeysetRecoveryResponseData struct {
}

// KeysetRecoveryResultData represents data for a keyset recovery success message
type KeysetRecoveryResultData struct {
	Error          string `json:"error,omitempty"`
	CommPubKey     string `json:"commPubKey,omitempty"`
	SigningKeyAddr string `json:"signingKeyAddr,omitempty"`
}

// generateID generates a simple unique ID for message correlation
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// AdminCommandRequestData represents data for an admin command request forwarded
// from the manager to the executor through the communication channel.
type AdminCommandRequestData struct {
	// CommandType identifies the admin command (e.g. "set_log_level", "get_log_level").
	CommandType string `json:"commandType"`
	// Data is the command-specific payload, opaque to the communication layer.
	Data json.RawMessage `json:"data"`
}

// AdminCommandResponseData represents data for an admin command response.
type AdminCommandResponseData struct {
	// Data is the command-specific response payload.
	Data json.RawMessage `json:"data"`
}

type Validatable interface {
	Validate() error
}

// extractData extracts request / response message data structs from an interface{}
func extractData[T any](data interface{}) (*T, error) {
	// Convert to JSON and back to exact type T
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data: %w", err)
	}

	var respData T
	if err := json.Unmarshal(jsonData, &respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	if v, ok := any(&respData).(Validatable); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	return &respData, nil
}
