package manager

import (
	"testing"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validManagerConfig returns a Config with distinct ports and paths so that
// Validate() passes. Tests mutate individual fields to trigger specific errors.
func validManagerConfig() *Config {
	return &Config{
		ChannelType:           "tcp",
		AdminChannelParams:    common.TcpChannelConnectionParams{Ip: "localhost", Port: 4002},
		LogServerTCPAddress:   common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000},
		DataLayerDBPath:       "/tmp/test-db",
		DeanonymizationReportPath: "/tmp/test-reports",
		InputWasmPath:         "/tmp/test-wasm",
		LogFileName:           "/tmp/test-manager.log",
		LogServerLogFile:      "/tmp/test-logserver.log",
		HandshakeTimeout:          5,
		BlockchainPollingInterval: 5,
		ReorgTimeout:              180,
		BlockchainConnectTimeout:  10,
		CommunicationParams:       common.CommunicationParams{RequestTimeoutSec: 30},
		AdminCommunicationParams:  common.CommunicationParams{RequestTimeoutSec: 30},
	}
}

func TestValidate_HappyPath(t *testing.T) {
	cfg := validManagerConfig()
	require.NoError(t, cfg.Validate())
}

func TestValidate_AdminPortEqualsLogServerPort_SameHost(t *testing.T) {
	cfg := validManagerConfig()
	cfg.AdminChannelParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000}
	cfg.LogServerTCPAddress = common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000}

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MANAGER_ADMIN_PORT and LOG_SERVER_PORT")
}

func TestValidate_AdminPortEqualsLogServerPort_DifferentHost(t *testing.T) {
	cfg := validManagerConfig()
	cfg.AdminChannelParams = common.TcpChannelConnectionParams{Ip: "10.0.0.1", Port: 5000}
	cfg.LogServerTCPAddress = common.TcpChannelConnectionParams{Ip: "10.0.0.2", Port: 5000}

	require.NoError(t, cfg.Validate())
}

func TestValidate_DataLayerDBPath_Equals_ReportsPath(t *testing.T) {
	cfg := validManagerConfig()
	cfg.DataLayerDBPath = "/tmp/shared"
	cfg.DeanonymizationReportPath = "/tmp/shared"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MANAGER_DATA_FOLDER and MANAGER_REPORTS_FOLDER")
}

func TestValidate_DataLayerDBPath_Equals_WasmPath(t *testing.T) {
	cfg := validManagerConfig()
	cfg.DataLayerDBPath = "/tmp/shared"
	cfg.InputWasmPath = "/tmp/shared"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MANAGER_DATA_FOLDER and MANAGER_INPUT_WASMS")
}

func TestValidate_LogFile_Equals_LogServerFile(t *testing.T) {
	cfg := validManagerConfig()
	cfg.LogFileName = "/tmp/same.log"
	cfg.LogServerLogFile = "/tmp/same.log"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MANAGER_LOG_FILE_NAME and LOG_SERVER_FILE_NAME")
}

func TestValidate_BothLogFilesEmpty_NoError(t *testing.T) {
	cfg := validManagerConfig()
	cfg.LogFileName = ""
	cfg.LogServerLogFile = ""

	require.NoError(t, cfg.Validate())
}

func TestValidate_RelativePathsResolveSame(t *testing.T) {
	cfg := validManagerConfig()
	cfg.DataLayerDBPath = "./data/../data"
	cfg.DeanonymizationReportPath = "./data"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MANAGER_DATA_FOLDER and MANAGER_REPORTS_FOLDER")
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := validManagerConfig()
	// Trigger both port and path collision
	cfg.AdminChannelParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000}
	cfg.LogServerTCPAddress = common.TcpChannelConnectionParams{Ip: "localhost", Port: 5000}
	cfg.LogFileName = "/tmp/same.log"
	cfg.LogServerLogFile = "/tmp/same.log"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MANAGER_ADMIN_PORT and LOG_SERVER_PORT")
	assert.Contains(t, err.Error(), "MANAGER_LOG_FILE_NAME and LOG_SERVER_FILE_NAME")
}

func TestValidate_EmptyOptionalPaths_NoError(t *testing.T) {
	cfg := validManagerConfig()
	cfg.InputWasmPath = ""
	cfg.LogFileName = ""
	cfg.LogServerLogFile = ""

	require.NoError(t, cfg.Validate())
}

func TestValidate_InvalidChannelType(t *testing.T) {
	cfg := validManagerConfig()
	cfg.ChannelType = "unix"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CHANNEL_TYPE")
	assert.Contains(t, err.Error(), "unix")
}

func TestValidate_VsockChannelType(t *testing.T) {
	cfg := validManagerConfig()
	cfg.ChannelType = "vsock"

	require.NoError(t, cfg.Validate())
}

func TestValidate_HandshakeTimeout_Zero(t *testing.T) {
	cfg := validManagerConfig()
	cfg.HandshakeTimeout = 0

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HANDSHAKE_TIMEOUT")
}

func TestValidate_HandshakeTimeout_Negative(t *testing.T) {
	cfg := validManagerConfig()
	cfg.HandshakeTimeout = -1

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HANDSHAKE_TIMEOUT")
}

func TestValidate_BlockchainPollingInterval_Zero(t *testing.T) {
	cfg := validManagerConfig()
	cfg.BlockchainPollingInterval = 0

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "BLOCKCHAIN_POLLING_INTERVAL")
}

func TestValidate_ReorgTimeout_Zero(t *testing.T) {
	cfg := validManagerConfig()
	cfg.ReorgTimeout = 0

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "REORG_TIMEOUT")
}

func TestValidate_BlockchainConnectTimeout_Zero_IsLegal(t *testing.T) {
	cfg := validManagerConfig()
	cfg.BlockchainConnectTimeout = 0

	require.NoError(t, cfg.Validate())
}

func TestValidate_BlockchainConnectTimeout_Negative(t *testing.T) {
	cfg := validManagerConfig()
	cfg.BlockchainConnectTimeout = -5

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "BLOCKCHAIN_CONNECT_TIMEOUT")
}

func TestValidate_RequestTimeoutSec_Zero(t *testing.T) {
	cfg := validManagerConfig()
	cfg.CommunicationParams.RequestTimeoutSec = 0

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC")
}

func TestValidate_AdminRequestTimeoutSec_Zero(t *testing.T) {
	cfg := validManagerConfig()
	cfg.AdminCommunicationParams.RequestTimeoutSec = 0

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MANAGER_ADMIN_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC")
}

func TestValidate_LogRotationEnabled_NoLogFile(t *testing.T) {
	cfg := validManagerConfig()
	cfg.LogServerRotationEnabled = true
	cfg.LogServerLogFile = ""

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_SERVER_FILE_ROTATION")
}

func TestValidate_LogRotationEnabled_WithLogFile(t *testing.T) {
	cfg := validManagerConfig()
	cfg.LogServerRotationEnabled = true
	cfg.LogServerLogFile = "/tmp/test-logserver.log"

	require.NoError(t, cfg.Validate())
}
