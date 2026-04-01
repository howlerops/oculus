package bridge

import "fmt"

func CreateBridge(cfg BridgeConfig) (Bridge, error) {
	switch cfg.Provider {
	case "anthropic", "": return NewAnthropicBridge(cfg), nil
	case "openai": return NewOpenAIBridge(cfg), nil
	case "ollama": cfg.BaseURL = "http://localhost:11434/v1"; return NewOpenAIBridge(cfg), nil
	case "claude-code", "codex", "gemini-cli": return NewCLIBridge(cfg), nil
	default: return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}

func AvailableProviders() []string {
	return []string{"anthropic", "openai", "ollama", "claude-code", "codex", "gemini-cli"}
}
