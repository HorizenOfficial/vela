package mockdb_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockDataLayer(t *testing.T) {
	dl := mockdb.NewMockDataLayer()
	assert.NotNil(t, dl, "NewMockDataLayer should return a non-nil instance")
}

func TestApplicationStateStore(t *testing.T) {
	createStore := func() storage.ApplicationStateStoreOld {
		return mockdb.NewMockDataLayer()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("StoreAndGetApplicationState", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		expectedState := &common.ApplicationState{
			ApplicationID:  "app-test-id-1",
			StateRoot:      []byte("some-test-root-hash-1"),
			EncryptedState: []byte{0x0A, 0x0B, 0x0C, 0x0D},
		}
		err := store.StoreApplicationState(ctx, expectedState)
		require.NoError(t, err, "StoreApplicationState should not error")
		actualState, err := store.GetApplicationState(ctx, expectedState.ApplicationID)
		require.NoError(t, err, "GetApplicationState for existing ID should not error")
		require.NotNil(t, actualState, "GetApplicationState should return a non-nil state")
		if diff := cmp.Diff(expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentApplicationState", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		_, err := store.GetApplicationState(ctx, "non-existent-app-id")
		require.Error(t, err, "Expected an error when getting a non-existent state")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetWASMBytecode", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		appID := "wasm-app-id-1"
		expectedBytecode := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03}
		err := store.StoreWASMBytecode(ctx, appID, expectedBytecode)
		require.NoError(t, err, "StoreWASMBytecode should not error")
		actualBytecode, err := store.GetWASMBytecode(ctx, appID)
		require.NoError(t, err, "GetWASMBytecode for existing ID should not error")
		require.NotNil(t, actualBytecode, "GetWASMBytecode should return non-nil bytecode")
		assert.True(t, bytes.Equal(expectedBytecode, actualBytecode), "Retrieved WASM bytecode mismatch")
	})

	t.Run("GetNonExistentWASMBytecode", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		_, err := store.GetWASMBytecode(ctx, "non-existent-wasm-id")
		require.Error(t, err, "Expected an error when getting non-existent WASM bytecode")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetDeanonymizationReport", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		expectedReport := &common.DeanonymizationReport{
			ApplicationID:   "deanon-1",
			ReportID:        "report-id-1",
			EncryptedReport: []byte("some-test-root-hash-1"),
		}
		err := store.StoreDeanonymizationReport(ctx, expectedReport)
		require.NoError(t, err, "StoreDeanonymizationReport should not error")
		actualReport, err := store.GetDeanonymizationReport(ctx, expectedReport.ReportID)
		require.NoError(t, err, "GetDeanonymizationReport for existing ID should not error")
		require.NotNil(t, actualReport, "GetDeanonymizationReport should return a non-nil report")
		if diff := cmp.Diff(expectedReport, actualReport); diff != "" {
			t.Errorf("Retrieved DeanonymizationReport mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentDeanonymizationReport", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		_, err := store.GetDeanonymizationReport(ctx, "non-existent-report-id")
		require.Error(t, err, "Expected an error when getting non-existent deanonymization report")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetUserKey", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		userID := "user-id-1"
		expectedPublicKey := []byte("some-public-key-bytes-1")
		err := store.StoreUserKey(ctx, userID, expectedPublicKey)
		require.NoError(t, err, "StoreUserKey should not error")
		actualPublicKey, err := store.GetUserKey(ctx, userID)
		require.NoError(t, err, "GetUserKey for existing ID should not error")
		require.NotNil(t, actualPublicKey, "GetUserKey should return non-nil public key")
		assert.True(t, bytes.Equal(expectedPublicKey, actualPublicKey), "Retrieved User Key mismatch")
	})

	t.Run("GetNonExistentUserKey", func(t *testing.T) {
		store := createStore()
		defer func() { require.NoError(t, store.Close(), "Store.Close() should not error") }()
		_, err := store.GetUserKey(ctx, "non-existent-user-id")
		require.Error(t, err, "Expected an error when getting non-existent user key")
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

		operations := map[string]func() error{
			"GetApplicationState": func() error {
				_, err := store.GetApplicationState(ctx, "any-id")
				return err
			},
			"StoreApplicationState": func() error {
				return store.StoreApplicationState(ctx, &common.ApplicationState{ApplicationID: "test-after-close"})
			},
			"GetWASMBytecode": func() error {
				_, err := store.GetWASMBytecode(ctx, "any-id")
				return err
			},
			"StoreWASMBytecode": func() error {
				return store.StoreWASMBytecode(ctx, "any-id", []byte("bytecode"))
			},
			"GetDeanonymizationReport": func() error {
				_, err := store.GetDeanonymizationReport(ctx, "any-id")
				return err
			},
			"StoreDeanonymizationReport": func() error {
				return store.StoreDeanonymizationReport(ctx, &common.DeanonymizationReport{ReportID: "test"})
			},
			"GetUserKey": func() error {
				_, err := store.GetUserKey(ctx, "any-id")
				return err
			},
			"StoreUserKey": func() error {
				return store.StoreUserKey(ctx, "test-user", []byte("key"))
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
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				appID := fmt.Sprintf("concurrent-app-%d", i)
				state := &common.ApplicationState{
					ApplicationID: appID,
					StateRoot:     []byte(fmt.Sprintf("root-%d", i)),
				}
				err := store.StoreApplicationState(ctx, state)
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()

		// Test concurrent reads
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				appID := fmt.Sprintf("concurrent-app-%d", i)
				state, err := store.GetApplicationState(ctx, appID)
				assert.NoError(t, err)
				assert.Equal(t, appID, state.ApplicationID)
				assert.Equal(t, []byte(fmt.Sprintf("root-%d", i)), state.StateRoot)
			}(i)
		}
		wg.Wait()
	})
}
