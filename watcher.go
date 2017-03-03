package distonic

import (
	git "github.com/libgit2/git2go"
)

type Watcher struct {
	dir  string
	url  string
	repo git.Repository
}

func NewWatcher(name, url string) *Watcher {
	return &Watcher{}
}
