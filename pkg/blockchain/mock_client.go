package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/crypto"
)

func SetupNewBlockChainClientConnected(client ChainClient, ProcessorContractAddress ethCommon.Address, TeeSignerAddress ethCommon.Address, ManagerAccount *bind.TransactOpts) *BlockChainClient {
	blockchainClient := NewBlockChainClient(ProcessorContractAddress, TeeSignerAddress, "", nil)
	blockchainClient.client = client

	blockchainClient.processorBoundContract = blockchainClient.processorEndpoint.Instance(blockchainClient.client, ProcessorContractAddress)
	blockchainClient.teeAuthBoundContract = blockchainClient.teeAuthEndpoint.Instance(blockchainClient.client, TeeSignerAddress)

	blockchainClient.account = ManagerAccount
	blockchainClient.connected = true

	return blockchainClient
}

// MockClient is a mock implementation of the blockchain client for testing
type MockClient struct {
	mu               sync.RWMutex
	requests         *orderedmap.OrderedMap[common.RequestIdType, *common.Request]
	pendingRequests  *orderedmap.OrderedMap[common.RequestIdType, *common.Request]
	failedRequests   *orderedmap.OrderedMap[common.RequestIdType, *common.Request]
	states           map[common.ApplicationIdType]*common.ApplicationState
	withdrawals      map[common.ApplicationIdType]*[]common.Withdrawal
	reports          map[common.RequestIdType]*common.DeanonymizationReport
	updatePayloads   map[common.RequestIdType]*common.UpdatePayload
	eventSubscribers []chan<- interface{}
	stateRoot        [32]byte
	*testutil.MockFunctions
}

// NewMockClient creates a new mock blockchain client
func NewMockClient() *MockClient {
	return &MockClient{
		requests:        orderedmap.NewOrderedMap[common.RequestIdType, *common.Request](),
		pendingRequests: orderedmap.NewOrderedMap[common.RequestIdType, *common.Request](),
		failedRequests:  orderedmap.NewOrderedMap[common.RequestIdType, *common.Request](),
		states:          make(map[common.ApplicationIdType]*common.ApplicationState),
		withdrawals:     make(map[common.ApplicationIdType]*[]common.Withdrawal),
		reports:         make(map[common.RequestIdType]*common.DeanonymizationReport),
		updatePayloads:  make(map[common.RequestIdType]*common.UpdatePayload),
		MockFunctions:   testutil.NewMockFunctions(),
	}
}

func (c *MockClient) SendRequestToChain(ctx context.Context, req *common.Request) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Generate a request ID if not provided
	emptyRequestId := common.RequestIdType{}
	if req.RequestID == emptyRequestId {
		id := testutil.GenerateRandomRequestID()
		req.RequestID = id
	}

	// Set timestamp if not provided
	if req.Timestamp == nil {
		req.Timestamp = new(big.Int).SetInt64(time.Now().Unix())
	}

	// Store the request
	c.requests.Set(req.RequestID, req)
	c.pendingRequests.Set(req.RequestID, req)

	return nil
}

// SubmitRequest submits a request to the blockchain according to the official interface
func (c *MockClient) SubmitRequest(ctx context.Context, protocolVersion uint8, applicationId common.ApplicationIdType, requestType common.RequestType, payload []byte, value *big.Int, maxFeeValue *big.Int) (common.RequestIdType, uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//prepare request
	req := &common.Request{
		ProtocolVersion: protocolVersion,
		ApplicationID:   applicationId,
		RequestType:     requestType,
		Payload:         payload,
		Value:           value,
		MaxFeeValue: maxFeeValue,
	}

	err := c.SendRequestToChain(ctx, req)
	if err != nil {
		return common.RequestIdType{}, 0, fmt.Errorf("failed to send request: %w", err)
	}

	return req.RequestID, 0, nil
}

// GetPendingRequests gets pending requests from the blockchain
func (c *MockClient) GetPendingRequests(ctx context.Context) ([]*common.Request, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	requests := make([]*common.Request, 0, c.pendingRequests.Len())
	for req := range c.pendingRequests.Values() {
		requests = append(requests, req)
	}

	return requests, nil
}

