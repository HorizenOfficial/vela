// Package storage defines the interfaces for persisting and retrieving application data,
// including application state, user keys, and deanonymization reports. It provides
// an abstraction layer for the underlying storage implementation.
package storage

import (
	"context"

	"github.com/horizen-pes/pkg/common"
)

// ApplicationStateStore defines the interface for managing versioned application state,
// including WebAssembly (WASM) bytecode and application data. It supports atomic
// storage of application data, rollbacks, and retrieval of specific versions.
type ApplicationStateStore interface {
	// Store atomically saves the application state and WASM bytecode for a given version.
	// This allows for reliable versioning and the ability to revert to a previous state if needed.
	Store(
		ctx context.Context,
		versionID []byte,
		stateArray []*common.ApplicationState,
		wasmArray []*common.WASMData,
	) error

	// Rollback reverts the application state to the specified versionID.
	Rollback(versionID []byte) error

	// LastVersionID returns the most recent version ID stored in the database.
	LastVersionID() ([]byte, error)

	// ListVersions returns a list of all stored version IDs.
	ListVersions() ([][]byte, error)

	// GetApplicationState retrieves the state of a specific application by its ID.
	GetApplicationState(ctx context.Context, applicationID common.ApplicationIdType) (*common.ApplicationState, error)
	// GetWASMBytecode retrieves the WASM bytecode for a specific application by its ID.
	GetWASMBytecode(ctx context.Context, applicationID common.ApplicationIdType) ([]byte, error)

	// Close releases any resources held by the data store.
	Close() error
}

// ApplicationReportStore defines the interface for managing deanonymization reports.
// This storage is not versioned
type ApplicationReportStore interface {
	// StoreDeanonymizationReport saves a new deanonymization report.
	StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error
	// GetDeanonymizationReport retrieves a deanonymization report by its ID.
	GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error)
}

// DataLayer is a composite interface that combines all the application's storage interfaces.
// It provides a single point of access to all data storage functionality.
type DataLayer interface {
	ApplicationStateStore
	ApplicationReportStore
}
