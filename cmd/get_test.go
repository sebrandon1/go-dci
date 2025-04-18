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
