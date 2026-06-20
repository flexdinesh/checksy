# 06 — `--verbose` + failure-reason details

## Parent

`.scratch/prd.md` (checksy v0.1)

## What to build

Make failures informative and add a `--verbose` mode. By default, failed rows show a short reason in the detail column. `--verbose` adds: the measurement method on every row, the full error text on failures, and the raw `/cdn-cgi/trace` body. Default output is unchanged when verbose is off.

## Acceptance criteria

- [ ] By default, failed rows show a short reason in the detail column
- [ ] `checksy --verbose` shows the measurement method on every row, full error text on failures, and the raw trace body
- [ ] Default output is unchanged when `--verbose` is off
- [ ] Verbose-vs-default behavior is table-driven and pinned

## Blocked by

- 01 (walking skeleton)
