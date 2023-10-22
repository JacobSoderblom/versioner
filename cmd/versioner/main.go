package main

import (
	"fmt"
	"os"
	"versioner/internal/add"
	"versioner/internal/config"
	"versioner/internal/context"
	"versioner/internal/tag"
	"versioner/internal/version"

	"github.com/alecthomas/kong"
	"github.com/go-git/go-git/v5"
)

type InitCmd struct{}

func (i *InitCmd) Run(ctx *context.Context) error {
	if err := config.Ensure(ctx.Wd()); err == nil {
		return config.AlreadyInitialized
	}

	fmt.Println("Creating configuration file for your project!")

	return config.Create(ctx)
}

type ChangelogCmd struct{}

func (i ChangelogCmd) Run(ctx *context.Context) error {
	return add.Add(ctx)
}

type VersionCmd struct{}

func (v VersionCmd) Run(ctx *context.Context) error {
	return version.Bump(ctx)
}

type TagCmd struct{}

func (t TagCmd) Run(ctx *context.Context) error {
	return tag.Tag(ctx)
}

type ChangesetCmd struct{}

func (c ChangesetCmd) Run(ctx *context.Context) error {
	changeset, err := add.GetLatestChangeSet(ctx)

	fmt.Print(changeset)

	return err
}

var cmd struct {
	Init      InitCmd      `cmd:"" help:"Initialize setup of project."`
	Add       ChangelogCmd `cmd:"" help:"Add changelog to your project"`
	Version   VersionCmd   `cmd:"" help:"Creates a new version based on existing changesets"`
	Tag       TagCmd       `cmd:"" help:"Creates a new tag of the current version"`
	Changeset ChangesetCmd `cmd:"" help:"Gets the changeset of current latest tagged version"`
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

	ctx := context.New(repo, wd)
	cli := kong.Parse(&cmd)
	// Call the Run() method of the selected parsed command.
	err = cli.Run(&ctx)
	cli.FatalIfErrorf(err)
}
