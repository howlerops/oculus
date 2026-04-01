package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/howlerops/oculus/pkg/orchestration"
	"github.com/howlerops/oculus/pkg/types"
)

// Update handles all bubbletea messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle permission dialog when active
	if m.state == StatePermission && m.permission != nil {
		return m.updatePermission(msg)
	}

	switch msg := msg.(type) {

	// Window resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.SetWidth(msg.Width)
		viewH := msg.Height - 8
		if viewH < 5 { viewH = 5 }
		m.viewport.viewport.Width = msg.Width
		m.viewport.viewport.Height = viewH
		m.viewport.width = msg.Width
		m.viewport.height = viewH
		m.statusBar.Width = msg.Width
		m.contextBar.Width = msg.Width
		m.markdown = NewMarkdownRenderer(msg.Width - 4)
		m.taskPanel.Width = msg.Width / 3
		return m, nil

	// Keyboard input
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			if m.state == StateLoading {
				m.cancel()
				m.state = StateChat
				m.streamBuffer = ""
				m.progress.Clear()
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyEnter:
			if m.state == StateLoading {
				return m, nil
			}
			if !m.input.IsSubmitted() {
				// Let input handle it
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				if m.input.IsSubmitted() {
					return m.submitInput()
				}
				return m, cmd
			}

		case tea.KeyTab:
			// Accept autocomplete selection
			if m.autocomplete != nil && m.autocomplete.Active {
				selected := m.autocomplete.GetSelected()
				if selected != "" {
					m.input.textarea.SetValue(selected + " ")
					m.input.textarea.CursorEnd()
					m.autocomplete.Dismiss()
				}
				return m, nil
			}

		case tea.KeyPgUp, tea.KeyPgDown:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

		// Forward to input
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		if m.input.IsSubmitted() {
			if m.autocomplete != nil {
				m.autocomplete.Dismiss()
			}
			return m.submitInput()
		}
		cmds = append(cmds, cmd)

		// Update autocomplete dropdown
		if m.autocomplete != nil {
			m.autocomplete.Update(m.input.Value())
		}

	// Streaming text from assistant
	case StreamTextMsg:
		m.streamBuffer += msg.Text
		// Render markdown incrementally
		if m.markdown != nil {
			rendered, _ := m.markdown.RenderStreaming(m.streamBuffer)
			m.viewport.UpdateStreamingContent(rendered)
		}
		return m, nil

	// Thinking block
	case StreamThinkingMsg:
		// Could show thinking indicator
		return m, nil

	// Tool execution started
	case ToolStartMsg:
		m.progress.StartTool(msg.ToolID, msg.ToolName)
		display := ToolCallDisplay{
			ToolName: msg.ToolName,
			Width:    m.width,
		}
		m.viewport.AppendRaw(display.View())
		return m, nil

	// Tool execution completed
	case ToolResultMsg:
		m.progress.CompleteTool(msg.ToolID)
		if msg.Result != "" {
			display := ToolCallDisplay{
				ToolName: msg.ToolID,
				Result:   msg.Result,
				IsError:  msg.IsError,
				Width:    m.width,
			}
			m.viewport.AppendRaw(display.View())
		}
		return m, nil

	// Query engine finished
	case ResponseMsg:
		m.state = StateChat
		m.progress.Clear()
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.messages = msg.Messages
			// Render the last assistant message with markdown
			m.renderLastAssistantMessage()
		}
		m.streamBuffer = ""

		// Update status bar with token counts
		totalIn, totalOut := m.countTokens()
		m.statusBar.InputTokens = totalIn
		m.statusBar.OutputTokens = totalOut

		return m, nil

	// Permission request from tool
	case PermissionRequestMsg:
		dialog := NewPermissionDialog(msg.Request)
		m.permission = &dialog
		m.state = StatePermission
		return m, nil

	// Task updates
	case TaskUpdateMsg:
		m.taskPanel.Tasks = msg.Tasks
		return m, nil

	// Error
	case ErrorMsg:
		m.err = msg.Err
		return m, nil

	// Spinner tick
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.progress.Tick()
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// submitInput processes the user's input
func (m Model) submitInput() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.input.Value())
	if input == "" {
		return m, nil
	}

	// Parse slash commands
	if strings.HasPrefix(input, "/") {
		return m.handleCommand(input)
	}

	m.input.Reset()

	// Add user message to viewport
	m.viewport.AppendMessage("user", input)
	m.messages = append(m.messages, types.NewUserMessage(input))

	// Start loading
	m.state = StateLoading
	m.streamBuffer = ""

	ctx, cancel := context.WithCancel(context.Background())
	m.ctx = ctx
	m.cancel = cancel

	// Launch query in background
	return m, m.runQuery(ctx)
}

