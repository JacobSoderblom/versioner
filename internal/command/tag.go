package command

import (
	"versioner/internal/context"
	"versioner/internal/tag"
)

type Tag struct{}

func (t Tag) Run(ctx *context.Context) error {
	return tag.Tag(ctx)
}
