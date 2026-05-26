package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
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
