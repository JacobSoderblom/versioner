package detect

import (
	"os"
	"path"
	"strings"
)

func Golang(wd string) (Project, error) {
	p := Project{}

	goModPath := path.Join(wd, "go.mod")

	_, err := os.Stat(goModPath)
	if err != nil {
		return p, err
	}

	b, err := os.ReadFile(goModPath)
	if err != nil {
		return p, err
	}

	goMod := string(b)
	moduleStr := "module "

	module := goMod[strings.Index(goMod, moduleStr)+len(moduleStr):]
	module = module[:strings.Index(module, "\n")]

	p.Name = module

	return p, nil
}
