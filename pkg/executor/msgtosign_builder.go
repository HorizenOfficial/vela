package executor

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/HorizenOfficial/vela/pkg/common"
)

type MsgToSignBuilder struct {
	msgArgs           abi.Arguments
	eventsArgs        abi.Arguments
	eventSubTypesArgs abi.Arguments
	withdrawalsArgs   abi.Arguments
}

type withdrawalTuple struct {
	TokenAddress ethCommon.Address
	Receiver     ethCommon.Address
	Amount       *big.Int
}

func NewMsgToSignBuilder() (*MsgToSignBuilder, error) {
	bytesArrayType, err := abi.NewType("bytes[]", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating byte array type: %w", err)
	}

	eventsArgs := abi.Arguments{{Type: bytesArrayType}}

	bytes32ArrayType, err := abi.NewType("bytes32[]", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating bytes32 array type: %w", err)
	}
	eventSubTypesArgs := abi.Arguments{{Type: bytes32ArrayType}}

	WithdrawalRequestArrayType, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "tokenAddress", Type: "address"},
		{Name: "receiver", Type: "address"},
		{Name: "amount", Type: "uint256"},
	})
	if err != nil {
		return nil, fmt.Errorf("failure creating withdrawal array type: %w", err)
	}
	withdrawalsArgs := abi.Arguments{{Type: WithdrawalRequestArrayType}}

	uint8Type, err := abi.NewType("uint8", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating uint8 type: %w", err)
	}
	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating string type: %w", err)
	}
	uint64Type, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating uint64 type: %w", err)
	}
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating uint256 type: %w", err)
	}
	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating bytes32 type: %w", err)
	}
	// Order must match UpdateEntryHash.entryHash in Solidity
	msgArgs := abi.Arguments{
		{Type: uint64Type},  // applicationId
		{Type: bytes32Type}, // prevStateRoot
		{Type: bytes32Type}, // newStateRoot
		{Type: bytes32Type}, // processedRequestId
		{Type: bytes32Type}, // userEventsHash
		{Type: bytes32Type}, // userEventSubTypesHash
		{Type: bytes32Type}, // appEventsHash
		{Type: bytes32Type}, // appEventSubTypesHash
		{Type: bytes32Type}, // withdrawalRequestsHash
		{Type: uint256Type}, // refundAmount
		{Type: uint256Type}, // applicationFee
		{Type: uint8Type},   // errorCode
		{Type: stringType},  // errorMsg
	}

	msgBuilder := &MsgToSignBuilder{msgArgs: msgArgs, eventsArgs: eventsArgs, eventSubTypesArgs: eventSubTypesArgs, withdrawalsArgs: withdrawalsArgs}
	return msgBuilder, nil
}

// BuildMsgHash builds the Ethereum personal_sign hash for a single update
// payload: TextHash(entryHash). This is what individually signed payloads sign.
// ProcessorEndpoint.stateUpdate verifies it through
// AbstractTeeAuthenticator.checkBatchSignature with a 1-entry array, whose
// digest is byte-identical to this one (see BuildBatchMsgHash).
func (b *MsgToSignBuilder) BuildMsgHash(updatePayload *common.UpdatePayload) ([]byte, error) {
	entryHash, err := b.buildEntryHash(updatePayload)
	if err != nil {
		return nil, err
	}
	return accounts.TextHash(entryHash), nil
}

// BuildBatchMsgHash builds the Ethereum personal_sign hash covering all batch
// entries: TextHash(entryHash_0 || entryHash_1 || ... || entryHash_N-1).
//
// The entry hashes are concatenated and TextHash-ed directly, WITHOUT an
// intermediate keccak256 over the concatenation. Two properties make this
// unambiguous and safe:
//   - Injectivity: every entry hash is a fixed 32-byte keccak256 output, so a
//     concatenation of N of them splits exactly one way. This relies on the
//     per-entry hash staying fixed-length — do not make buildEntryHash variable.
//   - Length binding: the personal_sign prefix commits to the total byte length
//     (32*N), so batches of different sizes can never collide.
//
// A consequence is that a single-entry batch hashes identically to BuildMsgHash of
// that entry, so single-request and batch submission share one signing scheme (and
// a 1-entry batch signature verifies on the single-request stateUpdate() path).
// The contract side must reconstruct the same digest: TextHash of the concatenated
// entry hashes using a DYNAMIC length prefix (32*N), not a fixed length.
func (b *MsgToSignBuilder) BuildBatchMsgHash(updatePayloads []*common.UpdatePayload) ([]byte, error) {
	if len(updatePayloads) == 0 {
		return nil, fmt.Errorf("no payloads to hash")
	}

	concatenated := make([]byte, 0, 32*len(updatePayloads))
	for i, payload := range updatePayloads {
		entryHash, err := b.buildEntryHash(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to hash batch entry %d: %w", i, err)
		}
		concatenated = append(concatenated, entryHash...)
	}

	return accounts.TextHash(concatenated), nil
}

