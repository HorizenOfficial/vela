package communication

import (
	"github.com/horizen-pes/pkg/common"
)

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
	// DeanonymizationRequestMessage represents a request for deanonymization
	DeanonymizationRequestMessage
	// DeanonymizationResponseMessage represents a response to a deanonymization request
	DeanonymizationResponseMessage
	// TODO: add getKeys message that goes from server to client
	// ErrorMessage represents an error message
	ErrorMessage
)

// Message represents a message exchanged between components
type Message struct {
	// Type is the type of message
	Type MessageType `json:"type"`
	// Data is the message data
	Data interface{} `json:"data"`
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
	// WasmModule is the WASM module to execute
	WasmModule []byte `json:"wasmModule"`
}

// DeanonymizationResponseData represents data for a deanonymization response message
type DeanonymizationResponseData struct {
	// Report is the deanonymization report
	Report *common.DeanonymizationReport `json:"report"`
}

// ErrorData represents data for an error message
type ErrorData struct {
	// Code is the error code
	Code string `json:"code"`
	// Message is the error message
	Message string `json:"message"`
}
