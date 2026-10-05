//go:build e2e

// Package e2e contains end-to-end tests that exercise the built gh-auth binary
// against real git (and, when available, gh) commands. Run with:
//
//	make test-e2e
//
// or:
//
//	go test -tags e2e ./e2e/...
package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// binPath builds the gh-auth binary once and returns its path.
func binPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "gh-auth")
	cmd := exec.Command("go", "build", "-o", bin, "github.com/jcchavezs/gh-auth/cmd/gh-auth")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("building binary: %v\n%s", err, out)
	}
	return bin
}

func TestStatusForwardsGhAuthStatus(t *testing.T) {
	if _, err := exec.LookPath("gh"); err != nil {
		t.Skip("gh not installed; skipping status e2e")
	}
	// gh authenticates from GH_TOKEN/GITHUB_TOKEN. Without it, `gh auth status`
	// reports no account and exits non-zero, so skip rather than fail.
	if os.Getenv("GH_TOKEN") == "" && os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("no GH_TOKEN/GITHUB_TOKEN; skipping status e2e")
	}

	bin := binPath(t)
	cmd := exec.Command(bin, "status")
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("status failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "github.com") {
		t.Errorf("expected status output to mention github.com, got:\n%s", out)
	}
}
