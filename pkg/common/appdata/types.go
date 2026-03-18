package appdata

import (
	"bytes"
	"encoding/binary"
	"fmt"

	ethCommon "github.com/ethereum/go-ethereum/common"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
)

const (
	//version of the AppData format - increment in case of future changes
	Version_1 = 1
	Version_2 = 2
	// KeyStore_KeySize is the size in bytes of a key in the appKeys map (20 bytes)
	KeyStore_KeySize = ethCommon.AddressLength
	// KeyStore_ValSize is the size in bytes of a value in the appKeys map (a PublicKeyP521)
	KeyStore_ValSize = 133
	// SeedStore_ValSize is the size in bytes of a seed (secp256k1 signature: 65 bytes)
	SeedStore_ValSize = 65
)

type AppData struct {
	version  uint8
	appNonce uint64
	appKeys  KeyStore
	appSeeds SeedStore
	appState []byte
}
type KeyStore map[ethCommon.Address]*cryptotypes.PublicKeyP521
type SeedStore map[ethCommon.Address][]byte

func NewAppData(initialAppState []byte) *AppData {
	return &AppData{
		version:  Version_2,
		appNonce: 0,
		appKeys:  make(map[ethCommon.Address]*cryptotypes.PublicKeyP521),
		appSeeds: make(SeedStore),
		appState: initialAppState,
	}
}

// Serialize converts the AppData struct into a byte slice.
// Version 2 format:
// - Version (uint8)
// - Nonce (uint64)
// - Number of keys (uint32)
// - Key-value pairs (KeyStore_KeySize + KeyStore_ValSize each)
// - Number of seeds (uint32)
// - Seed pairs (KeyStore_KeySize + SeedStore_ValSize each)
// - appState (variable length)
func (s *AppData) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Write the version
	if err := binary.Write(&buf, binary.BigEndian, s.version); err != nil {
		return nil, fmt.Errorf("failed to write version: %w", err)
	}

	// Write the nonce
	if err := binary.Write(&buf, binary.BigEndian, s.appNonce); err != nil {
		return nil, fmt.Errorf("failed to write nonce: %w", err)
	}

	// Write the number of keys as a uint32
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(s.appKeys))); err != nil {
		return nil, fmt.Errorf("failed to write key count: %w", err)
	}

	// Write each key-value pair
	// Note: the same order is not guaranted between differrent calls of this method!
	for k, v := range s.appKeys {
		keyBytes := k.Bytes()
		if len(keyBytes) != KeyStore_KeySize {
			return nil, fmt.Errorf("key '%s' has incorrect size: expected %d, got %d", k, KeyStore_KeySize, len(keyBytes))
		}
		buf.Write(keyBytes)
		buf.Write(v.Bytes())
	}

	// Write the number of seeds as a uint32
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(s.appSeeds))); err != nil {
		return nil, fmt.Errorf("failed to write seed count: %w", err)
	}

	// Write each seed key-value pair (address + 65-byte seed)
	for addr, seed := range s.appSeeds {
		buf.Write(addr.Bytes())
		if len(seed) != SeedStore_ValSize {
			return nil, fmt.Errorf("seed for %s has incorrect size: expected %d, got %d", addr, SeedStore_ValSize, len(seed))
		}
		buf.Write(seed)
	}

	// Write the appState
	buf.Write(s.appState)

	return buf.Bytes(), nil
}

func (s *AppData) GetAppState() []byte {
	return s.appState
}

func (s *AppData) SetAppState(appState []byte) {
	s.appState = appState
}

func (s *AppData) AddKey(user ethCommon.Address, key cryptotypes.PublicKeyP521) {
	s.appKeys[user] = &key
}

func (s *AppData) GetKeyStore() KeyStore {
	return s.appKeys
}

func (s *AppData) AddSeed(user ethCommon.Address, seed []byte) {
	seedCopy := make([]byte, len(seed))
	copy(seedCopy, seed)
	s.appSeeds[user] = seedCopy
}

