package stats

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Print writes a human-readable summary to w.
func Print(w io.Writer, s Summary) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Total lines:\t%d\n", s.Total)
	fmt.Fprintf(tw, "Matched:\t%d\n", s.Matched)
	fmt.Fprintf(tw, "Skipped:\t%d\n", s.Skipped)

	if s.Matched > 0 {
		fmt.Fprintf(tw, "Earliest:\t%s\n", s.Earliest.UTC().Format("2006-01-02T15:04:05Z"))
		fmt.Fprintf(tw, "Latest:\t%s\n", s.Latest.UTC().Format("2006-01-02T15:04:05Z"))
	}

	if len(s.Fields) > 0 {
		fmt.Fprintf(tw, "Fields:\t\n")
		keys := make([]string, 0, len(s.Fields))
		for k := range s.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(tw, "  %s:\t%d\n", k, s.Fields[k])
		}
	}

	tw.Flush()
}
