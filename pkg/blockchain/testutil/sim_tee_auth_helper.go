package testutil

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/horizen-pes/pkg/blockchain/contracts/processorendpoint"
	"github.com/horizen-pes/pkg/blockchain/contracts/tee"
	"github.com/horizen-pes/pkg/common"
	"github.com/stretchr/testify/require"
)

type SimTeeAuthenticatorHelper struct {
	t *testing.T

	teeContract         *tee.TeeAuthenticator
	teeContractInstance *bind.BoundContract
}

func NewSimTeeAuthenticatorHelper(t *testing.T, teeSignerAddress ethCommon.Address, nodeClient bind.ContractBackend) *SimTeeAuthenticatorHelper {

	teeContract := tee.NewTeeAuthenticator()

	teeContractInstance := teeContract.Instance(nodeClient, teeSignerAddress)

	return &SimTeeAuthenticatorHelper{
		t:                   t,
		teeContract:         teeContract,
		teeContractInstance: teeContractInstance,
	}
}

func (s *SimTeeAuthenticatorHelper) CheckSignature(payload *common.UpdatePayload) bool {

	events := make([][]byte, len(payload.Events))
	for i, event := range payload.Events {
		events[i] = event.EncryptedData
	}

	withdrawals := make([]tee.StructsWithdrawalRequest, len(payload.Withdrawals))
	for i, withdrawal := range payload.Withdrawals {
		amount := withdrawal.Amount

		withdrawals[i] = tee.StructsWithdrawalRequest{
			Receiver: withdrawal.DestinationAddress,
			Amount:   amount,
		}
	}

	params := s.teeContract.PackCheckSignature(
		processorendpoint.ApplicationIdToBindingType(payload.ApplicationID),
		payload.PrevStateRoot,
		payload.NewStateRoot,
		payload.RequestID,
		events,
		withdrawals,
		payload.Signature,
		payload.RefundAmount,
		payload.ApplicationFee,
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

func (s *SimTeeAuthenticatorHelper) GetTeeSigner() ethCommon.Address {

	result, err := bind.Call(s.teeContractInstance,
		&bind.CallOpts{Pending: false},
		s.teeContract.PackTeeSigner(),
		s.teeContract.UnpackTeeSigner)

	require.NoError(s.t, err)
	return result

}
