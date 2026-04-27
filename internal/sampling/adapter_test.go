package sampling

import (
	"testing"
)

func TestParseConfig_ValidRate(t *testing.T) {
	opts, err := ParseConfig(Config{RateStr: "0.25"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Rate != 0.25 {
		t.Errorf("expected rate 0.25, got %f", opts.Rate)
	}
}

func TestParseConfig_ValidEvery(t *testing.T) {
	opts, err := ParseConfig(Config{EveryStr: "5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Every != 5 {
		t.Errorf("expected every=5, got %d", opts.Every)
	}
}

func TestParseConfig_InvalidRate(t *testing.T) {
	_, err := ParseConfig(Config{RateStr: "abc"})
	if err == nil {
		t.Error("expected error for non-numeric rate")
	}
}

func TestParseConfig_RateOutOfRange(t *testing.T) {
	_, err := ParseConfig(Config{RateStr: "1.5"})
	if err == nil {
		t.Error("expected error for rate > 1.0")
	}
}

func TestParseConfig_RateNegative(t *testing.T) {
	_, err := ParseConfig(Config{RateStr: "-0.1"})
	if err == nil {
		t.Error("expected error for negative rate")
	}
}

func TestParseConfig_EveryLessThanOne(t *testing.T) {
	_, err := ParseConfig(Config{EveryStr: "0"})
	if err == nil {
		t.Error("expected error for every=0")
	}
}

func TestParseConfig_InvalidEvery(t *testing.T) {
	_, err := ParseConfig(Config{EveryStr: "abc"})
	if err == nil {
		t.Error("expected error for non-numeric every")
	}
}

func TestParseConfig_EmptyConfig(t *testing.T) {
	opts, err := ParseConfig(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Rate != 1.0 {
		t.Errorf("expected default rate 1.0, got %f", opts.Rate)
	}
	if opts.Every != 0 {
		t.Errorf("expected default every 0, got %d", opts.Every)
	}
}
