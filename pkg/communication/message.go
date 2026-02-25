package communication

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/horizen-pes/pkg/common"
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

	// ErrorMessage represents an error message
	ErrorMessage

	// KeyAttestationRequestMessage represents a request from manager to executor to create a key attestation
	KeyAttestationRequestMessage
	// KeyAttestationResponseMessage represents the executor's response with the attestation document
	KeyAttestationResponseMessage
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
}

func (prd *ProcessRequestData) Validate() error {
	if prd.Request == nil {
		return fmt.Errorf("Request is required")
	}
	if err := prd.Request.Validate(); err != nil {
		return fmt.Errorf("invalid Request: %w", err)
	}
	if prd.ApplicationState == nil {
		return fmt.Errorf("ApplicationState is required")
	}
	if len(prd.WasmModule) == 0 {
		return fmt.Errorf("WasmModule cannot be empty")
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

// DeployAppRequestData represents data for a deploy app request message
type DeployAppRequestData struct {
	// Request is the request to deploy an application
	Request *common.Request `json:"request"`
	// ApplicationState is the current state of the application
	ApplicationState *common.ApplicationState `json:"applicationState"`
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

// KeyAttestationResponseData represents data for a key attestation response message
type KeyAttestationResponseData struct {
	// Attestation is the raw attestation document bytes
	Attestation []byte `json:"attestation"`
}

// generateID generates a simple unique ID for message correlation
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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
