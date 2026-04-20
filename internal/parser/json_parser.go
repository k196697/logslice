package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// LogEntry represents a single parsed log line with its raw fields.
type LogEntry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
	Raw       string
}

// TimeFields is the ordered list of JSON keys tried when extracting a timestamp.
var TimeFields = []string{"time", "timestamp", "ts", "@timestamp"}

// ParseJSONLines reads newline-delimited JSON from r and returns a slice of
// LogEntry values. Lines that cannot be parsed are skipped with a warning.
func ParseJSONLines(r io.Reader) ([]LogEntry, error) {
	var entries []LogEntry
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		var fields map[string]interface{}
		if err := json.Unmarshal([]byte(line), &fields); err != nil {
			fmt.Printf("warning: skipping line %d (invalid JSON): %v\n", lineNum, err)
			continue
		}

		entry := LogEntry{
			Fields: fields,
			Raw:    line,
		}

		if ts, err := extractTimestamp(fields); err == nil {
			entry.Timestamp = ts
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return entries, nil
}

// extractTimestamp tries each known time field key and parses common formats.
func extractTimestamp(fields map[string]interface{}) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	for _, key := range TimeFields {
		val, ok := fields[key]
		if !ok {
			continue
		}

		switch v := val.(type) {
		case string:
			for _, f := range formats {
				if t, err := time.Parse(f, v); err == nil {
					return t, nil
				}
			}
		case float64:
			// Unix epoch seconds (possibly with sub-second precision)
			sec := int64(v)
			nsec := int64((v - float64(sec)) * 1e9)
			return time.Unix(sec, nsec).UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("no recognisable timestamp field found")
}
