HN - A Hacker News Pager
========================

HN is a simplistic Hacker News pager. It was initially desiged to compliment the Newsboat RSS reader in both style and usability, though it can also be used standalone.

Build Instructions
------------------

Like most Go programs, HN can be built by simply running `go build`. Optionally, it can be built and installed to your `GOPATH` by running `go install`.

Usage with Newsboat
-------------------

Assuming you are using the feed `https://hnrss.org/frontpage`, the following macro will allow you to open the comments to the currently selected feed when in article mode.

```
macro c pipe-to ~/.config/newsboat/scripts/hn-comments.sh -- "Open in HN viewer"
```

With the contents of `~/.config/newsboat/scripts/hn-comments.sh` as follows:

```bash
#!/bin/sh

input="$(cat - | grep -oP 'id=[0-9]+' | head -n1 | cut -d'=' -f2)"
exec hn "$input"
```
