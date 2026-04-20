package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"
)

// ParseCSVLines parses a CSV-formatted log stream into LogEntry slice.
// The first row is expected to be a header row defining field names.
// The timestampField specifies which column holds the timestamp.
func ParseCSVLines(r io.Reader, timestampField string) ([]LogEntry, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading CSV header: %w", err)
	}

	tsIndex := -1
	for i, col := range header {
		if strings.EqualFold(col, timestampField) {
			tsIndex = i
			break
		}
	}

	var entries []LogEntry
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// skip malformed rows
			continue
		}
		if len(row) != len(header) {
			continue
		}

		fields := make(map[string]interface{}, len(header))
		for i, col := range header {
			fields[col] = row[i]
		}

		var ts time.Time
		if tsIndex >= 0 {
			ts, _ = parseTimestampString(row[tsIndex])
		}

		entries = append(entries, LogEntry{
			Timestamp: ts,
			Fields:    fields,
			Raw:       strings.Join(row, ","),
		})
	}

	return entries, nil
}

// parseTimestampString attempts several common timestamp formats.
func parseTimestampString(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006/01/02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized timestamp format: %q", s)
}
