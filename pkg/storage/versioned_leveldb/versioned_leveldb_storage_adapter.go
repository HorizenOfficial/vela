package versioned_leveldb

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/horizen-pes/pkg/storage"
)

const (
	DefaultVersionsToKeep = 10 // default number of version to keep for rollback operations. Probably just 1 is enough
)

// It implements the storage.Storage interface on top of VersionedLDBKVStore.
type VersionedLevelDbStorageAdapter struct {
	dataBase    *VersionedLDBKVStore
	versionsKey []byte // Stored as a raw byte slice
}

// NewVersionedLevelDbStorageAdapter creates a new adapter with default versions to keep.
func NewVersionedLevelDbStorageAdapter(pathToDB string) (*VersionedLevelDbStorageAdapter, error) {
	return NewVersionedLevelDbStorageAdapterWithVersions(pathToDB, DefaultVersionsToKeep)
}

// NewVersionedLevelDbStorageAdapterWithVersions creates a new adapter with a specified number of versions to keep.
func NewVersionedLevelDbStorageAdapterWithVersions(pathToDB string, versionsToKeep int) (*VersionedLevelDbStorageAdapter, error) {
	dataBase, err := createDb(pathToDB, versionsToKeep)
	if err != nil {
		return nil, err
	}
	return &VersionedLevelDbStorageAdapter{
		dataBase:    dataBase,
		versionsKey: VersionsKey[:], // Use the global VersionsKey constant
	}, nil
}

// get(key: ByteArrayWrapper): Optional[ByteArrayWrapper] = dataBase.get(key).map(byteArrayToWrapper).asJava
func (s *VersionedLevelDbStorageAdapter) Get(key []byte) ([]byte, error) {
	val, err := s.dataBase.Db.Get(key, nil) // Direct access for single key get from underlying DB
	if err != nil {
		if errors.Is(err, leveldb_errors.ErrNotFound) {
			return nil, nil // Return nil, nil for not found (Go's Optional equivalent)
		}
		return nil, fmt.Errorf("failed to get key %x from storage: %w", key, err)
	}
	return val, nil
}

// getOrElse(key: ByteArrayWrapper, defaultValue: ByteArrayWrapper): ByteArrayWrapper = dataBase.getOrElse(key, defaultValue)
func (s *VersionedLevelDbStorageAdapter) GetOrElse(key []byte, defaultValue []byte) []byte {
	val, err := s.Get(key)
	if err != nil || val == nil { // If error or not found, return defaultValue
		return defaultValue
	}
	return val
}

// get(keys: JList[ByteArrayWrapper]): JList[JPair[ByteArrayWrapper, Optional[ByteArrayWrapper]]]
func (s *VersionedLevelDbStorageAdapter) GetBatch(keys [][]byte) ([]storage.KeyValuePair, error) {
	results := make([]storage.KeyValuePair, len(keys))
	for i, key := range keys {
		val, err := s.Get(key) // Reuse the single Get method
		if err != nil {
			return nil, fmt.Errorf("failed to get key %x in batch: %w", key, err)
		}
		results[i] = storage.KeyValuePair{Key: key, Value: val} // Value is nil if not found
	}
	return results, nil
}

// getAll: JList[JPair[ByteArrayWrapper, ByteArrayWrapper]]
func (s *VersionedLevelDbStorageAdapter) GetAll() ([]storage.KeyValuePair, error) {
	var allPairs []storage.KeyValuePair
	iter := s.GetIterator() // Use the adapter's GetIterator method
	defer iter.Release()

	for iter.Next() {
		// The iterator now returns copies, so we don't need to copy here.
		allPairs = append(allPairs, storage.KeyValuePair{Key: iter.Key(), Value: iter.Value()})
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error during GetAll: %w", err)
	}

	return allPairs, nil
}

// lastVersionID(): Optional[ByteArrayWrapper]
func (s *VersionedLevelDbStorageAdapter) LastVersionID() ([]byte, error) {
	versions, err := s.dataBase.Versions()
	if err != nil {
		return nil, fmt.Errorf("failed to get versions for LastVersionID: %w", err)
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found in the db")
	}
	return versions[0], nil // Head option is the first element (newest version)
}

// update(version: ByteArrayWrapper, toUpdate: JList[JPair[ByteArrayWrapper, ByteArrayWrapper]], toRemove: util.List[ByteArrayWrapper]): Unit
func (s *VersionedLevelDbStorageAdapter) Update(versionID []byte, toUpdate []storage.KeyValuePair, toRemove [][]byte) error {
	// key for storing version shall not be used as key in any key-value pair in VersionedLDBKVStore
	for _, pair := range toUpdate {
		if bytes.Equal(pair.Key, versionID) {
			return storageErrors.ErrInvalidParameter("Version ID cannot be used as a key in 'toUpdate'")
		}
	}
	for _, key := range toRemove {
		if bytes.Equal(key, versionID) {
			return storageErrors.ErrInvalidParameter("Version ID cannot be used as a key in 'toRemove'")
		}
	}
	// Check for duplicate keys in toUpdate
	seenKeys := make(map[string]struct{})
	for _, pair := range toUpdate {
		keyStr := string(pair.Key)
		if _, exists := seenKeys[keyStr]; exists {
			return storageErrors.ErrInvalidParameter(fmt.Sprintf("Duplicate key in 'toUpdate': %x", pair.Key))
		}
		seenKeys[keyStr] = struct{}{}
	}
	// Check for duplicate keys in toRemove
	seenKeys = make(map[string]struct{})
	for _, key := range toRemove {
		keyStr := string(key)
		if _, exists := seenKeys[keyStr]; exists {
			return storageErrors.ErrInvalidParameter(fmt.Sprintf("Duplicate key in 'toRemove': %x", key))
		}
		seenKeys[keyStr] = struct{}{}
	}

	err := s.dataBase.Update(toUpdate, toRemove, versionID)
	if err != nil {
		var storageErr *storageErrors.Error
		if errors.As(err, &storageErr) {
			return storageErr
		}
		return err
	}

	return nil
}

