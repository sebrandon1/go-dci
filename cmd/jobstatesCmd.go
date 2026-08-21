package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for jobstate command flags
var (
	getJobStatesJobIDFlag string
)

var getJobStatesCmd = &cobra.Command{
	Use:   "jobstates",
	Short: "Get job states, optionally filtered by job ID",
	Long: `Retrieve job state transitions from DCI. Job states track the lifecycle of a job
(e.g., new → pre-run → running → success/failure). Filter by --job-id to see
only the states for a specific job.`,
	Example: `  # List all job states
  dci jobstates

  # Filter by job
  dci jobstates --job-id <job-uuid>

  # Output as JSON
  dci jobstates --job-id <job-uuid> -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if getJobStatesJobIDFlag != "" {
			if err := validateResourceID(getJobStatesJobIDFlag, "job"); err != nil {
				return err
			}
		}


		if getJobStatesJobIDFlag != "" {
			printStatus("Getting job states for job ID: %s", getJobStatesJobIDFlag)
		} else {
			printStatus("Getting all job states...")
		}

		responses, err := dciClient.GetJobStates(cmd.Context(), getJobStatesJobIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get job states: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printJobStatesJSON(responses)
		}

		printJobStatesStdout(responses)

		return nil
	},
}

func printJobStatesStdout(responses []lib.JobStatesResponse) {
	totalStates := 0
	fmt.Println("---")
	for _, response := range responses {
		for _, js := range response.JobStates {
			fmt.Printf("ID: %s | Job ID: %s | Status: %s | Created: %s\n",
				js.ID, js.JobID, js.Status, js.CreatedAt)
			totalStates++
		}
	}

	if totalStates == 0 {
		fmt.Println("No job states found.")
		return
	}
	fmt.Printf("Total Job States: %d\n", totalStates)
}

func printJobStatesJSON(responses []lib.JobStatesResponse) error {
	// Flatten all job states into a single slice
	var allJobStates []lib.JobStateEntry
	for _, response := range responses {
		allJobStates = append(allJobStates, response.JobStates...)
	}

	jsonBytes, err := json.Marshal(map[string]any{
		"jobstates": allJobStates,
		"total":     len(allJobStates),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(getJobStatesCmd)

	// get job states flags
	getJobStatesCmd.PersistentFlags().StringVar(&getJobStatesJobIDFlag, "job-id", "", "Filter by Job ID (optional)")
}
