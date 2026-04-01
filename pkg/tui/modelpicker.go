package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/howlerops/oculus/pkg/auth"
	"github.com/howlerops/oculus/pkg/config"
)

// ModelPickerState tracks whether the picker is open
type ModelPickerState int

const (
	PickerClosed ModelPickerState = iota
	PickerOpen
)

// ModelEntry is a selectable model in the picker
type ModelEntry struct {
	Provider    string // "Anthropic", "OpenAI", "Ollama", "Claude CLI", etc.
	Name        string // display name
	ID          string // model ID to save
	Description string
	IsHeader    bool // section header, not selectable
}

// LensTarget is which setting we're editing
type LensTarget int

const (
	TargetDefault LensTarget = iota
	TargetFocus
	TargetScan
	TargetCraft
)

func (t LensTarget) String() string {
	switch t {
	case TargetFocus:
		return "Focus"
	case TargetScan:
		return "Scan"
	case TargetCraft:
		return "Craft"
	default:
		return "Default"
	}
}

// ModelPicker is the interactive model selection component
type ModelPicker struct {
	State      ModelPickerState
	Entries    []ModelEntry
	Cursor     int
	Target     LensTarget
	Width      int
	Height     int
	Scroll     int    // scroll offset for long lists
	CurrentID  string // currently active model for highlighting
}

// NewModelPicker creates a new model picker in closed state
func NewModelPicker() *ModelPicker {
	return &ModelPicker{
		State: PickerClosed,
	}
}

// Open activates the picker and loads available models
func (p *ModelPicker) Open(width, height int) {
	p.State = PickerOpen
	p.Width = width
	p.Height = height
	p.Cursor = 0
	p.Scroll = 0
	p.Target = TargetDefault
	p.Entries = p.loadModels()

	// Load current model for highlighting
	settings, _ := config.LoadSettings()
	if settings != nil {
		p.CurrentID = settings.Model
	}

	// Skip to first selectable entry
	for p.Cursor < len(p.Entries) && p.Entries[p.Cursor].IsHeader {
		p.Cursor++
	}
}

// Close deactivates the picker
func (p *ModelPicker) Close() {
	p.State = PickerClosed
}

func (p *ModelPicker) loadModels() []ModelEntry {
	var entries []ModelEntry
	providers := auth.DetectProviders()

	// Anthropic
	entries = append(entries, ModelEntry{IsHeader: true, Provider: "Anthropic", Name: "  Anthropic"})
	entries = append(entries,
		ModelEntry{Provider: "anthropic", Name: "Claude Opus 4", ID: "claude-opus-4-6", Description: "Most capable | $15/$75 per M tokens"},
		ModelEntry{Provider: "anthropic", Name: "Claude Sonnet 4", ID: "claude-sonnet-4-6", Description: "Balanced | $3/$15 per M tokens"},
		ModelEntry{Provider: "anthropic", Name: "Claude Haiku 4", ID: "claude-haiku-4-5-20251001", Description: "Fast & cheap | $0.80/$4 per M tokens"},
	)

	// OpenAI (if available)
	for _, prov := range providers {
		if prov.Name == "OpenAI" && prov.Available {
			entries = append(entries, ModelEntry{IsHeader: true, Provider: "OpenAI", Name: "  OpenAI"})
			entries = append(entries,
				ModelEntry{Provider: "openai", Name: "GPT-4o", ID: "gpt-4o", Description: "Flagship multimodal"},
				ModelEntry{Provider: "openai", Name: "GPT-4o Mini", ID: "gpt-4o-mini", Description: "Fast & affordable"},
				ModelEntry{Provider: "openai", Name: "o1", ID: "o1-preview", Description: "Reasoning model"},
			)
		}
	}

	// Ollama (fetch live)
	for _, prov := range providers {
		if prov.Name == "Ollama" && prov.Available {
			entries = append(entries, ModelEntry{IsHeader: true, Provider: "Ollama", Name: "  Ollama (local)"})
			models := fetchOllamaModelList()
			if len(models) > 0 {
				for _, m := range models {
					entries = append(entries, ModelEntry{Provider: "ollama", Name: m, ID: m, Description: "Local model"})
				}
			} else {
				entries = append(entries, ModelEntry{Provider: "ollama", Name: "(no models - run: ollama pull llama3)", ID: "", Description: "", IsHeader: true})
			}
		}
	}

	// CLI bridges
	for _, prov := range providers {
		if prov.Type == "cli" && prov.Available {
			providerID := ""
			switch prov.Name {
			case "Claude CLI":
				providerID = "claude-code"
			case "Codex CLI":
				providerID = "codex"
			case "Gemini CLI":
				providerID = "gemini-cli"
			}
			entries = append(entries, ModelEntry{IsHeader: true, Provider: prov.Name, Name: "  " + prov.Name})
			entries = append(entries, ModelEntry{
				Provider:    providerID,
				Name:        prov.Name + " (subscription)",
				ID:          providerID,
				Description: prov.Details,
			})
		}
	}

	return entries
}

