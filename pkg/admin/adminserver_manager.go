package admin

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/communication"
	"github.com/horizen-pes/pkg/logger"
)

// ManagerAdminCommandServer is the interface for the manager admin command server
type ManagerAdminCommandServer interface {
	Start(ctx context.Context, identityLogTag string) error
	Stop() error
	SetCmdHandler(handler ManagerCmdHandler)
}

// ManagerAdminClientConnection represents a client connection to the manager admin server
type ManagerAdminClientConnection struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	timeout  time.Time
	idLogTag string
	log      logger.Logger
}

// ManagerAdminServer is an implementation of the ManagerAdminCommandServer interface
type ManagerAdminServer struct {
	factory       communication.ConnectionFactory
	mu            sync.Mutex
	isRunning     bool
	listener      net.Listener
	clientMu      sync.Mutex
	client        *ManagerAdminClientConnection
	handler       ManagerCmdHandler
	shutdownChan  chan struct{}
	clientTimeout time.Duration
	log           logger.Logger
}

// NewManagerAdminServer creates a new manager admin server with the specified connection factory
func NewManagerAdminServer(factory communication.ConnectionFactory, communicationParams common.CommunicationParams, log logger.Logger) *ManagerAdminServer {
	return &ManagerAdminServer{
		factory:       factory,
		shutdownChan:  make(chan struct{}),
		clientTimeout: communicationParams.RequestTimeoutSec * time.Second,
		log:           log,
	}
}

// Start starts the manager admin server and begins listening for connections
func (s *ManagerAdminServer) Start(ctx context.Context, idLogTag string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("manager admin server already running")
	}

	// Create listener using factory
	listener, err := s.factory.CreateServerListener()
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	s.listener = listener
	s.isRunning = true

	// Start accepting connections in a goroutine
	go s.acceptConnections(ctx, idLogTag)

	return nil
}

// Stop stops the manager admin server
func (s *ManagerAdminServer) Stop() error {
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

	// Forcefully close any active client connection to unblock handlers
	s.clientMu.Lock()
	if s.client != nil {
		s.client.conn.Close()
	}
	s.clientMu.Unlock()

	return nil
}

// SetCmdHandler sets the handler for client requests
func (s *ManagerAdminServer) SetCmdHandler(handler ManagerCmdHandler) {
	s.handler = handler
}

// acceptConnections accepts incoming connections
func (s *ManagerAdminServer) acceptConnections(ctx context.Context, idLogTag string) {
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
			s.log.Warn("%s: Error accepting connection: %v", idLogTag, err)
			continue
		}

		// Handle new client connection
		go s.handleNewClient(ctx, conn, idLogTag)
	}
}

func (s *ManagerAdminServer) setupNewClient(client *ManagerAdminClientConnection) bool {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	if s.client != nil && time.Now().Before(s.client.timeout) {
		return false
	}

	s.client = client
	s.client.timeout = time.Now().Add(s.clientTimeout)

	return true
}

// handleNewClient handles a new client connection
func (s *ManagerAdminServer) handleNewClient(ctx context.Context, conn net.Conn, idLogTag string) {
	s.log.Info("%s: New client connected from %s", idLogTag, conn.RemoteAddr())
	client := &ManagerAdminClientConnection{
		conn:     conn,
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
		idLogTag: idLogTag,
		log:      s.log,
	}

	defer client.Close()

	if ok := s.setupNewClient(client); !ok {
		client.sendErrorResponse("INVALID_REQUEST", fmt.Errorf("server is busy"))
		return
	}

	defer func() {
		s.clientMu.Lock()
		defer s.clientMu.Unlock()
		if s.client == client {
			s.client = nil
		}
	}()

	select {
	case <-s.shutdownChan:
		return
	case <-ctx.Done():
		return
	default:
	}

	client.handleManagerCommand(ctx, s.handler)
}

// Close closes the client connection
func (c *ManagerAdminClientConnection) Close() {
	c.log.Info("%s: Closing client connection", c.idLogTag)
	if c.conn != nil {
		c.conn.Close()
	}
}

// sendMessage sends a message to the client
func (c *ManagerAdminClientConnection) sendMessage(msg AdminMessage) error {
	// Serialize message
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	c.log.Debug("%s: MsgBytes length before delimiter: %d", c.idLogTag, len(data))

	// Add newline delimiter
	data = append(data, communication.MsgDelimiter)
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

// handleManagerCommand handles commands initiated by the client
func (c *ManagerAdminClientConnection) handleManagerCommand(ctx context.Context, handler ManagerCmdHandler) {
	if handler == nil {
		c.sendErrorResponse("INTERNAL_ERROR", fmt.Errorf("no request handler set"))
		return
	}

	msgBytes, connectionClosedErr := communication.ReadMessageFromSocket(c.conn, c.reader, c.idLogTag, c.log)
	if connectionClosedErr != nil {
		return
	}

	// Parse and route the complete message
	var msg AdminMessage
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		c.log.Error("%s: Error parsing message: %v", c.idLogTag, err)
		c.sendErrorResponse("INVALID_REQUEST", fmt.Errorf("invalid message format: %w", err))
		return
	}

	c.log.Info("%s: Received message: Type=%v", c.idLogTag, msg.Type)

	switch msg.Type {
	case GetVersionRequestMessage:
		c.handleGetVersionRequest(ctx, handler)

	default:
		c.sendErrorResponse("UNKNOWN_REQUEST", fmt.Errorf("unknown request type: %v", msg.Type))
	}
}

// handleGetVersionRequest handles GetVersion messages
func (c *ManagerAdminClientConnection) handleGetVersionRequest(ctx context.Context, handler ManagerCmdHandler) {
	version, err := handler.GetVersion(ctx)
	if err != nil {
		c.sendErrorResponse("ERROR_GETTING_VERSION", err)
		return
	}

	response := AdminMessage{
		Type: AdminResponseMessage,
		Data: version,
	}

	if err := c.sendMessage(response); err != nil {
		c.log.Warn("%s: Failed to send GetVersion response: %v", c.idLogTag, err)
		return
	}
	c.log.Info("%s: GetVersion handled successfully", c.idLogTag)
}

// sendErrorResponse sends an error response
func (c *ManagerAdminClientConnection) sendErrorResponse(code string, err error) {
	c.log.Info("%s: Sending error response: Code=%s, Error=%v", c.idLogTag, code, err)
	response := AdminMessage{
		Type: AdminErrorMessage,
		Data: communication.ErrorData{
			Code:    code,
			Message: err.Error(),
		},
	}

	if sendErr := c.sendMessage(response); sendErr != nil {
		c.log.Warn("%s: Failed to send error response: %v", c.idLogTag, sendErr)
	}
}
