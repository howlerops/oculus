package ratelimit

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RateLimitInfo struct {
	RequestsLimit     int
	RequestsRemaining int
	RequestsReset     time.Time
	TokensLimit       int
	TokensRemaining   int
	TokensReset       time.Time
	RetryAfter        time.Duration
}

type Tracker struct {
	mu      sync.Mutex
	current *RateLimitInfo
	history []RateLimitInfo
}

func NewTracker() *Tracker { return &Tracker{} }

func (t *Tracker) ProcessHeaders(headers http.Header) *RateLimitInfo {
	info := &RateLimitInfo{}
	info.RequestsLimit, _ = strconv.Atoi(headers.Get("anthropic-ratelimit-requests-limit"))
	info.RequestsRemaining, _ = strconv.Atoi(headers.Get("anthropic-ratelimit-requests-remaining"))
	if reset := headers.Get("anthropic-ratelimit-requests-reset"); reset != "" {
		info.RequestsReset, _ = time.Parse(time.RFC3339, reset)
	}
	info.TokensLimit, _ = strconv.Atoi(headers.Get("anthropic-ratelimit-tokens-limit"))
	info.TokensRemaining, _ = strconv.Atoi(headers.Get("anthropic-ratelimit-tokens-remaining"))
	if reset := headers.Get("anthropic-ratelimit-tokens-reset"); reset != "" {
		info.TokensReset, _ = time.Parse(time.RFC3339, reset)
	}
	if ra := headers.Get("retry-after"); ra != "" {
		if secs, err := strconv.Atoi(ra); err == nil {
			info.RetryAfter = time.Duration(secs) * time.Second
		}
	}
	t.mu.Lock()
	t.current = info
	t.history = append(t.history, *info)
	t.mu.Unlock()
	return info
}

func (t *Tracker) GetCurrent() *RateLimitInfo {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.current
}

func (t *Tracker) IsNearLimit() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.current == nil {
		return false
	}
	return t.current.RequestsRemaining < 5 || t.current.TokensRemaining < 1000
}

func (t *Tracker) GetWarningMessage() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.current == nil {
		return ""
	}
	if t.current.RequestsRemaining < 5 {
		return fmt.Sprintf("Rate limit warning: %d requests remaining (resets %s)",
			t.current.RequestsRemaining, t.current.RequestsReset.Format("15:04:05"))
	}
	return ""
}
