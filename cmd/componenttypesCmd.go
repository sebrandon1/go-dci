package cmd

import (
	"encoding/json"
	"fmt"
	"log"

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
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if getComponentTypeIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting component type with ID: %s\n", getComponentTypeIDFlag)
		}

		response, err := client.GetComponentType(getComponentTypeIDFlag)
		if err != nil {
			fmt.Printf("Failed to get component type: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printComponentTypeJSON(response)
		} else {
			printComponentTypeStdout(response)
		}
	},
}

var createComponentTypeCmd = &cobra.Command{
	Use:   "create-componenttype",
	Short: "Create a new component type in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if createComponentTypeName == "" {
			fmt.Println("Error: --name is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Creating component type: %s\n", createComponentTypeName)
		}

		response, err := client.CreateComponentType(createComponentTypeName)
		if err != nil {
			fmt.Printf("Failed to create component type: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printComponentTypeJSON(response)
		} else {
			fmt.Println("Component type created successfully!")
			printComponentTypeStdout(response)
		}
	},
}

var updateComponentTypeCmd = &cobra.Command{
	Use:   "update-componenttype",
	Short: "Update an existing component type in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if updateComponentTypeIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateComponentTypeRequest{}
		if updateComponentTypeName != "" {
			updates.Name = updateComponentTypeName
		}
		if updateComponentTypeState != "" {
			updates.State = updateComponentTypeState
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Updating component type: %s\n", updateComponentTypeIDFlag)
		}

		response, err := client.UpdateComponentType(updateComponentTypeIDFlag, updates)
		if err != nil {
			fmt.Printf("Failed to update component type: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printComponentTypeJSON(response)
		} else {
			fmt.Println("Component type updated successfully!")
			printComponentTypeStdout(response)
		}
	},
}

var deleteComponentTypeCmd = &cobra.Command{
	Use:   "delete-componenttype",
	Short: "Delete a component type from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if deleteComponentTypeIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Deleting component type: %s\n", deleteComponentTypeIDFlag)
		}

		err = client.DeleteComponentType(deleteComponentTypeIDFlag)
		if err != nil {
			fmt.Printf("Failed to delete component type: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteComponentTypeIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("Component type deleted successfully!")
		}
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

func printComponentTypeJSON(response *lib.ComponentTypeResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(getComponentTypeCmd)
	rootCmd.AddCommand(createComponentTypeCmd)
	rootCmd.AddCommand(updateComponentTypeCmd)
	rootCmd.AddCommand(deleteComponentTypeCmd)

	// get component type flags
	getComponentTypeCmd.PersistentFlags().StringVar(&getComponentTypeIDFlag, "id", "", "Component Type ID (required)")
	getComponentTypeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create component type flags
	createComponentTypeCmd.PersistentFlags().StringVar(&createComponentTypeName, "name", "", "Component type name (required)")
	createComponentTypeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update component type flags
	updateComponentTypeCmd.PersistentFlags().StringVar(&updateComponentTypeIDFlag, "id", "", "Component Type ID to update (required)")
	updateComponentTypeCmd.PersistentFlags().StringVar(&updateComponentTypeName, "name", "", "New component type name")
	updateComponentTypeCmd.PersistentFlags().StringVar(&updateComponentTypeState, "state", "", "New component type state")
	updateComponentTypeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete component type flags
	deleteComponentTypeCmd.PersistentFlags().StringVar(&deleteComponentTypeIDFlag, "id", "", "Component Type ID to delete (required)")
	deleteComponentTypeCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
