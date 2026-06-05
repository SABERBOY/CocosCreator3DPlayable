package assets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResourceScriptEscapesScriptContext(t *testing.T) {
	script, err := ResourceScript([]Asset{
		{Path: `assets/payload".json`, Kind: KindJSON, Data: []byte(`"</script><script>window.__BAD=1</script>"`)},
		{Path: "assets/text.txt", Kind: KindText, Data: []byte(`plain </script> text`)},
		{Path: "assets/raw.bin", Kind: KindBinary, Data: []byte{0, 1, 2}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(script, "</script>") {
		t.Fatalf("resource script contains raw script terminator: %s", script)
	}
	if !strings.Contains(script, `\u003c/script\u003e`) {
		t.Fatalf("resource script did not escape angle brackets: %s", script)
	}
	if !strings.Contains(script, `assets/payload\".json`) {
		t.Fatalf("resource key was not quoted safely: %s", script)
	}
}

func TestResolveRejectsEscapingReferences(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "style.css"), []byte("body{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(root, "../secret.txt"); err == nil {
		t.Fatal("expected traversal reference to fail")
	}
	if _, err := Resolve(root, "https://example.test/style.css"); err == nil {
		t.Fatal("expected external reference to fail")
	}
	if _, err := Resolve(root, "style.css"); err != nil {
		t.Fatalf("expected local reference to resolve: %v", err)
	}
}