// GetNextPendingRequest gets the next pending request and the current stateRoot from the blockchain
func (c *MockClient) GetNextPendingRequest(ctx context.Context) (*common.Request, [32]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if f, ok := c.GetMockedFunc("GetNextPendingRequest"); ok {
		return f.(func(context.Context) (*common.Request, [32]byte, error))(ctx)
	}
	var req *common.Request
	if c.pendingRequests.Len() > 0 {
		req = c.pendingRequests.Front().Value
	}

	return req, c.stateRoot, nil

}

// MarkRequestFailed marks a request as failed
func (c *MockClient) MarkRequestFailed(ctx context.Context, requestID common.RequestIdType, requestFailure *apperrors.RequestFailure) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if f, ok:= c.GetMockedFunc("MarkRequestFailed"); ok {
		return f.(func(context.Context, common.RequestIdType, *apperrors.RequestFailure) error)(ctx, requestID, requestFailure)
	}

	if !c.pendingRequests.Has(requestID) {
		return fmt.Errorf("request not found: %s", requestID)
	}

	c.pendingRequests.Delete(requestID)
	req, _ := c.requests.Get(requestID)
	c.failedRequests.Set(requestID, req)

	return nil
}

func (c *MockClient) GetCompletedRequests() []*common.Request {
	c.mu.RLock()
	defer c.mu.RUnlock()

	completed := make([]*common.Request, 0)
	for id, req := range c.requests.AllFromFront() {
		if !c.pendingRequests.Has(id) {
			completed = append(completed, req)
		}
	}
	return completed
}

func (c *MockClient) GetFailedRequests() []*common.Request {
	c.mu.RLock()
	defer c.mu.RUnlock()

	failed := make([]*common.Request, 0)
	for req := range c.failedRequests.Values() {
		failed = append(failed, req)
	}
	return failed
}

func (c *MockClient) WaitForRequestCompletion(requestID common.RequestIdType, timeout time.Duration) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			c.mu.RLock()
			exists := c.pendingRequests.Has(requestID)
			failed := c.failedRequests.Has(requestID)
			c.mu.RUnlock()
			if failed {
				return fmt.Errorf("request %s has failed", requestID)
			}
			if !exists {
				return nil // Request completed
			}
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for request %s to complete", requestID)
		}
	}
}

// SubmitStateUpdate submits a state update to the blockchain
func (c *MockClient) SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if f, ok := c.GetMockedFunc("SubmitStateUpdate"); ok {
		return f.(func(context.Context, *common.UpdatePayload) error)(ctx, update)
	}
	// Complete the request if it exists
	if !c.pendingRequests.Has(update.RequestID) {
		return fmt.Errorf("request not found: %s", update.RequestID)
	}
	c.pendingRequests.Delete(update.RequestID)

	// Store update payload for separate verification by test suite
	c.updatePayloads[update.RequestID] = update

	// Update state
	c.states[update.ApplicationID] = &common.ApplicationState{
		ApplicationID:  update.ApplicationID,
		StateRoot:      update.NewStateRoot,
		EncryptedState: nil, // State is stored separately in the data layer
	}

	c.stateRoot = update.NewStateRoot

	// Emit events
	c.emitEvents(update.Events)

	// Store withdrawals, in real implementation this would be events or some other mechanism
	if len(update.Withdrawals) > 0 {
		withdrawals, exists := c.withdrawals[update.ApplicationID]
		if !exists {
			withdrawals = &[]common.Withdrawal{}
			c.withdrawals[update.ApplicationID] = withdrawals
		}
		*withdrawals = append(*withdrawals, update.Withdrawals...)
	}

	return nil
}

// GetRequestUpdatePayload gets the update payload for a request
func (c *MockClient) GetRequestUpdatePayload(ctx context.Context, requestID common.RequestIdType) (*common.UpdatePayload, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	update, exists := c.updatePayloads[requestID]
	if !exists {
		return nil, fmt.Errorf("update payload not found for request: %s", requestID)
	}
	return update, nil
}

// GetApplicationState gets the state of an application
func (c *MockClient) GetApplicationState(ctx context.Context, applicationID common.ApplicationIdType) (*common.ApplicationState, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state, exists := c.states[applicationID]
	if !exists {
		return nil, fmt.Errorf("application state not found: %d", applicationID)
	}

	return state, nil
}

