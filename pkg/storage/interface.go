// Package storage provides the interface and mock implementation for the data layer
package storage

import (
	"context"
	"github.com/horizen-pes/pkg/common"
)

// ApplicationStateStore defines the interface for the data layer
type ApplicationStateStore interface {
	// StoreApplicationState stores the state of an application
	StoreApplicationState(ctx context.Context, state *common.ApplicationState) error
	// GetApplicationState gets the state of an application
	GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error)
	// StoreWASMBytecode stores WASM bytecode for an application
	StoreWASMBytecode(ctx context.Context, applicationID string, bytecode []byte) error
	// GetWASMBytecode gets WASM bytecode for an application
	GetWASMBytecode(ctx context.Context, applicationID string) ([]byte, error)
	// StoreDeanonymizationReport stores a deanonymization report
	StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error
	// GetDeanonymizationReport gets a deanonymization report
	GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error)
	// StoreUserKey stores a user's public key
	StoreUserKey(ctx context.Context, userID string, publicKey []byte) error
	// GetUserKey retrieves a user's public key
	GetUserKey(ctx context.Context, userID string) ([]byte, error)
	// Close closes the data layer
	Close() error
}
