package tradelocker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const eventTradeLockerRequest = "tradelocker_request"

type Logger interface {
	Log(ctx context.Context, sessionID *uuid.UUID, eventType, method, path string, statusCode int, message string)
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     Logger
	configOnce sync.Once
	configErr  error
	columns    []Column
}

func NewClient(baseURL string, logger Logger) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger: logger,
	}
}

func parseSessionID(sessionID *string) *uuid.UUID {
	if sessionID == nil {
		return nil
	}
	parsed, err := uuid.Parse(*sessionID)
	if err != nil {
		return nil
	}
	return &parsed
}

func (c *Client) logRequest(ctx context.Context, sessionID *string, method, path string, statusCode int, message string) {
	if c.logger == nil {
		return
	}
	c.logger.Log(ctx, parseSessionID(sessionID), eventTradeLockerRequest, method, path, statusCode, message)
}

func (c *Client) doJSON(ctx context.Context, sessionID *string, method, path, accNum string, body interface{}, bearerToken string, out interface{}) (int, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if accNum != "" {
		req.Header.Set("accNum", accNum)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logRequest(ctx, sessionID, method, path, 0, err.Error())
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	msg := ""
	if resp.StatusCode >= 400 {
		msg = string(respBody)
		if len(msg) > 500 {
			msg = msg[:500]
		}
	}
	c.logRequest(ctx, sessionID, method, path, resp.StatusCode, msg)

	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("tradelocker error %d: %s", resp.StatusCode, string(respBody))
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return resp.StatusCode, fmt.Errorf("decode response: %w", err)
		}
	}

	return resp.StatusCode, nil
}
