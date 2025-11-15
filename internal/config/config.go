package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Theme    string `json:"theme"`
	Sound    bool   `json:"sound"`
	BlindMode bool  `json:"blind_mode"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Theme:    "catppuccin-mocha",
		Sound:    false,
		BlindMode: false,
	}
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Default(), nil
	}

	configDir := filepath.Join(homeDir, ".keyarch")
	configPath := filepath.Join(configDir, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), nil
	}

	return &cfg, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".keyarch")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
