// Package output handles writing formatted log lines to one or more destinations.
package output

import (
	"io"
	"os"
	"sync"
)

// Writer wraps one or more io.Writer destinations and provides thread-safe
// concurrent writes.
type Writer struct {
	mu      sync.Mutex
	targets []io.Writer
}

// New creates a Writer that fans out to the provided targets. If no targets are
// supplied, os.Stdout is used as the default destination.
func New(targets ...io.Writer) *Writer {
	if len(targets) == 0 {
		targets = []io.Writer{os.Stdout}
	}
	return &Writer{targets: targets}
}

// Write writes p to every target in order. The call is serialised so that
// concurrent goroutines (e.g. one per log source) do not interleave lines.
// It returns the number of bytes written to the first target and the first
// error encountered, if any.
func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var (
		n   int
		err error
	)
	for i, t := range w.targets {
		nn, werr := t.Write(p)
		if i == 0 {
			n = nn
		}
		if werr != nil && err == nil {
			err = werr
		}
	}
	return n, err
}

// WriteLine appends a newline to line and writes it via Write.
func (w *Writer) WriteLine(line string) error {
	_, err := w.Write([]byte(line + "\n"))
	return err
}
