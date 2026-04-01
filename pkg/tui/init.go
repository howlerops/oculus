package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/howlerops/oculus/pkg/lens"
	"github.com/howlerops/oculus/pkg/query"
)

// NewModel creates the full TUI model with all components wired
func NewModel(engine *query.Engine, lensManager *lens.Manager, systemPrompt interface{}) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ctx, cancel := context.WithCancel(context.Background())

	md := NewMarkdownRenderer(80)

	return Model{
		state:        StateChat,
		input:        NewInputModel(),
		viewport:     NewMessageViewport(80, 20),
		markdown:     md,
		progress:     NewToolProgressTracker(),
		taskPanel:    NewTaskPanel(),
		statusBar:    NewStatusBar(80),
		contextBar:   NewContextBar(200000, 80),
		spinner:      s,
		engine:       engine,
		lensManager:  lensManager,
		systemPrompt: systemPrompt,
		width:        80,
		height:       24,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Init returns the initial commands
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.input.Init(),
		m.spinner.Tick,
	)
}

// Run starts the TUI with the given engine and system prompt
func Run(engine *query.Engine, lensManager *lens.Manager, systemPrompt interface{}) error {
	m := NewModel(engine, lensManager, systemPrompt)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	m.program = p
	_, err := p.Run()
	return err
}
