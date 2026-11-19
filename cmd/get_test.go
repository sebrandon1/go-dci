package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestFindOcpVersionFromComponents(t *testing.T) {
	testCases := []struct {
		testComponents     []lib.Components
		expectedOcpVersion string
	}{
		{
			testComponents: []lib.Components{
				{
					Name: "OpenShift 4.14.2",
				},
			},
			expectedOcpVersion: "4.14.2",
		},
		{
			testComponents:     []lib.Components{},
			expectedOcpVersion: "",
		},
		{
			testComponents: []lib.Components{
				{
					Name: "OpenShift-4.19.0",
				},
			},
			expectedOcpVersion: "4.19.0",
		},
		{
			testComponents: []lib.Components{
				{
					Name: "OpenShift",
				},
			},
			expectedOcpVersion: "",
		},
		{
			testComponents: []lib.Components{
				{
					Name: "OpenShift Container Platform 4.20.1",
				},
			},
			expectedOcpVersion: "4.20.1",
		},
		{
			testComponents: []lib.Components{
				{
					Name: "some-other 4.1",
				},
			},
			expectedOcpVersion: "",
		},
		{
			testComponents: []lib.Components{
				{
					Name:    "OpenShift",
					Version: "4.21.0",
				},
			},
			expectedOcpVersion: "4.21.0",
		},
		{
			testComponents: []lib.Components{
				{
					Name: "OpenShift",
				},
				{
					Name: "OpenShift 4.22.0",
				},
			},
			expectedOcpVersion: "4.22.0",
		},
		{
			testComponents: []lib.Components{
				{
					Name:    "OpenShift",
					Version: "not-a-version",
				},
			},
			expectedOcpVersion: "",
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedOcpVersion, findOcpVersionFromComponents(testCase.testComponents))
	}
}

func TestIsCertsuiteJob(t *testing.T) {
	testCases := []struct {
		testComponents []lib.Components
		expectedOutput bool
	}{
		{
			testComponents: []lib.Components{
				{
					Name: "cnf-certification-test",
				},
			},
			expectedOutput: true,
		},
		{
			testComponents: []lib.Components{
				{
					Name: "some-other-test",
				},
			},
			expectedOutput: false,
		},
		{
			testComponents: []lib.Components{
				{
					Name: "certsuite",
				},
				{
					Name: "some-other-test",
				},
			},
			expectedOutput: true,
		},
		{
			testComponents: []lib.Components{
				{
					Name: "some-other-test",
				},
				{
					Name: "another-test",
				},
			},
			expectedOutput: false,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedOutput, isCertsuiteJob(testCase.testComponents))
	}
}

func TestCountOcpVersions(t *testing.T) {
	testCases := []struct {
		testJobsResponses []lib.JobsResponse
		expectedCount     map[string]int
	}{
		{
			testJobsResponses: []lib.JobsResponse{
				{
					Jobs: []lib.Job{
						{
							Components: []lib.Components{
								{
									Name: "cnf-certification-test",
								},
								{
									Name: "OpenShift 4.14.2",
								},
							},
						},
					},
				},
			},
			expectedCount: map[string]int{
				"4.14": 1,
			},
		},
		{
			testJobsResponses: []lib.JobsResponse{
				{
					Jobs: []lib.Job{
						{
							Components: []lib.Components{
								{
									Name: "certsuite",
								},
								{
									Name: "OpenShift 4.14.2",
								},
							},
						},
						{
							Components: []lib.Components{
								{
									Name: "OpenShift 4.14.2",
								},
								{
									Name: "OpenShift 4.14.3",
								},
							},
						},
					},
				},
			},
			expectedCount: map[string]int{
				"4.14": 1,
			},
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedCount, countOcpVersions(testCase.testJobsResponses))
	}
}

func TestExtractCommitVersion(t *testing.T) {
	testCases := []struct {
		testComponent  string
		expectedOutput string
	}{
		{
			testComponent:  "cnf-certification-test 4321testcommit",
			expectedOutput: "4321testcommit",
		},
		{
			testComponent:  "cnf-certification-test",
			expectedOutput: "unknown",
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedOutput, extractCommitVersion(testCase.testComponent))
	}
}

