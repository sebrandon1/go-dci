package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPrintJobStdout(t *testing.T) {
	response := &lib.JobResponse{
		Job: lib.Job{
			ID:         "job-123",
			Name:       "test-job",
			Status:     "success",
			State:      "active",
			TopicID:    "topic-456",
			RemoteciID: "remoteci-789",
			TeamID:     "team-abc",
			Comment:    "Test comment",
			Tags:       []string{"tag1", "tag2"},
			Duration:   3600,
			CreatedAt:  "2024-01-01T00:00:00.000000",
			UpdatedAt:  "2024-01-02T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printJobStdout(response)
	})
}

func TestPrintJobStdout_NoCommentNoTags(t *testing.T) {
	response := &lib.JobResponse{
		Job: lib.Job{
			ID:         "job-123",
			Name:       "test-job",
			Status:     "failure",
			State:      "active",
			TopicID:    "topic-456",
			RemoteciID: "remoteci-789",
			TeamID:     "team-abc",
			Duration:   120,
			CreatedAt:  "2024-01-01T00:00:00.000000",
			UpdatedAt:  "2024-01-02T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printJobStdout(response)
	})
}

func TestPrintJobJSON(t *testing.T) {
	response := &lib.JobResponse{
		Job: lib.Job{
			ID:         "job-123",
			Name:       "test-job",
			Status:     "success",
			State:      "active",
			TopicID:    "topic-456",
			RemoteciID: "remoteci-789",
			TeamID:     "team-abc",
			Comment:    "Test comment",
			Tags:       []string{"tag1", "tag2"},
			Duration:   3600,
			CreatedAt:  "2024-01-01T00:00:00.000000",
			UpdatedAt:  "2024-01-02T00:00:00.000000",
		},
	}
	if err := printJobJSON(response); err != nil {
		t.Fatalf("printJobJSON returned an error: %v", err)
	}
}

func TestPrintFilesStdout(t *testing.T) {
	response := &lib.FilesResponse{
		Meta: lib.Meta{Count: 2},
		Files: []lib.File{
			{
				ID:    "file-1",
				JobID: "job-123",
				Name:  "results.xml",
				Mime:  "application/xml",
				Size:  1024,
			},
			{
				ID:    "file-2",
				JobID: "job-123",
				Name:  "log.txt",
				Mime:  "text/plain",
				Size:  2048,
			},
		},
	}
	assert.NotPanics(t, func() {
		printFilesStdout(response)
	})
}

func TestPrintFilesStdout_Empty(t *testing.T) {
	response := &lib.FilesResponse{
		Meta:  lib.Meta{Count: 0},
		Files: []lib.File{},
	}
	assert.NotPanics(t, func() {
		printFilesStdout(response)
	})
}

func TestPrintFilesJSON(t *testing.T) {
	response := &lib.FilesResponse{
		Meta: lib.Meta{Count: 1},
		Files: []lib.File{
			{
				ID:    "file-1",
				JobID: "job-123",
				Name:  "results.xml",
				Mime:  "application/xml",
				Size:  1024,
			},
		},
	}
	if err := printFilesJSON(response); err != nil {
		t.Fatalf("printFilesJSON returned an error: %v", err)
	}
}

func TestGetJobCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lib.JobResponse{Job: lib.Job{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "test-job"}}); err != nil {
			t.Errorf("encode job response: %v", err)
		}
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getJobIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getJobIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getJobCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetJobCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getJobIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getJobIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getJobCmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestDeleteJobCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	deleteJobIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatStdout
	defer func() { deleteJobIDFlag = ""; yesFlag = false }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := deleteJobCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestUpdateJobCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(lib.JobResponse{Job: lib.Job{ID: "job-1", Comment: "retry", Tags: []string{"ocp"}}})
	})

	updateJobIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateJobComment = "retry"
	updateJobTags = "ocp, certsuite"
	t.Cleanup(func() { updateJobIDFlag, updateJobComment, updateJobTags = "", "", "" })

	assert.NoError(t, runCommand(updateJobCmd))
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.Equal(t, "/api/v1/jobs/550e8400-e29b-41d4-a716-446655440000", gotPath)
}

func TestUpdateJobCmd_JSONAndDryRun(t *testing.T) {
	updateJobIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateJobTags = "ocp"
	t.Cleanup(func() { updateJobIDFlag, updateJobTags = "", "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(updateJobCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.JobResponse{Job: lib.Job{ID: "job-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(updateJobCmd))
}

func TestUpdateJobCmd_InvalidID(t *testing.T) {
	updateJobIDFlag = "bad-id"
	t.Cleanup(func() { updateJobIDFlag = "" })

	err := runCommand(updateJobCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid job ID format")
}

func TestScheduleJobCmd_Success(t *testing.T) {
	var gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.CreateJobResponse{Job: lib.Job{ID: "job-1", TopicID: "topic-1"}})
	})

	scheduleJobTopicID = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { scheduleJobTopicID = "" })

	assert.NoError(t, runCommand(scheduleJobCmd))
	assert.Equal(t, "/api/v1/jobs/schedule", gotPath)
}

func TestScheduleJobCmd_JSONAndDryRun(t *testing.T) {
	scheduleJobTopicID = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { scheduleJobTopicID = "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(scheduleJobCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.CreateJobResponse{Job: lib.Job{ID: "job-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(scheduleJobCmd))
}

func TestScheduleJobCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	scheduleJobTopicID = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { scheduleJobTopicID = "" })

	assert.Error(t, runCommand(scheduleJobCmd))
}

func TestGetJobFilesCmd_Success(t *testing.T) {
	var gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(lib.FilesResponse{Files: []lib.File{{ID: "file-1", Name: "results.xml"}}})
	})

	jobFilesIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { jobFilesIDFlag = "" })

	assert.NoError(t, runCommand(getJobFilesCmd))
	assert.Equal(t, "/api/v1/jobs/550e8400-e29b-41d4-a716-446655440000/files", gotPath)
}

func TestGetJobFilesCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.FilesResponse{Files: []lib.File{{ID: "file-1"}}})
	})

	jobFilesIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { jobFilesIDFlag = "" })

	assert.NoError(t, runCommand(getJobFilesCmd))
}

func TestDeleteJobCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	deleteJobIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { deleteJobIDFlag = "" })

	assert.NoError(t, runCommand(deleteJobCmd))
}

func TestDeleteJobCmd_DryRunSkipsAPI(t *testing.T) {
	deleteJobIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	dryRunFlag = true
	t.Cleanup(func() { deleteJobIDFlag = ""; resetGlobalFlags() })

	assert.NoError(t, runCommand(deleteJobCmd))
}

func TestDeleteJobCmd_InvalidID(t *testing.T) {
	deleteJobIDFlag = "nope"
	t.Cleanup(func() { deleteJobIDFlag = "" })

	err := runCommand(deleteJobCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid job ID format")
}
