package executor

import (
	"math/big"
	"os"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validExecutorConfig returns a Config with valid defaults so that
// Validate() passes. Tests mutate individual fields to trigger specific errors.
func validExecutorConfig() *Config {
	return &Config{
		ChannelType:         "tcp",
		KeySetRecoveryType:  common.RecoveryTypeUnsafe,
		FuelPricePerUnit:    big.NewInt(1),
		MinFeePerRequest:    big.NewInt(10),
		CommunicationParams: common.CommunicationParams{RequestTimeoutSec: 30},
	}
}

func TestValidate_Type0_NoError(t *testing.T) {
	cfg := validExecutorConfig()
	require.NoError(t, cfg.Validate())
}

func TestValidate_KMS_MissingARN(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.KeySetRecoveryType = common.RecoveryTypeKMS
	cfg.KMSRegion = "eu-west-1"
	cfg.KMSKeyARN = ""

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "KMS key ARN")
}

func TestValidate_KMS_MissingRegion(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.KeySetRecoveryType = common.RecoveryTypeKMS
	cfg.KMSKeyARN = "arn:aws:kms:eu-west-1:123456789:key/test"
	cfg.KMSRegion = ""

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "KMS region")
}

func TestValidate_KMS_MissingBothARNAndRegion(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.KeySetRecoveryType = common.RecoveryTypeKMS
	cfg.KMSKeyARN = ""
	cfg.KMSRegion = ""

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "KMS key ARN")
	assert.Contains(t, err.Error(), "KMS region")
}

func TestValidate_KMS_ValidConfig(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.KeySetRecoveryType = common.RecoveryTypeKMS
	cfg.KMSKeyARN = "arn:aws:kms:eu-west-1:123456789:key/test"
	cfg.KMSRegion = "eu-west-1"

	require.NoError(t, cfg.Validate())
}

func TestValidate_InvalidChannelType(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.ChannelType = "unix"

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CHANNEL_TYPE")
	assert.Contains(t, err.Error(), "unix")
}

func TestValidate_VsockChannelType(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.ChannelType = "vsock"

	require.NoError(t, cfg.Validate())
}

func TestValidate_FuelPricePerUnit_Nil(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.FuelPricePerUnit = nil

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_FUEL_PRICE_PER_UNIT")
}

func TestValidate_FuelPricePerUnit_Zero(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.FuelPricePerUnit = big.NewInt(0)

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_FUEL_PRICE_PER_UNIT")
}

func TestValidate_FuelPricePerUnit_Negative(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.FuelPricePerUnit = big.NewInt(-1)

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_FUEL_PRICE_PER_UNIT")
}

func TestValidate_MinFeePerRequest_Nil(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.MinFeePerRequest = nil

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_MIN_FEE_PER_REQUEST")
}

func TestValidate_MinFeePerRequest_Negative(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.MinFeePerRequest = big.NewInt(-1)

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_MIN_FEE_PER_REQUEST")
}

func TestValidate_MinFeePerRequest_Zero_IsLegal(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.MinFeePerRequest = big.NewInt(0)

	require.NoError(t, cfg.Validate())
}

func TestValidate_RequestTimeoutSec_Zero(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.CommunicationParams.RequestTimeoutSec = 0

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC")
}

func TestValidate_RequestTimeoutSec_Negative(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.CommunicationParams.RequestTimeoutSec = -5

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC")
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.ChannelType = "invalid"
	cfg.FuelPricePerUnit = big.NewInt(0)
	cfg.CommunicationParams.RequestTimeoutSec = 0

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CHANNEL_TYPE")
	assert.Contains(t, err.Error(), "EXECUTOR_FUEL_PRICE_PER_UNIT")
	assert.Contains(t, err.Error(), "EXECUTOR_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC")
}

func TestValidate_MaxGuestMemoryBytes_InRange_NoError(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.MaxGuestMemoryBytes = 512 * 1024 * 1024
	require.NoError(t, cfg.Validate())
}

func TestValidate_MaxGuestMemoryBytes_Negative(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.MaxGuestMemoryBytes = -1

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_MAX_GUEST_MEMORY_BYTES")
}

func TestValidate_MaxGuestMemoryBytes_AboveCeiling(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.MaxGuestMemoryBytes = (2 << 30) + 1

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_MAX_GUEST_MEMORY_BYTES")
}

// loadConfigFromFile writes an executor.conf with the given content in a fresh
// temp working directory and runs LoadConfig against it, exercising the real
// file-parsing path (GetConfigVarInt64 etc.) rather than setting struct fields
// directly.
func loadConfigFromFile(t *testing.T, conf string) *Config {
	t.Helper()
	t.Chdir(t.TempDir())
	require.NoError(t, os.WriteFile(confFileName, []byte(conf), 0600))
	cfg, err := LoadConfig()
	require.NoError(t, err)
	return cfg
}

func TestLoadConfig_MaxGuestMemoryBytes_AboveCeilingInConfFile_Rejected(t *testing.T) {
	// Clear any ambient env overrides: empty env values fall through to the file.
	t.Setenv("EXECUTOR_MAX_GUEST_MEMORY_BYTES", "")
	t.Setenv("EXECUTOR_KEYSET_RECOVERY_TYPE", "")

	// 3 GiB — above the 2 GiB ceiling. The value must survive parsing intact
	// (regression: a 32-bit ParseInt silently replaced it with the default,
	// so Validate never saw it) and then be rejected by Validate.
	cfg := loadConfigFromFile(t,
		"EXECUTOR_KEYSET_RECOVERY_TYPE=0\nEXECUTOR_MAX_GUEST_MEMORY_BYTES=3221225472\n")
	require.Equal(t, int64(3221225472), cfg.MaxGuestMemoryBytes)

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_MAX_GUEST_MEMORY_BYTES")
}

