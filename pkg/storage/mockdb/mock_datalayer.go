package mockdb

import (
	"context"
	"fmt"
	"sync"

	"errors"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
)

// MockDataLayer is a mock implementation of the data layer for testing.
// It is safe for concurrent use.
type MockDataLayer struct {
	mutex     sync.RWMutex
	states    map[string]*common.ApplicationState
	bytecodes map[string][]byte
	reports   map[string]*common.DeanonymizationReport
	keys      map[string][]byte
	isClosed  bool
	versions [][]byte
	*testutil.MockFunctions
}

// NewMockDataLayer creates a new mock data layer.
func NewMockDataLayer() *MockDataLayer {

	return &MockDataLayer{
		states:    make(map[string]*common.ApplicationState),
		bytecodes: make(map[string][]byte),
		reports:   make(map[string]*common.DeanonymizationReport),
		keys:      make(map[string][]byte),
		versions:  make([][]byte, 0),
		MockFunctions: testutil.NewMockFunctions(),
	}
}

// checkClosed returns an error if the mock data layer is closed.
func (d *MockDataLayer) checkClosed() error {
	if d.isClosed {
		return storageErrors.ErrStorageIsClosed("mock data layer is closed")
	}
	return nil
}

// Store stores the state of an application.
// versionID is ignored in current mock implementation
func (d *MockDataLayer) Store(
	ctx context.Context,
	versionID []byte,
	stateArray []*common.ApplicationState,
	wasmArray []*common.WASMData,
) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if f, ok:= d.GetMockedFunc("Store"); ok {
		return f.(func(context.Context, []byte, []*common.ApplicationState, []*common.WASMData) (error))(ctx, versionID, stateArray, wasmArray)
	}
	if err := d.checkClosed(); err != nil {
		return err
	}

	for _, state := range stateArray {
		if state != nil {
			d.states[state.ApplicationID] = state
		}
	}

	for _, wasm := range wasmArray {
		if wasm != nil {
			d.bytecodes[wasm.ApplicationID] = wasm.Bytecode
		}
	}


	d.versions = append(d.versions, versionID)

	return nil
}

// GetApplicationState retrieves the state of an application.
func (d *MockDataLayer) GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok:= d.GetMockedFunc("GetApplicationState"); ok {
		return f.(func(context.Context, string) (*common.ApplicationState, error))(ctx, applicationID)
	}

	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	state, exists := d.states[applicationID]
	if !exists {
		return nil, storageErrors.ErrNotFound("application state not found: " + applicationID)
	}
	return state, nil
}

// GetWASMBytecode retrieves WASM bytecode for an application.
func (d *MockDataLayer) GetWASMBytecode(ctx context.Context, applicationID string) ([]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok:= d.GetMockedFunc("GetWASMBytecode"); ok {
		return f.(func(context.Context, string) ([]byte, error))(ctx, applicationID)
	}
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	bytecode, exists := d.bytecodes[applicationID]
	if !exists {
		return nil, storageErrors.ErrNotFound("wasm bytecode not found for application: " + applicationID)
	}
	return bytecode, nil
}

// StoreDeanonymizationReport stores a deanonymization report.
func (d *MockDataLayer) StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if f, ok:= d.GetMockedFunc("StoreDeanonymizationReport"); ok {
		return f.(func(context.Context, *common.DeanonymizationReport) (error))(ctx, report)
	}

	if err := d.checkClosed(); err != nil {
		return err
	}
	d.reports[report.ReportID] = report
	return nil
}

// GetDeanonymizationReport retrieves a deanonymization report.
func (d *MockDataLayer) GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	report, exists := d.reports[reportID]
	if !exists {
		return nil, storageErrors.ErrNotFound("deanonymization report not found: " + reportID)
	}
	return report, nil
}

// Close marks the mock data layer as closed.
func (d *MockDataLayer) Close() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if f, ok:= d.GetMockedFunc("Close"); ok {
		return f.(func() ( error))()
	}

	d.isClosed = true
	return nil
}

// Rollback is a mock implementation of the Rollback method.
// versionID is ignored in current mock implementation
func (d *MockDataLayer) Rollback(versionID []byte) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok:= d.GetMockedFunc("Rollback"); ok {
		return f.(func([]byte) (error))(versionID)
	}

	var initialState = [32]byte{}
	if string(versionID) == string(initialState[:]) {
		d.versions = make([][]byte, 0)
		d.states =   make(map[string]*common.ApplicationState)
		d.bytecodes = make(map[string][]byte)
		return nil
	}


	for i, v := range d.versions {
		if string(v) == string(versionID) {
			d.versions = d.versions[:i+1]
			for _, v := range d.states {
				copy(v.StateRoot[:32], versionID)
			}

			return nil
		}		
	}
	return fmt.Errorf("versionID not found: %x", versionID)
}

// LastVersionID is a mock implementation of the LastVersionID method.
func (d *MockDataLayer) LastVersionID() ([]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok:= d.GetMockedFunc("LastVersionID"); ok {
		return f.(func() ([]byte, error))()
	}

	if len(d.versions) == 0 {
		return nil, errors.New("no versions found in the db")
	}
	
	return d.versions[len(d.versions) - 1], nil
}

// ListVersions is a mock implementation of the ListVersions method.
func (d *MockDataLayer) ListVersions() ([][]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok:= d.GetMockedFunc("ListVersions"); ok {
		return f.(func() ([][]byte, error))()
	}
	
	lifoVersions := make([][]byte, len(d.versions))
	for i, v := range d.versions {
		lifoVersions[len(d.versions)-1-i] = v
	}
	return lifoVersions, nil
}

var _ storage.ApplicationStateStore = (*MockDataLayer)(nil)
var _ storage.ApplicationReportStore = (*MockDataLayer)(nil)

var _ storage.DataLayer = (*MockDataLayer)(nil)
