package changelog

import (
	"fmt"
	"os"
	"path"
	"strings"
	"versioner/internal/config"
	"versioner/internal/detect"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/pkg/errors"
)

var NoReleaseCandidates = errors.New("no release candidates")

const mdTemplate = `---
%s
---

%s

`

type Changelog struct {
	Confirmed bool
	Releases  map[string]string
	Summary   string
}

func Add() error {
	if err := config.Ensure(); err != nil {
		return err
	}

	project, err := detect.Golang()
	if err != nil {
		return err
	}

	changes, err := startTea(project)
	if err != nil {
		return err
	}

	var releases []string
	for k, v := range changes.Releases {
		releases = append(releases, fmt.Sprintf("%s: %s", k, v))
	}

	if len(releases) == 0 {
		return NoReleaseCandidates
	}

	content := fmt.Sprintf(mdTemplate, strings.Join(releases, "\n"), changes.Summary)

	filePath, err := generateUniquePath()
	if err != nil {
		return err
	}

	if err = os.WriteFile(filePath, []byte(content), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func generateUniquePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	name := namesgenerator.GetRandomName(0)

	p := path.Join(wd, config.Dir, fmt.Sprintf("%s.md", name))

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return p, nil
	}

	return generateUniquePath()
}
