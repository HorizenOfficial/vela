package boltdb_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	berrors "go.etcd.io/bbolt/errors"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/boltdb"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
)

var testBoltDBBaseDir string

func TestMain(m *testing.M) {
	fmt.Println("Running TestMain setup for BoltDBDataLayer integration tests...")
	var err error
	testBoltDBBaseDir, err = os.MkdirTemp("", "boltdb_test_dbs_")
	if err != nil {
		fmt.Printf("Failed to create BoltDB test base directory: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	fmt.Println("Running TestMain teardown for BoltDBDataLayer integration tests...")
	err = os.RemoveAll(testBoltDBBaseDir)
	if err != nil {
		fmt.Printf("Failed to clean up BoltDB test base directory %s: %v\n", testBoltDBBaseDir, err)
	}
	os.Exit(code)
}

func TestBoltDBDataLayer(t *testing.T) {
	createStore := func(t *testing.T) storage.ApplicationStateStore {
		tempDir, err := os.MkdirTemp(testBoltDBBaseDir, "boltdb-test-")
		require.NoError(t, err, "Failed to create temp directory for BoltDB")

		cfg := boltdb.BoltDBConfig{
			Path:    filepath.Join(tempDir, "test.db"),
			Timeout: 1 * time.Second,
		}
		dl, err := boltdb.NewBoltDBDataLayer(cfg)
		require.NoError(t, err, "Failed to create BoltDBDataLayer instance")

		t.Cleanup(func() {
			//if dl != nil {
			//	require.NoError(t, dl.Close(), "Closing BoltDB store should not error during cleanup")
			//}
			// Close the database connection.
			require.NoError(t, dl.Close(), "Closing BoltDB store should not error during cleanup")
			// Remove the temporary directory and all its contents.
			require.NoError(t, os.RemoveAll(tempDir), "Failed to remove BoltDB test directory: %s", tempDir)
		})

		return dl
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("StoreAndGetApplicationState", func(t *testing.T) {
		store := createStore(t)
		expectedState := &common.ApplicationState{
			ApplicationID:  "boltdb-app-id-1",
			StateRoot:      []byte("boltdb-root-hash-1"),
			EncryptedState: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		}
		err := store.StoreApplicationState(ctx, expectedState)
		require.NoError(t, err, "StoreApplicationState should not return an error")
		actualState, err := store.GetApplicationState(ctx, expectedState.ApplicationID)
		require.NoError(t, err, "GetApplicationState for existing ID should not return an error")
		require.NotNil(t, actualState, "GetApplicationState should return a non-nil state")
		if diff := cmp.Diff(expectedState, actualState); diff != "" {
			t.Errorf("Retrieved ApplicationState mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentApplicationState", func(t *testing.T) {
		store := createStore(t)
		_, err := store.GetApplicationState(ctx, "boltdb-non-existent-app-id")
		require.Error(t, err, "Expected an error when getting a non-existent application state")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetWASMBytecode", func(t *testing.T) {
		store := createStore(t)
		appID := "boltdb-wasm-app-id-1"
		expectedBytecode := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04}
		err := store.StoreWASMBytecode(ctx, appID, expectedBytecode)
		require.NoError(t, err, "StoreWASMBytecode should not return an error")
		actualBytecode, err := store.GetWASMBytecode(ctx, appID)
		require.NoError(t, err, "GetWASMBytecode for existing ID should not return an error")
		require.NotNil(t, actualBytecode, "GetWASMBytecode should return non-nil bytecode")
		assert.True(t, bytes.Equal(expectedBytecode, actualBytecode), "Retrieved WASM bytecode mismatch")
	})

	t.Run("GetNonExistentWASMBytecode", func(t *testing.T) {
		store := createStore(t)
		_, err := store.GetWASMBytecode(ctx, "boltdb-non-existent-wasm-id")
		require.Error(t, err, "Expected an error when getting non-existent WASM bytecode")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetDeanonymizationReport", func(t *testing.T) {
		store := createStore(t)
		expectedReport := &common.DeanonymizationReport{
			ApplicationID:   "id-001",
			ReportID:        "boltdb-report-id-1",
			EncryptedReport: []byte("some-test-root-hash-1"),
		}
		err := store.StoreDeanonymizationReport(ctx, expectedReport)
		require.NoError(t, err, "StoreDeanonymizationReport should not return an error")
		actualReport, err := store.GetDeanonymizationReport(ctx, expectedReport.ReportID)
		require.NoError(t, err, "GetDeanonymizationReport for existing ID should not return an error")
		require.NotNil(t, actualReport, "GetDeanonymizationReport should return a non-nil report")
		if diff := cmp.Diff(expectedReport, actualReport); diff != "" {
			t.Errorf("Retrieved DeanonymizationReport mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetNonExistentDeanonymizationReport", func(t *testing.T) {
		store := createStore(t)
		_, err := store.GetDeanonymizationReport(ctx, "boltdb-non-existent-report-id")
		require.Error(t, err, "Expected an error when getting non-existent deanonymization report")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("StoreAndGetUserKey", func(t *testing.T) {
		store := createStore(t)
		userID := "boltdb-user-id-1"
		expectedPublicKey := []byte("boltdb-public-key-bytes-1")
		err := store.StoreUserKey(ctx, userID, expectedPublicKey)
		require.NoError(t, err, "StoreUserKey should not return an error")
		actualPublicKey, err := store.GetUserKey(ctx, userID)
		require.NoError(t, err, "GetUserKey for existing ID should not return an error")
		require.NotNil(t, actualPublicKey, "GetUserKey should return non-nil public key")
		assert.True(t, bytes.Equal(expectedPublicKey, actualPublicKey), "Retrieved User Key mismatch")
	})

	t.Run("GetNonExistentUserKey", func(t *testing.T) {
		store := createStore(t)
		_, err := store.GetUserKey(ctx, "boltdb-non-existent-user-id")
		require.Error(t, err, "Expected an error when getting non-existent user key")
		var notFoundErr *storageErrors.Error
		if !errors.As(err, &notFoundErr) || notFoundErr.Code != storageErrors.NotFound {
			t.Errorf("Expected a 'not found' error, got: %T (%v)", err, err)
		}
	})

	t.Run("OperationsAfterClose", func(t *testing.T) {
		store := createStore(t)
		require.NoError(t, store.Close(), "Closing the BoltDB store should not return an error on first close")
		_ = store.Close()
		require.NoError(t, store.Close(), "Closing an already closed BoltDB store should not return an error")

		// Define a map of operations to test after the store is closed.
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
				require.Error(t, err, "Expected an error from a closed store: %s")
				assert.True(t, errors.Is(err, berrors.ErrDatabaseNotOpen), "Error should be ErrDatabaseNotOpen")
			})
		}
	})
}
