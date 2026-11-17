package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func TestScanSourceEndpointsMissingDirectory(t *testing.T) {
	_, err := scanSourceEndpoints(filepath.Join(t.TempDir(), "absent"), "BaseURL")
	if err == nil {
		t.Fatal("expected an error for a missing library directory")
	}
	if !strings.Contains(err.Error(), "failed to open library directory") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// captureStdout runs fn with os.Stdout redirected and returns what it wrote.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	old := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = old })

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("close write end of pipe: %v", err)
	}
	out := <-done
	_ = r.Close()

	return out
}

func TestNormalizeSpecPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"no params", "/api/v1/identity", "/api/v1/identity"},
		{"single param", "/api/v1/topics/{id}", "/api/v1/topics/{}"},
		{"multiple params", "/api/v1/jobs/{job_id}/files/{file_id}", "/api/v1/jobs/{}/files/{}"},
		{"trailing slash", "/api/v1/topics/", "/api/v1/topics"},
		{"param and trailing slash", "/topics/{id}/", "/topics/{}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSpecPath(tt.path); got != tt.want {
				t.Fatalf("normalizeSpecPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestLoadEndpointsFromYAML(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "endpoints.yaml")
		content := `version: "v1"
base_path: "/api/v1"
endpoints:
  - path: /topics
    method: GET
    description: Get all topics
    category: topics
  - path: /jobs/{id}
    method: DELETE
    description: Delete a job
    category: jobs
    deprecated: true
`
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatalf("write endpoints file: %v", err)
		}

		endpoints, err := loadEndpointsFromYAML(path)
		if err != nil {
			t.Fatalf("load endpoints: %v", err)
		}
		if len(endpoints) != 2 {
			t.Fatalf("expected two endpoints, got %d", len(endpoints))
		}
		if endpoints[0].Method != "GET" || endpoints[0].Path != "/topics" || endpoints[0].Category != "topics" {
			t.Fatalf("unexpected first endpoint: %+v", endpoints[0])
		}
		if !endpoints[1].Deprecated || endpoints[1].Method != "DELETE" {
			t.Fatalf("unexpected second endpoint: %+v", endpoints[1])
		}
	})

	t.Run("malformed YAML", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "bad.yaml")
		if err := os.WriteFile(path, []byte("endpoints: [oops"), 0600); err != nil {
			t.Fatalf("write endpoints file: %v", err)
		}

		_, err := loadEndpointsFromYAML(path)
		if err == nil {
			t.Fatal("expected a YAML parse error")
		}
		if !strings.Contains(err.Error(), "failed to parse YAML") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := loadEndpointsFromYAML(filepath.Join(t.TempDir(), "absent.yaml"))
		if err == nil {
			t.Fatal("expected a read error")
		}
		if !strings.Contains(err.Error(), "failed to read file") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGenerateReport(t *testing.T) {
	spec := []Endpoint{
		{Method: "GET", Path: "/topics", Description: "Get all topics", Category: "topics"},
		{Method: "GET", Path: "/topics/{id}", Description: "Get a topic", Category: "topics"},
		{Method: "DELETE", Path: "/jobs/{id}", Description: "Delete a job", Category: "jobs", Deprecated: true},
		{Method: "GET", Path: "/files/{id}", Description: "Get a file"},
	}
	impl := []ImplementedEndpoint{
		{Method: "GET", Path: "/topics", Function: "GetTopics"},
		{Method: "GET", Path: "/topics/{}", Function: "GetTopic"},
		{Method: "GET", Path: "/analytics", Function: "GetAnalytics"},
	}

	report := generateReport(spec, impl)

	if report.TotalSpec != 4 {
		t.Fatalf("TotalSpec = %d, want 4", report.TotalSpec)
	}
	if report.TotalImplemented != 2 {
		t.Fatalf("TotalImplemented = %d, want 2", report.TotalImplemented)
	}
	if len(report.Implemented) != 2 {
		t.Fatalf("len(Implemented) = %d, want 2", len(report.Implemented))
	}
	if len(report.Missing) != 1 || report.Missing[0].Path != "/files/{id}" {
		t.Fatalf("unexpected missing endpoints: %+v", report.Missing)
	}
	if len(report.MissingDeprecated) != 1 || report.MissingDeprecated[0].Path != "/jobs/{id}" {
		t.Fatalf("unexpected missing deprecated endpoints: %+v", report.MissingDeprecated)
	}
	if len(report.Extra) != 1 || report.Extra[0].Path != "/analytics" {
		t.Fatalf("unexpected extra endpoints: %+v", report.Extra)
	}

	topics, ok := report.ByCategory["topics"]
	if !ok {
		t.Fatal("expected a topics category")
	}
	if topics.Total != 2 || topics.Implemented != 2 {
		t.Fatalf("unexpected topics category: %+v", topics)
	}
	if other := report.ByCategory["other"]; other.Total != 1 || other.Implemented != 0 {
		t.Fatalf("expected the uncategorized endpoint under \"other\", got %+v", report.ByCategory)
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		name  string
		pct   float64
		width int
		want  string
	}{
		{"empty", 0, 10, "[          ]"},
		{"half", 50, 10, "[=====     ]"},
		{"full", 100, 10, "[==========]"},
		{"overflows width", 150, 10, "[==========]"},
		{"quarter", 25, 10, "[==        ]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := progressBar(tt.pct, tt.width); got != tt.want {
				t.Fatalf("progressBar(%v, %d) = %q, want %q", tt.pct, tt.width, got, tt.want)
			}
		})
	}
}

