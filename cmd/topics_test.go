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

func TestPrintTopicStdout(t *testing.T) {
	response := &lib.TopicResponse{
		Topic: lib.Topic{
			ID:             "topic-123",
			Name:           "OCP-4.14",
			ProductID:      "prod-456",
			State:          "active",
			ExportControl:  true,
			ComponentTypes: []string{"ocp", "certsuite"},
			CreatedAt:      "2024-01-01T00:00:00.000000",
			UpdatedAt:      "2024-06-01T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printTopicStdout(response)
	})
}

func TestPrintTopicStdout_NoComponentTypes(t *testing.T) {
	response := &lib.TopicResponse{
		Topic: lib.Topic{
			ID:            "topic-123",
			Name:          "OCP-4.14",
			ProductID:     "prod-456",
			State:         "active",
			ExportControl: false,
			CreatedAt:     "2024-01-01T00:00:00.000000",
			UpdatedAt:     "2024-06-01T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printTopicStdout(response)
	})
}

func TestPrintTopicJSON(t *testing.T) {
	response := &lib.TopicResponse{
		Topic: lib.Topic{
			ID:             "topic-123",
			Name:           "OCP-4.14",
			ProductID:      "prod-456",
			State:          "active",
			ExportControl:  true,
			ComponentTypes: []string{"ocp", "certsuite"},
			CreatedAt:      "2024-01-01T00:00:00.000000",
			UpdatedAt:      "2024-06-01T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printTopicJSON(response)
	})
}

func TestGetTopicCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.TopicResponse{Topic: lib.Topic{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "OCP-4.14"}})
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getTopicIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getTopicIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getTopicCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetTopicCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getTopicIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getTopicIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getTopicCmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestDeleteTopicCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	deleteTopicIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatStdout
	defer func() { deleteTopicIDFlag = ""; yesFlag = false }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := deleteTopicCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
