package fullstack

import (
	"context"
	"math/big"
	"sort"
	"sync"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela-common-go/subgraph"
	"github.com/HorizenOfficial/vela/pkg/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

// InProcessSubgraph implements subgraph.Client backed by in-memory data populated
// from the eventBroadcastingClient's SubmitStateUpdate interception. It replaces
// a real graph-node for fullstack tests, providing the wallet with the same query
// interface without any external dependencies.
type InProcessSubgraph struct {
	mu                 sync.RWMutex
	requestCompletions map[velacommon.RequestIdType]*subgraph.RequestCompleted
	deployCompletions  map[velacommon.RequestIdType]*subgraph.RequestCompleted
	userEvents         []subgraph.UserEvent
	withdrawals        []subgraph.OnChainWithdrawal
	nextSortKey        uint64
	nextBlockNumber    uint64
}

func NewInProcessSubgraph() *InProcessSubgraph {
	return &InProcessSubgraph{
		requestCompletions: make(map[velacommon.RequestIdType]*subgraph.RequestCompleted),
		deployCompletions:  make(map[velacommon.RequestIdType]*subgraph.RequestCompleted),
		nextBlockNumber:    1,
	}
}

// RecordStateUpdate is called by the eventBroadcastingClient after each
// SubmitStateUpdate (both successful and failed). It extracts RequestCompleted,
// UserEvent, and Withdrawal records from the update payload and stores them
// for later query.
func (s *InProcessSubgraph) RecordStateUpdate(update *common.UpdatePayload, isDeploy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	blockNumber := s.nextBlockNumber
	s.nextBlockNumber++

	// Determine status
	status := velacommon.RequestResultOK
	if update.ErrorCode != 0 {
		status = velacommon.RequestResultFailed
	}

	rc := &subgraph.RequestCompleted{
		ApplicationID:   update.ApplicationID,
		RequestID:       update.RequestID,
		Status:          status,
		ErrorCode:       update.ErrorCode,
		ApplicationFees: update.ApplicationFee.ToInt(),
		BlockNumber:     blockNumber,
	}

	if isDeploy {
		s.deployCompletions[update.RequestID] = rc
	} else {
		s.requestCompletions[update.RequestID] = rc
	}

	// Store user events
	for i, evt := range update.Events {
		s.nextSortKey++
		s.userEvents = append(s.userEvents, subgraph.UserEvent{
			ApplicationID: update.ApplicationID,
			RequestID:     update.RequestID,
			EventSubType:  evt.EventSubType,
			EncryptedData: evt.EncryptedData,
			BlockNumber:   blockNumber,
			LogIndex:      uint64(i),
			SortKey:       new(big.Int).SetUint64(s.nextSortKey),
		})
	}

	// Store withdrawals
	for _, w := range update.Withdrawals {
		s.withdrawals = append(s.withdrawals, subgraph.OnChainWithdrawal{
			ApplicationID: update.ApplicationID,
			RequestID:     update.RequestID,
			To:            w.DestinationAddress,
			TokenAddress:  w.TokenAddress,
			Amount:        w.Amount.ToInt(),
			BlockNumber:   blockNumber,
		})
	}
}

// InjectRequestCompleted inserts a synthetic RequestCompleted record directly
// into the in-memory store, bypassing the usual RecordStateUpdate path. Used by
// tests that need to exercise downstream consumers (e.g. the authority service's
// /getreport handler, which queries the subgraph for completion status) without
// driving a real request through manager/executor.
func (s *InProcessSubgraph) InjectRequestCompleted(appID velacommon.ApplicationIdType, requestID velacommon.RequestIdType, status velacommon.RequestResultStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	blockNumber := s.nextBlockNumber
	s.nextBlockNumber++

	s.requestCompletions[requestID] = &subgraph.RequestCompleted{
		ApplicationID: appID,
		RequestID:     requestID,
		Status:        status,
		BlockNumber:   blockNumber,
	}
}

// --- subgraph.Client implementation ---

func (s *InProcessSubgraph) HealthCheck(_ context.Context) error {
	return nil
}

func (s *InProcessSubgraph) GetRequestCompletedByID(_ context.Context, requestID velacommon.RequestIdType) (*subgraph.RequestCompleted, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rc, ok := s.requestCompletions[requestID]
	if !ok {
		return nil, nil // not found yet — wallet will retry
	}
	return rc, nil
}

func (s *InProcessSubgraph) GetDeployRequestCompletedByID(_ context.Context, requestID velacommon.RequestIdType) (*subgraph.RequestCompleted, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rc, ok := s.deployCompletions[requestID]
	if !ok {
		return nil, nil
	}
	return rc, nil
}

func (s *InProcessSubgraph) GetUserEvents(ctx context.Context, applicationID velacommon.ApplicationIdType, eventSubType [32]byte, limit int, before *big.Int) ([]subgraph.UserEvent, error) {
	var subTypes [][32]byte
	if eventSubType != ([32]byte{}) {
		subTypes = [][32]byte{eventSubType}
	}
	return s.GetUserEventsBySubTypes(ctx, applicationID, subTypes, limit, before)
}

func (s *InProcessSubgraph) GetUserEventsBySubTypes(_ context.Context, applicationID velacommon.ApplicationIdType, eventSubTypes [][32]byte, limit int, before *big.Int) ([]subgraph.UserEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	subtypeSet := make(map[[32]byte]bool, len(eventSubTypes))
	for _, st := range eventSubTypes {
		subtypeSet[st] = true
	}

	var filtered []subgraph.UserEvent
	for _, ev := range s.userEvents {
		if ev.ApplicationID != applicationID {
			continue
		}
		if len(subtypeSet) > 0 && !subtypeSet[ev.EventSubType] {
			continue
		}
		if before != nil && ev.SortKey.Cmp(before) >= 0 {
			continue
		}
		filtered = append(filtered, ev)
	}

	// Sort descending by SortKey (newest first) — matches real subgraph behavior
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].SortKey.Cmp(filtered[j].SortKey) > 0
	})

	if limit <= 0 {
		limit = 10
	}
	if limit > 1000 {
		limit = 1000
	}
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered, nil
}

