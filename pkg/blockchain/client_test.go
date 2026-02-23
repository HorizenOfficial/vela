package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/blockchain/testutil"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	commontestutil "github.com/horizen-pes/pkg/common/testutil"
	"github.com/stretchr/testify/require"
)

//go:generate mkdir -p ./contracts/mocktee
//go:generate solc --via-ir --combined-json abi,bin ../../contracts/contracts/mocks/MockTeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/MockTeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/MockTeeAuthenticatorAbi/combined.json --pkg mocktee --type MockTeeAuthenticator --out ./contracts/mocktee/MockTeeAuthenticator.go
//go:generate mkdir -p ./contracts/noattestationtee
//go:generate solc --via-ir --combined-json abi,bin ../../contracts/contracts/mocks/NoAttestationTeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/NoAttestationTeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/NoAttestationTeeAuthenticatorAbi/combined.json --pkg noattestationtee --type NoAttestationTeeAuthenticator --out ./contracts/noattestationtee/NoAttestationTeeAuthenticator.go
//go:generate mkdir -p ./contracts/authority
//go:generate solc --via-ir --combined-json abi,bin ../../contracts/contracts/AuthorityRegistry.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/AuthorityRegistryAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/AuthorityRegistryAbi/combined.json --pkg authority --type AuthorityRegistry --out ./contracts/authority/AuthorityRegistry.go

var (
	applicationId = common.NewApplicationId(1)
)

func SetupNewBlockChainClient(testHelper *testutil.SimTestHelper) *BlockChainClient {
	return SetupNewBlockChainClientConnected(testHelper.Client(), testHelper.ProcessorContractAddress, testHelper.TeeSignerAddress, testHelper.ManagerAccount)

}

func setupSimTestHelper(t *testing.T, autoMining bool, teePubSecp521r1 []byte) *testutil.SimTestHelper {
	return testutil.NewSimTestHelper(t, autoMining, true, nil, teePubSecp521r1)
}

func TestGetPendingRequests(t *testing.T) {

	testHelper := setupSimTestHelper(t, false, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	res, err := blockchainClient.getPendingRequests(context.Background())
	require.NoError(t, err)

	require.Equal(t, 0, len(res), "There should be zero pending request")

	currentStateRoot := testHelper.GetStateRoot()
	pendingRequest, stateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.Nil(t, pendingRequest, "There should be no pending request")
	require.Equal(t, currentStateRoot, stateRoot)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1203055)
	maxFeeValue := big.NewInt(100)
	payload := ethCommon.FromHex("0x001234")
	tx := testHelper.SubmitRequest(applicationId, common.Process, payload, transferValue, maxFeeValue)

	testHelper.MineBlock()

	// wait for transaction inclusion
	testHelper.WaitMined(tx)
	fmt.Println("Request was successfully included")

	res, err = blockchainClient.getPendingRequests(context.Background())

	require.NoError(t, err)
	require.Equal(t, 1, len(res), "There should be one pending request")

	request := res[0]
	require.Equal(t, testHelper.ProtocolVersion, request.ProtocolVersion, "Protocol version should match")
	require.Equal(t, applicationId, request.ApplicationID, "Application ID should match")
	require.Equal(t, common.Process, request.RequestType, "Request type should match")
	require.Equal(t, payload, request.Payload, "Payload should match")
	require.Equal(t, 1, request.Timestamp.ToInt().Sign(), "Timestamp should be set and positive")

	require.Equal(t, testHelper.Submitter.From, request.Sender, "Sender should match")
	require.Equal(t, 0, request.DepositAmount.ToInt().Cmp(transferValue), "Value should match")

	pendingRequest, stateRoot, err = blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.Equal(t, request.RequestID, pendingRequest.RequestID, "RequestID should match")
	require.Equal(t, testHelper.ProtocolVersion, pendingRequest.ProtocolVersion, "Protocol version should match")
	require.Equal(t, applicationId, pendingRequest.ApplicationID, "Application ID should match")
	require.Equal(t, common.Process, pendingRequest.RequestType, "Request type should match")
	require.Equal(t, payload, pendingRequest.Payload, "Payload should match")
	require.Equal(t, request.Timestamp, pendingRequest.Timestamp, "Timestamp should match")

	require.Equal(t, testHelper.Submitter.From, pendingRequest.Sender, "Sender should match")
	require.Equal(t, 0, pendingRequest.DepositAmount.ToInt().Cmp(transferValue), "Value should match")

	require.Equal(t, currentStateRoot, stateRoot)

}

