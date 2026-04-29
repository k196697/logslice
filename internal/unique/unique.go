package unique

// Options configures the unique field value extraction.
type Options struct {
	// Fields is the list of field names to extract unique values from.
	Fields []string
	// Limit caps the number of unique values per field (0 = unlimited).
	Limit int
}

// Result holds unique values per field.
type Result map[string][]string

// Entry is a minimal representation of a log entry for this package.
type Entry struct {
	Fields map[string]string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Limit: 0,
	}
}

// Run extracts unique values for each requested field across all entries.
func Run(entries []Entry, opts Options) (Result, error) {
	if len(opts.Fields) == 0 {
		return Result{}, nil
	}

	seen := make(map[string]map[string]struct{}, len(opts.Fields))
	order := make(map[string][]string, len(opts.Fields))

	for _, f := range opts.Fields {
		seen[f] = make(map[string]struct{})
		order[f] = []string{}
	}

	for _, e := range entries {
		for _, f := range opts.Fields {
			v, ok := e.Fields[f]
			if !ok {
				continue
			}
			if _, dup := seen[f][v]; dup {
				continue
			}
			if opts.Limit > 0 && len(order[f]) >= opts.Limit {
				continue
			}
			seen[f][v] = struct{}{}
			order[f] = append(order[f], v)
		}
	}

	result := make(Result, len(opts.Fields))
	for _, f := range opts.Fields {
		result[f] = order[f]
	}
	return result, nil
}
