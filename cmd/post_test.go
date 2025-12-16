package cmd

import (
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

	assert.NotPanics(t, func() {
		printCreateJobJSON(response)
	})
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

	assert.NotPanics(t, func() {
		printJobStateJSON(response)
	})
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

	assert.NotPanics(t, func() {
		printUploadFileJSON(response)
	})
}

