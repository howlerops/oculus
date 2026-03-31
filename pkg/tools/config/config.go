package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	appconfig "github.com/howlerops/oculus/pkg/config"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type ConfigTool struct {
	tool.BaseTool
}

func NewConfigTool() *ConfigTool {
	return &ConfigTool{
		BaseTool: tool.BaseTool{ToolName: "Config", ToolSearchHint: "settings configuration preference"},
	}
}

func (t *ConfigTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"setting": map[string]interface{}{"type": "string", "description": "Setting key (e.g. theme, model)"},
			"value":   map[string]interface{}{"description": "Value to set (omit for GET)"},
		},
		Required: []string{"setting"},
	}
}

func (t *ConfigTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Read or write Claude Code configuration settings.", nil
}

func (t *ConfigTool) Prompt(_ context.Context) (string, error) {
	return "Get or set Claude Code configuration settings.\n\nUsage:\n- Get: omit value parameter\n- Set: include value parameter\n\nSettings: theme, model, verbose, editorMode, permissions.defaultMode", nil
}

func (t *ConfigTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	setting, _ := input["setting"].(string)
	if setting == "" {
		return &tool.Result{Data: "Error: setting is required"}, nil
	}

	value, hasValue := input["value"]

	// Load settings
	settings, err := appconfig.LoadSettings()
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error loading settings: %v", err)}, nil
	}

	// Convert to map for dynamic access
	data, _ := json.Marshal(settings)
	var settingsMap map[string]interface{}
	json.Unmarshal(data, &settingsMap)

	if !hasValue || value == nil {
		// GET
		val := getNestedValue(settingsMap, setting)
		if val == nil {
			return &tool.Result{Data: fmt.Sprintf("Setting %q is not set", setting)}, nil
		}
		valStr, _ := json.Marshal(val)
		return &tool.Result{
			Data: fmt.Sprintf("Setting %q = %s", setting, string(valStr)),
		}, nil
	}

	// SET
	setNestedValue(settingsMap, setting, value)
	newData, _ := json.MarshalIndent(settingsMap, "", "  ")
	if err := os.WriteFile(appconfig.GetSettingsPath(), newData, 0o644); err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error saving settings: %v", err)}, nil
	}

	return &tool.Result{
		Data: fmt.Sprintf("Setting %q updated to %v", setting, value),
	}, nil
}

func getNestedValue(m map[string]interface{}, key string) interface{} {
	parts := strings.Split(key, ".")
	current := interface{}(m)
	for _, part := range parts {
		cm, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current = cm[part]
	}
	return current
}

func setNestedValue(m map[string]interface{}, key string, value interface{}) {
	parts := strings.Split(key, ".")
	current := m
	for i, part := range parts {
		if i == len(parts)-1 {
			current[part] = value
			return
		}
		next, ok := current[part].(map[string]interface{})
		if !ok {
			next = make(map[string]interface{})
			current[part] = next
		}
		current = next
	}
}
