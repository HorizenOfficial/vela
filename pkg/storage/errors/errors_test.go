package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/HorizenOfficial/vela/pkg/storage/errors"
)

// TestError verifies that the error constructors create errors with the correct code and message.
func TestError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		e := errors.NewError("test_code", "test message")
		assert.Equal(t, "test message", e.Error())
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		err := errors.ErrNotFound("not found")
		assert.Equal(t, errors.NotFound, err.Code)
		assert.Equal(t, "not found", err.Message)
	})

	t.Run("ErrInvalidParameter", func(t *testing.T) {
		err := errors.ErrInvalidParameter("invalid param")
		assert.Equal(t, errors.InvalidParameter, err.Code)
		assert.Equal(t, "invalid param", err.Message)
	})

	t.Run("ErrVersionAlreadyExists", func(t *testing.T) {
		versionID := []byte{0x01, 0x02, 0x03}
		err := errors.ErrVersionAlreadyExists(versionID)
		assert.Equal(t, errors.VersionAlreadyExists, err.Code)
		assert.Equal(t, fmt.Sprintf("version %x already exists", versionID), err.Message)
	})

	t.Run("ErrInconsistentState", func(t *testing.T) {
		err := errors.ErrInconsistentState("inconsistent state")
		assert.Equal(t, errors.InconsistentState, err.Code)
		assert.Equal(t, "inconsistent state", err.Message)
	})

	t.Run("ErrVersionNotFound", func(t *testing.T) {
		versionID := []byte{0x04, 0x05, 0x06}
		err := errors.ErrVersionNotFound(versionID)
		assert.Equal(t, errors.VersionNotFound, err.Code)
		assert.Equal(t, fmt.Sprintf("version %x not found", versionID), err.Message)
	})

	t.Run("ErrStorageIsClosed", func(t *testing.T) {
		err := errors.ErrStorageIsClosed("storage is closed")
		assert.Equal(t, errors.StorageIsClosed, err.Code)
		assert.Equal(t, "storage is closed", err.Message)
	})

	t.Run("ErrNoVersionInDb", func(t *testing.T) {
		err := errors.ErrNoVersionInDb("no versions found in the db")
		assert.Equal(t, errors.NoVersionInDb, err.Code)
		assert.Equal(t, "no versions found in the db", err.Message)
	})
}
