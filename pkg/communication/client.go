package communication

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
	storageErrors "github.com/horizen-pes/pkg/storage/errors"
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
	idLogTag        string
	log             logger.Logger
}

// NewClient creates a new client with the specified connection factory
func NewClient(factory ConnectionFactory, communicationParams common.CommunicationParams, log logger.Logger) *Client {
	return &Client{
		factory:         factory,
		pendingRequests: make(map[string]*PendingRequest),
		shutdown:        make(chan struct{}),
		reqTimeout:      communicationParams.RequestTimeoutSec * time.Second,
		log:             log,
	}
}

// Connect establishes a connection using the connection factory
func (c *Client) Connect(ctx context.Context, idLogTag string) error {
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
	c.idLogTag = idLogTag

	// Start the message reader goroutine
	go MessageReaderLoop(
		ctx,
		idLogTag,
		c.conn,
		c.reader,
		c.shutdown,
		func(ctx context.Context, msg Message) {
			c.routeIncomingMessage(ctx, msg)
		},
		func() {
			c.Close()
		},
		c.log,
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

	c.log.Warn("Client connection closed")
	return nil
}

// SetClientRequestHandler sets the handler for server-initiated requests
func (c *Client) SetClientRequestHandler(handler ClientRequestHandler) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.requestHandler = handler
}

// SendProcessRequest sends a process request and waits for response
// Returns UpdatePayload, ApplicationState, optional DeanonymizationReport, and error
func (c *Client) SendProcessRequest(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, *common.DeanonymizationReport, error) {
	msg := Message{
		ID:   generateID(),
		Type: ProcessRequestMessage,
		Data: ProcessRequestData{
			Request:          req,
			ApplicationState: appState,
			WasmModule:       wasmModule,
		},
	}

	respMsg, err := c.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to send process request: %w", err)
	}

	if respMsg.Type == ErrorMessage {
		errorData, err := extractData[ErrorData](respMsg.Data)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to extract server error data: %w", err)
		}
		return nil, nil, nil, fmt.Errorf("server error: %s", errorData.Message)
	}

	if respMsg.Type != ProcessResponseMessage {
		return nil, nil, nil, fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	respData, err := extractData[ProcessResponseData](respMsg.Data)
	if err != nil {
		return nil, nil, nil, err
	}

	return respData.UpdatePayload, respData.UpdatedApplicationState, respData.DeanonymizationReport, nil
}

// SendDeployApp sends a deploy app request and waits for response
func (c *Client) SendDeployApp(ctx context.Context, req *common.Request, appState *common.ApplicationState, wasmModule []byte) (*common.UpdatePayload, *common.ApplicationState, error) {
	uid := generateID()
	c.log.Debug("Generated UID: %s", uid)

	msg := Message{
		ID:   uid,
		Type: DeployAppRequestMessage,
		Data: DeployAppRequestData{
			Request:          req,
			ApplicationState: appState,
			WasmModule:       wasmModule,
		},
	}

	respMsg, err := c.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send deploy app request: %w", err)
	}

	if respMsg.Type == ErrorMessage {
		errorData, err := extractData[ErrorData](respMsg.Data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to extract server error data: %w", err)
		}
		return nil, nil, fmt.Errorf("server error: %s", errorData.Message)
	}

	if respMsg.Type != DeployAppResponseMessage {
		return nil, nil, fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	respData, err := extractData[DeployAppResponseData](respMsg.Data)
	if err != nil {
		return nil, nil, err
	}

	return respData.UpdatePayload, respData.ApplicationState, nil
}

// SendKeyAttestationRequest sends a key attestation request to the executor and waits for the response
func (c *Client) SendKeyAttestationRequest(ctx context.Context) ([]byte, error) {
	msg := Message{
		ID:   generateID(),
		Type: KeyAttestationRequestMessage,
	}

	respMsg, err := c.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send key attestation request: %w", err)
	}

	if respMsg.Type == ErrorMessage {
		errorData, err := extractData[ErrorData](respMsg.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		return nil, fmt.Errorf("key attestation failed [%s]: %s", errorData.Code, errorData.Message)
	}

	if respMsg.Type != KeyAttestationResponseMessage {
		return nil, fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	respData, err := extractData[KeyAttestationResponseData](respMsg.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to extract response data: %w", err)
	}

	return respData.Attestation, nil
}

// sendRequestAndWaitForResponse sends a request and waits for the response
func (c *Client) sendRequestAndWaitForResponse(ctx context.Context, msg Message) (Message, error) {
	// Create response channel
	responseChan := make(chan Message, 1)

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
		return Message{}, fmt.Errorf("%s: failed to send message: %w", c.idLogTag, err)
	}

	// Wait for response or timeout
	select {
	case response, ok := <-responseChan:
		if !ok {
			return Message{}, fmt.Errorf("request timeout or response channel closed")
		}
		return response, nil
	case <-ctx.Done():
		return Message{}, ctx.Err()
	case <-c.shutdown:
		return Message{}, fmt.Errorf("client is shutting down")
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
	c.log.Debug("%s: MagBytes length before delimiter: %d", c.idLogTag, len(data))

	// Add delimiter
	data = append(data, MsgDelimiter)
	c.log.Debug("%s: MagBytes length after delimiter: %d", c.idLogTag, len(data))

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
func (c *Client) routeIncomingMessage(ctx context.Context, msg Message) {
	c.pendingMu.Lock()
	pendingReq, exists := c.pendingRequests[msg.ID]
	c.pendingMu.Unlock()

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
			c.log.Warn("%s: Warning: response channel for request ID %s is full or closed, ignoring response\n", c.idLogTag, msg.ID)
		}
		return
	}

	// Handle server-initiated requests
	c.handleServerRequest(ctx, msg)
}

