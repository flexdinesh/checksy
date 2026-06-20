package check

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusOK
	StatusFail
)

type Kind string

const (
	KindHTTP Kind = "http"
	KindDNS  Kind = "dns"
	KindPing Kind = "ping"
)

type Result struct {
	Kind    Kind
	Label   string
	Status  Status
	Latency time.Duration
	Method  string
	Detail  string
	Err     error
}

type Runner func(ctx context.Context, timeout time.Duration) []Result

// Verdict is UP (StatusOK) iff at least one HTTP check succeeded. Ping and DNS
// are diagnostics and never flip the verdict.
func Verdict(results []Result) Status {
	for _, r := range results {
		if r.Kind == KindHTTP && r.Status == StatusOK {
			return StatusOK
		}
	}
	return StatusFail
}

// All runs every check concurrently. Each check gets its own timeout-derived
// context, so slow checks are bounded independently and total time is roughly
// the slowest check, not the sum.
func All(ctx context.Context, timeout time.Duration) []Result {
	return runConcurrently(ctx,
		timeout,
		func(ctx context.Context) Result {
			return HTTP(ctx, VerdictLabel, VerdictURL, http.StatusNoContent, timeout)
		},
		func(ctx context.Context) Result { return Ping(ctx, "1.1.1.1", realICMP, realTCPDial, timeout) },
		func(ctx context.Context) Result { return Ping(ctx, "8.8.8.8", realICMP, realTCPDial, timeout) },
		func(ctx context.Context) Result { return DNS(ctx, DNSTarget, systemLookup, timeout) },
	)
}

func runConcurrently(ctx context.Context, timeout time.Duration, checks ...func(context.Context) Result) []Result {
	results := make([]Result, len(checks))
	var wait sync.WaitGroup
	for index, checkFn := range checks {
		wait.Add(1)
		go func(i int, fn func(context.Context) Result) {
			defer wait.Done()
			checkCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			results[i] = fn(checkCtx)
		}(index, checkFn)
	}
	wait.Wait()
	return results
}
