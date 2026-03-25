package executor

import (
	"context"
	"crypto/ecdsa"
	"testing"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/appdata"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/crypto"
	"github.com/stretchr/testify/require"
)

// --- VerifySeed tests ---

func TestVerifySeed_Valid(t *testing.T) {
	privKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	msgHash := ethCrypto.Keccak256([]byte(SubtypeKeyMessage))
	seed, err := ethCrypto.Sign(msgHash, privKey)
	require.NoError(t, err)
	addr := ethCrypto.PubkeyToAddress(privKey.PublicKey)
	require.NoError(t, VerifySeed(seed, addr))
}

func TestVerifySeed_WrongSigner(t *testing.T) {
	privKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	msgHash := ethCrypto.Keccak256([]byte(SubtypeKeyMessage))
	seed, err := ethCrypto.Sign(msgHash, privKey)
	require.NoError(t, err)

	// Different key → different address
	otherKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	otherAddr := ethCrypto.PubkeyToAddress(otherKey.PublicKey)

	err = VerifySeed(seed, otherAddr)
	require.Error(t, err)
	require.Contains(t, err.Error(), "does not match")
}

func TestVerifySeed_WrongLength(t *testing.T) {
	err := VerifySeed(make([]byte, 64), ethCrypto.PubkeyToAddress(mustGenerateKey(t).PublicKey))
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid seed length")
}

func TestVerifySeed_InvalidSignature(t *testing.T) {
	addr := ethCrypto.PubkeyToAddress(mustGenerateKey(t).PublicKey)
	// 65 zero bytes are not a valid signature
	err := VerifySeed(make([]byte, 65), addr)
	require.Error(t, err)
}

// --- GenerateSubtype tests ---

func TestGenerateSubtype_Deterministic(t *testing.T) {
	seed := mustSeed(t)
	s1 := GenerateSubtype(seed, 1)
	s2 := GenerateSubtype(seed, 1)
	require.Equal(t, s1, s2)
}

func TestGenerateSubtype_DifferentIndex(t *testing.T) {
	seed := mustSeed(t)
	require.NotEqual(t, GenerateSubtype(seed, 1), GenerateSubtype(seed, 2))
}

func TestGenerateSubtype_DifferentSeed(t *testing.T) {
	seed1, seed2 := mustSeed(t), mustSeed(t)
	require.NotEqual(t, GenerateSubtype(seed1, 1), GenerateSubtype(seed2, 1))
}

func TestGenerateSubtype_Format(t *testing.T) {
	s := GenerateSubtype(mustSeed(t), 1)
	require.True(t, len(s) == 66 && s[:2] == "0x", "should be '0x' + 64 hex chars, got: %s", s)
	for _, c := range s[2:] {
		require.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'), "should be lowercase hex")
	}
}

// --- AllSubtypes tests ---

func TestAllSubtypes_Length(t *testing.T) {
	require.Len(t, AllSubtypes(mustSeed(t), 50), 50)
}

func TestAllSubtypes_Unique(t *testing.T) {
	subtypes := AllSubtypes(mustSeed(t), 50)
	seen := make(map[string]struct{}, 50)
	for _, st := range subtypes {
		_, dup := seen[st]
		require.False(t, dup, "duplicate subtype found: %s", st)
		seen[st] = struct{}{}
	}
}

func TestAllSubtypes_Indices(t *testing.T) {
	seed := mustSeed(t)
	subtypes := AllSubtypes(seed, 3)
	require.Equal(t, GenerateSubtype(seed, 1), subtypes[0])
	require.Equal(t, GenerateSubtype(seed, 2), subtypes[1])
	require.Equal(t, GenerateSubtype(seed, 3), subtypes[2])
}

// --- GenerateRandomSubtype tests ---

func TestGenerateRandomSubtype_InSet(t *testing.T) {
	seed := mustSeed(t)
	all := AllSubtypes(seed, DefaultSubtypeN)
	allSet := make(map[string]struct{}, len(all))
	for _, st := range all {
		allSet[st] = struct{}{}
	}
	for i := 0; i < 20; i++ {
		st, err := GenerateRandomSubtype(seed, DefaultSubtypeN)
		require.NoError(t, err)
		_, ok := allSet[st]
		require.True(t, ok, "generated subtype %s not in expected set", st)
	}
}

