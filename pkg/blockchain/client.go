package blockchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/processorendpoint"
	"github.com/HorizenOfficial/vela/pkg/blockchain/contracts/tee"
	"github.com/HorizenOfficial/vela/pkg/common"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	"github.com/HorizenOfficial/vela/pkg/crypto"
)

//go:generate mkdir -p ../../contract_abis
//go:generate mkdir -p ./contracts/processorendpoint
//go:generate solc --via-ir --optimize --combined-json abi,bin ../../contracts/contracts/ProcessorEndpoint.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/ProcessorEndpointAbi --overwrite
//go:generate sh -c "jq --indent 2 '.contracts[\"contracts/contracts/ProcessorEndpoint.sol:ProcessorEndpoint\"].abi' ../../contract_abis/ProcessorEndpointAbi/combined.json > ../../subgraphs/hcce/abis/ProcessorEndpoint.json"
//go:generate sh -c "jq -r '.contracts[\"contracts/contracts/ProcessorEndpoint.sol:ProcessorEndpoint\"].abi' ../../contract_abis/ProcessorEndpointAbi/combined.json > ../../contract_abis/ProcessorEndpointAbi/ProcessorEndpoint.abi"
//go:generate sh -c "jq -r '.contracts[\"contracts/contracts/ProcessorEndpoint.sol:ProcessorEndpoint\"].bin' ../../contract_abis/ProcessorEndpointAbi/combined.json > ../../contract_abis/ProcessorEndpointAbi/ProcessorEndpoint.bin"
//go:generate abigen --v2 --abi ../../contract_abis/ProcessorEndpointAbi/ProcessorEndpoint.abi --bin ../../contract_abis/ProcessorEndpointAbi/ProcessorEndpoint.bin --pkg processorendpoint --type ProcessorEndpoint --out ./contracts/processorendpoint/ProcessorEndpoint.go
//go:generate mkdir -p ./contracts/tee
//go:generate solc --via-ir --optimize --combined-json abi,bin ../../contracts/contracts/TeeAuthenticator.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/TeeAuthenticatorAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/TeeAuthenticatorAbi/combined.json --pkg tee --type TeeAuthenticator --out ./contracts/tee/TeeAuthenticator.go
//go:generate mkdir -p ./contracts/tokenallowlist
//go:generate solc --via-ir --optimize --combined-json abi,bin ../../contracts/contracts/TokenAllowlist.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/TokenAllowlistAbi --overwrite
//go:generate sh -c "jq -r '.contracts[\"contracts/contracts/TokenAllowlist.sol:TokenAllowlist\"].abi' ../../contract_abis/TokenAllowlistAbi/combined.json > ../../contract_abis/TokenAllowlistAbi/TokenAllowlist.abi"
//go:generate sh -c "jq -r '.contracts[\"contracts/contracts/TokenAllowlist.sol:TokenAllowlist\"].bin' ../../contract_abis/TokenAllowlistAbi/combined.json > ../../contract_abis/TokenAllowlistAbi/TokenAllowlist.bin"
//go:generate abigen --v2 --abi ../../contract_abis/TokenAllowlistAbi/TokenAllowlist.abi --bin ../../contract_abis/TokenAllowlistAbi/TokenAllowlist.bin --pkg tokenallowlist --type TokenAllowlist --out ./contracts/tokenallowlist/TokenAllowlist.go

type ChainClient interface {
	ethereum.BlockNumberReader
	ethereum.ChainReader
	ethereum.ChainStateReader
	ethereum.ContractCaller
	ethereum.GasEstimator
	ethereum.GasPricer
	ethereum.GasPricer1559
	ethereum.FeeHistoryReader
	ethereum.LogFilterer
	ethereum.PendingStateReader
	ethereum.PendingContractCaller
	ethereum.TransactionReader
	ethereum.TransactionSender
	ethereum.ChainIDReader
}

