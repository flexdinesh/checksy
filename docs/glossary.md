# Glossary

## Homebrew Tap

A Git repository that Homebrew can read formulae from. `checksy` uses the
`flexdinesh/homebrew-tap` tap for stable Homebrew installs.

## Stable Release

A SemVer Git tag on `main`, such as `v0.1.0`, that GoReleaser turns into GitHub
Release artifacts and a Homebrew formula update.

## GoReleaser Snapshot

A local or CI release dry run that builds archives and generated metadata without
publishing GitHub Release artifacts or updating the Homebrew tap.
