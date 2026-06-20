package args

import (
	"testing"
	"time"
)

func TestParseReturnsDefaults(t *testing.T) {
	result := Parse(nil)
	if !result.OK {
		t.Fatalf("Parse returned error: %v", result.Err)
	}
	if result.Options.Timeout != DefaultTimeout {
		t.Fatalf("expected Timeout %v, got %v", DefaultTimeout, result.Options.Timeout)
	}
	if result.Options.ExitCode || result.Options.Verbose || result.Options.Help || result.Options.Version {
		t.Fatalf("expected boolean flags to default false: %+v", result.Options)
	}
}

func TestParseSupportsFlags(t *testing.T) {
	result := Parse([]string{"--exit-code", "--verbose", "--help", "--version"})
	if !result.OK {
		t.Fatalf("Parse returned error: %v", result.Err)
	}
	if !result.Options.ExitCode || !result.Options.Verbose || !result.Options.Help || !result.Options.Version {
		t.Fatalf("expected all flags true: %+v", result.Options)
	}
}

func TestParseSupportsTimeout(t *testing.T) {
	tests := []struct {
		name string
		argv []string
		want time.Duration
	}{
		{name: "separate value", argv: []string{"--timeout", "2s"}, want: 2 * time.Second},
		{name: "equals value", argv: []string{"--timeout=1500ms"}, want: 1500 * time.Millisecond},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Parse(test.argv)
			if !result.OK {
				t.Fatalf("Parse returned error: %v", result.Err)
			}
			if result.Options.Timeout != test.want {
				t.Fatalf("expected Timeout %v, got %v", test.want, result.Options.Timeout)
			}
		})
	}
}

func TestParseRejectsInvalidTimeout(t *testing.T) {
	for _, argv := range [][]string{{"--timeout", "abc"}, {"--timeout"}} {
		if Parse(argv).OK {
			t.Fatalf("expected Parse(%v) to fail", argv)
		}
	}
}

func TestParseRejectsUnknownArgs(t *testing.T) {
	for _, argv := range [][]string{{"--raw"}, {"--bogus"}} {
		if Parse(argv).OK {
			t.Fatalf("expected Parse(%v) to fail", argv)
		}
	}
}