func TestPrintComponentsStdout(t *testing.T) {
	// Test that the function doesn't panic with valid input
	componentsResponses := []lib.ComponentsResponse{
		{
			Meta: lib.Meta{Count: 2},
			Components: []lib.Components{
				{
					ID:      "comp-1",
					Name:    "Test Component 1",
					Type:    "ocp",
					Version: "4.14.1",
					TopicID: "topic-123",
				},
				{
					ID:      "comp-2",
					Name:    "Test Component 2",
					Type:    "certsuite",
					Version: "v1.0.0",
					TopicID: "topic-456",
				},
			},
		},
	}

	// This test just verifies the function doesn't panic
	// In a production scenario, you might capture stdout and verify the output
	assert.NotPanics(t, func() {
		printComponentsStdout(componentsResponses)
	})
}

func TestPrintComponentsStdout_EmptyResponse(t *testing.T) {
	componentsResponses := []lib.ComponentsResponse{}

	assert.NotPanics(t, func() {
		printComponentsStdout(componentsResponses)
	})
}

func TestPrintComponentsJSON(t *testing.T) {
	componentsResponses := []lib.ComponentsResponse{
		{
			Meta: lib.Meta{Count: 1},
			Components: []lib.Components{
				{
					ID:      "comp-1",
					Name:    "Test Component",
					Type:    "test",
					Version: "1.0.0",
					TopicID: "topic-123",
				},
			},
		},
	}

	if err := printComponentsJSON(componentsResponses); err != nil {
		t.Fatalf("printComponentsJSON returned an error: %v", err)
	}
}

func TestPrintComponentsJSON_EmptyComponents(t *testing.T) {
	componentsResponses := []lib.ComponentsResponse{
		{
			Meta:       lib.Meta{Count: 0},
			Components: []lib.Components{},
		},
	}

	if err := printComponentsJSON(componentsResponses); err != nil {
		t.Fatalf("printComponentsJSON returned an error: %v", err)
	}
}

func TestPrintComponentsJSON_MultipleResponses(t *testing.T) {
	// Test that multiple ComponentsResponse objects are flattened correctly
	componentsResponses := []lib.ComponentsResponse{
		{
			Components: []lib.Components{
				{ID: "comp-1", Name: "Component 1"},
			},
		},
		{
			Components: []lib.Components{
				{ID: "comp-2", Name: "Component 2"},
				{ID: "comp-3", Name: "Component 3"},
			},
		},
	}

	if err := printComponentsJSON(componentsResponses); err != nil {
		t.Fatalf("printComponentsJSON returned an error: %v", err)
	}
}

func TestPrintIdentityStdout(t *testing.T) {
	identity := &lib.IdentityResponse{
		Identity: lib.Identity{
			ID:       "remoteci-123",
			Name:     "test-remoteci",
			Type:     "remoteci",
			TeamID:   "team-456",
			TeamName: "Test Team",
			State:    "active",
		},
	}

	assert.NotPanics(t, func() {
		printIdentityStdout(identity)
	})
}

func TestPrintIdentityStdout_WithAllFields(t *testing.T) {
	identity := &lib.IdentityResponse{
		Identity: lib.Identity{
			ID:       "user-789",
			Name:     "testuser",
			Type:     "user",
			Email:    "test@example.com",
			Fullname: "Test User",
			TeamID:   "team-456",
			TeamName: "Test Team",
			State:    "active",
		},
	}

	assert.NotPanics(t, func() {
		printIdentityStdout(identity)
	})
}

func TestPrintIdentityJSON(t *testing.T) {
	identity := &lib.IdentityResponse{
		Identity: lib.Identity{
			ID:       "remoteci-123",
			Name:     "test-remoteci",
			Type:     "remoteci",
			TeamID:   "team-456",
			TeamName: "Test Team",
			State:    "active",
		},
	}

	if err := printIdentityJSON(identity); err != nil {
		t.Fatalf("printIdentityJSON returned an error: %v", err)
	}
}

