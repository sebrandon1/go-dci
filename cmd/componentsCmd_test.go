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
	if err := printComponentJSON(response); err != nil {
		t.Fatalf("printComponentJSON returned an error: %v", err)
	}
}

func TestGetComponentCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lib.ComponentResponse{Component: lib.Components{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "test-component"}}); err != nil {
			t.Errorf("encode component response: %v", err)
		}
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

func TestCreateComponentCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.ComponentResponse{Component: lib.Components{ID: "comp-1", Name: "OCP 4.18", Tags: []string{"ga"}}})
	})

	createComponentName = "OCP 4.18"
	createComponentType = "ocp"
	createComponentTopicID = "550e8400-e29b-41d4-a716-446655440000"
	createComponentVersion = "4.18.0"
	t.Cleanup(func() {
		createComponentName, createComponentType = "", ""
		createComponentTopicID, createComponentVersion = "", ""
	})

	assert.NoError(t, runCommand(createComponentCmd))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v1/components", gotPath)
}

func TestCreateComponentCmd_JSONAndDryRun(t *testing.T) {
	createComponentName = "OCP 4.18"
	createComponentType = "ocp"
	createComponentTopicID = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() {
		createComponentName, createComponentType, createComponentTopicID = "", "", ""
	})

	dryRunFlag = true
	assert.NoError(t, runCommand(createComponentCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.ComponentResponse{Component: lib.Components{ID: "comp-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(createComponentCmd))
}

func TestUpdateComponentCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(lib.ComponentResponse{Component: lib.Components{ID: "comp-1"}})
	})

	updateComponentIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateComponentName = "OCP 4.18.1"
	updateComponentState = "active"
	updateComponentVersion = "4.18.1"
	updateComponentTags = "ga, stable"
	t.Cleanup(func() {
		updateComponentIDFlag, updateComponentName = "", ""
		updateComponentState, updateComponentVersion, updateComponentTags = "", "", ""
	})

	assert.NoError(t, runCommand(updateComponentCmd))
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.Equal(t, "/api/v1/components/550e8400-e29b-41d4-a716-446655440000", gotPath)
}

func TestUpdateComponentCmd_JSONAndDryRun(t *testing.T) {
	updateComponentIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateComponentTags = "ga"
	t.Cleanup(func() { updateComponentIDFlag, updateComponentTags = "", "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(updateComponentCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.ComponentResponse{Component: lib.Components{ID: "comp-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(updateComponentCmd))
}

func TestUpdateComponentCmd_InvalidID(t *testing.T) {
	updateComponentIDFlag = "bad"
	t.Cleanup(func() { updateComponentIDFlag = "" })

	err := runCommand(updateComponentCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid component ID format")
}

func TestDeleteComponentCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	deleteComponentIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { deleteComponentIDFlag = "" })

	assert.NoError(t, runCommand(deleteComponentCmd))
}
