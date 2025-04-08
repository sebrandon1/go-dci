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
)

var (
	ocpVersionsToLookFor = []string{"4.12", "4.13", "4.14", "4.15", "4.16", "4.17"}
)

var getOcpCountCmd = &cobra.Command{
	Use:   "ocpcount",
	Short: "Get the count of jobs for each OCP version",
	Run: func(cmd *cobra.Command, args []string) {
		// Pull the access key and secret key from the config file
		accessKey = GetConfigValue("accesskey")
		secretKey = GetConfigValue("secretkey")

		// Check if the values are empty
		if accessKey == "" || secretKey == "" {
			fmt.Println("Please set the access key and secret key using the 'config set' command")
			return
		}

		var jsonOutput lib.OcpJsonOutput

		if outputFormat != "json" {
			fmt.Printf("Getting all jobs from DCI that are %s days old\n", ageInDays)
		}

		// Convert ageInDays to an integer
		daysBackLimit, err := strconv.Atoi(ageInDays)
		if err != nil {
			panic(err)
		}

		// Use the client to get all jobs from DCI
		client := lib.NewClient(accessKey, secretKey)

		jobsResponses, err := client.GetJobs(daysBackLimit)
		if err != nil {
			// fmt.Printf("responses: %v\n", jobsResponses)
			panic(err)
		}

		// Create a count of the number of jobs for each OCP version
		ocpVersionCount := make(map[string]int)
		for _, job := range jobsResponses {
			for _, j := range job.Jobs {
				if isCertsuiteJob(j.Components) {

					// Get the OCP version from the components
					ocpVersion := findOcpVersionFromComponents(j.Components)

					// Check if the OCP version is in the list of versions to look for
					if ocpVersion != "" {
						// Check if the OCP version is in the list of versions to look for
						for _, v := range ocpVersionsToLookFor {
							if strings.Contains(ocpVersion, v) {
								ocpVersionCount[v]++
							}
						}
					}
				}
			}
		}

		// Print the OCP version count
		if outputFormat != "json" {
			for _, ocpVersion := range ocpVersionsToLookFor {

				for k, v := range ocpVersionCount {
					if ocpVersion == k {
						fmt.Printf("OCP Version: %s - Run Count: %d\n", k, v)
					}
				}
			}
		} else {
			// marshal the jsonOutput
			for _, ocpVersion := range ocpVersionsToLookFor {
				jo := lib.JsonOcpVersionCount{
					OcpVersion: ocpVersion,
					RunCount:   ocpVersionCount[ocpVersion],
				}

				jsonOutput.OcpVersions = append(jsonOutput.OcpVersions, jo)
			}

			jsonOutputBytes, err := json.Marshal(jsonOutput)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonOutputBytes))
		}
	},
}

var getJobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Get all jobs with a specific age in days",
	Run: func(cmd *cobra.Command, args []string) {
		// Pull the access key and secret key from the config file
		accessKey = GetConfigValue("accesskey")
		secretKey = GetConfigValue("secretkey")

		// Check if the values are empty
		if accessKey == "" || secretKey == "" {
			fmt.Println("Please set the access key and secret key using the 'config set' command")
			return
		}

		var jsonOutput lib.JobsJsonOutput

		// Initialize the counters
		certsuiteJobsCtr := 0
		totalJobsCtr := 0
		startRun := time.Now()

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
			fmt.Printf("failed to get jobs: %v\n", err)
			return
		}

		// Print the job IDs gathered from the response
		for _, job := range jobsResponses {
			for _, j := range job.Jobs {
				totalJobsCtr++ // Keep track of the total number of jobs

				for _, c := range j.Components {
					if strings.Contains(c.Name, "cnf-certification-test") || strings.Contains(c.Name, "certsuite") {
						// get the commit/version from the component name
						commit := "unknown"
						if parts := strings.Split(c.Name, " "); len(parts) > 1 {
							commit = parts[1]
						}

						// find out how long ago this job ran
						daysAgo, _ := time.Parse(dateFormat, j.CreatedAt)
						daysSince := time.Since(daysAgo).Hours() / 24
						if outputFormat != "json" {
							fmt.Printf("Job ID: %s  -  Certsuite Version: %s (Days Since: %f)\n", j.ID, commit, daysSince)
						}

						jo := lib.JsonCertsuiteInfo{
							ID:               j.ID,
							CertsuiteVersion: commit,
						}
						jsonOutput.Jobs = append(jsonOutput.Jobs, jo)

						certsuiteJobsCtr++
					}
				}
			}
		}

		if outputFormat != "json" {
			fmt.Printf("Total Certsuite Jobs: %d\n", certsuiteJobsCtr)
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

func isCertsuiteJob(components []lib.Components) bool {
	for _, c := range components {
		if strings.Contains(c.Name, "cnf-certification-test") || strings.Contains(c.Name, "certsuite") {
			return true
		}
	}

	return false
}

func findOcpVersionFromComponents(components []lib.Components) string {
	for _, c := range components {
		if strings.Contains(c.Name, "OpenShift") {
			return strings.Split(c.Name, " ")[1]
		}
	}

	return ""
}

func init() {
	rootCmd.AddCommand(getJobsCmd)
	rootCmd.AddCommand(getOcpCountCmd)

	// Bind the access key and secret key to the variables
	getJobsCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
	getJobsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "stdout", "Output format (json) - default is stdout")

	getOcpCountCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
	getOcpCountCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "stdout", "Output format (json) - default is stdout")
}
