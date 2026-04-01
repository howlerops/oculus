package config

import "testing"

func TestResolveModel(t *testing.T) {
	tests := []struct {
		input string
		want  string
		found bool
	}{
		{"opus", "claude-opus-4-6", true},
		{"sonnet", "claude-sonnet-4-6", true},
		{"haiku", "claude-haiku-4-5-20251001", true},
		{"o", "claude-opus-4-6", true},
		{"claude-sonnet-4-6", "claude-sonnet-4-6", true},
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
	cost := EstimateCost("claude-sonnet-4-6", 1000000, 100000)
	// 1M input * $3/M + 100k output * $15/M = $3 + $1.5 = $4.5
	if cost < 4.0 || cost > 5.0 {
		t.Errorf("EstimateCost = %f, want ~4.5", cost)
	}
}

func TestSupportsFeature(t *testing.T) {
	if !SupportsFeature("claude-opus-4-6", "thinking") {
		t.Error("Opus should support thinking")
	}
	if SupportsFeature("claude-haiku-4-5-20251001", "thinking") {
		t.Error("Haiku should not support thinking")
	}
}
