package source

import (
	"bufio"
	"context"
	"io"
)

// Source represents a named log source that emits lines.
type Source struct {
	Name   string
	reader io.Reader
}

// Line represents a single log line from a named source.
type Line struct {
	Source string
	Text   string
}

// New creates a new Source with the given name and reader.
func New(name string, r io.Reader) *Source {
	return &Source{Name: name, reader: r}
}

// Tail reads lines from the source and sends them to the returned channel.
// The channel is closed when the reader is exhausted or the context is cancelled.
func (s *Source) Tail(ctx context.Context) <-chan Line {
	ch := make(chan Line)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(s.reader)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- Line{Source: s.Name, Text: scanner.Text()}:
			}
		}
	}()
	return ch
}

// Merge fans in multiple source channels into a single Line channel.
func Merge(ctx context.Context, sources ...<-chan Line) <-chan Line {
	out := make(chan Line)
	for _, src := range sources {
		go func(c <-chan Line) {
			for {
				select {
				case <-ctx.Done():
					return
				case line, ok := <-c:
					if !ok {
						return
					}
					select {
					case out <- line:
					case <-ctx.Done():
						return
					}
				}
			}
		}(src)
	}
	return out
}
