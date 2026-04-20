// logslice is a CLI tool to filter and slice structured log files by time range or field values.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/logslice/internal/filter"
	"github.com/yourusername/logslice/internal/parser"
)

const timeLayout = time.RFC3339

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		fromStr    = flag.String("from", "", "Start of time range (RFC3339), e.g. 2024-01-01T00:00:00Z")
		toStr      = flag.String("to", "", "End of time range (RFC3339), e.g. 2024-01-02T00:00:00Z")
		tsField    = flag.String("ts-field", "time", "JSON field name to use as timestamp")
		fieldExpr  = flag.String("field", "", "Field filter expression in key=value format, e.g. level=error")
		inputFile  = flag.String("file", "", "Path to input log file (defaults to stdin)")
	)
	flag.Parse()

	// Parse optional time bounds.
	var from, to time.Time
	var err error

	if *fromStr != "" {
		from, err = time.Parse(timeLayout, *fromStr)
		if err != nil {
			return fmt.Errorf("invalid --from value %q: %w", *fromStr, err)
		}
	}
	if *toStr != "" {
		to, err = time.Parse(timeLayout, *toStr)
		if err != nil {
			return fmt.Errorf("invalid --to value %q: %w", *toStr, err)
		}
	}

	// Open input source.
	var input *os.File
	if *inputFile != "" {
		input, err = os.Open(*inputFile)
		if err != nil {
			return fmt.Errorf("opening file %q: %w", *inputFile, err)
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	// Parse JSON log lines.
	entries, err := parser.ParseJSONLines(input, *tsField)
	if err != nil {
		return fmt.Errorf("parsing log lines: %w", err)
	}

	// Build filter options.
	opts := filter.Options{
		From: from,
		To:   to,
	}

	if *fieldExpr != "" {
		key, value, ok := parseFieldExpr(*fieldExpr)
		if !ok {
			return fmt.Errorf("invalid --field expression %q: expected key=value", *fieldExpr)
		}
		opts.FieldKey = key
		opts.FieldValue = value
	}

	// Apply filters.
	filtered := filter.Apply(entries, opts)

	// Write matching entries to stdout.
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for _, entry := range filtered {
		line, err := entry.MarshalJSON()
		if err != nil {
			return fmt.Errorf("marshalling entry: %w", err)
		}
		writer.Write(line)
		writer.WriteByte('\n')
	}

	return nil
}

// parseFieldExpr splits a "key=value" expression into its components.
// Returns ok=false if the expression is not in the expected format.
func parseFieldExpr(expr string) (key, value string, ok bool) {
	for i, ch := range expr {
		if ch == '=' {
			return expr[:i], expr[i+1:], true
		}
	}
	return "", "", false
}
