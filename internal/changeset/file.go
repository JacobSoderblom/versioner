package changeset

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"versioner/internal/config"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/pkg/errors"
)

const mdTemplate = `---
%s
---

%s

`

func ParseChangesets() (Changesets, error) {
	changesets := []Changeset{}

	wd, err := os.Getwd()
	if err != nil {
		return changesets, err
	}

	configPath := path.Join(wd, config.Dir)

	var changesetPaths []string
	filepath.WalkDir(configPath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if filepath.Ext(d.Name()) == ".md" {
			changesetPaths = append(changesetPaths, s)
		}

		return nil
	})

	changesets, err = parseChangesets(changesetPaths)
	if err != nil {
		return changesets, err
	}

	return changesets, nil
}

func writeChangesetFile(content string) error {
	p, err := generateUniquePath()
	if err != nil {
		return err
	}
	return os.WriteFile(p, []byte(content), os.ModePerm)
}

func parseChangesets(paths []string) ([]Changeset, error) {
	var changesets []Changeset
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return changesets, err
		}

		c, err := parseChangeset(string(b), p)
		if err != nil {
			return changesets, err
		}

		c.path = p

		changesets = append(changesets, c)
	}

	return changesets, nil
}

func parseChangeset(str, file string) (Changeset, error) {
	if len(str) < 4 {
		return Changeset{}, errors.Wrap(ErrChangesetMalformated, fmt.Sprintf("could not parse changeset '%s'", file))
	}

	str = str[4:]
	endSectionIndex := strings.Index(str, "---") - 1

	if endSectionIndex == -1 || endSectionIndex > len(str)-1 {
		return Changeset{}, errors.Wrap(ErrChangesetMalformated, fmt.Sprintf("could not parse changeset '%s'", file))
	}

	conType := str[:endSectionIndex]

	str = str[endSectionIndex+4:]

	rest := removeEmptyStrings(strings.Split(str, "\n"))

	if len(rest) == 0 {
		return Changeset{}, errors.Wrap(ErrChangesetMalformated, fmt.Sprintf("could not parse changeset '%s'", file))
	}

	summary := rest[0]

	return Changeset{
		Type:     strings.ReplaceAll(conType, "!", ""),
		Summary:  summary,
		Breaking: strings.Contains(conType, "!"),
	}, nil
}

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}

	return r
}

func generateUniquePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	name := namesgenerator.GetRandomName(0)

	p := path.Join(wd, config.Dir, fmt.Sprintf("%s.md", name))

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return p, nil
	}

	return generateUniquePath()
}
