package executor

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/common"
)

type MsgToSignBuilder struct {
	msgArgs           abi.Arguments
	eventsArgs        abi.Arguments
	eventSubTypesArgs abi.Arguments
	withdrawalsArgs   abi.Arguments
}

type withdrawalTuple struct {
	Receiver ethCommon.Address
	Amount   *big.Int
}

func NewMsgToSignBuilder() (*MsgToSignBuilder, error) {
	bytesArrayType, err := abi.NewType("bytes[]", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating byte array type: %w", err)
	}

	eventsArgs := abi.Arguments{{Type: bytesArrayType}}

	stringArrayType, err := abi.NewType("string[]", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating string array type: %w", err)
	}
	eventSubTypesArgs := abi.Arguments{{Type: stringArrayType}}

	WithdrawalRequestArrayType, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "receiver", Type: "address"},
		{Name: "amount", Type: "uint256"},
	})
	if err != nil {
		return nil, fmt.Errorf("failure creating withdrawal array type: %w", err)
	}
	withdrawalsArgs := abi.Arguments{{Type: WithdrawalRequestArrayType}}

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
	boolType, err := abi.NewType("bool", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failure creating bool type: %w", err)
	}
	msgArgs := abi.Arguments{
		{Type: uint64Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: uint256Type},
		{Type: uint256Type},
		{Type: boolType},
	}

	msgBuilder := &MsgToSignBuilder{msgArgs: msgArgs, eventsArgs: eventsArgs, eventSubTypesArgs: eventSubTypesArgs, withdrawalsArgs: withdrawalsArgs}
	return msgBuilder, nil
}

func (b *MsgToSignBuilder) BuildMsgHash(updatePayload *common.UpdatePayload) ([]byte, error) {

	events := make([][]byte, len(updatePayload.Events))
	eventSubTypes := make([]string, len(updatePayload.Events))
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

	withdrawals := make([]withdrawalTuple, len(updatePayload.Withdrawals))

	for i, withdrawal := range updatePayload.Withdrawals {
		amount := withdrawal.Amount

		withdrawals[i] = withdrawalTuple{
			Receiver: withdrawal.DestinationAddress,
			Amount:   amount,
		}
	}

	encodedWithdrawal, err := b.withdrawalsArgs.Pack(withdrawals)
	if err != nil {
		return nil, fmt.Errorf("failed to encode withdrawals: %w", err)
	}
	withdrawalHash := ethCrypto.Keccak256(encodedWithdrawal)
	var withdrawalArr [32]byte = [32]byte(withdrawalHash)

	values := []interface{}{
		updatePayload.ApplicationID,
		updatePayload.PrevStateRoot,
		updatePayload.NewStateRoot,
		updatePayload.RequestID,
		eventArr,
		eventSubTypesArr,
		withdrawalArr,
		updatePayload.RefundAmount,
		updatePayload.ApplicationFee,
		updatePayload.ReportGenerated,
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
