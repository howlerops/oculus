package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ToolProgress tracks a running tool's progress
type ToolProgress struct {
	ToolName  string
	ToolID    string
	StartTime time.Time
	Status    string // "running", "completed", "error"
	Message   string
}

// ToolProgressTracker manages multiple concurrent tool progress indicators
type ToolProgressTracker struct {
	tools  map[string]*ToolProgress
	order  []string
	frame  int
	frames []string
}

func NewToolProgressTracker() *ToolProgressTracker {
	return &ToolProgressTracker{
		tools:  make(map[string]*ToolProgress),
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

func (t *ToolProgressTracker) StartTool(id, name string) {
	t.tools[id] = &ToolProgress{
		ToolName:  name,
		ToolID:    id,
		StartTime: time.Now(),
		Status:    "running",
	}
	t.order = append(t.order, id)
}

func (t *ToolProgressTracker) CompleteTool(id string) {
	if tool, ok := t.tools[id]; ok {
		tool.Status = "completed"
	}
}

func (t *ToolProgressTracker) ErrorTool(id string, msg string) {
	if tool, ok := t.tools[id]; ok {
		tool.Status = "error"
		tool.Message = msg
	}
}

func (t *ToolProgressTracker) Tick() {
	t.frame = (t.frame + 1) % len(t.frames)
}

func (t *ToolProgressTracker) HasRunning() bool {
	for _, tool := range t.tools {
		if tool.Status == "running" {
			return true
		}
	}
	return false
}

func (t *ToolProgressTracker) Clear() {
	t.tools = make(map[string]*ToolProgress)
	t.order = nil
}

func (t *ToolProgressTracker) View() string {
	var lines []string

	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	toolStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	doneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	for _, id := range t.order {
		tool, ok := t.tools[id]
		if !ok {
			continue
		}

		elapsed := time.Since(tool.StartTime).Round(time.Second)

		switch tool.Status {
		case "running":
			spinner := spinnerStyle.Render(t.frames[t.frame])
			name := toolStyle.Render(tool.ToolName)
			dur := timeStyle.Render(fmt.Sprintf("%s", elapsed))
			lines = append(lines, fmt.Sprintf("  %s %s %s", spinner, name, dur))
		case "completed":
			check := doneStyle.Render("✓")
			name := toolStyle.Render(tool.ToolName)
			dur := timeStyle.Render(fmt.Sprintf("%s", elapsed))
			lines = append(lines, fmt.Sprintf("  %s %s %s", check, name, dur))
		case "error":
			x := errStyle.Render("✗")
			name := toolStyle.Render(tool.ToolName)
			lines = append(lines, fmt.Sprintf("  %s %s %s", x, name, errStyle.Render(tool.Message)))
		}
	}

	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}
