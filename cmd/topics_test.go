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
	if err := printTopicJSON(response); err != nil {
		t.Fatalf("printTopicJSON returned an error: %v", err)
	}
}

func TestGetTopicCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lib.TopicResponse{Topic: lib.Topic{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "OCP-4.14"}}); err != nil {
			t.Errorf("encode topic response: %v", err)
		}
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

func TestCreateTopicCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.TopicResponse{Topic: lib.Topic{ID: "topic-1", Name: "OCP-4.18", ComponentTypes: []string{"ocp"}}})
	})

	createTopicName = "OCP-4.18"
	createTopicProductID = "550e8400-e29b-41d4-a716-446655440000"
	createTopicComponentTypes = "ocp, rhcos"
	t.Cleanup(func() { createTopicName, createTopicProductID, createTopicComponentTypes = "", "", "" })

	assert.NoError(t, runCommand(createTopicCmd))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v1/topics", gotPath)
}

func TestCreateTopicCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.TopicResponse{Topic: lib.Topic{ID: "topic-1", Name: "OCP-4.18"}})
	})

	createTopicName = "OCP-4.18"
	createTopicProductID = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { createTopicName, createTopicProductID = "", "" })

	assert.NoError(t, runCommand(createTopicCmd))
}

func TestCreateTopicCmd_DryRunSkipsAPI(t *testing.T) {
	createTopicName = "OCP-4.18"
	createTopicProductID = "550e8400-e29b-41d4-a716-446655440000"
	createTopicComponentTypes = "ocp,rhcos"
	dryRunFlag = true
	t.Cleanup(func() {
		createTopicName, createTopicProductID, createTopicComponentTypes = "", "", ""
		resetGlobalFlags()
	})

	assert.NoError(t, runCommand(createTopicCmd))
}

func TestCreateTopicCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	createTopicName = "OCP-4.18"
	createTopicProductID = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { createTopicName, createTopicProductID = "", "" })

	assert.Error(t, runCommand(createTopicCmd))
}

func TestUpdateTopicCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(lib.TopicResponse{Topic: lib.Topic{ID: "topic-1", Name: "OCP-4.18-GA"}})
	})

	updateTopicIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateTopicName = "OCP-4.18-GA"
	t.Cleanup(func() { updateTopicIDFlag, updateTopicName = "", "" })

	assert.NoError(t, runCommand(updateTopicCmd))
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.Equal(t, "/api/v1/topics/550e8400-e29b-41d4-a716-446655440000", gotPath)
}

func TestUpdateTopicCmd_JSONAndDryRun(t *testing.T) {
	updateTopicIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateTopicName = "OCP-4.18-GA"
	t.Cleanup(func() { updateTopicIDFlag, updateTopicName = "", "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(updateTopicCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.TopicResponse{Topic: lib.Topic{ID: "topic-1", Name: "OCP-4.18-GA"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(updateTopicCmd))
}

func TestUpdateTopicCmd_InvalidID(t *testing.T) {
	updateTopicIDFlag = "not-a-uuid"
	t.Cleanup(func() { updateTopicIDFlag = "" })

	err := runCommand(updateTopicCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid topic ID format")
}

func TestGetTopicComponentsCmd_Success(t *testing.T) {
	var gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(lib.ComponentsResponse{Components: []lib.Components{{ID: "comp-1", Name: "OCP 4.18"}}})
	})

	topicComponentsIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { topicComponentsIDFlag = "" })

	assert.NoError(t, runCommand(getTopicComponentsCmd))
	assert.Equal(t, "/api/v1/topics/550e8400-e29b-41d4-a716-446655440000/components", gotPath)
}

func TestGetTopicComponentsCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.ComponentsResponse{Components: []lib.Components{{ID: "comp-1"}}})
	})

	topicComponentsIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { topicComponentsIDFlag = "" })

	assert.NoError(t, runCommand(getTopicComponentsCmd))
}

func TestGetTopicComponentsCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	topicComponentsIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	t.Cleanup(func() { topicComponentsIDFlag = "" })

	assert.Error(t, runCommand(getTopicComponentsCmd))
}
