// Package filter provides log-level and key-based filtering for structured
// log lines emitted by logpipe sources.
//
// # Overview
//
// A Filter wraps an upstream channel of raw log lines and re-emits only those
// lines that satisfy all configured predicates.  Two kinds of predicates are
// supported:
//
//   - MinLevel – drop any JSON log line whose "level" field is below the
//     requested severity.  Lines that are not valid JSON, or that carry no
//     "level" field, are always passed through unchanged so that non-structured
//     output is never silently swallowed.
//
//   - KeyContains – keep only lines where a named JSON key contains a given
//     substring (case-sensitive).  As with MinLevel, non-JSON lines bypass
//     this check and are forwarded as-is.
//
// # Severity order
//
// The recognised level strings and their relative order are:
//
//	trace < debug < info < warn < error < fatal
//
// Comparisons are case-insensitive, so "INFO", "Info", and "info" are all
// treated identically.
//
// # Usage
//
//	lines := make(chan string)
//	f := filter.New(lines, filter.Options{
//	    MinLevel:    "warn",
//	    KeyContains: map[string]string{"service": "auth"},
//	})
//	for line := range f.Lines() {
//	    fmt.Println(line)
//	}
package filter
