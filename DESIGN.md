---
name: checksy
description: A compact terminal network receipt.
colors:
  primary-light: "#245ea8"
  primary-dark: "#73b7ff"
  secondary-light: "#7846a2"
  secondary-dark: "#c4a7ff"
  success-light: "#187649"
  success-dark: "#78dba9"
  failure-light: "#b12e45"
  failure-dark: "#ff8f9a"
  warning-light: "#96630b"
  warning-dark: "#efc66a"
  neutral-light: "#566378"
  neutral-dark: "#b1bccb"
components:
  primary-label-light:
    textColor: "{colors.primary-light}"
  primary-label-dark:
    textColor: "{colors.primary-dark}"
  secondary-label-light:
    textColor: "{colors.secondary-light}"
  secondary-label-dark:
    textColor: "{colors.secondary-dark}"
  status-ok-light:
    textColor: "{colors.success-light}"
  status-ok-dark:
    textColor: "{colors.success-dark}"
  status-fail-light:
    textColor: "{colors.failure-light}"
  status-fail-dark:
    textColor: "{colors.failure-dark}"
  recovery-light:
    textColor: "{colors.warning-light}"
  recovery-dark:
    textColor: "{colors.warning-dark}"
  quiet-label-light:
    textColor: "{colors.neutral-light}"
  quiet-label-dark:
    textColor: "{colors.neutral-dark}"
---

# Design System: checksy

## Overview

**Creative North Star: "Network receipt"**

A compact, calm terminal report. One verdict leads an open check ledger and closing network facts. Color clarifies labels and results; words and semantic glyphs retain meaning without color.

The terminal supplies the background, font, and text metrics. Output remains in scrollback after the command exits. Help, version, and argument errors use the same visual language.

**Key Characteristics:**
- Open aligned rows, one thin divider.
- Blue and violet labels; green, red, and amber meanings.
- Terminal-native type, background, and monochrome behavior.
- Width-aware wrapping; verbose details below the summary.

## Colors

Paired foreground shades support known light and dark themes; an unknown theme uses the terminal's own ANSI palette.

### Primary
- **Clear Blue** (`primary-light`, `primary-dark`): brand, network fact labels, flags, and verbose target labels.

### Secondary
- **Diagnostic Violet** (`secondary-light`, `secondary-dark`): check kinds, section headings, usage label, and version.

### Tertiary
- **Positive Green** (`success-light`, `success-dark`): `UP` and `✓ ok`.
- **Failure Red** (`failure-light`, `failure-dark`): `DOWN`, `✗ fail`, and the argument error label.
- **Recovery Amber** (`warning-light`, `warning-dark`): `• n/a`, no-checks message, and recovery hints.

### Neutral
- **Quiet Slate** (`neutral-light`, `neutral-dark`): divider and quiet column headings. Ordinary values and body text inherit the terminal foreground.

**The Terminal Ownership Rule.** Inherit the terminal background and ordinary foreground; never paint an enclosing surface.

`CHECKSY_THEME=light|dark` takes precedence over a recognized `COLORFGBG` background. Known themes select the paired tokens when color capabilities support them; ANSI-only output retains slots 4, 5, 2, 1, and 3, with default neutral. No terminal theme queries. Non-terminal writers, `NO_COLOR`, and `TERM=dumb` disable color.

## Typography

Use the terminal's existing font and cell metrics. There is no app font family, size ramp, line-height, or tracking token. Bold distinguishes the brand, verdict, section headings, usage label, and error label. Check kinds and ledger headings are uppercase; body values remain ordinary text.

**The Words Survive Rule.** Preserve status words and glyphs together so monochrome output remains understandable.

## Layout

Measure visible terminal columns, excluding ANSI sequences. Width comes from the terminal, then `COLUMNS`, with an 80-column fallback. Wrap to the available width.

The report starts with brand and verdict, then a thin divider and open rows. Wide rows have status, check, target, latency, and detail columns, separated by two spaces. Status and check widths are 6 and 5 columns. Target width is content-driven, capped at one third of the terminal; latency has a minimum width of 7 and aligns right. Cells wrap independently without losing their column positions.

When fewer than 10 columns remain for detail, each check becomes a header plus a two-space-indented latency/detail line. Available facts follow after one blank line; pair adjacent facts with a three-space minimum gap when they fit, otherwise stack and wrap. Unavailable facts are omitted.

Help uses a 16-column flag field at widths of at least 50; narrower output stacks flags and descriptions. Verbose Details and Trace follow the receipt, separated by blank lines.

## Elevation & Depth

Flat terminal text. No shadows, tonal surfaces, motion, or overlays. Order, alignment, blank lines, and bold establish hierarchy.

## Shapes

Open text rows and a single `─` divider, sized to the widest summary line and capped by terminal width. No enclosing box, rounded container, or raster decoration. The native `✓`, `✗`, `•`, and `—` communicate status, separation, and unavailable latency.

## Components

- **Verdict line:** bold blue brand, ordinary `• internet`, bold green `UP` or red `DOWN`. HTTP determines the verdict.
- **Check ledger:** violet check kinds, colored status glyphs and words, ordinary targets/details, and right-aligned latency in wide output. Narrow rows unfold beneath their headers.
- **Network facts:** blue labels with ordinary values; pair or stack to fit. Order: Public IP, Local IP, Gateway, Resolver.
- **Verbose appendix:** bold violet Details/Trace headings, blue target labels, wrapped full details and raw trace text. Summary stays first.
- **Help:** violet usage/group labels, blue brand/flags, ordinary descriptions, width-aware grouping.
- **Argument error:** blue brand and bold red error label, ordinary explanation, amber correction. Timeout errors suggest `checksy --timeout 2s`; others suggest help.
- **Version:** bold blue brand and violet version on one wrapped line.

## Do's and Don'ts

### Do:
- **Do** preserve terminal-owned background, font, and ordinary foreground.
- **Do** keep status words beside their semantic glyphs.
- **Do** wrap long values and errors within the available terminal width.
- **Do** use the same palette and hierarchy for report, help, version, and errors.

### Don't:
- **Don't** add enclosing grids, fixed backgrounds, app font metrics, motion, or raster assets.
- **Don't** rely on color alone to express a result.
- **Don't** move verbose details above the compact summary.
- **Don't** add interaction, an alternate screen, or a quit step.
