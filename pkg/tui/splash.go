package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const oculusASCII = `
   ___   ____  _   _ _     _   _ ____
  / _ \ / ___|| | | | |   | | | / ___|
 | | | | |    | | | | |   | | | \___ \
 | |_| | |___ | |_| | |___| |_| |___) |
  \___/ \____| \___/|_____|\_____/____/
`

const tagline = "AI Coding Assistant • by HowlerOps"

// RenderSplash returns the splash screen for interactive mode
func RenderSplash(width int) string {
	var sb strings.Builder

	// Logo in sky blue
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0ea5e9")).
		Bold(true)

	// Tagline in cyan
	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#22d3ee")).
		Italic(true)

	// Version in muted
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64748b"))

	sb.WriteString("\n")
	for _, line := range strings.Split(oculusASCII, "\n") {
		centered := centerText(line, width)
		sb.WriteString(logoStyle.Render(centered))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString(tagStyle.Render(centerText(tagline, width)))
	sb.WriteString("\n")
	sb.WriteString(versionStyle.Render(centerText("v0.5.1 • github.com/howlerops/oculus", width)))
	sb.WriteString("\n\n")

	return sb.String()
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	pad := (width - len(text)) / 2
	return strings.Repeat(" ", pad) + text
}
