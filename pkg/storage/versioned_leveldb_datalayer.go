package storage

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage/errors"
	"github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

// VersionedLevelDBDataLayer is an implementation of the ApplicationStateStore interface using Versioned LevelDB.
type VersionedLevelDBDataLayer struct {
	adapter *versioned_leveldb.VersionedLevelDbStorageAdapter
	isClosed bool // to track if the data layer is closed
}

// VersionedLevelDBConfig holds configuration parameters for opening a Versioned LevelDB database.
type VersionedLevelDBConfig struct {
	Path           string // Filesystem path for the database file (e.g., "data/my_versioned.db")
	VersionsToKeep int    // Number of versions to keep
}

// NewVersionedLevelDBDataLayer creates a new VersionedLevelDBDataLayer instance.
func NewVersionedLevelDBDataLayer(cfg VersionedLevelDBConfig) (*VersionedLevelDBDataLayer, error) {
	adapter, err := versioned_leveldb.NewVersionedLevelDbStorageAdapterWithVersions(cfg.Path, cfg.VersionsToKeep)
	if err != nil {
		return nil, fmt.Errorf("failed to create VersionedLevelDbStorageAdapter: %w", err)
	}
	return &VersionedLevelDBDataLayer{adapter: adapter, isClosed: false}, nil
}

// generateVersionID generates a unique version ID based on the key and data.
func generateVersionID(key []byte, data []byte) []byte {
	h := sha256.New()
	h.Write(key)
	h.Write(data)
	return h.Sum(nil)
}

// checkClosed checks if the data layer is closed and returns an error if it is.
func (vdl *VersionedLevelDBDataLayer) checkClosed() error {
	if vdl.isClosed {
		return errors.ErrStorageIsClosed("Versioned LevelDB data layer is closed")
	}
	return nil
}

// StoreApplicationState stores the state of an application.
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

	toUpdate := []common.KeyValuePair{{Key: key, Value: value}}
	toRemove := [][]byte{} // No removals for a store operation

	err = vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store application state in Versioned LevelDB: %w", err)
	}
	return nil
}

// GetApplicationState gets the state of an application.
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

// StoreWASMBytecode stores WASM bytecode for an application.
func (vdl *VersionedLevelDBDataLayer) StoreWASMBytecode(ctx context.Context, applicationID string, bytecode []byte) error {
	if err := vdl.checkClosed(); err != nil {
		return err
	}
	key := []byte(applicationID)
	versionID := generateVersionID(key, bytecode)

	toUpdate := []common.KeyValuePair{{Key: key, Value: bytecode}}
	toRemove := [][]byte{}

	err := vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store WASM bytecode in Versioned LevelDB: %w", err)
	}
	return nil
}

// GetWASMBytecode gets WASM bytecode for an application.
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

// StoreDeanonymizationReport stores a deanonymization report.
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

	toUpdate := []common.KeyValuePair{{Key: key, Value: value}}
	toRemove := [][]byte{}

	err = vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store deanonymization report in Versioned LevelDB: %w", err)
	}
	return nil
}

// GetDeanonymizationReport gets a deanonymization report.
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

// StoreUserKey stores a user's public key.
func (vdl *VersionedLevelDBDataLayer) StoreUserKey(ctx context.Context, userID string, publicKey []byte) error {
	if err := vdl.checkClosed(); err != nil {
		return err
	}
	key := []byte(userID)
	versionID := generateVersionID(key, publicKey)

	toUpdate := []common.KeyValuePair{{Key: key, Value: publicKey}}
	toRemove := [][]byte{}

	err := vdl.adapter.Update(versionID, toUpdate, toRemove)
	if err != nil {
		return fmt.Errorf("failed to store user key in Versioned LevelDB: %w", err)
	}
	return nil
}

// GetUserKey retrieves a user's public key.
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

// Close closes the Versioned LevelDB database connection.
func (vdl *VersionedLevelDBDataLayer) Close() error {
	if vdl.isClosed {
		return nil // Already closed, no-op
	}
	vdl.isClosed = true
	return vdl.adapter.Close()
}