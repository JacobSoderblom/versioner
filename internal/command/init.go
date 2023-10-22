package command

import (
	"encoding/json"
	"os"
	"path"
	"versioner/internal/config"
	"versioner/internal/context"

	"github.com/go-git/go-git/v5"
)

type Init struct{}

func (i Init) Run(ctx *context.Context) error {
	if err := config.Ensure(ctx.Wd()); err == nil {
		return config.AlreadyInitialized
	}

	configPath := path.Join(ctx.Wd(), config.Dir)
	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return err
	}

	configFilePath := path.Join(configPath, config.FileName)

	branch, err := i.getBaseBranch(ctx.Repo())
	if err != nil {
		return err
	}

	config := config.Configuration{
		BaseBranch: branch,
		Ignore:     []string{},
	}

	var b []byte
	if b, err = json.MarshalIndent(&config, "", "  "); err != nil {
		return err
	}

	if err = os.WriteFile(configFilePath, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (i Init) getBaseBranch(repo *git.Repository) (string, error) {
	conf, err := repo.Config()
	if err != nil {
		return "", err
	}

	main := conf.Branches["main"]
	master := conf.Branches["master"]

	if main != nil {
		return "main", nil
	}

	if master != nil {
		return "master", nil
	}

	branches := []string{}

	for branch := range conf.Branches {
		branches = append(branches, branch)
	}

	if len(branches) > 0 {
		return branches[0], nil
	}

	return "", nil
}