// --- Integration test: AssociateKey with seed + subtype generation ---

func TestAssociateKey_AndSubtypeGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	exec := newTestExecutor(t, NewMockRuntime(testLogger))

	// Generate secp256k1 signing key; derive sender address from it
	signingKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	sender := ethCrypto.PubkeyToAddress(signingKey.PublicKey)

	// Compute seed: sign keccak256(SubtypeKeyMessage) deterministically
	msgHash := ethCrypto.Keccak256([]byte(SubtypeKeyMessage))
	seed, err := ethCrypto.Sign(msgHash, signingKey)
	require.NoError(t, err)
	require.Len(t, seed, 65)

	// Generate P521 encryption key for the user
	p521Key, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	pubKeyBytes := p521Key.PublicKey().Bytes()
	require.Len(t, pubKeyBytes, 133)

	// Encrypt the seed with the user's P521 private key and the enclave's P521 public key
	encryptedSeed, err := crypto.Encrypt(p521Key, exec.keySet.CommunicationKey.PublicKey(), seed)
	require.NoError(t, err)

	// Build AssociateKey payload: P521 pubkey (133 bytes) + encrypted seed (93 bytes)
	payloadWithSeed := append(pubKeyBytes, encryptedSeed...)
	require.Len(t, payloadWithSeed, 226)

	appState := buildEncryptedAppState(t, exec, nil, nil)

	req := newProcessRequest()
	req.RequestType = common.AssociateKey
	req.Payload = payloadWithSeed
	req.Sender = sender

	// Execute AssociateKey request
	respPayload, newAppState, _, err := exec.HandleProcessRequest(context.Background(), req, appState, nil)
	require.NoError(t, err)
	require.NotNil(t, respPayload)
	require.Equal(t, uint8(0), respPayload.ErrorCode, "expected success, got error: %s", respPayload.ErrorMsg)

	// Deserialize new app data and verify seed was stored
	decryptedState, err := crypto.DecryptWithAES(exec.keySet.StateKey, newAppState.EncryptedState)
	require.NoError(t, err)
	ad, err := appdata.DeserializeAppData(decryptedState)
	require.NoError(t, err)

	storedSeed, hasSeed := ad.GetSeed(sender)
	require.True(t, hasSeed, "seed should be stored after AssociateKey")
	require.Equal(t, seed, storedSeed)

	// Now verify that encryptEvents generates a subtype from the seed
	// Register the P521 key in a fresh app state that has both key and seed
	keyStore := appdata.KeyStore{sender: p521Key.PublicKey()}
	seedStore := appdata.SeedStore{sender: seed}

	commKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	plainEvents := []common.PlainEvent{{
		UserID:       sender,
		EventSubType: "deposit", // WASM-provided subtype; should be overridden
		Data:         []byte(`{"type":"deposit"}`),
	}}

	e := &StatelessExecutor{log: testLogger}
	encryptedEvents, failure := e.encryptEvents(
		context.Background(),
		plainEvents,
		common.NewApplicationId(1),
		commKey,
		nil,
		keyStore,
		seedStore,
	)
	require.Nil(t, failure)
	require.Len(t, encryptedEvents, 1)

	// The subtype should be one of the N possible subtypes derived from the seed
	generatedSubtype := encryptedEvents[0].EventSubType
	allSubtypes := AllSubtypes(seed, DefaultSubtypeN)
	require.Contains(t, allSubtypes, generatedSubtype, "generated subtype should be in the expected set")
	require.NotEqual(t, "deposit", generatedSubtype, "WASM-provided subtype should be overridden")
}

// --- Event retrieval test: verify subtype set filtering ---

