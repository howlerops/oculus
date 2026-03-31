package config

import (
	"fmt"
	"strings"
)

// ValidationError represents a settings validation error
type ValidationError struct {
	Path    string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("settings.%s: %s", e.Path, e.Message)
}

// ValidateSettings checks settings for errors
func ValidateSettings(s *SettingsJson) []ValidationError {
	var errs []ValidationError

	// Validate permission mode
	if s.DefaultMode != "" {
		validModes := map[string]bool{
			"default": true, "acceptEdits": true, "bypassPermissions": true,
			"dontAsk": true, "plan": true, "auto": true,
		}
		if !validModes[s.DefaultMode] {
			errs = append(errs, ValidationError{
				Path:    "defaultMode",
				Message: fmt.Sprintf("invalid mode %q", s.DefaultMode),
			})
		}
	}

	// Validate hooks
	validEvents := map[string]bool{
		"PreToolUse": true, "PostToolUse": true, "Notification": true,
		"PreCompact": true, "PostCompact": true, "SessionStart": true,
		"Stop": true, "SubagentStop": true, "UserPromptSubmit": true,
	}
	for event := range s.Hooks {
		if !validEvents[event] {
			errs = append(errs, ValidationError{
				Path:    "hooks." + event,
				Message: "unknown hook event",
			})
		}
		for i, hook := range s.Hooks[event] {
			if hook.Command == "" {
				errs = append(errs, ValidationError{
					Path:    fmt.Sprintf("hooks.%s[%d].command", event, i),
					Message: "command is required",
				})
			}
		}
	}

	// Validate MCP servers
	for name, server := range s.MCPServers {
		if server.Command == "" && server.URL == "" {
			errs = append(errs, ValidationError{
				Path:    "mcpServers." + name,
				Message: "either command or url is required",
			})
		}
		if server.Transport != "" && server.Transport != "stdio" && server.Transport != "http" {
			errs = append(errs, ValidationError{
				Path:    "mcpServers." + name + ".transport",
				Message: fmt.Sprintf("invalid transport %q (use stdio or http)", server.Transport),
			})
		}
	}

	return errs
}

// ValidateToolName checks if a tool name is valid for rules
func ValidateToolName(name string) bool {
	if name == "" {
		return false
	}
	// Must be alphanumeric with optional (content) suffix
	base := name
	if idx := strings.Index(name, "("); idx >= 0 {
		base = name[:idx]
	}
	for _, c := range base {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}
