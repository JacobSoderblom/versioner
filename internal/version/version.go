package version

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"versioner/internal/config"
	"versioner/internal/tag"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	Major = "major"
	Minor = "minor"
	Patch = "patch"
	caser = cases.Title(language.English)
)

func Bump(repo *git.Repository) error {
	conf, err := config.Read()
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
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

	var changesets []string
	for _, p := range changesetPaths {
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}

		changesets = append(changesets, string(b))
	}

	bumps := getHighestBumps(changesets)

	_, _, ver, err := tag.FindLatest(repo, "")
	if err != nil {
		return err
	}

	var newVer semver.Version
	var project string

	for p, bump := range bumps {
		project = p
		newVer = bumpVersion(ver, bump)

		break
	}

	tag := fmt.Sprintf("v%s", newVer.String())

	p := path.Join(wd, "CHANGELOG.md")

	b, err := readChangelog(p, project)
	if err != nil {
		return err
	}

	entry := createEntry(tag, changesets)

	content := strings.Split(string(b), "\n")
	content = insert(content, 1, entry)
	text := strings.Join(content, "\n")

	if err = writeChangelog(p, []byte(text)); err != nil {
		return err
	}

	conf.NextVersion = tag

	if err = config.Set(conf); err != nil {
		return err
	}

	for _, cp := range changesetPaths {
		if err = os.Remove(cp); err != nil {
			return err
		}
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
			Amend: conf.AmendCommit,
		}); err != nil {
			return err
		}
	}

	return nil
}

func createEntry(version string, changesets []string) string {
	groupedChangelog := groupChangesets(changesets)

	entry := fmt.Sprintf("\n## %s\n\n", version)

	for bump, changes := range groupedChangelog {
		entry += fmt.Sprintf("### %s changes\n\n", caser.String(strings.Split(bump, ": ")[1]))
		for i, change := range changes {
			if i == len(changes)-1 {
				entry += fmt.Sprintf("- %s", change)
				continue
			}

			entry += fmt.Sprintf("- %s\n", change)
		}
	}

	return entry
}

func getHighestBumps(changesets []string) map[string]string {
	bumps := map[string]string{}

	for _, changeset := range changesets {

		if len(changeset) < 4 {
			continue
		}

		bb := getBumps(changeset)

		for _, bump := range bb {
			nameAndBump := strings.Split(bump, ": ")
			if len(nameAndBump) < 2 {
				continue
			}

			if prev, ok := bumps[nameAndBump[0]]; ok && !isBumpHigher(prev, nameAndBump[1]) {
				continue
			}

			bumps[nameAndBump[0]] = nameAndBump[1]
		}
	}

	return bumps
}

func isBumpHigher(base, comp string) bool {
	if base == Major {
		return false
	}

	if comp == Major {
		return true
	}

	if base == Minor && comp == Patch {
		return false
	}

	return true
}

func bumpVersion(ver *semver.Version, bump string) semver.Version {
	var newVer semver.Version

	switch bump {
	case Major:
		newVer = ver.IncMajor()
	case Minor:
		newVer = ver.IncMinor()
	default:
		newVer = ver.IncPatch()

	}
	return newVer
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

func groupChangesets(changesets []string) map[string][]string {
	g := map[string][]string{}

	for _, c := range changesets {
		text := getChangelogText(c)

		for _, bump := range getBumps(c) {
			g[bump] = append(g[bump], text)
		}
	}

	return g
}

func getBumps(changeset string) []string {
	bumpStr := changeset[4:]
	bumpStr = bumpStr[:strings.Index(bumpStr, "---")]

	bumps := removeEmptyStrings(strings.Split(bumpStr, "\n"))

	return bumps
}

func getChangelogText(changeset string) string {
	changeset = changeset[4:]
	changeset = changeset[strings.Index(changeset, "---")+5 : len(changeset)-2]

	return changeset
}

func readChangelog(filePath, project string) ([]byte, error) {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		if err = os.WriteFile(filePath, []byte(fmt.Sprintf("# %s\n", project)), os.ModePerm); err != nil {
			return []byte{}, err
		}
	}

	return os.ReadFile(filePath)
}

func writeChangelog(filePath string, content []byte) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return os.WriteFile(filePath, content, os.ModePerm)
}

func insert(a []string, index int, value string) []string {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}
