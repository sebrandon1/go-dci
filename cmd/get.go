package cmd

import (
	"encoding/json"
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
var outputFormat string

const (
	dateFormat = "2006-01-02T15:04:05.999999"
	tnfRegex   = `v(\d+\.)?(\d+\.)?(\*|\d+)$`
)

var getJobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Get all jobs with a specific age in days",
	Run: func(cmd *cobra.Command, args []string) {
		// Pull the access key and secret key from the config file
		accessKey = GetConfigValue("accesskey")
		secretKey = GetConfigValue("secretkey")

		var jsonOutput lib.JobsJsonOutput

		// Initialize the counters
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

		if outputFormat != "json" {
			fmt.Printf("Getting all jobs from DCI that are %s days old\n", ageInDays)
		}

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

				for _, c := range j.Components {
					if strings.Contains(c.Name, "cnf-certification-test") {
						// get the commit/version from the component name
						commit := strings.Split(c.Name, " ")[1]

						// find out how long ago this job ran
						daysAgo, _ := time.Parse(dateFormat, j.CreatedAt)
						daysSince := time.Since(daysAgo).Hours() / 24
						if outputFormat != "json" {
							fmt.Printf("Job ID: %s  -  TNF Version: %s (Days Since: %f)\n", j.ID, commit, daysSince)
						}

						jo := lib.JsonTNFInfo{
							ID:         j.ID,
							TNFVersion: commit,
						}
						jsonOutput.Jobs = append(jsonOutput.Jobs, jo)

						tnfJobsCtr++
					}
				}
			}
		}

		if outputFormat != "json" {
			fmt.Printf("Total TNF Jobs: %d\n", tnfJobsCtr)
			fmt.Printf("Total DCI Jobs: %d\n", totalJobsCtr)
			fmt.Printf("Total go-dci runtime: %v\n", time.Since(startRun))
		} else {
			// marshal the jsonOutput
			jsonOutputBytes, err := json.Marshal(jsonOutput)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonOutputBytes))
		}
	},
}

func init() {
	rootCmd.AddCommand(getJobsCmd)

	// Bind the access key and secret key to the variables
	getJobsCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
	getJobsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "stdout", "Output format (json) - default is stdout")
}
