package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected VaultAddr: %q", cfg.VaultAddr)
	}
	if cfg.AuthMethod != "token" {
		t.Errorf("unexpected AuthMethod: %q", cfg.AuthMethod)
	}
	if cfg.Redact != "none" {
		t.Errorf("unexpected Redact: %q", cfg.Redact)
	}
	if cfg.Output != "text" {
		t.Errorf("unexpected Output: %q", cfg.Output)
	}
	if cfg.ShowAll {
		t.Error("expected ShowAll to be false by default")
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected VaultAddr: %q", cfg.VaultAddr)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	content := `vault_addr: https://vault.example.com
auth_method: approle
redact: mask
output: json
show_all: true
kv_mounts:
  - secret
  - kv
`
	path := writeTempFile(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "https://vault.example.com" {
		t.Errorf("unexpected VaultAddr: %q", cfg.VaultAddr)
	}
	if cfg.AuthMethod != "approle" {
		t.Errorf("unexpected AuthMethod: %q", cfg.AuthMethod)
	}
	if cfg.Redact != "mask" {
		t.Errorf("unexpected Redact: %q", cfg.Redact)
	}
	if !cfg.ShowAll {
		t.Error("expected ShowAll to be true")
	}
	if len(cfg.KVMounts) != 2 {
		t.Errorf("expected 2 kv_mounts, got %d", len(cfg.KVMounts))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTempFile(t, ": invalid: yaml: [")
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv("VAULT_ADDR", "https://env.vault.com")
	t.Setenv("VAULTDIFF_REDACT", "redact")
	t.Setenv("VAULTDIFF_OUTPUT", "json")
	t.Setenv("VAULTDIFF_AUTH_METHOD", "kubernetes")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "https://env.vault.com" {
		t.Errorf("expected env override for VaultAddr, got %q", cfg.VaultAddr)
	}
	if cfg.Redact != "redact" {
		t.Errorf("expected env override for Redact, got %q", cfg.Redact)
	}
	if cfg.Output != "json" {
		t.Errorf("expected env override for Output, got %q", cfg.Output)
	}
	if cfg.AuthMethod != "kubernetes" {
		t.Errorf("expected env override for AuthMethod, got %q", cfg.AuthMethod)
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}
