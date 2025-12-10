package authorityservice

import (
	"log"
	"os"
	"strconv"

	"github.com/magiconair/properties"
)

const confFileName = "authorityservice.conf"

// Config groups runtime settings for the authority service.
type Config struct {
	// ListenAddress is the host:port to bind the HTTP server to.
	ListenAddress string
	// TLSCertFile and TLSKeyFile optionally enable TLS if both are provided.
	TLSCertFile string
	TLSKeyFile  string

	// ChainID is the expected chain id for requests (replay protection across chains).
	ChainID uint64

	// NonceTTLSeconds controls how long a generated nonce is considered valid.
	NonceTTLSeconds int64

	// Data layer configuration
	DataLayerType          string
	DataLayerDBPath        string
	DataLayerNumOfVersions int
}

// DefaultConfig builds a config using environment variables as overrides.
func DefaultConfig() *Config {
	listenAddr := os.Getenv("AUTHORITY_SERVICE_LISTEN_ADDRESS")
	if listenAddr == "" {
		listenAddr = ":8081"
	}

	tlsCert := os.Getenv("AUTHORITY_SERVICE_TLS_CERT")
	tlsKey := os.Getenv("AUTHORITY_SERVICE_TLS_KEY")

	chainID := uint64(0)
	if chainIDStr := os.Getenv("AUTHORITY_SERVICE_CHAIN_ID"); chainIDStr != "" {
		if val, err := strconv.ParseUint(chainIDStr, 10, 64); err == nil {
			chainID = val
		} else {
			log.Printf("Failed to parse AUTHORITY_SERVICE_CHAIN_ID: %v", err)
		}
	}

	nonceTTL := int64(300)
	if ttlStr := os.Getenv("AUTHORITY_SERVICE_NONCE_TTL"); ttlStr != "" {
		if val, err := strconv.ParseInt(ttlStr, 10, 64); err == nil {
			nonceTTL = val
		} else {
			log.Printf("Failed to parse AUTHORITY_SERVICE_NONCE_TTL: %v", err)
		}
	}

	dataPath := os.Getenv("MANAGER_REPORTS_FOLDER")
	if dataPath == "" {
		dataPath = "/tmp/horizen-pes-data/manager_db"
	}

	dataLayerType := os.Getenv("AUTHORITY_SERVICE_DATA_LAYER_TYPE")
	if dataLayerType == "" {
		dataLayerType = "versioned_leveldb"
	}

	numVersions := 10
	if vStr := os.Getenv("AUTHORITY_SERVICE_DATA_LAYER_VERSIONS"); vStr != "" {
		if val, err := strconv.Atoi(vStr); err == nil {
			numVersions = val
		} else {
			log.Printf("Failed to parse AUTHORITY_SERVICE_DATA_LAYER_VERSIONS: %v", err)
		}
	}

	return &Config{
		ListenAddress:          listenAddr,
		TLSCertFile:            tlsCert,
		TLSKeyFile:             tlsKey,
		ChainID:                chainID,
		NonceTTLSeconds:        nonceTTL,
		DataLayerType:          dataLayerType,
		DataLayerDBPath:        dataPath,
		DataLayerNumOfVersions: numVersions,
	}
}

// ReadConfig loads configuration from authorityservice.conf if present, otherwise falls back to defaults.
func ReadConfig() *Config {
	if !fileExists(confFileName) {
		log.Printf("File %s not found. Using default configuration for authority service.\n", confFileName)
		return DefaultConfig()
	}

	config, err := properties.LoadFile(confFileName, properties.UTF8)
	if err != nil {
		log.Fatalf("Couldn't load %s file.\n", confFileName)
	}

	chainID := config.GetUint64("ChainID", 0)

	return &Config{
		ListenAddress:          config.GetString("ListenAddress", ":8081"),
		TLSCertFile:            config.GetString("TLSCertFile", ""),
		TLSKeyFile:             config.GetString("TLSKeyFile", ""),
		ChainID:                chainID,
		NonceTTLSeconds:        config.GetInt64("NonceTTLSeconds", 300),
		DataLayerType:          config.GetString("DataLayerType", "versioned_leveldb"),
		DataLayerDBPath:        config.GetString("DataLayerDBPath", "/tmp/horizen-pes-data/manager_db"),
		DataLayerNumOfVersions: config.GetInt("DataLayerNumOfVersions", 10),
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
