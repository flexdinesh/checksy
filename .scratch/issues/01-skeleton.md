# 01 — Walking skeleton: HTTP verdict check, end-to-end

## Parent

`.scratch/prd.md` (checksy v0.1)

## What to build

The tracer bullet — a thin but complete path through every layer. Initialize the Go module and stand up the CLI shape: entrypoint, argument parsing, the `Check`/`Target`/`Result` domain types, the Bubble Tea program, and the display layer. Wire exactly one real check through it: the **HTTP verdict check** against `https://connectivitycheck.gstatic.com/generate_204`, which expects status 204.

The internet verdict is driven solely by this HTTP check: 204 ⇒ UP, otherwise DOWN. `checksy` renders a verdict header plus the one HTTP detail row; `checksy --exit-code` prints nothing and exits `0` (UP) or `1` (DOWN); argument errors exit `2`. The check executor is injected at the top-level `run()` seam so the verdict→exit-code logic is exercised deterministically with canned Results, independent of the real network. Ships with a default per-check timeout constant so a blackholed target can never hang the tool.

Mirror the `gitsy` sibling repo for module shape, package conventions, Bubble Tea model/update/view pattern, tone styling, and the hand-rolled argument parser.

## Acceptance criteria

- [ ] `checksy` runs and renders a verdict header (`internet UP` or `internet DOWN`) plus one detail row for the HTTP check to `connectivitycheck.gstatic.com/generate_204`
- [ ] Verdict is UP iff the HTTP check returns status 204; otherwise DOWN
- [ ] `checksy --exit-code` prints nothing and exits `0` when the HTTP check returns 204, `1` otherwise
- [ ] An argument error (e.g. unknown flag) exits `2`
- [ ] `checksy --help` prints usage listing every flag; `checksy --version` prints the version
- [ ] A blackholed/hung HTTP target terminates within the default per-check timeout (no infinite hang)
- [ ] The verdict→exit-code mapping is unit-tested via an injectable executor with canned Results, with no real network
- [ ] `go test ./...` passes; parser and formatter tests are table-driven and pin exact user-facing wording

## Blocked by

None — can start immediately.
