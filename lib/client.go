package lib

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
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
	maxRecords      = 50000
	defaultPageSize = 100
	// SHA-256 of empty string for unsigned GET requests
	emptyStringSHA256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	// Default HTTP timeouts
	defaultRequestTimeout      = 30 * time.Second
	defaultTLSHandshakeTimeout = 5 * time.Second
	defaultDialTimeout         = 10 * time.Second
)

type Client struct {
	BaseURL         string
	AccessKey       string
	SecretKey       string
	httpClient      *http.Client
	MaxRetries      int
	RequestTimeout  time.Duration
	TLSTimeout      time.Duration
	DialTimeout     time.Duration
}

func NewClient(accessKey, secretKey string) *Client {
	return &Client{
		BaseURL:        DCIURL,
		AccessKey:      accessKey,
		SecretKey:      secretKey,
		httpClient:     newDefaultHTTPClient(),
		MaxRetries:     3,
		RequestTimeout: defaultRequestTimeout,
		TLSTimeout:     defaultTLSHandshakeTimeout,
		DialTimeout:    defaultDialTimeout,
	}
}

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultRequestTimeout,
		Transport: &http.Transport{
			TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
			ResponseHeaderTimeout: defaultRequestTimeout,
			DialContext: (&net.Dialer{
				Timeout: defaultDialTimeout,
			}).DialContext,
		},
	}
}

func isRetryable(method string, statusCode int) bool {
	if statusCode < 500 {
		return false
	}
	return method == http.MethodGet || method == http.MethodDelete
}

func (c *Client) retryBackoff(ctx context.Context, attempt int) error {
	backoff := time.Duration(1<<uint(attempt)) * time.Second
	jitter := backoff / 4
	sleep := backoff + time.Duration(rand.Int64N(int64(2*jitter))) - jitter

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(sleep):
		return nil
	}
}

func (c *Client) doRequest(ctx context.Context, method, reqURL string, body []byte, headers map[string]string) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	for attempt := range c.MaxRetries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		var bodyReader io.Reader
		if body != nil {
			bodyReader = bytes.NewReader(body)
		}

		req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
		if err != nil {
			return nil, err
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		payloadHash := emptyStringSHA256
		if body != nil {
			hash := sha256.Sum256(body)
			payloadHash = hex.EncodeToString(hash[:])
		}

		signer := signerv4.NewSigner()
		creds := aws.Credentials{AccessKeyID: c.AccessKey, SecretAccessKey: c.SecretKey}
		if err := signer.SignHTTP(ctx, creds, req, payloadHash, serviceName, awsRegion, time.Now()); err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < c.MaxRetries-1 {
				if err := c.retryBackoff(ctx, attempt); err != nil {
					return nil, err
				}
				continue
			}
			return nil, lastErr
		}

		if !isRetryable(method, resp.StatusCode) {
			return resp, nil
		}

		lastResp = resp
		if attempt < c.MaxRetries-1 {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()

			if err := c.retryBackoff(ctx, attempt); err != nil {
				return nil, err
			}
			continue
		}
	}

	// Exhausted all retries — return last response or error
	if lastResp != nil {
		return lastResp, nil
	}

	return nil, lastErr
}

// doJSON is a convenience method for POST/PUT requests with a JSON body.
func (c *Client) doJSON(ctx context.Context, method, reqURL string, jsonBody []byte) (*http.Response, error) {
	return c.doRequest(ctx, method, reqURL, jsonBody, map[string]string{"Content-Type": "application/json"})
}

// buildFilter creates a DCI API where clause filter for a single field.
// Returns empty string if value is empty.
func buildFilter(field, value string) string {
	if value == "" {
		return ""
	}
	return field + ":" + value
}

