package helpers

import (
	"errors"
	"log"
	"os/exec"
)

func IsInsideGitRepo() bool {
	_, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	var exErr *exec.ExitError
	if !errors.As(err, &exErr) {
		log.Fatalf("%p. Exiting...", err)
	}

	return err == nil
}
