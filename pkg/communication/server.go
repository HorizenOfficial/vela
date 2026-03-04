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
)

// ClientConnection represents a connection to a client
type ClientConnection struct {
	conn            net.Conn
	connLock        sync.Mutex
	connected       bool
	reader          *bufio.Reader
	writer          *bufio.Writer
	pendingRequests map[string]*PendingRequest
	pendingMu       sync.Mutex
	shutdown        chan struct{}
	reqTimeout      time.Duration
	idLogTag        string
	log             logger.Logger
}

// Server is a unified server implementation of the ExecutorServer interface
type Server struct {
	factory           ConnectionFactory
	listener          net.Listener
	isRunning         bool
	mu                sync.Mutex
	client            *ClientConnection
	clientMu          sync.RWMutex
	handler           RequestHandler
	connectionHandler ConnectionHandler
	shutdownChan      chan struct{}
	reqTimeout        time.Duration
	log               logger.Logger
}

// NewServer creates a new server with the specified connection factory
func NewServer(factory ConnectionFactory, communicationParams common.CommunicationParams, log logger.Logger) *Server {
	return &Server{
		factory:      factory,
		shutdownChan: make(chan struct{}),
		reqTimeout:   communicationParams.RequestTimeoutSec * time.Second,
		log:          log,
	}
}

// Start starts the server and begins listening for connections
func (s *Server) Start(ctx context.Context, idLogTag string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("server already running")
	}

	// Create listener using factory
	listener, err := s.factory.CreateServerListener()
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.log.Info("%s: server listening on %s %s", idLogTag, listener.Addr().Network(), listener.Addr().String())

	s.listener = listener
	s.isRunning = true

	// Start accepting connections in a goroutine
	go s.acceptConnections(ctx, idLogTag)

	return nil
}

// Stop stops the server and closes all connections
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	s.isRunning = false
	close(s.shutdownChan)

	// Close listener
	if s.listener != nil {
		s.listener.Close()
	}

	// Close client connection
	s.clientMu.Lock()
	if s.client != nil {
		s.client.internalClose()
		s.client = nil
	}
	s.clientMu.Unlock()

	return nil
}

// SetRequestHandler sets the handler for client requests
func (s *Server) SetRequestHandler(handler RequestHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
}

// SetConnectionHandler sets the handler for new client connections.
func (s *Server) SetConnectionHandler(handler ConnectionHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connectionHandler = handler
}

// acceptConnections accepts incoming connections
func (s *Server) acceptConnections(ctx context.Context, idLogTag string) {
	for {
		select {
		case <-s.shutdownChan:
			return
		case <-ctx.Done():
			return
		default:
			// Continue accepting
		}

		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdownChan:
				return
			default:
				s.log.Warn("%s: Error accepting connection: %v", idLogTag, err)
				continue
			}
		}

		// Handle new client connection
		go s.handleNewClient(ctx, conn, idLogTag)
	}
}

// handleNewClient handles a new client connection
func (s *Server) handleNewClient(ctx context.Context, conn net.Conn, idLogTag string) {
	s.log.Info("%s: New client connected from %s", idLogTag, conn.RemoteAddr())

	client := &ClientConnection{
		conn:            conn,
		reader:          bufio.NewReader(conn),
		writer:          bufio.NewWriter(conn),
		connected:       true,
		pendingRequests: make(map[string]*PendingRequest),
		shutdown:        make(chan struct{}),
		reqTimeout:      s.reqTimeout,
		idLogTag:        idLogTag,
		log:             s.log,
	}

	// Set as current client (only one client supported)
	s.clientMu.Lock()
	if s.client != nil {
		s.client.internalClose()
	}
	s.client = client
	s.clientMu.Unlock()

	// Start client message handling
	go MessageReaderLoop(
		ctx,
		idLogTag,
		client.conn,
		client.reader,
		client.shutdown,
		func(ctx context.Context, msg Message) {
			client.routeIncomingMessage(ctx, msg, s)
		},
		client.internalClose,
		s.log,
	)
	go client.cleanupLoop()

	s.mu.Lock()
	handler := s.connectionHandler
	s.mu.Unlock()

	if handler != nil {
		// We use a goroutine to avoid blocking the accept loop
		go handler(ctx, client)
	}
}

