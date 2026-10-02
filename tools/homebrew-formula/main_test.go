package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratePinsSourceAndChecksum(t *testing.T) {
	archive := filepath.Join(t.TempDir(), "source.tar.gz")
	if err := os.WriteFile(archive, []byte("abc"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := generate(&out, "v0.1.5", archive); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`url "https://github.com/flexdinesh/checksy/archive/refs/tags/v0.1.5.tar.gz"`,
		`sha256 "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"`,
		`depends_on "go" => :build`,
		`github.com/flexdinesh/checksy/internal/version.Version=#{version}`,
		`system "go", "build", *std_go_args(ldflags:), "./cmd/checksy"`,
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("formula missing %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "releases/download/") || strings.Contains(out.String(), "bin.install") {
		t.Fatal("formula must build source, not install prebuilt binaries")
	}
}

func TestGenerateRejectsInvalidTagsWithoutOutput(t *testing.T) {
	for _, tag := range []string{"", "dev", "v0.1.5-rc.1", "0.1.5", "v01.1.5", `v0.1.5"`} {
		t.Run(tag, func(t *testing.T) {
			var out bytes.Buffer
			if err := generate(&out, tag, "unused"); err == nil {
				t.Fatal("expected invalid tag error")
			}
			if out.Len() != 0 {
				t.Fatal("invalid tag must not produce a formula")
			}
		})
	}
}

func TestGenerateRejectsMissingSourceWithoutOutput(t *testing.T) {
	var out bytes.Buffer
	if err := generate(&out, "v0.1.5", filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("expected missing source error")
	}
	if out.Len() != 0 {
		t.Fatal("missing source must not produce a formula")
	}
}
