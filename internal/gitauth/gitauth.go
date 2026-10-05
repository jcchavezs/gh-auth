// Package gitauth provides the core logic for inspecting the active GitHub
// account and configuring git globally based on it. It mirrors the behaviour
// of the original gh-auth bash script.
package gitauth

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// ErrGHNotInstalled is returned when the GitHub CLI (gh) is not available.
var ErrGHNotInstalled = errors.New("GitHub CLI not installed")

// Account holds the identity of the currently authenticated GitHub account.
type Account struct {
	Username string
	Email    string
}

// commandExists reports whether the given executable is in PATH.
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// GetCurrentAccount returns the login and email of the currently authenticated
// GitHub account using the gh CLI.
func GetCurrentAccount(ctx context.Context) (Account, error) {
	if !commandExists("gh") {
		return Account{}, ErrGHNotInstalled
	}

	username, err := ghAPIField(ctx, "user", ".login")
	if err != nil {
		return Account{}, fmt.Errorf("getting GitHub username: %w", err)
	}

	email, err := ghAPIField(ctx, "user", ".email")
	if err != nil {
		// A missing/private email is not fatal; the caller can prompt for it.
		email = ""
	}
	if email == "null" {
		email = ""
	}

	// Strip any "+tag" subaddress so e.g. "user+alias@domain.com" becomes
	// "user@domain.com".
	email = normalizeEmail(email)

	return Account{Username: username, Email: email}, nil
}

// emailSubaddress matches a "+tag" subaddress between the local part and the
// "@" of an email address.
var emailSubaddress = regexp.MustCompile(`\+[^@]*@`)

// normalizeEmail removes a "+tag" subaddress from an email address, keeping
// everything after the "@" intact.
func normalizeEmail(email string) string {
	return emailSubaddress.ReplaceAllString(email, "@")
}

// ghAPIField runs `gh api <endpoint> --jq <jq>` and returns the trimmed output.
func ghAPIField(ctx context.Context, endpoint, jq string) (string, error) {
	out, err := exec.CommandContext(ctx, "gh", "api", endpoint, "--jq", jq).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

var gpgKeyWithPrefix = regexp.MustCompile(`^[a-z0-9]+/([A-Fa-f0-9]+)$`)

// cleanGPGKey removes a common algorithm prefix (e.g. "rsa4096/") from a key id.
func cleanGPGKey(key string) string {
	if m := gpgKeyWithPrefix.FindStringSubmatch(key); m != nil {
		return m[1]
	}
	return key
}

// FindGPGKey looks for a secret GPG key whose uid matches the given email and
// returns the last 16 characters of its key id (the long key id). It returns an
// empty string when gpg is unavailable or no matching key is found.
func FindGPGKey(ctx context.Context, email string) string {
	if email == "" || !commandExists("gpg") {
		return ""
	}

	out, err := exec.CommandContext(ctx, "gpg", "--list-secret-keys", "--with-colons").Output()
	if err != nil {
		return ""
	}

	return parseSecretKeyForEmail(string(out), email)
}

// parseSecretKeyForEmail scans gpg colon-delimited secret key output and returns
// the long (16-char) key id of the first key whose uid contains email, or "".
func parseSecretKeyForEmail(colonOutput, email string) string {
	var currentKey string
	scanner := bufio.NewScanner(strings.NewReader(colonOutput))
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ":")
		switch fields[0] {
		case "sec":
			if len(fields) > 4 {
				currentKey = fields[4]
			}
		case "uid":
			if len(fields) > 9 && strings.Contains(fields[9], email) {
				if len(currentKey) > 16 {
					return currentKey[len(currentKey)-16:]
				}
				return currentKey
			}
		}
	}

	return ""
}

// runGit runs a git command, forwarding its stderr into the returned error when
// it fails so the caller surfaces git's own diagnostics.
func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return fmt.Errorf("%w: %s", err, msg)
		}
		return err
	}
	return nil
}

