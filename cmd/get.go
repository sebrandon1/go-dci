package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

var ageInDays string
var outputFormat string

const (
	dateFormat         = "2006-01-02T15:04:05.999999"
	OutputFormatJSON   = "json"
	OutputFormatStdout = "stdout"
)

var (
	ocpVersionsToLookFor = []string{"4.12", "4.13", "4.14", "4.15", "4.16", "4.17", "4.18", "4.19", "4.20"}
)

var getOcpCountCmd = &cobra.Command{
	Use:   "ocpcount",
	Short: "Get the count of jobs for each OCP version",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting all jobs from DCI that are %s days old\n", ageInDays)
		}

		daysBackLimit, err := strconv.Atoi(ageInDays)
		if err != nil {
			panic(err)
		}

		client := lib.NewClient(accessKey, secretKey)

		jobsResponses, err := client.GetJobs(daysBackLimit)
		if err != nil {
			panic(err)
		}

		ocpVersionCount := countOcpVersions(jobsResponses)

		if outputFormat != OutputFormatJSON {
			printOcpVersionCount(ocpVersionCount)
		} else {
			printOcpVersionCountJSON(ocpVersionCount)
		}
	},
}

var getJobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Get all jobs with a specific age in days",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		var jsonOutput lib.JobsJsonOutput

		certsuiteJobsCtr := 0
		totalJobsCtr := 0
		startRun := time.Now()

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting all jobs from DCI that are %s days old\n", ageInDays)
		}

		daysBackLimit, err := strconv.Atoi(ageInDays)
		if err != nil {
			panic(err)
		}

		jobsResponses, err := client.GetJobs(daysBackLimit)
		if err != nil {
			fmt.Printf("failed to get jobs: %v\n", err)
			return
		}

		for _, job := range jobsResponses {
			for _, j := range job.Jobs {
				totalJobsCtr++

				for _, c := range j.Components {
					if strings.Contains(c.Name, "cnf-certification-test") || strings.Contains(c.Name, "certsuite") {
						commit := extractCommitVersion(c.Name)
						daysSince := calculateDaysSince(j.CreatedAt)
						if outputFormat != OutputFormatJSON {
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

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Total Certsuite Jobs: %d\n", certsuiteJobsCtr)
			fmt.Printf("Total DCI Jobs: %d\n", totalJobsCtr)
			fmt.Printf("Total go-dci runtime: %v\n", time.Since(startRun))
		} else {
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

func getCredentials() (string, string, error) {
	accessKey := GetConfigValue("accesskey")
	secretKey := GetConfigValue("secretkey")
	if accessKey == "" || secretKey == "" {
		return "", "", errors.New("access key or secret key is not set")
	}
	return accessKey, secretKey, nil
}

func countOcpVersions(jobsResponses []lib.JobsResponse) map[string]int {
	ocpVersionCount := make(map[string]int)
	for _, job := range jobsResponses {
		for _, j := range job.Jobs {
			if isCertsuiteJob(j.Components) {
				ocpVersion := findOcpVersionFromComponents(j.Components)
				if ocpVersion != "" {
					for _, v := range ocpVersionsToLookFor {
						if strings.Contains(ocpVersion, v) {
							ocpVersionCount[v]++
							break
						}
					}
				}
			}
		}
	}
	return ocpVersionCount
}

func printOcpVersionCount(ocpVersionCount map[string]int) {
	for _, ocpVersion := range ocpVersionsToLookFor {
		if count, exists := ocpVersionCount[ocpVersion]; exists {
			fmt.Printf("OCP Version: %s - Run Count: %d\n", ocpVersion, count)
		}
	}
}

func printOcpVersionCountJSON(ocpVersionCount map[string]int) {
	var jsonOutput lib.OcpJsonOutput
	for _, ocpVersion := range ocpVersionsToLookFor {
		jo := lib.JsonOcpVersionCount{
			OcpVersion: ocpVersion,
			RunCount:   ocpVersionCount[ocpVersion],
		}
		jsonOutput.OcpVersions = append(jsonOutput.OcpVersions, jo)
	}

	jsonOutputBytes, err := json.Marshal(jsonOutput)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonOutputBytes))
}

func extractCommitVersion(componentName string) string {
	parts := strings.Split(componentName, " ")
	if len(parts) > 1 {
		return parts[1]
	}
	return "unknown"
}

func calculateDaysSince(createdAt string) float64 {
	createdTime, _ := time.Parse(dateFormat, createdAt)
	return time.Since(createdTime).Hours() / 24
}

func init() {
	rootCmd.AddCommand(getJobsCmd)
	rootCmd.AddCommand(getOcpCountCmd)

	getJobsCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
	getJobsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	getOcpCountCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
	getOcpCountCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
