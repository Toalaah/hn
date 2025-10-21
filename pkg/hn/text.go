package hn

import (
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/net/html"
)

type BlockType int

const (
	BlockTypeText BlockType = iota
	BlockTypeItalic
	BlockTypeQuote
	BlockTypeLink
	BlockTypeRaw
	// Define some non-comment-specific blocks so that we can re-use the render function below for the header.
	BlockTypeAuthor
	BlockTypeMetadata
)

type TextBlock struct {
	Type BlockType
	Text string
}

type TextBlocks []TextBlock

type TextStyle struct {
	Unselected *color.Color
	Selected   *color.Color
}

func (t TextStyle) Get(selected bool) *color.Color {
	if selected && t.Selected != nil {
		return t.Selected
	}
	return t.Unselected
}

var textStyles = struct {
	Text,
	Link,
	Quote,
	Author,
	Metadata TextStyle
	SelectedBg *color.Color
}{
	Text:       TextStyle{color.Set(color.FgWhite), nil},
	Link:       TextStyle{color.Set(color.FgRed), nil},
	Quote:      TextStyle{color.Set(color.Faint), nil},
	Author:     TextStyle{color.Set(color.FgHiYellow, color.Bold), nil},
	Metadata:   TextStyle{color.Set(color.Faint), color.Set(color.FgWhite)},
	SelectedBg: color.Set(color.BgBlue),
}

func (t TextBlocks) Render(state State) string {
	var (
		res string

		normal = textStyles.Text.Get(state.Selected)
		link   = textStyles.Link.Get(state.Selected)
		quote  = textStyles.Quote.Get(state.Selected)
		author = textStyles.Author.Get(state.Selected)
		meta   = textStyles.Metadata.Get(state.Selected)
	)

	write := func(w io.Writer, c *color.Color, text string) {
		normal.SetWriter(w)
		// c.SetWriter(w)
		c.Fprint(w, text)
		if state.Selected {
			textStyles.SelectedBg.SetWriter(w)
		}
		normal.SetWriter(w)
	}

	// Render comment w/ ansi markup
	{
		b := &strings.Builder{}
		for _, part := range t {
			switch part.Type {
			case BlockTypeText:
				fallthrough
			default:
				write(b, normal, part.Text)
			case BlockTypeLink:
				write(b, link, part.Text)
			case BlockTypeItalic:
				write(b, color.Set(color.Italic), part.Text)
			case BlockTypeQuote:
				write(b, quote, part.Text)
			case BlockTypeAuthor:
				write(b, author, part.Text)
			case BlockTypeMetadata:
				write(b, meta, part.Text)
			}
		}

		res = wordwrap.String(b.String(), state.TextWidth)
		i := indent.NewWriter(1, func(w io.Writer) {
			if state.Selected {
				textStyles.SelectedBg.SetWriter(w)
			}
			w.Write([]byte(strings.Repeat(" ", state.Depth*2)))
		})
		_, _ = i.Write([]byte(res))
		res = i.String()
		res = padding.String(res, uint(state.Width))
	}

	return res
}

func parseMarkupToBlocks(s string) []TextBlock {
	var (
		parts     = make([]TextBlock, 0)
		state     = BlockTypeText
		tokenizer = html.NewTokenizer(strings.NewReader(s))
	)
	for tokenType := tokenizer.Next(); tokenType != html.ErrorToken; tokenType = tokenizer.Next() {
		token := tokenizer.Token()
		switch tokenType {
		case html.StartTagToken:
			switch token.Data {
			case "p":
				state = BlockTypeText
				parts = append(parts, TextBlock{BlockTypeText, "\n"})
			case "i":
				state = BlockTypeItalic
			case "a":
				state = BlockTypeLink
			case "pre":
				state = BlockTypeRaw
				parts = append(parts, TextBlock{BlockTypeText, "\n"})
			case "code":
				break
			default:
				panic("unhandled start tag: " + token.Data)
			}
		case html.TextToken:
			if strings.HasPrefix(token.Data, ">") {
				state = BlockTypeQuote
			}
			parts = append(parts, TextBlock{state, token.Data})
		case html.EndTagToken:
			state = BlockTypeText
		}
	}
	return parts
}
