package storage

// KeyValuePair represents a key-value pair for storage operations.
type KeyValuePair struct {
	Key   []byte
	Value []byte
}

// StorageIterator defines an interface for iterating over key-value pairs in storage.
type StorageIterator interface {
	Next() bool
	// Key returns the key of the current entry. The slice is only valid until the next call to Next().
	Key() []byte
	// Value returns the value of the current entry. The slice is only valid until the next call to Next().
	Value() []byte
	Release()
	Error() error
}

// VersionedStorage defines the interface for a versioned key-value storage.
// Versioning is per-application: each app has its own independent version chain.
type VersionedStorage interface {
	Get(key []byte) ([]byte, error)
	GetOrElse(key []byte, defaultValue []byte) []byte
	GetBatch(keys [][]byte) ([]KeyValuePair, error)
	GetAll() ([]KeyValuePair, error)
	LastVersionID(appID uint64) ([]byte, error)
	Update(appID uint64, versionID []byte, toUpdate []KeyValuePair, toRemove [][]byte) error
	Rollback(appID uint64, versionID []byte) error
	RollbackVersions(appID uint64) ([][]byte, error)
	RollbackVersionsLimited(appID uint64, maxNumberOfItems int) ([][]byte, error)
	Close() error
	IsEmpty(appID uint64) (bool, error)
	NumberOfVersions(appID uint64) (int, error)
	GetIterator() StorageIterator
}
