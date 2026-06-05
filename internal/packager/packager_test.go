package packager

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPacksFixtureIntoSingleHTML(t *testing.T) {
	input := filepath.Join("..", "..", "testdata", "cocos3x", "basic")
	output := filepath.Join(t.TempDir(), "playable.html")

	report, err := Run(context.Background(), Options{
		InputDir:   input,
		OutputFile: output,
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.ResourceCount == 0 || report.ScriptCount == 0 {
		t.Fatalf("expected resources and scripts, got %+v", report)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	html := string(data)
	for _, forbidden := range []string{` src=`, ` href=`, `</script><script>window.__BAD`} {
		if strings.Contains(html, forbidden) {
			t.Fatalf("output contains forbidden fragment %q in:\n%s", forbidden, html)
		}
	}
	if !strings.Contains(html, `window.res`) || !strings.Contains(html, `System.import('index.js')`) {
		t.Fatalf("output does not look like a playable bundle:\n%s", html)
	}
}

func TestRunMinifyPreservesScriptEscaping(t *testing.T) {
	input := filepath.Join("..", "..", "testdata", "cocos3x", "basic")
	output := filepath.Join(t.TempDir(), "playable.html")
	if _, err := Run(context.Background(), Options{InputDir: input, OutputFile: output, Minify: true}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `</script><script>window.__BAD`) {
		t.Fatalf("minified output reintroduced raw script terminator:\n%s", string(data))
	}
}

func TestRunRefusesOverwriteWithoutForce(t *testing.T) {
	input := filepath.Join("..", "..", "testdata", "cocos3x", "basic")
	output := filepath.Join(t.TempDir(), "playable.html")
	if err := os.WriteFile(output, []byte("exists"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(context.Background(), Options{InputDir: input, OutputFile: output}); err == nil {
		t.Fatal("expected overwrite without force to fail")
	}
	if _, err := Run(context.Background(), Options{InputDir: input, OutputFile: output, Force: true}); err != nil {
		t.Fatalf("expected overwrite with force to succeed: %v", err)
	}
}

func TestRunRejectsTraversalScriptReference(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "web-mobile")
	if err := os.MkdirAll(filepath.Join(input, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		filepath.Join(root, "secret.txt"):              "secret",
		filepath.Join(input, "src", "import-map.json"): `{}`,
		filepath.Join(input, "index.js"):               `System.register([], function () { return { setters: [], execute: function () {} }; });`,
		filepath.Join(input, "index.html"):             `<html><body><script src="../secret.txt"></script><script>System.import('index.js')</script></body></html>`,
	}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	_, err := Run(context.Background(), Options{InputDir: input, OutputFile: filepath.Join(root, "out.html")})
	if err == nil || !strings.Contains(err.Error(), "escapes input directory") {
		t.Fatalf("expected traversal rejection, got %v", err)
	}
}
