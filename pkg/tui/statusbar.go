package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar shows model, branch, tokens, and session info in the footer
type StatusBar struct {
	Model        string
	GitBranch    string
	InputTokens  int
	OutputTokens int
	TotalCost    float64
	MaxContext   int
	SessionID    string
	Width        int
}

func NewStatusBar(width int) StatusBar {
	return StatusBar{Width: width, MaxContext: 200000}
}

func (s StatusBar) View() string {
	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("7")).
		Width(s.Width)

	var parts []string

	// Model
	if s.Model != "" {
		modelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true)
		// Short name
		short := s.Model
		switch {
		case strings.Contains(short, "opus"):
			short = "opus"
		case strings.Contains(short, "sonnet"):
			short = "sonnet"
		case strings.Contains(short, "haiku"):
			short = "haiku"
		}
		parts = append(parts, modelStyle.Render(short))
	}

	// Git branch
	if s.GitBranch != "" {
		branchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
		parts = append(parts, branchStyle.Render("⎇ "+s.GitBranch))
	}

	// Token count + context bar
	if s.InputTokens > 0 || s.OutputTokens > 0 {
		total := s.InputTokens + s.OutputTokens
		pct := float64(total) / float64(s.MaxContext) * 100

		color := "10" // green
		if pct > 70 {
			color = "11"
		} // yellow
		if pct > 85 {
			color = "9"
		} // red

		tokenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		tokens := fmt.Sprintf("tokens: %s", formatTokenCount(total))
		parts = append(parts, tokenStyle.Render(tokens))

		// Mini context bar
		barWidth := 10
		filled := int(pct / 100 * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}
		empty := barWidth - filled
		bar := tokenStyle.Render(strings.Repeat("█", filled)) +
			lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(strings.Repeat("░", empty))
		parts = append(parts, bar+tokenStyle.Render(fmt.Sprintf(" %.0f%%", pct)))
	}

	// Cost
	if s.TotalCost > 0 {
		costStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		parts = append(parts, costStyle.Render(fmt.Sprintf("$%.4f", s.TotalCost)))
	}

	content := " " + strings.Join(parts, "  │  ") + " "
	return barStyle.Render(content)
}

func formatTokenCount(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}
