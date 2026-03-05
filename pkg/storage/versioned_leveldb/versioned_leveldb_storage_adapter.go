package versioned_leveldb

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"

	storageErrors "github.com/HorizenOfficial/vela/pkg/storage/errors"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/HorizenOfficial/vela/pkg/storage"
)

const (
	DefaultVersionsToKeep = 10
)

type VersionedLevelDbStorageAdapter struct {
	dataBase    *VersionedLDBKVStore
	versionsKey []byte
	mu          sync.RWMutex
}

func NewVersionedLevelDbStorageAdapter(pathToDB string) (*VersionedLevelDbStorageAdapter, error) {
	return NewVersionedLevelDbStorageAdapterWithVersions(pathToDB, DefaultVersionsToKeep)
}

func NewVersionedLevelDbStorageAdapterWithVersions(pathToDB string, versionsToKeep int) (*VersionedLevelDbStorageAdapter, error) {
	dataBase, err := createDb(pathToDB, versionsToKeep)
	if err != nil {
		return nil, err
	}
	return &VersionedLevelDbStorageAdapter{
		dataBase:    dataBase,
		versionsKey: VersionsKey[:],
	}, nil
}

func (s *VersionedLevelDbStorageAdapter) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, err := s.dataBase.Db.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb_errors.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get key %x from storage: %w", key, err)
	}
	return val, nil
}

func (s *VersionedLevelDbStorageAdapter) GetOrElse(key []byte, defaultValue []byte) []byte {
	val, err := s.Get(key)
	if err != nil || val == nil {
		return defaultValue
	}
	return val
}

func (s *VersionedLevelDbStorageAdapter) GetBatch(keys [][]byte) ([]storage.KeyValuePair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	results := make([]storage.KeyValuePair, len(keys))
	for i, key := range keys {
		val, err := s.dataBase.Db.Get(key, nil)
		if err != nil {
			if !errors.Is(err, leveldb_errors.ErrNotFound) {
				return nil, fmt.Errorf("failed to get key %x in batch: %w", key, err)
			}
			val = nil
		}
		results[i] = storage.KeyValuePair{Key: key, Value: val}
	}
	return results, nil
}

func (s *VersionedLevelDbStorageAdapter) GetAll() ([]storage.KeyValuePair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var allPairs []storage.KeyValuePair
	iter := s.GetIterator()
	defer iter.Release()

	for iter.Next() {
		allPairs = append(allPairs, storage.KeyValuePair{Key: iter.Key(), Value: iter.Value()})
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error during GetAll: %w", err)
	}

	return allPairs, nil
}

func (s *VersionedLevelDbStorageAdapter) LastVersionID() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	versions, err := s.dataBase.Versions()
	if err != nil {
		return nil, fmt.Errorf("failed to get versions for LastVersionID: %w", err)
	}
	if len(versions) == 0 {
		return nil, storageErrors.ErrNoVersionInDb("no versions found in the db")
	}
	return versions[0], nil
}

func (s *VersionedLevelDbStorageAdapter) Update(versionID []byte, toUpdate []storage.KeyValuePair, toRemove [][]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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

	seenKeys := make(map[string]struct{})
	for _, pair := range toUpdate {
		keyStr := string(pair.Key)
		if _, exists := seenKeys[keyStr]; exists {
			return storageErrors.ErrInvalidParameter(fmt.Sprintf("Duplicate key in 'toUpdate': %x", pair.Key))
		}
		seenKeys[keyStr] = struct{}{}
	}

	for _, key := range toRemove {
		keyStr := string(key)
		if _, exists := seenKeys[keyStr]; exists {
			return storageErrors.ErrInvalidParameter(fmt.Sprintf("Duplicate key in 'toRemove': %x", key))
		}
		seenKeys[keyStr] = struct{}{}
	}

	return s.dataBase.Update(toUpdate, toRemove, versionID)
}

func (s *VersionedLevelDbStorageAdapter) Rollback(versionID []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dataBase.RollbackTo(versionID)
}

func (s *VersionedLevelDbStorageAdapter) RollbackVersions() ([][]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dataBase.Versions()
}

func (s *VersionedLevelDbStorageAdapter) RollbackVersionsLimited(maxNumberOfItems int) ([][]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	versions, err := s.dataBase.Versions()
	if err != nil {
		return nil, fmt.Errorf("failed to get versions for rollbackVersionsLimited: %w", err)
	}
	if maxNumberOfItems <= 0 {
		return [][]byte{}, nil
	}
	if maxNumberOfItems >= len(versions) {
		return versions, nil
	}
	return versions[:maxNumberOfItems], nil
}

func (s *VersionedLevelDbStorageAdapter) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dataBase.Close()
}

func createDb(path string, versionsToKeep int) (*VersionedLDBKVStore, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create database directory %s: %w", path, err)
	}
	options := &opt.Options{}
	db, err := leveldb.OpenFile(path, options)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB at %s: %w", path, err)
	}
	return NewVersionedLDBKVStore(db, versionsToKeep), nil
}

func (s *VersionedLevelDbStorageAdapter) IsEmpty() (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	versions, err := s.dataBase.Versions()
	if err != nil {
		return true, fmt.Errorf("error checking if empty: %w", err)
	}
	return len(versions) == 0, nil
}

func (s *VersionedLevelDbStorageAdapter) NumberOfVersions() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	versions, err := s.dataBase.Versions()
	if err != nil {
		return 0, fmt.Errorf("error getting number of versions: %w", err)
	}
	return len(versions), nil
}

func (s *VersionedLevelDbStorageAdapter) GetIterator() storage.StorageIterator {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dataBase.GetIterator()
}

var _ storage.VersionedStorage = (*VersionedLevelDbStorageAdapter)(nil)

