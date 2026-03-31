package tui

// ScreenType identifies which screen to show
type ScreenType string

const (
	ScreenChat          ScreenType = "chat"
	ScreenResumeChooser ScreenType = "resume_chooser"
	ScreenSettings      ScreenType = "settings"
	ScreenHelp          ScreenType = "help"
	ScreenOnboarding    ScreenType = "onboarding"
)

// ScreenRouter manages which screen is active
type ScreenRouter struct {
	Current  ScreenType
	Previous ScreenType
	History  []ScreenType
}

// NewScreenRouter creates a ScreenRouter defaulting to the chat screen
func NewScreenRouter() *ScreenRouter {
	return &ScreenRouter{Current: ScreenChat}
}

// Navigate pushes the current screen onto history and switches to the given screen
func (r *ScreenRouter) Navigate(screen ScreenType) {
	r.Previous = r.Current
	r.History = append(r.History, r.Current)
	r.Current = screen
}

// Back pops the most recent screen from history and makes it current
func (r *ScreenRouter) Back() {
	if len(r.History) > 0 {
		r.Current = r.History[len(r.History)-1]
		r.History = r.History[:len(r.History)-1]
	}
}
