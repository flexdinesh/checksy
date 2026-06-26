package ui

import (
	"errors"
	"testing"
	"time"

	"github.com/flexdinesh/checksy/internal/check"
)

func TestTitle(t *testing.T) {
	tests := []struct {
		name    string
		results []check.Result
		want    string
	}{
		{name: "up", results: []check.Result{{Kind: check.KindHTTP, Status: check.StatusOK}}, want: "checksy • internet UP"},
		{name: "down", results: []check.Result{{Kind: check.KindHTTP, Status: check.StatusFail}}, want: "checksy • internet DOWN"},
		{name: "empty is down", results: nil, want: "checksy • internet DOWN"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Title(test.results); got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}

func TestFormatLatency(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{name: "milliseconds", d: 12300 * time.Microsecond, want: "12.3ms"},
		{name: "sub-millisecond", d: 800 * time.Microsecond, want: "0.8ms"},
		{name: "one and a half", d: 1500 * time.Microsecond, want: "1.5ms"},
		{name: "zero is blank", d: 0, want: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FormatLatency(test.d); got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}

func TestFactsLine(t *testing.T) {
	tests := []struct {
		name  string
		facts check.Facts
		want  string
	}{
		{name: "all", facts: check.Facts{PublicIP: "203.0.113.42", LocalIP: "192.168.1.23", Gateway: "192.168.1.1", Resolver: "192.168.1.1"}, want: "ip 203.0.113.42\nlocal ip 192.168.1.23\ngateway 192.168.1.1\nresolver 192.168.1.1"},
		{name: "only ip", facts: check.Facts{PublicIP: "203.0.113.42"}, want: "ip 203.0.113.42"},
		{name: "only resolver", facts: check.Facts{Resolver: "192.168.1.1"}, want: "resolver 192.168.1.1"},
		{name: "neither is blank", facts: check.Facts{}, want: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := FactsLine(test.facts)
			if got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}

func TestDetail(t *testing.T) {
	tests := []struct {
		name    string
		result  check.Result
		verbose bool
		want    string
	}{
		{name: "ping success shows method", result: check.Result{Kind: check.KindPing, Status: check.StatusOK, Method: "icmp"}, want: "icmp"},
		{name: "ping tcp fallback shows method", result: check.Result{Kind: check.KindPing, Status: check.StatusOK, Method: "tcp"}, want: "tcp"},
		{name: "dns success shows resolved ip", result: check.Result{Kind: check.KindDNS, Status: check.StatusOK, Detail: "1.1.1.1"}, want: "1.1.1.1"},
		{name: "http success shows status code", result: check.Result{Kind: check.KindHTTP, Status: check.StatusOK, Detail: "204"}, want: "204"},
		{name: "failure shows short reason", result: check.Result{Kind: check.KindHTTP, Status: check.StatusFail, Detail: "timeout after 5s"}, want: "timeout after 5s"},
		{name: "ping success blank method is blank", result: check.Result{Kind: check.KindPing, Status: check.StatusOK, Method: ""}, want: ""},
		{name: "verbose failure shows method and full error text", result: check.Result{Kind: check.KindHTTP, Status: check.StatusFail, Method: "http", Detail: "timeout after 5s", Err: errors.New("Get http://example.com: context deadline exceeded")}, verbose: true, want: "http • Get http://example.com: context deadline exceeded"},
		{name: "verbose failure without err falls back to detail", result: check.Result{Kind: check.KindHTTP, Status: check.StatusFail, Detail: "timeout after 5s"}, verbose: true, want: "http • timeout after 5s"},
		{name: "verbose success still shows method for ping", result: check.Result{Kind: check.KindPing, Status: check.StatusOK, Method: "icmp"}, verbose: true, want: "icmp"},
		{name: "verbose dns success shows method and detail", result: check.Result{Kind: check.KindDNS, Status: check.StatusOK, Method: "system", Detail: "1.1.1.1"}, verbose: true, want: "system • 1.1.1.1"},
		{name: "verbose http success shows method and detail", result: check.Result{Kind: check.KindHTTP, Status: check.StatusOK, Method: "http", Detail: "204"}, verbose: true, want: "http • 204"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Detail(test.result, test.verbose); got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}
