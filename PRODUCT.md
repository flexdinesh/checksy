# checksy

<!-- impeccable:product-schema 1 -->

## Platform

cli

Terminal CLI only. Existing releases target macOS and Linux. No web or native mobile interface planned.

## Users

Computer users who need a quick check or peek at their network. Do not assume network expertise.

## Product Purpose

Run `checksy` to see whether the internet is working and where connectivity looks broken. Success means a quick, understandable verdict with useful diagnostic context.

## Positioning

A small Go CLI combining a fixed set of connectivity checks in one compact terminal report, with a silent exit-code mode for scripts.

## Operating Context

- Run once from a shell; output remains in scrollback and the command exits.
- Read the internet verdict, then individual checks and discovered network facts.
- Use `--verbose` for measurement methods, full errors, and raw trace details.
- Use `--exit-code` for silent script execution; `--timeout` controls per-check timeout.

## Capabilities and Constraints

- Preserve fixed concurrent checks: HTTP connectivity, ping-style reachability to `1.1.1.1` and `8.8.8.8`, and system DNS resolution.
- HTTP connectivity determines the internet verdict. Ping and DNS explain failures without changing it.
- Run unprivileged: attempt ICMP, fall back to TCP on permission errors, and show the method used.
- Preserve the non-interactive, one-shot terminal report; no alternate screen or quit step.
- Discover public IP, local IP, default gateway, and resolver when available. Omit unavailable facts.
- Preserve silent mode: exit `0` when up, `1` when down, `2` for argument errors. Normal report mode does not use the verdict as its exit code.

## Brand Commitments

- Name: `checksy`.
- Use **check**, **target**, **result**, and **terminal report**, following `CONTEXT.md`.
- Avoid **probe**, **test**, **scan**, **endpoint**, **host**, **destination**, **outcome**, **response**, **TUI**, and **dashboard** for those product concepts.

## Evidence on Hand

- `README.md`: usage, installation, checks, and verdict semantics.
- `CONTEXT.md`: product terminology.
- `docs/adr/0001-icmp-with-tcp-fallback.md`: unprivileged operation and honest measurement labels.
- `docs/adr/0002-one-shot-terminal-report.md`: terminal workflow and network facts.
- `internal/report/`, `internal/ui/`, `internal/check/`, and `cmd/checksy/`: implemented behavior and tests.

## Product Principles

- Make a quick network check easy for ordinary computer users.
- Keep the verdict clear and diagnostic details honest.
- Require no elevated privileges or interaction after checks finish.
- Preserve shell-friendly output and script behavior.
