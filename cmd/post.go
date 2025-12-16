package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for POST command flags
var (
	createJobTopicID     string
	createJobComponents  string
	createJobComment     string
	updateJobStateJobID  string
	updateJobStateStatus string
	updateJobStateComment string
	uploadFileJobID      string
	uploadFilePath       string
	uploadFileMimeType   string
)

var createJobCmd = &cobra.Command{
	Use:   "create-job",
	Short: "Create a new job in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if createJobTopicID == "" {
			fmt.Println("Error: --topic-id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		// Parse component IDs if provided
		var componentIDs []string
		if createJobComponents != "" {
			componentIDs = strings.Split(createJobComponents, ",")
			for i := range componentIDs {
				componentIDs[i] = strings.TrimSpace(componentIDs[i])
			}
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Creating job for topic ID: %s\n", createJobTopicID)
		}

		response, err := client.CreateJob(createJobTopicID, componentIDs, createJobComment)
		if err != nil {
			fmt.Printf("Failed to create job: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printCreateJobJSON(response)
		} else {
			printCreateJobStdout(response)
		}
	},
}

var updateJobStateCmd = &cobra.Command{
	Use:   "update-job-state",
	Short: "Update the state of a job (pre-run, running, success, failure, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if updateJobStateJobID == "" {
			fmt.Println("Error: --job-id is required")
			return
		}

		if updateJobStateStatus == "" {
			fmt.Println("Error: --status is required")
			return
		}

		// Validate status
		validStatuses := []string{"new", "pre-run", "running", "post-run", "success", "failure", "killed", "error"}
		isValid := false
		for _, s := range validStatuses {
			if updateJobStateStatus == s {
				isValid = true
				break
			}
		}
		if !isValid {
			fmt.Printf("Error: invalid status '%s'. Valid values: %s\n", updateJobStateStatus, strings.Join(validStatuses, ", "))
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Updating job %s to status: %s\n", updateJobStateJobID, updateJobStateStatus)
		}

		response, err := client.UpdateJobState(updateJobStateJobID, lib.JobState(updateJobStateStatus), updateJobStateComment)
		if err != nil {
			fmt.Printf("Failed to update job state: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printJobStateJSON(response)
		} else {
			printJobStateStdout(response)
		}
	},
}

var uploadFileCmd = &cobra.Command{
	Use:   "upload-file",
	Short: "Upload a file (e.g., test results) to a job in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if uploadFileJobID == "" {
			fmt.Println("Error: --job-id is required")
			return
		}

		if uploadFilePath == "" {
			fmt.Println("Error: --file is required")
			return
		}

		// Default mime type for test results
		if uploadFileMimeType == "" {
			uploadFileMimeType = "application/junit"
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Uploading file %s to job %s\n", uploadFilePath, uploadFileJobID)
		}

		response, err := client.UploadFile(uploadFileJobID, uploadFilePath, uploadFileMimeType)
		if err != nil {
			fmt.Printf("Failed to upload file: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printUploadFileJSON(response)
		} else {
			printUploadFileStdout(response)
		}
	},
}

func printCreateJobStdout(response *lib.CreateJobResponse) {
	fmt.Println("Job created successfully!")
	fmt.Println("---")
	fmt.Printf("Job ID:    %s\n", response.Job.ID)
	fmt.Printf("Topic ID:  %s\n", response.Job.TopicID)
	fmt.Printf("Status:    %s\n", response.Job.Status)
	fmt.Printf("State:     %s\n", response.Job.State)
	fmt.Printf("Created:   %s\n", response.Job.CreatedAt)
}

func printCreateJobJSON(response *lib.CreateJobResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printJobStateStdout(response *lib.JobStateResponse) {
	fmt.Println("Job state updated successfully!")
	fmt.Println("---")
	fmt.Printf("JobState ID: %s\n", response.JobState.ID)
	fmt.Printf("Job ID:      %s\n", response.JobState.JobID)
	fmt.Printf("Status:      %s\n", response.JobState.Status)
	if response.JobState.Comment != "" {
		fmt.Printf("Comment:     %s\n", response.JobState.Comment)
	}
	fmt.Printf("Created:     %s\n", response.JobState.CreatedAt)
}

func printJobStateJSON(response *lib.JobStateResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printUploadFileStdout(response *lib.UploadFileResponse) {
	fmt.Println("File uploaded successfully!")
	fmt.Println("---")
	fmt.Printf("File ID:   %s\n", response.File.ID)
	fmt.Printf("Job ID:    %s\n", response.File.JobID)
	fmt.Printf("Name:      %s\n", response.File.Name)
	fmt.Printf("MIME Type: %s\n", response.File.Mime)
	fmt.Printf("Size:      %s bytes\n", response.File.Size)
	fmt.Printf("Created:   %s\n", response.File.CreatedAt)
}

func printUploadFileJSON(response *lib.UploadFileResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(createJobCmd)
	rootCmd.AddCommand(updateJobStateCmd)
	rootCmd.AddCommand(uploadFileCmd)

	// create-job flags
	createJobCmd.PersistentFlags().StringVar(&createJobTopicID, "topic-id", "", "Topic ID for the job (required)")
	createJobCmd.PersistentFlags().StringVar(&createJobComponents, "components", "", "Comma-separated list of component IDs")
	createJobCmd.PersistentFlags().StringVar(&createJobComment, "comment", "", "Optional comment for the job")
	createJobCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update-job-state flags
	updateJobStateCmd.PersistentFlags().StringVar(&updateJobStateJobID, "job-id", "", "Job ID to update (required)")
	updateJobStateCmd.PersistentFlags().StringVar(&updateJobStateStatus, "status", "", "New status (pre-run, running, success, failure, etc.) (required)")
	updateJobStateCmd.PersistentFlags().StringVar(&updateJobStateComment, "comment", "", "Optional comment for the state change")
	updateJobStateCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// upload-file flags
	uploadFileCmd.PersistentFlags().StringVar(&uploadFileJobID, "job-id", "", "Job ID to attach the file to (required)")
	uploadFileCmd.PersistentFlags().StringVar(&uploadFilePath, "file", "", "Path to the file to upload (required)")
	uploadFileCmd.PersistentFlags().StringVar(&uploadFileMimeType, "mime", "application/junit", "MIME type of the file (default: application/junit)")
	uploadFileCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}

