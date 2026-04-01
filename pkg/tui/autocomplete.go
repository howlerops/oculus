package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/howlerops/oculus/pkg/commands"
)

// Autocomplete shows a dropdown of matching commands when typing /
type Autocomplete struct {
	Active     bool
	Query      string
	Matches    []*commands.Command
	Selected   int
	Registry   *commands.Registry
	MaxVisible int
}

func NewAutocomplete(registry *commands.Registry) *Autocomplete {
	return &Autocomplete{
		Registry:   registry,
		MaxVisible: 8,
	}
}

// Update refreshes matches based on current input
func (a *Autocomplete) Update(input string) {
	if !strings.HasPrefix(input, "/") || len(input) < 1 {
		a.Active = false
		a.Matches = nil
		return
	}

	query := strings.TrimPrefix(input, "/")
	a.Query = query
	a.Active = true
	a.Selected = 0

	if a.Registry == nil {
		a.Matches = nil
		return
	}

	// Get all commands and filter by prefix
	all := a.Registry.List()
	a.Matches = nil
	for _, cmd := range all {
		if query == "" || strings.HasPrefix(cmd.Name, query) {
			a.Matches = append(a.Matches, cmd)
		}
	}

	// Limit visible
	if len(a.Matches) > a.MaxVisible {
		a.Matches = a.Matches[:a.MaxVisible]
	}

	if a.Selected >= len(a.Matches) {
		a.Selected = 0
	}
}

// SelectNext moves selection down
func (a *Autocomplete) SelectNext() {
	if len(a.Matches) == 0 {
		return
	}
	a.Selected = (a.Selected + 1) % len(a.Matches)
}

// SelectPrev moves selection up
func (a *Autocomplete) SelectPrev() {
	if len(a.Matches) == 0 {
		return
	}
	a.Selected--
	if a.Selected < 0 {
		a.Selected = len(a.Matches) - 1
	}
}

// GetSelected returns the currently selected command name
func (a *Autocomplete) GetSelected() string {
	if len(a.Matches) == 0 || a.Selected >= len(a.Matches) {
		return ""
	}
	return "/" + a.Matches[a.Selected].Name
}

// View renders the autocomplete dropdown
func (a *Autocomplete) View() string {
	if !a.Active || len(a.Matches) == 0 {
		return ""
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#0ea5e9")).
		Padding(0, 1).
		Width(40)

	var lines []string
	for i, cmd := range a.Matches {
		name := "/" + cmd.Name
		desc := cmd.Description
		if len(desc) > 25 {
			desc = desc[:22] + "..."
		}

		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0ea5e9"))
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b"))

		if i == a.Selected {
			nameStyle = nameStyle.Bold(true).Background(lipgloss.Color("#1e293b"))
			descStyle = descStyle.Background(lipgloss.Color("#1e293b"))
		}

		line := nameStyle.Render(name) + " " + descStyle.Render(desc)
		lines = append(lines, line)
	}

	return boxStyle.Render(strings.Join(lines, "\n"))
}

// Dismiss hides the autocomplete
func (a *Autocomplete) Dismiss() {
	a.Active = false
	a.Matches = nil
}
