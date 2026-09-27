---
version: 1
slug: "internal-report-report-go"
primary_target: "internal/report/report.go"
related_targets: ["internal/report/terminal.go","internal/ui/ui.go","cmd/checksy/main.go","internal/args/args.go"]
---

# Terminal output

Mode: Operate. Audience: computer users making a quick network check.
Scope: report, help, version, and errors. Code-first; no raster assets or motion.

## Direction contract

THESIS: Network receipt: one verdict, a compact check ledger, and closing network context. Replace the enclosing grid with open aligned rows.

OWN-WORLD: Blue and violet labels, green success, red failure, amber recovery hints. Inherit terminal background and font; semantic glyphs and words survive monochrome.

STORY: Read the HTTP-based verdict, inspect four checks, then locate public/local IP, gateway, and resolver. Help explains existing flags; argument errors name the correction.

FIRST VIEWPORT: Brand and verdict share one line above one thin divider, a quiet column header, four open check rows, and a compact two-row facts footer. Narrow terminals wrap cells and stack facts. Signature behavior: verbose details unfold after the receipt without shifting its summary; no animation or interaction after completion.

FORM: Utility inspection card, grounded candidate 7; seed 6a702473; user selected Network receipt. Preserve all check, timeout, privilege, and exit-code semantics.

FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance

## Verification

Capture light/dark at 80 columns plus narrow 40-column and monochrome output. Cover UP, DOWN, mixed diagnostic failure, unavailable facts, IPv6, long errors, verbose trace, help, version, and argument errors. No HTML/CSS detector: terminal Go output is outside its supported surface.
