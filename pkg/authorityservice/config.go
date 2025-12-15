package authorityservice

import (
	"fmt"
	"os"

	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/crypto"
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

	// ReportsPath is the filesystem path where deanonymization reports are stored.
	ReportsPath string

	// RpcURL is the RPC URL used to connect to the blockchain node.
	RpcURL string
	// ProcessorAddress is the address of the ProcessorEndpoint contract.
	ProcessorAddress string
	// TeeAuthAddress is the address of the TeeAuthenticator contract.
	TeeAuthAddress string
	// PrivateKey is the Ethereum private key used to sign transactions (even though authority mostly reads).
	PrivateKey cryptotypes.PrivateKeySecp256k1

	// LogConsole is true if we want output on console
	LogConsole bool
	// LogConsoleLevel is the level of logging for the console
	LogConsoleLevel string
	// LogConsoleColor is true if the log output should be colored
	LogConsoleColor bool

	// LogFileName is the path to the log file. An empty string means no output on log file
	LogFileName string
	// LogFIleLevel is the level of logging for the console
	LogFileLevel string
}

func LoadConfig() (*Config, error) {
	fileProps := properties.NewProperties()
	if common.FileExists(confFileName) {
		p, err := properties.LoadFile(confFileName, properties.UTF8)
		if err != nil {
			return nil, err
		}
		fileProps = p
	}

	listenAddr := common.GetConfigVar("AUTHORITY_SERVICE_LISTEN_ADDRESS", ":8081", fileProps)
	tlsCert := common.GetConfigVar("AUTHORITY_SERVICE_TLS_CERT", "", fileProps)
	tlsKey := common.GetConfigVar("AUTHORITY_SERVICE_TLS_KEY", "", fileProps)

	chainID := uint64(common.GetConfigVarInt64("AUTHORITY_SERVICE_CHAIN_ID", 0, fileProps))
	nonceTTL := common.GetConfigVarInt64("AUTHORITY_SERVICE_NONCE_TTL", 300, fileProps)

	// We reuse MANAGER_REPORTS_FOLDER so manager and authority service point to the same folder by default
	reportsPath := common.GetConfigVar("MANAGER_REPORTS_FOLDER", "/tmp/horizen-pes-data/manager_reports", fileProps)

	rpcURL := common.GetConfigVar("CHAIN_RPC_PROTOCOL", "http", fileProps) + "://" +
		common.GetConfigVar("CHAIN_RPC_ADDRESS", "127.0.0.1", fileProps) + ":" +
		common.GetConfigVar("CHAIN_RPC_PORT", "8545", fileProps)

	processorAddress := common.GetConfigVar("CHAIN_PROCESSOR_ADDRESS", "", fileProps)
	teeAuthAddress := common.GetConfigVar("CHAIN_TEEAUTHENTICATOR_ADDRESS", "", fileProps)

	var privateKey *cryptotypes.PrivateKeySecp256k1
	privateKeyFromEnv := os.Getenv("AUTHORITY_SERVICE_KEY_SECP256")
	if privateKeyFromEnv == "" {
		privateKey, _ = crypto.GeneratePrivateKeySecp256k1()
	} else {
		privateKey, _ = crypto.ImportPrivateKeySecp256k1FromHex(privateKeyFromEnv)
	}

	logConsole := common.GetConfigVarBool("AUTHORITY_SERVICE_LOG_CONSOLE", true, fileProps)
	logConsoleLevel := common.GetConfigVar("AUTHORITY_SERVICE_LOG_CONSOLE_LEVEL", "info", fileProps)
	logConsoleColor := common.GetConfigVarBool("AUTHORITY_SERVICE_LOG_CONSOLE_COLOR", false, fileProps)
	logFileName := common.GetConfigVar("AUTHORITY_SERVICE_LOG_FILE_NAME", "", fileProps)
	logFileLevel := common.GetConfigVar("AUTHORITY_SERVICE_LOG_FILE_LEVEL", "info", fileProps)

	return &Config{
		ListenAddress:    listenAddr,
		TLSCertFile:      tlsCert,
		TLSKeyFile:       tlsKey,
		ChainID:          chainID,
		NonceTTLSeconds:  nonceTTL,
		ReportsPath:      reportsPath,
		RpcURL:           rpcURL,
		ProcessorAddress: processorAddress,
		TeeAuthAddress:   teeAuthAddress,
		PrivateKey:       *privateKey,
		LogConsole:       logConsole,
		LogConsoleLevel:  logConsoleLevel,
		LogConsoleColor:  logConsoleColor,
		LogFileName:      logFileName,
		LogFileLevel:     logFileLevel,
	}, nil
}
