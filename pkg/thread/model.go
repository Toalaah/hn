package thread

import (
	"cmp"
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mergestat/timediff"
	"github.com/toalaah/hn/pkg/hn"
)

type ThreadI interface {
	ID() int
	Text() string
	Parent() *ThreadI
	Children() []ThreadI
}

func NewModel(thread hn.Thread) *Model {
	n := thread.NumComments()
	m := &Model{
		KeyMap:   DefaultKeyMap(),
		Styles:   DefaultStyles(),
		Indent:   2,
		MaxWidth: 72,

		curRoot:    nil,
		thread:     thread,
		viewport:   viewport.New(0, 0),
		meta:       make([]threadMetadata, n),
		lastStatus: "",
	}
	if len(m.thread.Children) > 0 {
		m.curRoot = &m.thread.Children[0]
	}
	return m
}

var lookup = make(map[int]int)

func (m *Model) toIndex(t *hn.Thread) int {
	if i, ok := lookup[t.ID]; ok {
		return i
	}
	i, _ := m.thread.CommentIndex(t)
	i--
	lookup[t.ID] = i
	return i
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) textWidth() int {
	if m.MaxWidth == 0 {
		return m.viewport.Width
	}
	return min(m.MaxWidth, m.viewport.Width)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleInput(msg)
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		padding := 2
		m.viewport.Height = msg.Height - padding
		m.viewport.YPosition = 1
	}
	return m, nil
}

func (m *Model) viewHeader() string {
	var b strings.Builder
	b.WriteString(m.Styles.Title.Render(m.thread.Title))
	urlParts := strings.Split(m.thread.URL.Hostname(), ".")
	domain := urlParts[len(urlParts)-2] + "." + urlParts[len(urlParts)-1]
	b.WriteString(" " + m.Styles.Meta.Render(fmt.Sprintf("(%s)", domain)))
	return b.String()
}

func (m *Model) View() string {
	header := m.viewHeader()
	footer := cmp.Or(m.lastStatus, m.defaultStatus())
	m.viewport.SetContent(m.view())
	components := []string{
		header,
		m.viewport.View(),
		footer,
	}
	return strings.Join(components, "\n")
}

func (m *Model) threadView(t *hn.Thread, selected bool, isThreadChild bool, indent int) string {
	var b strings.Builder

	idx := m.toIndex(t)
	isCollapsed := m.meta[idx].collapsed
	// base style
	sty := lipgloss.NewStyle().MarginLeft(indent)
	w := m.textWidth()
	if selected {
		sty = sty.Inherit(m.Styles.Selected)
		isThreadChild = true
	} else if isThreadChild {
		sty = sty.Inherit(m.Styles.SubThread)
	} else {
		sty = sty.Inherit(m.Styles.Unselected)
	}

	// Always write author & title.
	b.WriteString(m.Styles.Author.Inherit(sty).MarginLeft(indent).Render(t.Author))
	b.WriteString(" " + m.Styles.Meta.Render(fmt.Sprintf("(%s)", timediff.TimeDiff(t.Date))))

	// Handle '[n more] / [-]'.
	if isCollapsed {
		b.WriteString(m.Styles.Meta.Render(fmt.Sprintf(" [%d more]", t.NumComments()+1)))
	} else {
		b.WriteString(m.Styles.Meta.Render(" [-]"))
	}
	b.WriteString("\n")

	// If comment is collapsed, we are done.
	if isCollapsed {
		m.meta[idx].height = lipgloss.Height(b.String())
		return b.String()
	}

	b.WriteString(sty.Width(w).Render(t.Text))
	m.meta[idx].height = lipgloss.Height(b.String())
	b.WriteString("\n")

	// Generate view for any futher children.
	if len(t.Children) > 0 {
		for i := range t.Children {
			child := &t.Children[i]
			b.WriteString(m.threadView(child, m.curRoot == child, isThreadChild, indent+m.Indent))
		}
	}

	return b.String()
}

func (m *Model) view() string {
	var b strings.Builder
	for i := range m.thread.Children {
		c := &m.thread.Children[i]
		b.WriteString(m.threadView(c, m.curRoot == c, m.curRoot == c, 0))
	}
	return b.String()
}

