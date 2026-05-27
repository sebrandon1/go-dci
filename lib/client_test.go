package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestClient(serverURL string) *Client {
	return &Client{BaseURL: serverURL, AccessKey: "testKey", SecretKey: "testSecret", httpClient: &http.Client{}}
}

func TestNewClient(t *testing.T) {
	client := NewClient("testAccessKey", "testSecretKey")

	assert.NotNil(t, client)
	assert.Equal(t, "testAccessKey", client.AccessKey)
	assert.Equal(t, "testSecretKey", client.SecretKey)
	assert.Equal(t, DCIURL, client.BaseURL)
	assert.NotNil(t, client.httpClient)
}

func TestGetComponents_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/components")

		response := ComponentsResponse{
			Meta:       Meta{Count: 0},
			Components: []Components{},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

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

	client := newTestClient(server.URL)

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

	client := newTestClient(server.URL)

	components, err := client.GetComponentsByTopicID(expectedTopicID)
	assert.NoError(t, err)
	assert.NotNil(t, components)
	assert.Len(t, components, 1)
	assert.Len(t, components[0].Components, 1)
	assert.Equal(t, expectedTopicID, components[0].Components[0].TopicID)
}

func TestFetchComponents_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	components, err := client.fetchComponents("", 100, 0)
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

	client := newTestClient(server.URL)

	components, err := client.fetchComponents("", 100, 0)
	assert.Error(t, err)
	assert.Empty(t, components.Components)
}

func TestComponentsResponse_Struct(t *testing.T) {
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

func TestGetIdentity_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/identity", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := IdentityResponse{
			Identity: Identity{
				ID:       "remoteci-123",
				Name:     "test-remoteci",
				Type:     "remoteci",
				TeamID:   "team-456",
				TeamName: "Test Team",
				State:    "active",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	identity, err := client.GetIdentity()
	assert.NoError(t, err)
	assert.NotNil(t, identity)
	assert.Equal(t, "remoteci-123", identity.Identity.ID)
	assert.Equal(t, "test-remoteci", identity.Identity.Name)
	assert.Equal(t, "remoteci", identity.Identity.Type)
	assert.Equal(t, "team-456", identity.Identity.TeamID)
	assert.Equal(t, "Test Team", identity.Identity.TeamName)
	assert.Equal(t, "active", identity.Identity.State)
}

func TestGetIdentity_UserType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := IdentityResponse{
			Identity: Identity{
				ID:       "user-789",
				Name:     "testuser",
				Type:     "user",
				Email:    "testuser@example.com",
				Fullname: "Test User",
				TeamID:   "team-456",
				TeamName: "Test Team",
				State:    "active",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	identity, err := client.GetIdentity()
	assert.NoError(t, err)
	assert.NotNil(t, identity)
	assert.Equal(t, "user", identity.Identity.Type)
	assert.Equal(t, "testuser@example.com", identity.Identity.Email)
	assert.Equal(t, "Test User", identity.Identity.Fullname)
}

func TestGetIdentity_AuthenticationFailed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		AccessKey:  "badKey",
		SecretKey:  "badSecret",
		httpClient: &http.Client{},
	}

	identity, err := client.GetIdentity()
	assert.Error(t, err)
	assert.Nil(t, identity)
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestGetIdentity_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("invalid json"))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	identity, err := client.GetIdentity()
	assert.Error(t, err)
	assert.Nil(t, identity)
}

func TestIdentityResponse_Struct(t *testing.T) {
	original := IdentityResponse{
		Identity: Identity{
			ID:        "test-id",
			Name:      "test-name",
			Type:      "remoteci",
			Email:     "test@example.com",
			Etag:      "etag-123",
			Fullname:  "Test Fullname",
			State:     "active",
			TeamID:    "team-id",
			TeamName:  "Team Name",
			Timezone:  "UTC",
			CreatedAt: "2024-01-01T00:00:00.000000",
			UpdatedAt: "2024-01-02T00:00:00.000000",
		},
	}

	jsonBytes, err := json.Marshal(original)
	assert.NoError(t, err)

	var decoded IdentityResponse
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original.Identity.ID, decoded.Identity.ID)
	assert.Equal(t, original.Identity.Name, decoded.Identity.Name)
	assert.Equal(t, original.Identity.Type, decoded.Identity.Type)
	assert.Equal(t, original.Identity.Email, decoded.Identity.Email)
	assert.Equal(t, original.Identity.TeamName, decoded.Identity.TeamName)
}