// TestLoadConfig_OutOfRangePorts_FallBackToDefaults covers values that do not fit
// the uint32 that ports and vsock CIDs are stored as.
//
// Widening the shared parser to the full int64 range (needed for
// EXECUTOR_MAX_GUEST_MEMORY_BYTES) removed the accidental guard these callers had:
// a 32-bit parse used to fail and fall back to the default. Without a check they
// are silently truncated instead — 2^32 becomes port 0, which makes a listener
// bind an arbitrary OS-assigned port rather than failing, so the manager dials the
// configured port and gets connection refused with nothing logged at startup.
func TestLoadConfig_OutOfRangePorts_FallBackToDefaults(t *testing.T) {
	for _, envVar := range []string{
		"EXECUTOR_MAX_GUEST_MEMORY_BYTES", "EXECUTOR_KEYSET_RECOVERY_TYPE",
		"CHANNEL_TYPE", "EXECUTOR_PORT", "EXECUTOR_KMS_PROXY_PORT", "LOG_SERVER_PORT",
	} {
		t.Setenv(envVar, "")
	}

	cfg := loadConfigFromFile(t, "EXECUTOR_KEYSET_RECOVERY_TYPE=0\n"+
		"CHANNEL_TYPE=tcp\n"+
		"EXECUTOR_PORT=4294967296\n"+
		"EXECUTOR_KMS_PROXY_PORT=4294967296\n"+
		"LOG_SERVER_PORT=4294967296\n")

	require.Equal(t, uint32(8000), cfg.KMSProxyPort, "KMS proxy port must fall back, not truncate to 0")

	tcpParams, ok := cfg.ChannelParams.(common.TcpChannelConnectionParams)
	require.True(t, ok, "expected TCP channel params")
	require.Equal(t, uint32(4000), tcpParams.Port, "executor port must fall back, not truncate to 0")

	logParams, ok := cfg.LogChannelParams.(common.TcpChannelConnectionParams)
	require.True(t, ok, "expected TCP log channel params")
	require.Equal(t, uint32(5000), logParams.Port, "log server port must fall back, not truncate to 0")
}

func TestLoadConfig_MaxGuestMemoryBytes_CeilingInConfFile_Accepted(t *testing.T) {
	t.Setenv("EXECUTOR_MAX_GUEST_MEMORY_BYTES", "")
	t.Setenv("EXECUTOR_KEYSET_RECOVERY_TYPE", "")

	// Exactly 2 GiB (2147483648) is the maximum legal value; it is above
	// math.MaxInt32, so it also only works with a full 64-bit parse.
	cfg := loadConfigFromFile(t,
		"EXECUTOR_KEYSET_RECOVERY_TYPE=0\nEXECUTOR_MAX_GUEST_MEMORY_BYTES=2147483648\n")
	require.Equal(t, int64(2147483648), cfg.MaxGuestMemoryBytes)
	require.NoError(t, cfg.Validate())
}

// TestValidate_TCPPortAboveRange covers the range check that GetConfigVarUint32
// alone cannot provide: 70000 fits in a uint32, so it survives parsing and would
// otherwise only fail later at bind time with a generic error.
func TestValidate_TCPPortAboveRange(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.ChannelType = "tcp"
	cfg.ChannelParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: 70000}

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_PORT")
}

// TestValidate_LogPortAboveRange is the same check for the log channel.
func TestValidate_LogPortAboveRange(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.ChannelType = "tcp"
	cfg.LogChannelParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: 70000}

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_SERVER_PORT")
}

// TestValidate_VsockPortAboveTCPRangeAccepted guards the other side of the rule:
// vsock ports are full uint32 values, so the TCP bound must not apply to them.
func TestValidate_VsockPortAboveTCPRangeAccepted(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.ChannelType = "vsock"
	cfg.ChannelParams = common.VSockChannelConnectionParams{CID: 3, Port: 70000}

	require.NoError(t, cfg.Validate())
}

// TestValidate_TCPPortZero covers the other unusable port. 0 parses cleanly through
// GetConfigVarUint32, so only Validate() can catch it: the executor would listen on an
// OS-assigned port while the manager dials the configured one from the same variable.
func TestValidate_TCPPortZero(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.ChannelType = "tcp"
	cfg.ChannelParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: 0}

	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EXECUTOR_PORT")
}

// TestValidate_LogPortZeroAccepted guards the deliberate exception: 0 on the log
// channel means "no TCP log listener" (logserver.StartLogServer), so it must not be
// rejected alongside the rendezvous port.
func TestValidate_LogPortZeroAccepted(t *testing.T) {
	cfg := validExecutorConfig()
	cfg.ChannelType = "tcp"
	cfg.ChannelParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: 4000}
	cfg.LogChannelParams = common.TcpChannelConnectionParams{Ip: "localhost", Port: 0}

	require.NoError(t, cfg.Validate())
}
