// Package storage defines the interfaces for persisting and retrieving application data,
// including application state and user keys. It provides
// an abstraction layer for the underlying storage implementation.
package storage

import (
	"context"

	"github.com/HorizenOfficial/vela/pkg/common"
)

// ApplicationStateStore defines the interface for managing versioned application state,
// including WebAssembly (WASM) bytecode and application data. It supports atomic
// storage of application data, rollbacks, and retrieval of specific versions.
// Versioning is per-application: each application maintains its own independent
// version chain, enabling rollback for one app without affecting others.
type ApplicationStateStore interface {
	// Store atomically saves the application state and WASM bytecode for a given version.
	// The version is filed under the application that produced the data.
	// Precondition: all items in stateArray and wasmArray must share the same ApplicationID
	// (sequential processing — one app per Store call).
	Store(
		ctx context.Context,
		versionID []byte,
		stateArray []*common.ApplicationState,
		wasmArray []*common.WASMData,
	) error

	// Rollback reverts the application state for the given app to the specified versionID.
	// Only ChangeSets belonging to this app are reverted.
	Rollback(appID common.ApplicationIdType, versionID []byte) error

	// LastVersionID returns the most recent version ID for the given application.
	LastVersionID(appID common.ApplicationIdType) ([]byte, error)

	// ListVersions returns a list of all stored version IDs for the given application,
	// ordered from most recent to oldest (LIFO).
	ListVersions(appID common.ApplicationIdType) ([][]byte, error)

	// GetApplicationState retrieves the state of a specific application by its ID.
	GetApplicationState(ctx context.Context, applicationID common.ApplicationIdType) (*common.ApplicationState, error)
	// GetWASMBytecode retrieves the WASM bytecode for a specific application by its ID.
	GetWASMBytecode(ctx context.Context, applicationID common.ApplicationIdType) ([]byte, error)

	// Close releases any resources held by the data store.
	Close() error
}

// EnclaveKeyStore defines the interface for managing enclave keys.
// This storage is not versioned.
type EnclaveKeyStore interface {
	// StoreEnclaveKeySetRecovery saves the enclave key set recovery data.
	StoreEnclaveKeySetRecovery(ctx context.Context, recoveryData *common.EnclaveKeySetRecovery) error
	// GetEnclaveKeySetRecovery retrieves the enclave key set recovery data.
	GetEnclaveKeySetRecovery(ctx context.Context) (*common.EnclaveKeySetRecovery, error)
}

// DataLayer is a composite interface that combines all the application's storage interfaces.
// It provides a single point of access to all data storage functionality.
type DataLayer interface {
	ApplicationStateStore
	EnclaveKeyStore
}
