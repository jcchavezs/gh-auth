package cli

import (
	"github.com/jcchavezs/gh-auth/internal/gitauth"
	"github.com/spf13/cobra"
)

func newSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch",
		Short: "Switch GitHub account and configure Git",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := gitauth.SwitchAccount(cmd.Context()); err != nil {
				return err
			}
			return configureGitFromGitHub(cmd)
		},
	}
}
