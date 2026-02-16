package kms

import (
	"context"
	"crypto/rand"
	"fmt"

	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/common/testutil"
)

// Compile-time check that MockKMSClient implements KMSClient
var _ KMSClient = (*MockKMSClient)(nil)

// Compile-time check that MockEnclaveHandle implements EnclaveHandle
var _ EnclaveHandle = (*MockEnclaveHandle)(nil)

// MockKMSClient implements KMSClient for testing without real AWS KMS.
// It simulates the KMS behavior by storing a "simulated" data key that
// can be "encrypted" and "decrypted" without actual cryptographic operations.
type MockKMSClient struct {
	*testutil.MockFunctions

	// SimulatedDataKey is the plaintext key that will be returned.
	// In real KMS, this would be the decrypted plaintext.
	SimulatedDataKey cryptotypes.AES256Key

	// SimulatedCiphertext is the "encrypted" form of the data key.
	// In real KMS, this would be encrypted under the KMS master key.
	SimulatedCiphertext []byte

	// Call tracking for assertions
	GenerateDataKeyCalled bool
	DecryptCalled         bool
	LastKeyARN            string
	LastAttestation       []byte
}

// NewMockKMSClient creates a new mock KMS client with random simulated keys.
func NewMockKMSClient() *MockKMSClient {
	var key cryptotypes.AES256Key
	_, _ = rand.Read(key[:])

	ciphertext := make([]byte, 64)
	_, _ = rand.Read(ciphertext)

	return &MockKMSClient{
		MockFunctions:       testutil.NewMockFunctions(),
		SimulatedDataKey:    key,
		SimulatedCiphertext: ciphertext,
	}
}

// GenerateDataKeyWithAttestation simulates KMS GenerateDataKey with attestation.
// If a mocked function is set, it will be called instead of the default behavior.
func (m *MockKMSClient) GenerateDataKeyWithAttestation(
	ctx context.Context,
	keyARN string,
	attestationDoc []byte,
) (*DataKeyOutput, error) {
	m.GenerateDataKeyCalled = true
	m.LastKeyARN = keyARN
	m.LastAttestation = attestationDoc

	// Check for mocked function override
	if f, ok := m.GetMockedFunc("GenerateDataKeyWithAttestation"); ok {
		return f.(func(context.Context, string, []byte) (*DataKeyOutput, error))(ctx, keyARN, attestationDoc)
	}

	// Validate inputs
	if keyARN == "" {
		return nil, fmt.Errorf("mock KMS: keyARN is required")
	}
	if len(attestationDoc) == 0 {
		return nil, fmt.Errorf("mock KMS: attestation document is required")
	}

	// Return simulated output
	// In mock mode, CiphertextForRecipient contains the plaintext directly
	// (since we don't have real RSA encryption in the mock)
	plaintext := make([]byte, len(m.SimulatedDataKey))
	copy(plaintext, m.SimulatedDataKey[:])
	return &DataKeyOutput{
		CiphertextBlob:         m.SimulatedCiphertext,
		CiphertextForRecipient: plaintext,
	}, nil
}

// DecryptWithAttestation simulates KMS Decrypt with attestation.
// If a mocked function is set, it will be called instead of the default behavior.
func (m *MockKMSClient) DecryptWithAttestation(
	ctx context.Context,
	ciphertext []byte,
	attestationDoc []byte,
) ([]byte, error) {
	m.DecryptCalled = true
	m.LastAttestation = attestationDoc

	// Check for mocked function override
	if f, ok := m.GetMockedFunc("DecryptWithAttestation"); ok {
		return f.(func(context.Context, []byte, []byte) ([]byte, error))(ctx, ciphertext, attestationDoc)
	}

	// Validate inputs
	if len(ciphertext) == 0 {
		return nil, fmt.Errorf("mock KMS: ciphertext is required")
	}
	if len(attestationDoc) == 0 {
		return nil, fmt.Errorf("mock KMS: attestation document is required")
	}

	// Return a copy to avoid shared-memory mutation in callers.
	plaintext := make([]byte, len(m.SimulatedDataKey))
	copy(plaintext, m.SimulatedDataKey[:])
	return plaintext, nil
}

// MockEnclaveHandle implements EnclaveHandle for testing outside of a real enclave.
// It simulates attestation and key decryption without actual Nitro hardware.
type MockEnclaveHandle struct {
	*testutil.MockFunctions

	// SimulatedAttestation is the attestation document returned by Attest.
	SimulatedAttestation []byte

	// Call tracking for assertions
	AttestCalled  bool
	DecryptCalled bool
	LastUserData  []byte
}

// NewMockEnclaveHandle creates a new mock enclave handle with random simulated attestation.
func NewMockEnclaveHandle() *MockEnclaveHandle {
	attestation := make([]byte, 128)
	_, _ = rand.Read(attestation)

	return &MockEnclaveHandle{
		MockFunctions:        testutil.NewMockFunctions(),
		SimulatedAttestation: attestation,
	}
}

// Attest simulates generating an attestation document.
// If a mocked function is set, it will be called instead of the default behavior.
func (m *MockEnclaveHandle) Attest(userData []byte) ([]byte, error) {
	m.AttestCalled = true
	m.LastUserData = userData

	// Check for mocked function override
	if f, ok := m.GetMockedFunc("Attest"); ok {
		return f.(func([]byte) ([]byte, error))(userData)
	}

	return m.SimulatedAttestation, nil
}

// DecryptKMSEnvelopedKey simulates decrypting KMS enveloped data.
// In the mock, we assume the input is already the plaintext (since MockKMSClient
// doesn't actually encrypt with RSA).
// If a mocked function is set, it will be called instead of the default behavior.
func (m *MockEnclaveHandle) DecryptKMSEnvelopedKey(envelopedKey []byte) ([]byte, error) {
	m.DecryptCalled = true

	// Check for mocked function override
	if f, ok := m.GetMockedFunc("DecryptKMSEnvelopedKey"); ok {
		return f.(func([]byte) ([]byte, error))(envelopedKey)
	}

	// Validate input
	if len(envelopedKey) == 0 {
		return nil, fmt.Errorf("mock enclave: enveloped key is empty")
	}

	// In mock mode, the "enveloped" key is actually plaintext
	return envelopedKey, nil
}
