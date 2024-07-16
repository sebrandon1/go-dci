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
