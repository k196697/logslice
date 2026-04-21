package pipeline

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/highlight"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
)

// Options holds the configuration for a pipeline run.
type Options struct {
	Format      string
	TimeField   string
	From        string
	To          string
	Fields      map[string]string
	OutputFmt   output.Format
	Highlights  []highlight.Rule
	TailN       int
}

// Runner executes the full parse → filter → highlight → output pipeline.
type Runner struct {
	opts Options
	out  io.Writer
}

// New creates a new pipeline Runner.
func New(opts Options, out io.Writer) *Runner {
	return &Runner{opts: opts, out: out}
}

// Run executes the pipeline reading lines from r.
func (r *Runner) Run(lines []string) error {
	if !parser.IsSupported(r.opts.Format) {
		return fmt.Errorf("unsupported format %q; supported: %s",
			r.opts.Format, parser.FormatDescription())
	}

	entries, err := parseLines(r.opts.Format, r.opts.TimeField, lines)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	filtered, err := applyFilter(r.opts, entries)
	if err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	hl := highlight.New(r.opts.Highlights)
	fmt := output.NewFormatter(r.opts.OutputFmt, r.opts.TimeField, hl)

	for _, e := range filtered {
		line, err := fmt.Format(e)
		if err != nil {
			return fmt.Errorf("format: %w", err)
		}
		if _, err := io.WriteString(r.out, line+"\n"); err != nil {
			return err
		}
	}
	return nil
}
