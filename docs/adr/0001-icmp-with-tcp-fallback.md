# ICMP ping with TCP-connect fallback

Real ICMP echo requires raw sockets, which need root on macOS (and a sysctl tweak on Linux). Since checksy must run unprivileged as `checksy`, raw ICMP would silently fail for ~every user.

We attempt ICMP via `golang.org/x/net/icmp`; on permission error we fall back to measuring a TCP `Dial` RTT to port 443 (or 53). The `Result.Method` field records which path ran (`icmp` or `tcp`) and is shown in the UI so the output stays honest. On macOS the TCP fallback is the common path — that's expected, not a bug.

## Considered options

- **Skip ICMP, measure TCP-connect latency only.** Rejected: loses the "ping" intent and the ICMP path on platforms where it does work unprivileged.
- **Require `sudo` for real ping.** Rejected: destroys the "just run `checksy`" ergonomics.
- **Attempt ICMP + TCP fallback (chosen).** Matches intent, cross-platform, honest via the method label.
