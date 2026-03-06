// check-swagger-alignment.go - Compares implemented endpoints against API spec
// Usage: go run scripts/check-swagger-alignment.go [options]
//
// This tool reads the endpoint definitions from a YAML file and compares
// them against the endpoints implemented in the lib/ directory.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// EndpointsFile represents the YAML file structure
type EndpointsFile struct {
	Version   string     `yaml:"version"`
	BasePath  string     `yaml:"base_path"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

// Endpoint represents an API endpoint
type Endpoint struct {
	Method      string `yaml:"method"`
	Path        string `yaml:"path"`
	Description string `yaml:"description"`
	Category    string `yaml:"category"`
	Deprecated  bool   `yaml:"deprecated"`
}

// ImplementedEndpoint represents an endpoint found in source code
type ImplementedEndpoint struct {
	Method     string
	Path       string
	SourceFile string
	Function   string
}

func main() {
	endpointsFile := flag.String("endpoints-file", "./scripts/dci-endpoints.yaml", "Path to endpoints YAML file")
	libPath := flag.String("lib-path", "./lib", "Path to lib directory")
	baseURLVar := flag.String("base-url-var", "DCIURL|BaseURL", "Variable name pattern for base URL in source (regex)")
	outputFormat := flag.String("output", "text", "Output format: text, json, markdown")
	flag.Parse()

	// Load endpoints from YAML
	fmt.Println("Loading endpoints from", *endpointsFile, "...")
	specEndpoints, err := loadEndpointsFromYAML(*endpointsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading endpoints: %v\n", err)
		os.Exit(1)
	}

	// Parse implemented endpoints from source
	fmt.Println("Scanning lib/ for implemented endpoints...")
	implEndpoints, err := scanSourceEndpoints(*libPath, *baseURLVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning source: %v\n", err)
		os.Exit(1)
	}

	// Compare and generate report
	report := generateReport(specEndpoints, implEndpoints)

	switch *outputFormat {
	case "json":
		printJSONReport(report)
	case "markdown":
		printMarkdownReport(report)
	default:
		printTextReport(report)
	}
}

func loadEndpointsFromYAML(filePath string) ([]Endpoint, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var endpointsFile EndpointsFile
	if err := yaml.Unmarshal(data, &endpointsFile); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return endpointsFile.Endpoints, nil
}

func scanSourceEndpoints(libPath, baseURLVar string) ([]ImplementedEndpoint, error) {
	var endpoints []ImplementedEndpoint

	// Patterns to match endpoint definitions
	// Matches: DCIURL + "/path" or c.BaseURL + "/path" or fmt.Sprintf patterns
	urlConcatPattern := regexp.MustCompile(`(?:` + baseURLVar + `)\s*\+\s*"(/[^"]+)"`)
	sprintfPattern := regexp.MustCompile(`fmt\.Sprintf\s*\(\s*"%s(/[^"]+)"`)
	// Match HTTP method from http.NewRequest or helper function calls
	methodPattern := regexp.MustCompile(`(?:http\.NewRequest|httpGetWithAWSAuth|httpGetSimpleWithAWSAuth|httpPostWithAWSAuth|httpPostFileWithAWSAuth)\s*\(\s*(?:"(GET|POST|PUT|DELETE|PATCH)")?`)

	err := filepath.Walk(libPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		scanner := bufio.NewScanner(file)
		var currentMethod string
		var currentFunc string

		for scanner.Scan() {
			line := scanner.Text()

			// Track current function name
			if funcMatch := regexp.MustCompile(`func\s+(?:\([^)]+\)\s+)?(\w+)`).FindStringSubmatch(line); len(funcMatch) > 1 {
				currentFunc = funcMatch[1]
				currentMethod = "" // Reset method for new function
			}

			// Look for HTTP method
			if methodMatch := methodPattern.FindStringSubmatch(line); len(methodMatch) > 1 && methodMatch[1] != "" {
				currentMethod = methodMatch[1]
			}

			// Detect method from function names if not explicit
			if currentMethod == "" {
				funcLower := strings.ToLower(currentFunc)
				if strings.HasPrefix(funcLower, "get") || strings.HasPrefix(funcLower, "fetch") {
					currentMethod = "GET"
				} else if strings.HasPrefix(funcLower, "create") || strings.HasPrefix(funcLower, "upload") || strings.HasPrefix(funcLower, "schedule") {
					currentMethod = "POST"
				} else if strings.HasPrefix(funcLower, "update") {
					currentMethod = "PUT"
				} else if strings.HasPrefix(funcLower, "delete") {
					currentMethod = "DELETE"
				}
			}

			// Look for URL patterns
			var urlPath string
			if matches := urlConcatPattern.FindStringSubmatch(line); len(matches) > 1 {
				urlPath = matches[1]
			} else if matches := sprintfPattern.FindStringSubmatch(line); len(matches) > 1 {
				urlPath = matches[1]
			}

			if urlPath != "" {
				// Normalize the path (replace %s, %d with {})
				normalizedPath := normalizePath(urlPath)
				method := currentMethod
				if method == "" {
					method = "GET" // Default assumption
				}

				// Check if this endpoint is already added (avoid duplicates)
				isDuplicate := false
				for _, ep := range endpoints {
					if ep.Method == method && ep.Path == normalizedPath && ep.Function == currentFunc {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					endpoints = append(endpoints, ImplementedEndpoint{
						Method:     method,
						Path:       normalizedPath,
						SourceFile: filepath.Base(path),
						Function:   currentFunc,
					})
				}
			}
		}

		return scanner.Err()
	})

	return endpoints, err
}

