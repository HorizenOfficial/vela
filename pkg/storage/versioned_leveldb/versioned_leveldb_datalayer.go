package versioned_leveldb

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/errors"
)

type VersionedLevelDBDataLayer struct {
	adapter  *VersionedLevelDbStorageAdapter
	isClosed bool
	mu       sync.RWMutex
}

// just two attributes for the time being; should we need more customization we can add them here
type VersionedLevelDBConfig struct {
	DBPath         string
	VersionsToKeep int
}

func NewVersionedLevelDBDataLayer(cfg VersionedLevelDBConfig) (*VersionedLevelDBDataLayer, error) {
	// Construct the final, specific path for this database instance.
	adapter, err := NewVersionedLevelDbStorageAdapterWithVersions(cfg.DBPath, cfg.VersionsToKeep)
	if err != nil {
		return nil, fmt.Errorf("failed to create VersionedLevelDbStorageAdapter: %w", err)
	}
	return &VersionedLevelDBDataLayer{adapter: adapter}, nil
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

func (vdl *VersionedLevelDBDataLayer) checkClosed() error {
	vdl.mu.RLock()
	defer vdl.mu.RUnlock()
	if vdl.isClosed {
		return errors.ErrStorageIsClosed("Versioned LevelDB data layer is closed")
	}
	return nil
}

func (vdl *VersionedLevelDBDataLayer) StoreApplicationState(ctx context.Context, state *common.ApplicationState) error {
	if err := vdl.checkClosed(); err != nil {
		return err
	}
	value, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal application state: %w", err)
	}

	key := []byte(state.ApplicationID)
	versionID := generateVersionID(key, value)

	toUpdate := []storage.KeyValuePair{{Key: key, Value: value}}
	toRemove := [][]byte{}

	err = vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store application state in Versioned LevelDB: %w", err)
	}
	return nil
}

func (vdl *VersionedLevelDBDataLayer) GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error) {
	if err := vdl.checkClosed(); err != nil {
		return nil, err
	}
	key := []byte(applicationID)
	value, err := vdl.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get application state from Versioned LevelDB: %w", err)
	}
	if value == nil {
		return nil, errors.ErrNotFound(fmt.Sprintf("application state not found: %s", applicationID))
	}

	state := &common.ApplicationState{}
	err = json.Unmarshal(value, state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal application state: %w", err)
	}
	return state, nil
}

func (vdl *VersionedLevelDBDataLayer) StoreWASMBytecode(ctx context.Context, applicationID string, bytecode []byte) error {
	if err := vdl.checkClosed(); err != nil {
		return err
	}
	key := []byte(applicationID)
	versionID := generateVersionID(key, bytecode)

	toUpdate := []storage.KeyValuePair{{Key: key, Value: bytecode}}
	toRemove := [][]byte{}

	err := vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store WASM bytecode in Versioned LevelDB: %w", err)
	}
	return nil
}

func (vdl *VersionedLevelDBDataLayer) GetWASMBytecode(ctx context.Context, applicationID string) ([]byte, error) {
	if err := vdl.checkClosed(); err != nil {
		return nil, err
	}
	key := []byte(applicationID)
	value, err := vdl.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get WASM bytecode from Versioned LevelDB: %w", err)
	}
	if value == nil {
		return nil, errors.ErrNotFound("wasm bytecode not found for application: " + applicationID)
	}
	return value, nil
}

func (vdl *VersionedLevelDBDataLayer) StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	if err := vdl.checkClosed(); err != nil {
		return err
	}
	value, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal deanonymization report: %w", err)
	}

	key := []byte(report.ReportID)
	versionID := generateVersionID(key, value)

	toUpdate := []storage.KeyValuePair{{Key: key, Value: value}}
	toRemove := [][]byte{}

	err = vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store deanonymization report in Versioned LevelDB: %w", err)
	}
	return nil
}

func (vdl *VersionedLevelDBDataLayer) GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error) {
	if err := vdl.checkClosed(); err != nil {
		return nil, err
	}
	key := []byte(reportID)
	value, err := vdl.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get deanonymization report from Versioned LevelDB: %w", err)
	}
	if value == nil {
		return nil, errors.ErrNotFound("deanonymization report not found: " + reportID)
	}

	report := &common.DeanonymizationReport{}
	err = json.Unmarshal(value, report)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal deanonymization report: %w", err)
	}
	return report, nil
}

func (vdl *VersionedLevelDBDataLayer) StoreUserKey(ctx context.Context, userID string, publicKey []byte) error {
	if err := vdl.checkClosed(); err != nil {
		return err
	}
	key := []byte(userID)
	versionID := generateVersionID(key, publicKey)

	toUpdate := []storage.KeyValuePair{{Key: key, Value: publicKey}}
	toRemove := [][]byte{}

	err := vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store user key in Versioned LevelDB: %w", err)
	}
	return nil
}

func (vdl *VersionedLevelDBDataLayer) GetUserKey(ctx context.Context, userID string) ([]byte, error) {
	if err := vdl.checkClosed(); err != nil {
		return nil, err
	}
	key := []byte(userID)
	value, err := vdl.adapter.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get user key from Versioned LevelDB: %w", err)
	}
	if value == nil {
		return nil, errors.ErrNotFound("public key not found for user: " + userID)
	}
	return value, nil
}

func (vdl *VersionedLevelDBDataLayer) Close() error {
	vdl.mu.Lock()
	defer vdl.mu.Unlock()
	if vdl.isClosed {
		return nil
	}
	vdl.isClosed = true
	return vdl.adapter.Close()
}

var _ storage.ApplicationStateStore = (*VersionedLevelDBDataLayer)(nil)

func (vdl *VersionedLevelDBDataLayer) NumberOfVersions() int {
	return vdl.adapter.NumberOfVersions()
}

func (vdl *VersionedLevelDBDataLayer) Rollback(versionID []byte) error {
	if err := vdl.checkClosed(); err != nil {
		return err
	}
	return vdl.adapter.Rollback(versionID)
}

func (vdl *VersionedLevelDBDataLayer) LastVersionID() ([]byte, error) {
	if err := vdl.checkClosed(); err != nil {
		return nil, err
	}
	return vdl.adapter.LastVersionID()
}
