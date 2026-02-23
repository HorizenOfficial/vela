package blockchain

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/horizen-pes/pkg/blockchain/contracts/processorendpoint"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
)

//go:generate mkdir -p ../../contract_abis
//go:generate mkdir -p ./contracts/processorendpoint
//go:generate solc --via-ir --combined-json abi,bin ../../contracts/contracts/ProcessorEndpoint.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/ProcessorEndpointAbi --overwrite
//go:generate sh -c "jq --indent 2 '.contracts[\"contracts/contracts/ProcessorEndpoint.sol:ProcessorEndpoint\"].abi' ../../contract_abis/ProcessorEndpointAbi/combined.json > ../../subgraphs/hcce/abis/ProcessorEndpoint.json"
//go:generate sh -c "jq -r '.contracts[\"contracts/contracts/ProcessorEndpoint.sol:ProcessorEndpoint\"].abi' ../../contract_abis/ProcessorEndpointAbi/combined.json > ../../contract_abis/ProcessorEndpointAbi/ProcessorEndpoint.abi"
//go:generate sh -c "jq -r '.contracts[\"contracts/contracts/ProcessorEndpoint.sol:ProcessorEndpoint\"].bin' ../../contract_abis/ProcessorEndpointAbi/combined.json > ../../contract_abis/ProcessorEndpointAbi/ProcessorEndpoint.bin"
//go:generate abigen --v2 --abi ../../contract_abis/ProcessorEndpointAbi/ProcessorEndpoint.abi --bin ../../contract_abis/ProcessorEndpointAbi/ProcessorEndpoint.bin --pkg processorendpoint --type ProcessorEndpoint --out ./contracts/processorendpoint/ProcessorEndpoint.go

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

type BlockChainClient struct {
	mu                     sync.RWMutex
	connected              bool
	processorAddress       ethCommon.Address
	rpcURL                 string
	processorBoundContract *bind.BoundContract
	processorEndpoint      *processorendpoint.ProcessorEndpoint
	client                 ChainClient
	privKey                *cryptotypes.PrivateKeySecp256k1
	account                *bind.TransactOpts
}

func NewCoreBlockChainClient(processor ethCommon.Address, rpcURL string, key *cryptotypes.PrivateKeySecp256k1) *BlockChainClient {
	return &BlockChainClient{
		processorAddress:  processor,
		rpcURL:            rpcURL,
		processorEndpoint: processorendpoint.NewProcessorEndpoint(),
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

func (c *BlockChainClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("already connected")
	}

	var err error
	c.client, err = ethclient.Dial(c.rpcURL)
	if err != nil {
		return fmt.Errorf("cannot connect to chain: %w", err)
	}

	c.processorBoundContract = c.processorEndpoint.Instance(c.client, c.processorAddress)

	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve chain ID: %w", err)
	}

	if c.privKey != nil {
		c.account = bind.NewKeyedTransactor(c.privKey.PrivateKey, chainID)
	}

	c.connected = true
	return nil
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

// getPendingRequests gets pending requests from the blockchain.
// It is intentionally kept unexported because core flow uses GetNextPendingRequest.
func (c *BlockChainClient) getPendingRequests(ctx context.Context) ([]*common.Request, error) {
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
			DepositAmount:   common.ToBig(request.DepositAmount),
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
		DepositAmount:   common.ToBig(request.DepositAmount),
		MaxFeeValue:     common.ToBig(request.MaxFeeValue),
	}

	return req, stateRoot, nil
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

func (c *BlockChainClient) MarkRequestFailed(ctx context.Context, requestID common.RequestIdType, requestFailure *apperrors.RequestFailure) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}
	if c.account == nil {
		return fmt.Errorf("client not configured for signing transactions")
	}

	c.account.Value = nil

	if requestFailure == nil {
		requestFailure = apperrors.New(apperrors.CodeInternalFallback, "internal error", nil)
	}

	solCode := uint8(requestFailure.Category())
	msg := requestFailure.ExternalMessage()

	return c.sendTxAndWaitMined(ctx, c.processorEndpoint.PackMarkRequestFailed(requestID, solCode, msg))
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

	events := make([][]byte, len(update.Events))
	eventSubTypes := make([]string, len(update.Events))
	for i, event := range update.Events {
		events[i] = event.EncryptedData
		eventSubTypes[i] = event.EventSubType
	}

	withdrawals := make([]processorendpoint.StructsWithdrawalRequest, len(update.Withdrawals))
	for i, withdrawal := range update.Withdrawals {
		amount := withdrawal.Amount.ToInt()
		withdrawals[i] = processorendpoint.StructsWithdrawalRequest{
			Receiver: withdrawal.DestinationAddress,
			Amount:   amount,
		}
	}

	params := c.processorEndpoint.PackStateUpdate(
		processorendpoint.ApplicationIdToBindingType(update.ApplicationID),
		update.PrevStateRoot,
		update.NewStateRoot,
		update.RequestID,
		events,
		eventSubTypes,
		withdrawals,
		update.RefundAmount.ToInt(),
		update.ApplicationFee.ToInt(),
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
