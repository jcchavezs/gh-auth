package gitauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jcchavezs/gh-auth/internal/cli"
)

// profileGitConfigFile is the name of the file, stored inside a profile's
// config directory, that caches the git identity resolved from GitHub so it can
// be reused across shells without re-querying GitHub.
const profileGitConfigFile = "gitconfig.json"

// ProfileGitConfig holds the git identity persisted for a profile.
type ProfileGitConfig struct {
	Name       string `json:"user.name"`
	Email      string `json:"user.email"`
	SigningKey string `json:"user.signingkey"`
	GPGSign    string `json:"commit.gpgsign"`
}

// profileGitConfigPath returns the path of the cached git config file for the
// profile whose gh config lives in configDir.
func profileGitConfigPath(configDir string) string {
	return filepath.Join(configDir, profileGitConfigFile)
}

// ReadProfileGitConfig reads the cached git identity for the profile whose gh
// config lives in configDir. It returns an error wrapping os.ErrNotExist when
// the file has not been written yet.
func ReadProfileGitConfig(configDir string) (ProfileGitConfig, bool, error) {
	data, err := os.ReadFile(profileGitConfigPath(configDir))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ProfileGitConfig{}, false, nil
		}

		return ProfileGitConfig{}, false, err
	}

	var cfg ProfileGitConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ProfileGitConfig{}, false, fmt.Errorf("parsing %s: %w", profileGitConfigFile, err)
	}
	return cfg, true, nil
}

// SaveProfileGitConfig persists the given identity for the profile whose gh
// config lives in configDir. commit.gpgsign is inherited from the global git
// config.
func SaveProfileGitConfig(configDir, username, email, signingKey string) (ProfileGitConfig, error) {
	cfg := ProfileGitConfig{
		Name:       username,
		Email:      email,
		SigningKey: signingKey,
		GPGSign:    gitConfigGet("commit.gpgsign"),
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return ProfileGitConfig{}, err
	}
	if err := os.WriteFile(profileGitConfigPath(configDir), data, 0o600); err != nil {
		return ProfileGitConfig{}, fmt.Errorf("writing %s: %w", profileGitConfigFile, err)
	}
	return cfg, nil
}

func GetAccountFromGithub(ctx context.Context, out io.Writer, in io.Reader) (ProfileGitConfig, error) {
	account, err := GetCurrentAccount(ctx)
	if err != nil {
		if errors.Is(err, ErrGHNotInstalled) {
			return ProfileGitConfig{}, err
		}

		return ProfileGitConfig{}, fmt.Errorf("failed to get GitHub account info: %w", err)
	}

	username, email := account.Username, account.Email
	if email == "" {
		_, _ = fmt.Fprintln(out, "Could not get email from GitHub account")
		_, _ = fmt.Fprintf(out, "Username: %s\n", username)
		_, _ = fmt.Fprintln(out)
		manualEmail, err := cli.Prompt(out, in, "Enter email manually: ")
		if err != nil {
			return ProfileGitConfig{}, err
		}
		email = manualEmail
	}

	gpgKey := FindGPGKey(ctx, email)
	if gpgKey == "" {
		_, _ = fmt.Fprintf(out, "No GPG key found for %s\n", email)
		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintln(out, "Available GPG keys:")
		_, _ = fmt.Fprintln(out, ListSecretKeys(ctx))
		_, _ = fmt.Fprintln(out)
		manualKey, err := cli.Prompt(out, in, "Enter GPG key ID (or press Enter to skip): ")
		if err != nil {
			return ProfileGitConfig{}, err
		}
		gpgKey = manualKey
	}

	return ProfileGitConfig{Name: username, Email: email, SigningKey: gpgKey}, nil
}
