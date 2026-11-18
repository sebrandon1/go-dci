package cmd

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/stretchr/testify/assert"
)

func TestPrintCreateJobStdout(t *testing.T) {
	response := &lib.CreateJobResponse{
		Job: lib.Job{
			ID:        "job-123",
			TopicID:   "topic-456",
			Status:    "new",
			State:     "active",
			CreatedAt: "2024-01-01T00:00:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		printCreateJobStdout(response)
	})
}

func TestPrintCreateJobJSON(t *testing.T) {
	response := &lib.CreateJobResponse{
		Job: lib.Job{
			ID:        "job-123",
			TopicID:   "topic-456",
			Status:    "new",
			State:     "active",
			CreatedAt: "2024-01-01T00:00:00.000000",
		},
	}

	if err := printCreateJobJSON(response); err != nil {
		t.Fatalf("printCreateJobJSON returned an error: %v", err)
	}
}

func TestPrintJobStateStdout(t *testing.T) {
	response := &lib.JobStateResponse{}
	response.JobState.ID = "jobstate-123"
	response.JobState.JobID = "job-456"
	response.JobState.Status = "running"
	response.JobState.Comment = "test comment"
	response.JobState.CreatedAt = "2024-01-01T00:00:00.000000"

	assert.NotPanics(t, func() {
		printJobStateStdout(response)
	})
}

func TestPrintJobStateStdout_NoComment(t *testing.T) {
	response := &lib.JobStateResponse{}
	response.JobState.ID = "jobstate-123"
	response.JobState.JobID = "job-456"
	response.JobState.Status = "success"
	response.JobState.CreatedAt = "2024-01-01T00:00:00.000000"

	assert.NotPanics(t, func() {
		printJobStateStdout(response)
	})
}

func TestPrintJobStateJSON(t *testing.T) {
	response := &lib.JobStateResponse{}
	response.JobState.ID = "jobstate-123"
	response.JobState.JobID = "job-456"
	response.JobState.Status = "running"
	response.JobState.CreatedAt = "2024-01-01T00:00:00.000000"

	if err := printJobStateJSON(response); err != nil {
		t.Fatalf("printJobStateJSON returned an error: %v", err)
	}
}

func TestPrintUploadFileStdout(t *testing.T) {
	response := &lib.UploadFileResponse{}
	response.File.ID = "file-123"
	response.File.JobID = "job-456"
	response.File.Name = "test-results.xml"
	response.File.Mime = "application/junit"
	response.File.Size = "1024"
	response.File.CreatedAt = "2024-01-01T00:00:00.000000"

	assert.NotPanics(t, func() {
		printUploadFileStdout(response)
	})
}

func TestPrintUploadFileJSON(t *testing.T) {
	response := &lib.UploadFileResponse{}
	response.File.ID = "file-123"
	response.File.JobID = "job-456"
	response.File.Name = "test-results.xml"
	response.File.Mime = "application/junit"
	response.File.Size = "1024"
	response.File.CreatedAt = "2024-01-01T00:00:00.000000"

	if err := printUploadFileJSON(response); err != nil {
		t.Fatalf("printUploadFileJSON returned an error: %v", err)
	}
}

func TestCreateJobCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.CreateJobResponse{Job: lib.Job{ID: "job-1", TopicID: "topic-1"}})
	})

	createJobTopicID = "550e8400-e29b-41d4-a716-446655440000"
	createJobComponents = "comp-1, comp-2"
	createJobComment = "nightly"
	t.Cleanup(func() { createJobTopicID, createJobComponents, createJobComment = "", "", "" })

	assert.NoError(t, runCommand(createJobCmd))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v1/jobs", gotPath)
}

func TestCreateJobCmd_JSONAndDryRun(t *testing.T) {
	createJobTopicID = "550e8400-e29b-41d4-a716-446655440000"
	createJobComponents = "comp-1"
	t.Cleanup(func() { createJobTopicID, createJobComponents = "", "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(createJobCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.CreateJobResponse{Job: lib.Job{ID: "job-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(createJobCmd))
}

func TestCreateJobCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	createJobTopicID = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { createJobTopicID = "" })

	assert.Error(t, runCommand(createJobCmd))
}

func TestUpdateJobStateCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jobstate": map[string]string{"id": "js-1", "job_id": "job-1", "status": "running"},
		})
	})

	updateJobStateJobID = "550e8400-e29b-41d4-a716-446655440000"
	updateJobStateStatus = "running"
	updateJobStateComment = "started"
	t.Cleanup(func() { updateJobStateJobID, updateJobStateStatus, updateJobStateComment = "", "", "" })

	assert.NoError(t, runCommand(updateJobStateCmd))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v1/jobstates", gotPath)
}

func TestUpdateJobStateCmd_JSONAndDryRun(t *testing.T) {
	updateJobStateJobID = "550e8400-e29b-41d4-a716-446655440000"
	updateJobStateStatus = "success"
	t.Cleanup(func() { updateJobStateJobID, updateJobStateStatus = "", "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(updateJobStateCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jobstate": map[string]string{"id": "js-1", "job_id": "job-1", "status": "success"},
		})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(updateJobStateCmd))
}

func TestUpdateJobStateCmd_InvalidStatus(t *testing.T) {
	updateJobStateJobID = "550e8400-e29b-41d4-a716-446655440000"
	updateJobStateStatus = "exploded"
	t.Cleanup(func() { updateJobStateJobID, updateJobStateStatus = "", "" })

	err := runCommand(updateJobStateCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
}

func TestUpdateJobStateCmd_InvalidJobID(t *testing.T) {
	updateJobStateJobID = "bad"
	updateJobStateStatus = "running"
	t.Cleanup(func() { updateJobStateJobID, updateJobStateStatus = "", "" })

	err := runCommand(updateJobStateCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid job ID format")
}

func TestUploadFileCmd_Success(t *testing.T) {
	var gotPath, gotJobID, gotMime string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotJobID = r.Header.Get("DCI-JOB-ID")
		gotMime = r.Header.Get("DCI-MIME")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"file": map[string]string{"id": "file-1", "job_id": "job-1", "name": "results.xml", "mime": "application/junit", "size": "12"},
		})
	})

	uploadFilePath = filepath.Join(t.TempDir(), "results.xml")
	if err := os.WriteFile(uploadFilePath, []byte("<testsuite/>"), 0600); err != nil {
		t.Fatalf("write upload fixture: %v", err)
	}
	uploadFileJobID = "550e8400-e29b-41d4-a716-446655440000"
	uploadFileMimeType = ""
	t.Cleanup(func() { uploadFilePath, uploadFileJobID, uploadFileMimeType = "", "", "" })

	assert.NoError(t, runCommand(uploadFileCmd))
	assert.Equal(t, "/api/v1/files", gotPath)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", gotJobID)
	assert.Equal(t, "application/junit", gotMime)
}

func TestUploadFileCmd_JSONAndDryRun(t *testing.T) {
	uploadFilePath = filepath.Join(t.TempDir(), "results.xml")
	if err := os.WriteFile(uploadFilePath, []byte("<testsuite/>"), 0600); err != nil {
		t.Fatalf("write upload fixture: %v", err)
	}
	uploadFileJobID = "550e8400-e29b-41d4-a716-446655440000"
	uploadFileMimeType = "text/xml"
	t.Cleanup(func() { uploadFileJobID, uploadFilePath, uploadFileMimeType = "", "", "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(uploadFileCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"file": map[string]string{"id": "file-1", "job_id": "job-1", "name": "results.xml"},
		})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(uploadFileCmd))
}

func TestUploadFileCmd_MissingFile(t *testing.T) {
	uploadFileJobID = "550e8400-e29b-41d4-a716-446655440000"
	uploadFilePath = filepath.Join(t.TempDir(), "absent.xml")
	t.Cleanup(func() { uploadFileJobID, uploadFilePath = "", "" })

	assert.Error(t, runCommand(uploadFileCmd))
}
