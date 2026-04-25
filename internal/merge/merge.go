package merge

import (
	"sort"
	"time"
)

// Entry represents a log entry for merging purposes.
type Entry struct {
	Timestamp time.Time
	Fields    map[string]string
	Raw       string
}

// Options controls merge behaviour.
type Options struct {
	// SortByTime sorts all merged entries by timestamp ascending.
	SortByTime bool
	// DeduplicateConsecutive removes back-to-back identical Raw lines.
	DeduplicateConsecutive bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		SortByTime:             true,
		DeduplicateConsecutive: false,
	}
}

// Run merges multiple slices of entries into one.
// When opts.SortByTime is true the result is sorted by Timestamp ascending.
// Entries with a zero Timestamp are placed at the end.
func Run(sources [][]Entry, opts Options) []Entry {
	total := 0
	for _, s := range sources {
		total += len(s)
	}

	merged := make([]Entry, 0, total)
	for _, s := range sources {
		merged = append(merged, s...)
	}

	if opts.SortByTime {
		sort.SliceStable(merged, func(i, j int) bool {
			zi := merged[i].Timestamp.IsZero()
			zj := merged[j].Timestamp.IsZero()
			if zi && zj {
				return false
			}
			if zi {
				return false
			}
			if zj {
				return true
			}
			return merged[i].Timestamp.Before(merged[j].Timestamp)
		})
	}

	if opts.DeduplicateConsecutive {
		merged = dedupConsecutive(merged)
	}

	return merged
}

func dedupConsecutive(entries []Entry) []Entry {
	if len(entries) == 0 {
		return entries
	}
	out := []Entry{entries[0]}
	for i := 1; i < len(entries); i++ {
		if entries[i].Raw != out[len(out)-1].Raw {
			out = append(out, entries[i])
		}
	}
	return out
}
