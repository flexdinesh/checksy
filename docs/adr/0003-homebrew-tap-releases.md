# Publish stable releases through the Homebrew tap

## Status

Accepted

## Context

`checksy` is a terminal CLI, and users should be able to install stable versions
without installing Go manually. The repository already uses Go-compatible SemVer release
tags and GoReleaser archives. The initial tap integration installed prebuilt
binaries through a cask, causing macOS Gatekeeper warnings because the binaries
were not Developer ID signed and notarized. Go installs worked because they
built the executable locally.

There is an existing `flexdinesh/homebrew-tap` repository with Homebrew CI. That
tap is the natural Homebrew distribution channel for personal CLI tools.

## Decision

Stable `checksy` releases use GoReleaser to build prebuilt macOS and Linux
archives for `amd64` and `arm64` and publish GitHub Release artifacts. The release
workflow separately generates `Formula/checksy.rb` in `flexdinesh/homebrew-tap`
from the tagged source archive and its SHA-256, removing `Casks/checksy.rb`.

The release workflow opens or updates a pull request against the tap instead of
pushing directly to tap `main`. The tap branch is deterministic per version, such
as `checksy-v0.1.0`, so rerunning a release updates the same pull request.

The formula builds locally with Go as a build dependency and has no Homebrew
runtime dependencies. This follows the working Go install approach without
requiring Apple signing and notarization credentials.

## Consequences

Users can install stable releases with `brew install flexdinesh/tap/checksy`.
Installations take longer because they compile source. Existing cask users must
uninstall the cask before installing the formula.

The source repository needs a `HOMEBREW_TAP_TOKEN` secret with write and pull
request access to `flexdinesh/homebrew-tap`.

The tap repository remains responsible for Homebrew-native style, audit, and
install and formula tests before an update is merged.

GoReleaser's formula publisher is deprecated. Formula generation and tap pull
requests therefore belong to the release workflow rather than GoReleaser.
