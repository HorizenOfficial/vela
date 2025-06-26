package communication

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/horizen-pes/pkg/common"
)

// TCPClient is a TCP implementation of the ExecutorClient interface
type TCPClient struct {
	addr     string
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	mu       sync.Mutex
	connLock sync.Mutex
}

// NewTCPClient creates a new TCP client
func NewTCPClient(addr string) *TCPClient {
	return &TCPClient{
		addr: addr,
	}
}

// Connect establishes a connection to the executor
func (c *TCPClient) Connect(ctx context.Context) error {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if c.conn != nil {
		return nil // Already connected
	}

	// Establish TCP connection
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return fmt.Errorf("failed to connect to TCP server: %w", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	return nil
}

// Close closes the connection to the executor
func (c *TCPClient) Close() error {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if c.conn == nil {
		return nil // Already closed
	}

	err := c.conn.Close()
	c.conn = nil
	c.reader = nil
	c.writer = nil
	return err
}

// ProcessRequest sends a request to the executor and returns the response
func (c *TCPClient) ProcessRequest(
	ctx context.Context,
	req *common.Request,
	appState *common.ApplicationState,
	wasmModule []byte,
) (*common.UpdatePayload, *common.ApplicationState, error) {
	// Create the message
	msg := Message{
		Type: ProcessRequestMessage,
		Data: ProcessRequestData{
			Request:          req,
			ApplicationState: appState,
			WasmModule:       wasmModule,
		},
	}

	// Send the message and get the response
	respMsg, err := c.sendMessage(ctx, msg)
	if err != nil {
		return nil, nil, err
	}

	// Check the response type
	if respMsg.Type == ErrorMessage {
		errData, ok := respMsg.Data.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("invalid error message format")
		}
		return nil, nil, fmt.Errorf("executor error: %s - %s", errData["code"], errData["message"])
	}

	if respMsg.Type != ProcessResponseMessage {
		return nil, nil, fmt.Errorf("unexpected response type: %d", respMsg.Type)
	}

	// Extract the response data
	respData, err := extractProcessResponseData(respMsg.Data)
	if err != nil {
		return nil, nil, err
	}

	return respData.UpdatePayload, respData.UpdatedApplicationState, nil
}

// DeployApp deploys a new application to the executor
func (c *TCPClient) DeployApp(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error) {
	// Create the message
	msg := Message{
		Type: DeployAppRequestMessage,
		Data: DeployAppRequestData{
			Request: req,
		},
	}

	// Send the message and get the response
	respMsg, err := c.sendMessage(ctx, msg)
	if err != nil {
		return nil, nil, err
	}

	// Check the response type
	if respMsg.Type == ErrorMessage {
		errData, ok := respMsg.Data.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("invalid error message format")
		}
		return nil, nil, fmt.Errorf("executor error: %s - %s", errData["code"], errData["message"])
	}

	if respMsg.Type != DeployAppResponseMessage {
		return nil, nil, fmt.Errorf("unexpected response type: %d", respMsg.Type)
	}

	// Extract the response data
	respData, err := extractDeployAppResponseData(respMsg.Data)
	if err != nil {
		return nil, nil, err
	}

	return respData.ApplicationState, respData.WasmModule, nil
}

// GenerateDeanonymizationReport generates a deanonymization report
func (c *TCPClient) GenerateDeanonymizationReport(
	ctx context.Context,
	req *common.Request,
	appState *common.ApplicationState,
	wasmModule []byte,
) (*common.DeanonymizationReport, error) {
	// Create the message
	msg := Message{
		Type: DeanonymizationRequestMessage,
		Data: DeanonymizationRequestData{
			Request:          req,
			ApplicationState: appState,
			WasmModule:       wasmModule,
		},
	}

	// Send the message and get the response
	respMsg, err := c.sendMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Check the response type
	if respMsg.Type == ErrorMessage {
		errData, ok := respMsg.Data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid error message format")
		}
		return nil, fmt.Errorf("executor error: %s - %s", errData["code"], errData["message"])
	}

	if respMsg.Type != DeanonymizationResponseMessage {
		return nil, fmt.Errorf("unexpected response type: %d", respMsg.Type)
	}

	// Extract the response data
	respData, err := extractDeanonymizationResponseData(respMsg.Data)
	if err != nil {
		return nil, err
	}

	return respData.Report, nil
}

// sendMessage sends a message to the executor and returns the response
func (c *TCPClient) sendMessage(ctx context.Context, msg Message) (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil, fmt.Errorf("not connected to server")
	}

	// Marshal the message
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Add a newline to separate messages
	msgBytes = append(msgBytes, '\n')

	// Write the message
	_, err = c.writer.Write(msgBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	err = c.writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to flush message: %w", err)
	}

	// Read the response
	respBytes, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal the response
	var respMsg Message
	if err := json.Unmarshal(respBytes, &respMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &respMsg, nil
}