func TestGetComponentTypes_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/componenttypes")

		response := ComponentTypesResponse{
			Meta:           Meta{Count: 0},
			ComponentTypes: []ComponentType{},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	componentTypes, err := client.GetComponentTypes()
	assert.NoError(t, err)
	assert.NotNil(t, componentTypes)
	assert.Len(t, componentTypes, 1)
	assert.Empty(t, componentTypes[0].ComponentTypes)
}

func TestGetComponentTypes_WithData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/componenttypes", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := ComponentTypesResponse{
			Meta: Meta{Count: 3},
			ComponentTypes: []ComponentType{
				{
					ID:    "ct-1",
					Name:  "ocp",
					State: "active",
				},
				{
					ID:    "ct-2",
					Name:  "certsuite",
					State: "active",
				},
				{
					ID:    "ct-3",
					Name:  "rhel",
					State: "active",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	componentTypes, err := client.GetComponentTypes()
	assert.NoError(t, err)
	assert.NotNil(t, componentTypes)
	assert.Len(t, componentTypes, 1)
	assert.Len(t, componentTypes[0].ComponentTypes, 3)
	assert.Equal(t, "ct-1", componentTypes[0].ComponentTypes[0].ID)
	assert.Equal(t, "ocp", componentTypes[0].ComponentTypes[0].Name)
	assert.Equal(t, "certsuite", componentTypes[0].ComponentTypes[1].Name)
}

func TestFetchComponentTypes_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	componentTypes, err := client.fetchComponentTypes(100, 0)
	assert.Error(t, err)
	assert.Empty(t, componentTypes.ComponentTypes)
}

func TestFetchComponentTypes_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("invalid json"))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	componentTypes, err := client.fetchComponentTypes(100, 0)
	assert.Error(t, err)
	assert.Empty(t, componentTypes.ComponentTypes)
}

func TestComponentTypesResponse_Struct(t *testing.T) {
	original := ComponentTypesResponse{
		Meta: Meta{Count: 2},
		ComponentTypes: []ComponentType{
			{
				ID:        "ct-1",
				Name:      "ocp",
				Etag:      "etag-123",
				State:     "active",
				CreatedAt: "2024-01-01T00:00:00.000000",
				UpdatedAt: "2024-01-02T00:00:00.000000",
			},
			{
				ID:    "ct-2",
				Name:  "certsuite",
				State: "active",
			},
		},
	}

	jsonBytes, err := json.Marshal(original)
	assert.NoError(t, err)

	var decoded ComponentTypesResponse
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original.Meta.Count, decoded.Meta.Count)
	assert.Len(t, decoded.ComponentTypes, 2)
	assert.Equal(t, original.ComponentTypes[0].ID, decoded.ComponentTypes[0].ID)
	assert.Equal(t, original.ComponentTypes[0].Name, decoded.ComponentTypes[0].Name)
	assert.Equal(t, original.ComponentTypes[1].Name, decoded.ComponentTypes[1].Name)
}

func TestCreateJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/jobs", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var reqBody CreateJobRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "topic-123", reqBody.TopicID)

		response := CreateJobResponse{
			Job: Job{
				ID:        "job-456",
				TopicID:   "topic-123",
				Status:    "new",
				State:     "active",
				CreatedAt: "2024-01-01T00:00:00.000000",
			},
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	response, err := client.CreateJob("topic-123", []string{"comp-1", "comp-2"}, "test comment")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "job-456", response.Job.ID)
	assert.Equal(t, "topic-123", response.Job.TopicID)
	assert.Equal(t, "new", response.Job.Status)
}

func TestCreateJob_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid topic_id"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	response, err := client.CreateJob("invalid-topic", nil, "")
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to create job")
}

func TestUpdateJobState_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/jobstates", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var reqBody UpdateJobStateRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "job-123", reqBody.JobID)
		assert.Equal(t, "running", reqBody.Status)

		response := JobStateResponse{}
		response.JobState.ID = "jobstate-789"
		response.JobState.JobID = "job-123"
		response.JobState.Status = "running"
		response.JobState.CreatedAt = "2024-01-01T00:00:00.000000"

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	response, err := client.UpdateJobState("job-123", JobStateRunning, "starting test run")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "jobstate-789", response.JobState.ID)
	assert.Equal(t, "job-123", response.JobState.JobID)
	assert.Equal(t, "running", response.JobState.Status)
}

