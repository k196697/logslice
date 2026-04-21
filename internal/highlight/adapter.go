package highlight

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ApplyToEntry formats a single filter.Entry as a colorized text line.
// Fields are printed as key=value pairs; the timestamp field is highlighted
// with Bold if the highlighter is enabled.
func ApplyToEntry(h *Highlighter, e filter.Entry, timeField string) string {
	var parts []string

	// Timestamp first
	if !e.Timestamp.IsZero() {
		ts := e.Timestamp.Format("2006-01-02T15:04:05Z07:00")
		if h != nil {
			ts = fmt.Sprintf("%s%s%s", Bold, ts, Reset)
		}
		parts = append(parts, fmt.Sprintf("%s=%s", timeField, ts))
	}

	for k, v := range e.Fields {
		val := fmt.Sprintf("%v", v)
		if h != nil {
			val = h.ColorizeField(k, val)
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, val))
	}

	return strings.Join(parts, " ")
}
