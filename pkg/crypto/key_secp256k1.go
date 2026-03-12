package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
)

/*
GeneratePrivateKeySecp256k1 generates a private key with Elliptic Curve secp256k1.
(This is the curve used in Bitcoin and Ethereum)
*/
func GeneratePrivateKeySecp256k1() (*cryptotypes.PrivateKeySecp256k1, error) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return &cryptotypes.PrivateKeySecp256k1{PrivateKey: key}, nil
}

// SavePrivateKeySecp256k1ToFileDER saves a secp256k1 private key to a file in DER format.
func SavePrivateKeySecp256k1ToFileDER(privKey *cryptotypes.PrivateKeySecp256k1, filename string) error {
	derBytes := crypto.FromECDSA(privKey.PrivateKey)
	err := os.WriteFile(filename, derBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// LoadPrivateKeySecp256k1FromFileDER loads a secp256k1 private key from a file in DER format.
func LoadPrivateKeySecp256k1FromFileDER(filename string) (*cryptotypes.PrivateKeySecp256k1, error) {
	derBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	key, err := crypto.ToECDSA(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &cryptotypes.PrivateKeySecp256k1{PrivateKey: key}, nil
}

// SavePrivateKeySecp256k1ToFilePEM saves a secp256k1 private key to a file in PEM format.
func SavePrivateKeySecp256k1ToFilePEM(privKey *cryptotypes.PrivateKeySecp256k1, filename string) error {
	derBytes := crypto.FromECDSA(privKey.PrivateKey)
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derBytes,
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

// LoadPrivateKeySecp256k1FromFilePEM loads a secp256k1 private key from a file in PEM format.
func LoadPrivateKeySecp256k1FromFilePEM(filename string) (*cryptotypes.PrivateKeySecp256k1, error) {
	pemBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	key, err := crypto.ToECDSA(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &cryptotypes.PrivateKeySecp256k1{PrivateKey: key}, nil
}

// ExportPrivateKeySecp256k1ToHex exports a secp256k1 private key to a 64-character hex string.
func ExportPrivateKeySecp256k1ToHex(privKey *cryptotypes.PrivateKeySecp256k1) string {
	return hex.EncodeToString(crypto.FromECDSA(privKey.PrivateKey))
}

// ExportPublicKeySecp256k1ToHex exports a secp256k1 public key to a hex string.
func ExportPublicKeySecp256k1ToHex(pubKey *cryptotypes.PublicKeySecp256k1) string {
	return hex.EncodeToString(crypto.FromECDSAPub(pubKey.PublicKey))
}

// ImportPrivateKeySecp256k1FromHex imports a secp256k1 private key from a 64-character hex string.
func ImportPrivateKeySecp256k1FromHex(hexKey string) (*cryptotypes.PrivateKeySecp256k1, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	key, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &cryptotypes.PrivateKeySecp256k1{PrivateKey: key}, nil
}

// ImportPublicKeySecp256k1FromHex imports a secp256k1 public key from a hex string.
func ImportPublicKeySecp256k1FromHex(hexKey string) (*cryptotypes.PublicKeySecp256k1, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	key, err := crypto.UnmarshalPubkey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public key: %w", err)
	}

	return &cryptotypes.PublicKeySecp256k1{PublicKey: key}, nil
}
