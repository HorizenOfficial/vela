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
}

// NewMockDataLayer creates a new mock data layer
func NewMockDataLayer() *MockDataLayer {
	return &MockDataLayer{
		states:   make(map[string]*common.ApplicationState),
		bytecode: make(map[string][]byte),
		reports:  make(map[string]*common.DeanonymizationReport),
	}
}

func (d *MockDataLayer) StoreApplicationState(ctx context.Context, state *common.ApplicationState) error {
	d.states[state.ApplicationID] = state
	return nil
}

func (d *MockDataLayer) GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error) {
	state, exists := d.states[applicationID]
	if !exists {
		return nil, ErrNotFound("application state not found: " + applicationID)
	}
	return state, nil
}

func (d *MockDataLayer) StoreWASMBytecode(ctx context.Context, applicationID string, bytecode []byte) error {
	d.bytecode[applicationID] = bytecode
	return nil
}

func (d *MockDataLayer) GetWASMBytecode(ctx context.Context, applicationID string) ([]byte, error) {
	bytecode, exists := d.bytecode[applicationID]
	if !exists {
		return nil, ErrNotFound("wasm bytecode not found for application: " + applicationID)
	}
	return bytecode, nil
}

func (d *MockDataLayer) StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	d.reports[report.ReportID] = report
	return nil
}

func (d *MockDataLayer) GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error) {
	report, exists := d.reports[reportID]
	if !exists {
		return nil, ErrNotFound("deanonymization report not found: " + reportID)
	}
	return report, nil
}

func (d *MockDataLayer) StoreUserKey(ctx context.Context, userID string, publicKey []byte) error {
	d.keys[userID] = publicKey
	return nil
}

func (d *MockDataLayer) GetUserKey(ctx context.Context, userID string) ([]byte, error) {
	publicKey, exists := d.keys[userID]
	if !exists {
		return nil, ErrNotFound("public key not found for user: " + userID)
	}
	return publicKey, nil
}

func (d *MockDataLayer) Close() error {
	return nil
}

func ErrNotFound(message string) *Error { return NewError("not found", message) }

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
