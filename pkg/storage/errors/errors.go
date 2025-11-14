package errors

import "fmt"

const (
	// NotFound is a code for when a resource is not found.
	NotFound = "not_found"
	// InvalidParameter is a code for when a parameter is invalid.
	InvalidParameter = "invalid_parameter"
	// VersionAlreadyExists is a code for when a version already exists.
	VersionAlreadyExists = "version_already_exists"
	// InconsistentState is a code for when the storage is in an inconsistent state.
	InconsistentState = "inconsistent_state"
	// VersionNotFound is a code for when a version is not found.
	VersionNotFound = "version_not_found"
	// StorageIsClosed is a code for when the storage is closed.
	StorageIsClosed = "storage_is_closed"
	// NoVersionInDb is a code for when the storage has not version yet.
	NoVersionInDb = "no_versions_found_in_the_db"
)

// Error is a custom error type used by the storage package.
type Error struct {
	Code    string
	Message string
}

// NewError creates a new instance of the custom Error type.
func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error method makes our custom Error type satisfy the error interface.
func (e *Error) Error() string {
	return e.Message
}

// IsNotFound checks if an error is a NotFound error.
func IsNotFound(err error) bool {
	e, ok := err.(*Error)
	return ok && e.Code == NotFound
}

// ErrNotFound creates a new Error instance with "not_found" code.
func ErrNotFound(message string) *Error { return NewError(NotFound, message) }

// ErrInvalidParameter creates a new Error instance with "invalid_parameter" code.
func ErrInvalidParameter(message string) *Error { return NewError(InvalidParameter, message) }

// ErrVersionAlreadyExists creates a new Error instance with "version_already_exists" code.
func ErrVersionAlreadyExists(versionID []byte) *Error {
	return NewError(VersionAlreadyExists, fmt.Sprintf("version %x already exists", versionID))
}

// ErrInconsistentState creates a new Error instance with "inconsistent_state" code.
func ErrInconsistentState(message string) *Error { return NewError(InconsistentState, message) }

// ErrVersionNotFound creates a new Error instance with "version_not_found" code.
func ErrVersionNotFound(versionID []byte) *Error {
	return NewError(VersionNotFound, fmt.Sprintf("version %x not found", versionID))
}

// ErrStorageIsClosed creates a new Error instance with "storage_is_closed" code.
func ErrStorageIsClosed(message string) *Error { return NewError(StorageIsClosed, message) }

// ErrNoVersionInDb creates a new Error instance with "no_versions_found_in_the_db" code.
func ErrNoVersionInDb(message string) *Error { return NewError(NoVersionInDb, message) }
