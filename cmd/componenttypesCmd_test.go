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

func TestPrintComponentTypeStdout(t *testing.T) {
	response := &lib.ComponentTypeResponse{
		ComponentType: lib.ComponentType{
			ID:        "ct-123",
			Name:      "ocp",
			State:     "active",
			CreatedAt: "2024-01-01T00:00:00.000000",
			UpdatedAt: "2024-06-01T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printComponentTypeStdout(response)
	})
}

func TestPrintComponentTypeJSON(t *testing.T) {
	response := &lib.ComponentTypeResponse{
		ComponentType: lib.ComponentType{
			ID:        "ct-123",
			Name:      "ocp",
			State:     "active",
			CreatedAt: "2024-01-01T00:00:00.000000",
			UpdatedAt: "2024-06-01T00:00:00.000000",
		},
	}
	assert.NotPanics(t, func() {
		printComponentTypeJSON(response)
	})
}

func TestGetComponentTypeCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.ComponentTypeResponse{ComponentType: lib.ComponentType{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "ocp"}})
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getComponentTypeIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getComponentTypeIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getComponentTypeCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetComponentTypeCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getComponentTypeIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getComponentTypeIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getComponentTypeCmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestDeleteComponentTypeCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	deleteComponentTypeIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatStdout
	defer func() { deleteComponentTypeIDFlag = ""; yesFlag = false }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := deleteComponentTypeCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
