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

func TestPrintUsersStdout(t *testing.T) {
	response := &lib.UsersResponse{
		Meta: lib.Meta{Count: 2},
		Users: []lib.User{
			{
				ID:        "user-abc-123",
				Name:      "jdoe",
				Fullname:  "Jane Doe",
				Email:     "jdoe@redhat.com",
				TeamID:    "team-456",
				State:     "active",
				CreatedAt: "2024-02-10T08:00:00.000000",
				UpdatedAt: "2024-09-05T17:30:00.000000",
			},
			{
				ID:        "user-def-789",
				Name:      "bsmith",
				Fullname:  "Bob Smith",
				Email:     "bsmith@partner.com",
				TeamID:    "team-012",
				State:     "active",
				CreatedAt: "2024-03-15T14:20:00.000000",
				UpdatedAt: "2024-10-01T10:00:00.000000",
			},
		},
	}

	assert.NotPanics(t, func() {
		printUsersStdout([]lib.UsersResponse{*response})
	})
}

func TestPrintUsersStdout_Empty(t *testing.T) {
	response := &lib.UsersResponse{
		Meta:  lib.Meta{Count: 0},
		Users: []lib.User{},
	}

	assert.NotPanics(t, func() {
		printUsersStdout([]lib.UsersResponse{*response})
	})
}

func TestPrintUsersJSON(t *testing.T) {
	response := &lib.UsersResponse{
		Meta: lib.Meta{Count: 1},
		Users: []lib.User{
			{
				ID:       "user-abc-123",
				Name:     "jdoe",
				Fullname: "Jane Doe",
				Email:    "jdoe@redhat.com",
				TeamID:   "team-456",
				State:    "active",
			},
		},
	}

	assert.NotPanics(t, func() {
		printUsersJSON([]lib.UsersResponse{*response})
	})
}

func TestPrintUserStdout(t *testing.T) {
	response := &lib.UserResponse{
		User: lib.User{
			ID:        "user-abc-123",
			Name:      "jdoe",
			Fullname:  "Jane Doe",
			Email:     "jdoe@redhat.com",
			TeamID:    "team-456",
			State:     "active",
			CreatedAt: "2024-02-10T08:00:00.000000",
			UpdatedAt: "2024-09-05T17:30:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		printUserStdout(response)
	})
}

func TestPrintUserJSON(t *testing.T) {
	response := &lib.UserResponse{
		User: lib.User{
			ID:        "user-abc-123",
			Name:      "jdoe",
			Fullname:  "Jane Doe",
			Email:     "jdoe@redhat.com",
			TeamID:    "team-456",
			State:     "active",
			CreatedAt: "2024-02-10T08:00:00.000000",
			UpdatedAt: "2024-09-05T17:30:00.000000",
		},
	}

	assert.NotPanics(t, func() {
		printUserJSON(response)
	})
}

func TestGetUserCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lib.UserResponse{User: lib.User{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "testuser"}})
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getUserIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getUserIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getUserCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}

func TestGetUserCmd_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	getUserIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	outputFormat = OutputFormatStdout
	defer func() { getUserIDFlag = "" }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := getUserCmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestDeleteUserCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	dciClient = lib.NewClient("test", "test")
	dciClient.BaseURL = server.URL + "/api/v1"
	defer func() { dciClient = nil }()

	deleteUserIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatStdout
	defer func() { deleteUserIDFlag = ""; yesFlag = false }()

	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := deleteUserCmd.RunE(cmd, []string{})
	assert.NoError(t, err)
}
