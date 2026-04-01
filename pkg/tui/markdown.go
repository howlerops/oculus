package tui

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// MarkdownRenderer renders markdown text for terminal display
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
	width    int
}

// NewMarkdownRenderer creates a renderer with the given width
func NewMarkdownRenderer(width int) *MarkdownRenderer {
	if width <= 0 {
		width = 80
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("notty"),
		glamour.WithWordWrap(width-4),
	)
	if err != nil {
		// Fallback: no markdown rendering
		return &MarkdownRenderer{width: width}
	}

	return &MarkdownRenderer{
		renderer: renderer,
		width:    width,
	}
}

// Render converts markdown to terminal-formatted text
func (mr *MarkdownRenderer) Render(markdown string) string {
	if mr.renderer == nil {
		return markdown // fallback: raw text
	}

	rendered, err := mr.renderer.Render(markdown)
	if err != nil {
		return markdown
	}

	// Trim trailing whitespace from glamour output
	return strings.TrimRight(rendered, "\n ")
}

// RenderStreaming renders markdown incrementally for streaming display.
// Returns the rendered portion and whether the content appears complete.
func (mr *MarkdownRenderer) RenderStreaming(partial string) (string, bool) {
	// Don't render mid-code-block (would break formatting)
	openBlocks := strings.Count(partial, "```")
	if openBlocks%2 != 0 {
		// Inside a code block - render everything before the last ``` normally
		// and the code block content as raw
		lastBlock := strings.LastIndex(partial, "```")
		rendered := mr.Render(partial[:lastBlock])
		codeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		return rendered + "\n" + codeStyle.Render(partial[lastBlock:]), false
	}

	return mr.Render(partial), true
}

// RenderInline renders a short markdown snippet without block-level formatting
func (mr *MarkdownRenderer) RenderInline(text string) string {
	text = renderBold(text)
	text = renderItalic(text)
	text = renderInlineCode(text)
	return text
}

func renderBold(text string) string {
	style := lipgloss.NewStyle().Bold(true)
	for {
		start := strings.Index(text, "**")
		if start == -1 {
			break
		}
		end := strings.Index(text[start+2:], "**")
		if end == -1 {
			break
		}
		inner := text[start+2 : start+2+end]
		text = text[:start] + style.Render(inner) + text[start+2+end+2:]
	}
	return text
}

func renderItalic(text string) string {
	style := lipgloss.NewStyle().Italic(true)
	for {
		start := strings.Index(text, "*")
		if start == -1 {
			break
		}
		// Skip if it's a ** (bold)
		if start+1 < len(text) && text[start+1] == '*' {
			break
		}
		end := strings.Index(text[start+1:], "*")
		if end == -1 {
			break
		}
		inner := text[start+1 : start+1+end]
		text = text[:start] + style.Render(inner) + text[start+1+end+1:]
	}
	return text
}

func renderInlineCode(text string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Background(lipgloss.Color("236"))
	for {
		start := strings.Index(text, "`")
		if start == -1 {
			break
		}
		// Skip code blocks
		if start+2 < len(text) && text[start:start+3] == "```" {
			break
		}
		end := strings.Index(text[start+1:], "`")
		if end == -1 {
			break
		}
		inner := text[start+1 : start+1+end]
		text = text[:start] + style.Render(inner) + text[start+1+end+1:]
	}
	return text
}
