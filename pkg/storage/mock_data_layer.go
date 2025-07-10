// Package storage provides the interface and mock implementation for the data layer
package storage

import (
	"context"

	"github.com/horizen-pes/pkg/common"
)

// MockDataLayer is a mock implementation of the data layer for testing
type MockDataLayer struct {
	states   map[string]*common.ApplicationState
	bytecode map[string][]byte
	reports  map[string]*common.DeanonymizationReport
	keys     map[string][]byte
	isClosed bool // to track if the mock is closed
}

// NewMockDataLayer creates a new mock data layer
func NewMockDataLayer() *MockDataLayer {
	return &MockDataLayer{
		states:   make(map[string]*common.ApplicationState),
		bytecode: make(map[string][]byte),
		reports:  make(map[string]*common.DeanonymizationReport),
		keys:     make(map[string][]byte),
		isClosed: false,
	}
}

func (d *MockDataLayer) checkClosed() error {
	if d.isClosed {
		return ErrStorageIsClosed("Mock data layer is closed")
	}
	return nil
}

func (d *MockDataLayer) StoreApplicationState(ctx context.Context, state *common.ApplicationState) error {
	if err := d.checkClosed(); err != nil {
		return err
	}
	d.states[state.ApplicationID] = state
	return nil
}

func (d *MockDataLayer) GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error) {
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	state, exists := d.states[applicationID]
	if !exists {
		return nil, ErrNotFound("application state not found: " + applicationID)
	}
	return state, nil
}

func (d *MockDataLayer) StoreWASMBytecode(ctx context.Context, applicationID string, bytecode []byte) error {
	if err := d.checkClosed(); err != nil {
		return err
	}
	d.bytecode[applicationID] = bytecode
	return nil
}

func (d *MockDataLayer) GetWASMBytecode(ctx context.Context, applicationID string) ([]byte, error) {
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	bytecode, exists := d.bytecode[applicationID]
	if !exists {
		return nil, ErrNotFound("wasm bytecode not found for application: " + applicationID)
	}
	return bytecode, nil
}

func (d *MockDataLayer) StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	if err := d.checkClosed(); err != nil {
		return err
	}
	d.reports[report.ReportID] = report
	return nil
}

func (d *MockDataLayer) GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error) {
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	report, exists := d.reports[reportID]
	if !exists {
		return nil, ErrNotFound("deanonymization report not found: " + reportID)
	}
	return report, nil
}

func (d *MockDataLayer) StoreUserKey(ctx context.Context, userID string, publicKey []byte) error {
	if err := d.checkClosed(); err != nil {
		return err
	}
	d.keys[userID] = publicKey
	return nil
}

func (d *MockDataLayer) GetUserKey(ctx context.Context, userID string) ([]byte, error) {
	if err := d.checkClosed(); err != nil {
		return nil, err
	}
	publicKey, exists := d.keys[userID]
	if !exists {
		return nil, ErrNotFound("public key not found for user: " + userID)
	}
	return publicKey, nil
}

func (d *MockDataLayer) Close() error {
	// we can close an already closed storage, no problems
	d.isClosed = true
	return nil
}

func ErrNotFound(message string) *Error        { return NewError("not found", message) }
func ErrStorageIsClosed(message string) *Error { return NewError("storage is closed", message) }

type Error struct {
	Code    string
	Message string
}

func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
