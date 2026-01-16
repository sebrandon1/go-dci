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

// GetComponentType retrieves a single component type by ID from the DCI API
func (c *Client) GetComponentType(componentTypeID string) (*ComponentTypeResponse, error) {
	url := fmt.Sprintf("%s/componenttypes/%s", c.BaseURL, componentTypeID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting component type: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get component type with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var componentType ComponentTypeResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&componentType); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &componentType, nil
}

// CreateComponentType creates a new component type in DCI
func (c *Client) CreateComponentType(name string) (*ComponentTypeResponse, error) {
	reqBody := CreateComponentTypeRequest{
		Name: name,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/componenttypes", c.BaseURL)
	httpResponse, err := c.httpPostWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating component type: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to create component type with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response ComponentTypeResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateComponentType updates an existing component type in DCI
func (c *Client) UpdateComponentType(componentTypeID string, updates UpdateComponentTypeRequest) (*ComponentTypeResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/componenttypes/%s", c.BaseURL, componentTypeID)
	httpResponse, err := c.httpPutWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating component type: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to update component type with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response ComponentTypeResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteComponentType deletes a component type from DCI
func (c *Client) DeleteComponentType(componentTypeID string) error {
	url := fmt.Sprintf("%s/componenttypes/%s", c.BaseURL, componentTypeID)
	httpResponse, err := c.httpDeleteWithAWSAuth(url)
	if err != nil {
		return fmt.Errorf("error deleting component type: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("failed to delete component type with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	return nil
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

// GetComponent retrieves a single component by ID from the DCI API
func (c *Client) GetComponent(componentID string) (*ComponentResponse, error) {
	url := fmt.Sprintf("%s/components/%s", c.BaseURL, componentID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting component: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get component with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var component ComponentResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&component); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &component, nil
}

// CreateComponent creates a new component in DCI
func (c *Client) CreateComponent(name, componentType, topicID, version string) (*ComponentResponse, error) {
	reqBody := CreateComponentRequest{
		Name:    name,
		Type:    componentType,
		TopicID: topicID,
		Version: version,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/components", c.BaseURL)
	httpResponse, err := c.httpPostWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating component: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to create component with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response ComponentResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateComponent updates an existing component in DCI
func (c *Client) UpdateComponent(componentID string, updates UpdateComponentRequest) (*ComponentResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/components/%s", c.BaseURL, componentID)
	httpResponse, err := c.httpPutWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating component: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to update component with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response ComponentResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteComponent deletes a component from DCI
func (c *Client) DeleteComponent(componentID string) error {
	url := fmt.Sprintf("%s/components/%s", c.BaseURL, componentID)
	httpResponse, err := c.httpDeleteWithAWSAuth(url)
	if err != nil {
		return fmt.Errorf("error deleting component: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("failed to delete component with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	return nil
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

	url := fmt.Sprintf("%s/jobstates", c.BaseURL)
	httpResponse, err := c.httpPostWithAWSAuth(url, jsonBody)
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

// CreateJobState creates a new job state entry
func (c *Client) CreateJobState(jobID string, status JobState, comment string) (*JobStateResponse, error) {
	reqBody := UpdateJobStateRequest{
		JobID:   jobID,
		Status:  string(status),
		Comment: comment,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/jobstates", c.BaseURL)
	httpResponse, err := c.httpPostWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating job state: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to create job state with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response JobStateResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetJobStates retrieves job states, optionally filtered by job ID
func (c *Client) GetJobStates(jobID string) (*JobStatesResponse, error) {
	url := fmt.Sprintf("%s/jobstates", c.BaseURL)
	if jobID != "" {
		url = fmt.Sprintf("%s?where=job_id:%s", url, jobID)
	}

	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting job states: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get job states with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response JobStatesResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetFile downloads a file by ID from DCI
func (c *Client) GetFile(fileID string) ([]byte, string, error) {
	url := fmt.Sprintf("%s/files/%s", c.BaseURL, fileID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, "", fmt.Errorf("error getting file: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, "", fmt.Errorf("failed to get file with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	content, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading file content: %w", err)
	}

	contentType := httpResponse.Header.Get("Content-Type")
	return content, contentType, nil
}

// DeleteFile deletes a file from DCI
func (c *Client) DeleteFile(fileID string) error {
	url := fmt.Sprintf("%s/files/%s", c.BaseURL, fileID)
	httpResponse, err := c.httpDeleteWithAWSAuth(url)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("failed to delete file with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	return nil
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

// GetRemoteCIs retrieves all remote CIs from DCI
func (c *Client) GetRemoteCIs() (*RemoteCIsResponse, error) {
	url := fmt.Sprintf("%s/remotecis", c.BaseURL)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting remote CIs: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get remote CIs with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response RemoteCIsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetRemoteCI retrieves a specific remote CI by ID
func (c *Client) GetRemoteCI(remoteciID string) (*RemoteCIResponse, error) {
	url := fmt.Sprintf("%s/remotecis/%s", c.BaseURL, remoteciID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting remote CI: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get remote CI with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response RemoteCIResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CreateRemoteCI creates a new remote CI in DCI
func (c *Client) CreateRemoteCI(name, teamID string) (*RemoteCIResponse, error) {
	reqBody := CreateRemoteCIRequest{
		Name:   name,
		TeamID: teamID,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/remotecis", c.BaseURL)
	httpResponse, err := c.httpPostWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating remote CI: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to create remote CI with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response RemoteCIResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateRemoteCI updates an existing remote CI in DCI
func (c *Client) UpdateRemoteCI(remoteciID string, updates UpdateRemoteCIRequest) (*RemoteCIResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/remotecis/%s", c.BaseURL, remoteciID)
	httpResponse, err := c.httpPutWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating remote CI: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to update remote CI with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response RemoteCIResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteRemoteCI deletes a remote CI from DCI
func (c *Client) DeleteRemoteCI(remoteciID string) error {
	url := fmt.Sprintf("%s/remotecis/%s", c.BaseURL, remoteciID)
	httpResponse, err := c.httpDeleteWithAWSAuth(url)
	if err != nil {
		return fmt.Errorf("error deleting remote CI: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("failed to delete remote CI with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	return nil
}

// GetTeams retrieves all teams from DCI
func (c *Client) GetTeams() (*TeamsResponse, error) {
	url := fmt.Sprintf("%s/teams", c.BaseURL)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting teams: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get teams with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response TeamsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetTeam retrieves a specific team by ID
func (c *Client) GetTeam(teamID string) (*TeamResponse, error) {
	url := fmt.Sprintf("%s/teams/%s", c.BaseURL, teamID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting team: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get team with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response TeamResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CreateTeam creates a new team in DCI
func (c *Client) CreateTeam(name string) (*TeamResponse, error) {
	reqBody := CreateTeamRequest{
		Name: name,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/teams", c.BaseURL)
	httpResponse, err := c.httpPostWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating team: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to create team with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response TeamResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateTeam updates an existing team in DCI
func (c *Client) UpdateTeam(teamID string, updates UpdateTeamRequest) (*TeamResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/teams/%s", c.BaseURL, teamID)
	httpResponse, err := c.httpPutWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating team: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to update team with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response TeamResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteTeam deletes a team from DCI
func (c *Client) DeleteTeam(teamID string) error {
	url := fmt.Sprintf("%s/teams/%s", c.BaseURL, teamID)
	httpResponse, err := c.httpDeleteWithAWSAuth(url)
	if err != nil {
		return fmt.Errorf("error deleting team: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("failed to delete team with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	return nil
}

// GetUsers retrieves all users from DCI
func (c *Client) GetUsers() (*UsersResponse, error) {
	url := fmt.Sprintf("%s/users", c.BaseURL)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get users with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response UsersResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetUser retrieves a specific user by ID
func (c *Client) GetUser(userID string) (*UserResponse, error) {
	url := fmt.Sprintf("%s/users/%s", c.BaseURL, userID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get user with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response UserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CreateUser creates a new user in DCI
func (c *Client) CreateUser(name, email, fullname, teamID, password string) (*UserResponse, error) {
	reqBody := CreateUserRequest{
		Name:     name,
		Email:    email,
		Fullname: fullname,
		TeamID:   teamID,
		Password: password,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/users", c.BaseURL)
	httpResponse, err := c.httpPostWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to create user with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response UserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateUser updates an existing user in DCI
func (c *Client) UpdateUser(userID string, updates UpdateUserRequest) (*UserResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	url := fmt.Sprintf("%s/users/%s", c.BaseURL, userID)
	httpResponse, err := c.httpPutWithAWSAuth(url, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to update user with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response UserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteUser deletes a user from DCI
func (c *Client) DeleteUser(userID string) error {
	url := fmt.Sprintf("%s/users/%s", c.BaseURL, userID)
	httpResponse, err := c.httpDeleteWithAWSAuth(url)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return fmt.Errorf("failed to delete user with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	return nil
}

// GetProducts retrieves all products from DCI
func (c *Client) GetProducts() (*ProductsResponse, error) {
	url := fmt.Sprintf("%s/products", c.BaseURL)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting products: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get products with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response ProductsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetProduct retrieves a specific product by ID
func (c *Client) GetProduct(productID string) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/products/%s", c.BaseURL, productID)
	httpResponse, err := httpGetSimpleWithAWSAuth(url, awsRegion, serviceName, c.AccessKey, c.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	defer func() {
		if cerr := httpResponse.Body.Close(); cerr != nil {
			fmt.Printf("Error closing response body: %v\n", cerr)
		}
	}()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("failed to get product with status code %d: %s", httpResponse.StatusCode, string(body))
	}

	var response ProductResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}