// SubscribeToEvents subscribes to events from the blockchain
func (c *MockClient) SubscribeToEvents(ctx context.Context, eventCh chan<- interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.eventSubscribers = append(c.eventSubscribers, eventCh)

	// Start a goroutine to close the channel when the context is done
	// not sure if this is needed before real testing, learning purposes
	go func() {
		<-ctx.Done()
		c.mu.Lock()
		defer c.mu.Unlock()

		for i, ch := range c.eventSubscribers {
			if ch == eventCh {
				c.eventSubscribers = append(c.eventSubscribers[:i], c.eventSubscribers[i+1:]...)
				break
			}
		}
	}()

	return nil
}

// GetWithdrawals gets withdrawal requests for an application
func (c *MockClient) GetWithdrawals(ctx context.Context, applicationID common.ApplicationIdType) (*[]common.Withdrawal, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	withdrawals, exists := c.withdrawals[applicationID]
	if !exists {
		return nil, fmt.Errorf("withdrawals not found for application: %d", applicationID)
	}

	return withdrawals, nil
}

// SubmitDeanonymizationReport submits a deanonymization report to the blockchain
func (c *MockClient) SubmitDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if f, ok := c.GetMockedFunc("SubmitDeanonymizationReport"); ok {
		return f.(func(context.Context, *common.DeanonymizationReport) error)(ctx, report)
	}

	// Complete the request if it exists
	if !c.pendingRequests.Has(report.ReportID) {
		return fmt.Errorf("request not found: %s", report.ReportID)
	}
	c.pendingRequests.Delete(report.ReportID)

	// store the report
	c.reports[report.ReportID] = report

	return nil
}

// GetDeanonymizationReport gets a deanonymization report
func (c *MockClient) GetDeanonymizationReport(ctx context.Context, reportID common.RequestIdType) (*common.DeanonymizationReport, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	report, exists := c.reports[reportID]
	if !exists {
		return nil, fmt.Errorf("deanonymization report not found: %s", reportID)
	}

	return report, nil
}

func (c *MockClient) GetUserEvents(ctx context.Context, privKey cryptotypes.PrivateKeyP521, applicationId common.ApplicationIdType, fromBlock uint64, toBlock uint64, filter func([]byte) bool, stopAtFirst bool) ([][]byte, error) {
	return [][]byte{}, nil
}

func (c *MockClient) GetRequestCompletedEvent(ctx context.Context, requestID common.RequestIdType, fromBlock uint64, toBlock uint64) (*common.RequestResult, error) {
	return nil, nil
}

func (c *MockClient) GetTeePublicKey(ctx context.Context) (*cryptotypes.PublicKeyP521, error) {
	key, err := crypto.GeneratePrivateKeyP521()
	return key.PublicKey(), err
}

// Close closes the blockchain client
func (c *MockClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if f, ok := c.GetMockedFunc("Close"); ok {
		return f.(func() error)()
	}
	// Close all event subscribers
	c.eventSubscribers = nil

	return nil
}

// Close closes the blockchain client
func (c *MockClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if f, ok := c.GetMockedFunc("Connect"); ok {
		return f.(func(context.Context) error)(ctx)
	}

	return nil
}

func (c *MockClient) ClearAllData() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.requests = orderedmap.NewOrderedMap[common.RequestIdType, *common.Request]()
	c.pendingRequests = orderedmap.NewOrderedMap[common.RequestIdType, *common.Request]()
	c.states = make(map[common.ApplicationIdType]*common.ApplicationState)
	c.withdrawals = make(map[common.ApplicationIdType]*[]common.Withdrawal)
	c.reports = make(map[common.RequestIdType]*common.DeanonymizationReport)
	c.failedRequests = orderedmap.NewOrderedMap[common.RequestIdType, *common.Request]()
	c.updatePayloads = make(map[common.RequestIdType]*common.UpdatePayload)
	c.stateRoot = [32]byte{}
	c.MockedFunctions = make(map[string]interface{})
}

// emitEvent emits an event to all subscribers
func (c *MockClient) emitEvents(events []common.Event) {
	for _, event := range events {
		for _, ch := range c.eventSubscribers {
			select {
			case ch <- event:
			default:
				// Skip if the channel is full
			}
		}
	}
}
