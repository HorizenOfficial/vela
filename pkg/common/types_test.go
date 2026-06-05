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
		TokenAddress: velacommon.ETH_TOKEN,
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
		TokenAddress: velacommon.ETH_TOKEN,
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

func TestRequestType_StringAndToUint8(t *testing.T) {
	cases := []struct {
		rt   RequestType
		name string
		val  uint8
	}{
		{Deploy, "deploy", 0},
		{Process, "process", 1},
		{Deanonymize, "deanonymize", 2},
		{AssociateKey, "associatekey", 3},
		{TrustProcess, "trustprocess", 4},
	}
	for _, c := range cases {
		require.Equal(t, c.name, c.rt.String())
		v, err := c.rt.ToUint8()
		require.NoError(t, err)
		require.Equal(t, c.val, v)
	}

	unknown := RequestType(99)
	require.Equal(t, "unknown", unknown.String())
	_, err := unknown.ToUint8()
	require.Error(t, err)
}
