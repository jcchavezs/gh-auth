// Package cli wires up the gh-auth cobra command tree.
package cli

import (
	"github.com/spf13/cobra"
)

// version is overridden at build time via -ldflags.
var version = "dev"

// ANSI escape codes used to match the original script's output styling.
const (
	colorGreen = "\033[0;32m"
	colorReset = "\033[0m"
	bold       = "\033[1m"
	normal     = "\033[0m"
)

// NewRootCmd builds the root gh-auth command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "gh-auth",
		Short:   "Authenticate gh and git with GitHub and configure git accordingly.",
		Long:    "Authenticate gh and git with GitHub and configure git accordingly.",
		Version: version,
		Example: `  gh-auth switch
  gh-auth sign-status`,
		SilenceUsage: true,
	}

	root.AddCommand(
		newSwitchCmd(),
		newSignStatusCmd(),
		newStatusCmd(),
	)

	return root
}
