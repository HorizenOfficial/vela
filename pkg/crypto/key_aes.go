package crypto

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/horizen-pes/pkg/common"
)

// GenerateAESKey generates a random AES-256 key.
func GenerateAESKey() (common.AES256Key, error) {
	var key common.AES256Key
	_, err := rand.Read(key[:])
	if err != nil {
		return common.AES256Key{}, fmt.Errorf("failed to generate random AES key: %w", err)
	}
	return key, nil
}

// SaveAESKeyToFile saves an AES key to a file in raw binary format.
func SaveAESKeyToFile(key *common.AES256Key, filename string) error {
	err := os.WriteFile(filename, key[:], 0600)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

// LoadAESKeyFromFile loads an AES key from a file in raw binary format.
func LoadAESKeyFromFile(filename string) (*common.AES256Key, error) {
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("invalid key size: %d", len(keyBytes))
	}
	var key common.AES256Key
	copy(key[:], keyBytes)
	return &key, nil
}

// SaveAESKeyToFilePEM saves an AES key to a file in PEM format.
func SaveAESKeyToFilePEM(key *common.AES256Key, filename string) error {
	pemBlock := &pem.Block{
		Type:  "AES-256 KEY",
		Bytes: key[:],
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	err = pem.Encode(file, pemBlock)
	if err != nil {
		return fmt.Errorf("failed to encode PEM: %w", err)
	}

	return nil
}

// LoadAESKeyFromFilePEM loads an AES key from a file in PEM format.
func LoadAESKeyFromFilePEM(filename string) (*common.AES256Key, error) {
	pemBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	if pemBlock.Type != "AES-256 KEY" {
		return nil, fmt.Errorf("unsupported PEM type: %s", pemBlock.Type)
	}

	if len(pemBlock.Bytes) != 32 {
		return nil, fmt.Errorf("invalid key size: %d", len(pemBlock.Bytes))
	}

	var key common.AES256Key
	copy(key[:], pemBlock.Bytes)
	return &key, nil
}
