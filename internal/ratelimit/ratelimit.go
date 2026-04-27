// Package ratelimit provides a line-rate limiter for log pipelines.
// It allows capping the number of log lines emitted per second to prevent
// downstream consumers from being overwhelmed by high-volume sources.
package ratelimit

import (
	"context"
	"time"
)

// Limiter wraps an input channel and throttles lines to a maximum rate.
type Limiter struct {
	rate     int           // max lines per second
	interval time.Duration // derived from rate
}

// New creates a Limiter that allows at most linesPerSecond lines through
// per second. A rate of zero or negative disables limiting.
func New(linesPerSecond int) *Limiter {
	var interval time.Duration
	if linesPerSecond > 0 {
		interval = time.Second / time.Duration(linesPerSecond)
	}
	return &Limiter{
		rate:     linesPerSecond,
		interval: interval,
	}
}

// Apply reads lines from in and forwards them to the returned channel,
// inserting a delay between lines to honour the configured rate. If the
// rate is zero the returned channel is simply the input channel.
func (l *Limiter) Apply(ctx context.Context, in <-chan string) <-chan string {
	if l.rate <= 0 {
		return in
	}

	out := make(chan string)
	go func() {
		defer close(out)
		ticker := time.NewTicker(l.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-in:
				if !ok {
					return
				}
				// Wait for the next tick before forwarding.
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
				select {
				case out <- line:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
