package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError represents an error from the Anthropic API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Type       string `json:"type"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Anthropic API error %d (%s): %s", e.StatusCode, e.Type, e.Message)
}

// IsRateLimited returns true for 429 errors
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == 429
}

// IsOverloaded returns true for 529 errors
func (e *APIError) IsOverloaded() bool {
	return e.StatusCode == 529
}

// IsPromptTooLong returns true when the prompt exceeds the context window
func (e *APIError) IsPromptTooLong() bool {
	return e.Type == "invalid_request_error" &&
		strings.Contains(e.Message, "prompt is too long")
}

// PromptTooLongError is returned when the prompt exceeds context limits
var PromptTooLongError = &APIError{
	StatusCode: 400,
	Type:       "invalid_request_error",
	Message:    "prompt is too long: your conversation has grown beyond the context window. Use /compact to summarize older messages.",
}

func parseAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var errResp struct {
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Type:       "unknown",
			Message:    string(body),
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Type:       errResp.Error.Type,
		Message:    errResp.Error.Message,
	}
}
