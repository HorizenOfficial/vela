package versioned_leveldb_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/storage"
	storageErrors "github.com/HorizenOfficial/vela/pkg/storage/errors"
	versionedDb "github.com/HorizenOfficial/vela/pkg/storage/versioned_leveldb"
)

var testVersionedLevelDBBaseDir string

// TestMain sets up a temporary base directory for all VersionedLevelDBDataLayer
// integration tests and cleans it up after all tests in the package have run.
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

// createStore is a helper function that creates a new LevelDBDataLayer instance
// for testing. It takes a database path and the number of versions to keep,
// and ensures that the store is closed and cleaned up after the test.
func createStore(t *testing.T, dbPath string, versionsToKeep int) *versionedDb.LevelDBDataLayer {
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

// TestVersionedLevelDBDataLayer is the main test suite for the LevelDBDataLayer.
// It covers all aspects of the data layer's functionality, including data
// persistence, versioning, error handling, and concurrency.
func TestVersionedLevelDBDataLayer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("StoreAndGetApplicationStateWithCorruptedData", func(t *testing.T) {
		// Verifies that attempting to retrieve an application state that has been
		// manually corrupted (e.g., invalid JSON) results in an error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		// Use t.Cleanup instead of defer for directory removal. createStore registers
		// dl.Close() via t.Cleanup; since t.Cleanup callbacks run LIFO, registering
		// RemoveAll first ensures it runs after Close, preventing "no such file"
		// errors from LevelDB trying to flush to a deleted directory.
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		appID := common.NewApplicationId(12345)

		// Manually insert corrupted data into the database.
		key := []byte(versionedDb.TestAppStatePrefix + appID.String())
		value := []byte("corrupted-json")
		versionID := versionedDb.GenerateVersionID_ForTest(key, value)
		toUpdate := []storage.KeyValuePair{{Key: key, Value: value}}
		toRemove := [][]byte{}
		err = store.GetAdapter_ForTest().Update(uint64(appID), versionID, toUpdate, toRemove)
		require.NoError(t, err, "Storing corrupted data should not fail")

		_, err = store.GetApplicationState(ctx, appID)
		require.Error(t, err, "Expected an error when getting corrupted application state")
	})

	t.Run("StoreAndGetApplicationState", func(t *testing.T) {
		// Tests the fundamental store and get operations for an ApplicationState,
		// ensuring data is retrieved correctly and without errors.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		expectedState := common.ApplicationState{
			ApplicationID:  common.NewApplicationId(5423),
			StateRoot:      sha256.Sum256([]byte("versioned-leveldb-root-hash-1")),
			EncryptedState: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		}
		versionID := versionedDb.GenerateVersionID_ForTest([]byte(expectedState.ApplicationID.String()), expectedState.StateRoot[:])
		err = store.Store(ctx, versionID, &expectedState)
		require.NoError(t, err, "Store should not return an error")
		actualState, err := store.GetApplicationState(ctx, expectedState.ApplicationID)
		require.NoError(t, err, "GetApplicationState for existing ID should not return an error")
		require.NotNil(t, actualState, "GetApplicationState should return a non-nil state")
		if diff := cmp.Diff(&expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentApplicationState", func(t *testing.T) {
		// Ensures that trying to get an application state with an ID that does not
		// exist returns a 'NotFound' error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		_, err = store.GetApplicationState(ctx, common.NewApplicationId(888))
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
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		appID := common.NewApplicationId(725677)
		expectedBytecode := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04}
		versionID := versionedDb.GenerateVersionID_ForTest([]byte(appID.String()), expectedBytecode)
		state := common.ApplicationState{
			ApplicationID: appID,
			StateRoot:     sha256.Sum256([]byte("wasm-test-root")),
		}
		wasm := common.WASMData{
			ApplicationID: appID,
			Bytecode:      expectedBytecode,
		}
		err = store.StoreWithWasm(ctx, versionID, &state, &wasm)
		require.NoError(t, err, "StoreWithWasm should not return an error")
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
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		_, err = store.GetWASMBytecode(ctx, 111)
		require.Error(t, err, "Expected an error when getting non-existent WASM bytecode")
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
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		require.NoError(t, store.Close(), "Closing the Versioned LevelDB store should not return an error on first close")

		operations := map[string]func() error{
			"GetApplicationState": func() error {
				_, err := store.GetApplicationState(ctx, 40)
				return err
			},
			"Store": func() error {
				return store.Store(ctx, []byte("version"), &common.ApplicationState{ApplicationID: common.NewApplicationId(8632)})
			},
			"GetWASMBytecode": func() error {
				_, err := store.GetWASMBytecode(ctx, 40)
				return err
			},
		}

		for name, op := range operations {
			t.Run(name, func(t *testing.T) {
				err := op()
				require.Error(t, err, "Expected an error from a closed store")
				t.Log("Got an err as expected:", err)
				var closedErr *storageErrors.Error
				if assert.True(t, errors.As(err, &closedErr), "Error should be a storage error") {
					assert.Equal(t, storageErrors.StorageIsClosed, closedErr.Code, "Error code should be StorageIsClosed")
				}
			})
		}
	})

	t.Run("StoreWithLargeValue", func(t *testing.T) {
		// Tests the database's ability to handle large values (1MB)
		// without errors or data corruption.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		largeValue := make([]byte, 1024*1024) // 1MB
		for i := range largeValue {
			largeValue[i] = byte(i % 256)
		}
		appID := common.NewApplicationId(4444444)
		versionID := versionedDb.GenerateVersionID_ForTest([]byte("large-value-app"), largeValue)
		state := &common.ApplicationState{ApplicationID: appID, StateRoot: sha256.Sum256([]byte("large-value-root"))}
		err = store.StoreWithWasm(ctx, versionID, state, &common.WASMData{ApplicationID: appID, Bytecode: largeValue})
		require.NoError(t, err, "Storing a large value should not produce an error")

		retrievedValue, err := store.GetWASMBytecode(ctx, appID)
		require.NoError(t, err, "Getting a large value should not produce an error")
		assert.True(t, bytes.Equal(largeValue, retrievedValue), "Retrieved large value should match the original")
	})

	t.Run("ReadWriteOnPersistentDB", func(t *testing.T) {
		// Tests data persistence across database instances. It writes data with one
		// instance, closes it, then verifies the data is readable by a new instance.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "persistent-db-test-")
		require.NoError(t, err, "Failed to create temp directory for persistent DB test")
		t.Cleanup(func() { os.RemoveAll(tempDir) })

		dbPath := filepath.Join(tempDir, "test.db")
		cfg := versionedDb.VersionedLevelDBConfig{
			DBPath:         dbPath,
			VersionsToKeep: 5,
		}

		// Create and use the first instance of the DB
		dl1, err := versionedDb.NewVersionedLevelDBDataLayer(cfg)
		require.NoError(t, err, "Failed to create first VersionedLevelDBDataLayer instance")

		expectedState := common.ApplicationState{
			ApplicationID:  common.NewApplicationId(8168),
			StateRoot:      sha256.Sum256([]byte("persistent-root-hash")),
			EncryptedState: []byte{0x0A, 0x0B, 0x0C, 0x0D, 0x0E},
		}
		versionID := versionedDb.GenerateVersionID_ForTest([]byte(expectedState.ApplicationID.String()), expectedState.StateRoot[:])
		err = dl1.Store(ctx, versionID, &expectedState)
		require.NoError(t, err, "Store on first instance should not return an error")

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
		if diff := cmp.Diff(&expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("ReopenWithDifferentVersionsToKeep", func(t *testing.T) {
		// Tests the version pruning mechanism. It stores more versions than the initial
		// limit, reopens with a smaller limit, and verifies a new write triggers
		// pruning to the new, lower limit.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "persistent-db-versions-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })

		dbPath := filepath.Join(tempDir, "test.db")

		// Create and use the first instance of the DB
		dl1 := createStore(t, dbPath, 10)

		appID := common.NewApplicationId(744567)
		appIDStr := appID.String()
		// Store 15 versions
		for i := 0; i < 15; i++ {
			state := common.ApplicationState{
				ApplicationID:  appID,
				StateRoot:      sha256.Sum256([]byte(fmt.Sprintf("root-hash-%d", i))),
				EncryptedState: []byte{byte(i)},
			}
			versionID := versionedDb.GenerateVersionID_ForTest([]byte(appIDStr), state.StateRoot[:])
			err := dl1.Store(ctx, versionID, &state)
			require.NoError(t, err, "Store on first instance should not return an error")
		}

		// We should have 10 versions now
		numVersions, err := dl1.GetAdapter_ForTest().NumberOfVersions(uint64(appID))
		require.NoError(t, err)
		require.Equal(t, 10, numVersions)

		// Close the first instance
		require.NoError(t, dl1.Close(), "Closing first DB instance should not error")

		// Create a second instance of the DB using the same path but with fewer versions to keep
		dl2 := createStore(t, dbPath, 5)
		defer dl2.Close()

		// The number of versions should still be 10 because pruning only happens on write
		numVersions, err = dl2.GetAdapter_ForTest().NumberOfVersions(uint64(appID))
		require.NoError(t, err)
		require.Equal(t, 10, numVersions)

		// Store one more version, this should trigger pruning
		state := common.ApplicationState{
			ApplicationID:  appID,
			StateRoot:      sha256.Sum256([]byte("root-hash-15")),
			EncryptedState: []byte{15},
		}
		versionID := versionedDb.GenerateVersionID_ForTest([]byte(appIDStr), state.StateRoot[:])
		err = dl2.Store(ctx, versionID, &state)
		require.NoError(t, err, "Store on second instance should not return an error")

		// Now we should have 5 versions
		numVersions, err = dl2.GetAdapter_ForTest().NumberOfVersions(uint64(appID))
		require.NoError(t, err)
		require.Equal(t, 5, numVersions)
	})

	t.Run("ReopenWithMoreVersionsToKeep", func(t *testing.T) {
		// Tests that increasing the number of versions to keep works as expected.
		// It stores up to an initial limit, reopens with a higher limit, and
		// verifies that it can then store more versions up to the new limit.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "persistent-db-more-versions-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })

		dbPath := filepath.Join(tempDir, "test.db")

		// Create and use the first instance of the DB
		dl1 := createStore(t, dbPath, 5)

		appID := common.NewApplicationId(985009)
		appIdStr := appID.String()
		var prunedVersionID []byte
		// Store 10 versions
		for i := 0; i < 10; i++ {
			state := common.ApplicationState{
				ApplicationID:  appID,
				StateRoot:      sha256.Sum256([]byte(fmt.Sprintf("root-hash-%d", i))),
				EncryptedState: []byte{byte(i)},
			}
			versionID := versionedDb.GenerateVersionID_ForTest([]byte(appIdStr), state.StateRoot[:])
			err := dl1.Store(ctx, versionID, &state)
			require.NoError(t, err, "Store on first instance should not return an error")

			if i == 0 {
				// Capture the version ID of the first state, which will be pruned.
				lastVersion, err := dl1.LastVersionID(appID)
				require.NoError(t, err)
				prunedVersionID = lastVersion
			}
		}

		// We should have 5 versions now
		numVersions, err := dl1.GetAdapter_ForTest().NumberOfVersions(uint64(appID))
		require.NoError(t, err)
		require.Equal(t, 5, numVersions)

		// Close the first instance
		require.NoError(t, dl1.Close(), "Closing first DB instance should not error")

		// Create a second instance of the DB using the same path but with more versions to keep
		dl2 := createStore(t, dbPath, 10)
		defer dl2.Close()

		// The number of versions should still be 5
		numVersions, err = dl2.GetAdapter_ForTest().NumberOfVersions(uint64(appID))
		require.NoError(t, err)
		require.Equal(t, 5, numVersions)

		// Store 5 more versions
		for i := 10; i < 15; i++ {
			state := common.ApplicationState{
				ApplicationID:  appID,
				StateRoot:      sha256.Sum256([]byte(fmt.Sprintf("root-hash-%d", i))),
				EncryptedState: []byte{byte(i)},
			}
			versionID := versionedDb.GenerateVersionID_ForTest([]byte(appIdStr), state.StateRoot[:])
			err = dl2.Store(ctx, versionID, &state)
			require.NoError(t, err, "Store on second instance should not return an error")
		}

		// Now we should have 10 versions
		numVersions, err = dl2.GetAdapter_ForTest().NumberOfVersions(uint64(appID))
		require.NoError(t, err)
		require.Equal(t, 10, numVersions)

		// Attempt to roll back to a version that was pruned by the first instance.
		// This should fail, as reopening does not resurrect pruned versions.
		err = dl2.Rollback(appID, prunedVersionID)
		require.Error(t, err, "Rollback to a pruned version should fail")
		var versionNotFoundErr *storageErrors.Error
		if assert.True(t, errors.As(err, &versionNotFoundErr), "Error should be a storage error") {
			assert.Equal(t, storageErrors.VersionNotFound, versionNotFoundErr.Code, "Error code should be VersionNotFound")
		}
	})

	t.Run("ListVersions", func(t *testing.T) {
		// Verifies that ListVersions returns the correct number of versions in the
		// correct order (newest first).
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "store-twice-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })

		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)

		// Initially, there should be no versions
		listAppID := common.NewApplicationId(54)
		versions, err := store.ListVersions(listAppID)
		if err != nil {
			t.Fatalf("ListVersions failed: %v", err)
		}
		if len(versions) != 0 {
			t.Fatalf("expected no versions, got %d", len(versions))
		}

		// Store two versions
		v1 := versionedDb.GenerateVersionID_ForTest([]byte("v1"), []byte("d1"))
		v2 := versionedDb.GenerateVersionID_ForTest([]byte("v2"), []byte("d2"))

		state := &common.ApplicationState{
			ApplicationID:  common.NewApplicationId(54),
			StateRoot:      sha256.Sum256([]byte("state1")),
			EncryptedState: []byte("encs1"),
		}
		wasm := &common.WASMData{
			ApplicationID: common.NewApplicationId(54),
			Bytecode:      []byte("wasm1"),
		}

		if err := store.StoreWithWasm(ctx, v1, state, wasm); err != nil {
			t.Fatalf("failed to store v1: %v", err)
		}
		if err := store.StoreWithWasm(ctx, v2, state, wasm); err != nil {
			t.Fatalf("failed to store v2: %v", err)
		}

		// Check that ListVersions returns them in reverse order (newest first)
		versions, err = store.ListVersions(listAppID)
		if err != nil {
			t.Fatalf("ListVersions failed: %v", err)
		}

		if len(versions) != 2 {
			t.Fatalf("expected 2 versions, got %d", len(versions))
		}

		if !bytes.Equal(versions[0], v2) {
			t.Errorf("expected latest version to be %s, got %s", v2, versions[0])
		}
		if !bytes.Equal(versions[1], v1) {
			t.Errorf("expected older version to be %s, got %s", v1, versions[1])
		}
	})

	t.Run("RollbackAndLastVersionID", func(t *testing.T) {
		// Tests the Rollback functionality and verifies that LastVersionID is
		// updated correctly after a rollback.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "rollback-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })

		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)

		// 2. Store first version
		appID := common.NewApplicationId(9999)

		// 1. Test LastVersionID on empty db
		_, err = store.LastVersionID(appID)
		require.Error(t, err, "Expected an error when getting last version from empty storage")
		appIDSlice := []byte(appID.String())
		state1 := common.ApplicationState{
			ApplicationID: appID,
			StateRoot:     sha256.Sum256([]byte("root1")),
		}
		versionID1 := versionedDb.GenerateVersionID_ForTest(appIDSlice, state1.StateRoot[:])
		err = store.Store(ctx, versionID1, &state1)
		require.NoError(t, err)
		v1, err := store.LastVersionID(appID)
		require.NoError(t, err)

		// 3. Store second version
		state2 := common.ApplicationState{
			ApplicationID: appID,
			StateRoot:     sha256.Sum256([]byte("root2")),
		}
		versionID2 := versionedDb.GenerateVersionID_ForTest(appIDSlice, state2.StateRoot[:])
		err = store.Store(ctx, versionID2, &state2)
		require.NoError(t, err)
		v2, err := store.LastVersionID(appID)
		require.NoError(t, err)
		assert.NotEqual(t, v1, v2)

		// 4. Rollback to first version
		err = store.Rollback(appID, v1)
		require.NoError(t, err)

		// 5. Check last version is now v1
		lastVersion, err := store.LastVersionID(appID)
		require.NoError(t, err)
		assert.Equal(t, v1, lastVersion)

		// 6. Check that state is rolled back
		retrievedState, err := store.GetApplicationState(ctx, appID)
		require.NoError(t, err)
		assert.Equal(t, state1.StateRoot, retrievedState.StateRoot)

		// 7. Test rollback on closed db
		require.NoError(t, store.Close())
		err = store.Rollback(appID, v1)
		require.Error(t, err)
		var closedErr *storageErrors.Error
		if assert.True(t, errors.As(err, &closedErr), "Error should be a storage error") {
			assert.Equal(t, storageErrors.StorageIsClosed, closedErr.Code, "Error code should be StorageIsClosed")
		}

		// 8. Test LastVersionID on closed db
		_, err = store.LastVersionID(appID)
		require.Error(t, err)
		if assert.True(t, errors.As(err, &closedErr), "Error should be a storage error") {
			assert.Equal(t, storageErrors.StorageIsClosed, closedErr.Code, "Error code should be StorageIsClosed")
		}
	})

	t.Run("StoreTwice", func(t *testing.T) {
		// Verifies that attempting to store the same version ID twice results in
		// a 'VersionAlreadyExists' error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "store-twice-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })

		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)

		state := common.ApplicationState{
			ApplicationID:  common.NewApplicationId(801267),
			StateRoot:      sha256.Sum256([]byte("root")),
			EncryptedState: []byte("state"),
		}
		versionID := versionedDb.GenerateVersionID_ForTest([]byte(state.ApplicationID.String()), state.StateRoot[:])
		err = store.Store(ctx, versionID, &state)
		require.NoError(t, err)

		err = store.Store(ctx, versionID, &state)
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.VersionAlreadyExists, storageErr.Code)
	})

	t.Run("GetApplicationStateWithCorruptedData", func(t *testing.T) {
		// Verifies that attempting to retrieve an application state that has been
		// manually corrupted (e.g., invalid JSON stored directly via the adapter)
		// results in an unmarshal error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "corrupted-data-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })

		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		appID := common.NewApplicationId(91)

		// Insert corrupted JSON directly via the adapter, bypassing Store validation.
		corruptedKey := []byte(versionedDb.TestAppStatePrefix + appID.String())
		adapter := store.GetAdapter_ForTest()
		corruptedData := []storage.KeyValuePair{{Key: corruptedKey, Value: []byte("not-valid-json")}}
		versionID := versionedDb.GenerateVersionID_ForTest([]byte(appID.String()), []byte("corrupted-json"))
		err = adapter.Update(uint64(appID), versionID, corruptedData, nil)
		require.NoError(t, err, "Inserting corrupted data via adapter should not fail")

		// Now, attempt to read the data as an ApplicationState.
		_, err = store.GetApplicationState(ctx, appID)
		require.Error(t, err, "Expected an error when getting corrupted application state")
		assert.Contains(t, err.Error(), "failed to unmarshal application state")
	})

	t.Run("NoKeyCollisionWithSameID", func(t *testing.T) {
		// Verifies that storing an ApplicationState and WASM bytecode with the same ID does not result in a key collision
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "collision-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)

		sharedID := common.NewApplicationId(5477)
		sharedIdStr := sharedID.String()
		expectedState := common.ApplicationState{
			ApplicationID:  sharedID,
			StateRoot:      sha256.Sum256([]byte("state-root")),
			EncryptedState: []byte{0x01, 0x02, 0x03},
		}
		expectedBytecode := []byte{0xDE, 0xAD, 0xBE, 0xEF}

		versionID1 := versionedDb.GenerateVersionID_ForTest([]byte(sharedIdStr), expectedState.StateRoot[:])
		err = store.Store(ctx, versionID1, &expectedState)
		require.NoError(t, err, "Store should not return an error")
		// Store again with WASM to verify no key collision between state and wasm prefixes.
		versionID2 := versionedDb.GenerateVersionID_ForTest([]byte(sharedIdStr), expectedBytecode)
		err = store.StoreWithWasm(ctx, versionID2, &expectedState, &common.WASMData{ApplicationID: sharedID, Bytecode: expectedBytecode})
		require.NoError(t, err, "StoreWithWasm should not return an error")

		actualState, err := store.GetApplicationState(ctx, sharedID)
		require.NoError(t, err, "GetApplicationState should not return an error")
		if diff := cmp.Diff(&expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}

		actualBytecode, err := store.GetWASMBytecode(ctx, sharedID)
		require.NoError(t, err, "GetWASMBytecode should not return an error")
		assert.True(t, bytes.Equal(expectedBytecode, actualBytecode), "Retrieved WASM bytecode mismatch")
	})

	t.Run("StoreNilStateShouldFail", func(t *testing.T) {
		// Verifies that passing a nil state results in an error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "nil-entry-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)

		err = store.Store(context.Background(), []byte("v1"), nil)
		require.Error(t, err, "expected error when state is nil")
		assert.Contains(t, err.Error(), "state must not be nil")
	})

	t.Run("ConcurrentReadWrite", func(t *testing.T) {
		// Tests the behavior of the data layer under concurrent read and write
		// operations to ensure thread safety.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "concurrency-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 20)

		initial_app_id := common.NewApplicationId(890)
		// Pre-populate with some data
		initialState := common.ApplicationState{
			ApplicationID: initial_app_id,
			StateRoot:     sha256.Sum256([]byte("initial-root")),
		}
		versionID := versionedDb.GenerateVersionID_ForTest([]byte(initialState.ApplicationID.String()), initialState.StateRoot[:])
		err = store.Store(ctx, versionID, &initialState)
		require.NoError(t, err)

		var wg sync.WaitGroup
		numGoroutines := 10

		// Writer goroutines
		for i := range numGoroutines {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				appID := common.NewApplicationId(uint64(i))
				state := common.ApplicationState{
					ApplicationID: appID,
					StateRoot:     sha256.Sum256([]byte(fmt.Sprintf("root-%d", i))),
				}
				versionID := versionedDb.GenerateVersionID_ForTest([]byte(appID.String()), state.StateRoot[:])
				err := store.Store(ctx, versionID, &state)
				if err != nil {
					var storageErr *storageErrors.Error
					if errors.As(err, &storageErr) && storageErr.Code == storageErrors.VersionAlreadyExists {
						// This is acceptable in a concurrent test
					} else {
						assert.NoError(t, err)
					}
				}
			}(i)
		}

		// Reader goroutines
		for range numGoroutines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Read initial state
				_, err := store.GetApplicationState(ctx, initial_app_id)
				assert.NoError(t, err)

				// Read one of the new states (might not exist yet)
				appID := common.NewApplicationId(5)
				_, _ = store.GetApplicationState(ctx, appID)
			}()
		}

		wg.Wait()

		// Verify all data is present
		for i := range numGoroutines {
			appID := common.NewApplicationId(uint64(i))
			_, err := store.GetApplicationState(ctx, appID)
			assert.NoError(t, err, "should be able to get state for %s", appID)
		}
	})

	t.Run("StoreWithMixedDataTypes", func(t *testing.T) {
		// Verifies that a single Store call can atomically save both ApplicationState
		// and WASMData for the same application under one version.
		// Per-app versioning requires all items in a single Store call to share the same ApplicationID.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "mixed-types-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		sharedAppID := common.NewApplicationId(653485)

		state := &common.ApplicationState{
			ApplicationID: sharedAppID,
			StateRoot:     sha256.Sum256([]byte("state-root")),
		}
		wasm := &common.WASMData{
			ApplicationID: sharedAppID,
			Bytecode:      []byte{0xCA, 0xFE},
		}
		versionID := versionedDb.GenerateVersionID_ForTest(state.StateRoot[:], wasm.Bytecode)

		err = store.StoreWithWasm(ctx, versionID, state, wasm)
		require.NoError(t, err)

		// Verify both were stored correctly
		retrievedState, err := store.GetApplicationState(ctx, sharedAppID)
		require.NoError(t, err)
		if diff := cmp.Diff(state, retrievedState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}

		retrievedWasm, err := store.GetWASMBytecode(ctx, sharedAppID)
		require.NoError(t, err)
		assert.Equal(t, wasm.Bytecode, retrievedWasm)
	})

	t.Run("RollbackAffectsAllDataTypes", func(t *testing.T) {
		// Verifies that a rollback correctly reverts changes across all versioned
		// data types (ApplicationState and WASMData) for a single app,
		// and that per-app rollback does not affect other apps.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "rollback-mixed-types-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)

		appAId := common.NewApplicationId(6)
		appBId := common.NewApplicationId(7)

		// Version 1: Initial state for AppA (state only)
		appAStateV1 := &common.ApplicationState{ApplicationID: appAId, StateRoot: sha256.Sum256([]byte("v1"))}
		v1A := versionedDb.GenerateVersionID_ForTest([]byte(appAId.String()), appAStateV1.StateRoot[:])
		err = store.Store(ctx, v1A, appAStateV1)
		require.NoError(t, err)

		// Version 2 for AppA: update state and add WASM (same app)
		appAStateV2 := &common.ApplicationState{ApplicationID: appAId, StateRoot: sha256.Sum256([]byte("v2"))}
		appAWasm := &common.WASMData{ApplicationID: appAId, Bytecode: []byte("wasm-a")}
		v2A := versionedDb.GenerateVersionID_ForTest(appAStateV2.StateRoot[:], appAWasm.Bytecode)
		err = store.StoreWithWasm(ctx, v2A, appAStateV2, appAWasm)
		require.NoError(t, err)

		// Separate store for AppB (independent version chain)
		appBState := &common.ApplicationState{ApplicationID: appBId, StateRoot: sha256.Sum256([]byte("b-v1"))}
		appBWasm := &common.WASMData{ApplicationID: appBId, Bytecode: []byte("wasm-b")}
		v1B := versionedDb.GenerateVersionID_ForTest([]byte(appBId.String()), appBWasm.Bytecode)
		err = store.StoreWithWasm(ctx, v1B, appBState, appBWasm)
		require.NoError(t, err)

		// Verify V2 state for AppA
		state, err := store.GetApplicationState(ctx, appAId)
		require.NoError(t, err)
		assert.Equal(t, appAStateV2.StateRoot, state.StateRoot)
		wasmA, err := store.GetWASMBytecode(ctx, appAId)
		require.NoError(t, err)
		assert.Equal(t, appAWasm.Bytecode, wasmA)

		// Verify AppB WASM is stored
		wasmB, err := store.GetWASMBytecode(ctx, appBId)
		require.NoError(t, err)
		assert.Equal(t, appBWasm.Bytecode, wasmB)

		// Rollback AppA to V1 — should not affect AppB
		err = store.Rollback(appAId, v1A)
		require.NoError(t, err)

		// Verify AppA state reverted to V1
		state, err = store.GetApplicationState(ctx, appAId)
		require.NoError(t, err)
		assert.Equal(t, appAStateV1.StateRoot, state.StateRoot, "AppA state should be reverted to V1")

		// Verify AppA WASM was also reverted (gone after rollback to v1 which had no WASM)
		_, err = store.GetWASMBytecode(ctx, appAId)
		require.Error(t, err, "WASM for AppA should be gone after rollback to v1")
		var notFoundErr *storageErrors.Error
		require.True(t, errors.As(err, &notFoundErr) && notFoundErr.Code == storageErrors.NotFound)

		// Verify AppB is unaffected by AppA's rollback
		wasmB, err = store.GetWASMBytecode(ctx, appBId)
		require.NoError(t, err, "AppB WASM should be unaffected by AppA rollback")
		assert.Equal(t, appBWasm.Bytecode, wasmB)
	})

	t.Run("StoreAndGetEnclaveKeySetRecovery", func(t *testing.T) {
		// Tests the store and get operations for EnclaveKeySetRecovery.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		expectedRecoveryData := &common.EnclaveKeySetRecovery{
			RecoveryType:       common.RecoveryTypeKMS,
			KeySetCiphertext:   []byte{0x01, 0x02, 0x03},
			RecoveryCiphertext: []byte{0x04, 0x05, 0x06},
		}
		err = store.StoreEnclaveKeySetRecovery(ctx, expectedRecoveryData)
		require.NoError(t, err, "StoreEnclaveKeySetRecovery should not return an error")
		actualRecoveryData, err := store.GetEnclaveKeySetRecovery(ctx)
		require.NoError(t, err, "GetEnclaveKeySetRecovery should not return an error")
		require.NotNil(t, actualRecoveryData, "GetEnclaveKeySetRecovery should return non-nil data")
		if diff := cmp.Diff(expectedRecoveryData, actualRecoveryData); diff != "" {
			t.Errorf("Retrieved EnclaveKeySetRecovery mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentEnclaveKeySetRecovery", func(t *testing.T) {
		// Ensures that trying to get non-existent recovery data returns a 'NotFound' error.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "versioned-leveldb-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)
		_, err = store.GetEnclaveKeySetRecovery(ctx)
		require.Error(t, err, "Expected an error when getting non-existent recovery data")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("ConcurrentRollbackAndWrite", func(t *testing.T) {
		// Tests for race conditions between writing new versions and rolling back to old ones.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "concurrent-rollback-write-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 20)

		// Pre-populate with an initial version
		initialVersionID := versionedDb.GenerateVersionID_ForTest([]byte("initial"), []byte("state"))
		err = store.Store(ctx, initialVersionID, &common.ApplicationState{ApplicationID: common.NewApplicationId(1)})
		require.NoError(t, err)

		var wg sync.WaitGroup
		const writerRoutines = 5
		const rollbackRoutines = 5

		// Writer goroutines
		for i := 0; i < writerRoutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				appID := common.NewApplicationId(uint64(i))
				state := common.ApplicationState{
					ApplicationID: appID,
					StateRoot:     sha256.Sum256([]byte(fmt.Sprintf("root-%d", i))),
				}
				versionID := versionedDb.GenerateVersionID_ForTest([]byte(appID.String()), state.StateRoot[:])
				// Errors are possible and acceptable (e.g., version already exists)
				_ = store.Store(ctx, versionID, &state)
			}(i)
		}

		// Rollback goroutines
		for i := 0; i < rollbackRoutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Attempt to rollback to the initial, stable version.
				// This may fail if other rollbacks are happening, which is acceptable.
				_ = store.Rollback(common.NewApplicationId(1), initialVersionID)
			}()
		}

		wg.Wait()

		// The final state is non-deterministic, but the test should complete without panicking.
		// We can check that the database is still in a valid state.
		_, err = store.ListVersions(common.NewApplicationId(1))
		assert.NoError(t, err, "ListVersions should not fail after concurrent operations")
	})

	t.Run("MultiApp_InterleavedUpdates", func(t *testing.T) {
		// Verifies that interleaved Store calls for different apps maintain
		// independent version chains and do not corrupt each other's state.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "multi-app-interleaved-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 10)

		app1 := common.NewApplicationId(100)
		app2 := common.NewApplicationId(200)
		app3 := common.NewApplicationId(300)

		// Interleave stores: app1-v1, app2-v1, app1-v2, app3-v1, app2-v2
		app1StateV1 := &common.ApplicationState{ApplicationID: app1, StateRoot: sha256.Sum256([]byte("app1-v1"))}
		v1app1 := versionedDb.GenerateVersionID_ForTest([]byte("app1"), app1StateV1.StateRoot[:])
		err = store.Store(ctx, v1app1, app1StateV1)
		require.NoError(t, err)

		app2StateV1 := &common.ApplicationState{ApplicationID: app2, StateRoot: sha256.Sum256([]byte("app2-v1"))}
		v1app2 := versionedDb.GenerateVersionID_ForTest([]byte("app2"), app2StateV1.StateRoot[:])
		err = store.Store(ctx, v1app2, app2StateV1)
		require.NoError(t, err)

		app1StateV2 := &common.ApplicationState{ApplicationID: app1, StateRoot: sha256.Sum256([]byte("app1-v2"))}
		v2app1 := versionedDb.GenerateVersionID_ForTest([]byte("app1-v2"), app1StateV2.StateRoot[:])
		err = store.Store(ctx, v2app1, app1StateV2)
		require.NoError(t, err)

		app3StateV1 := &common.ApplicationState{ApplicationID: app3, StateRoot: sha256.Sum256([]byte("app3-v1"))}
		v1app3 := versionedDb.GenerateVersionID_ForTest([]byte("app3"), app3StateV1.StateRoot[:])
		err = store.Store(ctx, v1app3, app3StateV1)
		require.NoError(t, err)

		app2StateV2 := &common.ApplicationState{ApplicationID: app2, StateRoot: sha256.Sum256([]byte("app2-v2"))}
		v2app2 := versionedDb.GenerateVersionID_ForTest([]byte("app2-v2"), app2StateV2.StateRoot[:])
		err = store.Store(ctx, v2app2, app2StateV2)
		require.NoError(t, err)

		// Verify each app has independent version chains
		versions1, err := store.ListVersions(app1)
		require.NoError(t, err)
		assert.Len(t, versions1, 2, "app1 should have 2 versions")

		versions2, err := store.ListVersions(app2)
		require.NoError(t, err)
		assert.Len(t, versions2, 2, "app2 should have 2 versions")

		versions3, err := store.ListVersions(app3)
		require.NoError(t, err)
		assert.Len(t, versions3, 1, "app3 should have 1 version")

		// Verify latest state per app
		s1, err := store.GetApplicationState(ctx, app1)
		require.NoError(t, err)
		assert.Equal(t, app1StateV2.StateRoot, s1.StateRoot)

		s2, err := store.GetApplicationState(ctx, app2)
		require.NoError(t, err)
		assert.Equal(t, app2StateV2.StateRoot, s2.StateRoot)

		s3, err := store.GetApplicationState(ctx, app3)
		require.NoError(t, err)
		assert.Equal(t, app3StateV1.StateRoot, s3.StateRoot)
	})

	t.Run("MultiApp_IndependentRollback", func(t *testing.T) {
		// Rolling back one app must not affect any other app's state or versions.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "multi-app-rollback-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 10)

		app1 := common.NewApplicationId(10)
		app2 := common.NewApplicationId(20)

		// Build 3 versions for app1, 2 versions for app2
		var app1Versions [][]byte
		for i := 1; i <= 3; i++ {
			s := &common.ApplicationState{
				ApplicationID: app1,
				StateRoot:     sha256.Sum256(fmt.Appendf(nil, "app1-v%d", i)),
			}
			vid := versionedDb.GenerateVersionID_ForTest(fmt.Appendf(nil, "app1-vid-%d", i), s.StateRoot[:])
			err = store.Store(ctx, vid, s)
			require.NoError(t, err)
			app1Versions = append(app1Versions, vid)
		}

		var app2Versions [][]byte
		for i := 1; i <= 2; i++ {
			s := &common.ApplicationState{
				ApplicationID: app2,
				StateRoot:     sha256.Sum256(fmt.Appendf(nil, "app2-v%d", i)),
			}
			vid := versionedDb.GenerateVersionID_ForTest(fmt.Appendf(nil, "app2-vid-%d", i), s.StateRoot[:])
			err = store.Store(ctx, vid, s)
			require.NoError(t, err)
			app2Versions = append(app2Versions, vid)
		}

		// Rollback app1 to v1
		err = store.Rollback(app1, app1Versions[0])
		require.NoError(t, err)

		// App1 should have 1 version left
		v1, err := store.ListVersions(app1)
		require.NoError(t, err)
		assert.Len(t, v1, 1)

		lastV1, err := store.LastVersionID(app1)
		require.NoError(t, err)
		assert.Equal(t, app1Versions[0], lastV1)

		// App2 must be completely unaffected
		v2, err := store.ListVersions(app2)
		require.NoError(t, err)
		assert.Len(t, v2, 2, "app2 versions must be unaffected by app1 rollback")

		lastV2, err := store.LastVersionID(app2)
		require.NoError(t, err)
		assert.Equal(t, app2Versions[1], lastV2)

		s2, err := store.GetApplicationState(ctx, app2)
		require.NoError(t, err)
		assert.Equal(t, sha256.Sum256(fmt.Appendf(nil, "app2-v%d", 2)), s2.StateRoot)
	})

	t.Run("MultiApp_VersionPruningPerApp", func(t *testing.T) {
		// Version pruning (maxVersionsToKeep) must operate per-app.
		// App with many updates should be pruned independently.
		const maxVersions = 3
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "multi-app-pruning-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), maxVersions)

		app1 := common.NewApplicationId(50)
		app2 := common.NewApplicationId(60)

		// Push 5 versions for app1 (exceeds maxVersions=3)
		for i := 1; i <= 5; i++ {
			s := &common.ApplicationState{
				ApplicationID: app1,
				StateRoot:     sha256.Sum256(fmt.Appendf(nil, "app1-v%d", i)),
			}
			vid := versionedDb.GenerateVersionID_ForTest(fmt.Appendf(nil, "app1-prune-%d", i), s.StateRoot[:])
			err = store.Store(ctx, vid, s)
			require.NoError(t, err)
		}

		// Push 2 versions for app2 (under maxVersions)
		for i := 1; i <= 2; i++ {
			s := &common.ApplicationState{
				ApplicationID: app2,
				StateRoot:     sha256.Sum256(fmt.Appendf(nil, "app2-v%d", i)),
			}
			vid := versionedDb.GenerateVersionID_ForTest(fmt.Appendf(nil, "app2-prune-%d", i), s.StateRoot[:])
			err = store.Store(ctx, vid, s)
			require.NoError(t, err)
		}

		// App1 should be pruned to maxVersions
		v1, err := store.ListVersions(app1)
		require.NoError(t, err)
		assert.Len(t, v1, maxVersions, "app1 should be pruned to maxVersions")

		// App2 should retain all 2 versions (under limit)
		v2, err := store.ListVersions(app2)
		require.NoError(t, err)
		assert.Len(t, v2, 2, "app2 should retain all versions (under limit)")
	})

	t.Run("MultiApp_StoreRejectsMismatchedAppIDs", func(t *testing.T) {
		// StoreWithWasm must reject state and wasm with different ApplicationIDs.
		tempDir, err := os.MkdirTemp(testVersionedLevelDBBaseDir, "multi-app-reject-test-")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(tempDir) })
		store := createStore(t, filepath.Join(tempDir, "test.db"), 5)

		app1 := common.NewApplicationId(1)
		app2 := common.NewApplicationId(2)

		s1 := &common.ApplicationState{ApplicationID: app1, StateRoot: sha256.Sum256([]byte("s1"))}
		w2 := &common.WASMData{ApplicationID: app2, Bytecode: []byte("wasm")}
		vid := versionedDb.GenerateVersionID_ForTest(s1.StateRoot[:], w2.Bytecode)

		err = store.StoreWithWasm(ctx, vid, s1, w2)
		require.Error(t, err, "StoreWithWasm should reject mismatched ApplicationIDs")
		assert.Contains(t, err.Error(), "inconsistent ApplicationID")

		// Nil state should fail
		err = store.Store(ctx, vid, nil)
		require.Error(t, err, "Store should reject nil state")
		assert.Contains(t, err.Error(), "state must not be nil")

		// Nil wasm in StoreWithWasm should fail
		err = store.StoreWithWasm(ctx, vid, s1, nil)
		require.Error(t, err, "StoreWithWasm should reject nil wasm")
		assert.Contains(t, err.Error(), "wasm must not be nil")
	})
}
