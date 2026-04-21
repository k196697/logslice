package transform

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// Op represents a transformation operation type.
type Op string

const (
	OpRename Op = "rename"
	OpDrop   Op = "drop"
	OpSet    Op = "set"
)

// Rule defines a single field transformation.
type Rule struct {
	Op    Op
	Field string
	Value string // used for rename (new name) and set (new value)
}

// ParseRule parses a transformation expression of the form:
//
//	rename:old=new
//	drop:fieldname
//	set:field=value
func ParseRule(expr string) (Rule, error) {
	parts := strings.SplitN(expr, ":", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("invalid transform expression %q: missing op prefix", expr)
	}
	op := Op(strings.ToLower(parts[0]))
	rest := parts[1]

	switch op {
	case OpDrop:
		if rest == "" {
			return Rule{}, fmt.Errorf("drop transform requires a field name")
		}
		return Rule{Op: OpDrop, Field: rest}, nil
	case OpRename, OpSet:
		kv := strings.SplitN(rest, "=", 2)
		if len(kv) != 2 || kv[0] == "" {
			return Rule{}, fmt.Errorf("%s transform requires field=value format", op)
		}
		return Rule{Op: op, Field: kv[0], Value: kv[1]}, nil
	default:
		return Rule{}, fmt.Errorf("unknown transform op %q", op)
	}
}

// Apply applies a slice of Rules to a single filter.Entry, returning a new entry.
func Apply(entry filter.Entry, rules []Rule) filter.Entry {
	out := filter.Entry{
		Timestamp: entry.Timestamp,
		Fields:    make(map[string]string, len(entry.Fields)),
	}
	for k, v := range entry.Fields {
		out.Fields[k] = v
	}

	for _, r := range rules {
		switch r.Op {
		case OpDrop:
			delete(out.Fields, r.Field)
		case OpRename:
			if val, ok := out.Fields[r.Field]; ok {
				out.Fields[r.Value] = val
				delete(out.Fields, r.Field)
			}
		case OpSet:
			out.Fields[r.Field] = r.Value
		}
	}
	return out
}

// Run applies rules to all entries and returns the transformed slice.
func Run(entries []filter.Entry, rules []Rule) []filter.Entry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]filter.Entry, len(entries))
	for i, e := range entries {
		out[i] = Apply(e, rules)
	}
	return out
}
