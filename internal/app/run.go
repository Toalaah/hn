package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/toalaah/hn/pkg/threadview"
)

func Run(t threadview.Thread) error {
	m, err := threadview.New(t,
		threadview.WithHeadSelectable(false),
		threadview.WithHideCollapsedChildren(true),
	)
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run()
	return err
}
