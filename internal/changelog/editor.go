package changelog

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/tcnksm/go-gitconfig"
)

func openEditor() (string, error) {
	editor, err := gitconfig.Global("core.editor")
	if err != nil {
		return "", err
	}

	f, err := os.CreateTemp("", "changelog-editor")
	if err != nil {
		return "", errors.Wrap(err, "could not create temp file")
	}

	defer os.Remove(f.Name())

	cmd := exec.Command("sh", "-c", editor+" "+f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("could not start %s", editor))
	}

	// read tmpfile
	b, err := os.ReadFile(f.Name())
	if err != nil {
		return "", errors.Wrap(err, "could not read temp file")
	}

	return string(b), nil
}
