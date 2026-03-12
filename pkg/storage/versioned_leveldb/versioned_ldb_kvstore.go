package versioned_leveldb

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/horizen-pes/pkg/storage"
)

const (
	ConstantsHashLength = 32
	ChangeSetPrefix     = 0x16
)

// versionsKeyForApp computes a deterministic per-app versions key: sha256("versions_" + appID).
func versionsKeyForApp(appID uint64) [ConstantsHashLength]byte {
	return sha256.Sum256([]byte(fmt.Sprintf("versions_%d", appID)))
}

type ChangeSet struct {
	InsertedKeys [][]byte
	Removed      []storage.KeyValuePair
	Altered      []storage.KeyValuePair
}

type ChangeSetSerializer struct{}

func (s *ChangeSetSerializer) ToBytes(cs ChangeSet) ([]byte, error) {
	return json.Marshal(cs)
}

func (s *ChangeSetSerializer) ParseBytes(data []byte) (ChangeSet, error) {
	var cs ChangeSet
	err := json.Unmarshal(data, &cs)
	return cs, err
}

type VersionedLDBKVStore struct {
	Db                  *leveldb.DB
	KeepVersions        int
	ChangeSetSerializer ChangeSetSerializer
	mu                  sync.RWMutex
}

func NewVersionedLDBKVStore(db *leveldb.DB, versionsToKeep int) *VersionedLDBKVStore {
	return &VersionedLDBKVStore{
		Db:                  db,
		KeepVersions:        versionsToKeep,
		ChangeSetSerializer: ChangeSetSerializer{},
	}
}

func (v *VersionedLDBKVStore) Update(appID uint64, toInsert []storage.KeyValuePair, toRemove [][]byte, versionID []byte) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if len(versionID) != ConstantsHashLength {
		return storageErrors.ErrInvalidParameter(fmt.Sprintf("Illegal version ID size: expected %d bytes, got %d", ConstantsHashLength, len(versionID)))
	}

	appVersionsKey := versionsKeyForApp(appID)

	tx, err := v.Db.OpenTransaction()
	if err != nil {
		return fmt.Errorf("failed to open transaction: %w", err)
	}
	defer tx.Discard()

	// NOTE: This duplicate check is global (not per-app). ChangeSet entries share the flat LevelDB
	// keyspace keyed by versionID. In practice, versionID = sha256(stateRoot || data) so cross-app
	// collisions are extremely unlikely. If per-app namespacing of changeset keys is needed in the
	// future, prefix the key with appID here and in RollbackTo.
	_, err = tx.Get(versionID, nil)
	if err == nil {
		return storageErrors.ErrVersionAlreadyExists(versionID)
	}
	if !errors.Is(err, leveldb_errors.ErrNotFound) {
		return fmt.Errorf("failed to check existing version ID: %w", err)
	}

	insertedKeys := make([][]byte, 0, len(toInsert))
	altered := make([]storage.KeyValuePair, 0)
	removed := make([]storage.KeyValuePair, 0)

	for _, kv := range toInsert {
		oldValue, err := tx.Get(kv.Key, nil)
		if err != nil {
			if errors.Is(err, leveldb_errors.ErrNotFound) {
				insertedKeys = append(insertedKeys, kv.Key)
			} else {
				return fmt.Errorf("failed to get old value for insert/update of key %x: %w", kv.Key, err)
			}
		} else {
			altered = append(altered, storage.KeyValuePair{Key: kv.Key, Value: oldValue})
		}
	}

	for _, k := range toRemove {
		oldValue, err := tx.Get(k, nil)
		if err != nil {
			if !errors.Is(err, leveldb_errors.ErrNotFound) {
				return fmt.Errorf("failed to get old value for removal of key %x: %w", k, err)
			}
		} else {
			removed = append(removed, storage.KeyValuePair{Key: k, Value: oldValue})
		}
	}

	changeSet := ChangeSet{
		InsertedKeys: insertedKeys,
		Removed:      removed,
		Altered:      altered,
	}

	currentVersionsBytes, err := tx.Get(appVersionsKey[:], nil)
	if err != nil && !errors.Is(err, leveldb_errors.ErrNotFound) {
		return fmt.Errorf("failed to get current versions list: %w", err)
	}

	var parsedCurrentVersions [][]byte
	if len(currentVersionsBytes) > 0 {
		if len(currentVersionsBytes)%ConstantsHashLength != 0 {
			return storageErrors.ErrInconsistentState("Corrupted versions key: length not a multiple of hash length")
		}
		for i := 0; i < len(currentVersionsBytes); i += ConstantsHashLength {
			parsedCurrentVersions = append(parsedCurrentVersions, currentVersionsBytes[i:i+ConstantsHashLength])
		}
	}

	newVersionsList := make([][]byte, 1, len(parsedCurrentVersions)+1)
	newVersionsList[0] = versionID
	newVersionsList = append(newVersionsList, parsedCurrentVersions...)

	versionsToShrink := [][]byte{}
	if len(newVersionsList) > v.KeepVersions {
		versionsToShrink = newVersionsList[v.KeepVersions:]
		newVersionsList = newVersionsList[0:v.KeepVersions]
	}
	updatedVersionsRaw := bytes.Join(newVersionsList, nil)

	batch := new(leveldb.Batch)
	batch.Put(appVersionsKey[:], updatedVersionsRaw)

	for _, verId := range versionsToShrink {
		batch.Delete(verId)
	}

	changeSetBytes, err := v.ChangeSetSerializer.ToBytes(changeSet)
	if err != nil {
		return fmt.Errorf("failed to serialize change set for version %x: %w", versionID, err)
	}
	prefixedChangeSetBytes := append([]byte{ChangeSetPrefix}, changeSetBytes...)
	batch.Put(versionID, prefixedChangeSetBytes)

	for _, kv := range toInsert {
		batch.Put(kv.Key, kv.Value)
	}
	for _, k := range toRemove {
		batch.Delete(k)
	}

	if err := tx.Write(batch, nil); err != nil {
		return fmt.Errorf("failed to write batch to transaction: %w", err)
	}

	return tx.Commit()
}

