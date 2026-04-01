package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/howlerops/oculus/pkg/commands"
	"github.com/howlerops/oculus/pkg/lens"
	"github.com/howlerops/oculus/pkg/query"
	"github.com/howlerops/oculus/pkg/types"
)

// AppState tracks the overall TUI state
type AppState string

const (
	StateChat       AppState = "chat"
	StatePermission AppState = "permission"
	StateLoading    AppState = "loading"
)

// Model is the main bubbletea model wiring all components together
type Model struct {
	// State
	state        AppState
	messages     []types.Message
	streamBuffer string
	err          error
	width        int
	height       int
	loadingStart time.Time

	// Sub-components
	input       InputModel
	viewport    MessageViewport
	markdown    *MarkdownRenderer
	progress    *ToolProgressTracker
	permission  *PermissionDialog
	modelPicker *ModelPicker
	taskPanel   TaskPanel
	statusBar   StatusBar
	contextBar  ContextBar
	spinner     spinner.Model

	// Commands
	CmdRegistry  *commands.Registry
	autocomplete *Autocomplete

	// Backend
	engine       *query.Engine
	lensManager  *lens.Manager
	systemPrompt interface{}

	// Context management
	ctx    context.Context
	cancel context.CancelFunc

	// Program reference for async message sending
	program *tea.Program
}
