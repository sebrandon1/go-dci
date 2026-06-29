package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
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
