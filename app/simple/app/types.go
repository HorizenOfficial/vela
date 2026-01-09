package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"unsafe"

	"github.com/horizen-pes/app/simple/utils"
)

const AddressLength = 20

const (
	MaxBigIntBytes  = 64
	MaxAddressBytes = AddressLength
)

type Address [AddressLength]byte

// HexToAddress converts a hex string to Address with validation.
func HexToAddress(s string) (Address, error) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	if len(s) != AddressLength*2 {
		return Address{}, fmt.Errorf("invalid address length: got %d hex chars, want %d", len(s), AddressLength*2)
	}

	data, err := hex.DecodeString(s)
	if err != nil {
		return Address{}, err
	}

	var address Address
	copy(address[:], data)
	return address, nil
}

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > AddressLength {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// Bytes returns a copy of address bytes (safe, immutable to caller)
func (a Address) Bytes() []byte {
	b := make([]byte, AddressLength)
	copy(b, a[:])
	return b
}

func (a Address) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	addr, err := HexToAddress(s)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	*a = addr
	return nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Hex())
}

func (a Address) String() string {
	return a.Hex()
}

// SerializeAndWriteResult handles common serialization and returns a WASM pointer.
func SerializeAndWriteResult(result any) *byte {
	reportJSON, err := json.Marshal(result)
	if err != nil {
		return utils.StringToPtr([]byte(WasmSerializationError))
	}
	return utils.StringToPtr(reportJSON)
}

// PtrToNonNegativeBigInt converts a WASM pointer and length representing the a big.Int value to a Go big.Int pointer.
// The byte slice is obtained with the (big.Int).Bytes() method, i.e. it represents the absolute value in big-endian byte order, so the value is always non-negative.
func PtrToNonNegativeBigInt(ptr *byte, length int32) *big.Int {
	if ptr == nil || length <= 0 || length > MaxBigIntBytes {
		return big.NewInt(0)
	}
	b := unsafe.Slice(ptr, length)
	return new(big.Int).SetBytes(b)
}

// PtrToAddress converts a WASM pointer and length to a ethereum address.
func PtrToAddress(ptr *byte, length int32) *Address {
	if ptr == nil || length <= 0 || length > MaxAddressBytes {
		return nil
	}
	var address Address
	address.SetBytes(unsafe.Slice(ptr, length))
	return &address
}

// ----- module internal types

// AccountState represents the state of a user account
type AccountState struct {
	Address Address  `json:"address"`
	Balance *big.Int `json:"balance"`
}

// ApplicationInternalState represents the internal state of the application
type ApplicationInternalState struct {
	AppID    int64                    `json:"appId"`
	Accounts map[string]*AccountState `json:"accounts"`
}

// WithdrawInstruction represents instructions for withdrawing funds
type WithdrawInstruction struct {
	To     Address  `json:"to"`
	Amount *big.Int `json:"amount"`
}

type CompareInstructions struct {
	TargetAddress Address `json:"targetAddress"`
}

// PayloadInstructions represents the deserialized payload instructions
type PayloadInstructions struct {
	Type            string               `json:"type"`
	CompareAccounts *CompareInstructions `json:"compare,omitempty"`
	Withdraw        *WithdrawInstruction `json:"withdraw,omitempty"`
}

// ReportPayloadInstructions represent a specific information on how to generate a report
// In this simple app its a custom tag to add to the report
type ReportPayloadInstructions struct {
	IncludeTag string `json:"tag,omitempty"`
}

// DeanonymizationReport represents the structure of the deanonymization report.
type DeanonymizationReport struct {
	Tag      string                   `json:"tag,omitempty"`
	Accounts map[string]*AccountState `json:"accounts"`
}

// --- Local replacements for Host types ---
// This is a deliberate design choice required by the WebAssembly architecture.
// The application communicates by serializing the host-side struct to JSON, passing it to the
// Wasm module, which then deserializes it into its own identical local struct.
// This maintains a clean separation between the two environments.
// The Wasm module is a separate, sandboxed program and should not import types directly from
// the host application's packages, even if they are defined exacltly the same way.
// Moreover we do use analogous but different types, for instance ethereum addresses in the Host
// and [20]byte array type in the guest (this is because tinygo does not support the full standard
// go runtime needed by go-ethereum).
// ---
// TODO: bigInt (math/big pkg) as of now is not fully supported by tinygo, so far we did not experience
// errors, but to be on the safe side we should consider using a different type in the Guest application,
// for instance 8xuint32 or 4xuint64 structs representing uint256 values

// LoadModuleResult is a local replacemente for wasmCommon.LoadModuleResult
type LoadModuleResult struct {
	State []byte   `json:"state"`
	Fuel  *big.Int `json:"fuel"`
	Error string   `json:"error,omitempty"`
}

// DepositResult is a local replacement for wasmCommon.DepositResult
type DepositResult struct {
	State  []byte       `json:"state"`
	Events []PlainEvent `json:"events"`
	Fuel   *big.Int     `json:"fuel"`
	Error  string       `json:"error,omitempty"`
}

// ProcessResult is a local replacement for wasmCommon.ProcessResult
type ProcessResult struct {
	State       []byte       `json:"state"`
	Events      []PlainEvent `json:"events"`
	Withdrawals []Withdrawal `json:"withdrawals"`
	Fuel        *big.Int     `json:"fuel"`
	Error       string       `json:"error,omitempty"`
}

// DeanonymizationResult is a local replacement for wasmCommon.DeanonymizationResult
type DeanonymizationResult struct {
	Report []byte   `json:"report"`
	Fuel   *big.Int `json:"fuel"`
	Error  string   `json:"error,omitempty"`
}

// PlainEvent is a local replacement for common.PlainEvent
type PlainEvent struct {
	UserID       Address `json:"userId"`
	EventSubType string  `json:"eventSubType"`
	Data         []byte  `json:"data"`
}

// Withdrawal is a local replacement for common.Withdrawal
type Withdrawal struct {
	DestinationAddress Address  `json:"destinationAddress"`
	Amount             *big.Int `json:"amount"`
}

type DepositEvent struct {
	Type    string   `json:"type"`
	Amount  *big.Int `json:"amount"`
	Balance *big.Int `json:"balance"`
	Nonce   uint64   `json:"nonce"`
}

type SenderEvent struct {
	Type    string   `json:"type"`
	To      Address  `json:"to"`
	Amount  *big.Int `json:"amount"`
	Balance *big.Int `json:"balance"`
	Nonce   uint64   `json:"nonce"`
}

type RecipientEvent struct {
	Type    string   `json:"type"`
	From    Address  `json:"from"`
	Amount  *big.Int `json:"amount"`
	Balance *big.Int `json:"balance"`
	Nonce   uint64   `json:"nonce"`
}

// WithdrawalEvent is a local replacement for wasmCommon.WithdrawalEvent
type WithdrawalEvent SenderEvent

const (
	WasmSerializationError = `{"error":"wasm serialization error"}`
)

type MemoryStats struct {
	MapSize              int64 `json:"mapSize"`
	CumulativeMemorySize int64 `json:"cumulativeMemorySize"`
}
