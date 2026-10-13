package cmd

import "testing"

func TestValueOrNotSet(t *testing.T) {
	if got := valueOrNotSet(""); got != "Not set" {
		t.Errorf(`valueOrNotSet("") = %q, want "Not set"`, got)
	}
	if got := valueOrNotSet("value"); got != "value" {
		t.Errorf(`valueOrNotSet("value") = %q, want "value"`, got)
	}
}
