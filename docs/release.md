# Releases

Stable releases are SemVer Git tags on `main`. Development installs use the
automatically updated `dev` branch.

## Current Policy

Stable releases are created manually from the latest code on `main` by running
the GitHub Actions release workflow. Each dispatch creates the next `v0.1.x`
release. If there are no `v0.1.x` release tags yet, the first dispatch creates
`v0.1.0`.

Examples:

```bash
v0.1.0
v0.1.1
v0.1.2
```

The workflow checks out the latest `main`, tests it, creates or reuses the tag,
runs GoReleaser, publishes macOS and Linux archives plus checksums as the latest
stable GitHub Release, and opens or updates a pull request against
`flexdinesh/homebrew-tap`.

Do not create a moving `latest` tag. Go already resolves `@latest` to the newest
SemVer tag.

Every push to `main` runs tests and builds the CLI, then advances the `dev`
branch to that tested commit. Superseded runs skip publishing; updates are
serialized and fast-forward only so an older run cannot roll `dev` back.
Failed CI leaves `dev` at its previous commit. Treat `dev` as an automated
installation channel; merge changes into `main` instead of pushing to `dev`.

Go resolves `@dev` from that branch, usually as a pseudo-version. This creates
no SemVer tags or GitHub Releases and leaves `@latest` on the newest stable
version. The old workflow triggered by pushes to `dev` has been removed.

## Installing

```bash
# Stable Homebrew install.
brew install flexdinesh/tap/checksy

# Alternative latest Go release.
go install github.com/flexdinesh/checksy/cmd/checksy@latest

# Specific release.
go install github.com/flexdinesh/checksy/cmd/checksy@v0.1.0

# Development release.
go install github.com/flexdinesh/checksy/cmd/checksy@dev
```

Go module proxies may briefly cache the previous `dev` revision. To fetch
directly from GitHub immediately after CI publishes:

```bash
GOPROXY=direct go install github.com/flexdinesh/checksy/cmd/checksy@dev
```

## Version Output

Local builds print a development version. Release archives get the version from
the release tag through GoReleaser linker flags. Go installs report the module
version: a stable tag, or usually a pseudo-version for `@dev`. If `dev` points
at a tagged stable commit, Go may report that stable version.

```bash
checksy --version
```

## Required Secret

The release workflow requires:

- `HOMEBREW_TAP_TOKEN`: a fine-grained GitHub token with contents write and pull request write access to `flexdinesh/homebrew-tap`.

The workflow also uses the built-in `GITHUB_TOKEN` to create tags and publish
the GitHub Release in this repository. The CI `Publish dev` job uses only the
built-in `GITHUB_TOKEN` with contents write permission to update `dev`.

## Homebrew

The Homebrew formula builds the tagged source archive locally. Go is a build
dependency, with no Homebrew runtime dependencies. This avoids the Gatekeeper
warnings caused by quarantined, unsigned prebuilt cask binaries on macOS.

After GoReleaser publishes the release, the workflow downloads its source archive
and runs `tools/homebrew-formula` to generate `Formula/checksy.rb` with that
archive's SHA-256. It opens a tap pull request containing the formula and removes
the old `Casks/checksy.rb`. GoReleaser publishes binary archives independently;
it no longer publishes a cask or uses its deprecated formula publisher.

Existing cask users must migrate once after the tap pull request merges:

```bash
brew uninstall --cask checksy
brew update
brew install flexdinesh/tap/checksy
```

The tap pull request branch is deterministic per version, such as
`checksy-v0.1.0`, so rerunning a failed release updates the same tap pull
request.

The release also replaces existing GitHub Release artifacts on rerun. If a
workflow publishes the GitHub Release but fails while opening the Homebrew tap
pull request, rerunning the same workflow on the same commit should reuse the tag,
refresh the release artifacts, and retry the tap pull request.

## Release Steps

1. Merge the release-ready code to `main`.
2. Run the **Release** workflow from GitHub Actions.
3. Confirm the workflow created or reused the expected `v0.1.x` tag.
4. Review the generated GitHub Release artifacts and checksums.
   Confirm it is marked **Latest**, and verify both Go install forms above.
5. Merge the generated `flexdinesh/homebrew-tap` pull request after tap CI passes.
6. Verify with `brew install flexdinesh/tap/checksy` and `checksy --version`.

## Verify Locally

```bash
go test ./...
go build ./cmd/checksy
goreleaser release --snapshot --clean
```

The workflow pins GoReleaser `v2.18.0`. The snapshot command verifies binary
archive generation without publishing. Verify formula generation separately:

```bash
curl --fail --location --retry 3 \
  https://github.com/flexdinesh/checksy/archive/refs/tags/v0.1.5.tar.gz \
  --output /tmp/checksy-source.tar.gz
go run ./tools/homebrew-formula -tag v0.1.5 \
  -archive /tmp/checksy-source.tar.gz > /tmp/checksy.rb
ruby -c /tmp/checksy.rb
```

The tap's macOS CI runs Homebrew style, audit, install, and formula tests before
the generated update is merged.

## Switching Minor Versions

Switch manually when `0.1.x` no longer feels right, for example when a release
is the first meaningful preview rather than just the next small change.

To switch, update `.github/workflows/release.yml` so the tag selector uses the
new minor line, such as `v0.2.*`, and starts at `v0.2.0`.

After that, releases should continue as:

```bash
v0.2.0
v0.2.1
v0.2.2
```
