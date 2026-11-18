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
	if err := printComponentTypeJSON(response); err != nil {
		t.Fatalf("printComponentTypeJSON returned an error: %v", err)
	}
}

func TestGetComponentTypeCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lib.ComponentTypeResponse{ComponentType: lib.ComponentType{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "ocp"}}); err != nil {
			t.Errorf("encode component type response: %v", err)
		}
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

func TestCreateComponentTypeCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.ComponentTypeResponse{ComponentType: lib.ComponentType{ID: "ct-1", Name: "ocp"}})
	})

	createComponentTypeName = "ocp"
	t.Cleanup(func() { createComponentTypeName = "" })

	assert.NoError(t, runCommand(createComponentTypeCmd))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v1/componenttypes", gotPath)
}

func TestCreateComponentTypeCmd_JSONAndDryRun(t *testing.T) {
	createComponentTypeName = "ocp"
	t.Cleanup(func() { createComponentTypeName = "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(createComponentTypeCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.ComponentTypeResponse{ComponentType: lib.ComponentType{ID: "ct-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(createComponentTypeCmd))
}

func TestUpdateComponentTypeCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(lib.ComponentTypeResponse{ComponentType: lib.ComponentType{ID: "ct-1", Name: "ocp-stable"}})
	})

	updateComponentTypeIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateComponentTypeName = "ocp-stable"
	updateComponentTypeState = "active"
	t.Cleanup(func() {
		updateComponentTypeIDFlag, updateComponentTypeName, updateComponentTypeState = "", "", ""
	})

	assert.NoError(t, runCommand(updateComponentTypeCmd))
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.Equal(t, "/api/v1/componenttypes/550e8400-e29b-41d4-a716-446655440000", gotPath)
}

func TestUpdateComponentTypeCmd_JSONAndDryRun(t *testing.T) {
	updateComponentTypeIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateComponentTypeState = "inactive"
	t.Cleanup(func() { updateComponentTypeIDFlag, updateComponentTypeState = "", "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(updateComponentTypeCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.ComponentTypeResponse{ComponentType: lib.ComponentType{ID: "ct-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(updateComponentTypeCmd))
}

func TestUpdateComponentTypeCmd_InvalidID(t *testing.T) {
	updateComponentTypeIDFlag = "bad"
	t.Cleanup(func() { updateComponentTypeIDFlag = "" })

	err := runCommand(updateComponentTypeCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid component type ID format")
}

func TestDeleteComponentTypeCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	deleteComponentTypeIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { deleteComponentTypeIDFlag = "" })

	assert.NoError(t, runCommand(deleteComponentTypeCmd))
}
