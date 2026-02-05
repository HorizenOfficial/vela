package admin

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockManagerCmdHandler is a mock implementation of the AdminCmdHandler interface.
type MockManagerCmdHandler struct {
	version   string
	err       error
	callCount int
	mu        sync.Mutex
}

func (m *MockManagerCmdHandler) ExecuteCommand(ctx context.Context, msg AdminMessage) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++

	if msg.Type != GetVersionRequestMessage {
		return nil, errors.New("unsupported command type")
	}

	return m.version, m.err
}

func (m *MockManagerCmdHandler) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func TestManagerAdminServer_StartStop(t *testing.T) {
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	defer listener.Close()

	factory := &MockConnectionFactory{listener: listener}
	server := NewAdminServer(factory, commParams, testLogger)

	// Test Start
	err := server.Start(context.Background(), "test")
	assert.NoError(t, err, "Start should not return an error")
	assert.True(t, server.isRunning, "server should be running after Start")
	assert.NotNil(t, server.listener, "listener should be set after Start")

	// Test starting again
	err = server.Start(context.Background(), "test")
	assert.Error(t, err, "Start should return an error if called again")

	// Test Stop
	err = server.Stop()
	assert.NoError(t, err, "Stop should not return an error")
	assert.False(t, server.isRunning, "server should not be running after Stop")

	// Test stopping again
	err = server.Stop()
	assert.NoError(t, err, "Stop should not return an error if called again")

	// Test that listener is closed
	_, err = listener.Accept()
	assert.Error(t, err, "listener should be closed after Stop")
}

func TestManagerAdminServer_HandleGetVersionSuccess(t *testing.T) {
	server := NewAdminServer(nil, commParams, testLogger)
	server.clientTimeout = 500 * time.Millisecond
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	handler := &MockManagerCmdHandler{version: "1.2.3"}
	server.SetCmdHandler(handler)

	go server.handleNewClient(context.Background(), serverConn, "test")

	// Send request from client side
	req := AdminMessage{Type: GetVersionRequestMessage}
	reqBytes, _ := json.Marshal(req)
	_, err := clientConn.Write(append(reqBytes, communication.MsgDelimiter))
	require.NoError(t, err)

	// Read response on client side
	respBytes, err := bufio.NewReader(clientConn).ReadBytes(communication.MsgDelimiter)
	require.NoError(t, err)

	var respMsg AdminMessage
	err = json.Unmarshal(respBytes, &respMsg)
	require.NoError(t, err)

	assert.Equal(t, AdminResponseMessage, respMsg.Type)
	assert.Equal(t, "1.2.3", respMsg.Data)
	assert.Equal(t, 1, handler.GetCallCount())
}

func TestManagerAdminServer_HandleGetVersionHandlerError(t *testing.T) {
	server := NewAdminServer(nil, commParams, testLogger)
	server.clientTimeout = 500 * time.Millisecond
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	handler := &MockManagerCmdHandler{err: errors.New("version error")}
	server.SetCmdHandler(handler)

	go server.handleNewClient(context.Background(), serverConn, "test")

	req := AdminMessage{Type: GetVersionRequestMessage}
	reqBytes, _ := json.Marshal(req)
	_, err := clientConn.Write(append(reqBytes, communication.MsgDelimiter))
	require.NoError(t, err)

	respBytes, err := bufio.NewReader(clientConn).ReadBytes(communication.MsgDelimiter)
	require.NoError(t, err)

	var respMsg AdminMessage
	json.Unmarshal(respBytes, &respMsg)

	assert.Equal(t, AdminErrorMessage, respMsg.Type)
	var errData communication.ErrorData
	dataBytes, _ := json.Marshal(respMsg.Data)
	json.Unmarshal(dataBytes, &errData)

	assert.Equal(t, "COMMAND_ERROR", errData.Code)
	assert.Equal(t, "version error", errData.Message)
	assert.Equal(t, 1, handler.GetCallCount())
}

func TestManagerAdminServer_HandleUnknownRequest(t *testing.T) {
	server := NewAdminServer(nil, commParams, testLogger)
	server.clientTimeout = 500 * time.Millisecond

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	handler := &MockManagerCmdHandler{version: "1.0.0"}
	server.SetCmdHandler(handler)

	go server.handleNewClient(context.Background(), serverConn, "test")

	req := AdminMessage{Type: 999} // Unknown type
	reqBytes, _ := json.Marshal(req)
	_, err := clientConn.Write(append(reqBytes, communication.MsgDelimiter))
	require.NoError(t, err)

	respBytes, err := bufio.NewReader(clientConn).ReadBytes(communication.MsgDelimiter)
	require.NoError(t, err)

	var respMsg AdminMessage
	json.Unmarshal(respBytes, &respMsg)

	assert.Equal(t, AdminErrorMessage, respMsg.Type)
	var errData communication.ErrorData
	dataBytes, _ := json.Marshal(respMsg.Data)
	json.Unmarshal(dataBytes, &errData)

	assert.Equal(t, "COMMAND_ERROR", errData.Code)
}

func TestManagerAdminServer_ServerBusy(t *testing.T) {
	server := NewAdminServer(nil, common.CommunicationParams{RequestTimeoutSec: 2}, testLogger)

	handler := &MockManagerCmdHandler{version: "1.0.0"}
	server.SetCmdHandler(handler)

	// Occupy the server with a "first" client
	_, serverConn1 := net.Pipe()
	defer serverConn1.Close()
	go server.handleNewClient(context.Background(), serverConn1, "test_busy_1")

	// Wait for the first client to be processed
	time.Sleep(50 * time.Millisecond)

	// Attempt to connect with a "second" client
	clientConn2, serverConn2 := net.Pipe()
	defer clientConn2.Close()

	go server.handleNewClient(context.Background(), serverConn2, "test_busy_2")

	// The second client should get a busy message
	respBytes, err := bufio.NewReader(clientConn2).ReadBytes(communication.MsgDelimiter)
	require.NoError(t, err)

	var respMsg AdminMessage
	json.Unmarshal(respBytes, &respMsg)

	assert.Equal(t, AdminErrorMessage, respMsg.Type)
	var errData communication.ErrorData
	dataBytes, _ := json.Marshal(respMsg.Data)
	json.Unmarshal(dataBytes, &errData)

	assert.Equal(t, "INVALID_REQUEST", errData.Code)
	assert.Equal(t, "server is busy", errData.Message)

	// Try again after the first client timed out
	time.Sleep(server.clientTimeout)

	require.NotNil(t, server.client) // client still pending

	clientConn3, serverConn3 := net.Pipe()
	defer clientConn3.Close()

	go server.handleNewClient(context.Background(), serverConn3, "test_not_busy")

	// Send request from third client side
	req := AdminMessage{Type: GetVersionRequestMessage}
	reqBytes, _ := json.Marshal(req)
	_, err = clientConn3.Write(append(reqBytes, communication.MsgDelimiter))
	require.NoError(t, err)

	// Read response on client side
	respBytes, err = bufio.NewReader(clientConn3).ReadBytes(communication.MsgDelimiter)
	require.NoError(t, err)

	err = json.Unmarshal(respBytes, &respMsg)
	require.NoError(t, err)

	assert.Equal(t, AdminResponseMessage, respMsg.Type)
}
