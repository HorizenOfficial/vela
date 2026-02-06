package kms

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// MockKMSClient Tests
// ============================================================================

func TestNewMockKMSClient(t *testing.T) {
	mock := NewMockKMSClient()

	require.NotNil(t, mock)
	require.NotNil(t, mock.MockFunctions)
	assert.NotEmpty(t, mock.SimulatedDataKey)
	assert.NotEmpty(t, mock.SimulatedCiphertext)
	assert.False(t, mock.GenerateDataKeyCalled)
	assert.False(t, mock.DecryptCalled)
}

func TestMockKMSClient_GenerateDataKey_Success(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	output, err := mock.GenerateDataKeyWithAttestation(
		ctx,
		"arn:aws:kms:us-east-1:123456789:key/test-key",
		[]byte("mock-attestation-document"),
	)

	require.NoError(t, err)
	require.NotNil(t, output)
	assert.True(t, mock.GenerateDataKeyCalled)
	assert.Equal(t, "arn:aws:kms:us-east-1:123456789:key/test-key", mock.LastKeyARN)
	assert.Equal(t, []byte("mock-attestation-document"), mock.LastAttestation)
	assert.Equal(t, mock.SimulatedCiphertext, output.CiphertextBlob)
	assert.Equal(t, mock.SimulatedDataKey[:], output.CiphertextForRecipient)
}

func TestMockKMSClient_GenerateDataKey_MissingKeyARN(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	output, err := mock.GenerateDataKeyWithAttestation(ctx, "", []byte("attestation"))

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "keyARN is required")
}

func TestMockKMSClient_GenerateDataKey_MissingAttestation(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	output, err := mock.GenerateDataKeyWithAttestation(ctx, "arn:aws:kms:test", nil)

	require.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "attestation document is required")
}

func TestMockKMSClient_GenerateDataKey_MockedFunction(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	expectedErr := errors.New("KMS service unavailable")
	mock.AddMockedFunc("GenerateDataKeyWithAttestation",
		func(ctx context.Context, keyARN string, attestation []byte) (*DataKeyOutput, error) {
			return nil, expectedErr
		})

	_, err := mock.GenerateDataKeyWithAttestation(ctx, "arn:aws:kms:test", []byte("attestation"))

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.True(t, mock.GenerateDataKeyCalled) // Still tracked even when mocked
}

func TestMockKMSClient_GenerateDataKey_MockedFunctionCustomOutput(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	customOutput := &DataKeyOutput{
		CiphertextBlob:         []byte("custom-ciphertext"),
		CiphertextForRecipient: []byte("custom-key-for-enclave"),
	}

	mock.AddMockedFunc("GenerateDataKeyWithAttestation",
		func(ctx context.Context, keyARN string, attestation []byte) (*DataKeyOutput, error) {
			return customOutput, nil
		})

	output, err := mock.GenerateDataKeyWithAttestation(ctx, "arn:aws:kms:test", []byte("attestation"))

	require.NoError(t, err)
	assert.Equal(t, customOutput, output)
}

func TestMockKMSClient_Decrypt_Success(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	plaintext, err := mock.DecryptWithAttestation(
		ctx,
		[]byte("encrypted-ciphertext"),
		[]byte("mock-attestation"),
	)

	require.NoError(t, err)
	assert.True(t, mock.DecryptCalled)
	assert.Equal(t, mock.SimulatedDataKey[:], plaintext)
	assert.Equal(t, []byte("mock-attestation"), mock.LastAttestation)
}

