package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for component command flags
var (
	getComponentIDFlag       string
	createComponentName      string
	createComponentType      string
	createComponentTopicID   string
	createComponentVersion   string
	updateComponentIDFlag    string
	updateComponentName      string
	updateComponentState     string
	updateComponentVersion   string
	updateComponentTags      string
	deleteComponentIDFlag    string
)

var getComponentCmd = &cobra.Command{
	Use:   "component",
	Short: "Get a specific component by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if getComponentIDFlag == "" {
			return fmt.Errorf("--id is required")
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting component with ID: %s\n", getComponentIDFlag)
		}

		response, err := client.GetComponent(cmd.Context(), getComponentIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get component: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentJSON(response)
		}

		printComponentStdout(response)

		return nil
	},
}

var createComponentCmd = &cobra.Command{
	Use:   "create-component",
	Short: "Create a new component in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if createComponentName == "" {
			return fmt.Errorf("--name is required")
		}

		if createComponentType == "" {
			return fmt.Errorf("--type is required")
		}

		if createComponentTopicID == "" {
			return fmt.Errorf("--topic-id is required")
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Creating component: %s\n", createComponentName)
		}

		response, err := client.CreateComponent(cmd.Context(), createComponentName, createComponentType, createComponentTopicID, createComponentVersion)
		if err != nil {
			return fmt.Errorf("failed to create component: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentJSON(response)
		}

		fmt.Println("Component created successfully!")
		printComponentStdout(response)

		return nil
	},
}

var updateComponentCmd = &cobra.Command{
	Use:   "update-component",
	Short: "Update an existing component in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if updateComponentIDFlag == "" {
			return fmt.Errorf("--id is required")
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateComponentRequest{}
		if updateComponentName != "" {
			updates.Name = updateComponentName
		}
		if updateComponentState != "" {
			updates.State = lib.ResourceState(updateComponentState)
		}
		if updateComponentVersion != "" {
			updates.Version = updateComponentVersion
		}
		if updateComponentTags != "" {
			tags := strings.Split(updateComponentTags, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			updates.Tags = tags
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Updating component: %s\n", updateComponentIDFlag)
		}

		response, err := client.UpdateComponent(cmd.Context(), updateComponentIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update component: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentJSON(response)
		}

		fmt.Println("Component updated successfully!")
		printComponentStdout(response)

		return nil
	},
}

var deleteComponentCmd = &cobra.Command{
	Use:   "delete-component",
	Short: "Delete a component from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if deleteComponentIDFlag == "" {
			return fmt.Errorf("--id is required")
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Deleting component: %s\n", deleteComponentIDFlag)
		}

		err = client.DeleteComponent(cmd.Context(), deleteComponentIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete component: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteComponentIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("Component deleted successfully!")
		}

		return nil
	},
}

func printComponentStdout(response *lib.ComponentResponse) {
	fmt.Println("---")
	fmt.Printf("ID:           %s\n", response.Component.ID)
	fmt.Printf("Name:         %s\n", response.Component.Name)
	fmt.Printf("Type:         %s\n", response.Component.Type)
	fmt.Printf("Version:      %s\n", response.Component.Version)
	fmt.Printf("Topic ID:     %s\n", response.Component.TopicID)
	fmt.Printf("State:        %s\n", response.Component.State)
	if response.Component.DisplayName != "" {
		fmt.Printf("Display Name: %s\n", response.Component.DisplayName)
	}
	if len(response.Component.Tags) > 0 {
		fmt.Printf("Tags:         %s\n", strings.Join(response.Component.Tags, ", "))
	}
	fmt.Printf("Created:      %s\n", response.Component.CreatedAt)
	fmt.Printf("Updated:      %s\n", response.Component.UpdatedAt)
}

func printComponentJSON(response *lib.ComponentResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(getComponentCmd)
	rootCmd.AddCommand(createComponentCmd)
	rootCmd.AddCommand(updateComponentCmd)
	rootCmd.AddCommand(deleteComponentCmd)

	// get component flags
	getComponentCmd.PersistentFlags().StringVar(&getComponentIDFlag, "id", "", "Component ID (required)")
	getComponentCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create component flags
	createComponentCmd.PersistentFlags().StringVar(&createComponentName, "name", "", "Component name (required)")
	createComponentCmd.PersistentFlags().StringVar(&createComponentType, "type", "", "Component type (required)")
	createComponentCmd.PersistentFlags().StringVar(&createComponentTopicID, "topic-id", "", "Topic ID (required)")
	createComponentCmd.PersistentFlags().StringVar(&createComponentVersion, "version", "", "Component version")
	createComponentCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update component flags
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentIDFlag, "id", "", "Component ID to update (required)")
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentName, "name", "", "New component name")
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentState, "state", "", "New component state")
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentVersion, "version", "", "New component version")
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentTags, "tags", "", "Comma-separated list of tags")
	updateComponentCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete component flags
	deleteComponentCmd.PersistentFlags().StringVar(&deleteComponentIDFlag, "id", "", "Component ID to delete (required)")
	deleteComponentCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
