package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/howlerops/oculus/pkg/query"
	"github.com/howlerops/oculus/pkg/types"
)

// Styles
var (
	userStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#0ea5e9")).Bold(true)  // sky blue
	assistantStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#22d3ee"))             // cyan
	toolStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b")).Italic(true) // slate
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))              // red
	headerStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#0ea5e9")).Bold(true)   // sky blue
)

// ResponseMsg is sent when the query engine finishes
type ResponseMsg struct {
	Messages []types.Message
	Err      error
}

// StreamTextMsg is sent for streamed text
type StreamTextMsg struct {
	Text string
}

// ToolUseMsg is sent when a tool starts
type ToolUseMsg struct {
	Name string
}

// Model is the bubbletea model for the TUI
type Model struct {
	textarea     textarea.Model
	spinner      spinner.Model
	messages     []types.Message
	streamBuffer string
	history      []string
	historyIdx   int
	engine       *query.Engine
	systemPrompt interface{}
	isLoading    bool
	err          error
	width        int
	height       int
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewModel creates the TUI model
func NewModel(engine *query.Engine, systemPrompt interface{}) Model {
	ta := textarea.New()
	ta.Placeholder = "Type your message..."
	ta.Focus()
	ta.SetHeight(1)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	s := spinner.New()
	s.Spinner = spinner.Dot

	ctx, cancel := context.WithCancel(context.Background())

	return Model{
		textarea:     ta,
		spinner:      s,
		engine:       engine,
		systemPrompt: systemPrompt,
		width:        80,
		height:       24,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			if m.isLoading {
				m.cancel()
				m.isLoading = false
				m.streamBuffer = ""
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyEnter:
			if m.isLoading {
				return m, nil
			}
			input := strings.TrimSpace(m.textarea.Value())
			if input == "" {
				return m, nil
			}

			m.textarea.Reset()
			m.history = append(m.history, input)
			m.historyIdx = len(m.history)

			if input == "/quit" || input == "/exit" {
				return m, tea.Quit
			}

			m.messages = append(m.messages, types.NewUserMessage(input))
			m.isLoading = true
			m.streamBuffer = ""

			ctx, cancel := context.WithCancel(context.Background())
			m.ctx = ctx
			m.cancel = cancel

			return m, m.runQuery(ctx, input)

		case tea.KeyUp:
			if len(m.history) > 0 && m.historyIdx > 0 {
				m.historyIdx--
				m.textarea.SetValue(m.history[m.historyIdx])
			}
			return m, nil

		case tea.KeyDown:
			if m.historyIdx < len(m.history)-1 {
				m.historyIdx++
				m.textarea.SetValue(m.history[m.historyIdx])
			} else {
				m.historyIdx = len(m.history)
				m.textarea.SetValue("")
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)

	case StreamTextMsg:
		m.streamBuffer += msg.Text
		return m, nil

	case ToolUseMsg:
		m.streamBuffer += fmt.Sprintf("\n[Tool: %s]\n", msg.Name)
		return m, nil

	case ResponseMsg:
		m.isLoading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.messages = msg.Messages
		}
		m.streamBuffer = ""
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var sb strings.Builder

	// Header
	sb.WriteString(headerStyle.Render("Oculus"))
	sb.WriteString("\n\n")

	// Messages
	for _, msg := range m.messages {
		switch msg.Kind {
		case "user":
			if msg.User != nil {
				for _, block := range msg.User.Content {
					if block.Type == types.ContentBlockText {
						sb.WriteString(userStyle.Render("You: "))
						sb.WriteString(block.Text)
						sb.WriteString("\n\n")
					}
				}
			}
		case "assistant":
			if msg.Assistant != nil {
				sb.WriteString(assistantStyle.Render("Claude: "))
				for _, block := range msg.Assistant.Content {
					if block.Type == types.ContentBlockText {
						sb.WriteString(block.Text)
					} else if block.Type == types.ContentBlockToolUse {
						sb.WriteString(toolStyle.Render(fmt.Sprintf("[Tool: %s]", block.Name)))
					}
				}
				sb.WriteString("\n\n")
			}
		}
	}

	// Streaming buffer
	if m.streamBuffer != "" {
		sb.WriteString(assistantStyle.Render("Claude: "))
		sb.WriteString(m.streamBuffer)
		sb.WriteString("\n")
	}

	// Loading indicator
	if m.isLoading {
		sb.WriteString(m.spinner.View())
		sb.WriteString(" Thinking...\n")
	}

	// Error
	if m.err != nil {
		sb.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		sb.WriteString("\n")
		m.err = nil
	}

	sb.WriteString("\n")
	sb.WriteString(m.textarea.View())

	return sb.String()
}

func (m Model) runQuery(ctx context.Context, _ string) tea.Cmd {
	return func() tea.Msg {
		handler := &tuiStreamHandler{}
		msgs, err := m.engine.RunQuery(ctx, m.messages, m.systemPrompt, handler)
		return ResponseMsg{Messages: msgs, Err: err}
	}
}

// tuiStreamHandler sends bubbletea messages for streaming
type tuiStreamHandler struct {
	program *tea.Program
}

func (h *tuiStreamHandler) OnText(text string) {
	if h.program != nil {
		h.program.Send(StreamTextMsg{Text: text})
	}
}

func (h *tuiStreamHandler) OnToolUseStart(id, name string) {
	if h.program != nil {
		h.program.Send(ToolUseMsg{Name: name})
	}
}

func (h *tuiStreamHandler) OnToolUseResult(id string, result interface{}) {}
func (h *tuiStreamHandler) OnThinking(text string)                        {}
func (h *tuiStreamHandler) OnComplete(stopReason types.StopReason, usage *types.Usage) {
}
func (h *tuiStreamHandler) OnError(err error) {}

// Run starts the TUI
func Run(engine *query.Engine, systemPrompt interface{}) error {
	m := NewModel(engine, systemPrompt)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
