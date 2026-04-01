package bridge

import "testing"

func TestCreateBridge(t *testing.T) {
	tests := []struct {
		provider string
		wantName string
		wantErr  bool
	}{
		{"anthropic", "anthropic", false},
		{"openai", "openai", false},
		{"ollama", "ollama", false},
		{"claude-code", "claude-code", false},
		{"codex", "codex", false},
		{"gemini-cli", "gemini-cli", false},
		{"", "anthropic", false},       // default
		{"unknown-provider", "", true}, // error
	}
	for _, tt := range tests {
		b, err := CreateBridge(BridgeConfig{Provider: tt.provider, APIKey: "test"})
		if tt.wantErr {
			if err == nil {
				t.Errorf("CreateBridge(%q) expected error", tt.provider)
			}
			continue
		}
		if err != nil {
			t.Errorf("CreateBridge(%q) error: %v", tt.provider, err)
		}
		if b.Name() != tt.wantName {
			t.Errorf("CreateBridge(%q).Name() = %q, want %q", tt.provider, b.Name(), tt.wantName)
		}
	}
}

func TestAnthropicBridgeAvailability(t *testing.T) {
	b := NewAnthropicBridge(BridgeConfig{APIKey: "test-key"})
	if !b.IsAvailable() {
		t.Error("should be available with API key")
	}

	b2 := NewAnthropicBridge(BridgeConfig{})
	if b2.IsAvailable() {
		t.Error("should not be available without API key")
	}
}

func TestOpenAIBridgeDefaults(t *testing.T) {
	b := NewOpenAIBridge(BridgeConfig{Provider: "openai"})
	if b.Name() != "openai" {
		t.Error("name should be openai")
	}
}

func TestCLIBridgeDetection(t *testing.T) {
	// These may or may not be available
	b := NewCLIBridge(BridgeConfig{Provider: "claude-code"})
	if b.Name() != "claude-code" {
		t.Error("name should be claude-code")
	}
	// IsAvailable depends on whether claude is installed
}