func (s *AppData) GetSeed(user ethCommon.Address) ([]byte, bool) {
	seed, exists := s.appSeeds[user]
	return seed, exists
}

func (s *AppData) GetSeedStore() SeedStore {
	return s.appSeeds
}

func (s *AppData) IncrementNonce() {
	s.appNonce = s.appNonce + 1
}

// DeserializeAppData converts a byte slice back into an AppData struct.
func DeserializeAppData(data []byte) (*AppData, error) {
	reader := bytes.NewReader(data)

	// Read the version
	var version uint8
	if err := binary.Read(reader, binary.BigEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read version: %w", err)
	}

	// Read the nonce
	var nonce uint64
	if err := binary.Read(reader, binary.BigEndian, &nonce); err != nil {
		return nil, fmt.Errorf("failed to read nonce: %w", err)
	}

	// Read the number of keys
	var numKeys uint32
	if err := binary.Read(reader, binary.BigEndian, &numKeys); err != nil {
		return nil, fmt.Errorf("failed to read key count: %w", err)
	}

	// Check for sufficient data for all key-value pairs
	expectedMapSize := int(numKeys) * (KeyStore_KeySize + KeyStore_ValSize)
	if reader.Len() < expectedMapSize {
		return nil, fmt.Errorf("insufficient data for key-value pairs: need %d, have %d", expectedMapSize, reader.Len())
	}

	// Read each key-value pair
	ks := make(KeyStore)
	for i := uint32(0); i < numKeys; i++ {
		keyBuf := make([]byte, KeyStore_KeySize)
		if _, err := reader.Read(keyBuf); err != nil { // Should not fail due to length check above
			return nil, fmt.Errorf("failed to read key %d: %w", i, err)
		}

		valBuf := make([]byte, KeyStore_ValSize)
		if _, err := reader.Read(valBuf); err != nil { // Should not fail
			return nil, fmt.Errorf("failed to read value for key %s: %w", string(keyBuf), err)
		}

		val, err := cryptotypes.NewPublicKeyP521(valBuf)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key for key %s: %w", string(keyBuf), err)
		}
		ks[ethCommon.BytesToAddress(keyBuf)] = val
	}

	// Read seed store — only present in Version_2 and above
	ss := make(SeedStore)
	switch version {
	case Version_1:
		// no seed store in v1; ss stays empty
	case Version_2:
		var numSeeds uint32
		if err := binary.Read(reader, binary.BigEndian, &numSeeds); err != nil {
			return nil, fmt.Errorf("failed to read seed count: %w", err)
		}
		expectedSeedSize := int(numSeeds) * (KeyStore_KeySize + SeedStore_ValSize)
		if reader.Len() < expectedSeedSize {
			return nil, fmt.Errorf("insufficient data for seed pairs: need %d, have %d", expectedSeedSize, reader.Len())
		}
		for i := uint32(0); i < numSeeds; i++ {
			addrBuf := make([]byte, KeyStore_KeySize)
			if _, err := reader.Read(addrBuf); err != nil {
				return nil, fmt.Errorf("failed to read seed address %d: %w", i, err)
			}
			seedBuf := make([]byte, SeedStore_ValSize)
			if _, err := reader.Read(seedBuf); err != nil {
				return nil, fmt.Errorf("failed to read seed value %d: %w", i, err)
			}
			ss[ethCommon.BytesToAddress(addrBuf)] = seedBuf
		}
	default:
		return nil, fmt.Errorf("unsupported AppData version: %d", version)
	}

	// The rest of the reader is the appState
	appState := make([]byte, reader.Len())
	if _, err := reader.Read(appState); err != nil {
		return nil, fmt.Errorf("failed to read appState: %w", err)
	}

	return &AppData{
		version:  version,
		appNonce: nonce,
		appKeys:  ks,
		appSeeds: ss,
		appState: appState,
	}, nil
}
