package versioned_leveldb_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/common"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

var testVersionedLevelDBBaseDir string

func TestMain(m *testing.M) {
	fmt.Println("Running TestMain setup for VersionedLevelDBDataLayer integration tests...")
	var err error
	testVersionedLevelDBBaseDir, err = os.MkdirTemp("", "versioned_leveldb_test_dbs_")
	if err != nil {
		fmt.Printf("Failed to create Versioned LevelDB test base directory: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	fmt.Println("Running TestMain teardown for VersionedLevelDBDataLayer integration tests...")
	err = os.RemoveAll(testVersionedLevelDBBaseDir)
	if err != nil {
		fmt.Printf("Failed to clean up Versioned LevelDB test base directory %s: %v\n", testVersionedLevelDBBaseDir, err)
	}
	os.Exit(code)
}

func TestVersionedLevelDBDataLayer(t *testing.T) {
	createStore := func(t *testing.T, dbPath string, versionsToKeep int) *versionedDb.VersionedLevelDBDataLayer {
		cfg := versionedDb.VersionedLevelDBConfig{
			DBPath:         dbPath,
			VersionsToKeep: versionsToKeep,
		}
		dl, err := versionedDb.NewVersionedLevelDBDataLayer(cfg)
		require.NoError(t, err, "Failed to create VersionedLevelDBDataLayer instance")

		t.Cleanup(func() {
			require.NoError(t, dl.Close(), "Closing Versioned LevelDB store should not error during cleanup")
		})

		return dl
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("NewVersionedLevelDBDataLayerWithInvalidPath", func(t *testing.T) {
		// Verifies that creating a new data layer with an invalid or inaccessible
		// path returns an error.
		_, err := versionedDb.NewVersionedLevelDBDataLayer(versionedDb.VersionedLevelDBConfig{
			DBPath:         "/invalid-path",
			VersionsToKeep: 5,
		})
		require.Error(t, err)
	})

	t.Run("StoreAndGetApplicationStateWithCorruptedData", func(t *testing.T) {
		// Verifies that attempting to retrieve an application state that has been
		// manually corrupted (e.g., invalid JSON) results in an error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		appID := "corrupted-app-id"

		// Manually insert corrupted data into the database.
		err = store.StoreWASMBytecode(ctx, appID, []byte("corrupted-json"))
		require.NoError(t, err, "Storing corrupted data should not fail")

		_, err = store.GetApplicationState(ctx, appID)
		require.Error(t, err, "Expected an error when getting corrupted application state")
	})

	t.Run("StoreAndGetApplicationState", func(t *testing.T) {
		// Tests the fundamental store and get operations for an ApplicationState,
		// ensuring data is retrieved correctly and without errors.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		expectedState := &common.ApplicationState{
			ApplicationID:  "versioned-leveldb-app-id-1",
			StateRoot:      []byte("versioned-leveldb-root-hash-1"),
			EncryptedState: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		}
		err = store.StoreApplicationState(ctx, expectedState)
		require.NoError(t, err, "StoreApplicationState should not return an error")
		actualState, err := store.GetApplicationState(ctx, expectedState.ApplicationID)
		require.NoError(t, err, "GetApplicationState for existing ID should not return an error")
		require.NotNil(t, actualState, "GetApplicationState should return a non-nil state")
		if diff := cmp.Diff(expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentApplicationState", func(t *testing.T) {
		// Ensures that trying to get an application state with an ID that does not
		// exist returns a 'NotFound' error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		_, err = store.GetApplicationState(ctx, "versioned-leveldb-non-existent-app-id")
		require.Error(t, err, "Expected an error when getting a non-existent application state")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetWASMBytecode", func(t *testing.T) {
		// Tests the fundamental store and get operations for WASM bytecode,
		// ensuring data is retrieved correctly and without errors.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		appID := "versioned-leveldb-wasm-app-id-1"
		expectedBytecode := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04}
		err = store.StoreWASMBytecode(ctx, appID, expectedBytecode)
		require.NoError(t, err, "StoreWASMBytecode should not return an error")
		actualBytecode, err := store.GetWASMBytecode(ctx, appID)
		require.NoError(t, err, "GetWASMBytecode for existing ID should not return an error")
		require.NotNil(t, actualBytecode, "GetWASMBytecode should return non-nil bytecode")
		assert.True(t, bytes.Equal(expectedBytecode, actualBytecode), "Retrieved WASM bytecode mismatch")
	})

	t.Run("GetNonExistentWASMBytecode", func(t *testing.T) {
		// Ensures that trying to get WASM bytecode with an ID that does not
		// exist returns a 'NotFound' error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		_, err = store.GetWASMBytecode(ctx, "versioned-leveldb-non-existent-wasm-id")
		require.Error(t, err, "Expected an error when getting non-existent WASM bytecode")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetDeanonymizationReport", func(t *testing.T) {
		// Tests the fundamental store and get operations for a DeanonymizationReport,
		// ensuring data is retrieved correctly and without errors.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		expectedReport := &common.DeanonymizationReport{
			ApplicationID:   "id-001",
			ReportID:        "versioned-leveldb-report-id-1",
			EncryptedReport: []byte("some-test-root-hash-1"),
		}
		err = store.StoreDeanonymizationReport(ctx, expectedReport)
		require.NoError(t, err, "StoreDeanonymizationReport should not return an error")
		actualReport, err := store.GetDeanonymizationReport(ctx, expectedReport.ReportID)
		require.NoError(t, err, "GetDeanonymizationReport for existing ID should not return an error")
		require.NotNil(t, actualReport, "GetDeanonymizationReport should return a non-nil report")
		if diff := cmp.Diff(expectedReport, actualReport); diff != "" {
			t.Errorf("Retrieved DeanonymizationReport mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentDeanonymizationReport", func(t *testing.T) {
		// Ensures that trying to get a deanonymization report with an ID that
		// does not exist returns a 'NotFound' error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		_, err = store.GetDeanonymizationReport(ctx, "versioned-leveldb-non-existent-report-id")
		require.Error(t, err, "Expected an error when getting non-existent deanonymization report")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetUserKey", func(t *testing.T) {
		// Tests the fundamental store and get operations for a user's public key,
		// ensuring data is retrieved correctly and without errors.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		userID := "versioned-leveldb-user-id-1"
		expectedPublicKey := []byte("versioned-leveldb-public-key-bytes-1")
		err = store.StoreUserKey(ctx, userID, expectedPublicKey)
		require.NoError(t, err, "StoreUserKey should not return an error")
		actualPublicKey, err := store.GetUserKey(ctx, userID)
		require.NoError(t, err, "GetUserKey for existing ID should not return an error")
		require.NotNil(t, actualPublicKey, "GetUserKey should return non-nil public key")
		assert.True(t, bytes.Equal(expectedPublicKey, actualPublicKey), "Retrieved User Key mismatch")
	})

	t.Run("GetNonExistentUserKey", func(t *testing.T) {
		// Ensures that trying to get a user key with an ID that does not
		// exist returns a 'NotFound' error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		_, err = store.GetUserKey(ctx, "versioned-leveldb-non-existent-user-id")
		require.Error(t, err, "Expected an error when getting non-existent user key")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("OperationsAfterClose", func(t *testing.T) {
		// Verifies that all data manipulation and retrieval operations return a
		// 'StorageIsClosed' error after the database instance has been closed.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		require.NoError(t, store.Close(), "Closing the Versioned LevelDB store should not return an error on first close")

		operations := map[string]func() error{
			"GetApplicationState": func() error {
				_, err := store.GetApplicationState(ctx, "any-id")
				return err
			},
			"StoreApplicationState": func() error {
				return store.StoreApplicationState(ctx, &common.ApplicationState{ApplicationID: "test-after-close"})
			},
			"GetWASMBytecode": func() error {
				_, err := store.GetWASMBytecode(ctx, "any-id")
				return err
			},
			"StoreWASMBytecode": func() error {
				return store.StoreWASMBytecode(ctx, "any-id", []byte("bytecode"))
			},
			"GetDeanonymizationReport": func() error {
				_, err := store.GetDeanonymizationReport(ctx, "any-id")
				return err
			},
			"StoreDeanonymizationReport": func() error {
				return store.StoreDeanonymizationReport(ctx, &common.DeanonymizationReport{ReportID: "test"})
			},
			"GetUserKey": func() error {
				_, err := store.GetUserKey(ctx, "any-id")
				return err
			},
			"StoreUserKey": func() error {
				return store.StoreUserKey(ctx, "test-user", []byte("key"))
			},
		}

		for name, op := range operations {
			t.Run(name, func(t *testing.T) {
				err := op()
				require.Error(t, err, "Expected an error from a closed store")
				var closedErr *storageErrors.Error
				if assert.True(t, errors.As(err, &closedErr), "Error should be a storage error") {
					assert.Equal(t, storageErrors.StorageIsClosed, closedErr.Code, "Error code should be StorageIsClosed")
				}
			})
		}
	})

	t.Run("StoreWithEmptyKeyAndValue", func(t *testing.T) {
		// Checks that the database correctly handles storing and retrieving an
		// entry with an empty key and an empty value.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		err = store.StoreUserKey(ctx, "", []byte{})
		require.NoError(t, err, "Storing an empty key and value should not produce an error")

		val, err := store.GetUserKey(ctx, "")
		require.NoError(t, err, "Getting an empty key should not produce an error")
		assert.Equal(t, []byte{}, val, "Expected an empty value")
	})

	t.Run("StoreWithLargeValue", func(t *testing.T) {
		// Tests the database's ability to handle large values (1MB)
		// without errors or data corruption.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		largeValue := make([]byte, 1024*1024) // 1MB
		for i := range largeValue {
			largeValue[i] = byte(i % 256)
		}

		err = store.StoreWASMBytecode(ctx, "large-value-app", largeValue)
		require.NoError(t, err, "Storing a large value should not produce an error")

		retrievedValue, err := store.GetWASMBytecode(ctx, "large-value-app")
		require.NoError(t, err, "Getting a large value should not produce an error")
		assert.True(t, bytes.Equal(largeValue, retrievedValue), "Retrieved large value should match the original")
	})

	t.Run("ReadWriteOnPersistentDB", func(t *testing.T) {
		// Tests data persistence across database instances. It writes data with one
		// instance, closes it, then verifies the data is readable by a new instance.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "persistent-db-test-")
		require.NoError(t, err, "Failed to create temp directory for persistent DB test")
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "test.db")
		cfg := versionedDb.VersionedLevelDBConfig{
			DBPath:         dbPath,
			VersionsToKeep: 5,
		}

		// Create and use the first instance of the DB
		dl1, err := versionedDb.NewVersionedLevelDBDataLayer(cfg)
		require.NoError(t, err, "Failed to create first VersionedLevelDBDataLayer instance")

		expectedState := &common.ApplicationState{
			ApplicationID:  "persistent-app-id",
			StateRoot:      []byte("persistent-root-hash"),
			EncryptedState: []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E},
		}
		err = dl1.StoreApplicationState(ctx, expectedState)
		require.NoError(t, err, "StoreApplicationState on first instance should not return an error")

		// Close the first instance
		require.NoError(t, dl1.Close(), "Closing first DB instance should not error")

		// Create a second instance of the DB using the same path
		dl2, err := versionedDb.NewVersionedLevelDBDataLayer(cfg)
		require.NoError(t, err, "Failed to create second VersionedLevelDBDataLayer instance")
		defer dl2.Close()

		// Verify that the data is still there
		actualState, err := dl2.GetApplicationState(ctx, expectedState.ApplicationID)
		require.NoError(t, err, "GetApplicationState on second instance should not return an error")
		require.NotNil(t, actualState, "GetApplicationState on second instance should return a non-nil state")
		if diff := cmp.Diff(expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("ReopenWithDifferentVersionsToKeep", func(t *testing.T) {
		// Tests the version pruning mechanism. It stores more versions than the initial
		// limit, reopens with a smaller limit, and verifies a new write triggers
		// pruning to the new, lower limit.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "persistent-db-versions-test-")
		require.NoError(t, err, "Failed to create temp directory for persistent DB test")
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "test.db")

		// Create and use the first instance of the DB
		dl1 := createStore(t, dbPath, 10)

		appID := "persistent-app-id-versions"
		// Store 15 versions
		for i := 0; i < 15; i++ {
			state := &common.ApplicationState{
				ApplicationID:  appID,
				StateRoot:      []byte(fmt.Sprintf("root-hash-%d", i)),
				EncryptedState: []byte{byte(i)},
			}
			err := dl1.StoreApplicationState(ctx, state)
			require.NoError(t, err, "StoreApplicationState on first instance should not return an error")
		}

		// We should have 10 versions now
		require.Equal(t, 10, dl1.NumberOfVersions())

		// Close the first instance
		require.NoError(t, dl1.Close(), "Closing first DB instance should not error")

		// Create a second instance of the DB using the same path but with fewer versions to keep
		dl2 := createStore(t, dbPath, 5)
		defer dl2.Close()

		// The number of versions should still be 10 because pruning only happens on write
		require.Equal(t, 10, dl2.NumberOfVersions())

		// Store one more version, this should trigger pruning
		state := &common.ApplicationState{
			ApplicationID:  appID,
			StateRoot:      []byte("root-hash-15"),
			EncryptedState: []byte{15},
		}
		err = dl2.StoreApplicationState(ctx, state)
		require.NoError(t, err, "StoreApplicationState on second instance should not return an error")

		// Now we should have 5 versions
		require.Equal(t, 5, dl2.NumberOfVersions())
	})

	t.Run("ReopenWithMoreVersionsToKeep", func(t *testing.T) {
		// Tests that increasing the number of versions to keep works as expected.
		// It stores up to an initial limit, reopens with a higher limit, and
		// verifies that it can then store more versions up to the new limit.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "persistent-db-more-versions-test-")
		require.NoError(t, err, "Failed to create temp directory for persistent DB test")
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "test.db")

		// Create and use the first instance of the DB
		dl1 := createStore(t, dbPath, 5)

		appID := "persistent-app-id-more-versions"
		var prunedVersionID []byte
		// Store 10 versions
		for i := 0; i < 10; i++ {
			state := &common.ApplicationState{
				ApplicationID:  appID,
				StateRoot:      []byte(fmt.Sprintf("root-hash-%d", i)),
				EncryptedState: []byte{byte(i)},
			}
			err := dl1.StoreApplicationState(ctx, state)
			require.NoError(t, err, "StoreApplicationState on first instance should not return an error")

			if i == 0 {
				// Capture the version ID of the first state, which will be pruned.
				lastVersion, err := dl1.LastVersionID()
				require.NoError(t, err)
				prunedVersionID = lastVersion
			}
		}

		// We should have 5 versions now
		require.Equal(t, 5, dl1.NumberOfVersions())

		// Close the first instance
		require.NoError(t, dl1.Close(), "Closing first DB instance should not error")

		// Create a second instance of the DB using the same path but with more versions to keep
		dl2 := createStore(t, dbPath, 10)
		defer dl2.Close()

		// The number of versions should still be 5
		require.Equal(t, 5, dl2.NumberOfVersions())

		// Store 5 more versions
		for i := 10; i < 15; i++ {
			state := &common.ApplicationState{
				ApplicationID:  appID,
				StateRoot:      []byte(fmt.Sprintf("root-hash-%d", i)),
				EncryptedState: []byte{byte(i)},
			}
			err := dl2.StoreApplicationState(ctx, state)
			require.NoError(t, err, "StoreApplicationState on second instance should not return an error")
		}

		// Now we should have 10 versions
		require.Equal(t, 10, dl2.NumberOfVersions())

		// Attempt to roll back to a version that was pruned by the first instance.
		// This should fail, as reopening does not resurrect pruned versions.
		err = dl2.Rollback(prunedVersionID)
		require.Error(t, err, "Rollback to a pruned version should fail")
		var versionNotFoundErr *storageErrors.Error
		if assert.True(t, errors.As(err, &versionNotFoundErr), "Error should be a storage error") {
			assert.Equal(t, storageErrors.VersionNotFound, versionNotFoundErr.Code, "Error code should be VersionNotFound")
		}
	})
}
