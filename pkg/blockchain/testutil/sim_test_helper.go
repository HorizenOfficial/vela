package testutil

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/params"
	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/authority"
	defaultauthority "github.com/HorizenOfficial/vela/pkg/blockchain/contracts/defaultauthoritychecker"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/mockerc20"
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

	// MockERC20 is populated by DeployMockERC20 and remains the zero value if
	// the test never deploys one. Only one MockERC20 per SimTestHelper — tests
	// that need multiple tokens should call DeployMockERC20 multiple times and
	// track the returned addresses themselves.
	MockERC20Address  ethCommon.Address
	mockERC20Contract *mockerc20.MockERC20
	mockERC20Instance *bind.BoundContract

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
	autoMineWG               sync.WaitGroup
}

func (s *SimTestHelper) GenerateNewUser() *bind.TransactOpts {
	opts, _ := s.generateNewUserWithKey()
	return opts
}

// generateNewUserWithKey is like GenerateNewUser but also returns the raw
// private key for callers that need to sign payloads directly (not just
// send contract txs through the TransactOpts). Used internally so the
// manager's private key can be stashed on SimTestHelper.ManagerPrivKey.
func (s *SimTestHelper) generateNewUserWithKey() (*bind.TransactOpts, *ecdsa.PrivateKey) {
	chainID := params.AllDevChainProtocolChanges.ChainID
	userPrivateKey, err := ethCrypto.GenerateKey()
	require.NoError(s.t, err, "failed to generate user private key")
	return bind.NewKeyedTransactor(userPrivateKey, chainID), userPrivateKey
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
	// Capture the manager's private key too — some tests need to sign
	// payloads (e.g. forge a state update for the InvalidStateRoot negative
	// test) rather than only submit contract txs.
	helper.ManagerAccount, helper.ManagerPrivKey = helper.generateNewUserWithKey()

	// Deployer funds user accounts (5 ETH each via CreateFundedAccount) AND
	// pays gas for test setup (contract deploys, role grants, allowlist ops).
	// 9 ETH was enough for single-user tests but runs out after the second
	// CreateFundedAccount in multi-user scenarios. 100 ETH gives generous
	// headroom without side effects — the value isn't asserted on anywhere.
	deployerInitialBalance, _ := new(big.Int).SetString("100000000000000000000", 10) // 100 ETH
	helper.sim = simulated.NewBackend(map[ethCommon.Address]ethTypes.Account{
		helper.Submitter.From:      {Balance: big.NewInt(9e18)},
		helper.Deployer.From:       {Balance: deployerInitialBalance},
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

		helper.autoMineWG.Add(1)
		go func() {
			defer helper.autoMineWG.Done()
			fmt.Println("Auto mining enabled")
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Shutting down auto mining")
					return
				case <-ticker.C:
					block := helper.sim.Commit()
					fmt.Println("Mined block: ", block)
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
		// ProcessorEndpoint.submitRequest reverts on Deploy — deploys must go
		// through submitDeployRequest. Fail here so the test points at the
		// caller rather than at a generic on-chain revert.
		require.Fail(s.t, "use SubmitDeployRequest for deploy requests; submitRequest reverts on-chain for Deploy")
		return nil
	case common.Process:
		reqType = 1
	case common.Deanonymize:
		reqType = 2
	case common.AssociateKey:
		reqType = 3
	default:
		panic("Unsupported request type")
	}

	// msg.value differs by token: ETH requests carry assetAmount + maxFeeValue
	// (business asset + fee in the same tx); ERC-20 requests carry maxFeeValue
	// only (business asset arrives via the contract's transferFrom pull).
	// Mirrors the check in ProcessorEndpoint.submitRequest.
	if tokenAddress == velacommon.NativeTokenAddress() {
		sender.Value = new(big.Int).Add(assetAmount, maxFeeValue)
	} else {
		sender.Value = new(big.Int).Set(maxFeeValue)
	}
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
	// Stop auto-mining FIRST and wait for the goroutine to exit, otherwise a
	// ticker firing between sim.Close() and ctx cancellation calls Commit() on
	// a nil SimulatedBeacon and panics.
	if s.autoMining {
		s.cancel()
		s.autoMineWG.Wait()
	}

	if s.sim != nil {
		err := s.sim.Close()
		require.NoError(s.t, err, "failed to close simulated backend")
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

// GrantDeployerRole grants the ProcessorEndpoint's DEPLOYER_ROLE to addr.
// Must be called from the ADMIN (set to s.Deployer in the constructor).
// Used by the wallet driver so a user's secp key can submit deploy requests.
func (s *SimTestHelper) GrantDeployerRole(addr ethCommon.Address) *ethTypes.Transaction {
	roleHash, err := bind.Call(s.processEndpointInstance,
		&bind.CallOpts{Pending: false},
		s.processEndpointContract.PackDEPLOYERROLE(),
		s.processEndpointContract.UnpackDEPLOYERROLE)
	require.NoError(s.t, err, "failed to read DEPLOYER_ROLE")

	tx, err := bind.Transact(
		s.processEndpointInstance,
		s.Deployer,
		s.processEndpointContract.PackGrantRole(roleHash, addr),
	)
	require.NoError(s.t, err, "failed to grant DEPLOYER_ROLE")
	return tx
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

// --- MockERC20 infrastructure ---

// DeployMockERC20 deploys a MockERC20 token contract on the simulated chain
// and stores its binding on the helper. Subsequent MintERC20 / ApproveERC20
// calls target this deployed instance. Returns the deployed token address.
//
// The mock's mint is permissionless and decimals are caller-chosen, so a
// single call is enough to stage any ERC-20 scenario. Tests needing multiple
// distinct tokens call this multiple times and track the returned addresses.
func (s *SimTestHelper) DeployMockERC20(name, symbol string, decimals uint8) ethCommon.Address {
	// A second call would silently overwrite the cached binding so that
	// MintERC20 / ApproveERC20 / BalanceOfERC20 target only the latest token,
	// breaking helpers for the previously deployed one. Fail fast instead.
	//
	// TODO: when a multi-token test legitimately hits this guard, drop the
	// cached MockERC20Address / mockERC20Contract / mockERC20Instance fields
	// and change the helpers to take a tokenAddress parameter, reconstructing
	// the binding internally. Callers track the returned addresses themselves.
	require.Nil(s.t, s.mockERC20Instance, "DeployMockERC20 already called")
	deployer := bind.DefaultDeployer(s.Deployer, s.sim.Client())
	contract := mockerc20.NewMockERC20()
	deployParams := bind.DeploymentParams{
		Contracts: []*bind.MetaData{&mockerc20.MockERC20MetaData},
		Inputs:    map[string][]byte{mockerc20.MockERC20MetaData.ID: contract.PackConstructor(name, symbol, decimals)},
	}
	deployRes, err := bind.LinkAndDeploy(&deployParams, deployer)
	require.NoError(s.t, err, "failed to deploy MockERC20")

	addr := deployRes.Addresses[mockerc20.MockERC20MetaData.ID]
	tx := deployRes.Txs[mockerc20.MockERC20MetaData.ID]
	s.sim.Commit()
	_, err = bind.WaitDeployed(context.Background(), s.sim.Client(), tx.Hash())
	require.NoError(s.t, err, "MockERC20 deployment not mined")

	s.MockERC20Address = addr
	s.mockERC20Contract = contract
	s.mockERC20Instance = contract.Instance(s.sim.Client(), addr)
	return addr
}

// MintERC20 mints `amount` tokens of the deployed MockERC20 to `to`. The
// mock's mint is permissionless; the Deployer account is used as msg.sender
// for consistency with other admin-style helpers.
func (s *SimTestHelper) MintERC20(to ethCommon.Address, amount *big.Int) *ethTypes.Transaction {
	require.NotNil(s.t, s.mockERC20Instance, "DeployMockERC20 must be called before MintERC20")
	tx, err := bind.Transact(
		s.mockERC20Instance,
		s.Deployer,
		s.mockERC20Contract.PackMint(to, amount),
	)
	require.NoError(s.t, err, "failed to mint ERC-20")
	return tx
}

// ApproveERC20 calls MockERC20.approve(spender, amount) from `owner`. The
// ProcessorEndpoint contract uses the resulting allowance to pull tokens via
// transferFrom when the user later calls submitRequest with an ERC-20
// assetAmount. Tests call this to pre-approve before a deposit, since the
// wallet's submitRequest path uses transferFrom and does not embed an
// EIP-2612 permit.
func (s *SimTestHelper) ApproveERC20(owner *bind.TransactOpts, spender ethCommon.Address, amount *big.Int) *ethTypes.Transaction {
	require.NotNil(s.t, s.mockERC20Instance, "DeployMockERC20 must be called before ApproveERC20")
	tx, err := bind.Transact(
		s.mockERC20Instance,
		owner,
		s.mockERC20Contract.PackApprove(spender, amount),
	)
	require.NoError(s.t, err, "failed to approve ERC-20")
	return tx
}

// BalanceOfERC20 reads the caller-visible balance for `account` on the deployed
// MockERC20. Useful for end-to-end assertions (e.g. confirming that a
// ClaimPendingPayments moved tokens back to the user).
func (s *SimTestHelper) BalanceOfERC20(account ethCommon.Address) *big.Int {
	require.NotNil(s.t, s.mockERC20Instance, "DeployMockERC20 must be called before BalanceOfERC20")
	bal, err := bind.Call(
		s.mockERC20Instance,
		&bind.CallOpts{Pending: false},
		s.mockERC20Contract.PackBalanceOf(account),
		s.mockERC20Contract.UnpackBalanceOf,
	)
	require.NoError(s.t, err, "failed to read ERC-20 balance")
	return bal
}

// --- ProcessorEndpoint ERC-20-adjacent helpers ---

// AddAllowedToken allowlists `tokenAddress` for ERC-20 deposits on the
// ProcessorEndpoint. Must be called before any submitRequest that targets the
// token, or the contract reverts with TokenNotAllowed. Uses the Deployer
// account (holder of the ADMIN role from the constructor).
func (s *SimTestHelper) AddAllowedToken(tokenAddress ethCommon.Address) *ethTypes.Transaction {
	tx, err := bind.Transact(
		s.processEndpointInstance,
		s.Deployer,
		s.processEndpointContract.PackAddAllowedToken(tokenAddress),
	)
	require.NoError(s.t, err, "failed to add allowed token")
	return tx
}

// GetPendingClaims reads pendingClaims[tokenAddress][payee] on the
// ProcessorEndpoint. Accumulates token refunds on state-update failures and
// withdrawal outputs; cleared to zero by Claim.
func (s *SimTestHelper) GetPendingClaims(tokenAddress, payee ethCommon.Address) *big.Int {
	amount, err := bind.Call(
		s.processEndpointInstance,
		&bind.CallOpts{Pending: false},
		s.processEndpointContract.PackPendingClaims(tokenAddress, payee),
		s.processEndpointContract.UnpackPendingClaims,
	)
	require.NoError(s.t, err, "failed to read pendingClaims")
	return amount
}

// Claim pulls any pending amount of `tokenAddress` for `payee` out of the
// ProcessorEndpoint and into the payee's balance. Payable by any caller; the
// Deployer account is used here so tests don't need to thread TransactOpts
// through for every test-side claim. For user-initiated claims, the wallet
// driver calls BlockChainClient.Claim from the user's own account.
func (s *SimTestHelper) Claim(tokenAddress, payee ethCommon.Address) *ethTypes.Transaction {
	tx, err := bind.Transact(
		s.processEndpointInstance,
		s.Deployer,
		s.processEndpointContract.PackClaim(tokenAddress, payee),
	)
	require.NoError(s.t, err, "failed to claim pending payments")
	return tx
}
