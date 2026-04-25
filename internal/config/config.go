// Package config provides configuration loading and validation for logpipe.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Source defines a single log source to tail.
type Source struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// Config holds the top-level logpipe configuration.
type Config struct {
	Sources  []Source `yaml:"sources"`
	MinLevel string   `yaml:"min_level"`
	Format   string   `yaml:"format"`
	KeyFilter string  `yaml:"key_filter"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks that required fields are present and values are acceptable.
func (c *Config) validate() error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("config: at least one source is required")
	}
	for i, s := range c.Sources {
		if s.Path == "" {
			return fmt.Errorf("config: source[%d] missing path", i)
		}
		if s.Name == "" {
			c.Sources[i].Name = s.Path
		}
	}
	validFormats := map[string]bool{"raw": true, "pretty": true, "": true}
	if !validFormats[c.Format] {
		return fmt.Errorf("config: invalid format %q (must be raw or pretty)", c.Format)
	}
	return nil
}
