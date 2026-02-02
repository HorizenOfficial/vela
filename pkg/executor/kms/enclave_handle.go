// Package kms provides AWS KMS integration for Nitro Enclaves key management.
package kms

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"

	"github.com/hf/nsm"
	"github.com/hf/nsm/request"
)

// NitroEnclaveHandle implements EnclaveHandle using the Nitro Secure Module (NSM).
// It provides attestation generation and KMS envelope decryption capabilities
// that only work inside an AWS Nitro Enclave.
type NitroEnclaveHandle struct {
	// rsaKey is the ephemeral RSA key used for KMS envelope encryption.
	// The public key is included in the attestation document.
	rsaKey *rsa.PrivateKey
}

// NewNitroEnclaveHandle creates a new NitroEnclaveHandle.
// It generates an ephemeral RSA-2048 key pair for KMS envelope encryption.
// The public key will be included in attestation documents sent to KMS.
func NewNitroEnclaveHandle() (*NitroEnclaveHandle, error) {
	// Generate ephemeral RSA key pair for KMS communication
	// KMS will encrypt the data key using this public key
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	return &NitroEnclaveHandle{
		rsaKey: rsaKey,
	}, nil
}

// Attest generates a Nitro Enclave attestation document.
// The attestation document is a CBOR-encoded structure that includes:
// - PCR values (PCR0 = enclave image hash, PCR1 = kernel hash, PCR2 = application hash)
// - The enclave's RSA public key (for KMS to encrypt the response)
// - Optional user data (can be used for nonce/challenge-response)
//
// This function only works inside a Nitro Enclave. Outside an enclave,
// it will fail with "NSM not available".
func (h *NitroEnclaveHandle) Attest(userData []byte) ([]byte, error) {
	// Open connection to NSM (Nitro Secure Module)
	session, err := nsm.OpenDefaultSession()
	if err != nil {
		return nil, fmt.Errorf("failed to open NSM session (not running in enclave?): %w", err)
	}
	defer session.Close()

	// Marshal the RSA public key in DER format for inclusion in attestation
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&h.rsaKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RSA public key: %w", err)
	}

	// Request attestation from NSM
	// - UserData: optional application-specific data (e.g., nonce)
	// - PublicKey: RSA public key for KMS to encrypt the response
	// - Nonce: optional additional randomness
	res, err := session.Send(&request.Attestation{
		UserData:  userData,
		PublicKey: publicKeyDER,
		Nonce:     nil,
	})
	if err != nil {
		return nil, fmt.Errorf("NSM attestation request failed: %w", err)
	}

	if res.Attestation == nil || res.Attestation.Document == nil {
		return nil, fmt.Errorf("NSM returned empty attestation document")
	}

	return res.Attestation.Document, nil
}

// DecryptKMSEnvelopedKey decrypts a KMS CiphertextForRecipient blob.
// KMS encrypts the data key using RSAES-OAEP with SHA-256 using the
// public key from the attestation document.
//
// This function uses the enclave's private RSA key to decrypt the envelope.
func (h *NitroEnclaveHandle) DecryptKMSEnvelopedKey(envelopedKey []byte) ([]byte, error) {
	if envelopedKey == nil || len(envelopedKey) == 0 {
		return nil, fmt.Errorf("enveloped key is empty")
	}

	// KMS uses RSAES-OAEP with SHA-256 for envelope encryption
	plaintext, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		h.rsaKey,
		envelopedKey,
		nil, // No label
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt KMS envelope: %w", err)
	}

	return plaintext, nil
}

// GetPublicKey returns the RSA public key in DER format.
// This can be used for debugging or verification purposes.
func (h *NitroEnclaveHandle) GetPublicKey() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&h.rsaKey.PublicKey)
}
