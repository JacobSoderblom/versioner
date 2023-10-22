package add

import (
	"fmt"
	"os"
	"path"
	"strings"
	"versioner/internal/config"
	"versioner/internal/context"
	"versioner/internal/detect"
	"versioner/internal/tag"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

func Add(ctx *context.Context) error {
	if err := config.Ensure(ctx.Wd()); err != nil {
		return err
	}

	project, err := detect.Run(ctx.Wd())
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

	return errors.Wrap(change.Save(ctx.Wd()), "could not save new changeset")
}

func GetLatestChangeSet(ctx *context.Context) (string, error) {
	p := path.Join(ctx.Wd(), "CHANGELOG.md")
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}

	_, _, ver, err := tag.FindLatest(ctx.Repo(), "")
	if err != nil {
		return "", err
	}

	tag := fmt.Sprintf("v%s", ver.String())

	tags, err := ctx.Repo().TagObjects()
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