func normalizePath(path string) string {
	// Replace %s, %d, %v with {}
	path = regexp.MustCompile(`%[sdv]`).ReplaceAllString(path, "{}")
	// Remove trailing slashes for comparison
	path = strings.TrimSuffix(path, "/")
	return path
}

func normalizeSpecPath(path string) string {
	// Replace {param_name} with {}
	path = regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(path, "{}")
	// Remove trailing slashes
	path = strings.TrimSuffix(path, "/")
	return path
}

// Report represents the comparison report
type Report struct {
	TotalSpec         int
	TotalImplemented  int
	Implemented       []Endpoint
	Missing           []Endpoint
	MissingDeprecated []Endpoint
	Extra             []ImplementedEndpoint
	ByCategory        map[string]CategoryReport
}

// CategoryReport represents coverage for a specific category
type CategoryReport struct {
	Total       int
	Implemented int
}

func generateReport(specEndpoints []Endpoint, implEndpoints []ImplementedEndpoint) Report {
	report := Report{
		TotalSpec:  len(specEndpoints),
		ByCategory: make(map[string]CategoryReport),
	}

	// Create a set of implemented paths for quick lookup
	implSet := make(map[string]bool)
	for _, impl := range implEndpoints {
		key := impl.Method + " " + impl.Path
		implSet[key] = true
	}

	// Create a set of spec paths for finding extras
	specSet := make(map[string]bool)

	for _, spec := range specEndpoints {
		normalizedSpecPath := normalizeSpecPath(spec.Path)
		key := spec.Method + " " + normalizedSpecPath
		specSet[key] = true

		// Update category stats
		category := spec.Category
		if category == "" {
			category = "other"
		}
		cr := report.ByCategory[category]
		cr.Total++
		if implSet[key] {
			cr.Implemented++
		}
		report.ByCategory[category] = cr

		// Check if implemented
		if implSet[key] {
			report.Implemented = append(report.Implemented, spec)
			report.TotalImplemented++
		} else {
			if spec.Deprecated {
				report.MissingDeprecated = append(report.MissingDeprecated, spec)
			} else {
				report.Missing = append(report.Missing, spec)
			}
		}
	}

	// Find extra implementations (in source but not in spec)
	for _, impl := range implEndpoints {
		key := impl.Method + " " + impl.Path
		if !specSet[key] {
			report.Extra = append(report.Extra, impl)
		}
	}

	return report
}

