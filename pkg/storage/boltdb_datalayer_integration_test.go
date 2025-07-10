package storage_test // Use a separate test package to only test exported functionalities

import (
	"bytes" // For comparing []byte slices
	"context"
	"errors" // For using errors.As to check custom error types
	"fmt"
	"os"            // For file system operations (creating/removing temp directories)
	"path/filepath" // For joining file paths
	"testing"
	"time" // For BoltDB timeout option

	"github.com/google/go-cmp/cmp"        // For deep comparison of structs
	"github.com/stretchr/testify/assert"  // For general assertions (Error, NoError, Nil, NotNil)
	"github.com/stretchr/testify/require" // For assertions that stop the test immediately on failure

	"github.com/horizen-pes/pkg/common"  // Your common types and interface
	"github.com/horizen-pes/pkg/storage" // Your storage package with BoltDBDataLayer and custom Error

	berrors "go.etcd.io/bbolt/errors"
)

// testBoltDBBaseDir is the base directory where all temporary BoltDB test files will be created.
var testBoltDBBaseDir string

// TestMain is a special function that runs before any tests in the package.
// It's used here to set up and tear down the global test environment (the base directory).
func TestMain(m *testing.M) {
	fmt.Println("Running TestMain setup for BoltDBDataLayer integration tests...")

	// Create a temporary base directory for all BoltDB test files.
	// This ensures tests are isolated from each other and from real data.
	var err error
	testBoltDBBaseDir, err = os.MkdirTemp("", "boltdb_test_dbs_")
	if err != nil {
		fmt.Printf("Failed to create BoltDB test base directory: %v\n", err)
		os.Exit(1) // Exit if setup fails
	}

	// Run all tests in the package.
	code := m.Run()

	// Clean up the base directory after all tests have completed.
	fmt.Println("Running TestMain teardown for BoltDBDataLayer integration tests...")
	err = os.RemoveAll(testBoltDBBaseDir)
	if err != nil {
		fmt.Printf("Failed to clean up BoltDB test base directory %s: %v\n", testBoltDBBaseDir, err)
	}

	os.Exit(code) // Exit with the result code from tests
}

