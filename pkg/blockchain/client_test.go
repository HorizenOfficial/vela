package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/blockchain/testutil"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/crypto"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
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
	blockchainClient := NewBlockChainClient(testHelper.ProcessorContractAddress, testHelper.TeeSignerAddress, "", nil)
	blockchainClient.client = testHelper.Client()

	blockchainClient.processorBoundContract = blockchainClient.processorEndpoint.Instance(blockchainClient.client, testHelper.ProcessorContractAddress)
	blockchainClient.teeAuthBoundContract = blockchainClient.teeAuthEndpoint.Instance(blockchainClient.client, testHelper.TeeSignerAddress)

	blockchainClient.account = testHelper.ManagerAccount
	blockchainClient.connected = true

	return blockchainClient

}

func setupSimTestHelper(t *testing.T, autoMining bool, teePubSecp521r1 []byte) *testutil.SimTestHelper {
	return testutil.NewSimTestHelper(t, autoMining, true, nil, teePubSecp521r1)
}

func TestGetPendingRequests(t *testing.T) {

	testHelper := setupSimTestHelper(t, false, nil)
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
	require.Equal(t, request.RequestID, pendingRequest.RequestID, "RequestID should match")
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

	testHelper := setupSimTestHelper(t, true, nil)
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
	requestId := res[0].RequestID

	err = blockchainClient.MarkRequestCompleted(context.Background(), requestId)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")

	// Test that completing the same request results in ProcessorEndpointInvalidRequestId
	err = blockchainClient.MarkRequestCompleted(context.Background(), requestId)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidRequestId")
	_, isReorgErr := err.(ReorgError)
	require.True(t, isReorgErr)

}

func TestMarkRequestFailed(t *testing.T) {

	testHelper := setupSimTestHelper(t, true, nil)
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

	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1000000)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, oldStateRoot, err := blockchainClient.GetNextPendingRequest(context.Background())
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res.ApplicationID, EncryptedData: []byte{0x04, 0x05, 0x06}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: "0x1234567890123456789012345678901234567890", Amount: 10},
	}

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID: res.ApplicationID,
		RequestID:     res.RequestID,
		PrevStateRoot: oldStateRoot,
		NewStateRoot:  [32]byte{0x04, 0x05, 0x06},
		Events:        events[:],
		Withdrawals:   withdrawals,
		Signature:     signature[:],
	}

	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)

	listOfRes, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(listOfRes), "There should be 0 pending request")

	// Test error - NonceTooLow
	blockchainClient.account.Nonce = big.NewInt(0)
	blockchainClient.account.GasLimit = 100000 //more than enough to avoid out of gas error
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg := err.(ReorgError)
	require.True(t, isReorg)

	blockchainClient.account.Nonce = nil //reset nonce to let it be fetched from the network
	blockchainClient.account.GasLimit = 0 //reset gas limit to let it be estimated

	// Test error - wrong application id

	payload.ApplicationID = "9999"
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidApplicationId")

	// Test error - wrong prev state root
	payload.ApplicationID = res.ApplicationID
	payload.PrevStateRoot = [32]byte{0x07, 0x08, 0xaa, 0xbb, 0xee}
	
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidStateRoot")

	// Test error - wrong request id
	payload.PrevStateRoot = payload.NewStateRoot
	payload.NewStateRoot = [32]byte{0x07, 0x08, 0x09}
	payload.RequestID = "9999999999999999999999999999999999999999999999999999999999999999"
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	_, isReorg = err.(ReorgError)
	require.True(t, isReorg)
	require.Contains(t, err.Error(), "ProcessorEndpointInvalidRequestId")

}

func TestGetUserEvents_StopAtFirst(t *testing.T) {
	//generate secp521r1 pair for TEE and user
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate tee private key")
	teePub := teeKey.PublicKey()

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate user private key")
	userPub := userKey.PublicKey()
	
	testHelper := setupSimTestHelper(t, true, teePub.Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// submit request and state update with message event
	messageSkipped := "test message skipped"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageSkipped, teeKey, userPub);
	message := "test message"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message, teeKey, userPub);

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, *applicationId, 0, 0, nil, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(userEvents), "There should be 1 user event")

	require.Equal(t, []byte(message), userEvents[0], "Decrypted message should match original")
}

func TestGetUserEvents_MultipleEvents(t *testing.T) {
	//generate secp521r1 pair for TEE and user
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate tee private key")
	teePub := teeKey.PublicKey()

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate user private key")
	userPub := userKey.PublicKey()
	
	testHelper := setupSimTestHelper(t, true, teePub.Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// submit request and state update with message event
	message1 := "test message 1"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message1, teeKey, userPub);
	message2 := "test message 2"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message2, teeKey, userPub);
	message3 := "test message 3"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message3, teeKey, userPub);

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, *applicationId, 0, 0, nil, false)
	require.NoError(t, err)
	require.Equal(t, 3, len(userEvents), "There should be 3 user event")

	// they are in reverse order
	require.Equal(t, []byte(message1), userEvents[2], "Decrypted message should match original (1)")
	require.Equal(t, []byte(message2), userEvents[1], "Decrypted message should match original (2)")
	require.Equal(t, []byte(message3), userEvents[0], "Decrypted message should match original (2)")

}

