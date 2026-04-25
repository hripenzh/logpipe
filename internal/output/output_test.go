package output_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/example/logpipe/internal/output"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	// Just ensure construction without targets does not panic.
	w := output.New()
	if w == nil {
		t.Fatal("expected non-nil Writer")
	}
}

func TestWrite_SingleTarget(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	_, err := w.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestWriteLine_AppendsNewline(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	if err := w.WriteLine("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "hello\n" {
		t.Fatalf("expected %q, got %q", "hello\n", got)
	}
}

func TestWrite_MultipleTargets(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	w := output.New(&buf1, &buf2)

	w.WriteLine("fan-out")

	for i, b := range []*bytes.Buffer{&buf1, &buf2} {
		if got := b.String(); got != "fan-out\n" {
			t.Errorf("target %d: expected %q, got %q", i, "fan-out\n", got)
		}
	}
}

func TestWrite_ConcurrentSafety(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			w.WriteLine("line")
		}()
	}
	wg.Wait()

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != goroutines {
		t.Fatalf("expected %d lines, got %d", goroutines, len(lines))
	}
}
