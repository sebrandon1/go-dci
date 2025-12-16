package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
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

	// Verify the function doesn't panic with valid input
	assert.NotPanics(t, func() {
		printComponentsJSON(componentsResponses)
	})
}

func TestPrintComponentsJSON_EmptyComponents(t *testing.T) {
	componentsResponses := []lib.ComponentsResponse{
		{
			Meta:       lib.Meta{Count: 0},
			Components: []lib.Components{},
		},
	}

	assert.NotPanics(t, func() {
		printComponentsJSON(componentsResponses)
	})
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

	assert.NotPanics(t, func() {
		printComponentsJSON(componentsResponses)
	})
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

	assert.NotPanics(t, func() {
		printIdentityJSON(identity)
	})
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

	assert.NotPanics(t, func() {
		printComponentTypesJSON(componentTypesResponses)
	})
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

	assert.NotPanics(t, func() {
		printComponentTypesJSON(componentTypesResponses)
	})
}