// defaultConnectTimeout is used when ConnectTimeout is zero.
const defaultConnectTimeout = 10 * time.Second

type BlockChainClient struct {
	mu                     sync.RWMutex
	connected              bool
	processorAddress       ethCommon.Address
	teeAuthAddress         ethCommon.Address
	rpcURL                 string
	connectTimeout         time.Duration
	processorBoundContract *bind.BoundContract
	processorEndpoint      *processorendpoint.ProcessorEndpoint
	teeAuthBoundContract   *bind.BoundContract
	teeAuthEndpoint        *tee.TeeAuthenticator
	client                 ChainClient
	privKey                *cryptotypes.PrivateKeySecp256k1
	account                *bind.TransactOpts
	chainID                *big.Int
}

func NewBlockChainClient(processor ethCommon.Address, teeAuthenticator ethCommon.Address, rpcURL string, key *cryptotypes.PrivateKeySecp256k1) *BlockChainClient {
	return &BlockChainClient{
		processorAddress:  processor,
		teeAuthAddress:    teeAuthenticator,
		rpcURL:            rpcURL,
		processorEndpoint: processorendpoint.NewProcessorEndpoint(),
		teeAuthEndpoint:   tee.NewTeeAuthenticator(),
		privKey:           key,
	}
}

// NewReadOnlyBlockChainClient builds a client configured only for read operations (no signing key required).
func NewReadOnlyBlockChainClient(processor ethCommon.Address, rpcURL string) *BlockChainClient {
	return &BlockChainClient{
		processorAddress:  processor,
		rpcURL:            rpcURL,
		processorEndpoint: processorendpoint.NewProcessorEndpoint(),
	}
}

// SetConnectTimeout overrides the default timeout for the dial and initial
// ChainID RPC call in Connect.
func (c *BlockChainClient) SetConnectTimeout(d time.Duration) error {
	if d <= 0 || d > 5*time.Minute {
		return fmt.Errorf("invalid connect timeout %v: must be greater than 0 and at most 5m", d)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connectTimeout = d
	return nil
}

// IsConnected returns true if the client has successfully connected to the blockchain.
func (c *BlockChainClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

func (c *BlockChainClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil // already connected, no-op
	}

	// Use a timeout so Connect does not hang indefinitely when the RPC
	// node is unreachable (go-ethereum's HTTP client retries forever
	// with a no-deadline context). The timeout covers both the dial and
	// the initial ChainID RPC call.
	timeout := c.connectTimeout
	if timeout == 0 {
		timeout = defaultConnectTimeout
	}
	connectCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Use DialContext so that the HTTP transport respects context cancellation.
	var err error
	c.client, err = ethclient.DialContext(connectCtx, c.rpcURL)
	if err != nil {
		return fmt.Errorf("cannot connect to chain: %w", err)
	}

	c.processorBoundContract = c.processorEndpoint.Instance(c.client, c.processorAddress)
	if c.teeAuthEndpoint != nil && c.teeAuthAddress != (ethCommon.Address{}) {
		c.teeAuthBoundContract = c.teeAuthEndpoint.Instance(c.client, c.teeAuthAddress)
	}

	chainID, err := c.client.ChainID(connectCtx)
	if err != nil {
		// Clean up so the next Connect() call starts fresh.
		// Close the underlying client if it supports it (ethclient.Client does).
		if closer, ok := c.client.(interface{ Close() }); ok {
			closer.Close()
		}
		c.client = nil
		c.processorBoundContract = nil
		c.teeAuthBoundContract = nil
		return fmt.Errorf("failed to retrieve chain ID: %w", err)
	}

	if c.privKey != nil {
		c.account = bind.NewKeyedTransactor(c.privKey.PrivateKey, chainID)
	}
	c.chainID = chainID

	c.connected = true
	return nil
}

// ChainID returns the connected chain ID.
func (c *BlockChainClient) ChainID(ctx context.Context) (*big.Int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("client not connected, call Connect first")
	}

	if c.chainID != nil {
		return new(big.Int).Set(c.chainID), nil
	}

	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chain ID: %w", err)
	}
	return chainID, nil
}

