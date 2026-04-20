// Package parser provides readers for structured log formats.
// It defines the shared LogEntry type used across JSON and CSV parsers.
package parser

import "time"

// LogEntry represents a single parsed log line with an optional timestamp
// and a map of all extracted fields.
type LogEntry struct {
	// Timestamp is the parsed time value from the designated timestamp field.
	// It is zero if no timestamp could be determined.
	Timestamp time.Time

	// Fields holds all key-value pairs extracted from the log line.
	Fields map[string]interface{}

	// Raw is the original unparsed line or row as read from the source.
	Raw string
}

// Format identifies the log file format to use when parsing.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// SupportedFormats returns the list of format identifiers recognised by logslice.
func SupportedFormats() []Format {
	return []Format{FormatJSON, FormatCSV}
}
