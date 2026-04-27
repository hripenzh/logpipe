package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/user/logpipe/internal/tail"
)

func TestFile_EmitsNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch, err := tail.File(ctx, f.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Write lines after tail has started.
	lines := []string{"first line", "second line", "third line"}
	for _, l := range lines {
		if _, err := f.WriteString(l + "\n"); err != nil {
			t.Fatal(err)
		}
	}

	for _, want := range lines {
		select {
		case got := <-ch:
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		case <-ctx.Done():
			t.Fatalf("timed out waiting for line %q", want)
		}
	}
}

func TestFile_ClosesChannelOnCancel(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())

	ch, err := tail.File(ctx, f.Name())
	if err != nil {
		t.Fatal(err)
	}

	cancel()

	select {
	case _, open := <-ch:
		if open {
			t.Error("expected channel to be closed after cancel")
		}
	case <-time.After(2 * time.Second):
		t.Error("channel was not closed within timeout")
	}
}

func TestFile_ReturnsErrorForMissingFile(t *testing.T) {
	ctx := context.Background()
	_, err := tail.File(ctx, "/nonexistent/path/to/file.log")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
