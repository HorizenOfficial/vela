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
}

// Server is a unified server implementation of the ExecutorServer interface
type Server struct {
	factory      ConnectionFactory
	listener     net.Listener
	isRunning    bool
	mu           sync.Mutex
	client       *ClientConnection
	clientMu     sync.RWMutex
	handler      RequestHandler
	shutdownChan chan struct{}
	reqTimeout   time.Duration
}

// NewServer creates a new server with the specified connection factory
func NewServer(factory ConnectionFactory) *Server {
	return &Server{
		factory:      factory,
		shutdownChan: make(chan struct{}),
		reqTimeout:   30 * time.Second,
	}
}

// Start starts the server and begins listening for connections
func (s *Server) Start(ctx context.Context) error {
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

	s.listener = listener
	s.isRunning = true

	// Start accepting connections in a goroutine
	go s.acceptConnections(ctx)

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
		s.client.close()
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

// SendGetUserKeys requests user keys from the connected client
func (s *Server) SendGetUserKeys(ctx context.Context, users []string) (map[string][]byte, error) {
	s.clientMu.RLock()
	client := s.client
	s.clientMu.RUnlock()

	if client == nil {
		return nil, fmt.Errorf("no client connected")
	}

	if len(users) == 0 {
		return map[string][]byte{}, nil
	}

	// Create a request message
	msg := Message{
		ID:   generateID(),
		Type: GetUserKeysRequestMessage,
		Data: GetUserKeysRequestData{
			Users: users,
		},
	}

	// Send request and wait for response
	respMsg, err := client.sendRequestAndWaitForResponse(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get user keys: %w", err)
	}

	if respMsg.Type == ErrorMessage {
		errorData, _ := extractData[ErrorData](respMsg.Data)
		return nil, fmt.Errorf("client error: %s", errorData.Message)
	}

	if respMsg.Type != GetUserKeysResponseMessage {
		return nil, fmt.Errorf("unexpected response type: %v", respMsg.Type)
	}

	respData, err := extractData[GetUserKeysResponseData](respMsg.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to extract response data: %w", err)
	}

	return respData.UserKeys, nil
}

// acceptConnections accepts incoming connections
func (s *Server) acceptConnections(ctx context.Context) {
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
				log.Printf("Server: Error accepting connection: %v", err)
				continue
			}
		}

		// Handle new client connection
		go s.handleNewClient(ctx, conn)
	}
}

// handleNewClient handles a new client connection
func (s *Server) handleNewClient(ctx context.Context, conn net.Conn) {
	log.Printf("Server: New client connected from %s", conn.RemoteAddr())

	client := &ClientConnection{
		conn:            conn,
		reader:          bufio.NewReader(conn),
		writer:          bufio.NewWriter(conn),
		connected:       true,
		pendingRequests: make(map[string]*PendingRequest),
		shutdown:        make(chan struct{}),
		reqTimeout:      s.reqTimeout,
	}

	// Set as current client (only one client supported)
	s.clientMu.Lock()
	if s.client != nil {
		s.client.close()
	}
	s.client = client
	s.clientMu.Unlock()

	// Start client message handling
	go client.messageReaderLoop(ctx, s)
	go client.cleanupLoop()
}

// close closes the client connection
func (c *ClientConnection) close() {
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
func (c *ClientConnection) sendRequestAndWaitForResponse(ctx context.Context, msg Message) (*Message, error) {
	// Create response channel
	responseChan := make(chan *Message, 1)
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
		return nil, fmt.Errorf("client connection is shutting down")
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

	// Add newline delimiter
	data = append(data, delimiter)

	// Write a message
	if _, err := c.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if err := c.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

// messageReaderLoop continuously reads messages from the connection
func (c *ClientConnection) messageReaderLoop(ctx context.Context, server *Server) {
	defer c.close()

	for {
		select {
		case <-c.shutdown:
			return
		case <-ctx.Done():
			return
		default:
			// Continue reading
		}

		// Read message with timeout
		log.Printf("Set Read operation time out")
		c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		msgBytes, err := c.reader.ReadBytes(delimiter)
		if err != nil {
			// Check if it's a timeout error
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Read operation time out!")
				continue // Continue reading on timeout
			}
			// Connection error, exit loop
			return
		}

		// Parse message
		var msg Message
		log.Printf("Server: message %d bytes: %x...%x", len(msgBytes), msgBytes[:16], msgBytes[len(msgBytes)-16:])
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Server: Error parsing message: %v", err)
			continue
		}

		log.Printf("Server: Received message: ID=%s, Type=%v", msg.ID, msg.Type)
		// Route message in a separate goroutine
		go c.routeIncomingMessage(ctx, &msg, server)
	}
}

