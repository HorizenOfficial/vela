package common

import (
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/stretchr/testify/require"
)

func TestRequestValidate_ZeroAssetAmountWithZeroToken(t *testing.T) {
	req := &Request{
		Timestamp:    NewBig(1000),
		TokenAddress: velacommon.NativeTokenAddress(),
		AssetAmount:  NewBig(0),
		MaxFeeValue:  NewBig(100),
	}
	require.NoError(t, req.Validate())
}

func TestRequestValidate_NonZeroAssetAmountWithNonZeroToken(t *testing.T) {
	req := &Request{
		Timestamp:    NewBig(1000),
		TokenAddress: ethCommon.HexToAddress("0xdead000000000000000000000000000000000001"),
		AssetAmount:  NewBig(500),
		MaxFeeValue:  NewBig(100),
	}
	require.NoError(t, req.Validate())
}

func TestRequestValidate_NonZeroAssetAmountWithZeroToken(t *testing.T) {
	req := &Request{
		Timestamp:    NewBig(1000),
		TokenAddress: velacommon.NativeTokenAddress(),
		AssetAmount:  NewBig(500),
		MaxFeeValue:  NewBig(100),
	}
	require.NoError(t, req.Validate())
}

func TestRequestValidate_ZeroAssetAmountWithNonZeroToken(t *testing.T) {
	req := &Request{
		Timestamp:    NewBig(1000),
		TokenAddress: ethCommon.HexToAddress("0xdead000000000000000000000000000000000001"),
		AssetAmount:  NewBig(0),
		MaxFeeValue:  NewBig(100),
	}
	err := req.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "tokenAddress must be zero address when assetAmount is zero")
}
