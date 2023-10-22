package version

import (
	"strings"
	"versioner/internal/changelog"
	"versioner/internal/changeset"
	"versioner/internal/config"
	"versioner/internal/context"
	"versioner/internal/tag"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

func Bump(ctx *context.Context) error {
	conf, err := config.Read(ctx.Wd())
	if err != nil {
		return err
	}

	w, err := ctx.Repo().Worktree()
	if err != nil {
		return err
	}

	cc, err := changeset.ParseChangesets(ctx.Wd())
	if err != nil {
		return errors.Wrap(err, "could not read changesets")
	}

	_, _, curr, err := tag.FindLatest(ctx.Repo(), "")
	if err != nil {
		return err
	}

	entry, err := changelog.NewEntry(*curr, cc)
	if err != nil {
		return err
	}

	c, err := changelog.Parse(ctx.Wd())
	if err != nil {
		return err
	}

	c.Add(entry)

	if err = c.Save(); err != nil {
		return err
	}

	conf.NextVersion = strings.Replace(entry.Version, "v", "", 1)
	if err = config.Set(ctx.Wd(), conf); err != nil {
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
