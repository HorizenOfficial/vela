package executor

import (
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/stretchr/testify/require"
)

// Verify that EventSubType is propagated into the signature hash.
func TestMsgToSignBuilder_IncludesEventSubType(t *testing.T) {
	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)

	// Base payload
	base := &common.UpdatePayload{
		ApplicationID:  common.NewApplicationId(1),
		RequestID:      common.RequestIdType{},
		PrevStateRoot:  [32]byte{1},
		NewStateRoot:   [32]byte{2},
		Events:         []common.Event{{EventSubType: "a", EncryptedData: []byte{0x01}}},
		Withdrawals:    nil,
		RefundAmount:   common.NewBig(0),
		ApplicationFee: common.NewBig(0),
	}

	h1, err := builder.BuildMsgHash(base)
	require.NoError(t, err)

	// Change only the subtype; hash must change.
	changed := *base
	changed.Events = []common.Event{{EventSubType: "b", EncryptedData: []byte{0x01}}}

	h2, err := builder.BuildMsgHash(&changed)
	require.NoError(t, err)

	require.NotEqual(t, h1, h2, "hash should change when event subtype changes")
}

// Verify that WithdrawalRequest.TokenAddress is included in the signature hash.
func TestMsgToSignBuilder_IncludesWithdrawalTokenAddress(t *testing.T) {
	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)

	ethAddr := ethCommon.HexToAddress("0x1234567890123456789012345678901234567890")
	tokenA := ethCommon.Address{}                                                       // ETH
	tokenB := ethCommon.HexToAddress("0xdead000000000000000000000000000000000001") // ERC-20

	base := &common.UpdatePayload{
		ApplicationID: common.NewApplicationId(1),
		PrevStateRoot: [32]byte{1},
		NewStateRoot:  [32]byte{2},
		Withdrawals: []common.Withdrawal{
			{TokenAddress: tokenA, DestinationAddress: ethAddr, Amount: common.NewBig(1000)},
		},
		RefundAmount:   common.NewBig(0),
		ApplicationFee: common.NewBig(0),
	}

	h1, err := builder.BuildMsgHash(base)
	require.NoError(t, err)

	// Change only the token address; hash must change.
	changed := *base
	changed.Withdrawals = []common.Withdrawal{
		{TokenAddress: tokenB, DestinationAddress: ethAddr, Amount: common.NewBig(1000)},
	}

	h2, err := builder.BuildMsgHash(&changed)
	require.NoError(t, err)

	require.NotEqual(t, h1, h2, "hash should change when withdrawal token address changes")
}

// Verify that an empty withdrawal list with token-aware tuples still hashes correctly.
func TestMsgToSignBuilder_EmptyWithdrawals(t *testing.T) {
	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)

	payload := &common.UpdatePayload{
		ApplicationID:  common.NewApplicationId(1),
		PrevStateRoot:  [32]byte{1},
		NewStateRoot:   [32]byte{2},
		Withdrawals:    []common.Withdrawal{},
		RefundAmount:   common.NewBig(0),
		ApplicationFee: common.NewBig(0),
	}

	h, err := builder.BuildMsgHash(payload)
	require.NoError(t, err)
	require.NotEmpty(t, h)
}
