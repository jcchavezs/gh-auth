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
	out := cmd.OutOrStdout()

	account, err := gitauth.GetCurrentAccount()
	if err != nil {
		if errors.Is(err, gitauth.ErrGHNotInstalled) {
			return err
		}
		return fmt.Errorf("failed to get GitHub account info: %w", err)
	}

	username, email := account.Username, account.Email

	if username == "" || email == "" {
		fprintln(out, "Could not get complete account information")
		fprintf(out, "Username: %s\n", username)
		fprintf(out, "Email: %s\n", email)
		fprintln(out)
		manualEmail, err := prompt(cmd, "Enter email manually: ")
		if err != nil {
			return err
		}
		email = manualEmail
	}

	gpgKey := gitauth.FindGPGKey(email)
	if gpgKey == "" {
		fprintf(out, "No GPG key found for %s\n", email)
		fprintln(out)
		fprintln(out, "Available GPG keys:")
		fprintln(out, gitauth.ListSecretKeys())
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
		fprintln(out, "GPG signing disabled")
	}

	displayKey := result.GPGKey
	if displayKey == "" {
		displayKey = "None"
	}
	fprintf(out, "%s✓%s Global GPG Key updated to %s%s (%s)%s\n",
		colorGreen, colorReset, bold, displayKey, result.Email, normal)

	return nil
}
