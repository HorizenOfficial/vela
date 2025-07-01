package blockchain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/common"
)

// MockClient is a mock implementation of the blockchain client for testing
type MockClient struct {
	mu               sync.RWMutex
	requests         map[string]*common.Request
	pendingRequests  map[string]*common.Request
	states           map[string]*common.ApplicationState
	withdrawals      map[string]*[]common.Withdrawal
	reports          map[string]*common.DeanonymizationReport
	publicKeys       map[string][]byte
	eventSubscribers []chan<- interface{}
}

// NewMockClient creates a new mock blockchain client
func NewMockClient() *MockClient {
	return &MockClient{
		requests:        make(map[string]*common.Request),
		pendingRequests: make(map[string]*common.Request),
		states:          make(map[string]*common.ApplicationState),
		publicKeys:      make(map[string][]byte),
	}
}

// SubmitRequest submits a request to the blockchain
func (c *MockClient) SubmitRequest(ctx context.Context, req *common.Request) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Generate a request ID if not provided
	if req.RequestID == "" {
		id, err := generateRandomID()
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
	c.requests[req.RequestID] = req
	c.pendingRequests[req.RequestID] = req

	return nil
}

// GetPendingRequests gets pending requests from the blockchain
func (c *MockClient) GetPendingRequests(ctx context.Context) ([]*common.Request, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	requests := make([]*common.Request, 0, len(c.pendingRequests))
	for _, req := range c.pendingRequests {
		requests = append(requests, req)
	}

	return requests, nil
}

// MarkRequestCompleted marks a request as completed
func (c *MockClient) MarkRequestCompleted(ctx context.Context, requestID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.pendingRequests[requestID]; !exists {
		return fmt.Errorf("request not found: %s", requestID)
	}

	delete(c.pendingRequests, requestID)

	return nil
}

// SubmitStateUpdate submits a state update to the blockchain
func (c *MockClient) SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ignore validation for simplicity, but in a real implementation, you would validate the signature and state root

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

// generateRandomID generates a random ID
func generateRandomID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
