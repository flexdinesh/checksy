# PRD: checksy v0.1

## Problem Statement

As a developer or operator, when my network feels off — pages won't load, a deploy hangs, a DNS lookup times out — I have no fast, trustworthy way to answer the single question that matters: *is the internet actually working, and if not, where is it broken?* I end up juggling `ping`, `dig`, and `curl` in three terminals, each with different output shapes, none of them summarizing the situation. And I can't put that answer into a script or a shell prompt that fails loudly when connectivity drops.

## Solution

`checksy` is a small Go CLI that runs a fixed set of connectivity checks against well-known public targets and reports the outcome in a compact Bubble Tea TUI — or, for scripting, exits with a deterministic pass/fail code. Running `checksy` prints a formatted verdict header (internet UP/DOWN, your public egress IP, your system resolver) plus a detail table of every check with its status and latency. Running `checksy --exit-code` stays silent and returns an exit code: `0` when the internet verdict is UP, `1` when DOWN, `2` on usage error.

## User Stories

1. As a developer, I want to run `checksy` with no arguments and see whether the internet is working, so that I don't have to remember a pile of network commands.
2. As a developer, I want a single verdict line (internet UP or DOWN) at the top of the output, so that I get the answer in under a second of reading.
3. As a developer, I want the verdict to reflect whether I can actually reach a real web endpoint, so that "UP" means I can use the internet, not just that I can ping an IP.
4. As an operator, I want to run `checksy --exit-code` in a script or shell prompt and get a pass/fail exit code, so that automation can react when connectivity drops.
5. As an operator, I want `--exit-code` to print nothing, so that it composes cleanly in scripts without polluting output.
6. As an operator, I want a DOWN verdict to exit `1` and a bad flag to exit `2`, so that I can distinguish "internet down" from "I misused the tool."
7. As a developer, I want to see the latency of each check, so that I can tell a flaky/slow link from a hard failure.
8. As a developer, I want to see ping latency to well-known public IPs (Cloudflare and Google), so that I can judge raw reachability.
9. As a developer, I want to see DNS resolution working, so that I can tell whether name resolution is the broken layer.
10. As a developer, I want to see the resolved IP for a DNS check, so that I can confirm the answer is sane.
11. As a developer, I want to see HTTP reachability to a purpose-built connectivity endpoint, so that the verdict reflects a real end-to-end web request.
12. As a developer, I want the HTTP verdict to expect a specific success response (HTTP 204 from a captive-check endpoint), so that captive portals are detected rather than mistaken for working internet.
13. As a developer, I want to see my public egress IP in the output, so that I can confirm which network/exit I'm actually on.
14. As a developer, I want to see the system resolver address, so that I know which DNS server my machine is configured to use.
15. As a developer, I want the "ping" check to show whether it used real ICMP or fell back to TCP, so that the latency number is honest about how it was measured.
16. As a developer, I want a hung or blackholed network to terminate the tool within a few seconds, so that `checksy` never hangs forever.
17. As a developer, I want to configure the per-check timeout via a flag, so that I can tune for slow links or demand fast failure.
18. As a developer, I want failures to show a short reason inline, so that I can start diagnosing without re-running anything.
19. As a developer, I want a `--verbose` mode that adds extra detail (measurement method, full error text, raw trace body), so that I can dig deeper when the default output isn't enough.
20. As a developer, I want a `--help` message listing every flag, so that I can discover the tool's options.
21. As a developer, I want a `--version` flag, so that I can tell which build of checksy is installed.
22. As a developer, I want the output to be well-spaced and color-coded (green for ok, red for fail), so that the state is readable at a glance.
23. As a developer, I want the TUI to run all checks concurrently, so that the tool finishes in the time of the slowest check, not the sum of all checks.
24. As a developer, I want to install checksy via `go install`, so that I can get it on my PATH the same way I install other Go CLIs.
25. As a developer, I want checksy to run unprivileged (no `sudo`), so that the "just run `checksy`" workflow is frictionless.

## Implementation Decisions

