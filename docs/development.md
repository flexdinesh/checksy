# Development

```bash
# Check out the repo.
git clone git@github.com:flexdinesh/checksy.git && cd checksy

# Run tests.
go test ./...

# Build the binary.
go build -o bin/checksy ./cmd/checksy

# Install the local build as a binary.
go install ./cmd/checksy
```

## Skipping Actions

`[skip ci]` can be used as a temporary escape hatch when a commit should skip
GitHub Actions, such as a docs-only change that should not run release
automation.

```bash
git commit -m "docs: update readme [skip ci]"
```
