package executor

import (
	"context"
	"errors"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/blockchain/testutil"
	"github.com/horizen-pes/pkg/common"
	commontestutil "github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/executor/kms"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndRestoreEnclaveKeySet_Type0(t *testing.T) {
	ctx := context.Background()

	// Test GenerateEnclaveKeySet with Type 0 (no KMS)
	generatedKeySet, recovery, err := GenerateEnclaveKeySet(ctx, 0, nil, nil, "")
	require.NoError(t, err, "GenerateEnclaveKeySet should not return an error")
	require.NotNil(t, generatedKeySet, "Generated key set should not be nil")
	require.NotNil(t, recovery, "Recovery data should not be nil")
	require.Equal(t, 0, recovery.RecoveryType, "Recovery type should be 0")

	// Test RestoreEnclaveKeySet with Type 0 (no KMS)
	restoredKeySet, err := RestoreEnclaveKeySet(ctx, recovery, nil, nil)
	require.NoError(t, err, "RestoreEnclaveKeySet should not return an error")
	require.NotNil(t, restoredKeySet, "Restored key set should not be nil")

	// Verify that the restored key set is identical to the generated key set
	generatedSerialized, err := generatedKeySet.Serialize()
	require.NoError(t, err, "Failed to serialize generated key set")

	restoredSerialized, err := restoredKeySet.Serialize()
	require.NoError(t, err, "Failed to serialize restored key set")

	require.Equal(t, generatedSerialized, restoredSerialized, "Restored key set should be identical to generated key set")
}

func TestGenerateEnclaveKeySet_UnsupportedType(t *testing.T) {
	ctx := context.Background()

	// Attempt to generate a keyset with an unsupported recovery type
	_, _, err := GenerateEnclaveKeySet(ctx, 99, nil, nil, "")

	// Verify that an error is returned
	require.Error(t, err, "GenerateEnclaveKeySet should return an error for unsupported recovery types")
	require.Contains(t, err.Error(), "unsupported recovery type: 99", "Error message should indicate unsupported type")
}

func TestRestoreEnclaveKeySet_UnsupportedType(t *testing.T) {
	ctx := context.Background()

	// Create a recovery struct with an unsupported type
	recovery := &common.EnclaveKeySetRecovery{
		RecoveryType: 99, // Unsupported type
	}

	// Attempt to restore the keyset
	_, err := RestoreEnclaveKeySet(ctx, recovery, nil, nil)

	// Verify that an error is returned
	require.Error(t, err, "RestoreEnclaveKeySet should return an error for unsupported recovery types")
	require.Contains(t, err.Error(), "unsupported recovery type: 99", "Error message should indicate unsupported type")
}

func TestCheckSignature(t *testing.T) {
	applicationId := common.NewApplicationId(1)
	execConfig, err := LoadConfig()
	require.NoError(t, err)
	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)

	ks, _ := CreateNewKeySet()

	executor := &StatelessExecutor{
		config:           execConfig,
		MsgToSignBuilder: builder,
		keySet:           ks,
		log:              testLogger,
	}
	executorAddress := ethCrypto.PubkeyToAddress(*executor.keySet.SigningKey.PublicKey().PublicKey)

	testHelper := testutil.NewSimTestHelper(t, false, false, &executorAddress, nil)
	defer testHelper.Close()

	events := [1]common.Event{{ApplicationID: applicationId, EncryptedData: []byte{0x07, 0x07, 0x07}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: ethcommon.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(1)},
	}

	updatePayload := &common.UpdatePayload{
		ApplicationID:  applicationId,
		RequestID:      commontestutil.GenerateRandomRequestID(),
		PrevStateRoot:  [32]byte{0x08, 0x05, 0x06},
		NewStateRoot:   [32]byte{0x04, 0x05, 0x06},
		Events:         events[:],
		Withdrawals:    withdrawals,
		RefundAmount:   big.NewInt(100),
		ApplicationFee: big.NewInt(100),
	}

	signature, err := executor.signUpdatePayload(updatePayload)
	require.NoError(t, err)

	// Sign the hash
	updatePayload.Signature = signature
	teeSignerContract := testHelper.GetSimTeeAuthenticatorHelper()

	result := teeSignerContract.CheckSignature(updatePayload)
	require.True(t, result, "Signature verification failed")
}

