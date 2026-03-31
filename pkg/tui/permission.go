package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PermissionRequest represents a tool asking for permission
type PermissionRequest struct {
	ToolName    string
	Description string
	Input       string
	IsReadOnly  bool
}

// PermissionResponse is the user's decision
type PermissionResponse struct {
	Decision string // "allow", "deny", "always_allow", "always_deny"
	ToolName string
}

// PermissionDialog is a bubbletea component for tool permission prompts
type PermissionDialog struct {
	Request  PermissionRequest
	Response chan PermissionResponse
	selected int
	options  []permOption
	width    int
}

type permOption struct {
	key   string
	label string
	value string
}

func NewPermissionDialog(req PermissionRequest) PermissionDialog {
	return PermissionDialog{
		Request:  req,
		Response: make(chan PermissionResponse, 1),
		options: []permOption{
			{key: "y", label: "Yes, allow once", value: "allow"},
			{key: "n", label: "No, deny", value: "deny"},
			{key: "a", label: "Always allow this tool", value: "always_allow"},
			{key: "d", label: "Always deny this tool", value: "always_deny"},
		},
		width: 60,
	}
}

func (d PermissionDialog) Init() tea.Cmd { return nil }

func (d PermissionDialog) Update(msg tea.Msg) (PermissionDialog, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			d.respond("allow")
		case "enter":
			d.respond(d.options[d.selected].value)
		case "n", "N":
			d.respond("deny")
		case "a", "A":
			d.respond("always_allow")
		case "d", "D":
			d.respond("always_deny")
		case "up", "k":
			if d.selected > 0 {
				d.selected--
			}
		case "down", "j":
			if d.selected < len(d.options)-1 {
				d.selected++
			}
		case "esc":
			d.respond("deny")
		}
	case tea.WindowSizeMsg:
		d.width = msg.Width
	}
	return d, nil
}

func (d *PermissionDialog) respond(decision string) {
	select {
	case d.Response <- PermissionResponse{Decision: decision, ToolName: d.Request.ToolName}:
	default:
	}
}

func (d PermissionDialog) View() string {
	var sb strings.Builder

	// Border style
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("3")).
		Padding(1, 2).
		Width(min(d.width-4, 70))

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
	sb.WriteString(titleStyle.Render("Permission Required"))
	sb.WriteString("\n\n")

	// Tool name
	toolNameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	sb.WriteString(toolNameStyle.Render(d.Request.ToolName))
	if d.Request.IsReadOnly {
		roStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
		sb.WriteString(" " + roStyle.Render("(read-only)"))
	}
	sb.WriteString("\n")

	// Description
	if d.Request.Description != "" {
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
		sb.WriteString(descStyle.Render(d.Request.Description))
		sb.WriteString("\n")
	}

	// Input preview
	if d.Request.Input != "" {
		input := d.Request.Input
		if len(input) > 200 {
			input = input[:197] + "..."
		}
		inputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		sb.WriteString("\n" + inputStyle.Render(input) + "\n")
	}

	sb.WriteString("\n")

	// Options
	for i, opt := range d.options {
		cursor := "  "
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
		if i == d.selected {
			cursor = "▸ "
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		}
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
		sb.WriteString(cursor + keyStyle.Render("["+opt.key+"]") + " " + style.Render(opt.label) + "\n")
	}

	return borderStyle.Render(sb.String())
}
