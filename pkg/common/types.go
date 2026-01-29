// Package common provide the data structures used in the application.
package common

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

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

func (rt RequestIdType) MarshalJSON() ([]byte, error) {
	s := hex.EncodeToString(rt[:])
	return []byte(`"0x` + s + `"`), nil
}

func (rt *RequestIdType) UnmarshalJSON(data []byte) error {
	// data is expected to be a hex string with a "0x" prefix in quotes representing an array of exactly 32 bytes
	// e.g. "0xab12...a8" (68 chars in total, prefix and start-end quotes included)
	if len(data) != 68 || data[0] != '"' || data[1] != '0' || data[2] != 'x' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid RequestIdType format")
	}

	b, err := hex.DecodeString(string(data[3 : len(data)-1]))
	if err != nil {
		return err
	}
	if len(b) != 32 {
		return fmt.Errorf("invalid RequestIdType length")
	}
	copy(rt[:], b)
	return nil
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
	// DepositAmount is the optional deposit value in WEI
	DepositAmount *big.Int `json:"depositAmount"`
	// MaxFeeValue is the maximum fee value reserved for fee payment
	MaxFeeValue *big.Int `json:"maxFeeValue"`
}

func (r *Request) Validate() error {
	if err := validateBigInt("timestamp", r.Timestamp, false); err != nil {
		return err
	}

	if err := validateBigInt("depositAmount", r.DepositAmount, true); err != nil {
		return err
	}

	if err := validateBigInt("maxFeeValue", r.MaxFeeValue, true); err != nil {
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
	// RefundAmount is the amount to refund in WEI
	RefundAmount *big.Int `json:"refundAmount"`
	// ApplicationFee is the fee charged for the application in WEI
	ApplicationFee *big.Int `json:"applicationFee"`
	// ReportGenerated indicates if a deanonymization report was generated
	ReportGenerated bool `json:"reportGenerated"`
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

// RequestResultStatus represents the final status of a request after the execution
type RequestResultStatus uint8

const (
	RequestResultOK RequestResultStatus = iota
	RequestResultFailed
	RequestResultFailedNotRefunded
	RequestResultUnknown
)

// RequestResult represents the result on chain of a request (eg successful or failed with its error)
type RequestResult struct {
	Status       RequestResultStatus
	ErrorCode    uint8
	ErrorMessage string
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
