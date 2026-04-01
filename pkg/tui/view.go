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

// View renders the full TUI filling the entire terminal
func (m Model) View() string {
	// Permission dialog takes over the screen
	if m.state == StatePermission && m.permission != nil {
		return padToHeight(m.permission.View(), m.width, m.height)
	}

	// Build bottom section (always visible): autocomplete + input + status bar
	var bottomParts []string

	// Autocomplete dropdown
	if m.autocomplete != nil && m.autocomplete.Active && len(m.autocomplete.Matches) > 0 {
		bottomParts = append(bottomParts, m.autocomplete.View())
	}

	// Input area
	bottomParts = append(bottomParts, m.input.View())

	// Status bar
	m.statusBar.Width = m.width
	m.statusBar.Model = m.getModelName()
	bottomParts = append(bottomParts, m.statusBar.View())

	bottom := strings.Join(bottomParts, "\n")
	bottomLines := strings.Count(bottom, "\n") + 1

	// Calculate how much space the main content area gets
	mainHeight := m.height - bottomLines
	if mainHeight < 1 {
		mainHeight = 1
	}

	// Build main content
	var mainContent string

	if len(m.messages) == 0 && m.viewport.content.Len() == 0 && m.state == StateChat {
		// Splash screen - center it vertically in the main area
		splash := RenderSplash(m.width)
		splashLines := strings.Count(splash, "\n") + 1
		topPad := (mainHeight - splashLines) / 3 // bias toward top third
		if topPad < 0 {
			topPad = 0
		}
		bottomPad := mainHeight - splashLines - topPad
		if bottomPad < 0 {
			bottomPad = 0
		}
		mainContent = strings.Repeat("\n", topPad) + splash + strings.Repeat("\n", bottomPad)
	} else {
		// Normal conversation view
		var parts []string

		// Header
		header := headerViewStyle.Render("◉ Oculus")
		if m.state == StateLoading {
			header += " " + m.spinner.View() + mutedViewStyle.Render(" thinking...")
		}
		parts = append(parts, header)

		// Viewport (scrollable messages)
		viewportHeight := mainHeight - 2 // minus header and padding
		if viewportHeight < 3 {
			viewportHeight = 3
		}
		m.viewport.viewport.Height = viewportHeight
		m.viewport.viewport.Width = m.width
		parts = append(parts, m.viewport.View())

		// Loading indicator in viewport area
		if m.state == StateLoading && m.streamBuffer == "" && !m.progress.HasRunning() {
			loadingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0ea5e9"))
			frames := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
			frame := frames[int(m.spinner.Spinner.FPS)%len(frames)]
			parts = append(parts, loadingStyle.Render(fmt.Sprintf("\n  %s Working...\n", frame)))
		}

		// Tool progress
		if m.progress.HasRunning() {
			parts = append(parts, m.progress.View())
		}

		// Streaming buffer
		if m.streamBuffer != "" && m.state == StateLoading {
			streamStyle := assistStyle.Width(m.width - 4)
			parts = append(parts, streamStyle.Render(m.streamBuffer))
		}

		// Error
		if m.err != nil {
			parts = append(parts, errViewStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		}

		// Task panel
		taskView := m.taskPanel.View()
		if taskView != "" {
			parts = append(parts, taskView)
		}

		mainContent = strings.Join(parts, "\n")

		// Pad main content to fill available height
		currentLines := strings.Count(mainContent, "\n") + 1
		if currentLines < mainHeight {
			mainContent += strings.Repeat("\n", mainHeight-currentLines)
		}
	}

	return mainContent + "\n" + bottom
}

// padToHeight ensures output fills the terminal
func padToHeight(content string, width, height int) string {
	lines := strings.Count(content, "\n") + 1
	if lines < height {
		content += strings.Repeat("\n", height-lines)
	}
	return content
}

// getModelName returns the current model name for the status bar
func (m Model) getModelName() string {
	if m.engine != nil {
		return m.engine.Model
	}
	return "unknown"
}
