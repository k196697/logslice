package filter

import (
	"time"
)

// LogEntry represents a single parsed log entry.
type LogEntry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
	Raw       string
}

// Options holds the filtering criteria for log entries.
type Options struct {
	From      *time.Time
	To        *time.Time
	FieldKey  string
	FieldVal  string
}

// Apply returns true if the given LogEntry matches all criteria in Options.
func Apply(entry LogEntry, opts Options) bool {
	if opts.From != nil && entry.Timestamp.Before(*opts.From) {
		return false
	}
	if opts.To != nil && entry.Timestamp.After(*opts.To) {
		return false
	}
	if opts.FieldKey != "" {
		val, ok := entry.Fields[opts.FieldKey]
		if !ok {
			return false
		}
		if opts.FieldVal != "" {
			strVal, ok := val.(string)
			if !ok {
				return false
			}
			if strVal != opts.FieldVal {
				return false
			}
		}
	}
	return true
}

// Run filters a slice of LogEntry values using the given Options.
func Run(entries []LogEntry, opts Options) []LogEntry {
	result := make([]LogEntry, 0, len(entries))
	for _, e := range entries {
		if Apply(e, opts) {
			result = append(result, e)
		}
	}
	return result
}
