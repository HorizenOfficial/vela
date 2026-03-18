package mockdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/storage"
	storageErrors "github.com/HorizenOfficial/vela/pkg/storage/errors"
)

// MockDataLayer is a mock implementation of the data layer for testing.
// It is safe for concurrent use. Versioning is per-application.
type MockDataLayer struct {
	mutex              sync.RWMutex
	states             map[common.ApplicationIdType]*common.ApplicationState
	bytecodes          map[common.ApplicationIdType][]byte
	keys               map[string][]byte
	enclaveKeyRecovery *common.EnclaveKeySetRecovery
	isClosed           bool
	versions           map[common.ApplicationIdType][][]byte
	*testutil.MockFunctions
}

// NewMockDataLayer creates a new mock data layer.
func NewMockDataLayer() *MockDataLayer {

	return &MockDataLayer{
		states:        make(map[common.ApplicationIdType]*common.ApplicationState),
		bytecodes:     make(map[common.ApplicationIdType][]byte),
		keys:          make(map[string][]byte),
		versions:      make(map[common.ApplicationIdType][][]byte),
		MockFunctions: testutil.NewMockFunctions(),
	}
}

// checkClosed returns an error if the mock data layer is closed.
// Caller must hold d.mutex (read or write) before calling.
func (d *MockDataLayer) checkClosed() error {
	if d.isClosed {
		return storageErrors.ErrStorageIsClosed("mock data layer is closed")
	}
	return nil
}

// Store saves the application state. state must not be nil.
func (d *MockDataLayer) Store(
	ctx context.Context,
	versionID []byte,
	state *common.ApplicationState,
) error {
	return d.storeInternal(ctx, versionID, state, nil)
}

// StoreWithWasm atomically saves the application state and WASM bytecode.
// Both state and wasm must not be nil and must share the same ApplicationID.
func (d *MockDataLayer) StoreWithWasm(
	ctx context.Context,
	versionID []byte,
	state *common.ApplicationState,
	wasm *common.WASMData,
) error {
	if wasm == nil {
		return fmt.Errorf("wasm must not be nil; use Store to save state without WASM")
	}
	return d.storeInternal(ctx, versionID, state, wasm)
}

func (d *MockDataLayer) storeInternal(
	ctx context.Context,
	versionID []byte,
	state *common.ApplicationState,
	wasm *common.WASMData,
) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if f, ok := d.GetMockedFunc("Store"); ok {
		return f.(func(context.Context, []byte, *common.ApplicationState, *common.WASMData) error)(ctx, versionID, state, wasm)
	}
	if err := d.checkClosed(); err != nil {
		return err
	}

	if state == nil {
		return fmt.Errorf("state must not be nil")
	}
	if wasm != nil && state.ApplicationID != wasm.ApplicationID {
		return fmt.Errorf("inconsistent ApplicationID: state has %d, wasm has %d", state.ApplicationID, wasm.ApplicationID)
	}

	appID := state.ApplicationID
	d.states[state.ApplicationID] = state

	if wasm != nil {
		d.bytecodes[wasm.ApplicationID] = wasm.Bytecode
	}

	d.versions[appID] = append(d.versions[appID], versionID)

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
// Only reverts version history and state root for the specified app.
func (d *MockDataLayer) Rollback(appID common.ApplicationIdType, versionID []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if f, ok := d.GetMockedFunc("Rollback"); ok {
		return f.(func(common.ApplicationIdType, []byte) error)(appID, versionID)
	}
	if err := d.checkClosed(); err != nil {
		return err
	}

	var initialState = [32]byte{}
	if string(versionID) == string(initialState[:]) {
		d.versions[appID] = make([][]byte, 0)
		delete(d.states, appID)
		delete(d.bytecodes, appID)
		return nil
	}

	appVersions := d.versions[appID]
	for i, v := range appVersions {
		if string(v) == string(versionID) {
			d.versions[appID] = appVersions[:i+1]
			// Only update the target app's state root.
			if state, exists := d.states[appID]; exists {
				copy(state.StateRoot[:32], versionID)
			}
			return nil
		}
	}
	return fmt.Errorf("versionID not found for app %d: %x", appID, versionID)
}

// LastVersionID is a mock implementation of the LastVersionID method.
func (d *MockDataLayer) LastVersionID(appID common.ApplicationIdType) ([]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok := d.GetMockedFunc("LastVersionID"); ok {
		return f.(func(common.ApplicationIdType) ([]byte, error))(appID)
	}
	if err := d.checkClosed(); err != nil {
		return nil, err
	}

	appVersions := d.versions[appID]
	if len(appVersions) == 0 {
		return nil, storageErrors.ErrNoVersionInDb("No version in db")
	}

	return appVersions[len(appVersions)-1], nil
}

// ListVersions is a mock implementation of the ListVersions method.
func (d *MockDataLayer) ListVersions(appID common.ApplicationIdType) ([][]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if f, ok := d.GetMockedFunc("ListVersions"); ok {
		return f.(func(common.ApplicationIdType) ([][]byte, error))(appID)
	}
	if err := d.checkClosed(); err != nil {
		return nil, err
	}

	appVersions := d.versions[appID]
	lifoVersions := make([][]byte, len(appVersions))
	for i, v := range appVersions {
		lifoVersions[len(appVersions)-1-i] = v
	}
	return lifoVersions, nil
}

var _ storage.ApplicationStateStore = (*MockDataLayer)(nil)

var _ storage.DataLayer = (*MockDataLayer)(nil)