func (v *VersionedLDBKVStore) RollbackTo(appID uint64, versionID []byte) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	appVersionsKey := versionsKeyForApp(appID)

	tx, err := v.Db.OpenTransaction()
	if err != nil {
		return fmt.Errorf("failed to open transaction for rollback: %w", err)
	}
	defer tx.Discard()

	currentVersionsBytes, err := tx.Get(appVersionsKey[:], nil)
	if err != nil {
		if errors.Is(err, leveldb_errors.ErrNotFound) {
			return storageErrors.ErrVersionNotFound(versionID)
		}
		return fmt.Errorf("failed to get current versions for rollback: %w", err)
	}

	var actualVersions [][]byte
	if len(currentVersionsBytes) == 0 {
		return storageErrors.ErrVersionNotFound(versionID)
	}
	if len(currentVersionsBytes)%ConstantsHashLength != 0 {
		return storageErrors.ErrInconsistentState("Corrupted versions key during rollback: length not a multiple of hash length")
	}
	for i := 0; i < len(currentVersionsBytes); i += ConstantsHashLength {
		actualVersions = append(actualVersions, currentVersionsBytes[i:i+ConstantsHashLength])
	}

	targetVersionFound := false
	versionsToRollBack := [][]byte{}
	updatedVersionsList := [][]byte{}

	for _, ver := range actualVersions {
		if bytes.Equal(ver, versionID) {
			targetVersionFound = true
			updatedVersionsList = append(updatedVersionsList, ver)
		} else if !targetVersionFound {
			versionsToRollBack = append(versionsToRollBack, ver)
		} else {
			updatedVersionsList = append(updatedVersionsList, ver)
		}
	}

	if !targetVersionFound {
		return storageErrors.ErrVersionNotFound(versionID)
	}

	batch := new(leveldb.Batch)

	for _, verId := range versionsToRollBack {
		changeSetBytesWithPrefix, err := tx.Get(verId, nil)
		if err != nil {
			if errors.Is(err, leveldb_errors.ErrNotFound) {
				return storageErrors.ErrInconsistentState(fmt.Sprintf("Inconsistent versioned storage state: ChangeSet for version %x not found during rollback", verId))
			}
			return fmt.Errorf("failed to get change set for version %x during rollback: %w", verId, err)
		}

		if len(changeSetBytesWithPrefix) == 0 || changeSetBytesWithPrefix[0] != ChangeSetPrefix {
			return storageErrors.ErrInconsistentState(fmt.Sprintf("Inconsistent versioned storage state: Malformed ChangeSet for version %x", verId))
		}
		changeSetData := changeSetBytesWithPrefix[1:]

		changeSet, err := v.ChangeSetSerializer.ParseBytes(changeSetData)
		if err != nil {
			return fmt.Errorf("failed to parse change set for version %x during rollback: %w", verId, err)
		}

		for _, k := range changeSet.InsertedKeys {
			batch.Delete(k)
		}
		for _, kv := range changeSet.Removed {
			batch.Put(kv.Key, kv.Value)
		}
		for _, kv := range changeSet.Altered {
			batch.Put(kv.Key, kv.Value)
		}
		batch.Delete(verId)
	}

	updatedVersionsRaw := bytes.Join(updatedVersionsList, nil)
	batch.Put(appVersionsKey[:], updatedVersionsRaw)

	if err := tx.Write(batch, nil); err != nil {
		return fmt.Errorf("failed to write rollback batch to transaction: %w", err)
	}

	return tx.Commit()
}