// GetKeysetRecovery sends a request to the client to get the keyset recovery data.
func (c *ClientConnection) GetKeysetRecovery(ctx context.Context) (bool, *common.EnclaveKeySetRecovery, error) {
	msg := Message{
		ID:   generateID(),
		Type: GetKeysetRecoveryRequestMessage,
		Data: GetKeysetRecoveryRequestData{},
	}

	c.log.Info("%s: Sending GetKeysetRecoveryRequestMessage (type %d) msg to Manager", c.idLogTag, GetKeysetRecoveryRequestMessage)
	respMsg, err := c.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return false, nil, err
	}

	if respMsg.Type == ErrorMessage {
		errorData, err := extractData[ErrorData](respMsg.Data)
		if err != nil {
			return false, nil, fmt.Errorf("failed to extract client error data: %w", err)
		}
		return false, nil, fmt.Errorf("client error: %s", errorData.Message)
	}

	if respMsg.Type != GetKeysetRecoveryResponseMessage {
		return false, nil, fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	respData, err := extractData[GetKeysetRecoveryResponseData](respMsg.Data)
	if err != nil {
		return false, nil, fmt.Errorf("failed to extract response data: %w", err)
	}

	return respData.DataFound, respData.KeySetRecovery, nil
}

// KeysetRecoveryResult sends a confirmation to the client for the keyset recovery.
func (c *ClientConnection) KeysetRecoveryResult(ctx context.Context, result error, commPubKey, signingKeyAddr string) error {
	var errMsg string
	if result != nil {
		errMsg = result.Error()
	}

	msg := Message{
		ID:   generateID(),
		Type: KeysetRecoveryResultMessage,
		Data: KeysetRecoveryResultData{
			Error:          errMsg,
			CommPubKey:     commPubKey,
			SigningKeyAddr: signingKeyAddr,
		},
	}

	c.log.Info("%s: sending key set recovery result to manager", c.idLogTag)
	err := c.sendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

// SetKeysetRecovery sends a request to the client to set the keyset recovery data.
func (c *ClientConnection) SetKeysetRecovery(ctx context.Context, recovery *common.EnclaveKeySetRecovery, commPubKey, signingKeyAddr string) error {
	msg := Message{
		ID:   generateID(),
		Type: SetKeysetRecoveryRequestMessage,
		Data: SetKeysetRecoveryRequestData{
			KeySetRecovery: recovery,
			CommPubKey:     commPubKey,
			SigningKeyAddr: signingKeyAddr,
		},
	}

	c.log.Info("%s: sending key set recovery to manager", c.idLogTag)
	respMsg, err := c.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return err
	}

	if respMsg.Type == ErrorMessage {
		errorData, err := extractData[ErrorData](respMsg.Data)
		if err != nil {
			return fmt.Errorf("failed to extract client error data: %w", err)
		}
		return fmt.Errorf("client error: %s", errorData.Message)
	}

	if respMsg.Type != SetKeysetRecoveryResponseMessage {
		return fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	_, err = extractData[SetKeysetRecoveryResponseData](respMsg.Data)
	if err != nil {
		return fmt.Errorf("failed to extract response data: %w", err)
	}

	return nil
}

// Close closes the client connection
func (c *ClientConnection) Close() {
	c.internalClose()
}

// internalClose closes the client connection
func (c *ClientConnection) internalClose() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if !c.connected {
		return
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
		c.conn.Close()
	}
}

