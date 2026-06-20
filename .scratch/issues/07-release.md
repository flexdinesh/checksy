# 07 — Release scaffolding: GoReleaser + GitHub Actions, v0.1.0

## Parent

`.scratch/prd.md` (checksy v0.1)

## What to build

Stand up release tooling mirroring the `gitsy` sibling repo: a GoReleaser config (darwin/linux, amd64/arm64, version stamping via ldflags) and a GitHub Actions release workflow triggered on tags. Cut the first release as `v0.1.0`, installable via `go install`.

## Acceptance criteria

- [ ] GoReleaser config builds darwin/linux, amd64+arm64, with version ldflags
- [ ] GitHub Actions release workflow mirrors gitsy
- [ ] First release tagged `v0.1.0` installs via `go install github.com/flexdinesh/checksy/cmd/checksy@v0.1.0`

## Blocked by

- 01 (walking skeleton)

## Type

HITL — cutting and verifying the release is a human act.