// GetAppEvents / GetAppEventsBySubTypes: stubs. InProcessSubgraph only records
// UserEvents (encrypted, per-user) from the state-update path; app-level
// events are not produced by the current fullstack test flows. If a future
// test needs them, extend RecordStateUpdate to also populate an appEvents
// slice and filter here identically to the UserEvent path.
func (s *InProcessSubgraph) GetAppEvents(_ context.Context, _ velacommon.ApplicationIdType, _ [32]byte, _ int, _ *big.Int) ([]subgraph.AppEvent, error) {
	return nil, nil
}

func (s *InProcessSubgraph) GetAppEventsBySubTypes(_ context.Context, _ velacommon.ApplicationIdType, _ [][32]byte, _ int, _ *big.Int) ([]subgraph.AppEvent, error) {
	return nil, nil
}

func (s *InProcessSubgraph) GetRefunds(_ context.Context, _ velacommon.ApplicationIdType, _ *velacommon.RequestIdType, _ int) ([]subgraph.OnChainRefund, error) {
	// Stub — not needed for core wallet flow (deposit/withdraw/getprivatebalance)
	return nil, nil
}

func (s *InProcessSubgraph) GetWithdrawals(_ context.Context, applicationID velacommon.ApplicationIdType, requestID *velacommon.RequestIdType, limit int) ([]subgraph.OnChainWithdrawal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []subgraph.OnChainWithdrawal
	for _, w := range s.withdrawals {
		if w.ApplicationID != applicationID {
			continue
		}
		if requestID != nil && w.RequestID != *requestID {
			continue
		}
		filtered = append(filtered, w)
	}

	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered, nil
}

func (s *InProcessSubgraph) GetClaimsExecuted(_ context.Context, _ ethCommon.Address, _ *ethCommon.Address, _ int) ([]subgraph.ClaimExecuted, error) {
	// Stub — not needed for core wallet flow
	return nil, nil
}
