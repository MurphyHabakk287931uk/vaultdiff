package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultdiff/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("unexpected default address: %s", cfg.Vault.Address)
	}
	if cfg.Diff.RedactMode != "none" {
		t.Errorf("unexpected default redact_mode: %s", cfg.Diff.RedactMode)
	}
	if cfg.Output.Format != "text" {
		t.Errorf("unexpected default format: %s", cfg.Output.Format)
	}
	if !cfg.Output.Color {
		t.Error("expected color to be true by default")
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	content := `
vault:
  address: https://vault.example.com
  token: s.test
  mounts:
    - kv
    - secret
diff:
  redact_mode: redact
  show_unchanged: true
output:
  format: json
  color: false
`
	path := writeTemp(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("got address %q", cfg.Vault.Address)
	}
	if cfg.Diff.RedactMode != "redact" {
		t.Errorf("got redact_mode %q", cfg.Diff.RedactMode)
	}
	if !cfg.Diff.ShowUnchanged {
		t.Error("expected show_unchanged true")
	}
	if cfg.Output.Format != "json" {
		t.Errorf("got format %q", cfg.Output.Format)
	}
	if len(cfg.Vault.Mounts) != 2 {
		t.Errorf("expected 2 mounts, got %d", len(cfg.Vault.Mounts))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ": invalid: yaml: [")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return path
}
