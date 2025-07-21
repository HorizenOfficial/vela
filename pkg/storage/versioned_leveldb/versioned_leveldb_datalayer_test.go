package versioned_leveldb_test

import (
	"bytes" // For comparing []byte slices
	"context"
	"errors" // For using errors.As to check custom error types
	"fmt"
	"os"            // For file system operations (creating/removing temp directories)
	"path/filepath" // For joining file paths
	"testing"

	"github.com/google/go-cmp/cmp"        // For deep comparison of structs
	"github.com/stretchr/testify/assert"  // For general assertions (Error, NoError, Nil, NotNil)
	"github.com/stretchr/testify/require" // For assertions that stop the test immediately on failure

	"github.com/horizen-pes/pkg/common"  // Your common types and interface
	"github.com/horizen-pes/pkg/storage" // Your storage package with BoltDBDataLayer and custom Error
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

// testVersionedLevelDBBaseDir is the base directory where all temporary Versioned LevelDB test files will be created.
var testVersionedLevelDBBaseDir string

// TestMain is a special function that runs before any tests in the package.
// It's used here to set up and tear down the global test environment (the base directory).
func TestMain(m *testing.M) {
	fmt.Println("Running TestMain setup for VersionedLevelDBDataLayer integration tests...")

	// Create a temporary base directory for all Versioned LevelDB test files.
	// This ensures tests are isolated from each other and from real data.
	var err error
	testVersionedLevelDBBaseDir, err = os.MkdirTemp("", "versioned_leveldb_test_dbs_")
	if err != nil {
		fmt.Printf("Failed to create Versioned LevelDB test base directory: %v\n", err)
		os.Exit(1) // Exit if setup fails
	}

	// Run all tests in the package.
	code := m.Run()

	// Clean up the base directory after all tests have completed.
	fmt.Println("Running TestMain teardown for VersionedLevelDBDataLayer integration tests...")
	err = os.RemoveAll(testVersionedLevelDBBaseDir)
	if err != nil {
		fmt.Printf("Failed to clean up Versioned LevelDB test base directory %s: %v\n", testVersionedLevelDBBaseDir, err)
	}

	os.Exit(code) // Exit with the result code from tests
}

// TestVersionedLevelDBDataLayer provides a comprehensive test suite for the VersionedLevelDBDataLayer implementation.
// It uses a factory function to create a new, isolated Versioned LevelDB instance for each subtest.
func TestVersionedLevelDBDataLayer(t *testing.T) {
	// createStore is a factory function that returns a new, clean VersionedLevelDBDataLayer instance.
	// Each call to createStore will create a new temporary Versioned LevelDB file.
	createStore := func() storage.ApplicationStateStore {
		// Create a unique temporary directory for this specific test instance.
		// This ensures complete isolation between different subtests.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err, "Failed to create temp directory for Versioned LevelDB")

		// Configure Versioned LevelDB to use a file within the temporary directory.
		cfg := versionedDb.VersionedLevelDBConfig{
			Path:           filepath.Join(tempDir, "test.db"), // The actual .db file
			VersionsToKeep: 5,                                 // Keep a small number of versions for testing
		}
		dl, err := versionedDb.NewVersionedLevelDBDataLayer(cfg)
		require.NoError(t, err, "Failed to create VersionedLevelDBDataLayer instance")

		// Use t.Cleanup to ensure the database is closed and its temporary directory is removed
		// after the current test (or subtest) finishes, regardless of pass/fail.
		t.Cleanup(func() {
			// Close the database connection.
			require.NoError(t, dl.Close(), "Closing Versioned LevelDB store should not error during cleanup")
			// Remove the temporary directory and all its contents.
			require.NoError(t, os.RemoveAll(tempDir), "Failed to remove Versioned LevelDB test directory: %s", tempDir)
		})

		return dl
	}

	// Context for database operations.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is canceled when tests complete

	// --- Test ApplicationState Operations ---
	t.Run("StoreAndGetApplicationState", func(t *testing.T) {
		store := createStore() // Get a fresh Versioned LevelDB instance for this subtest

		expectedState := &common.ApplicationState{
			ApplicationID:  "versioned-leveldb-app-id-1",
			StateRoot:      []byte("versioned-leveldb-root-hash-1"),
			EncryptedState: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		}

		err := store.StoreApplicationState(ctx, expectedState)
		require.NoError(t, err, "StoreApplicationState should not return an error")

		actualState, err := store.GetApplicationState(ctx, expectedState.ApplicationID)
		require.NoError(t, err, "GetApplicationState for existing ID should not return an error")
		require.NotNil(t, actualState, "GetApplicationState should return a non-nil state")

		// Compare the retrieved state with the expected one using go-cmp for deep equality.
		if diff := cmp.Diff(expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentApplicationState", func(t *testing.T) {
		store := createStore() // Get a fresh Versioned LevelDB instance

		_, err := store.GetApplicationState(ctx, "versioned-leveldb-non-existent-app-id")
		require.Error(t, err, "Expected an error when getting a non-existent application state")

		// Verify that the error is specifically our custom "not found" error.
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not_found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test WASM Bytecode Operations ---
	t.Run("StoreAndGetWASMBytecode", func(t *testing.T) {
		store := createStore()

		appID := "versioned-leveldb-wasm-app-id-1"
		expectedBytecode := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04}

		err := store.StoreWASMBytecode(ctx, appID, expectedBytecode)
		require.NoError(t, err, "StoreWASMBytecode should not return an error")

		actualBytecode, err := store.GetWASMBytecode(ctx, appID)
		require.NoError(t, err, "GetWASMBytecode for existing ID should not return an error")
		require.NotNil(t, actualBytecode, "GetWASMBytecode should return non-nil bytecode")

		// Use bytes.Equal for byte slice comparison.
		assert.True(t, bytes.Equal(expectedBytecode, actualBytecode), "Retrieved WASM bytecode mismatch")
	})

	t.Run("GetNonExistentWASMBytecode", func(t *testing.T) {
		store := createStore()

		_, err := store.GetWASMBytecode(ctx, "versioned-leveldb-non-existent-wasm-id")
		require.Error(t, err, "Expected an error when getting non-existent WASM bytecode")

		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not_found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test DeanonymizationReport Operations ---
	t.Run("StoreAndGetDeanonymizationReport", func(t *testing.T) {
		store := createStore()

		expectedReport := &common.DeanonymizationReport{
			ApplicationID:   "id-001",
			ReportID:        "versioned-leveldb-report-id-1",
			EncryptedReport: []byte("some-test-root-hash-1"),
		}

		err := store.StoreDeanonymizationReport(ctx, expectedReport)
		require.NoError(t, err, "StoreDeanonymizationReport should not return an error")

		actualReport, err := store.GetDeanonymizationReport(ctx, expectedReport.ReportID)
		require.NoError(t, err, "GetDeanonymizationReport for existing ID should not return an error")
		require.NotNil(t, actualReport, "GetDeanonymizationReport should return a non-nil report")

		if diff := cmp.Diff(expectedReport, actualReport); diff != "" {
			t.Errorf("Retrieved DeanonymizationReport mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentDeanonymizationReport", func(t *testing.T) {
		store := createStore()

		_, err := store.GetDeanonymizationReport(ctx, "versioned-leveldb-non-existent-report-id")
		require.Error(t, err, "Expected an error when getting non-existent deanonymization report")

		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not_found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test UserKey Operations ---
	t.Run("StoreAndGetUserKey", func(t *testing.T) {
		store := createStore()

		userID := "versioned-leveldb-user-id-1"
		expectedPublicKey := []byte("versioned-leveldb-public-key-bytes-1")

		err := store.StoreUserKey(ctx, userID, expectedPublicKey)
		require.NoError(t, err, "StoreUserKey should not return an error")

		actualPublicKey, err := store.GetUserKey(ctx, userID)
		require.NoError(t, err, "GetUserKey for existing ID should not return an error")
		require.NotNil(t, actualPublicKey, "GetUserKey should return non-nil public key")

		assert.True(t, bytes.Equal(expectedPublicKey, actualPublicKey), "Retrieved User Key mismatch")
	})

	t.Run("GetNonExistentUserKey", func(t *testing.T) {
		store := createStore()

		_, err := store.GetUserKey(ctx, "versioned-leveldb-non-existent-user-id")
		require.Error(t, err, "Expected an error when getting non-existent user key")

		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not_found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test Close Method and Operations After Close ---
	t.Run("OperationsAfterClose", func(t *testing.T) {
		store := createStore() // Get a fresh Versioned LevelDB instance for this subtest

		// Close the Versioned LevelDB store.
		err := store.Close()
		require.NoError(t, err, "Closing the Versioned LevelDB store should not return an error on first close")

		// Try to perform an operation (e.g., GetApplicationState) after closing.
		// This should result in an error from LevelDB.
		_, err = store.GetApplicationState(ctx, "any-id")
		require.Error(t, err, "Expected an error when getting application state from a closed store")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		// Try to perform a store operation (e.g., StoreApplicationState) after closing.
		someState := &common.ApplicationState{ApplicationID: "test-after-close", StateRoot: []byte("a")}
		err = store.StoreApplicationState(ctx, someState)
		require.Error(t, err, "Expected an error when storing application state to a closed store")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		// Add similar checks for other Store/Get methods after Close()
		_, err = store.GetWASMBytecode(ctx, "any-id")
		assert.Error(t, err, "Expected error when getting WASM bytecode from closed store")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		err = store.StoreDeanonymizationReport(ctx, &common.DeanonymizationReport{ReportID: "test"})
		assert.Error(t, err, "Expected error when storing report to closed store")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		_, err = store.GetUserKey(ctx, "any-id")
		assert.Error(t, err, "Expected error when getting user key from closed store")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		err = store.StoreUserKey(ctx, "test-user", []byte("key"))
		assert.Error(t, err, "Expected error when storing user key to closed store")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

	})
}
