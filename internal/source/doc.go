// Package source provides primitives for reading log lines from named
// io.Reader sources and multiplexing them into a single stream.
//
// A Source wraps an io.Reader with a human-readable name (e.g. a filename
// or service label). Calling Tail starts a goroutine that scans the reader
// line-by-line and emits source.Line values on a channel until the reader
// is exhausted or the supplied context is cancelled.
//
// Multiple sources can be combined with Merge, which fans all source
// channels into one channel for downstream consumers such as filters or
// formatters.
//
// Example:
//
//	s1 := source.New("app.log", appFile)
//	s2 := source.New("worker.log", workerFile)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	for line := range source.Merge(ctx, s1.Tail(ctx), s2.Tail(ctx)) {
//		fmt.Printf("[%s] %s\n", line.Source, line.Text)
//	}
package source
