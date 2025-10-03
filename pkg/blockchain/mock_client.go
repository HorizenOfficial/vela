package blockchain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/elliotchance/orderedmap/v3"
)

// MockClient is a mock implementation of the blockchain client for testing
type MockClient struct {
	mu               sync.RWMutex
	requests         *orderedmap.OrderedMap[string, *common.Request]
	pendingRequests  *orderedmap.OrderedMap[string, *common.Request]
	failedRequests   *orderedmap.OrderedMap[string, *common.Request]
	states           map[string]*common.ApplicationState
	withdrawals      map[string]*[]common.Withdrawal
	reports          map[string]*common.DeanonymizationReport
	updatePayloads   map[string]*common.UpdatePayload
	publicKeys       map[string][]byte
	eventSubscribers []chan<- interface{}
}

// NewMockClient creates a new mock blockchain client
func NewMockClient() *MockClient {
	return &MockClient{
		requests:        orderedmap.NewOrderedMap[string, *common.Request](),
		pendingRequests: orderedmap.NewOrderedMap[string, *common.Request](),
		failedRequests:  orderedmap.NewOrderedMap[string, *common.Request](),
		states:          make(map[string]*common.ApplicationState),
		withdrawals:     make(map[string]*[]common.Withdrawal),
		reports:         make(map[string]*common.DeanonymizationReport),
		updatePayloads:  make(map[string]*common.UpdatePayload),
		publicKeys:      make(map[string][]byte),
	}
}

// SubmitRequest submits a request to the blockchain
func (c *MockClient) SubmitRequest(ctx context.Context, req *common.Request) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Generate a request ID if not provided
	if req.RequestID == "" {
		id, err := GenerateRandomID()
		if err != nil {
			return fmt.Errorf("failed to generate request ID: %w", err)
		}
		req.RequestID = id
	}

	// Set timestamp if not provided
	if req.Timestamp == 0 {
		req.Timestamp = time.Now().Unix()
	}

	// Store the request
	c.requests.Set(req.RequestID, req)
	c.pendingRequests.Set(req.RequestID, req)
	

	return nil
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

	var req *common.Request
	if c.pendingRequests.Len() > 0 {
		req = c.pendingRequests.Front().Value
	}

	var stateRoot [32]byte

	appState, err := c.GetApplicationState(ctx, req.ApplicationID)
	if err != nil {
		return req, [32]byte{}, nil
	}

	if appState != nil {
		stateRoot = appState.StateRoot
	}
	return req, stateRoot, nil

}


// MarkRequestFailed marks a request as failed
func (c *MockClient) MarkRequestFailed(ctx context.Context, requestID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

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

func (c *MockClient) WaitForRequestCompletion(requestID string, timeout time.Duration) error {
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
func (c *MockClient) GetRequestUpdatePayload(ctx context.Context, requestID string) (*common.UpdatePayload, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	update, exists := c.updatePayloads[requestID]
	if !exists {
		return nil, fmt.Errorf("update payload not found for request: %s", requestID)
	}
	return update, nil
}

// GetApplicationState gets the state of an application
func (c *MockClient) GetApplicationState(ctx context.Context, applicationID string) (*common.ApplicationState, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state, exists := c.states[applicationID]
	if !exists {
		return nil, fmt.Errorf("application state not found: %s", applicationID)
	}

	return state, nil
}

// RegisterPublicKey registers a public key for an address
func (c *MockClient) RegisterPublicKey(ctx context.Context, address string, publicKey []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.publicKeys[address] = publicKey

	return nil
}

// GetPublicKey gets the public key for an address
func (c *MockClient) GetPublicKey(ctx context.Context, address string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	publicKey, exists := c.publicKeys[address]
	if !exists {
		return nil, fmt.Errorf("public key not found for address: %s", address)
	}

	return publicKey, nil
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
func (c *MockClient) GetWithdrawals(ctx context.Context, applicationID string) (*[]common.Withdrawal, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	withdrawals, exists := c.withdrawals[applicationID]
	if !exists {
		return nil, fmt.Errorf("withdrawals not found for application: %s", applicationID)
	}

	return withdrawals, nil
}

// SubmitDeanonymizationReport submits a deanonymization report to the blockchain
func (c *MockClient) SubmitDeanonymizationReport(ctx context.Context, report *common.DeanonymizationReport) error {
	c.mu.Lock()
	defer c.mu.Unlock()

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
func (c *MockClient) GetDeanonymizationReport(ctx context.Context, reportID string) (*common.DeanonymizationReport, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	report, exists := c.reports[reportID]
	if !exists {
		return nil, fmt.Errorf("deanonymization report not found: %s", reportID)
	}

	return report, nil
}

// Close closes the blockchain client
func (c *MockClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close all event subscribers
	c.eventSubscribers = nil

	return nil
}

// Close closes the blockchain client
func (c *MockClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}

func (c *MockClient) ClearAllData() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.requests = orderedmap.NewOrderedMap[string, *common.Request]()
	c.pendingRequests = orderedmap.NewOrderedMap[string, *common.Request]()
	c.states = make(map[string]*common.ApplicationState)
	c.publicKeys = make(map[string][]byte)
	c.withdrawals = make(map[string]*[]common.Withdrawal)
	c.reports = make(map[string]*common.DeanonymizationReport)
	c.failedRequests =orderedmap.NewOrderedMap[string, *common.Request]()
	c.updatePayloads = make(map[string]*common.UpdatePayload)
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

// GenerateRandomID generates a random ID
func GenerateRandomID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
