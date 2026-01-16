package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for topic command flags
var (
	getTopicIDFlag           string
	createTopicName          string
	createTopicProductID     string
	createTopicComponentTypes string
	updateTopicName          string
	deleteTopicIDFlag        string
	topicComponentsIDFlag    string
)

var getTopicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Get a specific topic by ID",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if getTopicIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting topic with ID: %s\n", getTopicIDFlag)
		}

		response, err := client.GetTopic(getTopicIDFlag)
		if err != nil {
			fmt.Printf("Failed to get topic: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printTopicJSON(response)
		} else {
			printTopicStdout(response)
		}
	},
}

var createTopicCmd = &cobra.Command{
	Use:   "create-topic",
	Short: "Create a new topic in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if createTopicName == "" {
			fmt.Println("Error: --name is required")
			return
		}

		if createTopicProductID == "" {
			fmt.Println("Error: --product-id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		// Parse component types if provided
		var componentTypes []string
		if createTopicComponentTypes != "" {
			componentTypes = strings.Split(createTopicComponentTypes, ",")
			for i := range componentTypes {
				componentTypes[i] = strings.TrimSpace(componentTypes[i])
			}
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Creating topic: %s\n", createTopicName)
		}

		response, err := client.CreateTopic(createTopicName, createTopicProductID, componentTypes)
		if err != nil {
			fmt.Printf("Failed to create topic: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printTopicJSON(response)
		} else {
			fmt.Println("Topic created successfully!")
			printTopicStdout(response)
		}
	},
}

var updateTopicCmd = &cobra.Command{
	Use:   "update-topic",
	Short: "Update an existing topic in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if getTopicIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateTopicRequest{}
		if updateTopicName != "" {
			updates.Name = updateTopicName
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Updating topic: %s\n", getTopicIDFlag)
		}

		response, err := client.UpdateTopic(getTopicIDFlag, updates)
		if err != nil {
			fmt.Printf("Failed to update topic: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printTopicJSON(response)
		} else {
			fmt.Println("Topic updated successfully!")
			printTopicStdout(response)
		}
	},
}

var deleteTopicCmd = &cobra.Command{
	Use:   "delete-topic",
	Short: "Delete a topic from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if deleteTopicIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Deleting topic: %s\n", deleteTopicIDFlag)
		}

		err = client.DeleteTopic(deleteTopicIDFlag)
		if err != nil {
			fmt.Printf("Failed to delete topic: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteTopicIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("Topic deleted successfully!")
		}
	},
}

var getTopicComponentsCmd = &cobra.Command{
	Use:   "topic-components",
	Short: "Get all components for a specific topic",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if topicComponentsIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting components for topic ID: %s\n", topicComponentsIDFlag)
		}

		componentsResponses, err := client.GetTopicComponents(topicComponentsIDFlag)
		if err != nil {
			fmt.Printf("Failed to get topic components: %v\n", err)
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

func printTopicStdout(response *lib.TopicResponse) {
	fmt.Println("---")
	fmt.Printf("ID:            %s\n", response.Topic.ID)
	fmt.Printf("Name:          %s\n", response.Topic.Name)
	fmt.Printf("Product ID:    %s\n", response.Topic.ProductID)
	fmt.Printf("State:         %s\n", response.Topic.State)
	fmt.Printf("Export Control: %t\n", response.Topic.ExportControl)
	if len(response.Topic.ComponentTypes) > 0 {
		fmt.Printf("Component Types: %s\n", strings.Join(response.Topic.ComponentTypes, ", "))
	}
	fmt.Printf("Created:       %s\n", response.Topic.CreatedAt)
	fmt.Printf("Updated:       %s\n", response.Topic.UpdatedAt)
}

func printTopicJSON(response *lib.TopicResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(getTopicCmd)
	rootCmd.AddCommand(createTopicCmd)
	rootCmd.AddCommand(updateTopicCmd)
	rootCmd.AddCommand(deleteTopicCmd)
	rootCmd.AddCommand(getTopicComponentsCmd)

	// get topic flags
	getTopicCmd.PersistentFlags().StringVar(&getTopicIDFlag, "id", "", "Topic ID (required)")
	getTopicCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create topic flags
	createTopicCmd.PersistentFlags().StringVar(&createTopicName, "name", "", "Topic name (required)")
	createTopicCmd.PersistentFlags().StringVar(&createTopicProductID, "product-id", "", "Product ID (required)")
	createTopicCmd.PersistentFlags().StringVar(&createTopicComponentTypes, "component-types", "", "Comma-separated list of component type names")
	createTopicCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update topic flags
	updateTopicCmd.PersistentFlags().StringVar(&getTopicIDFlag, "id", "", "Topic ID to update (required)")
	updateTopicCmd.PersistentFlags().StringVar(&updateTopicName, "name", "", "New topic name")
	updateTopicCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete topic flags
	deleteTopicCmd.PersistentFlags().StringVar(&deleteTopicIDFlag, "id", "", "Topic ID to delete (required)")
	deleteTopicCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// get topic components flags
	getTopicComponentsCmd.PersistentFlags().StringVar(&topicComponentsIDFlag, "id", "", "Topic ID (required)")
	getTopicComponentsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
