package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/crypto"
	"github.com/horizen-pes/pkg/blockchain/testutil"
	"github.com/stretchr/testify/require"
)

//go:generate mkdir -p ./contracts/mocktee
//go:generate solc --combined-json abi,bin ../../contracts/contracts/mocks/MockTeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/MockTeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/MockTeeAuthenticatorAbi/combined.json --pkg mocktee --type MockTeeAuthenticator --out ./contracts/mocktee/MockTeeAuthenticator.go
//go:generate mkdir -p ./contracts/tee
//go:generate solc --combined-json abi,bin ../../contracts/contracts/TeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/TeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/TeeAuthenticatorAbi/combined.json --pkg tee --type TeeAuthenticator --out ./contracts/tee/TeeAuthenticator.go
//go:generate mkdir -p ./contracts/authority
//go:generate solc --combined-json abi,bin ../../contracts/contracts/AuthorityRegistry.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/AuthorityRegistryAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/AuthorityRegistryAbi/combined.json --pkg authority --type AuthorityRegistry --out ./contracts/authority/AuthorityRegistry.go

var (
	applicationId = big.NewInt(1)
)

func SetupNewBlockChainClient(testHelper *testutil.SimTestHelper) *BlockChainClient {
	blockchainClient := NewBlockChainClient(testHelper.ProcessorContractAddress, testHelper.KeyRegistryAddress, "", nil)
	blockchainClient.client = testHelper.Client()

	blockchainClient.processorBoundContract = blockchainClient.processorEndpoint.Instance(blockchainClient.client, testHelper.ProcessorContractAddress)
	blockchainClient.keyRegistryBoundContract = blockchainClient.keyRegistryEndpoint.Instance(blockchainClient.client, testHelper.KeyRegistryAddress)

	blockchainClient.account = testHelper.ManagerAccount
	blockchainClient.connected = true

	return blockchainClient

}

func setupSimTestHelperManualMining(t *testing.T) *testutil.SimTestHelper {
	return testutil.NewSimTestHelper(t, false, true, nil)
}

func setupSimTestHelperAutoMining(t *testing.T) *testutil.SimTestHelper {
	return testutil.NewSimTestHelper(t, true, true, nil)
}


func TestGetPendingRequests(t *testing.T) {

	testHelper := setupSimTestHelperManualMining(t)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	res, err := blockchainClient.GetPendingRequests(context.Background())
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
	payload := ethCommon.FromHex("0x001234")
	tx := testHelper.SubmitRequest(applicationId, common.Process, payload, transferValue)

	testHelper.MineBlock()

	// wait for transaction inclusion
	testHelper.WaitMined(tx)
	fmt.Println("Request was successfully included")

	res, err = blockchainClient.GetPendingRequests(context.Background())

	require.NoError(t, err)
	require.Equal(t, 1, len(res), "There should be one pending request")

	request := res[0]
	require.Equal(t, strconv.Itoa(int(testHelper.ProtocolVersion)), request.ProtocolVersion, "Protocol version should match")
	require.Equal(t, applicationId.String(), request.ApplicationID, "Application ID should match")
	require.Equal(t, common.Process, request.RequestType, "Request type should match")
	require.Equal(t, payload, request.Payload, "Payload should match")
	require.Greater(t, request.Timestamp, int64(0), "Timestamp should match")

	require.Equal(t, testHelper.Submitter.From.String(), request.Sender, "Sender should match")
	require.Equal(t, transferValue.Uint64(), request.Value, "Value should match")

	pendingRequest, stateRoot, err = blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)
	require.Equal(t, request.RequestID, pendingRequest.RequestID, "Payload should match")
	require.Equal(t, strconv.Itoa(int(testHelper.ProtocolVersion)), pendingRequest.ProtocolVersion, "Protocol version should match")
	require.Equal(t, applicationId.String(), pendingRequest.ApplicationID, "Application ID should match")
	require.Equal(t, common.Process, pendingRequest.RequestType, "Request type should match")
	require.Equal(t, payload, pendingRequest.Payload, "Payload should match")
	require.Equal(t, request.Timestamp, pendingRequest.Timestamp, "Timestamp should match")

	require.Equal(t, testHelper.Submitter.From.String(), pendingRequest.Sender, "Sender should match")
	require.Equal(t, transferValue.Uint64(), pendingRequest.Value, "Value should match")


	require.Equal(t, currentStateRoot, stateRoot)


}

func TestMarkRequestCompleted(t *testing.T) {

	testHelper := setupSimTestHelperAutoMining(t)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(0)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	err = blockchainClient.MarkRequestCompleted(context.Background(), res[0].RequestID)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")

}

func TestMarkRequestFailed(t *testing.T) {

	testHelper := setupSimTestHelperAutoMining(t)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	transferValue := big.NewInt(1000000)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	err = blockchainClient.MarkRequestFailed(context.Background(), res[0].RequestID)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")
}

func TestSubmitStateUpdate(t *testing.T) {

	testHelper := setupSimTestHelperAutoMining(t)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1000000)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res[0].ApplicationID, EncryptedData: []byte{0x04, 0x05, 0x06}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: "0x1234567890123456789012345678901234567890", Amount: 10},
	}

	oldStateRoot := testHelper.GetStateRoot()

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID: res[0].ApplicationID,
		RequestID:     res[0].RequestID,
		PrevStateRoot: oldStateRoot,
		NewStateRoot:  [32]byte{0x04, 0x05, 0x06},
		Events:        events[:],
		Withdrawals:   withdrawals,
		Signature:     signature[:],
	}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be 0 pending request")
}

func TestGetPublicKey(t *testing.T) {

	testHelper := setupSimTestHelperManualMining(t)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	res, err := blockchainClient.GetPublicKey(context.Background(), testHelper.Submitter.From.Hex())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be no key")

	userEncryptKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userEncryptPubKey := userEncryptKey.PublicKey().Bytes()

	tx, err := testHelper.RegisterUserKey(testHelper.Submitter, userEncryptPubKey)
	require.NoError(t, err)

	testHelper.MineBlock()
	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err = blockchainClient.GetPublicKey(context.Background(), testHelper.Submitter.From.Hex())

	require.NoError(t, err)
	require.Equal(t, userEncryptPubKey, res)

}