// paginatedURL builds a URL with limit, offset, sort, and optional where filters.
func paginatedURL(base string, limit, offset int, filters ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		// Fallback: return base as-is if unparseable (should not happen)
		return base
	}

	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	q.Set("sort", "-created_at")

	for _, f := range filters {
		if f != "" {
			q.Add("where", f)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// paginate is a generic helper that handles the standard pagination loop.
func paginate[T any](ctx context.Context, fetch func(limit, offset int) (T, int, error)) ([]T, error) {
	var collection []T
	offset := 0

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		page, count, err := fetch(defaultPageSize, offset)
		if err != nil {
			return nil, err
		}

		collection = append(collection, page)

		if count < defaultPageSize {
			break
		}

		if offset >= maxRecords {
			break
		}

		offset += defaultPageSize
	}

	return collection, nil
}

// GetIdentity retrieves the authenticated user/remoteci identity from the DCI API
func (c *Client) GetIdentity(ctx context.Context) (*IdentityResponse, error) {
	httpResponse, err := c.doRequest(ctx, http.MethodGet, c.BaseURL+"/identity", nil, nil)
	if err != nil {
		return nil, err
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var identity IdentityResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&identity)
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

// GetComponentTypes retrieves all component types from the DCI API with pagination
func (c *Client) GetComponentTypes(ctx context.Context) ([]ComponentTypesResponse, error) {
	return c.GetComponentTypesByName(ctx, "")
}

// GetComponentTypesByName retrieves component types filtered by name (empty string for all)
func (c *Client) GetComponentTypesByName(ctx context.Context, name string) ([]ComponentTypesResponse, error) {
	return paginate(ctx, func(limit, offset int) (ComponentTypesResponse, int, error) {
		resp, err := c.fetchComponentTypes(ctx, name, limit, offset)
		return resp, len(resp.ComponentTypes), err
	})
}

// fetchComponentTypes is an internal helper to fetch component types with optional name filtering
func (c *Client) fetchComponentTypes(ctx context.Context, name string, requestLimit, offset int) (ComponentTypesResponse, error) {
	filter := buildFilter("name", name)
	reqURL := paginatedURL(c.BaseURL+"/componenttypes", requestLimit, offset, filter)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return ComponentTypesResponse{}, err
	}

	defer func() { _ = httpResponse.Body.Close() }()

	var componentTypes ComponentTypesResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&componentTypes)
	if err != nil {
		return ComponentTypesResponse{}, err
	}

	return componentTypes, nil
}

// GetComponentType retrieves a single component type by ID from the DCI API
func (c *Client) GetComponentType(ctx context.Context, componentTypeID string) (*ComponentTypeResponse, error) {
	reqURL := fmt.Sprintf("%s/componenttypes/%s", c.BaseURL, componentTypeID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting component type: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var componentType ComponentTypeResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&componentType); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &componentType, nil
}

