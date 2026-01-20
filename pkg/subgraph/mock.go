package subgraph

import (
	"context"
	"fmt"

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

func (m *MockClient) GetUserEvents(_ context.Context, applicationID common.ApplicationIdType, eventSubType string, limit int, skip int) ([]UserEvent, error) {
	all, ok := m.events[applicationID]
	if !ok {
		return nil, nil
	}
	if skip < 0 {
		return nil, fmt.Errorf("invalid skip %d", skip)
	}

	var filtered []UserEvent
	for _, ev := range all {
		if eventSubType != "" && ev.EventSubType != eventSubType {
			continue
		}
		filtered = append(filtered, ev)
	}

	if skip >= len(filtered) {
		return nil, nil
	}
	filtered = filtered[skip:]

	if limit <= 0 {
		return filtered, nil
	}
	if len(filtered) <= limit {
		return filtered, nil
	}
	if limit < 0 {
		return nil, fmt.Errorf("invalid limit %d", limit)
	}
	return filtered[:limit], nil
}
