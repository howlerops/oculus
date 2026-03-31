package components

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	PermissionBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("3")).
				Padding(1, 2)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))

	MutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	ToolNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true)
)

// titleCase capitalizes the first letter of a string
func titleCase(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// PermissionDialog renders a tool permission prompt
func PermissionDialog(toolName, description string, isReadOnly bool) string {
	var sb strings.Builder
	sb.WriteString(ToolNameStyle.Render(toolName))
	sb.WriteString("\n")
	sb.WriteString(description)
	sb.WriteString("\n\n")
	if isReadOnly {
		sb.WriteString(MutedStyle.Render("(read-only operation)"))
	} else {
		sb.WriteString(WarningStyle.Render("Allow this action?"))
	}
	sb.WriteString("\n")
	sb.WriteString(MutedStyle.Render("[y]es  [n]o  [a]lways allow  [d]eny always"))
	return PermissionBorder.Render(sb.String())
}

// ProgressBar renders a text-based progress bar of the given width
func ProgressBar(percent float64, width int) string {
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	empty := width - filled
	return fmt.Sprintf("[%s%s] %.0f%%", strings.Repeat("█", filled), strings.Repeat("░", empty), percent)
}

// TokenCounter renders a colored token usage display
func TokenCounter(input, output, total, max int) string {
	_ = input  // reserved for future use
	_ = output // reserved for future use
	pct := float64(total) / float64(max) * 100
	color := "10" // green
	if pct > 70 {
		color = "11" // yellow
	}
	if pct > 85 {
		color = "9" // red
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return style.Render(fmt.Sprintf("tokens: %dk/%dk (%.0f%%)", total/1000, max/1000, pct))
}

// SpinnerFrames holds the braille-dot spinner animation frames
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// MessageBubble renders a labelled chat message with role-appropriate styling
func MessageBubble(role, content string, width int) string {
	_ = width // reserved for future wrapping
	var style lipgloss.Style
	switch role {
	case "user":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	case "assistant":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	case "system":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
	default:
		style = lipgloss.NewStyle()
	}
	label := style.Render(titleCase(role) + ": ")
	return label + content
}
