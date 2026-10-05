package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewRootCmdSubcommands(t *testing.T) {
	root := NewRootCmd()

	got := map[string]bool{}
	for _, c := range root.Commands() {
		got[c.Name()] = true
	}

	for _, name := range []string{"switch", "sign-status", "status"} {
		if !got[name] {
			t.Errorf("missing subcommand %q (have %v)", name, got)
		}
	}
}

func TestValueOrNotSet(t *testing.T) {
	if got := valueOrNotSet(""); got != "Not set" {
		t.Errorf(`valueOrNotSet("") = %q, want "Not set"`, got)
	}
	if got := valueOrNotSet("value"); got != "value" {
		t.Errorf(`valueOrNotSet("value") = %q, want "value"`, got)
	}
}

func TestPrompt(t *testing.T) {
	cmd := &cobra.Command{}
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetIn(strings.NewReader("  hello world  \n"))

	got, err := prompt(cmd, "Enter: ")
	if err != nil {
		t.Fatalf("prompt returned error: %v", err)
	}
	if got != "hello world" {
		t.Errorf("prompt value = %q, want %q", got, "hello world")
	}
	if out.String() != "Enter: " {
		t.Errorf("prompt message = %q, want %q", out.String(), "Enter: ")
	}
}

func TestPromptEOF(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetIn(strings.NewReader("")) // no newline, immediate EOF

	got, err := prompt(cmd, "")
	if err != nil {
		t.Fatalf("prompt returned error on EOF: %v", err)
	}
	if got != "" {
		t.Errorf("prompt value = %q, want empty", got)
	}
}
