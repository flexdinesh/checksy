# checksy

checksy is a small Go CLI that runs a fixed set of connectivity checks against public targets and shows a compact verdict table so you can see whether the internet is working, and where it looks broken.

## Install

Stable Homebrew install:

```bash
brew install --cask flexdinesh/tap/checksy
```

Alternative stable install with Go:

`@latest` resolves to the newest stable SemVer tag, such as `v0.1.0`. There is no moving `latest` Git tag.

```bash
# Install the latest stable release.
go install github.com/flexdinesh/checksy/cmd/checksy@latest

# Install a specific stable release.
go install github.com/flexdinesh/checksy/cmd/checksy@v0.1.0
```

## Usage

```bash
# Check internet connectivity and show the verdict table.
checksy

# Run silently for scripts and exit 0 when internet is up, 1 when down.
checksy --exit-code

# Use a custom per-check timeout.
checksy --timeout 2s

# Show measurement methods, full failure text, and raw trace details.
checksy --verbose

# Show help.
checksy --help

# Show the installed version.
checksy --version
```

## Checks

checksy runs these checks concurrently:

- HTTP reachability to `https://connectivitycheck.gstatic.com/generate_204`, expecting `204`.
- Ping-style reachability to `1.1.1.1` and `8.8.8.8`, using ICMP when available and TCP-connect fallback when unprivileged.
- DNS resolution of `one.one.one.one` through the system resolver.
- Public egress IP discovery through Cloudflare's `/cdn-cgi/trace`.

The internet verdict is driven by the HTTP connectivity check. Ping and DNS rows are diagnostics; they explain what looks broken, but they do not flip the exit code.

## Development

See [docs/development.md](docs/development.md).

The ICMP fallback decision is documented in [docs/adr/0001-icmp-with-tcp-fallback.md](docs/adr/0001-icmp-with-tcp-fallback.md).

## Releases

Releases are created from `main`. See [docs/release.md](docs/release.md).
