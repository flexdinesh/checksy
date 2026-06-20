# One-shot terminal report instead of interactive TUI

CheckSy runs a fixed set of connectivity checks and exits. There is no live state to browse, filter, refresh, or interact with after the checks complete.

We render a pretty one-shot terminal report to stdout instead of launching an interactive Bubble Tea TUI or using the terminal alternate screen. The report keeps the compact verdict table, discovered facts, and verbose details, but it returns immediately after writing output.

## Considered options

- **Interactive TUI in the alternate screen.** Rejected: it makes users quit the interface after a fixed check run, hides output from shell scrollback, and adds interactivity without a current workflow.
- **Plain unstyled text table.** Rejected: CheckSy should still be easy to scan at a glance.
- **Pretty one-shot terminal report (chosen).** Matches the "run and show me the internet verdict" workflow while preserving polished terminal output and shell-friendly behavior.
