package changeset

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

var (
	Major = "major"
	Minor = "minor"
	Patch = "patch"
	None  = ""
	Types = ConventionalTypes{
		{
			Title:         "Bug fixes",
			Type:          "fix",
			CanBeBreaking: true,
			Level:         Patch,
		}, {
			Title:         "New features",
			Type:          "feat",
			CanBeBreaking: true,
			Level:         Minor,
		}, {
			Title:         "Documentation",
			Type:          "docs",
			CanBeBreaking: false,
			Level:         None,
		}, {
			Title:         "Refactoring",
			Type:          "refactor",
			CanBeBreaking: false,
			Level:         None,
		}, {
			Title:         "Automations",
			Type:          "ci",
			CanBeBreaking: false,
			Level:         None,
		}, {
			Title:         "Miscellaneous",
			Type:          "chore",
			CanBeBreaking: false,
			Level:         None,
		}, {
			Title:         "Revert",
			Type:          "revert",
			CanBeBreaking: true,
			Level:         Patch,
		},
	}
)

var (
	ErrConventionalTypeNotFound = errors.New("convetional type not found")
	ErrChangesetMalformated     = errors.New("changeset malformated")
)

type Changeset struct {
	Breaking bool
	Type     string
	Summary  string
	path     string
}

func (c Changeset) ConventionalType() (ConventionalType, error) {
	for _, t := range Types {
		if t.Type == c.Type {
			return t, nil
		}
	}

	return ConventionalType{}, errors.Wrap(ErrConventionalTypeNotFound, fmt.Sprintf("%s:", c.Type))
}

func (c Changeset) Save() error {
	release := c.Type
	if c.Breaking {
		release += "!"
	}

	content := fmt.Sprintf(mdTemplate, release, c.Summary)

	return writeChangesetFile(content)
}

func (c Changeset) Remove() error {
	if len(c.path) == 0 {
		return nil
	}

	return os.Remove(c.path)
}

type Changesets []Changeset

func (cc Changesets) HighestLevel() (string, error) {
	level := Patch

	for _, c := range cc {
		ct, err := c.ConventionalType()
		if err != nil {
			return level, errors.Wrap(err, "could not calculate highest semver level")
		}
		if isLevelHigher(level, ct.Level) {
			level = ct.Level
		}
	}

	return level, nil
}

func (cc Changesets) Filter(t string) []Changeset {
	changesets := []Changeset{}

	for _, c := range cc {
		if c.Type == t {
			changesets = append(changesets, c)
		}
	}

	return changesets
}

func (cc Changesets) Remove() error {
	for _, c := range cc {
		if err := c.Remove(); err != nil {
			return err
		}
	}

	return nil
}

type ConventionalType struct {
	Title         string
	Type          string
	CanBeBreaking bool
	Level         string
}

type ConventionalTypes []ConventionalType

func (c ConventionalTypes) CanBeBreaking(t string) bool {
	for _, ct := range c {
		if ct.Type == t {
			return ct.CanBeBreaking
		}
	}

	return false
}

func isLevelHigher(base, comp string) bool {
	if base == Major {
		return false
	}

	if comp == Major {
		return true
	}

	if base == Minor && comp == Patch {
		return false
	}

	if comp == None {
		return false
	}

	return true
}
