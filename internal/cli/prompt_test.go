package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrompt(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("  hello world  \n")

	got, err := Prompt(&out, in, "Enter: ")
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
	var out bytes.Buffer
	in := strings.NewReader("") // no newline, immediate EOF

	got, err := Prompt(&out, in, "")
	if err != nil {
		t.Fatalf("prompt returned error on EOF: %v", err)
	}
	if got != "" {
		t.Errorf("prompt value = %q, want empty", got)
	}
}
