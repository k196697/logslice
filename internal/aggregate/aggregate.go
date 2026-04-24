package aggregate

import (
	"fmt"
	"sort"

	"github.com/logslice/logslice/internal/filter"
)

// Options controls aggregation behaviour.
type Options struct {
	// GroupBy is the field name to group entries by.
	GroupBy string
	// CountField is the name of the synthetic count field in output.
	CountField string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		GroupBy:    "",
		CountField: "count",
	}
}

// Group holds a group key and the entries that belong to it.
type Group struct {
	Key     string
	Entries []filter.Entry
	Count   int
}

// Run groups entries by the value of opts.GroupBy field and returns one
// summary entry per group, with a synthetic count field attached.
func Run(entries []filter.Entry, opts Options) ([]filter.Entry, error) {
	if opts.GroupBy == "" {
		return nil, fmt.Errorf("aggregate: GroupBy field must not be empty")
	}

	countField := opts.CountField
	if countField == "" {
		countField = "count"
	}

	// Preserve insertion order of keys.
	order := []string{}
	groups := map[string]*Group{}

	for _, e := range entries {
		val, ok := e.Fields[opts.GroupBy]
		if !ok {
			val = ""
		}
		key := fmt.Sprintf("%v", val)
		if _, exists := groups[key]; !exists {
			groups[key] = &Group{Key: key}
			order = append(order, key)
		}
		groups[key].Entries = append(groups[key].Entries, e)
		groups[key].Count++
	}

	sort.Strings(order)

	result := make([]filter.Entry, 0, len(groups))
	for _, key := range order {
		g := groups[key]
		// Use the first entry as the representative, then attach count.
		rep := g.Entries[0]
		newFields := make(map[string]interface{}, len(rep.Fields)+1)
		for k, v := range rep.Fields {
			newFields[k] = v
		}
		newFields[countField] = g.Count
		rep.Fields = newFields
		result = append(result, rep)
	}
	return result, nil
}
