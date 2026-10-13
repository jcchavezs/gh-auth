package cmd

import "testing"

func TestNewRootCmdSubcommands(t *testing.T) {
	root := NewRootCmd()

	got := map[string]bool{}
	for _, c := range root.Commands() {
		got[c.Name()] = true
	}

	for _, name := range []string{"switch", "sign-status", "status", "shell"} {
		if !got[name] {
			t.Errorf("missing subcommand %q (have %v)", name, got)
		}
	}
}
