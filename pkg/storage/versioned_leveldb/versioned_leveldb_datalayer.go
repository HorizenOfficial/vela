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

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
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
	stateArray []*common.ApplicationState,
	wasmArray []*common.WASMData,
) error {
	if err := vdl.checkClosed("state store"); err != nil {
		return err
	}

	// Derive appID from stateArray. All items must share the same ApplicationID
	// (sequential processing — one app per Store call).
	appID, err := extractAppID(stateArray, wasmArray)
	if err != nil {
		return err
	}

	toUpdate := []storage.KeyValuePair{}
	toRemove := [][]byte{}

	for _, state := range stateArray {

		if state == nil {
			return fmt.Errorf("application state entry is nil")
		}

		value, err := json.Marshal(state)
		if err != nil {
			return fmt.Errorf("failed to marshal application state: %w", err)
		}

		key := []byte(appStatePrefix + state.ApplicationID.String())
		toUpdate = append(toUpdate, storage.KeyValuePair{Key: key, Value: value})
	}

	for _, wasm := range wasmArray {
		if wasm == nil {
			return fmt.Errorf("wasm entry is nil")
		}

		key := []byte(wasmPrefix + wasm.ApplicationID.String())
		toUpdate = append(toUpdate, storage.KeyValuePair{Key: key, Value: wasm.Bytecode})
	}

	err = vdl.adapter.Update(uint64(appID), versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store application state in Versioned LevelDB: %w", err)
	}
	return nil
}

// extractAppID derives the application ID from state and WASM arrays.
// All non-nil items must share the same ApplicationID.
// Returns an error if both arrays are empty/nil — with per-app versioning, every Store call
// must be associated with an app (this intentionally narrows the pre-multi-app API contract).
func extractAppID(stateArray []*common.ApplicationState, wasmArray []*common.WASMData) (common.ApplicationIdType, error) {
	var appID common.ApplicationIdType
	found := false

	for _, state := range stateArray {
		if state == nil {
			continue
		}
		if !found {
			appID = state.ApplicationID
			found = true
		} else if state.ApplicationID != appID {
			return 0, fmt.Errorf("inconsistent ApplicationID in stateArray: got %d and %d", appID, state.ApplicationID)
		}
	}

	for _, wasm := range wasmArray {
		if wasm == nil {
			continue
		}
		if !found {
			appID = wasm.ApplicationID
			found = true
		} else if wasm.ApplicationID != appID {
			return 0, fmt.Errorf("inconsistent ApplicationID in wasmArray: got %d and %d", appID, wasm.ApplicationID)
		}
	}

	if !found {
		return 0, fmt.Errorf("cannot determine ApplicationID: stateArray and wasmArray are both empty or nil")
	}
	return appID, nil
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
