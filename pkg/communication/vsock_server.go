package communication

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/mdlayher/vsock"
)

// VSockServer is a v-socket implementation of the ExecutorServer interface
type VSockServer struct {
	port         uint32
	handler      RequestHandler
	listener     *vsock.Listener
	mu           sync.Mutex
	isRunning    bool
	shutdownChan chan struct{}
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewVSockServer creates a new v-socket server
func NewVSockServer(port uint32) *VSockServer {
	return &VSockServer{
		port:         port,
		shutdownChan: make(chan struct{}),
	}
}

// Start starts the server
func (s *VSockServer) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("server is already running")
	}

	if s.handler == nil {
		return fmt.Errorf("request handler is not set")
	}

	// Create a new context with cancellation
	s.ctx, s.cancel = context.WithCancel(ctx)

	// Create a new v-socket listener
	listener, err := vsock.Listen(s.port, &vsock.Config{})
	if err != nil {
		return fmt.Errorf("failed to create v-socket listener: %w", err)
	}
	s.listener = listener

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting v-socket server on port %d", s.port)
		s.serve()
		close(s.shutdownChan)
	}()

	s.isRunning = true
	return nil
}

// Stop stops the server
func (s *VSockServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	// Cancel the context
	s.cancel()

	// Close the listener
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close v-socket listener: %w", err)
	}

	// Wait for the server to shutdown
	<-s.shutdownChan

	s.isRunning = false
	return nil
}

// SetRequestHandler sets the handler for incoming requests
func (s *VSockServer) SetRequestHandler(handler RequestHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
}

// serve handles incoming connections
func (s *VSockServer) serve() {
	for {
		// Accept a new connection
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				// Server is shutting down
				return
			default:
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
		}

		// Handle the connection in a goroutine
		go s.handleConnection(conn)
	}
}

// handleConnection handles a single connection
func (s *VSockServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		select {
		case <-s.ctx.Done():
			// Server is shutting down
			return
		default:
			// Read the message length
			msgLen, err := readUint32(conn)
			if err != nil {
				if err != io.EOF {
					log.Printf("Failed to read message length: %v", err)
				}
				return
			}

			// Read the message
			msgBytes := make([]byte, msgLen)
			if _, err := io.ReadFull(conn, msgBytes); err != nil {
				log.Printf("Failed to read message: %v", err)
				return
			}

			// Unmarshal the message
			var msg Message
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				s.sendErrorResponse(conn, "failed to unmarshal message", err)
				continue
			}

			// Handle the message
			respMsg, err := s.handleMessage(msg)
			if err != nil {
				log.Printf("Failed to handle message: %v", err)
				s.sendErrorResponse(conn, "failed to handle message", err)
				continue
			}

			// Marshal the response
			respBytes, err := json.Marshal(respMsg)
			if err != nil {
				log.Printf("Failed to marshal response: %v", err)
				s.sendErrorResponse(conn, "failed to marshal response", err)
				continue
			}

			// Write the response length
			respLen := uint32(len(respBytes))
			if err := writeUint32(conn, respLen); err != nil {
				log.Printf("Failed to write response length: %v", err)
				return
			}

			// Write the response
			if _, err := conn.Write(respBytes); err != nil {
				log.Printf("Failed to write response: %v", err)
				return
			}
		}
	}
}

// handleMessage handles a message and returns a response
func (s *VSockServer) handleMessage(msg Message) (*Message, error) {
	switch msg.Type {
	case ProcessRequestMessage:
		return s.handleProcessRequest(msg)
	case DeployAppRequestMessage:
		return s.handleDeployApp(msg)
	case DeanonymizationRequestMessage:
		return s.handleDeanonymization(msg)
	default:
		return nil, fmt.Errorf("unexpected message type: %d", msg.Type)
	}
}

// handleProcessRequest handles a process request message
func (s *VSockServer) handleProcessRequest(msg Message) (*Message, error) {
	// Extract the request data
	reqData, err := extractProcessRequestData(msg.Data)
	if err != nil {
		return nil, err
	}

	// Process the request
	updatePayload, applicationState, err := s.handler.ProcessRequest(s.ctx, reqData.Request, reqData.ApplicationState, reqData.WasmModule)
	if err != nil {
		return nil, err
	}

	// Create the response message
	respMsg := &Message{
		Type: ProcessResponseMessage,
		Data: ProcessResponseData{
			UpdatePayload:           updatePayload,
			UpdatedApplicationState: applicationState,
		},
	}

	return respMsg, nil
}

// handleDeployApp handles a deploy app request message
func (s *VSockServer) handleDeployApp(msg Message) (*Message, error) {
	// Extract the request data
	reqData, err := extractDeployAppRequestData(msg.Data)
	if err != nil {
		return nil, err
	}

	// Deploy the application
	appState, wasmModule, err := s.handler.DeployApp(s.ctx, reqData.Request)
	if err != nil {
		return nil, err
	}

	// Create the response message
	respMsg := &Message{
		Type: DeployAppResponseMessage,
		Data: DeployAppResponseData{
			ApplicationState: appState,
			WasmModule:       wasmModule,
		},
	}

	return respMsg, nil
}

// handleDeanonymization handles a deanonymization request message
func (s *VSockServer) handleDeanonymization(msg Message) (*Message, error) {
	// Extract the request data
	reqData, err := extractDeanonymizationRequestData(msg.Data)
	if err != nil {
		return nil, err
	}

	// Generate the deanonymization report
	report, err := s.handler.GenerateDeanonymizationReport(s.ctx, reqData.Request, reqData.ApplicationState, reqData.WasmModule)
	if err != nil {
		return nil, err
	}

	// Create the response message
	respMsg := &Message{
		Type: DeanonymizationResponseMessage,
		Data: DeanonymizationResponseData{
			Report: report,
		},
	}

	return respMsg, nil
}

// sendErrorResponse sends an error response to the client
func (s *VSockServer) sendErrorResponse(conn net.Conn, code string, err error) {
	// Create the error message
	errMsg := Message{
		Type: ErrorMessage,
		Data: ErrorData{
			Code:    code,
			Message: err.Error(),
		},
	}

	// Marshal the message
	respBytes, err := json.Marshal(errMsg)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		return
	}

	// Write the response length
	respLen := uint32(len(respBytes))
	if err := writeUint32(conn, respLen); err != nil {
		log.Printf("Failed to write error response length: %v", err)
		return
	}

	// Write the response
	if _, err := conn.Write(respBytes); err != nil {
		log.Printf("Failed to write error response: %v", err)
		return
	}
}
