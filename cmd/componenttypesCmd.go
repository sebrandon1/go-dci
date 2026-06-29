package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for component type command flags
var (
	getComponentTypeIDFlag    string
	createComponentTypeName   string
	updateComponentTypeIDFlag string
	updateComponentTypeName   string
	updateComponentTypeState  string
	deleteComponentTypeIDFlag string
)

var getComponentTypeCmd = &cobra.Command{
	Use:   "componenttype",
	Short: "Get a specific component type by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting component type with ID: %s\n", getComponentTypeIDFlag)

		response, err := client.GetComponentType(cmd.Context(), getComponentTypeIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get component type: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentTypeJSON(response)
		}

		printComponentTypeStdout(response)

		return nil
	},
}

var createComponentTypeCmd = &cobra.Command{
	Use:   "create-componenttype",
	Short: "Create a new component type in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		if dryRunFlag {
			printStatus("[DRY RUN] Would create component type: name=%s\n", createComponentTypeName)
			return nil
		}

		printStatus("Creating component type: %s\n", createComponentTypeName)

		response, err := client.CreateComponentType(cmd.Context(), createComponentTypeName)
		if err != nil {
			return fmt.Errorf("failed to create component type: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentTypeJSON(response)
		}

		printStatus("Component type created successfully!")
		printComponentTypeStdout(response)

		return nil
	},
}

var updateComponentTypeCmd = &cobra.Command{
	Use:   "update-componenttype",
	Short: "Update an existing component type in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateComponentTypeRequest{}
		if updateComponentTypeName != "" {
			updates.Name = updateComponentTypeName
		}
		if updateComponentTypeState != "" {
			updates.State = lib.ResourceState(updateComponentTypeState)
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would update component type: id=%s, name=%s, state=%s\n", updateComponentTypeIDFlag, updateComponentTypeName, updateComponentTypeState)
			return nil
		}

		printStatus("Updating component type: %s\n", updateComponentTypeIDFlag)

		response, err := client.UpdateComponentType(cmd.Context(), updateComponentTypeIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update component type: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printComponentTypeJSON(response)
		}

		printStatus("Component type updated successfully!")
		printComponentTypeStdout(response)

		return nil
	},
}

var deleteComponentTypeCmd = &cobra.Command{
	Use:   "delete-componenttype",
	Short: "Delete a component type from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would delete component type: id=%s\n", deleteComponentTypeIDFlag)
			return nil
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("component type", deleteComponentTypeIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			printStatus("Deletion canceled")
			return nil
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Deleting component type: %s\n", deleteComponentTypeIDFlag)

		err = client.DeleteComponentType(cmd.Context(), deleteComponentTypeIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete component type: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteComponentTypeIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			printStatus("Component type deleted successfully!")
		}

		return nil
	},
}

func printComponentTypeStdout(response *lib.ComponentTypeResponse) {
	fmt.Println("---")
	fmt.Printf("ID:       %s\n", response.ComponentType.ID)
	fmt.Printf("Name:     %s\n", response.ComponentType.Name)
	fmt.Printf("State:    %s\n", response.ComponentType.State)
	fmt.Printf("Created:  %s\n", response.ComponentType.CreatedAt)
	fmt.Printf("Updated:  %s\n", response.ComponentType.UpdatedAt)
}

func printComponentTypeJSON(response *lib.ComponentTypeResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(getComponentTypeCmd)
	rootCmd.AddCommand(createComponentTypeCmd)
	rootCmd.AddCommand(updateComponentTypeCmd)
	rootCmd.AddCommand(deleteComponentTypeCmd)

	// get component type flags
	getComponentTypeCmd.PersistentFlags().StringVar(&getComponentTypeIDFlag, "id", "", "Component Type ID")
	_ = getComponentTypeCmd.MarkPersistentFlagRequired("id")
	getComponentTypeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create component type flags
	createComponentTypeCmd.PersistentFlags().StringVar(&createComponentTypeName, "name", "", "Component type name")
	_ = createComponentTypeCmd.MarkPersistentFlagRequired("name")
	createComponentTypeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update component type flags
	updateComponentTypeCmd.PersistentFlags().StringVar(&updateComponentTypeIDFlag, "id", "", "Component Type ID to update")
	_ = updateComponentTypeCmd.MarkPersistentFlagRequired("id")
	updateComponentTypeCmd.PersistentFlags().StringVar(&updateComponentTypeName, "name", "", "New component type name")
	updateComponentTypeCmd.PersistentFlags().StringVar(&updateComponentTypeState, "state", "", "New component type state")
	updateComponentTypeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete component type flags
	deleteComponentTypeCmd.PersistentFlags().StringVar(&deleteComponentTypeIDFlag, "id", "", "Component Type ID to delete")
	_ = deleteComponentTypeCmd.MarkPersistentFlagRequired("id")
	deleteComponentTypeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
