package fullstack

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/HorizenOfficial/vela/pkg/blockchain"
	"github.com/HorizenOfficial/vela/pkg/common"
)

// eventBroadcastingClient wraps a real blockchain.Client (typically a BlockChainClient
// connected to the simulated backend) and intercepts SubmitStateUpdate to:
//   - broadcast events to the suite's eventChannel
//   - track request completion (for AssertRequestCompleted)
//   - store update payloads (for GetRequestUpdatePayload)
//   - accumulate withdrawals (for WaitForWithdrawal)
//   - notify the InProcessSubgraph via the onStateUpdate callback
//
// All other Client methods are delegated to the underlying client unchanged.
type eventBroadcastingClient struct {
	blockchain.Client // embedded real client — all methods delegated by default

	mu             sync.Mutex
	eventChannel   chan<- interface{}
	pendingIDs     map[common.RequestIdType]struct{}
	deployIDs      map[common.RequestIdType]struct{} // tracks which requests are deploy requests
	completedIDs   map[common.RequestIdType]struct{}
	failedIDs      map[common.RequestIdType]struct{}
	updatePayloads map[common.RequestIdType]*common.UpdatePayload
	withdrawals    map[common.ApplicationIdType][]common.Withdrawal

	// onStateUpdate is called after each successful SubmitStateUpdate.
	// Used by InProcessSubgraph to record data for wallet queries.
	onStateUpdate func(update *common.UpdatePayload, isDeploy bool)
}

func newEventBroadcastingClient(inner blockchain.Client, eventCh chan<- interface{}) *eventBroadcastingClient {
	return &eventBroadcastingClient{
		Client:         inner,
		eventChannel:   eventCh,
		pendingIDs:     make(map[common.RequestIdType]struct{}),
		deployIDs:      make(map[common.RequestIdType]struct{}),
		completedIDs:   make(map[common.RequestIdType]struct{}),
		failedIDs:      make(map[common.RequestIdType]struct{}),
		updatePayloads: make(map[common.RequestIdType]*common.UpdatePayload),
		withdrawals:    make(map[common.ApplicationIdType][]common.Withdrawal),
	}
}

// markPending registers a request as pending so AssertRequestCompleted can track it.
// isDeploy indicates whether this is a deploy request (needed by the subgraph to
// route to the correct completion map).
func (c *eventBroadcastingClient) markPending(requestID common.RequestIdType, isDeploy bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pendingIDs[requestID] = struct{}{}
	if isDeploy {
		c.deployIDs[requestID] = struct{}{}
	}
}

// SubmitStateUpdate delegates to the real client and then records the result
// for test observation (events, completions, withdrawals).
func (c *eventBroadcastingClient) SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error {
	// Delegate to the real blockchain client (on-chain transaction)
	err := c.Client.SubmitStateUpdate(ctx, update)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove from pending
	delete(c.pendingIDs, update.RequestID)

	// Check if this was a deploy request
	_, isDeploy := c.deployIDs[update.RequestID]

	if update.ErrorCode != 0 {
		// Mark as failed
		c.failedIDs[update.RequestID] = struct{}{}
		// Still notify the subgraph so it records the failure
		if c.onStateUpdate != nil {
			c.onStateUpdate(update, isDeploy)
		}
		return nil
	}

	// Mark as completed
	c.completedIDs[update.RequestID] = struct{}{}

	// Store the full update payload for GetRequestUpdatePayload
	c.updatePayloads[update.RequestID] = update

	// Broadcast events to the suite's event channel
	for _, event := range update.Events {
		select {
		case c.eventChannel <- event:
		default:
			// Channel full — drop to avoid blocking the manager
		}
	}

	// Accumulate withdrawals
	if len(update.Withdrawals) > 0 {
		c.withdrawals[update.ApplicationID] = append(
			c.withdrawals[update.ApplicationID],
			update.Withdrawals...,
		)
	}

	// Notify the in-process subgraph
	if c.onStateUpdate != nil {
		c.onStateUpdate(update, isDeploy)
	}

	return nil
}

// waitForRequestCompletion polls until the request is completed or failed, or times out.
func (c *eventBroadcastingClient) waitForRequestCompletion(requestID common.RequestIdType, timeout time.Duration) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			_, completed := c.completedIDs[requestID]
			_, failed := c.failedIDs[requestID]
			c.mu.Unlock()
			if failed {
				return fmt.Errorf("request %s has failed", requestID)
			}
			if completed {
				return nil
			}
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for request %s to complete", requestID)
		}
	}
}

// getUpdatePayload returns the stored UpdatePayload for a completed request.
func (c *eventBroadcastingClient) getUpdatePayload(requestID common.RequestIdType) (*common.UpdatePayload, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	update, exists := c.updatePayloads[requestID]
	if !exists {
		return nil, fmt.Errorf("update payload not found for request: %s", requestID)
	}
	return update, nil
}

// getWithdrawals returns the accumulated withdrawals for an application.
func (c *eventBroadcastingClient) getWithdrawals(appID common.ApplicationIdType) ([]common.Withdrawal, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	w, exists := c.withdrawals[appID]
	return w, exists
}
