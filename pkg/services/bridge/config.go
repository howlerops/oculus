package bridge

import (
	"os"
	"strconv"

	appconfig "github.com/jbeck018/claude-go/pkg/config"
)

const (
	DefaultBridgeHost   = "127.0.0.1"
	DefaultBridgePort   = 0 // auto-assign
	DefaultPollInterval = 1000
)

// LoadBridgeConfig reads bridge settings from settings.json and env
func LoadBridgeConfig() BridgeConfig {
	cfg := BridgeConfig{
		Host:         DefaultBridgeHost,
		PollInterval: DefaultPollInterval,
		Transport:    "websocket",
	}

	// Check env for remote mode
	if os.Getenv("CLAUDE_CODE_REMOTE") == "1" {
		cfg.Enabled = true
	}

	// Check settings
	settings, _ := appconfig.LoadSettings()
	if settings != nil && settings.Env != nil {
		if val, ok := settings.Env["CLAUDE_BRIDGE_ENABLED"]; ok && val == "1" {
			cfg.Enabled = true
		}
	}

	// Auth token from env
	cfg.AuthToken = os.Getenv("CLAUDE_BRIDGE_TOKEN")
	cfg.SessionID = os.Getenv("CLAUDE_SESSION_ID")

	if portStr := os.Getenv("CLAUDE_BRIDGE_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	return cfg
}

// IsBridgeEnabled checks if the bridge should be started
func IsBridgeEnabled() bool {
	return LoadBridgeConfig().Enabled
}
