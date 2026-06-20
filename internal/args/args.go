package args

import (
	"fmt"
	"strings"
	"time"
)

const DefaultTimeout = 5 * time.Second

const Usage = `Usage: checksy [options]

Check internet connectivity against public targets.

Options:
  --exit-code         Silence output and exit 0 (up) or 1 (down); argument errors exit 2
  --timeout <dur>     Per-check timeout (default: 5s)
  --verbose           Show method, full error text, and raw response bodies
  --help              Show this help message
  --version           Show package version
`

type Options struct {
	ExitCode bool
	Verbose  bool
	Timeout  time.Duration
	Help     bool
	Version  bool
}

type ParseResult struct {
	OK      bool
	Options Options
	Err     error
}

func Parse(argv []string) ParseResult {
	options := Options{
		Timeout: DefaultTimeout,
	}

	for index := 0; index < len(argv); index++ {
		arg := argv[index]

		switch {
		case arg == "--":
			continue
		case arg == "--exit-code":
			options.ExitCode = true
		case arg == "--verbose":
			options.Verbose = true
		case arg == "--help":
			options.Help = true
		case arg == "--version":
			options.Version = true
		case arg == "--timeout":
			value, ok := nextValue(argv, index)
			if !ok {
				return parseError("Missing value for --timeout")
			}
			parsed, err := time.ParseDuration(value)
			if err != nil {
				return parseError(fmt.Sprintf("Invalid --timeout value: %s", value))
			}
			options.Timeout = parsed
			index++
		case strings.HasPrefix(arg, "--timeout="):
			value := strings.TrimPrefix(arg, "--timeout=")
			parsed, err := time.ParseDuration(value)
			if err != nil {
				return parseError(fmt.Sprintf("Invalid --timeout value: %s", value))
			}
			options.Timeout = parsed
		default:
			return parseError(fmt.Sprintf("Unknown argument: %s", arg))
		}
	}

	return ParseResult{OK: true, Options: options}
}

func nextValue(argv []string, index int) (string, bool) {
	if index+1 >= len(argv) || strings.HasPrefix(argv[index+1], "--") {
		return "", false
	}
	return argv[index+1], true
}

func parseError(message string) ParseResult {
	return ParseResult{OK: false, Err: fmt.Errorf("%s", message)}
}
