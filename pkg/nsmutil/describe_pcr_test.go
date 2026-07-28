package nsmutil

import (
	"fmt"
	"testing"

	"github.com/hf/nsm/request"
	"github.com/hf/nsm/response"
	"github.com/stretchr/testify/require"
)

type mockSession struct {
	resp    response.Response
	sendErr error
	closed  bool
}

func (m *mockSession) Send(_ request.Request) (response.Response, error) {
	return m.resp, m.sendErr
}

func (m *mockSession) Close() error {
	m.closed = true
	return nil
}

func TestDescribePCRWithSession_Success(t *testing.T) {
	want := make([]byte, 48)
	for i := range want {
		want[i] = byte(i)
	}
	sess := &mockSession{resp: response.Response{DescribePCR: &response.DescribePCR{Data: want}}}

	got, err := DescribePCRWithSession(func() (Session, error) { return sess, nil }, 0)
	require.NoError(t, err)
	require.Equal(t, want, got)
	require.True(t, sess.closed, "session should be closed")
}

func TestDescribePCRWithSession_OpenerError(t *testing.T) {
	_, err := DescribePCRWithSession(func() (Session, error) {
		return nil, fmt.Errorf("no nsm device")
	}, 0)
	require.Error(t, err)
}

func TestDescribePCRWithSession_DeviceError(t *testing.T) {
	sess := &mockSession{resp: response.Response{Error: "InvalidIndex"}}
	_, err := DescribePCRWithSession(func() (Session, error) { return sess, nil }, 0)
	require.Error(t, err)
}

func TestDescribePCRWithSession_EmptyResponse(t *testing.T) {
	sess := &mockSession{resp: response.Response{}}
	_, err := DescribePCRWithSession(func() (Session, error) { return sess, nil }, 0)
	require.Error(t, err)
}
