package chatapp

import "testing"

func TestParseArgsRequiresName(t *testing.T) {
	t.Parallel()

	_, err := ParseArgs(nil)
	if err != ErrNameRequired {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestParseArgsUsesDefaults(t *testing.T) {
	t.Parallel()

	cfg, err := ParseArgs([]string{"-name", "Alice"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if cfg.Name != "Alice" {
		t.Fatalf("unexpected name: %q", cfg.Name)
	}

	if cfg.ListenAddr != defaultListenAddr {
		t.Fatalf("unexpected listen address: %q", cfg.ListenAddr)
	}

	if cfg.Mode() != "server" {
		t.Fatalf("unexpected mode: %q", cfg.Mode())
	}
}

func TestParseArgsTrimsName(t *testing.T) {
	t.Parallel()

	cfg, err := ParseArgs([]string{"-name", "  Alice  "})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if cfg.Name != "Alice" {
		t.Fatalf("unexpected name: %q", cfg.Name)
	}
}

func TestParseArgsRejectsWhitespaceOnlyName(t *testing.T) {
	t.Parallel()

	_, err := ParseArgs([]string{"-name", "   "})
	if err != ErrNameRequired {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestParseArgsClientMode(t *testing.T) {
	t.Parallel()

	cfg, err := ParseArgs([]string{"-name", "Bob", "-peer", "127.0.0.1:50051"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if cfg.Mode() != "client" {
		t.Fatalf("unexpected mode: %q", cfg.Mode())
	}
}
