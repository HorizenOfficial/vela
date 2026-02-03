package subgraph

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
)

type graphError struct {
	Message string `json:"message"`
}

type graphResponse[T any] struct {
	Data   T            `json:"data"`
	Errors []graphError `json:"errors"`
}

type client struct {
	endpoint   string
	httpClient *http.Client
	log        logger.Logger
}

// NewClient builds a subgraph client pointing to the given endpoint.
func NewClient(endpoint string) Client {
	return NewClientWithLogger(endpoint, nil)
}

// NewClientWithLogger builds a subgraph client pointing to the given endpoint and using the provided logger.
func NewClientWithLogger(endpoint string, log logger.Logger) Client {
	return &client{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		log: log,
	}
}

func (c *client) Logger() logger.Logger {
	return c.log
}

type healthCheckResponse struct {
	Meta *struct {
		HasIndexingErrors bool `json:"hasIndexingErrors"`
	} `json:"_meta"`
}

type requestCompletedResponse struct {
	RequestCompleteds []struct {
		RequestID       string      `json:"requestId"`
		Status          json.Number `json:"status"`
		ErrorCode       json.Number `json:"errorCode"`
		ErrorMessage    string      `json:"errorMessage"`
		ApplicationFees json.Number `json:"applicationFees"`
		BlockNumber     json.Number `json:"blockNumber"`
	} `json:"requestCompleteds"`
}

func (c *client) GetRequestCompletedByID(ctx context.Context, requestID common.RequestIdType) (*RequestCompleted, error) {
	query := `
query($requestId: Bytes!) {
  requestCompleteds(where: { requestId: $requestId }, first: 1) {
    requestId
    status
    errorCode
    errorMessage
    applicationFees
    blockNumber
  }
}`

	var resp graphResponse[requestCompletedResponse]
	if err := c.doGraphQL(ctx, query, map[string]interface{}{"requestId": "0x" + requestID.String()}, &resp); err != nil {
		return nil, err
	}
	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("subgraph returned errors: %v", resp.Errors[0].Message)
	}
	if len(resp.Data.RequestCompleteds) == 0 {
		return nil, nil
	}

	entity := resp.Data.RequestCompleteds[0]

	statusUint, err := strconv.ParseUint(entity.Status.String(), 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid status value %q: %w", entity.Status, err)
	}
	status, err := common.UInt8ToRequestResultStatus(uint8(statusUint))
	if err != nil {
		return nil, err
	}

	errorCodeUint, err := strconv.ParseUint(entity.ErrorCode.String(), 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid errorCode %q: %w", entity.ErrorCode, err)
	}

	appFees, ok := common.StringToBigInt(entity.ApplicationFees.String())
	if !ok {
		return nil, fmt.Errorf("invalid applicationFees %q", entity.ApplicationFees)
	}

	blockNumber, err := strconv.ParseUint(entity.BlockNumber.String(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid blockNumber %q: %w", entity.BlockNumber, err)
	}

	return &RequestCompleted{
		RequestID:       requestID,
		Status:          status,
		ErrorCode:       uint8(errorCodeUint),
		ErrorMessage:    entity.ErrorMessage,
		ApplicationFees: appFees,
		BlockNumber:     blockNumber,
	}, nil
}

func (c *client) HealthCheck(ctx context.Context) error {
	query := `
query HealthCheck {
  _meta {
    hasIndexingErrors
  }
}`

	var resp graphResponse[healthCheckResponse]
	if err := c.doGraphQL(ctx, query, map[string]interface{}{}, &resp); err != nil {
		return err
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("subgraph returned errors: %v", resp.Errors[0].Message)
	}
	if resp.Data.Meta == nil {
		return fmt.Errorf("subgraph health check returned empty meta")
	}
	if resp.Data.Meta.HasIndexingErrors {
		return fmt.Errorf("subgraph reports indexing errors")
	}
	return nil
}

