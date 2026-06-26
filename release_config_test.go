package checksy_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseIdentityIsCanonical(t *testing.T) {
	staleModulePath := "github.com/dineshpandiyan" + "/checksy"
	goMod := readFile(t, "go.mod")
	if !strings.Contains(goMod, "module github.com/flexdinesh/checksy") {
		t.Fatalf("go.mod should use github.com/flexdinesh/checksy:\n%s", goMod)
	}

	readme := readFile(t, "README.md")
	if !strings.Contains(readme, "go install github.com/flexdinesh/checksy/cmd/checksy@latest") {
		t.Fatalf("README should document the canonical Go install path")
	}

	for _, path := range goFiles(t, ".") {
		contents := readFile(t, path)
		if strings.Contains(contents, staleModulePath) {
			t.Fatalf("%s still references %s", path, staleModulePath)
		}
	}
}

func TestGoReleaserPackagesSnapshotsForSupportedPlatforms(t *testing.T) {
	config := readFile(t, ".goreleaser.yaml")
	for _, want := range []string{
		"project_name: checksy",
		"main: ./cmd/checksy",
		"binary: checksy",
		"CGO_ENABLED=0",
		"-trimpath",
		"darwin",
		"linux",
		"amd64",
		"arm64",
		"-X github.com/flexdinesh/checksy/internal/version.Version={{.Version}}",
		"checksums.txt",
		"replace_existing_artifacts: true",
	} {
		if !strings.Contains(config, want) {
			t.Fatalf(".goreleaser.yaml should contain %q", want)
		}
	}

	ci := readFile(t, ".github/workflows/ci.yml")
	for _, want := range []string{
		"go test ./...",
		"go build ./cmd/checksy",
		"version: v2.16.0",
		"args: release --snapshot --clean",
	} {
		if !strings.Contains(ci, want) {
			t.Fatalf("CI workflow should contain %q", want)
		}
	}
}

func TestStableReleaseWorkflowPublishesSemverTags(t *testing.T) {
	workflow := readFile(t, ".github/workflows/release.yml")
	for _, want := range []string{
		"workflow_dispatch",
		"contents: write",
		"ref: main",
		"fetch-depth: 0",
		"go test ./...",
		"git tag -l 'v0.1.*'",
		"next=\"v0.1.0\"",
		"git push origin",
		"version: v2.16.0",
		"args: release --clean",
	} {
		if !strings.Contains(workflow, want) {
			t.Fatalf("release workflow should contain %q", want)
		}
	}
}

func TestGoReleaserPublishesHomebrewTapPullRequest(t *testing.T) {
	config := readFile(t, ".goreleaser.yaml")
	for _, want := range []string{
		"homebrew_casks:",
		"name: checksy",
		"owner: flexdinesh",
		"name: homebrew-tap",
		"token: \"{{ .Env.HOMEBREW_TAP_TOKEN }}\"",
		"branch: \"checksy-{{ .Tag }}\"",
		"pull_request:",
		"enabled: true",
		"directory: Casks",
		"homepage: https://github.com/flexdinesh/checksy",
		"binaries:",
		"- checksy",
		"caveats:",
		"Run `checksy --help` to view available command-line options.",
	} {
		if !strings.Contains(config, want) {
			t.Fatalf(".goreleaser.yaml should contain %q", want)
		}
	}

	for _, deprecated := range []string{
		"brews:",
		"directory: Formula",
	} {
		if strings.Contains(config, deprecated) {
			t.Fatalf(".goreleaser.yaml should not contain deprecated formula config %q", deprecated)
		}
	}

	workflow := readFile(t, ".github/workflows/release.yml")
	for _, want := range []string{
		"version: v2.16.0",
		"HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}",
	} {
		if !strings.Contains(workflow, want) {
			t.Fatalf("release workflow should contain %q", want)
		}
	}
}

func TestReleaseDocsExplainHomebrewChannel(t *testing.T) {
	readme := readFile(t, "README.md")
	for _, want := range []string{
		"brew install --cask flexdinesh/tap/checksy",
		"go install github.com/flexdinesh/checksy/cmd/checksy@latest",
	} {
		if !strings.Contains(readme, want) {
			t.Fatalf("README should contain %q", want)
		}
	}

	releaseDoc := readFile(t, "docs/release.md")
	for _, want := range []string{
		"HOMEBREW_TAP_TOKEN",
		"v0.1.0",
		"brew install --cask flexdinesh/tap/checksy",
		"goreleaser release --snapshot --clean",
		"Do not create a moving `latest` tag",
	} {
		if !strings.Contains(releaseDoc, want) {
			t.Fatalf("docs/release.md should contain %q", want)
		}
	}

	adr := readFile(t, "docs/adr/0003-homebrew-tap-releases.md")
	if !strings.Contains(adr, "flexdinesh/homebrew-tap") {
		t.Fatalf("Homebrew release ADR should record the custom tap")
	}

	glossary := readFile(t, "docs/glossary.md")
	if !strings.Contains(glossary, "Homebrew Tap") {
		t.Fatalf("glossary should define the Homebrew tap")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(contents)
}

func goFiles(t *testing.T, root string) []string {
	t.Helper()
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch path {
			case ".git", ".scratch", "bin", "dist":
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return paths
}
