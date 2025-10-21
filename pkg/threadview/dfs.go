package threadview

func dfs(root, cur Thread, f func(root, cur Thread)) {
	f(root, cur)
	for _, c := range cur.Children() {
		dfs(cur, c, f)
	}
}

// NumNodes returns the number of children from the root thread t.
func NumNodes(t Thread) int {
	n := 0
	dfs(nil, t, func(root, cur Thread) { n++ })
	return n
}
