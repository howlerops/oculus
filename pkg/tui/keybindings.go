package tui

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/jbeck018/claude-go/pkg/config"
)

// KeyBinding associates a key chord with an action
type KeyBinding struct {
	Key     string `json:"key"`
	Action  string `json:"action"`
	Command string `json:"command,omitempty"`
}

// KeyBindingsConfig holds the full set of key bindings
type KeyBindingsConfig struct {
	Bindings []KeyBinding `json:"bindings"`
}

var defaultBindings = KeyBindingsConfig{
	Bindings: []KeyBinding{
		{Key: "ctrl+c", Action: "cancel"},
		{Key: "ctrl+d", Action: "exit"},
		{Key: "ctrl+l", Action: "clear_screen"},
		{Key: "up", Action: "history_prev"},
		{Key: "down", Action: "history_next"},
		{Key: "tab", Action: "autocomplete"},
		{Key: "enter", Action: "submit"},
		{Key: "escape", Action: "cancel_edit"},
	},
}

// LoadKeyBindings reads keybindings from disk, falling back to defaults
func LoadKeyBindings() KeyBindingsConfig {
	path := filepath.Join(config.GetClaudeConfigDir(), "keybindings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultBindings
	}
	var cfg KeyBindingsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultBindings
	}
	return mergeBindings(defaultBindings, cfg)
}

// SaveKeyBindings persists the given keybindings to disk
func SaveKeyBindings(cfg KeyBindingsConfig) error {
	path := filepath.Join(config.GetClaudeConfigDir(), "keybindings.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// mergeBindings overlays custom bindings on top of defaults, keyed by action
func mergeBindings(defaults, custom KeyBindingsConfig) KeyBindingsConfig {
	actionMap := make(map[string]KeyBinding)
	for _, b := range defaults.Bindings {
		actionMap[b.Action] = b
	}
	for _, b := range custom.Bindings {
		actionMap[b.Action] = b
	}
	var result KeyBindingsConfig
	for _, b := range actionMap {
		result.Bindings = append(result.Bindings, b)
	}
	return result
}
