package executor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/communication"
)

// mockServerConn is a minimal communication.ServerConnection used to drive
// StatelessExecutor.performHandshake in isolation.
type mockServerConn struct {
	found       bool
	recovery    *common.EnclaveKeySetRecovery
	getErr      error
	setCalled   bool
	resultCalls int
}

func (m *mockServerConn) GetKeysetRecovery(ctx context.Context) (bool, *common.EnclaveKeySetRecovery, error) {
	return m.found, m.recovery, m.getErr
}

func (m *mockServerConn) SetKeysetRecovery(ctx context.Context, recovery *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr, pcr0 string) error {
	m.setCalled = true
	return nil
}

func (m *mockServerConn) KeysetRecoveryResult(ctx context.Context, result error, commPubKey, signingKeyAddr, pcr0 string) error {
	m.resultCalls++
	return nil
}

func (m *mockServerConn) Close() {}

var _ communication.ServerConnection = (*mockServerConn)(nil)

func newContinuityExecutor(t *testing.T, expectExisting bool) *StatelessExecutor {
	t.Helper()
	builder, err := NewMsgToSignBuilder()
	require.NoError(t, err)
	return &StatelessExecutor{
		config: &Config{
			KeySetRecoveryType:   common.RecoveryTypeUnsafe,
			ExpectExistingKeyset: expectExisting,
		},
		MsgToSignBuilder: builder,
		log:              testLogger,
	}
}

// TestPerformHandshake_ExpectExisting_FoundFalse_Aborts verifies the R2/G3
// guard: with ExpectExistingKeyset set, a found=false response aborts the
// handshake with a typed error and does NOT generate or store a keyset.
func TestPerformHandshake_ExpectExisting_FoundFalse_Aborts(t *testing.T) {
	e := newContinuityExecutor(t, true)
	conn := &mockServerConn{found: false}

	keySet, err := e.performHandshake(context.Background(), conn)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnexpectedKeysetGeneration)
	require.Nil(t, keySet, "no keyset should be produced")
	require.False(t, conn.setCalled, "SetKeysetRecovery must not be called — nothing generated or stored")
	require.Nil(t, e.keySet, "executor keyset must remain unset")
}

// TestPerformHandshake_FirstInstall_FoundFalse_Generates verifies the
// first-install path is unchanged when the guard is off: found=false generates
// a new keyset and sends it to the manager for storage.
func TestPerformHandshake_FirstInstall_FoundFalse_Generates(t *testing.T) {
	e := newContinuityExecutor(t, false)
	conn := &mockServerConn{found: false}

	keySet, err := e.performHandshake(context.Background(), conn)

	require.NoError(t, err)
	require.NotNil(t, keySet, "a new keyset should be generated")
	require.True(t, conn.setCalled, "SetKeysetRecovery must be called to store the new keyset")
}

// TestPerformHandshake_ExpectExisting_FoundTrue_Restores verifies the guard
// does not interfere with the recovery path: found=true still restores.
func TestPerformHandshake_ExpectExisting_FoundTrue_Restores(t *testing.T) {
	// Generate a real keyset + recovery blob to feed back as "found".
	_, recovery, err := GenerateEnclaveKeySet(context.Background(), common.RecoveryTypeUnsafe, nil, nil, "")
	require.NoError(t, err)

	e := newContinuityExecutor(t, true)
	conn := &mockServerConn{found: true, recovery: recovery}

	keySet, err := e.performHandshake(context.Background(), conn)

	require.NoError(t, err)
	require.NotNil(t, keySet, "keyset should be restored from recovery data")
	require.False(t, conn.setCalled, "restore path must not store a new keyset")
	require.Equal(t, 1, conn.resultCalls, "restore path should confirm via KeysetRecoveryResult")
}
