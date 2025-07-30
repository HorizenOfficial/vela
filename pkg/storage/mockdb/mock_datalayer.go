package mockdb

import (
	"context"
	"sync"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/errors"
)

// MockDataLayer is a mock implementation of the data layer for testing.
// It is safe for concurrent use.
type MockDataLayer struct {
	mu       sync.RWMutex
	states   map[string]*common.ApplicationState
	bytecode map[string][]byte
	reports  map[string]*common.DeanonymizationReport
	keys     map[string][]byte
	isClosed bool
}

// NewMockDataLayer creates a new mock data layer.
func NewMockDataLayer() *MockDataLayer {
	return &MockDataLayer{
		states:   make(map[string]*common.ApplicationState),
		bytecode: make(map[string][]byte),
		reports:  make(map[string]*common.DeanonymizationReport),
		keys:     make(map[string][]byte),
	}
}

// checkClosed returns an error if the mock data layer is closed.
func (d *MockDataLayer) checkClosed() error {
	if d.isClosed {
		return errors.ErrStorageIsClosed("mock data layer is closed")
	}
	return nil
}

// Store stores the state of an application.
func (d *MockDataLayer) Store(
	ctx context.Context,
	versionID []byte,
	stateArray []*common.ApplicationState,
	wasmArray []*common.WASMData,
) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.checkClosed(); err != nil {
		return err
	}

	for _, state := range stateArray {
		d.states[state.ApplicationID] = state
	}

	for _, wasm := range wasmArray {
		d.bytecode[wasm.ApplicationID] = wasm.Bytecode
	}

	return nil
}

// GetApplicationState retrieves the state of an application.
func (d *MockDataLayer) GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	state, exists := d.states[applicationID]
	if !exists {
		return nil, errors.ErrNotFound("application state not found: " + applicationID)
	}
	return state, nil
}

// GetWASMBytecode retrieves WASM bytecode for an application.
func (d *MockDataLayer) GetWASMBytecode(ctx context.Context, applicationID string) ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	bytecode, exists := d.bytecode[applicationID]
	if !exists {
		return nil, errors.ErrNotFound("wasm bytecode not found for application: " + applicationID)
	}
	return bytecode, nil
}

// StoreDeanonymizationReport stores a deanonymization report.
func (d *MockDataLayer) StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.checkClosed(); err != nil {
		return err
	}
	d.reports[report.ReportID] = report
	return nil
}

// GetDeanonymizationReport retrieves a deanonymization report.
func (d *MockDataLayer) GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	report, exists := d.reports[reportID]
	if !exists {
		return nil, errors.ErrNotFound("deanonymization report not found: " + reportID)
	}
	return report, nil
}

// StoreUserKey stores a user's public key.
func (d *MockDataLayer) StoreUserKey(ctx context.Context, userID string, publicKey []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.checkClosed(); err != nil {
		return err
	}
	d.keys[userID] = publicKey
	return nil
}

// GetUserKey retrieves a user's public key.
func (d *MockDataLayer) GetUserKey(ctx context.Context, userID string) ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	publicKey, exists := d.keys[userID]
	if !exists {
		return nil, errors.ErrNotFound("public key not found for user: " + userID)
	}
	return publicKey, nil
}

// Close marks the mock data layer as closed.
func (d *MockDataLayer) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.isClosed = true
	return nil
}

// Rollback is a mock implementation of the Rollback method.
func (d *MockDataLayer) Rollback(versionID []byte) error {
	return nil
}

// LastVersionID is a mock implementation of the LastVersionID method.
func (d *MockDataLayer) LastVersionID() ([]byte, error) {
	return []byte("mock_version_id"), nil
}

var _ storage.ApplicationStateStore = (*MockDataLayer)(nil)
var _ storage.ApplicationUserKeyStore = (*MockDataLayer)(nil)
var _ storage.ApplicationReportStore = (*MockDataLayer)(nil)

var _ storage.DataLayer = (*MockDataLayer)(nil)
