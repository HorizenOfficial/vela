package testutil

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/params"
	"github.com/horizen-pes/pkg/blockchain/contracts/authority"
	"github.com/horizen-pes/pkg/blockchain/contracts/keyregistry"
	"github.com/horizen-pes/pkg/blockchain/contracts/mocktee"
	"github.com/horizen-pes/pkg/blockchain/contracts/processorendpoint"
	"github.com/horizen-pes/pkg/blockchain/contracts/tee"
	"github.com/horizen-pes/pkg/common"
	"github.com/stretchr/testify/require"
)

type RequestStatus uint8

const (
	RequestPending RequestStatus = iota
	RequestCompleted
	RequestFailedRefunded
	RequestFailedNotRefunded
)

type SimTestHelper struct {
	t *testing.T

	ProtocolVersion uint8
	sim             *simulated.Backend

	processEndpointContract     *processorendpoint.ProcessorEndpoint
	processEndpointInstance     *bind.BoundContract
	keyRegistryContract         *keyregistry.KeyRegistry
	keyRegistryContractInstance *bind.BoundContract

	ProcessorContractAddress ethCommon.Address
	KeyRegistryAddress       ethCommon.Address
	TeeSignerAddress         ethCommon.Address
	Deployer                 *bind.TransactOpts
	Submitter                *bind.TransactOpts
	ManagerAccount           *bind.TransactOpts
	autoMining               bool
	cancel                   context.CancelFunc
}

func (s *SimTestHelper) GenerateNewUser() *bind.TransactOpts {
	// Since we are using a simulated backend, we will get the chain ID
	// from the same place that the simulated backend gets it.
	chainID := params.AllDevChainProtocolChanges.ChainID

	userPrivateKey, err := ethCrypto.GenerateKey()
	require.NoError(s.t, err, "failed to generate user private key")
	return bind.NewKeyedTransactor(userPrivateKey, chainID)
}


func (s *SimTestHelper) Client() simulated.Client {
	return s.sim.Client()
}

func (s *SimTestHelper) setupContracts(useMockContracts bool, teeSigner *ethCommon.Address) {
	// use the default deployer: it simply creates, signs and submits the deployment transactions
	deployer := bind.DefaultDeployer(s.Deployer, s.sim.Client())

	var tx *ethTypes.Transaction

	if useMockContracts {
		teeDeployParams := bind.DeploymentParams{
			Contracts: []*bind.MetaData{&mocktee.MockTeeAuthenticatorMetaData},
		}

		// create and submit the contract deployment
		deployRes, err := bind.LinkAndDeploy(&teeDeployParams, deployer)
		require.NoError(s.t, err)

		s.TeeSignerAddress, tx = deployRes.Addresses[mocktee.MockTeeAuthenticatorMetaData.ID], deployRes.Txs[mocktee.MockTeeAuthenticatorMetaData.ID]
	} else {
		require.NotNil(s.t, teeSigner, "teeSigner address must be provided when not using mock contracts")
		teeContract := *tee.NewTeeAuthenticator()
		constructorInput := teeContract.PackConstructor(s.Deployer.From, *teeSigner)
		teeDeployParams := bind.DeploymentParams{
			Contracts: []*bind.MetaData{&tee.TeeAuthenticatorMetaData},
			Inputs:    map[string][]byte{tee.TeeAuthenticatorMetaData.ID: constructorInput},
		}

		// create and submit the contract deployment
		deployRes, err := bind.LinkAndDeploy(&teeDeployParams, deployer)
		require.NoError(s.t, err)

		s.TeeSignerAddress, tx = deployRes.Addresses[tee.TeeAuthenticatorMetaData.ID], deployRes.Txs[tee.TeeAuthenticatorMetaData.ID]
	}

	s.sim.Commit()
	// wait for the pending contract to be deployed on-chain
	_, err := bind.WaitDeployed(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err)
	fmt.Printf("Tee authenticator contract deployed at address 0x%x\n", s.TeeSignerAddress)

	authorityContract := *authority.NewAuthorityRegistry()

	constructorInput := authorityContract.PackConstructor(s.Deployer.From)

	deployParams := bind.DeploymentParams{
		Contracts: []*bind.MetaData{&authority.AuthorityRegistryMetaData},
		Inputs:    map[string][]byte{authority.AuthorityRegistryMetaData.ID: constructorInput},
	}

	deployRes, err := bind.LinkAndDeploy(&deployParams, deployer)
	require.NoError(s.t, err)

	authorityAddress, tx := deployRes.Addresses[authority.AuthorityRegistryMetaData.ID], deployRes.Txs[authority.AuthorityRegistryMetaData.ID]

	s.sim.Commit()

	_, err = bind.WaitDeployed(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err)
	fmt.Printf("Authority contract deployed at address 0x%x\n", authorityAddress)

	contract := *processorendpoint.NewProcessorEndpoint()

	constructorInput = contract.PackConstructor(s.TeeSignerAddress, authorityAddress, s.ManagerAccount.From)
	// set up params to deploy an instance of the ProcessorEndpoint contract
	deployParams = bind.DeploymentParams{
		Contracts: []*bind.MetaData{&processorendpoint.ProcessorEndpointMetaData},
		Inputs:    map[string][]byte{processorendpoint.ProcessorEndpointMetaData.ID: constructorInput},
	}

	// create and submit the contract deployment
	deployRes, err = bind.LinkAndDeploy(&deployParams, deployer)
	require.NoError(s.t, err)

	s.ProcessorContractAddress, tx = deployRes.Addresses[processorendpoint.ProcessorEndpointMetaData.ID], deployRes.Txs[processorendpoint.ProcessorEndpointMetaData.ID]

	// call Commit to make the simulated backend mine a block
	s.sim.Commit()

	// wait for the pending contract to be deployed on-chain
	_, err = bind.WaitDeployed(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err)
	fmt.Printf("Processor Endpoint contract deployed at address 0x%x\n", s.ProcessorContractAddress)

	deployParams = bind.DeploymentParams{
		Contracts: []*bind.MetaData{&keyregistry.KeyRegistryMetaData},
	}

	// create and submit the contract deployment
	deployRes, err = bind.LinkAndDeploy(&deployParams, deployer)
	require.NoError(s.t, err)

	s.KeyRegistryAddress, tx = deployRes.Addresses[keyregistry.KeyRegistryMetaData.ID], deployRes.Txs[keyregistry.KeyRegistryMetaData.ID]

	s.sim.Commit()

	// wait for the pending contract to be deployed on-chain
	_, err = bind.WaitDeployed(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err)

	fmt.Printf("Key Registry contract deployed at address 0x%x\n", s.KeyRegistryAddress)

}

