package api

import (
	"math"
	"math/rand"
	"time"
)

// retryDelay calculates exponential backoff with jitter
func retryDelay(attempt int) time.Duration {
	base := math.Pow(2, float64(attempt)) * 1000 // base ms
	jitter := rand.Float64() * 1000               // up to 1s jitter
	delay := time.Duration(base+jitter) * time.Millisecond

	// Cap at 60 seconds
	if delay > 60*time.Second {
		delay = 60 * time.Second
	}
	return delay
}
