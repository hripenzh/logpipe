package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/formatter"
	"github.com/user/logpipe/internal/output"
	"github.com/user/logpipe/internal/pipeline"
	"github.com/user/logpipe/internal/source"
)

func makeSource(lines ...source.Line) <-chan source.Line {
	ch := make(chan source.Line, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return ch
}

func TestPipeline_PassesMatchingLines(t *testing.T) {
	var buf bytes.Buffer
	out := output.New(&buf)
	fmt := formatter.New(formatter.Raw, false)
	filt := filter.New(filter.Config{})

	ch := makeSource(
		source.Line{Source: "app", Text: `{"level":"info","msg":"hello"}`},
	)

	p := pipeline.New(pipeline.Config{
		Source: ch, Filter: filt, Formatter: fmt, Output: out,
	})

	ctx := context.Background()
	if err := p.Run(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected output to contain 'hello', got: %s", buf.String())
	}
}

func TestPipeline_FiltersLowLevelLines(t *testing.T) {
	var buf bytes.Buffer
	out := output.New(&buf)
	fmt := formatter.New(formatter.Raw, false)
	filt := filter.New(filter.Config{MinLevel: "warn"})

	ch := makeSource(
		source.Line{Source: "app", Text: `{"level":"debug","msg":"verbose"}`},
		source.Line{Source: "app", Text: `{"level":"warn","msg":"important"}`},
	)

	p := pipeline.New(pipeline.Config{
		Source: ch, Filter: filt, Formatter: fmt, Output: out,
	})

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if strings.Contains(got, "verbose") {
		t.Errorf("expected debug line to be filtered out")
	}
	if !strings.Contains(got, "important") {
		t.Errorf("expected warn line to pass through")
	}
}

func TestPipeline_RespectsContextCancel(t *testing.T) {
	var buf bytes.Buffer
	out := output.New(&buf)
	fmt := formatter.New(formatter.Raw, false)
	filt := filter.New(filter.Config{})

	// unbuffered channel that never sends — pipeline should block until cancel
	ch := make(chan source.Line)

	p := pipeline.New(pipeline.Config{
		Source: ch, Filter: filt, Formatter: fmt, Output: out,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := p.Run(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}
