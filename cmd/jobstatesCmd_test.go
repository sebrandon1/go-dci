package cmd

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/stretchr/testify/assert"
)

func TestPrintJobStatesStdout(t *testing.T) {
	responses := []lib.JobStatesResponse{
		{
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
		},
	}
	assert.NotPanics(t, func() {
		printJobStatesStdout(responses)
	})
}

func TestPrintJobStatesStdout_Empty(t *testing.T) {
	responses := []lib.JobStatesResponse{
		{
			Meta:      lib.Meta{Count: 0},
			JobStates: []lib.JobStateEntry{},
		},
	}
	assert.NotPanics(t, func() {
		printJobStatesStdout(responses)
	})
}

func TestPrintJobStatesJSON(t *testing.T) {
	responses := []lib.JobStatesResponse{
		{
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
		},
	}
	if err := printJobStatesJSON(responses); err != nil {
		t.Fatalf("printJobStatesJSON returned an error: %v", err)
	}
}

func TestGetJobStatesCmd_FilteredByJob(t *testing.T) {
	var gotQuery string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(lib.JobStatesResponse{JobStates: []lib.JobStateEntry{{ID: "js-1", JobID: "job-1", Status: "running"}}})
	})

	getJobStatesJobIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { getJobStatesJobIDFlag = "" })

	assert.NoError(t, runCommand(getJobStatesCmd))
	assert.Contains(t, gotQuery, "job_id")
}

func TestGetJobStatesCmd_AllJobStates(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.JobStatesResponse{JobStates: []lib.JobStateEntry{{ID: "js-1"}}})
	})

	assert.NoError(t, runCommand(getJobStatesCmd))
}

func TestGetJobStatesCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.JobStatesResponse{JobStates: []lib.JobStateEntry{{ID: "js-1"}}})
	})

	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(getJobStatesCmd))
}

func TestGetJobStatesCmd_InvalidJobID(t *testing.T) {
	getJobStatesJobIDFlag = "bad"
	t.Cleanup(func() { getJobStatesJobIDFlag = "" })

	err := runCommand(getJobStatesCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid job ID format")
}

func TestGetJobStatesCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	assert.Error(t, runCommand(getJobStatesCmd))
}
