package subgraph

import (
	"context"
	"math/big"
	"sort"
	"strings"

	"github.com/horizen-pes/pkg/common"
)

// MockClient provides canned responses for tests.
type MockClient struct {
	requests map[common.RequestIdType]*RequestCompleted
	events   map[common.ApplicationIdType][]UserEvent
}

func NewMockClient() *MockClient {
	return &MockClient{
		requests: make(map[common.RequestIdType]*RequestCompleted),
		events:   make(map[common.ApplicationIdType][]UserEvent),
	}
}

func (m *MockClient) WithRequestCompleted(rc *RequestCompleted) *MockClient {
	if rc != nil {
		m.requests[rc.RequestID] = rc
	}
	return m
}

func (m *MockClient) WithUserEvents(appID common.ApplicationIdType, events []UserEvent) *MockClient {
	m.events[appID] = events
	return m
}

func (m *MockClient) GetRequestCompletedByID(_ context.Context, requestID common.RequestIdType) (*RequestCompleted, error) {
	if rc, ok := m.requests[requestID]; ok {
		return rc, nil
	}
	return nil, nil
}

func (m *MockClient) HealthCheck(context.Context) error {
	return nil
}

func (m *MockClient) GetUserEvents(_ context.Context, applicationID common.ApplicationIdType, eventSubType string, limit int, before *big.Int) ([]UserEvent, error) {
	all, ok := m.events[applicationID]
	if !ok {
		return nil, nil
	}

	var filtered []UserEvent
	trimmedSubType := strings.TrimSpace(eventSubType)
	for _, ev := range all {
		if trimmedSubType != "" && ev.EventSubType != trimmedSubType {
			continue
		}
		if before != nil && userEventSortKey(ev).Cmp(before) >= 0 {
			continue
		}
		filtered = append(filtered, ev)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return userEventSortKey(filtered[i]).Cmp(userEventSortKey(filtered[j])) > 0
	})

	if limit <= 0 {
		limit = 10
	}
	if limit > 1000 {
		limit = 1000
	}
	if len(filtered) <= limit {
		return filtered, nil
	}
	return filtered[:limit], nil
}
