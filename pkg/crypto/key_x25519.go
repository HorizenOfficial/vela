package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/horizen-pes/pkg/common"
)

/*
GeneratePrivateKey25519 generates a private key with Elliptic Curve 25519.
*/
func GeneratePrivateKey25519() (*common.PrivateKey25519, error) {
	curve := ecdh.X25519()
	newKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return &common.PrivateKey25519{newKey}, nil
}

// SavePrivateKey25519ToFileDER saves a 25519 private key to a file in PKCS #8, ASN.1 DER format.
func SavePrivateKey25519ToFileDER(privKey *common.PrivateKey25519, filename string) error {
	derBytes, err := x509.MarshalPKCS8PrivateKey(privKey.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	err = os.WriteFile(filename, derBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// LoadPrivateKey25519FromFileDER loads a 25519 private key from a file in PKCS #8, ASN.1 DER format.
func LoadPrivateKey25519FromFileDER(filename string) (*common.PrivateKey25519, error) {
	derBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	key, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	ecdhKey, ok := key.(*ecdh.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an ecdh.PrivateKey")
	}

	return &common.PrivateKey25519{ecdhKey}, nil
}

// SavePrivateKey25519ToFilePEM saves a 25519 private key to a file in PEM format.
func SavePrivateKey25519ToFilePEM(privKey *common.PrivateKey25519, filename string) error {
	derBytes, err := x509.MarshalPKCS8PrivateKey(privKey.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
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

// LoadPrivateKey25519FromFilePEM loads a 25519 private key from a file in PEM format.
func LoadPrivateKey25519FromFilePEM(filename string) (*common.PrivateKey25519, error) {
	pemBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	ecdhKey, ok := key.(*ecdh.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an ecdh.PrivateKey")
	}

	return &common.PrivateKey25519{ecdhKey}, nil
}

// ExportPrivateKey25519ToHex exports a 25519 private key to a hex string.
func ExportPrivateKey25519ToHex(privKey *common.PrivateKey25519) string {
	return hex.EncodeToString(privKey.Bytes())
}

// ExportPublicKey25519ToHex exports a 25519 public key to a hex string.
func ExportPublicKey25519ToHex(pubKey *common.PublicKey25519) string {
	return hex.EncodeToString(pubKey.Bytes())
}

// ImportPrivateKey25519FromHex imports a 25519 private key from a hex string.
func ImportPrivateKey25519FromHex(hexKey string) (*common.PrivateKey25519, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	curve := ecdh.X25519()
	key, err := curve.NewPrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	return &common.PrivateKey25519{key}, nil
}
