package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPHookConfig defines an HTTP-based hook
type HTTPHookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method,omitempty"` // default POST
	Headers map[string]string `json:"headers,omitempty"`
	Timeout int               `json:"timeout,omitempty"` // ms
}

// ExecuteHTTPHook sends hook data to an HTTP endpoint
func ExecuteHTTPHook(ctx context.Context, cfg HTTPHookConfig, input HookInput) (*HookOutput, error) {
	method := cfg.Method
	if method == "" {
		method = "POST"
	}

	timeout := 30 * time.Second
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Millisecond
	}

	body, _ := json.Marshal(input)

	// SSRF guard: block internal IPs
	if isInternalURL(cfg.URL) {
		return nil, fmt.Errorf("SSRF blocked: internal URL %s", cfg.URL)
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, method, cfg.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))

	var output HookOutput
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		json.Unmarshal(respBody, &output) //nolint:errcheck
	} else {
		output.Message = fmt.Sprintf("HTTP hook returned %d: %s", resp.StatusCode, string(respBody))
	}

	return &output, nil
}

func isInternalURL(url string) bool {
	internals := []string{"localhost", "127.0.0.1", "0.0.0.0", "::1", "169.254.", "10.", "172.16.", "192.168."}
	lower := strings.ToLower(url)
	for _, prefix := range internals {
		if strings.Contains(lower, prefix) {
			return true
		}
	}
	return false
}
