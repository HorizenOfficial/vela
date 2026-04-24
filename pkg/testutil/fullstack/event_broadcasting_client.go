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
// connected to the simulated backend) and intercepts:
//   - GetNextPendingRequest — to cache each request's RequestType, so
//     SubmitStateUpdate can classify the completion the same way the contract
//     does (DeployRequestCompleted vs RequestCompleted).
//   - SubmitStateUpdate — to broadcast events, track completions, store
//     update payloads, accumulate withdrawals, and notify the in-process
//     subgraph. Errors from the underlying call are also captured on a
//     channel so negative-path tests (e.g., TEE attestation rejection) can
//     assert on the exact revert reason rather than waiting for an outer
//     timeout.
//
// All other Client methods are delegated to the underlying client unchanged.
type eventBroadcastingClient struct {
	blockchain.Client // embedded real client — all methods delegated by default

	mu             sync.Mutex
	eventChannel   chan<- interface{}
	pendingIDs     map[common.RequestIdType]struct{}
	requestTypes   map[common.RequestIdType]common.RequestType // cached from GetNextPendingRequest
	completedIDs   map[common.RequestIdType]struct{}
	failedIDs      map[common.RequestIdType]struct{}
	updatePayloads map[common.RequestIdType]*common.UpdatePayload
	withdrawals    map[common.ApplicationIdType][]common.Withdrawal

	// stateUpdateErrors buffers errors returned by the underlying
	// SubmitStateUpdate. Non-blocking send: if the channel is full, older
	// errors are kept and the new one is dropped (so the manager's polling
	// loop is never blocked). Tests drain via WaitForStateUpdateError.
	stateUpdateErrors chan error

	// onStateUpdate is called after each successful SubmitStateUpdate.
	// Used by InProcessSubgraph to record data for wallet queries.
	onStateUpdate func(update *common.UpdatePayload, isDeploy bool)
}

func newEventBroadcastingClient(inner blockchain.Client, eventCh chan<- interface{}) *eventBroadcastingClient {
	return &eventBroadcastingClient{
		Client:            inner,
		eventChannel:      eventCh,
		pendingIDs:        make(map[common.RequestIdType]struct{}),
		requestTypes:      make(map[common.RequestIdType]common.RequestType),
		completedIDs:      make(map[common.RequestIdType]struct{}),
		failedIDs:         make(map[common.RequestIdType]struct{}),
		updatePayloads:    make(map[common.RequestIdType]*common.UpdatePayload),
		withdrawals:       make(map[common.ApplicationIdType][]common.Withdrawal),
		stateUpdateErrors: make(chan error, 16),
	}
}

// markPending registers a request as pending so AssertRequestCompleted can
// track it. Request-type tracking is populated automatically via
// GetNextPendingRequest — callers do not need to supply it.
func (c *eventBroadcastingClient) markPending(requestID common.RequestIdType) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pendingIDs[requestID] = struct{}{}
}

// GetNextPendingRequest delegates to the real client and caches the request's
// RequestType so SubmitStateUpdate can later classify the completion
// canonically (matches the contract's own DeployRequestCompleted vs
// RequestCompleted emission logic).
func (c *eventBroadcastingClient) GetNextPendingRequest(ctx context.Context) (*common.Request, [32]byte, error) {
	req, stateRoot, err := c.Client.GetNextPendingRequest(ctx)
	if err != nil {
		return req, stateRoot, err
	}
	if req != nil {
		c.mu.Lock()
		c.requestTypes[req.RequestID] = req.RequestType
		c.mu.Unlock()
	}
	return req, stateRoot, nil
}

// SubmitStateUpdate delegates to the real client and then records the result
// for test observation (events, completions, withdrawals).
func (c *eventBroadcastingClient) SubmitStateUpdate(ctx context.Context, update *common.UpdatePayload) error {
	// Delegate to the real blockchain client (on-chain transaction)
	err := c.Client.SubmitStateUpdate(ctx, update)
	if err != nil {
		// Non-blocking capture so tests can assert on the specific revert
		// reason. Still propagate the error to the manager unchanged.
		select {
		case c.stateUpdateErrors <- err:
		default:
		}
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove from pending
	delete(c.pendingIDs, update.RequestID)

	// Classify using the cached RequestType (captured from GetNextPendingRequest).
	// Matches the contract's DeployRequestCompleted vs RequestCompleted emission
	// logic at ProcessorEndpoint.sol:401.
	reqType, ok := c.requestTypes[update.RequestID]
	isDeploy := ok && reqType == common.Deploy
	delete(c.requestTypes, update.RequestID)

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

// waitForStateUpdateError blocks up to `timeout` for a SubmitStateUpdate
// error to arrive from the underlying blockchain client. Returns the error
// if one occurred in that window, or (nil, false) on timeout. Used by
// negative-path tests to pin down the specific on-chain revert reason
// (e.g. InvalidSignature on a TEE-attestation-rejection test).
func (c *eventBroadcastingClient) waitForStateUpdateError(timeout time.Duration) (error, bool) {
	select {
	case err := <-c.stateUpdateErrors:
		return err, true
	case <-time.After(timeout):
		return nil, false
	}
}
