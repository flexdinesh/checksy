package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/flexdinesh/checksy/internal/check"
)

// Title renders the verdict header line, e.g. "checksy • internet UP".
func Title(results []check.Result) string {
	verdict := "DOWN"
	if check.Verdict(results) == check.StatusOK {
		verdict = "UP"
	}
	return "checksy • internet " + verdict
}

// FormatLatency renders a duration as milliseconds with one decimal place.
// A zero duration (e.g. a failed check with no measured latency) renders blank.
func FormatLatency(d time.Duration) string {
	if d <= 0 {
		return ""
	}
	return fmt.Sprintf("%.1fms", float64(d.Microseconds())/1000.0)
}

// FactsLine renders discovered network facts for the header.
// Blank fields are omitted; the line is empty when nothing was discovered.
func FactsLine(facts check.Facts) string {
	parts := []string{}
	if facts.PublicIP != "" {
		parts = append(parts, "ip "+facts.PublicIP)
	}
	if facts.LocalIP != "" {
		parts = append(parts, "local ip "+facts.LocalIP)
	}
	if facts.Gateway != "" {
		parts = append(parts, "gateway "+facts.Gateway)
	}
	if facts.Resolver != "" {
		parts = append(parts, "resolver "+facts.Resolver)
	}
	return strings.Join(parts, "\n")
}

// Detail shapes the polymorphic detail column for a result. By default a
// successful ping shows its method (icmp/tcp), DNS shows the resolved IP, and
// HTTP shows the status code; a failed row shows a short reason.
func Detail(r check.Result, verbose bool) string {
	if verbose {
		return verboseDetail(r)
	}
	if r.Status == check.StatusFail {
		return r.Detail
	}
	if r.Kind == check.KindPing {
		return r.Method
	}
	return r.Detail
}

func verboseDetail(r check.Result) string {
	method := r.Method
	if method == "" {
		method = string(r.Kind)
	}

	detail := r.Detail
	if r.Status == check.StatusFail && r.Err != nil {
		detail = r.Err.Error()
	}
	if detail == "" {
		return method
	}
	if method == "" {
		return detail
	}
	return method + " • " + detail
}
