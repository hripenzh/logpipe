package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the logpipe binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	bin := filepath.Join(tmpDir, "logpipe")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return bin
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logpipe-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestMain_MissingConfig(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--config", "/nonexistent/logpipe.yaml")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config")
	}
	if !strings.Contains(string(out), "failed to load config") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestMain_NoSourcesConfig(t *testing.T) {
	bin := buildBinary(t)
	cfgPath := writeTempConfig(t, "sources: []\noutput:\n  format: raw\n")
	cmd := exec.Command(bin, "--config", cfgPath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for empty sources")
	}
	if !strings.Contains(string(out), "no sources") {
		t.Errorf("expected 'no sources' error, got: %s", out)
	}
}