// gitConfigSet sets a value in the global git config.
func gitConfigSet(key, value string) error {
	return runGit("config", "--global", key, value)
}

// gitConfigUnset unsets a value in the global git config, ignoring the error
// returned when the key is not set.
func gitConfigUnset(key string) {
	_ = runGit("config", "--global", "--unset", key)
}

// gitConfigGet returns a value from the global git config, or "" when unset.
func gitConfigGet(key string) string {
	out, err := exec.Command("git", "config", "--global", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}


// ListSecretKeys returns a human-readable listing of available secret GPG keys,
// limited to the sec/uid lines, matching the script's fallback output.
func ListSecretKeys(ctx context.Context) string {
	if !commandExists("gpg") {
		return "No GPG keys found"
	}
	out, err := exec.CommandContext(ctx, "gpg", "--list-secret-keys", "--keyid-format", "LONG").Output()
	if err != nil {
		return "No GPG keys found"
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "sec") || strings.HasPrefix(trimmed, "uid") {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return "No GPG keys found"
	}
	return strings.Join(lines, "\n")
}

// ConfigResult describes the outcome of configuring git from the GitHub account.
type ConfigResult struct {
	Username string
	Email    string
	GPGKey   string
}

// ConfigureGit configures global git user.name, user.email and GPG signing
// based on the provided account and GPG key. When gpgKey is empty, signing is
// disabled. The caller is responsible for resolving the account and key
// (including any interactive prompting).
func ConfigureGit(username, email, gpgKey string) (ConfigResult, error) {
	gpgKey = strings.TrimSpace(gpgKey)
	if gpgKey != "" {
		gpgKey = cleanGPGKey(gpgKey)
	}

	if err := gitConfigSet("user.name", username); err != nil {
		return ConfigResult{}, fmt.Errorf("setting user.name: %w", err)
	}
	if err := gitConfigSet("user.email", email); err != nil {
		return ConfigResult{}, fmt.Errorf("setting user.email: %w", err)
	}

	if gpgKey != "" {
		if err := gitConfigSet("user.signingkey", gpgKey); err != nil {
			return ConfigResult{}, fmt.Errorf("setting user.signingkey: %w", err)
		}
		if err := gitConfigSet("commit.gpgsign", "true"); err != nil {
			return ConfigResult{}, fmt.Errorf("enabling commit.gpgsign: %w", err)
		}
	} else {
		gitConfigUnset("commit.gpgsign")
		gitConfigUnset("user.signingkey")
	}

	return ConfigResult{Username: username, Email: email, GPGKey: gpgKey}, nil
}

// SwitchAccount runs `gh auth switch`, which shows gh's interactive menu.
func SwitchAccount(ctx context.Context) error {
	if !commandExists("gh") {
		return ErrGHNotInstalled
	}
	cmd := exec.CommandContext(ctx, "gh", "auth", "switch")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AuthStatus runs `gh auth status`, forwarding its output to the given streams.
func AuthStatus(ctx context.Context, stdout, stderr io.Writer) error {
	if !commandExists("gh") {
		return ErrGHNotInstalled
	}
	cmd := exec.CommandContext(ctx, "gh", "auth", "status")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

// SignStatus captures the current git signing configuration.
type SignStatus struct {
	Account     *Account
	Name        string
	Email       string
	SigningKey  string
	GPGSign     string
	GHInstalled bool
	AccountErr  error
}

// GetSignStatus gathers the current GitHub account and global git signing
// config.
func GetSignStatus(ctx context.Context) SignStatus {
	status := SignStatus{
		GHInstalled: commandExists("gh"),
		Name:        gitConfigGet("user.name"),
		Email:       gitConfigGet("user.email"),
		SigningKey:  gitConfigGet("user.signingkey"),
		GPGSign:     gitConfigGet("commit.gpgsign"),
	}

	if status.GHInstalled {
		acc, err := GetCurrentAccount(ctx)
		if err != nil {
			status.AccountErr = err
		} else {
			status.Account = &acc
		}
	}

	return status
}
