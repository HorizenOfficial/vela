package blockchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/horizen-pes/pkg/blockchain/contracts/processorendpoint"
	"github.com/horizen-pes/pkg/blockchain/contracts/tee"
	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/crypto"
)

//go:generate mkdir -p ../../contract_abis
//go:generate mkdir -p ./contracts/processorendpoint
//go:generate solc --combined-json abi,bin ../../contracts/contracts/ProcessorEndpoint.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/ProcessorEndpointAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/ProcessorEndpointAbi/combined.json --pkg processorendpoint --type ProcessorEndpoint --out ./contracts/processorendpoint/ProcessorEndpoint.go

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
	teeAuthAddress         ethCommon.Address
	rpcURL                 string
	processorBoundContract *bind.BoundContract
	processorEndpoint      *processorendpoint.ProcessorEndpoint
	teeAuthBoundContract   *bind.BoundContract
	teeAuthEndpoint        *tee.TeeAuthenticator
	client                 ChainClient
	privKey                *cryptotypes.PrivateKeySecp256k1
	account                *bind.TransactOpts
}

func toRequestType(i uint8) common.RequestType {
	switch i {
	case 0:
		return common.Deploy
	case 1:
		return common.Process
	case 2:
		return common.Deanonymize
	case 3:
		return common.AssociateKey
	default:
		return ""
	}
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

	c.account = bind.NewKeyedTransactor(c.privKey.PrivateKey, chainID)

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
		//TODO check that all big.Int can fit in a Uint64. If not, the specific request should be marked as failed
		requestId := hex.EncodeToString(request.RequestId[:])
		req := &common.Request{
			ProtocolVersion: strconv.FormatUint(uint64(request.ProtocolVersion), 10),
			ApplicationID:   request.ApplicationId.String(),
			RequestID:       requestId,
			RequestType:     toRequestType(request.RequestType),
			Payload:         request.Payload,
			Timestamp:       request.Timestamp.Int64(),
			Sender:          request.Sender.String(),
			Value:           request.Value.Uint64(),
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

	requestId := hex.EncodeToString(request.RequestId[:])
	//TODO check that all big.Int can fit in a Uint64. If not, the specific request should be marked as failed
	req := &common.Request{
		ProtocolVersion: strconv.FormatUint(uint64(request.ProtocolVersion), 10),
		ApplicationID:   request.ApplicationId.String(),
		RequestID:       requestId,
		RequestType:     toRequestType(request.RequestType),
		Payload:         request.Payload,
		Timestamp:       request.Timestamp.Int64(),
		Sender:          request.Sender.String(),
		Value:           request.Value.Uint64(),
	}

	return req, stateRoot, nil
}

func (c *BlockChainClient) sendTxAndWaitMined(ctx context.Context, data []byte) error {
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

func (c *BlockChainClient) MarkRequestCompleted(ctx context.Context, requestID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}

	reqId, err := common.RequestIdStringTo32Byte(requestID)
	if err != nil {
		return fmt.Errorf("invalid request ID %s: %w", requestID, err)
	}

	c.account.Value = nil
	return c.sendTxAndWaitMined(ctx, c.processorEndpoint.PackMarkRequestCompleted(reqId))

}

func (c *BlockChainClient) MarkRequestFailed(ctx context.Context, requestID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}

	reqId, err := common.RequestIdStringTo32Byte(requestID)
	if err != nil {
		return fmt.Errorf("invalid request ID %s: %w", requestID, err)
	}

	c.account.Value = nil
	return c.sendTxAndWaitMined(ctx, c.processorEndpoint.PackMarkRequestFailed(reqId))

}

