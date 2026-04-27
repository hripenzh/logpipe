package pipeline_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/formatter"
	"github.com/yourorg/logpipe/internal/output"
	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/source"
)

// TestPipeline_FullIntegration wires together a real source, filter, formatter,
// and output to verify the full pipeline processes and delivers log lines.
func TestPipeline_FullIntegration(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"server started"}`,
		`{"level":"debug","msg":"verbose detail"}`,
		`{"level":"error","msg":"something failed"}`,
		`not json at all`,
	}

	src := source.New("app", strings.NewReader(strings.Join(lines, "\n")))

	f, err := filter.New(filter.Config{
		MinLevel: "info",
	})
	if err != nil {
		t.Fatalf("filter.New: %v", err)
	}

	fmt, err := formatter.New(formatter.Config{
		Format: "raw",
	})
	if err != nil {
		t.Fatalf("formatter.New: %v", err)
	}

	var mu sync.Mutex
	var received []string

	out := output.New(output.WriterFunc(func(line string) error {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, line)
		return nil
	}))

	p := pipeline.New(src, f, fmt, out)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.Run(ctx); err != nil && err != context.DeadlineExceeded {
		t.Fatalf("pipeline.Run: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	// debug line should be filtered out; info, error, and non-JSON should pass
	if len(received) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(received), received)
	}

	for _, want := range []string{"server started", "something failed", "not json at all"} {
		found := false
		for _, line := range received {
			if strings.Contains(line, want) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected output to contain %q, got: %v", want, received)
		}
	}
}
