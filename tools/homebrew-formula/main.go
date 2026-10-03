package main

import (
	"crypto/sha256"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"text/template"
)

//go:embed checksy.rb.tmpl
var formulaTemplate string

func main() {
	tag := flag.String("tag", "", "Stable release tag")
	archivePath := flag.String("archive", "", "Downloaded source archive")
	flag.Parse()
	if err := generate(os.Stdout, *tag, *archivePath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generate(out io.Writer, tag, archivePath string) error {
	if !regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`).MatchString(tag) {
		return fmt.Errorf("invalid stable release tag: %q", tag)
	}
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, archive); err != nil {
		return err
	}
	formula, err := template.New("checksy").Parse(formulaTemplate)
	if err != nil {
		return err
	}
	return formula.Execute(out, struct {
		Tag    string
		SHA256 string
	}{tag, fmt.Sprintf("%x", hash.Sum(nil))})
}
