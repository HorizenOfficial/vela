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

// ChainID is the chain ID used by the simulated backend. Exposed so tests can
// pass it to signature-building helpers that must match what the authority
// service enforces on /getreport.
var ChainID = params.AllDevChainProtocolChanges.ChainID

// FullStackSystemTestSuite is the real-chain system test suite. It uses a
// go-ethereum simulated.Backend via SimTestHelper for real contract interactions.
// Requests are submitted on-chain, the manager polls the chain and processes them,
// and the executor runs WASM — the full production path minus network and TEE attestation.
type FullStackSystemTestSuite struct {
	*testutil.TestSuiteCore                                     // shared infrastructure
	simHelper     *blockchainTestutil.SimTestHelper              // simulated chain + contracts
	wrappedClient *eventBroadcastingClient                      // intercepts SubmitStateUpdate
	subgraphImpl  *InProcessSubgraph                            // in-process subgraph for wallet queries
	authority     *InProcessAuthority                           // in-process authority HTTP service
	userAccounts  map[ethCommon.Address]*bind.TransactOpts       // funded test accounts
	eventChannel  chan interface{}                              // stashed so RestartAll can rebuild the wrappedClient with the same sink
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
	return newFullStackSuiteInternal(t, appType, mgrConfig, execConfig, keySet, recoveryData, mgrLogCfg, excLogCfg, nil)
}

// NewFullStackSystemTestSuiteWithTeeSigner is the attestation-rejection variant
// constructor: it registers teeSignerOverride as the contract's expected TEE
// signer while the executor still uses its REAL signing key. Because the two
// addresses differ, every on-chain stateUpdate's signature recovers to an
// address that does NOT match the registered signer, and
// teeAuthenticator.checkSignature returns false, which triggers
// ProcessorEndpoint's InvalidSignature revert. Use this to prove the
// TEE-verify wiring is actually enforced end-to-end — a regression that
// removed or short-circuited the check would pass every other fullstack test
// but fail this one.
//
// The executor's real comm pubkey is still registered on-chain, so wallet
// payload encryption works normally and the test exercises the rejection
// path (not a pre-flight encryption failure).
func NewFullStackSystemTestSuiteWithTeeSigner(
	t *testing.T,
	appType string,
	mgrLogCfg *logger.Config,
	excLogCfg *logger.Config,
	teeSignerOverride *ethCommon.Address,
) *FullStackSystemTestSuite {
	mgrConfig, err := manager.LoadConfig()
	require.NoError(t, err)
	execConfig, err := executor.LoadConfig()
	require.NoError(t, err)
	ctx := context.Background()
	keySet, newRecoveryData, err := executor.GenerateEnclaveKeySet(ctx, execConfig.KeySetRecoveryType, nil, nil, "")
	require.NoError(t, err)
	return newFullStackSuiteInternal(t, appType, mgrConfig, execConfig, keySet, newRecoveryData, mgrLogCfg, excLogCfg, teeSignerOverride)
}

