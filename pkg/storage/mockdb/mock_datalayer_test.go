package mockdb_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/HorizenOfficial/vela/pkg/common"
	storageErrors "github.com/HorizenOfficial/vela/pkg/storage/errors"
	"github.com/HorizenOfficial/vela/pkg/storage/mockdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMockDataLayer verifies that the constructor for MockDataLayer returns a non-nil instance.
func TestNewMockDataLayer(t *testing.T) {
	dl := mockdb.NewMockDataLayer()
	assert.NotNil(t, dl, "NewMockDataLayer should return a non-nil instance")
}

// TestApplicationStateStore tests the full functionality of the MockDataLayer,
// covering all the methods of the ApplicationStateStore
// and ApplicationReportStore interfaces. It checks for correct data storage and
// retrieval, error handling for non-existent data, behavior after closing the store,
// and concurrent access.
func TestApplicationStateStore(t *testing.T) {
	createStore := func() *mockdb.MockDataLayer {
		return mockdb.NewMockDataLayer()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("StoreAndGet", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()

		expectedState := common.ApplicationState{
			ApplicationID:  common.NewApplicationId(1),
			StateRoot:      sha256.Sum256([]byte("some-test-root-hash-1")),
			EncryptedState: []byte{0x0A, 0x0B, 0x0C, 0x0D},
		}
		expectedBytecode := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03}

		err := store.Store(
			ctx,
			[]byte("version-1"),
			[]*common.ApplicationState{&expectedState},
			[]*common.WASMData{{ApplicationID: expectedState.ApplicationID, Bytecode: expectedBytecode}},
		)
		require.NoError(t, err, "Store should not error")

		actualState, err := store.GetApplicationState(ctx, expectedState.ApplicationID)
		require.NoError(t, err, "GetApplicationState for existing ID should not error")
		require.NotNil(t, actualState, "GetApplicationState should return a non-nil state")
		if diff := cmp.Diff(&expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}

		actualBytecode, err := store.GetWASMBytecode(ctx, expectedState.ApplicationID)
		require.NoError(t, err, "GetWASMBytecode for existing ID should not error")
		require.NotNil(t, actualBytecode, "GetWASMBytecode should return non-nil bytecode")
		assert.True(t, bytes.Equal(expectedBytecode, actualBytecode), "Retrieved WASM bytecode mismatch")
	})

	t.Run("GetNonExistentApplicationState", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		_, err := store.GetApplicationState(ctx, 45)
		require.Error(t, err, "Expected an error when getting a non-existent state")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("GetNonExistentWASMBytecode", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		_, err := store.GetWASMBytecode(ctx, 653)
		require.Error(t, err, "Expected an error when getting non-existent WASM bytecode")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("CloseStore", func(t *testing.T) {
		store := createStore()
		err := store.Close()
		assert.NoError(t, err, "Closing the store should not error")
	})

	t.Run("OperationsAfterClose", func(t *testing.T) {
		store := createStore()
		require.NoError(t, store.Close(), "Closing the store should not return an error on first close")
		err := store.Close()
		assert.NoError(t, err, "Error is not expected when trying to close an already closed store")
		any_id := common.NewApplicationId(999)

		operations := map[string]func() error{
			"GetApplicationState": func() error {
				_, err := store.GetApplicationState(ctx, any_id)
				return err
			},
			"Store": func() error {
				return store.Store(ctx, []byte("version-1"), nil, nil)
			},
			"GetWASMBytecode": func() error {
				_, err := store.GetWASMBytecode(ctx, any_id)
				return err
			},
			"StoreEnclaveKeySetRecovery": func() error {
				return store.StoreEnclaveKeySetRecovery(ctx, &common.EnclaveKeySetRecovery{})
			},
			"GetEnclaveKeySetRecovery": func() error {
				_, err := store.GetEnclaveKeySetRecovery(ctx)
				return err
			},
			"Rollback": func() error {
				return store.Rollback(any_id, []byte("some-version"))
			},
			"LastVersionID": func() error {
				_, err := store.LastVersionID(any_id)
				return err
			},
			"ListVersions": func() error {
				_, err := store.ListVersions(any_id)
				return err
			},
		}

		for name, op := range operations {
			t.Run(name, func(t *testing.T) {
				err := op()
				require.Error(t, err, "Expected an error from a closed store")
				var closedErr *storageErrors.Error
				if assert.True(t, errors.As(err, &closedErr), "Error should be a storage error") {
					assert.Equal(t, storageErrors.StorageIsClosed, closedErr.Code, "Error code should be StorageIsClosed")
				}
			})
		}
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()

		var wg sync.WaitGroup
		numGoroutines := 50

		// Test concurrent writes
		for i := range numGoroutines {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				appID := common.NewApplicationId(uint64(i))
				state := common.ApplicationState{
					ApplicationID: appID,
					StateRoot:     sha256.Sum256([]byte(fmt.Sprintf("root-%d", i))),
				}
				err := store.Store(ctx, []byte(fmt.Sprintf("version-%d", i)), []*common.ApplicationState{&state}, nil)
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()

		// Test concurrent reads
		for i := range numGoroutines {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				appID := common.NewApplicationId(uint64(i))
				state, err := store.GetApplicationState(ctx, appID)
				assert.NoError(t, err)
				assert.Equal(t, appID, state.ApplicationID)
				assert.Equal(t, sha256.Sum256([]byte(fmt.Sprintf("root-%d", i))), state.StateRoot)
			}(i)
		}
		wg.Wait()
	})

	t.Run("StoreAndGetEnclaveKeySetRecovery", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		expectedRecoveryData := &common.EnclaveKeySetRecovery{
			RecoveryType:       common.RecoveryTypeKMS,
			KeySetCiphertext:   []byte{0x01, 0x02, 0x03},
			RecoveryCiphertext: []byte{0x04, 0x05, 0x06},
		}
		err := store.StoreEnclaveKeySetRecovery(ctx, expectedRecoveryData)
		require.NoError(t, err, "StoreEnclaveKeySetRecovery should not error")
		actualRecoveryData, err := store.GetEnclaveKeySetRecovery(ctx)
		require.NoError(t, err, "GetEnclaveKeySetRecovery for existing data should not error")
		require.NotNil(t, actualRecoveryData, "GetEnclaveKeySetRecovery should return a non-nil data")
		if diff := cmp.Diff(expectedRecoveryData, actualRecoveryData); diff != "" {
			t.Errorf("Retrieved EnclaveKeySetRecovery mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentEnclaveKeySetRecovery", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		_, err := store.GetEnclaveKeySetRecovery(ctx)
		require.Error(t, err, "Expected an error when getting non-existent recovery data")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreArrays", func(t *testing.T) {
		ctx := context.Background()
		dataLayer := mockdb.NewMockDataLayer()

		// Store per-app data in separate calls (Store enforces single-app-per-call).
		state1 := &common.ApplicationState{
			ApplicationID:  common.NewApplicationId(1),
			StateRoot:      sha256.Sum256([]byte("root1")),
			EncryptedState: []byte("state1"),
		}
		wasm1 := &common.WASMData{
			ApplicationID: common.NewApplicationId(1),
			Bytecode:      []byte("wasm1"),
		}
		err := dataLayer.Store(ctx, []byte("version-1"), []*common.ApplicationState{state1}, []*common.WASMData{wasm1})
		require.NoError(t, err)

		state2 := &common.ApplicationState{
			ApplicationID:  common.NewApplicationId(2),
			StateRoot:      sha256.Sum256([]byte("root2")),
			EncryptedState: []byte("state2"),
		}
		err = dataLayer.Store(ctx, []byte("version-2"), []*common.ApplicationState{state2}, nil)
		require.NoError(t, err)

		wasm3 := &common.WASMData{
			ApplicationID: common.NewApplicationId(3),
			Bytecode:      []byte("wasm3"),
		}
		err = dataLayer.Store(ctx, []byte("version-3"), nil, []*common.WASMData{wasm3})
		require.NoError(t, err)

		// Verify ApplicationStates were stored correctly
		actualState1, err := dataLayer.GetApplicationState(ctx, common.NewApplicationId(1))
		require.NoError(t, err)
		assert.Equal(t, state1, actualState1)

		actualState2, err := dataLayer.GetApplicationState(ctx, common.NewApplicationId(2))
		require.NoError(t, err)
		assert.Equal(t, state2, actualState2)

		// Verify WASMData were stored correctly
		actualWasm1, err := dataLayer.GetWASMBytecode(ctx, common.NewApplicationId(1))
		require.NoError(t, err)
		assert.Equal(t, wasm1.Bytecode, actualWasm1)

		actualWasm3, err := dataLayer.GetWASMBytecode(ctx, common.NewApplicationId(3))
		require.NoError(t, err)
		assert.Equal(t, wasm3.Bytecode, actualWasm3)

		// App with only WASM data has no state
		_, err = dataLayer.GetApplicationState(ctx, 3)
		assert.Error(t, err)

		// App with only state data has no WASM
		_, err = dataLayer.GetWASMBytecode(ctx, 2)
		assert.Error(t, err)

		// Mixed-app Store calls are rejected
		mixedState := []*common.ApplicationState{state1, state2}
		err = dataLayer.Store(ctx, []byte("version-bad"), mixedState, nil)
		assert.Error(t, err, "Store should reject mixed ApplicationIDs")
	})
}
