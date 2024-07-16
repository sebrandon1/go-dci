package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	signerv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

const (
	// https://doc.distributed-ci.io/python-dciauth/#using-postman
	DCIURL      = "https://api.distributed-ci.io/api/v1"
	awsRegion   = "BHS3"
	serviceName = "api"
	dateFormat  = "2006-01-02T15:04:05.999999"
	maxRecords  = 50000
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

		defer httpResponse.Body.Close()

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

		httpResponse, err := HttpGetWithAWSAuth(c.BaseURL+"/jobs", awsRegion, serviceName, c.AccessKey, c.SecretKey, requestLimit, offset)
		if err != nil {
			fmt.Printf("Error getting jobs: %s\n", err)
			return nil, err
		}

		defer httpResponse.Body.Close()

		var jobs JobsResponse
		err = json.NewDecoder(httpResponse.Body).Decode(&jobs)
		if err != nil {
			fmt.Printf("Error decoding the response: %s\n", err)
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

func HttpGetWithAWSAuth(url, region, serviceName, accessKey, secretKey string, limit, offset int) (*http.Response, error) {
	// Created a signed request with aws-sdk-go
	// https://docs.aws.amazon.com/sdk-for-go/api/aws/signer/v4/

	// Create a new signer
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	signer := signerv4.NewSigner(creds)

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
	_, err = signer.Sign(req, nil, serviceName, region, time.Now())
	if err != nil {
		return nil, err
	}

	// Send the request
	client := &http.Client{}

	// Perform the requests and adjust the offset based on the response
	return client.Do(req)
}
