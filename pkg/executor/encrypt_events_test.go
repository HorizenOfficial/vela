package executor

import (
	"context"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/appdata"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
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

// When PlainEvent.RecipientPubKey is set, encryptEvents must use it as the
// ECIES receiver and ignore keyStore[UserID] entirely. The encrypted payload
// must round-trip with the *explicit* recipient's private key — proving the
// explicit-recipient path was taken, not the keyStore fallback.
func TestEncryptEvents_ExplicitRecipientPubKey_UsedAndRoundTrips(t *testing.T) {
	executorKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	// Explicit recipient — NOT in keyStore.
	explicitKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	// A different key registered in keyStore for the same UserID — to prove
	// the explicit path is taken, not the keyStore fallback. If the fallback
	// fired, the ciphertext would only decrypt with this keyStore key.
	wrongKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	userAddr := ethCommon.HexToAddress("0xadd0000000000000000000000000000000000002")
	keyStore := appdata.KeyStore{
		userAddr: wrongKey.PublicKey(),
	}

	plaintext := []byte("hint-for-scheduler")
	plain := []common.PlainEvent{{
		UserID:          userAddr,
		EventSubType:    [32]byte{0xCC},
		Data:            plaintext,
		RecipientPubKey: explicitKey.PublicKey().Bytes(),
	}}

	executor := &StatelessExecutor{log: testLogger}
	encrypted, failure, err := executor.encryptEvents(context.Background(), plain, common.NewApplicationId(1), executorKey, nil, keyStore, appdata.SeedStore{})
	require.Nil(t, failure)
	require.Nil(t, err)
	require.Len(t, encrypted, 1)
	require.NotEmpty(t, encrypted[0].EncryptedData)

	// Decrypt with the explicit recipient's private key — must succeed.
	decrypted, err := crypto.Decrypt(executorKey.PublicKey(), explicitKey, encrypted[0].EncryptedData)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)

	// Decrypt with the keyStore-registered (wrong) key — must fail. This proves
	// the keyStore fallback was NOT taken.
	_, err = crypto.Decrypt(executorKey.PublicKey(), wrongKey, encrypted[0].EncryptedData)
	require.Error(t, err)
}

// A non-empty but malformed RecipientPubKey must surface as a typed
// CodeParsingKeyError failure, not silently fall back to keyStore.
func TestEncryptEvents_InvalidRecipientPubKey_ReturnsParsingKeyError(t *testing.T) {
	executorKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	userAddr := ethCommon.HexToAddress("0xadd0000000000000000000000000000000000003")
	// keyStore has a valid entry for this UserID — but we should NOT fall back
	// to it; a malformed explicit key must be treated as an error.
	validKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	keyStore := appdata.KeyStore{
		userAddr: validKey.PublicKey(),
	}

	plain := []common.PlainEvent{{
		UserID:          userAddr,
		EventSubType:    [32]byte{0xDD},
		Data:            []byte("payload"),
		RecipientPubKey: []byte{0x04, 0xDE, 0xAD, 0xBE, 0xEF}, // wrong length, definitely not a valid P521 SEC1 point
	}}

	executor := &StatelessExecutor{log: testLogger}
	encrypted, failure, err := executor.encryptEvents(context.Background(), plain, common.NewApplicationId(1), executorKey, nil, keyStore, appdata.SeedStore{})
	require.Nil(t, err)
	require.Nil(t, encrypted)
	require.NotNil(t, failure)
	require.Equal(t, apperrors.CodeParsingKeyError.Code, failure.RequestError.Code)
}

// When RecipientPubKey is nil/empty, the existing keyStore[UserID] lookup must
// run unchanged — regression coverage for the vela-nova flow.
func TestEncryptEvents_NilRecipientPubKey_FallsBackToKeyStore(t *testing.T) {
	executorKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	userAddr := ethCommon.HexToAddress("0xadd0000000000000000000000000000000000004")
	keyStore := appdata.KeyStore{
		userAddr: userKey.PublicKey(),
	}

	plaintext := []byte("vela-nova-style-event")
	plain := []common.PlainEvent{{
		UserID:       userAddr,
		EventSubType: [32]byte{0xEE},
		Data:         plaintext,
		// RecipientPubKey intentionally unset (nil).
	}}

	executor := &StatelessExecutor{log: testLogger}
	encrypted, failure, err := executor.encryptEvents(context.Background(), plain, common.NewApplicationId(1), executorKey, nil, keyStore, appdata.SeedStore{})
	require.Nil(t, failure)
	require.Nil(t, err)
	require.Len(t, encrypted, 1)

	// The encrypted payload must decrypt with the keyStore-registered user key.
	decrypted, err := crypto.Decrypt(executorKey.PublicKey(), userKey, encrypted[0].EncryptedData)
	require.NoError(t, err)
	require.Equal(t, plaintext, decrypted)
}
