package command

import (
	"versioner/internal/config"
	"versioner/internal/context"
	"versioner/internal/detect"
	"versioner/internal/tui"

	"github.com/pkg/errors"
)

type Add struct{}

func (a Add) Run(ctx *context.Context) error {
	if err := config.Ensure(ctx.Wd()); err != nil {
		return err
	}

	project, err := detect.Run(ctx.Wd())
	if err != nil {
		return err
	}

	change, abort, err := tui.NewAddProgram(project)
	if err != nil {
		return err
	}

	if abort {
		return nil
	}

	return errors.Wrap(change.Save(ctx.Wd()), "could not save new changeset")
}
