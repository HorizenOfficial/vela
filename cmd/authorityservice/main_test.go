package main

import (
	"context"
	"math/big"
	"testing"

	"github.com/horizen-pes/pkg/blockchain"
	"github.com/stretchr/testify/require"
)

func TestEnsureChainIDMatch(t *testing.T) {
	bc := blockchain.NewMockClient()
	bc.SetChainID(big.NewInt(42))

	err := ensureChainID(context.Background(), bc, 42, "http://example-rpc")
	require.NoError(t, err)
}

func TestEnsureChainIDMismatch(t *testing.T) {
	bc := blockchain.NewMockClient()
	bc.SetChainID(big.NewInt(1))

	err := ensureChainID(context.Background(), bc, 42, "http://example-rpc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "mismatch")
}
