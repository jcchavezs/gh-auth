package gitauth

import "testing"

func TestCleanGPGKey(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no prefix", "C9FA703B07F5B3FD", "C9FA703B07F5B3FD"},
		{"rsa prefix", "rsa4096/C9FA703B07F5B3FD", "C9FA703B07F5B3FD"},
		{"ed25519 prefix", "ed25519/ABCDEF0123456789", "ABCDEF0123456789"},
		{"empty", "", ""},
		{"unrelated slash", "foo/bar", "foo/bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanGPGKey(tt.in); got != tt.want {
				t.Errorf("cleanGPGKey(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no subaddress", "user@domain.com", "user@domain.com"},
		{"spam subaddress", "user+spam@domain.com", "user@domain.com"},
		{"arbitrary subaddress", "foo+anything@gmail.com", "foo@gmail.com"},
		{"empty tag", "foo+@gmail.com", "foo@gmail.com"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeEmail(tt.in); got != tt.want {
				t.Errorf("normalizeEmail(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseSecretKeyForEmail(t *testing.T) {
	const colons = `sec:u:4096:1:ABCDEF0123456789:1700000000::u:::scESC::::::23::0:
fpr:::::::::1234567890ABCDEF1234567890ABCDEF01234567:
uid:u::::1700000000::HASH::Jane Doe <jane@example.com>::::::::::0:
sec:u:4096:1:0011223344556677:1700000000::u:::scESC::::::23::0:
uid:u::::1700000000::HASH2::John Roe <john+tag@domain.com>::::::::::0:
`

	tests := []struct {
		name  string
		email string
		want  string
	}{
		{"exact match returns long id", "jane@example.com", "ABCDEF0123456789"},
		{"second key by substring", "john+tag@domain.com", "0011223344556677"},
		{"partial domain match", "domain.com", "0011223344556677"},
		{"no match", "nobody@nowhere.com", ""},
		{"empty", "", "ABCDEF0123456789"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSecretKeyForEmail(colons, tt.email); got != tt.want {
				t.Errorf("parseSecretKeyForEmail(email=%q) = %q, want %q", tt.email, got, tt.want)
			}
		})
	}
}

func TestParseSecretKeyForEmailShortKey(t *testing.T) {
	// A key id of 16 or fewer chars should be returned verbatim.
	const colons = `sec:u:4096:1:SHORTKEY12345678:1700000000::u:::scESC::::::23::0:
uid:u::::1700000000::HASH::Jane Doe <jane@example.com>::::::::::0:
`
	if got := parseSecretKeyForEmail(colons, "jane@example.com"); got != "SHORTKEY12345678" {
		t.Errorf("got %q, want %q", got, "SHORTKEY12345678")
	}
}
