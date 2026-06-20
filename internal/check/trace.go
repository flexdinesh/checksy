package check

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// TraceURL returns the public egress IP via Cloudflare's cdn-cgi/trace.
const TraceURL = "https://www.cloudflare.com/cdn-cgi/trace"

// Facts are discovered details shown in the verdict header, not the table.
type Facts struct {
	PublicIP  string
	Resolver  string
	TraceBody string
}

// Discover fetches the public egress IP and reads the system resolver. Either
// field is left blank on failure — the header degrades gracefully.
func Discover(ctx context.Context, timeout time.Duration) Facts {
	traceCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	body := fetchTraceBody(traceCtx, TraceURL)
	return Facts{
		PublicIP:  parseTrace(body),
		Resolver:  systemResolver("/etc/resolv.conf"),
		TraceBody: body,
	}
}

// fetchTraceBody GETs the trace body and returns it raw.
func fetchTraceBody(ctx context.Context, url string) string {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// parseTrace extracts the ip= value from a cdn-cgi/trace body.
func parseTrace(body string) string {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if ip, ok := strings.CutPrefix(line, "ip="); ok {
			return ip
		}
	}
	return ""
}

// systemResolver reads the first nameserver from a resolv.conf-style file.
func systemResolver(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if ns, ok := strings.CutPrefix(line, "nameserver "); ok {
			ns = strings.TrimSpace(ns)
			if ns != "" {
				return ns
			}
		}
	}
	return ""
}