// CreateComponentType creates a new component type in DCI
func (c *Client) CreateComponentType(ctx context.Context, name string) (*ComponentTypeResponse, error) {
	reqBody := CreateComponentTypeRequest{
		Name: name,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/componenttypes", c.BaseURL)
	httpResponse, err := c.doJSON(ctx, http.MethodPost, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating component type: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response ComponentTypeResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateComponentType updates an existing component type in DCI
func (c *Client) UpdateComponentType(ctx context.Context, componentTypeID string, updates UpdateComponentTypeRequest) (*ComponentTypeResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/componenttypes/%s", c.BaseURL, componentTypeID)
	httpResponse, err := c.doJSON(ctx, http.MethodPut, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating component type: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response ComponentTypeResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteComponentType deletes a component type from DCI
func (c *Client) DeleteComponentType(ctx context.Context, componentTypeID string) error {
	reqURL := fmt.Sprintf("%s/componenttypes/%s", c.BaseURL, componentTypeID)
	httpResponse, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting component type: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return formatHTTPError(httpResponse.StatusCode, body)
	}

	return nil
}

func (c *Client) GetTopics(ctx context.Context) ([]TopicsResponse, error) {
	return c.GetTopicsByName(ctx, "")
}

// GetTopicsByName retrieves topics filtered by name (empty string for all)
func (c *Client) GetTopicsByName(ctx context.Context, name string) ([]TopicsResponse, error) {
	return paginate(ctx, func(limit, offset int) (TopicsResponse, int, error) {
		resp, err := c.fetchTopics(ctx, name, limit, offset)
		return resp, len(resp.Topics), err
	})
}

// fetchTopics is an internal helper to fetch topics with optional name filtering
func (c *Client) fetchTopics(ctx context.Context, name string, requestLimit, offset int) (TopicsResponse, error) {
	filter := buildFilter("name", name)
	reqURL := paginatedURL(c.BaseURL+"/topics", requestLimit, offset, filter)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return TopicsResponse{}, err
	}

	defer func() { _ = httpResponse.Body.Close() }()

	var topics TopicsResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&topics)
	if err != nil {
		return TopicsResponse{}, err
	}

	return topics, nil
}

// GetTopic retrieves a single topic by ID from the DCI API
func (c *Client) GetTopic(ctx context.Context, topicID string) (*TopicResponse, error) {
	reqURL := fmt.Sprintf("%s/topics/%s", c.BaseURL, topicID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting topic: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var topic TopicResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&topic); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &topic, nil
}

// CreateTopic creates a new topic in DCI
func (c *Client) CreateTopic(ctx context.Context, name, productID string, componentTypes []string) (*TopicResponse, error) {
	reqBody := CreateTopicRequest{
		Name:           name,
		ProductID:      productID,
		ComponentTypes: componentTypes,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	httpResponse, err := c.doJSON(ctx, http.MethodPost, c.BaseURL+"/topics", jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating topic: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response TopicResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateTopic updates an existing topic in DCI
func (c *Client) UpdateTopic(ctx context.Context, topicID string, updates UpdateTopicRequest) (*TopicResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/topics/%s", c.BaseURL, topicID)
	httpResponse, err := c.doJSON(ctx, http.MethodPut, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating topic: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response TopicResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteTopic deletes a topic from DCI
func (c *Client) DeleteTopic(ctx context.Context, topicID string) error {
	reqURL := fmt.Sprintf("%s/topics/%s", c.BaseURL, topicID)
	httpResponse, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting topic: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return formatHTTPError(httpResponse.StatusCode, body)
	}

	return nil
}

// GetTopicComponents retrieves all components for a specific topic using the topic components endpoint
func (c *Client) GetTopicComponents(ctx context.Context, topicID string) ([]ComponentsResponse, error) {
	return paginate(ctx, func(limit, offset int) (ComponentsResponse, int, error) {
		resp, err := c.fetchTopicComponents(ctx, topicID, limit, offset)
		return resp, len(resp.Components), err
	})
}

// fetchTopicComponents is an internal helper to fetch components for a topic with pagination
func (c *Client) fetchTopicComponents(ctx context.Context, topicID string, requestLimit, offset int) (ComponentsResponse, error) {
	reqURL := paginatedURL(fmt.Sprintf("%s/topics/%s/components", c.BaseURL, topicID), requestLimit, offset)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return ComponentsResponse{}, fmt.Errorf("error getting topic components: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	var components ComponentsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&components); err != nil {
		return ComponentsResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	return components, nil
}

func (c *Client) GetJobs(ctx context.Context, daysBackLimit int) ([]JobsResponse, error) {
	var jobCollection []JobsResponse

	// Default values to page through the results
	requestLimit := defaultPageSize
	offset := 0

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		outOfDateRangeJobReturned := false

		jobs, err := c.fetchJobs(ctx, requestLimit, offset)
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

func (c *Client) GetJobsByDate(ctx context.Context, startDate, endDate time.Time) ([]JobsResponse, error) {
	var jobCollection []JobsResponse

	// Default values to page through the results
	requestLimit := defaultPageSize
	offset := 0

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		jobs, err := c.fetchJobs(ctx, requestLimit, offset)
		if err != nil {
			return nil, err
		}

		for _, job := range jobs.Jobs {
			// Parse the created at date
			createdAt, err := time.Parse(dateFormat, job.CreatedAt)
			if err != nil {
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
func (c *Client) GetJob(ctx context.Context, jobID string) (*JobResponse, error) {
	reqURL := fmt.Sprintf("%s/jobs/%s", c.BaseURL, jobID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting job: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var job JobResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&job); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &job, nil
}

// UpdateJob updates an existing job in DCI
func (c *Client) UpdateJob(ctx context.Context, jobID string, updates UpdateJobRequest) (*JobResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/jobs/%s", c.BaseURL, jobID)
	httpResponse, err := c.doJSON(ctx, http.MethodPut, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating job: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response JobResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteJob deletes a job from DCI
func (c *Client) DeleteJob(ctx context.Context, jobID string) error {
	reqURL := fmt.Sprintf("%s/jobs/%s", c.BaseURL, jobID)
	httpResponse, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting job: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return formatHTTPError(httpResponse.StatusCode, body)
	}

	return nil
}

// ScheduleJob schedules a job with auto-selected components for a topic
func (c *Client) ScheduleJob(ctx context.Context, topicID string) (*CreateJobResponse, error) {
	reqBody := ScheduleJobRequest{
		TopicID: topicID,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/jobs/schedule", c.BaseURL)
	httpResponse, err := c.doJSON(ctx, http.MethodPost, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error scheduling job: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response CreateJobResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetJobFiles retrieves all files for a specific job
func (c *Client) GetJobFiles(ctx context.Context, jobID string) (*FilesResponse, error) {
	reqURL := fmt.Sprintf("%s/jobs/%s/files", c.BaseURL, jobID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting job files: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response FilesResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

func (c *Client) fetchJobs(ctx context.Context, requestLimit, offset int) (JobsResponse, error) {
	// Get jobs from the API
	httpResponse, err := c.doRequest(ctx, http.MethodGet, paginatedURL(c.BaseURL+"/jobs", requestLimit, offset), nil, nil)
	if err != nil {
		return JobsResponse{}, err
	}

	defer func() { _ = httpResponse.Body.Close() }()

	// Decode the response into JobsResponse
	var jobs JobsResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&jobs)
	if err != nil {
		return JobsResponse{}, err
	}

	return jobs, nil
}

// GetComponents retrieves all components from the DCI API with pagination
func (c *Client) GetComponents(ctx context.Context) ([]ComponentsResponse, error) {
	return c.GetComponentsFiltered(ctx, "", "", "")
}

// GetComponentsByTopicID retrieves components filtered by topic ID (empty string for all)
func (c *Client) GetComponentsByTopicID(ctx context.Context, topicID string) ([]ComponentsResponse, error) {
	return c.GetComponentsFiltered(ctx, topicID, "", "")
}

// GetComponentsFiltered retrieves components with optional filters for topic, type, and name
func (c *Client) GetComponentsFiltered(ctx context.Context, topicID, componentType, name string) ([]ComponentsResponse, error) {
	return paginate(ctx, func(limit, offset int) (ComponentsResponse, int, error) {
		resp, err := c.fetchComponents(ctx, topicID, componentType, name, limit, offset)
		return resp, len(resp.Components), err
	})
}

// GetComponent retrieves a single component by ID from the DCI API
func (c *Client) GetComponent(ctx context.Context, componentID string) (*ComponentResponse, error) {
	reqURL := fmt.Sprintf("%s/components/%s", c.BaseURL, componentID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting component: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var component ComponentResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&component); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &component, nil
}

// CreateComponent creates a new component in DCI
func (c *Client) CreateComponent(ctx context.Context, name, componentType, topicID, version string) (*ComponentResponse, error) {
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

	reqURL := fmt.Sprintf("%s/components", c.BaseURL)
	httpResponse, err := c.doJSON(ctx, http.MethodPost, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating component: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response ComponentResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateComponent updates an existing component in DCI
func (c *Client) UpdateComponent(ctx context.Context, componentID string, updates UpdateComponentRequest) (*ComponentResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/components/%s", c.BaseURL, componentID)
	httpResponse, err := c.doJSON(ctx, http.MethodPut, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating component: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response ComponentResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteComponent deletes a component from DCI
func (c *Client) DeleteComponent(ctx context.Context, componentID string) error {
	reqURL := fmt.Sprintf("%s/components/%s", c.BaseURL, componentID)
	httpResponse, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting component: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return formatHTTPError(httpResponse.StatusCode, body)
	}

	return nil
}

// fetchComponents is an internal helper to fetch components with optional filtering
func (c *Client) fetchComponents(ctx context.Context, topicID, componentType, name string, requestLimit, offset int) (ComponentsResponse, error) {
	filters := []string{
		buildFilter("topic_id", topicID),
		buildFilter("type", componentType),
		buildFilter("name", name),
	}

	reqURL := paginatedURL(c.BaseURL+"/components", requestLimit, offset, filters...)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return ComponentsResponse{}, err
	}

	defer func() { _ = httpResponse.Body.Close() }()

	var components ComponentsResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&components)
	if err != nil {
		return ComponentsResponse{}, err
	}

	return components, nil
}

// CreateJob creates a new job in DCI
func (c *Client) CreateJob(ctx context.Context, topicID string, componentIDs []string, comment string) (*CreateJobResponse, error) {
	reqBody := CreateJobRequest{
		TopicID:    topicID,
		Components: componentIDs,
		Comment:    comment,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	httpResponse, err := c.doJSON(ctx, http.MethodPost, c.BaseURL+"/jobs", jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating job: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response CreateJobResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateJobState updates the state of a job (pre-run, running, success, failure, etc.)
func (c *Client) UpdateJobState(ctx context.Context, jobID string, status JobState, comment string) (*JobStateResponse, error) {
	reqBody := UpdateJobStateRequest{
		JobID:   jobID,
		Status:  string(status),
		Comment: comment,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/jobstates", c.BaseURL)
	httpResponse, err := c.doJSON(ctx, http.MethodPost, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating job state: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response JobStateResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetJobStates retrieves job states, optionally filtered by job ID
// fetchJobStates is an internal helper to fetch job states with optional job ID filtering
func (c *Client) fetchJobStates(ctx context.Context, jobID string, requestLimit, offset int) (JobStatesResponse, error) {
	var filter string
	if jobID != "" {
		filter = "job_id:" + jobID
	}

	reqURL := paginatedURL(c.BaseURL+"/jobstates", requestLimit, offset, filter)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return JobStatesResponse{}, err
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return JobStatesResponse{}, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response JobStatesResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return JobStatesResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}

// GetJobStates retrieves job states with pagination support
func (c *Client) GetJobStates(ctx context.Context, jobID string) ([]JobStatesResponse, error) {
	return paginate(ctx, func(limit, offset int) (JobStatesResponse, int, error) {
		resp, err := c.fetchJobStates(ctx, jobID, limit, offset)
		if err != nil {
			return JobStatesResponse{}, 0, err
		}
		return resp, len(resp.JobStates), nil
	})
}

// GetFile downloads a file by ID from DCI
func (c *Client) GetFile(ctx context.Context, fileID string) ([]byte, string, error) {
	reqURL := fmt.Sprintf("%s/files/%s", c.BaseURL, fileID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error getting file: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, "", formatHTTPError(httpResponse.StatusCode, body)
	}

	content, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading file content: %w", err)
	}

	contentType := httpResponse.Header.Get("Content-Type")
	return content, contentType, nil
}

// DeleteFile deletes a file from DCI
func (c *Client) DeleteFile(ctx context.Context, fileID string) error {
	reqURL := fmt.Sprintf("%s/files/%s", c.BaseURL, fileID)
	httpResponse, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return formatHTTPError(httpResponse.StatusCode, body)
	}

	return nil
}

// UploadFile uploads a file (e.g., test results) to a job in DCI
func (c *Client) UploadFile(ctx context.Context, jobID, filePath, mimeType string) (*UploadFileResponse, error) {
	// Read the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	fileName := filepath.Base(filePath)

	return c.uploadFileContent(ctx, c.BaseURL+"/files", fileContent, jobID, fileName, mimeType)
}

// UploadFileContent uploads file content directly (without reading from disk) to a job in DCI
func (c *Client) UploadFileContent(ctx context.Context, jobID, fileName, mimeType string, content []byte) (*UploadFileResponse, error) {
	return c.uploadFileContent(ctx, c.BaseURL+"/files", content, jobID, fileName, mimeType)
}

// uploadFileContent is the shared implementation for file uploads
func (c *Client) uploadFileContent(ctx context.Context, reqURL string, content []byte, jobID, fileName, mimeType string) (*UploadFileResponse, error) {
	headers := map[string]string{
		"DCI-JOB-ID":    jobID,
		"DCI-NAME":      fileName,
		"DCI-MIME":      mimeType,
		"Content-Type":  "application/octet-stream",
	}

	httpResponse, err := c.doRequest(ctx, http.MethodPost, reqURL, content, headers)
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response UploadFileResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetRemoteCIs retrieves all remote CIs from DCI
func (c *Client) GetRemoteCIs(ctx context.Context) (*RemoteCIsResponse, error) {
	reqURL := fmt.Sprintf("%s/remotecis", c.BaseURL)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting remote CIs: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response RemoteCIsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetRemoteCI retrieves a specific remote CI by ID
func (c *Client) GetRemoteCI(ctx context.Context, remoteciID string) (*RemoteCIResponse, error) {
	reqURL := fmt.Sprintf("%s/remotecis/%s", c.BaseURL, remoteciID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting remote CI: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response RemoteCIResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CreateRemoteCI creates a new remote CI in DCI
func (c *Client) CreateRemoteCI(ctx context.Context, name, teamID string) (*RemoteCIResponse, error) {
	reqBody := CreateRemoteCIRequest{
		Name:   name,
		TeamID: teamID,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/remotecis", c.BaseURL)
	httpResponse, err := c.doJSON(ctx, http.MethodPost, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating remote CI: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response RemoteCIResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateRemoteCI updates an existing remote CI in DCI
func (c *Client) UpdateRemoteCI(ctx context.Context, remoteciID string, updates UpdateRemoteCIRequest) (*RemoteCIResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/remotecis/%s", c.BaseURL, remoteciID)
	httpResponse, err := c.doJSON(ctx, http.MethodPut, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating remote CI: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response RemoteCIResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteRemoteCI deletes a remote CI from DCI
func (c *Client) DeleteRemoteCI(ctx context.Context, remoteciID string) error {
	reqURL := fmt.Sprintf("%s/remotecis/%s", c.BaseURL, remoteciID)
	httpResponse, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting remote CI: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return formatHTTPError(httpResponse.StatusCode, body)
	}

	return nil
}

// GetTeams retrieves all teams from DCI
func (c *Client) GetTeams(ctx context.Context) ([]TeamsResponse, error) {
	return c.GetTeamsFiltered(ctx, "")
}

// GetTeamsFiltered retrieves teams with optional name filter
func (c *Client) GetTeamsFiltered(ctx context.Context, name string) ([]TeamsResponse, error) {
	return paginate(ctx, func(limit, offset int) (TeamsResponse, int, error) {
		resp, err := c.fetchTeams(ctx, name, limit, offset)
		return resp, len(resp.Teams), err
	})
}

// fetchTeams is an internal helper to fetch teams with optional name filtering
func (c *Client) fetchTeams(ctx context.Context, name string, requestLimit, offset int) (TeamsResponse, error) {
	filter := buildFilter("name", name)
	reqURL := paginatedURL(c.BaseURL+"/teams", requestLimit, offset, filter)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return TeamsResponse{}, fmt.Errorf("error getting teams: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return TeamsResponse{}, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response TeamsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return TeamsResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}

// GetTeam retrieves a specific team by ID
func (c *Client) GetTeam(ctx context.Context, teamID string) (*TeamResponse, error) {
	reqURL := fmt.Sprintf("%s/teams/%s", c.BaseURL, teamID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting team: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response TeamResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CreateTeam creates a new team in DCI
func (c *Client) CreateTeam(ctx context.Context, name string) (*TeamResponse, error) {
	reqBody := CreateTeamRequest{
		Name: name,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/teams", c.BaseURL)
	httpResponse, err := c.doJSON(ctx, http.MethodPost, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating team: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response TeamResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateTeam updates an existing team in DCI
func (c *Client) UpdateTeam(ctx context.Context, teamID string, updates UpdateTeamRequest) (*TeamResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/teams/%s", c.BaseURL, teamID)
	httpResponse, err := c.doJSON(ctx, http.MethodPut, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating team: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response TeamResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteTeam deletes a team from DCI
func (c *Client) DeleteTeam(ctx context.Context, teamID string) error {
	reqURL := fmt.Sprintf("%s/teams/%s", c.BaseURL, teamID)
	httpResponse, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting team: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return formatHTTPError(httpResponse.StatusCode, body)
	}

	return nil
}

// GetUsers retrieves all users from DCI
func (c *Client) GetUsers(ctx context.Context) ([]UsersResponse, error) {
	return c.GetUsersFiltered(ctx, "")
}

// GetUsersFiltered retrieves users with optional name filter
func (c *Client) GetUsersFiltered(ctx context.Context, name string) ([]UsersResponse, error) {
	return paginate(ctx, func(limit, offset int) (UsersResponse, int, error) {
		resp, err := c.fetchUsers(ctx, name, limit, offset)
		return resp, len(resp.Users), err
	})
}

// fetchUsers is an internal helper to fetch users with optional name filtering
func (c *Client) fetchUsers(ctx context.Context, name string, requestLimit, offset int) (UsersResponse, error) {
	filter := buildFilter("name", name)
	reqURL := paginatedURL(c.BaseURL+"/users", requestLimit, offset, filter)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return UsersResponse{}, fmt.Errorf("error getting users: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return UsersResponse{}, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response UsersResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return UsersResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}

// GetUser retrieves a specific user by ID
func (c *Client) GetUser(ctx context.Context, userID string) (*UserResponse, error) {
	reqURL := fmt.Sprintf("%s/users/%s", c.BaseURL, userID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response UserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CreateUser creates a new user in DCI
func (c *Client) CreateUser(ctx context.Context, name, email, fullname, teamID, password string) (*UserResponse, error) {
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

	reqURL := fmt.Sprintf("%s/users", c.BaseURL)
	httpResponse, err := c.doJSON(ctx, http.MethodPost, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response UserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// UpdateUser updates an existing user in DCI
func (c *Client) UpdateUser(ctx context.Context, userID string, updates UpdateUserRequest) (*UserResponse, error) {
	jsonBody, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	reqURL := fmt.Sprintf("%s/users/%s", c.BaseURL, userID)
	httpResponse, err := c.doJSON(ctx, http.MethodPut, reqURL, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response UserResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// DeleteUser deletes a user from DCI
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	reqURL := fmt.Sprintf("%s/users/%s", c.BaseURL, userID)
	httpResponse, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return formatHTTPError(httpResponse.StatusCode, body)
	}

	return nil
}

// GetProducts retrieves all products from DCI
func (c *Client) GetProducts(ctx context.Context) (*ProductsResponse, error) {
	reqURL := fmt.Sprintf("%s/products", c.BaseURL)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting products: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response ProductsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetProduct retrieves a specific product by ID
func (c *Client) GetProduct(ctx context.Context, productID string) (*ProductResponse, error) {
	reqURL := fmt.Sprintf("%s/products/%s", c.BaseURL, productID)
	httpResponse, err := c.doRequest(ctx, http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, formatHTTPError(httpResponse.StatusCode, body)
	}

	var response ProductResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}
