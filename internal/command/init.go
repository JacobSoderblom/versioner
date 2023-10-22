package command

import (
	"versioner/internal/config"
	"versioner/internal/context"
)

type Init struct{}

func (i Init) Run(ctx *context.Context) error {
	return config.Create(ctx)
}
