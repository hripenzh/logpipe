// Package formatter provides log line rendering for logpipe.
//
// Three output formats are supported:
//
//   - FormatRaw    — lines are emitted as-is, with an optional source prefix.
//   - FormatPretty — JSON log lines are parsed and rendered in a human-friendly
//     single-line format: timestamp level message key=value …
//     Non-JSON lines are passed through with the source prefix.
//   - FormatJSON   — JSON log lines are re-emitted with an injected "_source"
//     field so downstream consumers can identify the origin.
//
// Usage:
//
//	f := formatter.New(formatter.FormatPretty, "api-server")
//	formatted := f.Format(rawLine)
//	fmt.Println(formatted)
package formatter
