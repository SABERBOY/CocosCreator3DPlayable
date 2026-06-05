package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackCommand(t *testing.T) {
	output := filepath.Join(t.TempDir(), "playable.html")
	cmd := newRootCommand()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{
		"pack",
		"--input", filepath.Join("..", "..", "testdata", "cocos3x", "basic"),
		"--output", output,
	})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "wrote ") {
		t.Fatalf("unexpected stdout: %s", stdout.String())
	}
}
