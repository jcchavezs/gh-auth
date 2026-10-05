package cli

import (
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
	status := gitauth.GetSignStatus()

	switch {
	case !status.GHInstalled:
		fprintln(out, "GitHub CLI not installed")
	case status.Account != nil:
		fprintf(out, "%s✓%s Active GitHub account: %s%s <%s>%s\n",
			colorGreen, colorReset, bold, status.Account.Username, status.Account.Email, normal)
	default:
		fprintln(out, "No active GitHub account or GitHub CLI not authenticated")
	}

	fprintf(out, "- Name: %s%s%s\n", bold, valueOrNotSet(status.Name), normal)
	fprintf(out, "- Email: %s%s%s\n", bold, valueOrNotSet(status.Email), normal)
	fprintf(out, "- GPG Key: %s%s%s\n", bold, valueOrNotSet(status.SigningKey), normal)
	fprintf(out, "- GPG Signing: %s%s%s\n", bold, valueOrNotSet(status.GPGSign), normal)
}
