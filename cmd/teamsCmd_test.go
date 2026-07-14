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

func TestPrintTeamsStdout(t *testing.T) {
	response := &lib.TeamsResponse{
		Meta: lib.Meta{Count: 2},
		Teams: []lib.Team{
			{
				ID:        "team-abc-123",
				Name:      "Red Hat Partner Team",
				Country:   "US",
				External:  true,
				State:     "active",
				CreatedAt: "2024-01-05T12:00:00.000000",
				UpdatedAt: "2024-07-22T09:30:00.000000",
			},
			{
				ID:        "team-def-456",
				Name:      "Internal QE Team",
				Country:   "CZ",
				External:  false,
				State:     "active",
				CreatedAt: "2024-02-18T15:45:00.000000",
				UpdatedAt: "2024-08-10T11:00:00.000000",
			},
		},
	}

	assert.NotPanics(t, func() {
		printTeamsStdout([]lib.TeamsResponse{*response})
	})
}

func TestPrintTeamsStdout_Empty(t *testing.T) {
	response := &lib.TeamsResponse{
		Meta:  lib.Meta{Count: 0},
		Teams: []lib.Team{},
	}

	assert.NotPanics(t, func() {
		printTeamsStdout([]lib.TeamsResponse{*response})
	})
}

func TestPrintTeamsJSON(t *testing.T) {
	response := &lib.TeamsResponse{
		Meta: lib.Meta{Count: 1},
		Teams: []lib.Team{
			{
				ID:       "team-abc-123",
				Name:     "Red Hat Partner Team",
				External: true,
				State:    "active",
			},
		},
	}

	assert.NotPanics(t, func() {
		printTeamsJSON([]lib.TeamsResponse{*response})
	})
}

func TestPrintTeamStdout(t *testing.T) {
	response := &lib.TeamResponse{
		Team: lib.Team{
			ID:        "team-abc-123",
			Name:      "Red Hat Partner Team",
			Country:   "US",
			External:  true,
			State:     "active",
			CreatedAt: "2024-01-05T12:00:00.000000",
			UpdatedAt: "2024-07-22T09:30:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		printTeamStdout(response)
	})
}

func TestPrintTeamJSON(t *testing.T) {
	response := &lib.TeamResponse{
		Team: lib.Team{
			ID:        "team-abc-123",
			Name:      "Red Hat Partner Team",
			Country:   "US",
			External:  true,
			State:     "active",
			CreatedAt: "2024-01-05T12:00:00.000000",
			UpdatedAt: "2024-07-22T09:30:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		printTeamJSON(response)
	})
}

func TestGetTeamCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.TeamResponse{Team: lib.Team{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "test-team"}})
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getTeamIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getTeamIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getTeamCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetTeamCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getTeamIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getTeamIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getTeamCmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestDeleteTeamCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	deleteTeamIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatStdout
	defer func() { deleteTeamIDFlag = ""; yesFlag = false }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := deleteTeamCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
