package executor

import (
	"encoding/json"
	"fmt"

	"github.com/horizen-pes/pkg/common/crypto"
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

// EnclaveKeySetRecovery contains the data needed to recover the EnclaveKeySet.
type EnclaveKeySetRecovery struct {
	// RecoveryType is the type of recovery data.
	RecoveryType int
	// KeySetCiphertext is the encrypted EnclaveKeySet.
	KeySetCiphertext []byte
	// RecoveryCiphertext is the master key used to encrypt the EnclaveKeySet.
	RecoveryCiphertext []byte
}
