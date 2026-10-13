package cmd

import (
	"github.com/spf13/cobra"
)

// version is overridden at build time via -ldflags.
var version = "dev"

// NewRootCmd builds the root gh-auth command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "gh-auth",
		Short:   "Authenticate gh and git with GitHub and configure git accordingly.",
		Long:    "Authenticate gh and git with GitHub and configure git accordingly.",
		Version: version,
		Example: `  gh-auth switch
  gh-auth sign-status
  gh-auth shell work`,
		SilenceUsage: true,
	}

	root.AddCommand(
		newSwitchCmd(),
		newSignStatusCmd(),
		newStatusCmd(),
		newShellCmd(),
	)

	return root
}
