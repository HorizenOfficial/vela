package communication

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPServer is a TCP implementation of the ExecutorServer interface
type TCPServer struct {
	addr         string
	handler      RequestHandler
	listener     net.Listener
	mu           sync.Mutex
	isRunning    bool
	shutdownChan chan struct{}
	clients      map[net.Conn]struct{}
	clientsMu    sync.Mutex
}

// NewTCPServer creates a new TCP server
func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		addr:         addr,
		shutdownChan: make(chan struct{}),
		clients:      make(map[net.Conn]struct{}),
	}
}

// Start starts the server
func (s *TCPServer) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("server is already running")
	}

	if s.handler == nil {
		return fmt.Errorf("request handler is not set")
	}

	// Create a TCP listener
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to create TCP listener: %w", err)
	}
	s.listener = listener

	// Start accepting connections in a goroutine
	go func() {
		log.Printf("Starting TCP server on %s", s.addr)
		for {
			conn, err := listener.Accept()
			if err != nil {
				// Check if the server is shutting down
				select {
				case <-s.shutdownChan:
					return
				default:
					log.Printf("Error accepting connection: %v", err)
					continue
				}
			}

			// Add the client to the clients map
			s.clientsMu.Lock()
			s.clients[conn] = struct{}{}
			s.clientsMu.Unlock()

			// Handle the connection in a goroutine
			go s.handleConnection(ctx, conn)
		}
	}()

	s.isRunning = true
	return nil
}

// Stop stops the server
func (s *TCPServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	// Close the listener
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	// Signal that the server is shutting down
	close(s.shutdownChan)

	// Close all client connections
	s.clientsMu.Lock()
	for conn := range s.clients {
		conn.Close()
	}
	s.clients = make(map[net.Conn]struct{})
	s.clientsMu.Unlock()

	s.isRunning = false
	return nil
}

// SetRequestHandler sets the handler for incoming requests
func (s *TCPServer) SetRequestHandler(handler RequestHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
}

// handleConnection handles a client connection
func (s *TCPServer) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		// Remove the client from the clients map
		s.clientsMu.Lock()
		delete(s.clients, conn)
		s.clientsMu.Unlock()

		// Close the connection
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// Check if the server is shutting down
		select {
		case <-s.shutdownChan:
			return
		case <-ctx.Done():
			return
		default:
			// Continue processing
		}

		// Read a message from the client
		msgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		// Unmarshal the message
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			s.sendErrorResponse(writer, "failed to unmarshal message", err)
			continue
		}

		// Handle the message based on its type
		switch msg.Type {
		case ProcessRequestMessage:
			s.handleProcess(ctx, writer, msg)
		case DeployAppRequestMessage:
			s.handleDeploy(ctx, writer, msg)
		case DeanonymizationRequestMessage:
			s.handleDeanonymize(ctx, writer, msg)
		default:
			s.sendErrorResponse(writer, "unexpected message type", fmt.Errorf("expected ProcessRequestMessage, got %d", msg.Type))
		}
	}
}

// handleProcess handles process requests
func (s *TCPServer) handleProcess(ctx context.Context, writer *bufio.Writer, msg Message) {
	// Extract the request data
	reqData, err := extractProcessRequestData(msg.Data)
	if err != nil {
		s.sendErrorResponse(writer, "failed to extract request data", err)
		return
	}

	// Process the request
	updatePayload, updatedState, err := s.handler.ProcessRequest(ctx, reqData.Request, reqData.ApplicationState, reqData.WasmModule)
	if err != nil {
		s.sendErrorResponse(writer, "failed to process request", err)
		return
	}

	// Create the response message
	respMsg := Message{
		Type: ProcessResponseMessage,
		Data: ProcessResponseData{
			UpdatePayload:           updatePayload,
			UpdatedApplicationState: updatedState,
		},
	}

	// Send the response
	s.sendResponse(writer, respMsg)
}

// handleDeploy handles deploy requests
func (s *TCPServer) handleDeploy(ctx context.Context, writer *bufio.Writer, msg Message) {
	// Extract the request data
	reqData, err := extractDeployAppRequestData(msg.Data)
	if err != nil {
		s.sendErrorResponse(writer, "failed to extract request data", err)
		return
	}

	// Deploy the application
	applicationState, wasmModule, err := s.handler.DeployApp(ctx, reqData.Request)
	if err != nil {
		s.sendErrorResponse(writer, "failed to deploy application", err)
		return
	}

	// Create the response message
	respMsg := Message{
		Type: DeployAppResponseMessage,
		Data: DeployAppResponseData{
			ApplicationState: applicationState,
			WasmModule:       wasmModule,
		},
	}

	// Send the response
	s.sendResponse(writer, respMsg)
}

// handleDeanonymize handles deanonymization requests
func (s *TCPServer) handleDeanonymize(ctx context.Context, writer *bufio.Writer, msg Message) {
	// Extract the request data
	reqData, err := extractDeanonymizationRequestData(msg.Data)
	if err != nil {
		s.sendErrorResponse(writer, "failed to extract request data", err)
		return
	}

	// Generate the deanonymization report
	report, err := s.handler.GenerateDeanonymizationReport(ctx, reqData.Request, reqData.ApplicationState, reqData.WasmModule)
	if err != nil {
		s.sendErrorResponse(writer, "failed to generate deanonymization report", err)
		return
	}

	// Create the response message
	respMsg := Message{
		Type: DeanonymizationResponseMessage,
		Data: DeanonymizationResponseData{
			Report: report,
		},
	}

	// Send the response
	s.sendResponse(writer, respMsg)
}

// sendResponse sends a response to the client
func (s *TCPServer) sendResponse(writer *bufio.Writer, msg Message) {
	// Marshal the message
	respBytes, err := json.Marshal(msg)
	if err != nil {
		s.sendErrorResponse(writer, "failed to marshal response", err)
		return
	}

	// Add a newline to separate messages
	respBytes = append(respBytes, '\n')

	// Write the response
	_, err = writer.Write(respBytes)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Printf("Error flushing response: %v", err)
		return
	}
}

// sendErrorResponse sends an error response to the client
func (s *TCPServer) sendErrorResponse(writer *bufio.Writer, code string, err error) {
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

	// Add a newline to separate messages
	respBytes = append(respBytes, '\n')

	// Write the response
	_, err = writer.Write(respBytes)
	if err != nil {
		log.Printf("Error writing error response: %v", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Printf("Error flushing error response: %v", err)
		return
	}
}
