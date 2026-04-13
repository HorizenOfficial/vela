// Package common provide the data structures used in the application.
package common

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

// Type aliases for types that moved to vela-common-go/common.
type ApplicationIdType = velacommon.ApplicationIdType
type RequestIdType = velacommon.RequestIdType
type RequestResultStatus = velacommon.RequestResultStatus

const (
	RequestResultOK      = velacommon.RequestResultOK
	RequestResultFailed  = velacommon.RequestResultFailed
	RequestResultUnknown = velacommon.RequestResultUnknown
)

func NewApplicationId(id uint64) ApplicationIdType {
	return velacommon.NewApplicationId(id)
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
	Timestamp *Big `json:"timestamp"`
	// Sender is the address of the sender
	Sender ethCommon.Address `json:"sender"`
	// TokenAddress is the address of the business asset token (0x0 = ETH)
	TokenAddress ethCommon.Address `json:"tokenAddress"`
	// AssetAmount is the business asset amount (replaces DepositAmount)
	AssetAmount *Big `json:"assetAmount"`
	// MaxFeeValue is the maximum fee value reserved for fee payment (always ETH)
	MaxFeeValue *Big `json:"maxFeeValue"`
}

func (r *Request) Validate() error {
	if err := validateBigInt("timestamp", r.Timestamp.ToInt(), false); err != nil {
		return err
	}

	if err := validateBigInt("assetAmount", r.AssetAmount.ToInt(), true); err != nil {
		return err
	}

	// If assetAmount is zero, tokenAddress must be the zero address
	if r.AssetAmount.ToInt().Sign() == 0 && r.TokenAddress != (ethCommon.Address{}) {
		return fmt.Errorf("tokenAddress must be zero address when assetAmount is zero")
	}

	if err := validateBigInt("maxFeeValue", r.MaxFeeValue.ToInt(), true); err != nil {
		return err
	}
	return nil
}

// Event represents an event to be emitted
type Event struct {
	// ApplicationID is the ID of the application
	ApplicationID ApplicationIdType `json:"applicationId"`
	// UserID is the ID of the user associated with the event
	UserID ethCommon.Address `json:"userId"`
	// EventSubType is the optional subtype used for filtering
	EventSubType string `json:"eventSubType"`
	// EncryptedData is the encrypted event data
	EncryptedData []byte `json:"encryptedData"`
}

func (e Event) String() string {
	return fmt.Sprintf("Event{ApplicationID: %d, UserID: %s, EventSubType: %s, EncryptedData: %s}",
		e.ApplicationID, e.UserID.Hex(), e.EventSubType, hex.EncodeToString(e.EncryptedData))
}

// Withdrawal represents a withdrawal from the system
type Withdrawal struct {
	// TokenAddress is the address of the token to withdraw (0x0 = ETH)
	TokenAddress ethCommon.Address `json:"tokenAddress"`
	// DestinationAddress is the address to send the funds to
	DestinationAddress ethCommon.Address `json:"destinationAddress"`
	// Amount is the amount to withdraw
	Amount *Big `json:"amount"`
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
	// RefundAmount is the amount to refund in WEI
	RefundAmount *Big `json:"refundAmount"`
	// ApplicationFee is the fee charged for the application in WEI
	ApplicationFee *Big `json:"applicationFee"`
	// ErrorCode is the error code (0 for success, non-zero for error)
	ErrorCode uint8 `json:"errorCode"`
	// ErrorMsg is the error message (empty for success)
	ErrorMsg string `json:"errorMsg"`

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
	// Authority is the entity requesting the report
	Authority ethCommon.Address `json:"authority"`
}

// DecryptedReport represents a decrypted deanonymization report
type DecryptedReport struct {
	ApplicationID   ApplicationIdType `json:"applicationId"`
	RequestID       RequestIdType     `json:"requestId"`
	ReportDataBytes []byte            `json:"reportDataBytes"`
}

// PlainEvent represents an emitted event before encryption.
type PlainEvent struct {
	// UserID is the address of the user associated with the event
	UserID ethCommon.Address `json:"userId"`
	// EventSubType is the optional subtype used for filtering
	EventSubType string `json:"eventSubType"`
	// Data is the encrypted event data
	Data []byte `json:"data"`
}

// RequestResult represents the result on chain of a request (eg successful or failed with its error)
type RequestResult struct {
	Status       RequestResultStatus
	ErrorCode    uint8
	ErrorMessage string
}

// RecoveryType represents the keyset recovery mechanism.
type RecoveryType int

const (
	RecoveryTypeUnsafe RecoveryType = 0
	RecoveryTypeKMS    RecoveryType = 1
)

// EnclaveKeySetRecovery contains the data needed to recover the EnclaveKeySet.
type EnclaveKeySetRecovery struct {
	// RecoveryType is the type of recovery data.
	RecoveryType RecoveryType `json:"recoveryType"`
	// KeySetCiphertext is the encrypted EnclaveKeySet.
	KeySetCiphertext []byte `json:"keySetCiphertext"`
	// RecoveryCiphertext is the cryptographic data needed to recover the EnclaveKeySet.
	RecoveryCiphertext []byte `json:"recoveryCiphertext"`
}

type ChannelConnectionParams interface {
	IsChannelConnectionParams()
}

type VSockChannelConnectionParams struct {
	CID  uint32
	Port uint32
}

type CommunicationParams struct {
	RequestTimeoutSec time.Duration
}

func (VSockChannelConnectionParams) IsChannelConnectionParams() {}

// NewVSockChannelConnectionParams parses "cid:port"
func NewVSockChannelConnectionParams(s string) (*VSockChannelConnectionParams, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid vsock address %q, expected cid:port", s)
	}

	// Parse CID
	cid64, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid CID %q: %w", parts[0], err)
	}

	// Parse port
	port64, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid port %q: %w", parts[1], err)
	}

	return &VSockChannelConnectionParams{
		CID:  uint32(cid64),
		Port: uint32(port64),
	}, nil
}

type TcpChannelConnectionParams struct {
	Ip   string
	Port uint32
}

func (TcpChannelConnectionParams) IsChannelConnectionParams() {}

func (f TcpChannelConnectionParams) Url() string {
	return fmt.Sprintf("%s:%d", f.Ip, f.Port)
}

// NewTcpChannelConnectionParams parses "ip:port"
func NewTcpChannelConnectionParams(addr string) (*TcpChannelConnectionParams, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid address %q, expected ip:port", addr)
	}

	// Parse port
	port64, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid port %q: %w", parts[1], err)
	}

	return &TcpChannelConnectionParams{
		Ip:   parts[0],
		Port: uint32(port64),
	}, nil
}
