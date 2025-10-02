package blockchain

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/horizen-pes/pkg/blockchain/contracts/keyregistry"
	"github.com/horizen-pes/pkg/blockchain/contracts/processorendpoint"
	"github.com/horizen-pes/pkg/blockchain/contracts/tee"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/crypto"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
)

//go:generate mkdir -p ../../contract_abis
//go:generate mkdir -p ./contracts/processorendpoint
//go:generate solc --combined-json abi,bin ../../contracts/contracts/ProcessorEndpoint.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/ProcessorEndpointAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/ProcessorEndpointAbi/combined.json --pkg processorendpoint --type ProcessorEndpoint --out ./contracts/processorendpoint/ProcessorEndpoint.go
//go:generate mkdir -p ./contracts/keyregistry
//go:generate solc --combined-json abi,bin ../../contracts/contracts/KeyRegistry.sol --base-path ../.. --include-path ../../contracts/node_modules --pretty-json -o ../../contract_abis/KeyRegistryAbi --overwrite
//go:generate abigen --v2 --combined-json ../../contract_abis/KeyRegistryAbi/combined.json --pkg keyregistry --type KeyRegistry --out ./contracts/keyregistry/KeyRegistry.go

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
	mu                       sync.RWMutex
	connected                bool
	processorAddress         ethCommon.Address
	keyRegistryAddress       ethCommon.Address
	teeAuthAddress			 ethCommon.Address
	rpcURL                   string
	processorBoundContract   *bind.BoundContract
	processorEndpoint        *processorendpoint.ProcessorEndpoint
	keyRegistryBoundContract *bind.BoundContract
	keyRegistryEndpoint      *keyregistry.KeyRegistry
	teeAuthBoundContract     *bind.BoundContract
	teeAuthEndpoint          *tee.TeeAuthenticator
	client                   ChainClient
	privKey                  *cryptotypes.PrivateKeySecp256k1
	account                  *bind.TransactOpts
}

func toRequestType(i uint8) common.RequestType {
	switch i {
	case 0:
		return common.Deploy
	case 1:
		return common.Process
	case 2:
		return common.Deanonymize
	default:
		return ""
	}
}

func NewBlockChainClient(processor ethCommon.Address, keyRegistry ethCommon.Address, teeAuthenticator ethCommon.Address, rpcURL string, key *cryptotypes.PrivateKeySecp256k1) *BlockChainClient {
	return &BlockChainClient{
		processorAddress:    processor,
		keyRegistryAddress:  keyRegistry,
		teeAuthAddress:		 teeAuthenticator,
		rpcURL:              rpcURL,
		processorEndpoint:   processorendpoint.NewProcessorEndpoint(),
		keyRegistryEndpoint: keyregistry.NewKeyRegistry(),
		teeAuthEndpoint:	 tee.NewTeeAuthenticator(),
		privKey:             key,
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
	c.keyRegistryBoundContract = c.keyRegistryEndpoint.Instance(c.client, c.keyRegistryAddress)

	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve chain ID: %w", err)
	}

	c.account = bind.NewKeyedTransactor(c.privKey.PrivateKey, chainID)

	c.connected = true
	return nil
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
		return nil, fmt.Errorf("call returned error: %w", err)
	}

	output := make([]*common.Request, 0, len(listOfRequests))
	for _, request := range listOfRequests {
		//TODO check that all big.Int can fit in a Uint64. If not, the specific request should be marked as failed
		req := &common.Request{
			ProtocolVersion: strconv.FormatUint(uint64(request.ProtocolVersion), 10),
			ApplicationID:   request.ApplicationId.String(),
			RequestID:       request.RequestId.String(),
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

func (c *BlockChainClient) MarkRequestCompleted(ctx context.Context, requestID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}

	reqId, ok := common.StringToBigInt(requestID)
	if !ok {
		return fmt.Errorf("invalid request ID: %s", requestID)
	}

	tx, err := bind.Transact(c.processorBoundContract, c.account, c.processorEndpoint.PackMarkRequestCompleted(reqId))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// wait for transaction inclusion
	if _, err := bind.WaitMined(ctx, c.client, tx.Hash()); err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %w", err)
	}
	return nil

}

func (c *BlockChainClient) MarkRequestFailed(ctx context.Context, requestID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}

	reqId, ok := common.StringToBigInt(requestID)
	if !ok {
		return fmt.Errorf("invalid request ID: %s", requestID)
	}

	tx, err := bind.Transact(c.processorBoundContract, c.account, c.processorEndpoint.PackMarkRequestFailed(reqId))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// wait for transaction inclusion
	if _, err := bind.WaitMined(ctx, c.client, tx.Hash()); err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %w", err)
	}
	return nil

}

func (c *BlockChainClient) SubmitDeanonymizationReport(ctx context.Context, update *common.DeanonymizationReport) error {
	// This is the only thing that has to be done on the blockchain for deanonymization reports
	return c.MarkRequestCompleted(ctx, update.ReportID)
}

