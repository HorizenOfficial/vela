package communication

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/mdlayher/vsock"
)

// VSockClient is a v-socket implementation of the ExecutorClient interface
type VSockClient struct {
	cid      uint32
	port     uint32
	conn     *vsock.Conn
	mu       sync.Mutex
	connLock sync.Mutex
}

// NewVSockClient creates a new v-socket client
func NewVSockClient(cidStr, portStr string) *VSockClient {
	cidUint, _ := strconv.ParseUint(cidStr, 10, 32)
	portUint, _ := strconv.ParseUint(portStr, 10, 32)

	cid := uint32(cidUint)
	port := uint32(portUint)
	return &VSockClient{
		cid:  cid,
		port: port,
	}
}

// Connect establishes a connection to the executor
func (c *VSockClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	// Connect to the v-socket
	conn, err := vsock.Dial(c.cid, c.port, &vsock.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to executor: %w", err)
	}
	// Set a read deadline for the response (e.g., 5 seconds)
	deadline := 5 // seconds
	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(deadline) * time.Second)); err != nil {
		conn.Close()
		return fmt.Errorf("failed to set read deadline: %w", err)
	}

	c.conn = conn
	return nil
}

// Close closes the connection to the executor
func (c *VSockClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	return err
}

// ProcessRequest sends a request to the executor and returns the response
func (c *VSockClient) ProcessRequest(
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

	// Send the message
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
func (c *VSockClient) DeployApp(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error) {
	// Create the message
	msg := Message{
		Type: DeployAppRequestMessage,
		Data: DeployAppRequestData{
			Request: req,
		},
	}

	// Send the message
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
func (c *VSockClient) GenerateDeanonymizationReport(
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

	// Send the message
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
func (c *VSockClient) sendMessage(ctx context.Context, msg Message) (*Message, error) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// Check if the connection is established
	if c.conn == nil {
		return nil, fmt.Errorf("not connected to executor")
	}

	// Set a deadline for the operation
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(30 * time.Second)
	}
	if err := c.conn.SetDeadline(deadline); err != nil {
		return nil, fmt.Errorf("failed to set deadline: %w", err)
	}

	// Marshal the message
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Write the message length
	msgLen := uint32(len(msgBytes))
	if err := writeUint32(c.conn, msgLen); err != nil {
		return nil, fmt.Errorf("failed to write message length: %w", err)
	}

	// Write the message
	if _, err := c.conn.Write(msgBytes); err != nil {
		return nil, fmt.Errorf("failed to write message: %w", err)
	}

	// Read the response length
	respLen, err := readUint32(c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read response length: %w", err)
	}

	// Read the response
	respBytes := make([]byte, respLen)
	if _, err := io.ReadFull(c.conn, respBytes); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unmarshal the response
	var respMsg Message
	if err := json.Unmarshal(respBytes, &respMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &respMsg, nil
}

func writeUint32(w io.Writer, n uint32) error {
	return binary.Write(w, binary.BigEndian, n)
}

func readUint32(r io.Reader) (uint32, error) {
	var n uint32
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}
