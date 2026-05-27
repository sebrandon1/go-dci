package cmd

import (
	"encoding/json"
	"fmt"
	"log"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			if getJobStatesJobIDFlag != "" {
				fmt.Printf("Getting job states for job ID: %s\n", getJobStatesJobIDFlag)
			} else {
				fmt.Println("Getting all job states...")
			}
		}

		response, err := client.GetJobStates(getJobStatesJobIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get job states: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			printJobStatesJSON(response)
		} else {
			printJobStatesStdout(response)
		}

		return nil
	},
}

func printJobStatesStdout(response *lib.JobStatesResponse) {
	if len(response.JobStates) == 0 {
		fmt.Println("No job states found.")
		return
	}
	fmt.Println("---")
	for _, js := range response.JobStates {
		fmt.Printf("ID: %s | Job ID: %s | Status: %s | Created: %s\n",
			js.ID, js.JobID, js.Status, js.CreatedAt)
	}
	fmt.Printf("Total Job States: %d\n", len(response.JobStates))
}

func printJobStatesJSON(response *lib.JobStatesResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(getJobStatesCmd)

	// get job states flags
	getJobStatesCmd.PersistentFlags().StringVar(&getJobStatesJobIDFlag, "job-id", "", "Filter by Job ID (optional)")
	getJobStatesCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
