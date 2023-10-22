package main

import (
	"fmt"
	"os"
	"versioner/internal/command"
	"versioner/internal/context"

	"github.com/alecthomas/kong"
	"github.com/go-git/go-git/v5"
)

var cmd struct {
	Init    command.Init    `cmd:"" help:"Initialize setup of project."`
	Add     command.Add     `cmd:"" help:"Add changelog to your project"`
	Version command.Version `cmd:"" help:"Creates a new version based on existing changesets"`
	Tag     command.Tag     `cmd:"" help:"Creates a new tag of the current version"`
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
