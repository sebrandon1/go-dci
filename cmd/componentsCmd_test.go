package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
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
	assert.NotPanics(t, func() {
		printComponentJSON(response)
	})
}
