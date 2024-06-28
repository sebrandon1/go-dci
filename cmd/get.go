package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

var accessKey string
var secretKey string
var ageInDays string

var getJobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Get all jobs with a specific age in days",
	Run: func(cmd *cobra.Command, args []string) {
		// Use the client to get all jobs from DCI
		client := lib.NewClient(accessKey, secretKey)

		// Convert ageInDays to an integer
		daysBackLimit, err := strconv.Atoi(ageInDays)
		if err != nil {
			panic(err)
		}

		jobsResponses, err := client.GetJobs(daysBackLimit)
		if err != nil {
			panic(err)
		}

		// Print the job IDs gathered from the response
		for _, job := range jobsResponses {
			for _, j := range job.Jobs {

				fmt.Printf("Job Name: %s - Created At: %s\n", j.Name, j.CreatedAt)

				for _, tag := range j.Tags {
					if strings.Contains(tag, "tnf") {
						fmt.Printf("Job ID: %s  -  TNF Version: %s", j.ID, tag)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getJobsCmd)

	// Bind the access key and secret key to the variables
	getJobsCmd.PersistentFlags().StringVarP(&accessKey, "access-key", "a", "", "AWS access key")
	getJobsCmd.PersistentFlags().StringVarP(&secretKey, "secret-key", "s", "", "AWS secret key")
	getJobsCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
}
