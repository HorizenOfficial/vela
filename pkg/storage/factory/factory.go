package factory

import (
	"fmt"
	"strings"

	"github.com/horizen-pes/pkg/storage"
	"github.com/horizen-pes/pkg/storage/mockdb"
	versionedDb "github.com/horizen-pes/pkg/storage/versioned_leveldb"
)

// DataLayerConfig captures the minimal parameters needed to build a data layer.
type DataLayerConfig struct {
	Type        string
	DBPath      string
	NumVersions int
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
		levelCfg := versionedDb.VersionedLevelDBConfig{
			DBPath:         cfg.DBPath,
			VersionsToKeep: cfg.NumVersions,
		}
		return versionedDb.NewVersionedLevelDBDataLayer(levelCfg)
	default:
		return nil, fmt.Errorf("unknown data layer type: %s", cfg.Type)
	}
}
