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
	appStatePrefix              = "appstate_"
	wasmPrefix                  = "wasm_"
	deanonymizationReportPrefix = "deanon_"
	userKeyPrefix               = "userkey_"
)

// just two attributes for the time being; should we need more customization we can add them here
type VersionedLevelDBConfig struct {
	DBPath         string
	VersionsToKeep int
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
		key := []byte(appStatePrefix + state.ApplicationID)
		toUpdate = append(toUpdate, storage.KeyValuePair{Key: key, Value: value})
	}

	for _, wasm := range wasmArray {
		if wasm == nil {
			return fmt.Errorf("wasm entry is nil")
		}

		key := []byte(wasmPrefix + wasm.ApplicationID)
		toUpdate = append(toUpdate, storage.KeyValuePair{Key: key, Value: wasm.Bytecode})
	}

	err := vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store application state in Versioned LevelDB: %w", err)
	}
	return nil
}

func (vdl *VersionedLevelDBAppStateStore) GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error) {
	if err := vdl.checkClosed("state store"); err != nil {
		return nil, err
	}
	key := []byte(appStatePrefix + applicationID)
	value, err := vdl.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get application state from Versioned LevelDB: %w", err)
	}
	if value == nil {
		return nil, storageErrors.ErrNotFound(fmt.Sprintf("application state not found: %s", applicationID))
	}

	state := &common.ApplicationState{}
	err = json.Unmarshal(value, state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal application state: %w", err)
	}
	return state, nil
}

func (vdl *VersionedLevelDBAppStateStore) GetWASMBytecode(ctx context.Context, applicationID string) ([]byte, error) {
	if err := vdl.checkClosed("state store"); err != nil {
		return nil, err
	}
	key := []byte(wasmPrefix + applicationID)
	value, err := vdl.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get WASM bytecode from Versioned LevelDB: %w", err)
	}
	if value == nil {
		return nil, storageErrors.ErrNotFound("wasm bytecode not found for application: " + applicationID)
	}
	return value, nil
}

func (vdl *VersionedLevelDBAppStateStore) Rollback(versionID []byte) error {
	if err := vdl.checkClosed("state store"); err != nil {
		return err
	}
	return vdl.adapter.Rollback(versionID)
}

func (vdl *VersionedLevelDBAppStateStore) LastVersionID() ([]byte, error) {
	if err := vdl.checkClosed("state store"); err != nil {
		return nil, err
	}
	return vdl.adapter.LastVersionID()
}

func (vdl *VersionedLevelDBAppStateStore) ListVersions() ([][]byte, error) {
	if err := vdl.checkClosed("state store"); err != nil {
		return nil, err
	}
	return vdl.adapter.RollbackVersions()
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

// LevelDBUserKeyStore is a non-versioned store for user keys.
type LevelDBUserKeyStore struct {
	Adapter *LevelDbStorageAdapter
	closableStore
}

func NewLevelDBUserKeyStore(dbPath string) (*LevelDBUserKeyStore, error) {
	adapter, err := NewLevelDbStorageAdapter(dbPath)
	if err != nil {
		return nil, err
	}
	return &LevelDBUserKeyStore{Adapter: adapter}, nil
}

func (s *LevelDBUserKeyStore) StoreUserKey(ctx context.Context, userID string, publicKey []byte) error {
	if err := s.checkClosed("user key store"); err != nil {
		return err
	}
	key := []byte(userKeyPrefix + userID)
	err := s.Adapter.Put(key, publicKey)
	if err != nil {
		return fmt.Errorf("failed to store user key in LevelDB: %w", err)
	}
	return nil
}

func (s *LevelDBUserKeyStore) GetUserKey(ctx context.Context, userID string) ([]byte, error) {
	if err := s.checkClosed("user key store"); err != nil {
		return nil, err
	}
	key := []byte(userKeyPrefix + userID)
	value, err := s.Adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get user key from LevelDB: %w", err)
	}
	if value == nil {
		return nil, storageErrors.ErrNotFound("public key not found for user: " + userID)
	}
	return value, nil
}

func (s *LevelDBUserKeyStore) Close() error {
	return s.close(s.Adapter.Close)
}

var _ storage.ApplicationUserKeyStore = (*LevelDBUserKeyStore)(nil)

// LevelDBReportStore is a non-versioned store for reports.
type LevelDBReportStore struct {
	Adapter *LevelDbStorageAdapter
	closableStore
}

func NewLevelDBReportStore(dbPath string) (*LevelDBReportStore, error) {
	adapter, err := NewLevelDbStorageAdapter(dbPath)
	if err != nil {
		return nil, err
	}
	return &LevelDBReportStore{Adapter: adapter}, nil
}

func (s *LevelDBReportStore) StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	if err := s.checkClosed("report store"); err != nil {
		return err
	}
	value, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal deanonymization report: %w", err)
	}
	key := []byte(deanonymizationReportPrefix + report.ReportID)
	err = s.Adapter.Put(key, value)
	if err != nil {
		return fmt.Errorf("failed to store deanonymization report in LevelDB: %w", err)
	}
	return nil
}

func (s *LevelDBReportStore) GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error) {
	if err := s.checkClosed("report store"); err != nil {
		return nil, err
	}
	key := []byte(deanonymizationReportPrefix + reportID)
	value, err := s.Adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get deanonymization report from LevelDB: %w", err)
	}
	if value == nil {
		return nil, storageErrors.ErrNotFound("deanonymization report not found: " + reportID)
	}
	report := &common.DeanonymizationReport{}
	err = json.Unmarshal(value, report)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deanonymization report: %w", err)
	}
	return report, nil
}

func (s *LevelDBReportStore) Close() error {
	return s.close(s.Adapter.Close)
}

var _ storage.ApplicationReportStore = (*LevelDBReportStore)(nil)

// LevelDBDataLayer implements the DataLayer interface using LevelDB.
type LevelDBDataLayer struct {
	*VersionedLevelDBAppStateStore
	*LevelDBUserKeyStore
	*LevelDBReportStore
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

	userKeyStore, err := NewLevelDBUserKeyStore(filepath.Join(cfg.DBPath, "userkeys"))
	if err != nil {
		return nil, err
	}

	reportStore, err := NewLevelDBReportStore(filepath.Join(cfg.DBPath, "reports"))
	if err != nil {
		return nil, err
	}

	return &LevelDBDataLayer{
		VersionedLevelDBAppStateStore: appStateStore,
		LevelDBUserKeyStore:           userKeyStore,
		LevelDBReportStore:            reportStore,
	}, nil
}

func (d *LevelDBDataLayer) Close() error {
	if err := d.VersionedLevelDBAppStateStore.Close(); err != nil {
		return err
	}
	if err := d.LevelDBUserKeyStore.Close(); err != nil {
		return err
	}
	if err := d.LevelDBReportStore.Close(); err != nil {
		return err
	}
	return nil
}

var _ storage.DataLayer = (*LevelDBDataLayer)(nil)
