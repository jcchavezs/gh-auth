package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/jcchavezs/gh-auth/internal/gitauth"
)

// ErrShellNotSet is returned when the SHELL environment variable is not set.
var ErrShellNotSet = errors.New("$SHELL is not set")

// RunShell starts $SHELL with GH_CONFIG_DIR pointed at configDir, forwarding the
// current process stdio so the shell is interactive.
//
// The git identity is read from the profile's cached gitconfig.json when
// available; when it is missing or refreshProfile is true, it is resolved from
// GitHub and cached for subsequent shells. When the profile has a signing key,
// commit signing is forced through GIT_CONFIG_* so it applies to every git
// invocation regardless of how global options are ordered on the command line.
//
// When command is non-empty it is run instead of an interactive shell and the
// informational banner is skipped so the command's output stays clean.
func RunShell(ctx context.Context, profile string, refreshProfile bool, command []string, out io.Writer, in io.Reader) error {
	if os.Getenv("GHAUTH_PROFILE") != "" {
		return errors.New("profile is already enabled")
	}

	var name string
	var cmdArgs []string
	if len(command) > 0 {
		name, cmdArgs = command[0], command[1:]
	} else {
		name = os.Getenv("SHELL")
		if name == "" {
			return ErrShellNotSet
		}
	}

	configDir, err := gitauth.ProfileConfigDir(profile)
	if err != nil {
		return err
	}

	if err := os.Setenv("GH_CONFIG_DIR", configDir); err != nil {
		return err
	}

	var username, email, gpgKey string

	cfg, found, err := gitauth.ReadProfileGitConfig(configDir)
	if err != nil {
		return err
	}

	if found && !refreshProfile {
		username, email, gpgKey = cfg.Name, cfg.Email, cfg.SigningKey
	} else {
		account, err := gitauth.GetAccountFromGithub(ctx, out, in)
		if err != nil {
			return err
		}
		if _, err := gitauth.SaveProfileGitConfig(configDir, account.Name, account.Email, account.SigningKey); err != nil {
			return err
		}
		username, email, gpgKey = account.Name, account.Email, account.SigningKey
	}

	cmd := exec.CommandContext(ctx, name, cmdArgs...)
	cmd.Env = append(os.Environ(),
		"GH_CONFIG_DIR="+configDir,
		"GHAUTH_PROFILE="+profile,
		"GIT_AUTHOR_NAME="+username,
		"GIT_AUTHOR_EMAIL="+email,
		"GIT_COMMITTER_NAME="+username,
		"GIT_COMMITTER_EMAIL="+email,
	)
	if gpgKey != "" {
		// Force signing via GIT_CONFIG_* so it applies to every commit without
		// rewriting the command line (which breaks on `git <global-opts> commit`).
		cmd.Env = append(cmd.Env,
			"GIT_CONFIG_COUNT=2",
			"GIT_CONFIG_KEY_0=user.signingkey",
			"GIT_CONFIG_VALUE_0="+gpgKey,
			"GIT_CONFIG_KEY_1=commit.gpgsign",
			"GIT_CONFIG_VALUE_1=true",
		)
	}

	if len(command) == 0 {
		// Write a message indicating the shell is starting and the profile being used as well as the name, email and signing key.
		bold := color.New(color.Bold).SprintFunc()
		_, _ = fmt.Fprintf(out, "Starting shell with profile %s.\n\n", bold(profile))
		_, _ = fmt.Fprintf(out, "%s%s\n", bold("Git author/committer name: "), username)
		_, _ = fmt.Fprintf(out, "%s%s\n", bold("Git author/committer email: "), email)
		if gpgKey != "" {
			_, _ = fmt.Fprintf(out, "%s%s\n", bold("Git signing key: "), gpgKey)
		}
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