func (c *BlockChainClient) SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.connected == false {
		return fmt.Errorf("client not connected, call Connect first")
	}

	reqId, ok := common.StringToBigInt(update.RequestID)
	if !ok {
		return fmt.Errorf("invalid request ID: %s", update.RequestID)
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

	tx, err := bind.Transact(c.processorBoundContract, c.account, params)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// wait for transaction inclusion
	if _, err := bind.WaitMined(ctx, c.client, tx.Hash()); err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %w", err)
	}
	return nil

}

func (c *BlockChainClient) GetPublicKey(ctx context.Context, address string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return nil, fmt.Errorf("client not connected, call Connect first")
	}

	pubKey, err := bind.Call(c.keyRegistryBoundContract,
		&bind.CallOpts{Pending: false},
		c.keyRegistryEndpoint.PackGetPK(ethCommon.HexToAddress(address)),
		c.keyRegistryEndpoint.UnpackGetPK)

	if err != nil {
		return nil, fmt.Errorf("call returned error: %w", err)
	}

	return pubKey, nil
}

func (c *BlockChainClient) RegisterPK(ctx context.Context, publicKey []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return fmt.Errorf("client not connected, call Connect first")
	}

	tx, err := bind.Transact(c.processorBoundContract, c.account, c.keyRegistryEndpoint.PackRegisterPK(publicKey))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %w", err)
	}

	// wait for transaction inclusion
	if _, err := bind.WaitMined(ctx, c.client, tx.Hash()); err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %w", err)
	}
	return nil

}

// Close closes the blockchain client
func (c *BlockChainClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.client = nil
	c.processorBoundContract = nil
	c.keyRegistryBoundContract = nil
	c.account = nil

	c.connected = false
	return nil
}

// GetPrivateBalance scans UserEvent logs backwards and returns all the events that are decryptable in the last block in which a decryptable event is present. 
// If f!=nil, events should be decryptable and f should return true
func (c *BlockChainClient) GetUserEvents(ctx context.Context, privKey cryptotypes.PrivateKeyP521, applicationId big.Int, fromBlock uint64, toBlock uint64, f func([]byte) bool) ([][]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	EMPTY := [][]byte{}

	if !c.connected {
		return EMPTY, fmt.Errorf("client not connected, call Connect first")
	}

	contractAddr := c.processorAddress
	//check from block
	if fromBlock == 0 {
		latestBlock, err := c.client.BlockByNumber(ctx, nil)
		if err != nil {
			return EMPTY, fmt.Errorf("failed to get latest block: %w", err)
		}
		fromBlock = latestBlock.NumberU64()
	}

	parsedABI, err := processorendpoint.ProcessorEndpointMetaData.ParseABI()
	if err != nil {
		return EMPTY, fmt.Errorf("cannot parse ABI: %w", err)
	}

	//retrieve tee public key (needed to decrypt)
	pubSecp521r1, err := bind.Call(c.teeAuthBoundContract,
		&bind.CallOpts{Pending: false},
		c.teeAuthEndpoint.PackPubSecp521r1(),
		c.teeAuthEndpoint.UnpackPubSecp521r1)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve pubSecp521r1: %w", err)
	}
	importedPubSecp521r1, err := crypto.ImportPublicKeyP521FromHex(hex.EncodeToString(pubSecp521r1))
	if err != nil {
		return nil, fmt.Errorf("cannot import pubSecp521r1: %w", err)
	}

	//needed for event filter
	userEventSig := parsedABI.Events["UserEvent"].ID
	appIdHash := ethCommon.BigToHash(&applicationId)
	topicsHash := [][]ethCommon.Hash{{userEventSig}, {appIdHash}}

	for blockNumber := fromBlock; blockNumber >= toBlock; blockNumber-- {
		query := ethereum.FilterQuery{
			Addresses: []ethCommon.Address{contractAddr},
			FromBlock: new(big.Int).SetUint64(blockNumber),
			ToBlock:   new(big.Int).SetUint64(blockNumber),
			Topics:    topicsHash,
		} //in this way we avoid problems for too bigs interval and we avoid to invert the sort by block

		logs, err := c.client.FilterLogs(ctx, query)
		if err != nil {
			return EMPTY, fmt.Errorf("failed to filter logs: %w", err)
		}
		
		var events [][]byte
		for _, vLog := range logs {
			event, err := c.processorEndpoint.UnpackUserEventEvent(&vLog)
			if err != nil {
				continue
			}
			decrypted, err := crypto.Decrypt(importedPubSecp521r1, &privKey, event.EncryptedData)
			if err != nil && (f == nil || f(decrypted)) {
				//found decryptable event that pass filter function
				events = append(events, decrypted)
			}
		}
		if len(events) > 0 { //at least one event is found in this block
			return events, nil 
		}
	}
	return EMPTY, errors.New("no event found")
}
