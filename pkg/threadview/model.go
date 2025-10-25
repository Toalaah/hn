package threadview

import (
	"cmp"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	KeyMap KeyMap

	viewport              viewport.Model
	curRoot               Thread
	head                  Thread
	lastStatus            string
	headSelectable        bool
	hideCollapsedChildren bool
	numNodes              int
	meta                  []metadata
}

type metadata struct {
	node      Thread
	collapsed bool
	visible   bool
	height    int
}

// DisplayStateMsg is used to inform a node of its current state. The node is always notified of its state prior to calling its `View()` method.
type DisplayStateMsg struct {
	// Whether this node is currently marked as being collapsed/minimized.
	Collapsed bool
	// Whether this node is currently selected i.e the "active" node.
	Selected bool
	// Whether this node is a descendent of the currently selected node.
	Subthread bool
	// The depth of this node, from the thread root.
	Depth int
	// The current width of the threadview's viewport.
	Width int
}

func New(t Thread, opts ...Option) (*Model, error) {
	if t == nil {
		return nil, errors.New("thread is nil")
	}
	n := NumNodes(t)
	m := &Model{
		KeyMap:         DefaultKeyMap(),
		viewport:       viewport.New(0, 0),
		curRoot:        t,
		head:           t,
		headSelectable: true,
		numNodes:       n,
		meta:           make([]metadata, n),
	}

	i := 0
	dfs(nil, m.head, func(root, cur Thread) {
		m.meta[i] = metadata{
			visible:   true,
			node:      cur,
			collapsed: false,
			height:    0,
		}
		i++
	})

	for _, opt := range opts {
		opt(m)
	}

	if !m.headSelectable {
		children := t.Children()
		if len(children) > 0 {
			m.curRoot = children[0]
		} else {
			m.curRoot = nil
		}
	}

	return m, nil
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	// Handle own updates.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleInput(msg)
	case tea.MouseMsg:
		if tea.MouseEvent(msg).Button == tea.MouseButtonLeft {
			m.curRoot = m.getThreadFromYPos(msg.Y)
		}
	case tea.WindowSizeMsg:
		padding := 1
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - padding
		m.viewport.YPosition = 0
	case CopyTextResultMsg:
		if msg.Error == nil {
			m.lastStatus = "Contents copied to clipboard"
		} else {
			m.lastStatus = fmt.Sprintf("Failed to copy to clipboard: %s", msg.Error.Error())
		}
		return m, ClearStatusAfter(1250 * time.Millisecond)
	case ClearStatusMsg:
		m.lastStatus = ""
	}
	// Propagate messages.
	v, cmd := m.viewport.Update(msg)
	m.viewport = v
	cmds = append(cmds, cmd)
	dfs(nil, m.head, func(root, cur Thread) {
		_, cmd := cur.Update(msg)
		cmds = append(cmds, cmd)
	})
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	threads := m.threadView(m.head, DisplayStateMsg{
		Collapsed: m.meta[m.threadIndex(m.head)].collapsed,
		Selected:  m.curRoot == m.head,
		Subthread: false,
		Depth:     0,
		Width:     m.viewport.Width,
	})
	threads = strings.Trim(threads, "\n")
	footer := cmp.Or(m.lastStatus, m.defaultStatus())
	m.viewport.SetContent(threads)
	return strings.Join([]string{m.viewport.View(), footer}, "\n")
}

func (m *Model) threadView(t Thread, state DisplayStateMsg) string {
	var b strings.Builder

	idx := m.threadIndex(t)
	if !m.meta[idx].visible {
		return ""
	}

	t.Update(state)
	s := t.View()
	if s == "" {
		m.meta[idx].height = 0
	} else {
		b.WriteString(s)
		m.meta[idx].height = lipgloss.Height(s)
	}
	b.WriteString("\n")

	state.Depth++
	state.Subthread = cmp.Or(state.Subthread, state.Selected)
	visible := !(m.hideCollapsedChildren && state.Collapsed)
	children := t.Children()
	for _, c := range children {
		cidx := m.threadIndex(c)
		m.meta[cidx].visible = visible
		state.Selected = c == m.curRoot
		state.Collapsed = m.meta[cidx].collapsed
		s := m.threadView(c, state)

		b.WriteString(s)
	}

	return b.String()
}

