package versioned_leveldb_test // Use a separate test package

import (
	"bytes" // For comparing []byte slices
	"context"
	"crypto/sha256"    // For generating unique version IDs
		errors "errors" // For errors.As
	"fmt"
	"os" // For file system operations (temp directories)
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	"github.com/horizen-pes/pkg/storage/versioned_leveldb"

	// For joining file paths
	"testing"
	"time" // For time.Now() in version ID generation

	"github.com/google/go-cmp/cmp"        // For deep comparison of structs
	"github.com/stretchr/testify/assert"  // For general assertions (Error, NoError, Nil, NotNil)
	"github.com/stretchr/testify/require" // For assertions that stop the test immediately on failure
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/horizen-pes/pkg/common" // Your common types and interface
)

// Global base directory for all temporary LevelDB test files.
var testLevelDBVersionedBaseDir string

// TestMain runs before any tests in the package. Used for global setup/teardown.
func TestMain(m *testing.M) {
	fmt.Println("Running TestMain setup for VersionedLevelDbStorageAdapter integration tests...")

	var err error
	// Create a unique temporary base directory for all test DBs.
	testLevelDBVersionedBaseDir, err = os.MkdirTemp("", "versioned_leveldb_test_dbs_")
	if err != nil {
		fmt.Printf("Failed to create base directory for versioned LevelDB tests: %v\n", err)
		os.Exit(1)
	}

	// Run all tests.
	code := m.Run()

	// Clean up the base directory after all tests are done.
	fmt.Println("Running TestMain teardown for VersionedLevelDbStorageAdapter integration tests...")
	err = os.RemoveAll(testLevelDBVersionedBaseDir)
	if err != nil {
		fmt.Printf("Failed to clean up base directory %s: %v\n", testLevelDBVersionedBaseDir, err)
	}

	os.Exit(code) // Exit with the test result code.
}

// Helper function to generate a unique version ID for tests.
// Ensures the ID has the correct length (ConstantsHashLength).
func generateVersionID(suffix string) []byte {
	hash := sha256.Sum256([]byte("version-" + suffix + "-" + time.Now().String()))
	return hash[:] // Return a slice from the array
}

// Helper to create a KeyValuePair for testing.
func createTestKVPair(keySuffix, valueSuffix string) common.KeyValuePair {
	return common.KeyValuePair{
		Key:   []byte("test_key_" + keySuffix),
		Value: []byte("test_value_" + valueSuffix),
	}
}

// Helper to check if an error is a specific storage.Error with a given code.
func isStorageErrorWithCode(err error, code string) bool {
	var se *storageErrors.Error
	return errors.As(err, &se) && se.Code == code
}

