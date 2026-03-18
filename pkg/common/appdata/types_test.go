package appdata

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/HorizenOfficial/vela/pkg/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppDataSerializationDeserialization(t *testing.T) {
	// 1. Create a new AppData with some test data
	appStateSize := 128
	appState := make([]byte, appStateSize)
	_, err := rand.Read(appState)
	assert.NoError(t, err)

	appData := NewAppData(appState)

	// 2. Add a couple of keys
	addr1 := ethCommon.Address{1}
	sk1, err := crypto.GeneratePrivateKeyP521()
	assert.NoError(t, err)
	pk1 := sk1.PublicKey()
	appData.AddKey(addr1, *pk1)

	addr2 := ethCommon.Address{2}
	sk2, err := crypto.GeneratePrivateKeyP521()
	assert.NoError(t, err)
	pk2 := sk2.PublicKey()
	appData.AddKey(addr2, *pk2)

	// 3. Add seeds
	seed1 := make([]byte, SeedStore_ValSize)
	_, err = rand.Read(seed1)
	assert.NoError(t, err)
	appData.AddSeed(addr1, seed1)

	seed2 := make([]byte, SeedStore_ValSize)
	_, err = rand.Read(seed2)
	assert.NoError(t, err)
	appData.AddSeed(addr2, seed2)

	// 4. increment nonce
	appData.IncrementNonce()

	// 5. Serialize it
	serializedData, err := appData.Serialize()
	assert.NoError(t, err)

	// 6. Deserialize it
	deserializedAppData, err := DeserializeAppData(serializedData)
	assert.NoError(t, err)

	// 7. Check that they are the same
	assert.Equal(t, uint8(Version_2), deserializedAppData.version, "Version should be Version_2")
	assert.Equal(t, appData.appNonce, deserializedAppData.appNonce, "Nonce should be the same")
	assert.Equal(t, appData.appState, deserializedAppData.appState, "appState should be the same")
	assert.Equal(t, len(appData.appKeys), len(deserializedAppData.appKeys), "Number of keys should be the same")
	for addr, pk := range appData.appKeys {
		deserializedPk, ok := deserializedAppData.appKeys[addr]
		assert.True(t, ok, "Key should be present in deserialized map")
		assert.Equal(t, pk.Bytes(), deserializedPk.Bytes(), "Public keys should be the same")
	}

	// Check seeds round-trip
	assert.Equal(t, len(appData.appSeeds), len(deserializedAppData.appSeeds), "Number of seeds should be the same")
	for addr, seed := range appData.appSeeds {
		deserializedSeed, ok := deserializedAppData.appSeeds[addr]
		assert.True(t, ok, "Seed should be present in deserialized map")
		assert.Equal(t, seed, deserializedSeed, "Seeds should be equal")
	}
}

// buildV1Bytes manually constructs a Version_1 serialized AppData blob (no seed section).
func buildV1Bytes(appState []byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(Version_1)                                    // version
	binary.Write(&buf, binary.BigEndian, uint64(0))             // nonce
	binary.Write(&buf, binary.BigEndian, uint32(0))             // numKeys
	buf.Write(appState)
	return buf.Bytes()
}

func TestAppDataDeserializeVersion1BackwardCompat(t *testing.T) {
	appState := []byte{0xAA, 0xBB, 0xCC}
	v1Bytes := buildV1Bytes(appState)

	deserialized, err := DeserializeAppData(v1Bytes)
	require.NoError(t, err)
	assert.Equal(t, uint8(Version_1), deserialized.version)
	assert.Equal(t, appState, deserialized.appState)
	assert.Empty(t, deserialized.appSeeds, "V1 data should produce empty seed store")
	assert.Empty(t, deserialized.appKeys, "V1 data with no keys should produce empty key store")
}
