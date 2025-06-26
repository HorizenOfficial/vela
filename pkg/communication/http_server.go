package communication

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// HTTPServer is an HTTP implementation of the ExecutorServer interface
type HTTPServer struct {
	addr         string
	handler      RequestHandler
	server       *http.Server
	mu           sync.Mutex
	isRunning    bool
	shutdownChan chan struct{}
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(addr string) *HTTPServer {
	return &HTTPServer{
		addr:         addr,
		shutdownChan: make(chan struct{}),
	}
}

// Start starts the server
func (s *HTTPServer) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("server is already running")
	}

	if s.handler == nil {
		return fmt.Errorf("request handler is not set")
	}

	// Create a new HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/entrypoint", s.entrypoint)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", s.addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
		close(s.shutdownChan)
	}()

	s.isRunning = true
	return nil
}

// Stop stops the server
func (s *HTTPServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*1000000000) // 5 seconds
	defer cancel()

	// Shutdown the server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	// Wait for the server to shutdown
	<-s.shutdownChan

	s.isRunning = false
	return nil
}

// SetRequestHandler sets the handler for incoming requests
func (s *HTTPServer) SetRequestHandler(handler RequestHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
}

func (s *HTTPServer) entrypoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorResponse(w, "failed to read request body", err)
		return
	}

	// Unmarshal the message
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		s.sendErrorResponse(w, "failed to unmarshal message", err)
		return
	}

	switch msg.Type {
	case ProcessRequestMessage:
		s.handleProcess(r.Context(), w, msg)
	case DeployAppRequestMessage:
		s.handleDeploy(r.Context(), w, msg)
	case DeanonymizationRequestMessage:
		s.handleDeanonymize(r.Context(), w, msg)
	default:
		s.sendErrorResponse(w, "unexpected message type", fmt.Errorf("expected ProcessRequestMessage, got %d", msg.Type))
		return
	}
}

// handleProcess handles process requests
func (s *HTTPServer) handleProcess(context context.Context, w http.ResponseWriter, msg Message) {
	// Extract the request data
	reqData, err := extractProcessRequestData(msg.Data)
	if err != nil {
		s.sendErrorResponse(w, "failed to extract request data", err)
		return
	}

	// Process the request
	updatePayload, updatedState, err := s.handler.ProcessRequest(context, reqData.Request, reqData.ApplicationState, reqData.WasmModule)
	if err != nil {
		s.sendErrorResponse(w, "failed to process request", err)
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
	s.sendResponse(w, respMsg)
}

// handleDeploy handles deploy requests
func (s *HTTPServer) handleDeploy(context context.Context, w http.ResponseWriter, msg Message) {
	// Extract the request data
	reqData, err := extractDeployAppRequestData(msg.Data)
	if err != nil {
		s.sendErrorResponse(w, "failed to extract request data", err)
		return
	}

	// Deploy the application
	applicationState, wasmModule, err := s.handler.DeployApp(context, reqData.Request)
	if err != nil {
		s.sendErrorResponse(w, "failed to deploy application", err)
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
	s.sendResponse(w, respMsg)
}

// handleDeanonymize handles deanonymization requests
func (s *HTTPServer) handleDeanonymize(context context.Context, w http.ResponseWriter, msg Message) {
	// Extract the request data
	reqData, err := extractDeanonymizationRequestData(msg.Data)
	if err != nil {
		s.sendErrorResponse(w, "failed to extract request data", err)
		return
	}

	// Generate the deanonymization report
	report, err := s.handler.GenerateDeanonymizationReport(context, reqData.Request, reqData.ApplicationState, reqData.WasmModule)
	if err != nil {
		s.sendErrorResponse(w, "failed to generate deanonymization report", err)
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
	s.sendResponse(w, respMsg)
}

// sendResponse sends a response to the client
func (s *HTTPServer) sendResponse(w http.ResponseWriter, msg Message) {
	// Marshal the message
	respBytes, err := json.Marshal(msg)
	if err != nil {
		s.sendErrorResponse(w, "failed to marshal response", err)
		return
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// sendErrorResponse sends an error response to the client
func (s *HTTPServer) sendErrorResponse(w http.ResponseWriter, code string, err error) {
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// extractProcessRequestData extracts ProcessRequestData from an interface{}
func extractProcessRequestData(data interface{}) (*ProcessRequestData, error) {
	// Convert to JSON and back to handle the case where data is a map[string]interface{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	var reqData ProcessRequestData
	if err := json.Unmarshal(jsonData, &reqData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
	}

	return &reqData, nil
}

// extractDeployAppRequestData extracts DeployAppRequestData from an interface{}
func extractDeployAppRequestData(data interface{}) (*DeployAppRequestData, error) {
	// Convert to JSON and back to handle the case where data is a map[string]interface{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	var reqData DeployAppRequestData
	if err := json.Unmarshal(jsonData, &reqData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
	}

	return &reqData, nil
}

// extractDeanonymizationRequestData extracts DeanonymizationRequestData from an interface{}
func extractDeanonymizationRequestData(data interface{}) (*DeanonymizationRequestData, error) {
	// Convert to JSON and back to handle the case where data is a map[string]interface{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	var reqData DeanonymizationRequestData
	if err := json.Unmarshal(jsonData, &reqData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
	}

	return &reqData, nil
}
