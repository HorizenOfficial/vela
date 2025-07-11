package errors

import "fmt"

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

// ErrNotFound creates a new Error instance with "not_found" code.
func ErrNotFound(message string) *Error { return NewError("not_found", message) }

// ErrInvalidParameter creates a new Error instance with "invalid_parameter" code.
func ErrInvalidParameter(message string) *Error { return NewError("invalid_parameter", message) }

// ErrVersionAlreadyExists creates a new Error instance with "version_already_exists" code.
func ErrVersionAlreadyExists(versionID []byte) *Error {
	return NewError("version_already_exists", fmt.Sprintf("version %x already exists", versionID))
}

// ErrInconsistentState creates a new Error instance with "inconsistent_state" code.
func ErrInconsistentState(message string) *Error { return NewError("inconsistent_state", message) }

// ErrVersionNotFound creates a new Error instance with "version_not_found" code.
func ErrVersionNotFound(versionID []byte) *Error {
	return NewError("version_not_found", fmt.Sprintf("version %x not found", versionID))
}

// ErrStorageIsClosed creates a new Error instance with "storage_is_closed" code.
func ErrStorageIsClosed(message string) *Error { return NewError("storage_is_closed", message) }
