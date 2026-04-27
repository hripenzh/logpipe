// Package config handles loading and validating logpipe configuration files.
package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Source represents a single log source defined in the config file.
type Source struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	Command string `yaml:"command"`
}

// Config is the top-level configuration structure.
type Config struct {
	Sources  []Source `yaml:"sources"`
	MinLevel string   `yaml:"min_level"`
	Format   string   `yaml:"format"`
	// RateLimit caps lines forwarded per second (0 = unlimited).
	RateLimit int `yaml:"rate_limit"`
}

// Load reads and validates a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: cannot read file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: invalid YAML: %w", err)
	}

	if len(cfg.Sources) == 0 {
		return nil, errors.New("config: at least one source must be defined")
	}

	if cfg.Format == "" {
		cfg.Format = "pretty"
	}

	if cfg.RateLimit < 0 {
		return nil, errors.New("config: rate_limit must be zero or positive")
	}

	for i, s := range cfg.Sources {
		if s.Name == "" {
			return nil, fmt.Errorf("config: source[%d] missing name", i)
		}
		if s.Path == "" && s.Command == "" {
			return nil, fmt.Errorf("config: source %q must specify path or command", s.Name)
		}
	}

	return &cfg, nil
}
