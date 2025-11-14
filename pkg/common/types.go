// Package common provide the data structures used in the application.
package common

import (
	"encoding/hex"
	"fmt"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

type ApplicationIdType uint64

func NewApplicationId(id int64) ApplicationIdType {
	if id < 0 {
		panic("ApplicationIdType cannot be negative")
	}
	return ApplicationIdType(id)
}

func (aid ApplicationIdType) String() string {
	return fmt.Sprintf("%d", uint64(aid))
}

func (aid ApplicationIdType) ToHash() ethCommon.Hash {
	return ethCommon.BytesToHash(new(big.Int).SetUint64(uint64(aid)).Bytes())
}

// RequestType represents the type of request being sent to the TEE
type RequestType uint8

const (
	// Deploy represents a request to deploy a new application
	Deploy RequestType = iota
	// Process is used for processing a batch of requests
	Process
	// Deanonymize is used for deanonymization requests
	Deanonymize
	// AssociateKey records an association between an Ethereum address and a Secp521r1_PubKey
	AssociateKey
)


func (rt RequestType) String() string {
	switch rt {
	case Deploy:
		return "deploy"
	case Process:
		return "process"
	case Deanonymize:
		return "deanonymize"
	case AssociateKey:
		return "associatekey"
	default:
		return "unknown"
	}
}


type RequestIdType [32]byte

func (rt RequestIdType) String() string {
	return hex.EncodeToString(rt[:])
}


// Request represents a request to the system
type Request struct {
	// ProtocolVersion is the version of the protocol being used
	ProtocolVersion uint8 `json:"protocolVersion"`
	// ApplicationID is the ID of the application (empty for Deploy)
	ApplicationID ApplicationIdType `json:"applicationId"`
	// RequestID is a unique identifier for the request
	RequestID RequestIdType `json:"requestId"`
	// RequestType is the type of request
	RequestType RequestType `json:"requestType"`
	// Payload is the payload for the request
	// All payloads except the one for AssociateKey are encrypted
	Payload []byte `json:"payload"`
	// Timestamp is the time the request was submitted
	Timestamp *big.Int `json:"timestamp"`
	// Sender is the address of the sender
	Sender ethCommon.Address `json:"sender"`
	// Value is the optional deposit value in WEI
	Value *big.Int `json:"value"`
}

// Event represents an event to be emitted
type Event struct {
	// ApplicationID is the ID of the application
	ApplicationID ApplicationIdType `json:"applicationId"`
	// UserID is the ID of the user associated with the event
	UserID ethCommon.Address `json:"userId"`
	// EncryptedData is the encrypted event data
	EncryptedData []byte `json:"encryptedData"`
}

// Withdrawal represents a withdrawal from the system
type Withdrawal struct {
	// DestinationAddress is the address to send the funds to
	DestinationAddress ethCommon.Address `json:"destinationAddress"`
	// Amount is the amount to withdraw in WEI
	Amount *big.Int `json:"amount"`
}

// UpdatePayload represents an update to the state
type UpdatePayload struct {
	// ApplicationID is the ID of the application
	ApplicationID ApplicationIdType `json:"applicationId"`
	// RequestID is the ID of the request being processed
	RequestID RequestIdType `json:"requestId"`
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
	ApplicationID ApplicationIdType `json:"applicationId"`
	// StateRoot is the root hash of the state
	StateRoot [32]byte `json:"stateRoot"`
	// EncryptedState is the encrypted state data
	EncryptedState []byte `json:"encryptedState"`
}

type WASMData struct {
	// ApplicationID is the ID of the application
	ApplicationID ApplicationIdType `json:"applicationId"`
	// Bytecode is the wasm bytecode
	Bytecode []byte `json:"bytecode"`
}

// DeanonymizationReport represents a report for deanonymization
type DeanonymizationReport struct {
	// ApplicationID is the ID of the application
	ApplicationID ApplicationIdType `json:"applicationId"`
	// ReportID is a unique identifier for the report
	ReportID RequestIdType `json:"reportId"`
	// EncryptedReport is the encrypted report data
	EncryptedReport []byte `json:"encryptedReport"`
}

// DecryptedReport represents a decrypted deanonymization report
type DecryptedReport struct {
	ApplicationID   ApplicationIdType `json:"applicationId"`
	RequestID       RequestIdType `json:"requestId"`
	ReportDataBytes []byte `json:"reportDataBytes"`
}

// PlainEvent represents an emitted event before encryption.
type PlainEvent struct {
	// UserID is the address of the user associated with the event
	UserID ethCommon.Address `json:"userId"`
	// Data is the encrypted event data
	Data []byte `json:"data"`
}

// RequestResultStatus represents the final status of a request after the execution
type RequestResultStatus uint8

const (
	RequestResultOK RequestResultStatus = iota
	RequestResultFailed
	RequestResultFailedNotRefunded
	RequestResultUnknown
)

// RequestResult represents the result on chain of a request (eg successful or failed with its error )
type RequestResult struct {
	Status        RequestResultStatus
	FailureReason string
}

// EnclaveKeySetRecovery contains the data needed to recover the EnclaveKeySet.
type EnclaveKeySetRecovery struct {
	// RecoveryType is the type of recovery data.
	RecoveryType int `json:"recoveryType"`
	// KeySetCiphertext is the encrypted EnclaveKeySet.
	KeySetCiphertext []byte `json:"keySetCiphertext"`
	// RecoveryCiphertext is the cryptographic data needed to recover the EnclaveKeySet.
	RecoveryCiphertext []byte `json:"recoveryCiphertext"`
}
