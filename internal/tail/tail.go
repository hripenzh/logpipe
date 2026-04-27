// Package tail provides utilities for watching a file and emitting new lines
// as they are appended, similar to `tail -f`.
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// PollInterval is the duration between file size checks when watching.
const PollInterval = 100 * time.Millisecond

// File tails a file, sending each new line to the returned channel.
// The channel is closed when ctx is cancelled or an unrecoverable error occurs.
func File(ctx context.Context, path string) (<-chan string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Seek to end so we only emit lines written after this call.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}

	ch := make(chan string, 64)

	go func() {
		defer close(ch)
		defer f.Close()

		reader := bufio.NewReader(f)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadString('\n')
			if len(line) > 0 {
				// Strip trailing newline before sending.
				if len(line) > 0 && line[len(line)-1] == '\n' {
					line = line[:len(line)-1]
				}
				select {
				case ch <- line:
				case <-ctx.Done():
					return
				}
			}

			if err != nil {
				if err == io.EOF {
					// No new data yet; wait before retrying.
					select {
					case <-time.After(PollInterval):
					case <-ctx.Done():
						return
					}
					continue
				}
				// Unrecoverable read error.
				return
			}
		}
	}()

	return ch, nil
}
