package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputModel is an enhanced text input with multi-line and history
type InputModel struct {
	textarea      textarea.Model
	history       []string
	historyIdx    int
	searching     bool
	searchQuery   string
	searchResults []int
	searchIdx     int
	submitted     bool
	width         int
}


// filterEscapes strips terminal escape/OSC sequences from text.
// Handles ESC[...m (CSI), ESC]...ST (OSC), and bare ]11;... responses.
func filterEscapes(s string) string {
	var result []byte
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		// ESC-prefixed sequences
		if b[i] == 0x1b && i+1 < len(b) {
			if b[i+1] == '[' {
				// CSI: ESC [ params letter
				j := i + 2
				for j < len(b) && ((b[j] >= '0' && b[j] <= '9') || b[j] == ';') {
					j++
				}
				if j < len(b) {
					j++ // skip final letter
				}
				i = j - 1
				continue
			}
			if b[i+1] == ']' {
				// OSC: ESC ] ... BEL or ESC backslash
				j := i + 2
				for j < len(b) && b[j] != 0x07 && b[j] != 0x1b {
					j++
				}
				if j < len(b) && b[j] == 0x07 {
					i = j
					continue
				}
				if j+1 < len(b) && b[j] == 0x1b && b[j+1] == '\\' {
					i = j + 1
					continue
				}
				i = j - 1
				continue
			}
			i++ // skip ESC + next char
			continue
		}
		// Bare OSC response: ]11;rgb:... backslash
		if b[i] == ']' && i+1 < len(b) && b[i+1] >= '0' && b[i+1] <= '9' {
			j := i + 1
			for j < len(b) && b[j] != 0x07 && b[j] != '\\' && b[j] != '\n' {
				j++
			}
			if j < len(b) && (b[j] == 0x07 || b[j] == '\\') {
				i = j
				continue
			}
		}
		result = append(result, b[i])
	}
	return string(result)
}

func NewInputModel() InputModel {
	ta := textarea.New()
	ta.Placeholder = "Type your message... (Enter to send, Shift+Enter for newline)"
	ta.Focus()
	ta.SetHeight(3)
	ta.SetWidth(80)
	ta.ShowLineNumbers = false
	ta.CharLimit = 0 // unlimited

	return InputModel{
		textarea:   ta,
		historyIdx: -1,
		width:      80,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Filter out escape sequence responses (OSC, CSI) that some terminals inject
		if msg.Type == tea.KeyRunes {
			s := string(msg.Runes)
			if strings.Contains(s, "]11;") || strings.Contains(s, "]10;") ||
				strings.Contains(s, "rgb:") || strings.HasPrefix(s, "\x1b") {
				return m, nil // swallow escape response
			}
		}

		// Handle search mode
		if m.searching {
			return m.updateSearch(msg)
		}

		switch msg.Type {
		case tea.KeyEnter:
			// Check for alt+enter (newline)
			if msg.Alt {
				var cmd tea.Cmd
				m.textarea, cmd = m.textarea.Update(msg)
				return m, cmd
			}
			// Submit
			value := strings.TrimSpace(m.textarea.Value())
			if value != "" {
				m.submitted = true
				m.history = append(m.history, value)
				m.historyIdx = len(m.history)
			}
			return m, nil

		case tea.KeyUp:
			// History navigation (only when on first line)
			if m.textarea.Line() == 0 {
				if m.historyIdx > 0 {
					m.historyIdx--
					m.textarea.SetValue(m.history[m.historyIdx])
					m.textarea.CursorEnd()
				}
				return m, nil
			}

		case tea.KeyDown:
			// History navigation (only when on last line)
			lines := strings.Count(m.textarea.Value(), "\n")
			if m.textarea.Line() == lines {
				if m.historyIdx < len(m.history)-1 {
					m.historyIdx++
					m.textarea.SetValue(m.history[m.historyIdx])
				} else {
					m.historyIdx = len(m.history)
					m.textarea.SetValue("")
				}
				return m, nil
			}

		case tea.KeyCtrlR:
			// Enter search mode
			m.searching = true
			m.searchQuery = ""
			m.searchResults = nil
			m.searchIdx = 0
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.textarea.SetWidth(msg.Width - 4)
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m InputModel) updateSearch(msg tea.KeyMsg) (InputModel, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc, tea.KeyCtrlC:
		m.searching = false
		return m, nil
	case tea.KeyEnter:
		// Accept search result
		if len(m.searchResults) > 0 && m.searchIdx < len(m.searchResults) {
			idx := m.searchResults[m.searchIdx]
			m.textarea.SetValue(m.history[idx])
			m.historyIdx = idx
		}
		m.searching = false
		return m, nil
	case tea.KeyCtrlR:
		// Next search result
		if m.searchIdx < len(m.searchResults)-1 {
			m.searchIdx++
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.updateSearchResults()
		}
		return m, nil
	default:
		if msg.Type == tea.KeyRunes {
			m.searchQuery += string(msg.Runes)
			m.updateSearchResults()
		}
		return m, nil
	}
}

func (m *InputModel) updateSearchResults() {
	m.searchResults = nil
	m.searchIdx = 0
	if m.searchQuery == "" {
		return
	}
	query := strings.ToLower(m.searchQuery)
	// Search backwards through history
	for i := len(m.history) - 1; i >= 0; i-- {
		if strings.Contains(strings.ToLower(m.history[i]), query) {
			m.searchResults = append(m.searchResults, i)
		}
	}
}

// Value returns the current input text
func (m InputModel) Value() string { return filterEscapes(m.textarea.Value()) }

// Reset clears the input
func (m *InputModel) Reset() {
	m.textarea.Reset()
	m.submitted = false
}

// IsSubmitted checks if the user pressed Enter
func (m InputModel) IsSubmitted() bool { return m.submitted }

// SetWidth updates the input width
func (m *InputModel) SetWidth(w int) {
	m.width = w
	m.textarea.SetWidth(w - 4)
}

func (m InputModel) View() string {
	if m.searching {
		searchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
		prompt := searchStyle.Render("(reverse-i-search)") + "`" + m.searchQuery + "`: "
		if len(m.searchResults) > 0 && m.searchIdx < len(m.searchResults) {
			idx := m.searchResults[m.searchIdx]
			return prompt + m.history[idx]
		}
		return prompt
	}

	// Show line count for multi-line
	value := m.textarea.Value()
	lines := strings.Count(value, "\n") + 1

	view := m.textarea.View()
	if lines > 1 {
		lineInfo := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(
			fmt.Sprintf(" (%d lines)", lines))
		view += lineInfo
	}
	return view
}
