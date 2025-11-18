package executor

import (
	"testing"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/blockchain/testutil"
	"github.com/horizen-pes/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndRestoreEnclaveKeySet(t *testing.T) {
	// Test GenerateEnclaveKeySet
	generatedKeySet, recovery, err := GenerateEnclaveKeySet(0)
	require.NoError(t, err, "GenerateEnclaveKeySet should not return an error")
	require.NotNil(t, generatedKeySet, "Generated key set should not be nil")
	require.NotNil(t, recovery, "Recovery data should not be nil")

	// Test RestoreEnclaveKeySet
	restoredKeySet, err := RestoreEnclaveKeySet(recovery)
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
	// Attempt to generate a keyset with an unsupported recovery type
	_, _, err := GenerateEnclaveKeySet(99)

	// Verify that an error is returned
	require.Error(t, err, "GenerateEnclaveKeySet should return an error for unsupported recovery types")
	require.Contains(t, err.Error(), "unsupported recovery type: 99", "Error message should indicate unsupported type")
}

func TestRestoreEnclaveKeySet_UnsupportedType(t *testing.T) {
	// Create a recovery struct with an unsupported type
	recovery := &common.EnclaveKeySetRecovery{
		RecoveryType: 99, // Unsupported type
	}

	// Attempt to restore the keyset
	_, err := RestoreEnclaveKeySet(recovery)

	// Verify that an error is returned
	require.Error(t, err, "RestoreEnclaveKeySet should return an error for unsupported recovery types")
	require.Contains(t, err.Error(), "unsupported recovery type: 99", "Error message should indicate unsupported type")
}

func TestCheckSignature(t *testing.T) {
	applicationId := "1"
	execConfig := DefaultConfig()

	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)

	qqq, _ := CreateNewKeySet()

	executor := &StatelessExecutor{
		config:           execConfig,
		MsgToSignBuilder: builder,
		keySet:           qqq,
		log:              testLogger,
	}
	executorAddress := ethCrypto.PubkeyToAddress(*executor.keySet.SigningKey.PublicKey().PublicKey)

	testHelper := testutil.NewSimTestHelper(t, false, false, &executorAddress, nil)
	defer testHelper.Close()

	events := [1]common.Event{{ApplicationID: applicationId, EncryptedData: []byte{0x07, 0x07, 0x07}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: "0x1234567890123456789012345678901234567890", Amount: 1},
	}

	updatePayload := &common.UpdatePayload{
		ApplicationID: applicationId,
		RequestID:     "7890",
		PrevStateRoot: [32]byte{0x08, 0x05, 0x06},
		NewStateRoot:  [32]byte{0x04, 0x05, 0x06},
		Events:        events[:],
		Withdrawals:   withdrawals,
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

	execConfig := DefaultConfig()

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
