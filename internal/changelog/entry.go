package changelog

import (
	"fmt"
	"strings"
	"versioner/internal/changeset"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

type Section struct {
	Title   string
	Content []string
}

type Entry struct {
	Version  string
	Sections []Section
}

func NewEntry(curr semver.Version, cc changeset.Changesets) (Entry, error) {
	next, err := cc.HighestLevel()
	if err != nil {
		return Entry{}, errors.Wrap(err, "could not create a new entry")
	}

	newVer := bumpVersion(curr, next)

	ss := []Section{}

	for _, c := range cc {
		ct, err := c.ConventionalType()
		if err != nil {
			return Entry{}, errors.Wrap(err, "could not create a new entry")
		}

		title := ct.Title

		if c.Breaking {
			title = "Breaking changes"
		}

		s := getOrCreateSection(title, ss)
		s.Content = append(s.Content, c.Summary)

		ss = append(ss, s)
	}

	e := Entry{
		Version:  newVer.String(),
		Sections: ss,
	}

	return e, nil
}

func (e Entry) Markdown() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s\n", e.Version))

	for _, s := range e.Sections {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("### %s\n", s.Title))

		for _, c := range s.Content {
			sb.WriteString("\n")
			sb.WriteString(c)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func parseEntry(str string) (Entry, error) {
	e := Entry{}

	sCount := 0
	sections := []Section{}

	ss := removeEmptyStrings(strings.Split(str, "\n"))
	for _, s := range ss {
		if strings.Contains(s, "## ") {
			ver, err := semver.NewVersion(strings.Replace(s, "## ", "", 1))
			if err != nil {
				return Entry{}, err
			}

			e.Version = ver.String()
		}

		if strings.Contains(s, "### ") {
			if sCount > 0 {
				sCount++
			}

			title := strings.Replace(s, "### ", "", 1)
			sections = append(sections, Section{Title: title})
		}

		if len(sections) > 0 {
			sections[sCount].Content = append(sections[sCount].Content, s)
		}

	}

	e.Sections = sections

	return e, nil
}

func bumpVersion(ver semver.Version, bump string) semver.Version {
	var newVer semver.Version

	switch bump {
	case changeset.Major:
		newVer = ver.IncMajor()
	case changeset.Minor:
		newVer = ver.IncMinor()
	default:
		newVer = ver.IncPatch()

	}
	return newVer
}

func getOrCreateSection(title string, ss []Section) Section {
	for _, s := range ss {
		if s.Title == title {
			return s
		}
	}

	return Section{Title: title}
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