// LatestBlockNumber returns the latest block number from the chain.
func (c *BlockChainClient) LatestBlockNumber(ctx context.Context) (uint64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return 0, fmt.Errorf("client not connected, call Connect first")
	}
	return c.client.BlockNumber(ctx)
}

func (c *BlockChainClient) UnpackProcessorEndpointError(chainErr error) error {
	if chainErr == nil {
		return nil
	}

	raw, hasRevertErrorData := ethclient.RevertErrorData(chainErr)
	if !hasRevertErrorData || len(raw) == 0 {
		return fmt.Errorf("call returned error: %w", chainErr)
	}
	rawUnpackedErr, unpack_err := c.processorEndpoint.UnpackError(raw)
	if unpack_err != nil {
		return fmt.Errorf("call returned unknown error: %w", chainErr)
	}
	return fmt.Errorf("contract revert: %T", rawUnpackedErr)

}

func (c *BlockChainClient) UnpackProcessorEndpointErrorAndCheckForReorg(chainErr error) error {
	if chainErr == nil {
		return nil
	}

	if strings.Contains(chainErr.Error(), "nonce too low") {
		return ReorgError{causedBy: chainErr}
	}

	unpackedError := c.UnpackProcessorEndpointError(chainErr)

	if strings.Contains(unpackedError.Error(), "ProcessorEndpointInvalidApplicationId") ||
		strings.Contains(unpackedError.Error(), "ProcessorEndpointInvalidStateRoot") ||
		strings.Contains(unpackedError.Error(), "ProcessorEndpointInvalidRequestId") {
		return ReorgError{causedBy: unpackedError}
	}
	return unpackedError

}

// GetPendingRequests gets pending requests from the blockchain
func (c *BlockChainClient) GetPendingRequests(ctx context.Context) ([]*common.Request, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("client not connected, call Connect first")
	}

	listOfRequests, err := bind.Call(c.processorBoundContract,
		&bind.CallOpts{Pending: false},
		c.processorEndpoint.PackGetPendingRequests(),
		c.processorEndpoint.UnpackGetPendingRequests)

	if err != nil {
		return nil, c.UnpackProcessorEndpointErrorAndCheckForReorg(err)
	}

	output := make([]*common.Request, 0, len(listOfRequests))
	for _, request := range listOfRequests {
		req := &common.Request{
			ProtocolVersion: request.ProtocolVersion,
			ApplicationID:   processorendpoint.ApplicationIdFromBindingType(request.ApplicationId),
			RequestID:       request.RequestId,
			RequestType:     common.RequestType(request.RequestType),
			Payload:         request.Payload,
			Timestamp:       common.ToBig(request.Timestamp),
			Sender:          request.Sender,
			Facilitator:     request.Facilitator,
			TokenAddress:    request.TokenAddress,
			AssetAmount:     common.ToBig(request.AssetAmount),
			MaxFeeValue:     common.ToBig(request.MaxFeeValue),
		}

		output = append(output, req)
	}
	return output, nil
}

func (c *BlockChainClient) GetNextPendingRequest(ctx context.Context) (*common.Request, [32]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, [32]byte{}, fmt.Errorf("client not connected, call Connect first")
	}

	output, err := bind.Call(c.processorBoundContract,
		&bind.CallOpts{Pending: false},
		c.processorEndpoint.PackGetNextPendingRequest(),
		c.processorEndpoint.UnpackGetNextPendingRequest)

	if err != nil {
		return nil, [32]byte{}, c.UnpackProcessorEndpointErrorAndCheckForReorg(err)
	}

	stateRoot := output.Arg1
	if !output.Success {
		return nil, stateRoot, nil
	}

	request := output.Arg0

	req := &common.Request{
		ProtocolVersion: request.ProtocolVersion,
		ApplicationID:   processorendpoint.ApplicationIdFromBindingType(request.ApplicationId),
		RequestID:       common.RequestIdType(request.RequestId),
		RequestType:     common.RequestType(request.RequestType),
		Payload:         request.Payload,
		Timestamp:       common.ToBig(request.Timestamp),
		Sender:          request.Sender,
		Facilitator:     request.Facilitator,
		TokenAddress:    request.TokenAddress,
		AssetAmount:     common.ToBig(request.AssetAmount),
		MaxFeeValue:     common.ToBig(request.MaxFeeValue),
	}

	return req, stateRoot, nil
}

