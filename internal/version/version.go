package version

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"versioner/internal/config"

	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	Major = "major"
	Minor = "minor"
	Patch = "patch"
)

func Bump(repo *git.Repository) error {
	if err := config.Ensure(); err != nil {
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

	for project, bump := range bumps {
		_, _, ver, err := findLatestSemverTag(repo, project)
		if err != nil {
			return err
		}

		newVer := bumpVersion(ver, bump)

		fmt.Println(fmt.Sprintf("new version will be '%s'", newVer.String()))
	}

	return nil
}

func findLatestSemverTag(repo *git.Repository, tagToFind string) (string, plumbing.Hash, *semver.Version, error) {
	tagList := make(map[plumbing.Hash]string)

	tags, err := repo.Tags()
	if err != nil {
		return "", plumbing.ZeroHash, nil, err
	}

	for ref, err := tags.Next(); err == nil; ref, err = tags.Next() {
		tagName := ref.Name().Short()
		if !strings.Contains(tagName, tagToFind) {
			continue
		}

		obj, err := repo.TagObject(ref.Hash())
		if err == nil {
			tagList[obj.Target] = tagName
		} else {
			tagList[ref.Hash()] = tagName
		}
	}
	tags.Close()

	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return "", plumbing.ZeroHash, nil, err
	}
	defer iter.Close()

	for ref, err := iter.Next(); err == nil; ref, err = iter.Next() {
		tag, found := tagList[ref.Hash]
		if found {
			version, err := semver.NewVersion(tag)
			if err == nil {
				return tag, ref.Hash, version, nil
			}
		}
	}

	version, err := semver.NewVersion("0.0.0")
	if err != nil {
		return "", plumbing.ZeroHash, nil, err
	}

	return "", plumbing.ZeroHash, version, nil
}

func getHighestBumps(changesets []string) map[string]string {
	bumps := map[string]string{}

	for _, changeset := range changesets {

		if len(changeset) < 4 {
			continue
		}

		projectsStr := changeset[4:]
		projectsStr = projectsStr[:strings.Index(projectsStr, "---")]

		projects := removeEmptyStrings(strings.Split(projectsStr, "\n"))

		for _, project := range projects {
			nameAndBump := strings.Split(project, ": ")
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
	if base == "major" {
		return false
	}

	if comp == "major" {
		return true
	}

	if base == "minor" && comp == "patch" {
		return false
	}

	return true
}

func bumpVersion(ver *semver.Version, bump string) semver.Version {
	var newVer semver.Version

	switch bump {
	case "major":
		newVer = ver.IncMajor()
	case "minor":
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
