package formatter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/formatter"
)

func TestFormatter_RawPassesThrough(t *testing.T) {
	f := formatter.New(formatter.FormatRaw, "")
	input := `{"level":"info","msg":"hello"}`
	if got := f.Format(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestFormatter_RawAddsSourcePrefix(t *testing.T) {
	f := formatter.New(formatter.FormatRaw, "app")
	input := "plain text"
	got := f.Format(input)
	if !strings.HasPrefix(got, "[app] ") {
		t.Errorf("expected source prefix, got %q", got)
	}
}

func TestFormatter_PrettyNonJSON(t *testing.T) {
	f := formatter.New(formatter.FormatPretty, "svc")
	input := "not json at all"
	got := f.Format(input)
	if !strings.Contains(got, "[svc]") {
		t.Errorf("expected source label in output, got %q", got)
	}
	if !strings.Contains(got, input) {
		t.Errorf("expected original line in output, got %q", got)
	}
}

func TestFormatter_PrettyJSON(t *testing.T) {
	f := formatter.New(formatter.FormatPretty, "")
	input := `{"time":"2024-01-15T10:30:00Z","level":"error","msg":"boom"}`
	got := f.Format(input)
	if !strings.Contains(got, "ERROR") {
		t.Errorf("expected level in output, got %q", got)
	}
	if !strings.Contains(got, "boom") {
		t.Errorf("expected message in output, got %q", got)
	}
	if !strings.Contains(got, "10:30:00") {
		t.Errorf("expected formatted time in output, got %q", got)
	}
}

func TestFormatter_PrettyJSONWithSource(t *testing.T) {
	f := formatter.New(formatter.FormatPretty, "worker")
	input := `{"level":"warn","msg":"slow query","latency":"200ms"}`
	got := f.Format(input)
	if !strings.HasPrefix(got, "[worker]") {
		t.Errorf("expected source prefix, got %q", got)
	}
	if !strings.Contains(got, "latency=200ms") {
		t.Errorf("expected extra fields, got %q", got)
	}
}

func TestFormatter_JSONAddsSource(t *testing.T) {
	f := formatter.New(formatter.FormatJSON, "api")
	input := `{"level":"info","msg":"started"}`
	got := f.Format(input)
	if !strings.Contains(got, `"_source":"api"`) {
		t.Errorf("expected _source field, got %q", got)
	}
}

func TestFormatter_JSONNoSourceUnchanged(t *testing.T) {
	f := formatter.New(formatter.FormatJSON, "")
	input := `{"level":"debug","msg":"trace"}`
	if got := f.Format(input); got != input {
		t.Errorf("expected unchanged output, got %q", got)
	}
}
