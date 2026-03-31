package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// DiffDisplay renders a colored unified diff
type DiffDisplay struct {
	FilePath string
	OldText  string
	NewText  string
	Width    int
}

// RenderDiff creates a colored diff view
func RenderDiff(filePath, oldText, newText string, width int) string {
	if width <= 0 {
		width = 80
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	sb.WriteString(headerStyle.Render("--- " + filePath + " (before)"))
	sb.WriteString("\n")
	sb.WriteString(headerStyle.Render("+++ " + filePath + " (after)"))
	sb.WriteString("\n")

	// Simple line-by-line diff
	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")

	removeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	contextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	// Find changed region
	commonPrefix := 0
	for commonPrefix < len(oldLines) && commonPrefix < len(newLines) && oldLines[commonPrefix] == newLines[commonPrefix] {
		commonPrefix++
	}

	commonSuffix := 0
	for commonSuffix < len(oldLines)-commonPrefix && commonSuffix < len(newLines)-commonPrefix {
		oi := len(oldLines) - 1 - commonSuffix
		ni := len(newLines) - 1 - commonSuffix
		if oldLines[oi] != newLines[ni] {
			break
		}
		commonSuffix++
	}

	// Show context before
	start := commonPrefix - 3
	if start < 0 {
		start = 0
	}
	for i := start; i < commonPrefix; i++ {
		sb.WriteString(contextStyle.Render(" "+oldLines[i]) + "\n")
	}

	// Show removed lines
	for i := commonPrefix; i < len(oldLines)-commonSuffix; i++ {
		line := oldLines[i]
		if len(line) > width-2 {
			line = line[:width-5] + "..."
		}
		sb.WriteString(removeStyle.Render("-"+line) + "\n")
	}

	// Show added lines
	for i := commonPrefix; i < len(newLines)-commonSuffix; i++ {
		line := newLines[i]
		if len(line) > width-2 {
			line = line[:width-5] + "..."
		}
		sb.WriteString(addStyle.Render("+"+line) + "\n")
	}

	// Show context after
	end := len(oldLines) - commonSuffix + 3
	if end > len(oldLines) {
		end = len(oldLines)
	}
	for i := len(oldLines) - commonSuffix; i < end; i++ {
		sb.WriteString(contextStyle.Render(" "+oldLines[i]) + "\n")
	}

	return sb.String()
}

// RenderEditResult formats an Edit tool result with diff colors
func RenderEditResult(filePath, oldString, newString string, width int) string {
	return RenderDiff(filePath, oldString, newString, width)
}
