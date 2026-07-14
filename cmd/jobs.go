package cmd

import (
	"encoding/json"
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(getJobIDFlag, "job"); err != nil {
			return err
		}


		printStatus("Getting job with ID: %s\n", getJobIDFlag)

		response, err := dciClient.GetJob(cmd.Context(), getJobIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get job: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printJobJSON(response)
		}

		printJobStdout(response)

		return nil
	},
}

var updateJobCmd = &cobra.Command{
	Use:   "update-job",
	Short: "Update an existing job in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(updateJobIDFlag, "job"); err != nil {
			return err
		}


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

		if dryRunFlag {
			printStatus("[DRY RUN] Would update job: id=%s, comment=%q, tags=%v\n", updateJobIDFlag, updateJobComment, updates.Tags)
			return nil
		}

		printStatus("Updating job: %s\n", updateJobIDFlag)

		response, err := dciClient.UpdateJob(cmd.Context(), updateJobIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update job: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printJobJSON(response)
		}

		printStatus("Job updated successfully!")
		printJobStdout(response)

		return nil
	},
}

var deleteJobCmd = &cobra.Command{
	Use:   "delete-job",
	Short: "Delete a job from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(deleteJobIDFlag, "job"); err != nil {
			return err
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would delete job: id=%s\n", deleteJobIDFlag)
			return nil
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("job", deleteJobIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			printStatus("Deletion canceled")
			return nil
		}

		printStatus("Deleting job: %s\n", deleteJobIDFlag)

		err = dciClient.DeleteJob(cmd.Context(), deleteJobIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete job: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteJobIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			printStatus("Job deleted successfully!")
		}

		return nil
	},
}

var scheduleJobCmd = &cobra.Command{
	Use:   "schedule-job",
	Short: "Schedule a job with auto-selected components",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(scheduleJobTopicID, "topic"); err != nil {
			return err
		}


		if dryRunFlag {
			printStatus("[DRY RUN] Would schedule job: topic-id=%s\n", scheduleJobTopicID)
			return nil
		}

		printStatus("Scheduling job for topic: %s\n", scheduleJobTopicID)

		response, err := dciClient.ScheduleJob(cmd.Context(), scheduleJobTopicID)
		if err != nil {
			return fmt.Errorf("failed to schedule job: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printCreateJobJSON(response)
		}

		printStatus("Job scheduled successfully!")
		printCreateJobStdout(response)

		return nil
	},
}

var getJobFilesCmd = &cobra.Command{
	Use:   "job-files",
	Short: "Get all files for a specific job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(jobFilesIDFlag, "job"); err != nil {
			return err
		}


		printStatus("Getting files for job ID: %s\n", jobFilesIDFlag)

		response, err := dciClient.GetJobFiles(cmd.Context(), jobFilesIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get job files: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printFilesJSON(response)
		}

		printFilesStdout(response)

		return nil
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

func printJobJSON(response *lib.JobResponse) error {
	lib.RedactJob(&response.Job)
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
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

func printFilesJSON(response *lib.FilesResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(getJobCmd)
	rootCmd.AddCommand(updateJobCmd)
	rootCmd.AddCommand(deleteJobCmd)
	rootCmd.AddCommand(scheduleJobCmd)
	rootCmd.AddCommand(getJobFilesCmd)

	// get job flags
	getJobCmd.PersistentFlags().StringVar(&getJobIDFlag, "id", "", "Job ID")
	_ = getJobCmd.MarkPersistentFlagRequired("id")

	// update job flags
	updateJobCmd.PersistentFlags().StringVar(&updateJobIDFlag, "id", "", "Job ID to update")
	_ = updateJobCmd.MarkPersistentFlagRequired("id")
	updateJobCmd.PersistentFlags().StringVar(&updateJobComment, "comment", "", "New comment for the job")
	updateJobCmd.PersistentFlags().StringVar(&updateJobTags, "tags", "", "Comma-separated list of tags")

	// delete job flags
	deleteJobCmd.PersistentFlags().StringVar(&deleteJobIDFlag, "id", "", "Job ID to delete")
	_ = deleteJobCmd.MarkPersistentFlagRequired("id")

	// schedule job flags
	scheduleJobCmd.PersistentFlags().StringVar(&scheduleJobTopicID, "topic-id", "", "Topic ID for the job")
	_ = scheduleJobCmd.MarkPersistentFlagRequired("topic-id")

	// get job files flags
	getJobFilesCmd.PersistentFlags().StringVar(&jobFilesIDFlag, "id", "", "Job ID")
	_ = getJobFilesCmd.MarkPersistentFlagRequired("id")
}
