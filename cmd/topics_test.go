package cmd

import (
	"testing"

	"github.com/sebrandon1/go-dci/lib"
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
	assert.NotPanics(t, func() {
		printTopicJSON(response)
	})
}
