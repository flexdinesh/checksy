# 05 — `--timeout` flag + timeout-failure path

## Parent

`.scratch/prd.md` (checksy v0.1)

## What to build

Promote the skeleton's default per-check timeout into a configurable `--timeout <duration>` flag. Every check runner honors a context deadline derived from it. A check that exceeds the timeout is marked failed with detail `timeout after Xs`. Default remains `5s` when the flag is absent; invalid duration values are an argument error (exit `2`).

## Acceptance criteria

- [ ] `checksy --timeout 2s` bounds every check to 2s
- [ ] A check exceeding the timeout is marked failed with detail `timeout after Xs`
- [ ] Default timeout is `5s` when the flag is absent
- [ ] An invalid duration value exits `2`
- [ ] The timeout-failure path is tested via local fixtures (a deliberately slow local server/listener)

## Blocked by

- 01 (walking skeleton)
