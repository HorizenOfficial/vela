package communication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/horizen-pes/pkg/common"
)

// HTTPClient is an HTTP implementation of the ExecutorClient interface
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Connect establishes a connection to the executor
func (c *HTTPClient) Connect(ctx context.Context) error {
	// For HTTP, there's no connection to establish
	return nil
}

// Close closes the connection to the executor
func (c *HTTPClient) Close() error {
	// For HTTP, there's no connection to close
	c.httpClient.CloseIdleConnections()
	return nil
}

// ProcessRequest sends a request to the executor and returns the response
func (c *HTTPClient) ProcessRequest(
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
	respMsg, err := c.sendMessage(ctx, "/entrypoint", msg)
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
func (c *HTTPClient) DeployApp(ctx context.Context, req *common.Request) (*common.ApplicationState, []byte, error) {
	// Create the message
	msg := Message{
		Type: DeployAppRequestMessage,
		Data: DeployAppRequestData{
			Request: req,
		},
	}

	// Send the message
	respMsg, err := c.sendMessage(ctx, "/entrypoint", msg)
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
func (c *HTTPClient) GenerateDeanonymizationReport(
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
	respMsg, err := c.sendMessage(ctx, "/entrypoint", msg)
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
func (c *HTTPClient) sendMessage(ctx context.Context, path string, msg Message) (*Message, error) {
	// Marshal the message
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(msgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("executor returned non-OK status: %d - %s", resp.StatusCode, string(respBytes))
	}

	// Unmarshal the response
	var respMsg Message
	if err := json.Unmarshal(respBytes, &respMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &respMsg, nil
}

// extractProcessResponseData extracts ProcessResponseData from an interface{}
func extractProcessResponseData(data interface{}) (*ProcessResponseData, error) {
	// Convert to JSON and back to handle the case where data is a map[string]interface{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data: %w", err)
	}

	var respData ProcessResponseData
	if err := json.Unmarshal(jsonData, &respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	return &respData, nil
}

// extractDeployAppResponseData extracts DeployAppResponseData from an interface{}
func extractDeployAppResponseData(data interface{}) (*DeployAppResponseData, error) {
	// Convert to JSON and back to handle the case where data is a map[string]interface{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data: %w", err)
	}

	var respData DeployAppResponseData
	if err := json.Unmarshal(jsonData, &respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	return &respData, nil
}

// extractDeanonymizationResponseData extracts DeanonymizationResponseData from an interface{}
func extractDeanonymizationResponseData(data interface{}) (*DeanonymizationResponseData, error) {
	// Convert to JSON and back to handle the case where data is a map[string]interface{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data: %w", err)
	}

	var respData DeanonymizationResponseData
	if err := json.Unmarshal(jsonData, &respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	return &respData, nil
}
