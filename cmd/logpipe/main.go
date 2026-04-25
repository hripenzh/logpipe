package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/logpipe/internal/config"
	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/formatter"
	"github.com/yourorg/logpipe/internal/source"
)

func main() {
	configPath := flag.String("config", "logpipe.yaml", "path to config file")
	minLevel := flag.String("level", "", "minimum log level (debug, info, warn, error)")
	outputFmt := flag.String("format", "", "output format: raw or pretty (overrides config)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logpipe: failed to load config: %v\n", err)
		os.Exit(1)
	}

	if *minLevel != "" {
		cfg.Filter.MinLevel = *minLevel
	}
	if *outputFmt != "" {
		cfg.Output.Format = *outputFmt
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	sources := make([]*source.Source, 0, len(cfg.Sources))
	for _, sc := range cfg.Sources {
		s, err := source.New(sc.Name, sc.Command)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logpipe: failed to create source %q: %v\n", sc.Name, err)
			os.Exit(1)
		}
		sources = append(sources, s)
	}

	merged := source.Merge(sources...)

	f := filter.New(cfg.Filter.MinLevel, cfg.Filter.KeyContains)
	fmt := formatter.New(cfg.Output.Format)

	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-merged:
			if !ok {
				return
			}
			if !f.Pass(entry.Line) {
				continue
			}
			formatted := fmt.Format(entry.Source, entry.Line)
			fmt.Fprintln(os.Stdout, formatted)
		}
	}
}
