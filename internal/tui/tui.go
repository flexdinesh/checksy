package tui

import (
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/flexdinesh/checksy/internal/check"
	"github.com/flexdinesh/checksy/internal/ui"
)

const minWidth = 40

// Run launches the Bubble Tea program rendering the given results and facts,
// blocking until the user quits. verbose switches on extra detail.
func Run(out io.Writer, results []check.Result, facts check.Facts, verbose bool) error {
	program := tea.NewProgram(NewModel(results, facts, verbose), tea.WithAltScreen(), tea.WithOutput(out))
	_, err := program.Run()
	return err
}

type Model struct {
	results []check.Result
	facts   check.Facts
	verbose bool
	table   table.Model
	width   int
	height  int
}

func NewModel(results []check.Result, facts check.Facts, verbose bool) Model {
	model := Model{
		results: results,
		facts:   facts,
		verbose: verbose,
		table: table.New(
			table.WithFocused(true),
			table.WithStyles(tableStyles()),
		),
		width:  80,
		height: 24,
	}
	model.updateTable()
	return model
}

func (model Model) Init() tea.Cmd { return nil }

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" || msg.String() == "Q" {
			return model, tea.Quit
		}
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		model.updateTable()
		return model, nil
	}
	var cmd tea.Cmd
	model.table, cmd = model.table.Update(message)
	return model, cmd
}

func (model Model) View() string {
	width := max(model.width, minWidth)
	sections := []string{frameStyle(width, 1).Render(titleStyle(width).Render(ui.Title(model.results)))}
	if facts := ui.FactsLine(model.facts); facts != "" {
		sections = append(sections, frameStyle(width, 1).Render(factsStyle(width).Render(facts)))
	}
	if model.verbose && model.facts.TraceBody != "" {
		sections = append(sections, frameStyle(width, 1).Render(traceStyle(width).Render(model.facts.TraceBody)))
	}
	headerHeight := lipgloss.Height(strings.Join(sections, "\n"))
	tableHeight := max(4, model.height-headerHeight-1)
	sections = append(sections, model.renderTable(width, tableHeight))
	return strings.Join(sections, "\n")
}

func (model *Model) updateTable() {
	width := max(model.width, minWidth)
	model.table.SetColumns(columns(width))
	model.table.SetRows(model.rows())
	model.table.SetWidth(width - 4)
	model.table.SetHeight(max(4, len(model.results)+2))
}

func (model Model) rows() []table.Row {
	rows := make([]table.Row, 0, len(model.results))
	for _, result := range model.results {
		rows = append(rows, table.Row{
			result.Label,
			string(result.Kind),
			statusCell(result),
			ui.FormatLatency(result.Latency),
			ui.Detail(result, model.verbose),
		})
	}
	if len(rows) == 0 {
		rows = append(rows, table.Row{"", "", "", "", "no checks ran"})
	}
	return rows
}

func statusCell(result check.Result) string {
	switch result.Status {
	case check.StatusOK:
		return toneStyle("green", false).Render("✓")
	case check.StatusFail:
		return toneStyle("red", false).Render("✗")
	default:
		return toneStyle("gray", false).Render("•")
	}
}

func columns(width int) []table.Column {
	content := width - 4
	target := content * 30 / 100
	if target < 12 {
		target = 12
	}
	kind := 6
	status := 4
	latency := 10
	detail := content - target - kind - status - latency
	if detail < 8 {
		detail = 8
	}
	return []table.Column{
		{Title: "TARGET", Width: target},
		{Title: "CHECK", Width: kind},
		{Title: "STATUS", Width: status},
		{Title: "LATENCY", Width: latency},
		{Title: "DETAIL", Width: detail},
	}
}

func (model Model) renderTable(width int, height int) string {
	return frameStyle(width, height).Render(model.table.View())
}

func tableStyles() table.Styles {
	styles := table.DefaultStyles()
	styles.Header = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		BorderBottom(true)
	styles.Cell = lipgloss.NewStyle()
	styles.Selected = styles.Cell
	return styles
}

func frameStyle(width int, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(max(1, width-2)).
		Height(max(1, height)).
		Padding(0, 1)
}

func titleStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true).
		MaxWidth(max(1, width-4))
}

func factsStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		MaxWidth(max(1, width-4))
}

func traceStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		MaxWidth(max(1, width-4))
}

func toneStyle(tone string, bold bool) lipgloss.Style {
	style := lipgloss.NewStyle().Bold(bold)
	switch tone {
	case "green":
		return style.Foreground(lipgloss.Color("10"))
	case "red":
		return style.Foreground(lipgloss.Color("9"))
	case "gray":
		return style.Foreground(lipgloss.Color("8"))
	default:
		return style
	}
}
