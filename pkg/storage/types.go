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
type VersionedStorage interface {
	Get(key []byte) ([]byte, error)
	GetOrElse(key []byte, defaultValue []byte) []byte
	GetBatch(keys [][]byte) ([]KeyValuePair, error)
	GetAll() ([]KeyValuePair, error)
	LastVersionID() ([]byte, error)
	Update(versionID []byte, toUpdate []KeyValuePair, toRemove [][]byte) error
	Rollback(versionID []byte) error
	RollbackVersions() ([][]byte, error)
	RollbackVersionsLimited(maxNumberOfItems int) ([][]byte, error)
	Close() error
	IsEmpty() (bool, error)
	NumberOfVersions() (int, error)
	GetIterator() StorageIterator
}
