package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

const (
	Dir      = ".versioner"
	FileName = "config.json"
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

func Ensure(wd string) error {
	if _, err := os.Stat(path.Join(wd, Dir, FileName)); os.IsNotExist(err) {
		return NotInitialized
	}

	return nil
}

func Read(wd string) (Configuration, error) {
	var config Configuration

	if err := Ensure(wd); err != nil {
		return config, err
	}

	b, err := os.ReadFile(path.Join(wd, Dir, FileName))
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(b, &config); err != nil {
		return config, err
	}

	return config, nil
}

func Set(wd string, conf Configuration) error {
	if err := Ensure(wd); err != nil {
		return err
	}

	configPath := path.Join(wd, Dir, FileName)

	if err := os.Remove(configPath); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&conf, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}
