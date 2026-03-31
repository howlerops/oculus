package toolsummary

import (
	"fmt"
	"strings"
)

// GenerateToolUseSummary creates a compact summary of tool use.
func GenerateToolUseSummary(toolName string, input map[string]interface{}, output string) string {
	switch toolName {
	case "Bash":
		cmd, _ := input["command"].(string)
		if len(cmd) > 60 {
			cmd = cmd[:57] + "..."
		}
		lines := strings.Split(output, "\n")
		if len(lines) > 3 {
			output = strings.Join(lines[:3], "\n") + fmt.Sprintf("\n... (%d more lines)", len(lines)-3)
		}
		return fmt.Sprintf("$ %s\n%s", cmd, output)
	case "Read":
		path, _ := input["file_path"].(string)
		lines := strings.Count(output, "\n")
		return fmt.Sprintf("Read %s (%d lines)", path, lines)
	case "Edit":
		path, _ := input["file_path"].(string)
		return fmt.Sprintf("Edited %s", path)
	case "Write":
		path, _ := input["file_path"].(string)
		return fmt.Sprintf("Wrote %s", path)
	case "Glob":
		pattern, _ := input["pattern"].(string)
		matches := strings.Count(output, "\n") + 1
		return fmt.Sprintf("Glob %s -> %d files", pattern, matches)
	case "Grep":
		pattern, _ := input["pattern"].(string)
		return fmt.Sprintf("Grep %q", pattern)
	case "Agent":
		desc, _ := input["description"].(string)
		return fmt.Sprintf("Agent: %s", desc)
	default:
		if len(output) > 100 {
			output = output[:97] + "..."
		}
		return fmt.Sprintf("%s: %s", toolName, output)
	}
}
