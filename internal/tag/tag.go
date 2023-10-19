package tag

import (
	"versioner/internal/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

var (
	NoVersionAvailable = errors.New("there is no version available for tagging")
	TagAlreadyExist    = errors.New("tag already exist")
)

func Tag(repo *git.Repository) error {
	conf, err := config.Read()
	if err != nil {
		return err
	}

	if len(conf.NextVersion) == 0 {
		return NoVersionAvailable
	}

	if err = tagExists(conf.NextVersion, repo); err != nil {
		return err
	}

	h, err := repo.Head()
	if err != nil {
		return err
	}

	if _, err = repo.CreateTag(conf.NextVersion, h.Hash(), &git.CreateTagOptions{Message: conf.NextVersion}); err != nil {
		return err
	}

	return nil
}

func tagExists(tag string, r *git.Repository) error {
	tags, err := r.TagObjects()
	if err != nil {
		return err
	}

	res := false
	tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return nil
		}
		return nil
	})

	if res {
		return TagAlreadyExist
	}

	return nil
}
