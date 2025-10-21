package threadview

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	// Navigate up/down in current sub-thread.
	Up   key.Binding
	Down key.Binding
	// Move pager.
	PageUp   key.Binding
	PageDown key.Binding
	// Shorthand for navigating to first sub-thread.
	Top key.Binding
	// Shorthand for navigating to last sub-thread.
	Bottom key.Binding
	// Navigate to next/previous thread.
	Next key.Binding
	Prev key.Binding
	// Jump back to root of current thread.
	Root key.Binding
	// Expand/minimize current thread.
	ToggleFold key.Binding
	// Snap viewport to currently selected thread.
	ResetView key.Binding
	// Issue copy command to current thread.
	Copy key.Binding
	// Quit out of view.
	Quit key.Binding
}

// DefaultKeyMap returns the default key bindings for a new threadview model.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:         key.NewBinding(key.WithKeys("k")),
		Down:       key.NewBinding(key.WithKeys("j")),
		PageUp:     key.NewBinding(key.WithKeys("u", "ctrl+u")),
		PageDown:   key.NewBinding(key.WithKeys("d", "ctrl+d")),
		Top:        key.NewBinding(key.WithKeys("g")),
		Bottom:     key.NewBinding(key.WithKeys("G")),
		Next:       key.NewBinding(key.WithKeys("n")),
		Prev:       key.NewBinding(key.WithKeys("p")),
		Root:       key.NewBinding(key.WithKeys("r")),
		ToggleFold: key.NewBinding(key.WithKeys("tab")),
		ResetView:  key.NewBinding(key.WithKeys("z")),
		Copy:       key.NewBinding(key.WithKeys("y")),
		Quit:       key.NewBinding(key.WithKeys("h", "q")),
	}
}
