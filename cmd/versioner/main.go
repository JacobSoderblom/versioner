package main

import (
	"fmt"
	"os"
	"versioner/internal/changelog"
	"versioner/internal/config"
	"versioner/internal/version"

	"github.com/alecthomas/kong"
	"github.com/go-git/go-git/v5"
)

type Context struct {
	Debug bool
	Repo  *git.Repository
}

type InitCmd struct{}

func (i *InitCmd) Run(ctx *Context) error {
	if err := config.Ensure(); err == nil {
		return config.AlreadyInitialized
	}

	fmt.Println("Creating configuration file for your project!")

	return config.Create(ctx.Repo)
}

type ChangelogCmd struct{}

func (i ChangelogCmd) Run(ctx *Context) error {
	return changelog.Add()
}

type VersionCmd struct{}

func (v VersionCmd) Run(ctx *Context) error {
	return version.Bump(ctx.Repo)
}

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Init    InitCmd      `cmd:"" help:"Initialize setup of project."`
	Add     ChangelogCmd `cmd:"" help:"Add changelog to your project"`
	Version VersionCmd   `cmd:"" help:"Creates a new version based on existing changesets"`
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	repo, err := git.PlainOpen(wd)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err = ctx.Run(&Context{Debug: cli.Debug, Repo: repo})
	ctx.FatalIfErrorf(err)
}