func TestMockKMSClient_Decrypt_MissingCiphertext(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	_, err := mock.DecryptWithAttestation(ctx, nil, []byte("attestation"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext is required")
}

func TestMockKMSClient_Decrypt_MissingAttestation(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	_, err := mock.DecryptWithAttestation(ctx, []byte("ciphertext"), nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "attestation document is required")
}

func TestMockKMSClient_Decrypt_MockedFunctionError(t *testing.T) {
	ctx := context.Background()
	mock := NewMockKMSClient()

	expectedErr := errors.New("access denied - PCR mismatch")
	mock.AddMockedFunc("DecryptWithAttestation",
		func(ctx context.Context, ciphertext, attestation []byte) ([]byte, error) {
			return nil, expectedErr
		})

	_, err := mock.DecryptWithAttestation(ctx, []byte("ciphertext"), []byte("attestation"))

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// ============================================================================
// MockEnclaveHandle Tests
// ============================================================================

func TestNewMockEnclaveHandle(t *testing.T) {
	mock := NewMockEnclaveHandle()

	require.NotNil(t, mock)
	require.NotNil(t, mock.MockFunctions)
	assert.NotEmpty(t, mock.SimulatedAttestation)
	assert.False(t, mock.AttestCalled)
	assert.False(t, mock.DecryptCalled)
}

func TestMockEnclaveHandle_Attest_Success(t *testing.T) {
	mock := NewMockEnclaveHandle()

	doc, err := mock.Attest([]byte("user-data"))

	require.NoError(t, err)
	assert.True(t, mock.AttestCalled)
	assert.Equal(t, []byte("user-data"), mock.LastUserData)
	assert.Equal(t, mock.SimulatedAttestation, doc)
}

func TestMockEnclaveHandle_Attest_NilUserData(t *testing.T) {
	mock := NewMockEnclaveHandle()

	doc, err := mock.Attest(nil)

	require.NoError(t, err)
	assert.True(t, mock.AttestCalled)
	assert.Nil(t, mock.LastUserData)
	assert.Equal(t, mock.SimulatedAttestation, doc)
}

func TestMockEnclaveHandle_Attest_MockedFunctionError(t *testing.T) {
	mock := NewMockEnclaveHandle()

	expectedErr := errors.New("NSM not available - not running in enclave")
	mock.AddMockedFunc("Attest",
		func(userData []byte) ([]byte, error) {
			return nil, expectedErr
		})

	_, err := mock.Attest(nil)

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.True(t, mock.AttestCalled)
}

func TestMockEnclaveHandle_Attest_MockedFunctionCustomDoc(t *testing.T) {
	mock := NewMockEnclaveHandle()

	customDoc := []byte("custom-attestation-document-with-pcr0")
	mock.AddMockedFunc("Attest",
		func(userData []byte) ([]byte, error) {
			return customDoc, nil
		})

	doc, err := mock.Attest(nil)

	require.NoError(t, err)
	assert.Equal(t, customDoc, doc)
}

func TestMockEnclaveHandle_DecryptKMSEnvelopedKey_Success(t *testing.T) {
	mock := NewMockEnclaveHandle()
	inputKey := []byte("enveloped-key-from-kms")

	plaintext, err := mock.DecryptKMSEnvelopedKey(inputKey)

	require.NoError(t, err)
	assert.True(t, mock.DecryptCalled)
	// In mock mode, input passes through as output
	assert.Equal(t, inputKey, plaintext)
}

func TestMockEnclaveHandle_DecryptKMSEnvelopedKey_EmptyInput(t *testing.T) {
	mock := NewMockEnclaveHandle()

	_, err := mock.DecryptKMSEnvelopedKey(nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "enveloped key is empty")
}

func TestMockEnclaveHandle_DecryptKMSEnvelopedKey_MockedFunctionError(t *testing.T) {
	mock := NewMockEnclaveHandle()

	expectedErr := errors.New("RSA decryption failed")
	mock.AddMockedFunc("DecryptKMSEnvelopedKey",
		func(envelopedKey []byte) ([]byte, error) {
			return nil, expectedErr
		})

	_, err := mock.DecryptKMSEnvelopedKey([]byte("some-key"))

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// ============================================================================
// Integration-style Tests (Mock KMS + Mock Enclave working together)
// ============================================================================

func TestMockKMSAndEnclave_GenerateAndDecryptFlow(t *testing.T) {
	ctx := context.Background()
	mockKMS := NewMockKMSClient()
	mockEnclave := NewMockEnclaveHandle()

	// Step 1: Generate attestation
	attestationDoc, err := mockEnclave.Attest(nil)
	require.NoError(t, err)
	assert.True(t, mockEnclave.AttestCalled)

	// Step 2: Generate data key with attestation
	dataKeyOutput, err := mockKMS.GenerateDataKeyWithAttestation(
		ctx,
		"arn:aws:kms:us-east-1:123456789:key/my-key",
		attestationDoc,
	)
	require.NoError(t, err)
	assert.True(t, mockKMS.GenerateDataKeyCalled)

	// Step 3: Decrypt the enveloped key
	decryptedKey, err := mockEnclave.DecryptKMSEnvelopedKey(dataKeyOutput.CiphertextForRecipient)
	require.NoError(t, err)
	assert.True(t, mockEnclave.DecryptCalled)

	// Verify we got a valid 32-byte key
	assert.Len(t, decryptedKey, 32)
	assert.Equal(t, mockKMS.SimulatedDataKey[:], decryptedKey)

	// Step 4: Verify CiphertextBlob is stored (for recovery)
	assert.NotEmpty(t, dataKeyOutput.CiphertextBlob)
}

func TestMockKMSAndEnclave_RestoreFlow(t *testing.T) {
	ctx := context.Background()
	mockKMS := NewMockKMSClient()
	mockEnclave := NewMockEnclaveHandle()

	// Simulate stored ciphertext from previous generation
	storedCiphertext := mockKMS.SimulatedCiphertext

	// Step 1: Generate fresh attestation (new enclave instance)
	attestationDoc, err := mockEnclave.Attest(nil)
	require.NoError(t, err)

	// Step 2: Decrypt stored ciphertext using KMS with attestation
	envelopedKey, err := mockKMS.DecryptWithAttestation(ctx, storedCiphertext, attestationDoc)
	require.NoError(t, err)
	assert.True(t, mockKMS.DecryptCalled)

	// Step 3: Unwrap the key using enclave
	masterKey, err := mockEnclave.DecryptKMSEnvelopedKey(envelopedKey)
	require.NoError(t, err)

	// Verify we recovered the same key
	assert.Equal(t, mockKMS.SimulatedDataKey[:], masterKey)
}
