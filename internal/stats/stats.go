package stats

import (
	"time"
)

// Summary holds aggregate statistics for a set of log entries.
type Summary struct {
	Total     int
	Matched   int
	Skipped   int
	Earliest  time.Time
	Latest    time.Time
	Fields    map[string]int // field name -> occurrence count
}

// Entry is a minimal interface for a log entry used by the stats package.
type Entry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
}

// Compute calculates a Summary from a slice of entries.
// total is the total number of lines processed (including skipped).
func Compute(entries []Entry, total int) Summary {
	s := Summary{
		Total:   total,
		Matched: len(entries),
		Skipped: total - len(entries),
		Fields:  make(map[string]int),
	}

	for i, e := range entries {
		if i == 0 || e.Timestamp.Before(s.Earliest) {
			s.Earliest = e.Timestamp
		}
		if i == 0 || e.Timestamp.After(s.Latest) {
			s.Latest = e.Timestamp
		}
		for k := range e.Fields {
			s.Fields[k]++
		}
	}

	return s
}
