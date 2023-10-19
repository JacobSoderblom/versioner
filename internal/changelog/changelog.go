package changelog

import (
	"fmt"
	"os"
	"path"
	"strings"
	"versioner/internal/config"
	"versioner/internal/detect"
	"versioner/internal/tag"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

	project, err := detect.Run()
	if err != nil {
		return err
	}

	changes, err := startTea(project)
	if err != nil {
		return err
	}

	if !changes.Confirmed {
		return nil
	}

	var releases []string
	for p, b := range changes.Releases {
		releases = append(releases, fmt.Sprintf("%s: %s", p, b))
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

func GetLatestChangeSet(repo *git.Repository) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	p := path.Join(wd, "CHANGELOG.md")
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}

	_, _, ver, err := tag.FindLatest(repo, "")
	if err != nil {
		return "", err
	}

	tag := fmt.Sprintf("v%s", ver.String())

	tags, err := repo.TagObjects()
	if err != nil {
		return "", err
	}

	foundCurrent := false
	prevTag := ""
	tags.ForEach(func(t *object.Tag) error {
		if foundCurrent {
			prevTag = t.Name
			// just to stop the loop
			return errors.New("found tag")
		}

		if t.Name == tag {
			foundCurrent = true
		}

		return nil
	})

	c := string(b)

	currTagTitle := fmt.Sprintf("## %s", tag)
	prevTagTitle := fmt.Sprintf("## %s", prevTag)

	currI := strings.Index(c, currTagTitle)
	nextI := strings.Index(c, prevTagTitle)

	if currI == -1 {
		return "", nil
	}

	if nextI > len(c)-1 {
		return "", nil
	}

	return c[currI:nextI], nil
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
