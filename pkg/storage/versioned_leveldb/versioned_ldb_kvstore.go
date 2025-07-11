package versioned_leveldb

import (
	"bytes" // Although goleveldb doesn't take context, it's good practice for wrapping methods
	"crypto/sha256"
	"encoding/json"
	"errors" // For errors.Is
	"fmt"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"

	// For io.Closer interface if needed (though not directly here)
	// For file system operations (like os.ErrNotExist, if needed)
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors" // Alias to avoid conflict with your custom Error type
	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/horizen-pes/pkg/common"
)

const (
	ConstantsHashLength = 32   // Assuming SHA256 hash length for version IDs
	ChangeSetPrefix     = 0x16 // From Scala's `ChangeSetPrefix: Byte = 0x16`
)

// VersionsKey is the key under which the concatenated list of version IDs is stored.
// It's initialized once to a SHA256 hash of the string "versions".
var VersionsKey [ConstantsHashLength]byte

// init function runs once when the package is imported.
func init() {
	VersionsKey = sha256.Sum256([]byte("versions"))
}

// ChangeSet represents the changes made in a version.
// This struct will be serialized and stored alongside each version ID.
type ChangeSet struct {
	InsertedKeys [][]byte              `json:"inserted_keys"` // Keys that were newly inserted
	Removed      []common.KeyValuePair `json:"removed"`       // Original key-value pairs that were removed
	Altered      []common.KeyValuePair `json:"altered"`       // Original key-value pairs whose values were changed
}

// ChangeSetSerializer provides methods to serialize/deserialize ChangeSet structs.
type ChangeSetSerializer struct{}

// ToBytes marshals a ChangeSet struct into JSON bytes.
func (s *ChangeSetSerializer) ToBytes(cs ChangeSet) ([]byte, error) {
	return json.Marshal(cs)
}

// ParseBytes unmarshals JSON bytes into a ChangeSet struct.
func (s *ChangeSetSerializer) ParseBytes(data []byte) (ChangeSet, error) {
	var cs ChangeSet
	err := json.Unmarshal(data, &cs)
	return cs, err
}

// VersionedLDBKVStore is the Go equivalent of Scala's VersionedLDBKVStore.
// It provides versioning capabilities on top of a LevelDB instance.
type VersionedLDBKVStore struct {
	Db                  *leveldb.DB
	KeepVersions        int // The number of versions to keep (Scala's `val keepVersions`)
	ChangeSetSerializer ChangeSetSerializer
}

// NewVersionedLDBKVStore creates a new instance of VersionedLDBKVStore.
// It takes an already opened LevelDB database connection and the number of versions to keep.
func NewVersionedLDBKVStore(db *leveldb.DB, versionsToKeep int) *VersionedLDBKVStore {
	return &VersionedLDBKVStore{
		Db:                  db,
		KeepVersions:        versionsToKeep,
		ChangeSetSerializer: ChangeSetSerializer{},
	}
}

