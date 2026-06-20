package check

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTraceExtractsPublicIP(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "typical body", body: "fl=abc\nh=www.cloudflare.com\nip=203.0.113.42\nts=123\nvisit_scheme=https\nuag=checksy\ncolo=SFO\n", want: "203.0.113.42"},
		{name: "ip first line", body: "ip=198.51.100.7\n", want: "198.51.100.7"},
		{name: "no ip line", body: "fl=abc\nh=example.com\n", want: ""},
		{name: "empty body", body: "", want: ""},
		{name: "ip with trailing whitespace", body: "ip=203.0.113.42\r\n", want: "203.0.113.42"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := parseTrace(test.body); got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}

func TestSystemResolverReadsFirstNameserver(t *testing.T) {
	body := "# generated\nnameserver 192.168.1.1\nnameserver 8.8.8.8\n"
	path := filepath.Join(t.TempDir(), "resolv.conf")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if got := systemResolver(path); got != "192.168.1.1" {
		t.Fatalf("expected 192.168.1.1, got %q", got)
	}
}

func TestSystemResolverEmptyWhenNoNameserver(t *testing.T) {
	body := "# generated\noptions edns0\n"
	path := filepath.Join(t.TempDir(), "resolv.conf")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if got := systemResolver(path); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestSystemResolverEmptyWhenFileMissing(t *testing.T) {
	if got := systemResolver(filepath.Join(t.TempDir(), "nope")); got != "" {
		t.Fatalf("expected empty for missing file, got %q", got)
	}
}
