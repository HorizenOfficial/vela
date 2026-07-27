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
//   - GetNextPendingRequest / GetPendingRequestsWithStateRoot — to cache each
//     request's RequestType, so SubmitStateUpdate / SubmitBatchStateUpdate can classify the completion
//     the same way the contract does (DeployRequestCompleted vs
//     RequestCompleted), and to track TRUSTPROCESS requests for test discovery.
//   - SubmitStateUpdate / SubmitBatchStateUpdate — to broadcast events, track
//     completions, store update payloads, accumulate withdrawals, and notify
//     the in-process subgraph. Errors from the underlying call are also
//     captured on a channel so negative-path tests (e.g., TEE attestation
//     rejection) can assert on the exact revert reason rather than waiting for
//     an outer timeout.
//
// All other Client methods are delegated to the underlying client unchanged.
type eventBroadcastingClient struct {
	blockchain.Client // embedded real client — all methods delegated by default

	mu             sync.Mutex
	eventChannel   chan<- interface{}
	pendingIDs     map[common.RequestIdType]struct{}
	requestTypes   map[common.RequestIdType]common.RequestType // cached from GetNextPendingRequest / GetPendingRequestsWithStateRoot
	completedIDs   map[common.RequestIdType]struct{}
	failedIDs      map[common.RequestIdType]struct{}
	updatePayloads map[common.RequestIdType]*common.UpdatePayload
	withdrawals    map[common.ApplicationIdType][]common.Withdrawal

	// trustProcessIDs records, in fetch order, the request IDs the manager pulled
	// that were of type TrustProcess. These requests are enqueued on-chain by the
	// trigger (not submitted by the test), so this is how a test discovers them.
	trustProcessIDs []common.RequestIdType

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
		c.cacheRequestLocked(req)
		c.mu.Unlock()
	}
	return req, stateRoot, nil
}

// GetPendingRequestsWithStateRoot delegates to the real client (the contract selects
// the application and returns up to maxCount of its pending requests) and caches each
// request's RequestType so a later SubmitBatchStateUpdate/SubmitStateUpdate can
// classify the completion canonically. This is the batch-path counterpart of
// GetNextPendingRequest — tests that drive the manager through batching rely on it to
// populate the same request-type and TRUSTPROCESS bookkeeping.
func (c *eventBroadcastingClient) GetPendingRequestsWithStateRoot(ctx context.Context, maxCount uint64) (common.ApplicationIdType, []*common.Request, [32]byte, error) {
	appID, requests, stateRoot, err := c.Client.GetPendingRequestsWithStateRoot(ctx, maxCount)
	if err != nil {
		return appID, requests, stateRoot, err
	}
	if len(requests) > 0 {
		c.mu.Lock()
		for _, req := range requests {
			if req != nil {
				c.cacheRequestLocked(req)
			}
		}
		c.mu.Unlock()
	}
	return appID, requests, stateRoot, nil
}

// cacheRequestLocked records a fetched request's RequestType so a later state update
// can classify its completion the same way the contract does, and tracks TRUSTPROCESS
// requests (enqueued on-chain by a trigger, not submitted by the test) so tests can
// discover them. Callers must hold c.mu.
func (c *eventBroadcastingClient) cacheRequestLocked(req *common.Request) {
	c.requestTypes[req.RequestID] = req.RequestType
	if req.RequestType == common.TrustProcess && !containsRequestID(c.trustProcessIDs, req.RequestID) {
		c.trustProcessIDs = append(c.trustProcessIDs, req.RequestID)
	}
}

func containsRequestID(ids []common.RequestIdType, id common.RequestIdType) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
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
	c.recordUpdateLocked(update)
	return nil
}

// SubmitBatchStateUpdate delegates the batch to the real client and then records
// each entry's result for test observation, mirroring SubmitStateUpdate. The batch
// is atomic on-chain, so bookkeeping runs only after the underlying call succeeds;
// on error the reason is captured on stateUpdateErrors (like SubmitStateUpdate) and
// no entry is recorded. Batches never contain deploy requests — deploys are always
// processed individually — so every entry classifies as a non-deploy completion.
func (c *eventBroadcastingClient) SubmitBatchStateUpdate(ctx context.Context, updates []*common.UpdatePayload, batchSignature []byte) error {
	// Delegate to the real blockchain client (single on-chain transaction)
	err := c.Client.SubmitBatchStateUpdate(ctx, updates, batchSignature)
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
	for _, update := range updates {
		c.recordUpdateLocked(update)
	}
	return nil
}

// recordUpdateLocked records the effect of a successfully-submitted update payload
// for test observation: it clears pending state, classifies the completion, and
// tracks failures/completions, stored payloads, broadcast events, and withdrawals.
// Callers must hold c.mu.
func (c *eventBroadcastingClient) recordUpdateLocked(update *common.UpdatePayload) {
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
		return
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

// waitForRequestFailed blocks until the given request has been marked FAILED by
// the executor (a signed error stateUpdate), or times out. The negative
// counterpart of waitForRequestCompletion.
func (c *eventBroadcastingClient) waitForRequestFailed(requestID common.RequestIdType, timeout time.Duration) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			_, failed := c.failedIDs[requestID]
			_, completed := c.completedIDs[requestID]
			c.mu.Unlock()
			if failed {
				return nil
			}
			if completed {
				return fmt.Errorf("request %s completed successfully, expected failure", requestID)
			}
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for request %s to fail", requestID)
		}
	}
}

// waitForTrustProcess blocks until a TRUSTPROCESS request (enqueued on-chain by
// a trigger and pulled by the manager) has completed, returning its request ID
// and stored UpdatePayload. Fails if such a request completes with an error.
func (c *eventBroadcastingClient) waitForTrustProcess(timeout time.Duration) (common.RequestIdType, *common.UpdatePayload, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for _, id := range c.trustProcessIDs {
				if _, failed := c.failedIDs[id]; failed {
					c.mu.Unlock()
					return id, nil, fmt.Errorf("TRUSTPROCESS request %s has failed", id)
				}
				if _, completed := c.completedIDs[id]; completed {
					payload := c.updatePayloads[id]
					c.mu.Unlock()
					return id, payload, nil
				}
			}
			c.mu.Unlock()
		case <-timeoutCh:
			return common.RequestIdType{}, nil, fmt.Errorf("timeout waiting for a TRUSTPROCESS request to complete")
		}
	}
}

// waitForFailedTrustProcess blocks until a TRUSTPROCESS request (enqueued
// on-chain by a trigger and pulled by the manager) has been marked FAILED by the
// executor (a signed error stateUpdate), returning its request ID. This is the
// negative counterpart of waitForTrustProcess.
func (c *eventBroadcastingClient) waitForFailedTrustProcess(timeout time.Duration) (common.RequestIdType, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for _, id := range c.trustProcessIDs {
				if _, failed := c.failedIDs[id]; failed {
					c.mu.Unlock()
					return id, nil
				}
			}
			c.mu.Unlock()
		case <-timeoutCh:
			return common.RequestIdType{}, fmt.Errorf("timeout waiting for a TRUSTPROCESS request to fail")
		}
	}
}

// trustProcessCount returns how many distinct TRUSTPROCESS requests the manager
// has pulled so far. Tests use it to assert the trigger loop terminated (exactly
// one) rather than re-enqueuing endlessly.
func (c *eventBroadcastingClient) trustProcessCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.trustProcessIDs)
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