// SubmitRequest submits a request to the ProcessorEndpoint smart contract using a common.Request.
func (c *BlockChainClient) SubmitRequest(ctx context.Context, protocolVersion uint8, applicationId *big.Int, requestType common.RequestType, payload []byte, value *big.Int) (string, uint64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return "", 0, fmt.Errorf("client not connected, call Connect first")
	}

	reqType, err := requestType.ToUint8()
	if err != nil {
		return "", 0, fmt.Errorf("invalid request type: %w", err)
	}

	// Pack the transaction data using the generated binding
	data := c.processorEndpoint.PackSubmitRequest(protocolVersion, applicationId, reqType, payload, value)
	// Set the value for the transaction (msg.value)
	c.account.Value = value

    // Send the transaction
    tx, err := bind.Transact(c.processorBoundContract, c.account, data)
	c.account.Value = nil
    if err != nil {
        return "", 0, fmt.Errorf("failed to submit transaction: %w", c.UnpackProcessorEndpointError(err))
    }

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, c.client, tx.Hash())
	if err != nil {
		return "", 0, fmt.Errorf("error waiting for tx inclusion: %w", err)
	}
	if receipt.Status != 1 {
		return "", 0, fmt.Errorf("transaction failed")
	}

	// Parse the returned requestId from the transaction receipt logs
	for _, vLog := range receipt.Logs {
		event, err := c.processorEndpoint.UnpackRequestSubmittedEvent(vLog)
		if err == nil {
			return common.RequestId32ByteToString(event.RequestId), receipt.BlockNumber.Uint64(), nil
		}
	}

	return "", 0, fmt.Errorf("requestId not found in logs")
}

func (c *BlockChainClient) SubmitDeanonymizationReport(ctx context.Context, update *common.DeanonymizationReport) error {
	// This is the only thing that has to be done on the blockchain for deanonymization reports
	return c.MarkRequestCompleted(ctx, update.ReportID)
}

func (c *BlockChainClient) SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}

	reqId, err := common.RequestIdStringTo32Byte(update.RequestID)
	if err != nil {
		return fmt.Errorf("invalid request ID %s: %w", update.RequestID, err)
	}

	appId, ok := common.StringToBigInt(update.ApplicationID)
	if !ok {
		return fmt.Errorf("invalid application ID: %s", update.ApplicationID)
	}

	events := make([][]byte, len(update.Events))
	for i, event := range update.Events {
		events[i] = event.EncryptedData
	}

	withdrawals := make([]processorendpoint.StructsWithdrawalRequest, len(update.Withdrawals))
	for i, withdrawal := range update.Withdrawals {
		amount := new(big.Int).SetUint64(withdrawal.Amount)
		withdrawals[i] = processorendpoint.StructsWithdrawalRequest{
			Receiver: ethCommon.HexToAddress(withdrawal.DestinationAddress),
			Amount:   amount,
		}
	}

	params := c.processorEndpoint.PackStateUpdate(
		appId,
		update.PrevStateRoot,
		update.NewStateRoot,
		reqId,
		events,
		withdrawals,
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

// GetUserEvents scans UserEvent logs backwards and returns all the events that are decryptable with the given key
// privKey: user key that will be used to decrypt events
// applicationId: filter events by the given applicationId
// fromBlock: block from which the function search events
// toBlock: block until which the function search events. Note that fromBlock >= toBlock (backwards search)
// f: optional filter function for decrypted events
// stopAtFirst: bool flag to stop at first found event
func (c *BlockChainClient) GetUserEvents(ctx context.Context, privKey cryptotypes.PrivateKeyP521, applicationId big.Int, fromBlock uint64, toBlock uint64, filter func([]byte) bool, stopAtFirst bool) ([][]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	EMPTY := [][]byte{}

	if !c.connected {
		return EMPTY, fmt.Errorf("client not connected, call Connect first")
	}

	contractAddr := c.processorAddress
	//check from block

	fromBlock, err := c.checkQueryFromBlock(ctx, fromBlock, toBlock)
	if err != nil {
		return EMPTY, err
	}

	//retrieve tee public key (needed to decrypt)
	importedPubSecp521r1, err := c.GetTeePublicKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get public key: %w", err)
	}

	//needed for event filter
	userEventSig := c.processorEndpoint.GetEventID(processorendpoint.ProcessorEndpointUserEventEventName)

	appIdHash := ethCommon.BigToHash(&applicationId)
	topicsHash := [][]ethCommon.Hash{{userEventSig}, {appIdHash}}

	var events [][]byte
	query := ethereum.FilterQuery{
		Addresses: []ethCommon.Address{contractAddr},
		FromBlock: new(big.Int).SetUint64(toBlock),
		ToBlock:   new(big.Int).SetUint64(fromBlock),
		Topics:    topicsHash,
	}

	logs, err := c.client.FilterLogs(ctx, query)
	if err != nil {
		return EMPTY, fmt.Errorf("failed to filter logs: %w", err)
	}

	for i := len(logs) - 1; i >= 0; i-- { //backwards search
		vLog := logs[i]
		if !vLog.Removed {
			event, err := c.processorEndpoint.UnpackUserEventEvent(&vLog)
			if err != nil {
				continue
			}
			decrypted, err := crypto.Decrypt(importedPubSecp521r1, &privKey, event.EncryptedData)
			if err == nil && (filter == nil || filter(decrypted)) {
				//found decryptable event that pass filter function
				events = append(events, decrypted)
				if stopAtFirst {
					return events, nil
				}
			}

		}
	}
	return events, nil
}