func NewSimTestHelper(t *testing.T, autoMining bool, useMockContracts bool, teeSigner *ethCommon.Address) *SimTestHelper {

	helper := &SimTestHelper{
		t:                       t,
		ProtocolVersion:         uint8(0),
		processEndpointContract: processorendpoint.NewProcessorEndpoint(),
		keyRegistryContract:     keyregistry.NewKeyRegistry(),
	}

	helper.Submitter = helper.GenerateNewUser()
	helper.Deployer = helper.GenerateNewUser()
	helper.ManagerAccount = helper.GenerateNewUser()

	helper.sim = simulated.NewBackend(map[ethCommon.Address]ethTypes.Account{
		helper.Submitter.From:      {Balance: big.NewInt(9e18)},
		helper.Deployer.From:       {Balance: big.NewInt(9e18)},
		helper.ManagerAccount.From: {Balance: big.NewInt(9e18)},
	})

	helper.setupContracts(useMockContracts, teeSigner)

	helper.processEndpointInstance = helper.processEndpointContract.Instance(helper.sim.Client(), helper.ProcessorContractAddress)

	helper.keyRegistryContractInstance = helper.keyRegistryContract.Instance(helper.sim.Client(), helper.KeyRegistryAddress)

	if autoMining {
		go func() {
			fmt.Println("Auto mining enabled")
			ctx, cancel := context.WithCancel(context.Background())
			helper.cancel = cancel
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					block := helper.sim.Commit()
					fmt.Println("Mined block: ", block)
				case <-ctx.Done():
					fmt.Println("Shutting down auto mining")
					return
				}
			}

		}()
	}
	helper.autoMining = autoMining

	return helper
}

func (s *SimTestHelper) SubmitRequestFromUser(applicationId *big.Int, requestType common.RequestType, payload []byte, value *big.Int, sender *bind.TransactOpts) *ethTypes.Transaction {
	if payload == nil {
		payload = ethCommon.FromHex("0x00")
	}

	var reqType uint8
	switch requestType {
	case common.Deploy:
		reqType = 0
	case common.Process:
		reqType = 1
	case common.Deanonymize:
		reqType = 2
	default:
		panic("Unsupported request type")
	}

	sender.Value = value
	tx, err := bind.Transact(s.processEndpointInstance, sender, s.processEndpointContract.PackSubmitRequest(s.ProtocolVersion, applicationId, reqType, payload, value))
	require.NoError(s.t, err, "failed to submit transaction")
	sender.Value = big.NewInt(0)
	return tx
}