func TestDumpKeys(t *testing.T) {

	execConfig, err := LoadConfig()
	require.NoError(t, err)

	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)

	ks, err := CreateNewKeySet()
	require.NoError(t, err)

	executor := &StatelessExecutor{
		config:           execConfig,
		MsgToSignBuilder: builder,
		keySet:           ks,
		log:              testLogger,
	}

	executor.DumpPublicKeys()
	executor.DumpPrivateKeys()
}

// ============================================================================
// Type 1 (KMS) Tests using mocks
// ============================================================================

func TestGenerateAndRestoreEnclaveKeySet_Type1(t *testing.T) {
	ctx := context.Background()
	mockKMS := kms.NewMockKMSClient()
	mockEnclave := kms.NewMockEnclaveHandle()

	// Test GenerateEnclaveKeySet with Type 1 (KMS)
	generatedKeySet, recovery, err := GenerateEnclaveKeySet(
		ctx,
		1,
		mockKMS,
		mockEnclave,
		"arn:aws:kms:us-east-1:123456789:key/test-key",
	)
	require.NoError(t, err, "GenerateEnclaveKeySet should not return an error")
	require.NotNil(t, generatedKeySet, "Generated key set should not be nil")
	require.NotNil(t, recovery, "Recovery data should not be nil")
	require.Equal(t, 1, recovery.RecoveryType, "Recovery type should be 1")

	// Verify KMS was called
	require.True(t, mockKMS.GenerateDataKeyCalled, "KMS GenerateDataKey should have been called")
	require.True(t, mockEnclave.AttestCalled, "Enclave Attest should have been called")

	// Verify RecoveryCiphertext is the KMS CiphertextBlob (not plaintext key)
	require.Equal(t, mockKMS.SimulatedCiphertext, recovery.RecoveryCiphertext,
		"RecoveryCiphertext should be the KMS CiphertextBlob")

	// Test RestoreEnclaveKeySet with Type 1 (KMS)
	// Reset call tracking
	mockKMS.DecryptCalled = false
	mockEnclave.AttestCalled = false
	mockEnclave.DecryptCalled = false

	restoredKeySet, err := RestoreEnclaveKeySet(ctx, recovery, mockKMS, mockEnclave)
	require.NoError(t, err, "RestoreEnclaveKeySet should not return an error")
	require.NotNil(t, restoredKeySet, "Restored key set should not be nil")

	// Verify KMS Decrypt was called for restoration
	require.True(t, mockKMS.DecryptCalled, "KMS Decrypt should have been called")
	require.True(t, mockEnclave.AttestCalled, "Enclave Attest should have been called for restoration")

	// Verify that the restored key set is identical to the generated key set
	generatedSerialized, err := generatedKeySet.Serialize()
	require.NoError(t, err, "Failed to serialize generated key set")

	restoredSerialized, err := restoredKeySet.Serialize()
	require.NoError(t, err, "Failed to serialize restored key set")

	require.Equal(t, generatedSerialized, restoredSerialized,
		"Restored key set should be identical to generated key set")
}

func TestGenerateEnclaveKeySet_Type1_MissingKMSClient(t *testing.T) {
	ctx := context.Background()
	mockEnclave := kms.NewMockEnclaveHandle()

	// Attempt to generate with nil KMS client
	_, _, err := GenerateEnclaveKeySet(ctx, 1, nil, mockEnclave, "arn:aws:kms:test")

	require.Error(t, err, "Should return error when KMS client is nil")
	require.Contains(t, err.Error(), "KMS client is required")
}

func TestGenerateEnclaveKeySet_Type1_MissingEnclaveHandle(t *testing.T) {
	ctx := context.Background()
	mockKMS := kms.NewMockKMSClient()

	// Attempt to generate with nil enclave handle
	_, _, err := GenerateEnclaveKeySet(ctx, 1, mockKMS, nil, "arn:aws:kms:test")

	require.Error(t, err, "Should return error when enclave handle is nil")
	require.Contains(t, err.Error(), "enclave handle is required")
}