// Update handles keypresses in the picker
func (p *ModelPicker) Update(msg tea.KeyMsg) (selected bool, result string) {
	switch msg.String() {
	case "esc", "q":
		p.Close()
		return false, ""

	case "up", "k":
		p.Cursor--
		for p.Cursor >= 0 && p.Entries[p.Cursor].IsHeader {
			p.Cursor--
		}
		if p.Cursor < 0 {
			p.Cursor = 0
			for p.Cursor < len(p.Entries) && p.Entries[p.Cursor].IsHeader {
				p.Cursor++
			}
		}
		p.ensureVisible()

	case "down", "j":
		p.Cursor++
		for p.Cursor < len(p.Entries) && p.Entries[p.Cursor].IsHeader {
			p.Cursor++
		}
		if p.Cursor >= len(p.Entries) {
			p.Cursor = len(p.Entries) - 1
			for p.Cursor > 0 && p.Entries[p.Cursor].IsHeader {
				p.Cursor--
			}
		}
		p.ensureVisible()

	case "tab":
		// Cycle through lens targets and update current highlight
		p.Target = (p.Target + 1) % 4
		p.updateCurrentID()

	case " ": // spacebar - select without closing
		if p.Cursor >= 0 && p.Cursor < len(p.Entries) && !p.Entries[p.Cursor].IsHeader {
			entry := p.Entries[p.Cursor]
			if entry.ID != "" {
				result := p.saveSelection(entry)
				p.CurrentID = entry.ID
				return false, result // don't close, just update
			}
		}

	case "enter":
		if p.Cursor >= 0 && p.Cursor < len(p.Entries) && !p.Entries[p.Cursor].IsHeader {
			entry := p.Entries[p.Cursor]
			if entry.ID == "" {
				return false, ""
			}
			result := p.saveSelection(entry)
			p.Close()
			return true, result
		}
	}

	return false, ""
}

func (p *ModelPicker) ensureVisible() {
	maxVisible := p.Height - 8
	if maxVisible < 5 {
		maxVisible = 5
	}
	if p.Cursor < p.Scroll {
		p.Scroll = p.Cursor
	}
	if p.Cursor >= p.Scroll+maxVisible {
		p.Scroll = p.Cursor - maxVisible + 1
	}
}

func (p *ModelPicker) saveSelection(entry ModelEntry) string {
	settings, _ := config.LoadSettings()
	if settings == nil {
		settings = &config.SettingsJson{}
	}

	switch p.Target {
	case TargetDefault:
		settings.Model = entry.ID
	case TargetFocus, TargetScan, TargetCraft:
		if settings.Lenses == nil {
			settings.Lenses = &config.LensSettings{}
		}
		lmc := &config.LensModelConfig{Model: entry.ID, Provider: entry.Provider}
		switch p.Target {
		case TargetFocus:
			settings.Lenses.Focus = lmc
		case TargetScan:
			settings.Lenses.Scan = lmc
		case TargetCraft:
			settings.Lenses.Craft = lmc
		}
	}

	data, _ := json.Marshal(settings)
	os.WriteFile(config.GetSettingsPath(), data, 0o644)
	config.InvalidateSettingsCache()

	return fmt.Sprintf("Model set to: %s (%s) [%s]", entry.Name, entry.Provider, p.Target.String())
}

