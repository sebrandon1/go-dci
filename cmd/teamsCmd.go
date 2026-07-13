package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for team command flags
var (
	getTeamIDFlag       string
	createTeamNameFlag  string
	updateTeamIDFlag    string
	updateTeamNameFlag  string
	updateTeamStateFlag string
	deleteTeamIDFlag    string
	teamsNameFilter     string
)

var getTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Get all teams from DCI, optionally filtered by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		if teamsNameFilter != "" {
			printStatus("Getting teams matching name: %s\n", teamsNameFilter)
		} else {
			printStatus("Getting teams...")
		}

		response, err := client.GetTeamsFiltered(cmd.Context(), teamsNameFilter)
		if err != nil {
			return fmt.Errorf("failed to get teams: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printTeamsJSON(response)
		}

		printTeamsStdout(response)

		return nil
	},
}

var getTeamCmd = &cobra.Command{
	Use:   "team",
	Short: "Get a specific team by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(getTeamIDFlag, "team"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting team with ID: %s\n", getTeamIDFlag)

		response, err := client.GetTeam(cmd.Context(), getTeamIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get team: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printTeamJSON(response)
		}

		printTeamStdout(response)

		return nil
	},
}

var createTeamCmd = &cobra.Command{
	Use:   "create-team",
	Short: "Create a new team in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		if dryRunFlag {
			printStatus("[DRY RUN] Would create team: name=%s\n", createTeamNameFlag)
			return nil
		}

		printStatus("Creating team: %s\n", createTeamNameFlag)

		response, err := client.CreateTeam(cmd.Context(), createTeamNameFlag)
		if err != nil {
			return fmt.Errorf("failed to create team: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printTeamJSON(response)
		}

		printStatus("Team created successfully!")
		printTeamStdout(response)

		return nil
	},
}

var updateTeamCmd = &cobra.Command{
	Use:   "update-team",
	Short: "Update an existing team in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(updateTeamIDFlag, "team"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateTeamRequest{}
		if updateTeamNameFlag != "" {
			updates.Name = updateTeamNameFlag
		}
		if updateTeamStateFlag != "" {
			updates.State = lib.ResourceState(updateTeamStateFlag)
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would update team: id=%s, name=%s, state=%s\n", updateTeamIDFlag, updateTeamNameFlag, updateTeamStateFlag)
			return nil
		}

		printStatus("Updating team: %s\n", updateTeamIDFlag)

		response, err := client.UpdateTeam(cmd.Context(), updateTeamIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update team: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printTeamJSON(response)
		}

		printStatus("Team updated successfully!")
		printTeamStdout(response)

		return nil
	},
}

var deleteTeamCmd = &cobra.Command{
	Use:   "delete-team",
	Short: "Delete a team from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(deleteTeamIDFlag, "team"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would delete team: id=%s\n", deleteTeamIDFlag)
			return nil
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("team", deleteTeamIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			printStatus("Deletion canceled")
			return nil
		}


		client := lib.NewClient(accessKey, secretKey)

		printStatus("Deleting team: %s\n", deleteTeamIDFlag)

		err = client.DeleteTeam(cmd.Context(), deleteTeamIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete team: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteTeamIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			printStatus("Team deleted successfully!")
		}

		return nil
	},
}

func printTeamsStdout(responses []lib.TeamsResponse) {
	totalTeams := 0
	for _, resp := range responses {
		totalTeams += len(resp.Teams)
	}

	if totalTeams == 0 {
		fmt.Println("No teams found.")
		return
	}

	fmt.Println("---")
	for _, resp := range responses {
		for _, team := range resp.Teams {
			fmt.Printf("ID: %s | Name: %s | State: %s | External: %v\n",
				team.ID, team.Name, team.State, team.External)
		}
	}
	fmt.Printf("Total Teams: %d\n", totalTeams)
}

func printTeamsJSON(responses []lib.TeamsResponse) error {
	// Flatten all teams from paginated responses
	var allTeams []lib.Team
	for _, resp := range responses {
		allTeams = append(allTeams, resp.Teams...)
	}

	jsonBytes, err := json.Marshal(map[string]interface{}{
		"teams": allTeams,
		"total": len(allTeams),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func printTeamStdout(response *lib.TeamResponse) {
	fmt.Println("---")
	fmt.Printf("ID:        %s\n", response.Team.ID)
	fmt.Printf("Name:      %s\n", response.Team.Name)
	fmt.Printf("Country:   %s\n", response.Team.Country)
	fmt.Printf("External:  %v\n", response.Team.External)
	fmt.Printf("State:     %s\n", response.Team.State)
	fmt.Printf("Created:   %s\n", response.Team.CreatedAt)
	fmt.Printf("Updated:   %s\n", response.Team.UpdatedAt)
}

func printTeamJSON(response *lib.TeamResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(getTeamsCmd)
	rootCmd.AddCommand(getTeamCmd)
	rootCmd.AddCommand(createTeamCmd)
	rootCmd.AddCommand(updateTeamCmd)
	rootCmd.AddCommand(deleteTeamCmd)

	// get teams flags
	getTeamsCmd.PersistentFlags().StringVarP(&teamsNameFilter, "name", "n", "", "Filter teams by name")
	getTeamsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// get team flags
	getTeamCmd.PersistentFlags().StringVar(&getTeamIDFlag, "id", "", "Team ID")
	_ = getTeamCmd.MarkPersistentFlagRequired("id")
	getTeamCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create team flags
	createTeamCmd.PersistentFlags().StringVar(&createTeamNameFlag, "name", "", "Team name")
	_ = createTeamCmd.MarkPersistentFlagRequired("name")
	createTeamCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update team flags
	updateTeamCmd.PersistentFlags().StringVar(&updateTeamIDFlag, "id", "", "Team ID to update")
	_ = updateTeamCmd.MarkPersistentFlagRequired("id")
	updateTeamCmd.PersistentFlags().StringVar(&updateTeamNameFlag, "name", "", "New team name")
	updateTeamCmd.PersistentFlags().StringVar(&updateTeamStateFlag, "state", "", "New team state")
	updateTeamCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete team flags
	deleteTeamCmd.PersistentFlags().StringVar(&deleteTeamIDFlag, "id", "", "Team ID to delete")
	_ = deleteTeamCmd.MarkPersistentFlagRequired("id")
	deleteTeamCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