func TestUpdateJobState_ToSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody UpdateJobStateRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "success", reqBody.Status)

		response := JobStateResponse{}
		response.JobState.ID = "jobstate-101"
		response.JobState.JobID = reqBody.JobID
		response.JobState.Status = "success"

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	response, err := client.UpdateJobState("job-123", JobStateSuccess, "")
	assert.NoError(t, err)
	assert.Equal(t, "success", response.JobState.Status)
}

func TestUpdateJobState_ToFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody UpdateJobStateRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "failure", reqBody.Status)

		response := JobStateResponse{}
		response.JobState.ID = "jobstate-102"
		response.JobState.JobID = reqBody.JobID
		response.JobState.Status = "failure"
		response.JobState.Comment = reqBody.Comment

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	response, err := client.UpdateJobState("job-123", JobStateFailure, "tests failed")
	assert.NoError(t, err)
	assert.Equal(t, "failure", response.JobState.Status)
}

func TestUpdateJobState_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "job not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	response, err := client.UpdateJobState("nonexistent-job", JobStateRunning, "")
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to update job state")
}

func TestUploadFileContent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/files", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "job-123", r.Header.Get("DCI-JOB-ID"))
		assert.Equal(t, "test-results.xml", r.Header.Get("DCI-NAME"))
		assert.Equal(t, "application/junit", r.Header.Get("DCI-MIME"))

		response := UploadFileResponse{}
		response.File.ID = "file-456"
		response.File.JobID = "job-123"
		response.File.Name = "test-results.xml"
		response.File.Mime = "application/junit"
		response.File.Size = "1024"
		response.File.CreatedAt = "2024-01-01T00:00:00.000000"

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	content := []byte("<testsuites><testsuite></testsuite></testsuites>")
	response, err := client.UploadFileContent("job-123", "test-results.xml", "application/junit", content)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "file-456", response.File.ID)
	assert.Equal(t, "job-123", response.File.JobID)
	assert.Equal(t, "test-results.xml", response.File.Name)
	assert.Equal(t, "application/junit", response.File.Mime)
}

func TestUploadFileContent_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid job_id"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	response, err := client.UploadFileContent("invalid-job", "test.xml", "application/junit", []byte("content"))
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to upload file")
}

func TestCreateJobRequest_Struct(t *testing.T) {
	original := CreateJobRequest{
		TopicID:    "topic-123",
		Components: []string{"comp-1", "comp-2"},
		Comment:    "test comment",
	}

	jsonBytes, err := json.Marshal(original)
	assert.NoError(t, err)

	var decoded CreateJobRequest
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original.TopicID, decoded.TopicID)
	assert.Equal(t, original.Components, decoded.Components)
	assert.Equal(t, original.Comment, decoded.Comment)
}

func TestJobStateConstants(t *testing.T) {
	assert.Equal(t, JobState("new"), JobStateNew)
	assert.Equal(t, JobState("pre-run"), JobStatePreRun)
	assert.Equal(t, JobState("running"), JobStateRunning)
	assert.Equal(t, JobState("post-run"), JobStatePostRun)
	assert.Equal(t, JobState("success"), JobStateSuccess)
	assert.Equal(t, JobState("failure"), JobStateFailure)
	assert.Equal(t, JobState("killed"), JobStateKilled)
	assert.Equal(t, JobState("error"), JobStateError)
}

