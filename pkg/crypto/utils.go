package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/horizen-pes/pkg/common"
	"golang.org/x/crypto/hkdf"
)

const (
	// AES-256
	aesKeySize = 32
)

/*
GeneratePrivateKeyP521 generates a private key with Elliptic Curve Diffie-Hellman over NIST P-521 curve,
also known as secp521r1.
(This can be used for encryption/decription)
*/
func GeneratePrivateKeyP521() (*common.PrivateKeyP521, error) {
	curve := ecdh.P521()
	newKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return &common.PrivateKeyP521{newKey}, nil
}

// SavePrivateKeyP521ToFileDER saves a P521 private key to a file in PKCS #8, ASN.1 DER format.
func SavePrivateKeyP521ToFileDER(privKey *common.PrivateKeyP521, filename string) error {
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
func LoadPrivateKeyP521FromFileDER(filename string) (*common.PrivateKeyP521, error) {
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

	return &common.PrivateKeyP521{ecdhKey}, nil
}

// SavePrivateKeyP521ToFilePEM saves a P521 private key to a file in PEM format.
func SavePrivateKeyP521ToFilePEM(privKey *common.PrivateKeyP521, filename string) error {
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
func LoadPrivateKeyP521FromFilePEM(filename string) (*common.PrivateKeyP521, error) {
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

	return &common.PrivateKeyP521{ecdhKey}, nil
}

/*
GeneratePrivateKey25519 generates a private key with Elliptic Curve 25519.
(This is the curve used in Ethereum)
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

/*
Encrypt encrypts a generic byte array.
Sender and receiver keys must use Elliptic Curve Diffie-Hellman over NIST curve (25519 is not suitable).
A shared derived AES key is generated with the aboves, and the message is encrypted with AES-256-GCM.
A random nonce is used for the encryption and prepended to the generated bytes.
*/
func Encrypt(senderPrivKey *common.PrivateKeyP521, receiverPubKey *common.PublicKeyP521, message []byte) ([]byte, error) {
	sharedKey, err := senderPrivKey.ECDH(receiverPubKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute shared key: %w", err)
	}
	sharedAESKey, err := deriveAES256Key(sharedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive AES key: %w", err)
	}

	return EncryptWithAES(sharedAESKey, message)
}

// EncryptWithAES encrypts a generic byte array with a given AES-256 key.
func EncryptWithAES(key common.AES256Key, message []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create new cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create new GCM: %w", err)
	}
	// generate IV (nonce) for AES-256-GCM
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, message, nil)

	// Prepend nonce at the ciphertext
	return append(nonce, ciphertext...), nil
}

/*
Decrypt decrypts a chipertext.
Sender and receiver keys must use Elliptic Curve Diffie-Hellman over NIST curve (25519 is not suitable).
The input array must contain the nonce prepended.
*/
func Decrypt(senderPubKey *common.PublicKeyP521, receiverPrivKey *common.PrivateKeyP521, message []byte) ([]byte, error) {
	sharedKey, err := receiverPrivKey.ECDH(senderPubKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute shared key: %w", err)
	}
	sharedAESKey, err := deriveAES256Key(sharedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive AES key: %w", err)
	}
	return DecryptWithAES(sharedAESKey, message)
}

// DecryptWithAES decrypts a chipertext with a given AES-256 key.
func DecryptWithAES(key common.AES256Key, message []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create new cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create new GCM: %w", err)
	}
	nonceSize := aesgcm.NonceSize()
	if len(message) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := message[:nonceSize], message[nonceSize:]

	plainMessage, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message: %w", err)
	}
	return plainMessage, nil
}

func deriveAES256Key(secret []byte) (common.AES256Key, error) {
	h := hkdf.New(sha256.New, secret, nil, nil)
	key := make([]byte, aesKeySize)
	if _, err := io.ReadFull(h, key); err != nil {
		return common.AES256Key{}, fmt.Errorf("failed to derive key: %w", err)
	}
	return common.AES256Key(key), nil
}
