package executor

import (
	"context"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/communication"
	"github.com/stretchr/testify/require"
)

// The guest execution bound only guarantees forward progress if it fires BEFORE the
// manager stops waiting. The manager's patience is configured on the parent EC2, out of
// the enclave's sight, so the handshake carries it and the executor validates against
// that — not against its own, unrelated, request timeout. These tests pin that the
// handshake is where the check happens and that a rejected handshake persists nothing.

// scriptedServerConnection is a ServerConnection whose manager reply is fixed up
// front, so performHandshake can be driven without a real channel.
type scriptedServerConnection struct {
	reply *communication.GetKeysetRecoveryResponseData

	setKeysetCalled bool
	resultReported  bool
	resultErr       error
}

func (c *scriptedServerConnection) GetKeysetRecovery(context.Context) (*communication.GetKeysetRecoveryResponseData, error) {
	return c.reply, nil
}

func (c *scriptedServerConnection) SetKeysetRecovery(context.Context, *common.EnclaveKeySetRecovery, string, string) error {
	c.setKeysetCalled = true
	return nil
}

func (c *scriptedServerConnection) KeysetRecoveryResult(_ context.Context, result error, _, _ string) error {
	c.resultReported = true
	c.resultErr = result
	return nil
}

func (c *scriptedServerConnection) Close() {}

// newHandshakeExecutor builds an executor that generates a fresh keyset on a first
// connection (unsafe recovery type: no KMS needed) with the given guest bound.
func newHandshakeExecutor(guestTimeoutMs int64) *StatelessExecutor {
	return &StatelessExecutor{
		config: &Config{
			KeySetRecoveryType:      common.RecoveryTypeUnsafe,
			GuestExecutionTimeoutMs: guestTimeoutMs,
		},
		log: testLogger,
	}
}

// firstConnectionReply is what a manager with no stored keyset answers, reporting
// the given request timeout.
func firstConnectionReply(requestTimeoutMs int64) *communication.GetKeysetRecoveryResponseData {
	return &communication.GetKeysetRecoveryResponseData{
		DataFound:        false,
		RequestTimeoutMs: requestTimeoutMs,
	}
}

func TestHandshakeValidatesGuestBoundAgainstManagerTimeout(t *testing.T) {
	const guestTimeoutMs = 10_000
	marginMs := communication.ExecutionBudgetMargin.Milliseconds()

	requireRejected := func(t *testing.T, conn *scriptedServerConnection, err error) {
		t.Helper()
		require.Error(t, err)
		for _, name := range []string{"EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS", "MANAGER_COMMUNICATION_PARAMS_REQUEST_TIMEOUT_SEC"} {
			require.Contains(t, err.Error(), name, "the operator must be told which two settings disagree")
		}
		require.True(t, conn.resultReported, "the manager must be told why the handshake failed")
		require.Error(t, conn.resultErr)
		require.Contains(t, conn.resultErr.Error(), "EXECUTOR_GUEST_EXECUTION_TIMEOUT_MS")
		require.False(t, conn.setKeysetCalled, "a rejected handshake must not generate and persist a keyset")
	}

	t.Run("manager patience below the guest bound is rejected", func(t *testing.T) {
		conn := &scriptedServerConnection{reply: firstConnectionReply(5_000)}
		_, err := newHandshakeExecutor(guestTimeoutMs).performHandshake(context.Background(), conn)
		requireRejected(t, conn, err)
	})

	t.Run("patience inside the reply margin is rejected", func(t *testing.T) {
		// 11 s of patience minus the 2 s reply margin leaves 9 s: the budget still
		// beats a 10 s bound, so this is the silent window a bare comparison misses.
		conn := &scriptedServerConnection{reply: firstConnectionReply(guestTimeoutMs + marginMs - 1_000)}
		_, err := newHandshakeExecutor(guestTimeoutMs).performHandshake(context.Background(), conn)
		requireRejected(t, conn, err)
	})

	t.Run("patience equal to bound plus margin is rejected", func(t *testing.T) {
		// Equal lets the two deadlines race; the bound must be strictly below.
		conn := &scriptedServerConnection{reply: firstConnectionReply(guestTimeoutMs + marginMs)}
		_, err := newHandshakeExecutor(guestTimeoutMs).performHandshake(context.Background(), conn)
		requireRejected(t, conn, err)
	})

	t.Run("adequate patience completes the handshake", func(t *testing.T) {
		conn := &scriptedServerConnection{reply: firstConnectionReply(30_000)}
		keySet, err := newHandshakeExecutor(guestTimeoutMs).performHandshake(context.Background(), conn)
		require.NoError(t, err)
		require.NotNil(t, keySet)
		require.True(t, conn.setKeysetCalled, "a first connection must hand the new keyset to the manager")
	})

	t.Run("the default bound is resolved before comparing", func(t *testing.T) {
		// 0 means the 10 s default, which 5 s of patience cannot cover.
		conn := &scriptedServerConnection{reply: firstConnectionReply(5_000)}
		_, err := newHandshakeExecutor(0).performHandshake(context.Background(), conn)
		requireRejected(t, conn, err)
	})

	t.Run("a manager that does not report its patience is served", func(t *testing.T) {
		// An older manager omits the field. Refusing it would break every mixed-version
		// deployment for a check that was previously not performed at all.
		conn := &scriptedServerConnection{reply: firstConnectionReply(0)}
		_, err := newHandshakeExecutor(guestTimeoutMs).performHandshake(context.Background(), conn)
		require.NoError(t, err)
		require.True(t, conn.setKeysetCalled)
	})
}
