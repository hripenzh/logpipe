package filter

import (
	"encoding/json"
	"strings"
)

// Level represents a log severity level.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// levelOrder maps levels to numeric priority for comparison.
var levelOrder = map[Level]int{
	LevelDebug: 0,
	LevelInfo:  1,
	LevelWarn:  2,
	LevelError: 3,
}

// Options holds the filtering criteria applied to each log line.
type Options struct {
	// MinLevel filters out log entries below this severity.
	MinLevel Level
	// KeyContains filters entries where the given key's value contains the substring.
	KeyContains map[string]string
}

// Filter evaluates a raw JSON log line against the provided Options.
// It returns true if the line passes all active filters.
func Filter(line []byte, opts Options) bool {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		// Non-JSON lines pass through unfiltered.
		return true
	}

	if opts.MinLevel != "" {
		if !passesLevelFilter(entry, opts.MinLevel) {
			return false
		}
	}

	for key, substr := range opts.KeyContains {
		val, ok := entry[key]
		if !ok {
			return false
		}
		strVal, ok := val.(string)
		if !ok {
			return false
		}
		if !strings.Contains(strVal, substr) {
			return false
		}
	}

	return true
}

func passesLevelFilter(entry map[string]interface{}, minLevel Level) bool {
	raw, ok := entry["level"]
	if !ok {
		return true
	}
	strLevel, ok := raw.(string)
	if !ok {
		return true
	}
	entryLevel := Level(strings.ToLower(strLevel))
	minPriority, knownMin := levelOrder[minLevel]
	entryPriority, knownEntry := levelOrder[entryLevel]
	if !knownMin || !knownEntry {
		return true
	}
	return entryPriority >= minPriority
}
