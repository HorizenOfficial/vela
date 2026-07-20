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
