package hn

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Represents a HN thread.
type Thread struct {
	ID       int       `json:"id"`
	Author   string    `json:"author"`
	Date     time.Time `json:"created_at"`
	Text     string    `json:"text"`
	Parent   *Thread   `json:"-"`
	Title    string    `json:"title,omitempty"`
	URL      *url.URL  `json:"url"`
	Points   int       `json:"points"`
	Children []Thread  `json:"children"`
}

func (t *Thread) NumComments() int {
	var numCommentsFromNode func(tt *Thread) int
	numCommentsFromNode = func(tt *Thread) int {
		i := 1
		for _, cc := range tt.Children {
			i += numCommentsFromNode(&cc)
		}
		return i
	}
	i := 0
	for _, c := range t.Children {
		i += numCommentsFromNode(&c)
	}
	return i
}

func dfs(root, cur *Thread, f func(root, cur *Thread)) {
	f(root, cur)
	for i := range cur.Children {
		dfs(cur, &cur.Children[i], f)
	}
}

func (t *Thread) CommentIndex(target *Thread) (int, bool) {
	return commentIndex(t, target, 0)
}

func commentIndex(cur, target *Thread, n int) (int, bool) {
	if cur == target {
		return n, true
	}
	found := false
	acc := 0
	for i := range cur.Children {
		if n, found = commentIndex(&cur.Children[i], target, n+1); found {
			return n, found
		}
		acc += n
	}
	return n, false
}

func initNodes(t *Thread) {
	dfs(nil, t, func(root, cur *Thread) {
		cur.Parent = root
		text, err := strconv.Unquote("`" + cur.Text + "`")
		if err != nil {
			text = cur.Text
		}
		text = strings.ReplaceAll(text, "\\n", "\n")
		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ReplaceAll(text, "<p>", "\n")
		text = strings.ReplaceAll(text, "  ", " ")
		text = strings.TrimSpace(text)
		text = html.UnescapeString(text)
		cur.Text = text
	})
}

func NewThreadFromData(body []byte) (Thread, error) {
	t := Thread{}
	if err := json.Unmarshal(body, &t); err != nil {
		return Thread{}, err
	}
	initNodes(&t)
	return t, nil
}

func NewThread(id int) (Thread, error) {
	resp, err := http.Get(fmt.Sprintf("https://hn.algolia.com/api/v1/items/%d", id))
	if err != nil {
		return Thread{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Thread{}, err
	}
	return NewThreadFromData(body)
}

func (t *Thread) UnmarshalJSON(data []byte) error {
	type Dummy Thread

	tmp := struct {
		URL string `json:"url"`
		*Dummy
	}{Dummy: (*Dummy)(t)}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if url, err := url.Parse(tmp.URL); err != nil {
		return err
	} else {
		t.URL = url
	}

	return nil
}
