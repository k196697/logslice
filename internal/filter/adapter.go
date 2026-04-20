// Package filter provides log entry filtering by time range and field values.
// This file bridges the parser output into filter.LogEntry types.
package filter

import (
	"time"
)

// ParsedEntry is the minimal interface expected from the parser layer.
type ParsedEntry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
	Raw       string
}

// FromParsed converts a slice of ParsedEntry into a slice of LogEntry
// suitable for use with the filter package.
func FromParsed(parsed []ParsedEntry) []LogEntry {
	out := make([]LogEntry, 0, len(parsed))
	for _, p := range parsed {
		out = append(out, LogEntry{
			Timestamp: p.Timestamp,
			Fields:    p.Fields,
			Raw:       p.Raw,
		})
	}
	return out
}

// ToParsed converts a slice of LogEntry back into ParsedEntry,
// useful for downstream formatting or output stages.
func ToParsed(entries []LogEntry) []ParsedEntry {
	out := make([]ParsedEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, ParsedEntry{
			Timestamp: e.Timestamp,
			Fields:    e.Fields,
			Raw:       e.Raw,
		})
	}
	return out
}
