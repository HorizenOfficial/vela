package blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/horizen-pes/pkg/blockchain/contracts/authority"
	"github.com/horizen-pes/pkg/blockchain/contracts/processorendpoint"
	"github.com/horizen-pes/pkg/blockchain/contracts/tee"
	"github.com/horizen-pes/pkg/common"
	"github.com/stretchr/testify/require"
)

//go:generate mkdir -p ./contracts/tee
//go:generate solc --combined-json abi,bin ../../contracts/contracts/mocks/MockTeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules -o ./contracts/tee/MockTeeAuthenticator.json --overwrite
//go:generate abigen --v2 --combined-json ./contracts/tee/MockTeeAuthenticator.json/combined.json --pkg tee --type MockTeeAuthenticator --out ./contracts/tee/MockTeeAuthenticator.go
//go:generate mkdir -p ./contracts/authority
//go:generate solc --combined-json abi,bin ../../contracts/contracts/AuthorityRegistry.sol --base-path ../.. --include-path ../../contracts/node_modules -o ./contracts/authority/AuthorityRegistry.json --overwrite
//go:generate abigen --v2 --combined-json ./contracts/authority/AuthorityRegistry.json/combined.json --pkg authority --type AuthorityRegistry --out ./contracts/authority/AuthorityRegistry.go

const (
	defaultProtocolVersion = uint8(0)
)

var (
	applicationId = big.NewInt(1)
)

func setupTestAccounts(t *testing.T) (submitter *bind.TransactOpts, deployer *bind.TransactOpts, manager *bind.TransactOpts) {
	// Since we are using a simulated backend, we will get the chain ID
	// from the same place that the simulated backend gets it.
	chainID := params.AllDevChainProtocolChanges.ChainID

	submitterPrivateKey, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate submitter private key")
	submitter = bind.NewKeyedTransactor(submitterPrivateKey, chainID)

	deployerPrivateKey, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate deployer private key")
	deployer = bind.NewKeyedTransactor(deployerPrivateKey, chainID)

	managerPrivateKey, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate manager private key")
	manager = bind.NewKeyedTransactor(managerPrivateKey, chainID)

	return
}

func setupRequestContractClient(sim *simulated.Backend, processorAddress ethCommon.Address, keyRegistryAddress ethCommon.Address, manager *bind.TransactOpts) *RequestContractClient {
	blockchainClient := NewRequestContractClient(processorAddress, keyRegistryAddress, "", nil)
	blockchainClient.client = sim.Client()

	blockchainClient.processorBoundContract = blockchainClient.processorEndpoint.Instance(blockchainClient.client, processorAddress)
	blockchainClient.keyRegistryBoundContract = blockchainClient.keyRegistryEndpoint.Instance(blockchainClient.client, keyRegistryAddress)

	blockchainClient.account = manager
	blockchainClient.connected = true


	return blockchainClient
}

func TestGetPendingRequests(t *testing.T) {

	submitter, deployer, manager := setupTestAccounts(t)

	sim := simulated.NewBackend(map[ethCommon.Address]ethTypes.Account{
		submitter.From: {Balance: big.NewInt(9e18)},
		deployer.From:  {Balance: big.NewInt(9e18)},
		manager.From:   {Balance: big.NewInt(9e18)},
	})

	processorAddress, keyRegistryAddress := setupContracts(t, sim, deployer, manager.From)

	contractEndpoint := processorendpoint.NewProcessorEndpoint()
	instance := contractEndpoint.Instance(sim.Client(), processorAddress)

	blockchainClient := setupRequestContractClient(sim, processorAddress, keyRegistryAddress, manager)

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	require.Equal(t, 0, len(res), "There should be zero pending request")

	//*****************************************************
	// submit request
	tx, err := bind.Transact(instance, submitter, contractEndpoint.PackSubmitRequest(defaultProtocolVersion, applicationId, uint8(1), ethCommon.FromHex("0x00"), big.NewInt(0)))
	require.NoError(t, err, "failed to submit transaction")

	// call Commit to make the simulated backend mine a block
	sim.Commit()
	// wait for transaction inclusion
	_, err = bind.WaitMined(context.Background(), sim.Client(), tx.Hash())
	require.NoError(t, err, "error waiting for tx inclusion")

	fmt.Println("Request was successfully included")

	res, err = blockchainClient.GetPendingRequests(context.Background())

	require.NoError(t, err)
	require.Equal(t, 1, len(res), "There should be one pending request")
	request := res[0]
	require.Equal(t, "0", request.ProtocolVersion, "Protocol version should match")
	require.Equal(t, "1", request.ApplicationID, "Application ID should match")
	require.Equal(t, common.Process, request.RequestType, "Request type should match")
	require.Equal(t, ethCommon.FromHex("0x00"), request.Payload, "Payload should match")
	require.Greater(t, request.Timestamp, int64(0), "Timestamp should match")
	require.Equal(t, submitter.From.String(), request.Sender, "Sender should match")

}

