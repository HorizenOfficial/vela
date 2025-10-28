package executor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndRestoreEnclaveKeySet(t *testing.T) {
	// Test GenerateEnclaveKeySet
	generatedKeySet, recovery, err := GenerateEnclaveKeySet()
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
