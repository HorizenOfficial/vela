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

func TestVersionedLDBKVStore(t *testing.T) {
	t.Run("UpdateWithInvalidVersionIDSize", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		err := kvStore.Update(nil, nil, []byte("invalid-version-id"))
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.InvalidParameter, storageErr.Code)
	})

	t.Run("UpdateWhenVersionIDExists", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		versionID := sha256.Sum256([]byte("version1"))
		err := kvStore.Update(nil, nil, versionID[:])
		require.NoError(t, err)

		err = kvStore.Update(nil, nil, versionID[:])
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.VersionAlreadyExists, storageErr.Code)
	})

	t.Run("RollbackToNonExistentVersion", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		versionID := sha256.Sum256([]byte("non-existent-version"))
		err := kvStore.RollbackTo(versionID[:])
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.VersionNotFound, storageErr.Code)
	})

	t.Run("VersionsWithCorruptedData", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		// Manually put corrupted data into the versions key
		err := kvStore.Db.Put(versioned_leveldb.VersionsKey[:], []byte("corrupted-data"), nil)
		require.NoError(t, err)

		_, err = kvStore.Versions()
		require.Error(t, err)
		var storageErr *storageErrors.Error
		require.True(t, errors.As(err, &storageErr))
		assert.Equal(t, storageErrors.InconsistentState, storageErr.Code)
	})

	t.Run("VersionIDExists", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		versionID := sha256.Sum256([]byte("version1"))
		err := kvStore.Update(nil, nil, versionID[:])
		require.NoError(t, err)

		exists, err := kvStore.VersionIDExists(versionID[:])
		require.NoError(t, err)
		assert.True(t, exists)

		nonExistentVersionID := sha256.Sum256([]byte("non-existent-version"))
		exists, err = kvStore.VersionIDExists(nonExistentVersionID[:])
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("GetIterator", func(t *testing.T) {
		kvStore, cleanup := createKVStore(t, 10)
		defer cleanup()

		versionID := sha256.Sum256([]byte("version1"))
		key := []byte("key1")
		value := []byte("value1")
		err := kvStore.Update([]storage.KeyValuePair{{Key: key, Value: value}}, nil, versionID[:])
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

		err := kvStore.Update(nil, nil, vID1[:])
		require.NoError(t, err)
		err = kvStore.Update(nil, nil, vID2[:])
		require.NoError(t, err)
		err = kvStore.Update(nil, nil, vID3[:])
		require.NoError(t, err)

		// Manually delete the change set for vID2
		err = kvStore.Db.Delete(vID2[:], nil)
		require.NoError(t, err)

		err = kvStore.RollbackTo(vID1[:])
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
		err := kvStore.Update(toUpdate, nil, versionID[:])
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

		err := kvStore.Update([]storage.KeyValuePair{{Key: key1, Value: value1}}, nil, vID1[:])
		require.NoError(t, err)
		err = kvStore.Update([]storage.KeyValuePair{{Key: key2, Value: value2}}, nil, vID2[:])
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
			assert.NotEqual(t, versioned_leveldb.VersionsKey[:], key)
			assert.NotEqual(t, vID1[:], key)
			assert.NotEqual(t, vID2[:], key)
		}
	})
}

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
