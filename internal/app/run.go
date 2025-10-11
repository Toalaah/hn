package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/toalaah/hn/pkg/hn"
	"github.com/toalaah/hn/pkg/thread"
)

func Run(t hn.Thread) error {
	m := thread.New(t,
		thread.WithWidth(80),
	)
	p := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
