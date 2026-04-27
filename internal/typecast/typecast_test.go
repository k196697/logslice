package typecast

import (
	"testing"
)

func makeEntry(ts int64, fields map[string]interface{}) Entry {
	return Entry{Timestamp: ts, Fields: fields}
}

func TestParseRule_Valid(t *testing.T) {
	for _, tc := range []struct{ expr, field, typ string }{
		{"status:int", "status", "int"},
		{"latency:float", "latency", "float"},
		{"ok:bool", "ok", "bool"},
		{"msg:string", "msg", "string"},
	} {
		r, err := ParseRule(tc.expr)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.expr, err)
		}
		if r.Field != tc.field || r.TargetType != tc.typ {
			t.Errorf("got %+v, want field=%s type=%s", r, tc.field, tc.typ)
		}
	}
}

func TestParseRule_MissingColon(t *testing.T) {
	_, err := ParseRule("statusint")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseRule_UnsupportedType(t *testing.T) {
	_, err := ParseRule("field:uuid")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestParseRule_EmptyField(t *testing.T) {
	_, err := ParseRule(":int")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestRun_CastsToInt(t *testing.T) {
	entries := []Entry{makeEntry(1, map[string]interface{}{"status": "200"})}
	opts := Options{Rules: []Rule{{Field: "status", TargetType: "int"}}}
	out := Run(entries, opts)
	if v, ok := out[0].Fields["status"].(int64); !ok || v != 200 {
		t.Errorf("expected int64(200), got %T %v", out[0].Fields["status"], out[0].Fields["status"])
	}
}

func TestRun_CastsToFloat(t *testing.T) {
	entries := []Entry{makeEntry(1, map[string]interface{}{"latency": "1.23"})}
	opts := Options{Rules: []Rule{{Field: "latency", TargetType: "float"}}}
	out := Run(entries, opts)
	if v, ok := out[0].Fields["latency"].(float64); !ok || v != 1.23 {
		t.Errorf("expected float64(1.23), got %T %v", out[0].Fields["latency"], out[0].Fields["latency"])
	}
}

func TestRun_CastsToBool(t *testing.T) {
	entries := []Entry{makeEntry(1, map[string]interface{}{"ok": "true"})}
	opts := Options{Rules: []Rule{{Field: "ok", TargetType: "bool"}}}
	out := Run(entries, opts)
	if v, ok := out[0].Fields["ok"].(bool); !ok || !v {
		t.Errorf("expected bool(true), got %T %v", out[0].Fields["ok"], out[0].Fields["ok"])
	}
}

func TestRun_SkipsMissingField(t *testing.T) {
	entries := []Entry{makeEntry(1, map[string]interface{}{"msg": "hello"})}
	opts := Options{Rules: []Rule{{Field: "status", TargetType: "int"}}}
	out := Run(entries, opts)
	if _, ok := out[0].Fields["status"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestRun_InvalidValueSkipped(t *testing.T) {
	entries := []Entry{makeEntry(1, map[string]interface{}{"status": "not-a-number"})}
	opts := Options{Rules: []Rule{{Field: "status", TargetType: "int"}}}
	out := Run(entries, opts)
	if v := out[0].Fields["status"]; v != "not-a-number" {
		t.Errorf("expected original value preserved, got %v", v)
	}
}

func TestRun_EmptyInput(t *testing.T) {
	out := Run(nil, Options{Rules: []Rule{{Field: "x", TargetType: "int"}}})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
