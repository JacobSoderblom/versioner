package command

import (
	"versioner/internal/add"
	"versioner/internal/context"
)

type Add struct{}

func (a Add) Run(ctx *context.Context) error {
	return add.Add(ctx)
}
