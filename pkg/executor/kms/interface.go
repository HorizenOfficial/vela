// Package kms provides interfaces and implementations for AWS KMS integration
// with Nitro Enclaves attestation for secure key management.
package kms

import (
	"context"
)

// DataKeyOutput contains the result of a GenerateDataKey operation from KMS.
type DataKeyOutput struct {
	// CiphertextBlob is the KMS-encrypted data key.
	// This can be safely stored and later decrypted using KMS.
	CiphertextBlob []byte

	// CiphertextForRecipient is the data key encrypted for the enclave's RSA public key.
	// This can only be decrypted inside the enclave using its private key.
	CiphertextForRecipient []byte
}

// KMSClient abstracts AWS KMS operations for testability.
// The real implementation uses the Nitro Enclaves SDK to attach attestation documents
// to KMS requests, ensuring that only authorized enclaves can access the keys.
type KMSClient interface {
	// GenerateDataKeyWithAttestation generates a new AES-256 data key using AWS KMS.
	// The attestation document is attached to the request, allowing KMS to verify
	// that the request originates from an authorized enclave (by checking PCR values).
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - keyARN: The ARN of the KMS key to use for encryption
	//   - attestationDoc: The Nitro Enclave attestation document
	//
	// Returns:
	//   - DataKeyOutput containing the encrypted data key
	//   - error if the operation fails (e.g., KMS unavailable, PCR mismatch)
	GenerateDataKeyWithAttestation(ctx context.Context, keyARN string, attestationDoc []byte) (*DataKeyOutput, error)

	// DecryptWithAttestation decrypts a KMS ciphertext using attestation.
	// The attestation document is attached to prove the request comes from an authorized enclave.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - ciphertext: The KMS CiphertextBlob to decrypt
	//   - attestationDoc: The Nitro Enclave attestation document
	//
	// Returns:
	//   - CiphertextForRecipient (encrypted for enclave's RSA key)
	//   - error if the operation fails
	DecryptWithAttestation(ctx context.Context, ciphertext []byte, attestationDoc []byte) ([]byte, error)
}

// EnclaveHandle abstracts Nitro Enclave operations for testability.
// The real implementation uses the Nitro Security Module (NSM) to generate
// attestation documents and decrypt data encrypted for the enclave.
type EnclaveHandle interface {
	// Attest generates an attestation document from the Nitro Security Module.
	// The attestation document contains:
	//   - PCR values (measurements of the enclave image)
	//   - A freshly generated RSA public key
	//   - Optional user data
	//
	// This document can be verified by AWS KMS to ensure requests come from
	// an authorized enclave running the expected code.
	//
	// Parameters:
	//   - userData: Optional data to include in the attestation (can be nil)
	//
	// Returns:
	//   - The attestation document bytes
	//   - error if attestation fails (e.g., not running in an enclave)
	Attest(userData []byte) ([]byte, error)

	// DecryptKMSEnvelopedKey decrypts data that was encrypted for this enclave.
	// When KMS returns CiphertextForRecipient, it's encrypted using the RSA public key
	// from the attestation document. Only the enclave with the corresponding private key
	// can decrypt it.
	//
	// Parameters:
	//   - envelopedKey: The CiphertextForRecipient from KMS response
	//
	// Returns:
	//   - The decrypted plaintext (e.g., the AES-256 data key)
	//   - error if decryption fails
	DecryptKMSEnvelopedKey(envelopedKey []byte) ([]byte, error)
}
