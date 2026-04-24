// Package formatter provides log line formatting for logpipe output.
package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Format controls how log lines are rendered.
type Format string

const (
	FormatPretty Format = "pretty"
	FormatJSON   Format = "json"
	FormatRaw    Format = "raw"
)

// Formatter renders log lines for terminal output.
type Formatter struct {
	format    Format
	source    string
	timeField string
	levelField string
	messageField string
}

// New creates a Formatter with the given format and source label.
func New(format Format, source string) *Formatter {
	return &Formatter{
		format:       format,
		source:       source,
		timeField:    "time",
		levelField:   "level",
		messageField: "msg",
	}
}

// Format renders a single log line.
func (f *Formatter) Format(line string) string {
	switch f.format {
	case FormatPretty:
		return f.pretty(line)
	case FormatJSON:
		return f.jsonOut(line)
	default:
		if f.source != "" {
			return fmt.Sprintf("[%s] %s", f.source, line)
		}
		return line
	}
}

func (f *Formatter) pretty(line string) string {
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		// Non-JSON: pass through with source prefix
		if f.source != "" {
			return fmt.Sprintf("[%s] %s", f.source, line)
		}
		return line
	}

	level := strings.ToUpper(fmt.Sprintf("%v", fields[f.levelField]))
	msg := fmt.Sprintf("%v", fields[f.messageField])
	ts := formatTime(fields[f.timeField])

	var extra []string
	for k, v := range fields {
		if k == f.timeField || k == f.levelField || k == f.messageField {
			continue
		}
		extra = append(extra, fmt.Sprintf("%s=%v", k, v))
	}

	base := fmt.Sprintf("%s %-5s %s", ts, level, msg)
	if len(extra) > 0 {
		base += "  " + strings.Join(extra, " ")
	}
	if f.source != "" {
		return fmt.Sprintf("[%s] %s", f.source, base)
	}
	return base
}

func (f *Formatter) jsonOut(line string) string {
	if f.source == "" {
		return line
	}
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return line
	}
	fields["_source"] = f.source
	b, err := json.Marshal(fields)
	if err != nil {
		return line
	}
	return string(b)
}

func formatTime(v interface{}) string {
	if v == nil {
		return time.Now().Format("15:04:05")
	}
	s := fmt.Sprintf("%v", v)
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Format("15:04:05")
		}
	}
	return s
}
