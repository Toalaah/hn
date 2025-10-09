package thread

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Title      lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	SubThread  lipgloss.Style
	Author     lipgloss.Style
	Meta       lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Underline(true),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("#828282")),
		SubThread:  lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D1D1")),
		Author:     lipgloss.NewStyle().Bold(true),
		Meta:       lipgloss.NewStyle().Faint(true),
	}
}
