package executor

import (
	"encoding/json"
	"fmt"

	"github.com/horizen-pes/pkg/common/crypto"
	pesCrypto "github.com/horizen-pes/pkg/crypto"
)

// EnclaveKeySet contains all the keys used by the executor.
type EnclaveKeySet struct {
	// Key used for encrypting/decrypting the state exchanged between executor and manager.
	StateKey crypto.AES256Key
	// Key used for signing the messages (update payload) sent from the executor to the tee autenticator through the manager.
	SigningKey crypto.PrivateKeySecp256k1
	// Key used for encrypting/decrypting the payloads, events and reports.
	// Sender and receiver keys must Elliptic Curve Diffie-Hellman over NIST curve
	CommunicationKey crypto.PrivateKeyP521
}

// enclaveKeySetJSON is a helper struct for JSON marshaling/unmarshaling
type enclaveKeySetJSON struct {
	StateKey         crypto.AES256Key `json:"StateKey"`
	SigningKey       []byte           `json:"SigningKey"`
	CommunicationKey []byte           `json:"CommunicationKey"`
}

// MarshalJSON returns the JSON representation of the EnclaveKeySet.
func (ks *EnclaveKeySet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&enclaveKeySetJSON{
		StateKey:         ks.StateKey,
		SigningKey:       ks.SigningKey.Bytes(),
		CommunicationKey: ks.CommunicationKey.Bytes(),
	})
}

// UnmarshalJSON deserializes the JSON representation of the EnclaveKeySet.
func (ks *EnclaveKeySet) UnmarshalJSON(data []byte) error {
	var jks enclaveKeySetJSON
	if err := json.Unmarshal(data, &jks); err != nil {
		return fmt.Errorf("failed to unmarshal enclave key set JSON: %w", err)
	}

	signingKey, err := pesCrypto.NewPrivateKeySecp256k1FromBytes(jks.SigningKey)
	if err != nil {
		return fmt.Errorf("failed to deserialize signing key: %w", err)
	}
	communicationKey, err := pesCrypto.NewPrivateKeyP521FromBytes(jks.CommunicationKey)
	if err != nil {
		return fmt.Errorf("failed to deserialize communication key: %w", err)
	}

	ks.StateKey = jks.StateKey
	ks.SigningKey = *signingKey
	ks.CommunicationKey = *communicationKey

	return nil
}

// Serialize returns the JSON representation of the EnclaveKeySet.
func (ks *EnclaveKeySet) Serialize() ([]byte, error) {
	serialized, err := json.Marshal(ks)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize enclave key set: %w", err)
	}
	return serialized, nil
}

// DeserializeEnclaveKeySet deserializes the JSON representation of the EnclaveKeySet.
func DeserializeEnclaveKeySet(data []byte) (*EnclaveKeySet, error) {
	var ks EnclaveKeySet
	err := json.Unmarshal(data, &ks)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize enclave key set: %w", err)
	}
	return &ks, nil
}
