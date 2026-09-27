package main

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/flexdinesh/checksy/internal/check"
	"github.com/flexdinesh/checksy/internal/report"
)

func TestRunExitCodeIsSilentAndZeroWhenUp(t *testing.T) {
	runner := func(context.Context, time.Duration) []check.Result {
		return []check.Result{{Kind: check.KindHTTP, Status: check.StatusOK}}
	}
	var out bytes.Buffer
	code := run([]string{"--exit-code"}, runner, &out)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("expected no output in --exit-code mode, got %q", out.String())
	}
}

func TestRunExitCodeIsOneAndSilentWhenDown(t *testing.T) {
	runner := func(context.Context, time.Duration) []check.Result {
		return []check.Result{{Kind: check.KindHTTP, Status: check.StatusFail}}
	}
	var out bytes.Buffer
	code := run([]string{"--exit-code"}, runner, &out)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("expected no output in --exit-code mode, got %q", out.String())
	}
}

func TestRunExitsTwoOnUnknownFlag(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"--bogus"}, func(context.Context, time.Duration) []check.Result { return nil }, &out)
	if code != 2 {
		t.Fatalf("expected exit code 2 for unknown flag, got %d", code)
	}
}

func TestRunHelpPrintsUsage(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"--help"}, func(context.Context, time.Duration) []check.Result { return nil }, &out)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(out.String(), "Usage:") {
		t.Fatalf("expected usage text, got %q", out.String())
	}
}

func TestRunVersionPrintsVersion(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"--version"}, func(context.Context, time.Duration) []check.Result { return nil }, &out)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.HasPrefix(out.String(), "checksy ") {
		t.Fatalf("expected version line, got %q", out.String())
	}
}

func TestRunRendersResultsWithDiscoveredFacts(t *testing.T) {
	runner := func(context.Context, time.Duration) []check.Result {
		return []check.Result{{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Detail: "204"}}
	}
	discover := func(context.Context, time.Duration) check.Facts {
		return check.Facts{PublicIP: "203.0.113.42", Resolver: "192.168.1.1"}
	}
	render := func(out io.Writer, results []check.Result, facts check.Facts, verbose bool) error {
		if verbose {
			t.Fatal("expected verbose false")
		}
		if len(results) != 1 || results[0].Label != "gstatic.com" {
			t.Fatalf("unexpected results: %+v", results)
		}
		if facts.PublicIP != "203.0.113.42" || facts.Resolver != "192.168.1.1" {
			t.Fatalf("unexpected facts: %+v", facts)
		}
		_, err := io.WriteString(out, "rendered")
		return err
	}

	var out bytes.Buffer
	code := runWithDeps(nil, runner, discover, render, &out, io.Discard)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if out.String() != "rendered" {
		t.Fatalf("expected renderer output, got %q", out.String())
	}
}

func TestRunPassesTimeoutAndVerboseToDependencies(t *testing.T) {
	var runnerTimeout time.Duration
	var discoverTimeout time.Duration
	var renderedVerbose bool

	runner := func(_ context.Context, timeout time.Duration) []check.Result {
		runnerTimeout = timeout
		return []check.Result{{Kind: check.KindHTTP, Status: check.StatusOK}}
	}
	discover := func(_ context.Context, timeout time.Duration) check.Facts {
		discoverTimeout = timeout
		return check.Facts{}
	}
	render := func(_ io.Writer, _ []check.Result, _ check.Facts, verbose bool) error {
		renderedVerbose = verbose
		return nil
	}

	code := runWithDeps([]string{"--timeout", "2s", "--verbose"}, runner, discover, render, io.Discard, io.Discard)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if runnerTimeout != 2*time.Second {
		t.Fatalf("expected runner timeout 2s, got %v", runnerTimeout)
	}
	if discoverTimeout != 2*time.Second {
		t.Fatalf("expected discover timeout 2s, got %v", discoverTimeout)
	}
	if !renderedVerbose {
		t.Fatal("expected verbose true to reach renderer")
	}
}

func TestArgumentErrorsUseStderrWithoutRunningChecks(t *testing.T) {
	for _, argv := range [][]string{{"--bogus"}, {"--timeout"}, {"--timeout", "abc"}, {"--exit-code", "--bogus"}} {
		var out, errOut bytes.Buffer
		runner := func(context.Context, time.Duration) []check.Result {
			t.Fatal("argument error must not run network checks")
			return nil
		}
		code := runWithDeps(argv, runner, check.Discover, report.Run, &out, &errOut)
		if code != 2 || out.Len() != 0 {
			t.Fatalf("want exit 2 and empty stdout; got %d, %q", code, out.String())
		}
		got := errOut.String()
		if !strings.Contains(got, "error") || !strings.Contains(got, "checksy --") {
			t.Fatalf("error needs reason and recovery: %q", got)
		}
	}
}

func TestReportModePreservesSuccessfulExitWhenInternetDown(t *testing.T) {
	runner := func(context.Context, time.Duration) []check.Result {
		return []check.Result{{Kind: check.KindHTTP, Status: check.StatusFail, Label: "gstatic.com", Detail: "timeout after 5s"}}
	}
	discover := func(context.Context, time.Duration) check.Facts { return check.Facts{} }
	var out, errOut bytes.Buffer
	code := runWithDeps(nil, runner, discover, report.Run, &out, &errOut)
	if code != 0 || !strings.Contains(out.String(), "internet DOWN") || errOut.Len() != 0 {
		t.Fatalf("normal report mode should display DOWN and exit 0; got %d, stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}
