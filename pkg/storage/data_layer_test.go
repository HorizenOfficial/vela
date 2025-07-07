package storage_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMockDataLayer simply verifies the mock can be created.
func TestNewMockDataLayer(t *testing.T) {
	dl := storage.NewMockDataLayer()
	assert.NotNil(t, dl, "NewMockDataLayer should return a non-nil instance")
}

func TestApplicationStateStore(t *testing.T) {
	// createStore provides a fresh instance of the store for each test scenario.
	createStore := func() storage.ApplicationStateStore {
		return storage.NewMockDataLayer()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is canceled when tests complete

	// --- Test ApplicationState Operations ---
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

		var notFoundErr *storage.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test WASM Bytecode Operations ---
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

		var notFoundErr *storage.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test DeanonymizationReport Operations ---
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

		var notFoundErr *storage.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test UserKey Operations ---
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

		var notFoundErr *storage.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != "not found" {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	// --- Test Close Method ---
	t.Run("CloseStore", func(t *testing.T) {
		store := createStore()
		// No defer Close() here, as we're testing it explicitly
		err := store.Close()
		assert.NoError(t, err, "Closing the store should not error")
	})
}
