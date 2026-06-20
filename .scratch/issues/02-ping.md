# 02 — Ping check (ICMP + TCP fallback) + method detail

## Parent

`.scratch/prd.md` (checksy v0.1)

## What to build

Add the **Ping** check for targets `1.1.1.1` and `8.8.8.8`. Real ICMP echo needs raw sockets, which require privileges the tool will not assume, so the runner attempts ICMP and falls back to measuring a TCP-connect RTT to a well-known port on permission failure (see ADR `docs/adr/0001-icmp-with-tcp-fallback.md`). Each ping `Result` carries a `Method` (`icmp` or `tcp`) shown in the detail column so the latency number is honest about how it was measured. Ping is a diagnostic — it never flips the exit code.

## Acceptance criteria

- [ ] `checksy` shows ping rows for `1.1.1.1` and `8.8.8.8`, each with a latency
- [ ] Each ping row shows the method used (`icmp` or `tcp`) in the detail column
- [ ] When ICMP is unavailable (e.g. unprivileged on macOS), the check falls back to TCP-connect RTT and still reports a latency with method `tcp`
- [ ] A timed-out or unreachable ping target shows as failed with a short reason
- [ ] The TCP-fallback path is tested against a local TCP listener; method selection is tested for both the ICMP-success and fallback paths

## Blocked by

- 01 (walking skeleton)