func TestMarkRequestFailed(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	transferValue := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.getPendingRequests(context.Background())
	require.NoError(t, err)

	failure := apperrors.New(
		apperrors.CodeSubmittingStateUpdateFailed,
		"test failure",
		nil,
	)

	err = blockchainClient.MarkRequestFailed(context.Background(), res[0].RequestID, failure)
	require.NoError(t, err)

	res, err = blockchainClient.getPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")
}

func TestSubmitStateUpdate(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1000000)
	maxFeeValue := big.NewInt(100)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue, maxFeeValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, oldStateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res.ApplicationID, EncryptedData: []byte{0x04, 0x05, 0x06}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: ethCommon.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: common.NewBig(10)},
	}

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID:  res.ApplicationID,
		RequestID:      res.RequestID,
		PrevStateRoot:  oldStateRoot,
		NewStateRoot:   [32]byte{0x04, 0x05, 0x06},
		Events:         events[:],
		Withdrawals:    withdrawals,
		Signature:      signature[:],
		RefundAmount:   common.NewBig(90),
		ApplicationFee: common.NewBig(10), // 90 + 10 = 100 == maxFeeValue
	}

	// =========================================================
	// Case 1: refund + applicationFees != maxFeeValue -> InvalidValue
	// =========================================================
	payloadWrongSum := *payload // copy value
	payloadWrongSum.RefundAmount = common.NewBig(80)
	payloadWrongSum.ApplicationFee = common.NewBig(10)

	err = blockchainClient.SubmitStateUpdate(context.Background(), &payloadWrongSum)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidValue")

	// =========================================================
	// Case 2: applicationFees < minFeePerRequest but correct sum -> InvalidValue
	// =========================================================
	payloadWrongFee := *payload
	payloadWrongFee.RefundAmount = common.NewBig(100)
	payloadWrongFee.ApplicationFee = common.NewBig(0)

	err = blockchainClient.SubmitStateUpdate(context.Background(), &payloadWrongFee)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidValue")

	// =========================================================
	// Case OK: refund=90, applicationFees=10 -> success
	// =========================================================
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)

	listOfRes, err := blockchainClient.getPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(listOfRes), "There should be 0 pending request")

	// Test error - NonceTooLow
	blockchainClient.account.Nonce = big.NewInt(0)
	blockchainClient.account.GasLimit = 100000 //more than enough to avoid out of gas error
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg := err.(ReorgError)
	require.True(t, isReorg)

	blockchainClient.account.Nonce = nil  //reset nonce to let it be fetched from the network
	blockchainClient.account.GasLimit = 0 //reset gas limit to let it be estimated

	// Test error - wrong application id
	payload.ApplicationID = 9999
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidApplicationId")

	// Test error - wrong prev state root
	payload.ApplicationID = res.ApplicationID
	payload.PrevStateRoot = [32]byte{0x07, 0x08, 0xaa, 0xbb, 0xee}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	fmt.Printf("DEBUG update error: %v\n", err)

	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidStateRoot")

	// Test error - wrong request id
	payload.PrevStateRoot = payload.NewStateRoot
	payload.NewStateRoot = [32]byte{0x07, 0x08, 0x09}
	payload.RequestID = commontestutil.GenerateRandomRequestID()
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidRequestId")
}
