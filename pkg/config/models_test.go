package config

import "testing"

func TestResolveModel(t *testing.T) {
	tests := []struct {
		input string
		want  string
		found bool
	}{
		{"opus", "claude-opus-4-20250514", true},
		{"sonnet", "claude-sonnet-4-20250514", true},
		{"haiku", "claude-haiku-4-20250506", true},
		{"o", "claude-opus-4-20250514", true},
		{"claude-sonnet-4-20250514", "claude-sonnet-4-20250514", true},
		{"unknown-model", "", false},
	}
	for _, tt := range tests {
		info, found := ResolveModel(tt.input)
		if found != tt.found {
			t.Errorf("ResolveModel(%q) found=%v, want %v", tt.input, found, tt.found)
		}
		if found && info.ID != tt.want {
			t.Errorf("ResolveModel(%q) = %s, want %s", tt.input, info.ID, tt.want)
		}
	}
}

func TestEstimateCost(t *testing.T) {
	cost := EstimateCost("claude-sonnet-4-20250514", 1000000, 100000)
	// 1M input * $3/M + 100k output * $15/M = $3 + $1.5 = $4.5
	if cost < 4.0 || cost > 5.0 {
		t.Errorf("EstimateCost = %f, want ~4.5", cost)
	}
}

func TestSupportsFeature(t *testing.T) {
	if !SupportsFeature("claude-opus-4-20250514", "thinking") {
		t.Error("Opus should support thinking")
	}
	if SupportsFeature("claude-haiku-4-20250506", "thinking") {
		t.Error("Haiku should not support thinking")
	}
}
