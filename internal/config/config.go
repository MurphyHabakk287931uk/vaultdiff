package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultdiff configuration.
type Config struct {
	Vault  VaultConfig  `yaml:"vault"`
	Diff   DiffConfig   `yaml:"diff"`
	Output OutputConfig `yaml:"output"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address string   `yaml:"address"`
	Token   string   `yaml:"token"`
	Mounts  []string `yaml:"mounts"`
}

// DiffConfig holds diff behaviour settings.
type DiffConfig struct {
	RedactMode    string `yaml:"redact_mode"`
	ShowUnchanged bool   `yaml:"show_unchanged"`
}

// OutputConfig holds output formatting settings.
type OutputConfig struct {
	Format string `yaml:"format"`
	Color  bool   `yaml:"color"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Vault: VaultConfig{
			Address: "http://127.0.0.1:8200",
			Mounts:  []string{"secret"},
		},
		Diff: DiffConfig{
			RedactMode:    "none",
			ShowUnchanged: false,
		},
		Output: OutputConfig{
			Format: "text",
			Color:  true,
		},
	}
}

// Load reads a YAML config file from path and merges it over defaults.
// If path is empty the default config is returned.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	return cfg, nil
}
