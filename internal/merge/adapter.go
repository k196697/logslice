package merge

import (
	"time"

	"github.com/user/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slices to merge Entry slices.
func FromFilterEntries(sources [][]filter.Entry) [][]Entry {
	out := make([][]Entry, len(sources))
	for i, src := range sources {
		out[i] = make([]Entry, len(src))
		for j, fe := range src {
			out[i][j] = fromFilterEntry(fe)
		}
	}
	return out
}

// ToFilterEntries converts merged entries back to filter.Entry slice.
func ToFilterEntries(entries []Entry) []filter.Entry {
	out := make([]filter.Entry, len(entries))
	for i, e := range entries {
		out[i] = toFilterEntry(e)
	}
	return out
}

func fromFilterEntry(fe filter.Entry) Entry {
	fields := make(map[string]string, len(fe.Fields))
	for k, v := range fe.Fields {
		fields[k] = v
	}
	var ts time.Time
	if fe.Timestamp != nil {
		ts = *fe.Timestamp
	}
	return Entry{
		Timestamp: ts,
		Fields:    fields,
		Raw:       fe.Raw,
	}
}

func toFilterEntry(e Entry) filter.Entry {
	fields := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	var ts *time.Time
	if !e.Timestamp.IsZero() {
		t := e.Timestamp
		ts = &t
	}
	return filter.Entry{
		Timestamp: ts,
		Fields:    fields,
		Raw:       e.Raw,
	}
}
