package flatten

import (
	"fmt"
	"strings"
)

// Options controls how nested fields are flattened.
type Options struct {
	// Separator is placed between key segments (default: ".").
	Separator string
	// MaxDepth limits recursion; 0 means unlimited.
	MaxDepth int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator: ".",
		MaxDepth:  0,
	}
}

// Entry is a log entry with string-keyed fields.
type Entry struct {
	Timestamp int64
	Fields    map[string]interface{}
}

// Run flattens nested map fields in each entry according to opts.
func Run(entries []Entry, opts Options) []Entry {
	if opts.Separator == "" {
		opts.Separator = "."
	}
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		flat := make(map[string]interface{})
		flattenMap("", e.Fields, flat, opts.Separator, opts.MaxDepth, 0)
		out = append(out, Entry{Timestamp: e.Timestamp, Fields: flat})
	}
	return out
}

func flattenMap(prefix string, src map[string]interface{}, dst map[string]interface{}, sep string, maxDepth, depth int) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		switch child := v.(type) {
		case map[string]interface{}:
			if maxDepth > 0 && depth+1 >= maxDepth {
				dst[key] = formatValue(child)
			} else {
				flattenMap(key, child, dst, sep, maxDepth, depth+1)
			}
		default:
			dst[key] = v
		}
	}
}

func formatValue(v interface{}) string {
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}
