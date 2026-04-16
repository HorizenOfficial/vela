package testutil

import (
	"context"
	"crypto/ecdsa"
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
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/authority"
	defaultauthority "github.com/HorizenOfficial/vela/pkg/blockchain/contracts/defaultauthoritychecker"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/mocktee"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/processorendpoint"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/noattestationtee"
	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/stretchr/testify/require"
)

type SimTestHelper struct {
	t *testing.T

	ProtocolVersion uint8
	sim             *simulated.Backend

	processEndpointContract *processorendpoint.ProcessorEndpoint
	processEndpointInstance *bind.BoundContract

	ProcessorContractAddress ethCommon.Address
	TeeSignerAddress         ethCommon.Address
	AuthorityAddress         ethCommon.Address
	DefaultAuthorityAddress  ethCommon.Address
	Deployer                 *bind.TransactOpts
	Submitter                *bind.TransactOpts
	ManagerAccount           *bind.TransactOpts
	ManagerPrivKey           *ecdsa.PrivateKey
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

func (s *SimTestHelper) setupContracts(useMockContracts bool, teeSigner *ethCommon.Address, teePubSecp521r1 []byte) {
	// use the default deployer: it simply creates, signs and submits the deployment transactions
	deployer := bind.DefaultDeployer(s.Deployer, s.sim.Client())

	var tx *ethTypes.Transaction

	if useMockContracts {
		teeDeployParams := bind.DeploymentParams{
			Contracts: []*bind.MetaData{&mocktee.MockTeeAuthenticatorMetaData},
			Inputs: map[string][]byte{
				mocktee.MockTeeAuthenticatorMetaData.ID: mocktee.NewMockTeeAuthenticator().PackConstructor(
					*teeSigner, teePubSecp521r1,
				),
			},
		}

		// create and submit the contract deployment
		deployRes, err := bind.LinkAndDeploy(&teeDeployParams, deployer)
		require.NoError(s.t, err)

		s.TeeSignerAddress, tx = deployRes.Addresses[mocktee.MockTeeAuthenticatorMetaData.ID], deployRes.Txs[mocktee.MockTeeAuthenticatorMetaData.ID]
	} else {
		require.NotNil(s.t, teeSigner, "teeSigner address must be provided when not using mock contracts")
		teeContract := *noattestationtee.NewNoAttestationTeeAuthenticator()

		constructorInput := teeContract.PackConstructor(s.Deployer.From, *teeSigner, teePubSecp521r1)
		teeDeployParams := bind.DeploymentParams{
			Contracts: []*bind.MetaData{&noattestationtee.NoAttestationTeeAuthenticatorMetaData},
			Inputs:    map[string][]byte{noattestationtee.NoAttestationTeeAuthenticatorMetaData.ID: constructorInput},
		}

		// create and submit the contract deployment
		deployRes, err := bind.LinkAndDeploy(&teeDeployParams, deployer)
		require.NoError(s.t, err)

		s.TeeSignerAddress, tx = deployRes.Addresses[noattestationtee.NoAttestationTeeAuthenticatorMetaData.ID], deployRes.Txs[noattestationtee.NoAttestationTeeAuthenticatorMetaData.ID]
	}

	s.sim.Commit()
	// wait for the pending contract to be deployed on-chain
	_, err := bind.WaitDeployed(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err)
	fmt.Printf("Tee authenticator contract deployed at address 0x%x\n", s.TeeSignerAddress)

	// 1) Deploy DefaultAuthority first
	defaultAuthContract := *defaultauthority.NewDefaultAuthority()
	defaultAuthInput := defaultAuthContract.PackConstructor(s.Deployer.From)

	defaultDeployParams := bind.DeploymentParams{
		Contracts: []*bind.MetaData{&defaultauthority.DefaultAuthorityMetaData},
		Inputs:    map[string][]byte{defaultauthority.DefaultAuthorityMetaData.ID: defaultAuthInput},
	}

	defaultDeployRes, err := bind.LinkAndDeploy(&defaultDeployParams, deployer)
	require.NoError(s.t, err)

	s.DefaultAuthorityAddress = defaultDeployRes.Addresses[defaultauthority.DefaultAuthorityMetaData.ID]
	s.sim.Commit()

	// 2) Deploy AuthorityRegistry with (owner, defaultAuthority)
	authorityContract := *authority.NewAuthorityRegistry()

	constructorInput := authorityContract.PackConstructor(
		s.Deployer.From,  // owner
		s.DefaultAuthorityAddress,  // default authority contract
	)

	deployParams := bind.DeploymentParams{
		Contracts: []*bind.MetaData{&authority.AuthorityRegistryMetaData},
		Inputs:    map[string][]byte{authority.AuthorityRegistryMetaData.ID: constructorInput},
	}

	deployRes, err := bind.LinkAndDeploy(&deployParams, deployer)
	require.NoError(s.t, err)

	s.AuthorityAddress, tx = deployRes.Addresses[authority.AuthorityRegistryMetaData.ID], deployRes.Txs[authority.AuthorityRegistryMetaData.ID]
	s.sim.Commit()


	_, err = bind.WaitDeployed(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err)
	fmt.Printf("Authority contract deployed at address 0x%x\n", s.AuthorityAddress)

	contract := *processorendpoint.NewProcessorEndpoint()

	constructorInput = contract.PackConstructor(s.TeeSignerAddress, s.AuthorityAddress, s.ManagerAccount.From, s.Deployer.From, big.NewInt(5))
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

}

func NewSimTestHelper(t *testing.T, autoMining bool, useMockContracts bool, teeSigner *ethCommon.Address, teePubSecp521r1 []byte) *SimTestHelper {

	helper := &SimTestHelper{
		t:                       t,
		ProtocolVersion:         uint8(0),
		processEndpointContract: processorendpoint.NewProcessorEndpoint(),
	}

	helper.Submitter = helper.GenerateNewUser()
	helper.Deployer = helper.GenerateNewUser()
	helper.ManagerAccount = helper.GenerateNewUser()

	helper.sim = simulated.NewBackend(map[ethCommon.Address]ethTypes.Account{
		helper.Submitter.From:      {Balance: big.NewInt(9e18)},
		helper.Deployer.From:       {Balance: big.NewInt(9e18)},
		helper.ManagerAccount.From: {Balance: big.NewInt(9e18)},
	})

	if teeSigner == nil {
		//generate mock tee signer
		teeSignerKey, err := ethCrypto.GenerateKey()
		require.NoError(t, err, "failed to generate tee signer private key")
		teeSignerAddress := ethCrypto.PubkeyToAddress(teeSignerKey.PublicKey)
		teeSigner = &teeSignerAddress
	}
	if teePubSecp521r1 == nil {
		//generate mock secp521r1 pk
		teePubSecp521r1 = make([]byte, 133)
	}
	helper.setupContracts(useMockContracts, teeSigner, teePubSecp521r1)

	helper.processEndpointInstance = helper.processEndpointContract.Instance(helper.sim.Client(), helper.ProcessorContractAddress)

	if autoMining {
		ctx, cancel := context.WithCancel(context.Background())
		helper.cancel = cancel

		go func() {
			fmt.Println("Auto mining enabled")
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

func (s *SimTestHelper) SubmitRequestFromUser(applicationId common.ApplicationIdType, requestType common.RequestType, payload []byte, tokenAddress ethCommon.Address, assetAmount *big.Int, maxFeeValue *big.Int, sender *bind.TransactOpts) *ethTypes.Transaction {
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
	case common.AssociateKey:
		reqType = 3
	default:
		panic("Unsupported request type")
	}

	sender.Value = new(big.Int).Add(assetAmount, maxFeeValue)
	tx, err := bind.Transact(s.processEndpointInstance, sender, s.processEndpointContract.PackSubmitRequest(s.ProtocolVersion, processorendpoint.ApplicationIdToBindingType(applicationId), reqType, payload, tokenAddress, assetAmount, maxFeeValue))
	require.NoError(s.t, err, "failed to submit transaction")
	sender.Value = big.NewInt(0)
	return tx
}

func (s *SimTestHelper) SubmitRequest(applicationId common.ApplicationIdType, requestType common.RequestType, payload []byte, tokenAddress ethCommon.Address, assetAmount *big.Int, maxFeeValue *big.Int) *ethTypes.Transaction {
	return s.SubmitRequestFromUser(applicationId, requestType, payload, tokenAddress, assetAmount, maxFeeValue, s.Submitter)
}

func (s *SimTestHelper) SubmitDeployRequestFromUser(payload []byte, maxFeeValue *big.Int, sender *bind.TransactOpts) *ethTypes.Transaction {
	if payload == nil {
		payload = ethCommon.FromHex("0x00")
	}
	sender.Value = new(big.Int).Set(maxFeeValue)
	tx, err := bind.Transact(s.processEndpointInstance, sender, s.processEndpointContract.PackSubmitDeployRequest(s.ProtocolVersion, payload))
	require.NoError(s.t, err, "failed to submit deploy transaction")
	sender.Value = big.NewInt(0)
	return tx
}

func (s *SimTestHelper) SubmitDeployRequest(payload []byte, maxFeeValue *big.Int) *ethTypes.Transaction {
	return s.SubmitDeployRequestFromUser(payload, maxFeeValue, s.Deployer)
}

func (s *SimTestHelper) MineBlock() ethCommon.Hash {
	require.False(s.t, s.autoMining, "auto mining is enabled, cannot manually mine blocks")
	return s.sim.Commit()
}

func (s *SimTestHelper) WaitMined(tx *ethTypes.Transaction) {
	_, err := bind.WaitMined(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err, "error waiting for tx inclusion")
}

func (s *SimTestHelper) GetStateRoot(applicationId common.ApplicationIdType) [32]byte {
	oldStateRoot, err := bind.Call(s.processEndpointInstance,
		&bind.CallOpts{Pending: false},
		s.processEndpointContract.PackApplicationStateRoots(uint64(applicationId)),
		s.processEndpointContract.UnpackApplicationStateRoots)
	require.NoError(s.t, err)
	return oldStateRoot
}

func (s *SimTestHelper) GetAppCustody(applicationId common.ApplicationIdType, tokenAddress ethCommon.Address) *big.Int {
	funds, err := bind.Call(s.processEndpointInstance,
		&bind.CallOpts{Pending: false},
		s.processEndpointContract.PackAppCustody(uint64(applicationId), tokenAddress),
		s.processEndpointContract.UnpackAppCustody)
	require.NoError(s.t, err)
	return funds
}

func (s *SimTestHelper) GetRequest(requestID common.RequestIdType) processorendpoint.RequestByIdOutput {
	request, err := bind.Call(s.processEndpointInstance,
		&bind.CallOpts{Pending: false},
		s.processEndpointContract.PackRequestById(requestID),
		s.processEndpointContract.UnpackRequestById)
	require.NoError(s.t, err)
	return request
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

func (s *SimTestHelper) GetDeployRequestSubmittedEvent(tx *ethTypes.Transaction) *processorendpoint.ProcessorEndpointDeployRequestSubmitted {
	receipt, err := s.GetTxReceipt(tx)
	require.NoError(s.t, err, "error getting transaction receipt")
	require.GreaterOrEqual(s.t, len(receipt.Logs), 1, "There should be at least one log for DeployRequestSubmitted")
	event := processorendpoint.ProcessorEndpointDeployRequestSubmitted{}
	err = s.processEndpointInstance.UnpackLog(&event,
		processorendpoint.ProcessorEndpointDeployRequestSubmittedEventName, *receipt.Logs[0])
	require.NoError(s.t, err, "error unpacking DeployRequestSubmittedEvent")
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

func (s *SimTestHelper) AddAuthority(applicationId *big.Int, newAuthority ethCommon.Address) *ethTypes.Transaction {
    defaultAuthContract := defaultauthority.NewDefaultAuthority()
    defaultAuthInstance := defaultAuthContract.Instance(s.Client(), s.DefaultAuthorityAddress)

    tx, err := bind.Transact(
        defaultAuthInstance,
        s.Deployer,
        defaultAuthContract.PackAddAllowedAuthority(applicationId, newAuthority),
    )
    require.NoError(s.t, err, "failed to send transaction")

    return tx
}

func (s *SimTestHelper) GetSimTeeAuthenticatorHelper() *SimTeeAuthenticatorHelper {
	return NewSimTeeAuthenticatorHelper(s.t, s.TeeSignerAddress, s.sim.Client())
}
