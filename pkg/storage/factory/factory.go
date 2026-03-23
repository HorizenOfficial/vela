package factory

import (
	"fmt"
	"strings"

	"github.com/HorizenOfficial/vela/pkg/storage"
	"github.com/HorizenOfficial/vela/pkg/storage/mockdb"
	versionedDb "github.com/HorizenOfficial/vela/pkg/storage/versioned_leveldb"
)

// DataLayerConfig captures the minimal parameters needed to build a data layer.
type DataLayerConfig struct {
	Type        string
	DBPath      string
	NumVersions int // Maximum number of historical versions to keep per application.
}

const (
	DataLayerTypeVersionedLevelDB = "versioned_leveldb"
	DataLayerTypeMockDB           = "mockdb"
)

// NewDataLayer constructs a DataLayer based on the provided configuration.
func NewDataLayer(cfg DataLayerConfig) (storage.DataLayer, error) {
	switch cfg.Type {
	case DataLayerTypeMockDB:
		return mockdb.NewMockDataLayer(), nil
	case DataLayerTypeVersionedLevelDB:
		if strings.TrimSpace(cfg.DBPath) == "" {
			return nil, fmt.Errorf("data layer path is empty")
		}
		if cfg.NumVersions <= 0 {
			return nil, fmt.Errorf("NumVersions must be positive, got %d", cfg.NumVersions)
		}
		levelCfg := versionedDb.VersionedLevelDBConfig{
			DBPath:         cfg.DBPath,
			VersionsToKeep: cfg.NumVersions,
		}
		return versionedDb.NewVersionedLevelDBDataLayer(levelCfg)
	default:
		return nil, fmt.Errorf("unknown data layer type: %s", cfg.Type)
	}
}
