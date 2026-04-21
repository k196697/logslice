package dedup

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// Options controls deduplication behavior.
type Options struct {
	// Fields to use as the dedup key. If empty, all fields are used.
	Fields []string
	// Consecutive, when true, only deduplicates adjacent duplicate entries.
	Consecutive bool
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Fields:      nil,
		Consecutive: false,
	}
}

// Run removes duplicate log entries based on the provided options.
// It returns a new slice with duplicates removed, preserving order.
func Run(entries []filter.Entry, opts Options) []filter.Entry {
	if len(entries) == 0 {
		return entries
	}

	seen := make(map[string]struct{})
	result := make([]filter.Entry, 0, len(entries))
	lastKey := ""

	for _, e := range entries {
		key := hashEntry(e, opts.Fields)
		if opts.Consecutive {
			if key == lastKey {
				continue
			}
			lastKey = key
			result = append(result, e)
		} else {
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			result = append(result, e)
		}
	}

	return result
}

// hashEntry produces a stable hash key for an entry.
func hashEntry(e filter.Entry, fields []string) string {
	var sb strings.Builder

	if len(fields) == 0 {
		// Use all fields in sorted order for a stable key.
		keys := make([]string, 0, len(e.Fields))
		for k := range e.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(k)
			sb.WriteByte('=')
			sb.WriteString(e.Fields[k])
			sb.WriteByte(';')
		}
	} else {
		for _, f := range fields {
			sb.WriteString(f)
			sb.WriteByte('=')
			sb.WriteString(e.Fields[f])
			sb.WriteByte(';')
		}
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return fmt.Sprintf("%x", sum)
}
