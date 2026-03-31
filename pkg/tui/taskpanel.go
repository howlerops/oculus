package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TaskItem represents a task in the panel
type TaskItem struct {
	ID          string
	Subject     string
	Status      string // "pending", "in_progress", "completed"
	Description string
}

// TaskPanel displays the current task list
type TaskPanel struct {
	Tasks    []TaskItem
	Visible  bool
	Width    int
	Selected int
}

func NewTaskPanel() TaskPanel {
	return TaskPanel{Visible: true, Width: 40}
}

func (p *TaskPanel) AddTask(task TaskItem) {
	p.Tasks = append(p.Tasks, task)
}

func (p *TaskPanel) UpdateTask(id, status string) {
	for i := range p.Tasks {
		if p.Tasks[i].ID == id {
			p.Tasks[i].Status = status
		}
	}
}

func (p *TaskPanel) RemoveTask(id string) {
	for i, t := range p.Tasks {
		if t.ID == id {
			p.Tasks = append(p.Tasks[:i], p.Tasks[i+1:]...)
			return
		}
	}
}

func (p *TaskPanel) Toggle() { p.Visible = !p.Visible }

func (p TaskPanel) CompletedCount() int {
	n := 0
	for _, t := range p.Tasks {
		if t.Status == "completed" {
			n++
		}
	}
	return n
}

func (p TaskPanel) View() string {
	if !p.Visible || len(p.Tasks) == 0 {
		return ""
	}

	var sb strings.Builder

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(0, 1).
		Width(p.Width)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	sb.WriteString(titleStyle.Render(fmt.Sprintf("Tasks (%d/%d)", p.CompletedCount(), len(p.Tasks))))
	sb.WriteString("\n")

	for _, task := range p.Tasks {
		var icon string
		var style lipgloss.Style
		switch task.Status {
		case "completed":
			icon = "✓"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Strikethrough(true)
		case "in_progress":
			icon = "▸"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
		default:
			icon = "○"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
		}

		subject := task.Subject
		maxLen := p.Width - 6
		if len(subject) > maxLen {
			subject = subject[:maxLen-3] + "..."
		}

		sb.WriteString(fmt.Sprintf(" %s %s\n", icon, style.Render(subject)))
	}

	return borderStyle.Render(sb.String())
}
