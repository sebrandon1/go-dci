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
var topicID string

const (
	dateFormat         = "2006-01-02T15:04:05.999999"
	OutputFormatJSON   = "json"
	OutputFormatStdout = "stdout"
)

var (
	ocpVersionsToLookFor = []string{"4.12", "4.13", "4.14", "4.15", "4.16", "4.17", "4.18", "4.19", "4.20"}
)

var getTopicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "Get all topics from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Println("Getting all topics from DCI")
		}

		topicsResponses, err := client.GetTopics()
		if err != nil {
			fmt.Printf("failed to get topics: %v\n", err)
			return
		}

		totalTopics := 0
		for _, tr := range topicsResponses {
			totalTopics += len(tr.Topics)
		}

		if outputFormat == OutputFormatJSON {
			printTopicsJSON(topicsResponses)
		} else {
			printTopicsStdout(topicsResponses)
			fmt.Printf("Total Topics: %d\n", totalTopics)
		}
	},
}

var getComponentTypesCmd = &cobra.Command{
	Use:   "componenttypes",
	Short: "Get all component types from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Println("Getting all component types from DCI")
		}

		componentTypesResponses, err := client.GetComponentTypes()
		if err != nil {
			fmt.Printf("failed to get component types: %v\n", err)
			return
		}

		totalComponentTypes := 0
		for _, ctr := range componentTypesResponses {
			totalComponentTypes += len(ctr.ComponentTypes)
		}

		if outputFormat == OutputFormatJSON {
			printComponentTypesJSON(componentTypesResponses)
		} else {
			printComponentTypesStdout(componentTypesResponses)
			fmt.Printf("Total Component Types: %d\n", totalComponentTypes)
		}
	},
}

var getIdentityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Verify authentication and display current identity information",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		identity, err := client.GetIdentity()
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printIdentityJSON(identity)
		} else {
			printIdentityStdout(identity)
		}
	},
}

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
			fmt.Printf("Error: invalid age value '%s': %v\n", ageInDays, err)
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		jobsResponses, err := client.GetJobs(daysBackLimit)
		if err != nil {
			fmt.Printf("Error getting jobs: %v\n", err)
			return
		}

		ocpVersionCount := countOcpVersions(jobsResponses)

		if outputFormat != OutputFormatJSON {
			printOcpVersionCount(ocpVersionCount)
		} else {
			printOcpVersionCountJSON(ocpVersionCount)
		}
	},
}

var getComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Get all components, optionally filtered by topic ID",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		var componentsResponses []lib.ComponentsResponse

		if topicID != "" {
			if outputFormat != OutputFormatJSON {
				fmt.Printf("Getting components for topic ID: %s\n", topicID)
			}
			componentsResponses, err = client.GetComponentsByTopicID(topicID)
		} else {
			if outputFormat != OutputFormatJSON {
				fmt.Println("Getting all components from DCI")
			}
			componentsResponses, err = client.GetComponents()
		}

		if err != nil {
			fmt.Printf("failed to get components: %v\n", err)
			return
		}

		totalComponents := 0
		for _, cr := range componentsResponses {
			totalComponents += len(cr.Components)
		}

		if outputFormat == OutputFormatJSON {
			printComponentsJSON(componentsResponses)
		} else {
			printComponentsStdout(componentsResponses)
			fmt.Printf("Total Components: %d\n", totalComponents)
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
			fmt.Printf("Error: invalid age value '%s': %v\n", ageInDays, err)
			return
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
				fmt.Printf("Error marshaling JSON output: %v\n", err)
				return
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

func printComponentsStdout(componentsResponses []lib.ComponentsResponse) {
	for _, cr := range componentsResponses {
		for _, c := range cr.Components {
			fmt.Printf("ID: %s | Name: %s | Type: %s | Version: %s | TopicID: %s\n",
				c.ID, c.Name, c.Type, c.Version, c.TopicID)
		}
	}
}

func printComponentsJSON(componentsResponses []lib.ComponentsResponse) {
	// Flatten all components into a single slice
	var allComponents []lib.Components
	for _, cr := range componentsResponses {
		allComponents = append(allComponents, cr.Components...)
	}

	output := struct {
		Components []lib.Components `json:"components"`
	}{
		Components: allComponents,
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printIdentityStdout(identity *lib.IdentityResponse) {
	fmt.Println("Authentication successful!")
	fmt.Println("---")
	fmt.Printf("ID:       %s\n", identity.Identity.ID)
	fmt.Printf("Name:     %s\n", identity.Identity.Name)
	fmt.Printf("Type:     %s\n", identity.Identity.Type)
	if identity.Identity.Email != "" {
		fmt.Printf("Email:    %s\n", identity.Identity.Email)
	}
	if identity.Identity.Fullname != "" {
		fmt.Printf("Fullname: %s\n", identity.Identity.Fullname)
	}
	if identity.Identity.TeamName != "" {
		fmt.Printf("Team:     %s\n", identity.Identity.TeamName)
	}
	if identity.Identity.TeamID != "" {
		fmt.Printf("Team ID:  %s\n", identity.Identity.TeamID)
	}
	if identity.Identity.State != "" {
		fmt.Printf("State:    %s\n", identity.Identity.State)
	}
}

func printIdentityJSON(identity *lib.IdentityResponse) {
	jsonBytes, err := json.Marshal(identity)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printComponentTypesStdout(componentTypesResponses []lib.ComponentTypesResponse) {
	for _, ctr := range componentTypesResponses {
		for _, ct := range ctr.ComponentTypes {
			fmt.Printf("ID: %s | Name: %s | State: %s\n", ct.ID, ct.Name, ct.State)
		}
	}
}

func printComponentTypesJSON(componentTypesResponses []lib.ComponentTypesResponse) {
	// Flatten all component types into a single slice
	var allComponentTypes []lib.ComponentType
	for _, ctr := range componentTypesResponses {
		allComponentTypes = append(allComponentTypes, ctr.ComponentTypes...)
	}

	output := struct {
		ComponentTypes []lib.ComponentType `json:"componenttypes"`
	}{
		ComponentTypes: allComponentTypes,
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printTopicsStdout(topicsResponses []lib.TopicsResponse) {
	for _, tr := range topicsResponses {
		for _, t := range tr.Topics {
			fmt.Printf("ID: %s | Name: %s | Product: %s | State: %s\n",
				t.ID, t.Name, t.Product.Name, t.State)
		}
	}
}

func printTopicsJSON(topicsResponses []lib.TopicsResponse) {
	// Flatten all topics into a single slice using the Topic struct
	type Topic struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		ProductID     string   `json:"product_id"`
		ProductName   string   `json:"product_name"`
		State         string   `json:"state"`
		ExportControl bool     `json:"export_control"`
		CreatedAt     string   `json:"created_at,omitempty"`
		UpdatedAt     string   `json:"updated_at,omitempty"`
	}

	var allTopics []Topic
	for _, tr := range topicsResponses {
		for _, t := range tr.Topics {
			allTopics = append(allTopics, Topic{
				ID:            t.ID,
				Name:          t.Name,
				ProductID:     t.ProductID,
				ProductName:   t.Product.Name,
				State:         t.State,
				ExportControl: t.ExportControl,
				CreatedAt:     t.CreatedAt,
				UpdatedAt:     t.UpdatedAt,
			})
		}
	}

	output := struct {
		Topics []Topic `json:"topics"`
	}{
		Topics: allTopics,
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(getJobsCmd)
	rootCmd.AddCommand(getOcpCountCmd)
	rootCmd.AddCommand(getComponentsCmd)
	rootCmd.AddCommand(getIdentityCmd)
	rootCmd.AddCommand(getComponentTypesCmd)
	rootCmd.AddCommand(getTopicsCmd)

	getJobsCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
	getJobsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	getOcpCountCmd.PersistentFlags().StringVarP(&ageInDays, "age", "d", "", "Age in days")
	getOcpCountCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	getComponentsCmd.PersistentFlags().StringVarP(&topicID, "topic", "t", "", "Filter components by topic ID")
	getComponentsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	getIdentityCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	getComponentTypesCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	getTopicsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
