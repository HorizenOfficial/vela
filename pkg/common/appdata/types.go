package appdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	//version of the AppData format - increment in case of future changes
	Version_1 = 1
	// WasmFingerprintSize is the byte length of SHA-256 fingerprint.
	WasmFingerprintSize = 32
	// Store_KeySize is the size in bytes of an address key in the appKeys/appEventSeeds maps (20 bytes)
	Store_KeySize = ethCommon.AddressLength
	// KeyStore_ValSize is the size in bytes of a value in the appKeys map (a PublicKeyP521)
	KeyStore_ValSize = 133
	// SeedStore_ValSize is the size in bytes of a seed (secp256k1 signature: 65 bytes)
	SeedStore_ValSize = 65
)

type AppData struct {
	version         uint8
	appNonce        uint64
	wasmFingerprint [WasmFingerprintSize]byte
	appKeys         KeyStore
	appEventSeeds SeedStore
	appState        []byte
}
type KeyStore map[ethCommon.Address]*cryptotypes.PublicKeyP521
type SeedStore map[ethCommon.Address][]byte

func NewAppData(initialAppState []byte) *AppData {
	return &AppData{
		version:         Version_1,
		appNonce:        0,
		wasmFingerprint: [WasmFingerprintSize]byte{},
		appKeys:         make(KeyStore),
		appEventSeeds: make(SeedStore),
		appState:        initialAppState,
	}
}

// Serialize converts the AppData struct into a byte slice.
// The format is:
// - Version (uint8)
// - Nonce (uint64)
// - WASM fingerprint ([32]byte)
// - Number of keys (uint32)
// - Key-value pairs (Store_KeySize + KeyStore_ValSize each)
// - Number of seeds (uint32)
// - Seed pairs (Store_KeySize + SeedStore_ValSize each)
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

	// Write the wasm fingerprint
	if _, err := buf.Write(s.wasmFingerprint[:]); err != nil {
		return nil, fmt.Errorf("failed to write wasm fingerprint: %w", err)
	}

	// Write the number of keys as a uint32
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(s.appKeys))); err != nil {
		return nil, fmt.Errorf("failed to write key count: %w", err)
	}

	// Write each key-value pair
	// Note: the same order is not guaranted between differrent calls of this method!
	for k, v := range s.appKeys {
		keyBytes := k.Bytes()
		if len(keyBytes) != Store_KeySize {
			return nil, fmt.Errorf("key '%s' has incorrect size: expected %d, got %d", k, Store_KeySize, len(keyBytes))
		}
		valBytes := v.Bytes()
		if len(valBytes) != KeyStore_ValSize {
			return nil, fmt.Errorf("value for key '%s' has incorrect size: expected %d, got %d", k, KeyStore_ValSize, len(valBytes))
		}
		buf.Write(keyBytes)
		buf.Write(valBytes)
	}

	// Write the number of seeds as a uint32
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(s.appEventSeeds))); err != nil {
		return nil, fmt.Errorf("failed to write seed count: %w", err)
	}

	// Write each seed key-value pair (address + 65-byte seed)
	for addr, seed := range s.appEventSeeds {
		addrBytes := addr.Bytes()
		if len(addrBytes) != Store_KeySize {
			return nil, fmt.Errorf("seed address '%s' has incorrect size: expected %d, got %d", addr, Store_KeySize, len(addrBytes))
		}
		buf.Write(addrBytes)
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

func (s *AppData) AddSeed(user ethCommon.Address, seed []byte) error {
	if len(seed) != SeedStore_ValSize {
		return fmt.Errorf("invalid seed size for %s: expected %d, got %d", user, SeedStore_ValSize, len(seed))
	}
	seedCopy := make([]byte, len(seed))
	copy(seedCopy, seed)
	s.appEventSeeds[user] = seedCopy
	return nil
}

func (s *AppData) GetSeed(user ethCommon.Address) ([]byte, bool) {
	seed, exists := s.appEventSeeds[user]
	return seed, exists
}

func (s *AppData) GetEventSeedStore() SeedStore {
	return s.appEventSeeds
}

func (s *AppData) IncrementNonce() {
	s.appNonce = s.appNonce + 1
}

func (s *AppData) GetWasmFingerprint() [WasmFingerprintSize]byte {
	return s.wasmFingerprint
}

func (s *AppData) SetWasmFingerprint(fingerprint [WasmFingerprintSize]byte) {
	s.wasmFingerprint = fingerprint
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

	// Read the wasm fingerprint
	var wasmFingerprint [WasmFingerprintSize]byte
	if _, err := io.ReadFull(reader, wasmFingerprint[:]); err != nil {
		return nil, fmt.Errorf("failed to read wasm fingerprint: %w", err)
	}

	if version != Version_1 {
		return nil, fmt.Errorf("unsupported AppData version: %d", version)
	}

	// Read the number of keys
	var numKeys uint32
	if err := binary.Read(reader, binary.BigEndian, &numKeys); err != nil {
		return nil, fmt.Errorf("failed to read key count: %w", err)
	}

	// Check for sufficient data for all key-value pairs
	expectedMapSize := int(numKeys) * (Store_KeySize + KeyStore_ValSize)
	if reader.Len() < expectedMapSize {
		return nil, fmt.Errorf("insufficient data for key-value pairs: need %d, have %d", expectedMapSize, reader.Len())
	}

	// Read each key-value pair
	ks := make(KeyStore)
	for i := uint32(0); i < numKeys; i++ {
		keyBuf := make([]byte, Store_KeySize)
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

	// Read seed store
	ss := make(SeedStore)
	var numSeeds uint32
	if err := binary.Read(reader, binary.BigEndian, &numSeeds); err != nil {
		return nil, fmt.Errorf("failed to read seed count: %w", err)
	}
	expectedSeedSize := int(numSeeds) * (Store_KeySize + SeedStore_ValSize)
	if reader.Len() < expectedSeedSize {
		return nil, fmt.Errorf("insufficient data for seed pairs: need %d, have %d", expectedSeedSize, reader.Len())
	}
	for i := uint32(0); i < numSeeds; i++ {
		addrBuf := make([]byte, Store_KeySize)
		if _, err := reader.Read(addrBuf); err != nil {
			return nil, fmt.Errorf("failed to read seed address %d: %w", i, err)
		}
		seedBuf := make([]byte, SeedStore_ValSize)
		if _, err := reader.Read(seedBuf); err != nil {
			return nil, fmt.Errorf("failed to read seed value %d: %w", i, err)
		}
		ss[ethCommon.BytesToAddress(addrBuf)] = seedBuf
	}

	// The rest of the reader is the appState
	appState := make([]byte, reader.Len())
	if _, err := reader.Read(appState); err != nil {
		return nil, fmt.Errorf("failed to read appState: %w", err)
	}

	return &AppData{
		version:         version,
		appNonce:        nonce,
		wasmFingerprint: wasmFingerprint,
		appKeys:         ks,
		appEventSeeds: ss,
		appState:        appState,
	}, nil
}