// sendRequestAndWaitForResponse sends a request and waits for the response
func (c *ClientConnection) sendRequestAndWaitForResponse(ctx context.Context, msg Message) (Message, error) {
	// Create response channel
	responseChan := make(chan Message, 1)
	pendingReq := &PendingRequest{
		ResponseChan: responseChan,
		Timeout:      time.Now().Add(c.reqTimeout),
	}

	// Store pending request
	c.pendingMu.Lock()
	c.pendingRequests[msg.ID] = pendingReq
	c.pendingMu.Unlock()

	// Ensure cleanup on exit
	defer func() {
		c.pendingMu.Lock()
		delete(c.pendingRequests, msg.ID)
		c.pendingMu.Unlock()
	}()

	// Send message
	if err := c.sendMessage(msg); err != nil {
		return Message{}, fmt.Errorf("failed to send message: %w", err)
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
		return Message{}, fmt.Errorf("client connection is shutting down")
	}
}

// sendMessage sends a message to the client
func (c *ClientConnection) sendMessage(msg Message) error {
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
	c.log.Debug("%s: MsgBytes length before delimiter: %d", c.idLogTag, len(data))

	// Add newline delimiter
	data = append(data, MsgDelimiter)
	c.log.Debug("%s: MsgBytes length after delimiter: %d", c.idLogTag, len(data))

	// Write a message
	if _, err := c.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if err := c.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

// routeIncomingMessage routes incoming messages to appropriate handlers
func (c *ClientConnection) routeIncomingMessage(ctx context.Context, msg Message, server *Server) {
	// Check if this is a response to a pending request
	c.pendingMu.Lock()
	pendingReq, exists := c.pendingRequests[msg.ID]
	c.pendingMu.Unlock()

	if exists {
		// This is a response to a pending request
		select {
		case pendingReq.ResponseChan <- msg:
		default:
			// Channel is full or closed, ignore
		}
		return
	}

	// This is a client-initiated request
	c.handleClientRequest(ctx, msg, server)
}

// handleClientRequest handles requests initiated by the client
func (c *ClientConnection) handleClientRequest(ctx context.Context, msg Message, server *Server) {
	server.mu.Lock()
	handler := server.handler
	server.mu.Unlock()

	if handler == nil {
		c.sendErrorResponse(msg.ID, "NO_HANDLER", fmt.Errorf("no request handler set"))
		return
	}

	switch msg.Type {
	case ProcessRequestMessage:
		c.handleProcessRequest(ctx, msg, handler)
	case DeployAppRequestMessage:
		c.handleDeployAppRequest(ctx, msg, handler)
	case BuildErrorPayloadRequestMessage:
		c.handleBuildErrorPayloadRequest(ctx, msg, handler)
	case KeyAttestationRequestMessage:
		c.handleKeyAttestationRequest(ctx, msg, handler)
	default:
		c.sendErrorResponse(msg.ID, "UNKNOWN_REQUEST", fmt.Errorf("unknown request type: %v", msg.Type))
	}
}

// handleProcessRequest handles ProcessRequest messages
func (c *ClientConnection) handleProcessRequest(ctx context.Context, msg Message, handler RequestHandler) {
	reqData, err := extractData[ProcessRequestData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	updatePayload, updatedState, deanonymizationReport, err := handler.HandleProcessRequest(ctx, reqData.Request, reqData.ApplicationState, reqData.WasmModule)
	if err != nil {
		c.sendErrorResponse(msg.ID, "HANDLER_ERROR", err)
		return
	}

	response := Message{
		ID:   msg.ID,
		Type: ProcessResponseMessage,
		Data: ProcessResponseData{
			UpdatePayload:           updatePayload,
			UpdatedApplicationState: updatedState,
			DeanonymizationReport:   deanonymizationReport,
		},
	}

	if err := c.sendMessage(response); err != nil {
		c.log.Warn("%s: Failed to send HandleProcessRequest response: %v", c.idLogTag, err)
	}
	c.log.Info("%s: ProcessRequest handled successfully, ID=%s", c.idLogTag, msg.ID)
}

// handleDeployAppRequest handles DeployApp messages
func (c *ClientConnection) handleDeployAppRequest(ctx context.Context, msg Message, handler RequestHandler) {
	reqData, err := extractData[DeployAppRequestData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	updatePayload, appState, err := handler.HandleDeployApp(ctx, reqData.Request, reqData.ApplicationState)
	if err != nil {
		c.sendErrorResponse(msg.ID, "HANDLER_ERROR", err)
		return
	}

	response := Message{
		ID:   msg.ID,
		Type: DeployAppResponseMessage,
		Data: DeployAppResponseData{
			UpdatePayload:    updatePayload,
			ApplicationState: appState,
		},
	}

	if err := c.sendMessage(response); err != nil {
		c.log.Warn("%s: Failed to send HandleDeployApp response: %v", c.idLogTag, err)
	}
	c.log.Info("%s: DeployApp handled successfully, ID=%s", c.idLogTag, msg.ID)
}

func (c *ClientConnection) handleBuildErrorPayloadRequest(ctx context.Context, msg Message, handler RequestHandler) {
	reqData, err := extractData[BuildErrorPayloadRequestData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	updatePayload, err := handler.HandleBuildErrorPayloadRequest(ctx, reqData.Request, reqData.StateRoot, reqData.Failure)
	if err != nil {
		c.sendErrorResponse(msg.ID, "HANDLER_ERROR", err)
		return
	}

	response := Message{
		ID:   msg.ID,
		Type: BuildErrorPayloadResponseMessage,
		Data: BuildErrorPayloadResponseData{
			UpdatePayload: updatePayload,
		},
	}

	if err := c.sendMessage(response); err != nil {
		c.log.Warn("%s: Failed to send HandleBuildErrorPayloadRequest response: %v", c.idLogTag, err)
	}
	c.log.Info("%s: BuildErrorPayloadRequest handled successfully, ID=%s", c.idLogTag, msg.ID)
}

// handleKeyAttestationRequest handles KeyAttestationRequest messages
func (c *ClientConnection) handleKeyAttestationRequest(ctx context.Context, msg Message, handler RequestHandler) {
	attestation, err := handler.HandleKeyAttestationRequest(ctx)
	if err != nil {
		c.sendErrorResponse(msg.ID, "ATTESTATION_ERROR", err)
		return
	}

	response := Message{
		ID:   msg.ID,
		Type: KeyAttestationResponseMessage,
		Data: KeyAttestationResponseData{
			Attestation: attestation,
		},
	}

	if err := c.sendMessage(response); err != nil {
		c.log.Warn("%s: Failed to send HandleKeyAttestationRequest response: %v", c.idLogTag, err)
	}
	c.log.Info("%s: KeyAttestationRequest handled successfully, ID=%s", c.idLogTag, msg.ID)
}

// sendErrorResponse sends an error response
func (c *ClientConnection) sendErrorResponse(requestID string, code string, err error) {
	c.log.Info("%s: Sending error response: ID=%s, Code=%s, Error=%v", c.idLogTag, requestID, code, err)
	response := Message{
		ID:   requestID,
		Type: ErrorMessage,
		Data: ErrorData{
			Code:    code,
			Message: err.Error(),
		},
	}

	if sendErr := c.sendMessage(response); sendErr != nil {
		c.log.Warn("%s: Failed to send error response: %v", c.idLogTag, sendErr)
	}
}

// cleanupLoop periodically cleans up timed-out requests
func (c *ClientConnection) cleanupLoop() {
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
func (c *ClientConnection) cleanupTimedOutRequests() {
	now := time.Now()
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()

	for id, req := range c.pendingRequests {
		if now.After(req.Timeout) {
			c.log.Warn("%s: Cleaning up timed out request: %s", c.idLogTag, id)
			close(req.ResponseChan)
			delete(c.pendingRequests, id)
		}
	}
}
