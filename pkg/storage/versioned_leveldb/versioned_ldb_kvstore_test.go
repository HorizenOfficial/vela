package versioned_leveldb_test

import (
	"crypto/sha256"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	"github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

var testAppID uint64 = 1

// createKVStore is a helper function that sets up a new VersionedLDBKVStore
// for testing. It creates a temporary directory for the LevelDB database and
// returns the store instance and a cleanup function to be called with defer.
func createKVStore(t *testing.T, versionsToKeep int) (*versioned_leveldb.VersionedLDBKVStore, func()) {
	tempDir, err := os.MkdirTemp("", "kvstore-test-")
	require.NoError(t, err)

	db, err := leveldb.OpenFile(filepath.Join(tempDir, "db"), nil)
	require.NoError(t, err)

	kvStore := versioned_leveldb.NewVersionedLDBKVStore(db, versionsToKeep)

	cleanup := func() {
		require.NoError(t, kvStore.Close())
		require.NoError(t, os.RemoveAll(tempDir))
	}

	return kvStore, cleanup
}

// TestVersionedLDBKVStore tests the core functionality of the VersionedLDBKVStore,
// including updates, rollbacks, version management, and error handling.
func TestVersionedLDBKVStore(t *testing.T) {
	t.Run("UpdateWithInvalidVersionIDSize", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		err := kvStore.Update(testAppID, nil, nil, []byte("invalid-version-id"))
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.InvalidParameter, storageErr.Code)
	})

	t.Run("UpdateWhenVersionIDExists", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		versionID := sha256.Sum256([]byte("version1"))
		err := kvStore.Update(testAppID, nil, nil, versionID[:])
		require.NoError(t, err)

		err = kvStore.Update(testAppID, nil, nil, versionID[:])
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.VersionAlreadyExists, storageErr.Code)
	})

	t.Run("RollbackTo", func(t *testing.T) {
		dir := t.TempDir()
		db, err := leveldb.OpenFile(filepath.Join(dir, "testdb"), nil)
		require.NoError(t, err)
		defer db.Close()

		store := versioned_leveldb.NewVersionedLDBKVStore(db, 10)

		// v1
		v1 := sha256.Sum256([]byte("v1"))
		err = store.Update(testAppID, []storage.KeyValuePair{{Key: []byte("k1"), Value: []byte("v1")}}, nil, v1[:])
		require.NoError(t, err)

		// v2
		v2 := sha256.Sum256([]byte("v2"))
		err = store.Update(testAppID, []storage.KeyValuePair{{Key: []byte("k2"), Value: []byte("v2")}}, nil, v2[:])
		require.NoError(t, err)

		// v3
		v3 := sha256.Sum256([]byte("v3"))
		err = store.Update(testAppID, []storage.KeyValuePair{{Key: []byte("k3"), Value: []byte("v3")}}, nil, v3[:])
		require.NoError(t, err)

		// Sanity check: all keys should exist
		val, err := db.Get([]byte("k1"), nil)
		require.NoError(t, err)
		require.Equal(t, []byte("v1"), val)

		val, err = db.Get([]byte("k2"), nil)
		require.NoError(t, err)
		require.Equal(t, []byte("v2"), val)

		val, err = db.Get([]byte("k3"), nil)
		require.NoError(t, err)
		require.Equal(t, []byte("v3"), val)

		// Rollback to v2
		err = store.RollbackTo(testAppID, v2[:])
		require.NoError(t, err)

		// Expected: k1 and k2 present, k3 removed
		_, err = db.Get([]byte("k1"), nil)
		require.NoError(t, err, "expected k1 to remain after rollback to v2")

		_, err = db.Get([]byte("k2"), nil)
		require.NoError(t, err, "expected k2 to remain after rollback to v2")

		_, err = db.Get([]byte("k3"), nil)
		require.Error(t, err, "expected k3 to be removed after rollback to v2")
	})

	t.Run("RollbackToNonExistentVersion", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		versionID := sha256.Sum256([]byte("non-existent-version"))
		err := kvStore.RollbackTo(testAppID, versionID[:])
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.VersionNotFound, storageErr.Code)
	})

	t.Run("VersionsWithCorruptedData", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		// Manually put corrupted data into the per-app versions key
		appVersionsKey := versioned_leveldb.VersionsKeyForApp_ForTest(testAppID)
		err := kvStore.Db.Put(appVersionsKey[:], []byte("corrupted-data"), nil)
		require.NoError(t, err)

		_, err = kvStore.Versions(testAppID)
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.InconsistentState, storageErr.Code)
	})

	t.Run("VersionIDExists", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		versionID := sha256.Sum256([]byte("version1"))
		err := kvStore.Update(testAppID, nil, nil, versionID[:])
		require.NoError(t, err)

		exists, err := kvStore.VersionIDExists(testAppID, versionID[:])
		require.NoError(t, err)
		assert.True(t, exists)

		nonExistentVersionID := sha256.Sum256([]byte("non-existent-version"))
		exists, err = kvStore.VersionIDExists(testAppID, nonExistentVersionID[:])
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("GetIterator", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		versionID := sha256.Sum256([]byte("version1"))
		key := []byte("key1")
		value := []byte("value1")
		err := kvStore.Update(testAppID, []storage.KeyValuePair{{Key: key, Value: value}}, nil, versionID[:])
		require.NoError(t, err)

		iter := kvStore.GetIterator()
		defer iter.Release()

		assert.True(t, iter.Next())
		assert.Equal(t, key, iter.Key())
		assert.Equal(t, value, iter.Value())
		assert.False(t, iter.Next())
		assert.NoError(t, iter.Error())
	})

	t.Run("RollbackToWithInconsistentState", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		vID1 := sha256.Sum256([]byte("version1"))
		vID2 := sha256.Sum256([]byte("version2"))
		vID3 := sha256.Sum256([]byte("version3"))

		err := kvStore.Update(testAppID, nil, nil, vID1[:])
		require.NoError(t, err)
		err = kvStore.Update(testAppID, nil, nil, vID2[:])
		require.NoError(t, err)
		err = kvStore.Update(testAppID, nil, nil, vID3[:])
		require.NoError(t, err)

		// Manually delete the change set for vID2
		err = kvStore.Db.Delete(vID2[:], nil)
		require.NoError(t, err)

		err = kvStore.RollbackTo(testAppID, vID1[:])
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.InconsistentState, storageErr.Code)
	})

	t.Run("RollbackToWithMalformedChangeSet", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		vID1 := sha256.Sum256([]byte("version1"))
		vID2 := sha256.Sum256([]byte("version2"))

		err := kvStore.Update(testAppID, nil, nil, vID1[:])
		require.NoError(t, err)
		err = kvStore.Update(testAppID, nil, nil, vID2[:])
		require.NoError(t, err)

		// Manually put a malformed change set for vID2
		err = kvStore.Db.Put(vID2[:], []byte("malformed-change-set"), nil)
		require.NoError(t, err)

		err = kvStore.RollbackTo(testAppID, vID1[:])
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.InconsistentState, storageErr.Code)
	})

	t.Run("UpdateWithLargeNumberOfKeys", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		const numKeys = 1000
		toUpdate := make([]storage.KeyValuePair, numKeys)
		for i := 0; i < numKeys; i++ {
			toUpdate[i] = storage.KeyValuePair{
				Key:   []byte("key" + strconv.Itoa(i)),
				Value: []byte("value" + strconv.Itoa(i)),
			}
		}

		versionID := sha256.Sum256([]byte("large-update"))
		err := kvStore.Update(testAppID, toUpdate, nil, versionID[:])
		require.NoError(t, err)

		// Verify a few keys
		val, err := kvStore.Db.Get([]byte("key10"), nil)
		require.NoError(t, err)
		assert.Equal(t, []byte("value10"), val)

		val, err = kvStore.Db.Get([]byte("key999"), nil)
		require.NoError(t, err)
		assert.Equal(t, []byte("value999"), val)
	})

	t.Run("GetIteratorWithExcludedKeys", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		vID1 := sha256.Sum256([]byte("version1"))
		vID2 := sha256.Sum256([]byte("version2"))
		key1 := []byte("key1")
		value1 := []byte("value1")
		key2 := []byte("key2")
		value2 := []byte("value2")

		err := kvStore.Update(testAppID, []storage.KeyValuePair{{Key: key1, Value: value1}}, nil, vID1[:])
		require.NoError(t, err)
		err = kvStore.Update(testAppID, []storage.KeyValuePair{{Key: key2, Value: value2}}, nil, vID2[:])
		require.NoError(t, err)

		iter := kvStore.GetIterator()
		defer iter.Release()

		var retrievedKeys [][]byte
		for iter.Next() {
			retrievedKeys = append(retrievedKeys, iter.Key())
		}
		require.NoError(t, iter.Error())

		assert.Len(t, retrievedKeys, 2)
		for _, key := range retrievedKeys {
			assert.NotEqual(t, vID1[:], key)
			assert.NotEqual(t, vID2[:], key)
		}
	})

	t.Run("ErrorIterator", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "kvstore-test-error-iterator-")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		db, err := leveldb.OpenFile(filepath.Join(tempDir, "db"), nil)
		require.NoError(t, err)

		kvStore := versioned_leveldb.NewVersionedLDBKVStore(db, 10)

		// Close the database to force an error when creating an iterator
		require.NoError(t, kvStore.Close())

		iter := kvStore.GetIterator()

		assert.False(t, iter.Next())
		assert.Nil(t, iter.Key())
		assert.Nil(t, iter.Value())
		assert.Error(t, iter.Error())
		iter.Release()
	})

	t.Run("UpdateWithMixedOperationsAndRollback", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		// Pre-existing data
		keyToAlter := []byte("key-to-alter")
		originalValue := []byte("original-value")
		keyToRemove := []byte("key-to-remove")
		valueToRemove := []byte("value-to-remove")

		v1 := sha256.Sum256([]byte("v1"))
		err := kvStore.Update(testAppID, []storage.KeyValuePair{
			{Key: keyToAlter, Value: originalValue},
			{Key: keyToRemove, Value: valueToRemove},
		}, nil, v1[:])
		require.NoError(t, err)

		// The mixed update
		keyToInsert := []byte("key-to-insert")
		valueToInsert := []byte("value-to-insert")
		alteredValue := []byte("altered-value")

		v2 := sha256.Sum256([]byte("v2"))
		err = kvStore.Update(
			testAppID,
			[]storage.KeyValuePair{ // ToUpdate
				{Key: keyToInsert, Value: valueToInsert},
				{Key: keyToAlter, Value: alteredValue},
			},
			[][]byte{keyToRemove}, // ToRemove
			v2[:],
		)
		require.NoError(t, err)

		// Verify state after update
		val, err := kvStore.Db.Get(keyToInsert, nil)
		require.NoError(t, err)
		assert.Equal(t, valueToInsert, val)

		val, err = kvStore.Db.Get(keyToAlter, nil)
		require.NoError(t, err)
		assert.Equal(t, alteredValue, val)

		_, err = kvStore.Db.Get(keyToRemove, nil)
		assert.Equal(t, leveldb.ErrNotFound, err)

		// Rollback the mixed update
		err = kvStore.RollbackTo(testAppID, v1[:])
		require.NoError(t, err)

		// Verify state after rollback
		_, err = kvStore.Db.Get(keyToInsert, nil)
		assert.Equal(t, leveldb.ErrNotFound, err, "keyToInsert should be gone after rollback")

		val, err = kvStore.Db.Get(keyToAlter, nil)
		require.NoError(t, err)
		assert.Equal(t, originalValue, val, "keyToAlter should have its original value")

		val, err = kvStore.Db.Get(keyToRemove, nil)
		require.NoError(t, err)
		assert.Equal(t, valueToRemove, val, "keyToRemove should be restored after rollback")
	})

	t.Run("RollbackMultipleLevels", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		// Create 5 versions
		versions := make([][]byte, 5)
		for i := 0; i < 5; i++ {
			v := sha256.Sum256([]byte("v" + strconv.Itoa(i+1)))
			versions[i] = v[:]
			key := []byte("key" + strconv.Itoa(i+1))
			value := []byte("value" + strconv.Itoa(i+1))
			err := kvStore.Update(testAppID, []storage.KeyValuePair{{Key: key, Value: value}}, nil, versions[i])
			require.NoError(t, err)
		}

		// Rollback from v5 to v2
		err := kvStore.RollbackTo(testAppID, versions[1]) // versions[1] is v2
		require.NoError(t, err)

		// Verify versions list
		currentVersions, err := kvStore.Versions(testAppID)
		require.NoError(t, err)
		require.Len(t, currentVersions, 2)
		assert.Equal(t, versions[1], currentVersions[0]) // Newest should be v2
		assert.Equal(t, versions[0], currentVersions[1]) // Oldest should be v1

		// Verify data state
		// Keys from v1 and v2 should exist
		_, err = kvStore.Db.Get([]byte("key1"), nil)
		require.NoError(t, err)
		_, err = kvStore.Db.Get([]byte("key2"), nil)
		require.NoError(t, err)

		// Keys from v3, v4, v5 should not exist
		_, err = kvStore.Db.Get([]byte("key3"), nil)
		assert.Equal(t, leveldb.ErrNotFound, err)
		_, err = kvStore.Db.Get([]byte("key4"), nil)
		assert.Equal(t, leveldb.ErrNotFound, err)
		_, err = kvStore.Db.Get([]byte("key5"), nil)
		assert.Equal(t, leveldb.ErrNotFound, err)
	})

	t.Run("VersionPruningOnUpdate", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 3) // Keep only 3 versions
		defer cleanup()

		// Add 5 versions
		versions := make([][]byte, 5)
		for i := 0; i < 5; i++ {
			v := sha256.Sum256([]byte("v" + strconv.Itoa(i)))
			versions[i] = v[:]
			key := []byte("key" + strconv.Itoa(i))
			value := []byte("value" + strconv.Itoa(i))
			err := kvStore.Update(testAppID, []storage.KeyValuePair{{Key: key, Value: value}}, nil, versions[i])
			require.NoError(t, err)

			// Check version count after the 4th and 5th updates
			if i >= 3 {
				currentVersions, err := kvStore.Versions(testAppID)
				require.NoError(t, err)
				assert.Len(t, currentVersions, 3, "Should have pruned to 3 versions")
			}
		}

		// Verify the correct versions were kept (the last 3)
		currentVersions, err := kvStore.Versions(testAppID)
		require.NoError(t, err)
		require.Len(t, currentVersions, 3)
		assert.Equal(t, versions[4], currentVersions[0]) // v4 (newest)
		assert.Equal(t, versions[3], currentVersions[1]) // v3
		assert.Equal(t, versions[2], currentVersions[2]) // v2 (oldest)

		// Verify that data from pruned versions (v0, v1) still exists
		_, err = kvStore.Db.Get([]byte("key0"), nil)
		require.NoError(t, err, "Data from pruned version v0 should still exist")
		_, err = kvStore.Db.Get([]byte("key1"), nil)
		require.NoError(t, err, "Data from pruned version v1 should still exist")
	})

	t.Run("MultiApp_IndependentVersionChains", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		app1 := uint64(1)
		app2 := uint64(2)

		// Store data under app1
		v1 := sha256.Sum256([]byte("app1-v1"))
		err := kvStore.Update(app1, []storage.KeyValuePair{{Key: []byte("a1_key"), Value: []byte("a1_val")}}, nil, v1[:])
		require.NoError(t, err)

		// Store data under app2
		v2 := sha256.Sum256([]byte("app2-v1"))
		err = kvStore.Update(app2, []storage.KeyValuePair{{Key: []byte("a2_key"), Value: []byte("a2_val")}}, nil, v2[:])
		require.NoError(t, err)

		// Each app should have exactly 1 version
		versions1, err := kvStore.Versions(app1)
		require.NoError(t, err)
		assert.Len(t, versions1, 1)

		versions2, err := kvStore.Versions(app2)
		require.NoError(t, err)
		assert.Len(t, versions2, 1)

		// Rolling back app1 should not affect app2
		err = kvStore.RollbackTo(app1, v1[:])
		require.NoError(t, err)

		versions2After, err := kvStore.Versions(app2)
		require.NoError(t, err)
		assert.Len(t, versions2After, 1, "app2 versions must be unaffected by app1 rollback")

		// app2 data should still be accessible
		val, err := kvStore.Db.Get([]byte("a2_key"), nil)
		require.NoError(t, err)
		assert.Equal(t, []byte("a2_val"), val)
	})
}

// TestChangeSetSerializer tests the serialization and deserialization of ChangeSet objects.
func TestChangeSetSerializer(t *testing.T) {
	t.Run("ToBytesAndParseBytes", func(t *testing.T) {
		serializer := versioned_leveldb.ChangeSetSerializer{}
		cs := versioned_leveldb.ChangeSet{
			InsertedKeys: [][]byte{[]byte("key1")},
			Removed:      []storage.KeyValuePair{{Key: []byte("key2"), Value: []byte("value2")}},
			Altered:      []storage.KeyValuePair{{Key: []byte("key3"), Value: []byte("value3")}},
		}

		data, err := serializer.ToBytes(cs)
		require.NoError(t, err)

		parsedCs, err := serializer.ParseBytes(data)
		require.NoError(t, err)
		assert.Equal(t, cs, parsedCs)
	})

	t.Run("ParseBytesWithInvalidData", func(t *testing.T) {
		serializer := versioned_leveldb.ChangeSetSerializer{}
		_, err := serializer.ParseBytes([]byte("invalid-json"))
		require.Error(t, err)
	})
}