// TestBoltDBDataLayer provides a comprehensive test suite for the BoltDBDataLayer implementation.
// It uses a factory function to create a new, isolated BoltDB instance for each subtest.
func TestBoltDBDataLayer(t *testing.T) {
	// createStore is a factory function that returns a new, clean BoltDBDataLayer instance.
	// Each call to createStore will create a new temporary BoltDB file.
	createStore := func() storage.ApplicationStateStore {
		// Create a unique temporary directory for this specific test instance.
		// This ensures complete isolation between different subtests.
		tempDir, err := os.MkdirTemp(testBoltDBBaseDir, "boltdb-test-")
		require.NoError(t, err, "Failed to create temp directory for BoltDB")

		// Configure BoltDB to use a file within the temporary directory.
		cfg := storage.BoltDBConfig{
			Path:    filepath.Join(tempDir, "test.db"), // The actual .db file
			Timeout: 1 * time.Second,                   // Timeout for opening the DB
		}
		dl, err := storage.NewBoltDBDataLayer(cfg)
		require.NoError(t, err, "Failed to create BoltDBDataLayer instance")

		// Use t.Cleanup to ensure the database is closed and its temporary directory is removed
		// after the current test (or subtest) finishes, regardless of pass/fail.
		t.Cleanup(func() {
			// Close the database connection.
			require.NoError(t, dl.Close(), "Closing BoltDB store should not error during cleanup")
			// Remove the temporary directory and all its contents.
			require.NoError(t, os.RemoveAll(tempDir), "Failed to remove BoltDB test directory: %s", tempDir)
		})

		return dl
	}

	// Context for database operations.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is canceled when tests complete

	// --- Test ApplicationState Operations ---
	t.Run("StoreAndGetApplicationState", func(t *testing.T) {
		store := createStore() // Get a fresh BoltDB instance for this subtest

		expectedState := &common.ApplicationState{
			ApplicationID:  "boltdb-app-id-1",
			StateRoot:      []byte("boltdb-root-hash-1"),
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
		store := createStore() // Get a fresh BoltDB instance

		_, err := store.GetApplicationState(ctx, "boltdb-non-existent-app-id")
		require.Error(t, err, "Expected an error when getting a non-existent application state")

		// Verify that the error is specifically our custom "not found" error.
		var notFoundErr *storage.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test WASM Bytecode Operations ---
	t.Run("StoreAndGetWASMBytecode", func(t *testing.T) {
		store := createStore()

		appID := "boltdb-wasm-app-id-1"
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

		_, err := store.GetWASMBytecode(ctx, "boltdb-non-existent-wasm-id")
		require.Error(t, err, "Expected an error when getting non-existent WASM bytecode")

		var notFoundErr *storage.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test DeanonymizationReport Operations ---
	t.Run("StoreAndGetDeanonymizationReport", func(t *testing.T) {
		store := createStore()

		expectedReport := &common.DeanonymizationReport{
			ApplicationID:   "id-001",
			ReportID:        "boltdb-report-id-1",
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

		_, err := store.GetDeanonymizationReport(ctx, "boltdb-non-existent-report-id")
		require.Error(t, err, "Expected an error when getting non-existent deanonymization report")

		var notFoundErr *storage.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test UserKey Operations ---
	t.Run("StoreAndGetUserKey", func(t *testing.T) {
		store := createStore()

		userID := "boltdb-user-id-1"
		expectedPublicKey := []byte("boltdb-public-key-bytes-1")

		err := store.StoreUserKey(ctx, userID, expectedPublicKey)
		require.NoError(t, err, "StoreUserKey should not return an error")

		actualPublicKey, err := store.GetUserKey(ctx, userID)
		require.NoError(t, err, "GetUserKey for existing ID should not return an error")
		require.NotNil(t, actualPublicKey, "GetUserKey should return non-nil public key")

		assert.True(t, bytes.Equal(expectedPublicKey, actualPublicKey), "Retrieved User Key mismatch")
	})

	t.Run("GetNonExistentUserKey", func(t *testing.T) {
		store := createStore()

		_, err := store.GetUserKey(ctx, "boltdb-non-existent-user-id")
		require.Error(t, err, "Expected an error when getting non-existent user key")

		var notFoundErr *storage.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test Close Method and Operations After Close ---
	t.Run("OperationsAfterClose", func(t *testing.T) {
		store := createStore() // Get a fresh BoltDB instance for this subtest

		// Close the BoltDB store.
		err := store.Close()
		require.NoError(t, err, "Closing the BoltDB store should not return an error on first close")

		// Subsequent Close() calls should also return an error for bbolt.
		err = store.Close()
		assert.NoError(t, err, "Error not expected when trying to close an already closed store")

		// Try to perform an operation (e.g., GetApplicationState) after closing.
		// This should result in an error from BoltDB.
		_, err = store.GetApplicationState(ctx, "any-id")
		require.Error(t, err, "Expected an error when getting application state from a closed store")
		//ErrDatabaseNotOpen is returned when a DB instance is accessed before it is opened or after it is closed.
		assert.Contains(t, err.Error(), "not open", "Error message should indicate store is closed")
		assert.True(t, errors.Is(err, berrors.ErrDatabaseNotOpen), "qqq")

		// Try to perform a store operation (e.g., StoreApplicationState) after closing.
		someState := &common.ApplicationState{ApplicationID: "test-after-close", StateRoot: []byte("a")}
		err = store.StoreApplicationState(ctx, someState)
		require.Error(t, err, "Expected an error when storing application state to a closed store")
		assert.Contains(t, err.Error(), "not open", "Error message should indicate store is closed")
		assert.True(t, errors.Is(err, berrors.ErrDatabaseNotOpen), "qqq")

		// Add similar checks for other Store/Get methods after Close()
		_, err = store.GetWASMBytecode(ctx, "any-id")
		assert.Error(t, err, "Expected error when getting WASM bytecode from closed store")
		assert.Contains(t, err.Error(), "not open", "Error message should indicate store is closed")
		assert.True(t, errors.Is(err, berrors.ErrDatabaseNotOpen), "qqq")

		err = store.StoreDeanonymizationReport(ctx, &common.DeanonymizationReport{ReportID: "test"})
		assert.Error(t, err, "Expected error when storing report to closed store")
		assert.Contains(t, err.Error(), "not open", "Error message should indicate store is closed")
		assert.True(t, errors.Is(err, berrors.ErrDatabaseNotOpen), "qqq")

		_, err = store.GetUserKey(ctx, "any-id")
		assert.Error(t, err, "Expected error when getting user key from closed store")
		assert.Contains(t, err.Error(), "not open", "Error message should indicate store is closed")
		assert.True(t, errors.Is(err, berrors.ErrDatabaseNotOpen), "qqq")

		err = store.StoreUserKey(ctx, "test-user", []byte("key"))
		assert.Error(t, err, "Expected error when storing user key to closed store")
		assert.Contains(t, err.Error(), "not open", "Error message should indicate store is closed")
		assert.True(t, errors.Is(err, berrors.ErrDatabaseNotOpen), "qqq")
	})
}