func (s *SimTestHelper) SubmitRequest(applicationId *big.Int, requestType common.RequestType, payload []byte, value *big.Int) *ethTypes.Transaction {
	return s.SubmitRequestFromUser(applicationId, requestType, payload, value, s.Submitter)
}

func (s *SimTestHelper) MineBlock() ethCommon.Hash {
	require.False(s.t, s.autoMining, "auto mining is enabled, cannot manually mine blocks")
	return s.sim.Commit()
}

func (s *SimTestHelper) WaitMined(tx *ethTypes.Transaction) {
	_, err := bind.WaitMined(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err, "error waiting for tx inclusion")
}

func (s *SimTestHelper) GetStateRoot() [32]byte {
	oldStateRoot, err := bind.Call(s.processEndpointInstance,
		&bind.CallOpts{Pending: false},
		s.processEndpointContract.PackStateRoot(),
		s.processEndpointContract.UnpackStateRoot)
	require.NoError(s.t, err)
	return oldStateRoot
}

func (s *SimTestHelper) GetRequest(requestID string) processorendpoint.RequestByIdOutput {
	reqId, err := common.RequestIdStringTo32Byte(requestID)
	require.NoError(s.t, err, "invalid request ID")
	request, err := bind.Call(s.processEndpointInstance,
		&bind.CallOpts{Pending: false},
		s.processEndpointContract.PackRequestById(reqId),
		s.processEndpointContract.UnpackRequestById)
	require.NoError(s.t, err)
	return request
}

func (s *SimTestHelper) RegisterUserKey(sender *bind.TransactOpts, userKey []byte) (*ethTypes.Transaction, error) {

	tx, err := bind.Transact(s.keyRegistryContractInstance, sender, s.keyRegistryContract.PackRegisterPK(userKey[:]))

	return tx, err

}

func (s *SimTestHelper) GetTxReceipt(tx *ethTypes.Transaction) (*ethTypes.Receipt, error) {

	receipt, err := s.sim.Client().TransactionReceipt(context.Background(), tx.Hash())

	return receipt, err
}

func (s *SimTestHelper) GetRequestSubmittedEvent(tx *ethTypes.Transaction) *processorendpoint.ProcessorEndpointRequestSubmitted {

	receipt, err := s.GetTxReceipt(tx)
	require.NoError(s.t, err, "error getting transaction receipt")
	require.Equal(s.t, 1, len(receipt.Logs), "There should be one RequestSubmittedEvent")
	event := processorendpoint.ProcessorEndpointRequestSubmitted{}
	err = s.processEndpointInstance.UnpackLog(&event,
		processorendpoint.ProcessorEndpointRequestSubmittedEventName, *receipt.Logs[0])
	require.NoError(s.t, err, "error unpacking RequestSubmittedEvent")
	return &event
}

func (s *SimTestHelper) Close() {
	if s.sim != nil {
		err := s.sim.Close()
		require.NoError(s.t, err, "failed to close simulated backend")
	}

	if s.autoMining {
		s.cancel()
	}
}


func (s *SimTestHelper) TransferFunds(sender *bind.TransactOpts, toAddress ethCommon.Address, value *big.Int) *ethTypes.Transaction {

	// Fetch the account's nonce
	nonce, err := s.sim.Client().PendingNonceAt(context.Background(), sender.From)
	require.NoError(s.t, err, "failed to retrieve account nonce")

	gasLimit := uint64(21000) // The standard gas limit for an Ether transfer

	// Get gas price
	gasPrice, err := s.sim.Client().SuggestGasPrice(context.Background())
	require.NoError(s.t, err, "failed to retrieve gas price")

	// Create and sign the transaction
	baseTx := &ethTypes.LegacyTx{
		To:       &toAddress,
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Value:    value,
		Data:     nil}
	rawTx := ethTypes.NewTx(baseTx)

	signedTx, err := sender.Signer(sender.From, rawTx)
	require.NoError(s.t, err, "failed to sign transaction")

	// Broadcast the transaction
	err = s.sim.Client().SendTransaction(context.Background(), signedTx)
	require.NoError(s.t, err, "failed to send transaction")
	return signedTx
}

func (s *SimTestHelper) GetSimTeeAuthenticatorHelper() *SimTeeAuthenticatorHelper {
	return NewSimTeeAuthenticatorHelper(s.t, s.TeeSignerAddress, s.sim.Client())
}
