package executor

import (
	"context"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/appdata"
	"github.com/HorizenOfficial/vela/pkg/crypto"
	"github.com/stretchr/testify/require"
)

// Ensure encryptEvents copies EventSubType and keeps UserID as recipient.
func TestEncryptEventsCopiesSubType(t *testing.T) {
	privKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	userAddr := ethCommon.HexToAddress("0xadd0000000000000000000000000000000000001")
	keyStore := appdata.KeyStore{
		userAddr: privKey.PublicKey(),
	}

	subType := [32]byte{0xAA, 0xBB}
	plain := []common.PlainEvent{{
		UserID:       userAddr,
		EventSubType: subType,
		Data:         []byte("hello"),
	}}

	executor := &StatelessExecutor{log: testLogger}
	encrypted, failure, err := executor.encryptEvents(context.Background(), plain, common.NewApplicationId(1), privKey, nil, keyStore, appdata.SeedStore{})
	require.Nil(t, failure)
	require.Nil(t, err)
	require.Len(t, encrypted, 1)
	require.Equal(t, subType, encrypted[0].EventSubType)
	require.Equal(t, userAddr, encrypted[0].UserID)
	require.NotEmpty(t, encrypted[0].EncryptedData)
}
