package changelog

import (
	"fmt"
	"os"
	"path"
	"strings"
	"versioner/internal/detect"

	"github.com/pkg/errors"
)

const commentStr = "[//]: # entry"

type Changelog struct {
	Title   string
	Entries []Entry
	Path    string
}

func Parse(wd string) (Changelog, error) {
	p, err := detect.Run(wd)
	if err != nil {
		return Changelog{}, errors.Wrap(err, "could not get changelog")
	}

	filePath := path.Join(wd, "CHANGELOG.md")
	title := fmt.Sprintf("# %s\n", p.Name)

	_, err = os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		if err = os.WriteFile(filePath, []byte(title), os.ModePerm); err != nil {
			return Changelog{}, errors.Wrap(err, "could not get changelog")
		}
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		return Changelog{}, errors.Wrap(err, "could not get changelog")
	}

	c := Changelog{
		Title: p.Name,
		Path:  filePath,
	}

	content := string(b)
	content = strings.ReplaceAll(content, title, "")
	entryStrs := removeEmptyStrings(strings.Split(content, commentStr))

	for _, eStr := range entryStrs {
		e, err := parseEntry(eStr)
		if err != nil {
			return Changelog{}, err
		}

		c.Entries = append(c.Entries, e)
	}

	return c, nil
}

func (c Changelog) Markdown() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n", c.Title))

	for _, e := range c.Entries {
		sb.WriteString("\n")
		sb.WriteString(e.Markdown())
	}

	return sb.String()
}

func (c Changelog) Save() error {
	content := c.Markdown()

	return errors.Wrap(os.WriteFile(c.Path, []byte(content), os.ModePerm), "could not save Changelog.md")
}

func (c *Changelog) Add(e Entry) {
	c.Entries = append([]Entry{e}, c.Entries...)
}