// GetPendingRequestsWithStateRoot fetches up to maxCount pending requests for the
// application selected by the contract, together with its applicationId and on-chain
// state root.
//
// STUB: the ProcessorEndpoint contract does not yet expose the batch view
// (per-application queues + round-robin selection, see docs/design/BATCH_EXECUTION.md
// section 4). Until it does, this delegates to GetNextPendingRequest, so a batch
// always contains at most one request. maxCount is accepted for API compatibility
// but ignored.
func (c *BlockChainClient) GetPendingRequestsWithStateRoot(ctx context.Context, maxCount uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
	req, stateRoot, err := c.GetNextPendingRequest(ctx)
	if err != nil {
		return 0, nil, [32]byte{}, err
	}
	if req == nil {
		return 0, nil, stateRoot, nil
	}
	return req.ApplicationID, []*common.Request{req}, stateRoot, nil
}

// SubmitBatchStateUpdate submits a batch of per-request update payloads together with
// a single batch signature covering all entry hashes.
//
// STUB: the ProcessorEndpoint contract does not yet expose batchStateUpdate() (see
// docs/design/BATCH_EXECUTION.md section 3.2). Until it does, this submits the batch
// through the single-request stateUpdate() path. This works only for a size-1 batch:
// the batch payloads are unsigned individually, but a 1-entry batch hashes identically
// to the single-request message (MsgToSignBuilder.BuildBatchMsgHash), so the batch
// signature verifies on-chain when attached to the lone payload. A batch of >1 cannot
// be replayed this way and is rejected loudly rather than silently dropping entries.
func (c *BlockChainClient) SubmitBatchStateUpdate(ctx context.Context, updates []*common.UpdatePayload, batchSignature []byte) error {
	if len(updates) != 1 {
		return fmt.Errorf("SubmitBatchStateUpdate stub supports only size-1 batches, got %d: batchStateUpdate() is not yet supported by the ProcessorEndpoint contract", len(updates))
	}
	// The single payload is unsigned; the batch signature (== single-request
	// signature for a 1-entry batch) is what stateUpdate() verifies on-chain.
	updates[0].Signature = batchSignature
	return c.SubmitStateUpdate(ctx, updates[0])
}

func (c *BlockChainClient) sendTxAndWaitMined(ctx context.Context, data []byte) error {
	if c.account == nil {
		return fmt.Errorf("client not configured for signing transactions")
	}

	tx, err := bind.Transact(c.processorBoundContract, c.account, data)
	if err != nil {
		return c.UnpackProcessorEndpointErrorAndCheckForReorg(err)
	}

	// wait for transaction inclusion
	receipt, err := bind.WaitMined(ctx, c.client, tx.Hash())
	if err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed")
	}
	return nil
}

