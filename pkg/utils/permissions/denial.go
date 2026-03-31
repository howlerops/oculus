package permissions

import (
	"sync"
	"time"
)

// DenialTracker tracks permission denials for auto-escalation
type DenialTracker struct {
	mu        sync.Mutex
	denials   map[string][]time.Time
	threshold int
	window    time.Duration
}

// NewDenialTracker creates a DenialTracker with the given threshold and window
func NewDenialTracker(threshold int, window time.Duration) *DenialTracker {
	return &DenialTracker{
		denials:   make(map[string][]time.Time),
		threshold: threshold,
		window:    window,
	}
}

// RecordDenial records a denial for a tool
func (d *DenialTracker) RecordDenial(toolName string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.denials[toolName] = append(d.denials[toolName], time.Now())
}

// ShouldEscalate checks if denials exceed threshold within window
func (d *DenialTracker) ShouldEscalate(toolName string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	cutoff := time.Now().Add(-d.window)
	var recent []time.Time
	for _, t := range d.denials[toolName] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	d.denials[toolName] = recent
	return len(recent) >= d.threshold
}

// GetDenialCount returns recent denial count for a tool
func (d *DenialTracker) GetDenialCount(toolName string) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	cutoff := time.Now().Add(-d.window)
	count := 0
	for _, t := range d.denials[toolName] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all denials
func (d *DenialTracker) Reset() {
	d.mu.Lock()
	d.denials = make(map[string][]time.Time)
	d.mu.Unlock()
}
