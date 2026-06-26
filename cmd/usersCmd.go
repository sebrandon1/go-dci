package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
)

// Variables for user command flags
var (
	getUserIDFlag         string
	createUserNameFlag    string
	createUserEmailFlag   string
	createUserFullnameFlag string
	createUserTeamIDFlag  string
	createUserPasswordFlag string
	updateUserIDFlag      string
	updateUserNameFlag    string
	updateUserEmailFlag   string
	updateUserFullnameFlag string
	updateUserStateFlag   string
	deleteUserIDFlag      string
)

var getUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get all users from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Println("Getting users...")
		}

		response, err := client.GetUsers(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get users: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printUsersJSON(response)
		}

		printUsersStdout(response)

		return nil
	},
}

var getUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get a specific user by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if getUserIDFlag == "" {
			return fmt.Errorf("--id is required")
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting user with ID: %s\n", getUserIDFlag)
		}

		response, err := client.GetUser(cmd.Context(), getUserIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get user: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printUserJSON(response)
		}

		printUserStdout(response)

		return nil
	},
}

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if createUserNameFlag == "" {
			return fmt.Errorf("--name is required")
		}
		if createUserEmailFlag == "" {
			return fmt.Errorf("--email is required")
		}
		if createUserTeamIDFlag == "" {
			return fmt.Errorf("--team-id is required")
		}
		if createUserPasswordFlag == "" {
			return fmt.Errorf("--password is required")
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Creating user: %s\n", createUserNameFlag)
		}

		response, err := client.CreateUser(
			cmd.Context(),
			createUserNameFlag,
			createUserEmailFlag,
			createUserFullnameFlag,
			createUserTeamIDFlag,
			createUserPasswordFlag,
		)
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printUserJSON(response)
		}

		fmt.Println("User created successfully!")
		printUserStdout(response)

		return nil
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update-user",
	Short: "Update an existing user in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if updateUserIDFlag == "" {
			return fmt.Errorf("--id is required")
		}

		client := lib.NewClient(accessKey, secretKey)

		updates := lib.UpdateUserRequest{}
		if updateUserNameFlag != "" {
			updates.Name = updateUserNameFlag
		}
		if updateUserEmailFlag != "" {
			updates.Email = updateUserEmailFlag
		}
		if updateUserFullnameFlag != "" {
			updates.Fullname = updateUserFullnameFlag
		}
		if updateUserStateFlag != "" {
			updates.State = lib.ResourceState(updateUserStateFlag)
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Updating user: %s\n", updateUserIDFlag)
		}

		response, err := client.UpdateUser(cmd.Context(), updateUserIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			return printUserJSON(response)
		}

		fmt.Println("User updated successfully!")
		printUserStdout(response)

		return nil
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Delete a user from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if deleteUserIDFlag == "" {
			return fmt.Errorf("--id is required")
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("user", deleteUserIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			if outputFormat != OutputFormatJSON {
				fmt.Println("Deletion canceled")
			}
			return nil
		}


		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Deleting user: %s\n", deleteUserIDFlag)
		}

		err = client.DeleteUser(cmd.Context(), deleteUserIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete user: %v", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteUserIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("User deleted successfully!")
		}

		return nil
	},
}

func printUsersStdout(response *lib.UsersResponse) {
	if len(response.Users) == 0 {
		fmt.Println("No users found.")
		return
	}
	fmt.Println("---")
	for _, user := range response.Users {
		fmt.Printf("ID: %s | Name: %s | Email: %s | State: %s\n",
			user.ID, user.Name, user.Email, user.State)
	}
	fmt.Printf("Total Users: %d\n", len(response.Users))
}

func printUsersJSON(response *lib.UsersResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func printUserStdout(response *lib.UserResponse) {
	fmt.Println("---")
	fmt.Printf("ID:       %s\n", response.User.ID)
	fmt.Printf("Name:     %s\n", response.User.Name)
	fmt.Printf("Fullname: %s\n", response.User.Fullname)
	fmt.Printf("Email:    %s\n", response.User.Email)
	fmt.Printf("Team ID:  %s\n", response.User.TeamID)
	fmt.Printf("State:    %s\n", response.User.State)
	fmt.Printf("Created:  %s\n", response.User.CreatedAt)
	fmt.Printf("Updated:  %s\n", response.User.UpdatedAt)
}

func printUserJSON(response *lib.UserResponse) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

func init() {
	rootCmd.AddCommand(getUsersCmd)
	rootCmd.AddCommand(getUserCmd)
	rootCmd.AddCommand(createUserCmd)
	rootCmd.AddCommand(updateUserCmd)
	rootCmd.AddCommand(deleteUserCmd)

	// get users flags
	getUsersCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// get user flags
	getUserCmd.PersistentFlags().StringVar(&getUserIDFlag, "id", "", "User ID (required)")
	getUserCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// create user flags
	createUserCmd.PersistentFlags().StringVar(&createUserNameFlag, "name", "", "Username (required)")
	createUserCmd.PersistentFlags().StringVar(&createUserEmailFlag, "email", "", "User email (required)")
	createUserCmd.PersistentFlags().StringVar(&createUserFullnameFlag, "fullname", "", "User full name")
	createUserCmd.PersistentFlags().StringVar(&createUserTeamIDFlag, "team-id", "", "Team ID (required)")
	createUserCmd.PersistentFlags().StringVar(&createUserPasswordFlag, "password", "", "User password (required)")
	createUserCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// update user flags
	updateUserCmd.PersistentFlags().StringVar(&updateUserIDFlag, "id", "", "User ID to update (required)")
	updateUserCmd.PersistentFlags().StringVar(&updateUserNameFlag, "name", "", "New username")
	updateUserCmd.PersistentFlags().StringVar(&updateUserEmailFlag, "email", "", "New email")
	updateUserCmd.PersistentFlags().StringVar(&updateUserFullnameFlag, "fullname", "", "New full name")
	updateUserCmd.PersistentFlags().StringVar(&updateUserStateFlag, "state", "", "New user state")
	updateUserCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")

	// delete user flags
	deleteUserCmd.PersistentFlags().StringVar(&deleteUserIDFlag, "id", "", "User ID to delete (required)")
	deleteUserCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", OutputFormatStdout, "Output format (json) - default is stdout")
}
