package context

import "github.com/go-git/go-git/v5"

type Context struct {
	repo *git.Repository
	wd   string
}

func New(repo *git.Repository, wd string) Context {
	return Context{
		wd:   wd,
		repo: repo,
	}
}

func (c Context) Wd() string {
	return c.wd
}

func (c Context) Repo() *git.Repository {
	return c.repo
}
