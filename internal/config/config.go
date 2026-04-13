// Package config loads vaultdiff configuration from a YAML file and environment
// variables, providing sensible defaults.
package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for vaultdiff.
type Config struct {
	VaultAddr  string   `yaml:"vault_addr"`
	AuthMethod string   `yaml:"auth_method"`
	KVMounts   []string `yaml:"kv_mounts"`
	Redact     string   `yaml:"redact"`
	Output     string   `yaml:"output"`
	ShowAll    bool     `yaml:"show_all"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		VaultAddr:  "http://127.0.0.1:8200",
		AuthMethod: "token",
		KVMounts:   []string{"secret"},
		Redact:     "none",
		Output:     "text",
		ShowAll:    false,
	}
}

// Load reads a YAML config file at path and merges it over the defaults.
// If path is empty the defaults are returned unchanged.
// Environment variables take final precedence over file values.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		applyEnv(&cfg)
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, errors.New("config file not found: " + path)
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	applyEnv(&cfg)
	return cfg, nil
}

// applyEnv overrides config fields with environment variable values when set.
func applyEnv(cfg *Config) {
	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULTDIFF_AUTH_METHOD"); v != "" {
		cfg.AuthMethod = v
	}
	if v := os.Getenv("VAULTDIFF_REDACT"); v != "" {
		cfg.Redact = v
	}
	if v := os.Getenv("VAULTDIFF_OUTPUT"); v != "" {
		cfg.Output = v
	}
}
