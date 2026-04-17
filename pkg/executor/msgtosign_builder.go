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
	// Order must match AbstractTeeAuthenticator.checkSignature in Solidity
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

func (b *MsgToSignBuilder) BuildMsgHash(updatePayload *common.UpdatePayload) ([]byte, error) {

	events := make([][]byte, len(updatePayload.Events))
	eventSubTypes := make([][32]byte, len(updatePayload.Events))
	for i, event := range updatePayload.Events {
		events[i] = event.EncryptedData
		copy(eventSubTypes[i][:], event.EventSubType)
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
		copy(appEventSubTypes[i][:], appEvent.EventSubType)
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

	// Order must match msgArgs above and AbstractTeeAuthenticator.checkSignature
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

	hash := ethCrypto.Keccak256(encoded)

	hash = accounts.TextHash(hash)
	return hash, nil

}