func TestEventRetrieval_BySubtypeSet(t *testing.T) {
	// User A has a seed; user B does not
	seedA := mustSeed(t)
	userAAddr := ethCrypto.PubkeyToAddress(mustGenerateKey(t).PublicKey)
	userBAddr := ethCrypto.PubkeyToAddress(mustGenerateKey(t).PublicKey)

	// Generate P521 keys for both users
	userAKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userBKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	keyStore := appdata.KeyStore{
		userAAddr: userAKey.PublicKey(),
		userBAddr: userBKey.PublicKey(),
	}
	// Only user A has a seed
	seedStore := appdata.SeedStore{userAAddr: seedA}

	commKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	plainEvents := []common.PlainEvent{
		{UserID: userAAddr, EventSubType: "deposit", Data: []byte("a1")},
		{UserID: userAAddr, EventSubType: "withdrawal", Data: []byte("a2")},
		{UserID: userBAddr, EventSubType: "deposit", Data: []byte("b1")},
	}

	e := &StatelessExecutor{log: testLogger}
	encrypted, failure := e.encryptEvents(
		context.Background(),
		plainEvents,
		common.NewApplicationId(1),
		commKey,
		nil,
		keyStore,
		seedStore,
	)
	require.Nil(t, failure)
	require.Len(t, encrypted, 3)

	// All user A events should have a subtype in AllSubtypes(seedA, N)
	userASubtypes := AllSubtypes(seedA, DefaultSubtypeN)
	subtypeSet := make(map[string]struct{}, len(userASubtypes))
	for _, st := range userASubtypes {
		subtypeSet[st] = struct{}{}
	}

	for _, ev := range encrypted {
		if ev.UserID == userAAddr {
			_, ok := subtypeSet[ev.EventSubType]
			require.True(t, ok, "user A event subtype %s should be in expected set", ev.EventSubType)
			require.NotEqual(t, "deposit", ev.EventSubType)
			require.NotEqual(t, "withdrawal", ev.EventSubType)
		} else if ev.UserID == userBAddr {
			// User B has no seed; original subtype is preserved
			require.Equal(t, "deposit", ev.EventSubType)
		}
	}
}

// --- helpers ---

func mustGenerateKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	return k
}

func mustSeed(t *testing.T) []byte {
	t.Helper()
	k, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	msgHash := ethCrypto.Keccak256([]byte(SubtypeKeyMessage))
	seed, err := ethCrypto.Sign(msgHash, k)
	require.NoError(t, err)
	return seed
}

// Verify that AssociateKeyRequest with the wrong seed signer is rejected
func TestAssociateKey_WrongSeedSigner(t *testing.T) {
	exec := newTestExecutor(t, NewMockRuntime(testLogger))

	// Sender is one address but seed is signed by a different key
	realSender := ethCrypto.PubkeyToAddress(mustGenerateKey(t).PublicKey)
	differentKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	msgHash := ethCrypto.Keccak256([]byte(SubtypeKeyMessage))
	wrongSeed, err := ethCrypto.Sign(msgHash, differentKey)
	require.NoError(t, err)

	p521Key, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	// Encrypt the wrong seed with the user's P521 key and the enclave's P521 public key
	encryptedWrongSeed, err := crypto.Encrypt(p521Key, exec.keySet.CommunicationKey.PublicKey(), wrongSeed)
	require.NoError(t, err)

	payloadWithSeed := append(p521Key.PublicKey().Bytes(), encryptedWrongSeed...)

	appState := buildEncryptedAppState(t, exec, nil, nil)
	req := newProcessRequest()
	req.RequestType = common.AssociateKey
	req.Payload = payloadWithSeed
	req.Sender = realSender

	respPayload, _, _, err := exec.HandleProcessRequest(context.Background(), req, appState, nil)
	require.NoError(t, err) // no system error; returns an error payload
	require.NotNil(t, respPayload)
	require.NotEqual(t, uint8(0), respPayload.ErrorCode, "should reject mismatched seed signer")
}

// Verify that without a seed the original WASM-provided subtype is preserved
func TestEncryptEvents_NoSeed_PreservesSubtype(t *testing.T) {
	p521Key, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userAddr := ethCrypto.PubkeyToAddress(mustGenerateKey(t).PublicKey)

	keyStore := appdata.KeyStore{userAddr: p521Key.PublicKey()}
	seedStore := appdata.SeedStore{} // no seed

	commKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	plain := []common.PlainEvent{{UserID: userAddr, EventSubType: "original_type", Data: []byte("x")}}

	e := &StatelessExecutor{log: testLogger}
	encrypted, failure := e.encryptEvents(context.Background(), plain, common.NewApplicationId(1), commKey, nil, keyStore, seedStore)
	require.Nil(t, failure)
	require.Equal(t, "original_type", encrypted[0].EventSubType, "subtype should be preserved when no seed is registered")

	_ = commontestutil.GenerateRandomRequestID() // keep import used
}
