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
	msgArgs         abi.Arguments
	eventsArgs      abi.Arguments
	withdrawalsArgs abi.Arguments
}

type withdrawalTuple struct {
	Receiver ethCommon.Address
	Amount   *big.Int
}

func NewMsgToSignBuilder() *MsgToSignBuilder {
	uint256Type, _ := abi.NewType("uint256", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	bytesArrayType, _ := abi.NewType("bytes[]", "", nil)
	eventsArgs := abi.Arguments{{Type: bytesArrayType}}
	WithdrawalRequestArrayType, _ := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{Name: "receiver", Type: "address"},
		{Name: "amount", Type: "uint256"},
	})

	withdrawalsArgs := abi.Arguments{{Type: WithdrawalRequestArrayType}}
	msgArgs := abi.Arguments{
		{Type: uint256Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: uint256Type},
		{Type: bytes32Type},
		{Type: bytes32Type},
	}

	msgBuilder := &MsgToSignBuilder{msgArgs: msgArgs, eventsArgs: eventsArgs, withdrawalsArgs: withdrawalsArgs}
	return msgBuilder
}

func (b *MsgToSignBuilder) BuildMsgHash(updatePayload *common.UpdatePayload) ([]byte, error) {

	reqId, ok := common.StringToBigInt(updatePayload.RequestID)
	if !ok {
		return nil, fmt.Errorf("invalid request ID: %s", updatePayload.RequestID)
	}

	appId, ok := common.StringToBigInt(updatePayload.ApplicationID)
	if !ok {
		return nil, fmt.Errorf("invalid application ID: %s", updatePayload.ApplicationID)
	}

	events := make([][]byte, len(updatePayload.Events))
	for i, event := range updatePayload.Events {
		events[i] = event.EncryptedData
	}

	encodedEvents, err := b.eventsArgs.Pack(events)
	if err != nil {
		return nil, fmt.Errorf("failed to encode events: %w", err)
	}
	eventsHash := ethCrypto.Keccak256(encodedEvents)
	var eventArr [32]byte = [32]byte(eventsHash)

	withdrawals := make([]withdrawalTuple, len(updatePayload.Withdrawals))

	for i, withdrawal := range updatePayload.Withdrawals {
		amount := new(big.Int).SetUint64(withdrawal.Amount)

		withdrawals[i] = withdrawalTuple{
			Receiver: ethCommon.HexToAddress(withdrawal.DestinationAddress),
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
		appId,
		updatePayload.PrevStateRoot,
		updatePayload.NewStateRoot,
		reqId,
		eventArr,
		withdrawalArr,
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