func TestGetUserEvents_WithFilter(t *testing.T) {
	//generate secp521r1 pair for TEE and user
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate tee private key")
	teePub := teeKey.PublicKey()

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate user private key")
	userPub := userKey.PublicKey()

	testHelper := setupSimTestHelper(t, true, teePub.Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// submit request and state update with message event
	messageTrue := "test message - true"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageTrue, teeKey, userPub);
	messageFalse := "test message - false"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageFalse, teeKey, userPub);

	//filter function
	filter := func(data []byte) bool {
    	return strings.Contains(string(data), "true")
	}

	//retrieve and decrypt user events
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, *applicationId, 0, 0, filter, false)
	require.NoError(t, err)
	require.Equal(t, 1, len(userEvents), "There should be 1 user event")
	require.Equal(t, []byte(messageTrue), userEvents[0], "Decrypted message should match the one that passes the filter")

}

func TestGetUserEvents_OtherUsersEvents(t *testing.T) {
	//generate secp521r1 pair for TEE and user
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate tee private key")
	teePub := teeKey.PublicKey()

	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate user private key")
	userPub := userKey.PublicKey()

	//generate key for another user
	otherUserKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err, "failed to generate other user private key")
	otherUserPub := otherUserKey.PublicKey()
	
	testHelper := setupSimTestHelper(t, true, teePub.Bytes())
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// submit request and state update with message event
	message := "test message"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, message, teeKey, userPub);
	messageOther := "test message - for other user"
	_submitRequestAndStateUpdateWithEncryptedMessageEvent(t, blockchainClient, testHelper, messageOther, teeKey, otherUserPub);

	//retrieve and decrypt user events
	// for user
	userEvents, err := blockchainClient.GetUserEvents(context.Background(), *userKey, *applicationId, 0, 0, nil, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(userEvents), "There should be 1 user event")
	require.Equal(t, []byte(message), userEvents[0], "Decrypted message should match original")

	// for other user
	otherUserEvents, err := blockchainClient.GetUserEvents(context.Background(), *otherUserKey, *applicationId, 0, 0, nil, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(otherUserEvents), "There should be 1 user event (other)")
	require.Equal(t, []byte(messageOther), otherUserEvents[0], "Decrypted message should match original (other)")
}

func _submitRequestAndStateUpdateWithEncryptedMessageEvent(t *testing.T, blockchainClient *BlockChainClient, testHelper *testutil.SimTestHelper, 
	message string, senderPrivKey *cryptotypes.PrivateKeyP521, receiverPubKey *cryptotypes.PublicKeyP521,
) {
	// submit request and state update
	transferValue := big.NewInt(1000000)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)
	testHelper.WaitMined(tx)
	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	//encrypt event payload with TEE private key and user public key
	encryptedMessage, err := crypto.Encrypt(senderPrivKey, receiverPubKey, []byte(message))
	require.NoError(t, err)

	events := [1]common.Event{{ApplicationID: res[0].ApplicationID, EncryptedData: encryptedMessage}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: "0x1234567890123456789012345678901234567890", Amount: 0},
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

	//complete state update
	err = blockchainClient.SubmitStateUpdate(context.Background(), payload)
	require.NoError(t, err)
}

func TestSubmitRequest(t *testing.T) {
	// mock private key for the client
	testHelper := setupSimTestHelper(t, true, nil)
	defer testHelper.Close()

	blockchainClient := SetupNewBlockChainClient(testHelper)

	// Prepare a request
	protocolVersion := uint8(0)
	applicationId := big.NewInt(1)
	requestType := common.Deploy
	payload := []byte("test-payload")
	value := big.NewInt(1)

	// Submit the request
	requestId, blockNumber, err := blockchainClient.SubmitRequest(context.Background(), protocolVersion, applicationId, requestType, payload, value)
	require.NoError(t, err)
	// Get pending requests
	pending, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	//check block number > 0
	require.True(t, blockNumber > 0, "Returned block number shouldn't be 0")
	// Check that the submitted request is present and matches
	found := false
	// Convert types for comparison
	protocolVersionStr := strconv.FormatUint(uint64(protocolVersion), 10)
	applicationIdStr := applicationId.String()
	valueUint := value.Uint64()
	
	for _, r := range pending {
		if r.RequestID == requestId {
			found = true

			if  r.ProtocolVersion != protocolVersionStr || r.ApplicationID != applicationIdStr || r.RequestType != requestType || string(r.Payload) != string(payload) || r.Value != valueUint {
				t.Errorf(
					"Request fields do not match: got {protocolVersion:%+v, applicationId:%+v, requestType:%+v, payload:%+v, value:%+v}, want {protocolVersion:%+v, applicationId:%+v, requestType:%+v, payload:%+v, value:%+v}",
					r.ProtocolVersion, r.ApplicationID, r.RequestType, string(r.Payload), r.Value,
					protocolVersionStr, applicationIdStr, requestType, string(payload), valueUint,
				)
			}
		}
	}
	if !found {
		t.Errorf("Submitted request not found in pending requests")
	}
}