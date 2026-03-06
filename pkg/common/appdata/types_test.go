package appdata

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/crypto"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestAppDataSerializationDeserialization(t *testing.T) {
	// 1. Create a new AppData with some test data
	appStateSize := 128
	appState := make([]byte, appStateSize)
	_, err := rand.Read(appState)
	assert.NoError(t, err)

	appData := NewAppData(appState)
	fingerprint := sha256.Sum256([]byte("mock-wasm-module"))
	appData.SetWasmFingerprint(fingerprint)

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

	// 3. increment nonce
	appData.IncrementNonce()

	// 3. Serialize it
	serializedData, err := appData.Serialize()
	assert.NoError(t, err)

	// 4. Deserialize it
	deserializedAppData, err := DeserializeAppData(serializedData)
	assert.NoError(t, err)

	// 5. Check that they are the same
	assert.Equal(t, appData.version, deserializedAppData.version, "Version should be the same")
	assert.Equal(t, appData.appNonce, deserializedAppData.appNonce, "Nonce should be the same")
	assert.Equal(t, appData.GetWasmFingerprint(), deserializedAppData.GetWasmFingerprint(), "wasm fingerprint should be the same")
	assert.Equal(t, appData.appState, deserializedAppData.appState, "appState should be the same")
	assert.Equal(t, len(appData.appKeys), len(deserializedAppData.appKeys), "Number of keys should be the same")
	for addr, pk := range appData.appKeys {
		deserializedPk, ok := deserializedAppData.appKeys[addr]
		assert.True(t, ok, "Key should be present in deserialized map")
		assert.Equal(t, pk.Bytes(), deserializedPk.Bytes(), "Public keys should be the same")
	}

	// 6. Verify binary layout includes fingerprint after nonce
	const nonceOffset = 1
	const fingerprintOffset = nonceOffset + 8
	assert.GreaterOrEqual(t, len(serializedData), fingerprintOffset+WasmFingerprintSize)
	assert.Equal(t, fingerprint[:], serializedData[fingerprintOffset:fingerprintOffset+WasmFingerprintSize])
}
