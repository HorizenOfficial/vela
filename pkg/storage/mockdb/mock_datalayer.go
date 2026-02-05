package mockdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
)

// MockDataLayer is a mock implementation of the data layer for testing.
// It is safe for concurrent use.
type MockDataLayer struct {
	mutex              sync.RWMutex
	states             map[common.ApplicationIdType]*common.ApplicationState
	bytecodes          map[common.ApplicationIdType][]byte
	keys               map[string][]byte
	enclaveKeyRecovery *common.EnclaveKeySetRecovery
	isClosed           bool
	versions           [][]byte
	*testutil.MockFunctions
}

// NewMockDataLayer creates a new mock data layer.
func NewMockDataLayer() *MockDataLayer {

	return &MockDataLayer{
		states:        make(map[common.ApplicationIdType]*common.ApplicationState),
		bytecodes:     make(map[common.ApplicationIdType][]byte),
		keys:          make(map[string][]byte),
		versions:      make([][]byte, 0),
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

	if f, ok := d.GetMockedFunc("Store"); ok {
		return f.(func(context.Context, []byte, []*common.ApplicationState, []*common.WASMData) error)(ctx, versionID, stateArray, wasmArray)
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
func (d *MockDataLayer) GetApplicationState(ctx context.Context, applicationID common.ApplicationIdType) (*common.ApplicationState, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok := d.GetMockedFunc("GetApplicationState"); ok {
		return f.(func(context.Context, common.ApplicationIdType) (*common.ApplicationState, error))(ctx, applicationID)
	}

	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	state, exists := d.states[applicationID]
	if !exists {
		return nil, storageErrors.ErrNotFound("application state not found: " + applicationID.String())
	}
	return state, nil
}

// GetWASMBytecode retrieves WASM bytecode for an application.
func (d *MockDataLayer) GetWASMBytecode(ctx context.Context, applicationID common.ApplicationIdType) ([]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok := d.GetMockedFunc("GetWASMBytecode"); ok {
		return f.(func(context.Context, common.ApplicationIdType) ([]byte, error))(ctx, applicationID)
	}
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	bytecode, exists := d.bytecodes[applicationID]
	if !exists {
		return nil, storageErrors.ErrNotFound("wasm bytecode not found for application: " + applicationID.String())
	}
	return bytecode, nil
}

// StoreEnclaveKeySetRecovery stores the enclave key set recovery data.
func (d *MockDataLayer) StoreEnclaveKeySetRecovery(ctx context.Context, recoveryData *common.EnclaveKeySetRecovery) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if f, ok := d.GetMockedFunc("StoreEnclaveKeySetRecovery"); ok {
		return f.(func(context.Context, *common.EnclaveKeySetRecovery) error)(ctx, recoveryData)
	}

	if err := d.checkClosed(); err != nil {
		return err
	}
	d.enclaveKeyRecovery = recoveryData
	return nil
}

// GetEnclaveKeySetRecovery retrieves the enclave key set recovery data.
func (d *MockDataLayer) GetEnclaveKeySetRecovery(ctx context.Context) (*common.EnclaveKeySetRecovery, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok := d.GetMockedFunc("GetEnclaveKeySetRecovery"); ok {
		return f.(func(context.Context) (*common.EnclaveKeySetRecovery, error))(ctx)
	}

	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	if d.enclaveKeyRecovery == nil {
		return nil, storageErrors.ErrNotFound("enclave key set recovery data not found")
	}
	return d.enclaveKeyRecovery, nil
}

// Close marks the mock data layer as closed.
func (d *MockDataLayer) Close() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if f, ok := d.GetMockedFunc("Close"); ok {
		return f.(func() error)()
	}

	d.isClosed = true
	return nil
}

// Rollback is a mock implementation of the Rollback method.
// For now, only the StateRoot is restored.
func (d *MockDataLayer) Rollback(versionID []byte) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok := d.GetMockedFunc("Rollback"); ok {
		return f.(func([]byte) error)(versionID)
	}

	var initialState = [32]byte{}
	if string(versionID) == string(initialState[:]) {
		d.versions = make([][]byte, 0)
		d.states = make(map[common.ApplicationIdType]*common.ApplicationState)
		d.bytecodes = make(map[common.ApplicationIdType][]byte)
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

	if f, ok := d.GetMockedFunc("LastVersionID"); ok {
		return f.(func() ([]byte, error))()
	}

	if len(d.versions) == 0 {
		return nil, storageErrors.ErrNoVersionInDb("No version in db")
	}

	return d.versions[len(d.versions)-1], nil
}

// ListVersions is a mock implementation of the ListVersions method.
func (d *MockDataLayer) ListVersions() ([][]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok := d.GetMockedFunc("ListVersions"); ok {
		return f.(func() ([][]byte, error))()
	}

	lifoVersions := make([][]byte, len(d.versions))
	for i, v := range d.versions {
		lifoVersions[len(d.versions)-1-i] = v
	}
	return lifoVersions, nil
}

var _ storage.ApplicationStateStore = (*MockDataLayer)(nil)

var _ storage.DataLayer = (*MockDataLayer)(nil)
