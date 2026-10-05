package main

import (
	"os"

	"github.com/jcchavezs/gh-auth/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