// isVersionExist checks if a version exists in the dataBase.versions.
func (s *VersionedLevelDbStorageAdapter) isVersionExist(versionForSearch []byte) (bool, error) {
	versions, err := s.dataBase.Versions()
	if err != nil {
		return false, err
	}
	for _, v := range versions {
		if bytes.Equal(v, versionForSearch) {
			return true, nil
		}
	}
	return false, nil
}

// rollback(versionID: ByteArrayWrapper): Unit
func (s *VersionedLevelDbStorageAdapter) Rollback(versionID []byte) error {
	err := s.dataBase.RollbackTo(versionID)
	if err != nil {
		var storageErr *storageErrors.Error
		if errors.As(err, &storageErr) {
			return storageErr
		}
		return err
	}
	return nil
}

// rollbackVersions(): JList[ByteArrayWrapper]
func (s *VersionedLevelDbStorageAdapter) RollbackVersions() ([][]byte, error) {
	versions, err := s.dataBase.Versions()
	if err != nil {
		return nil, fmt.Errorf("failed to get versions for rollbackVersions: %w", err)
	}
	return versions, nil
}

// rollbackVersions(maxNumberOfItems: Int): JList[ByteArrayWrapper]
func (s *VersionedLevelDbStorageAdapter) RollbackVersionsLimited(maxNumberOfItems int) ([][]byte, error) {
	versions, err := s.dataBase.Versions()
	if err != nil {
		return nil, fmt.Errorf("failed to get versions for rollbackVersionsLimited: %w", err)
	}
	if maxNumberOfItems <= 0 {
		// Return empty slice if maxNumberOfItems is zero or negative
		return [][]byte{}, nil // Replace VersionType with actual element type
	}

	if maxNumberOfItems >= len(versions) {
		// Return full slice if maxNumberOfItems is too large
		return versions, nil
	}
	return versions[:maxNumberOfItems], nil
}

// close(): Unit
func (s *VersionedLevelDbStorageAdapter) Close() error {
	return s.dataBase.Close()
}

// createDb is a private helper function to create and open a LevelDB database.
// It ensures the directory exists and then opens the LevelDB instance.
func createDb(path string, versionsToKeep int) (*VersionedLDBKVStore, error) {
	// Creates the directory if it doesn't exist.
	err := os.MkdirAll(path, 0755) // 0755 permissions: rwx for owner, rx for group and others
	if err != nil {
		return nil, fmt.Errorf("failed to create database directory %s: %w", path, err)
	}

	// `leveldb.OpenFile` by default will create the database if it doesn't exist.
	// We don't need a `CreateIfMissing` field in `opt.Options`.
	options := &opt.Options{} // Using default options, or customize as needed
	db, err := leveldb.OpenFile(path, options)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB at %s: %w", path, err)
	}
	return NewVersionedLDBKVStore(db, versionsToKeep), nil
}

/*
// createDb(path: File): VersionedLDBKVStore
func createDb(path string, versionsToKeep int) (*VersionedLDBKVStore, error) {
	// path.mkdirs() equivalent
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create database directory %s: %w", path, err)
	}

	options := &opt.Options{
		CreateIfMissing: true,
	}
	db, err := leveldb.OpenFile(path, options)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB at %s: %w", path, err)
	}
	return NewVersionedLDBKVStore(db, versionsToKeep), nil
}
*/

// isEmpty: Boolean
func (s *VersionedLevelDbStorageAdapter) IsEmpty() bool {
	versions, err := s.dataBase.Versions()
	if err != nil {
		// Log the error or handle it more robustly in a real app
		fmt.Printf("Error checking if empty: %v\n", err)
		return true // Defensive: if we can't get versions, assume empty or error state
	}
	return len(versions) == 0
}

// numberOfVersions: Int
func (s *VersionedLevelDbStorageAdapter) NumberOfVersions() int {
	versions, err := s.dataBase.Versions()
	if err != nil {
		// Log the error
		fmt.Printf("Error getting number of versions: %v\n", err)
		return 0 // Defensive
	}
	return len(versions)
}

// getIterator(): StorageIterator
func (s *VersionedLevelDbStorageAdapter) GetIterator() storage.StorageIterator {
	return s.dataBase.GetIterator()
}
