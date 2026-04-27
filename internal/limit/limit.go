package limit

// Options controls how many log entries are returned.
type Options struct {
	// MaxEntries is the maximum number of entries to keep.
	// A value of 0 means no limit.
	MaxEntries int

	// FromEnd, when true, keeps the last MaxEntries entries instead of the first.
	FromEnd bool
}

// DefaultOptions returns an Options with no limit applied.
func DefaultOptions() Options {
	return Options{}
}

// Entry is a minimal representation of a log entry used by this package.
type Entry struct {
	Fields map[string]string
	TimestampNano int64
}

// Run applies the limit to the provided entries.
// If opts.MaxEntries is 0, all entries are returned unchanged.
func Run(entries []Entry, opts Options) []Entry {
	if opts.MaxEntries <= 0 || len(entries) == 0 {
		return entries
	}

	if opts.MaxEntries >= len(entries) {
		return entries
	}

	if opts.FromEnd {
		return entries[len(entries)-opts.MaxEntries:]
	}

	return entries[:opts.MaxEntries]
}
