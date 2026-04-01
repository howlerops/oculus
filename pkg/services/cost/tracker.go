package cost

import (
	"fmt"
	"sync"
)

// Model pricing per million tokens (USD)
var modelPricing = map[string]struct {
	InputPerMillion  float64
	OutputPerMillion float64
}{
	"claude-sonnet-4-6":  {3.0, 15.0},
	"claude-opus-4-6":   {15.0, 75.0},
	"claude-haiku-4-5-20251001":  {0.80, 4.0},
}

// Tracker tracks token usage and cost across a session
type Tracker struct {
	mu           sync.Mutex
	InputTokens  int
	OutputTokens int
	CacheReads   int
	CacheWrites  int
	Model        string
}

// NewTracker creates a new cost tracker
func NewTracker(model string) *Tracker {
	return &Tracker{Model: model}
}

// Add records token usage from an API response
func (t *Tracker) Add(inputTokens, outputTokens, cacheReads, cacheWrites int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.InputTokens += inputTokens
	t.OutputTokens += outputTokens
	t.CacheReads += cacheReads
	t.CacheWrites += cacheWrites
}

// TotalCostUSD returns the estimated total cost
func (t *Tracker) TotalCostUSD() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	pricing, ok := modelPricing[t.Model]
	if !ok {
		pricing = modelPricing["claude-sonnet-4-6"]
	}

	inputCost := float64(t.InputTokens) / 1_000_000 * pricing.InputPerMillion
	outputCost := float64(t.OutputTokens) / 1_000_000 * pricing.OutputPerMillion
	return inputCost + outputCost
}

// Summary returns a formatted summary string
func (t *Tracker) Summary() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	cost := t.TotalCostUSD()
	return fmt.Sprintf(
		"Token Usage:\n  Input:  %d\n  Output: %d\n  Cache reads: %d\n  Cache writes: %d\n  Estimated cost: $%.4f",
		t.InputTokens, t.OutputTokens, t.CacheReads, t.CacheWrites, cost,
	)
}