// newFullStackSuiteInternal is the shared body for the public constructors.
// teeSignerOverride being nil means "register the real signer derived from
// the executor's signing key" (the normal path); non-nil means "register
// this override instead, while still using the real keyset for the
// executor" — the attestation-rejection path.
func newFullStackSuiteInternal(
	t *testing.T,
	appType string,
	mgrConfig *manager.Config,
	execConfig *executor.Config,
	keySet *executor.EnclaveKeySet,
	recoveryData *common.EnclaveKeySetRecovery,
	mgrLogCfg *logger.Config,
	excLogCfg *logger.Config,
	teeSignerOverride *ethCommon.Address,
) *FullStackSystemTestSuite {
	// Deploy the mock TEE authenticator pre-populated with the executor's real
	// signing address and P521 communication pubkey. The MockTeeAuthenticator
	// has no setter (its constructor is the only way to set pubSecp521r1), so
	// the contract state must be correct from construction. Without this,
	// GetTeePublicKey returns 133 zero bytes and wallet flows that encrypt
	// payloads (RegisterUser, Deposit) fail with "invalid public key".
	teeSignerAddr := ethCrypto.PubkeyToAddress(*keySet.SigningKey.PublicKey().PublicKey)
	// Default: deploy the MockTeeAuthenticator (checkSignature always
	// returns true — convenient for the happy-path suite). The rejection
	// variant needs the NoAttestationTeeAuthenticator, which inherits the
	// real ECDSA recover+compare logic, so we switch contracts when an
	// override signer is supplied.
	useMockTeeAuth := true
	if teeSignerOverride != nil {
		// The contract is told to expect the override, but the executor
		// keeps its real signing key — so every stateUpdate signature
		// recovers to an address that does NOT match the registered signer,
		// and the authenticator rejects it. The real comm pubkey is kept
		// so payload encryption still works (we want to reach the sig
		// check, not fail earlier at encrypt).
		teeSignerAddr = *teeSignerOverride
		useMockTeeAuth = false
	}
	teePubKeyBytes := keySet.CommunicationKey.PublicKey().Bytes()

	// Create SimTestHelper with auto-mining enabled (mines every 1s)
	autoMining := true
	simHelper := blockchainTestutil.NewSimTestHelper(t, autoMining, useMockTeeAuth, &teeSignerAddr, teePubKeyBytes)

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

	// Enable reports so the manager writes deanonymization reports into a temp
	// dir that the in-process authority service can read from.
	testutil.EnableReports(mgrConfig)

	// Build the shared core infrastructure, passing the wrapped client as blockchain.Client
	core := testutil.NewTestSuiteCore(t, appType, mgrConfig, execConfig, wrappedClient, keySet, recoveryData, mgrLogCfg, excLogCfg)

	// Override the event channel — NewTestSuiteCore created one, but we want
	// the one connected to the wrappedClient
	core.SetEventChannel(eventChannel)

	// Start the in-process authority service. It shares reports and artifacts
	// dirs with the manager (so wallet-uploaded WASM lands where the manager
	// reads from) and queries the in-process subgraph for completion status.
	authority := NewInProcessAuthority(t, ChainID.Uint64(), core.GetReportsPath(), core.GetArtifactsPath(), subgraphImpl)

	return &FullStackSystemTestSuite{
		TestSuiteCore: core,
		simHelper:     simHelper,
		wrappedClient: wrappedClient,
		subgraphImpl:  subgraphImpl,
		authority:     authority,
		userAccounts:  make(map[ethCommon.Address]*bind.TransactOpts),
		eventChannel:  eventChannel,
	}
}

