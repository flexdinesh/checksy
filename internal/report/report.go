package report

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/flexdinesh/checksy/internal/check"
	"github.com/flexdinesh/checksy/internal/ui"
)

// Run writes a one-shot terminal report using the output's color and width capabilities.
func Run(out io.Writer, results []check.Result, facts check.Facts, verbose bool) error {
	_, err := fmt.Fprintln(out, newDisplay(out).view(results, facts, verbose))
	return err
}

func View(results []check.Result, facts check.Facts, verbose bool) string {
	return newDisplay(os.Stdout).view(results, facts, verbose)
}

func (d display) view(results []check.Result, facts check.Facts, verbose bool) string {
	verdict, color := "DOWN", d.failure
	if check.Verdict(results) == check.StatusOK {
		verdict, color = "UP", d.success
	}
	title := d.primary.Bold(true).Render("checksy") + " • internet " + color.Bold(true).Render(verdict)
	rows := d.resultRows(results)
	ruleWidth := ansi.StringWidth(title)
	for _, line := range rows {
		ruleWidth = max(ruleWidth, ansi.StringWidth(line))
	}
	lines := []string{d.wrap(title), d.neutral.Render(strings.Repeat("─", min(ruleWidth, d.width)))}
	lines = append(lines, rows...)
	if network := d.networkFacts(facts); len(network) > 0 {
		lines = append(lines, "")
		lines = append(lines, network...)
	}
	if verbose {
		var details []string
		for _, result := range results {
			if detail := ui.Detail(result, true); detail != ui.Detail(result, false) {
				details = append(details, d.wrap(d.primary.Render(result.Label)+"  "+detail))
			}
		}
		if len(details) > 0 {
			lines = append(lines, "", d.secondary.Bold(true).Render("Details"))
			lines = append(lines, details...)
		}
		if facts.TraceBody != "" {
			lines = append(lines, "", d.secondary.Bold(true).Render("Trace"), d.wrap(strings.TrimRight(facts.TraceBody, "\n")))
		}
	}
	return strings.Join(lines, "\n")
}

func (d display) resultRows(results []check.Result) []string {
	if len(results) == 0 {
		return []string{d.wrap(d.warning.Render("• no checks ran"))}
	}
	targetWidth, latencyWidth := 6, 7
	for _, result := range results {
		targetWidth = max(targetWidth, ansi.StringWidth(result.Label))
		latencyWidth = max(latencyWidth, ansi.StringWidth(ui.FormatLatency(result.Latency)))
	}
	targetWidth = min(targetWidth, max(6, d.width/3))
	widths := []int{6, 5, targetWidth, latencyWidth, d.width - 6 - 5 - targetWidth - latencyWidth - 8}
	if widths[4] < 10 {
		return d.narrowResults(results)
	}
	lines := columns([]string{
		d.neutral.Render("STATUS"), d.secondary.Render("CHECK"), d.neutral.Render("TARGET"), d.neutral.Render("LATENCY"), d.neutral.Render("DETAIL"),
	}, widths)
	for _, result := range results {
		latency := ui.FormatLatency(result.Latency)
		if latency == "" {
			latency = "—"
		}
		lines = append(lines, columns([]string{
			d.status(result.Status),
			d.secondary.Render(strings.ToUpper(string(result.Kind))),
			result.Label, latency, ui.Detail(result, false),
		}, widths)...)
	}
	return lines
}

func (d display) narrowResults(results []check.Result) []string {
	var lines []string
	for _, result := range results {
		header := d.status(result.Status) + "  " + d.secondary.Render(strings.ToUpper(string(result.Kind))) + "  " + result.Label
		lines = append(lines, d.wrap(header))
		var detail []string
		if latency := ui.FormatLatency(result.Latency); latency != "" {
			detail = append(detail, latency)
		}
		if text := ui.Detail(result, false); text != "" {
			detail = append(detail, text)
		}
		if len(detail) > 0 {
			lines = append(lines, d.indent(strings.Join(detail, "  •  ")))
		}
	}
	return lines
}

func (d display) status(status check.Status) string {
	switch status {
	case check.StatusOK:
		return d.success.Render("✓ ok")
	case check.StatusFail:
		return d.failure.Render("✗ fail")
	default:
		return d.warning.Render("• n/a")
	}
}

func (d display) networkFacts(facts check.Facts) []string {
	var entries []string
	for _, fact := range []struct{ label, value string }{
		{"Public IP", facts.PublicIP}, {"Local IP", facts.LocalIP},
		{"Gateway", facts.Gateway}, {"Resolver", facts.Resolver},
	} {
		if fact.value != "" {
			entries = append(entries, d.primary.Render(fact.label)+" "+fact.value)
		}
	}
	var lines []string
	leftWidth := 0
	for i := 0; i < len(entries); i += 2 {
		leftWidth = max(leftWidth, ansi.StringWidth(entries[i]))
	}
	for i := 0; i < len(entries); {
		if i+1 < len(entries) && leftWidth+ansi.StringWidth(entries[i+1])+3 <= d.width {
			gap := strings.Repeat(" ", leftWidth-ansi.StringWidth(entries[i])+3)
			lines = append(lines, entries[i]+gap+entries[i+1])
			i += 2
		} else {
			lines = append(lines, d.wrap(entries[i]))
			i++
		}
	}
	return lines
}

// columns wraps cells independently, retaining their positions on continuation lines.
func columns(cells []string, widths []int) []string {
	wrapped := make([][]string, len(cells))
	height := 1
	for i, cell := range cells {
		wrapped[i] = strings.Split(ansi.Wrap(cell, widths[i], ""), "\n")
		height = max(height, len(wrapped[i]))
	}
	var lines []string
	for row := 0; row < height; row++ {
		var line strings.Builder
		for col, width := range widths {
			cell := ""
			if row < len(wrapped[col]) {
				cell = wrapped[col][row]
			}
			if col == 3 {
				line.WriteString(strings.Repeat(" ", max(0, width-ansi.StringWidth(cell))))
			}
			line.WriteString(cell)
			if col < len(widths)-1 {
				if col != 3 {
					line.WriteString(strings.Repeat(" ", max(0, width-ansi.StringWidth(cell))))
				}
				line.WriteString("  ")
			}
		}
		lines = append(lines, strings.TrimRight(line.String(), " "))
	}
	return lines
}
