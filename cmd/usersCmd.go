package cmd

import (
	"encoding/json"
	"fmt"
	"log"

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
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Println("Getting users...")
		}

		response, err := client.GetUsers()
		if err != nil {
			fmt.Printf("Failed to get users: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printUsersJSON(response)
		} else {
			printUsersStdout(response)
		}
	},
}

var getUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get a specific user by ID",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if getUserIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Getting user with ID: %s\n", getUserIDFlag)
		}

		response, err := client.GetUser(getUserIDFlag)
		if err != nil {
			fmt.Printf("Failed to get user: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printUserJSON(response)
		} else {
			printUserStdout(response)
		}
	},
}

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if createUserNameFlag == "" {
			fmt.Println("Error: --name is required")
			return
		}
		if createUserEmailFlag == "" {
			fmt.Println("Error: --email is required")
			return
		}
		if createUserTeamIDFlag == "" {
			fmt.Println("Error: --team-id is required")
			return
		}
		if createUserPasswordFlag == "" {
			fmt.Println("Error: --password is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Creating user: %s\n", createUserNameFlag)
		}

		response, err := client.CreateUser(
			createUserNameFlag,
			createUserEmailFlag,
			createUserFullnameFlag,
			createUserTeamIDFlag,
			createUserPasswordFlag,
		)
		if err != nil {
			fmt.Printf("Failed to create user: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printUserJSON(response)
		} else {
			fmt.Println("User created successfully!")
			printUserStdout(response)
		}
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update-user",
	Short: "Update an existing user in DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if updateUserIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
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
			updates.State = updateUserStateFlag
		}

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Updating user: %s\n", updateUserIDFlag)
		}

		response, err := client.UpdateUser(updateUserIDFlag, updates)
		if err != nil {
			fmt.Printf("Failed to update user: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			printUserJSON(response)
		} else {
			fmt.Println("User updated successfully!")
			printUserStdout(response)
		}
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Delete a user from DCI",
	Run: func(cmd *cobra.Command, args []string) {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			fmt.Println(err)
			return
		}

		if deleteUserIDFlag == "" {
			fmt.Println("Error: --id is required")
			return
		}

		client := lib.NewClient(accessKey, secretKey)

		if outputFormat != OutputFormatJSON {
			fmt.Printf("Deleting user: %s\n", deleteUserIDFlag)
		}

		err = client.DeleteUser(deleteUserIDFlag)
		if err != nil {
			fmt.Printf("Failed to delete user: %v\n", err)
			return
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteUserIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println("User deleted successfully!")
		}
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

func printUsersJSON(response *lib.UsersResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
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

func printUserJSON(response *lib.UserResponse) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
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
