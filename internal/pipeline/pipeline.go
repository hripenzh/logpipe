// Package pipeline wires together sources, filters, formatters, and output
// into a single processing pipeline for logpipe.
package pipeline

import (
	"context"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/formatter"
	"github.com/user/logpipe/internal/output"
	"github.com/user/logpipe/internal/source"
)

// Pipeline reads log lines from a merged source, applies a filter,
// formats each line, and writes the result to the configured output.
type Pipeline struct {
	src    <-chan source.Line
	filt   *filter.Filter
	fmt    *formatter.Formatter
	out    *output.Output
}

// Config holds the dependencies needed to construct a Pipeline.
type Config struct {
	Source    <-chan source.Line
	Filter    *filter.Filter
	Formatter *formatter.Formatter
	Output    *output.Output
}

// New creates a Pipeline from the provided Config.
func New(cfg Config) *Pipeline {
	return &Pipeline{
		src:  cfg.Source,
		filt: cfg.Filter,
		fmt:  cfg.Formatter,
		out:  cfg.Output,
	}
}

// Run processes log lines until the context is cancelled or the source
// channel is closed. It returns any write error encountered.
func (p *Pipeline) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-p.src:
			if !ok {
				return nil
			}
			if !p.filt.Apply(line.Text) {
				continue
			}
			formatted := p.fmt.Format(line.Source, line.Text)
			if err := p.out.WriteLine(formatted); err != nil {
				return err
			}
		}
	}
}
