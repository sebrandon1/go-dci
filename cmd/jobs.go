package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for job command flags
var (
	getJobIDFlag       string
	updateJobIDFlag    string
	updateJobComment   string
	updateJobTags      string
	deleteJobIDFlag    string
	scheduleJobTopicID string
	jobFilesIDFlag     string
)

var getJobCmd = &cobra.Command{
	Use:   "job",
	Short: "Get a specific job by ID",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if getJobIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting job with ID: %s\n", getJobIDFlag)
		}

		response, err := client.GetJob(getJobIDFlag)
		if err != nil {
			fmt.Printf("Failed to get job: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printJobJSON(response)
		} else {
			printJobStdout(response)
		}
	},
}

var updateJobCmd = &cobra.Command{
	Use:   "update-job",
	Short: "Update an existing job in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if updateJobIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateJobRequest{}
		if updateJobComment != "" {
			updates.Comment = updateJobComment
		}
		if updateJobTags != "" {
			tags := strings.Split(updateJobTags, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			updates.Tags = tags
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Updating job: %s\n", updateJobIDFlag)
		}

		response, err := client.UpdateJob(updateJobIDFlag, updates)
		if err != nil {
			fmt.Printf("Failed to update job: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printJobJSON(response)
		} else {
			fmt.Println("Job updated successfully!")
			printJobStdout(response)
		}
	},
}

var deleteJobCmd = &cobra.Command{
	Use:   "delete-job",
	Short: "Delete a job from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if deleteJobIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Deleting job: %s\n", deleteJobIDFlag)
		}

		err = client.DeleteJob(deleteJobIDFlag)
		if err != nil {
			fmt.Printf("Failed to delete job: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteJobIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("Job deleted successfully!")
		}
	},
}

var scheduleJobCmd = &cobra.Command{
	Use:   "schedule-job",
	Short: "Schedule a job with auto-selected components",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if scheduleJobTopicID == "" {
			fmt.Println("Error: --topic-id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Scheduling job for topic: %s\n", scheduleJobTopicID)
		}

		response, err := client.ScheduleJob(scheduleJobTopicID)
		if err != nil {
			fmt.Printf("Failed to schedule job: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printCreateJobJSON(response)
		} else {
			fmt.Println("Job scheduled successfully!")
			printCreateJobStdout(response)
		}
	},
}

var getJobFilesCmd = &cobra.Command{
	Use:   "job-files",
	Short: "Get all files for a specific job",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if jobFilesIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting files for job ID: %s\n", jobFilesIDFlag)
		}

		response, err := client.GetJobFiles(jobFilesIDFlag)
		if err != nil {
			fmt.Printf("Failed to get job files: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printFilesJSON(response)
		} else {
			printFilesStdout(response)
		}
	},
}

func printJobStdout(response *lib.JobResponse) {
	fmt.Println("---")
	fmt.Printf("ID:           %s\n", response.Job.ID)
	fmt.Printf("Name:         %s\n", response.Job.Name)
	fmt.Printf("Status:       %s\n", response.Job.Status)
	fmt.Printf("State:        %s\n", response.Job.State)
	fmt.Printf("Topic ID:     %s\n", response.Job.TopicID)
	fmt.Printf("RemoteCI ID:  %s\n", response.Job.RemoteciID)
	fmt.Printf("Team ID:      %s\n", response.Job.TeamID)
	if response.Job.Comment != "" {
		fmt.Printf("Comment:      %s\n", response.Job.Comment)
	}
	if len(response.Job.Tags) > 0 {
		fmt.Printf("Tags:         %s\n", strings.Join(response.Job.Tags, ", "))
	}
	fmt.Printf("Duration:     %d seconds\n", response.Job.Duration)
	fmt.Printf("Created:      %s\n", response.Job.CreatedAt)
	fmt.Printf("Updated:      %s\n", response.Job.UpdatedAt)
}

func printJobJSON(response *lib.JobResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printFilesStdout(response *lib.FilesResponse) {
	if len(response.Files) == 0 {
		fmt.Println("No files found.")
		return
	}
	fmt.Println("---")
	for _, file := range response.Files {
		fmt.Printf("ID: %s | Name: %s | MIME: %s | Size: %d bytes\n",
			file.ID, file.Name, file.Mime, file.Size)
	}
	fmt.Printf("Total Files: %d\n", len(response.Files))
}

func printFilesJSON(response *lib.FilesResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(getJobCmd)
	rootCmd.AddCommand(updateJobCmd)
	rootCmd.AddCommand(deleteJobCmd)
	rootCmd.AddCommand(scheduleJobCmd)
	rootCmd.AddCommand(getJobFilesCmd)

	// get job flags
	getJobCmd.PersistentFlags().StringVar(&getJobIDFlag, "id", "", "Job ID (required)")
	getJobCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update job flags
	updateJobCmd.PersistentFlags().StringVar(&updateJobIDFlag, "id", "", "Job ID to update (required)")
	updateJobCmd.PersistentFlags().StringVar(&updateJobComment, "comment", "", "New comment for the job")
	updateJobCmd.PersistentFlags().StringVar(&updateJobTags, "tags", "", "Comma-separated list of tags")
	updateJobCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete job flags
	deleteJobCmd.PersistentFlags().StringVar(&deleteJobIDFlag, "id", "", "Job ID to delete (required)")
	deleteJobCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// schedule job flags
	scheduleJobCmd.PersistentFlags().StringVar(&scheduleJobTopicID, "topic-id", "", "Topic ID for the job (required)")
	scheduleJobCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// get job files flags
	getJobFilesCmd.PersistentFlags().StringVar(&jobFilesIDFlag, "id", "", "Job ID (required)")
	getJobFilesCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
