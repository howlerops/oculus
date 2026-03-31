package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/howlerops/oculus/pkg/types"
)

func TestCreateMessageStream(t *testing.T) {
	// Mock SSE server that returns a simple text response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("expected API key 'test-key', got %q", r.Header.Get("X-API-Key"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got %q", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)

		events := []string{
			`event: message_start` + "\n" + `data: {"type":"message_start","message":{"role":"assistant","content":[]}}` + "\n\n",
			`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n",
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"4"}}` + "\n\n",
			`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n",
			`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":1}}` + "\n\n",
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}` + "\n\n",
		}

		flusher := w.(http.Flusher)
		for _, event := range events {
			fmt.Fprint(w, event)
			flusher.Flush()
		}
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	var collectedText strings.Builder
	var finalStopReason types.StopReason

	err := client.CreateMessageStream(context.Background(), MessageRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Messages: []MessageParam{
			{Role: "user", Content: "What is 2+2?"},
		},
	}, func(event types.StreamEvent) error {
		if event.Delta != nil {
			if text, ok := event.Delta["text"].(string); ok {
				collectedText.WriteString(text)
			}
			if sr, ok := event.Delta["stop_reason"].(string); ok && sr != "" {
				finalStopReason = types.StopReason(sr)
			}
		}
		if event.StopReason != "" {
			finalStopReason = event.StopReason
		}
		return nil
	})

	if err != nil {
		t.Fatalf("CreateMessageStream failed: %v", err)
	}

	if collectedText.String() != "4" {
		t.Errorf("expected text '4', got %q", collectedText.String())
	}

	if finalStopReason != types.StopReasonEndTurn {
		t.Errorf("expected stop reason 'end_turn', got %q", finalStopReason)
	}
}

func TestCreateMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"msg_123","type":"message","role":"assistant","content":[{"type":"text","text":"4"}],"model":"claude-sonnet-4-20250514","stop_reason":"end_turn","usage":{"input_tokens":10,"output_tokens":1}}`)
	}))
	defer server.Close()

	client := NewClient(ClientConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	resp, err := client.CreateMessage(context.Background(), MessageRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 100,
		Messages: []MessageParam{
			{Role: "user", Content: "What is 2+2?"},
		},
	})

	if err != nil {
		t.Fatalf("CreateMessage failed: %v", err)
	}

	if resp.StopReason != "end_turn" {
		t.Errorf("expected stop_reason 'end_turn', got %q", resp.StopReason)
	}

	if len(resp.Content) != 1 || resp.Content[0].Text != "4" {
		t.Errorf("expected content '4', got %+v", resp.Content)
	}
}

func TestAPIErrorParsing(t *testing.T) {
	// Test the error parsing directly without going through retry
	apiErr := &APIError{
		StatusCode: 429,
		Type:       "rate_limit_error",
		Message:    "Rate limited",
	}

	if !apiErr.IsRateLimited() {
		t.Error("expected IsRateLimited=true for 429")
	}
	if apiErr.IsOverloaded() {
		t.Error("expected IsOverloaded=false for 429")
	}
	if apiErr.IsPromptTooLong() {
		t.Error("expected IsPromptTooLong=false for rate limit")
	}

	promptErr := &APIError{
		StatusCode: 400,
		Type:       "invalid_request_error",
		Message:    "prompt is too long: your conversation has grown beyond the context window",
	}
	if !promptErr.IsPromptTooLong() {
		t.Error("expected IsPromptTooLong=true")
	}

	overloadErr := &APIError{StatusCode: 529, Type: "overloaded_error", Message: "Overloaded"}
	if !overloadErr.IsOverloaded() {
		t.Error("expected IsOverloaded=true for 529")
	}

	// Test error string formatting
	errStr := apiErr.Error()
	if !strings.Contains(errStr, "429") || !strings.Contains(errStr, "Rate limited") {
		t.Errorf("unexpected error string: %s", errStr)
	}
}
