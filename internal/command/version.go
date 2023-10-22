package command

import (
	"versioner/internal/context"
	"versioner/internal/version"
)

type Version struct{}

func (v Version) Run(ctx *context.Context) error {
	return version.Bump(ctx)
}
