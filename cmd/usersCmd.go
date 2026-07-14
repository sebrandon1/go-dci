package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-dci/lib"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
	usersNameFilter       string
)

var getUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get all users from DCI, optionally filtered by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		if usersNameFilter != "" {
			printStatus("Getting users matching name: %s\n", usersNameFilter)
		} else {
			printStatus("Getting users...")
		}

		response, err := client.GetUsersFiltered(cmd.Context(), usersNameFilter)
		if err != nil {
			return fmt.Errorf("failed to get users: %w", err)
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
		if err := validateResourceID(getUserIDFlag, "user"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		printStatus("Getting user with ID: %s\n", getUserIDFlag)

		response, err := client.GetUser(cmd.Context(), getUserIDFlag)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
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
		password := createUserPasswordFlag
		if password == "" {
			if term.IsTerminal(int(os.Stdin.Fd())) {
				var err error
				password, err = readPassword("Enter password: ")
				if err != nil {
					return err
				}
				if password == "" {
					return fmt.Errorf("password cannot be empty")
				}
			} else {
				return fmt.Errorf("password is required: use --password flag or run interactively")
			}
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		client := lib.NewClient(accessKey, secretKey)

		if dryRunFlag {
			printStatus("[DRY RUN] Would create user: name=%s, email=%s, fullname=%s, team-id=%s\n", createUserNameFlag, createUserEmailFlag, createUserFullnameFlag, createUserTeamIDFlag)
			return nil
		}

		printStatus("Creating user: %s\n", createUserNameFlag)

		response, err := client.CreateUser(
			cmd.Context(),
			createUserNameFlag,
			createUserEmailFlag,
			createUserFullnameFlag,
			createUserTeamIDFlag,
			password,
		)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printUserJSON(response)
		}

		printStatus("User created successfully!")
		printUserStdout(response)

		return nil
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update-user",
	Short: "Update an existing user in DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(updateUserIDFlag, "user"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
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

		if dryRunFlag {
			printStatus("[DRY RUN] Would update user: id=%s, name=%s, email=%s, fullname=%s, state=%s\n", updateUserIDFlag, updateUserNameFlag, updateUserEmailFlag, updateUserFullnameFlag, updateUserStateFlag)
			return nil
		}

		printStatus("Updating user: %s\n", updateUserIDFlag)

		response, err := client.UpdateUser(cmd.Context(), updateUserIDFlag, updates)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			return printUserJSON(response)
		}

		printStatus("User updated successfully!")
		printUserStdout(response)

		return nil
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Delete a user from DCI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateResourceID(deleteUserIDFlag, "user"); err != nil {
			return err
		}

		accessKey, secretKey, err := getCredentials()
		if err != nil {
			return err
		}

		if dryRunFlag {
			printStatus("[DRY RUN] Would delete user: id=%s\n", deleteUserIDFlag)
			return nil
		}

		// Confirm deletion
		confirmed, err := confirmDeletion("user", deleteUserIDFlag)
		if err != nil {
			return err
		}
		if !confirmed {
			printStatus("Deletion canceled")
			return nil
		}


		client := lib.NewClient(accessKey, secretKey)

		printStatus("Deleting user: %s\n", deleteUserIDFlag)

		err = client.DeleteUser(cmd.Context(), deleteUserIDFlag)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		if outputFormat == OutputFormatJSON {
			result := map[string]string{"status": "deleted", "id": deleteUserIDFlag}
			jsonBytes, _ := json.Marshal(result)
			fmt.Println(string(jsonBytes))
		} else {
			printStatus("User deleted successfully!")
		}

		return nil
	},
}

func printUsersStdout(responses []lib.UsersResponse) {
	totalUsers := 0
	for _, resp := range responses {
		totalUsers += len(resp.Users)
	}

	if totalUsers == 0 {
		fmt.Println("No users found.")
		return
	}

	fmt.Println("---")
	for _, resp := range responses {
		for _, user := range resp.Users {
			fmt.Printf("ID: %s | Name: %s | Email: %s | State: %s\n",
				user.ID, user.Name, user.Email, user.State)
		}
	}
	fmt.Printf("Total Users: %d\n", totalUsers)
}

func printUsersJSON(responses []lib.UsersResponse) error {
	// Flatten all users from paginated responses
	var allUsers []lib.User
	for _, resp := range responses {
		allUsers = append(allUsers, resp.Users...)
	}

	jsonBytes, err := json.Marshal(map[string]interface{}{
		"users": allUsers,
		"total": len(allUsers),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
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
		return fmt.Errorf("failed to marshal JSON: %w", err)
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
	getUsersCmd.PersistentFlags().StringVarP(&usersNameFilter, "name", "n", "", "Filter users by name")

	// get user flags
	getUserCmd.PersistentFlags().StringVar(&getUserIDFlag, "id", "", "User ID")
	_ = getUserCmd.MarkPersistentFlagRequired("id")

	// create user flags
	createUserCmd.PersistentFlags().StringVar(&createUserNameFlag, "name", "", "Username")
	_ = createUserCmd.MarkPersistentFlagRequired("name")
	createUserCmd.PersistentFlags().StringVar(&createUserEmailFlag, "email", "", "User email")
	_ = createUserCmd.MarkPersistentFlagRequired("email")
	createUserCmd.PersistentFlags().StringVar(&createUserFullnameFlag, "fullname", "", "User full name")
	createUserCmd.PersistentFlags().StringVar(&createUserTeamIDFlag, "team-id", "", "Team ID")
	_ = createUserCmd.MarkPersistentFlagRequired("team-id")
	createUserCmd.PersistentFlags().StringVar(&createUserPasswordFlag, "password", "", "User password (prompts interactively if omitted)")

	// update user flags
	updateUserCmd.PersistentFlags().StringVar(&updateUserIDFlag, "id", "", "User ID to update")
	_ = updateUserCmd.MarkPersistentFlagRequired("id")
	updateUserCmd.PersistentFlags().StringVar(&updateUserNameFlag, "name", "", "New username")
	updateUserCmd.PersistentFlags().StringVar(&updateUserEmailFlag, "email", "", "New email")
	updateUserCmd.PersistentFlags().StringVar(&updateUserFullnameFlag, "fullname", "", "New full name")
	updateUserCmd.PersistentFlags().StringVar(&updateUserStateFlag, "state", "", "New user state")

	// delete user flags
	deleteUserCmd.PersistentFlags().StringVar(&deleteUserIDFlag, "id", "", "User ID to delete")
	_ = deleteUserCmd.MarkPersistentFlagRequired("id")
}
