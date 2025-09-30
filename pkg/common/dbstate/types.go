package dbstate

import (
	"bytes"
	"encoding/binary"
	"fmt"

	ethCommon "github.com/ethereum/go-ethereum/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
)

const (
	// keySize is the size in bytes of a key in the keyStore map (20 bytes)
	keySize = ethCommon.AddressLength
	// valSize is the size in bytes of a PublicKeyP521
	valSize = 133
)

type DBState struct {
	keyStore KeyStore
	appState []byte
}
type KeyStore map[ethCommon.Address]*cryptotypes.PublicKeyP521

func NewDBState(initialAppState []byte) *DBState {
	return &DBState{
		keyStore: make(map[ethCommon.Address]*cryptotypes.PublicKeyP521),
		appState: initialAppState,
	}
}

// Serialize converts the DBState struct into a byte slice.
// The format is:
// - Number of keys (uint32)
// - Key-value pairs (keySize + valSize each)
// - appState (variable length)
func (s *DBState) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Write the number of keys as a uint32
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(s.keyStore))); err != nil {
		return nil, fmt.Errorf("failed to write key count: %w", err)
	}

	// Write each key-value pair
	for k, v := range s.keyStore {
		keyBytes := k.Bytes()
		if len(keyBytes) != keySize {
			return nil, fmt.Errorf("key '%s' has incorrect size: expected %d, got %d", k, keySize, len(keyBytes))
		}
		buf.Write(keyBytes)
		buf.Write(v.Bytes())
	}

	// Write the appState
	buf.Write(s.appState)

	return buf.Bytes(), nil
}

// DBState struct and its methods...

func (s *DBState) GetAppState() []byte {
	return s.appState
}

func (s *DBState) SetAppState(appState []byte) {
	s.appState = appState
}

func (s *DBState) AddKey(user ethCommon.Address, key cryptotypes.PublicKeyP521) {
	s.keyStore[user] = &key
}

func (s *DBState) GetKeyStore() KeyStore {
	return s.keyStore
}

// DeserializeDBState converts a byte slice back into a DBState struct.
func DeserializeDBState(data []byte) (*DBState, error) {
	reader := bytes.NewReader(data)

	// Read the number of keys
	var numKeys uint32
	if err := binary.Read(reader, binary.BigEndian, &numKeys); err != nil {
		return nil, fmt.Errorf("failed to read key count: %w", err)
	}

	// Check for sufficient data for all key-value pairs
	expectedMapSize := int(numKeys) * (keySize + valSize)
	if reader.Len() < expectedMapSize {
		return nil, fmt.Errorf("insufficient data for key-value pairs: need %d, have %d", expectedMapSize, reader.Len())
	}

	// Read each key-value pair
	ks := make(KeyStore)
	for i := uint32(0); i < numKeys; i++ {
		keyBuf := make([]byte, keySize)
		if _, err := reader.Read(keyBuf); err != nil { // Should not fail due to length check above
			return nil, fmt.Errorf("failed to read key %d: %w", i, err)
		}

		valBuf := make([]byte, valSize)
		if _, err := reader.Read(valBuf); err != nil { // Should not fail
			return nil, fmt.Errorf("failed to read value for key %s: %w", string(keyBuf), err)
		}

		val, err := cryptotypes.NewPublicKeyP521(valBuf)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key for key %s: %w", string(keyBuf), err)
		}
		ks[ethCommon.BytesToAddress(keyBuf)] = val
	}

	// The rest of the reader is the appState
	appState := make([]byte, reader.Len())
	if _, err := reader.Read(appState); err != nil {
		return nil, fmt.Errorf("failed to read appState: %w", err)
	}

	return &DBState{
		keyStore: ks,
		appState: appState,
	}, nil
}
