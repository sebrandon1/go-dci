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
	getComponentIDFlag     string
	createComponentName    string
	createComponentType    string
	createComponentTopicID string
	createComponentVersion string
	updateComponentIDFlag  string
	updateComponentName    string
	updateComponentState   string
	updateComponentVersion string
	updateComponentTags    string
	deleteComponentIDFlag  string
)

var getComponentCmd = &cobra.Command{
	Use:   "component",
	Short: "Get a specific component by ID",
	Long:  `Retrieve detailed information about a single DCI component by its UUID.`,
	Example: `  dci component --id <component-uuid>
  dci component --id <component-uuid> -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(getComponentIDFlag, "component"); err != nil {
			return err
		}


		printStatus("Getting component with ID: %s\n", getComponentIDFlag)

		response, err := dciClient.GetComponent(cmd.Context(), getComponentIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get component: %w", err)
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
	Long: `Create a new versioned component in DCI associated with a topic.
The component type must exist in DCI before creation.`,
	Example: `  dci create-component --name "OCP 4.17.3" --type ocp --topic-id <topic-uuid> --version 4.17.3
  dci create-component --name "OCP 4.17.3" --type ocp --topic-id <topic-uuid> --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if dryRunFlag {
			printStatus("[DRY RUN] Would create component: name=%s, type=%s, topic-id=%s, version=%s\n", createComponentName, createComponentType, createComponentTopicID, createComponentVersion)
			return nil
		}

		printStatus("Creating component: %s\n", createComponentName)

		response, err := dciClient.CreateComponent(cmd.Context(), createComponentName, createComponentType, createComponentTopicID, createComponentVersion)
		if err != nil {
			return fmt.Errorf("failed to create component: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentJSON(response)
		}

		printStatus("Component created successfully!")
		printComponentStdout(response)

		return nil
	},
}

var updateComponentCmd = &cobra.Command{
	Use:   "update-component",
	Short: "Update an existing component in DCI",
	Long:  `Update mutable fields (name, state, version, tags) of a DCI component by UUID.`,
	Example: `  dci update-component --id <component-uuid> --state active
  dci update-component --id <component-uuid> --version 4.17.4 --tags "ga,stable"
  dci update-component --id <component-uuid> --name "OCP 4.17.4" --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(updateComponentIDFlag, "component"); err != nil {
			return err
		}


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

		if dryRunFlag {
			printStatus("[DRY RUN] Would update component: id=%s, name=%s, state=%s, version=%s, tags=%v\n", updateComponentIDFlag, updateComponentName, updateComponentState, updateComponentVersion, updates.Tags)
			return nil
		}

		printStatus("Updating component: %s\n", updateComponentIDFlag)

		response, err := dciClient.UpdateComponent(cmd.Context(), updateComponentIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update component: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentJSON(response)
		}

		printStatus("Component updated successfully!")
		printComponentStdout(response)

		return nil
	},
}

var deleteComponentCmd = &cobra.Command{
	Use:   "delete-component",
	Short: "Delete a component from DCI",
	Long: `Permanently delete a DCI component by its UUID. Prompts for confirmation
unless --yes is set or output is JSON (automation mode).`,
	Example: `  dci delete-component --id <component-uuid>
  dci delete-component --id <component-uuid> --yes
  dci delete-component --id <component-uuid> --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(deleteComponentIDFlag, "component"); err != nil {
			return err
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would delete component: id=%s\n", deleteComponentIDFlag)
			return nil
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("component", deleteComponentIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			printStatus("Deletion canceled")
			return nil
		}

		printStatus("Deleting component: %s\n", deleteComponentIDFlag)

		err = dciClient.DeleteComponent(cmd.Context(), deleteComponentIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete component: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteComponentIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			printStatus("Component deleted successfully!")
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
		return fmt.Errorf("failed to marshal JSON: %w", err)
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
	getComponentCmd.PersistentFlags().StringVar(&getComponentIDFlag, "id", "", "Component ID")
	_ = getComponentCmd.MarkPersistentFlagRequired("id")

	// create component flags
	createComponentCmd.PersistentFlags().StringVar(&createComponentName, "name", "", "Component name")
	_ = createComponentCmd.MarkPersistentFlagRequired("name")
	createComponentCmd.PersistentFlags().StringVar(&createComponentType, "type", "", "Component type")
	_ = createComponentCmd.MarkPersistentFlagRequired("type")
	createComponentCmd.PersistentFlags().StringVar(&createComponentTopicID, "topic-id", "", "Topic ID")
	_ = createComponentCmd.MarkPersistentFlagRequired("topic-id")
	createComponentCmd.PersistentFlags().StringVar(&createComponentVersion, "version", "", "Component version")

	// update component flags
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentIDFlag, "id", "", "Component ID to update")
	_ = updateComponentCmd.MarkPersistentFlagRequired("id")
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentName, "name", "", "New component name")
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentState, "state", "", "New component state")
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentVersion, "version", "", "New component version")
	updateComponentCmd.PersistentFlags().StringVar(&updateComponentTags, "tags", "", "Comma-separated list of tags")

	// delete component flags
	deleteComponentCmd.PersistentFlags().StringVar(&deleteComponentIDFlag, "id", "", "Component ID to delete")
	_ = deleteComponentCmd.MarkPersistentFlagRequired("id")
}
