package admin

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// logLevelHandler implements AdminCmdHandler using a real logger,
// mirroring the SetLogLevel/GetLogLevel logic from the Manager.
type logLevelHandler struct {
	log logger.Logger
}

func (h *logLevelHandler) ExecuteCommand(ctx context.Context, msg AdminMessage) (interface{}, error) {
	switch msg.Type {
	case SetLogLevelRequestMessage:
		var req SetLogLevelRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			return nil, fmt.Errorf("invalid request data: %w", err)
		}
		if err := ValidateTarget(req.Target, "manager"); err != nil {
			return nil, err
		}
		znl, ok := h.log.(*logger.ZeroNetworkLogger)
		if !ok {
			return nil, fmt.Errorf("SetLogLevel is only supported with the ZeroNetworkLogger")
		}
		if req.Level == "" {
			return nil, fmt.Errorf("invalid log level: level must not be empty; supported levels: %s", SupportedLogLevels)
		}
		if err := znl.SetLevel(req.Level); err != nil {
			return nil, fmt.Errorf("invalid log level '%s'; supported levels: %s", req.Level, SupportedLogLevels)
		}
		return struct {
			Success bool   `json:"success"`
			Level   string `json:"level"`
		}{Success: true, Level: req.Level}, nil

	case GetLogLevelRequestMessage:
		var req GetLogLevelRequest
		if msg.Data != nil && string(msg.Data) != "null" {
			if err := json.Unmarshal(msg.Data, &req); err != nil {
				return nil, fmt.Errorf("invalid request data: %w", err)
			}
		}
		if err := ValidateTarget(req.Target, "manager"); err != nil {
			return nil, err
		}
		znl, ok := h.log.(*logger.ZeroNetworkLogger)
		if !ok {
			return nil, fmt.Errorf("GetLogLevel is only supported with the ZeroNetworkLogger")
		}
		return znl.GetLevel(), nil

	default:
		return nil, fmt.Errorf("unsupported command type: %v", msg.Type)
	}
}

// TestIntegration_AdminServer_SetGetLogLevel_RealTCP exercises a full TCP round-trip
// with a real AdminServer and ZeroNetworkLogger.
func TestIntegration_AdminServer_SetGetLogLevel_RealTCP(t *testing.T) {
	// 1. Start a dummy TCP listener to act as the log server (sink for ZeroNetworkLogger).
	logSink, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer logSink.Close()

	// Accept and discard connections in background
	go func() {
		for {
			conn, err := logSink.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				buf := make([]byte, 4096)
				for {
					if _, err := c.Read(buf); err != nil {
						c.Close()
						return
					}
				}
			}(conn)
		}
	}()

	logSinkAddr := logSink.Addr().(*net.TCPAddr)

	// 2. Create a real ZeroNetworkLogger pointing at the log sink.
	znl := logger.NewZeroNetworkLogger(&logger.Config{
		RemoteLogNetwork: "tcp",
		RemoteLogParams: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(logSinkAddr.Port),
		},
		NetworkLevel: "info",
	})
	defer znl.Close()

	// 3. Start a real AdminServer on a random TCP port.
	adminListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	adminAddr := adminListener.Addr().String()
	factory := &MockConnectionFactory{listener: adminListener}
	server := NewAdminServer(factory, common.CommunicationParams{RequestTimeoutSec: 5}, testLogger)
	server.SetCmdHandler(&logLevelHandler{log: znl})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = server.Start(ctx, "integration-test")
	require.NoError(t, err)
	defer server.Stop()

	// Helper: connect to admin server, send command, read response.
	sendAdminCommand := func(msg AdminMessage) AdminMessage {
		t.Helper()
		conn, err := net.DialTimeout("tcp", adminAddr, 2*time.Second)
		require.NoError(t, err)
		defer conn.Close()

		reqBytes, err := json.Marshal(msg)
		require.NoError(t, err)
		_, err = conn.Write(append(reqBytes, communication.MsgDelimiter))
		require.NoError(t, err)

		respBytes, err := bufio.NewReader(conn).ReadBytes(communication.MsgDelimiter)
		require.NoError(t, err)

		var resp AdminMessage
		require.NoError(t, json.Unmarshal(respBytes, &resp))
		return resp
	}

	// 4. GetLogLevel - should return "info" (the initial level).
	resp := sendAdminCommand(AdminMessage{Type: GetLogLevelRequestMessage})
	assert.Equal(t, AdminResponseMessage, resp.Type)
	var level string
	require.NoError(t, json.Unmarshal(resp.Data, &level))
	assert.Equal(t, "info", level)

	// Allow single-client cleanup
	time.Sleep(100 * time.Millisecond)

	// 5. SetLogLevel to "debug".
	setData, err := json.Marshal(SetLogLevelRequest{Level: "debug", Target: "manager"})
	require.NoError(t, err)
	resp = sendAdminCommand(AdminMessage{Type: SetLogLevelRequestMessage, Data: setData})
	assert.Equal(t, AdminResponseMessage, resp.Type)
	var setResp struct {
		Success bool   `json:"success"`
		Level   string `json:"level"`
	}
	require.NoError(t, json.Unmarshal(resp.Data, &setResp))
	assert.True(t, setResp.Success)
	assert.Equal(t, "debug", setResp.Level)

	// Allow single-client cleanup
	time.Sleep(100 * time.Millisecond)

	// 6. GetLogLevel again - should now return "debug".
	resp = sendAdminCommand(AdminMessage{Type: GetLogLevelRequestMessage})
	assert.Equal(t, AdminResponseMessage, resp.Type)
	require.NoError(t, json.Unmarshal(resp.Data, &level))
	assert.Equal(t, "debug", level)

	// 7. Verify the logger itself reflects the change.
	assert.Equal(t, "debug", znl.GetLevel())

	// Allow single-client cleanup
	time.Sleep(100 * time.Millisecond)

	// 8. SetLogLevel with wrong target should be rejected.
	setData, err = json.Marshal(SetLogLevelRequest{Level: "warn", Target: "executor"})
	require.NoError(t, err)
	resp = sendAdminCommand(AdminMessage{Type: SetLogLevelRequestMessage, Data: setData})
	assert.Equal(t, AdminErrorMessage, resp.Type)
	var errData communication.ErrorData
	require.NoError(t, json.Unmarshal(resp.Data, &errData))
	assert.Equal(t, "COMMAND_ERROR", errData.Code)
	assert.Contains(t, errData.Message, "this is the manager admin server")

	// Verify level was NOT changed by the rejected command.
	assert.Equal(t, "debug", znl.GetLevel())
}
