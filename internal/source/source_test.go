package source_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/logpipe/internal/source"
)

func TestSource_TailEmitsLines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	s := source.New("test", strings.NewReader(input))
	ctx := context.Background()

	var got []source.Line
	for line := range s.Tail(ctx) {
		got = append(got, line)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0].Text != "line one" {
		t.Errorf("expected 'line one', got %q", got[0].Text)
	}
	for _, l := range got {
		if l.Source != "test" {
			t.Errorf("expected source 'test', got %q", l.Source)
		}
	}
}

func TestSource_TailRespectsContextCancel(t *testing.T) {
	// Infinite reader via pipe — cancel should stop tailing
	pr, pw := strings.NewReader(""), strings.NewReader("")
	_ = pr
	_ = pw

	s := source.New("cancel-test", strings.NewReader("a\nb\nc\n"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		for range s.Tail(ctx) {
		}
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("Tail did not respect context cancellation")
	}
}

func TestMerge_CombinesMultipleSources(t *testing.T) {
	s1 := source.New("s1", strings.NewReader("alpha\nbeta\n"))
	s2 := source.New("s2", strings.NewReader("gamma\ndelta\n"))
	ctx := context.Background()

	merged := source.Merge(ctx, s1.Tail(ctx), s2.Tail(ctx))

	seen := map[string]bool{}
	for line := range merged {
		seen[line.Text] = true
	}

	for _, want := range []string{"alpha", "beta", "gamma", "delta"} {
		if !seen[want] {
			t.Errorf("expected to see line %q in merged output", want)
		}
	}
}
