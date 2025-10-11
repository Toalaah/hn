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
		Title:      lipgloss.NewStyle(),
		Selected:   lipgloss.NewStyle().Foreground(lipgloss.Color("000")).Background(lipgloss.Color("006")),
		Unselected: lipgloss.NewStyle().Foreground(lipgloss.Color("007")),
		SubThread:  lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Bold(true),
		Author:     lipgloss.NewStyle(),
		Meta:       lipgloss.NewStyle().Faint(true),
	}
}
