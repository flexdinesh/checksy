package report

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/flexdinesh/checksy/internal/args"
	"github.com/muesli/termenv"
)

type display struct {
	width                                                  int
	primary, secondary, success, failure, warning, neutral lipgloss.Style
}

func newDisplay(out io.Writer) display {
	r := lipgloss.NewRenderer(out)
	width := 80
	if columns, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && columns > 0 {
		width = columns
	}
	file, ok := out.(interface{ Fd() uintptr })
	if ok && term.IsTerminal(file.Fd()) {
		if columns, _, err := term.GetSize(file.Fd()); err == nil && columns > 0 {
			width = columns
		}
	} else {
		r.SetColorProfile(termenv.Ascii)
	}
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		r.SetColorProfile(termenv.Ascii)
	}
	return themedDisplay(r, width, terminalTheme())
}

// Terminal palette colors follow the user's theme without querying stdin.
// Known themes additionally receive paired truecolor/256-color shades.
func themedDisplay(r *lipgloss.Renderer, width int, theme string) display {
	colors := []lipgloss.TerminalColor{
		lipgloss.Color("4"), lipgloss.Color("5"), lipgloss.Color("2"),
		lipgloss.Color("1"), lipgloss.Color("3"), lipgloss.NoColor{},
	}
	if (theme == "light" || theme == "dark") && r.ColorProfile() != termenv.ANSI {
		r.SetHasDarkBackground(theme == "dark")
		colors = []lipgloss.TerminalColor{
			lipgloss.AdaptiveColor{Light: "#245ea8", Dark: "#73b7ff"},
			lipgloss.AdaptiveColor{Light: "#7846a2", Dark: "#c4a7ff"},
			lipgloss.AdaptiveColor{Light: "#187649", Dark: "#78dba9"},
			lipgloss.AdaptiveColor{Light: "#b12e45", Dark: "#ff8f9a"},
			lipgloss.AdaptiveColor{Light: "#96630b", Dark: "#efc66a"},
			lipgloss.AdaptiveColor{Light: "#566378", Dark: "#b1bccb"},
		}
	}
	return display{
		width:   max(1, width),
		primary: r.NewStyle().Foreground(colors[0]), secondary: r.NewStyle().Foreground(colors[1]),
		success: r.NewStyle().Foreground(colors[2]), failure: r.NewStyle().Foreground(colors[3]),
		warning: r.NewStyle().Foreground(colors[4]), neutral: r.NewStyle().Foreground(colors[5]),
	}
}

func terminalTheme() string {
	if theme := os.Getenv("CHECKSY_THEME"); theme == "light" || theme == "dark" {
		return theme
	}
	parts := strings.Split(os.Getenv("COLORFGBG"), ";")
	if len(parts) > 1 {
		if background, err := strconv.Atoi(parts[len(parts)-1]); err == nil && background >= 0 && background <= 15 {
			_, _, lightness := termenv.ConvertToRGB(termenv.ANSIColor(background)).Hsl()
			if lightness >= 0.5 {
				return "light"
			}
			return "dark"
		}
	}
	return ""
}

func (d display) wrap(text string) string {
	return ansi.Wrap(text, d.width, "")
}

func (d display) indent(text string) string {
	if d.width <= 2 {
		return d.wrap(text)
	}
	return "  " + strings.ReplaceAll(ansi.Wrap(text, d.width-2, ""), "\n", "\n  ")
}

// Help writes grouped usage with the same labels and palette as the report.
func Help(out io.Writer) error {
	d := newDisplay(out)
	var lines []string
	for _, line := range strings.Split(strings.TrimRight(args.Usage, "\n"), "\n") {
		switch {
		case strings.HasPrefix(line, "Usage:"):
			line = d.secondary.Bold(true).Render("Usage:") + " " + d.primary.Bold(true).Render("checksy") + " [options]"
		case strings.HasPrefix(line, "  --"):
			flag, description, ok := strings.Cut(strings.TrimSpace(line), "  ")
			if ok {
				if d.width >= 50 {
					for _, row := range columns([]string{d.primary.Render(flag), strings.TrimSpace(description)}, []int{16, d.width - 20}) {
						lines = append(lines, "  "+row)
					}
				} else {
					lines = append(lines, d.indent(d.primary.Render(flag)), d.indent(strings.TrimSpace(description)))
				}
				continue
			}
			lines = append(lines, d.indent(line))
			continue
		case line == "Report" || line == "Automation" || line == "Information":
			line = d.secondary.Bold(true).Render(line)
		}
		lines = append(lines, d.wrap(line))
	}
	_, err := fmt.Fprintln(out, strings.Join(lines, "\n"))
	return err
}

func Version(out io.Writer, version string) error {
	d := newDisplay(out)
	_, err := fmt.Fprintln(out, d.wrap(d.primary.Bold(true).Render("checksy")+" "+d.secondary.Render(version)))
	return err
}

// Error writes an argument error and a short recovery hint to the supplied error stream.
func Error(out io.Writer, err error) error {
	d := newDisplay(out)
	message := "Invalid arguments"
	if err != nil {
		message = err.Error()
	}
	hint := "Run checksy --help for available options."
	if strings.Contains(message, "--timeout") {
		hint = "Try checksy --timeout 2s."
	}
	_, writeErr := fmt.Fprintln(out, d.wrap(d.primary.Bold(true).Render("checksy")+" • "+d.failure.Bold(true).Render("error"))+"\n"+d.wrap(message)+"\n"+d.wrap(d.warning.Render(hint)))
	return writeErr
}
