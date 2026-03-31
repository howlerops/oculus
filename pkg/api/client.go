package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/howlerops/oculus/pkg/types"
)

const (
	DefaultBaseURL    = "https://api.anthropic.com"
	DefaultAPIVersion = "2023-06-01"
	DefaultMaxTokens  = 16384
	MessagesEndpoint  = "/v1/messages"
)

// ClientConfig holds API client configuration
type ClientConfig struct {
	APIKey     string
	BaseURL    string
	APIVersion string
	MaxRetries int
	HTTPClient *http.Client
}

// Client is the Anthropic API client
type Client struct {
	config     ClientConfig
	httpClient *http.Client
}

// NewClient creates a new Anthropic API client
func NewClient(cfg ClientConfig) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	if cfg.APIVersion == "" {
		cfg.APIVersion = DefaultAPIVersion
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Minute,
		}
	}

	return &Client{
		config:     cfg,
		httpClient: httpClient,
	}
}

// MessageRequest is the request body for the messages API
type MessageRequest struct {
	Model         string                 `json:"model"`
	MaxTokens     int                    `json:"max_tokens"`
	Messages      []MessageParam         `json:"messages"`
	System        interface{}            `json:"system,omitempty"` // string or []SystemBlock
	Stream        bool                   `json:"stream"`
	StopSequences []string               `json:"stop_sequences,omitempty"`
	Temperature   *float64               `json:"temperature,omitempty"`
	TopP          *float64               `json:"top_p,omitempty"`
	Tools         []ToolParam            `json:"tools,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Thinking      *ThinkingConfig        `json:"thinking,omitempty"`
}

// SystemBlock for multi-part system prompts with cache control
type SystemBlock struct {
	Type         string        `json:"type"`
	Text         string        `json:"text"`
	CacheControl *CacheControl `json:"cache_control,omitempty"`
}

// CacheControl for prompt caching
type CacheControl struct {
	Type string `json:"type"` // "ephemeral"
}

// ThinkingConfig for extended thinking
type ThinkingConfig struct {
	Type         string `json:"type"`          // "enabled"
	BudgetTokens int    `json:"budget_tokens"`
}

// MessageParam is a message in the request
type MessageParam struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // string or []ContentBlockParam
}

// ContentBlockParam is a content block in the request
type ContentBlockParam struct {
	Type         string                 `json:"type"`
	Text         string                 `json:"text,omitempty"`
	ID           string                 `json:"id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Input        map[string]interface{} `json:"input,omitempty"`
	ToolUseID    string                 `json:"tool_use_id,omitempty"`
	Content      interface{}            `json:"content,omitempty"`
	IsError      bool                   `json:"is_error,omitempty"`
	CacheControl *CacheControl          `json:"cache_control,omitempty"`
}

// ToolParam defines a tool for the API
type ToolParam struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"input_schema"`
}

// MessageResponse is the non-streaming response
type MessageResponse struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Role       string               `json:"role"`
	Content    []types.ContentBlock `json:"content"`
	Model      string               `json:"model"`
	StopReason string               `json:"stop_reason"`
	Usage      types.Usage          `json:"usage"`
}

// StreamCallback receives streaming events
type StreamCallback func(event types.StreamEvent) error

// CreateMessage sends a non-streaming message request
func (c *Client) CreateMessage(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
	req.Stream = false

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+MessagesEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.doWithRetry(ctx, httpReq, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	var msgResp MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &msgResp, nil
}

// CreateMessageStream sends a streaming message request
func (c *Client) CreateMessageStream(ctx context.Context, req MessageRequest, callback StreamCallback) error {
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+MessagesEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.doWithRetry(ctx, httpReq, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseAPIError(resp)
	}

	return c.parseSSEStream(resp.Body, callback)
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.config.APIKey)
	req.Header.Set("anthropic-version", c.config.APIVersion)
	req.Header.Set("anthropic-beta", "prompt-caching-2024-07-31")
}

func (c *Client) parseSSEStream(reader io.Reader, callback StreamCallback) error {
	scanner := bufio.NewScanner(reader)
	// Increase buffer size for large events
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var eventType string
	var dataLines []string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			// Empty line = end of event
			if eventType != "" && len(dataLines) > 0 {
				data := strings.Join(dataLines, "\n")
				event, err := parseStreamEvent(eventType, data)
				if err != nil {
					return fmt.Errorf("parse event %s: %w", eventType, err)
				}
				if err := callback(event); err != nil {
					return err
				}
			}
			eventType = ""
			dataLines = nil
			continue
		}

		if strings.HasPrefix(line, "event: ") {
			eventType = strings.TrimPrefix(line, "event: ")
		} else if strings.HasPrefix(line, "data: ") {
			dataLines = append(dataLines, strings.TrimPrefix(line, "data: "))
		}
	}

	return scanner.Err()
}

func parseStreamEvent(eventType string, data string) (types.StreamEvent, error) {
	event := types.StreamEvent{
		Type: types.StreamEventType(eventType),
	}

	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return event, err
	}

	return event, nil
}

func (c *Client) doWithRetry(ctx context.Context, req *http.Request, body []byte) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := retryDelay(attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
			// Reset body for retry
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Don't retry on success or client errors (except 429, 529)
		if resp.StatusCode < 500 && resp.StatusCode != 429 && resp.StatusCode != 529 {
			return resp, nil
		}

		// On last attempt, return the error response so caller can parse it
		if attempt == c.config.MaxRetries {
			lastResp = resp
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
			break
		}

		// Retryable error - close and try again
		resp.Body.Close()
		lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Return last response if available (for error parsing)
	if lastResp != nil {
		return lastResp, parseAPIError(lastResp)
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
