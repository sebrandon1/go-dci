package lib

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	signerv4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

const (
	// https://doc.distributed-ci.io/python-dciauth/#using-postman
	DCIURL      = "https://api.distributed-ci.io/api/v1"
	awsRegion   = "BHS3"
	serviceName = "api"
	dateFormat  = "2006-01-02T15:04:05.999999"
	maxRecords  = 50000
	// SHA-256 of empty string for unsigned GET requests
	emptyStringSHA256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

type Client struct {
	BaseURL   string
	AccessKey string
	SecretKey string
}

func NewClient(accessKey, secretKey string) *Client {
	return &Client{
		BaseURL:   DCIURL,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetIdentity retrieves the authenticated user/remoteci identity from the DCI API
func (c *Client) GetIdentity() (*IdentityResponse, error) {
	httpResponse, err := httpGetSimpleWithAWSAuth(c.BaseURL+"/identity", awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		fmt.Printf("Error getting identity: %s\n", err)
		return nil, err
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed with status code: %d", httpResponse.StatusCode)
	}

	var identity IdentityResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&identity)
	if err != nil {
		fmt.Printf("Error decoding the response: %s\n", err)
		return nil, err
	}

	return &identity, nil
}

// httpGetSimpleWithAWSAuth performs an authenticated GET request without pagination parameters
func httpGetSimpleWithAWSAuth(url, region, svcName, accessKey, secretKey string) (*http.Response, error) {
	signer := signerv4.NewSigner()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Sign the request
	creds := aws.Credentials{AccessKeyID: accessKey, SecretAccessKey: secretKey}
	if err := signer.SignHTTP(context.Background(), creds, req, emptyStringSHA256, svcName, region, time.Now()); err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

// GetComponentTypes retrieves all component types from the DCI API with pagination
func (c *Client) GetComponentTypes() ([]ComponentTypesResponse, error) {
	var componentTypesCollection []ComponentTypesResponse

	requestLimit := 100
	offset := 0

	for {
		componentTypes, err := c.fetchComponentTypes(requestLimit, offset)
		if err != nil {
			return nil, err
		}

		componentTypesCollection = append(componentTypesCollection, componentTypes)

		// If the number of component types returned is less than the request limit, we have reached the end
		if len(componentTypes.ComponentTypes) < requestLimit {
			break
		}

		// If we have reached the maximum number of records, we can stop the loop
		if offset >= maxRecords {
			break
		}

		// Increment the offset
		offset += requestLimit
	}

	return componentTypesCollection, nil
}

// fetchComponentTypes is an internal helper to fetch component types with pagination
func (c *Client) fetchComponentTypes(requestLimit, offset int) (ComponentTypesResponse, error) {
	httpResponse, err := HttpGetWithAWSAuth(c.BaseURL+"/componenttypes", awsRegion, serviceName, c.AccessKey, c.SecretKey, requestLimit, offset)
	if err != nil {
		fmt.Printf("Error getting component types: %s\n", err)
		return ComponentTypesResponse{}, err
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	var componentTypes ComponentTypesResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&componentTypes)
	if err != nil {
		fmt.Printf("Error decoding the response: %s\n", err)
		return ComponentTypesResponse{}, err
	}

	return componentTypes, nil
}

func (c *Client) GetTopics() ([]TopicsResponse, error) {
	var topicsCollection []TopicsResponse

	requestLimit := 100
	offset := 0
	maxRecords := 50000

	for {
		httpResponse, err := HttpGetWithAWSAuth(c.BaseURL+"/topics", awsRegion, serviceName, c.AccessKey, c.SecretKey, 100, 0)
		if err != nil {
			fmt.Printf("Error getting topics: %s\n", err)
			return nil, err
		}

		defer func() {
			err := httpResponse.Body.Close()
			if err != nil {
				fmt.Printf("Error closing the response body: %s\n", err)
			}
		}()

		var topics TopicsResponse
		err = json.NewDecoder(httpResponse.Body).Decode(&topics)
		if err != nil {
			fmt.Printf("Error decoding the response: %s\n", err)
			return nil, err
		}

		topicsCollection = append(topicsCollection, topics)

		offset += requestLimit

		// If the number of topics returned is less than the request limit, we have reached the end
		if len(topics.Topics) < requestLimit {
			break
		}

		// If we have reached the maximum number of records, we can stop the loop
		if len(topics.Topics) >= maxRecords {
			break
		}

		// Increment the offset
		offset += requestLimit
	}

	return topicsCollection, nil
}

// GetTopic retrieves a single topic by ID from the DCI API
func (c *Client) GetTopic(topicID string) (*TopicResponse, error) {
	url := fmt.Sprintf("%s/topics/%s", c.BaseURL, topicID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting topic: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get topic with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var topic TopicResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&topic); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &topic, nil
}

// CreateTopic creates a new topic in DCI
func (c *Client) CreateTopic(name, productID string, componentTypes []string) (*TopicResponse, error) {
	reqBody := CreateTopicRequest{
		Name:           name,
		ProductID:      productID,
		ComponentTypes: componentTypes,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	httpResponse, err := c.httpPostWithAWSAuth(c.BaseURL+"/topics", jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating topic: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to create topic with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response TopicResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateTopic updates an existing topic in DCI
func (c *Client) UpdateTopic(topicID string, updates UpdateTopicRequest) (*TopicResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/topics/%s", c.BaseURL, topicID)
	httpResponse, err := c.httpPutWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating topic: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to update topic with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response TopicResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteTopic deletes a topic from DCI
func (c *Client) DeleteTopic(topicID string) error {
	url := fmt.Sprintf("%s/topics/%s", c.BaseURL, topicID)
	httpResponse, err := c.httpDeleteWithAWSAuth(url)
	if err != nil {
		return fmt.Errorf("error deleting topic: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("failed to delete topic with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	return nil
}

// GetTopicComponents retrieves all components for a specific topic using the topic components endpoint
func (c *Client) GetTopicComponents(topicID string) ([]ComponentsResponse, error) {
	var componentsCollection []ComponentsResponse

	requestLimit := 100
	offset := 0

	for {
		url := fmt.Sprintf("%s/topics/%s/components", c.BaseURL, topicID)
		httpResponse, err := HttpGetWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey, requestLimit, offset)
		if err != nil {
			return nil, fmt.Errorf("error getting topic components: %w", err)
		}

		defer func() {
			if cerr := httpResponse.Body.Close(); cerr != nil {
				fmt.Printf("Error closing response body: %v\n", cerr)
			}
		}()

		var components ComponentsResponse
		if err := json.NewDecoder(httpResponse.Body).Decode(&components); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}

		componentsCollection = append(componentsCollection, components)

		// If fewer results than limit, we've reached the end
		if len(components.Components) < requestLimit {
			break
		}

		// If we've reached max records, stop
		if offset >= maxRecords {
			break
		}

		offset += requestLimit
	}

	return componentsCollection, nil
}

func (c *Client) GetJobs(daysBackLimit int) ([]JobsResponse, error) {
	var jobCollection []JobsResponse

	// Default values to page through the results
	requestLimit := 100
	offset := 0

	for {
		outOfDateRangeJobReturned := false

		jobs, err := c.fetchJobs(requestLimit, offset)
		if err != nil {
			return nil, err
		}

		jobCollection = append(jobCollection, jobs)

		// Increment the offset
		offset += requestLimit

		// Check if the job is out of the date range
		for _, job := range jobs.Jobs {
			// Parse the created at date
			createdAt, err := time.Parse(dateFormat, job.CreatedAt)
			if err != nil {
				fmt.Printf("Error parsing the created at date: %s\n", err)
				continue
			}

			// If the job is out of the date range, we can stop the loop
			if time.Since(createdAt).Hours() > float64(daysBackLimit*24) {
				outOfDateRangeJobReturned = true
				break
			}
		}

		// If the number of jobs returned is less than the request limit, we have reached the end
		if len(jobs.Jobs) < requestLimit {
			break
		}

		// If we have reached the end, we can stop the loop
		if outOfDateRangeJobReturned {
			break
		}

		// If we have reached the maximum number of records, we can stop the loop
		if len(jobCollection) >= maxRecords {
			break
		}
	}

	return jobCollection, nil
}

func (c *Client) GetJobsByDate(startDate, endDate time.Time) ([]JobsResponse, error) {
	var jobCollection []JobsResponse

	// Default values to page through the results
	requestLimit := 100
	offset := 0

	for {
		jobs, err := c.fetchJobs(requestLimit, offset)
		if err != nil {
			return nil, err
		}

		for _, job := range jobs.Jobs {
			// Parse the created at date
			createdAt, err := time.Parse(dateFormat, job.CreatedAt)
			if err != nil {
				fmt.Printf("Error parsing the created at date: %s\n", err)
				continue
			}

			// If the job is within the date range, add it to jobCollection
			if createdAt.After(startDate) && createdAt.Before(endDate) {
				jobCollection = append(jobCollection, jobs)
				break
			}
		}

		// Increment the offset for the next page
		offset += requestLimit

		// If the number of jobs returned is less than the request limit, we have reached the end
		if len(jobs.Jobs) < requestLimit {
			break
		}

		// If we have reached the maximum number of records, we can stop the loop
		if len(jobCollection) >= maxRecords {
			break
		}
	}

	return jobCollection, nil
}

// GetJob retrieves a single job by ID from the DCI API
func (c *Client) GetJob(jobID string) (*JobResponse, error) {
	url := fmt.Sprintf("%s/jobs/%s", c.BaseURL, jobID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting job: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get job with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var job JobResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&job); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &job, nil
}

// UpdateJob updates an existing job in DCI
func (c *Client) UpdateJob(jobID string, updates UpdateJobRequest) (*JobResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/jobs/%s", c.BaseURL, jobID)
	httpResponse, err := c.httpPutWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating job: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to update job with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response JobResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteJob deletes a job from DCI
func (c *Client) DeleteJob(jobID string) error {
	url := fmt.Sprintf("%s/jobs/%s", c.BaseURL, jobID)
	httpResponse, err := c.httpDeleteWithAWSAuth(url)
	if err != nil {
		return fmt.Errorf("error deleting job: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("failed to delete job with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	return nil
}

// ScheduleJob schedules a job with auto-selected components for a topic
func (c *Client) ScheduleJob(topicID string) (*CreateJobResponse, error) {
	reqBody := ScheduleJobRequest{
		TopicID: topicID,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/jobs/schedule", c.BaseURL)
	httpResponse, err := c.httpPostWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error scheduling job: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to schedule job with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response CreateJobResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetJobFiles retrieves all files for a specific job
func (c *Client) GetJobFiles(jobID string) (*FilesResponse, error) {
	url := fmt.Sprintf("%s/jobs/%s/files", c.BaseURL, jobID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting job files: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get job files with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response FilesResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

func HttpGetWithAWSAuth(url, region, serviceName, accessKey, secretKey string, limit, offset int) (*http.Response, error) {
	// Create signer using aws-sdk-go-v2 v4 signer
	signer := signerv4.NewSigner()

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Build the query string
	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	q.Add("sort", "-created_at") // Sort by created_at in descending order
	req.URL.RawQuery = q.Encode()

	// Sign the request
	// For GET with empty body use the precomputed empty payload hash
	creds := aws.Credentials{AccessKeyID: accessKey, SecretAccessKey: secretKey}
	if err := signer.SignHTTP(context.Background(), creds, req, emptyStringSHA256, serviceName, region, time.Now()); err != nil {
		return nil, err
	}

	// Send the request
	client := &http.Client{}

	// Perform the requests and adjust the offset based on the response
	return client.Do(req)
}

func (c *Client) fetchJobs(requestLimit, offset int) (JobsResponse, error) {
	// Get jobs from the API
	httpResponse, err := HttpGetWithAWSAuth(c.BaseURL+"/jobs", awsRegion, serviceName, c.AccessKey, c.SecretKey, requestLimit, offset)
	if err != nil {
		fmt.Printf("Error getting jobs: %s\n", err)
		return JobsResponse{}, err
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	// Decode the response into JobsResponse
	var jobs JobsResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&jobs)
	if err != nil {
		fmt.Printf("Error decoding the response: %s\n", err)
		return JobsResponse{}, err
	}

	return jobs, nil
}

// GetComponents retrieves all components from the DCI API with pagination
func (c *Client) GetComponents() ([]ComponentsResponse, error) {
	var componentsCollection []ComponentsResponse

	requestLimit := 100
	offset := 0

	for {
		components, err := c.fetchComponents("", requestLimit, offset)
		if err != nil {
			return nil, err
		}

		componentsCollection = append(componentsCollection, components)

		// If the number of components returned is less than the request limit, we have reached the end
		if len(components.Components) < requestLimit {
			break
		}

		// If we have reached the maximum number of records, we can stop the loop
		if offset >= maxRecords {
			break
		}

		// Increment the offset
		offset += requestLimit
	}

	return componentsCollection, nil
}

// GetComponentsByTopicID retrieves components filtered by topic ID
func (c *Client) GetComponentsByTopicID(topicID string) ([]ComponentsResponse, error) {
	var componentsCollection []ComponentsResponse

	requestLimit := 100
	offset := 0

	for {
		components, err := c.fetchComponents(topicID, requestLimit, offset)
		if err != nil {
			return nil, err
		}

		componentsCollection = append(componentsCollection, components)

		// If the number of components returned is less than the request limit, we have reached the end
		if len(components.Components) < requestLimit {
			break
		}

		// If we have reached the maximum number of records, we can stop the loop
		if offset >= maxRecords {
			break
		}

		// Increment the offset
		offset += requestLimit
	}

	return componentsCollection, nil
}

// fetchComponents is an internal helper to fetch components with optional topic filtering
func (c *Client) fetchComponents(topicID string, requestLimit, offset int) (ComponentsResponse, error) {
	url := c.BaseURL + "/components"

	httpResponse, err := httpGetComponentsWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey, topicID, requestLimit, offset)
	if err != nil {
		fmt.Printf("Error getting components: %s\n", err)
		return ComponentsResponse{}, err
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	var components ComponentsResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&components)
	if err != nil {
		fmt.Printf("Error decoding the response: %s\n", err)
		return ComponentsResponse{}, err
	}

	return components, nil
}

// httpGetComponentsWithAWSAuth performs an authenticated GET request for components with optional topic filtering
func httpGetComponentsWithAWSAuth(url, region, svcName, accessKey, secretKey, topicID string, limit, offset int) (*http.Response, error) {
	signer := signerv4.NewSigner()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Build the query string
	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	q.Add("sort", "-created_at")

	// Add topic_id filter if provided
	if topicID != "" {
		q.Add("where", "topic_id:"+topicID)
	}

	req.URL.RawQuery = q.Encode()

	// Sign the request
	creds := aws.Credentials{AccessKeyID: accessKey, SecretAccessKey: secretKey}
	if err := signer.SignHTTP(context.Background(), creds, req, emptyStringSHA256, svcName, region, time.Now()); err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

// CreateJob creates a new job in DCI
func (c *Client) CreateJob(topicID string, componentIDs []string, comment string) (*CreateJobResponse, error) {
	reqBody := CreateJobRequest{
		TopicID:    topicID,
		Components: componentIDs,
		Comment:    comment,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	httpResponse, err := c.httpPostWithAWSAuth(c.BaseURL+"/jobs", jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating job: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to create job with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response CreateJobResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateJobState updates the state of a job (pre-run, running, success, failure, etc.)
func (c *Client) UpdateJobState(jobID string, status JobState, comment string) (*JobStateResponse, error) {
	reqBody := UpdateJobStateRequest{
		JobID:   jobID,
		Status:  string(status),
		Comment: comment,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	httpResponse, err := c.httpPostWithAWSAuth(c.BaseURL+"/jobstates", jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating job state: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to update job state with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response JobStateResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UploadFile uploads a file (e.g., test results) to a job in DCI
func (c *Client) UploadFile(jobID, filePath, mimeType string) (*UploadFileResponse, error) {
	// Read the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	fileName := filepath.Base(filePath)

	httpResponse, err := c.httpPostFileWithAWSAuth(c.BaseURL+"/files", fileContent, jobID, fileName, mimeType)
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to upload file with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response UploadFileResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UploadFileContent uploads file content directly (without reading from disk) to a job in DCI
func (c *Client) UploadFileContent(jobID, fileName, mimeType string, content []byte) (*UploadFileResponse, error) {
	httpResponse, err := c.httpPostFileWithAWSAuth(c.BaseURL+"/files", content, jobID, fileName, mimeType)
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to upload file with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response UploadFileResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// httpPostWithAWSAuth performs an authenticated POST request with JSON body
func (c *Client) httpPostWithAWSAuth(url string, jsonBody []byte) (*http.Response, error) {
	signer := signerv4.NewSigner()

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Calculate SHA256 hash of the body for signing
	hash := sha256.Sum256(jsonBody)
	payloadHash := hex.EncodeToString(hash[:])

	// Sign the request
	creds := aws.Credentials{AccessKeyID: c.AccessKey, SecretAccessKey: c.SecretKey}
	if err := signer.SignHTTP(context.Background(), creds, req, payloadHash, serviceName, awsRegion, time.Now()); err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

// httpPutWithAWSAuth performs an authenticated PUT request with JSON body
func (c *Client) httpPutWithAWSAuth(url string, jsonBody []byte) (*http.Response, error) {
	signer := signerv4.NewSigner()

	req, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Calculate SHA256 hash of the body for signing
	hash := sha256.Sum256(jsonBody)
	payloadHash := hex.EncodeToString(hash[:])

	// Sign the request
	creds := aws.Credentials{AccessKeyID: c.AccessKey, SecretAccessKey: c.SecretKey}
	if err := signer.SignHTTP(context.Background(), creds, req, payloadHash, serviceName, awsRegion, time.Now()); err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

// httpDeleteWithAWSAuth performs an authenticated DELETE request
func (c *Client) httpDeleteWithAWSAuth(url string) (*http.Response, error) {
	signer := signerv4.NewSigner()

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	// Sign the request
	creds := aws.Credentials{AccessKeyID: c.AccessKey, SecretAccessKey: c.SecretKey}
	if err := signer.SignHTTP(context.Background(), creds, req, emptyStringSHA256, serviceName, awsRegion, time.Now()); err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

// httpPostFileWithAWSAuth performs an authenticated POST request for file uploads
func (c *Client) httpPostFileWithAWSAuth(url string, content []byte, jobID, fileName, mimeType string) (*http.Response, error) {
	signer := signerv4.NewSigner()

	req, err := http.NewRequest("POST", url, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	// Set DCI-specific headers for file upload
	req.Header.Set("DCI-JOB-ID", jobID)
	req.Header.Set("DCI-NAME", fileName)
	req.Header.Set("DCI-MIME", mimeType)
	req.Header.Set("Content-Type", "application/octet-stream")

	// Calculate SHA256 hash of the body for signing
	hash := sha256.Sum256(content)
	payloadHash := hex.EncodeToString(hash[:])

	// Sign the request
	creds := aws.Credentials{AccessKeyID: c.AccessKey, SecretAccessKey: c.SecretKey}
	if err := signer.SignHTTP(context.Background(), creds, req, payloadHash, serviceName, awsRegion, time.Now()); err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}
