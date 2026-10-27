package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/jcchavezs/gh-auth/internal/cli"
	"github.com/jcchavezs/gh-auth/internal/gitauth"
	"github.com/jcchavezs/gh-auth/internal/shell"
	"github.com/spf13/cobra"
)

// profileConfigured reports whether configDir is an existing directory that
// already holds a gh configuration, i.e. it contains both config.yml and
// hosts.yml.
func profileConfigured(configDir string) bool {
	if info, err := os.Stat(configDir); err != nil || !info.IsDir() {
		return false
	}
	for _, name := range []string{"config.yml", "hosts.yml"} {
		if _, err := os.Stat(filepath.Join(configDir, name)); err != nil {
			return false
		}
	}
	return true
}

// newShellCmd creates the 'shell' command that starts a shell with GH_CONFIG_DIR
// pointed at the specified profile's config directory.
func newShellCmd() *cobra.Command {
	var refreshProfile bool

	cmd := &cobra.Command{
		Use:   "shell [profile] [-- command [args...]]",
		Short: "Start a shell (or run a command) with gh pointed at a profile's config directory",
		Long: `Start a shell with GH_CONFIG_DIR pointed at the specified profile's config directory, allowing you to work with multiple GitHub profiles seamlessly.

If a command is given after "--", it is run instead of an interactive shell,
e.g. gh-auth shell work -- git commit -m "msg".

Based on https://github.com/cli/cli/issues/12145
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if dash := cmd.ArgsLenAtDash(); dash >= 0 {
				args = args[:dash]
			}
			return cobra.MaximumNArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var command []string
			if dash := cmd.ArgsLenAtDash(); dash >= 0 {
				command = args[dash:]
				args = args[:dash]
			}

			var profile string
			if len(args) == 1 {
				profile = args[0]
			} else {
				entered, err := cli.Prompt(cmd.OutOrStdout(), cmd.InOrStdin(), "Profile name: ")
				if err != nil {
					return err
				}
				profile = entered
			}

			if profile == "" {
				return errors.New("profile name cannot be empty")
			}

			configDir, err := gitauth.ProfileConfigDir(profile)
			if err != nil {
				return err
			}

			if !profileConfigured(configDir) {
				err := os.MkdirAll(configDir, 0755)
				if err != nil {
					return err
				}

				if err := gitauth.AuthLogin(cmd.Context(), configDir); err != nil {
					return err
				}
			}

			return shell.RunShell(cmd.Context(), profile, refreshProfile, command, cmd.OutOrStdout(), cmd.InOrStdin())
		},
	}

	cmd.Flags().BoolVar(&refreshProfile, "refresh-profile", false,
		"Re-resolve the git identity from GitHub and refresh the profile's cached config")

	return cmd
}
