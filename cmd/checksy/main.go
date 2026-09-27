package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/flexdinesh/checksy/internal/args"
	"github.com/flexdinesh/checksy/internal/check"
	"github.com/flexdinesh/checksy/internal/report"
	"github.com/flexdinesh/checksy/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:], check.All, os.Stdout))
}

func run(argv []string, runner check.Runner, out io.Writer) int {
	return runWithDeps(argv, runner, check.Discover, report.Run, out, os.Stderr)
}

type discoverer func(context.Context, time.Duration) check.Facts

type renderer func(io.Writer, []check.Result, check.Facts, bool) error

func runWithDeps(argv []string, runner check.Runner, discover discoverer, render renderer, out, errOut io.Writer) int {
	parsed := args.Parse(argv)
	if !parsed.OK {
		report.Error(errOut, parsed.Err)
		return 2
	}
	options := parsed.Options

	if options.Help {
		if err := report.Help(out); err != nil {
			return 1
		}
		return 0
	}
	if options.Version {
		if err := report.Version(out, version.String()); err != nil {
			return 1
		}
		return 0
	}

	ctx := context.Background()
	results := runner(ctx, options.Timeout)

	if options.ExitCode {
		if check.Verdict(results) == check.StatusOK {
			return 0
		}
		return 1
	}

	facts := discover(ctx, options.Timeout)
	if err := render(out, results, facts, options.Verbose); err != nil {
		return 1
	}
	return 0
}