// SubmitRequest submits a request to the ProcessorEndpoint smart contract using a common.Request.
func (c *BlockChainClient) SubmitRequest(ctx context.Context, protocolVersion uint8, applicationId common.ApplicationIdType, requestType common.RequestType, payload []byte, tokenAddress ethCommon.Address, assetAmount *big.Int, maxFeeValue *big.Int) (common.RequestIdType, uint64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return common.RequestIdType{}, 0, fmt.Errorf("client not connected, call Connect first")
	}
	if c.account == nil {
		return common.RequestIdType{}, 0, fmt.Errorf("client not configured for signing transactions")
	}

	reqType := uint8(requestType)

	// Pack the transaction data using the generated binding
	data := c.processorEndpoint.PackSubmitRequest(protocolVersion, processorendpoint.ApplicationIdToBindingType(applicationId), reqType, payload, tokenAddress, assetAmount, maxFeeValue)
	// Set the value for the transaction (msg.value).
	// For ETH requests: msg.value = assetAmount + maxFeeValue (carries both business asset and fee).
	// For ERC-20 requests: msg.value = maxFeeValue only (business asset arrives via transferFrom).
	if tokenAddress == velacommon.ETH_TOKEN {
		c.account.Value = new(big.Int).Add(assetAmount, maxFeeValue)
	} else {
		c.account.Value = new(big.Int).Set(maxFeeValue)
	}

	// Send the transaction
	tx, err := bind.Transact(c.processorBoundContract, c.account, data)
	c.account.Value = nil
	if err != nil {
		return common.RequestIdType{}, 0, fmt.Errorf("failed to submit transaction: %w", c.UnpackProcessorEndpointError(err))
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, c.client, tx.Hash())
	if err != nil {
		return common.RequestIdType{}, 0, fmt.Errorf("error waiting for tx inclusion: %w", err)
	}
	if receipt.Status != 1 {
		return common.RequestIdType{}, 0, fmt.Errorf("transaction failed")
	}

	// Parse the returned requestId from the transaction receipt logs
	for _, vLog := range receipt.Logs {
		event, err := c.processorEndpoint.UnpackRequestSubmittedEvent(vLog)
		if err == nil {
			return event.RequestId, receipt.BlockNumber.Uint64(), nil
		}
	}

	return common.RequestIdType{}, 0, fmt.Errorf("requestId not found in logs")
}

// SubmitDeployRequest submits a deploy request to the ProcessorEndpoint smart contract.
func (c *BlockChainClient) SubmitDeployRequest(ctx context.Context, protocolVersion uint8, payload []byte, maxFeeValue *big.Int) (common.ApplicationIdType, common.RequestIdType, uint64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return common.ApplicationIdType(0), common.RequestIdType{}, 0, fmt.Errorf("client not connected, call Connect first")
	}
	if c.account == nil {
		return common.ApplicationIdType(0), common.RequestIdType{}, 0, fmt.Errorf("client not configured for signing transactions")
	}
	// Pack the transaction data using the generated binding
	data := c.processorEndpoint.PackSubmitDeployRequest(protocolVersion, payload)
	// Set the value for the transaction (msg.value = maxFeeValue)
	c.account.Value = new(big.Int).Set(maxFeeValue)

	// Send the transaction
	tx, err := bind.Transact(c.processorBoundContract, c.account, data)
	c.account.Value = nil
	if err != nil {
		return common.ApplicationIdType(0), common.RequestIdType{}, 0, fmt.Errorf("failed to submit transaction: %w", c.UnpackProcessorEndpointError(err))
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, c.client, tx.Hash())
	if err != nil {
		return common.ApplicationIdType(0), common.RequestIdType{}, 0, fmt.Errorf("error waiting for tx inclusion: %w", err)
	}
	if receipt.Status != 1 {
		return common.ApplicationIdType(0), common.RequestIdType{}, 0, fmt.Errorf("transaction failed")
	}

	// Parse the returned requestId and applicationId from the transaction receipt logs
	for _, vLog := range receipt.Logs {
		event, err := c.processorEndpoint.UnpackDeployRequestSubmittedEvent(vLog)
		if err == nil {
			return common.NewApplicationId(event.ApplicationId), event.RequestId, receipt.BlockNumber.Uint64(), nil
		}
	}

	return common.ApplicationIdType(0), common.RequestIdType{}, 0, fmt.Errorf("requestId not found in logs")
}

