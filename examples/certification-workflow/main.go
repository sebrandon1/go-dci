// Example: Certification Workflow with go-dci Library
//
// This example demonstrates a complete certification workflow:
// - Creating a new job for a topic
// - Updating job state through the lifecycle
// - Uploading test results
// - Querying job status
//
// Usage:
//
//	export GO_DCI_ACCESSKEY="your-access-key"
//	export GO_DCI_SECRETKEY="your-secret-key"
//	go run main.go --topic-id <topic-uuid>
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sebrandon1/go-dci/lib"
)

func main() {
	// Parse command line arguments
	topicID := flag.String("topic-id", "", "Topic ID for the certification job (required)")
	resultsFile := flag.String("results", "", "Path to test results file (optional)")
	dryRun := flag.Bool("dry-run", false, "Simulate workflow without creating job")
	flag.Parse()

	if *topicID == "" {
		fmt.Println("Usage: go run main.go --topic-id <topic-uuid> [--results <file>] [--dry-run]")
		fmt.Println("\nTo find your topic ID, run the basic-usage example first.")
		os.Exit(1)
	}

	// Get credentials
	accessKey := os.Getenv("GO_DCI_ACCESSKEY")
	secretKey := os.Getenv("GO_DCI_SECRETKEY")

	if accessKey == "" || secretKey == "" {
		log.Fatal("GO_DCI_ACCESSKEY and GO_DCI_SECRETKEY environment variables are required")
	}

	// Initialize client
	client := lib.NewClient(accessKey, secretKey)

	fmt.Println("=== DCI Certification Workflow ===\n")

	// Step 1: Verify authentication
	fmt.Println("1. Verifying authentication...")
	identity, err := client.GetIdentity()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	fmt.Printf("   Authenticated as: %s\n\n", identity.Identity.Name)

	// Step 2: Get topic information
	fmt.Println("2. Getting topic information...")
	topic, err := client.GetTopic(*topicID)
	if err != nil {
		log.Fatalf("Failed to get topic: %v", err)
	}
	fmt.Printf("   Topic: %s\n", topic.Topic.Name)
	fmt.Printf("   ID: %s\n", topic.Topic.ID)
	fmt.Printf("   State: %s\n\n", topic.Topic.State)

	// Step 3: Get available components for this topic
	fmt.Println("3. Getting available components...")
	componentsResp, err := client.GetTopicComponents(*topicID)
	if err != nil {
		log.Fatalf("Failed to get components: %v", err)
	}

	var componentIDs []string
	for _, resp := range componentsResp {
		for _, comp := range resp.Components {
			fmt.Printf("   - %s (%s)\n", comp.Name, comp.Type)
			componentIDs = append(componentIDs, comp.ID)
		}
	}
	fmt.Println()

	if *dryRun {
		fmt.Println("=== Dry Run Mode - Stopping Here ===")
		fmt.Println("\nIn production, the following steps would be executed:")
		fmt.Println("4. Create a new job")
		fmt.Println("5. Update job state to 'pre-run'")
		fmt.Println("6. Update job state to 'running'")
		fmt.Println("7. Upload test results (if provided)")
		fmt.Println("8. Update job state to 'success' or 'failure'")
		return
	}

	// Step 4: Create a new job
	fmt.Println("4. Creating new certification job...")
	job, err := client.CreateJob(*topicID, componentIDs, "Certification run via go-dci example")
	if err != nil {
		log.Fatalf("Failed to create job: %v", err)
	}
	fmt.Printf("   Job created: %s\n", job.Job.ID)
	fmt.Printf("   Status: %s\n\n", job.Job.Status)

	jobID := job.Job.ID

	// Step 5: Update job state to pre-run
	fmt.Println("5. Updating job state to 'pre-run'...")
	_, err = client.UpdateJobState(jobID, lib.JobStatePreRun, "Starting pre-run checks")
	if err != nil {
		log.Fatalf("Failed to update job state: %v", err)
	}
	fmt.Println("   State updated to 'pre-run'\n")

	// Simulate some work
	time.Sleep(1 * time.Second)

	// Step 6: Update job state to running
	fmt.Println("6. Updating job state to 'running'...")
	_, err = client.UpdateJobState(jobID, lib.JobStateRunning, "Running certification tests")
	if err != nil {
		log.Fatalf("Failed to update job state: %v", err)
	}
	fmt.Println("   State updated to 'running'\n")

	// Step 7: Upload test results (if provided)
	if *resultsFile != "" {
		fmt.Println("7. Uploading test results...")
		uploadResp, err := client.UploadFile(jobID, *resultsFile, "application/junit")
		if err != nil {
			log.Printf("   Warning: Failed to upload file: %v\n", err)
		} else {
			fmt.Printf("   Uploaded: %s\n", uploadResp.File.Name)
			fmt.Printf("   File ID: %s\n", uploadResp.File.ID)
		}
		fmt.Println()
	} else {
		fmt.Println("7. Skipping file upload (no results file provided)\n")
	}

	// Simulate test execution
	time.Sleep(1 * time.Second)

	// Step 8: Update job state to success
	fmt.Println("8. Updating job state to 'success'...")
	_, err = client.UpdateJobState(jobID, lib.JobStateSuccess, "All certification tests passed")
	if err != nil {
		log.Fatalf("Failed to update job state: %v", err)
	}
	fmt.Println("   State updated to 'success'\n")

	// Step 9: Verify final job state
	fmt.Println("9. Verifying final job state...")
	finalJob, err := client.GetJob(jobID)
	if err != nil {
		log.Fatalf("Failed to get job: %v", err)
	}
	fmt.Printf("   Job ID: %s\n", finalJob.Job.ID)
	fmt.Printf("   Status: %s\n", finalJob.Job.Status)
	fmt.Printf("   Created: %s\n", finalJob.Job.CreatedAt)
	fmt.Printf("   Updated: %s\n", finalJob.Job.UpdatedAt)

	// Get job states history
	fmt.Println("\n   State History:")
	states, err := client.GetJobStates(jobID)
	if err != nil {
		log.Printf("   Could not get job states: %v\n", err)
	} else {
		for _, state := range states.JobStates {
			fmt.Printf("   - %s at %s\n", state.Status, state.CreatedAt)
			if state.Comment != "" {
				fmt.Printf("     Comment: %s\n", state.Comment)
			}
		}
	}

	fmt.Println("\n=== Certification Workflow Complete ===")
}
