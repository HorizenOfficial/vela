package versioned_leveldb_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	"github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

var testLevelDBVersionedBaseDir string
var testAdapterAppID uint64 = 1

// generateVersionID is a helper function to create a unique version identifier for testing.
func generateVersionID(suffix string) []byte {
	hash := sha256.Sum256([]byte("version-" + suffix + "-" + time.Now().String()))
	return hash[:]
}

// createTestKVPair is a helper function to create a key-value pair for testing.
func createTestKVPair(keySuffix, valueSuffix string) storage.KeyValuePair {
	return storage.KeyValuePair{
		Key:   []byte("test_key_" + keySuffix),
		Value: []byte("test_value_" + valueSuffix),
	}
}

// isStorageErrorWithCode is a helper function to check if an error is a
// specific type of storage error.
func isStorageErrorWithCode(err error, code string) bool {
	var se *storageErrors.Error
	return errors.As(err, &se) && se.Code == code
}

// TestVersionedLevelDbStorageAdapter is the main test suite for the
// VersionedLevelDbStorageAdapter. It covers all aspects of the adapter's
// functionality, including data retrieval, updates, rollbacks, and error handling.
func TestVersionedLevelDbStorageAdapter(t *testing.T) {
	createAdapter := func(t *testing.T, versionsToKeep int) storage.VersionedStorage {
		tempDir, err := os.MkdirTemp(testLevelDBVersionedBaseDir, "adapter-test-")
		require.NoError(t, err, "Failed to create temp directory for adapter DB")

		adapter, err := versioned_leveldb.NewVersionedLevelDbStorageAdapterWithVersions(tempDir, versionsToKeep)
		require.NoError(t, err, "Failed to create VersionedLevelDbStorageAdapter")

		t.Cleanup(func() {
			cleanupErr := adapter.Close()
			if cleanupErr != nil && !errors.Is(cleanupErr, leveldb.ErrClosed) {
				t.Errorf("Cleanup failed: Unexpected error during adapter.Close(): %v", cleanupErr)
			}
			require.NoError(t, os.RemoveAll(tempDir), "Failed to remove adapter temp directory: %s", tempDir)
		})
		return adapter
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("NewVersionedLevelDbStorageAdapter", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(testLevelDBVersionedBaseDir, "adapter-test-default-")
		require.NoError(t, err, "Failed to create temp directory for adapter DB")

		adapter, err := versioned_leveldb.NewVersionedLevelDbStorageAdapter(tempDir)
		require.NoError(t, err, "Failed to create VersionedLevelDbStorageAdapter with default versions")

		t.Cleanup(func() {
			cleanupErr := adapter.Close()
			if cleanupErr != nil && !errors.Is(cleanupErr, leveldb.ErrClosed) {
				t.Errorf("Cleanup failed: Unexpected error during adapter.Close(): %v", cleanupErr)
			}
			require.NoError(t, os.RemoveAll(tempDir), "Failed to remove adapter temp directory: %s", tempDir)
		})

		// Check that we can perform a basic operation
		err = adapter.Update(testAdapterAppID, generateVersionID("default-adapter-test"), []storage.KeyValuePair{createTestKVPair("dat", "val")}, nil)
		require.NoError(t, err, "Update on default adapter should succeed")
		versions, err := adapter.RollbackVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Len(t, versions, 1, "Should have one version after update")
	})

	t.Run("GetOrElseWithError", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		require.NoError(t, adapter.Close())

		defaultValue := []byte("default")
		val := adapter.GetOrElse([]byte("any-key"), defaultValue)
		assert.Equal(t, defaultValue, val)
	})

	t.Run("GetAndGetOrElse", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		testKey := []byte("mykey")
		testValue := []byte("myvalue")
		defaultValue := []byte("default")
		versionID := generateVersionID("get")

		val, err := adapter.Get(testKey)
		assert.NoError(t, err, "Get should not error for non-existent key")
		assert.Nil(t, val, "Expected nil value for non-existent key")

		valOrElse := adapter.GetOrElse(testKey, defaultValue)
		assert.True(t, bytes.Equal(defaultValue, valOrElse), "Expected default value for non-existent key")

		err = adapter.Update(testAdapterAppID, versionID, []storage.KeyValuePair{{Key: testKey, Value: testValue}}, nil)
		require.NoError(t, err, "Update should succeed")

		val, err = adapter.Get(testKey)
		assert.NoError(t, err, "Get should not error for existent key")
		assert.True(t, bytes.Equal(testValue, val), "Expected stored value for existent key")

		valOrElse = adapter.GetOrElse(testKey, defaultValue)
		assert.True(t, bytes.Equal(testValue, valOrElse), "Expected stored value for existent key")
	})

	t.Run("GetBatch", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		kv1 := createTestKVPair("1", "A")
		kv2 := createTestKVPair("2", "B")
		kv3 := createTestKVPair("3", "C")
		versionID1 := generateVersionID("batch1")
		versionID2 := generateVersionID("batch2")

		err := adapter.Update(testAdapterAppID, versionID1, []storage.KeyValuePair{kv1, kv2}, nil)
		require.NoError(t, err, "Update batch 1 should succeed")

		keysToGet := [][]byte{kv1.Key, []byte("non-existent"), kv2.Key, kv3.Key}
		results, err := adapter.GetBatch(keysToGet)
		require.NoError(t, err, "GetBatch should not error")
		require.Len(t, results, len(keysToGet), "Result count should match requested keys count")

		expectedResults := []storage.KeyValuePair{
			{Key: kv1.Key, Value: kv1.Value},
			{Key: []byte("non-existent"), Value: nil},
			{Key: kv2.Key, Value: kv2.Value},
			{Key: kv3.Key, Value: nil},
		}

		if diff := cmp.Diff(expectedResults, results, cmp.Comparer(bytes.Equal)); diff != "" {
			t.Errorf("GetBatch results mismatch (-want +got):\n%s", diff)
		}

		err = adapter.Update(testAdapterAppID, versionID2, []storage.KeyValuePair{kv3}, nil)
		require.NoError(t, err, "Update batch 2 should succeed")
		val, err := adapter.Get(kv3.Key)
		assert.NoError(t, err, "Get for kv3 should succeed")
		assert.True(t, bytes.Equal(kv3.Value, val), "kv3 should be updated")
	})

	t.Run("GetAll", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		kv1 := createTestKVPair("A", "1")
		kv2 := createTestKVPair("B", "2")
		versionID1 := generateVersionID("all1")
		versionID2 := generateVersionID("all2")

		err := adapter.Update(testAdapterAppID, versionID1, []storage.KeyValuePair{kv1}, nil)
		require.NoError(t, err, "Update for kv1 should succeed")
		err = adapter.Update(testAdapterAppID, versionID2, []storage.KeyValuePair{kv2}, nil)
		require.NoError(t, err, "Update for kv2 should succeed")

		allPairs, err := adapter.GetAll()
		require.NoError(t, err, "GetAll should not error")

		expectedPairs := []storage.KeyValuePair{kv1, kv2}
		sortKVPairs(allPairs)
		sortKVPairs(expectedPairs)

		if diff := cmp.Diff(expectedPairs, allPairs, cmp.Comparer(bytes.Equal)); diff != "" {
			t.Errorf("GetAll results mismatch (-want +got):\n%s", diff)
		}

		for _, pair := range allPairs {
			assert.False(t, bytes.Equal(pair.Key, versioned_leveldb.VersionsKey[:]), "GetAll should not return VersionsKey")
		}
	})

	t.Run("LastVersionID", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		initialID, err := adapter.LastVersionID(testAdapterAppID)
		assert.Nil(t, initialID)
		assert.Error(t, err, "Expected error when getting last version from empty storage")

		vID1 := generateVersionID("last1")
		vID2 := generateVersionID("last2")
		vID3 := generateVersionID("last3")

		err = adapter.Update(testAdapterAppID, vID1, []storage.KeyValuePair{createTestKVPair("l1", "1")}, nil)
		require.NoError(t, err, "Update vID1 should succeed")
		lastID, err := adapter.LastVersionID(testAdapterAppID)
		assert.NoError(t, err, "LastVersionID should not error after first update")
		assert.True(t, bytes.Equal(vID1, lastID), "Expected vID1 as last version")

		err = adapter.Update(testAdapterAppID, vID2, []storage.KeyValuePair{createTestKVPair("l2", "2")}, nil)
		require.NoError(t, err, "Update vID2 should succeed")
		lastID, err = adapter.LastVersionID(testAdapterAppID)
		assert.NoError(t, err, "LastVersionID should not error after second update")
		assert.True(t, bytes.Equal(vID2, lastID), "Expected vID2 as last version")

		err = adapter.Update(testAdapterAppID, vID3, []storage.KeyValuePair{createTestKVPair("l3", "3")}, nil)
		require.NoError(t, err, "Update vID3 should succeed")
		lastID, err = adapter.LastVersionID(testAdapterAppID)
		assert.NoError(t, err, "LastVersionID should not error after third update")
		assert.True(t, bytes.Equal(vID3, lastID), "Expected vID3 as last version")
	})

	t.Run("UpdateValidations", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		baseVersion := generateVersionID("base")
		err := adapter.Update(testAdapterAppID, baseVersion, []storage.KeyValuePair{createTestKVPair("v1", "val1")}, nil)
		require.NoError(t, err, "Base update should succeed")

		err = adapter.Update(testAdapterAppID, baseVersion, []storage.KeyValuePair{createTestKVPair("v2", "val2")}, nil)
		assert.Error(t, err, "Expected error when version ID already exists")
		assert.True(t, isStorageErrorWithCode(err, storageErrors.VersionAlreadyExists), "Expected 'version_already_exists' error")

		invalidUpdate := []storage.KeyValuePair{createTestKVPair("dup", "A"), createTestKVPair("dup", "B")}
		err = adapter.Update(testAdapterAppID, generateVersionID("dup_update"), invalidUpdate, nil)
		assert.Error(t, err, "Expected error for duplicate key in toUpdate")
		assert.True(t, isStorageErrorWithCode(err, storageErrors.InvalidParameter), "Expected 'invalid_parameter' error")

		invalidRemove := [][]byte{[]byte("dup_rem"), []byte("dup_rem")}
		err = adapter.Update(testAdapterAppID, generateVersionID("dup_remove"), nil, invalidRemove)
		assert.Error(t, err, "Expected error for duplicate key in toRemove")
		assert.True(t, isStorageErrorWithCode(err, storageErrors.InvalidParameter), "Expected 'invalid_parameter' error")

		verIDAsKey := generateVersionID("id_as_key")
		err = adapter.Update(testAdapterAppID, verIDAsKey, []storage.KeyValuePair{{Key: verIDAsKey, Value: []byte("val")}}, nil)
		assert.Error(t, err, "Expected error when version ID is used as data key in toUpdate")
		assert.True(t, isStorageErrorWithCode(err, storageErrors.InvalidParameter), "Expected 'invalid_parameter' error")

		verIDAsKey = generateVersionID("id_as_key_2")
		err = adapter.Update(testAdapterAppID, verIDAsKey, nil, [][]byte{verIDAsKey})
		assert.Error(t, err, "Expected error when version ID is used as data key in toRemove")
		assert.True(t, isStorageErrorWithCode(err, storageErrors.InvalidParameter), "Expected 'invalid_parameter' error")
	})

	t.Run("UpdateShrinksVersions", func(t *testing.T) {
		adapter := createAdapter(t, 2)
		vIDs := make([][]byte, 4)
		for i := 0; i < 4; i++ {
			vIDs[i] = generateVersionID(fmt.Sprintf("shrink_%d", i))
			err := adapter.Update(testAdapterAppID, vIDs[i], []storage.KeyValuePair{createTestKVPair(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))}, nil)
			require.NoError(t, err, fmt.Sprintf("Update for vID%d should succeed", i))
		}

		currentVersions, err := adapter.RollbackVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Len(t, currentVersions, 2, "Expected only 2 versions to be kept")
		assert.True(t, bytes.Equal(vIDs[3], currentVersions[0]), "Expected newest version to be vID3")
		assert.True(t, bytes.Equal(vIDs[2], currentVersions[1]), "Expected next newest version to be vID2")

		val0, err := adapter.Get(createTestKVPair("k0", "").Key)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(createTestKVPair("k0", "v0").Value, val0), "Data for shrunk version vID0 should still exist")
	})

	t.Run("RollbackSuccess", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		vID1 := generateVersionID("rb1")
		vID2 := generateVersionID("rb2")
		vID3 := generateVersionID("rb3")

		key1 := []byte("shared_key")
		key2 := []byte("key_v2")

		err := adapter.Update(testAdapterAppID, vID1, []storage.KeyValuePair{{Key: key1, Value: []byte("val1")}}, nil)
		require.NoError(t, err, "Update vID1 should succeed")

		err = adapter.Update(testAdapterAppID, vID2, []storage.KeyValuePair{{Key: key1, Value: []byte("val2_updated")}, {Key: key2, Value: []byte("val2_new")}}, nil)
		require.NoError(t, err, "Update vID2 should succeed")

		err = adapter.Update(testAdapterAppID, vID3, []storage.KeyValuePair{{Key: key1, Value: []byte("val3_updated_again")}}, [][]byte{key2})
		require.NoError(t, err, "Update vID3 should succeed")

		err = adapter.Rollback(testAdapterAppID, vID2)
		require.NoError(t, err, "Rollback to vID2 should succeed")

		val, err := adapter.Get(key1)
		require.NoError(t, err)
		assert.True(t, bytes.Equal([]byte("val2_updated"), val), "Key1 should be vID2 value after rollback")
		val, err = adapter.Get(key2)
		require.NoError(t, err)
		assert.True(t, bytes.Equal([]byte("val2_new"), val), "Key2 should exist (restored from vID2) after rollback")

		versions, err := adapter.RollbackVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Len(t, versions, 2, "Expected 2 versions after rollback")
		assert.True(t, bytes.Equal(vID2, versions[0]), "Newest version should be vID2")
		assert.True(t, bytes.Equal(vID1, versions[1]), "Next version should be vID1")
	})

	t.Run("RollbackNonExistentVersion", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		err := adapter.Rollback(testAdapterAppID, generateVersionID("non_exist"))
		assert.Error(t, err, "Expected error when rolling back to non-existent version")
		assert.True(t, isStorageErrorWithCode(err, storageErrors.VersionNotFound), "Expected 'version_not_found' error")
	})

	t.Run("RollbackVersionsAndLimited", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		vIDs := make([][]byte, 5)
		for i := 0; i < 5; i++ {
			vIDs[i] = generateVersionID(fmt.Sprintf("list_%d", i))
			err := adapter.Update(testAdapterAppID, vIDs[i], []storage.KeyValuePair{createTestKVPair(fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))}, nil)
			require.NoError(t, err)
		}

		allVersions, err := adapter.RollbackVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Len(t, allVersions, 5, "Expected 5 versions in total")
		assert.True(t, bytes.Equal(vIDs[4], allVersions[0]), "Expected newest version first")
		assert.True(t, bytes.Equal(vIDs[0], allVersions[4]), "Expected oldest version last")

		limitedVersions, err := adapter.RollbackVersionsLimited(testAdapterAppID, 3)
		require.NoError(t, err)
		assert.Len(t, limitedVersions, 3, "Expected 3 versions for limited list")
		assert.True(t, bytes.Equal(vIDs[4], limitedVersions[0]), "Expected newest version first in limited list")
		assert.True(t, bytes.Equal(vIDs[2], limitedVersions[2]), "Expected 3rd newest version in limited list")

		limitedVersionsZero, err := adapter.RollbackVersionsLimited(testAdapterAppID, 0)
		require.NoError(t, err)
		assert.Len(t, limitedVersionsZero, 0, "Expected 0 versions for limited list with 0 limit")

		limitedVersionsTooMany, err := adapter.RollbackVersionsLimited(testAdapterAppID, 100)
		require.NoError(t, err)
		assert.Len(t, limitedVersionsTooMany, 5, "Expected all versions if limit is too high")
	})

	t.Run("OperationsAfterClose", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		require.NoError(t, adapter.Close(), "Closing the adapter should not return an error on first close")

		err := adapter.Close()
		assert.Error(t, err, "Expected an error when trying to close an already closed adapter")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate adapter is already closed")

		_, err = adapter.Get([]byte("any_key"))
		require.Error(t, err, "Expected an error when getting from a closed adapter")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		err = adapter.Update(testAdapterAppID, generateVersionID("after_close"), []storage.KeyValuePair{createTestKVPair("ac", "val")}, nil)
		require.Error(t, err, "Expected an error when updating a closed adapter")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		err = adapter.Rollback(testAdapterAppID, generateVersionID("ac_rb"))
		require.Error(t, err, "Expected an error when rolling back on a closed adapter")
		assert.Contains(t, err.Error(), "closed", "Error message should indicate store is closed")

		iter := adapter.GetIterator()
		assert.Nil(t, iter.Key(), "Iterator from closed adapter should return nil key initially")
		assert.Error(t, iter.Error(), "Iterator from closed adapter should have an error")
		assert.Contains(t, iter.Error().Error(), "closed", "Iterator error message should indicate store is closed")
		iter.Release()
	})

	t.Run("IsEmptyAndNumberOfVersions", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		empty, err := adapter.IsEmpty(testAdapterAppID)
		require.NoError(t, err)
		assert.True(t, empty, "Expected empty storage initially")

		numVersions, err := adapter.NumberOfVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Zero(t, numVersions, "Expected 0 versions initially")

		vID1 := generateVersionID("empty_check_1")
		err = adapter.Update(testAdapterAppID, vID1, []storage.KeyValuePair{createTestKVPair("e1", "1")}, nil)
		require.NoError(t, err, "Update should succeed")

		empty, err = adapter.IsEmpty(testAdapterAppID)
		require.NoError(t, err)
		assert.False(t, empty, "Expected non-empty storage after update")

		numVersions, err = adapter.NumberOfVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Equal(t, 1, numVersions, "Expected 1 version after update")

		vID2 := generateVersionID("empty_check_2")
		err = adapter.Update(testAdapterAppID, vID2, []storage.KeyValuePair{createTestKVPair("e2", "2")}, nil)
		require.NoError(t, err, "Update should succeed")

		empty, err = adapter.IsEmpty(testAdapterAppID)
		require.NoError(t, err)
		assert.False(t, empty, "Expected non-empty storage after second update")

		numVersions, err = adapter.NumberOfVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Equal(t, 2, numVersions, "Expected 2 versions after second update")

		err = adapter.Rollback(testAdapterAppID, vID1)
		require.NoError(t, err, "Rollback should succeed")

		empty, err = adapter.IsEmpty(testAdapterAppID)
		require.NoError(t, err)
		assert.False(t, empty, "Expected non-empty storage after rollback (to 1 version)")

		numVersions, err = adapter.NumberOfVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Equal(t, 1, numVersions, "Expected 1 version after rollback")
	})

	t.Run("Iterator", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		kv1 := createTestKVPair("it_a", "val_a")
		kv2 := createTestKVPair("it_b", "val_b")
		versionID1 := generateVersionID("iter1")
		versionID2 := generateVersionID("iter2")

		err := adapter.Update(testAdapterAppID, versionID1, []storage.KeyValuePair{kv1}, nil)
		require.NoError(t, err)
		err = adapter.Update(testAdapterAppID, versionID2, []storage.KeyValuePair{kv2}, nil)
		require.NoError(t, err)

		iter := adapter.GetIterator()
		defer iter.Release()

		var retrievedPairs []storage.KeyValuePair
		for iter.Next() {
			retrievedPairs = append(retrievedPairs, storage.KeyValuePair{Key: iter.Key(), Value: iter.Value()})
		}
		require.NoError(t, iter.Error(), "Iterator should not have an error")

		expectedPairs := []storage.KeyValuePair{kv1, kv2}
		sortKVPairs(retrievedPairs)
		sortKVPairs(expectedPairs)

		if diff := cmp.Diff(expectedPairs, retrievedPairs, cmp.Comparer(bytes.Equal)); diff != "" {
			t.Errorf("Iterator retrieved pairs mismatch (-want +got):\n%s", diff)
		}

		for _, pair := range retrievedPairs {
			assert.False(t, bytes.Equal(pair.Key, versioned_leveldb.VersionsKey[:]), "Iterator should not return VersionsKey")
			assert.False(t, bytes.Equal(pair.Key, versionID1), "Iterator should not return versionID1 as a data key")
			assert.False(t, bytes.Equal(pair.Key, versionID2), "Iterator should not return versionID2 as a data key")
		}
	})

	t.Run("RollbackVersionsLimitedWithNegativeLimit", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		vID1 := generateVersionID("neg_limit_1")
		err := adapter.Update(testAdapterAppID, vID1, []storage.KeyValuePair{createTestKVPair("nl1", "1")}, nil)
		require.NoError(t, err)

		limitedVersions, err := adapter.RollbackVersionsLimited(testAdapterAppID, -1)
		require.NoError(t, err)
		assert.Len(t, limitedVersions, 0, "Expected 0 versions for limited list with negative limit")
	})

	t.Run("IsEmptyAfterRollbackToLastVersion", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		empty, err := adapter.IsEmpty(testAdapterAppID)
		require.NoError(t, err)
		assert.True(t, empty, "Expected empty storage initially")

		vID1 := generateVersionID("empty_rb_1")
		err = adapter.Update(testAdapterAppID, vID1, []storage.KeyValuePair{createTestKVPair("er1", "1")}, nil)
		require.NoError(t, err)

		empty, err = adapter.IsEmpty(testAdapterAppID)
		require.NoError(t, err)
		assert.False(t, empty, "Expected non-empty storage after update")

		// This rollback should leave the DB with one version
		err = adapter.Rollback(testAdapterAppID, vID1)
		require.NoError(t, err)

		empty, err = adapter.IsEmpty(testAdapterAppID)
		require.NoError(t, err)
		assert.False(t, empty, "Expected non-empty storage after rolling back to the only version")

		numVersions, err := adapter.NumberOfVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Equal(t, 1, numVersions, "Expected 1 version after rollback")
	})

	t.Run("GetAllWithEmptyDatabase", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		allPairs, err := adapter.GetAll()
		require.NoError(t, err, "GetAll should not error on an empty database")
		assert.Empty(t, allPairs, "GetAll on an empty database should return an empty slice")
	})

	t.Run("UpdateWithEmptyPayload", func(t *testing.T) {
		adapter := createAdapter(t, 10)
		vID := generateVersionID("empty_payload")

		// An update with no changes should still create a new version.
		err := adapter.Update(testAdapterAppID, vID, nil, nil)
		require.NoError(t, err, "Update with nil slices should succeed")

		numVersions, err := adapter.NumberOfVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Equal(t, 1, numVersions, "Should create a version even with empty payload")

		lastID, err := adapter.LastVersionID(testAdapterAppID)
		require.NoError(t, err)
		assert.Equal(t, vID, lastID, "The new version ID should be the last one")

		// A second update with empty slices should also succeed.
		vID2 := generateVersionID("empty_payload_2")
		err = adapter.Update(testAdapterAppID, vID2, []storage.KeyValuePair{}, [][]byte{})
		require.NoError(t, err, "Update with empty slices should succeed")

		numVersions, err = adapter.NumberOfVersions(testAdapterAppID)
		require.NoError(t, err)
		assert.Equal(t, 2, numVersions, "Should create a second version")
	})
}

// sortKVPairs is a helper function to sort a slice of KeyValuePair objects by key.
func sortKVPairs(pairs []storage.KeyValuePair) {
	sort.Slice(pairs, func(i, j int) bool {
		return bytes.Compare(pairs[i].Key, pairs[j].Key) < 0
	})
}
