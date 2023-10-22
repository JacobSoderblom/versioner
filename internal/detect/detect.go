package detect

import (
	"errors"
	"os"
	"strings"
)

type Project struct {
	Name string
	Path string
}

func Run(wd string) (Project, error) {
	project, err := Golang(wd)
	if !errors.Is(err, os.ErrNotExist) || err == nil {
		return project, err
	}

	strs := strings.Split(wd, "/")

	return Project{
		Name: strs[len(strs)-1],
		Path: "/",
	}, nil
}