// Update performs a versioned update operation.
// It inserts/updates `toInsert` pairs, removes `toRemove` keys, and records these changes
// under the given `versionID`. It also manages the versions history.
func (v *VersionedLDBKVStore) Update(toInsert []common.KeyValuePair, toRemove [][]byte, versionID []byte) error {
	// Validate the version ID size.
	if len(versionID) != ConstantsHashLength {
		return storageErrors.ErrInvalidParameter(fmt.Sprintf("Illegal version ID size: expected %d bytes, got %d", ConstantsHashLength, len(versionID)))
	}

	// Create a snapshot for consistent reads of the current state before modifications.
	snapshot, err := v.Db.GetSnapshot()
	if err != nil {
		return fmt.Errorf("failed to get snapshot: %w", err)
	}
	defer snapshot.Release() // Ensure the snapshot is released when the function exits.

	// Check if the version ID is already used.
	// We use the snapshot object directly for this read.
	_, err = snapshot.Get(versionID, nil) // CORRECTED: Call Get on snapshot directly, no ReadOptions needed
	if err == nil {
		return storageErrors.ErrVersionAlreadyExists(versionID) // Version ID is already present, cannot overwrite.
	}
	if !errors.Is(err, leveldb_errors.ErrNotFound) {
		return fmt.Errorf("failed to check existing version ID: %w", err)
	}

	// Prepare the ChangeSet by recording the original state of keys being modified or removed.
	insertedKeys := make([][]byte, 0, len(toInsert))
	altered := make([]common.KeyValuePair, 0)
	removed := make([]common.KeyValuePair, 0)

	// Process insertions/updates to determine `insertedKeys` and `altered` entries.
	for _, kv := range toInsert {
		oldValue, err := snapshot.Get(kv.Key, nil) // CORRECTED: Read from snapshot directly
		if err != nil {
			if errors.Is(err, leveldb_errors.ErrNotFound) {
				insertedKeys = append(insertedKeys, kv.Key) // Key was not present, so it's an insertion.
			} else {
				return fmt.Errorf("failed to get old value for insert/update of key %x: %w", kv.Key, err)
			}
		} else {
			// Key was present, so it's an alteration. Record the original key-value pair.
			altered = append(altered, common.KeyValuePair{Key: kv.Key, Value: oldValue})
		}
	}

	// Process removals to determine `removed` entries.
	for _, k := range toRemove {
		oldValue, err := snapshot.Get(k, nil) // CORRECTED: Read from snapshot directly
		if err != nil {
			if !errors.Is(err, leveldb_errors.ErrNotFound) {
				return fmt.Errorf("failed to get old value for removal of key %x: %w", k, err)
			}
			// If key not found, it means it didn't exist to be removed.
		} else {
			// Key was present, so it's a removal. Record the original key-value pair.
			removed = append(removed, common.KeyValuePair{Key: k, Value: oldValue})
		}
	}

	// Create the final ChangeSet.
	changeSet := ChangeSet{
		InsertedKeys: insertedKeys,
		Removed:      removed,
		Altered:      altered,
	}

	// --- Manage Version History ---

	// Get the current list of versions. This is typically read without a snapshot,
	// reflecting the most recent state of the versions metadata key.
	currentVersionsBytes, err := v.Db.Get(VersionsKey[:], nil)
	if err != nil && !errors.Is(err, leveldb_errors.ErrNotFound) {
		return fmt.Errorf("failed to get current versions list: %w", err)
	}

	// Parse current versions into a slice of byte slices.
	var parsedCurrentVersions [][]byte
	if len(currentVersionsBytes) > 0 {
		if len(currentVersionsBytes)%ConstantsHashLength != 0 {
			return storageErrors.ErrInconsistentState("Corrupted versions key: length not a multiple of hash length")
		}
		for i := 0; i < len(currentVersionsBytes); i += ConstantsHashLength {
			parsedCurrentVersions = append(parsedCurrentVersions, currentVersionsBytes[i:i+ConstantsHashLength])
		}
	}

	// Add the new version ID at the beginning (newer version first).
	newVersionsList := make([][]byte, 1, len(parsedCurrentVersions)+1)
	newVersionsList[0] = versionID
	newVersionsList = append(newVersionsList, parsedCurrentVersions...)

	// Determine versions to shrink (oldest ones beyond `keepVersions`).
	versionsToShrink := [][]byte{}
	if len(newVersionsList) > v.KeepVersions {
		versionsToShrink = newVersionsList[v.KeepVersions:]
		newVersionsList = newVersionsList[0:v.KeepVersions] // Keep only the desired number of versions.
	}
	updatedVersionsRaw := bytes.Join(newVersionsList, nil) // Concatenate remaining versions for storage.

	// --- Prepare Write Batch ---
	batch := new(leveldb.Batch) // Create a new LevelDB write batch.

	// 1. Update the `VersionsKey` with the new, shrunk list of versions.
	batch.Put(VersionsKey[:], updatedVersionsRaw)

	// 2. Delete the ChangeSets for the versions that are being shrunk/removed from history.
	for _, verId := range versionsToShrink {
		batch.Delete(verId)
	}

	// 3. Store the current version's ChangeSet.
	changeSetBytes, err := v.ChangeSetSerializer.ToBytes(changeSet)
	if err != nil {
		return fmt.Errorf("failed to serialize change set for version %x: %w", versionID, err)
	}
	// Prepend the ChangeSetPrefix byte as per Scala's `ChangeSetPrefix +: ChangeSetSerializer.toBytes(changeSet)`.
	prefixedChangeSetBytes := append([]byte{ChangeSetPrefix}, changeSetBytes...)
	batch.Put(versionID, prefixedChangeSetBytes) // Store change set data under its version ID key.

	// 4. Apply the actual data changes (inserts/updates and removals).
	for _, kv := range toInsert {
		batch.Put(kv.Key, kv.Value)
	}
	for _, k := range toRemove {
		batch.Delete(k)
	}

	// --- Execute Write Batch ---
	err = v.Db.Write(batch, nil) // Execute the batch write with default options.
	if err != nil {
		return fmt.Errorf("failed to write batch to LevelDB: %w", err)
	}

	return nil
}

