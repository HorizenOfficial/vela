package versioned_leveldb

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/storage"
	storageErrors "github.com/HorizenOfficial/vela/pkg/storage/errors"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
)

// prefixes, useful for avoiding collisions when using the same key
const (
	appStatePrefix           = "appstate_"
	wasmPrefix               = "wasm_"
	enclaveKeyRecoveryPrefix = "enclavekeyrecovery_"
)

// VersionedLevelDBConfig holds configuration for the versioned LevelDB storage.
type VersionedLevelDBConfig struct {
	DBPath         string
	VersionsToKeep int // Maximum number of historical versions to keep per application.
}

// generateVersionID creates a unique version identifier for a key-value pair.
// It computes the SHA-256 hash of the key concatenated with the value,
// ensuring that any change to the value results in a new version ID.
func generateVersionID(key []byte, data []byte) []byte {
	h := sha256.New()
	h.Write(key)
	h.Write(data)
	return h.Sum(nil)
}

type closableStore struct {
	isClosed bool
	mutex    sync.RWMutex
}

func (c *closableStore) checkClosed(storeName string) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.isClosed {
		return storageErrors.ErrStorageIsClosed(storeName + " is closed")
	}
	return nil
}

func (c *closableStore) close(closeFunc func() error) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isClosed {
		return nil
	}
	c.isClosed = true
	return closeFunc()
}

// VersionedLevelDBAppStateStore is a versioned store for application state.
type VersionedLevelDBAppStateStore struct {
	adapter *VersionedLevelDbStorageAdapter
	closableStore
}

func NewVersionedLevelDBAppStateStore(cfg VersionedLevelDBConfig) (*VersionedLevelDBAppStateStore, error) {
	adapter, err := NewVersionedLevelDbStorageAdapterWithVersions(cfg.DBPath, cfg.VersionsToKeep)
	if err != nil {
		return nil, fmt.Errorf("failed to create VersionedLevelDbStorageAdapter: %w", err)
	}
	return &VersionedLevelDBAppStateStore{adapter: adapter}, nil
}

func (vdl *VersionedLevelDBAppStateStore) Store(
	ctx context.Context,
	versionID []byte,
	state *common.ApplicationState,
	wasm *common.WASMData,
) error {
	if err := vdl.checkClosed("state store"); err != nil {
		return err
	}

	if state == nil && wasm == nil {
		return fmt.Errorf("cannot determine ApplicationID: both state and wasm are nil")
	}
	if state != nil && wasm != nil && state.ApplicationID != wasm.ApplicationID {
		return fmt.Errorf("inconsistent ApplicationID: state has %d, wasm has %d", state.ApplicationID, wasm.ApplicationID)
	}

	var appID common.ApplicationIdType
	if state != nil {
		appID = state.ApplicationID
	} else {
		appID = wasm.ApplicationID
	}

	toUpdate := []storage.KeyValuePair{}
	toRemove := [][]byte{}

	if state != nil {
		value, err := json.Marshal(state)
		if err != nil {
			return fmt.Errorf("failed to marshal application state: %w", err)
		}
		key := []byte(appStatePrefix + state.ApplicationID.String())
		toUpdate = append(toUpdate, storage.KeyValuePair{Key: key, Value: value})
	}

	if wasm != nil {
		key := []byte(wasmPrefix + wasm.ApplicationID.String())
		toUpdate = append(toUpdate, storage.KeyValuePair{Key: key, Value: wasm.Bytecode})
	}

	err := vdl.adapter.Update(uint64(appID), versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store application state in Versioned LevelDB: %w", err)
	}
	return nil
}

func (vdl *VersionedLevelDBAppStateStore) GetApplicationState(ctx context.Context, applicationID common.ApplicationIdType) (*common.ApplicationState, error) {
	if err := vdl.checkClosed("state store"); err != nil {
		return nil, err
	}
	key := []byte(appStatePrefix + applicationID.String())
	value, err := vdl.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get application state from Versioned LevelDB: %w", err)
	}
	if value == nil {
		return nil, storageErrors.ErrNotFound(fmt.Sprintf("application state not found: %d", applicationID))
	}

	state := &common.ApplicationState{}
	err = json.Unmarshal(value, state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal application state: %w", err)
	}
	return state, nil
}

func (vdl *VersionedLevelDBAppStateStore) GetWASMBytecode(ctx context.Context, applicationID common.ApplicationIdType) ([]byte, error) {
	if err := vdl.checkClosed("state store"); err != nil {
		return nil, err
	}
	key := []byte(wasmPrefix + applicationID.String())
	value, err := vdl.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get WASM bytecode from Versioned LevelDB: %w", err)
	}
	if value == nil {
		return nil, storageErrors.ErrNotFound("wasm bytecode not found for application: " + applicationID.String())
	}
	return value, nil
}

func (vdl *VersionedLevelDBAppStateStore) Rollback(appID common.ApplicationIdType, versionID []byte) error {
	if err := vdl.checkClosed("state store"); err != nil {
		return err
	}
	return vdl.adapter.Rollback(uint64(appID), versionID)
}

func (vdl *VersionedLevelDBAppStateStore) LastVersionID(appID common.ApplicationIdType) ([]byte, error) {
	if err := vdl.checkClosed("state store"); err != nil {
		return nil, err
	}
	return vdl.adapter.LastVersionID(uint64(appID))
}

