package config

import (
	"os"
	"strings"
)

// IsEnvTruthy returns true if an env var is set to a truthy value
func IsEnvTruthy(value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	return v == "1" || v == "true" || v == "yes"
}

// IsBareMode returns true if --bare flag or CLAUDE_CODE_BARE is set
func IsBareMode() bool {
	return IsEnvTruthy(os.Getenv("CLAUDE_CODE_BARE"))
}

// IsRemoteMode returns true if running in remote/CCR mode
func IsRemoteMode() bool {
	return IsEnvTruthy(os.Getenv("CLAUDE_CODE_REMOTE"))
}

// IsTestMode returns true during tests
func IsTestMode() bool {
	return os.Getenv("GO_TEST") == "1" || strings.HasSuffix(os.Args[0], ".test")
}

// GetAPIKey returns the Anthropic API key from env
func GetAPIKey() string {
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return key
	}
	return os.Getenv("CLAUDE_API_KEY")
}

// GetModel returns the model to use, with fallback
func GetModel() string {
	if model := os.Getenv("ANTHROPIC_MODEL"); model != "" {
		return model
	}
	return "claude-sonnet-4-20250514"
}
