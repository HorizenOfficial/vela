package boltdb

import (
	"context"
	"encoding/json" // For marshaling/unmarshaling Go structs to/from []byte
	"errors"        // For errors.As
	"fmt"
	"time" // For BoltDB Open options

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage"
	storageErrors "github.com/horizen-pes/pkg/storage/errors" // Alias for custom errors
	"go.etcd.io/bbolt"                                        // The bbolt library
)

// Define bucket names for different data types.
// These are constants and should be defined as byte slices once for efficiency.
var (
	appStatesBucket    = []byte("application_states")
	wasmBytecodeBucket = []byte("wasm_bytecode")
	reportsBucket      = []byte("deanonymization_reports")
	userKeysBucket     = []byte("user_keys")
)

// BoltDBDataLayer is an implementation of the ApplicationStateStore interface using BoltDB.
type BoltDBDataLayer struct {
	db *bbolt.DB
}

// BoltDBConfig holds configuration parameters for opening a BoltDB database.
type BoltDBConfig struct {
	Path    string        // Filesystem path for the database file (e.g., "data/my.db")
	Timeout time.Duration // Timeout for opening the database file
}

// NewBoltDBDataLayer creates a new BoltDBDataLayer instance.
// It opens the BoltDB file and ensures all necessary buckets exist.
func NewBoltDBDataLayer(cfg BoltDBConfig) (*BoltDBDataLayer, error) {
	// Open the BoltDB file with specified path, file permissions, and options.
	// 0600 grants read/write permissions only to the owner.
	db, err := bbolt.Open(cfg.Path, 0600, &bbolt.Options{Timeout: cfg.Timeout})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB at %s: %w", cfg.Path, err)
	}

	// Create buckets if they don't exist. This is done within a write transaction.
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{appStatesBucket, wasmBytecodeBucket, reportsBucket, userKeysBucket} {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return fmt.Errorf("create bucket %q: %w", bucket, err)
			}
		}
		return nil
	})
	if err != nil {
		// If bucket creation fails, close the opened database before returning the error.
		db.Close()
		return nil, fmt.Errorf("failed to create BoltDB buckets: %w", err)
	}

	return &BoltDBDataLayer{db: db}, nil
}

// StoreApplicationState stores the state of an application in the 'application_states' bucket.
func (bdl *BoltDBDataLayer) StoreApplicationState(ctx context.Context, state *common.ApplicationState) error {
	// Serialize the ApplicationState struct into JSON bytes.
	value, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal application state: %w", err)
	}

	// Perform the write operation within a BoltDB write transaction.
	err = bdl.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(appStatesBucket)
		// Put the state using ApplicationID as the key.
		return b.Put([]byte(state.ApplicationID), value)
	})
	if err != nil {
		return fmt.Errorf("failed to store application state in BoltDB: %w", err)
	}
	return nil
}

// GetApplicationState retrieves the state of an application from the 'application_states' bucket.
func (bdl *BoltDBDataLayer) GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error) {
	state := &common.ApplicationState{} // Initialize a new struct to unmarshal into

	// Perform the read operation within a BoltDB read-only transaction.
	err := bdl.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(appStatesBucket)
		// Retrieve the value by ApplicationID.
		value := b.Get([]byte(applicationID))
		if value == nil {
			// If key is not found, return our custom ErrNotFound.
			return storageErrors.ErrNotFound(fmt.Sprintf("application state not found: %s", applicationID))
		}
		// Unmarshal the JSON bytes back into the ApplicationState struct.
		return json.Unmarshal(value, state)
	})
	if err != nil {
		return nil, handleGetError(err, "application state")
	}
	return state, nil
}

// StoreWASMBytecode stores WASM bytecode for an application in the 'wasm_bytecode' bucket.
func (bdl *BoltDBDataLayer) StoreWASMBytecode(ctx context.Context, applicationID string, bytecode []byte) error {
	err := bdl.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(wasmBytecodeBucket)
		return b.Put([]byte(applicationID), bytecode)
	})
	if err != nil {
		return fmt.Errorf("failed to store WASM bytecode in BoltDB: %w", err)
	}
	return nil
}

