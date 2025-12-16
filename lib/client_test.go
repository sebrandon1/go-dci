package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("testAccessKey", "testSecretKey")

	assert.NotNil(t, client)
	assert.Equal(t, "testAccessKey", client.AccessKey)
	assert.Equal(t, "testSecretKey", client.SecretKey)
	assert.Equal(t, DCIURL, client.BaseURL)
}

func TestGetComponents_EmptyResponse(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path
		assert.Contains(t, r.URL.Path, "/components")

		// Return an empty components response
		response := ComponentsResponse{
			Meta:       Meta{Count: 0},
			Components: []Components{},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

	components, err := client.GetComponents()
	assert.NoError(t, err)
	assert.NotNil(t, components)
	assert.Len(t, components, 1)
	assert.Empty(t, components[0].Components)
}

func TestGetComponents_WithData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ComponentsResponse{
			Meta: Meta{Count: 2},
			Components: []Components{
				{
					ID:      "comp-1",
					Name:    "OpenShift 4.14",
					Type:    "ocp",
					Version: "4.14.1",
					TopicID: "topic-123",
				},
				{
					ID:      "comp-2",
					Name:    "certsuite v1.0.0",
					Type:    "certsuite",
					Version: "v1.0.0",
					TopicID: "topic-123",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

	components, err := client.GetComponents()
	assert.NoError(t, err)
	assert.NotNil(t, components)
	assert.Len(t, components, 1)
	assert.Len(t, components[0].Components, 2)
	assert.Equal(t, "comp-1", components[0].Components[0].ID)
	assert.Equal(t, "OpenShift 4.14", components[0].Components[0].Name)
}

func TestGetComponentsByTopicID(t *testing.T) {
	expectedTopicID := "topic-456"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the topic_id filter is in the query
		whereParam := r.URL.Query().Get("where")
		assert.Contains(t, whereParam, "topic_id:"+expectedTopicID)

		response := ComponentsResponse{
			Meta: Meta{Count: 1},
			Components: []Components{
				{
					ID:      "comp-filtered",
					Name:    "Filtered Component",
					Type:    "test",
					Version: "1.0.0",
					TopicID: expectedTopicID,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

	components, err := client.GetComponentsByTopicID(expectedTopicID)
	assert.NoError(t, err)
	assert.NotNil(t, components)
	assert.Len(t, components, 1)
	assert.Len(t, components[0].Components, 1)
	assert.Equal(t, expectedTopicID, components[0].Components[0].TopicID)
}

func TestFetchComponents_Error(t *testing.T) {
	// Create a server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

	components, err := client.fetchComponents("", 100, 0)
	// The error handling depends on how the API returns errors
	// For an empty body with 500, json decode will fail
	assert.Error(t, err)
	assert.Empty(t, components.Components)
}

func TestFetchComponents_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("invalid json"))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

	components, err := client.fetchComponents("", 100, 0)
	assert.Error(t, err)
	assert.Empty(t, components.Components)
}

func TestComponentsResponse_Struct(t *testing.T) {
	// Test JSON marshaling/unmarshaling of ComponentsResponse
	original := ComponentsResponse{
		Meta: Meta{Count: 1},
		Components: []Components{
			{
				ID:                   "test-id",
				Name:                 "Test Component",
				Type:                 "test-type",
				Version:              "1.0.0",
				TopicID:              "topic-id",
				CanonicalProjectName: "test-project",
				DisplayName:          "Test Display Name",
				State:                "active",
			},
		},
	}

	jsonBytes, err := json.Marshal(original)
	assert.NoError(t, err)

	var decoded ComponentsResponse
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original.Meta.Count, decoded.Meta.Count)
	assert.Len(t, decoded.Components, 1)
	assert.Equal(t, original.Components[0].ID, decoded.Components[0].ID)
	assert.Equal(t, original.Components[0].Name, decoded.Components[0].Name)
	assert.Equal(t, original.Components[0].Type, decoded.Components[0].Type)
}

