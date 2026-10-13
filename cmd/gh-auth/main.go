package main

import (
	"os"

	"github.com/jcchavezs/gh-auth/internal/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
