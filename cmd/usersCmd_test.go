package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
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
		printUsersStdout(response)
	})
}

func TestPrintUsersStdout_Empty(t *testing.T) {
	response := &lib.UsersResponse{
		Meta:  lib.Meta{Count: 0},
		Users: []lib.User{},
	}

	assert.NotPanics(t, func() {
		printUsersStdout(response)
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
		printUsersJSON(response)
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
