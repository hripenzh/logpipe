package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
sources:
  - name: app
    path: /var/log/app.log
min_level: info
format: pretty
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 1 || cfg.Sources[0].Name != "app" {
		t.Errorf("unexpected sources: %+v", cfg.Sources)
	}
	if cfg.Format != "pretty" {
		t.Errorf("expected format pretty, got %q", cfg.Format)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nope.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoSources(t *testing.T) {
	path := writeTemp(t, `sources: []\n`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty sources")
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	path := writeTemp(t, `
sources:
  - path: /tmp/app.log
format: json
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestLoad_DefaultsNameToPath(t *testing.T) {
	path := writeTemp(t, `
sources:
  - path: /tmp/app.log
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Sources[0].Name != "/tmp/app.log" {
		t.Errorf("expected name to default to path, got %q", cfg.Sources[0].Name)
	}
}
