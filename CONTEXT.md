# checksy

A small Go CLI that verifies internet connectivity by running a set of connectivity checks against public targets and rendering a one-shot terminal report (or returning a pass/fail exit code).

## Language

**Check**:
A single connectivity measurement against a target using a specific protocol (ICMP, DNS, TCP, HTTP).
_Avoid_: probe, test, scan

**Target**:
The network endpoint a check runs against — an IP address or hostname.
_Avoid_: endpoint, host, destination

**Result**:
The outcome of running a check: status (ok/fail), latency, and any discovered detail.
_Avoid_: outcome, response

**Terminal report**:
The non-interactive, pretty terminal output written once after checks complete.
_Avoid_: TUI, dashboard, screen