// buildEntryHash computes keccak256 of the ABI-encoded payload fields, without
// the Ethereum message prefix. It is the per-entry hash: prefixed and signed
// directly for single payloads, or aggregated into the batch hash.
func (b *MsgToSignBuilder) buildEntryHash(updatePayload *common.UpdatePayload) ([]byte, error) {

	events := make([][]byte, len(updatePayload.Events))
	eventSubTypes := make([][32]byte, len(updatePayload.Events))
	for i, event := range updatePayload.Events {
		events[i] = event.EncryptedData
		eventSubTypes[i] = event.EventSubType
	}

	encodedEvents, err := b.eventsArgs.Pack(events)
	if err != nil {
		return nil, fmt.Errorf("failed to encode events: %w", err)
	}
	eventsHash := ethCrypto.Keccak256(encodedEvents)
	var eventArr [32]byte = [32]byte(eventsHash)

	encodedEventSubTypes, err := b.eventSubTypesArgs.Pack(eventSubTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode event subtypes: %w", err)
	}
	eventSubTypesHash := ethCrypto.Keccak256(encodedEventSubTypes)
	var eventSubTypesArr [32]byte = [32]byte(eventSubTypesHash)

	appEvents := make([][]byte, len(updatePayload.AppEvents))
	appEventSubTypes := make([][32]byte, len(updatePayload.AppEvents))
	for i, appEvent := range updatePayload.AppEvents {
		appEvents[i] = appEvent.Data
		appEventSubTypes[i] = appEvent.EventSubType
	}

	encodedAppEvents, err := b.eventsArgs.Pack(appEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to encode app events: %w", err)
	}
	appEventsHash := ethCrypto.Keccak256(encodedAppEvents)
	var appEventsArr [32]byte = [32]byte(appEventsHash)

	encodedAppEventSubTypes, err := b.eventSubTypesArgs.Pack(appEventSubTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode app event subtypes: %w", err)
	}
	appEventSubTypesHash := ethCrypto.Keccak256(encodedAppEventSubTypes)
	var appEventSubTypesArr [32]byte = [32]byte(appEventSubTypesHash)

	withdrawals := make([]withdrawalTuple, len(updatePayload.Withdrawals))

	for i, withdrawal := range updatePayload.Withdrawals {
		amount := withdrawal.Amount.ToInt()

		withdrawals[i] = withdrawalTuple{
			TokenAddress: withdrawal.TokenAddress,
			Receiver:     withdrawal.DestinationAddress,
			Amount:       amount,
		}
	}

	encodedWithdrawal, err := b.withdrawalsArgs.Pack(withdrawals)
	if err != nil {
		return nil, fmt.Errorf("failed to encode withdrawals: %w", err)
	}
	withdrawalHash := ethCrypto.Keccak256(encodedWithdrawal)
	var withdrawalArr [32]byte = [32]byte(withdrawalHash)

	// Order must match msgArgs above and UpdateEntryHash.entryHash
	values := []interface{}{
		updatePayload.ApplicationID,       // applicationId
		updatePayload.PrevStateRoot,       // prevStateRoot
		updatePayload.NewStateRoot,        // newStateRoot
		updatePayload.RequestID,           // processedRequestId
		eventArr,                          // userEventsHash
		eventSubTypesArr,                  // userEventSubTypesHash
		appEventsArr,                      // appEventsHash
		appEventSubTypesArr,               // appEventSubTypesHash
		withdrawalArr,                     // withdrawalRequestsHash
		updatePayload.RefundAmount.ToInt(),  // refundAmount
		updatePayload.ApplicationFee.ToInt(), // applicationFee
		updatePayload.ErrorCode,           // errorCode
		updatePayload.ErrorMsg,            // errorMsg
	}

	// Encoding parameters
	encoded, err := b.msgArgs.Pack(values...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode parameters: %w", err)
	}

	return ethCrypto.Keccak256(encoded), nil
}