func TestMarkRequestCompleted(t *testing.T) {

	submitter, deployer, manager := setupTestAccounts(t)

	sim := simulated.NewBackend(map[ethCommon.Address]ethTypes.Account{
		submitter.From: {Balance: big.NewInt(9e18)},
		manager.From:   {Balance: big.NewInt(9e18)},
		deployer.From:  {Balance: big.NewInt(9e18)},
	})
	contract := *processorendpoint.NewProcessorEndpoint()

	processorAddress, keyRegistryAddress := setupContracts(t, sim, deployer, manager.From)

	instance := contract.Instance(sim.Client(), processorAddress)
	blockchainClient := setupRequestContractClient(sim, processorAddress, keyRegistryAddress, manager)


	//*****************************************************
	// submit request

	tx, err := bind.Transact(instance, submitter, contract.PackSubmitRequest(defaultProtocolVersion, applicationId, uint8(1), ethCommon.FromHex("0x00"), big.NewInt(0)))
	require.NoError(t, err, "failed to submit transaction")

	// call Commit to make the simulated backend mine a block
	sim.Commit()

	// wait for transaction inclusion
	_, err = bind.WaitMined(context.Background(), sim.Client(), tx.Hash())
	require.NoError(t, err, "error waiting for tx inclusion")

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)


	time.AfterFunc(1*time.Second, func() {
		sim.Commit()
	})

	err = blockchainClient.MarkRequestCompleted(context.Background(), res[0].RequestID)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")

}

func TestMarkRequestFailed(t *testing.T) {

	submitter, deployer, manager := setupTestAccounts(t)

	sim := simulated.NewBackend(map[ethCommon.Address]ethTypes.Account{
		submitter.From: {Balance: big.NewInt(9e18)},
		manager.From:   {Balance: big.NewInt(9e18)},
		deployer.From:  {Balance: big.NewInt(9e18)},
	})
	contract := *processorendpoint.NewProcessorEndpoint()

	processorAddress, keyRegistryAddress := setupContracts(t, sim, deployer, manager.From)

	instance := contract.Instance(sim.Client(), processorAddress)

	blockchainClient := setupRequestContractClient(sim, processorAddress, keyRegistryAddress, manager)


	submitter.Value = big.NewInt(1000000)
	tx, err := bind.Transact(instance, submitter, contract.PackSubmitRequest(defaultProtocolVersion, applicationId, uint8(1), ethCommon.FromHex("0x00"), submitter.Value))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %v", err))
	}

	// call Commit to make the simulated backend mine a block
	sim.Commit()
	// wait for transaction inclusion
	_, err = bind.WaitMined(context.Background(), sim.Client(), tx.Hash())
	require.NoError(t, err, "error waiting for tx inclusion")

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	time.AfterFunc(1*time.Second, func() {
		sim.Commit()
	})

	err = blockchainClient.MarkRequestFailed(context.Background(), res[0].RequestID)
	require.NoError(t, err)

	res, err = blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be zero pending request")
}

