// Package sink provides writers that consume formatted log lines
// and direct them to one or more destinations such as files or stdout.
package sink

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Sink reads formatted lines from a channel and writes them to a writer.
type Sink struct {
	w io.Writer
}

// New creates a Sink that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Sink {
	if w == nil {
		w = os.Stdout
	}
	return &Sink{w: w}
}

// Run reads lines from in until it is closed or ctx is cancelled.
// Each line is written to the underlying writer followed by a newline.
// Run returns the first write error encountered, or nil.
func (s *Sink) Run(ctx context.Context, in <-chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-in:
			if !ok {
				return nil
			}
			if _, err := fmt.Fprintln(s.w, line); err != nil {
				return err
			}
		}
	}
}

// FileWriter opens a file at path for appending (creating it if necessary)
// and returns a writer backed by that file, along with a close function.
func FileWriter(path string) (io.Writer, func() error, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("sink: open %q: %w", path, err)
	}
	return f, f.Close, nil
}