func printTextReport(report Report) {
	fmt.Println()
	fmt.Println("=============================================================")
	fmt.Println("                   API COVERAGE REPORT")
	fmt.Println("                go-dci vs DCI API Spec")
	fmt.Println("=============================================================")
	fmt.Println()

	coverage := float64(0)
	if report.TotalSpec > 0 {
		coverage = float64(report.TotalImplemented) / float64(report.TotalSpec) * 100
	}
	fmt.Println("SUMMARY")
	fmt.Println("-------")
	fmt.Printf("Total API Endpoints:     %d\n", report.TotalSpec)
	fmt.Printf("Implemented:             %d (%.1f%%)\n", report.TotalImplemented, coverage)
	fmt.Printf("Missing:                 %d\n", len(report.Missing))
	fmt.Printf("Missing (Deprecated):    %d\n", len(report.MissingDeprecated))
	fmt.Printf("Extra (not in spec):     %d\n", len(report.Extra))
	fmt.Println()

	// By category
	if len(report.ByCategory) > 0 {
		fmt.Println("BY CATEGORY")
		fmt.Println("-----------")

		// Sort categories by name
		var categories []string
		for cat := range report.ByCategory {
			categories = append(categories, cat)
		}
		sort.Strings(categories)

		for _, cat := range categories {
			cr := report.ByCategory[cat]
			pct := float64(0)
			if cr.Total > 0 {
				pct = float64(cr.Implemented) / float64(cr.Total) * 100
			}
			bar := progressBar(pct, 20)
			fmt.Printf("%-20s %s %2d/%-3d %5.1f%%\n", cat, bar, cr.Implemented, cr.Total, pct)
		}
		fmt.Println()
	}

	// Missing endpoints
	if len(report.Missing) > 0 {
		fmt.Println("MISSING ENDPOINTS")
		fmt.Println("-----------------")
		for _, ep := range report.Missing {
			fmt.Printf("- %-6s %-30s (%s)\n", ep.Method, ep.Path, ep.Description)
		}
		fmt.Println()
	}

	// Missing deprecated
	if len(report.MissingDeprecated) > 0 {
		fmt.Println("MISSING ENDPOINTS (Deprecated - Low Priority)")
		fmt.Println("----------------------------------------------")
		for _, ep := range report.MissingDeprecated {
			fmt.Printf("- %-6s %s\n", ep.Method, ep.Path)
		}
		fmt.Println()
	}

	// Extra implementations
	if len(report.Extra) > 0 {
		fmt.Println("EXTRA IMPLEMENTATIONS (Not in spec)")
		fmt.Println("------------------------------------")
		for _, ep := range report.Extra {
			fmt.Printf("- %-6s %-30s (%s)\n", ep.Method, ep.Path, ep.Function)
		}
		fmt.Println()
	}

	fmt.Println("=============================================================")
}

func progressBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("=", filled) + strings.Repeat(" ", width-filled) + "]"
}

func printJSONReport(report Report) {
	data, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(data))
}

func printMarkdownReport(report Report) {
	coverage := float64(0)
	if report.TotalSpec > 0 {
		coverage = float64(report.TotalImplemented) / float64(report.TotalSpec) * 100
	}

	fmt.Println("# API Coverage Report")
	fmt.Println()
	fmt.Println("## Summary")
	fmt.Println()
	fmt.Printf("| Metric | Value |\n")
	fmt.Printf("|--------|-------|\n")
	fmt.Printf("| Total API Endpoints | %d |\n", report.TotalSpec)
	fmt.Printf("| Implemented | %d (%.1f%%) |\n", report.TotalImplemented, coverage)
	fmt.Printf("| Missing | %d |\n", len(report.Missing))
	fmt.Printf("| Missing (Deprecated) | %d |\n", len(report.MissingDeprecated))
	fmt.Println()

	if len(report.Missing) > 0 {
		fmt.Println("## Missing Endpoints")
		fmt.Println()
		fmt.Println("| Method | Path | Description |")
		fmt.Println("|--------|------|-------------|")
		for _, ep := range report.Missing {
			fmt.Printf("| %s | %s | %s |\n", ep.Method, ep.Path, ep.Description)
		}
		fmt.Println()
	}
}