func TestGetTopics_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/topics")
		assert.Equal(t, "GET", r.Method)

		response := TopicsResponse{
			Meta: Meta{Count: 1},
			Topics: []struct {
				ComponentTypes         []string `json:"component_types,omitempty"`
				ComponentTypesOptional []any    `json:"component_types_optional,omitempty"`
				CreatedAt              string   `json:"created_at,omitempty"`
				Data                   struct {
				} `json:"data,omitempty"`
				Etag          string  `json:"etag,omitempty"`
				ExportControl bool    `json:"export_control,omitempty"`
				ID            string  `json:"id,omitempty"`
				Name          string  `json:"name,omitempty"`
				NextTopic     any     `json:"next_topic,omitempty"`
				NextTopicID   any     `json:"next_topic_id,omitempty"`
				Product       Product `json:"product,omitempty"`
				ProductID     string  `json:"product_id,omitempty"`
				State         string  `json:"state,omitempty"`
				UpdatedAt     string  `json:"updated_at,omitempty"`
			}{
				{
					ID:        "topic-1",
					Name:      "OCP-4.14",
					ProductID: "prod-1",
					State:     "active",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	topics, err := client.GetTopics()
	assert.NoError(t, err)
	assert.NotNil(t, topics)
	assert.Len(t, topics, 1)
	assert.Len(t, topics[0].Topics, 1)
	assert.Equal(t, "topic-1", topics[0].Topics[0].ID)
	assert.Equal(t, "OCP-4.14", topics[0].Topics[0].Name)
}

func TestGetTopics_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	topics, err := client.GetTopics()
	assert.Error(t, err)
	assert.Nil(t, topics)
}

func TestGetTopic_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/topics/topic-123")
		assert.Equal(t, "GET", r.Method)

		response := TopicResponse{Topic: Topic{ID: "topic-123", Name: "OCP-4.14", ProductID: "prod-1"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTopic("topic-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "topic-123", result.Topic.ID)
	assert.Equal(t, "OCP-4.14", result.Topic.Name)
}

func TestGetTopic_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTopic("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get topic")
}

func TestCreateTopic_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqBody CreateTopicRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "OCP-4.15", reqBody.Name)
		assert.Equal(t, "prod-1", reqBody.ProductID)
		assert.Equal(t, []string{"ocp", "certsuite"}, reqBody.ComponentTypes)

		response := TopicResponse{Topic: Topic{ID: "new-topic", Name: "OCP-4.15", ProductID: "prod-1"}}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateTopic("OCP-4.15", "prod-1", []string{"ocp", "certsuite"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-topic", result.Topic.ID)
}

func TestCreateTopic_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateTopic("bad", "bad", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create topic")
}

func TestUpdateTopic_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Contains(t, r.URL.Path, "/topics/topic-123")

		response := TopicResponse{Topic: Topic{ID: "topic-123", Name: "updated-topic"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateTopic("topic-123", UpdateTopicRequest{Name: "updated-topic"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "updated-topic", result.Topic.Name)
}

func TestUpdateTopic_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateTopic("nonexistent", UpdateTopicRequest{Name: "x"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update topic")
}

func TestDeleteTopic_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/topics/topic-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteTopic("topic-123")
	assert.NoError(t, err)
}

func TestDeleteTopic_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteTopic("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete topic")
}

func TestGetTopicComponents_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/topics/topic-123/components")
		assert.Equal(t, "GET", r.Method)

		response := ComponentsResponse{
			Meta: Meta{Count: 1},
			Components: []Components{
				{ID: "comp-1", Name: "OCP 4.14", TopicID: "topic-123"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTopicComponents("topic-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Len(t, result[0].Components, 1)
	assert.Equal(t, "comp-1", result[0].Components[0].ID)
}

func TestGetTopicComponents_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTopicComponents("topic-123")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFetchTopics_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TopicsResponse{
			Meta: Meta{Count: 1},
			Topics: []struct {
				ComponentTypes         []string `json:"component_types,omitempty"`
				ComponentTypesOptional []any    `json:"component_types_optional,omitempty"`
				CreatedAt              string   `json:"created_at,omitempty"`
				Data                   struct {
				} `json:"data,omitempty"`
				Etag          string  `json:"etag,omitempty"`
				ExportControl bool    `json:"export_control,omitempty"`
				ID            string  `json:"id,omitempty"`
				Name          string  `json:"name,omitempty"`
				NextTopic     any     `json:"next_topic,omitempty"`
				NextTopicID   any     `json:"next_topic_id,omitempty"`
				Product       Product `json:"product,omitempty"`
				ProductID     string  `json:"product_id,omitempty"`
				State         string  `json:"state,omitempty"`
				UpdatedAt     string  `json:"updated_at,omitempty"`
			}{
				{ID: "topic-1", Name: "test"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchTopics(100, 0)
	assert.NoError(t, err)
	assert.Len(t, result.Topics, 1)
	assert.Equal(t, "topic-1", result.Topics[0].ID)
}

func TestFetchTopics_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchTopics(100, 0)
	assert.Error(t, err)
	assert.Empty(t, result.Topics)
}

func TestFetchTopics_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("invalid json"))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchTopics(100, 0)
	assert.Error(t, err)
	assert.Empty(t, result.Topics)
}

func TestFetchTopicComponents_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ComponentsResponse{
			Meta:       Meta{Count: 1},
			Components: []Components{{ID: "comp-1", Name: "test-comp"}},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchTopicComponents("topic-123", 100, 0)
	assert.NoError(t, err)
	assert.Len(t, result.Components, 1)
	assert.Equal(t, "comp-1", result.Components[0].ID)
}

func TestFetchTopicComponents_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchTopicComponents("topic-123", 100, 0)
	assert.Error(t, err)
	assert.Empty(t, result.Components)
}

func TestFetchTopicComponents_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("invalid json"))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchTopicComponents("topic-123", 100, 0)
	assert.Error(t, err)
	assert.Empty(t, result.Components)
}

func TestGetJobs_Success(t *testing.T) {
	recentTimestamp := time.Now().Format("2006-01-02T15:04:05.999999")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/jobs")

		response := JobsResponse{
			Meta: Meta{Count: 1},
			Jobs: []Job{
				{
					ID:        "job-1",
					CreatedAt: recentTimestamp,
					Status:    "success",
					State:     "active",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	jobs, err := client.GetJobs(7)
	assert.NoError(t, err)
	assert.NotNil(t, jobs)
	assert.Len(t, jobs, 1)
	assert.Len(t, jobs[0].Jobs, 1)
	assert.Equal(t, "job-1", jobs[0].Jobs[0].ID)
}

func TestGetJobs_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	jobs, err := client.GetJobs(7)
	assert.Error(t, err)
	assert.Nil(t, jobs)
}

func TestGetJobsByDate_Success(t *testing.T) {
	// Use UTC because time.Parse with dateFormat produces UTC times
	now := time.Now().UTC()
	jobCreatedAt := now.Format("2006-01-02T15:04:05.999999")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := JobsResponse{
			Meta: Meta{Count: 1},
			Jobs: []Job{
				{
					ID:        "job-date-1",
					CreatedAt: jobCreatedAt,
					Status:    "success",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	jobs, err := client.GetJobsByDate(now.Add(-time.Hour), now.Add(time.Hour))
	assert.NoError(t, err)
	assert.NotNil(t, jobs)
	assert.Len(t, jobs, 1)
	assert.Equal(t, "job-date-1", jobs[0].Jobs[0].ID)
}

func TestGetJobsByDate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	jobs, err := client.GetJobsByDate(time.Now(), time.Now().Add(time.Hour))
	assert.Error(t, err)
	assert.Nil(t, jobs)
}

func TestGetJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/jobs/job-123")
		assert.Equal(t, "GET", r.Method)

		response := JobResponse{Job: Job{ID: "job-123", Status: "success", State: "active"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetJob("job-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "job-123", result.Job.ID)
	assert.Equal(t, "success", result.Job.Status)
}

func TestGetJob_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetJob("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get job")
}

func TestUpdateJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Contains(t, r.URL.Path, "/jobs/job-123")

		response := JobResponse{Job: Job{ID: "job-123", Comment: "updated comment"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateJob("job-123", UpdateJobRequest{Comment: "updated comment"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "updated comment", result.Job.Comment)
}

func TestUpdateJob_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateJob("nonexistent", UpdateJobRequest{Comment: "x"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update job")
}

func TestDeleteJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/jobs/job-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteJob("job-123")
	assert.NoError(t, err)
}

func TestDeleteJob_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteJob("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete job")
}

func TestScheduleJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Path, "/jobs/schedule")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqBody ScheduleJobRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "topic-123", reqBody.TopicID)

		response := CreateJobResponse{Job: Job{ID: "scheduled-job", TopicID: "topic-123", Status: "new"}}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.ScheduleJob("topic-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "scheduled-job", result.Job.ID)
}

func TestScheduleJob_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "bad request"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.ScheduleJob("bad-topic")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to schedule job")
}

func TestGetJobFiles_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/jobs/job-123/files")
		assert.Equal(t, "GET", r.Method)

		response := FilesResponse{
			Meta: Meta{Count: 1},
			Files: []File{
				{ID: "file-1", JobID: "job-123", Name: "results.xml", Mime: "application/xml"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetJobFiles("job-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Files, 1)
	assert.Equal(t, "file-1", result.Files[0].ID)
	assert.Equal(t, "results.xml", result.Files[0].Name)
}

func TestGetJobFiles_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetJobFiles("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get job files")
}

func TestGetJobStates_WithJobID_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/jobstates")
		assert.Contains(t, r.URL.RawQuery, "where=job_id%3Ajob-123")

		response := JobStatesResponse{
			Meta: Meta{Count: 1},
			JobStates: []JobStateEntry{
				{ID: "js-1", JobID: "job-123", Status: "running"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetJobStates("job-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.JobStates, 1)
	assert.Equal(t, "js-1", result.JobStates[0].ID)
}

func TestGetJobStates_EmptyJobID_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/jobstates")
		assert.NotContains(t, r.URL.RawQuery, "where=")

		response := JobStatesResponse{
			Meta: Meta{Count: 2},
			JobStates: []JobStateEntry{
				{ID: "js-1", JobID: "job-1", Status: "running"},
				{ID: "js-2", JobID: "job-2", Status: "success"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetJobStates("")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.JobStates, 2)
}

func TestGetJobStates_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "server error"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetJobStates("job-123")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get job states")
}

func TestFetchJobs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := JobsResponse{
			Meta: Meta{Count: 1},
			Jobs: []Job{{ID: "job-1", Status: "success"}},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchJobs(100, 0)
	assert.NoError(t, err)
	assert.Len(t, result.Jobs, 1)
	assert.Equal(t, "job-1", result.Jobs[0].ID)
}

func TestFetchJobs_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchJobs(100, 0)
	assert.Error(t, err)
	assert.Empty(t, result.Jobs)
}

func TestFetchJobs_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("invalid json"))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.fetchJobs(100, 0)
	assert.Error(t, err)
	assert.Empty(t, result.Jobs)
}

func TestGetComponent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/components/comp-123")
		assert.Equal(t, "GET", r.Method)

		response := ComponentResponse{Component: Components{ID: "comp-123", Name: "OCP 4.14", Type: "ocp", Version: "4.14.1"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetComponent("comp-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "comp-123", result.Component.ID)
	assert.Equal(t, "OCP 4.14", result.Component.Name)
}

func TestGetComponent_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetComponent("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get component")
}

func TestCreateComponent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqBody CreateComponentRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "OCP 4.15", reqBody.Name)
		assert.Equal(t, "ocp", reqBody.Type)
		assert.Equal(t, "topic-123", reqBody.TopicID)
		assert.Equal(t, "4.15.0", reqBody.Version)

		response := ComponentResponse{Component: Components{ID: "new-comp", Name: "OCP 4.15", Type: "ocp"}}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateComponent("OCP 4.15", "ocp", "topic-123", "4.15.0")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-comp", result.Component.ID)
}

func TestCreateComponent_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateComponent("bad", "bad", "bad", "bad")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create component")
}

func TestUpdateComponent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Contains(t, r.URL.Path, "/components/comp-123")

		response := ComponentResponse{Component: Components{ID: "comp-123", Name: "updated-comp"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateComponent("comp-123", UpdateComponentRequest{Name: "updated-comp"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "updated-comp", result.Component.Name)
}

func TestUpdateComponent_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateComponent("nonexistent", UpdateComponentRequest{Name: "x"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update component")
}

func TestDeleteComponent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/components/comp-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteComponent("comp-123")
	assert.NoError(t, err)
}

func TestDeleteComponent_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteComponent("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete component")
}

func TestGetComponentType_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/componenttypes/ct-123")
		assert.Equal(t, "GET", r.Method)

		response := ComponentTypeResponse{ComponentType: ComponentType{ID: "ct-123", Name: "ocp", State: "active"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetComponentType("ct-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ct-123", result.ComponentType.ID)
	assert.Equal(t, "ocp", result.ComponentType.Name)
}

func TestGetComponentType_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetComponentType("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get component type")
}

func TestCreateComponentType_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		response := ComponentTypeResponse{ComponentType: ComponentType{ID: "new-ct", Name: "new-type"}}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateComponentType("new-type")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-ct", result.ComponentType.ID)
}

func TestCreateComponentType_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateComponentType("bad")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create component type")
}

func TestUpdateComponentType_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Contains(t, r.URL.Path, "/componenttypes/ct-123")

		response := ComponentTypeResponse{ComponentType: ComponentType{ID: "ct-123", Name: "updated-ct"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateComponentType("ct-123", UpdateComponentTypeRequest{Name: "updated-ct"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "updated-ct", result.ComponentType.Name)
}

func TestUpdateComponentType_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateComponentType("nonexistent", UpdateComponentTypeRequest{Name: "x"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update component type")
}

func TestDeleteComponentType_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/componenttypes/ct-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteComponentType("ct-123")
	assert.NoError(t, err)
}

func TestDeleteComponentType_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteComponentType("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete component type")
}

func TestGetRemoteCIs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/remotecis")
		assert.Equal(t, "GET", r.Method)

		response := RemoteCIsResponse{
			Meta: Meta{Count: 2},
			RemoteCIs: []RemoteCI{
				{ID: "rci-1", Name: "remoteci-1", TeamID: "team-1"},
				{ID: "rci-2", Name: "remoteci-2", TeamID: "team-2"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetRemoteCIs()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.RemoteCIs, 2)
	assert.Equal(t, "rci-1", result.RemoteCIs[0].ID)
}

func TestGetRemoteCIs_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "server error"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetRemoteCIs()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get remote CIs")
}

func TestGetRemoteCI_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/remotecis/rci-123")
		assert.Equal(t, "GET", r.Method)

		response := RemoteCIResponse{RemoteCI: RemoteCI{ID: "rci-123", Name: "test-remoteci", TeamID: "team-1"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetRemoteCI("rci-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "rci-123", result.RemoteCI.ID)
	assert.Equal(t, "test-remoteci", result.RemoteCI.Name)
}

func TestGetRemoteCI_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetRemoteCI("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get remote CI")
}

func TestCreateRemoteCI_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqBody CreateRemoteCIRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "new-remoteci", reqBody.Name)
		assert.Equal(t, "team-1", reqBody.TeamID)

		response := RemoteCIResponse{RemoteCI: RemoteCI{ID: "new-rci", Name: "new-remoteci", TeamID: "team-1"}}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateRemoteCI("new-remoteci", "team-1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-rci", result.RemoteCI.ID)
}

func TestCreateRemoteCI_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateRemoteCI("bad", "bad")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create remote CI")
}

func TestUpdateRemoteCI_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Contains(t, r.URL.Path, "/remotecis/rci-123")

		response := RemoteCIResponse{RemoteCI: RemoteCI{ID: "rci-123", Name: "updated-remoteci"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateRemoteCI("rci-123", UpdateRemoteCIRequest{Name: "updated-remoteci"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "updated-remoteci", result.RemoteCI.Name)
}

func TestUpdateRemoteCI_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateRemoteCI("nonexistent", UpdateRemoteCIRequest{Name: "x"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update remote CI")
}

func TestDeleteRemoteCI_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/remotecis/rci-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteRemoteCI("rci-123")
	assert.NoError(t, err)
}

func TestDeleteRemoteCI_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteRemoteCI("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete remote CI")
}

func TestGetTeams_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/teams")
		assert.Equal(t, "GET", r.Method)

		response := TeamsResponse{
			Meta: Meta{Count: 2},
			Teams: []Team{
				{ID: "team-1", Name: "Team Alpha"},
				{ID: "team-2", Name: "Team Beta"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTeams()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Teams, 2)
	assert.Equal(t, "team-1", result.Teams[0].ID)
	assert.Equal(t, "Team Alpha", result.Teams[0].Name)
}

func TestGetTeams_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "server error"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTeams()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get teams")
}

func TestGetTeam_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/teams/team-123")
		assert.Equal(t, "GET", r.Method)

		response := TeamResponse{Team: Team{ID: "team-123", Name: "Test Team"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTeam("team-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "team-123", result.Team.ID)
	assert.Equal(t, "Test Team", result.Team.Name)
}

func TestGetTeam_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTeam("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get team")
}

func TestCreateTeam_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqBody CreateTeamRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "New Team", reqBody.Name)

		response := TeamResponse{Team: Team{ID: "new-team", Name: "New Team"}}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateTeam("New Team")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-team", result.Team.ID)
}

func TestCreateTeam_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateTeam("bad")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create team")
}

func TestUpdateTeam_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Contains(t, r.URL.Path, "/teams/team-123")

		response := TeamResponse{Team: Team{ID: "team-123", Name: "Updated Team"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateTeam("team-123", UpdateTeamRequest{Name: "Updated Team"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Team", result.Team.Name)
}

func TestUpdateTeam_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateTeam("nonexistent", UpdateTeamRequest{Name: "x"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update team")
}

func TestDeleteTeam_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/teams/team-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteTeam("team-123")
	assert.NoError(t, err)
}

func TestDeleteTeam_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteTeam("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete team")
}

func TestGetUsers_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/users")
		assert.Equal(t, "GET", r.Method)

		response := UsersResponse{
			Meta: Meta{Count: 2},
			Users: []User{
				{ID: "user-1", Name: "alice", Email: "alice@example.com"},
				{ID: "user-2", Name: "bob", Email: "bob@example.com"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetUsers()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, "user-1", result.Users[0].ID)
	assert.Equal(t, "alice", result.Users[0].Name)
}

func TestGetUsers_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "server error"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetUsers()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get users")
}

func TestGetUser_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/users/user-123")
		assert.Equal(t, "GET", r.Method)

		response := UserResponse{User: User{ID: "user-123", Name: "testuser", Email: "test@example.com"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetUser("user-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user-123", result.User.ID)
	assert.Equal(t, "testuser", result.User.Name)
}

func TestGetUser_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetUser("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get user")
}

func TestCreateUser_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqBody CreateUserRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "newuser", reqBody.Name)
		assert.Equal(t, "new@example.com", reqBody.Email)
		assert.Equal(t, "New User", reqBody.Fullname)
		assert.Equal(t, "team-1", reqBody.TeamID)

		response := UserResponse{User: User{ID: "new-user", Name: "newuser", Email: "new@example.com"}}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateUser("newuser", "new@example.com", "New User", "team-1", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-user", result.User.ID)
}

func TestCreateUser_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "invalid"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.CreateUser("bad", "bad", "bad", "bad", "bad")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestUpdateUser_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Contains(t, r.URL.Path, "/users/user-123")

		response := UserResponse{User: User{ID: "user-123", Name: "updated-user", Email: "updated@example.com"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateUser("user-123", UpdateUserRequest{Name: "updated-user", Email: "updated@example.com"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "updated-user", result.User.Name)
}

func TestUpdateUser_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.UpdateUser("nonexistent", UpdateUserRequest{Name: "x"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update user")
}

func TestDeleteUser_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/users/user-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteUser("user-123")
	assert.NoError(t, err)
}

func TestDeleteUser_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteUser("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete user")
}

func TestGetProducts_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/products")
		assert.Equal(t, "GET", r.Method)

		response := ProductsResponse{
			Meta: Meta{Count: 2},
			Products: []Product{
				{ID: "prod-1", Name: "RHEL", Label: "rhel"},
				{ID: "prod-2", Name: "OpenShift", Label: "ocp"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetProducts()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Products, 2)
	assert.Equal(t, "prod-1", result.Products[0].ID)
	assert.Equal(t, "RHEL", result.Products[0].Name)
}

func TestGetProducts_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "server error"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetProducts()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get products")
}

func TestGetProduct_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/products/prod-123")
		assert.Equal(t, "GET", r.Method)

		response := ProductResponse{Product: Product{ID: "prod-123", Name: "OpenShift", Label: "ocp"}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetProduct("prod-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "prod-123", result.Product.ID)
	assert.Equal(t, "OpenShift", result.Product.Name)
}

func TestGetProduct_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetProduct("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get product")
}

func TestGetFile_Success(t *testing.T) {
	expectedContent := []byte("file content bytes here")
	expectedContentType := "application/xml"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/files/file-123")
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", expectedContentType)
		_, err := w.Write(expectedContent)
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	content, contentType, err := client.GetFile("file-123")
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, content)
	assert.Equal(t, expectedContentType, contentType)
}

func TestGetFile_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	content, contentType, err := client.GetFile("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, content)
	assert.Empty(t, contentType)
	assert.Contains(t, err.Error(), "failed to get file")
}

func TestDeleteFile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Contains(t, r.URL.Path, "/files/file-123")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteFile("file-123")
	assert.NoError(t, err)
}

func TestDeleteFile_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"error": "not found"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	err := client.DeleteFile("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete file")
}

func TestUploadFile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/files", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "job-123", r.Header.Get("DCI-JOB-ID"))
		assert.Equal(t, "test.xml", r.Header.Get("DCI-NAME"))
		assert.Equal(t, "application/xml", r.Header.Get("DCI-MIME"))

		response := UploadFileResponse{}
		response.File.ID = "file-456"
		response.File.JobID = "job-123"
		response.File.Name = "test.xml"
		response.File.Mime = "application/xml"

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.xml")
	err := os.WriteFile(filePath, []byte("<test>content</test>"), 0644)
	assert.NoError(t, err)

	client := newTestClient(server.URL)
	result, err := client.UploadFile("job-123", filePath, "application/xml")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "file-456", result.File.ID)
	assert.Equal(t, "job-123", result.File.JobID)
	assert.Equal(t, "test.xml", result.File.Name)
}

func TestUploadFile_FileNotFound(t *testing.T) {
	client := &Client{BaseURL: "http://localhost", AccessKey: "testKey", SecretKey: "testSecret", httpClient: &http.Client{}}
	result, err := client.UploadFile("job-123", "/nonexistent/path/file.xml", "application/xml")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error reading file")
}

func TestUploadFile_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "bad request"}`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.xml")
	err := os.WriteFile(filePath, []byte("content"), 0644)
	assert.NoError(t, err)

	client := newTestClient(server.URL)
	result, err := client.UploadFile("bad-job", filePath, "application/xml")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to upload file")
}

