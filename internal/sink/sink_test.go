package sink_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/sink"
)

func sendLines(lines []string) <-chan string {
	ch := make(chan string, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return ch
}

func TestSink_WritesLinesToWriter(t *testing.T) {
	var buf bytes.Buffer
	s := sink.New(&buf)

	in := sendLines([]string{"line one", "line two", "line three"})
	if err := s.Run(context.Background(), in); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	for _, want := range []string{"line one", "line two", "line three"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q; got:\n%s", want, got)
		}
	}
}

func TestSink_RespectsContextCancel(t *testing.T) {
	var buf bytes.Buffer
	s := sink.New(&buf)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	in := make(chan string) // never sends
	if err := s.Run(ctx, in); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSink_DefaultsToStdout(t *testing.T) {
	// Just ensure New(nil) does not panic.
	s := sink.New(nil)
	if s == nil {
		t.Fatal("expected non-nil sink")
	}
}

func TestFileWriter_CreatesAndWritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	w, close, err := sink.FileWriter(path)
	if err != nil {
		t.Fatalf("FileWriter error: %v", err)
	}
	defer close()

	s := sink.New(w)
	in := sendLines([]string{"hello from file"})
	if err := s.Run(context.Background(), in); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	_ = close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), "hello from file") {
		t.Errorf("expected file to contain 'hello from file', got: %s", data)
	}
}

func TestFileWriter_ErrorOnBadPath(t *testing.T) {
	_, _, err := sink.FileWriter("/nonexistent/dir/test.log")
	if err == nil {
		t.Fatal("expected error for bad path, got nil")
	}
}
