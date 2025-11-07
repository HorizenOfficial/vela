package executor

import (
	"math/big"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/blockchain/testutil"
	"github.com/horizen-pes/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestCheckSignature(t *testing.T) {
	applicationId := common.NewApplicationId(1)
	execConfig := DefaultConfig()

	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)

	executor := &StatelessExecutor{
		config:           execConfig,
		MsgToSignBuilder: builder,
	}
	executorAddress := ethCrypto.PubkeyToAddress(execConfig.SignatureKey.PrivateKey.PublicKey)

	testHelper := testutil.NewSimTestHelper(t, false, false, &executorAddress, nil)
	defer testHelper.Close()

	events := [1]common.Event{{ApplicationID: applicationId, EncryptedData: []byte{0x07, 0x07, 0x07}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: ethCommon.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(1)},
	}

	updatePayload := &common.UpdatePayload{
		ApplicationID: applicationId,
		RequestID:     "7890",
		PrevStateRoot: [32]byte{0x08, 0x05, 0x06},
		NewStateRoot:  [32]byte{0x04, 0x05, 0x06},
		Events:        events[:],
		Withdrawals:   withdrawals,
	}

	signature, err := executor.signUpdatePayload(updatePayload)
	require.NoError(t, err)

	// Sign the hash
	updatePayload.Signature = signature
	teeSignerContract := testHelper.GetSimTeeAuthenticatorHelper()

	result := teeSignerContract.CheckSignature(updatePayload)
	require.True(t, result, "Signature verification failed")
}
