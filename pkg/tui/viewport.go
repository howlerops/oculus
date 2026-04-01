package tui

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// titleCase capitalizes the first letter of s (replaces deprecated strings.Title).
func titleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// MessageViewport is a scrollable message display
type MessageViewport struct {
	viewport viewport.Model
	content  *strings.Builder
	width    int
	height   int
	ready    bool
}

func NewMessageViewport(width, height int) MessageViewport {
	vp := viewport.New(width, height)
	vp.MouseWheelEnabled = true
	return MessageViewport{
		viewport: vp,
		content:  &strings.Builder{},
		width:    width,
		height:   height,
		ready:    true,
	}
}

func (m MessageViewport) Init() tea.Cmd { return nil }

func (m MessageViewport) Update(msg tea.Msg) (MessageViewport, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 6 // leave room for input + status
		m.viewport.Width = m.width
		m.viewport.Height = m.height
	}
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// AppendMessage adds a rendered message to the viewport
func (m *MessageViewport) AppendMessage(role, content string) {
	roleStyle := lipgloss.NewStyle().Bold(true)
	switch role {
	case "user":
		roleStyle = roleStyle.Foreground(lipgloss.Color("12"))
	case "assistant":
		roleStyle = roleStyle.Foreground(lipgloss.Color("10"))
	case "system":
		roleStyle = roleStyle.Foreground(lipgloss.Color("8"))
	case "tool":
		roleStyle = roleStyle.Foreground(lipgloss.Color("3"))
	}

	if m.content.Len() > 0 {
		m.content.WriteString("\n")
	}
	m.content.WriteString(roleStyle.Render(titleCase(role) + ": "))
	m.content.WriteString(content)
	m.content.WriteString("\n")

	m.viewport.SetContent(m.content.String())
	m.viewport.GotoBottom()
}

// AppendRaw adds raw text to the viewport
func (m *MessageViewport) AppendRaw(text string) {
	m.content.WriteString(text)
	m.viewport.SetContent(m.content.String())
	m.viewport.GotoBottom()
}

// UpdateStreamingContent updates the last message (for streaming)
func (m *MessageViewport) UpdateStreamingContent(fullContent string) {
	// Replace content from last role marker to end
	current := m.content.String()
	lastNewline := strings.LastIndex(current, "\n\n")
	if lastNewline >= 0 {
		m.content.Reset()
		m.content.WriteString(current[:lastNewline+2])
		m.content.WriteString(fullContent)
	} else {
		m.content.Reset()
		m.content.WriteString(fullContent)
	}
	m.viewport.SetContent(m.content.String())
	m.viewport.GotoBottom()
}

// Clear resets the viewport
func (m *MessageViewport) Clear() {
	m.content.Reset()
	m.viewport.SetContent("")
}

// ScrollPercent returns how far scrolled (0-100)
func (m MessageViewport) ScrollPercent() float64 {
	return m.viewport.ScrollPercent() * 100
}

// AtBottom checks if scrolled to bottom
func (m MessageViewport) AtBottom() bool {
	return m.viewport.AtBottom()
}

func (m MessageViewport) View() string {
	return m.viewport.View()
}
