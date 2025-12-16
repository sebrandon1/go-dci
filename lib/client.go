package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
