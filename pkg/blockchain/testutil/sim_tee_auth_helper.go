package testutil

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/processorendpoint"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/tee"
	"github.com/HorizenOfficial/vela/pkg/common"
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

	userEvents := make([][]byte, len(payload.Events))
	userEventSubTypes := make([][32]byte, len(payload.Events))
	for i, event := range payload.Events {
		userEvents[i] = event.EncryptedData
		userEventSubTypes[i] = event.EventSubType
	}

	appEvents := make([][]byte, len(payload.AppEvents))
	appEventSubTypes := make([][32]byte, len(payload.AppEvents))
	for i, appEvent := range payload.AppEvents {
		appEvents[i] = appEvent.Data
		appEventSubTypes[i] = appEvent.EventSubType
	}

	withdrawals := make([]tee.StructsWithdrawalRequest, len(payload.Withdrawals))
	for i, withdrawal := range payload.Withdrawals {
		amount := withdrawal.Amount.ToInt()

		withdrawals[i] = tee.StructsWithdrawalRequest{
			Receiver: withdrawal.DestinationAddress,
			Amount:   amount,
		}
	}

	sigParams := tee.StructsSignatureParams{
		ApplicationId:      processorendpoint.ApplicationIdToBindingType(payload.ApplicationID),
		PrevStateRoot:      payload.PrevStateRoot,
		NewStateRoot:       payload.NewStateRoot,
		ProcessedRequestId: payload.RequestID,
		UserEvents: tee.StructsEventData{
			Events:   userEvents,
			SubTypes: userEventSubTypes,
		},
		AppEvents: tee.StructsEventData{
			Events:   appEvents,
			SubTypes: appEventSubTypes,
		},
		WithdrawalRequests: withdrawals,
		RefundAmount:       payload.RefundAmount.ToInt(),
		ApplicationFee:     payload.ApplicationFee.ToInt(),
		ErrorCode:          payload.ErrorCode,
		ErrorMsg:           payload.ErrorMsg,
	}

	params := s.teeContract.PackCheckSignature(sigParams, payload.Signature)

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
