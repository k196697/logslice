package tail

import (
	"strings"
	"testing"
)

func lines(n int, total int) string {
	var sb strings.Builder
	for i := 1; i <= total; i++ {
		if i > 1 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strings.Repeat("x", i)) // unique content per line
	}
	_ = n
	return sb.String()
}

func TestReadLastN_FewerLinesThanN(t *testing.T) {
	input := "line1\nline2\nline3"
	got, err := ReadLastN(strings.NewReader(input), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "line1" || got[2] != "line3" {
		t.Errorf("unexpected content: %v", got)
	}
}

func TestReadLastN_ExactN(t *testing.T) {
	input := "a\nb\nc"
	got, err := ReadLastN(strings.NewReader(input), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestReadLastN_MoreLinesThanN(t *testing.T) {
	input := "one\ntwo\nthree\nfour\nfive"
	got, err := ReadLastN(strings.NewReader(input), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
	if got[0] != "three" || got[1] != "four" || got[2] != "five" {
		t.Errorf("wrong last lines: %v", got)
	}
}

func TestReadLastN_EmptyInput(t *testing.T) {
	got, err := ReadLastN(strings.NewReader(""), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestReadLastN_SkipsBlankLines(t *testing.T) {
	input := "alpha\n\nbeta\n\ngamma"
	got, err := ReadLastN(strings.NewReader(input), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 non-blank lines, got %d: %v", len(got), got)
	}
}

func TestReadLastN_InvalidN(t *testing.T) {
	_, err := ReadLastN(strings.NewReader("hello"), 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.N != 20 {
		t.Errorf("expected default N=20, got %d", opts.N)
	}
}