func TestSubmitStateUpdate(t *testing.T) {

	submitter, deployer, manager := setupTestAccounts(t)

	sim := simulated.NewBackend(map[ethCommon.Address]ethTypes.Account{
		submitter.From: {Balance: big.NewInt(9e18)},
		manager.From:   {Balance: big.NewInt(9e18)},
		deployer.From:  {Balance: big.NewInt(9e18)},
	})

	processorAddress, keyRegistryAddress := setupContracts(t, sim, deployer, manager.From)

	contract := *processorendpoint.NewProcessorEndpoint()
	instance := contract.Instance(sim.Client(), processorAddress)

	blockchainClient := setupRequestContractClient(sim, processorAddress, keyRegistryAddress, manager)


	//*****************************************************
	// submit request
	submitter.Value = big.NewInt(1000000)
	tx, err := bind.Transact(instance, submitter, contract.PackSubmitRequest(defaultProtocolVersion, applicationId, uint8(1), ethCommon.FromHex("0x00"), submitter.Value))
	require.NoError(t, err)

	// call Commit to make the simulated backend mine a block
	sim.Commit()
	// wait for transaction inclusion
	_, err = bind.WaitMined(context.Background(), sim.Client(), tx.Hash())
	require.NoError(t, err, "error waiting for tx inclusion")

	res, err := blockchainClient.GetPendingRequests(context.Background())
	require.NoError(t, err)

	time.AfterFunc(1*time.Second, func() {
		sim.Commit()
	})

	events := [1]common.Event{{ApplicationID: res[0].ApplicationID, EncryptedData: []byte{0x04, 0x05, 0x06}}}
	withdrawals := []common.Withdrawal{
		{DestinationAddress: "0x1234567890123456789012345678901234567890", Amount: "10"},
	}

	oldStateRoot, err := bind.Call(instance,
		&bind.CallOpts{Pending: false},
		contract.PackStateRoot(),
		contract.UnpackStateRoot)
	require.NoError(t, err)

	signature := [65]byte{}
	payload := &common.UpdatePayload{
		ApplicationID: res[0].ApplicationID,
		RequestID:     res[0].RequestID,
		PrevStateRoot: oldStateRoot,
		NewStateRoot:  []byte{0x04, 0x05, 0x06},
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

	submitter, deployer, manager := setupTestAccounts(t)

	sim := simulated.NewBackend(map[ethCommon.Address]ethTypes.Account{
		submitter.From: {Balance: big.NewInt(9e18)},
		deployer.From:  {Balance: big.NewInt(9e18)},
		manager.From:   {Balance: big.NewInt(9e18)},
	})

	processorAddress, keyRegistryAddress := setupContracts(t, sim, deployer, manager.From)

	contractEndpoint := processorendpoint.NewKeyRegistry()
	instance := contractEndpoint.Instance(sim.Client(), keyRegistryAddress)

	blockchainClient := setupRequestContractClient(sim, processorAddress, keyRegistryAddress, manager)


	res, err := blockchainClient.GetPublicKey(context.Background(), submitter.From.Hex())
	require.NoError(t, err)
	require.Equal(t, 0, len(res), "There should be no key")

	//*****************************************************
	// submit request
	submitterKey := [133]byte{1, 2, 3}

	tx, err := bind.Transact(instance, submitter, contractEndpoint.PackRegisterPK(submitterKey[:]))
	require.NoError(t, err, "failed to submit transaction")

	// call Commit to make the simulated backend mine a block
	sim.Commit()
	// wait for transaction inclusion
	_, err = bind.WaitMined(context.Background(), sim.Client(), tx.Hash())
	require.NoError(t, err, "error waiting for tx inclusion")

	fmt.Println("Request was successfully included")

	res, err = blockchainClient.GetPublicKey(context.Background(), submitter.From.Hex())

	require.NoError(t, err)
	require.Equal(t, submitterKey[:], res)

}

func setupContracts(t *testing.T, sim *simulated.Backend, deployerSigner *bind.TransactOpts, processorAddress ethCommon.Address) (processorContractAddress ethCommon.Address, keyRegistryAddress ethCommon.Address) {
	// use the default deployer: it simply creates, signs and submits the deployment transactions
	deployer := bind.DefaultDeployer(deployerSigner, sim.Client())

	teeDeployParams := bind.DeploymentParams{
		Contracts: []*bind.MetaData{&tee.MockTeeAuthenticatorMetaData},
	}

	// create and submit the contract deployment
	deployRes, err := bind.LinkAndDeploy(&teeDeployParams, deployer)
	require.NoError(t, err)

	teeAddress, tx := deployRes.Addresses[tee.MockTeeAuthenticatorMetaData.ID], deployRes.Txs[tee.MockTeeAuthenticatorMetaData.ID]

	sim.Commit()
	// wait for the pending contract to be deployed on-chain
	_, err = bind.WaitDeployed(context.Background(), sim.Client(), tx.Hash());
	require.NoError(t, err)
	fmt.Printf("Tee authenticator contract deployed at address 0x%x\n", teeAddress)

	authorityContract := *authority.NewAuthorityRegistry()

	constructorInput := authorityContract.PackConstructor(deployerSigner.From)

	deployParams := bind.DeploymentParams{
		Contracts: []*bind.MetaData{&authority.AuthorityRegistryMetaData},
		Inputs:    map[string][]byte{authority.AuthorityRegistryMetaData.ID: constructorInput},
	}

	deployRes, err = bind.LinkAndDeploy(&deployParams, deployer)
	require.NoError(t, err)

	authorityAddress, tx := deployRes.Addresses[authority.AuthorityRegistryMetaData.ID], deployRes.Txs[authority.AuthorityRegistryMetaData.ID]

	sim.Commit()
	
	_, err = bind.WaitDeployed(context.Background(), sim.Client(), tx.Hash());
	require.NoError(t, err)
	fmt.Printf("Authority contract deployed at address 0x%x\n", authorityAddress)

	contract := *processorendpoint.NewProcessorEndpoint()

	constructorInput = contract.PackConstructor(teeAddress, authorityAddress, processorAddress)
	// set up params to deploy an instance of the ProcessorEndpoint contract
	deployParams = bind.DeploymentParams{
		Contracts: []*bind.MetaData{&processorendpoint.ProcessorEndpointMetaData},
		Inputs:    map[string][]byte{processorendpoint.ProcessorEndpointMetaData.ID: constructorInput},
	}

	// create and submit the contract deployment
	deployRes, err = bind.LinkAndDeploy(&deployParams, deployer)
	require.NoError(t, err)

	processorContractAddress, tx = deployRes.Addresses[processorendpoint.ProcessorEndpointMetaData.ID], deployRes.Txs[processorendpoint.ProcessorEndpointMetaData.ID]

	// call Commit to make the simulated backend mine a block
	sim.Commit()

	// wait for the pending contract to be deployed on-chain
	 _, err = bind.WaitDeployed(context.Background(), sim.Client(), tx.Hash());
	 require.NoError(t, err)
	fmt.Printf("Processor Endpoint contract deployed at address 0x%x\n", processorContractAddress)

	deployParams = bind.DeploymentParams{
		Contracts: []*bind.MetaData{&processorendpoint.KeyRegistryMetaData},
	}

	// create and submit the contract deployment
	deployRes, err = bind.LinkAndDeploy(&deployParams, deployer)
	require.NoError(t, err)

	keyRegistryAddress, tx = deployRes.Addresses[processorendpoint.KeyRegistryMetaData.ID], deployRes.Txs[processorendpoint.KeyRegistryMetaData.ID]

	sim.Commit()

	// wait for the pending contract to be deployed on-chain
	 _, err = bind.WaitDeployed(context.Background(), sim.Client(), tx.Hash())
	require.NoError(t, err)

	fmt.Printf("Key Registry contract deployed at address 0x%x\n", keyRegistryAddress)

	return

}
