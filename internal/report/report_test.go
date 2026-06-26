package report_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/flexdinesh/checksy/internal/check"
	"github.com/flexdinesh/checksy/internal/report"
)

func TestRunWritesOneShotTerminalReport(t *testing.T) {
	results := []check.Result{{
		Kind:    check.KindHTTP,
		Label:   "gstatic.com",
		Status:  check.StatusOK,
		Latency: 142700 * time.Microsecond,
		Detail:  "204",
	}}
	facts := check.Facts{PublicIP: "203.0.113.42", LocalIP: "192.168.1.23", Gateway: "192.168.1.1", Resolver: "192.168.1.1"}

	var out bytes.Buffer
	if err := report.Run(&out, results, facts, false); err != nil {
		t.Fatalf("expected report to render, got %v", err)
	}

	got := out.String()
	for _, want := range []string{
		"checksy • internet UP",
		"ip 203.0.113.42",
		"local ip 192.168.1.23",
		"gateway 192.168.1.1",
		"resolver 192.168.1.1",
		"TARGET",
		"gstatic.com",
		"http",
		"✓",
		"142.7ms",
		"204",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected report to contain %q, got:\n%s", want, got)
		}
	}
}

func TestRunUsesCompactStatusMarkers(t *testing.T) {
	results := []check.Result{
		{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Detail: "204"},
		{Kind: check.KindDNS, Label: "one.one.one.one", Status: check.StatusFail, Detail: "lookup failed"},
	}

	var out bytes.Buffer
	if err := report.Run(&out, results, check.Facts{}, false); err != nil {
		t.Fatalf("expected report to render, got %v", err)
	}

	got := out.String()
	for _, want := range []string{"✓", "✗"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected report to contain status marker %q, got:\n%s", want, got)
		}
	}
}
