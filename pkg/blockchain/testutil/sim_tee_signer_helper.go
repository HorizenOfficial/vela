package testutil

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/horizen-pes/pkg/blockchain/contracts/tee"
	"github.com/horizen-pes/pkg/common"
	"github.com/stretchr/testify/require"
)

type SimTeeSignerHelper struct {
	t *testing.T

	teeContract         *tee.TeeAuthenticator
	teeContractInstance *bind.BoundContract

}

func NewSimTeeSignerHelper(t *testing.T, teeSignerAddress ethCommon.Address, nodeClient simulated.Client) *SimTeeSignerHelper {

	teeContract := tee.NewTeeAuthenticator()

	teeContractInstance := teeContract.Instance(nodeClient, teeSignerAddress)

	return &SimTeeSignerHelper{
		t:                   t,
		teeContract:         teeContract,
		teeContractInstance: teeContractInstance,
	}	
}

func (s *SimTeeSignerHelper) CheckSignature(payload *common.UpdatePayload) bool {

	reqId, ok := common.StringToBigInt(payload.RequestID)
	require.True(s.t, ok, "invalid request ID: %s", payload.RequestID)

	appId, ok := common.StringToBigInt(payload.ApplicationID)
	require.True(s.t, ok, "invalid application ID: %s", payload.ApplicationID)


	events := make([][]byte, len(payload.Events))
	for i, event := range payload.Events {
		events[i] = event.EncryptedData
	}

	withdrawals := make([]tee.StructsWithdrawalRequest, len(payload.Withdrawals))
	for i, withdrawal := range payload.Withdrawals {
		amount := new(big.Int).SetUint64(withdrawal.Amount)
		require.True(s.t, ok, "invalid amount: %s", withdrawal.Amount)

		withdrawals[i] = tee.StructsWithdrawalRequest{
			Receiver: ethCommon.HexToAddress(withdrawal.DestinationAddress),
			Amount:   amount,
		}
	}

	params := s.teeContract.PackCheckSignature(
		appId,
		payload.PrevStateRoot,
		payload.NewStateRoot,
		reqId,
		events,
		withdrawals,
		payload.Signature,
	)

	
	result, err := bind.Call(s.teeContractInstance,
		&bind.CallOpts{Pending: false},
		params,
		s.teeContract.UnpackCheckSignature)

	if err != nil {
		raw, hasRevertErrorData := ethclient.RevertErrorData(err)
		if !hasRevertErrorData {
				s.t.Fatalf("expected call error to contain revert error data.")
		}
		rawUnpackedErr, err := s.teeContract.UnpackError(raw)
		if err != nil {
				s.t.Fatalf("expected to unpack error")
		}
		fmt.Printf("Revert error: %#v\n", rawUnpackedErr)
		require.NoError(s.t, err)
	}
	return result

}

func (s *SimTeeSignerHelper) GetTeeSigner() ethCommon.Address {

	result, err := bind.Call(s.teeContractInstance,
		&bind.CallOpts{Pending: false},
		s.teeContract.PackTeeSigner(),
		s.teeContract.UnpackTeeSigner)

	require.NoError(s.t, err)
	return result

}
