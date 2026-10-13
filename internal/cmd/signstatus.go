package cmd

import (
	"fmt"

	"github.com/fatih/color"
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

	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	switch {
	case !status.GHInstalled:
		_, _ = fmt.Fprintln(out, "GitHub CLI not installed")
	case status.Account != nil:
		_, _ = fmt.Fprintf(out, "%s Active GitHub account: %s\n",
			green("✓"), bold(fmt.Sprintf("%s <%s>", status.Account.Username, status.Account.Email)))
	default:
		_, _ = fmt.Fprintln(out, "No active GitHub account or GitHub CLI not authenticated")
	}

	_, _ = fmt.Fprintf(out, "- Name: %s\n", bold(valueOrNotSet(status.Name)))
	_, _ = fmt.Fprintf(out, "- Email: %s\n", bold(valueOrNotSet(status.Email)))
	_, _ = fmt.Fprintf(out, "- GPG Key: %s\n", bold(valueOrNotSet(status.SigningKey)))
	_, _ = fmt.Fprintf(out, "- GPG Signing: %s\n", bold(valueOrNotSet(status.GPGSign)))
}