func (c *BlockChainClient) SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}
	if c.account == nil {
		return fmt.Errorf("client not configured for signing transactions")
	}
	userEvents := make([][]byte, len(update.Events))
	userEventSubTypes := make([][32]byte, len(update.Events))
	for i, event := range update.Events {
		userEvents[i] = event.EncryptedData
		userEventSubTypes[i] = event.EventSubType
	}
	userEventData := processorendpoint.StructsEventData{
		Events:   userEvents,
		SubTypes: userEventSubTypes,
	}

	appEvents := make([][]byte, len(update.AppEvents))
	appEventSubTypes := make([][32]byte, len(update.AppEvents))
	for i, appEvent := range update.AppEvents {
		appEvents[i] = appEvent.Data
		appEventSubTypes[i] = appEvent.EventSubType
	}
	appEventData := processorendpoint.StructsEventData{
		Events:   appEvents,
		SubTypes: appEventSubTypes,
	}

	withdrawals := make([]processorendpoint.StructsWithdrawalRequest, len(update.Withdrawals))
	for i, withdrawal := range update.Withdrawals {
		amount := withdrawal.Amount.ToInt()
		withdrawals[i] = processorendpoint.StructsWithdrawalRequest{
			TokenAddress: withdrawal.TokenAddress,
			Receiver:     withdrawal.DestinationAddress,
			Amount:       amount,
		}
	}

	params := c.processorEndpoint.PackStateUpdate(
		processorendpoint.ApplicationIdToBindingType(update.ApplicationID),
		update.PrevStateRoot,
		update.NewStateRoot,
		update.RequestID,
		userEventData,
		appEventData,
		withdrawals,
		update.RefundAmount.ToInt(),
		update.ApplicationFee.ToInt(),
		update.ErrorCode,
		update.ErrorMsg,
		update.Signature,
	)

	c.account.Value = nil
	return c.sendTxAndWaitMined(ctx, params)

}

// Close closes the blockchain client
func (c *BlockChainClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.client = nil
	c.processorBoundContract = nil
	c.account = nil

	c.connected = false
	return nil
}

// GetPendingPayments returns the pending payment balance for the given address.
func (c *BlockChainClient) GetPendingClaims(ctx context.Context, tokenAddress ethCommon.Address, addr ethCommon.Address) (*big.Int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("client not connected, call Connect first")
	}

	amount, err := bind.Call(c.processorBoundContract,
		&bind.CallOpts{Pending: false},
		c.processorEndpoint.PackPendingClaims(tokenAddress, addr),
		c.processorEndpoint.UnpackPendingClaims)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve pending claims: %w", err)
	}
	return amount, nil
}

// Claim calls claim on the ProcessorEndpoint contract for the given token and payee.
func (c *BlockChainClient) Claim(ctx context.Context, tokenAddress ethCommon.Address, payee ethCommon.Address) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.account == nil {
		return fmt.Errorf("client not configured for signing transactions")
	}
	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}

	c.account.Value = nil
	return c.sendTxAndWaitMined(ctx, c.processorEndpoint.PackClaim(tokenAddress, payee))
}

func (c *BlockChainClient) GetTeePublicKey(ctx context.Context) (*cryptotypes.PublicKeyP521, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.connected {
		return nil, fmt.Errorf("client not connected, call Connect first")
	}
	if c.teeAuthBoundContract == nil || c.teeAuthEndpoint == nil {
		return nil, fmt.Errorf("tee authenticator contract not configured")
	}

	pubSecp521r1, err := bind.Call(c.teeAuthBoundContract,
		&bind.CallOpts{Pending: false},
		c.teeAuthEndpoint.PackGetPubSecp521r1(),
		c.teeAuthEndpoint.UnpackGetPubSecp521r1)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve pubSecp521r1: %w", err)
	}
	return crypto.ImportPublicKeyP521FromHex(hex.EncodeToString(pubSecp521r1))
}
