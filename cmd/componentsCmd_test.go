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

func TestPrintComponentStdout(t *testing.T) {
	response := &lib.ComponentResponse{
		Component: lib.Components{
			ID:          "comp-123",
			Name:        "Test Component",
			Type:        "ocp",
			Version:     "4.14.1",
			TopicID:     "topic-456",
			State:       "active",
			DisplayName: "Test Component Display",
			Tags:        []string{"latest", "stable"},
			CreatedAt:   "2024-01-01T00:00:00.000000",
			UpdatedAt:   "2024-06-01T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printComponentStdout(response)
	})
}

func TestPrintComponentStdout_NoDisplayNameNoTags(t *testing.T) {
	response := &lib.ComponentResponse{
		Component: lib.Components{
			ID:        "comp-123",
			Name:      "Test Component",
			Type:      "ocp",
			Version:   "4.14.1",
			TopicID:   "topic-456",
			State:     "active",
			CreatedAt: "2024-01-01T00:00:00.000000",
			UpdatedAt: "2024-06-01T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printComponentStdout(response)
	})
}

func TestPrintComponentJSON(t *testing.T) {
	response := &lib.ComponentResponse{
		Component: lib.Components{
			ID:          "comp-123",
			Name:        "Test Component",
			Type:        "ocp",
			Version:     "4.14.1",
			TopicID:     "topic-456",
			State:       "active",
			DisplayName: "Test Component Display",
			Tags:        []string{"latest", "stable"},
			CreatedAt:   "2024-01-01T00:00:00.000000",
			UpdatedAt:   "2024-06-01T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printComponentJSON(response)
	})
}

func TestGetComponentCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.ComponentResponse{Component: lib.Components{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "test-component"}})
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getComponentIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getComponentIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getComponentCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetComponentCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getComponentIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getComponentIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getComponentCmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestDeleteComponentCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	deleteComponentIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatStdout
	defer func() { deleteComponentIDFlag = ""; yesFlag = false }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := deleteComponentCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
