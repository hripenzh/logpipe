// Package ratelimit provides rate-limiting for log line streams.
//
// # Overview
//
// When tailing high-volume log sources it can be useful to cap the number of
// lines forwarded downstream per second.  The [Limiter] type wraps any
// string channel and inserts the necessary delays to honour a configured
// lines-per-second ceiling.
//
// # Usage
//
//	lim := ratelimit.New(100) // at most 100 lines/sec
//	throttled := lim.Apply(ctx, rawLines)
//
// A rate of 0 (or negative) disables throttling entirely and returns the
// original channel unchanged, making it safe to always construct a Limiter
// from user-supplied configuration without special-casing the "unlimited"
// scenario.
package ratelimit
