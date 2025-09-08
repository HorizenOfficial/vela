package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	commonEth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/horizen-pes/pkg/blockchain/contracts/keyregistry"
	"github.com/horizen-pes/pkg/blockchain/contracts/processorendpoint"
	"github.com/horizen-pes/pkg/common"
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
	processorAddress         commonEth.Address
	keyRegistryAddress       commonEth.Address
	rpcURL                   string
	processorBoundContract   *bind.BoundContract
	processorEndpoint        *processorendpoint.ProcessorEndpoint
	keyRegistryBoundContract *bind.BoundContract
	keyRegistryEndpoint      *keyregistry.KeyRegistry
	client                   ChainClient
	privKey                  *common.PrivateKeySecp256k1
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

func stringToBigInt(s string) (*big.Int, bool) {
	i, ok := new(big.Int).SetString(s, 10)
	return i, ok
}

func NewBlockChainClient(processor commonEth.Address, keyRegistry commonEth.Address, rpcURL string, key *common.PrivateKeySecp256k1) *BlockChainClient {
	return &BlockChainClient{
		processorAddress:    processor,
		keyRegistryAddress:  keyRegistry,
		rpcURL:              rpcURL,
		processorEndpoint:   processorendpoint.NewProcessorEndpoint(),
		keyRegistryEndpoint: keyregistry.NewKeyRegistry(),
		privKey:             key,
	}
}

func (c *BlockChainClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("Already connected")
	}

	var err error
	c.client, err = ethclient.Dial(c.rpcURL)
	if err != nil {
		return fmt.Errorf("Cannot connect to chain: %v", err)
	}

	c.processorBoundContract = c.processorEndpoint.Instance(c.client, c.processorAddress)
	c.keyRegistryBoundContract = c.keyRegistryEndpoint.Instance(c.client, c.keyRegistryAddress)

	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("Failed to retrieve chain ID: %v", err)
	}

	c.account = bind.NewKeyedTransactor(c.privKey.PrivateKey, chainID)

	c.connected = true
	return nil
}

// GetPendingRequests gets pending requests from the blockchain
func (c *BlockChainClient) GetPendingRequests(ctx context.Context) ([]*common.Request, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.connected == false {
		return nil, fmt.Errorf("client not connected, call Connect first")
	}

	listOfRequests, err := bind.Call(c.processorBoundContract,
		&bind.CallOpts{Pending: false},
		c.processorEndpoint.PackGetPendingRequests(),
		c.processorEndpoint.UnpackGetPendingRequests)

	if err != nil {
		return nil, fmt.Errorf("call returned error: %v", err)
	}

	output := make([]*common.Request, 0, len(listOfRequests))
	for _, request := range listOfRequests {
		//TODO check that all big.Int can fit in a int64. If not, the specific request should be marked as failed
		req := &common.Request{
			ProtocolVersion: strconv.FormatUint(uint64(request.ProtocolVersion), 10),
			ApplicationID:   request.ApplicationId.String(),
			RequestID:       request.RequestId.String(),
			RequestType:     toRequestType(request.RequestType),
			Payload:         request.Payload,
			Timestamp:       request.Timestamp.Int64(),
			Sender:          request.Sender.String(),
			Value: 		     request.Value.Int64(),
		}

		output = append(output, req)
	}
	return output, nil
}

func (c *BlockChainClient) MarkRequestCompleted(ctx context.Context, requestID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.connected == false {
		return fmt.Errorf("client not connected, call Connect first")
	}

	reqId, ok := stringToBigInt(requestID)
	if !ok {
		return fmt.Errorf("invalid request ID: %s", requestID)
	}

	tx, err := bind.Transact(c.processorBoundContract, c.account, c.processorEndpoint.PackMarkRequestCompleted(reqId))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	// wait for transaction inclusion
	if _, err := bind.WaitMined(ctx, c.client, tx.Hash()); err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %v", err)
	}
	return nil

}

func (c *BlockChainClient) MarkRequestFailed(ctx context.Context, requestID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.connected == false {
		return fmt.Errorf("client not connected, call Connect first")
	}

	reqId, ok := stringToBigInt(requestID)
	if !ok {
		return fmt.Errorf("invalid request ID: %s", requestID)
	}

	tx, err := bind.Transact(c.processorBoundContract, c.account, c.processorEndpoint.PackMarkRequestFailed(reqId))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	// wait for transaction inclusion
	if _, err := bind.WaitMined(ctx, c.client, tx.Hash()); err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %v", err)
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

	reqId, ok := stringToBigInt(update.RequestID)
	if !ok {
		return fmt.Errorf("invalid request ID: %s", update.RequestID)
	}

	appId, ok := stringToBigInt(update.ApplicationID)
	if !ok {
		return fmt.Errorf("invalid application ID: %s", update.ApplicationID)
	}

	events := make([][]byte, len(update.Events))
	for i, event := range update.Events {
		events[i] = event.EncryptedData
	}

	withdrawals := make([]processorendpoint.StructsWithdrawalRequest, len(update.Withdrawals))
	for i, withdrawal := range update.Withdrawals {
		amount, ok := stringToBigInt(withdrawal.Amount)
		if !ok {
			return fmt.Errorf("invalid amount: %s", withdrawal.Amount)
		}
		withdrawals[i] = processorendpoint.StructsWithdrawalRequest{
			Receiver: commonEth.HexToAddress(withdrawal.DestinationAddress),
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
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	// wait for transaction inclusion
	if _, err := bind.WaitMined(ctx, c.client, tx.Hash()); err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %v", err)
	}
	return nil

}

func (c *BlockChainClient) GetPublicKey(ctx context.Context, address string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.connected == false {
		return nil, fmt.Errorf("client not connected, call Connect first")
	}

	pubKey, err := bind.Call(c.keyRegistryBoundContract,
		&bind.CallOpts{Pending: false},
		c.keyRegistryEndpoint.PackGetPK(commonEth.HexToAddress(address)),
		c.keyRegistryEndpoint.UnpackGetPK)

	if err != nil {
		return nil, fmt.Errorf("call returned error: %v", err)
	}

	return pubKey, nil
}

func (c *BlockChainClient) RegisterPK(ctx context.Context, publicKey []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.connected == false {
		return fmt.Errorf("client not connected, call Connect first")
	}

	tx, err := bind.Transact(c.processorBoundContract, c.account, c.keyRegistryEndpoint.PackRegisterPK(publicKey))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	// wait for transaction inclusion
	if _, err := bind.WaitMined(ctx, c.client, tx.Hash()); err != nil {
		return fmt.Errorf("error waiting for tx inclusion: %v", err)
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
