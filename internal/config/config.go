package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

const (
	Dir      = ".versioner"
	fileName = "config.json"
)

var (
	NotInitialized     = errors.New("versioner is not initialized")
	AlreadyInitialized = errors.New("versioner is already initialized")
)

type Configuration struct {
	BaseBranch  string   `json:"baseBranch,omitempty"`
	Ignore      []string `json:"ignore,omitempty"`
	Commit      bool     `json:"commit,omitempty"`
	CommitMsg   string   `json:"commitMsg,omitempty"`
	AmendCommit bool     `json:"amendCommit,omitempty"`
	NextVersion string   `json:"nextVersion,omitempty"`
}

func Create(repo *git.Repository) error {
	if err := Ensure(); err == nil {
		return AlreadyInitialized
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := path.Join(wd, Dir)
	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return err
	}

	configFilePath := path.Join(configPath, fileName)

	branch, err := getBaseBranch(repo)
	if err != nil {
		return err
	}

	config := Configuration{
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

func Ensure() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err = os.Stat(path.Join(wd, Dir, fileName)); os.IsNotExist(err) {
		return NotInitialized
	}

	return nil
}

func Read() (Configuration, error) {
	var config Configuration

	if err := Ensure(); err != nil {
		return config, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return config, err
	}

	b, err := os.ReadFile(path.Join(wd, Dir, fileName))
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(b, &config); err != nil {
		return config, err
	}

	return config, nil
}

func Set(conf Configuration) error {
	if err := Ensure(); err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := path.Join(wd, Dir, fileName)

	if err = os.Remove(configPath); err != nil {
		return err
	}

	var b []byte
	if b, err = json.MarshalIndent(&conf, "", "  "); err != nil {
		return err
	}

	if err = os.WriteFile(configPath, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func getBaseBranch(repo *git.Repository) (string, error) {
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
