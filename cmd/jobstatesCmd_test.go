package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/stretchr/testify/assert"
)

func TestPrintJobStatesStdout(t *testing.T) {
	response := &lib.JobStatesResponse{
		Meta: lib.Meta{Count: 2},
		JobStates: []lib.JobStateEntry{
			{
				ID:        "js-1",
				JobID:     "job-123",
				Status:    "running",
				Comment:   "Job started",
				CreatedAt: "2024-01-01T00:00:00.000000",
			},
			{
				ID:        "js-2",
				JobID:     "job-123",
				Status:    "success",
				CreatedAt: "2024-01-01T01:00:00.000000",
			},
		},
	}
	assert.NotPanics(t, func() {
		printJobStatesStdout(response)
	})
}

func TestPrintJobStatesStdout_Empty(t *testing.T) {
	response := &lib.JobStatesResponse{
		Meta:      lib.Meta{Count: 0},
		JobStates: []lib.JobStateEntry{},
	}
	assert.NotPanics(t, func() {
		printJobStatesStdout(response)
	})
}

func TestPrintJobStatesJSON(t *testing.T) {
	response := &lib.JobStatesResponse{
		Meta: lib.Meta{Count: 1},
		JobStates: []lib.JobStateEntry{
			{
				ID:        "js-1",
				JobID:     "job-123",
				Status:    "success",
				Comment:   "All tests passed",
				CreatedAt: "2024-01-01T01:00:00.000000",
			},
		},
	}
	assert.NotPanics(t, func() {
		printJobStatesJSON(response)
	})
}