// RollbackTo rolls storage state back to the specified checkpoint (versionID).
// It reverts changes from newer versions until the target version is reached.
func (v *VersionedLDBKVStore) RollbackTo(versionID []byte) error {
	// Create a snapshot for consistent reads during the rollback process.
	snapshot, err := v.Db.GetSnapshot()
	if err != nil {
		return fmt.Errorf("failed to get snapshot for rollback: %w", err)
	}
	defer snapshot.Release() // Ensure snapshot is released.

	// Get the current list of versions (outside the snapshot, reflecting latest metadata).
	currentVersionsBytes, err := v.Db.Get(VersionsKey[:], nil)
	if err != nil {
		if errors.Is(err, leveldb_errors.ErrNotFound) {
			return storageErrors.ErrVersionNotFound(versionID) // No versions exist if VersionsKey is not found.
		}
		return fmt.Errorf("failed to get current versions for rollback: %w", err)
	}

	// Parse current versions into a slice of byte slices.
	var actualVersions [][]byte
	if len(currentVersionsBytes) == 0 {
		return storageErrors.ErrVersionNotFound(versionID) // No versions to roll back from.
	}
	if len(currentVersionsBytes)%ConstantsHashLength != 0 {
		return storageErrors.ErrInconsistentState("Corrupted versions key during rollback: length not a multiple of hash length")
	}
	for i := 0; i < len(currentVersionsBytes); i += ConstantsHashLength {
		actualVersions = append(actualVersions, currentVersionsBytes[i:i+ConstantsHashLength])
	}

	// Identify versions to roll back (newer than target) and versions to keep (target and older).
	targetVersionFound := false
	versionsToRollBack := [][]byte{}  // Versions to revert (from newest to target-1)
	updatedVersionsList := [][]byte{} // Versions to keep (from target to oldest)

	for _, ver := range actualVersions {
		if bytes.Equal(ver, versionID) {
			targetVersionFound = true
			updatedVersionsList = append(updatedVersionsList, ver) // This is the target version, keep it.
		} else if !targetVersionFound {
			versionsToRollBack = append(versionsToRollBack, ver) // This version is newer than the target, needs to be rolled back.
		} else {
			updatedVersionsList = append(updatedVersionsList, ver) // This version is older than the target, keep it.
		}
	}

	if !targetVersionFound {
		return storageErrors.ErrVersionNotFound(versionID) // Target version not found in the history.
	}

	// Create a write batch for the rollback operations.
	batch := new(leveldb.Batch)

	// Revert all changes from the newest version down to the one just before the targeted version.
	for _, verId := range versionsToRollBack {
		// Get the ChangeSet for the version from the snapshot.
		changeSetBytesWithPrefix, err := snapshot.Get(verId, nil) // CORRECTED: Call Get on snapshot directly
		if err != nil {
			if errors.Is(err, leveldb_errors.ErrNotFound) {
				return storageErrors.ErrInconsistentState(fmt.Sprintf("Inconsistent versioned storage state: ChangeSet for version %x not found during rollback", verId))
			}
			return fmt.Errorf("failed to get change set for version %x during rollback: %w", verId, err)
		}

		// Validate and parse the ChangeSet.
		if len(changeSetBytesWithPrefix) == 0 || changeSetBytesWithPrefix[0] != ChangeSetPrefix {
			return storageErrors.ErrInconsistentState(fmt.Sprintf("Inconsistent versioned storage state: Malformed ChangeSet for version %x", verId))
		}
		changeSetData := changeSetBytesWithPrefix[1:] // Remove the prefix.

		changeSet, err := v.ChangeSetSerializer.ParseBytes(changeSetData)
		if err != nil {
			return fmt.Errorf("failed to parse change set for version %x during rollback: %w", verId, err)
		}

		// Apply the inverse of the changes:
		// 1. Delete keys that were newly inserted in this version.
		for _, k := range changeSet.InsertedKeys {
			batch.Delete(k)
		}
		// 2. Put back keys and values that were removed in this version.
		for _, kv := range changeSet.Removed {
			batch.Put(kv.Key, kv.Value)
		}
		// 3. Put back the original values for keys that were altered in this version.
		for _, kv := range changeSet.Altered {
			batch.Put(kv.Key, kv.Value)
		}
		// 4. Delete the ChangeSet entry itself for this version.
		batch.Delete(verId)
	}

	// Update the `VersionsKey` with the new list of versions (only the ones being kept).
	updatedVersionsRaw := bytes.Join(updatedVersionsList, nil) // Concatenate remaining versions.
	batch.Put(VersionsKey[:], updatedVersionsRaw)

	// Execute the rollback batch write.
	err = v.Db.Write(batch, nil)
	if err != nil {
		return fmt.Errorf("failed to write rollback batch to LevelDB: %w", err)
	}

	return nil
}