// routeIncomingMessage routes incoming messages to appropriate handlers
func (c *ClientConnection) routeIncomingMessage(ctx context.Context, msg *Message, server *Server) {
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
func (c *ClientConnection) handleClientRequest(ctx context.Context, msg *Message, server *Server) {
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
	case DeanonymizationRequestMessage:
		c.handleDeanonymizationRequest(ctx, msg, handler)
	default:
		c.sendErrorResponse(msg.ID, "UNKNOWN_REQUEST", fmt.Errorf("unknown request type: %v", msg.Type))
	}
}

// handleProcessRequest handles ProcessRequest messages
func (c *ClientConnection) handleProcessRequest(ctx context.Context, msg *Message, handler RequestHandler) {
	reqData, err := extractData[ProcessRequestData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	updatePayload, updatedState, err := handler.HandleProcessRequest(ctx, reqData.Request, reqData.ApplicationState, reqData.SenderKey, reqData.WasmModule)
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
		},
	}

	if err := c.sendMessage(response); err != nil {
		log.Printf("Server: Failed to send HandleProcessRequest response: %v", err)
	}
	log.Printf("Server: ProcessRequest handled successfully, ID=%s", msg.ID)
}

// handleDeployAppRequest handles DeployApp messages
func (c *ClientConnection) handleDeployAppRequest(ctx context.Context, msg *Message, handler RequestHandler) {
	reqData, err := extractData[DeployAppRequestData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	updatePayload, appState, err := handler.HandleDeployApp(ctx, reqData.Request)
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
		log.Printf("Server: Failed to send HandleDeployApp response: %v", err)
	}
	log.Printf("Server: DeployApp handled successfully, ID=%s", msg.ID)
}

// handleDeanonymizationRequest handles deanonymization messages
func (c *ClientConnection) handleDeanonymizationRequest(ctx context.Context, msg *Message, handler RequestHandler) {
	reqData, err := extractData[DeanonymizationRequestData](msg.Data)
	if err != nil {
		c.sendErrorResponse(msg.ID, "INVALID_REQUEST", err)
		return
	}

	report, err := handler.HandleGenerateDeanonymizationReport(ctx, reqData.Request, reqData.ApplicationState, reqData.SenderKey, reqData.WasmModule)
	if err != nil {
		c.sendErrorResponse(msg.ID, "HANDLER_ERROR", err)
		return
	}

	response := Message{
		ID:   msg.ID,
		Type: DeanonymizationResponseMessage,
		Data: DeanonymizationResponseData{
			Report: report,
		},
	}

	if err := c.sendMessage(response); err != nil {
		log.Printf("Server: Failed to send deanonymization response: %v", err)
	}
	log.Printf("Server: Deanonymization handled successfully, ID=%s", msg.ID)
}

// sendErrorResponse sends an error response
func (c *ClientConnection) sendErrorResponse(requestID string, code string, err error) {
	log.Printf("Server: Sending error response: ID=%s, Code=%s, Error=%v", requestID, code, err)
	response := Message{
		ID:   requestID,
		Type: ErrorMessage,
		Data: ErrorData{
			Code:    code,
			Message: err.Error(),
		},
	}

	if sendErr := c.sendMessage(response); sendErr != nil {
		log.Printf("Server: Failed to send error response: %v", sendErr)
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
			log.Printf("Server: Cleaning up timed out request: %s", id)
			close(req.ResponseChan)
			delete(c.pendingRequests, id)
		}
	}
}
