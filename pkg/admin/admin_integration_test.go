package admin

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	"github.com/HorizenOfficial/vela/pkg/communication"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCommRequestHandler implements communication.RequestHandler for integration testing.
// It handles admin commands (set/get log level) using a real logger.
type mockCommRequestHandler struct {
	log logger.Logger
}

func (h *mockCommRequestHandler) HandleProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
	return nil, nil, nil, apperrors.New(apperrors.CodeInternalFallback, "not implemented")
}

func (h *mockCommRequestHandler) HandleDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState) (*common.UpdatePayload, *common.ApplicationState, error) {
	return nil, nil, apperrors.New(apperrors.CodeInternalFallback, "not implemented")
}

func (h *mockCommRequestHandler) HandleAdminCommand(ctx context.Context, cmdType string, data json.RawMessage) (json.RawMessage, error) {
	switch cmdType {
	case AdminCmdSetLogLevel:
		var req SetLogLevelRequest
		if data != nil && string(data) != "null" {
			if err := json.Unmarshal(data, &req); err != nil {
				return nil, fmt.Errorf("invalid request data: %w", err)
			}
		}
		result, err := HandleSetLogLevel(h.log, "integration-test", req.Level)
		if err != nil {
			return nil, err
		}
		resp, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		return resp, nil
	case AdminCmdGetLogLevel:
		level, err := HandleGetLogLevel(h.log, "integration-test")
		if err != nil {
			return nil, err
		}
		resp, err := json.Marshal(level)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		return resp, nil
	default:
		return nil, fmt.Errorf("unsupported admin command type: %s", cmdType)
	}
}

// logLevelHandler implements AdminCmdHandler using a real logger,
type logLevelHandler struct {
	log logger.Logger
}

func (h *logLevelHandler) ExecuteCommand(ctx context.Context, msg AdminMessage) (interface{}, error) {
	switch msg.Type {
	case SetLogLevelRequestMessage:
		if err := ValidateTarget(msg.Target); err != nil {
			return nil, err
		}
		var req SetLogLevelRequest
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			return nil, fmt.Errorf("invalid request data: %w", err)
		}
		return HandleSetLogLevel(h.log, "integration-test", req.Level)

	case GetLogLevelRequestMessage:
		if err := ValidateTarget(msg.Target); err != nil {
			return nil, err
		}
		return HandleGetLogLevel(h.log, "integration-test")

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
	// The admin server handles one client at a time, so we retry the
	// dial until the previous connection is cleaned up.
	sendAdminCommand := func(msg AdminMessage) AdminMessage {
		t.Helper()
		var conn net.Conn
		require.Eventually(t, func() bool {
			c, dialErr := net.DialTimeout("tcp", adminAddr, 200*time.Millisecond)
			if dialErr != nil {
				return false
			}
			conn = c
			return true
		}, 2*time.Second, 20*time.Millisecond, "admin server did not become ready")
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

	// 5. SetLogLevel to "debug".
	setData, err := json.Marshal(SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	resp = sendAdminCommand(AdminMessage{Type: SetLogLevelRequestMessage, Target: "manager", Data: setData})
	assert.Equal(t, AdminResponseMessage, resp.Type)
	var setResp struct {
		Level string `json:"level"`
	}
	require.NoError(t, json.Unmarshal(resp.Data, &setResp))
	assert.Equal(t, "debug", setResp.Level)

	// 6. GetLogLevel again - should now return "debug".
	resp = sendAdminCommand(AdminMessage{Type: GetLogLevelRequestMessage})
	assert.Equal(t, AdminResponseMessage, resp.Type)
	require.NoError(t, json.Unmarshal(resp.Data, &level))
	assert.Equal(t, "debug", level)

	// 7. Verify the logger itself reflects the change.
	assert.Equal(t, "debug", znl.GetLevel())

	// 8. SetLogLevel with unknown target should be rejected.
	setData, err = json.Marshal(SetLogLevelRequest{Level: "warn"})
	require.NoError(t, err)
	resp = sendAdminCommand(AdminMessage{Type: SetLogLevelRequestMessage, Target: "foobar", Data: setData})
	assert.Equal(t, AdminErrorMessage, resp.Type)
	var errData communication.ErrorData
	require.NoError(t, json.Unmarshal(resp.Data, &errData))
	assert.Equal(t, "COMMAND_ERROR", errData.Code)
	assert.Contains(t, errData.Message, "unknown target 'foobar'")

	// Verify level was NOT changed by the rejected command.
	assert.Equal(t, "debug", znl.GetLevel())
}

// TestIntegration_AdminForwarding_RealTCP verifies that admin commands can be
// forwarded through the communication channel (Manager→Executor) using a real
// TCP communication.Server and communication.Client pair.
func TestIntegration_AdminForwarding_RealTCP(t *testing.T) {
	// 1. Start a dummy TCP listener to act as the log server (sink for ZeroNetworkLogger).
	logSink, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer logSink.Close()
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

	// 2. Create a real ZeroNetworkLogger (simulating the executor's logger).
	znl := logger.NewZeroNetworkLogger(&logger.Config{
		RemoteLogNetwork: "tcp",
		RemoteLogParams: common.TcpChannelConnectionParams{
			Ip:   "127.0.0.1",
			Port: uint32(logSinkAddr.Port),
		},
		NetworkLevel: "info",
	})
	defer znl.Close()

	// 3. Create a mock RequestHandler that handles admin commands using the real logger.
	handler := &mockCommRequestHandler{log: znl}

	// 4. Start a real communication.Server + Client on a random TCP port.
	commListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	commAddr := commListener.Addr().String()
	commListener.Close() // free the port so the factory can rebind it
	commFactory := communication.NewTCPConnectionFactory(commAddr)
	commParams := common.CommunicationParams{RequestTimeoutSec: 5}

	server := communication.NewServer(commFactory, commParams, testLogger)
	server.SetRequestHandler(handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = server.Start(ctx, "IntegrationServer")
	require.NoError(t, err)
	defer server.Stop()

	client := communication.NewClient(commFactory, commParams, testLogger)
	err = client.Connect(ctx, "IntegrationClient")
	require.NoError(t, err)
	defer client.Close()

	time.Sleep(100 * time.Millisecond)

	// 5. GetLogLevel via ForwardAdminCommand — should return "info".
	respData, err := client.ForwardAdminCommand(ctx, AdminCmdGetLogLevel, nil)
	require.NoError(t, err)
	var level string
	require.NoError(t, json.Unmarshal(respData, &level))
	assert.Equal(t, "info", level)

	// 6. SetLogLevel to "debug" via ForwardAdminCommand.
	setReqData, err := json.Marshal(SetLogLevelRequest{Level: "debug"})
	require.NoError(t, err)
	respData, err = client.ForwardAdminCommand(ctx, AdminCmdSetLogLevel, setReqData)
	require.NoError(t, err)
	var setResp SetLogLevelResponse
	require.NoError(t, json.Unmarshal(respData, &setResp))
	assert.Equal(t, "debug", setResp.Level)

	// 7. GetLogLevel again — should now return "debug".
	respData, err = client.ForwardAdminCommand(ctx, AdminCmdGetLogLevel, nil)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(respData, &level))
	assert.Equal(t, "debug", level)

	// 8. Verify the underlying logger actually changed.
	assert.Equal(t, "debug", znl.GetLevel())
}
