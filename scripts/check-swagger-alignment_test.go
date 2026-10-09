package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanSourceEndpointsRejectsSymlinksOutsideRoot(t *testing.T) {
	libDir := t.TempDir()
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "outside.go")
	if err := os.WriteFile(outsideFile, []byte("package lib\n"), 0600); err != nil {
		t.Fatalf("write outside source file: %v", err)
	}

	if err := os.Symlink(outsideFile, filepath.Join(libDir, "client.go")); err != nil {
		t.Skipf("cannot create symlink in this environment: %v", err)
	}

	if _, err := scanSourceEndpoints(libDir, "BaseURL"); err == nil {
		t.Fatal("expected source scan to reject a symlink that escapes its root")
	}
}