func (m *Model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch {
	case key.Matches(msg, m.KeyMap.Up):
		if p, ok := m.curRoot.Parent(); ok && !(p == m.head && !m.headSelectable) {
			m.curRoot = p
		}
	case key.Matches(msg, m.KeyMap.Down):
		if c := m.curRoot.Children(); len(c) > 0 {
			m.curRoot = c[0]
		}
	case key.Matches(msg, m.KeyMap.PageUp):
		m.viewport.ScrollUp(m.viewport.Height / 4)
	case key.Matches(msg, m.KeyMap.PageDown):
		m.viewport.ScrollDown(m.viewport.Height / 4)
	case key.Matches(msg, m.KeyMap.Top):
		m.curRoot = m.head.Children()[0]
	case key.Matches(msg, m.KeyMap.Bottom):
		threads := m.head.Children()
		m.curRoot = threads[len(threads)-1]
	case key.Matches(msg, m.KeyMap.Next):
		m.nextThread()
	case key.Matches(msg, m.KeyMap.Prev):
		m.prevThread()
	case key.Matches(msg, m.KeyMap.Root):
		for p, ok := m.curRoot.Parent(); ok && !(p == m.head && !m.headSelectable); p, ok = m.curRoot.Parent() {
			m.curRoot = p
		}
	case key.Matches(msg, m.KeyMap.ToggleFold):
		i := m.threadIndex(m.curRoot)
		c := m.meta[i].collapsed
		m.meta[i].collapsed = !c
		if m.hideCollapsedChildren {
			dfs(nil, m.curRoot, func(root, cur Thread) {
				m.meta[m.threadIndex(cur)].visible = c
			})
		}
	case key.Matches(msg, m.KeyMap.ResetView):
		m.seekToCurrentRoot()
	case key.Matches(msg, m.KeyMap.Copy):
		_, cmd := m.curRoot.Update(CopyTextMsg{})
		cmds = append(cmds, cmd)
	case key.Matches(msg, m.KeyMap.Quit):
		return m, tea.Quit
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) getYOffsetForThread(t Thread) int {
	y := 0
	for i := range m.threadIndex(t) {
		y += m.meta[i].height
	}
	return y
}

func (m *Model) getThreadFromYPos(y int) Thread {
	y = y + m.viewport.YOffset
	lo, hi := 0, 0
	for _, m := range m.meta {
		if !m.visible {
			continue
		}
		hi += m.height
		if lo <= y && y < hi {
			return m.node
		}
		lo = hi
	}
	// Should never happen
	return m.curRoot
}

func (m *Model) seekToCurrentRoot() {
	m.viewport.SetYOffset(m.getYOffsetForThread(m.curRoot))
}

func (m *Model) getParentOrTopThread(t Thread) Thread {
	p, ok := t.Parent()
	if !ok {
		return t
	}
	return p
}

func (m *Model) navigateSubThread(n int) {
	p := m.getParentOrTopThread(m.curRoot)
	children := p.Children()
	i := -1
	for j := range children {
		if children[j] == m.curRoot {
			i = j
			break
		}
	}
	if i == -1 {
		return
	}
	l := len(children)
	m.curRoot = children[clamp(i+n, 0, l-1)]
}

func (m *Model) nextThread() { m.navigateSubThread(1) }
func (m *Model) prevThread() { m.navigateSubThread(-1) }

type ClearStatusMsg struct{}
type CopyTextMsg struct{}
type CopyTextResultMsg struct{ Error error }

func ClearStatusAfter(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return ClearStatusMsg{}
	}
}

func (m *Model) defaultStatus() string {
	numNodes := m.numNodes
	if !m.headSelectable {
		numNodes--
	}
	right := fmt.Sprintf("(%d rows) #%d %d/%d %d%%",
		m.meta[m.threadIndex(m.curRoot)].height,
		m.curRoot.ID(),
		m.threadIndex(m.curRoot),
		numNodes, // Head is "virtual"
		int(m.viewport.ScrollPercent()*100),
	)
	left := fmt.Sprintf(
		"%s:Next %s:Prev %s:Down %s:Up %s:Fold %s:Quit",
		strings.Join(m.KeyMap.Next.Keys(), ","),
		strings.Join(m.KeyMap.Prev.Keys(), ","),
		strings.Join(m.KeyMap.Down.Keys(), ","),
		strings.Join(m.KeyMap.Up.Keys(), ","),
		strings.Join(m.KeyMap.ToggleFold.Keys(), ","),
		strings.Join(m.KeyMap.Quit.Keys(), ","),
	)
	padding := max(m.viewport.Width-lipgloss.Width(right)-lipgloss.Width(left), 0)
	return lipgloss.NewStyle().Background(lipgloss.Color("0")).Render(left + strings.Repeat(" ", padding) + right)
}

func (m *Model) threadIndex(t Thread) int {
	for i := range m.meta {
		node := m.meta[i].node
		if node != nil && node == t {
			return i
		}
	}
	panic("could not determine thread index")
}

type Option func(*Model)

func WithKeys(k KeyMap) Option {
	return func(m *Model) {
		m.KeyMap = k
	}
}

func WithHeadSelectable(b bool) Option {
	return func(m *Model) {
		m.headSelectable = b
	}
}

func WithHideCollapsedChildren(b bool) Option {
	return func(m *Model) {
		m.hideCollapsedChildren = b
	}
}
