package communication

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/horizen-pes/pkg/common"
	"time"
)

// MessageType represents the type of message being sent
type MessageType int

const delimiter = byte('\n')

const (
	// ProcessRequestMessage represents a request to process an action
	ProcessRequestMessage MessageType = iota
	// ProcessResponseMessage represents a response to a process request
	ProcessResponseMessage
	// DeployAppRequestMessage represents a request to deploy an application
	DeployAppRequestMessage
	// DeployAppResponseMessage represents a response to a deploy app request
	DeployAppResponseMessage
	// DeanonymizationRequestMessage represents a request for deanonymization
	DeanonymizationRequestMessage
	// DeanonymizationResponseMessage represents a response to a deanonymization request
	DeanonymizationResponseMessage
	// GetUserKeysRequestMessage represents a request from server to client for user keys
	GetUserKeysRequestMessage
	// GetUserKeysResponseMessage represents a response to a HandleGetUserKeys request
	GetUserKeysResponseMessage
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
	ResponseChan chan *Message
	Timeout      time.Time
}

// ProcessRequestData represents data for a process request message
type ProcessRequestData struct {
	// Request is the request to process
	Request *common.Request `json:"request"`
	// ApplicationState is the current state of the application
	ApplicationState *common.ApplicationState `json:"applicationState"`
	// SenderKey is the public key of the sender
	SenderKey []byte `json:"senderKey"`
	// WasmModule is the WASM module to execute
	WasmModule []byte `json:"wasmModule"`
}

// ProcessResponseData represents data for a process response message
type ProcessResponseData struct {
	// UpdatePayload is the update payload
	UpdatePayload *common.UpdatePayload `json:"updatePayload"`
	// UpdatedApplicationState is the updated application state
	UpdatedApplicationState *common.ApplicationState `json:"updatedApplicationState"`
}

// DeployAppRequestData represents data for a deploy app request message
type DeployAppRequestData struct {
	// Request is the request to deploy an application
	Request *common.Request `json:"request"`
}

// DeployAppResponseData represents data for a deploy app response message
type DeployAppResponseData struct {
	// UpdatePayload is the update payload
	UpdatePayload *common.UpdatePayload `json:"updatePayload"`
	// ApplicationState initialized application state
	ApplicationState *common.ApplicationState `json:"applicationState"`
	// WasmModule is the WASM module for the deployed application
	WasmModule []byte `json:"wasmModule"`
}

// DeanonymizationRequestData represents data for a deanonymization request message
type DeanonymizationRequestData struct {
	// Request is the request for deanonymization
	Request *common.Request `json:"request"`
	// ApplicationState is the current state of the application
	ApplicationState *common.ApplicationState `json:"applicationState"`
	// SenderKey is the public key of the sender
	SenderKey []byte `json:"senderKey"`
	// WasmModule is the WASM module to execute
	WasmModule []byte `json:"wasmModule"`
}

// DeanonymizationResponseData represents data for a deanonymization response message
type DeanonymizationResponseData struct {
	// Report is the deanonymization report
	Report *common.DeanonymizationReport `json:"report"`
}

// GetUserKeysRequestData represents data for a HandleGetUserKeys request message
type GetUserKeysRequestData struct {
	// Users is the list of users to get keys for
	Users []string `json:"users"`
}

// GetUserKeysResponseData represents data for a HandleGetUserKeys response message
type GetUserKeysResponseData struct {
	// UserKeys is a map of user ID to their public key
	UserKeys map[string][]byte `json:"userKeys"`
}

// ErrorData represents data for an error message
type ErrorData struct {
	// Code is the error code
	Code string `json:"code"`
	// Message is the error message
	Message string `json:"message"`
}

// generateID generates a simple unique ID for message correlation
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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

	return &respData, nil
}
