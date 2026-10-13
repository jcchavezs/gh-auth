package gitauth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfileGitConfigPath(t *testing.T) {
	got := profileGitConfigPath("/tmp/profile")
	want := filepath.Join("/tmp/profile", profileGitConfigFile)
	if got != want {
		t.Errorf("profileGitConfigPath = %q, want %q", got, want)
	}
}

func TestReadProfileGitConfig(t *testing.T) {
	t.Run("missing file reports not found without error", func(t *testing.T) {
		dir := t.TempDir()

		cfg, found, err := ReadProfileGitConfig(dir)
		if err != nil {
			t.Fatalf("ReadProfileGitConfig returned error: %v", err)
		}
		if found {
			t.Errorf("found = true, want false")
		}
		if cfg != (ProfileGitConfig{}) {
			t.Errorf("cfg = %+v, want zero value", cfg)
		}
	})

	t.Run("valid file is parsed", func(t *testing.T) {
		dir := t.TempDir()
		const content = `{
  "user.name": "Jane Doe",
  "user.email": "jane@example.com",
  "user.signingkey": "ABCDEF0123456789",
  "commit.gpgsign": "true"
}`
		if err := os.WriteFile(filepath.Join(dir, profileGitConfigFile), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}

		cfg, found, err := ReadProfileGitConfig(dir)
		if err != nil {
			t.Fatalf("ReadProfileGitConfig returned error: %v", err)
		}
		if !found {
			t.Fatalf("found = false, want true")
		}
		want := ProfileGitConfig{
			Name:       "Jane Doe",
			Email:      "jane@example.com",
			SigningKey: "ABCDEF0123456789",
			GPGSign:    "true",
		}
		if cfg != want {
			t.Errorf("cfg = %+v, want %+v", cfg, want)
		}
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, profileGitConfigFile), []byte("not json"), 0o600); err != nil {
			t.Fatal(err)
		}

		_, found, err := ReadProfileGitConfig(dir)
		if err == nil {
			t.Fatalf("ReadProfileGitConfig returned nil error, want parse error")
		}
		if found {
			t.Errorf("found = true, want false")
		}
	})
}

func TestSaveProfileGitConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()

	saved, err := SaveProfileGitConfig(dir, "Jane Doe", "jane@example.com", "ABCDEF0123456789")
	if err != nil {
		t.Fatalf("SaveProfileGitConfig returned error: %v", err)
	}
	if saved.Name != "Jane Doe" || saved.Email != "jane@example.com" || saved.SigningKey != "ABCDEF0123456789" {
		t.Errorf("saved identity = %+v, want name/email/signingkey populated", saved)
	}

	read, found, err := ReadProfileGitConfig(dir)
	if err != nil {
		t.Fatalf("ReadProfileGitConfig returned error: %v", err)
	}
	if !found {
		t.Fatalf("found = false after save, want true")
	}
	if read != saved {
		t.Errorf("round-trip mismatch: read %+v, want %+v", read, saved)
	}
}