// TestVersionedLevelDbStorageAdapter is the main test suite for the adapter.
func TestVersionedLevelDbStorageAdapter(t *testing.T) {
	// createAdapter is a factory function to get a new, isolated adapter instance for each subtest.
	createAdapter := func() common.VersionedStorage {
		// Create a unique temporary directory for each adapter instance's DB file.
		tempDir, err := os.MkdirTemp(testLevelDBVersionedBaseDir, "adapter-test-")
		require.NoError(t, err, "Failed to create temp directory for adapter DB")

		// Create the adapter with default versions to keep.
		adapter, err := versioned_leveldb.NewVersionedLevelDbStorageAdapter(tempDir)
		require.NoError(t, err, "Failed to create VersionedLevelDbStorageAdapter")

		// Ensure adapter is closed and its temp directory removed after the test.
		t.Cleanup(func() {
			//			require.NoError(t, adapter.Close(), "Failed to close adapter during cleanup")
			cleanupErr := adapter.Close() // Get the error from Close()

			// If cleanupErr is not nil, check if it's the specific LevelDB "closed" error
			if cleanupErr != nil {
				// Use errors.Is to check if the error chain contains LevelDB's ErrClosed
				if errors.Is(cleanupErr, leveldb.ErrClosed) {
					// This is the expected "already closed" error on subsequent closes from LevelDB.
					// We explicitly ignore it, as it's not a true cleanup failure.
					assert.True(t, errors.Is(cleanupErr, leveldb.ErrClosed), "Cleanup: Expected LevelDB ErrClosed on double close")
				} else {
					// If it's any other error, then it's an unexpected failure in cleanup.
					// Fail the test explicitly with an error.
					t.Errorf("Cleanup failed: Unexpected error during adapter.Close(): %v", cleanupErr)
				}
			}
			require.NoError(t, os.RemoveAll(tempDir), "Failed to remove adapter temp directory: %s", tempDir)
		})
		return adapter
	}

	// Context for operations.
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Test basic Get and GetOrElse ---
	t.Run("GetAndGetOrElse", func(t *testing.T) {
		adapter := createAdapter()
		testKey := []byte("mykey")
		testValue := []byte("myvalue")
		defaultValue := []byte("default")
		versionID := generateVersionID("get")

		// Initial state: Key should not exist
		val, err := adapter.Get(testKey)
		assert.NoError(t, err, "Get should not error for non-existent key")
		assert.Nil(t, val, "Expected nil value for non-existent key")

		valOrElse := adapter.GetOrElse(testKey, defaultValue)
		assert.True(t, bytes.Equal(defaultValue, valOrElse), "Expected default value for non-existent key")

		// Store a value using update (as per versioned storage logic)
		err = adapter.Update(versionID, []common.KeyValuePair{{Key: testKey, Value: testValue}}, nil)
		require.NoError(t, err, "Update should succeed")

		// Key should now exist
		val, err = adapter.Get(testKey)
		assert.NoError(t, err, "Get should not error for existent key")
		assert.True(t, bytes.Equal(testValue, val), "Expected stored value for existent key")

		valOrElse = adapter.GetOrElse(testKey, defaultValue)
		assert.True(t, bytes.Equal(testValue, valOrElse), "Expected stored value for existent key")
	})

	// --- Test GetBatch ---
	t.Run("GetBatch", func(t *testing.T) {
		adapter := createAdapter()
		kv1 := createTestKVPair("1", "A")
		kv2 := createTestKVPair("2", "B")
		kv3 := createTestKVPair("3", "C")
		versionID1 := generateVersionID("batch1")
		versionID2 := generateVersionID("batch2")

		// Store some items
		err := adapter.Update(versionID1, []common.KeyValuePair{kv1, kv2}, nil)
		require.NoError(t, err, "Update batch 1 should succeed")

		// Get a mix of existing and non-existent keys
		keysToGet := [][]byte{kv1.Key, []byte("non-existent"), kv2.Key, kv3.Key}
		results, err := adapter.GetBatch(keysToGet)
		require.NoError(t, err, "GetBatch should not error")
		require.Len(t, results, len(keysToGet), "Result count should match requested keys count")

		// Expected results (value will be nil for non-existent)
		expectedResults := []common.KeyValuePair{
			{Key: kv1.Key, Value: kv1.Value},
			{Key: []byte("non-existent"), Value: nil},
			{Key: kv2.Key, Value: kv2.Value},
			{Key: kv3.Key, Value: nil},
		}

		if diff := cmp.Diff(expectedResults, results, cmp.Comparer(bytes.Equal)); diff != "" {
			t.Errorf("GetBatch results mismatch (-want +got):\n%s", diff)
		}

		// Update kv3
		err = adapter.Update(versionID2, []common.KeyValuePair{kv3}, nil)
		require.NoError(t, err, "Update batch 2 should succeed")
		val, err := adapter.Get(kv3.Key)
		assert.NoError(t, err, "Get for kv3 should succeed")
		assert.True(t, bytes.Equal(kv3.Value, val), "kv3 should be updated")
	})

	// --- Test GetAll ---
	t.Run("GetAll", func(t *testing.T) {
		adapter := createAdapter()
		kv1 := createTestKVPair("A", "1")
		kv2 := createTestKVPair("B", "2")
		versionID1 := generateVersionID("all1")
		versionID2 := generateVersionID("all2")

		err := adapter.Update(versionID1, []common.KeyValuePair{kv1}, nil)
		require.NoError(t, err, "Update for kv1 should succeed")
		err = adapter.Update(versionID2, []common.KeyValuePair{kv2}, nil)
		require.NoError(t, err, "Update for kv2 should succeed")

		// GetAll should return only the data keys, excluding versions metadata
		allPairs, err := adapter.GetAll()
		require.NoError(t, err, "GetAll should not error")

		expectedPairs := []common.KeyValuePair{kv1, kv2} // Order might vary based on LevelDB's internal key sort
		// Sort both slices to ensure order-independent comparison for GetAll
		sortKVPairs(allPairs)
		sortKVPairs(expectedPairs)

		if diff := cmp.Diff(expectedPairs, allPairs, cmp.Comparer(bytes.Equal)); diff != "" {
			t.Errorf("GetAll results mismatch (-want +got):\n%s", diff)
		}

		// Verify metadata keys are not included
		for _, pair := range allPairs {
			assert.False(t, bytes.Equal(pair.Key, versioned_leveldb.VersionsKey[:]), "GetAll should not return VersionsKey")
			// We don't check for specific versionIDs as data keys because they are stored under their hash
			// but only if they are the actual change set key.
			// The adapter should correctly filter out keys that represent change sets for versions.
		}
	})

	// --- Test LastVersionID ---
	t.Run("LastVersionID", func(t *testing.T) {
		adapter := createAdapter()
		initialID, err := adapter.LastVersionID()
		assert.True(t, initialID == nil)
		assert.Error(t, err, "Expected error when getting last version from empty storage")

		vID1 := generateVersionID("last1")
		vID2 := generateVersionID("last2")
		vID3 := generateVersionID("last3")

		err = adapter.Update(vID1, []common.KeyValuePair{createTestKVPair("l1", "1")}, nil)
		require.NoError(t, err, "Update vID1 should succeed")
		lastID, err := adapter.LastVersionID()
		assert.NoError(t, err, "LastVersionID should not error after first update")
		assert.True(t, bytes.Equal(vID1, lastID), "Expected vID1 as last version")

		err = adapter.Update(vID2, []common.KeyValuePair{createTestKVPair("l2", "2")}, nil)
		require.NoError(t, err, "Update vID2 should succeed")
		lastID, err = adapter.LastVersionID()
		assert.NoError(t, err, "LastVersionID should not error after second update")
		assert.True(t, bytes.Equal(vID2, lastID), "Expected vID2 as last version")

		err = adapter.Update(vID3, []common.KeyValuePair{createTestKVPair("l3", "3")}, nil)
		require.NoError(t, err, "Update vID3 should succeed")
		lastID, err = adapter.LastVersionID()
		assert.NoError(t, err, "LastVersionID should not error after third update")
		assert.True(t, bytes.Equal(vID3, lastID), "Expected vID3 as last version")
	})

	// --- Test Update logic and validation ---
	t.Run("UpdateValidations", func(t *testing.T) {
		adapter := createAdapter()
		baseVersion := generateVersionID("base")
		err := adapter.Update(baseVersion, []common.KeyValuePair{createTestKVPair("v1", "val1")}, nil)
		require.NoError(t, err, "Base update should succeed")

		// Case: Version ID already exists
		err = adapter.Update(baseVersion, []common.KeyValuePair{createTestKVPair("v2", "val2")}, nil)
		assert.Error(t, err, "Expected error when version ID already exists")
		assert.True(t, isStorageErrorWithCode(err, "version_already_exists"), "Expected 'version_already_exists' error")

		// Case: Duplicate key in toUpdate
		invalidUpdate := []common.KeyValuePair{createTestKVPair("dup", "A"), createTestKVPair("dup", "B")}
		err = adapter.Update(generateVersionID("dup_update"), invalidUpdate, nil)
		assert.Error(t, err, "Expected error for duplicate key in toUpdate")
		assert.True(t, isStorageErrorWithCode(err, "invalid_parameter"), "Expected 'invalid_parameter' error")
		assert.Contains(t, err.Error(), "Duplicate key in 'toUpdate'", "Error message should mention duplicate key")

		// Case: Duplicate key in toRemove
		invalidRemove := [][]byte{[]byte("dup_rem"), []byte("dup_rem")}
		err = adapter.Update(generateVersionID("dup_remove"), nil, invalidRemove)
		assert.Error(t, err, "Expected error for duplicate key in toRemove")
		assert.True(t, isStorageErrorWithCode(err, "invalid_parameter"), "Expected 'invalid_parameter' error")
		assert.Contains(t, err.Error(), "Duplicate key in 'toRemove'", "Error message should mention duplicate key")

		// Case: Version ID used as data key
		verIDAsKey := generateVersionID("id_as_key")
		err = adapter.Update(verIDAsKey, []common.KeyValuePair{{Key: verIDAsKey, Value: []byte("val")}}, nil)
		assert.Error(t, err, "Expected error when version ID is used as data key in toUpdate")
		assert.True(t, isStorageErrorWithCode(err, "invalid_parameter"), "Expected 'invalid_parameter' error")
		assert.Contains(t, err.Error(), "Version ID cannot be used as a key in 'toUpdate'", "Error message should mention version ID as key")

		err = adapter.Update(verIDAsKey, nil, [][]byte{verIDAsKey})
		assert.Error(t, err, "Expected error when version ID is used as data key in toRemove")
		assert.True(t, isStorageErrorWithCode(err, "invalid_parameter"), "Expected 'invalid_parameter' error")
		assert.Contains(t, err.Error(), "Version ID cannot be used as a key in 'toRemove'", "Error message should mention version ID as key")
	})

	/*
		t.Run("UpdateShrinksVersions", func(t *testing.T) {
			adapterWithFewVersions := createAdapter()
			// Override versions to keep for this specific test
			adapterWithFewVersionsImpl, ok := adapterWithFewVersions.(*storage.VersionedLevelDbStorageAdapter)
			require.True(t, ok, "Failed to cast adapter to concrete type")
			adapterWithFewVersionsImpl.database.keepVersions = 2 // Keep only 2 versions

			vIDs := make([][]byte, 4)
			for i := 0; i < 4; i++ {
				vIDs[i] = generateVersionID(fmt.Sprintf("shrink_%d", i))
				err := adapterWithFewVersions.Update(vIDs[i], []common.KeyValuePair{createTestKVPair(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))}, nil)
				require.NoError(t, err, fmt.Sprintf("Update for vID%d should succeed", i))
			}

			// After 4 updates, only the last 2 versions should remain
			currentVersions, err := adapterWithFewVersions.RollbackVersions()
			require.NoError(t, err)
			assert.Len(t, currentVersions, 2, "Expected only 2 versions to be kept")
			assert.True(t, bytes.Equal(vIDs[3], currentVersions[0]), "Expected newest version to be vID3")
			assert.True(t, bytes.Equal(vIDs[2], currentVersions[1]), "Expected next newest version to be vID2")

			// The data associated with vID0 and vID1 should still exist in the DB (they are not rolled back, just history shrunk)
			val0, err := adapterWithFewVersions.Get(createTestKVPair("k0", "").Key)
			assert.NoError(t, err)
			assert.True(t, bytes.Equal(createTestKVPair("k0", "v0").Value, val0), "Data for shrunk version vID0 should still exist")
		})
	*/

	// --- Test Rollback ---
	t.Run("RollbackSuccess", func(t *testing.T) {
		adapter := createAdapter()
		vID1 := generateVersionID("rb1")
		vID2 := generateVersionID("rb2")
		vID3 := generateVersionID("rb3") // Target version

		key1 := []byte("shared_key")
		key2 := []byte("key_v2")
		key3 := []byte("key_v3")

		// Initial state (vID1)
		err := adapter.Update(vID1, []common.KeyValuePair{
			{Key: key1, Value: []byte("val1")},
		}, nil)
		require.NoError(t, err, "Update vID1 should succeed")

		// State after vID2 (key1 changed, key2 added, key3 removed)
		err = adapter.Update(vID2, []common.KeyValuePair{
			{Key: key1, Value: []byte("val2_updated")}, // key1 altered
			{Key: key2, Value: []byte("val2_new")},     // key2 inserted
		}, nil)
		require.NoError(t, err, "Update vID2 should succeed")

		// State after vID3 (key1 changed again, key2 removed)
		err = adapter.Update(vID3, []common.KeyValuePair{
			{Key: key1, Value: []byte("val3_updated_again")}, // key1 altered again
		}, [][]byte{key2}) // key2 removed
		require.NoError(t, err, "Update vID3 should succeed")

		// Verify current state before rollback
		val, err := adapter.Get(key1)
		require.NoError(t, err)
		assert.True(t, bytes.Equal([]byte("val3_updated_again"), val), "Key1 should be vID3 value")
		val, err = adapter.Get(key2)
		require.NoError(t, err)
		assert.Nil(t, val, "Key2 should be nil (removed in vID3)")
		val, err = adapter.Get(key3)
		require.NoError(t, err)
		assert.Nil(t, val, "Key3 should be nil (never added)")

		// Perform rollback to vID2
		err = adapter.Rollback(vID2)
		require.NoError(t, err, "Rollback to vID2 should succeed")

		// Verify state after rollback to vID2
		val, err = adapter.Get(key1)
		require.NoError(t, err)
		assert.True(t, bytes.Equal([]byte("val2_updated"), val), "Key1 should be vID2 value after rollback")
		val, err = adapter.Get(key2)
		require.NoError(t, err)
		assert.True(t, bytes.Equal([]byte("val2_new"), val), "Key2 should exist (restored from vID2) after rollback")

		// Verify versions after rollback (vID3 should be gone, vID1 and vID2 should remain)
		versions, err := adapter.RollbackVersions()
		require.NoError(t, err)
		assert.Len(t, versions, 2, "Expected 2 versions after rollback")
		assert.True(t, bytes.Equal(vID2, versions[0]), "Newest version should be vID2")
		assert.True(t, bytes.Equal(vID1, versions[1]), "Next version should be vID1")
	})

	t.Run("RollbackNonExistentVersion", func(t *testing.T) {
		adapter := createAdapter()
		err := adapter.Rollback(generateVersionID("non_exist"))
		assert.Error(t, err, "Expected error when rolling back to non-existent version")
		assert.True(t, isStorageErrorWithCode(err, "version_not_found"), "Expected 'version_not_found' error")
	})

	// --- Test RollbackVersions and RollbackVersionsLimited ---
	t.Run("RollbackVersionsAndLimited", func(t *testing.T) {
		adapter := createAdapter()
		vIDs := make([][]byte, 5)
		for i := 0; i < 5; i++ {
			vIDs[i] = generateVersionID(fmt.Sprintf("list_%d", i))
			err := adapter.Update(vIDs[i], []common.KeyValuePair{createTestKVPair(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))}, nil)
			require.NoError(t, err)
		}

		allVersions, err := adapter.RollbackVersions()
		require.NoError(t, err)
		assert.Len(t, allVersions, 5, "Expected 5 versions in total")
		// Verify order: newest first
		assert.True(t, bytes.Equal(vIDs[4], allVersions[0]), "Expected newest version first")
		assert.True(t, bytes.Equal(vIDs[0], allVersions[4]), "Expected oldest version last")

		// Test RollbackVersionsLimited
		limitedVersions, err := adapter.RollbackVersionsLimited(3)
		require.NoError(t, err)
		assert.Len(t, limitedVersions, 3, "Expected 3 versions for limited list")
		assert.True(t, bytes.Equal(vIDs[4], limitedVersions[0]), "Expected newest version first in limited list")
		assert.True(t, bytes.Equal(vIDs[2], limitedVersions[2]), "Expected 3rd newest version in limited list")

		limitedVersionsZero, err := adapter.RollbackVersionsLimited(0)
		require.NoError(t, err)
		assert.Len(t, limitedVersionsZero, 0, "Expected 0 versions for limited list with 0 limit")

		limitedVersionsTooMany, err := adapter.RollbackVersionsLimited(100)
		require.NoError(t, err)
		assert.Len(t, limitedVersionsTooMany, 5, "Expected all versions if limit is too high")
	})

	// --- Test Close and OperationsAfterClose ---
	t.Run("OperationsAfterClose", func(t *testing.T) {
		adapter := createAdapter() // Get a fresh adapter

		// Close the adapter
		err := adapter.Close()
		require.NoError(t, err, "Closing the adapter should not return an error on first close")

		// Subsequent Close() calls should also return an error for bbolt
		err = adapter.Close()
		assert.Error(t, err, "Expected an error when trying to close an already closed adapter")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate adapter is already closed")

		// Try to perform an operation after closing (Get)
		_, err = adapter.Get([]byte("any_key"))
		require.Error(t, err, "Expected an error when getting from a closed adapter")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		// Try to perform an operation after closing (Update)
		err = adapter.Update(generateVersionID("after_close"), []common.KeyValuePair{createTestKVPair("ac", "val")}, nil)
		require.Error(t, err, "Expected an error when updating a closed adapter")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		// Try to perform an operation after closing (Rollback)
		err = adapter.Rollback(generateVersionID("ac_rb"))
		require.Error(t, err, "Expected an error when rolling back on a closed adapter")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		// Try to perform an operation after closing (GetIterator)
		iter := adapter.GetIterator()
		assert.Nil(t, iter.Key(), "Iterator from closed adapter should return nil key initially")
		assert.Error(t, iter.Error(), "Iterator from closed adapter should have an error")
		assert.Contains(t, iter.Error().Error(), "closed", "Iterator error message should indicate store is closed")
		iter.Release() // Release regardless
	})

	// --- Test IsEmpty and NumberOfVersions ---
	t.Run("IsEmptyAndNumberOfVersions", func(t *testing.T) {
		adapter := createAdapter()

		assert.True(t, adapter.IsEmpty(), "Expected empty storage initially")
		assert.Zero(t, adapter.NumberOfVersions(), "Expected 0 versions initially")

		vID1 := generateVersionID("empty_check_1")
		err := adapter.Update(vID1, []common.KeyValuePair{createTestKVPair("e1", "1")}, nil)
		require.NoError(t, err, "Update should succeed")

		assert.False(t, adapter.IsEmpty(), "Expected non-empty storage after update")
		assert.Equal(t, 1, adapter.NumberOfVersions(), "Expected 1 version after update")

		vID2 := generateVersionID("empty_check_2")
		err = adapter.Update(vID2, []common.KeyValuePair{createTestKVPair("e2", "2")}, nil)
		require.NoError(t, err, "Update should succeed")

		assert.False(t, adapter.IsEmpty(), "Expected non-empty storage after second update")
		assert.Equal(t, 2, adapter.NumberOfVersions(), "Expected 2 versions after second update")

		// Test after rollback that reduces versions (but not to empty)
		err = adapter.Rollback(vID1) // Rollback to vID1, removing vID2
		require.NoError(t, err, "Rollback should succeed")
		assert.False(t, adapter.IsEmpty(), "Expected non-empty storage after rollback (to 1 version)")
		assert.Equal(t, 1, adapter.NumberOfVersions(), "Expected 1 version after rollback")
	})

	// --- Test Iterator ---
	t.Run("Iterator", func(t *testing.T) {
		adapter := createAdapter()
		kv1 := createTestKVPair("it_a", "val_a")
		kv2 := createTestKVPair("it_b", "val_b")
		versionID1 := generateVersionID("iter1")
		versionID2 := generateVersionID("iter2")

		err := adapter.Update(versionID1, []common.KeyValuePair{kv1}, nil)
		require.NoError(t, err)
		err = adapter.Update(versionID2, []common.KeyValuePair{kv2}, nil)
		require.NoError(t, err)

		iter := adapter.GetIterator()
		defer iter.Release() // Always release iterators

		var retrievedPairs []common.KeyValuePair
		for iter.Next() {
			retrievedPairs = append(retrievedPairs, common.KeyValuePair{Key: iter.Key(), Value: iter.Value()})
		}
		require.NoError(t, iter.Error(), "Iterator should not have an error")

		expectedPairs := []common.KeyValuePair{kv1, kv2}
		// The iterator will return keys in sorted order.
		// Ensure your expected and actual lists are sorted for comparison.
		sortKVPairs(retrievedPairs)
		sortKVPairs(expectedPairs)

		if diff := cmp.Diff(expectedPairs, retrievedPairs, cmp.Comparer(bytes.Equal)); diff != "" {
			t.Errorf("Iterator retrieved pairs mismatch (-want +got):\n%s", diff)
		}

		// Ensure metadata keys (VersionsKey, versionID change sets) are NOT returned by iterator
		for _, pair := range retrievedPairs {
			assert.False(t, bytes.Equal(pair.Key, versioned_leveldb.VersionsKey[:]), "Iterator should not return VersionsKey")
			// Add more specific checks if you have other known internal keys
			assert.False(t, bytes.Equal(pair.Key, versionID1), "Iterator should not return versionID1 as a data key")
			assert.False(t, bytes.Equal(pair.Key, versionID2), "Iterator should not return versionID2 as a data key")
		}
	})
}

// Helper function to sort a slice of KeyValuePair by Key.
// Useful for ensuring consistent comparison order in tests where order is not guaranteed by iteration.
func sortKVPairs(pairs []common.KeyValuePair) {
	// Simple bubble sort for demonstration; use sort.Slice for larger datasets
	for i := 0; i < len(pairs)-1; i++ {
		for j := i + 1; j < len(pairs); j++ {
			if bytes.Compare(pairs[i].Key, pairs[j].Key) > 0 {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
}
