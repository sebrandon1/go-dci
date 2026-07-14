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

func TestPrintRemoteCIsStdout(t *testing.T) {
	responses := []lib.RemoteCIsResponse{
		{
			Meta: lib.Meta{Count: 2},
			RemoteCIs: []lib.RemoteCI{
				{
					ID:        "remoteci-abc-123",
					Name:      "partner-lab-remoteci",
					TeamID:    "team-456",
					State:     "active",
					CreatedAt: "2024-03-10T09:00:00.000000",
					UpdatedAt: "2024-08-15T16:30:00.000000",
				},
				{
					ID:        "remoteci-def-789",
					Name:      "staging-remoteci",
					TeamID:    "team-012",
					State:     "active",
					CreatedAt: "2024-04-20T11:15:00.000000",
					UpdatedAt: "2024-09-01T08:45:00.000000",
				},
			},
		},
	}

	assert.NotPanics(t, func() {
		printRemoteCIsStdout(responses)
	})
}

func TestPrintRemoteCIsStdout_Empty(t *testing.T) {
	responses := []lib.RemoteCIsResponse{
		{
			Meta:      lib.Meta{Count: 0},
			RemoteCIs: []lib.RemoteCI{},
		},
	}

	assert.NotPanics(t, func() {
		printRemoteCIsStdout(responses)
	})
}

func TestPrintRemoteCIsJSON(t *testing.T) {
	responses := []lib.RemoteCIsResponse{
		{
			Meta: lib.Meta{Count: 1},
			RemoteCIs: []lib.RemoteCI{
				{
					ID:     "remoteci-abc-123",
					Name:   "partner-lab-remoteci",
					TeamID: "team-456",
					State:  "active",
				},
			},
		},
	}

	assert.NotPanics(t, func() {
		err := printRemoteCIsJSON(responses)
		assert.NoError(t, err)
	})
}

func TestPrintRemoteCIStdout(t *testing.T) {
	response := &lib.RemoteCIResponse{
		RemoteCI: lib.RemoteCI{
			ID:        "remoteci-abc-123",
			Name:      "partner-lab-remoteci",
			TeamID:    "team-456",
			State:     "active",
			CreatedAt: "2024-03-10T09:00:00.000000",
			UpdatedAt: "2024-08-15T16:30:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		printRemoteCIStdout(response)
	})
}

func TestPrintRemoteCIJSON(t *testing.T) {
	response := &lib.RemoteCIResponse{
		RemoteCI: lib.RemoteCI{
			ID:        "remoteci-abc-123",
			Name:      "partner-lab-remoteci",
			TeamID:    "team-456",
			State:     "active",
			CreatedAt: "2024-03-10T09:00:00.000000",
			UpdatedAt: "2024-08-15T16:30:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		err := printRemoteCIJSON(response)
		assert.NoError(t, err)
	})
}

func TestGetRemoteCICmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.RemoteCIResponse{RemoteCI: lib.RemoteCI{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "test-remoteci"}})
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getRemoteCIsCmd_IDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getRemoteCIsCmd_IDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getRemoteCICmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetRemoteCICmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getRemoteCIsCmd_IDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getRemoteCIsCmd_IDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getRemoteCICmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestDeleteRemoteCICmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	deleteRemoteCIIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatStdout
	defer func() { deleteRemoteCIIDFlag = ""; yesFlag = false }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := deleteRemoteCICmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
