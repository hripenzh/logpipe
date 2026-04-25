// Package pipeline provides the core processing pipeline for logpipe.
//
// A Pipeline connects a merged log source to a filter, formatter, and output
// writer. It is the central coordination point that ties together all other
// internal packages.
//
// Typical usage:
//
//	ch := source.Merge(ctx, sources...)
//	p := pipeline.New(pipeline.Config{
//		Source:    ch,
//		Filter:    filter.New(filter.Config{MinLevel: "warn"}),
//		Formatter: formatter.New(formatter.Pretty, true),
//		Output:    output.New(os.Stdout),
//	})
//	if err := p.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package pipeline