type userEventsResponse struct {
	UserEvents []struct {
		ApplicationID string `json:"applicationId"`
		RequestID     string `json:"requestId"`
		EventSubType  string `json:"eventSubType"`
		EncryptedData string `json:"encryptedData"`
		BlockNumber   string `json:"blockNumber"`
		LogIndex      string `json:"logIndex"`
		SortKey       string `json:"sortKey"`
	} `json:"userEvents"`
}

func (c *client) GetUserEvents(ctx context.Context, applicationID common.ApplicationIdType, eventSubType string, limit int, before *big.Int) ([]UserEvent, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 1000 {
		limit = 1000
	}

	variables := map[string]interface{}{
		"applicationId": fmt.Sprintf("%d", uint64(applicationID)),
		"limit":         limit,
	}

	varDefs := ""
	whereParts := []string{"applicationId: $applicationId"}
	if strings.TrimSpace(eventSubType) != "" {
		varDefs += ", $eventSubType: Bytes!"
		variables["eventSubType"] = eventSubType
		whereParts = append(whereParts, "eventSubType: $eventSubType")
	}
	if before != nil {
		varDefs += ", $before: BigInt!"
		variables["before"] = before.String()
		whereParts = append(whereParts, "sortKey_lt: $before")
	}

	query := fmt.Sprintf(`
query($applicationId: BigInt!, $limit: Int!%s) {
  userEvents(
    where: { %s }
    orderBy: sortKey
    orderDirection: desc
    first: $limit
  ) {
    applicationId
    requestId
    eventSubType
    encryptedData
    blockNumber
    logIndex
    sortKey
  }
}`, varDefs, strings.Join(whereParts, ", "))

	var resp graphResponse[userEventsResponse]
	if err := c.doGraphQL(ctx, query, variables, &resp); err != nil {
		return nil, err
	}
	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("subgraph returned errors: %v", resp.Errors[0].Message)
	}

	events := make([]UserEvent, 0, len(resp.Data.UserEvents))
	for _, entity := range resp.Data.UserEvents {
		reqID, err := parseRequestID(entity.RequestID)
		if err != nil {
			return nil, fmt.Errorf("invalid requestId %q: %w", entity.RequestID, err)
		}

		data, err := decodeHex(entity.EncryptedData)
		if err != nil {
			return nil, fmt.Errorf("invalid encryptedData for request %s: %w", reqID.String(), err)
		}

		blockNumber, err := strconv.ParseUint(entity.BlockNumber, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid blockNumber %q: %w", entity.BlockNumber, err)
		}

		logIndex, err := strconv.ParseUint(entity.LogIndex, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid logIndex %q: %w", entity.LogIndex, err)
		}

		sortKey, ok := common.StringToBigInt(entity.SortKey)
		if !ok {
			return nil, fmt.Errorf("invalid sortKey %q", entity.SortKey)
		}

		events = append(events, UserEvent{
			ApplicationID: applicationID,
			RequestID:     reqID,
			EventSubType:  entity.EventSubType,
			EncryptedData: data,
			BlockNumber:   blockNumber,
			LogIndex:      logIndex,
			SortKey:       sortKey,
		})
	}

	return events, nil
}

func (c *client) doGraphQL(ctx context.Context, query string, variables map[string]interface{}, dest any) error {
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal graphql payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("graphql request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		blob, _ := io.ReadAll(res.Body)
		return fmt.Errorf("graphql request returned status %d: %s", res.StatusCode, string(blob))
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("failed to decode graphql response: %w", err)
	}

	return nil
}

func parseRequestID(hexString string) (common.RequestIdType, error) {
	trimmed := strings.TrimPrefix(hexString, "0x")
	b, err := common.RequestIdStringTo32Byte(trimmed)
	if err != nil {
		return common.RequestIdType{}, err
	}
	return common.RequestIdType(b), nil
}

func decodeHex(hexString string) ([]byte, error) {
	trimmed := strings.TrimPrefix(hexString, "0x")
	if trimmed == "" {
		return nil, nil
	}
	return hex.DecodeString(trimmed)
}
