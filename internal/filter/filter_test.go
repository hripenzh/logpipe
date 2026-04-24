package filter

import (
	"testing"
)

func TestFilter_NonJSONPassesThrough(t *testing.T) {
	line := []byte("plain text log line")
	if !Filter(line, Options{MinLevel: LevelError}) {
		t.Error("expected non-JSON line to pass through")
	}
}

func TestFilter_MinLevelFiltersLower(t *testing.T) {
	line := []byte(`{"level":"debug","msg":"verbose output"}`)
	if Filter(line, Options{MinLevel: LevelWarn}) {
		t.Error("expected debug line to be filtered out when min level is warn")
	}
}

func TestFilter_MinLevelAllowsEqual(t *testing.T) {
	line := []byte(`{"level":"warn","msg":"something degraded"}`)
	if !Filter(line, Options{MinLevel: LevelWarn}) {
		t.Error("expected warn line to pass when min level is warn")
	}
}

func TestFilter_MinLevelAllowsHigher(t *testing.T) {
	line := []byte(`{"level":"error","msg":"something broke"}`)
	if !Filter(line, Options{MinLevel: LevelInfo}) {
		t.Error("expected error line to pass when min level is info")
	}
}

func TestFilter_KeyContainsMatch(t *testing.T) {
	line := []byte(`{"level":"info","service":"auth-service","msg":"login"}`)
	opts := Options{
		KeyContains: map[string]string{"service": "auth"},
	}
	if !Filter(line, opts) {
		t.Error("expected line to pass key-contains filter")
	}
}

func TestFilter_KeyContainsMismatch(t *testing.T) {
	line := []byte(`{"level":"info","service":"payment-service","msg":"charge"}`)
	opts := Options{
		KeyContains: map[string]string{"service": "auth"},
	}
	if Filter(line, opts) {
		t.Error("expected line to be filtered out by key-contains")
	}
}

func TestFilter_KeyContainsMissingKey(t *testing.T) {
	line := []byte(`{"level":"info","msg":"no service field"}`)
	opts := Options{
		KeyContains: map[string]string{"service": "auth"},
	}
	if Filter(line, opts) {
		t.Error("expected line to be filtered out when key is absent")
	}
}

func TestFilter_UnknownLevelPassesThrough(t *testing.T) {
	line := []byte(`{"level":"trace","msg":"very verbose"}`)
	if !Filter(line, Options{MinLevel: LevelError}) {
		t.Error("expected unknown level to pass through")
	}
}
