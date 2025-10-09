package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/toalaah/hn/pkg/hn"
	"github.com/toalaah/hn/pkg/thread"
)

func Run(t hn.Thread) error {
	p := tea.NewProgram(thread.NewModel(t), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