- **Module & layout.** New Go module `github.com/flexdinesh/checksy`, Go toolchain matching gitsy. Standard CLI shape mirroring the `gitsy` sibling repo: entrypoint in `cmd/checksy`, behavior split into focused `internal/` packages — argument parsing, the check domain (types + per-protocol runners), display row/verdict shaping, and the Bubble Tea program. Version stamping package for release builds.
- **Domain model.** Three nouns: a **Check** (a single connectivity measurement against a target via a protocol — ICMP, DNS, or HTTP); a **Target** (the network endpoint — an IP or host); a **Result** (the outcome: status ok/fail, latency, and discovered detail). Mirrors gitsy's Repo→Result→Row shape.
- **Check set.** A fixed, non-redundant set covering three layers: **Ping** (ICMP with TCP-connect fallback) to `1.1.1.1` and `8.8.8.8`; **DNS** resolution of `one.one.one.one` via the system resolver; **HTTP** to `https://connectivitycheck.gstatic.com/generate_204` expecting status 204 (the verdict check), plus `https://www.cloudflare.com/cdn-cgi/trace` to discover the public egress IP. Standalone TCP and UDP checks are intentionally excluded — TCP is subsumed by the ping fallback, UDP by DNS.
- **ICMP fallback.** Real ICMP echo needs raw sockets (root on macOS). checksy attempts ICMP and falls back to measuring TCP-connect RTT to a well-known port on permission failure. Each ping Result records and displays the method (`icmp` or `tcp`) so output stays honest. See ADR 0001.
- **Verdict & exit code.** The internet verdict is driven solely by HTTP: at least one HTTP target returning the expected 204 status ⇒ UP (exit 0); all HTTP targets failing ⇒ DOWN (exit 1). Ping and DNS results are diagnostics that explain *why* on failure but never flip the exit code. Usage/argument errors exit 2. A captive portal returning a non-204 body to the captive-check endpoint correctly yields DOWN.
- **Run model.** One-shot: run all checks once, render, wait for quit. No live-refresh in v1 (a `--watch` mode is a future addition that the Result model already supports without rework). `--exit-code` is the same run, minus the TUI, silent on stdout/stderr.
- **Concurrency.** All checks run concurrently via the Bubble Tea command-batching pattern inherited from gitsy; a concurrency cap is retained for consistency though the check count is small. Total runtime is bounded by the slowest check, not the sum.
- **Timeouts.** Every check is bounded by a per-check timeout (default 5s), configurable via `--timeout`. A timed-out check is a failure with a `timeout after Xs` detail. No separate overall deadline is needed because checks run concurrently.
- **Discovered facts.** The public egress IP is parsed from the `/cdn-cgi/trace` response body. The system resolver address is read from the host resolver configuration on Unix. These appear in the verdict header, not as table rows. Resolver reading is Unix-only in v1.
- **Output shape.** A verdict header (internet UP/DOWN, public IP, resolver) followed by a detail table with one row per check: target, check kind, status, latency, and a polymorphic detail column (method for ping, resolved IP for DNS, status code for HTTP, marker for the trace row). Latency is shown in milliseconds at one decimal place. Status is binary ok/fail; there is no "slow" state in v1.
- **Display library.** Bubble Tea (program/model/update/view), Bubbles (table, spinner), Lip Gloss (frames, tones), and a runewidth helper for terminal-width-safe layout — the same dependency set as gitsy. Tone styling (green ok, red fail, dim diagnostics) is reused.
- **Flags.** `--exit-code` (silent verdict mode), `--timeout <duration>` (per-check timeout), `--verbose` (extra detail), `--help`, `--version`. Hand-rolled parser matching gitsy's style, supporting both `--flag value` and `--flag=value`.
- **Release.** GoReleaser config mirroring gitsy (darwin/linux, amd64/arm64, version stamping via ldflags), GitHub Actions release workflow, first release tagged `v0.1.0`.

## Testing Decisions

- **Only external behavior is tested.** No tests assert private helpers or internal call graphs. Stable, user-facing wording is treated as part of the contract (as in gitsy), so formatter/parser tests pin exact strings.
- **Seam 1 — CLI / `run(argv, executor)`.** The highest seam. The check-executor is injected so the exit-code verdict logic (HTTP 204 ⇒ 0, else 1; usage error ⇒ 2) is exercised deterministically with canned Results, with no real network. Prior art: `gitsy/cmd/gitsy/main_test.go`.
- **Seam 2 — `internal/args`.** Pure table-driven parser tests covering every flag and error path. Prior art: `gitsy/internal/args/args_test.go`.
- **Seam 3 — `internal/ui`.** Pure table-driven formatting tests: verdict string, row tones, latency formatting, polymorphic detail. Prior art: `gitsy/internal/ui/ui_test.go`.
- **Seam 4 — `internal/check` runners.** Runners accept a target, so tests point them at local fixtures: an `httptest.Server` for HTTP and the trace-body parser, a local TCP listener for ping-fallback RTT, and local DNS where feasible. Response-parsing helpers (trace body ⇒ public IP, status ⇒ ok/fail) are pure unit tests. **No assertions against the real public internet in unit tests** — that path is verified manually, since it is inherently flaky.
- **Seam 5 — `internal/tui`.** The Bubble Tea model is fed synthetic Results via messages and its `View()` is asserted, fully offline. Prior art: `gitsy/internal/tui/tui_test.go` via `NewModel`.

## Out of Scope

- Live-refresh / `--watch` mode.
- A standalone TCP port-reachability check (subsumed by the ping TCP fallback).
- A standalone UDP check (subsumed by DNS).
- A direct DNS query to `1.1.1.1:53` ("deep" resolver check) — deferred to post-v1.
- A dedicated plain-HTTP captive-portal endpoint (e.g. `captive.apple.com`) — the 204-expectation on the existing HTTP target already catches most captive portals; deferred.
- Windows support for reading the system resolver config (Unix-only in v1).
- Configurable check targets / custom check sets (the v1 set is fixed and compiled in).
- A "slow link" warning state / latency-threshold coloring beyond plain ok/fail.
- JSON/machine-readable output format.

## Further Notes

- Vocabulary is defined in `CONTEXT.md` (Check / Target / Result); the ICMP-fallback decision is recorded in `docs/adr/0001-icmp-with-tcp-fallback.md`.
- The sibling repo `gitsy` (github.com/flexdinesh/gitsy) is the canonical reference for CLI shape, package conventions, Bubble Tea patterns, tone styling, GoReleaser config, and release workflow. checksy intentionally mirrors it.
- The verdict semantics — "ping can be red while the exit code is 0" — is deliberate and surprising; it is documented inline at the decision site in code rather than via a separate ADR (it is a one-line behavioral change to reverse).
