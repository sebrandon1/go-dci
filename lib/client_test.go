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

func TestGetIdentity_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request path
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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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
		BaseURL:   server.URL,
		AccessKey: "badKey",
		SecretKey: "badSecret",
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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

	identity, err := client.GetIdentity()
	assert.Error(t, err)
	assert.Nil(t, identity)
}

func TestIdentityResponse_Struct(t *testing.T) {
	// Test JSON marshaling/unmarshaling of IdentityResponse
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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

// POST operation tests

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

	client := &Client{
		BaseURL:   server.URL,
		AccessKey: "testKey",
		SecretKey: "testSecret",
	}

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