func TestGenerateEnclaveKeySet_Type1_MissingKeyARN(t *testing.T) {
	ctx := context.Background()
	mockKMS := kms.NewMockKMSClient()
	mockEnclave := kms.NewMockEnclaveHandle()

	// Attempt to generate with empty key ARN
	_, _, err := GenerateEnclaveKeySet(ctx, 1, mockKMS, mockEnclave, "")

	require.Error(t, err, "Should return error when key ARN is empty")
	require.Contains(t, err.Error(), "KMS key ARN is required")
}

func TestGenerateEnclaveKeySet_Type1_AttestationFailure(t *testing.T) {
	ctx := context.Background()
	mockKMS := kms.NewMockKMSClient()
	mockEnclave := kms.NewMockEnclaveHandle()

	// Simulate attestation failure
	mockEnclave.AddMockedFunc("Attest", func(userData []byte) ([]byte, error) {
		return nil, errors.New("NSM not available - not running in enclave")
	})

	_, _, err := GenerateEnclaveKeySet(ctx, 1, mockKMS, mockEnclave, "arn:aws:kms:test")

	require.Error(t, err, "Should return error when attestation fails")
	require.Contains(t, err.Error(), "failed to generate attestation")
}

func TestGenerateEnclaveKeySet_Type1_KMSFailure(t *testing.T) {
	ctx := context.Background()
	mockKMS := kms.NewMockKMSClient()
	mockEnclave := kms.NewMockEnclaveHandle()

	// Simulate KMS failure
	mockKMS.AddMockedFunc("GenerateDataKeyWithAttestation",
		func(ctx context.Context, keyARN string, attestation []byte) (*kms.DataKeyOutput, error) {
			return nil, errors.New("KMS service unavailable")
		})

	_, _, err := GenerateEnclaveKeySet(ctx, 1, mockKMS, mockEnclave, "arn:aws:kms:test")

	require.Error(t, err, "Should return error when KMS fails")
	require.Contains(t, err.Error(), "failed to generate data key from KMS")
}

func TestRestoreEnclaveKeySet_Type1_KMSDecryptFailure(t *testing.T) {
	ctx := context.Background()
	mockKMS := kms.NewMockKMSClient()
	mockEnclave := kms.NewMockEnclaveHandle()

	// First generate a valid recovery
	keySet, recovery, err := GenerateEnclaveKeySet(ctx, 1, mockKMS, mockEnclave, "arn:aws:kms:test")
	require.NoError(t, err)
	require.NotNil(t, keySet)

	// Simulate KMS decrypt failure on restore (e.g., PCR mismatch)
	mockKMS.AddMockedFunc("DecryptWithAttestation",
		func(ctx context.Context, ciphertext, attestation []byte) ([]byte, error) {
			return nil, errors.New("access denied - PCR mismatch")
		})

	_, err = RestoreEnclaveKeySet(ctx, recovery, mockKMS, mockEnclave)

	require.Error(t, err, "Should return error when KMS decrypt fails")
	require.Contains(t, err.Error(), "failed to decrypt master key from KMS")
}

func TestRestoreEnclaveKeySet_Type1_MissingDependencies(t *testing.T) {
	ctx := context.Background()
	mockKMS := kms.NewMockKMSClient()
	mockEnclave := kms.NewMockEnclaveHandle()

	// Generate valid Type 1 recovery
	_, recovery, err := GenerateEnclaveKeySet(ctx, 1, mockKMS, mockEnclave, "arn:aws:kms:test")
	require.NoError(t, err)

	t.Run("missing KMS client", func(t *testing.T) {
		_, err := RestoreEnclaveKeySet(ctx, recovery, nil, mockEnclave)
		require.Error(t, err)
		require.Contains(t, err.Error(), "KMS client is required")
	})

	t.Run("missing enclave handle", func(t *testing.T) {
		_, err := RestoreEnclaveKeySet(ctx, recovery, mockKMS, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "enclave handle is required")
	})
}
