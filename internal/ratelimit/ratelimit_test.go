package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/ratelimit"
)

func feedLines(lines []string) <-chan string {
	ch := make(chan string, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return ch
}

func TestLimiter_ZeroRatePassesThrough(t *testing.T) {
	lim := ratelimit.New(0)
	in := feedLines([]string{"a", "b", "c"})
	out := lim.Apply(context.Background(), in)

	var got []string
	for line := range out {
		got = append(got, line)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestLimiter_ThrottlesOutput(t *testing.T) {
	// 10 lines/sec → ~100 ms between lines; 3 lines should take ≥200 ms.
	lim := ratelimit.New(10)
	in := feedLines([]string{"x", "y", "z"})

	start := time.Now()
	out := lim.Apply(context.Background(), in)
	var got []string
	for line := range out {
		got = append(got, line)
	}
	elapsed := time.Since(start)

	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if elapsed < 150*time.Millisecond {
		t.Errorf("expected throttling, but elapsed was only %s", elapsed)
	}
}

func TestLimiter_RespectsContextCancel(t *testing.T) {
	lim := ratelimit.New(1) // very slow: 1 line/sec

	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		ch <- "line"
	}
	close(ch)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	out := lim.Apply(ctx, ch)
	var count int
	for range out {
		count++
	}
	// At 1 line/sec with a 300 ms window we expect at most 1 line through.
	if count > 2 {
		t.Errorf("expected ≤2 lines before cancel, got %d", count)
	}
}

func TestLimiter_PreservesLineContent(t *testing.T) {
	lim := ratelimit.New(100)
	lines := []string{"hello", "world", "foo"}
	out := lim.Apply(context.Background(), feedLines(lines))

	var got []string
	for l := range out {
		got = append(got, l)
	}
	for i, want := range lines {
		if got[i] != want {
			t.Errorf("line %d: want %q, got %q", i, want, got[i])
		}
	}
}
