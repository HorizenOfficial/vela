package crypto

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"

	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
)

/*
GeneratePrivateKeyP521 generates a private key with Elliptic Curve Diffie-Hellman over NIST P-521 curve,
also known as secp521r1.
(This can be used for encryption/decription)
*/
func GeneratePrivateKeyP521() (*cryptotypes.PrivateKeyP521, error) {
	curve := ecdh.P521()
	newKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return &cryptotypes.PrivateKeyP521{PrivateKey: newKey}, nil
}

// SavePrivateKeyP521ToFileDER saves a P521 private key to a file in PKCS #8, ASN.1 DER format.
func SavePrivateKeyP521ToFileDER(privKey *cryptotypes.PrivateKeyP521, filename string) error {
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

// LoadPrivateKeyP521FromFileDER loads a P521 private key from a file in PKCS #8, ASN.1 DER format.
func LoadPrivateKeyP521FromFileDER(filename string) (*cryptotypes.PrivateKeyP521, error) {
	derBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	key, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an ecdsa.PrivateKey")
	}

	ecdhKey, err := ecdsaKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to ECDH key: %w", err)
	}

	return &cryptotypes.PrivateKeyP521{PrivateKey: ecdhKey}, nil
}

// SavePrivateKeyP521ToFilePEM saves a P521 private key to a file in PEM format.
func SavePrivateKeyP521ToFilePEM(privKey *cryptotypes.PrivateKeyP521, filename string) error {
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

// LoadPrivateKeyP521FromFilePEM loads a P521 private key from a file in PEM format.
func LoadPrivateKeyP521FromFilePEM(filename string) (*cryptotypes.PrivateKeyP521, error) {
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

	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an ecdsa.PrivateKey")
	}

	ecdhKey, err := ecdsaKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to ECDH key: %w", err)
	}

	return &cryptotypes.PrivateKeyP521{PrivateKey: ecdhKey}, nil
}

// ExportPrivateKeyP521ToHex exports a P521 private key to a hex string.
func ExportPrivateKeyP521ToHex(privKey *cryptotypes.PrivateKeyP521) string {
	return hex.EncodeToString(privKey.Bytes())
}

// ExportPublicKeyP521ToHex exports a P521 public key to a hex string.
func ExportPublicKeyP521ToHex(pubKey *cryptotypes.PublicKeyP521) string {
	return hex.EncodeToString(pubKey.Bytes())
}

// ImportPrivateKeyP521FromHex imports a P521 private key from a hex string.
func ImportPrivateKeyP521FromHex(hexKey string) (*cryptotypes.PrivateKeyP521, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	curve := ecdh.P521()
	key, err := curve.NewPrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	return &cryptotypes.PrivateKeyP521{PrivateKey: key}, nil
}

// ImportPublicKeyP521FromHex imports a P521 public key from a hex string.
func ImportPublicKeyP521FromHex(hexKey string) (*cryptotypes.PublicKeyP521, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	curve := ecdh.P521()
	key, err := curve.NewPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	return &cryptotypes.PublicKeyP521{PublicKey: key}, nil
}
