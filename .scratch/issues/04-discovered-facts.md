# 04 — Discovered facts: public IP + system resolver header

## Parent

`.scratch/prd.md` (checksy v0.1)

## What to build

Enrich the verdict header with two discovered facts: the **public egress IP**, parsed from the body of `https://www.cloudflare.com/cdn-cgi/trace`, and the **system resolver address**, read from the host resolver configuration on Unix. These render in the header, not as table rows. Both degrade gracefully — if either is unavailable, the header omits the field without failing the run.

## Acceptance criteria

- [ ] The verdict header shows the public egress IP parsed from the Cloudflare `/cdn-cgi/trace` body
- [ ] The verdict header shows the system resolver address on Unix
- [ ] If either fact is unavailable, the header omits that field without failing the run
- [ ] The trace-body → IP parse is unit-tested against fixture bodies

## Blocked by

- 01 (walking skeleton)