func (m *Model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.KeyMap.Up):
		if !m.meta[m.toIndex(m.curRoot)].collapsed {
			if m.curRoot.Parent.Parent != nil {
				m.curRoot = cmp.Or(m.curRoot.Parent, &m.thread.Children[0])
				m.seekToCurrentRoot()
			}
		}
	case key.Matches(msg, m.KeyMap.Down):
		if !m.meta[m.toIndex(m.curRoot)].collapsed && len(m.curRoot.Children) > 0 {
			m.curRoot = &m.curRoot.Children[0]
			m.seekToCurrentRoot()
		}
	case key.Matches(msg, m.KeyMap.PageUp):
		m.viewport.ScrollUp(m.viewport.Height / 4)
	case key.Matches(msg, m.KeyMap.PageDown):
		m.viewport.ScrollDown(m.viewport.Height / 4)
	case key.Matches(msg, m.KeyMap.Top):
		if len(m.thread.Children) > 0 {
			m.curRoot = &m.thread.Children[0]
			m.seekToCurrentRoot()
		}
	case key.Matches(msg, m.KeyMap.Bottom):
		if l := len(m.thread.Children); l > 0 {
			m.curRoot = &m.thread.Children[l-1]
			m.seekToCurrentRoot()
		}
	case key.Matches(msg, m.KeyMap.Next):
		m.nextThread()
		m.seekToCurrentRoot()
	case key.Matches(msg, m.KeyMap.Prev):
		m.prevThread()
		m.seekToCurrentRoot()
	case key.Matches(msg, m.KeyMap.Root):
		for m.curRoot.Parent.Parent != nil {
			m.curRoot = m.curRoot.Parent
			m.seekToCurrentRoot()
		}
	case key.Matches(msg, m.KeyMap.ExpandCurrent):
		m.meta[m.toIndex(m.curRoot)].collapsed = false
		m.seekToCurrentRoot()
	case key.Matches(msg, m.KeyMap.CollapseCurrent):
		m.meta[m.toIndex(m.curRoot)].collapsed = true
	case key.Matches(msg, m.KeyMap.ExpandAll):
		for i := range m.meta {
			m.meta[i].collapsed = false
		}
		m.seekToCurrentRoot()
	case key.Matches(msg, m.KeyMap.CollapseAll):
		for i := range m.thread.Children {
			m.meta[m.toIndex(&m.thread.Children[i])].collapsed = true
		}
	case key.Matches(msg, m.KeyMap.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.KeyMap.ResetView):
		m.seekToCurrentRoot()
	case key.Matches(msg, m.KeyMap.Copy):
		go func() {
			if m.curRoot != nil {
				clipboard.WriteAll(m.curRoot.Text)
			}
		}()
		m.setEphemeralStatus("Contents copied to clipboard", time.Second)
	}
	return m, nil
}

func (m *Model) seekToCurrentRoot() {
	idx := m.toIndex(m.curRoot)
	y := 0
	for i := range idx {
		y += m.meta[i].height
	}
	m.viewport.SetYOffset(y)
}

func (m *Model) navigateSubThread(n int) {
	p := cmp.Or(m.curRoot.Parent, &m.thread)
	i := -1
	for j := range p.Children {
		if &p.Children[j] == m.curRoot {
			i = j
			break
		}
	}
	if i == -1 {
		panic("could not find self, should not happen")
	}
	l := len(p.Children)
	m.curRoot = &p.Children[clamp(i+n, 0, l-1)]
}

func (m *Model) nextThread() { m.navigateSubThread(1) }
func (m *Model) prevThread() { m.navigateSubThread(-1) }

func (m *Model) setEphemeralStatus(s string, d time.Duration) {
	m.lastStatus = s
	go func() {
		time.Sleep(d)
		m.lastStatus = ""
	}()
}

func (m *Model) defaultStatus() string {
	numComments := m.thread.NumComments()
	right := fmt.Sprintf("%d/%d %d%%", m.toIndex(m.curRoot)+1, numComments, int(m.viewport.ScrollPercent()*100))
	padding := m.viewport.Width - lipgloss.Width(right)
	left := lipgloss.NewStyle().Width(padding).Render("n:next p:prev j:down k:up q:quit ?:help")
	return left + right
}
