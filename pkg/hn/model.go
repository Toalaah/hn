package hn

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/mergestat/timediff"
	"github.com/toalaah/hn/pkg/threadview"
)

type Story struct {
	Id        int       `json:"id"`
	Author    string    `json:"author"`
	Date      time.Time `json:"created_at"`
	TextRaw   string    `json:"text"`
	Title     string    `json:"title,omitempty"`
	Points    int       `json:"points"`
	URL       *url.URL  `json:"url"`
	Children_ []Story   `json:"children"`
	parent    *Story    `json:"-"`

	textParts []TextBlock
	state     State
}

type State struct {
	threadview.DisplayStateMsg
	TextWidth int
}

func (t *Story) Init() tea.Cmd { return nil }
func (t *Story) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case threadview.DisplayStateMsg:
		t.state = State{msg, 72}
		// Direct replies should not be indented.
		t.state.Depth = max(0, t.state.Depth-1)
	case threadview.CopyTextMsg:
		cmds = append(cmds, func() tea.Msg {
			err := clipboard.WriteAll(t.Text())
			return threadview.CopyTextResultMsg{Error: err}
		})
	}
	return t, tea.Batch(cmds...)
}

func (t *Story) CommentView() string {
	var b strings.Builder
	// Header
	{
		relDate := fmt.Sprintf("(%s)", timediff.TimeDiff(t.Date))
		numComments := "[-]"
		if t.state.Collapsed {
			numComments = fmt.Sprintf("[%d more]", threadview.NumNodes(t))
		}
		header := TextBlocks{
			{Type: BlockTypeAuthor, Text: t.Author},
			{Type: BlockTypeText, Text: " "},
			{Type: BlockTypeMetadata, Text: relDate},
			{Type: BlockTypeText, Text: " "},
			{Type: BlockTypeMetadata, Text: numComments},
			{Type: BlockTypeText, Text: "\n"},
		}
		// Trim newline if comment is collapsed
		if t.state.Collapsed {
			header = header[:len(header)-1]
		}
		b.WriteString(header.Render(t.state))
	}
	// Comment text
	if !t.state.Collapsed {
		b.WriteString(TextBlocks(t.textParts).Render(t.state))
	}
	return b.String()
}

func (t *Story) View() string {
	defer color.Unset()
	if t.parent == nil {
		return ""
	}
	return t.CommentView()
}
