package executor

import (
	"math/big"
	"testing"

	"github.com/horizen-pes/pkg/common"
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
		RefundAmount:   common.ToBig(big.NewInt(0)),
		ApplicationFee: common.ToBig(big.NewInt(0)),
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