// Versions returns a sequence of all available version IDs, newest first.
func (v *VersionedLDBKVStore) Versions() ([][]byte, error) {
	// Read the VersionsKey directly, not from a snapshot, as it's metadata.
	versionsBytes, err := v.Db.Get(VersionsKey[:], nil)
	if err != nil {
		if errors.Is(err, leveldb_errors.ErrNotFound) {
			return [][]byte{}, nil // No versions yet, return empty slice.
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

// VersionIDExists checks if a specific version ID exists in the store's history.
func (v *VersionedLDBKVStore) VersionIDExists(versionID []byte) (bool, error) {
	versions, err := v.Versions() // Retrieves all versions.
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

// GetIterator returns a StorageIterator for the underlying LevelDB.
// It creates an iterator with a snapshot to ensure a consistent view of the database
// at the time the iterator is created. The caller must Release() the iterator.
func (v *VersionedLDBKVStore) GetIterator() common.StorageIterator {
	snapshot, err := v.Db.GetSnapshot()
	if err != nil {
		// In a real application, you might want to log this error or return a nil iterator
		// and expect callers to check the error. For this implementation, we panic or
		// return a dummy iterator if a snapshot can't be obtained, as it's a critical failure.
		//panic(fmt.Sprintf("Failed to get snapshot for iterator: %v", err)) // Or handle gracefully
		// If the DB is closed, getting a snapshot will fail.
		// Return a dummy iterator that immediately reports the "closed" error.
		return &errorIterator{err: fmt.Errorf("leveldb: %w", err)} // Wrap the original error

	}
	// Get a list of excluded keys (VersionsKey and all current version IDs which store change sets)
	// This makes the iterator only yield actual data keys.
	excludedKeysMap := make(map[[ConstantsHashLength]byte]struct{}) // CORRECTED: Map key type
	excludedKeysMap[VersionsKey] = struct{}{}                       // CORRECTED: Use array directly for map key

	currentVersions, err := v.Versions() // Get current version IDs (might error if DB closed)
	if err != nil {
		// If Versions() returns an error (e.g., DB closed), handle it gracefully for the iterator.
		return &errorIterator{err: fmt.Errorf("failed to get version list for iterator filtering: %w", err)}
	}
	for _, verID := range currentVersions {
		if len(verID) == ConstantsHashLength {
			var keyArr [ConstantsHashLength]byte
			copy(keyArr[:], verID)
			excludedKeysMap[keyArr] = struct{}{} // CORRECTED: Use array for map key
		} else {
			// This indicates a severe inconsistency in stored version IDs.
			// Log an error and potentially return a failing iterator or panic in production.
			fmt.Printf("WARNING: Version ID %x has unexpected length %d for iterator exclusion\n", verID, len(verID))
		}
	}

	return &ldbStorageIterator{
		iter:         snapshot.NewIterator(nil, nil),
		excludedKeys: excludedKeysMap,
	}
}

// Close closes the underlying LevelDB database.
func (v *VersionedLDBKVStore) Close() error {
	return v.Db.Close()
}

// --- LDB Storage Iterator Implementation ---

// ldbStorageIterator implements the common.StorageIterator interface for LevelDB.

type ldbStorageIterator struct {
	iter         iterator.Iterator                      // The underlying LevelDB iterator.
	excludedKeys map[[ConstantsHashLength]byte]struct{} // CORRECTED: Map key type
}

// Next advances the iterator to the next key-value pair.
func (it *ldbStorageIterator) Next() bool {
	for it.iter.Next() {
		currentKey := it.iter.Key()
		// Check if the current key has the hash length. Only hash keys are metadata.
		if len(currentKey) == ConstantsHashLength {
			var keyArr [ConstantsHashLength]byte
			copy(keyArr[:], currentKey)
			if _, isExcluded := it.excludedKeys[keyArr]; isExcluded { // CORRECTED: Use array for map lookup
				continue // This key is excluded, try the next one.
			}
		}
		// If key is not of hash length OR not found in excludedKeys (i.e., it's a data key), return true.
		return true
	}
	return false
}

// Key returns the key of the current key-value pair.
// It returns a copy of the key, as the underlying slice is only valid until the next call to Next().
func (it *ldbStorageIterator) Key() []byte {
	key := it.iter.Key()
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)
	return keyCopy
}

// Value returns the value of the current key-value pair.
// It returns a copy of the value, as the underlying slice is only valid until the next call to Next().
func (it *ldbStorageIterator) Value() []byte {
	val := it.iter.Value()
	valCopy := make([]byte, len(val))
	copy(valCopy, val)
	return valCopy
}

// Release releases the iterator's resources. This MUST be called when done with the iterator.
func (it *ldbStorageIterator) Release() {
	it.iter.Release()
}

// Error returns the last error encountered by the iterator.
func (it *ldbStorageIterator) Error() error {
	return it.iter.Error()
}

// --- Error Iterator for Closed DB ---

// errorIterator is a dummy iterator that always reports an error,
// used when the underlying DB is closed and an iterator cannot be created.

type errorIterator struct {
	err error
}

func (ei *errorIterator) Next() bool {
	return false // Cannot advance
}

func (ei *errorIterator) Key() []byte {
	return nil // No key
}

func (ei *errorIterator) Value() []byte {
	return nil // No value
}

func (ei *errorIterator) Release() {
	// Nothing to release
}

func (ei *errorIterator) Error() error {
	return ei.err // Always return the error it was initialized with
}