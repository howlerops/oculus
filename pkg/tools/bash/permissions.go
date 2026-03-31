package bash

import (
	"path/filepath"
	"strings"

	"github.com/howlerops/oculus/pkg/types"
)

// BashPermissionCheck evaluates bash command against permission rules
func BashPermissionCheck(command string, ctx types.ToolPermissionContext) types.PermissionResult {
	// Bypass mode allows everything
	if ctx.Mode == types.PermissionModeBypassPermissions {
		return types.PermissionResult{
			Behavior:     types.PermissionAllow,
			UpdatedInput: map[string]interface{}{"command": command},
		}
	}

	// Check deny rules
	for _, patterns := range ctx.AlwaysDenyRules {
		for _, pattern := range patterns {
			if matchBashRule(command, pattern) {
				return types.PermissionResult{
					Behavior: types.PermissionDeny,
					Message:  "Command denied by rule: " + pattern,
				}
			}
		}
	}

	// Check allow rules
	for _, patterns := range ctx.AlwaysAllowRules {
		for _, pattern := range patterns {
			if matchBashRule(command, pattern) {
				return types.PermissionResult{
					Behavior:     types.PermissionAllow,
					UpdatedInput: map[string]interface{}{"command": command},
				}
			}
		}
	}

	// Read-only commands in acceptEdits mode
	if ctx.Mode == types.PermissionModeAcceptEdits && IsReadOnlyCommand(command) {
		return types.PermissionResult{
			Behavior:     types.PermissionAllow,
			UpdatedInput: map[string]interface{}{"command": command},
		}
	}

	// DontAsk mode allows everything
	if ctx.Mode == types.PermissionModeDontAsk {
		return types.PermissionResult{
			Behavior:     types.PermissionAllow,
			UpdatedInput: map[string]interface{}{"command": command},
		}
	}

	// Default: ask
	return types.PermissionResult{
		Behavior: types.PermissionAsk,
		Message:  "Permission required for: " + truncateCmd(command, 60),
	}
}

// matchBashRule checks if a command matches a Bash permission rule pattern
func matchBashRule(command, pattern string) bool {
	// Pattern formats: "Bash", "Bash(cmd)", "Bash(cmd *)", "Bash(git *)"
	if !strings.HasPrefix(pattern, "Bash") {
		return false
	}

	// Just "Bash" matches all
	if pattern == "Bash" {
		return true
	}

	// Extract content from Bash(content)
	if !strings.HasPrefix(pattern, "Bash(") || !strings.HasSuffix(pattern, ")") {
		return false
	}
	content := pattern[5 : len(pattern)-1]

	// Exact match
	if command == content {
		return true
	}

	// Wildcard matching
	if strings.Contains(content, "*") {
		matched, _ := filepath.Match(content, command)
		if matched {
			return true
		}
		// Also try matching just the command prefix
		if strings.HasSuffix(content, " *") {
			prefix := strings.TrimSuffix(content, " *")
			if strings.HasPrefix(command, prefix+" ") || command == prefix {
				return true
			}
		}
		if strings.HasSuffix(content, "*") {
			prefix := strings.TrimSuffix(content, "*")
			if strings.HasPrefix(command, prefix) {
				return true
			}
		}
	}

	// Prefix match
	if strings.HasPrefix(command, content) {
		return true
	}

	return false
}

func truncateCmd(cmd string, max int) string {
	if len(cmd) <= max {
		return cmd
	}
	return cmd[:max-3] + "..."
}
