// Package common provide the data structures used in the application.
package common

import (
	"math/big"
)

// RequestType represents the type of request being sent to the TEE
type RequestType string

const (
	// AssociateKey records an association between an Ethereum address and a Secp521r1_PubKey
	AssociateKey RequestType = "associatekey"
	// Deploy represents a request to deploy a new application
	Deploy RequestType = "deploy"
	// Process is used for processing a batch of requests
	Process RequestType = "process"
	// Deanonymize is used for deanonymization requests
	Deanonymize RequestType = "deanonymize"
)

// Request represents a request to the system
type Request struct {
	// ProtocolVersion is the version of the protocol being used
	ProtocolVersion string `json:"protocolVersion"`
	// ApplicationID is the ID of the application (empty for Deploy)
	ApplicationID string `json:"applicationId"`
	// RequestID is a unique identifier for the request
	RequestID string `json:"requestId"`
	// RequestType is the type of request
	RequestType RequestType `json:"requestType"`
	// Payload is the payload for the request
	// All payloads except the one for AssociateKey are encrypted
	Payload []byte `json:"payload"`
	// Timestamp is the time the request was submitted
	Timestamp int64 `json:"timestamp"`
	// Sender is the address of the sender
	Sender string `json:"sender"`
	// Value is the optional deposit value in WEI
	Value uint64 `json:"value"`
}

// Event represents an event to be emitted
type Event struct {
	// ApplicationID is the ID of the application
	ApplicationID string `json:"applicationId"`
	// UserID is the ID of the user associated with the event
	UserID string `json:"userId"`
	// EncryptedData is the encrypted event data
	EncryptedData []byte `json:"encryptedData"`
}

// Withdrawal represents a withdrawal from the system
type Withdrawal struct {
	// DestinationAddress is the address to send the funds to
	DestinationAddress string `json:"destinationAddress"`
	// Amount is the amount to withdraw in WEI
	Amount uint64 `json:"amount"`
}

// UpdatePayload represents an update to the state
type UpdatePayload struct {
	// ApplicationID is the ID of the application
	ApplicationID string `json:"applicationId"`
	// RequestID is the ID of the request being processed
	RequestID string `json:"requestId"`
	// PrevStateRoot is the previous state root
	PrevStateRoot [32]byte `json:"prevStateRoot"`
	// NewStateRoot is the new state root
	NewStateRoot [32]byte `json:"newStateRoot"`
	// Events is a list of events to emit
	Events []Event `json:"events"`
	// Withdrawals is a list of withdrawals to execute
	Withdrawals []Withdrawal `json:"withdrawals"`
	// Signature is the TEE signature
	Signature []byte `json:"signature"`
}

// ApplicationState represents the state of an application
type ApplicationState struct {
	// ApplicationID is the ID of the application
	ApplicationID string `json:"applicationId"`
	// StateRoot is the root hash of the state
	StateRoot [32]byte `json:"stateRoot"`
	// EncryptedState is the encrypted state data
	EncryptedState []byte `json:"encryptedState"`
}

type WASMData struct {
	// ApplicationID is the ID of the application
	ApplicationID string `json:"applicationId"`
	// Bytecode is the wasm bytecode
	Bytecode []byte `json:"bytecode"`
}

// DeanonymizationReport represents a report for deanonymization
type DeanonymizationReport struct {
	// ApplicationID is the ID of the application
	ApplicationID string `json:"applicationId"`
	// ReportID is a unique identifier for the report
	ReportID string `json:"reportId"`
	// EncryptedReport is the encrypted report data
	EncryptedReport []byte `json:"encryptedReport"`
}

// PlainEvent represents an emitted event before encryption.
type PlainEvent struct {
	// UserID is the ID of the user associated with the event
	UserID string `json:"userId"`
	// Data is the encrypted event data
	Data []byte `json:"data"`
}

func StringToBigInt(s string) (*big.Int, bool) {
	i, ok := new(big.Int).SetString(s, 10)
	return i, ok
}
