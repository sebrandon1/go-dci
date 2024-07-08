package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

var accessKey string
var secretKey string
var ageInDays string

const (
	dateFormat = "2006-01-02T15:04:05.999999"
)

var getJobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Get all jobs with a specific age in days",
	Run: func(cmd *cobra.Command, args []string) {
		// Pull the access key and secret key from the config file
		accessKey = GetConfigValue("accesskey")
		secretKey = GetConfigValue("secretkey")

		tnfJobsCtr := 0
		totalJobsCtr := 0
		startRun := time.Now()

		// Check if the values are empty
		if accessKey == "" || secretKey == "" {
			fmt.Println("Please set the access key and secret key using the 'config set' command")
			return
		}

		// Use the client to get all jobs from DCI
		client := lib.NewClient(accessKey, secretKey)

		fmt.Printf("Getting all jobs from DCI that are %s days old\n", ageInDays)
		// fmt.Printf("Access Key: %s\n", accessKey)
		// fmt.Printf("Secret Key: %s\n", secretKey)

		// Convert ageInDays to an integer
		daysBackLimit, err := strconv.Atoi(ageInDays)
		if err != nil {
			panic(err)
		}

		jobsResponses, err := client.GetJobs(daysBackLimit)
		if err != nil {
			// fmt.Printf("responses: %v\n", jobsResponses)
			panic(err)
		}

		// Print the job IDs gathered from the response
		for _, job := range jobsResponses {
			for _, j := range job.Jobs {
				totalJobsCtr++ // Keep track of the total number of jobs

				// fmt.Printf("Job Name: %s - Created At: %s\n", j.Name, j.CreatedAt)

				for _, tag := range j.Tags {
					if strings.Contains(tag, "tnf") {
						// find out how long ago this job ran
						daysAgo, _ := time.Parse(dateFormat, j.CreatedAt)
						daysSince := time.Since(daysAgo).Hours() / 24
						fmt.Printf("Job ID: %s  -  TNF Version: %s (Days Since: %f)\n", j.ID, tag, daysSince)
						tnfJobsCtr++
					}
				}
			}
		}

		fmt.Printf("Total TNF Jobs: %d\n", tnfJobsCtr)
		fmt.Printf("Total DCI Jobs: %d\n", totalJobsCtr)
		fmt.Printf("Total go-dci runtime: %v\n", time.Since(startRun))
	},
}

func init() {
	rootCmd.AddCommand(getJobsCmd)

	// Bind the access key and secret key to the variables
	getJobsCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
}
