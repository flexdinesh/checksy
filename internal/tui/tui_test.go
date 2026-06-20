package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/flexdinesh/checksy/internal/check"
)

func TestViewRendersUpResults(t *testing.T) {
	results := []check.Result{{
		Kind:    check.KindHTTP,
		Label:   "gstatic.com",
		Status:  check.StatusOK,
		Latency: 142700 * time.Microsecond,
		Detail:  "204",
	}}
	view := NewModel(results, check.Facts{}, false).View()
	for _, want := range []string{"internet UP", "gstatic.com", "142.7ms", "204", "✓"} {
		if !strings.Contains(view, want) {
			t.Errorf("expected view to contain %q", want)
		}
	}
}

func TestViewRendersDownResults(t *testing.T) {
	results := []check.Result{{
		Kind:   check.KindHTTP,
		Label:  "gstatic.com",
		Status: check.StatusFail,
		Detail: "timeout after 5s",
	}}
	view := NewModel(results, check.Facts{}, false).View()
	for _, want := range []string{"internet DOWN", "gstatic.com", "✗"} {
		if !strings.Contains(view, want) {
			t.Errorf("expected view to contain %q", want)
		}
	}
}

func TestViewRendersFactsInHeader(t *testing.T) {
	results := []check.Result{{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Detail: "204"}}
	facts := check.Facts{PublicIP: "203.0.113.42", Resolver: "192.168.1.1"}
	view := NewModel(results, facts, false).View()
	for _, want := range []string{"203.0.113.42", "192.168.1.1"} {
		if !strings.Contains(view, want) {
			t.Errorf("expected view to contain %q", want)
		}
	}
}

func TestViewOmitsFactsWhenBlank(t *testing.T) {
	results := []check.Result{{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Detail: "204"}}
	view := NewModel(results, check.Facts{}, false).View()
	if strings.Contains(view, "ip ") || strings.Contains(view, "resolver ") {
		t.Fatalf("expected no facts line when blank, got:\n%s", view)
	}
}

func TestViewRendersPingMethodInDetail(t *testing.T) {
	results := []check.Result{{Kind: check.KindPing, Label: "1.1.1.1", Status: check.StatusOK, Method: "icmp", Latency: 12 * time.Millisecond}}
	view := NewModel(results, check.Facts{}, false).View()
	if !strings.Contains(view, "icmp") {
		t.Fatalf("expected view to show ping method icmp, got:\n%s", view)
	}
}

func TestViewRendersTraceBodyInVerbose(t *testing.T) {
	results := []check.Result{{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Detail: "204"}}
	facts := check.Facts{TraceBody: "fl=abc\nip=203.0.113.42\ncolo=SFO\n"}
	view := NewModel(results, facts, true).View()
	if !strings.Contains(view, "fl=abc") {
		t.Fatalf("expected view to show trace body in verbose mode, got:\n%s", view)
	}
}

func TestViewRendersMethodForEveryRowInVerbose(t *testing.T) {
	results := []check.Result{
		{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Method: "http", Detail: "204"},
		{Kind: check.KindDNS, Label: "one.one.one.one", Status: check.StatusOK, Method: "system", Detail: "1.1.1.1"},
		{Kind: check.KindPing, Label: "1.1.1.1", Status: check.StatusOK, Method: "tcp"},
	}
	view := NewModel(results, check.Facts{}, true).View()
	for _, want := range []string{"http • 204", "system • 1.1.1.1", "tcp"} {
		if !strings.Contains(view, want) {
			t.Errorf("expected view to contain %q", want)
		}
	}
}

func TestViewHidesTraceBodyWhenNotVerbose(t *testing.T) {
	results := []check.Result{{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Detail: "204"}}
	facts := check.Facts{TraceBody: "fl=abc\nip=203.0.113.42\n"}
	view := NewModel(results, facts, false).View()
	if strings.Contains(view, "fl=abc") {
		t.Fatalf("expected trace body hidden in non-verbose mode, got:\n%s", view)
	}
}
