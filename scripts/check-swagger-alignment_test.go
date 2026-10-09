package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanSourceEndpointsReadsFilesWithinRoot(t *testing.T) {
	libDir := t.TempDir()
	source := `package lib
func (c *Client) GetJob(ctx context.Context) {
	reqURL := BaseURL + "/jobs"
	c.doJSON(ctx, http.MethodGet, reqURL, nil)
}`
	if err := os.WriteFile(filepath.Join(libDir, "client.go"), []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	endpoints, err := scanSourceEndpoints(libDir, "BaseURL")
	if err != nil {
		t.Fatalf("scan source endpoints: %v", err)
	}
	if len(endpoints) != 1 {
		t.Fatalf("expected one endpoint, got %d", len(endpoints))
	}
	if endpoints[0].Method != "GET" || endpoints[0].Path != "/jobs" {
		t.Fatalf("unexpected endpoint: %+v", endpoints[0])
	}
}

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