func (c *BlockChainClient) GetTeePublicKey(ctx context.Context) (*cryptotypes.PublicKeyP521, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.connected {
		return nil, fmt.Errorf("client not connected, call Connect first")
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

func (c *BlockChainClient) checkQueryFromBlock(ctx context.Context, fromBlock uint64, toBlock uint64) (uint64, error) {
	if fromBlock == 0 {
		latestBlock, err := c.client.BlockNumber(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to get latest block number: %w", err)
		}
		fromBlock = latestBlock
	}

	if fromBlock < toBlock {
		return 0, fmt.Errorf("fromBlock should be >= than toBlock")
	}

	return fromBlock, nil
}

// GetRequestCompletedEvent looks for the RequestComleted event for the given request in the given block range and returns if the request was successful or failed
// requestID: identifier of the request
// fromBlock: block from which the function search events
// toBlock: block until which the function search events. Note that fromBlock >= toBlock (backwards search)
func (c *BlockChainClient) GetRequestCompletedEvent(ctx context.Context, requestID string, fromBlock uint64, toBlock uint64) (*common.RequestResult, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("client not connected, call Connect first")
	}

	fromBlock, err := c.checkQueryFromBlock(ctx, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}

	reqId, err := common.RequestIdStringTo32Byte(requestID)
	if err != nil {
		return nil, fmt.Errorf("invalid request ID %s: %w", requestID, err)
	}
	reqIdHash := ethCommon.BytesToHash(reqId[:])

	eventSig := c.processorEndpoint.GetEventID(processorendpoint.ProcessorEndpointRequestCompletedEventName)
	topicsHash := [][]ethCommon.Hash{{eventSig}, {reqIdHash}}

	query := ethereum.FilterQuery{
		Addresses: []ethCommon.Address{c.processorAddress},
		FromBlock: new(big.Int).SetUint64(toBlock),
		ToBlock:   new(big.Int).SetUint64(fromBlock),
		Topics:    topicsHash,
	}

	logs, err := c.client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	valid_logs := make([]types.Log, 0)

	for _, log := range logs {
		if !log.Removed {
			valid_logs = append(valid_logs, log)
		}
	}

	if len(valid_logs) == 0 {
		return nil, nil
	}

	if len(valid_logs) > 1 {
		return nil, fmt.Errorf("found more than 1 log for requestID: %s", requestID)
	}

	event, err := c.processorEndpoint.UnpackRequestCompletedEvent(&valid_logs[0])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack log: %w", err)
	}

	status, err := common.UInt8ToRequestResultStatus(event.Status)
	if err != nil {
		return nil, fmt.Errorf("unknown status: %w", err)
	}

	return &common.RequestResult{Status: status}, nil
}
