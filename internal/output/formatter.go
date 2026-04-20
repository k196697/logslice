package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Format controls the output format for log entries.
type Format string

const (
	FormatJSON    Format = "json"
	FormatText    Format = "text"
	FormatCompact Format = "compact"
)

// Formatter writes log entries to a writer in the configured format.
type Formatter struct {
	Format    Format
	Writer    io.Writer
	TimeField string
}

// NewFormatter creates a Formatter with sensible defaults.
func NewFormatter(w io.Writer, format Format, timeField string) *Formatter {
	if timeField == "" {
		timeField = "timestamp"
	}
	return &Formatter{Writer: w, Format: format, TimeField: timeField}
}

// Write outputs a single log entry in the configured format.
func (f *Formatter) Write(entry parser.LogEntry) error {
	switch f.Format {
	case FormatText:
		return f.writeText(entry)
	case FormatCompact:
		return f.writeCompact(entry)
	default:
		return f.writeJSON(entry)
	}
}

func (f *Formatter) writeJSON(entry parser.LogEntry) error {
	b, err := json.Marshal(entry.Fields)
	if err != nil {
		return fmt.Errorf("output: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(f.Writer, "%s\n", b)
	return err
}

func (f *Formatter) writeText(entry parser.LogEntry) error {
	ts := entry.Timestamp.Format(time.RFC3339)
	parts := []string{ts}
	for k, v := range entry.Fields {
		if k == f.TimeField {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	_, err := fmt.Fprintf(f.Writer, "%s\n", strings.Join(parts, " "))
	return err
}

func (f *Formatter) writeCompact(entry parser.LogEntry) error {
	level, _ := entry.Fields["level"].(string)
	msg, _ := entry.Fields["msg"].(string)
	if msg == "" {
		msg, _ = entry.Fields["message"].(string)
	}
	ts := entry.Timestamp.Format(time.RFC3339)
	_, err := fmt.Fprintf(f.Writer, "[%s] %-5s %s\n", ts, strings.ToUpper(level), msg)
	return err
}
