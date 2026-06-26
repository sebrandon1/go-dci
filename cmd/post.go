package cmd

import (
	"encoding/json"
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
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

		if dryRunFlag {
			printStatus("[DRY RUN] Would create job: topic-id=%s, components=%v, comment=%q\n", createJobTopicID, componentIDs, createJobComment)
			return nil
		}

		printStatus("Creating job for topic ID: %s\n", createJobTopicID)

		response, err := client.CreateJob(cmd.Context(), createJobTopicID, componentIDs, createJobComment)
		if err != nil {
			return fmt.Errorf("failed to create job: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printCreateJobJSON(response)
		}

		printCreateJobStdout(response)

		return nil
	},
}

var updateJobStateCmd = &cobra.Command{
	Use:   "update-job-state",
	Short: "Update the state of a job (pre-run, running, success, failure, etc.)",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
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
			return fmt.Errorf("invalid status '%s'. Valid values: %s", updateJobStateStatus, strings.Join(validStatuses, ", "))
		}

		client := lib.NewClient(accessKey, secretKey)

		if dryRunFlag {
			printStatus("[DRY RUN] Would update job state: job-id=%s, status=%s, comment=%q\n", updateJobStateJobID, updateJobStateStatus, updateJobStateComment)
			return nil
		}

		printStatus("Updating job %s to status: %s\n", updateJobStateJobID, updateJobStateStatus)

		response, err := client.UpdateJobState(cmd.Context(), updateJobStateJobID, lib.JobState(updateJobStateStatus), updateJobStateComment)
		if err != nil {
			return fmt.Errorf("failed to update job state: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printJobStateJSON(response)
		}

		printJobStateStdout(response)

		return nil
	},
}

var uploadFileCmd = &cobra.Command{
	Use:   "upload-file",
	Short: "Upload a file (e.g., test results) to a job in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		// Default mime type for test results
		if uploadFileMimeType == "" {
			uploadFileMimeType = "application/junit"
		}

		client := lib.NewClient(accessKey, secretKey)

		if dryRunFlag {
			printStatus("[DRY RUN] Would upload file: job-id=%s, file=%s, mime=%s\n", uploadFileJobID, uploadFilePath, uploadFileMimeType)
			return nil
		}

		printStatus("Uploading file %s to job %s\n", uploadFilePath, uploadFileJobID)

		response, err := client.UploadFile(cmd.Context(), uploadFileJobID, uploadFilePath, uploadFileMimeType)
		if err != nil {
			return fmt.Errorf("failed to upload file: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printUploadFileJSON(response)
		}

		printUploadFileStdout(response)

		return nil
	},
}

func printCreateJobStdout(response *lib.CreateJobResponse) {
	printStatus("Job created successfully!")
	fmt.Println("---")
	fmt.Printf("Job ID:    %s\n", response.Job.ID)
	fmt.Printf("Topic ID:  %s\n", response.Job.TopicID)
	fmt.Printf("Status:    %s\n", response.Job.Status)
	fmt.Printf("State:     %s\n", response.Job.State)
	fmt.Printf("Created:   %s\n", response.Job.CreatedAt)
}

func printCreateJobJSON(response *lib.CreateJobResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func printJobStateStdout(response *lib.JobStateResponse) {
	printStatus("Job state updated successfully!")
	fmt.Println("---")
	fmt.Printf("JobState ID: %s\n", response.JobState.ID)
	fmt.Printf("Job ID:      %s\n", response.JobState.JobID)
	fmt.Printf("Status:      %s\n", response.JobState.Status)
	if response.JobState.Comment != "" {
		fmt.Printf("Comment:     %s\n", response.JobState.Comment)
	}
	fmt.Printf("Created:     %s\n", response.JobState.CreatedAt)
}

func printJobStateJSON(response *lib.JobStateResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func printUploadFileStdout(response *lib.UploadFileResponse) {
	printStatus("File uploaded successfully!")
	fmt.Println("---")
	fmt.Printf("File ID:   %s\n", response.File.ID)
	fmt.Printf("Job ID:    %s\n", response.File.JobID)
	fmt.Printf("Name:      %s\n", response.File.Name)
	fmt.Printf("MIME Type: %s\n", response.File.Mime)
	fmt.Printf("Size:      %s bytes\n", response.File.Size)
	fmt.Printf("Created:   %s\n", response.File.CreatedAt)
}

func printUploadFileJSON(response *lib.UploadFileResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(createJobCmd)
	rootCmd.AddCommand(updateJobStateCmd)
	rootCmd.AddCommand(uploadFileCmd)

	// create-job flags
	createJobCmd.PersistentFlags().StringVar(&createJobTopicID, "topic-id", "", "Topic ID for the job")
	_ = createJobCmd.MarkPersistentFlagRequired("topic-id")
	createJobCmd.PersistentFlags().StringVar(&createJobComponents, "components", "", "Comma-separated list of component IDs")
	createJobCmd.PersistentFlags().StringVar(&createJobComment, "comment", "", "Optional comment for the job")
	createJobCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update-job-state flags
	updateJobStateCmd.PersistentFlags().StringVar(&updateJobStateJobID, "job-id", "", "Job ID to update")
	_ = updateJobStateCmd.MarkPersistentFlagRequired("job-id")
	updateJobStateCmd.PersistentFlags().StringVar(&updateJobStateStatus, "status", "", "New status (pre-run, running, success, failure, etc.)")
	_ = updateJobStateCmd.MarkPersistentFlagRequired("status")
	updateJobStateCmd.PersistentFlags().StringVar(&updateJobStateComment, "comment", "", "Optional comment for the state change")
	updateJobStateCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// upload-file flags
	uploadFileCmd.PersistentFlags().StringVar(&uploadFileJobID, "job-id", "", "Job ID to attach the file to")
	_ = uploadFileCmd.MarkPersistentFlagRequired("job-id")
	uploadFileCmd.PersistentFlags().StringVar(&uploadFilePath, "file", "", "Path to the file to upload")
	_ = uploadFileCmd.MarkPersistentFlagRequired("file")
	uploadFileCmd.PersistentFlags().StringVar(&uploadFileMimeType, "mime", "application/junit", "MIME type of the file (default: application/junit)")
	uploadFileCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}

