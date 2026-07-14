package executor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/admin"
	"github.com/stretchr/testify/require"
)

// TestHandleAdminCommand_Shutdown verifies the executor acks a shutdown command
// and closes its ShutdownRequested channel so main() can exit cleanly.
func TestHandleAdminCommand_Shutdown(t *testing.T) {
	exec := &StatelessExecutor{
		log:        testLogger,
		shutdownCh: make(chan struct{}),
	}

	// Not yet signalled.
	select {
	case <-exec.ShutdownRequested():
		t.Fatal("shutdown should not be signalled before the command")
	default:
	}

	resp, err := exec.HandleAdminCommand(context.Background(), admin.AdminCmdShutdown, nil)
	require.NoError(t, err)

	var sr admin.ShutdownResponse
	require.NoError(t, json.Unmarshal(resp, &sr))
	require.True(t, sr.Stopping, "shutdown ack should report Stopping=true")

	// The channel must now be closed.
	select {
	case <-exec.ShutdownRequested():
	default:
		t.Fatal("ShutdownRequested channel should be closed after the shutdown command")
	}

	// Idempotent: a second request must not panic (double close).
	require.NotPanics(t, func() { exec.requestShutdown() })
}
