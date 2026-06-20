package check

import (
	"context"
	"net"
	"time"
)

// Lookup resolves a host to IP addresses. Injectable so DNS result-shaping is
// testable without the real resolver.
type Lookup func(ctx context.Context, host string) ([]string, error)

// DNSTarget is the hostname checksy resolves to test the system resolver.
const DNSTarget = "one.one.one.one"

// DNS resolves host via the given resolver and reports the first resolved IP
// and the resolution latency. DNS is a diagnostic — it never flips the verdict.
func DNS(ctx context.Context, host string, lookup Lookup, timeout time.Duration) Result {
	start := time.Now()
	addrs, err := lookup(ctx, host)
	latency := time.Since(start)
	if err != nil {
		return Result{Kind: KindDNS, Label: host, Status: StatusFail, Latency: latency, Method: "system", Detail: shortErr(err, timeout), Err: err}
	}
	if len(addrs) == 0 {
		return Result{Kind: KindDNS, Label: host, Status: StatusFail, Latency: latency, Method: "system", Detail: "no addresses"}
	}
	return Result{
		Kind:    KindDNS,
		Label:   host,
		Status:  StatusOK,
		Latency: latency,
		Method:  "system",
		Detail:  addrs[0],
	}
}

// systemLookup resolves host using the default (system) resolver.
func systemLookup(ctx context.Context, host string) ([]string, error) {
	return net.DefaultResolver.LookupHost(ctx, host)
}
