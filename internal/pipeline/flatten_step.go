package pipeline

import (
	"logslice/internal/filter"
	"logslice/internal/flatten"
)

// ParseFlattenConfig returns a flatten.Options pointer if flattening is
// requested, or nil if the feature is disabled.
func ParseFlattenConfig(separator string, maxDepth int, enabled bool) *flatten.Options {
	if !enabled {
		return nil
	}
	opts := flatten.DefaultOptions()
	if separator != "" {
		opts.Separator = separator
	}
	if maxDepth > 0 {
		opts.MaxDepth = maxDepth
	}
	return &opts
}

// applyFlatten runs the flatten step when opts is non-nil and returns the
// (possibly unchanged) entry slice together with any error.
func applyFlatten(entries []filter.Entry, opts *flatten.Options) ([]filter.Entry, error) {
	if opts == nil {
		return entries, nil
	}
	result, err := flatten.RunFromConfig(entries, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
