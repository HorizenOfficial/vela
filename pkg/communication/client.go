package communication

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/common"
)

// Client is a unified client implementation of the ExecutorClient interface
type Client struct {
	conn            net.Conn
	connLock        sync.Mutex
	connected       bool
	reader          *bufio.Reader
	writer          *bufio.Writer
	pendingRequests map[string]*PendingRequest
	pendingMu       sync.Mutex
	shutdown        chan struct{}
	reqTimeout      time.Duration
	factory         ConnectionFactory
	requestHandler  ClientRequestHandler
}

// NewClient creates a new client with the specified connection factory
func NewClient(factory ConnectionFactory) *Client {
	return &Client{
		factory:         factory,
		pendingRequests: make(map[string]*PendingRequest),
		shutdown:        make(chan struct{}),
		reqTimeout:      30 * time.Second,
	}
}

// Connect establishes a connection using the connection factory
func (c *Client) Connect(ctx context.Context) error {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if c.connected {
		return fmt.Errorf("already connected")
	}

	// Create connection using factory
	conn, err := c.factory.CreateClientConnection()
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	c.connected = true

	// Start the message reader goroutine
	go MessageReaderLoop(
		ctx,
		"Client",
		c.conn,
		c.reader,
		c.shutdown,
		func(ctx context.Context, msg *Message) {
			c.routeIncomingMessage(ctx, msg)
		},
		func() {
			c.Close()
		},
	)

	// Start the cleanup goroutine for timed-out requests
	go c.cleanupLoop()

	return nil
}

// Close closes the connection and stops all goroutines
func (c *Client) Close() error {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if !c.connected {
		return nil
	}

	c.connected = false
	close(c.shutdown)

	// Cancel all pending requests
	c.pendingMu.Lock()
	for id, req := range c.pendingRequests {
		close(req.ResponseChan)
		delete(c.pendingRequests, id)
	}
	c.pendingMu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}

	log.Println("Client connection closed")
	return nil
}

// SetClientRequestHandler sets the handler for server-initiated requests
func (c *Client) SetClientRequestHandler(handler ClientRequestHandler) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.requestHandler = handler
}

// SendProcessRequest sends a process request and waits for response
func (c *Client) SendProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	msg := Message{
		ID:   generateID(),
		Type: ProcessRequestMessage,
		Data: ProcessRequestData{
			Request:          req,
			ApplicationState: appState,
			SenderKey:        senderKey,
			WasmModule:       wasmModule,
		},
	}

	respMsg, err := c.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return nil, nil, err
	}

	if respMsg.Type == ErrorMessage {
		errorData, _ := extractData[ErrorData](respMsg.Data)
		return nil, nil, fmt.Errorf("server error: %s", errorData.Message)
	}

	if respMsg.Type != ProcessResponseMessage {
		return nil, nil, fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	respData, err := extractData[ProcessResponseData](respMsg.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract response data: %w", err)
	}

	return respData.UpdatePayload, respData.UpdatedApplicationState, nil
}

// SendDeployApp sends a deploy app request and waits for response
func (c *Client) SendDeployApp(ctx context.Context, req *common.Request) (*common.UpdatePayload, *common.ApplicationState, error) {
	uid := generateID()
	log.Printf("Generated UID: %s", uid)

	msg := Message{
		ID:   uid,
		Type: DeployAppRequestMessage,
		Data: DeployAppRequestData{
			Request: req,
		},
	}

	respMsg, err := c.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return nil, nil, err
	}

	if respMsg.Type == ErrorMessage {
		errorData, _ := extractData[ErrorData](respMsg.Data)
		return nil, nil, fmt.Errorf("server error: %s", errorData.Message)
	}

	if respMsg.Type != DeployAppResponseMessage {
		return nil, nil, fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	respData, err := extractData[DeployAppResponseData](respMsg.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract response data: %w", err)
	}

	return respData.UpdatePayload, respData.ApplicationState, nil
}

// SendGenerateDeanonymizationReport sends a deanonymization request and waits for response
func (c *Client) SendGenerateDeanonymizationReport(ctx context.Context, req *common.Request, appState *common.ApplicationState, senderKey []byte, wasmModule []byte) (*common.DeanonymizationReport, error) {
	msg := Message{
		ID:   generateID(),
		Type: DeanonymizationRequestMessage,
		Data: DeanonymizationRequestData{
			Request:          req,
			ApplicationState: appState,
			SenderKey:        senderKey,
			WasmModule:       wasmModule,
		},
	}

	respMsg, err := c.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return nil, err
	}

	if respMsg.Type == ErrorMessage {
		errorData, _ := extractData[ErrorData](respMsg.Data)
		return nil, fmt.Errorf("server error: %s", errorData.Message)
	}

	if respMsg.Type != DeanonymizationResponseMessage {
		return nil, fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	respData, err := extractData[DeanonymizationResponseData](respMsg.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to extract response data: %w", err)
	}

	return respData.Report, nil
}

