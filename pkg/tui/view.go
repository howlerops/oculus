package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	userMsgStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#0ea5e9")).Bold(true)
	assistStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#22d3ee"))
	errViewStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))
	headerViewStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0ea5e9")).Bold(true)
	mutedViewStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b"))
)

// View renders the full TUI
func (m Model) View() string {
	var sections []string

	// Permission dialog takes over the screen
	if m.state == StatePermission && m.permission != nil {
		return m.permission.View()
	}

	// Show splash if no messages and no viewport content
	if len(m.messages) == 0 && m.viewport.content.Len() == 0 && m.state == StateChat {
		splash := RenderSplash(m.width)
		return splash + "\n" + m.input.View() + "\n" + m.statusBar.View()
	}

	// Header
	header := headerViewStyle.Render("◉ Oculus")
	if m.state == StateLoading {
		header += " " + m.spinner.View() + mutedViewStyle.Render(" thinking...")
	}
	sections = append(sections, header)

	// Main content area - scrollable viewport
	viewportHeight := m.height - 6 // header(1) + input(3) + status(1) + padding(1)
	if viewportHeight < 5 {
		viewportHeight = 5
	}
	m.viewport.viewport.Height = viewportHeight
	sections = append(sections, m.viewport.View())

	// Tool progress (if any tools running)
	if m.progress.HasRunning() {
		sections = append(sections, m.progress.View())
	}

	// Streaming buffer (live text not yet in viewport)
	if m.streamBuffer != "" && m.state == StateLoading {
		streamStyle := assistStyle.Width(m.width - 4)
		sections = append(sections, streamStyle.Render(m.streamBuffer))
	}

	// Error display
	if m.err != nil {
		sections = append(sections, errViewStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	// Task panel (inline, if visible and has tasks)
	taskView := m.taskPanel.View()
	if taskView != "" {
		sections = append(sections, taskView)
	}

	// Input area
	sections = append(sections, m.input.View())

	// Status bar + context bar footer
	m.statusBar.Model = m.getModelName()
	footer := m.statusBar.View()
	sections = append(sections, footer)

	return strings.Join(sections, "\n")
}

// getModelName returns the current model name for the status bar
func (m Model) getModelName() string {
	if m.engine != nil {
		return m.engine.Model
	}
	return "unknown"
}
