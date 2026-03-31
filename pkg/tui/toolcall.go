package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ToolCallDisplay renders a tool invocation with name badge, input, and result
type ToolCallDisplay struct {
	ToolName  string
	Input     map[string]interface{}
	Result    string
	IsError   bool
	Collapsed bool
	Width     int
}

// ToolBadgeColors maps tool names to colors
var ToolBadgeColors = map[string]string{
	"Bash":      "1",  // red
	"Read":      "4",  // blue
	"Edit":      "3",  // yellow
	"Write":     "2",  // green
	"Glob":      "6",  // cyan
	"Grep":      "5",  // magenta
	"Agent":     "13", // bright magenta
	"WebFetch":  "14", // bright cyan
	"WebSearch": "11", // bright yellow
}

func (t ToolCallDisplay) View() string {
	var sb strings.Builder
	width := t.Width
	if width <= 0 {
		width = 80
	}

	// Tool name badge
	color := "8" // default gray
	if c, ok := ToolBadgeColors[t.ToolName]; ok {
		color = c
	}

	badgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(color)).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 1)

	sb.WriteString(badgeStyle.Render(t.ToolName))

	// Input summary
	inputSummary := formatToolInput(t.ToolName, t.Input)
	if inputSummary != "" {
		inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		sb.WriteString(" " + inputStyle.Render(inputSummary))
	}
	sb.WriteString("\n")

	// Result
	if t.Result != "" && !t.Collapsed {
		resultStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			PaddingLeft(2).
			Width(width - 4)

		if t.IsError {
			resultStyle = resultStyle.Foreground(lipgloss.Color("9"))
		}

		result := t.Result
		lines := strings.Split(result, "\n")
		maxLines := 20
		if len(lines) > maxLines {
			result = strings.Join(lines[:maxLines], "\n")
			result += fmt.Sprintf("\n... (%d more lines)", len(lines)-maxLines)
		}
		sb.WriteString(resultStyle.Render(result))
		sb.WriteString("\n")
	} else if t.Collapsed && t.Result != "" {
		collapsedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		lineCount := strings.Count(t.Result, "\n") + 1
		sb.WriteString("  " + collapsedStyle.Render(fmt.Sprintf("(%d lines - press Enter to expand)", lineCount)) + "\n")
	}

	return sb.String()
}

func formatToolInput(toolName string, input map[string]interface{}) string {
	if input == nil {
		return ""
	}

	switch toolName {
	case "Bash":
		if cmd, ok := input["command"].(string); ok {
			if len(cmd) > 60 {
				cmd = cmd[:57] + "..."
			}
			return "$ " + cmd
		}
	case "Read":
		if path, ok := input["file_path"].(string); ok {
			return path
		}
	case "Edit":
		if path, ok := input["file_path"].(string); ok {
			return path
		}
	case "Write":
		if path, ok := input["file_path"].(string); ok {
			return path
		}
	case "Glob":
		if pattern, ok := input["pattern"].(string); ok {
			return pattern
		}
	case "Grep":
		if pattern, ok := input["pattern"].(string); ok {
			return fmt.Sprintf("/%s/", pattern)
		}
	case "Agent":
		if desc, ok := input["description"].(string); ok {
			return desc
		}
	case "WebFetch":
		if url, ok := input["url"].(string); ok {
			return url
		}
	case "WebSearch":
		if query, ok := input["query"].(string); ok {
			return fmt.Sprintf("%q", query)
		}
	}

	// Generic: show first string key
	for _, v := range input {
		if s, ok := v.(string); ok && len(s) < 60 {
			return s
		}
	}
	return ""
}