// updatePermission handles input during permission dialog
func (m Model) updatePermission(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.permission == nil {
		m.state = StateChat
		return m, nil
	}

	updated, cmd := m.permission.Update(msg)
	m.permission = &updated

	// Check if response was sent
	select {
	case resp := <-m.permission.Response:
		m.state = StateChat
		m.permission = nil
		_ = resp // Would send back to permission system
		return m, cmd
	default:
	}

	return m, cmd
}

// runQuery launches the query engine in a goroutine.
// Uses lensManager for intelligent routing when available, falling back to engine.
func (m Model) runQuery(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		handler := &TUIStreamHandler{Program: m.program}
		var msgs []types.Message
		var err error
		if m.lensManager != nil {
			msgs, err = m.lensManager.RunQuery(ctx, m.messages, m.systemPrompt, handler)
		} else {
			msgs, err = m.engine.RunQuery(ctx, m.messages, m.systemPrompt, handler)
		}
		return ResponseMsg{Messages: msgs, Err: err}
	}
}

// renderLastAssistantMessage renders the last assistant message with markdown
func (m *Model) renderLastAssistantMessage() {
	for i := len(m.messages) - 1; i >= 0; i-- {
		msg := m.messages[i]
		if msg.Kind == "assistant" && msg.Assistant != nil {
			var text string
			for _, block := range msg.Assistant.Content {
				if block.Type == types.ContentBlockText {
					text += block.Text
				}
			}
			if text != "" && m.markdown != nil {
				rendered := m.markdown.Render(text)
				m.viewport.AppendMessage("assistant", rendered)
			}
			return
		}
	}
}

// countTokens totals input/output tokens from all messages
func (m Model) countTokens() (int, int) {
	var totalIn, totalOut int
	for _, msg := range m.messages {
		if msg.Kind == "assistant" && msg.Assistant != nil && msg.Assistant.Usage != nil {
			totalIn += msg.Assistant.Usage.InputTokens
			totalOut += msg.Assistant.Usage.OutputTokens
		}
	}
	return totalIn, totalOut
}

// handleCommand dispatches slash commands dynamically
func (m Model) handleCommand(input string) (tea.Model, tea.Cmd) {
	m.input.Reset()

	// Add system message so splash screen clears
	m.messages = append(m.messages, types.NewSystemMessage(types.SystemMsgInformational, input))

	parts := strings.SplitN(strings.TrimPrefix(input, "/"), " ", 2)
	cmdName := parts[0]
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	// Built-in TUI-specific commands (need special handling)
	switch cmdName {
	case "quit", "exit", "q":
		return m, tea.Quit
	case "ralph":
		m.state = StateLoading
		m.viewport.AppendMessage("system", "Starting Ralph loop: "+args)
		return m, m.startRalph(args)
	case "plan":
		m.state = StateLoading
		m.viewport.AppendMessage("system", "Starting consensus planning: "+args)
		return m, m.startPlan(args)
	}

	// Try the command registry (all 56+ commands)
	if m.CmdRegistry != nil {
		if cmd := m.CmdRegistry.Find(cmdName); cmd != nil {
			output, continueConv, err := cmd.Run(m.ctx, args)
			if err != nil {
				m.viewport.AppendMessage("system", "Error: "+err.Error())
			} else if output != "" {
				m.viewport.AppendMessage("system", output)
			}
			if continueConv {
				// Send as user message to the LLM
				m.messages = append(m.messages, types.NewUserMessage(output))
				m.state = StateLoading
				return m, m.runQuery(m.ctx)
			}
			return m, nil
		}
	}

	m.viewport.AppendMessage("system", "Unknown command: /"+cmdName+"\nType /help for available commands")
	return m, nil
}

func (m Model) startRalph(task string) tea.Cmd {
	return func() tea.Msg {
		cfg := orchestration.RalphConfig{Task: task}
		err := orchestration.RalphLoop(m.ctx, cfg, m.lensManager)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ResponseMsg{Messages: m.messages}
	}
}

func (m Model) startPlan(task string) tea.Cmd {
	return func() tea.Msg {
		result, err := orchestration.PlanConsensus(m.ctx, task, m.lensManager)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		status := "Consensus reached"
		if !result.Converged {
			status = fmt.Sprintf("No consensus after %d rounds", result.Rounds)
		}
		return StreamTextMsg{Text: fmt.Sprintf("\n%s:\n\n%s", status, result.FinalPlan)}
	}
}
