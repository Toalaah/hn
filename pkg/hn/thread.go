package hn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/toalaah/hn/pkg/threadview"
)

func (t *Story) ID() int                           { return t.Id }
func (t *Story) Parent() (threadview.Thread, bool) { return t.parent, t.parent != nil }
func (t *Story) Children() []threadview.Thread {
	res := make([]threadview.Thread, len(t.Children_))
	for i := range t.Children_ {
		res[i] = &t.Children_[i]
	}
	return res
}

func (t *Story) Text() string {
	var b strings.Builder
	for _, p := range t.textParts {
		b.WriteString(p.Text)
	}
	return b.String()
}

func dfs(root, cur *Story, f func(root, cur *Story)) {
	f(root, cur)
	for i := range cur.Children_ {
		dfs(cur, &cur.Children_[i], f)
	}
}

func (t *Story) CommentIndex(target threadview.Thread) (int, bool) {
	return commentIndex(t, target.(*Story), 0)
}

func commentIndex(cur, target *Story, n int) (int, bool) {
	if cur == target {
		return n, true
	}
	found := false
	for i := range cur.Children_ {
		if n, found = commentIndex(&cur.Children_[i], target, n+1); found {
			return n, found
		}
	}
	return n, false
}

func initNodes(t *Story) {
	dfs(nil, t, func(root, cur *Story) {
		cur.parent = root
		cur.textParts = parseMarkupToBlocks(cur.TextRaw)
	})
}

func NewThreadFromData(body []byte) (*Story, error) {
	t := &Story{}
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, err
	}
	initNodes(t)
	return t, nil
}

func NewThread(id int) (*Story, error) {
	resp, err := http.Get(fmt.Sprintf("https://hn.algolia.com/api/v1/items/%d", id))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(resp.Status)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return NewThreadFromData(body)
}

func (t *Story) UnmarshalJSON(data []byte) error {
	type Dummy Story

	tmp := struct {
		URL string `json:"url"`
		*Dummy
	}{Dummy: (*Dummy)(t)}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if url, err := url.Parse(tmp.URL); err != nil {
		return err
	} else if url.Scheme != "" {
		t.URL = url
	}

	return nil
}
