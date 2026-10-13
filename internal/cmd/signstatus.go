package cmd

import (
	"fmt"

	"github.com/jcchavezs/gh-auth/internal/cli"
	"github.com/jcchavezs/gh-auth/internal/gitauth"
	"github.com/spf13/cobra"
)

func newSignStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sign-status",
		Short: "Show current account and Git signing configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			showSignStatus(cmd)
			return nil
		},
	}
}

// valueOrNotSet returns v, or "Not set" when v is empty.
func valueOrNotSet(v string) string {
	if v == "" {
		return "Not set"
	}
	return v
}

func showSignStatus(cmd *cobra.Command) {
	out := cmd.OutOrStdout()
	status := gitauth.GetSignStatus(cmd.Context())

	switch {
	case !status.GHInstalled:
		_, _ = fmt.Fprintln(out, "GitHub CLI not installed")
	case status.Account != nil:
		_, _ = fmt.Fprintf(out, "%s✓%s Active GitHub account: %s%s <%s>%s\n",
			cli.ColorGreen, cli.ColorReset, cli.Bold, status.Account.Username, status.Account.Email, cli.ColorReset)
	default:
		_, _ = fmt.Fprintln(out, "No active GitHub account or GitHub CLI not authenticated")
	}

	_, _ = fmt.Fprintf(out, "- Name: %s%s%s\n", cli.Bold, valueOrNotSet(status.Name), cli.ColorReset)
	_, _ = fmt.Fprintf(out, "- Email: %s%s%s\n", cli.Bold, valueOrNotSet(status.Email), cli.ColorReset)
	_, _ = fmt.Fprintf(out, "- GPG Key: %s%s%s\n", cli.Bold, valueOrNotSet(status.SigningKey), cli.ColorReset)
	_, _ = fmt.Fprintf(out, "- GPG Signing: %s%s%s\n", cli.Bold, valueOrNotSet(status.GPGSign), cli.ColorReset)
}