// sendRequestAndWaitForResponse sends a request and waits for the response
func (c *Client) sendRequestAndWaitForResponse(ctx context.Context, msg Message) (*Message, error) {
	// Create response channel
	responseChan := make(chan *Message, 1)

	// Store pending request
	c.pendingMu.Lock()
	c.pendingRequests[msg.ID] = &PendingRequest{
		ResponseChan: responseChan,
		Timeout:      time.Now().Add(c.reqTimeout),
	}
	c.pendingMu.Unlock()

	// Ensure cleanup on exit
	defer func() {
		c.pendingMu.Lock()
		delete(c.pendingRequests, msg.ID)
		c.pendingMu.Unlock()
	}()

	// Send message
	if err := c.sendMessage(msg); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Wait for response or timeout
	select {
	case response, ok := <-responseChan:
		if !ok {
			return nil, fmt.Errorf("request timeout or response channel closed")
		}
		return response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.shutdown:
		return nil, fmt.Errorf("client is shutting down")
	}
}

// sendMessage sends a message to the server
func (c *Client) sendMessage(msg Message) error {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if !c.connected {
		return fmt.Errorf("not connected")
	}

	// Serialize message
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	log.Printf("MagBytes length before delimiter: %d", len(data))

	// Add delimiter
	data = append(data, delimiter)
	log.Printf("MagBytes length after delimiter: %d", len(data))

	// Write message with delimiter
	if _, err := c.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if err := c.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

// routeIncomingMessage routes incoming messages to appropriate handlers
func (c *Client) routeIncomingMessage(ctx context.Context, msg *Message) {
	log.Printf("Client: acquiring lock, msg.type=%d", msg.Type)
	c.pendingMu.Lock()
	pendingReq, exists := c.pendingRequests[msg.ID]
	c.pendingMu.Unlock()
	log.Printf("Client: lock released. exists=%d", exists)

	if exists {
		/*
			Select on the response channel to avoid blocking.
			If the response channel is ready, send the response.
			If the channel is full or closed, ignore the response.
		*/
		select {
		case pendingReq.ResponseChan <- msg:
			// Channel is open and has room, send the response
		default:
			// Channel is full or closed, ignore and log the issue
			log.Printf("Client: Warning: response channel for request ID %s is full or closed, ignoring response\n", msg.ID)
		}
		return
	}

	// Handle server-initiated requests
	c.handleServerRequest(ctx, msg)
}

// handleServerRequest handles requests initiated by the server
func (c *Client) handleServerRequest(ctx context.Context, msg *Message) {
	switch msg.Type {
	case GetUserKeysRequestMessage:
		c.handleGetUserKeysRequest(ctx, msg)
	default:
		log.Printf("Client: Warning: unknown message type: %v\n", msg.Type)
	}
}

// handleGetUserKeysRequest handles GetUserKeys requests from server
func (c *Client) handleGetUserKeysRequest(ctx context.Context, msg *Message) {
	if c.requestHandler == nil {
		c.sendErrorResponse(msg.ID, "NO_HANDLER", fmt.Errorf("no request handler set"))
		return
	}

	reqData, err := extractData[GetUserKeysRequestData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	userKeys, err := c.requestHandler.HandleGetUserKeys(ctx, reqData.Users)
	if err != nil {
		c.sendErrorResponse(msg.ID, "HANDLER_ERROR", err)
		return
	}

	response := Message{
		ID:   msg.ID,
		Type: GetUserKeysResponseMessage,
		Data: GetUserKeysResponseData{
			UserKeys: userKeys,
		},
	}

	if err := c.sendMessage(response); err != nil {
		log.Printf("Client: Failed to send HandleGetUserKeys response: %v", err)
	}
}

// sendErrorResponse sends an error response
func (c *Client) sendErrorResponse(requestID string, code string, err error) {
	log.Printf("Client: Sending error response: ID=%s, Code=%s, Message=%s", requestID, code, err.Error())
	response := Message{
		ID:   requestID,
		Type: ErrorMessage,
		Data: ErrorData{
			Code:    code,
			Message: err.Error(),
		},
	}

	if sendErr := c.sendMessage(response); sendErr != nil {
		log.Printf("Client: Failed to send error response: %v", sendErr)
	}
}

// cleanupLoop periodically cleans up timed-out requests
func (c *Client) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupTimedOutRequests()
		case <-c.shutdown:
			return
		}
	}
}

// cleanupTimedOutRequests removes requests that have timed out
func (c *Client) cleanupTimedOutRequests() {
	now := time.Now()
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()

	for id, req := range c.pendingRequests {
		if now.After(req.Timeout) {
			close(req.ResponseChan)
			delete(c.pendingRequests, id)
		}
	}
}
