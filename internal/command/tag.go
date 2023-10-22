package command

import (
	"errors"
	"strings"
	"versioner/internal/config"
	"versioner/internal/context"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	NoVersionAvailable = errors.New("there is no version available for tagging")
	TagAlreadyExist    = errors.New("tag already exist")
)

type Tag struct{}

func (t Tag) Run(ctx *context.Context) error {
	conf, err := config.Read(ctx.Wd())
	if err != nil {
		return err
	}

	if len(conf.NextVersion) == 0 {
		return nil
	}

	if err = t.tagExists(conf.NextVersion, ctx.Repo()); err != nil {
		return err
	}

	h, err := ctx.Repo().Head()
	if err != nil {
		return err
	}

	if _, err = ctx.Repo().CreateTag(conf.NextVersion, h.Hash(), &git.CreateTagOptions{Message: conf.NextVersion}); err != nil {
		return err
	}

	return nil
}

func (t Tag) tagExists(tag string, r *git.Repository) error {
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

func findLatestTag(repo *git.Repository, tagToFind string) (string, plumbing.Hash, *semver.Version, error) {
	tagList := make(map[plumbing.Hash]string)

	tags, err := repo.Tags()
	if err != nil {
		return "", plumbing.ZeroHash, nil, err
	}

	for ref, err := tags.Next(); err == nil; ref, err = tags.Next() {
		tagName := ref.Name().Short()
		if !strings.Contains(tagName, tagToFind) {
			continue
		}

		obj, err := repo.TagObject(ref.Hash())
		if err == nil {
			tagList[obj.Target] = tagName
		} else {
			tagList[ref.Hash()] = tagName
		}
	}
	tags.Close()

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return "", plumbing.ZeroHash, nil, err
	}
	defer iter.Close()

	for ref, err := iter.Next(); err == nil; ref, err = iter.Next() {
		tag, found := tagList[ref.Hash]
		if found {
			version, err := semver.NewVersion(tag)
			if err == nil {
				return tag, ref.Hash, version, nil
			}
		}
	}

	version, err := semver.NewVersion("0.0.0")
	if err != nil {
		return "", plumbing.ZeroHash, nil, err
	}

	return "", plumbing.ZeroHash, version, nil
}