func (vdl *VersionedLevelDBAppStateStore) ListVersions(appID common.ApplicationIdType) ([][]byte, error) {
	if err := vdl.checkClosed("state store"); err != nil {
		return nil, err
	}
	return vdl.adapter.RollbackVersions(uint64(appID))
}

func (vdl *VersionedLevelDBAppStateStore) Close() error {
	return vdl.close(vdl.adapter.Close)
}

// adapter returns the underlying VersionedLevelDbStorageAdapter instance. This is intended for testing purposes only.
func (vdl *VersionedLevelDBAppStateStore) getAdapter() *VersionedLevelDbStorageAdapter {
	return vdl.adapter
}

var _ storage.ApplicationStateStore = (*VersionedLevelDBAppStateStore)(nil)

// LevelDbStorageAdapter provides basic key-value storage without versioning.
type LevelDbStorageAdapter struct {
	db *leveldb.DB
	mu sync.RWMutex
}

func NewLevelDbStorageAdapter(pathToDB string) (*LevelDbStorageAdapter, error) {
	err := os.MkdirAll(pathToDB, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create database directory %s: %w", pathToDB, err)
	}
	db, err := leveldb.OpenFile(pathToDB, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB at %s: %w", pathToDB, err)
	}
	return &LevelDbStorageAdapter{db: db}, nil
}

func (s *LevelDbStorageAdapter) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, err := s.db.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb_errors.ErrNotFound) {
			return nil, nil // Return nil, nil for not found
		}
		return nil, fmt.Errorf("failed to get key %x from storage: %w", key, err)
	}
	return val, nil
}

func (s *LevelDbStorageAdapter) Put(key []byte, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Put(key, value, nil)
}

func (s *LevelDbStorageAdapter) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Close()
}

// LevelDBEnclaveKeyStore is a non-versioned store for enclave key recovery data.
type LevelDBEnclaveKeyStore struct {
	adapter *LevelDbStorageAdapter
	closableStore
}

func NewLevelDBEnclaveKeyStore(dbPath string) (*LevelDBEnclaveKeyStore, error) {
	adapter, err := NewLevelDbStorageAdapter(dbPath)
	if err != nil {
		return nil, err
	}
	return &LevelDBEnclaveKeyStore{adapter: adapter}, nil
}

func (s *LevelDBEnclaveKeyStore) StoreEnclaveKeySetRecovery(ctx context.Context, recoveryData *common.EnclaveKeySetRecovery) error {
	if err := s.checkClosed("enclave key store"); err != nil {
		return err
	}
	value, err := json.Marshal(recoveryData)
	if err != nil {
		return fmt.Errorf("failed to marshal enclave key set recovery data: %w", err)
	}
	key := []byte(enclaveKeyRecoveryPrefix)
	err = s.adapter.Put(key, value)
	if err != nil {
		return fmt.Errorf("failed to store enclave key set recovery data in LevelDB: %w", err)
	}
	return nil
}

func (s *LevelDBEnclaveKeyStore) GetEnclaveKeySetRecovery(ctx context.Context) (*common.EnclaveKeySetRecovery, error) {
	if err := s.checkClosed("enclave key store"); err != nil {
		return nil, err
	}
	key := []byte(enclaveKeyRecoveryPrefix)
	value, err := s.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get enclave key set recovery data from LevelDB: %w", err)
	}
	if value == nil {
		return nil, storageErrors.ErrNotFound("enclave key set recovery data not found")
	}
	recoveryData := &common.EnclaveKeySetRecovery{}
	err = json.Unmarshal(value, recoveryData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal enclave key set recovery data: %w", err)
	}
	return recoveryData, nil
}

func (s *LevelDBEnclaveKeyStore) Close() error {
	return s.close(s.adapter.Close)
}

var _ storage.EnclaveKeyStore = (*LevelDBEnclaveKeyStore)(nil)

// LevelDBDataLayer implements the DataLayer interface using LevelDB.
type LevelDBDataLayer struct {
	*VersionedLevelDBAppStateStore
	*LevelDBEnclaveKeyStore
}

func NewVersionedLevelDBDataLayer(cfg VersionedLevelDBConfig) (*LevelDBDataLayer, error) {
	appStateStoreCfg := VersionedLevelDBConfig{
		DBPath:         filepath.Join(cfg.DBPath, "state"),
		VersionsToKeep: cfg.VersionsToKeep,
	}

	// this DB uses versions
	appStateStore, err := NewVersionedLevelDBAppStateStore(appStateStoreCfg)
	if err != nil {
		return nil, err
	}

	enclaveKeyStore, err := NewLevelDBEnclaveKeyStore(filepath.Join(cfg.DBPath, "enclave_keys"))
	if err != nil {
		return nil, err
	}

	return &LevelDBDataLayer{
		VersionedLevelDBAppStateStore: appStateStore,
		LevelDBEnclaveKeyStore:        enclaveKeyStore,
	}, nil
}

func (d *LevelDBDataLayer) Close() error {
	errState := d.VersionedLevelDBAppStateStore.Close()
	errEnclave := d.LevelDBEnclaveKeyStore.Close()
	return errors.Join(errState, errEnclave)
}

var _ storage.DataLayer = (*LevelDBDataLayer)(nil)
