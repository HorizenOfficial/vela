package fullstack

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela-common-go/subgraph"
	"github.com/HorizenOfficial/vela/pkg/blockchain"
	blockchainTestutil "github.com/HorizenOfficial/vela/pkg/blockchain/testutil"
	"github.com/HorizenOfficial/vela/pkg/common"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	"github.com/HorizenOfficial/vela/pkg/executor"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/manager"
	"github.com/HorizenOfficial/vela/pkg/testutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

// FullStackSystemTestSuite is the real-chain system test suite. It uses a
// go-ethereum simulated.Backend via SimTestHelper for real contract interactions.
// Requests are submitted on-chain, the manager polls the chain and processes them,
// and the executor runs WASM — the full production path minus network and TEE attestation.
type FullStackSystemTestSuite struct {
	*testutil.TestSuiteCore                                     // shared infrastructure
	simHelper     *blockchainTestutil.SimTestHelper              // simulated chain + contracts
	wrappedClient *eventBroadcastingClient                      // intercepts SubmitStateUpdate
	subgraphImpl  *InProcessSubgraph                            // in-process subgraph for wallet queries
	userAccounts  map[ethCommon.Address]*bind.TransactOpts       // funded test accounts
}

// NewFullStackSystemTestSuite creates a system test suite backed by a real
// simulated Ethereum chain.
func NewFullStackSystemTestSuite(t *testing.T, appType string, mgrLogCfg *logger.Config, excLogCfg *logger.Config) *FullStackSystemTestSuite {
	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)
	ctx := context.Background()
	keySet, newRecoveryData, err := executor.GenerateEnclaveKeySet(ctx, execConfig.KeySetRecoveryType, nil, nil, "")
	require.NoError(t, err)
	return NewFullStackSystemTestSuiteWithConfigs(t, appType, mgrConfig, execConfig, keySet, newRecoveryData, mgrLogCfg, excLogCfg)
}

func NewFullStackSystemTestSuiteWithConfigs(
	t *testing.T,
	appType string,
	mgrConfig *manager.Config,
	execConfig *executor.Config,
	keySet *executor.EnclaveKeySet,
	recoveryData *common.EnclaveKeySetRecovery,
	mgrLogCfg *logger.Config,
	excLogCfg *logger.Config,
) *FullStackSystemTestSuite {
	// Create SimTestHelper with auto-mining enabled (mines every 1s)
	simHelper := blockchainTestutil.NewSimTestHelper(t, true, true, nil, nil)

	// Create a real BlockChainClient connected to the simulated backend
	realClient := blockchain.SetupNewBlockChainClientConnected(
		simHelper.Client(),
		simHelper.ProcessorContractAddress,
		simHelper.TeeSignerAddress,
		simHelper.ManagerAccount,
	)

	// Create in-process subgraph
	subgraphImpl := NewInProcessSubgraph()

	// Create event channel and wrap the client
	eventChannel := make(chan interface{}, 100)
	wrappedClient := newEventBroadcastingClient(realClient, eventChannel)

	// Wire the subgraph to receive state updates from the manager
	wrappedClient.onStateUpdate = subgraphImpl.RecordStateUpdate

	// Reduce blockchain polling interval for faster test execution
	mgrConfig.BlockchainPollingInterval = 1

	// Build the shared core infrastructure, passing the wrapped client as blockchain.Client
	core := testutil.NewTestSuiteCore(t, appType, mgrConfig, execConfig, wrappedClient, keySet, recoveryData, mgrLogCfg, excLogCfg)

	// Override the event channel — NewTestSuiteCore created one, but we want
	// the one connected to the wrappedClient
	core.SetEventChannel(eventChannel)

	return &FullStackSystemTestSuite{
		TestSuiteCore: core,
		simHelper:     simHelper,
		wrappedClient: wrappedClient,
		subgraphImpl:  subgraphImpl,
		userAccounts:  make(map[ethCommon.Address]*bind.TransactOpts),
	}
}

// CreateFundedAccount generates a new secp256k1 key, funds the derived address
// with ETH, and registers it for on-chain transactions via SubmitRequest.
// It returns the address and the raw private key so the caller can register
// it with CryptoHelper for seed computation and payload encryption.
func (s *FullStackSystemTestSuite) CreateFundedAccount() (ethCommon.Address, *cryptotypes.PrivateKeySecp256k1, error) {
	privKey, err := ethCrypto.GenerateKey()
	if err != nil {
		return ethCommon.Address{}, nil, fmt.Errorf("failed to generate key: %w", err)
	}

	chainID := params.AllDevChainProtocolChanges.ChainID
	account := bind.NewKeyedTransactor(privKey, chainID)
	addr := account.From

	// Fund the account with 5 ETH
	fundAmount := new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	s.simHelper.TransferFunds(s.simHelper.Deployer, addr, fundAmount)

	s.userAccounts[addr] = account
	return addr, &cryptotypes.PrivateKeySecp256k1{PrivateKey: privKey}, nil
}

// RegisterAccount registers an existing private key as a funded account.
// Use this when you need a specific key (e.g., from CryptoHelper) to also
// be usable for on-chain transactions.
func (s *FullStackSystemTestSuite) RegisterAccount(privKey *ecdsa.PrivateKey) ethCommon.Address {
	chainID := params.AllDevChainProtocolChanges.ChainID
	account := bind.NewKeyedTransactor(privKey, chainID)
	addr := account.From

	// Fund the account with 5 ETH
	fundAmount := new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	s.simHelper.TransferFunds(s.simHelper.Deployer, addr, fundAmount)

	s.userAccounts[addr] = account
	return addr
}

