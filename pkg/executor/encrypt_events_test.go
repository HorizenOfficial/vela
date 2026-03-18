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

	plain := []common.PlainEvent{{
		UserID:       userAddr,
		EventSubType: "test_subtype",
		Data:         []byte("hello"),
	}}

	executor := &StatelessExecutor{log: testLogger}
	encrypted, failure := executor.encryptEvents(context.Background(), plain, common.NewApplicationId(1), privKey, nil, keyStore, appdata.SeedStore{})
	require.Nil(t, failure)
	require.Len(t, encrypted, 1)
	require.Equal(t, "test_subtype", encrypted[0].EventSubType)
	require.Equal(t, userAddr, encrypted[0].UserID)
	require.NotEmpty(t, encrypted[0].EncryptedData)
}
