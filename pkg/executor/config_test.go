package executor

import (
	"math/big"
	"testing"

	"github.com/horizen-pes/pkg/common"
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