func sampleReport() Report {
	return Report{
		TotalSpec:        4,
		TotalImplemented: 2,
		Implemented: []Endpoint{
			{Method: "GET", Path: "/topics", Category: "topics"},
			{Method: "GET", Path: "/topics/{id}", Category: "topics"},
		},
		Missing: []Endpoint{
			{Method: "GET", Path: "/files/{id}", Description: "Get a file"},
		},
		MissingDeprecated: []Endpoint{
			{Method: "DELETE", Path: "/jobs/{id}", Deprecated: true},
		},
		Extra: []ImplementedEndpoint{
			{Method: "GET", Path: "/analytics", Function: "GetAnalytics"},
		},
		ByCategory: map[string]CategoryReport{
			"topics": {Total: 2, Implemented: 2},
			"other":  {Total: 1, Implemented: 0},
		},
	}
}

func TestPrintTextReport(t *testing.T) {
	out := captureStdout(t, func() { printTextReport(sampleReport()) })

	for _, want := range []string{
		"API COVERAGE REPORT",
		"Total API Endpoints:     4",
		"Implemented:             2 (50.0%)",
		"Missing:                 1",
		"Missing (Deprecated):    1",
		"Extra (not in spec):     1",
		"BY CATEGORY",
		"topics",
		"[====================]",
		"MISSING ENDPOINTS",
		"MISSING ENDPOINTS (Deprecated - Low Priority)",
		"EXTRA IMPLEMENTATIONS (Not in spec)",
		"/analytics",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("text report missing %q:\n%s", want, out)
		}
	}

	// Categories are sorted, so "other" is printed before "topics".
	if strings.Index(out, "other") > strings.Index(out, "topics") {
		t.Fatalf("categories are not sorted:\n%s", out)
	}
}

func TestPrintTextReportEmptySpec(t *testing.T) {
	out := captureStdout(t, func() { printTextReport(Report{}) })

	if !strings.Contains(out, "Implemented:             0 (0.0%)") {
		t.Fatalf("expected a zero coverage summary, got:\n%s", out)
	}
	if strings.Contains(out, "BY CATEGORY") {
		t.Fatalf("expected no category section for an empty report, got:\n%s", out)
	}
}

func TestPrintJSONReport(t *testing.T) {
	want := sampleReport()
	out := captureStdout(t, func() { printJSONReport(want) })

	var got Report
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("unmarshal JSON report: %v\n%s", err, out)
	}
	if got.TotalSpec != want.TotalSpec || got.TotalImplemented != want.TotalImplemented {
		t.Fatalf("unexpected totals: %+v", got)
	}
	if len(got.Missing) != 1 || len(got.MissingDeprecated) != 1 || len(got.Extra) != 1 {
		t.Fatalf("unexpected report contents: %+v", got)
	}
	if got.ByCategory["topics"].Implemented != 2 {
		t.Fatalf("unexpected categories: %+v", got.ByCategory)
	}
}

func TestPrintMarkdownReport(t *testing.T) {
	out := captureStdout(t, func() { printMarkdownReport(sampleReport()) })

	if !strings.HasPrefix(out, "# API Coverage Report") {
		t.Fatalf("markdown report does not start with a heading:\n%s", out)
	}
	for _, want := range []string{
		"| Total API Endpoints | 4 |",
		"| Implemented | 2 (50.0%) |",
		"| Missing | 1 |",
		"| Missing (Deprecated) | 1 |",
		"## Missing Endpoints",
		"| GET | /files/{id} | Get a file |",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("markdown report missing %q:\n%s", want, out)
		}
	}
}

func TestPrintMarkdownReportWithoutMissing(t *testing.T) {
	report := Report{TotalSpec: 1, TotalImplemented: 1, ByCategory: map[string]CategoryReport{}}
	out := captureStdout(t, func() { printMarkdownReport(report) })

	if !strings.Contains(out, "| Implemented | 1 (100.0%) |") {
		t.Fatalf("unexpected summary:\n%s", out)
	}
	if strings.Contains(out, "## Missing Endpoints") {
		t.Fatalf("expected no missing section:\n%s", out)
	}
}
