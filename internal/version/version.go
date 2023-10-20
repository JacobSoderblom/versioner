package version

import (
	"strings"
	"versioner/internal/changelog"
	"versioner/internal/changeset"
	"versioner/internal/config"
	"versioner/internal/tag"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

func Bump(repo *git.Repository) error {
	conf, err := config.Read()
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	cc, err := changeset.ParseChangesets()
	if err != nil {
		return errors.Wrap(err, "could not read changesets")
	}

	_, _, curr, err := tag.FindLatest(repo, "")
	if err != nil {
		return err
	}

	entry, err := changelog.NewEntry(*curr, cc)
	if err != nil {
		return err
	}

	c, err := changelog.Parse()
	if err != nil {
		return err
	}

	c.Add(entry)

	if err = c.Save(); err != nil {
		return err
	}

	conf.NextVersion = strings.Replace(entry.Version, "v", "", 1)
	if err = config.Set(conf); err != nil {
		return err
	}

	if err = cc.Remove(); err != nil {
		return err
	}

	if conf.Commit {
		if err = w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
			return err
		}

		msg := "New version"
		if len(conf.CommitMsg) > 0 {
			msg = conf.CommitMsg
		}

		if _, err = w.Commit(msg, &git.CommitOptions{
			All:   true,
			Amend: conf.AmendCommit,
		}); err != nil {
			return err
		}
	}

	return nil
}
