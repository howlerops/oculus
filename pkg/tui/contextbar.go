package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ContextBar displays context window usage as a progress bar
type ContextBar struct {
	UsedTokens   int
	MaxTokens    int
	Width        int
	ShowDetailed bool
}

func NewContextBar(maxTokens, width int) ContextBar {
	if maxTokens == 0 {
		maxTokens = 200000
	}
	return ContextBar{MaxTokens: maxTokens, Width: width}
}

func (c ContextBar) Percent() float64 {
	if c.MaxTokens == 0 {
		return 0
	}
	return float64(c.UsedTokens) / float64(c.MaxTokens) * 100
}

func (c ContextBar) Level() string {
	pct := c.Percent()
	if pct >= 85 {
		return "critical"
	}
	if pct >= 70 {
		return "warning"
	}
	return "ok"
}

func (c ContextBar) View() string {
	pct := c.Percent()
	barWidth := c.Width - 20
	if barWidth < 10 {
		barWidth = 10
	}

	filled := int(pct / 100 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}
	empty := barWidth - filled

	// Color based on level
	var filledColor, labelColor string
	switch c.Level() {
	case "critical":
		filledColor = "9"
		labelColor = "9"
	case "warning":
		filledColor = "11"
		labelColor = "11"
	default:
		filledColor = "10"
		labelColor = "10"
	}

	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(filledColor))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(labelColor))

	bar := filledStyle.Render(strings.Repeat("█", filled)) + emptyStyle.Render(strings.Repeat("░", empty))

	label := labelStyle.Render(fmt.Sprintf("ctx: %.0f%%", pct))

	if c.ShowDetailed {
		detail := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(
			fmt.Sprintf(" (%s/%s)", formatTokenCount(c.UsedTokens), formatTokenCount(c.MaxTokens)))
		return fmt.Sprintf("%s [%s] %s", label, bar, detail)
	}

	return fmt.Sprintf("%s [%s]", label, bar)
}
