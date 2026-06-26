# Releases

Releases are SemVer Git tags on `main`.

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

The workflow creates the tag, runs GoReleaser, publishes macOS and Linux
archives plus checksums, and opens or updates a pull request against
`flexdinesh/homebrew-tap`.

Do not create a moving `latest` tag. Go already resolves `@latest` to the newest
SemVer tag.

Pushes to `dev` still run automatic GoReleaser snapshot releases through CI.

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

## Version Output

Local builds print a development version. Release builds get the version from
the release tag through GoReleaser linker flags.

```bash
checksy --version
```

## Required Secret

The release workflow requires:

- `HOMEBREW_TAP_TOKEN`: a fine-grained GitHub token with contents write and pull request write access to `flexdinesh/homebrew-tap`.

The workflow also uses the built-in `GITHUB_TOKEN` to create tags and publish
the GitHub Release in this repository.

## Homebrew

The Homebrew formula installs prebuilt release archives instead of building from
source. `checksy` does not declare Homebrew runtime dependencies because the
released binary contains the connectivity-checking implementation.

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
5. Merge the generated `flexdinesh/homebrew-tap` pull request after tap CI passes.
6. Verify with `brew install flexdinesh/tap/checksy` and `checksy --version`.

## Verify Locally

```bash
go test ./...
go build ./cmd/checksy
goreleaser release --snapshot --clean
```

The workflows pin GoReleaser `v2.9.0` because GoReleaser deprecated formula
publishing through `brews` in later versions. The snapshot command remains useful
locally with newer GoReleaser versions because it verifies archive and formula
generation without publishing.

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
