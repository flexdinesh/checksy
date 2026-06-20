# 03 — DNS check + resolved-IP detail

## Parent

`.scratch/prd.md` (checksy v0.1)

## What to build

Add the **DNS** check, which resolves `one.one.one.one` via the system resolver and reports resolution latency. The detail column shows the resolved IP. DNS is a diagnostic — it never flips the exit code. A failed resolution shows as failed with a short reason.

## Acceptance criteria

- [ ] `checksy` shows a DNS row for resolving `one.one.one.one` with a latency
- [ ] The DNS detail column shows the resolved IP
- [ ] A failed resolution shows as failed with a short reason
- [ ] Response-handling helpers are pure unit tests (no real network)

## Blocked by

- 01 (walking skeleton)