func TestPrintComponentTypesStdout(t *testing.T) {
	componentTypesResponses := []lib.ComponentTypesResponse{
		{
			Meta: lib.Meta{Count: 2},
			ComponentTypes: []lib.ComponentType{
				{
					ID:    "ct-1",
					Name:  "ocp",
					State: "active",
				},
				{
					ID:    "ct-2",
					Name:  "certsuite",
					State: "active",
				},
			},
		},
	}

	assert.NotPanics(t, func() {
		printComponentTypesStdout(componentTypesResponses)
	})
}

func TestPrintComponentTypesStdout_EmptyResponse(t *testing.T) {
	componentTypesResponses := []lib.ComponentTypesResponse{}

	assert.NotPanics(t, func() {
		printComponentTypesStdout(componentTypesResponses)
	})
}

func TestPrintComponentTypesJSON(t *testing.T) {
	componentTypesResponses := []lib.ComponentTypesResponse{
		{
			Meta: lib.Meta{Count: 2},
			ComponentTypes: []lib.ComponentType{
				{
					ID:    "ct-1",
					Name:  "ocp",
					State: "active",
				},
				{
					ID:    "ct-2",
					Name:  "certsuite",
					State: "active",
				},
			},
		},
	}

	if err := printComponentTypesJSON(componentTypesResponses); err != nil {
		t.Fatalf("printComponentTypesJSON returned an error: %v", err)
	}
}

func TestPrintComponentTypesJSON_MultipleResponses(t *testing.T) {
	componentTypesResponses := []lib.ComponentTypesResponse{
		{
			ComponentTypes: []lib.ComponentType{
				{ID: "ct-1", Name: "ocp"},
			},
		},
		{
			ComponentTypes: []lib.ComponentType{
				{ID: "ct-2", Name: "certsuite"},
				{ID: "ct-3", Name: "rhel"},
			},
		},
	}

	if err := printComponentTypesJSON(componentTypesResponses); err != nil {
		t.Fatalf("printComponentTypesJSON returned an error: %v", err)
	}
}

func TestPrintTopicsStdout(t *testing.T) {
	topicsResponses := []lib.TopicsResponse{
		{
			Meta: lib.Meta{Count: 1},
			Topics: []lib.Topic{
				{
					ID:        "topic-123",
					Name:      "OCP-4.14",
					ProductID: "prod-456",
					Product:   lib.Product{Name: "OpenShift"},
					State:     "active",
				},
			},
		},
	}
	assert.NotPanics(t, func() {
		printTopicsStdout(topicsResponses)
	})
}

func TestPrintTopicsStdout_EmptyResponse(t *testing.T) {
	topicsResponses := []lib.TopicsResponse{}
	assert.NotPanics(t, func() {
		printTopicsStdout(topicsResponses)
	})
}

func TestPrintTopicsJSON(t *testing.T) {
	topicsResponses := []lib.TopicsResponse{
		{
			Meta: lib.Meta{Count: 1},
			Topics: []lib.Topic{
				{
					ID:            "topic-123",
					Name:          "OCP-4.14",
					ProductID:     "prod-456",
					Product:       lib.Product{Name: "OpenShift"},
					State:         "active",
					ExportControl: true,
					CreatedAt:     "2024-01-01T00:00:00.000000",
					UpdatedAt:     "2024-06-01T00:00:00.000000",
				},
			},
		},
	}
	if err := printTopicsJSON(topicsResponses); err != nil {
		t.Fatalf("printTopicsJSON returned an error: %v", err)
	}
}

func TestCalculateDaysSince(t *testing.T) {
	yesterday := time.Now().Add(-24 * time.Hour).Format(dateFormat)
	days, err := calculateDaysSince(yesterday)
	assert.NoError(t, err)
	assert.InDelta(t, 1.0, days, 0.5)
}

func TestCalculateDaysSince_InvalidTimestamp(t *testing.T) {
	days, err := calculateDaysSince("not-a-timestamp")
	assert.Error(t, err)
	assert.Equal(t, 0.0, days)
	assert.Contains(t, err.Error(), "failed to parse timestamp")
}

func TestGetIdentityCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lib.IdentityResponse{Identity: lib.Identity{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "test-remoteci", Type: "remoteci"}}); err != nil {
			t.Errorf("encode identity response: %v", err)
		}
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	outputFormat = OutputFormatStdout
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getIdentityCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetIdentityCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	outputFormat = OutputFormatStdout
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getIdentityCmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestRunGeneralJobsOutput_Stdout(t *testing.T) {
	responses := []lib.JobsResponse{
		{
			Jobs: []lib.Job{
				{
					ID:        "job-1",
					Name:      "nightly",
					Status:    "success",
					State:     "active",
					TopicID:   "topic-abc",
					CreatedAt: "2026-01-01T00:00:00.000000",
				},
				{
					ID:        "job-2",
					Name:      "regression",
					Status:    "failure",
					State:     "active",
					TopicID:   "topic-abc",
					CreatedAt: "2026-01-02T00:00:00.000000",
				},
			},
		},
	}

	outputFormat = OutputFormatStdout
	defer func() { outputFormat = OutputFormatStdout }()

	assert.NotPanics(t, func() {
		err := runGeneralJobsOutput(responses)
		assert.NoError(t, err)
	})
}

func TestRunGeneralJobsOutput_JSON(t *testing.T) {
	responses := []lib.JobsResponse{
		{
			Jobs: []lib.Job{
				{
					ID:      "job-1",
					Name:    "nightly",
					Status:  "success",
					State:   "active",
					TopicID: "topic-abc",
				},
			},
		},
	}

	outputFormat = OutputFormatJSON
	defer func() { outputFormat = OutputFormatStdout }()

	assert.NotPanics(t, func() {
		err := runGeneralJobsOutput(responses)
		assert.NoError(t, err)
	})
}

func TestRunGeneralJobsOutput_TopicName(t *testing.T) {
	responses := []lib.JobsResponse{
		{
			Jobs: []lib.Job{
				{
					ID:      "job-1",
					TopicID: "topic-abc",
				},
			},
		},
	}

	outputFormat = OutputFormatStdout
	defer func() { outputFormat = OutputFormatStdout }()

	assert.NotPanics(t, func() {
		err := runGeneralJobsOutput(responses)
		assert.NoError(t, err)
	})
}

func TestRunCertsuiteJobsOutput_Stdout(t *testing.T) {
	responses := []lib.JobsResponse{
		{
			Jobs: []lib.Job{
				{
					ID:        "job-1",
					CreatedAt: time.Now().Add(-24 * time.Hour).Format(dateFormat),
					Components: []lib.Components{
						{Name: "certsuite v1.2.3"},
					},
				},
				{
					ID:        "job-2",
					CreatedAt: time.Now().Add(-48 * time.Hour).Format(dateFormat),
					Components: []lib.Components{
						{Name: "some-other-component"},
					},
				},
			},
		},
	}

	outputFormat = OutputFormatStdout
	defer func() { outputFormat = OutputFormatStdout }()

	assert.NotPanics(t, func() {
		err := runCertsuiteJobsOutput(responses, time.Now())
		assert.NoError(t, err)
	})
}

func TestRunCertsuiteJobsOutput_JSON(t *testing.T) {
	responses := []lib.JobsResponse{
		{
			Jobs: []lib.Job{
				{
					ID:        "job-1",
					CreatedAt: time.Now().Add(-24 * time.Hour).Format(dateFormat),
					Components: []lib.Components{
						{Name: "cnf-certification-test v4.0.0"},
					},
				},
			},
		},
	}

	outputFormat = OutputFormatJSON
	defer func() { outputFormat = OutputFormatStdout }()

	assert.NotPanics(t, func() {
		err := runCertsuiteJobsOutput(responses, time.Now())
		assert.NoError(t, err)
	})
}

func TestGetTopicsCmd_Success(t *testing.T) {
	var gotQuery string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(lib.TopicsResponse{Topics: []lib.Topic{{ID: "topic-1", Name: "OCP-4.18"}}})
	})

	nameFilter = "OCP-4.18"
	t.Cleanup(func() { nameFilter = "" })

	assert.NoError(t, runCommand(getTopicsCmd))
	assert.Contains(t, gotQuery, "OCP-4.18")
}

func TestGetTopicsCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.TopicsResponse{Topics: []lib.Topic{{ID: "topic-1", Name: "OCP-4.18"}}})
	})

	nameFilter = ""
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(getTopicsCmd))
}

func TestGetTopicsCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	assert.Error(t, runCommand(getTopicsCmd))
}

func TestGetComponentTypesCmd_Success(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.ComponentTypesResponse{ComponentTypes: []lib.ComponentType{{ID: "ct-1", Name: "ocp"}}})
	})

	nameFilter = "ocp"
	t.Cleanup(func() { nameFilter = "" })

	assert.NoError(t, runCommand(getComponentTypesCmd))
}

func TestGetComponentTypesCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.ComponentTypesResponse{ComponentTypes: []lib.ComponentType{{ID: "ct-1"}}})
	})

	nameFilter = ""
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(getComponentTypesCmd))
}

func TestGetComponentsCmd_Filters(t *testing.T) {
	var gotQuery string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(lib.ComponentsResponse{Components: []lib.Components{{ID: "comp-1", Name: "OCP 4.18"}}})
	})

	topicID = "550e8400-e29b-41d4-a716-446655440000"
	typeFilter = "ocp"
	nameFilter = "4.18"
	t.Cleanup(func() { topicID, typeFilter, nameFilter = "", "", "" })

	assert.NoError(t, runCommand(getComponentsCmd))
	assert.Contains(t, gotQuery, "topic_id")
	assert.Contains(t, gotQuery, "type")
}

func TestGetComponentsCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.ComponentsResponse{Components: []lib.Components{{ID: "comp-1"}}})
	})

	topicID, typeFilter, nameFilter = "", "", ""
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(getComponentsCmd))
}

func TestGetComponentsCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	assert.Error(t, runCommand(getComponentsCmd))
}

func jobsFixture() []lib.JobsResponse {
	return []lib.JobsResponse{{Jobs: []lib.Job{
		{ID: "job-1", Name: "j1", Status: "success", State: "active", TopicID: "topic-1", CreatedAt: "2024-01-01T00:00:00.000000"},
	}}}
}

func TestGetJobsCmd_ByAge(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jobsFixture()[0])
	})

	ageInDays = "7"
	startDate, endDate = "", ""
	certsuiteFilter = false
	t.Cleanup(func() { ageInDays = "" })

	assert.NoError(t, runCommand(getJobsCmd))
}

func TestGetJobsCmd_ByDateRange(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jobsFixture()[0])
	})

	ageInDays = ""
	startDate, endDate = "2024-01-01", "2024-01-31"
	certsuiteFilter = false
	t.Cleanup(func() { startDate, endDate = "", "" })

	assert.NoError(t, runCommand(getJobsCmd))
}

func TestGetJobsCmd_Certsuite(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.JobsResponse{Jobs: []lib.Job{{
			ID:        "job-1",
			CreatedAt: "2024-01-01T00:00:00.000000",
			Components: []lib.Components{
				{Name: "cnf-certification-test 1.0.0"},
				{Name: "OpenShift 4.18"},
			},
		}}})
	})

	ageInDays = "7"
	startDate, endDate = "", ""
	certsuiteFilter = true
	t.Cleanup(func() { ageInDays = ""; certsuiteFilter = false })

	assert.NoError(t, runCommand(getJobsCmd))
}

func TestGetJobsCmd_ValidationErrors(t *testing.T) {
	tests := []struct {
		name            string
		age, start, end string
		wantContains    string
	}{
		{"age and dates", "7", "2024-01-01", "2024-01-31", "mutually exclusive"},
		{"nothing", "", "", "", "either --age or --start-date"},
		{"only start", "", "2024-01-01", "", "both --start-date and --end-date"},
		{"only end", "", "", "2024-01-31", "both --start-date and --end-date"},
		{"bad age", "seven", "", "", "invalid age value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ageInDays, startDate, endDate = tt.age, tt.start, tt.end
			certsuiteFilter = false
			t.Cleanup(func() { ageInDays, startDate, endDate = "", "", "" })

			err := runCommand(getJobsCmd)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantContains)
		})
	}
}

