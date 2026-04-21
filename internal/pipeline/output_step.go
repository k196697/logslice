package pipeline

import (
	"io"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/highlight"
	"github.com/yourorg/logslice/internal/output"
)

// outputOptions holds configuration for the output step of the pipeline.
type outputOptions struct {
	// Formatter controls the output format (JSON, text, compact, etc.).
	Formatter *output.Formatter

	// Highlighter applies color rules to matching fields or keywords.
	// May be nil if highlighting is disabled.
	Highlighter *highlight.Highlighter

	// TimeField is the name of the field used to represent the timestamp
	// in formatted output.
	TimeField string
}

// writeEntries writes a slice of filter entries to the given writer using
// the configured formatter and optional highlighter.
//
// Each entry is formatted according to the output format. If a highlighter
// is provided, colorization rules are applied to the formatted line before
// it is written.
//
// Returns the number of entries written and any write error encountered.
func writeEntries(w io.Writer, entries []filter.Entry, opts outputOptions) (int, error) {
	if opts.Formatter == nil {
		return 0, nil
	}

	written := 0
	for _, entry := range entries {
		line, err := opts.Formatter.Format(entry)
		if err != nil {
			// Skip entries that cannot be formatted rather than aborting
			// the entire output stream.
			continue
		}

		if opts.Highlighter != nil {
			line = applyHighlighting(opts.Highlighter, entry, line)
		}

		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return written, err
		}
		written++
	}

	return written, nil
}

// applyHighlighting applies color rules from the highlighter to a formatted
// log line. Field-level colorization is attempted first using the entry's
// raw fields; keyword-level colorization is applied to the full line string.
func applyHighlighting(h *highlight.Highlighter, entry filter.Entry, line string) string {
	// Apply field-value based colorization using structured entry data.
	for field, value := range entry.Fields {
		strVal, ok := value.(string)
		if !ok {
			continue
		}
		colored := highlight.ApplyToEntry(h, field, strVal)
		if colored != strVal {
			// Replace the plain value with the colorized version in the line.
			// This is a best-effort substitution for text-based formats.
			line = replaceFirst(line, strVal, colored)
		}
	}

	// Apply keyword-level colorization to the full line.
	line = h.ColorizeLine(line)

	return line
}

// replaceFirst replaces the first occurrence of old with new in s.
// Returns s unchanged if old is not found.
func replaceFirst(s, old, newVal string) string {
	if old == "" || old == newVal {
		return s
	}
	idx := indexString(s, old)
	if idx < 0 {
		return s
	}
	return s[:idx] + newVal + s[idx+len(old):]
}

// indexString returns the index of the first instance of substr in s,
// or -1 if substr is not present.
func indexString(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