func (v *VersionedLDBKVStore) Versions(appID uint64) ([][]byte, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.versionsUnlocked(appID)
}

// versionsUnlocked reads the version list for an app without acquiring the lock.
// Caller must hold at least a read lock.
func (v *VersionedLDBKVStore) versionsUnlocked(appID uint64) ([][]byte, error) {
	appVersionsKey := versionsKeyForApp(appID)
	versionsBytes, err := v.Db.Get(appVersionsKey[:], nil)
	if err != nil {
		if errors.Is(err, leveldb_errors.ErrNotFound) {
			return [][]byte{}, nil
		}
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	var versions [][]byte
	if len(versionsBytes)%ConstantsHashLength != 0 {
		return nil, storageErrors.ErrInconsistentState("Corrupted versions key: length not a multiple of hash length")
	}
	for i := 0; i < len(versionsBytes); i += ConstantsHashLength {
		versions = append(versions, versionsBytes[i:i+ConstantsHashLength])
	}
	return versions, nil
}

// allKnownAppVersionKeys performs a full DB scan to discover all per-app versions keys
// and changeset entries. It identifies internal metadata by structure:
//   - 32-byte key with value starting with ChangeSetPrefix → changeset entry
//   - 32-byte key with value length divisible by 32 (not ChangeSetPrefix) → versions key
//
// This heuristic is safe because all application data keys use string prefixes
// ("appstate_<id>", "wasm_<id>", etc.) which are never exactly 32 bytes long.
// If a new key schema introduces 32-byte data keys, this function must be updated.
//
// Caller must hold at least v.mu.RLock to ensure consistency with the iterator snapshot.
// This is used by GetIterator to filter out internal metadata.
func (v *VersionedLDBKVStore) allKnownAppVersionKeys() (map[string]struct{}, [][]byte, error) {
	excluded := make(map[string]struct{})
	var allVersionIDs [][]byte

	// Scan all keys in the DB to find per-app versions keys.
	// A versions key stores concatenated 32-byte version hashes (no ChangeSetPrefix).
	// We identify them by checking: value is non-empty, length % 32 == 0, and first byte != ChangeSetPrefix.
	iter := v.Db.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		// Skip data keys (app state, wasm, etc.) — they don't start with ChangeSetPrefix
		// and aren't 32-byte aligned. But ChangeSet entries DO start with ChangeSetPrefix.
		// Versions keys contain raw concatenated 32-byte hashes.
		if len(key) == ConstantsHashLength && len(value) > 0 {
			// Could be a versions key or a changeset key.
			if value[0] == ChangeSetPrefix {
				// This is a ChangeSet entry (version ID -> changeset data).
				// Exclude it from the iterator.
				excluded[string(key)] = struct{}{}
				allVersionIDs = append(allVersionIDs, append([]byte{}, key...))
			} else if len(value)%ConstantsHashLength == 0 {
				// This looks like a versions key (concatenated 32-byte version hashes).
				excluded[string(key)] = struct{}{}
				// Also parse and exclude the version IDs it references.
				for i := 0; i < len(value); i += ConstantsHashLength {
					vid := value[i : i+ConstantsHashLength]
					excluded[string(vid)] = struct{}{}
					allVersionIDs = append(allVersionIDs, append([]byte{}, vid...))
				}
			}
		}
	}
	if err := iter.Error(); err != nil {
		return nil, nil, fmt.Errorf("failed to scan for app version keys: %w", err)
	}

	return excluded, allVersionIDs, nil
}