func TestGetJobsCmd_InvalidDates(t *testing.T) {
	startDate, endDate = "01-01-2024", "2024-01-31"
	certsuiteFilter = false
	t.Cleanup(func() { startDate, endDate = "", "" })

	err := runCommand(getJobsCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start-date")
}

func TestGetJobsCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	ageInDays = "7"
	startDate, endDate = "", ""
	certsuiteFilter = false
	t.Cleanup(func() { ageInDays = "" })

	assert.Error(t, runCommand(getJobsCmd))
}

func TestGetOcpCountCmd_Success(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.JobsResponse{Jobs: []lib.Job{{
			ID:        "job-1",
			CreatedAt: "2024-01-01T00:00:00.000000",
			Components: []lib.Components{
				{Name: "cnf-certification-test 1.0.0"},
				{Name: "OpenShift 4.18"},
			},
		}}})
	})

	ageInDays = "7"
	t.Cleanup(func() { ageInDays = "" })

	assert.NoError(t, runCommand(getOcpCountCmd))
}

func TestGetOcpCountCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jobsFixture()[0])
	})

	ageInDays = "7"
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { ageInDays = "" })

	assert.NoError(t, runCommand(getOcpCountCmd))
}

func TestGetOcpCountCmd_InvalidAge(t *testing.T) {
	ageInDays = "yesterday"
	t.Cleanup(func() { ageInDays = "" })

	err := runCommand(getOcpCountCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid age value")
}

func TestGetOcpCountCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	ageInDays = "7"
	t.Cleanup(func() { ageInDays = "" })

	assert.Error(t, runCommand(getOcpCountCmd))
}
func TestGetOCPVersions(t *testing.T) {
	t.Run("defaults when unset", func(t *testing.T) {
		t.Setenv("OCP_VERSIONS_TO_TRACK", "")
		assert.Equal(t, []string{"4.12", "4.13", "4.14", "4.15", "4.16", "4.17", "4.18", "4.19", "4.20"}, getOCPVersions())
	})

	t.Run("parses and trims comma separated list", func(t *testing.T) {
		t.Setenv("OCP_VERSIONS_TO_TRACK", "4.16, 4.18 ,4.20")
		assert.Equal(t, []string{"4.16", "4.18", "4.20"}, getOCPVersions())
	})

	t.Run("single value", func(t *testing.T) {
		t.Setenv("OCP_VERSIONS_TO_TRACK", "4.18")
		assert.Equal(t, []string{"4.18"}, getOCPVersions())
	})
}

func TestPrintOcpVersionCount(t *testing.T) {
	first := ocpVersionsToLookFor[0]
	last := ocpVersionsToLookFor[len(ocpVersionsToLookFor)-1]

	out := captureStdout(t, func() {
		printOcpVersionCount(map[string]int{first: 3})
	})

	assert.Contains(t, out, "OCP Version: "+first+" - Run Count: 3")
	assert.NotContains(t, out, "OCP Version: "+last)
}

func TestPrintOcpVersionCountJSON(t *testing.T) {
	first := ocpVersionsToLookFor[0]

	out := captureStdout(t, func() {
		if err := printOcpVersionCountJSON(map[string]int{first: 3}); err != nil {
			t.Fatalf("printOcpVersionCountJSON returned an error: %v", err)
		}
	})

	var got lib.OcpJsonOutput
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("unmarshal ocp count JSON: %v\n%s", err, out)
	}
	if len(got.OcpVersions) != len(ocpVersionsToLookFor) {
		t.Fatalf("got %d versions, want %d", len(got.OcpVersions), len(ocpVersionsToLookFor))
	}
	assert.Equal(t, first, got.OcpVersions[0].OcpVersion)
	assert.Equal(t, 3, got.OcpVersions[0].RunCount)
	assert.Equal(t, 0, got.OcpVersions[1].RunCount)
}
