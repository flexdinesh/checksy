package report

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/flexdinesh/checksy/internal/check"
	"github.com/muesli/termenv"
)

func testDisplay(width int, theme string, profile termenv.Profile) display {
	r := lipgloss.NewRenderer(io.Discard)
	r.SetColorProfile(profile)
	return themedDisplay(r, width, theme)
}

func TestReceiptPreservesVerdictAndDiagnosticFailure(t *testing.T) {
	results := []check.Result{
		{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Latency: 12 * time.Millisecond, Detail: "204"},
		{Kind: check.KindDNS, Label: "one.one.one.one", Status: check.StatusFail, Detail: "lookup failed"},
	}
	got := testDisplay(80, "", termenv.Ascii).view(results, check.Facts{PublicIP: "203.0.113.42"}, false)
	for _, want := range []string{"internet UP", "✓ ok", "✗ fail", "lookup failed", "Public IP 203.0.113.42"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q:\n%s", want, got)
		}
	}
	if strings.Index(got, "gstatic.com") > strings.Index(got, "Public IP") {
		t.Fatal("network context should follow checks")
	}
}

func TestReceiptFitsTerminalWidth(t *testing.T) {
	results := []check.Result{
		{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusFail, Detail: "timeout after 5s", Err: errors.New("Get https://connectivitycheck.gstatic.com/generate_204: context deadline exceeded")},
		{Kind: check.KindDNS, Label: "one.one.one.one", Status: check.StatusOK, Latency: 123456 * time.Microsecond, Detail: "2606:4700:4700::1111"},
	}
	facts := check.Facts{
		PublicIP: "2001:db8:1234:5678:abcd:ef01:2345:6789", LocalIP: "2001:db8::2",
		Gateway: "fe80::1234", Resolver: "2606:4700:4700::1111",
		TraceBody: "ip=2001:db8:1234:5678:abcd:ef01:2345:6789\ncolo=SYD\n",
	}
	for _, width := range []int{20, 40, 60, 80, 120} {
		for _, theme := range []string{"light", "dark"} {
			got := testDisplay(width, theme, termenv.TrueColor).view(results, facts, true)
			for _, line := range strings.Split(got, "\n") {
				if size := ansi.StringWidth(line); size > width {
					t.Errorf("%s/%d: line occupies %d columns: %q", theme, width, size, line)
				}
			}
			for _, want := range []string{"DOWN", "Trace", "Details", "colo=SYD", "123.5ms"} {
				if !strings.Contains(ansi.Strip(got), want) {
					t.Errorf("%s/%d: missing %q", theme, width, want)
				}
			}
		}
	}
}

func TestReceiptThemesPreservePlainContent(t *testing.T) {
	results := []check.Result{{Kind: check.KindHTTP, Label: "gstatic.com", Status: check.StatusOK, Detail: "204"}}
	plain := testDisplay(80, "", termenv.Ascii).view(results, check.Facts{}, false)
	light := testDisplay(80, "light", termenv.TrueColor).view(results, check.Facts{}, false)
	dark := testDisplay(80, "dark", termenv.TrueColor).view(results, check.Facts{}, false)
	if light == dark {
		t.Fatal("light and dark output should use distinct shades")
	}
	for _, got := range []string{light, dark} {
		if ansi.Strip(got) != plain {
			t.Fatal("color must not change the report's text or geometry")
		}
	}
}

func TestUnavailableFactsAreOmitted(t *testing.T) {
	got := testDisplay(80, "", termenv.Ascii).view(nil, check.Facts{}, false)
	for _, want := range []string{"internet DOWN", "no checks ran"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
	for _, absent := range []string{"Public IP", "Local IP", "Gateway", "Resolver", "Trace"} {
		if strings.Contains(got, absent) {
			t.Errorf("invented unavailable fact %q", absent)
		}
	}
}

func TestTerminalTheme(t *testing.T) {
	for _, tc := range []struct{ override, fgbg, want string }{
		{"", "0;15", "light"}, {"", "15;0", "dark"}, {"dark", "0;15", "dark"},
		{"light", "15;0", "light"}, {"", "bad", ""}, {"auto", "15;0", "dark"},
	} {
		t.Run(tc.override+"/"+tc.fgbg, func(t *testing.T) {
			t.Setenv("CHECKSY_THEME", tc.override)
			t.Setenv("COLORFGBG", tc.fgbg)
			if got := terminalTheme(); got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestPipedOutputStaysPlainAndReadable(t *testing.T) {
	t.Setenv("COLUMNS", "40")
	t.Setenv("CLICOLOR_FORCE", "1")
	t.Setenv("COLORTERM", "truecolor")
	t.Setenv("CHECKSY_THEME", "light")
	var out bytes.Buffer
	if err := Help(&out); err != nil {
		t.Fatal(err)
	}
	for _, flag := range []string{"--help", "--version", "--timeout <dur>", "--verbose", "--exit-code"} {
		if !strings.Contains(out.String(), flag) {
			t.Errorf("help omits %q", flag)
		}
	}
	if err := Version(&out, "v1.2.3"); err != nil {
		t.Fatal(err)
	}
	if err := Error(&out, errors.New("Missing value for --timeout")); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if strings.Contains(got, "\x1b") {
		t.Fatal("piped output contains ANSI escapes")
	}
	for _, want := range []string{"Report", "Automation", "Information", "checksy v1.2.3", "Missing value for --timeout", "checksy --timeout 2s"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
	for _, line := range strings.Split(got, "\n") {
		if ansi.StringWidth(line) > 40 {
			t.Errorf("line exceeds piped width: %q", line)
		}
	}
}
