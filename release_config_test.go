package checksy_test

import (
	"os"
	"os/exec"
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

func TestGoReleaserPackagesStableReleasesForSupportedPlatforms(t *testing.T) {
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
		"prerelease: \"false\"",
		"make_latest: true",
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
	} {
		if !strings.Contains(ci, want) {
			t.Fatalf("CI workflow should contain %q", want)
		}
	}
}

func TestDevWorkflowPublishesOnlyAfterMainPassesCI(t *testing.T) {
	workflow := readFile(t, ".github/workflows/ci.yml")
	for _, want := range []string{
		"needs: test",
		"if: github.event_name == 'push' && github.ref == 'refs/heads/main'",
		"contents: write",
		"group: publish-dev",
		"cancel-in-progress: false",
		"fetch-depth: 0",
	} {
		if !strings.Contains(workflow, want) {
			t.Fatalf("CI workflow should contain %q", want)
		}
	}
	for _, obsolete := range []string{"- dev", "refs/heads/dev'", "goreleaser", "git tag"} {
		if strings.Contains(workflow, obsolete) {
			t.Fatalf("CI workflow should not contain %q", obsolete)
		}
	}
}

func TestDevPublicationTracksMainWithoutChangingStableTags(t *testing.T) {
	workflow := readFile(t, ".github/workflows/ci.yml")
	_, step, ok := strings.Cut(workflow, "      - name: Publish dev branch\n")
	if !ok {
		t.Fatal("missing dev publication step")
	}
	_, block, ok := strings.Cut(step, "        run: |\n")
	if !ok {
		t.Fatal("missing dev publication script")
	}
	var script strings.Builder
	for _, line := range strings.Split(block, "\n") {
		if line != "" && !strings.HasPrefix(line, "          ") {
			break
		}
		script.WriteString(strings.TrimPrefix(line, "          ") + "\n")
	}

	root := t.TempDir()
	run := func(dir, name string, args ...string) string {
		t.Helper()
		cmd := exec.Command(name, args...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v: %v\n%s", name, args, err, output)
		}
		return strings.TrimSpace(string(output))
	}
	remote := filepath.Join(root, "origin.git")
	checkout := filepath.Join(root, "checkout")
	run(root, "git", "init", "--bare", "--initial-branch=main", remote)
	run(root, "git", "clone", remote, checkout)
	run(checkout, "git", "config", "user.name", "Release test")
	run(checkout, "git", "config", "user.email", "release@example.com")
	run(checkout, "git", "config", "commit.gpgsign", "false")
	run(checkout, "git", "config", "tag.gpgsign", "false")
	run(checkout, "git", "commit", "--allow-empty", "-m", "Stable release")
	stable := run(checkout, "git", "rev-parse", "HEAD")
	run(checkout, "git", "tag", "v0.1.0")
	run(checkout, "git", "push", "origin", "main", "refs/tags/v0.1.0")

	assertDev := func(want string) {
		t.Helper()
		if got := run(remote, "git", "rev-parse", "refs/heads/dev"); got != want {
			t.Fatalf("dev = %s, want %s", got, want)
		}
		if got := run(remote, "git", "show-ref", "--tags"); got != stable+" refs/tags/v0.1.0" {
			t.Fatalf("stable tags changed: %s", got)
		}
	}
	// First publication creates the branch, including when main is tagged.
	run(checkout, "bash", "-c", script.String())
	assertDev(stable)

	run(checkout, "git", "commit", "--allow-empty", "-m", "Development change")
	latest := run(checkout, "git", "rev-parse", "HEAD")
	run(checkout, "git", "push", "origin", "main")

	// A queued old run must not publish while the newer main awaits CI.
	run(checkout, "git", "checkout", "--detach", stable)
	run(checkout, "bash", "-c", script.String())
	assertDev(stable)

	// Detached Actions checkouts publish, and reruns are idempotent.
	run(checkout, "git", "checkout", "--detach", latest)
	run(checkout, "bash", "-c", script.String())
	assertDev(latest)
	run(checkout, "bash", "-c", script.String())
	assertDev(latest)

	// Rerunning an older workflow cannot roll dev back.
	run(checkout, "git", "checkout", "--detach", stable)
	run(checkout, "bash", "-c", script.String())
	assertDev(latest)

	// A conflicting update to dev must be rejected, never force-pushed away.
	run(checkout, "git", "checkout", "--detach", latest)
	run(checkout, "git", "commit", "--allow-empty", "-m", "Concurrent dev update")
	concurrent := run(checkout, "git", "rev-parse", "HEAD")
	run(checkout, "git", "push", "origin", "HEAD:refs/heads/dev")
	run(checkout, "git", "checkout", "--detach", latest)
	cmd := exec.Command("bash", "-c", script.String())
	cmd.Dir = checkout
	if output, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("expected conflicting dev update to fail:\n%s", output)
	}
	assertDev(concurrent)
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
		"version: v2.18.0",
		"args: release --clean",
	} {
		if !strings.Contains(workflow, want) {
			t.Fatalf("release workflow should contain %q", want)
		}
	}
	if strings.Contains(workflow, "  push:") || strings.Contains(workflow, "  pull_request:") {
		t.Fatal("stable releases must only run via workflow_dispatch")
	}
}

func TestStableReleasePublishesHomebrewSourceFormulaPullRequest(t *testing.T) {
	config := readFile(t, ".goreleaser.yaml")
	for _, deprecated := range []string{
		"brews:",
		"homebrew_casks:",
	} {
		if strings.Contains(config, deprecated) {
			t.Fatalf("GoReleaser must not publish a prebuilt Homebrew package: %q", deprecated)
		}
	}

	workflow := readFile(t, ".github/workflows/release.yml")
	for _, want := range []string{
		"version: v2.18.0",
		"repository: flexdinesh/homebrew-tap",
		"token: ${{ secrets.HOMEBREW_TAP_TOKEN }}",
		"https://github.com/flexdinesh/checksy/archive/refs/tags/${RELEASE_TAG}.tar.gz",
		"go run ./tools/homebrew-formula -tag \"$RELEASE_TAG\"",
		"-archive \"$RUNNER_TEMP/checksy-source.tar.gz\" > homebrew-tap/Formula/checksy.rb",
		"rm -f homebrew-tap/Casks/checksy.rb",
		"uses: peter-evans/create-pull-request@v8",
		"path: homebrew-tap",
		"base: main",
		"branch: checksy-${{ steps.tag.outputs.tag }}",
	} {
		if !strings.Contains(workflow, want) {
			t.Fatalf("release workflow should contain %q", want)
		}
	}
}

func TestReleaseDocsExplainHomebrewChannel(t *testing.T) {
	readme := readFile(t, "README.md")
	for _, want := range []string{
		"brew install flexdinesh/tap/checksy",
		"brew uninstall --cask checksy",
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
		"brew install flexdinesh/tap/checksy",
		"brew uninstall --cask checksy",
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
