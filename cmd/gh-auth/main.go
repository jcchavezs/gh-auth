package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/jcchavezs/gh-auth/internal/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		// Propagate the exit status of a command run through `shell -- <command>`.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() > 0 {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
