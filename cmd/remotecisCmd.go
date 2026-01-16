package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for remoteci command flags
var (
	getRemoteCIsCmd_IDFlag       string
	createRemoteCINameFlag       string
	createRemoteCITeamIDFlag     string
	updateRemoteCIIDFlag         string
	updateRemoteCINameFlag       string
	updateRemoteCIStateFlag      string
	deleteRemoteCIIDFlag         string
)

var getRemoteCIsCmd = &cobra.Command{
	Use:   "remotecis",
	Short: "Get all remote CIs from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Println("Getting remote CIs...")
		}

		response, err := client.GetRemoteCIs()
		if err != nil {
			fmt.Printf("Failed to get remote CIs: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printRemoteCIsJSON(response)
		} else {
			printRemoteCIsStdout(response)
		}
	},
}

var getRemoteCICmd = &cobra.Command{
	Use:   "remoteci",
	Short: "Get a specific remote CI by ID",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if getRemoteCIsCmd_IDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting remote CI with ID: %s\n", getRemoteCIsCmd_IDFlag)
		}

		response, err := client.GetRemoteCI(getRemoteCIsCmd_IDFlag)
		if err != nil {
			fmt.Printf("Failed to get remote CI: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printRemoteCIJSON(response)
		} else {
			printRemoteCIStdout(response)
		}
	},
}

var createRemoteCICmd = &cobra.Command{
	Use:   "create-remoteci",
	Short: "Create a new remote CI in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if createRemoteCINameFlag == "" {
			fmt.Println("Error: --name is required")
			return
		}
		if createRemoteCITeamIDFlag == "" {
			fmt.Println("Error: --team-id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Creating remote CI: %s\n", createRemoteCINameFlag)
		}

		response, err := client.CreateRemoteCI(createRemoteCINameFlag, createRemoteCITeamIDFlag)
		if err != nil {
			fmt.Printf("Failed to create remote CI: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printRemoteCIJSON(response)
		} else {
			fmt.Println("Remote CI created successfully!")
			printRemoteCIStdout(response)
		}
	},
}

var updateRemoteCICmd = &cobra.Command{
	Use:   "update-remoteci",
	Short: "Update an existing remote CI in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if updateRemoteCIIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateRemoteCIRequest{}
		if updateRemoteCINameFlag != "" {
			updates.Name = updateRemoteCINameFlag
		}
		if updateRemoteCIStateFlag != "" {
			updates.State = updateRemoteCIStateFlag
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Updating remote CI: %s\n", updateRemoteCIIDFlag)
		}

		response, err := client.UpdateRemoteCI(updateRemoteCIIDFlag, updates)
		if err != nil {
			fmt.Printf("Failed to update remote CI: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printRemoteCIJSON(response)
		} else {
			fmt.Println("Remote CI updated successfully!")
			printRemoteCIStdout(response)
		}
	},
}

var deleteRemoteCICmd = &cobra.Command{
	Use:   "delete-remoteci",
	Short: "Delete a remote CI from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if deleteRemoteCIIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Deleting remote CI: %s\n", deleteRemoteCIIDFlag)
		}

		err = client.DeleteRemoteCI(deleteRemoteCIIDFlag)
		if err != nil {
			fmt.Printf("Failed to delete remote CI: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteRemoteCIIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("Remote CI deleted successfully!")
		}
	},
}

func printRemoteCIsStdout(response *lib.RemoteCIsResponse) {
	if len(response.RemoteCIs) == 0 {
		fmt.Println("No remote CIs found.")
		return
	}
	fmt.Println("---")
	for _, rci := range response.RemoteCIs {
		fmt.Printf("ID: %s | Name: %s | Team ID: %s | State: %s\n",
			rci.ID, rci.Name, rci.TeamID, rci.State)
	}
	fmt.Printf("Total Remote CIs: %d\n", len(response.RemoteCIs))
}

func printRemoteCIsJSON(response *lib.RemoteCIsResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func printRemoteCIStdout(response *lib.RemoteCIResponse) {
	fmt.Println("---")
	fmt.Printf("ID:       %s\n", response.RemoteCI.ID)
	fmt.Printf("Name:     %s\n", response.RemoteCI.Name)
	fmt.Printf("Team ID:  %s\n", response.RemoteCI.TeamID)
	fmt.Printf("State:    %s\n", response.RemoteCI.State)
	fmt.Printf("Created:  %s\n", response.RemoteCI.CreatedAt)
	fmt.Printf("Updated:  %s\n", response.RemoteCI.UpdatedAt)
}

func printRemoteCIJSON(response *lib.RemoteCIResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}

func init() {
	rootCmd.AddCommand(getRemoteCIsCmd)
	rootCmd.AddCommand(getRemoteCICmd)
	rootCmd.AddCommand(createRemoteCICmd)
	rootCmd.AddCommand(updateRemoteCICmd)
	rootCmd.AddCommand(deleteRemoteCICmd)

	// get remote CIs flags
	getRemoteCIsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// get remote CI flags
	getRemoteCICmd.PersistentFlags().StringVar(&getRemoteCIsCmd_IDFlag, "id", "", "Remote CI ID (required)")
	getRemoteCICmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create remote CI flags
	createRemoteCICmd.PersistentFlags().StringVar(&createRemoteCINameFlag, "name", "", "Remote CI name (required)")
	createRemoteCICmd.PersistentFlags().StringVar(&createRemoteCITeamIDFlag, "team-id", "", "Team ID (required)")
	createRemoteCICmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update remote CI flags
	updateRemoteCICmd.PersistentFlags().StringVar(&updateRemoteCIIDFlag, "id", "", "Remote CI ID to update (required)")
	updateRemoteCICmd.PersistentFlags().StringVar(&updateRemoteCINameFlag, "name", "", "New remote CI name")
	updateRemoteCICmd.PersistentFlags().StringVar(&updateRemoteCIStateFlag, "state", "", "New remote CI state")
	updateRemoteCICmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete remote CI flags
	deleteRemoteCICmd.PersistentFlags().StringVar(&deleteRemoteCIIDFlag, "id", "", "Remote CI ID to delete (required)")
	deleteRemoteCICmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
