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
	assert.NotPanics(t, func() {
		printJobJSON(response)
	})
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
	assert.NotPanics(t, func() {
		printFilesJSON(response)
	})
}

func TestGetJobCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.JobResponse{Job: lib.Job{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "test-job"}})
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
