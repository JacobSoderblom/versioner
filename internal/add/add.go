package add

import (
	"fmt"
	"os"
	"path"
	"strings"
	"versioner/internal/config"
	"versioner/internal/detect"
	"versioner/internal/tag"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

func Add() error {
	if err := config.Ensure(); err != nil {
		return err
	}

	project, err := detect.Run()
	if err != nil {
		return err
	}

	change, abort, err := startTea(project)
	if err != nil {
		return err
	}

	if abort {
		return nil
	}

	return errors.Wrap(change.Save(), "could not save new changeset")
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
