// Package common provide the data structures used in the application.
package common

import (
	"crypto/ecdh"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

// RequestType represents the type of request being sent to the TEE
type RequestType string

const (
	// Deploy represents a request to deploy a new application
	Deploy RequestType = "deploy"
	// Process is used for processing a batch of requests
	Process RequestType = "process"
	// Deanonymize is used for deanonymization requests
	Deanonymize RequestType = "deanonymize"
)

// PublicKeyP521 is a public key with Elliptic Curve Diffie-Hellman over NIST P-521 curve, also known as secp521r1.
// (This can be used for encryption/decription)
type PublicKeyP521 struct {
	*ecdh.PublicKey
}

// PrivateKeyP521 is a private key with Elliptic Curve Diffie-Hellman over NIST P-521 curve, also known as secp521r1.
// (This can be used for encryption/decription)
type PrivateKeyP521 struct {
	*ecdh.PrivateKey
}

// PublicKey returns the public key associated with a private key.
func (p *PrivateKeyP521) PublicKey() *PublicKeyP521 {
	return &PublicKeyP521{p.PrivateKey.PublicKey()}
}

// PublicKey25519 is a public key with Elliptic Curve 25519.
type PublicKey25519 struct {
	*ecdh.PublicKey
}

// PrivateKey25519 is a private key with with Elliptic Curve 25519.
type PrivateKey25519 struct {
	*ecdh.PrivateKey
}

// PublicKey returns the public key associated with a private key.
func (p *PrivateKey25519) PublicKey() *PublicKey25519 {
	return &PublicKey25519{p.PrivateKey.PublicKey()}
}

// PublicKeySecp256k1 is a public key with Elliptic Curve secp256k1.
// (This is the curve used in Bitcoin and Ethereum)
type PublicKeySecp256k1 struct {
	*ecdsa.PublicKey
}

// Address returns the Ethereum address of the public key.
func (p *PublicKeySecp256k1) Address() string {
	return crypto.PubkeyToAddress(*p.PublicKey).Hex()
}

// PrivateKeySecp256k1 is a private key with Elliptic Curve secp256k1.
// (This is the curve used in Bitcoin and Ethereum)
type PrivateKeySecp256k1 struct {
	*ecdsa.PrivateKey
}

// PublicKey returns the public key associated with a private key.
func (p *PrivateKeySecp256k1) PublicKey() *PublicKeySecp256k1 {
	return &PublicKeySecp256k1{&p.PrivateKey.PublicKey}
}

// Sign calculates an Ethereum signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func (p *PrivateKeySecp256k1) Sign(digest []byte) ([]byte, error) {
	return crypto.Sign(digest, p.PrivateKey)
}

// Represents a AES-256-GCM key
type AES256Key [32]byte

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
	// Payload is the encrypted payload for the request
	Payload []byte `json:"payload"`
	// Timestamp is the time the request was submitted
	Timestamp int64 `json:"timestamp"`
	// Sender is the address of the sender
	Sender string `json:"sender"`
	// Value is the optional deposit value in WEI
	Value int64 `json:"value"`
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
	Amount string `json:"amount"`
}

// UpdatePayload represents an update to the state
type UpdatePayload struct {
	// ApplicationID is the ID of the application
	ApplicationID string `json:"applicationId"`
	// RequestID is the ID of the request being processed
	RequestID string `json:"requestId"`
	// PrevStateRoot is the previous state root
	PrevStateRoot []byte `json:"prevStateRoot"`
	// NewStateRoot is the new state root
	NewStateRoot []byte `json:"newStateRoot"`
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
	StateRoot []byte `json:"stateRoot"`
	// EncryptedState is the encrypted state data
	EncryptedState []byte `json:"encryptedState"`
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
