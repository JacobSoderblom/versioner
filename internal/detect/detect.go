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

func Run() (Project, error) {
	project, err := Golang()
	if !errors.Is(err, os.ErrNotExist) || err == nil {
		return project, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return project, nil
	}

	strs := strings.Split(wd, "/")

	return Project{
		Name: strs[len(strs)-1],
		Path: "/",
	}, nil
}
