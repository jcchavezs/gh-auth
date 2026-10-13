package cmd

import (
	"fmt"

	"github.com/fatih/color"
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

// configureGitFromGitHub resolves the current GitHub account and matching GPG
// key (prompting the user when information is missing) and writes the global
// git configuration.
func configureGitFromGitHub(cmd *cobra.Command) error {
	account, err := gitauth.GetAccountFromGithub(cmd.Context(), cmd.OutOrStdout(), cmd.InOrStdin())
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	result, err := gitauth.ConfigureGit(account.Name, account.Email, account.SigningKey)
	if err != nil {
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	if result.GPGKey == "" {
		_, _ = fmt.Fprintf(out, "%s Git identity updated: %s (no signing key)\n",
			green("✓"), green(fmt.Sprintf("%s <%s>", result.Username, result.Email)))
	} else {
		_, _ = fmt.Fprintf(out, "%s Git identity updated: %s, GPG key %s\n",
			green("✓"), green(fmt.Sprintf("%s <%s>", result.Username, result.Email)), result.GPGKey)
	}

	return nil
}