// GetWASMBytecode gets WASM bytecode for an application from the 'wasm_bytecode' bucket.
func (bdl *BoltDBDataLayer) GetWASMBytecode(ctx context.Context, applicationID string) ([]byte, error) {
	var bytecode []byte // Declare slice to hold copied data
	err := bdl.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(wasmBytecodeBucket)
		val := b.Get([]byte(applicationID))
		if val == nil {
			return storageErrors.ErrNotFound("wasm bytecode not found for application: " + applicationID)
		}
		// IMPORTANT: Copy the slice's content. BoltDB slices are only valid within the transaction.
		bytecode = make([]byte, len(val))
		copy(bytecode, val)
		return nil
	})
	if err != nil {
		return nil, handleGetError(err, "WASM bytecode")
	}
	return bytecode, nil
}

// StoreDeanonymizationReport stores a deanonymization report in the 'deanonymization_reports' bucket.
func (bdl *BoltDBDataLayer) StoreDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	value, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal deanonymization report: %w", err)
	}
	err = bdl.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(reportsBucket)
		return b.Put([]byte(report.ReportID), value)
	})
	if err != nil {
		return fmt.Errorf("failed to store deanonymization report in BoltDB: %w", err)
	}
	return nil
}

// GetDeanonymizationReport gets a deanonymization report from the 'deanonymization_reports' bucket.
func (bdl *BoltDBDataLayer) GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error) {
	report := &common.DeanonymizationReport{}
	err := bdl.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(reportsBucket)
		value := b.Get([]byte(reportID))
		if value == nil {
			return storageErrors.ErrNotFound("deanonymization report not found: " + reportID)
		}
		return json.Unmarshal(value, report)
	})
	if err != nil {
		return nil, handleGetError(err, "deanonymization report")
	}
	return report, nil
}

// StoreUserKey stores a user's public key in the 'user_keys' bucket.
func (bdl *BoltDBDataLayer) StoreUserKey(ctx context.Context, userID string, publicKey []byte) error {
	err := bdl.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(userKeysBucket)
		return b.Put([]byte(userID), publicKey)
	})
	if err != nil {
		return fmt.Errorf("failed to store user key in BoltDB: %w", err)
	}
	return nil
}

// GetUserKey retrieves a user's public key from the 'user_keys' bucket.
func (bdl *BoltDBDataLayer) GetUserKey(ctx context.Context, userID string) ([]byte, error) {
	var publicKey []byte
	err := bdl.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(userKeysBucket)
		val := b.Get([]byte(userID))
		if val == nil {
			return storageErrors.ErrNotFound("public key not found for user: " + userID)
		}
		// IMPORTANT: Copy the slice's content. BoltDB slices are only valid within the transaction.
		publicKey = make([]byte, len(val))
		copy(publicKey, val)
		return nil
	})
	if err != nil {
		return nil, handleGetError(err, "user key")
	}
	return publicKey, nil
}

// Close closes the BoltDB database connection.
func (bdl *BoltDBDataLayer) Close() error {
	// bbolt.DB.Close() returns an error if the database is not open or if there's an I/O issue.
	return bdl.db.Close()
}

// handleGetError simplifies error handling for all Get operations.
// It checks for the custom ErrNotFound and wraps other errors with more context.
func handleGetError(err error, objectName string) error {
	var nfErr *storageErrors.Error
	if errors.As(err, &nfErr) && nfErr.Code == storageErrors.NotFound {
		return nfErr // Return the original "not found" error directly
	}
	// For any other error, wrap it with a more descriptive message.
	return fmt.Errorf("failed to get %s from BoltDB: %w", objectName, err)
}

var _ storage.ApplicationStateStore = (*BoltDBDataLayer)(nil)