// GetTransactOpts returns the TransactOpts for a previously created funded account.
func (s *FullStackSystemTestSuite) GetTransactOpts(addr ethCommon.Address) (*bind.TransactOpts, error) {
	opts, exists := s.userAccounts[addr]
	if !exists {
		return nil, fmt.Errorf("no funded account for address %s — call CreateFundedAccount first", addr)
	}
	return opts, nil
}

// SubmitRequest submits a request on-chain via the simulated backend.
// Unlike the mock suite, the contract assigns the requestID (and applicationID
// for deploy requests). This method updates req.RequestID and req.ApplicationID
// in place with the contract-assigned values.
func (s *FullStackSystemTestSuite) SubmitRequest(req *common.Request) error {
	sender, err := s.GetTransactOpts(req.Sender)
	if err != nil {
		return err
	}

	switch req.RequestType {
	case common.Deploy:
		tx := s.simHelper.SubmitDeployRequestFromUser(req.Payload, req.MaxFeeValue.ToInt(), sender)
		s.simHelper.WaitMined(tx)
		event := s.simHelper.GetDeployRequestSubmittedEvent(tx)
		req.RequestID = event.RequestId
		req.ApplicationID = common.NewApplicationId(event.ApplicationId)

	default:
		assetAmount := big.NewInt(0)
		if req.AssetAmount != nil {
			assetAmount = req.AssetAmount.ToInt()
		}

		tx := s.simHelper.SubmitRequestFromUser(
			req.ApplicationID,
			req.RequestType,
			req.Payload,
			req.TokenAddress,
			assetAmount,
			req.MaxFeeValue.ToInt(),
			sender,
		)
		s.simHelper.WaitMined(tx)
		event := s.simHelper.GetRequestSubmittedEvent(tx)
		req.RequestID = event.RequestId
	}

	// Register as pending in the wrapper for completion tracking
	s.wrappedClient.markPending(req.RequestID, req.RequestType == common.Deploy)

	return nil
}

// AssertRequestCompleted polls until the manager has processed the request
// (observed via the SubmitStateUpdate interception in the wrapped client).
func (s *FullStackSystemTestSuite) AssertRequestCompleted(requestID common.RequestIdType, timeout time.Duration) error {
	return s.wrappedClient.waitForRequestCompletion(requestID, timeout)
}

// WaitForAppStateInBlockchain polls the on-chain state root until it becomes non-zero
// for the given application (indicating the app has been deployed and state updated).
func (s *FullStackSystemTestSuite) WaitForAppStateInBlockchain(appID common.ApplicationIdType, timeout time.Duration) (*common.ApplicationState, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			stateRoot := s.simHelper.GetStateRoot(appID)
			if stateRoot != [32]byte{} {
				return &common.ApplicationState{
					ApplicationID: appID,
					StateRoot:     stateRoot,
				}, nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for app state %d in blockchain", appID)
		}
	}
}

// WaitForWithdrawal waits for a withdrawal to be recorded by the manager's
// SubmitStateUpdate (tracked via the event-broadcasting wrapper).
func (s *FullStackSystemTestSuite) WaitForWithdrawal(appID common.ApplicationIdType, timeout time.Duration) (*common.Withdrawal, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			withdrawals, exists := s.wrappedClient.getWithdrawals(appID)
			if exists && len(withdrawals) > 0 {
				return &withdrawals[0], nil
			}
		case <-timeoutCh:
			return nil, fmt.Errorf("timeout waiting for withdrawal for app %d", appID)
		}
	}
}

// GetRequestUpdatePayload returns the UpdatePayload for a completed request,
// as intercepted by the event-broadcasting wrapper.
func (s *FullStackSystemTestSuite) GetRequestUpdatePayload(reqId common.RequestIdType) (*common.UpdatePayload, error) {
	return s.wrappedClient.getUpdatePayload(reqId)
}

// GetAppCustody queries the on-chain appCustody for an application and token.
func (s *FullStackSystemTestSuite) GetAppCustody(appID common.ApplicationIdType, tokenAddress ethCommon.Address) *big.Int {
	return s.simHelper.GetAppCustody(appID, tokenAddress)
}

// GetDeployerAddress returns the address of the pre-authorized deployer account.
// Use this as the Sender for deploy requests, since the ProcessorEndpoint
// contract requires deployer authorization.
func (s *FullStackSystemTestSuite) GetDeployerAddress() ethCommon.Address {
	// Register the Deployer in userAccounts if not already present
	addr := s.simHelper.Deployer.From
	if _, exists := s.userAccounts[addr]; !exists {
		s.userAccounts[addr] = s.simHelper.Deployer
	}
	return addr
}

// GetSubgraphClient returns the in-process subgraph as a subgraph.Client interface.
// Pass this to wallet commands that need subgraph queries (getprivatebalance, etc.).
func (s *FullStackSystemTestSuite) GetSubgraphClient() subgraph.Client {
	return s.subgraphImpl
}

// GetSimTestHelper returns the underlying SimTestHelper for advanced test-side operations.
func (s *FullStackSystemTestSuite) GetSimTestHelper() *blockchainTestutil.SimTestHelper {
	return s.simHelper
}

// Cleanup tears down the fullstack suite: stops manager/executor, cleans temp
// dirs, and closes the simulated backend.
func (s *FullStackSystemTestSuite) Cleanup() error {
	s.CleanupCore()
	s.simHelper.Close()
	return nil
}
