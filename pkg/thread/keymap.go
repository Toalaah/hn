package thread

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/toalaah/hn/pkg/hn"
)

type KeyMap struct {
	// Navigate up/down in current sub-thread.
	Up   key.Binding
	Down key.Binding

	// Move pager
	PageUp   key.Binding
	PageDown key.Binding

	// Shorthand for navigating to first sub-thread.
	Top key.Binding
	// Shorthand for navigating to last sub-thread.
	Bottom key.Binding

	// Navigate to next/previous sub-thread
	Next key.Binding
	Prev key.Binding

	// Jump back to root of current sub-thread
	Root key.Binding

	// Expand/minimize current sub-thread.
	ExpandCurrent   key.Binding
	CollapseCurrent key.Binding

	// Expand/minimize all sub-threads.
	ExpandAll   key.Binding
	CollapseAll key.Binding

	Quit key.Binding

	// Snap pager to currently selected comment
	ResetView key.Binding
	// Copy current comment text
	Copy key.Binding
}

type threadMetadata struct {
	collapsed bool
	height    int
}

type Model struct {
	KeyMap   KeyMap
	Styles   Styles
	Indent   int
	MaxWidth int

	viewport   viewport.Model
	curRoot    *hn.Thread
	thread     hn.Thread
	meta       []threadMetadata
	lastStatus string
}

// DefaultKeyMap is the default key bindings for the table.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top"),
		),
		Next: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next"),
		),
		Prev: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "prev"),
		),
		Root: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "root"),
		),
		ExpandCurrent: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "expand"),
		),
		CollapseCurrent: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "collapse"),
		),
		ExpandAll: key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "expand all"),
		),
		CollapseAll: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "collapse all"),
		),
		Up: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "down"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		PageDown:  key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "page down")),
		PageUp:    key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "page down")),
		ResetView: key.NewBinding(key.WithKeys("z"), key.WithHelp("z", "reset view")),
		Copy:      key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "copy")),
	}
}