// handleServerRequest handles requests initiated by the server
func (c *Client) handleServerRequest(ctx context.Context, msg Message) {
	switch msg.Type {
	case GetKeysetRecoveryRequestMessage:
		c.log.Info("%s: got message GetKeysetRecoveryRequestMessage type: %v", c.idLogTag, msg.Type)
		c.handleGetKeysetRecoveryRequest(ctx, msg)
	case SetKeysetRecoveryRequestMessage:
		c.log.Info("%s: got message SetKeysetRecoveryRequestMessage type: %v", c.idLogTag, msg.Type)
		c.handleSetKeysetRecoveryRequest(ctx, msg)
	case KeysetRecoveryResultMessage:
		c.handleKeysetRecoveryResult(ctx, msg)
	default:
		c.log.Warn("%s: Warning: unknown message type: %v", c.idLogTag, msg.Type)
	}
}

// handleGetKeysetRecoveryRequest handles Executor messages from the server
func (c *Client) handleGetKeysetRecoveryRequest(ctx context.Context, msg Message) {
	c.log.Info("%s: entering %s", c.idLogTag, common.FnName())
	if c.requestHandler == nil {
		c.sendErrorResponse(msg.ID, "NO_HANDLER", fmt.Errorf("no request handler set"))
		return
	}

	recv, err := c.requestHandler.HandleGetKeysetRecoveryRequest(ctx)

	var dataFound bool
	var respRecv *common.EnclaveKeySetRecovery

	if err != nil {
		if storageErrors.IsNotFound(err) {
			c.log.Warn("%s: KeysetRecovery not found in data layer: %v", c.idLogTag, err)
			dataFound = false
			respRecv = nil
		} else {
			c.log.Error("%s: Unexpected error from datalayer: %v", c.idLogTag, err)
			c.sendErrorResponse(msg.ID, "Unexpected error from dataLayer ", err)
			return
		}
	} else {
		c.log.Warn("%s: KeysetRecovery found in data layer", c.idLogTag)
		dataFound = true
		respRecv = recv
	}

	response := Message{
		ID:   msg.ID,
		Type: GetKeysetRecoveryResponseMessage,
		Data: GetKeysetRecoveryResponseData{
			DataFound:      dataFound,
			KeySetRecovery: respRecv,
		},
	}

	c.log.Info("%s: Sending GetKeysetRecoveryResponseMessage to executor", c.idLogTag)
	err = c.sendMessage(response)
	if err != nil {
		c.log.Warn("%s: Failed to send GetKeysetRecoveryResponseMessage response: %v", c.idLogTag, err)
	}
}

func (c *Client) handleKeysetRecoveryResult(ctx context.Context, msg Message) {
	if c.requestHandler == nil {
		c.sendErrorResponse(msg.ID, "NO_HANDLER", fmt.Errorf("no request handler set"))
		return
	}

	reqData, err := extractData[KeysetRecoveryResultData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	var result error
	if reqData.Error != "" {
		result = fmt.Errorf("%s", reqData.Error)
	}

	err = c.requestHandler.HandleKeysetRecoveryResult(ctx, result, reqData.CommPubKey, reqData.SigningKeyAddr)
	if err != nil {
		c.sendErrorResponse(msg.ID, "HANDLER_ERROR", err)
		return
	}
}

// handleSetKeysetRecoveryRequest handles a request to set the keyset recovery data
func (c *Client) handleSetKeysetRecoveryRequest(ctx context.Context, msg Message) {
	if c.requestHandler == nil {
		c.sendErrorResponse(msg.ID, "NO_HANDLER", fmt.Errorf("no request handler set"))
		return
	}

	reqData, err := extractData[SetKeysetRecoveryRequestData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	err = c.requestHandler.HandleSetKeysetRecoveryRequest(ctx, reqData.KeySetRecovery, reqData.CommPubKey, reqData.SigningKeyAddr)
	if err != nil {
		c.sendErrorResponse(msg.ID, "HANDLER_ERROR", err)
		return
	}

	response := Message{
		ID:   msg.ID,
		Type: SetKeysetRecoveryResponseMessage,
		Data: SetKeysetRecoveryResponseData{},
	}

	if err := c.sendMessage(response); err != nil {
		c.log.Error("%s: Failed to send HandleSetKeysetRecoveryRequest response: %v", c.idLogTag, err)
		return
	}
	c.log.Info("%s: SetKeysetRecoveryRequest handled successfully, ID=%s", c.idLogTag, msg.ID)
}

// sendErrorResponse sends an error response
func (c *Client) sendErrorResponse(requestID string, code string, err error) {
	c.log.Info("%s: Sending error response: ID=%s, Code=%s, Message=%s", c.idLogTag, requestID, code, err.Error())
	response := Message{
		ID:   requestID,
		Type: ErrorMessage,
		Data: ErrorData{
			Code:    code,
			Message: err.Error(),
		},
	}

	if sendErr := c.sendMessage(response); sendErr != nil {
		c.log.Error("%s: Failed to send error response: %v", c.idLogTag, sendErr)
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
