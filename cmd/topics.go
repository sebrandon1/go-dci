package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for topic command flags
var (
	getTopicIDFlag            string
	createTopicName           string
	createTopicProductID      string
	createTopicComponentTypes string
	updateTopicName           string
	deleteTopicIDFlag         string
	topicComponentsIDFlag     string
)

var getTopicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Get a specific topic by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(getTopicIDFlag, "topic"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting topic with ID: %s\n", getTopicIDFlag)

		response, err := client.GetTopic(cmd.Context(), getTopicIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get topic: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printTopicJSON(response)
		}

		printTopicStdout(response)

		return nil
	},
}

var createTopicCmd = &cobra.Command{
	Use:   "create-topic",
	Short: "Create a new topic in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
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

		if dryRunFlag {
			printStatus("[DRY RUN] Would create topic: name=%s, product-id=%s, component-types=%v\n", createTopicName, createTopicProductID, componentTypes)
			return nil
		}

		printStatus("Creating topic: %s\n", createTopicName)

		response, err := client.CreateTopic(cmd.Context(), createTopicName, createTopicProductID, componentTypes)
		if err != nil {
			return fmt.Errorf("failed to create topic: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printTopicJSON(response)
		}

		printStatus("Topic created successfully!")
		printTopicStdout(response)

		return nil
	},
}

var updateTopicCmd = &cobra.Command{
	Use:   "update-topic",
	Short: "Update an existing topic in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateTopicRequest{}
		if updateTopicName != "" {
			updates.Name = updateTopicName
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would update topic: id=%s, name=%s\n", getTopicIDFlag, updateTopicName)
			return nil
		}

		printStatus("Updating topic: %s\n", getTopicIDFlag)

		response, err := client.UpdateTopic(cmd.Context(), getTopicIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update topic: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printTopicJSON(response)
		}

		printStatus("Topic updated successfully!")
		printTopicStdout(response)

		return nil
	},
}

var deleteTopicCmd = &cobra.Command{
	Use:   "delete-topic",
	Short: "Delete a topic from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(deleteTopicIDFlag, "topic"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would delete topic: id=%s\n", deleteTopicIDFlag)
			return nil
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("topic", deleteTopicIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			printStatus("Deletion canceled")
			return nil
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Deleting topic: %s\n", deleteTopicIDFlag)

		err = client.DeleteTopic(cmd.Context(), deleteTopicIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete topic: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteTopicIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			printStatus("Topic deleted successfully!")
		}

		return nil
	},
}

var getTopicComponentsCmd = &cobra.Command{
	Use:   "topic-components",
	Short: "Get all components for a specific topic",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting components for topic ID: %s\n", topicComponentsIDFlag)

		componentsResponses, err := client.GetTopicComponents(cmd.Context(), topicComponentsIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get topic components: %v", err)
		}

		totalComponents := 0
		for _, cr := range componentsResponses {
			totalComponents += len(cr.Components)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentsJSON(componentsResponses)
		}

		printComponentsStdout(componentsResponses)
		fmt.Printf("Total Components: %d\n", totalComponents)

		return nil
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

func printTopicJSON(response *lib.TopicResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(getTopicCmd)
	rootCmd.AddCommand(createTopicCmd)
	rootCmd.AddCommand(updateTopicCmd)
	rootCmd.AddCommand(deleteTopicCmd)
	rootCmd.AddCommand(getTopicComponentsCmd)

	// get topic flags
	getTopicCmd.PersistentFlags().StringVar(&getTopicIDFlag, "id", "", "Topic ID")
	_ = getTopicCmd.MarkPersistentFlagRequired("id")
	getTopicCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create topic flags
	createTopicCmd.PersistentFlags().StringVar(&createTopicName, "name", "", "Topic name")
	_ = createTopicCmd.MarkPersistentFlagRequired("name")
	createTopicCmd.PersistentFlags().StringVar(&createTopicProductID, "product-id", "", "Product ID")
	_ = createTopicCmd.MarkPersistentFlagRequired("product-id")
	createTopicCmd.PersistentFlags().StringVar(&createTopicComponentTypes, "component-types", "", "Comma-separated list of component type names")
	createTopicCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update topic flags
	updateTopicCmd.PersistentFlags().StringVar(&getTopicIDFlag, "id", "", "Topic ID to update")
	_ = updateTopicCmd.MarkPersistentFlagRequired("id")
	updateTopicCmd.PersistentFlags().StringVar(&updateTopicName, "name", "", "New topic name")
	updateTopicCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete topic flags
	deleteTopicCmd.PersistentFlags().StringVar(&deleteTopicIDFlag, "id", "", "Topic ID to delete")
	_ = deleteTopicCmd.MarkPersistentFlagRequired("id")
	deleteTopicCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// get topic components flags
	getTopicComponentsCmd.PersistentFlags().StringVar(&topicComponentsIDFlag, "id", "", "Topic ID")
	_ = getTopicComponentsCmd.MarkPersistentFlagRequired("id")
	getTopicComponentsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
