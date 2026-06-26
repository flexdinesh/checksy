# Publish stable releases through the Homebrew tap

## Status

Accepted

## Context

`checksy` is a terminal CLI, and users should be able to install stable versions
without installing Go. The repository already uses Go-compatible SemVer release
tags and GoReleaser archives, but it did not update a Homebrew cask.

There is an existing `flexdinesh/homebrew-tap` repository with Homebrew CI. That
tap is the natural Homebrew distribution channel for personal CLI tools.

## Decision

Stable `checksy` releases use GoReleaser to build prebuilt macOS and Linux
archives for `amd64` and `arm64`, publish GitHub Release artifacts, and generate
`Casks/checksy.rb` in `flexdinesh/homebrew-tap`.

The release workflow opens or updates a pull request against the tap instead of
pushing directly to tap `main`. The tap branch is deterministic per version, such
as `checksy-v0.1.0`, so rerunning a release updates the same pull request.

The cask has no runtime dependencies because the release archive contains the
complete `checksy` binary.

## Consequences

Users can install stable releases with `brew install --cask flexdinesh/tap/checksy`.

The source repository needs a `HOMEBREW_TAP_TOKEN` secret with write and pull
request access to `flexdinesh/homebrew-tap`.

The tap repository remains responsible for Homebrew-native style, audit, and
install checks before a cask update is merged.

GoReleaser's formula publisher is deprecated in newer versions, so `checksy`
uses the supported cask publisher.
