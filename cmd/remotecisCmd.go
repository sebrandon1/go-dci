package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for remoteci command flags
var (
	getRemoteCIsCmd_IDFlag   string
	createRemoteCINameFlag   string
	createRemoteCITeamIDFlag string
	updateRemoteCIIDFlag     string
	updateRemoteCINameFlag   string
	updateRemoteCIStateFlag  string
	deleteRemoteCIIDFlag     string
)

var getRemoteCIsCmd = &cobra.Command{
	Use:   "remotecis",
	Short: "Get all remote CIs from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting remote CIs...")

		response, err := client.GetRemoteCIs(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get remote CIs: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printRemoteCIsJSON(response)
		}

		printRemoteCIsStdout(response)

		return nil
	},
}

var getRemoteCICmd = &cobra.Command{
	Use:   "remoteci",
	Short: "Get a specific remote CI by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting remote CI with ID: %s\n", getRemoteCIsCmd_IDFlag)

		response, err := client.GetRemoteCI(cmd.Context(), getRemoteCIsCmd_IDFlag)
		if err != nil {
			return fmt.Errorf("failed to get remote CI: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printRemoteCIJSON(response)
		}

		printRemoteCIStdout(response)

		return nil
	},
}

var createRemoteCICmd = &cobra.Command{
	Use:   "create-remoteci",
	Short: "Create a new remote CI in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		if dryRunFlag {
			printStatus("[DRY RUN] Would create remote CI: name=%s, team-id=%s\n", createRemoteCINameFlag, createRemoteCITeamIDFlag)
			return nil
		}

		printStatus("Creating remote CI: %s\n", createRemoteCINameFlag)

		response, err := client.CreateRemoteCI(cmd.Context(), createRemoteCINameFlag, createRemoteCITeamIDFlag)
		if err != nil {
			return fmt.Errorf("failed to create remote CI: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printRemoteCIJSON(response)
		}

		printStatus("Remote CI created successfully!")
		printRemoteCIStdout(response)

		return nil
	},
}

var updateRemoteCICmd = &cobra.Command{
	Use:   "update-remoteci",
	Short: "Update an existing remote CI in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateRemoteCIRequest{}
		if updateRemoteCINameFlag != "" {
			updates.Name = updateRemoteCINameFlag
		}
		if updateRemoteCIStateFlag != "" {
			updates.State = lib.ResourceState(updateRemoteCIStateFlag)
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would update remote CI: id=%s, name=%s, state=%s\n", updateRemoteCIIDFlag, updateRemoteCINameFlag, updateRemoteCIStateFlag)
			return nil
		}

		printStatus("Updating remote CI: %s\n", updateRemoteCIIDFlag)

		response, err := client.UpdateRemoteCI(cmd.Context(), updateRemoteCIIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update remote CI: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printRemoteCIJSON(response)
		}

		printStatus("Remote CI updated successfully!")
		printRemoteCIStdout(response)

		return nil
	},
}

var deleteRemoteCICmd = &cobra.Command{
	Use:   "delete-remoteci",
	Short: "Delete a remote CI from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would delete remote CI: id=%s\n", deleteRemoteCIIDFlag)
			return nil
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("remote CI", deleteRemoteCIIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			printStatus("Deletion canceled")
			return nil
		}


		client := lib.NewClient(accessKey, secretKey)

		printStatus("Deleting remote CI: %s\n", deleteRemoteCIIDFlag)

		err = client.DeleteRemoteCI(cmd.Context(), deleteRemoteCIIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete remote CI: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteRemoteCIIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			printStatus("Remote CI deleted successfully!")
		}

		return nil
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

func printRemoteCIsJSON(response *lib.RemoteCIsResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
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

func printRemoteCIJSON(response *lib.RemoteCIResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
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
	getRemoteCICmd.PersistentFlags().StringVar(&getRemoteCIsCmd_IDFlag, "id", "", "Remote CI ID")
	_ = getRemoteCICmd.MarkPersistentFlagRequired("id")
	getRemoteCICmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create remote CI flags
	createRemoteCICmd.PersistentFlags().StringVar(&createRemoteCINameFlag, "name", "", "Remote CI name")
	_ = createRemoteCICmd.MarkPersistentFlagRequired("name")
	createRemoteCICmd.PersistentFlags().StringVar(&createRemoteCITeamIDFlag, "team-id", "", "Team ID")
	_ = createRemoteCICmd.MarkPersistentFlagRequired("team-id")
	createRemoteCICmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update remote CI flags
	updateRemoteCICmd.PersistentFlags().StringVar(&updateRemoteCIIDFlag, "id", "", "Remote CI ID to update")
	_ = updateRemoteCICmd.MarkPersistentFlagRequired("id")
	updateRemoteCICmd.PersistentFlags().StringVar(&updateRemoteCINameFlag, "name", "", "New remote CI name")
	updateRemoteCICmd.PersistentFlags().StringVar(&updateRemoteCIStateFlag, "state", "", "New remote CI state")
	updateRemoteCICmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete remote CI flags
	deleteRemoteCICmd.PersistentFlags().StringVar(&deleteRemoteCIIDFlag, "id", "", "Remote CI ID to delete")
	_ = deleteRemoteCICmd.MarkPersistentFlagRequired("id")
	deleteRemoteCICmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