// RestartAll stops and rebuilds the manager + executor (via TestSuiteCore)
// AND rebuilds the fullstack-specific blockchain client stack (a fresh
// BlockChainClient re-dialling the same simulated backend, wrapped in a
// fresh eventBroadcastingClient sharing the original event channel and
// subgraph sink). Used by tests that verify cross-restart keyset recovery.
//
// The simulated chain, subgraph, authority, SimTestHelper, and
// user-accounts map are NOT rebuilt — their state is valid across the
// restart and callers typically want to inspect pre/post state through
// the same handles.
func (s *FullStackSystemTestSuite) RestartAll() error {
	// Fresh blockchain client: manager.Stop() closed the previous one by
	// design. Build a new BlockChainClient dialing the same simulated
	// backend (same contract addresses, same manager account), then wrap
	// with a fresh eventBroadcastingClient that preserves the existing
	// event channel and subgraph sink.
	realClient := blockchain.SetupNewBlockChainClientConnected(
		s.simHelper.Client(),
		s.simHelper.ProcessorContractAddress,
		s.simHelper.TeeSignerAddress,
		s.simHelper.ManagerAccount,
	)
	freshWrapped := newEventBroadcastingClient(realClient, s.eventChannel)
	freshWrapped.onStateUpdate = s.subgraphImpl.RecordStateUpdate
	s.wrappedClient = freshWrapped

	return s.TestSuiteCore.RestartCore(freshWrapped)
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

	// Fund the account with 5 ETH. Wait for the tx to be mined before
	// returning — otherwise the caller's first user-signed tx can hit
	// "insufficient funds for transfer" because the funding is still pending.
	fundAmount := new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	s.simHelper.WaitMined(s.simHelper.TransferFunds(s.simHelper.Deployer, addr, fundAmount))

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

	// Fund the account with 5 ETH. Wait for the tx to be mined before
	// returning — otherwise the caller's first user-signed tx can hit
	// "insufficient funds for transfer" because the funding is still pending.
	fundAmount := new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	s.simHelper.WaitMined(s.simHelper.TransferFunds(s.simHelper.Deployer, addr, fundAmount))

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

// HasAccount reports whether the given address has been registered and funded
// via CreateFundedAccount / RegisterAccount. Unambiguous boolean — use this
// instead of error-testing GetTransactOpts when the caller only wants to
// branch on registration state.
func (s *FullStackSystemTestSuite) HasAccount(addr ethCommon.Address) bool {
	_, exists := s.userAccounts[addr]
	return exists
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
	s.wrappedClient.markPending(req.RequestID)

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

// WaitForStateUpdateError blocks up to `timeout` for the underlying
// blockchain client's SubmitStateUpdate to return an error (typically an
// on-chain revert reason). The second return value reports whether an
// error arrived in that window. Intended for negative-path tests that
// need to pin down the specific revert reason rather than waiting for an
// outer polling timeout.
func (s *FullStackSystemTestSuite) WaitForStateUpdateError(timeout time.Duration) (error, bool) {
	return s.wrappedClient.waitForStateUpdateError(timeout)
}

// AssertNoStateUpdateErrors fails the test if any SubmitStateUpdate errored
// during its execution. Use as a defensive safety net in happy-path tests:
// it catches silent retry-and-recover regressions where the manager hits a
// transient stateUpdate revert and retries until success, leaving outer
// (balance/custody) assertions green but the protocol invariant broken.
// Non-blocking: reads only what's already buffered, returns immediately.
func (s *FullStackSystemTestSuite) AssertNoStateUpdateErrors(t *testing.T) {
	t.Helper()
	select {
	case err := <-s.wrappedClient.stateUpdateErrors:
		t.Fatalf("unexpected stateUpdate error during happy-path test "+
			"(indicates a silent retry-and-recover regression): %v", err)
	default:
	}
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

// GetBlockchainClient returns the event-broadcasting blockchain client wired to
// the simulated backend. Pass this into wallet command constructors so they
// share the same chain view as the manager — requests submitted here land in
// the same pool the manager polls, and the wrapper intercepts SubmitStateUpdate
// (manager → chain) to feed the in-process subgraph.
func (s *FullStackSystemTestSuite) GetBlockchainClient() blockchain.Client {
	return s.wrappedClient
}

// GetSubgraph returns the concrete *InProcessSubgraph, giving tests access to
// helpers not part of the subgraph.Client interface (e.g. InjectRequestCompleted
// used by the authority smoke test).
func (s *FullStackSystemTestSuite) GetSubgraph() *InProcessSubgraph {
	return s.subgraphImpl
}

// GetAuthorityServiceURL returns the base URL of the in-process authority
// service (e.g. http://127.0.0.1:xxxxx). Use with wallet requestreport/downloadreport.
func (s *FullStackSystemTestSuite) GetAuthorityServiceURL() string {
	return s.authority.URL()
}

// GetAuthorityAddress returns the Ethereum address of the in-process authority.
// When tests write synthetic DeanonymizationReport files, the Authority field
// must equal this address for /getreport to succeed.
func (s *FullStackSystemTestSuite) GetAuthorityAddress() ethCommon.Address {
	return s.authority.AuthorityAddress()
}

// GetAuthorityKey returns the secp256k1 private key matching GetAuthorityAddress.
// Tests use it to sign /getreport challenges.
func (s *FullStackSystemTestSuite) GetAuthorityKey() *ecdsa.PrivateKey {
	return s.authority.AuthorityKey()
}

// GrantDeployerRole grants the ProcessorEndpoint DEPLOYER_ROLE to addr so it
// can submit deploy requests on its own behalf. Used by the wallet driver to
// authorize a freshly-created user account to deploy applications.
func (s *FullStackSystemTestSuite) GrantDeployerRole(addr ethCommon.Address) {
	tx := s.simHelper.GrantDeployerRole(addr)
	s.simHelper.WaitMined(tx)
}

// RegisterAuthority allowlists the in-process authority for the given application
// on-chain via DefaultAuthority.addAllowedAuthority. Call after DeployApp;
// required before any wallet flow that submits an on-chain Deanonymize request
// (the ProcessorEndpoint contract checks this allowlist).
//
// Not needed for direct HTTP tests against /nonce and /getreport — those paths
// only verify the signature against the report file and subgraph state.
func (s *FullStackSystemTestSuite) RegisterAuthority(appID common.ApplicationIdType) {
	tx := s.simHelper.AddAuthority(new(big.Int).SetUint64(uint64(appID)), s.authority.AuthorityAddress())
	s.simHelper.WaitMined(tx)
}

// GetSimTestHelper returns the underlying SimTestHelper for advanced test-side operations.
func (s *FullStackSystemTestSuite) GetSimTestHelper() *blockchainTestutil.SimTestHelper {
	return s.simHelper
}

// Cleanup tears down the fullstack suite: stops the authority service, manager,
// executor, cleans temp dirs, and closes the simulated backend.
func (s *FullStackSystemTestSuite) Cleanup() error {
	if s.authority != nil {
		if err := s.authority.Close(); err != nil {
			return err
		}
	}
	s.CleanupCore()
	s.simHelper.Close()
	return nil
}
