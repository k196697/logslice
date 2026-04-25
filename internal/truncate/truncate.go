package truncate

import (
	"errors"
	"strings"
)

// DefaultOptions returns a safe default Options.
func DefaultOptions() Options {
	return Options{
		MaxLength: 256,
		Suffix:    "...",
	}
}

// Options controls how field values are truncated.
type Options struct {
	// Fields lists the field names to truncate. If empty, all string fields are truncated.
	Fields []string
	// MaxLength is the maximum number of runes allowed before truncation.
	MaxLength int
	// Suffix is appended when a value is truncated.
	Suffix string
}

// Entry is a log entry with named string fields.
type Entry struct {
	Fields map[string]string
}

// Run applies truncation to all entries according to opts.
func Run(entries []Entry, opts Options) ([]Entry, error) {
	if opts.MaxLength <= 0 {
		return nil, errors.New("truncate: MaxLength must be greater than zero")
	}

	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		result = append(result, truncateEntry(e, opts))
	}
	return result, nil
}

func truncateEntry(e Entry, opts Options) Entry {
	out := Entry{Fields: make(map[string]string, len(e.Fields))}
	for k, v := range e.Fields {
		if shouldTruncate(k, opts.Fields) {
			out.Fields[k] = truncateValue(v, opts.MaxLength, opts.Suffix)
		} else {
			out.Fields[k] = v
		}
	}
	return out
}

func shouldTruncate(field string, fields []string) bool {
	if len(fields) == 0 {
		return true
	}
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

func truncateValue(v string, maxLen int, suffix string) string {
	runes := []rune(v)
	if len(runes) <= maxLen {
		return v
	}
	cut := maxLen - len([]rune(suffix))
	if cut < 0 {
		cut = 0
	}
	return string(runes[:cut]) + suffix
}

// ParseFields splits a comma-separated list of field names.
func ParseFields(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
