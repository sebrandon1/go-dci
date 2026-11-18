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

	if err := printUsersJSON([]lib.UsersResponse{*response}); err != nil {
		t.Fatalf("printUsersJSON returned an error: %v", err)
	}
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

	if err := printUserJSON(response); err != nil {
		t.Fatalf("printUserJSON returned an error: %v", err)
	}
}

func TestGetUserCmd_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lib.UserResponse{User: lib.User{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "testuser"}}); err != nil {
			t.Errorf("encode user response: %v", err)
		}
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

func TestGetUsersCmd_Success(t *testing.T) {
	var gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewEncoder(w).Encode(lib.UsersResponse{Users: []lib.User{{ID: "user-1", Name: "alice"}}})
	})

	usersNameFilter = "alice"
	t.Cleanup(func() { usersNameFilter = "" })

	assert.NoError(t, runCommand(getUsersCmd))
	assert.Equal(t, "/api/v1/users", gotPath)
}

func TestGetUsersCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.UsersResponse{Users: []lib.User{{ID: "user-1"}}})
	})

	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(getUsersCmd))
}

func TestGetUsersCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	assert.Error(t, runCommand(getUsersCmd))
}

func TestCreateUserCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.UserResponse{User: lib.User{ID: "user-1", Name: "alice"}})
	})

	createUserNameFlag = "alice"
	createUserEmailFlag = "alice@example.com"
	createUserTeamIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	createUserPasswordFlag = "s3cr3t"
	t.Cleanup(func() {
		createUserNameFlag, createUserEmailFlag = "", ""
		createUserTeamIDFlag, createUserPasswordFlag = "", ""
	})

	assert.NoError(t, runCommand(createUserCmd))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/v1/users", gotPath)
}

func TestCreateUserCmd_JSONAndDryRun(t *testing.T) {
	createUserNameFlag = "alice"
	createUserEmailFlag = "alice@example.com"
	createUserTeamIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	createUserPasswordFlag = "s3cr3t"
	t.Cleanup(func() {
		createUserNameFlag, createUserEmailFlag = "", ""
		createUserTeamIDFlag, createUserPasswordFlag = "", ""
	})

	dryRunFlag = true
	assert.NoError(t, runCommand(createUserCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lib.UserResponse{User: lib.User{ID: "user-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(createUserCmd))
}

func TestCreateUserCmd_RequiresPasswordWhenNotInteractive(t *testing.T) {
	createUserNameFlag = "alice"
	createUserPasswordFlag = ""
	t.Cleanup(func() { createUserNameFlag, createUserPasswordFlag = "", "" })

	err := runCommand(createUserCmd)
	if err == nil {
		t.Skip("stdin is a terminal in this environment")
	}
	assert.Contains(t, err.Error(), "password is required")
}

func TestCreateUserCmd_ServerError(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	})

	createUserNameFlag = "alice"
	createUserEmailFlag = "alice@example.com"
	createUserTeamIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	createUserPasswordFlag = "s3cr3t"
	t.Cleanup(func() {
		createUserNameFlag, createUserEmailFlag = "", ""
		createUserTeamIDFlag, createUserPasswordFlag = "", ""
	})

	assert.Error(t, runCommand(createUserCmd))
}

func TestUpdateUserCmd_Success(t *testing.T) {
	var gotMethod, gotPath string
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		_ = json.NewEncoder(w).Encode(lib.UserResponse{User: lib.User{ID: "user-1", Email: "new@example.com"}})
	})

	updateUserIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateUserEmailFlag = "new@example.com"
	updateUserStateFlag = "active"
	t.Cleanup(func() { updateUserIDFlag, updateUserEmailFlag, updateUserStateFlag = "", "", "" })

	assert.NoError(t, runCommand(updateUserCmd))
	assert.Equal(t, http.MethodPut, gotMethod)
	assert.Equal(t, "/api/v1/users/550e8400-e29b-41d4-a716-446655440000", gotPath)
}

func TestUpdateUserCmd_JSONAndDryRun(t *testing.T) {
	updateUserIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	updateUserNameFlag = "alice2"
	updateUserFullnameFlag = "Alice Smith"
	t.Cleanup(func() { updateUserIDFlag, updateUserNameFlag, updateUserFullnameFlag = "", "", "" })

	dryRunFlag = true
	assert.NoError(t, runCommand(updateUserCmd))

	dryRunFlag = false
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(lib.UserResponse{User: lib.User{ID: "user-1"}})
	})
	outputFormat = OutputFormatJSON
	assert.NoError(t, runCommand(updateUserCmd))
}

func TestUpdateUserCmd_InvalidID(t *testing.T) {
	updateUserIDFlag = "nope"
	t.Cleanup(func() { updateUserIDFlag = "" })

	err := runCommand(updateUserCmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID format")
}

func TestDeleteUserCmd_JSONOutput(t *testing.T) {
	newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	deleteUserIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	yesFlag = true
	outputFormat = OutputFormatJSON
	t.Cleanup(func() { deleteUserIDFlag = "" })

	assert.NoError(t, runCommand(deleteUserCmd))
}

func TestDeleteUserCmd_DryRunSkipsAPI(t *testing.T) {
	deleteUserIDFlag = "550e8400-e29b-41d4-a716-446655440000"
	dryRunFlag = true
	t.Cleanup(func() { deleteUserIDFlag = ""; resetGlobalFlags() })

	assert.NoError(t, runCommand(deleteUserCmd))
}
