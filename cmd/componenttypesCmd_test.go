package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
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
