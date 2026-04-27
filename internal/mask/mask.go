package mask

import (
	"fmt"
	"strings"
)

// DefaultOptions returns a default Options value.
func DefaultOptions() Options {
	return Options{
		Char:   "*",
		Fields: []string{},
	}
}

// Options configures the masking behaviour.
type Options struct {
	// Fields to partially mask.
	Fields []string
	// Char is the replacement character (default "*").
	Char string
	// KeepPrefix is the number of leading characters to preserve.
	KeepPrefix int
	// KeepSuffix is the number of trailing characters to preserve.
	KeepSuffix int
}

// Entry is a single log record understood by this package.
type Entry struct {
	Fields map[string]string
}

// Run applies partial masking to the specified fields of each entry.
func Run(entries []Entry, opts Options) ([]Entry, error) {
	if len(opts.Fields) == 0 {
		return entries, nil
	}
	if opts.Char == "" {
		opts.Char = "*"
	}

	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		masked := Entry{Fields: make(map[string]string, len(e.Fields))}
		for k, v := range e.Fields {
			masked.Fields[k] = v
		}
		for _, field := range opts.Fields {
			val, ok := masked.Fields[field]
			if !ok {
				continue
			}
			masked.Fields[field] = maskValue(val, opts)
		}
		out = append(out, masked)
	}
	return out, nil
}

func maskValue(s string, opts Options) string {
	n := len(s)
	pre := opts.KeepPrefix
	suf := opts.KeepSuffix
	if pre < 0 {
		pre = 0
	}
	if suf < 0 {
		suf = 0
	}
	if pre+suf >= n {
		// Nothing to mask.
		return s
	}
	midLen := n - pre - suf
	return fmt.Sprintf("%s%s%s",
		s[:pre],
		strings.Repeat(opts.Char, midLen),
		s[n-suf:],
	)
}
