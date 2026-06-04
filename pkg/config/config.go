package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds the user-specific settings
type Config struct {
	GitHubUsername string `json:"github_username"`
	GitHubToken    string `json:"github_token"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Saves to ~/.git-wrap.json
	return filepath.Join(home, ".git-wrap.json"), nil
}

// Load reads the configuration file from the home directory
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err // Config doesn't exist yet
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the configuration file to the home directory
func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600) // Read/Write permissions for the owner only
}
