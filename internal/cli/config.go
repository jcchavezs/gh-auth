package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/jcchavezs/gh-auth/internal/gitauth"
	"github.com/spf13/cobra"
)

// prompt reads a single line of input from the command's input stream after
// writing the given message to its output stream.
func prompt(cmd *cobra.Command, message string) (string, error) {
	fprint(cmd.OutOrStdout(), message)
	reader := bufio.NewReader(cmd.InOrStdin())
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// configureGitFromGitHub resolves the current GitHub account and matching GPG
// key (prompting the user when information is missing) and writes the global
// git configuration.
func configureGitFromGitHub(cmd *cobra.Command) error {
	ctx := cmd.Context()
	out := cmd.OutOrStdout()

	account, err := gitauth.GetCurrentAccount(ctx)
	if err != nil {
		if errors.Is(err, gitauth.ErrGHNotInstalled) {
			return err
		}
		return fmt.Errorf("failed to get GitHub account info: %w", err)
	}

	username, email := account.Username, account.Email

	if email == "" {
		fprintln(out, "Could not get email from GitHub account")
		fprintf(out, "Username: %s\n", username)
		fprintln(out)
		manualEmail, err := prompt(cmd, "Enter email manually: ")
		if err != nil {
			return err
		}
		email = manualEmail
	}

	gpgKey := gitauth.FindGPGKey(ctx, email)
	if gpgKey == "" {
		fprintf(out, "No GPG key found for %s\n", email)
		fprintln(out)
		fprintln(out, "Available GPG keys:")
		fprintln(out, gitauth.ListSecretKeys(ctx))
		fprintln(out)
		manualKey, err := prompt(cmd, "Enter GPG key ID (or press Enter to skip): ")
		if err != nil {
			return err
		}
		gpgKey = manualKey
	}

	result, err := gitauth.ConfigureGit(username, email, gpgKey)
	if err != nil {
		return err
	}

	if result.GPGKey == "" {
		fprintf(out, "%s✓%s Git identity updated: %s%s <%s>%s (no signing key)\n",
			colorGreen, colorReset, colorGreen, result.Username, result.Email, colorReset)
	} else {
		fprintf(out, "%s✓%s Git identity updated: %s%s <%s>%s, GPG key %s\n",
			colorGreen, colorReset, colorGreen, result.Username, result.Email, colorReset, result.GPGKey)
	}

	return nil
}
