package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/flexdinesh/checksy/internal/check"
	"github.com/flexdinesh/checksy/internal/ui"
)

// Run writes a one-shot terminal report for the given results and facts.
// verbose switches on extra detail.
func Run(out io.Writer, results []check.Result, facts check.Facts, verbose bool) error {
	_, err := fmt.Fprintln(out, View(results, facts, verbose))
	return err
}

func View(results []check.Result, facts check.Facts, verbose bool) string {
	sections := []string{titleStyle().Render(ui.Title(results))}
	if facts := ui.FactsLine(facts); facts != "" {
		sections = append(sections, factsStyle().Render(facts))
	}
	if verbose && facts.TraceBody != "" {
		sections = append(sections, traceStyle().Render(strings.TrimRight(facts.TraceBody, "\n")))
	}
	sections = append(sections, resultsTable(results, verbose))
	return strings.Join(sections, "\n")
}

func resultsTable(results []check.Result, verbose bool) string {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(borderStyle()).
		StyleFunc(tableStyle).
		Headers("TARGET", "CHECK", "STATUS", "LATENCY", "DETAIL")

	for _, result := range results {
		t.Row(
			result.Label,
			string(result.Kind),
			statusCell(result),
			ui.FormatLatency(result.Latency),
			ui.Detail(result, verbose),
		)
	}
	if len(results) == 0 {
		t.Row("", "", "", "", "no checks ran")
	}

	return t.String()
}

func statusCell(result check.Result) string {
	switch result.Status {
	case check.StatusOK:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✓")
	case check.StatusFail:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("✗")
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("•")
	}
}

func tableStyle(row, col int) lipgloss.Style {
	style := lipgloss.NewStyle().Padding(0, 1)
	if row == table.HeaderRow {
		return style.Foreground(lipgloss.Color("8")).Bold(true).Align(lipgloss.Center)
	}
	return style
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
}

func factsStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
}

func traceStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
}

func borderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
}
