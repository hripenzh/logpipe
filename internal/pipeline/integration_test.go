package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/formatter"
	"github.com/user/logpipe/internal/output"
	"github.com/user/logpipe/internal/pipeline"
	"github.com/user/logpipe/internal/source"
)

// TestPipeline_FullIntegration exercises source → filter → formatter → output
// with multiple sources and a key-contains filter to confirm end-to-end wiring.
func TestPipeline_FullIntegration(t *testing.T) {
	var buf bytes.Buffer
	out := output.New(&buf)
	fmt := formatter.New(formatter.Raw, true) // prefix source name
	filt := filter.New(filter.Config{
		KeyContains: map[string]string{"env": "prod"},
	})

	lines := []source.Line{
		{Source: "api", Text: `{"level":"info","env":"prod","msg":"request ok"}`},
		{Source: "worker", Text: `{"level":"info","env":"staging","msg":"job done"}`},
		{Source: "api", Text: `{"level":"error","env":"prod","msg":"timeout"}`},
	}

	ch := makeSource(lines...)
	p := pipeline.New(pipeline.Config{
		Source: ch, Filter: filt, Formatter: fmt, Output: out,
	})

	if err := p.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()

	if !strings.Contains(got, "request ok") {
		t.Errorf("expected prod api line in output")
	}
	if strings.Contains(got, "job done") {
		t.Errorf("expected staging worker line to be filtered out")
	}
	if !strings.Contains(got, "timeout") {
		t.Errorf("expected prod error line in output")
	}
	// source prefix should appear
	if !strings.Contains(got, "api") {
		t.Errorf("expected source prefix 'api' in output")
	}
}
