package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/horizen-pes/pkg/manager"
	boltDb "github.com/horizen-pes/pkg/storage/boltdb"
	"github.com/horizen-pes/pkg/storage/mockdb"
	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

func TestCreateDataLayer(t *testing.T) {
	t.Run("should return mockdb instance when type is mockdb", func(t *testing.T) {
		config := &manager.Config{DataLayerType: "mockdb"}
		dl, err := createDataLayer(config)
		require.NoError(t, err)
		require.NotNil(t, dl)
		_, ok := dl.(*mockdb.MockDataLayer)
		assert.True(t, ok, "Expected a *mockdb.MockDataLayer instance")
	})

	t.Run("should return error when path is empty for non-mockdb types", func(t *testing.T) {
		config := &manager.Config{
			DataLayerType:   "versioned_leveldb",
			DataLayerDBPath: "  ", // Empty/whitespace path
		}
		_, err := createDataLayer(config)
		require.Error(t, err)
		assert.Equal(t, "data layer path is empty", err.Error())
	})

	t.Run("should create versioned_leveldb successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		dbFileName := "test_manager.db"
		dbPath := filepath.Join(tempDir, dbFileName)

		config := &manager.Config{
			DataLayerType:          "versioned_leveldb",
			DataLayerDBPath:        dbPath,
			DataLayerNumOfVersions: 5,
		}

		dl, err := createDataLayer(config)
		require.NoError(t, err)
		require.NotNil(t, dl)

		// Check that the concrete type is correct
		_, ok := dl.(*versionedDb.VersionedLevelDBDataLayer)
		assert.True(t, ok, "Expected a *versionedDb.VersionedLevelDBDataLayer instance")

		// Verify that the database directory was created
		_, err = os.Stat(dbPath)
		assert.NoError(t, err, "Database directory should have been created")

		// Clean up the data layer
		err = dl.Close()
		assert.NoError(t, err, "Closing the data layer should not produce an error")
	})

	t.Run("should return not implemented error for boltdb", func(t *testing.T) {
		tempDir := t.TempDir()
		dbFileName := "test_manager.db"
		dbPath := filepath.Join(tempDir, dbFileName)

		config := &manager.Config{
			DataLayerType:   "boltdb",
			DataLayerDBPath: dbPath,
		}
		dl, err := createDataLayer(config)
		require.NoError(t, err)
		require.NotNil(t, dl)

		// Check that the concrete type is correct
		_, ok := dl.(*boltDb.BoltDBDataLayer)
		assert.True(t, ok, "Expected a *boltDb.BoltDBDataLayer instance")

		// Verify that the database directory was created
		_, err = os.Stat(dbPath)
		assert.NoError(t, err, "Database directory should have been created")

		// Clean up the data layer
		err = dl.Close()
		assert.NoError(t, err, "Closing the data layer should not produce an error")
	})

	t.Run("should return unknown type error for other types", func(t *testing.T) {
		tempDir := t.TempDir()
		config := &manager.Config{
			DataLayerType:   "foo_db",
			DataLayerDBPath: tempDir,
		}
		_, err := createDataLayer(config)
		require.Error(t, err)
		assert.Equal(t, fmt.Sprintf("unknown data layer type: %s", config.DataLayerType), err.Error())
	})
}
