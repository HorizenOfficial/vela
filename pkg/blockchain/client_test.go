package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/crypto"
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

func TestGetPendingRequests(t *testing.T) {

	testHelper := NewSimTestHelper(t, false, true, nil)
	defer testHelper.Close()

	blockchainClient := testHelper.SetupNewBlockChainClient()

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	require.Equal(t, 0, len(res), "There should be zero pending request")

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1203055)
	payload := ethCommon.FromHex("0x001234")
	tx := testHelper.SubmitRequest(applicationId, common.Process, payload, transferValue)

	// call Commit to make the simulated backend mine a block
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

	oldbalance, err := blockchainClient.client.BalanceAt(context.Background(), testHelper.deployer.From, nil)
	require.NoError(t, err)
	fmt.Println("KeyRegistry balance: ", oldbalance)
	tx = testHelper.TransferFunds(testHelper.Submitter, testHelper.deployer.From, big.NewInt(10000000000000000))
	testHelper.MineBlock()
	testHelper.WaitMined(tx)
	balance, err := testHelper.sim.Client().BalanceAt(context.Background(), testHelper.deployer.From, nil)
	require.NoError(t, err)
	fmt.Println("KeyRegistry balance: ", balance)
	require.Equal(t, oldbalance.Add(oldbalance, big.NewInt(10000000000000000)).Int64(), balance.Int64(), "KeyRegistry balance should match")
}

func TestMarkRequestCompleted(t *testing.T) {

	testHelper := NewSimTestHelper(t, false, true, nil)
	defer testHelper.Close()

	blockchainClient := testHelper.SetupNewBlockChainClient()

	//*****************************************************
	// submit request
	transferValue := big.NewInt(0)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// call Commit to make the simulated backend mine a block
	testHelper.MineBlock()

	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	time.AfterFunc(1*time.Second, func() {
		testHelper.MineBlock()
	})

	err = blockchainClient.MarkRequestCompleted(context.Background(), res[0].RequestID)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")

}

func TestMarkRequestFailed(t *testing.T) {
	testHelper := NewSimTestHelper(t, false, true, nil)
	defer testHelper.Close()

	blockchainClient := testHelper.SetupNewBlockChainClient()

	transferValue := big.NewInt(1000000)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// call Commit to make the simulated backend mine a block
	testHelper.MineBlock()
	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	time.AfterFunc(1*time.Second, func() {
		testHelper.MineBlock()
	})

	err = blockchainClient.MarkRequestFailed(context.Background(), res[0].RequestID)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")
}

func TestSubmitStateUpdate(t *testing.T) {

	testHelper := NewSimTestHelper(t, false, true, nil)
	defer testHelper.Close()

	blockchainClient := testHelper.SetupNewBlockChainClient()

	//*****************************************************
	// submit request
	transferValue := big.NewInt(1000000)
	tx := testHelper.SubmitRequest(applicationId, common.Process, nil, transferValue)

	// call Commit to make the simulated backend mine a block
	testHelper.MineBlock()
	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	time.AfterFunc(1*time.Second, func() {
		testHelper.MineBlock()
	})

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

	testHelper := NewSimTestHelper(t, false, true, nil)
	defer testHelper.Close()

	blockchainClient := testHelper.SetupNewBlockChainClient()

	res, err := blockchainClient.GetPublicKey(context.Background(), testHelper.Submitter.From.Hex())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be no key")

	userEncryptKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userEncryptPubKey := userEncryptKey.PublicKey().Bytes()

	tx, err := testHelper.RegisterUserKey(testHelper.Submitter, userEncryptPubKey)
	require.NoError(t, err)

	// call Commit to make the simulated backend mine a block
	testHelper.MineBlock()
	// wait for transaction inclusion
	testHelper.WaitMined(tx)

	res, err = blockchainClient.GetPublicKey(context.Background(), testHelper.Submitter.From.Hex())

	require.NoError(t, err)
	require.Equal(t, userEncryptPubKey, res)

}