func (v *VersionedLDBKVStore) VersionIDExists(appID uint64, versionID []byte) (bool, error) {
	versions, err := v.Versions(appID)
	if err != nil {
		return false, err
	}
	for _, ver := range versions {
		if bytes.Equal(ver, versionID) {
			return true, nil
		}
	}
	return false, nil
}

func (v *VersionedLDBKVStore) GetIterator() storage.StorageIterator {
	v.mu.RLock()
	defer v.mu.RUnlock()

	snapshot, err := v.Db.GetSnapshot()
	if err != nil {
		return &errorIterator{err: fmt.Errorf("leveldb: %w", err)}
	}

	excludedKeysMap, _, err := v.allKnownAppVersionKeys()
	if err != nil {
		snapshot.Release()
		return &errorIterator{err: fmt.Errorf("failed to get version keys for iterator filtering: %w", err)}
	}

	return &ldbStorageIterator{
		iter:         snapshot.NewIterator(nil, nil),
		excludedKeys: excludedKeysMap,
		snapshot:     snapshot,
	}
}

func (v *VersionedLDBKVStore) Close() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.Db.Close()
}

type ldbStorageIterator struct {
	iter         iterator.Iterator
	excludedKeys map[string]struct{}
	snapshot     *leveldb.Snapshot
}

func (it *ldbStorageIterator) Next() bool {
	for it.iter.Next() {
		if _, isExcluded := it.excludedKeys[string(it.iter.Key())]; isExcluded {
			continue
		}
		return true
	}
	return false
}

func (it *ldbStorageIterator) Key() []byte {
	key := it.iter.Key()
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)
	return keyCopy
}

func (it *ldbStorageIterator) Value() []byte {
	val := it.iter.Value()
	valCopy := make([]byte, len(val))
	copy(valCopy, val)
	return valCopy
}

func (it *ldbStorageIterator) Release() {
	it.iter.Release()
	it.snapshot.Release()
}

func (it *ldbStorageIterator) Error() error {
	return it.iter.Error()
}

type errorIterator struct {
	err error
}

func (ei *errorIterator) Next() bool {
	return false
}

func (ei *errorIterator) Key() []byte {
	return nil
}

func (ei *errorIterator) Value() []byte {
	return nil
}

func (ei *errorIterator) Release() {}

func (ei *errorIterator) Error() error {
	return ei.err
}
