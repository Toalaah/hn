package threadview

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Thread interface {
	tea.Model
	// ID returns a unique identifier for this thread object. It is up to the implementer to ensure that IDs do not collide.
	ID() int
	// Parent returns the thread's parent. If the parent is invalid (e.g nil if the thread for this `Parent` was called is the document root), the second return value should be set to false, otherwise true.
	Parent() (Thread, bool)
	// Children returns a list of this thread's sub-threads i.e children.
	Children() []Thread
}
