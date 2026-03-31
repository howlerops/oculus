package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// GlobalConfig holds persistent user configuration
type GlobalConfig struct {
	HasAcceptedTrustDialog bool                    `json:"hasAcceptedTrustDialog,omitempty"`
	Projects               map[string]ProjectConfig `json:"projects,omitempty"`
	NumConversations       int                     `json:"numConversations,omitempty"`
}

type ProjectConfig struct {
	AllowedTools []string `json:"allowedTools,omitempty"`
	DeniedTools  []string `json:"deniedTools,omitempty"`
}

var (
	globalConfig   *GlobalConfig
	globalConfigMu sync.RWMutex
)

func getGlobalConfigPath() string {
	return filepath.Join(GetOculusDir(), "config.json")
}

// GetGlobalConfig loads or returns cached global config
func GetGlobalConfig() *GlobalConfig {
	globalConfigMu.RLock()
	if globalConfig != nil {
		defer globalConfigMu.RUnlock()
		return globalConfig
	}
	globalConfigMu.RUnlock()

	globalConfigMu.Lock()
	defer globalConfigMu.Unlock()

	data, err := os.ReadFile(getGlobalConfigPath())
	if err != nil {
		globalConfig = &GlobalConfig{}
		return globalConfig
	}

	var cfg GlobalConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		globalConfig = &GlobalConfig{}
		return globalConfig
	}

	globalConfig = &cfg
	return globalConfig
}

// SaveGlobalConfig persists the global config to disk
func SaveGlobalConfig(cfg *GlobalConfig) error {
	globalConfigMu.Lock()
	defer globalConfigMu.Unlock()

	globalConfig = cfg

	dir := filepath.Dir(getGlobalConfigPath())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getGlobalConfigPath(), data, 0o644)
}