// View renders the picker
func (p *ModelPicker) View() string {
	if p.State != PickerOpen {
		return ""
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0ea5e9"))
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#22d3ee")).MarginTop(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#e2e8f0"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#0ea5e9")).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b"))
	tabStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b"))
	activeTabStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0ea5e9")).Bold(true).Underline(true)

	var sb strings.Builder

	// Title
	sb.WriteString(titleStyle.Render("Model Picker") + "\n")

	// Tab bar
	tabs := []string{"Default", "Focus", "Scan", "Craft"}
	var tabParts []string
	for i, tab := range tabs {
		if LensTarget(i) == p.Target {
			tabParts = append(tabParts, activeTabStyle.Render("["+tab+"]"))
		} else {
			tabParts = append(tabParts, tabStyle.Render(" "+tab+" "))
		}
	}
	sb.WriteString(strings.Join(tabParts, "  ") + "  (Tab to switch)\n\n")

	// Model list
	maxVisible := p.Height - 8
	if maxVisible < 5 {
		maxVisible = 5
	}
	end := p.Scroll + maxVisible
	if end > len(p.Entries) {
		end = len(p.Entries)
	}

	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true)

	for i := p.Scroll; i < end; i++ {
		entry := p.Entries[i]
		if entry.IsHeader {
			sb.WriteString(headerStyle.Render(entry.Name) + "\n")
			continue
		}

		isActive := entry.ID == p.CurrentID
		cursor := "  "
		style := itemStyle
		if i == p.Cursor {
			cursor = "▸ "
			style = selectedStyle
		}

		line := cursor + style.Render(entry.Name)
		if entry.Description != "" {
			line += " " + descStyle.Render(entry.Description)
		}
		if isActive {
			line += " " + activeStyle.Render("● active")
		}
		sb.WriteString(line + "\n")
	}

	// Scroll indicator
	if len(p.Entries) > maxVisible {
		sb.WriteString(descStyle.Render(fmt.Sprintf("\n  %d/%d models (scroll with up/down)", p.Cursor+1, len(p.Entries))))
	}

	sb.WriteString("\n" + descStyle.Render("  Space: select • Enter: select & close • Tab: switch lens • Esc: close"))

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#0ea5e9")).
		Padding(1, 2).
		Width(p.Width - 4)

	return border.Render(sb.String())
}

// updateCurrentID refreshes the current model ID based on the active target
func (p *ModelPicker) updateCurrentID() {
	settings, _ := config.LoadSettings()
	if settings == nil {
		p.CurrentID = ""
		return
	}
	switch p.Target {
	case TargetDefault:
		p.CurrentID = settings.Model
	case TargetFocus:
		if settings.Lenses != nil && settings.Lenses.Focus != nil {
			p.CurrentID = settings.Lenses.Focus.Model
		} else {
			p.CurrentID = ""
		}
	case TargetScan:
		if settings.Lenses != nil && settings.Lenses.Scan != nil {
			p.CurrentID = settings.Lenses.Scan.Model
		} else {
			p.CurrentID = ""
		}
	case TargetCraft:
		if settings.Lenses != nil && settings.Lenses.Craft != nil {
			p.CurrentID = settings.Lenses.Craft.Model
		} else {
			p.CurrentID = ""
		}
	}
}

func fetchOllamaModelList() []string {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil
	}

	var names []string
	for _, m := range result.Models {
		names = append(names, m.Name)
	}
	return names
}
